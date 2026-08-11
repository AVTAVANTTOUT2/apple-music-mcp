// Package logging provides structured logging that writes exclusively to stderr.
// stdout is reserved for the MCP protocol transport.
package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents the severity of a log message.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

// Logger writes structured logs to stderr (or a configured writer).
type Logger struct {
	mu     sync.Mutex
	writer io.Writer
	level  Level
}

// NewLogger creates a logger that writes to stderr.
func NewLogger(level Level) *Logger {
	return &Logger{
		writer: os.Stderr,
		level:  level,
	}
}

// NewFileLogger creates a logger that writes to a file.
func NewFileLogger(level Level, path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", path, err)
	}
	return &Logger{
		writer: f,
		level:  level,
	}, nil
}

// SetLevel changes the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info logs an informational message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	levelStr := levelNames[level]
	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(l.writer, "[%s] %-5s %s\n", timestamp, levelStr, msg)
}

// Close closes the underlying writer if it's closable (e.g., a file).
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if closer, ok := l.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Writer returns the underlying writer for use with other libraries.
func (l *Logger) Writer() io.Writer {
	return l.writer
}
