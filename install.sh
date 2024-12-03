#!/bin/bash

set -euo pipefail

# Constants
REPO="edsonjaramillo/cryptid"
BINARY="cryptid"
GITHUB_API="https://api.github.com/repos/${REPO}"
INSTALL_DIR="/usr/local/bin"

# OS/Arch detection
get_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        FreeBSD*) echo "freebsd" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "Unsupported operating system"; exit 1 ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "Unsupported architecture"; exit 1 ;;
    esac
}

# Get latest release version
get_latest_version() {
    if ! curl -sL "${GITHUB_API}/releases/latest" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4; then
        echo "Failed to fetch latest version"
        exit 1
    fi
}

# Install binary
install() {
    local os arch version tmp
    os=$(get_os)
    arch=$(get_arch)
    version=$(get_latest_version)
    tmp=$(mktemp -d)
    
    local archive="${BINARY}-${version}-${os}-${arch}.tar.gz"
    local url="https://github.com/${REPO}/releases/download/${version}/${archive}"
    
    # Download and extract
    if ! curl -sL "$url" -o "${tmp}/${archive}"; then
        echo "Download failed"
        exit 1
    fi
    
    tar -xzf "${tmp}/${archive}" -C "$tmp"
    
    # Install binary
    if [ "$os" = "windows" ]; then
        mv "${tmp}/${BINARY}.exe" "${INSTALL_DIR}/${BINARY}.exe"
    else
        mv "${tmp}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
        chmod +x "${INSTALL_DIR}/${BINARY}"
    fi
    
    rm -rf "$tmp"
    echo "Successfully installed ${BINARY} to ${INSTALL_DIR}"
}

# Check permissions
if [ "$(id -u)" -ne 0 ] && [ "$(get_os)" != "windows" ]; then
    echo "Please run as root"
    exit 1
fi

# Run installation
install