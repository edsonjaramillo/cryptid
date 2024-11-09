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
        *) error "Unsupported operating system" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture" ;;
    esac
}

# Get latest release version
get_latest_version() {
    if ! curl -sL "${GITHUB_API}/releases/latest" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4; then
        error "Failed to fetch latest version"
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
        error "Download failed"
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
    error "This script must be run as root (try using sudo)"
fi

# Run installation
install