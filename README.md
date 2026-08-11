# apple-music-mcp

Native, secure, zero-dependency MCP server to control Apple Music on macOS.

**No Python. No Node. No runtime to install.** Just a single binary.

[![CI](https://github.com/AVTAVANTTOUT2/apple-music-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/AVTAVANTTOUT2/apple-music-mcp/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AVTAVANTTOUT2/apple-music-mcp)](https://goreportcard.com/report/github.com/AVTAVANTTOUT2/apple-music-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## What It Does

apple-music-mcp connects AI assistants (Cursor, Claude Desktop, Codex, VS Code) directly to the Music.app on your Mac. Your AI can now control playback, search your library, manage playlists, set ratings and favorites — all without leaving the natural language conversation.

## Quick Install

```bash
# macOS 13+ (Apple Silicon or Intel)
curl -fsSL https://raw.githubusercontent.com/AVTAVANTTOUT2/apple-music-mcp/main/scripts/install.sh | bash
```

The installer script is versioned in the repository (`scripts/install.sh`). Release assets ship the compiled binaries (`.tar.gz` + `checksums.txt`); use the command above for one-line install.

Or download the binary directly from [Releases](https://github.com/AVTAVANTTOUT2/apple-music-mcp/releases).

## Requirements

- macOS 13 (Ventura) or later
- Apple Music.app (pre-installed on macOS)
- Automation permission for your terminal/client
- No Xcode, no Homebrew, no Python, no Node.js required

## Client Configuration

### Cursor

```json
{
  "mcpServers": {
    "apple-music": {
      "command": "/usr/local/bin/apple-music-mcp",
      "args": ["serve"]
    }
  }
}
```

### Claude Desktop

```json
{
  "mcpServers": {
    "apple-music": {
      "command": "/usr/local/bin/apple-music-mcp",
      "args": ["serve"]
    }
  }
}
```

Configuration file: `~/Library/Application Support/Claude/claude_desktop_config.json`

## macOS Permissions

When you first use apple-music-mcp, macOS will ask:

1. **Automation permission**: Allow [your terminal/client] to control Music.app
   - Settings → Privacy & Security → Automation

No other permissions required. We never ask for your Apple ID password.

## Commands

```bash
apple-music-mcp serve          # Start the MCP server (stdio)
apple-music-mcp doctor         # Run diagnostics
apple-music-mcp capabilities   # Show what's available on your system
apple-music-mcp version        # Print version info
apple-music-mcp configure      # Print MCP client config
```

## Features

### Playback Control
- ✅ Play, pause, toggle, stop
- ✅ Next/previous track, restart current
- ✅ Seek to position
- ✅ Volume, mute
- ✅ Shuffle (off/songs/albums)
- ✅ Repeat (off/one/all)

### Library & Metadata
- ✅ Current track info (name, artist, album, duration, genre, etc.)
- ✅ Artwork detection
- ✅ Search library by track, album, artist, playlist
- ✅ Browse recently added/played
- ✅ List tracks, albums, artists, genres

### Playlists
- ✅ List all playlists
- ✅ Create, rename, delete playlists
- ✅ Add/remove tracks
- ✅ Duplicate playlists
- ✅ Folder playlists

### Favorites & Ratings
- ✅ Favorite/unfavorite tracks and albums
- ✅ Suggest less (dislike)
- ✅ Set rating (0-100 scale, maps to 0-5 stars)

### AirPlay
- ✅ List AirPlay devices
- ✅ Select output device
- ✅ Device volume control

### User Interface
- ✅ Reveal item in Music.app
- ✅ Open Apple Music URLs

### Known Limitations (v0.1.0)
- ❌ Up Next queue — not available via Apple Events (requires Accessibility or MusicKit fallback)
- ❌ Apple Music catalog search — requires MusicKit developer token
- ❌ Autoplay toggle — not scriptable

See [docs/capability-matrix.md](docs/capability-matrix.md) for the full breakdown.

## MCP Tools

| Tool | Description |
|------|-------------|
| `music_get_state` | Current playback state and track info |
| `music_playback` | Control playback (play, pause, skip, seek, reveal) |
| `music_preferences` | Volume, shuffle, repeat, AirPlay output |
| `music_search` | Search library for tracks, albums, artists, playlists |
| `music_playlists` | Full playlist management |
| `music_library` | Browse library (recent, albums, artists, genres) |
| `music_favorites` | Manage favorites, ratings, dislikes |
| `music_queue` | Up Next queue (not yet available) |
| `music_capabilities` | System capability report |
| `music_doctor` | Non-destructive diagnostics |

## Example Prompts

- "What's playing right now?"
- "Pause the music."
- "Play Random Access Memories by Daft Punk."
- "Search my library for jazz tracks."
- "Create a playlist called Focus and add these tracks."
- "Set volume to 40%."
- "What are my playlists?"
- "Favorite this track."
- "Show me my recently added music."

## Read-Only Mode

```bash
APPLE_MUSIC_MCP_READ_ONLY=1 apple-music-mcp serve
```

Disables all mutations: no play, no delete, no favorite changes.

## Architecture

```
apple-music-mcp
├── cmd/apple-music-mcp/    # CLI entry point
├── internal/
│   ├── domain/             # Core types and interfaces
│   ├── backend/            # Backend abstraction
│   │   └── musicapp/       # Apple Events (osascript) backend
│   ├── applescript/        # Safe AppleScript execution + parsing
│   ├── mcpserver/          # MCP protocol server
│   ├── config/             # Configuration (env vars only)
│   ├── logging/            # Structured logging (stderr)
│   └── version/            # Build-time version info
└── scripts/                # Static AppleScript files
```

### Design Principles

- **Local-first**: Everything runs on your Mac, no cloud dependency
- **Zero telemetry**: No analytics, no tracking, no phone home
- **No secrets stored**: No Apple ID, no passwords, no tokens in plaintext
- **stdout is sacred**: Only JSON-RPC on stdout, all logs to stderr
- **Safe by default**: Input validated through schemas, never interpolated into scripts
- **Honest capabilities**: What's impossible is documented, not simulated

## Security

See [SECURITY.md](SECURITY.md) and [docs/threat-model.md](docs/threat-model.md) for the full security analysis.

## Development

```bash
git clone https://github.com/AVTAVANTTOUT2/apple-music-mcp.git
cd apple-music-mcp
go build -o apple-music-mcp ./cmd/apple-music-mcp/
go test -race ./...
go vet ./...
```

### Live tests (Music.app)

Run the full JARVIS workflow against the real Music.app (search, play, favorite, playlist):

```bash
APPLE_MUSIC_MCP_LIVE_TESTS=1 ./apple-music-mcp test-live
```

Requires Music.app, automation permission, and at least one track matching the search query (default: `Werenoi`).

### JARVIS / Agent integration

See [docs/jarvis-integration.md](docs/jarvis-integration.md) for MCP tool reference, example prompts, and troubleshooting.

## License

MIT — see [LICENSE](LICENSE)

---

Built with Apple Events, Go, and the [Model Context Protocol](https://modelcontextprotocol.io).
