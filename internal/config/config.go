// Package config manages the application configuration.
// All configuration comes from environment variables — zero hardcoded values.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all configuration for the Apple Music MCP server.
type Config struct {
	// ReadOnly mode: disables all mutations (play, delete, favorite, etc.)
	ReadOnly bool

	// Verbose enables debug logging to stderr.
	Verbose bool

	// LogFile is the path to an optional log file. Empty means stderr only.
	LogFile string

	// Transport is the MCP transport: "stdio" or "http".
	Transport string

	// HTTPAddress is the listen address for HTTP transport (127.0.0.1:PORT).
	HTTPAddress string

	// ScriptsDir overrides the directory containing AppleScript files.
	ScriptsDir string

	// TimeoutSeconds is the maximum time for an AppleScript call.
	TimeoutSeconds int

	// AllowQueueMutationTests enables destructive queue tests.
	AllowQueueMutationTests bool

	// AllowPlaybackTests enables playback tests.
	AllowPlaybackTests bool

	// AllowLiveTests enables real Music.app tests.
	AllowLiveTests bool

	// ConfigDir is the path to the configuration directory.
	ConfigDir string

	// LogDir is the path to the log directory.
	LogDir string

	// CacheDir is the path to the cache directory.
	CacheDir string
}

// DefaultConfig returns a Config with safe defaults.
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()

	return &Config{
		ReadOnly:                getEnvBool("APPLE_MUSIC_MCP_READ_ONLY", false),
		Verbose:                 getEnvBool("APPLE_MUSIC_MCP_VERBOSE", false),
		LogFile:                 os.Getenv("APPLE_MUSIC_MCP_LOG_FILE"),
		Transport:               getEnvStr("APPLE_MUSIC_MCP_TRANSPORT", "stdio"),
		HTTPAddress:             getEnvStr("APPLE_MUSIC_MCP_HTTP_ADDRESS", "127.0.0.1:9020"),
		ScriptsDir:              os.Getenv("APPLE_MUSIC_MCP_SCRIPTS_DIR"),
		TimeoutSeconds:          getEnvInt("APPLE_MUSIC_MCP_TIMEOUT", 15),
		AllowLiveTests:          getEnvBool("APPLE_MUSIC_MCP_LIVE_TESTS", false),
		AllowQueueMutationTests: getEnvBool("APPLE_MUSIC_MCP_QUEUE_MUTATION_TESTS", false),
		AllowPlaybackTests:      getEnvBool("APPLE_MUSIC_MCP_ALLOW_PLAYBACK_TESTS", false),
		ConfigDir:               filepath.Join(home, "Library", "Application Support", "apple-music-mcp"),
		LogDir:                  filepath.Join(home, "Library", "Logs", "apple-music-mcp"),
		CacheDir:                filepath.Join(home, "Library", "Caches", "apple-music-mcp"),
	}
}

// LoadFromEnv loads configuration from environment variables over defaults.
func LoadFromEnv() *Config {
	return DefaultConfig()
}

func getEnvStr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v == "1" || v == "true" || v == "yes" || v == "TRUE" || v == "YES"
}

func getEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	var result int
	if _, err := fmt.Sscanf(v, "%d", &result); err != nil {
		return defaultVal
	}
	return result
}
