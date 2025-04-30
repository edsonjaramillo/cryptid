#!/usr/bin/env bash

# Strict mode
set -euo pipefail

# Constants
readonly REPO="edsonjaramillo/hyde"
readonly BINARY="hyde"
readonly GITHUB_API="https://api.github.com/repos/${REPO}"

# --- Configuration ---
# Allow overriding install directory via environment variable
# Default: /usr/local/bin if root, $HOME/.local/bin if not root
DEFAULT_INSTALL_DIR=""
# Check if running as root (UID 0)
if [[ "$(id -u)" -eq 0 ]]; then
    DEFAULT_INSTALL_DIR="/usr/local/bin"
else
    # Ensure HOME is set and non-empty for non-root users
    if [[ -z "${HOME:-}" ]]; then
        echo "ERROR: HOME environment variable is not set." >&2
        exit 1
    fi
    DEFAULT_INSTALL_DIR="${HOME}/.local/bin"
fi
INSTALL_DIR="${HYDE_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
# --- End Configuration ---

# Global variables
OS=""
ARCH=""
VERSION=""
TMP_DIR=""
CHECKSUM_TOOL=""

# --- Helper Functions ---

# Basic logging
info() {
    echo "INFO: $1"
}

error() {
    echo "ERROR: $1" >&2
    exit 1
}

# Cleanup temporary directory on exit
cleanup() {
    if [[ -n "${TMP_DIR:-}" && -d "${TMP_DIR}" ]]; then
        info "Cleaning up temporary directory: ${TMP_DIR}"
        rm -rf "${TMP_DIR}"
    fi
}
# Set trap for cleanup
trap cleanup EXIT INT TERM ERR

# Check for necessary commands
check_deps() {
    info "Checking dependencies..."
    local missing_deps=()
    # Removed windows-specific checks
    local deps=("curl" "tar" "mktemp" "uname" "mkdir" "mv" "chmod" "rm" "id" "grep" "cut" "awk")
    
    # Check basic dependencies
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &>/dev/null; then
            missing_deps+=("$dep")
        fi
    done

    # Check for checksum tool (sha256sum preferred)
    if command -v sha256sum &>/dev/null; then
        CHECKSUM_TOOL="sha256sum"
    elif command -v shasum &>/dev/null; then
        # Check if shasum supports -a 256 (macOS default shasum might not need it for SHA256, but standardizing)
        if shasum -a 256 < /dev/null &> /dev/null; then
           CHECKSUM_TOOL="shasum -a 256"
        elif shasum < /dev/null &> /dev/null; then # Fallback if -a 256 fails but shasum exists
            info "Using 'shasum' without '-a 256'. Assuming default is SHA256 or compatible."
            CHECKSUM_TOOL="shasum"
        else
             missing_deps+=("sha256sum or compatible shasum")
        fi
    else
        missing_deps+=("sha256sum or shasum")
    fi
    
    # Check for install command (optional but preferred)
    if ! command -v install &>/dev/null; then
       info "Optional command 'install' not found. Will use 'mv' and 'chmod'."
    fi
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        error "Missing required dependencies: ${missing_deps[*]}. Please install them and try again."
    fi

    # Check for jq (optional but preferred)
    if ! command -v jq &>/dev/null; then
        info "Optional command 'jq' not found. Using fallback for GitHub API parsing (less reliable)."
    fi
    info "Dependencies OK."
}

# Detect Operating System (Linux or macOS only)
get_os() {
    case "$(uname -s)" in
        Darwin*)  OS="darwin" ;;
        Linux*)   OS="linux" ;;
        *) error "Unsupported operating system: $(uname -s). Only Linux and macOS are supported." ;;
    esac
}

# Detect Architecture
get_arch() {
    case "$(uname -m)" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $(uname -m)." ;;
    esac
}

