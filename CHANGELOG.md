# Changelog

All notable changes to apple-music-mcp will be documented in this file.

## [Unreleased]

### Fixed

- **Install URL (404)**: README and `install.sh` now point to `raw.githubusercontent.com/.../scripts/install.sh` instead of a non-existent release asset; future releases also attach `install.sh` as an asset
- **osascript argv**: `.applescript` files are invoked without `-l AppleScript` / `--`, fixing all JSON-based MCP tools (`music_search`, `music_playlists`, `music_favorites`, etc.)
- **Locale decimals**: French macOS decimal commas in AppleScript JSON output are normalized before parsing
- **play_track mapping**: MCP `play_track` / `play_album` / `play_artist` actions are mapped to AppleScript `play` handlers

### Added

- **`test-live` command**: Full JARVIS workflow (search → play → favorite → playlist → cleanup)
- **JARVIS documentation**: [docs/jarvis-integration.md](docs/jarvis-integration.md)
- **Unit tests**: osascript argv builder, JSON locale normalization, playback action mapping

## [0.1.0] — 2026-08-11

### Initial Release

- **MCP Server**: Full Model Context Protocol server (stdio transport, protocol 2024-11-05)
- **Playback Control**: Play, pause, toggle, stop, next, previous, restart, seek
- **Player State**: Comprehensive metadata (track, artist, album, duration, genre, year, etc.)
- **Preferences**: Volume, mute, shuffle, repeat, AirPlay output
- **Search**: Library search by track, album, artist, playlist
- **Playlists**: Full CRUD, add/remove tracks, copy, folder management
- **Library**: Browse tracks, albums, artists, genres, recently added, recently played
- **Favorites**: Favorite/unfavorite, dislike, ratings (0-100 scale)
- **Capabilities**: Dynamic capability detection per system
- **Doctor**: Non-destructive system diagnostics
- **CLI**: serve, doctor, capabilities, version, configure commands
- **Security**: Input validation, static scripts, no interpolation, read-only mode
- **Zero dependencies**: Single binary, no runtime required
- **macOS 13+**: Apple Silicon (arm64) and Intel (amd64), Universal binary
- **Tests**: Unit tests (MCP server, domain), race condition detection

### Known Limitations

- Up Next queue not available via Apple Events
- Apple Music catalog search requires MusicKit (not yet integrated)
- No Windows or Linux support (macOS-only by design)
