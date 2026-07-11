#!/usr/bin/env bash

set -euo pipefail

# GitHub repository details
OWNER="silverAndroid"
REPO="gradle-for-agents"

# Determine OS and Architecture
OS_RAW=$(uname -s)
ARCH_RAW=$(uname -m)

case "$OS_RAW" in
    Darwin)
        OS="darwin"
        ;;
    Linux)
        OS="linux"
        ;;
    *)
        echo "ERROR: Unsupported Operating System: $OS_RAW" >&2
        exit 1
        ;;
esac

case "$ARCH_RAW" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "ERROR: Unsupported Architecture: $ARCH_RAW" >&2
        exit 1
        ;;
esac

echo "Detecting latest release of ${OWNER}/${REPO}..."

# Retrieve latest release tag without hitting GitHub API rate limits
LATEST_URL=$(curl -fsSL -o /dev/null -w "%{url_effective}" "https://github.com/${OWNER}/${REPO}/releases/latest")
TAG=$(basename "$LATEST_URL")

# Fallback: if we couldn't resolve the tag (e.g. no releases yet or URL structure change), try GitHub API
if [[ "$TAG" == "latest" ]] || [[ -z "$TAG" ]]; then
    TAG=$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

if [[ -z "$TAG" ]]; then
    echo "ERROR: Could not resolve the latest release tag." >&2
    exit 1
fi

VERSION="${TAG#v}"
echo "Latest version found: ${TAG} (${VERSION})"

# Create temporary directory for downloads
TEMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "$TEMP_DIR"
}
trap clean_up EXIT

# Binaries to download and install
BINARIES=("gradle-for-agents" "gfa")

for binary in "${BINARIES[@]}"; do
    ARCHIVE_NAME="${binary}_${VERSION}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${ARCHIVE_NAME}"
    
    echo "Downloading ${binary} from ${DOWNLOAD_URL}..."
    if ! curl -fsSL "$DOWNLOAD_URL" -o "${TEMP_DIR}/${ARCHIVE_NAME}"; then
        echo "ERROR: Failed to download ${binary}." >&2
        exit 1
    fi
    
    echo "Extracting ${binary}..."
    tar -xzf "${TEMP_DIR}/${ARCHIVE_NAME}" -C "$TEMP_DIR" "${binary}"
done

# Installation target
INSTALL_DIR="/usr/local/bin"
USE_SUDO=false

if [ ! -w "$INSTALL_DIR" ]; then
    # Try to use sudo if it's not writable and user is not root
    if [ "$EUID" -ne 0 ] && command -v sudo >/dev/null 2>&1; then
        echo "Note: ${INSTALL_DIR} is not writable. Attempting to install with sudo..."
        USE_SUDO=true
    else
        # Fall back to user local directory
        INSTALL_DIR="${HOME}/.local/bin"
        echo "Note: ${INSTALL_DIR} is not writable. Installing to ${INSTALL_DIR} instead..."
        mkdir -p "$INSTALL_DIR"
    fi
fi

for binary in "${BINARIES[@]}"; do
    echo "Installing ${binary} to ${INSTALL_DIR}..."
    if [ "$USE_SUDO" = true ]; then
        sudo cp "${TEMP_DIR}/${binary}" "${INSTALL_DIR}/${binary}"
        sudo chmod +x "${INSTALL_DIR}/${binary}"
    else
        cp "${TEMP_DIR}/${binary}" "${INSTALL_DIR}/${binary}"
        chmod +x "${INSTALL_DIR}/${binary}"
    fi
done

echo "Successfully installed gradle-for-agents and gfa!"

# Verify if INSTALL_DIR is in PATH
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "WARNING: ${INSTALL_DIR} is not in your PATH."
    echo "To be able to run 'gfa' and 'gradle-for-agents', please add it to your shell profile."
    echo "For example, add the following line to your ~/.zshrc or ~/.bashrc:"
    echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
    echo ""
fi
