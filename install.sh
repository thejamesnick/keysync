#!/bin/bash
set -e

REPO_OWNER="thejamesnick"
REPO_NAME="keysync"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="keysync"

echo "🔐 Installing KeySync..."

# 1. Detect OS & Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$OS" == "darwin" ]; then
    if [ "$ARCH" == "arm64" ]; then
        TARGET="darwin-arm64"
    else
        TARGET="darwin-amd64"
    fi
elif [ "$OS" == "linux" ]; then
    if [ "$ARCH" == "aarch64" ]; then
        TARGET="linux-arm64"
    else
        TARGET="linux-amd64"
    fi
else
    echo "❌ Usage: Unsupported OS ($OS)"
    exit 1
fi

echo "  🔍 Detected: $OS/$ARCH (Target: $TARGET)"

# 2. Find Latest Release URL
# Uses GitHub API to get the download URL for the specific asset
ASSET_URL=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | \
    grep "browser_download_url" | \
    grep "$TARGET" | \
    cut -d '"' -f 4)

if [ -z "$ASSET_URL" ]; then
    echo "❌ Error: Could not find release asset for $TARGET."
    echo "   Ensure a release exists at github.com/$REPO_OWNER/$REPO_NAME/releases"
    exit 1
fi

echo "  📥 Downloading from GitHub..."
curl -sL -o "$BINARY_NAME" "$ASSET_URL"
chmod +x "$BINARY_NAME"

# 3. Install
echo "  📦 Moving to $INSTALL_DIR (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
fi

echo "  ✅ Installed! Run 'keysync help' to get started."
