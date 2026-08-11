#!/bin/bash
set -euo pipefail

# install.sh — Install apple-music-mcp from GitHub Releases
# Usage: curl -fsSL https://raw.githubusercontent.com/AVTAVANTTOUT2/apple-music-mcp/main/scripts/install.sh | bash

REPO="AVTAVANTTOUT2/apple-music-mcp"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"

echo "Installing apple-music-mcp..."

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    arm64)   ASSET_ARCH="arm64" ;;
    x86_64)  ASSET_ARCH="amd64" ;;
    *)
        # Try universal binary
        ASSET_ARCH="universal"
        ;;
esac

OS="darwin"
BINARY="apple-music-mcp-${OS}-${ASSET_ARCH}"
ARCHIVE="${BINARY}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/${VERSION}/download/${ARCHIVE}"

echo "Downloading ${DOWNLOAD_URL}..."
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$DOWNLOAD_URL" -o "$TMPDIR/$ARCHIVE"

# Verify checksum if available
if curl -fsSL "https://github.com/${REPO}/releases/${VERSION}/download/checksums.txt" -o "$TMPDIR/checksums.txt" 2>/dev/null; then
    echo "Verifying checksum..."
    EXPECTED=$(grep "$ARCHIVE" "$TMPDIR/checksums.txt" | awk '{print $1}')
    ACTUAL=$(shasum -a 256 "$TMPDIR/$ARCHIVE" | awk '{print $1}')
    if [ "$EXPECTED" != "$ACTUAL" ]; then
        echo "ERROR: Checksum verification failed!"
        echo "Expected: $EXPECTED"
        echo "Got:      $ACTUAL"
        exit 1
    fi
    echo "Checksum OK"
fi

# Extract
tar xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR/"

# Install
mkdir -p "$INSTALL_DIR"
cp "$TMPDIR/$BINARY" "$INSTALL_DIR/apple-music-mcp"
chmod +x "$INSTALL_DIR/apple-music-mcp"

echo ""
echo "apple-music-mcp installed to $INSTALL_DIR/apple-music-mcp"
echo ""
echo "Run 'apple-music-mcp doctor' to verify your setup."
echo ""
echo "To configure your MCP client:"
echo "  apple-music-mcp configure"
