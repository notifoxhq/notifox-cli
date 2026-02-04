#!/bin/sh
# Install notifox CLI from GitHub releases.
# Usage: curl -fsSL https://raw.githubusercontent.com/notifoxhq/notifox-cli/main/install.sh | sh
#    or: ./install.sh

set -e

REPO="notifoxhq/notifox-cli"
BINARY="notifox"
INSTALL_DIR="${NOTIFOX_INSTALL_DIR:-$HOME/.local/bin}"
# Use /usr/local/bin if we have write access and no custom dir requested
if [ -z "$NOTIFOX_INSTALL_DIR" ] && [ -w /usr/local/bin ] 2>/dev/null; then
  INSTALL_DIR="/usr/local/bin"
fi

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    echo "For Windows, run: irm https://notifox.com/install.ps1 | iex"
    echo "Or download from: https://github.com/$REPO/releases"
    exit 1
    ;;
esac

# Download URL (latest release)
VERSION="${NOTIFOX_VERSION:-latest}"
if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/$REPO/releases/latest/download/notifox-cli_${OS}_${ARCH}.tar.gz"
else
  URL="https://github.com/$REPO/releases/download/${VERSION}/notifox-cli_${OS}_${ARCH}.tar.gz"
fi

echo "Installing notifox to $INSTALL_DIR"
echo "  OS: $OS  Arch: $ARCH"
echo "  URL: $URL"

mkdir -p "$INSTALL_DIR"
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$URL" -o "$tmpdir/archive.tar.gz"
elif command -v wget >/dev/null 2>&1; then
  wget -q "$URL" -O "$tmpdir/archive.tar.gz"
else
  echo "Need curl or wget to download."
  exit 1
fi

tar -xzf "$tmpdir/archive.tar.gz" -C "$tmpdir"
mv "$tmpdir/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

# Ensure install dir is on PATH
case ":$PATH:" in
  *:"$INSTALL_DIR":*) ;;
  *)
    echo ""
    echo "Add notifox to your PATH:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo "Add the line above to your shell profile (~/.bashrc, ~/.zshrc, etc.)"
    ;;
esac

echo "Installed: $INSTALL_DIR/$BINARY"
"$INSTALL_DIR/$BINARY" version 2>/dev/null || true
