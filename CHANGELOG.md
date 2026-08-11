# Changelog

All notable changes to apple-music-mcp will be documented in this file.

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
