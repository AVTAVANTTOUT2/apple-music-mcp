// Package applescript provides safe, parameterized execution of AppleScript
// and JavaScript for Automation (JXA) scripts. Input is passed as arguments
// to static scripts — never interpolated into script source.
package applescript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AVTAVANTTOUT2/apple-music-mcp/internal/logging"
)

const (
	// DefaultTimeout is the maximum execution time for an AppleScript call.
	DefaultTimeout = 15 * time.Second

	// MaxArgLength prevents excessively long arguments.
	MaxArgLength = 4096

	// MaxArgs prevents argument flooding.
	MaxArgs = 50
)

// Executor runs AppleScript and JXA scripts safely.
type Executor struct {
	timeout time.Duration
	logger  *logging.Logger
}

// NewExecutor creates a new AppleScript executor with the default timeout.
func NewExecutor(logger *logging.Logger) *Executor {
	return &Executor{
		timeout: DefaultTimeout,
		logger:  logger,
	}
}

// SetTimeout overrides the default timeout.
func (e *Executor) SetTimeout(d time.Duration) {
	e.timeout = d
}

// RunAppleScript executes a static AppleScript file with the given arguments.
// The script path must be a file in the scripts/ directory, not user-provided.
// Arguments are passed via osascript's argv mechanism, never interpolated.
func (e *Executor) RunAppleScript(ctx context.Context, scriptPath string, args ...string) (string, error) {
	return e.runScript(ctx, scriptPath, "AppleScript", args...)
}

// RunJXA executes a static JXA (JavaScript for Automation) file with the given arguments.
func (e *Executor) RunJXA(ctx context.Context, scriptPath string, args ...string) (string, error) {
	return e.runScript(ctx, scriptPath, "JavaScript", args...)
}

// runScript executes an osascript command with the given language.
func (e *Executor) runScript(ctx context.Context, scriptPath string, language string, args ...string) (string, error) {
	if err := validatePath(scriptPath); err != nil {
		return "", fmt.Errorf("invalid script path %q: %w", scriptPath, err)
	}
	if err := validateArgs(args); err != nil {
		return "", fmt.Errorf("invalid arguments for script %q: %w", scriptPath, err)
	}

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmdArgs := buildOsascriptArgs(scriptPath, language, args)
	cmd := exec.CommandContext(ctx, "osascript", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	e.logger.Debug("executing %s script: %s with %d args", language, scriptPath, len(args))

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		e.logger.Error("%s execution failed after %v: %v (stderr: %s)", language, elapsed, err, stderrStr)
		return "", fmt.Errorf("%s execution failed: %w — %s", language, err, stderrStr)
	}

	result := strings.TrimSpace(stdout.String())
	e.logger.Debug("%s execution succeeded in %v, output length: %d", language, elapsed, len(result))
	return result, nil
}

// RunAppleScriptString executes an inline AppleScript command.
// This MUST ONLY be used for commands whose arguments have been validated.
// NEVER pass user input directly to this method.
func (e *Executor) RunAppleScriptString(ctx context.Context, script string) (string, error) {
	if strings.TrimSpace(script) == "" {
		return "", fmt.Errorf("script must not be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "osascript", "-e", script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		e.logger.Error("inline AppleScript failed after %v: %v (stderr: %s)", elapsed, err, stderrStr)
		return "", fmt.Errorf("AppleScript execution failed: %w — %s", err, stderrStr)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// buildOsascriptArgs builds the argv slice for osascript.
//
// .applescript files must be invoked as:
//
//	osascript /path/to/script.applescript arg1 arg2
//
// Using "-l AppleScript" together with "--" makes osascript pass "--" as argv[1]
// to the script, which breaks every handler that reads item 1 of argv.
//
// JXA (.js) files keep the language flag but never use "--".
func buildOsascriptArgs(scriptPath, language string, args []string) []string {
	ext := strings.ToLower(filepath.Ext(scriptPath))
	if ext == ".applescript" {
		cmdArgs := []string{scriptPath}
		return append(cmdArgs, args...)
	}

	cmdArgs := []string{"-l", language, scriptPath}
	return append(cmdArgs, args...)
}

// ParseJSONOutput parses JSON output from a script that returns a JSON string.
// AppleScript on French macOS locales may emit decimal commas in numbers; these
// are normalized to JSON-compatible decimal points before unmarshaling.
func ParseJSONOutput(output string, target interface{}) error {
	if output == "" {
		return fmt.Errorf("empty script output")
	}
	normalized := NormalizeJSONNumbers(output)
	if err := json.Unmarshal([]byte(normalized), target); err != nil {
		return fmt.Errorf("failed to parse script JSON output: %w (raw: %.200s)", err, output)
	}
	return nil
}

// validatePath ensures the script path is not dangerous.
func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("script path must not be empty")
	}
	// Reject paths that try to escape or use shell metacharacters
	if strings.Contains(path, "..") {
		return fmt.Errorf("script path must not contain '..'")
	}
	for _, c := range path {
		if c == '|' || c == ';' || c == '&' || c == '$' || c == '`' || c == '\'' || c == '"' || c == '\n' || c == '\r' {
			return fmt.Errorf("script path contains invalid character: %c", c)
		}
	}
	return nil
}

// validateArgs ensures arguments are safe and within limits.
func validateArgs(args []string) error {
	if len(args) > MaxArgs {
		return fmt.Errorf("too many arguments: %d (max %d)", len(args), MaxArgs)
	}
	for i, arg := range args {
		if len(arg) > MaxArgLength {
			return fmt.Errorf("argument %d too long: %d chars (max %d)", i, len(arg), MaxArgLength)
		}
		// Reject null bytes and other dangerous control characters
		for _, c := range arg {
			if c == 0 || c == '\x00' {
				return fmt.Errorf("argument %d contains null byte", i)
			}
		}
	}
	return nil
}

// EscapeAppleScriptString escapes a string for safe use in inline AppleScript.
// Usage: ONLY for trusted, validated identifiers, never for raw user input.
func EscapeAppleScriptString(s string) string {
	// Escape backslashes and double quotes for AppleScript string literals
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// SanitizeArg removes dangerous characters from a single argument.
func SanitizeArg(s string) string {
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}