# Get latest release version from GitHub API
get_latest_version() {
    info "Fetching latest release version from GitHub..."
    local api_url="${GITHUB_API}/releases/latest"
    local response
    
    if command -v jq &>/dev/null; then
        # Use jq for reliable parsing
        response=$(curl -fsSL "${api_url}" | jq -r '.tag_name') || {
            error "Failed to fetch latest version using curl/jq. Check network or API rate limits."
        }
    else
        # Fallback using grep/cut (less reliable if tag_name isn't the first)
        response=$(curl -fsSL "${api_url}" | grep '"tag_name":' | head -n 1 | cut -d'"' -f4) || {
             error "Failed to fetch latest version using curl/grep. Check network or API rate limits."
        }
        if [[ -z "$response" ]]; then
            error "Failed to parse latest version using grep/cut. Install 'jq' for better reliability."
        fi
    fi

    if [[ -z "$response" ]]; then
        error "Could not determine latest release version."
    fi
    VERSION="$response"
    info "Latest version: ${VERSION}"
}

# Download file with progress (if not suppressed) and error handling
download_file() {
    local url="$1"
    local dest="$2"
    info "Downloading ${url} to ${dest}"
    # Use -f to fail on server errors, -L to follow redirects, -o to save output
    if curl -fLo "$dest" "$url"; then
        info "Download successful."
    else
        error "Download failed for ${url}"
    fi
}

# Verify checksum
verify_checksum() {
    local archive_path="$1"
    local checksum_file_path="$2"
    local archive_filename
    archive_filename=$(basename "$archive_path")

    info "Verifying checksum for ${archive_filename}..."
    
    # Extract the expected checksum for our specific archive from the checksums file
    # Assumes format like: <checksum> <space(s)> <filename> (or <checksum> *<filename> for sha256sum)
    # Handle potential leading '*' from sha256sum output format if needed when comparing
    local expected_checksum
    expected_checksum=$(grep "${archive_filename}" "$checksum_file_path" | awk '{print $1}')

    if [[ -z "$expected_checksum" ]]; then
        error "Could not find checksum for ${archive_filename} in ${checksum_file_path}"
    fi

    # Calculate checksum of the downloaded file
    local calculated_checksum
    # $CHECKSUM_TOOL might contain arguments (like shasum -a 256), handle correctly
    # Use process substitution to avoid temp files and handle potential tool output variations
    calculated_checksum=$( $CHECKSUM_TOOL < "$archive_path" | awk '{print $1}' )
    
    # Compare checksums
    if [[ "$calculated_checksum" == "$expected_checksum" ]]; then
        info "Checksum verification successful."
    else
        error "Checksum mismatch! Expected '${expected_checksum}', but got '${calculated_checksum}'. Aborting installation."
    fi
}

