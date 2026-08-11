// Package scripts embeds the AppleScript files into the binary.
package scripts

import "embed"

// FS contains all AppleScript files needed at runtime.
//
//go:embed get_state.applescript
//go:embed playback.applescript
//go:embed preferences.applescript
//go:embed search.applescript
//go:embed playlists.applescript
//go:embed library.applescript
//go:embed favorites.applescript
var FS embed.FS
