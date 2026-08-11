# Contributing

Thanks for your interest in improving apple-music-mcp!

## Getting Started

```bash
git clone https://github.com/AVTAVANTTOUT2/apple-music-mcp.git
cd apple-music-mcp
go build ./...
go test -race ./...
```

## Development Environment

- Go 1.24+
- macOS 13+ with Music.app
- Automation permission for your terminal

## Before Submitting

1. Format: `gofmt -w .`
2. Vet: `go vet ./...`
3. Test: `go test -race ./...`
4. Build: `go build -o apple-music-mcp ./cmd/apple-music-mcp/`

## Pull Request Checklist

- [ ] Tests pass locally
- [ ] No new warnings from `go vet`
- [ ] Code is properly formatted
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] No secrets or personal data committed

## Architecture

See `docs/capability-matrix.md` and `docs/architecture.md` for the full design.

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