# Install binary
install_binary() {
    info "Preparing installation..."
    
    # Binary name is always the same now
    local binary_name="${BINARY}" 
    
    # Assume tar.gz format for Linux/macOS
    local archive_ext="tar.gz" 
    local archive_name="${BINARY}-${VERSION}-${OS}-${ARCH}.${archive_ext}"
    local checksum_filename="${BINARY}-${VERSION}-checksums.sha256" # Assumed name

    local download_url="https://github.com/${REPO}/releases/download/${VERSION}/${archive_name}"
    local checksum_url="https://github.com/${REPO}/releases/download/${VERSION}/${checksum_filename}"

    # Create temporary directory
    # Handle potential mktemp differences between Linux/macOS
    TMP_DIR=$(mktemp -d -t hyde-install.XXXXXX) || error "Failed to create temporary directory."
    info "Created temporary directory: ${TMP_DIR}"

    local archive_path="${TMP_DIR}/${archive_name}"
    local checksum_path="${TMP_DIR}/${checksum_filename}"
    local extracted_binary_path="${TMP_DIR}/${binary_name}"

    # Download archive and checksum file
    download_file "$download_url" "$archive_path"
    download_file "$checksum_url" "$checksum_path"

    # Verify checksum
    verify_checksum "$archive_path" "$checksum_path"

    # Extract archive
    info "Extracting ${archive_name}..."
    case "$archive_ext" in
        tar.gz) tar -xzf "$archive_path" -C "$TMP_DIR" || error "Failed to extract archive ${archive_path}" ;;
        # zip) unzip -q "$archive_path" -d "$TMP_DIR" ;; # Keep if zip might be used
        *) error "Unsupported archive extension: ${archive_ext}" ;;
    esac

    if [[ ! -f "$extracted_binary_path" ]]; then
        error "Extraction failed: Binary '${binary_name}' not found in the archive at ${TMP_DIR}."
    fi
    info "Extraction successful."

    # Ensure install directory exists
    info "Ensuring install directory exists: ${INSTALL_DIR}"
    # Use sudo only if we're trying to create dir outside HOME and don't have write perms on parent
    # Check if the target directory exists and is writable OR if its parent exists and is writable
    if [[ ! -d "$INSTALL_DIR" ]]; then
        local parent_dir
        parent_dir=$(dirname "$INSTALL_DIR")
        # Check if we need sudo to create the directory
        if [[ ! "$INSTALL_DIR" =~ ^"$HOME"/ ]] && [[ ! -w "$parent_dir" ]]; then
             info "Attempting to create directory with sudo: ${INSTALL_DIR}"
             if ! sudo mkdir -p "$INSTALL_DIR"; then
                 error "Failed to create install directory: ${INSTALL_DIR}. Try creating it manually or check permissions."
             fi
        else
             # Create directory normally (might fail if permissions are wrong on parent, handled later)
             mkdir -p "$INSTALL_DIR" || error "Failed to create install directory: ${INSTALL_DIR}. Check parent directory permissions."
        fi
    fi

    # Check write permissions for the final install step right before installing
    if [[ ! -w "$INSTALL_DIR" ]]; then
        error "No write permission for install directory: ${INSTALL_DIR}. Please run with 'sudo', check directory permissions, or set HYDE_INSTALL_DIR to a writable location (e.g., in your home directory)."
    fi
    
    local install_path="${INSTALL_DIR}/${binary_name}"
    
    # Install binary using 'install' or fallback to 'mv'/'chmod'
    info "Installing ${binary_name} to ${install_path}..."
    if command -v install &>/dev/null; then
         # Use install command (preferred, handles permissions)
         # Check if sudo is needed to write to the final location
         if [[ ! -w "$INSTALL_DIR" ]] || { [[ -e "$install_path" ]] && [[ ! -w "$install_path" ]]; }; then
             info "Using sudo to install binary..."
             sudo install -m 0755 "$extracted_binary_path" "$install_path" || error "Failed to install binary to ${install_path} using 'sudo install'."
         else
             install -m 0755 "$extracted_binary_path" "$install_path" || error "Failed to install binary to ${install_path} using 'install' command."
         fi
    else
         # Fallback to mv/chmod
         info "Using 'mv' and 'chmod' to install binary..."
         if [[ ! -w "$INSTALL_DIR" ]] || { [[ -e "$install_path" ]] && [[ ! -w "$install_path" ]]; }; then
             info "Using sudo to move and set permissions..."
             sudo mv -f "$extracted_binary_path" "$install_path" || error "Failed to move binary to ${install_path} using 'sudo mv'."
             sudo chmod 755 "$install_path" || error "Failed to make binary executable using 'sudo chmod': ${install_path}."
         else
             mv -f "$extracted_binary_path" "$install_path" || error "Failed to move binary to ${install_path}."
             chmod 755 "$install_path" || error "Failed to make binary executable: ${install_path}."
         fi
    fi

    # Final success message
    echo ""
    info "Successfully installed ${BINARY} (version ${VERSION}) to ${install_path}"
    
    # Check if INSTALL_DIR is in PATH
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*) 
            info "${INSTALL_DIR} is already in your PATH." 
            ;;
        *) 
            echo ""
            echo "==> WARNING: ${INSTALL_DIR} is not in your PATH."
            echo "    You should add it to your PATH environment variable."
            echo "    Add the following line to your shell profile (e.g., ~/.bashrc, ~/.zshrc, ~/.profile):"
            echo ""
            echo "      export PATH=\"${INSTALL_DIR}:\$PATH\""
            echo ""
            echo "    Then, restart your shell or run 'source <your_profile_file>'."
            ;;
    esac
    echo ""
    echo "You can now try running '${BINARY}'."
}

# --- Main Execution ---
main() {
    get_os
    get_arch
    check_deps # Check dependencies after determining OS
    get_latest_version
    install_binary
}

# Run main function
main

exit 0