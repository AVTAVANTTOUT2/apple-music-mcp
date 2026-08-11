// Package version provides build-time version information.
package version

import (
	"fmt"
	"runtime"
)

// These values are set at build time via ldflags.
var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info returns a formatted version string.
func Info() string {
	return fmt.Sprintf("apple-music-mcp v%s (commit: %s, built: %s, go: %s, os: %s/%s)",
		Version, Commit, BuildDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// Short returns just the version number.
func Short() string {
	return Version
}
