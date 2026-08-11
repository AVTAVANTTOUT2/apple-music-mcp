# Security Policy

## Reporting a Vulnerability

If you discover a security issue in apple-music-mcp, please report it privately.

**Do not file a public issue.**

Contact: [security contact] (to be configured by the repository owner)

## Security Model

apple-music-mcp is designed with security as a first-class constraint. Here's how:

### What We NEVER Do

- Never ask for or store your Apple ID password
- Never store secrets in plaintext
- Never expose arbitrary AppleScript execution
- Never expose arbitrary shell command execution
- Never use private Apple APIs
- Never modify SIP, TCC, or Gatekeeper
- Never send telemetry or analytics
- Never intercept or read Safari cookies without explicit consent

### What We ALWAYS Do

- Pass all user input as separate arguments to static scripts (never interpolated)
- Validate all input through strict schemas
- Reserve stdout exclusively for MCP JSON-RPC protocol
- Write all logs to stderr or a configured log file
- Use read-only mode by disabling mutations via `APPLE_MUSIC_MCP_READ_ONLY=1`
- Annotate destructive operations with `destructiveHint`
- Store credentials in macOS Keychain (if ever needed for MusicKit)
- Verify SHA-256 checksums before installation

### Surface Area

- **Transport**: stdio (default). HTTP transport disabled by default, if enabled listens only on 127.0.0.1.
- **Backend**: Apple Events (osascript) with static, pre-written scripts.
- **No network**: The native backend requires zero network access.

### macOS Permissions

The server requires only:
1. **Automation** permission: Allow your terminal/client to control Music.app

No Accessibility permission is required for the core features. It may be needed for future Up Next queue support.

## Threat Model

See [docs/threat-model.md](docs/threat-model.md) for the complete threat model.

## Dependencies

Go standard library + MCP Go SDK. All dependencies are pinned with checksums in go.sum.
