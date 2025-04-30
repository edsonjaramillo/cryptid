#!/usr/bin/env bash

# Strict mode
set -euo pipefail

# Constants
readonly REPO="edsonjaramillo/hyde"
readonly BINARY="hyde"
readonly COMPLETIONS_FILENAME="hyde_completions.bash" # Bash completions filename
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

# --- Helper Functions ---

# Basic logging
info() {
    echo "INFO: $1"
}

error() {
    echo "ERROR: $1" >&2
    exit 1
}

warning() {
    echo "WARNING: $1" >&2
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
    # Removed tar, sha256sum/shasum. Added install (more strongly recommended now).
    local deps=("curl" "mktemp" "uname" "mkdir" "mv" "chmod" "rm" "id" "grep" "cut" "install")

    # Check basic dependencies
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &>/dev/null; then
            missing_deps+=("$dep")
        fi
    done

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        # 'install' is crucial now for setting permissions correctly, especially for completions.
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
    info "Downloading ${url} -> ${dest}"
    # Use -f to fail on server errors, -L to follow redirects, -o to save output
    if curl -fLo "$dest" "$url"; then
        info "Download successful: $(basename "$dest")"
    else
        error "Download failed for ${url}"
    fi
}

# Determine appropriate Bash completions directory
get_bash_completion_dir() {
    local completions_dir=""
    # Check if running as root or installing to system location
    if [[ "$(id -u)" -eq 0 ]] || [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
        # System-wide completions directories (check common locations)
        if [[ -d "/usr/share/bash-completion/completions" ]]; then
            completions_dir="/usr/share/bash-completion/completions"
        elif [[ -d "/etc/bash_completion.d" ]]; then
            completions_dir="/etc/bash_completion.d"
        else
             warning "Could not find standard system-wide bash completions directory."
             warning "Checked /usr/share/bash-completion/completions and /etc/bash_completion.d"
             warning "Bash completions will not be installed."
             echo "" # Return empty string
             return
        fi
    else
        # User-specific completions directory (use XDG standard if possible)
        local user_data_dir="${XDG_DATA_HOME:-$HOME/.local/share}"
        if [[ -n "$user_data_dir" ]]; then
             completions_dir="${user_data_dir}/bash-completion/completions"
        else
             warning "Could not determine user data directory (XDG_DATA_HOME or ~/.local/share)."
             warning "Bash completions will not be installed."
             echo "" # Return empty string
             return
        fi
    fi
    echo "$completions_dir"
}


# Ensure directory exists, creating if necessary (using sudo if needed)
ensure_dir_exists() {
    local dir_path="$1"
    local use_sudo=false

    # Check if directory exists
    if [[ -d "$dir_path" ]]; then
        # Check if it's writable
        if [[ ! -w "$dir_path" ]]; then
            use_sudo=true # Might need sudo to write file later, even if dir exists
        fi
    else
        # Directory doesn't exist, check if we need sudo to create its parent
        local parent_dir
        parent_dir=$(dirname "$dir_path")
        # Need sudo if parent isn't writable OR if installing outside home dir as non-root
        if [[ ! -w "$parent_dir" ]] || { [[ "$(id -u)" -ne 0 ]] && [[ ! "$dir_path" =~ ^"$HOME"/ ]]; }; then
           use_sudo=true
        fi
    fi

    # Create the directory
    if [[ "$use_sudo" == true ]]; then
        info "Attempting to create/ensure directory with sudo: ${dir_path}"
        if ! sudo mkdir -p "$dir_path"; then
            error "Failed to create directory (with sudo): ${dir_path}. Check permissions or create it manually."
        fi
    else
        info "Ensuring directory exists: ${dir_path}"
        if ! mkdir -p "$dir_path"; then
             error "Failed to create directory: ${dir_path}. Check permissions."
        fi
    fi

    # Return whether sudo might be needed for file operations within this dir
    if [[ ! -w "$dir_path" ]]; then
        echo "sudo_required"
    else
        echo "writable"
    fi
}

# Install file using 'install' command (handles permissions and sudo)
install_file() {
    local source_path="$1"
    local dest_path="$2"
    local permissions="$3" # e.g., 0755 for executable, 0644 for data file
    local dest_dir_status="$4" # "sudo_required" or "writable"

    info "Installing $(basename "$source_path") to ${dest_path} with permissions ${permissions}..."

    if [[ "$dest_dir_status" == "sudo_required" ]]; then
        info "Using sudo to install file..."
        if ! sudo install -m "$permissions" "$source_path" "$dest_path"; then
            error "Failed to install file using 'sudo install': ${dest_path}"
        fi
    else
         if ! install -m "$permissions" "$source_path" "$dest_path"; then
             # If install fails without sudo, maybe permissions changed? Or dest file exists and isn't writable? Try sudo.
             warning "Install command failed. Retrying with sudo..."
             if ! sudo install -m "$permissions" "$source_path" "$dest_path"; then
                  error "Failed to install file (even with sudo): ${dest_path}"
             fi
         fi
    fi
    info "Successfully installed $(basename "$dest_path")"
}


# Main installation function
install_hyde() {
    info "Preparing installation..."

    # Binary source filename on GitHub releases
    local binary_source_filename="${BINARY}-${VERSION}-${OS}-${ARCH}"
    local binary_download_url="https://github.com/${REPO}/releases/download/${VERSION}/${binary_source_filename}"
    local completions_download_url="https://github.com/${REPO}/releases/download/${VERSION}/${COMPLETIONS_FILENAME}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d -t hyde-install.XXXXXX) || error "Failed to create temporary directory."
    info "Created temporary directory: ${TMP_DIR}"

    local downloaded_binary_path="${TMP_DIR}/${binary_source_filename}"
    local downloaded_completions_path="${TMP_DIR}/${COMPLETIONS_FILENAME}"

    # --- Download Files ---
    download_file "$binary_download_url" "$downloaded_binary_path"
    download_file "$completions_download_url" "$downloaded_completions_path"

    # Basic check if downloaded files exist
    if [[ ! -f "$downloaded_binary_path" ]]; then
        error "Downloaded binary not found at ${downloaded_binary_path}. Aborting."
    fi
    if [[ ! -f "$downloaded_completions_path" ]]; then
        # Warn but continue? Or error out? Let's warn for now.
        warning "Downloaded completions file not found at ${downloaded_completions_path}. Skipping completions installation."
        local completions_available=false
    else
        local completions_available=true
    fi

    # --- Install Binary ---
    local binary_install_path="${INSTALL_DIR}/${BINARY}" # Final name is just 'hyde'
    info "Preparing to install binary to ${binary_install_path}..."
    local bin_dir_status
    bin_dir_status=$(ensure_dir_exists "$INSTALL_DIR")
    install_file "$downloaded_binary_path" "$binary_install_path" "0755" "$bin_dir_status" # Executable permissions


    # --- Install Bash Completions ---
    local completions_install_path=""
    if [[ "$completions_available" == true ]]; then
        local target_completions_dir
        target_completions_dir=$(get_bash_completion_dir)

        if [[ -n "$target_completions_dir" ]]; then
            completions_install_path="${target_completions_dir}/${BINARY}" # Standard naming convention
            info "Preparing to install bash completions to ${completions_install_path}..."
            local comp_dir_status
            comp_dir_status=$(ensure_dir_exists "$target_completions_dir")
            install_file "$downloaded_completions_path" "$completions_install_path" "0644" "$comp_dir_status" # Read permissions
        else
            warning "Skipping Bash completions installation as target directory could not be determined or found."
            warning "You might need to install the 'bash-completion' package."
        fi
    fi

    # --- Final Messages ---
    echo ""
    info "--------------------------------------------------"
    info "${BINARY} (version ${VERSION}) installed successfully!"
    info "  Binary: ${binary_install_path}"
    if [[ -n "$completions_install_path" ]]; then
       info "  Bash Completions: ${completions_install_path}"
    fi
    info "--------------------------------------------------"
    echo ""

    # Check if binary install directory is in PATH
    case ":$PATH:" in
        *":${INSTALL_DIR}:"*)
            info "${INSTALL_DIR} is already in your PATH."
            ;;
        *)
            warning "${INSTALL_DIR} is not in your PATH."
            warning "You should add it to your PATH environment variable."
            warning "Add the following line to your shell profile (e.g., ~/.bashrc, ~/.zshrc, ~/.profile):"
            warning ""
            warning "  export PATH=\"${INSTALL_DIR}:\$PATH\""
            warning ""
            warning "Then, restart your shell or run 'source <your_profile_file>'."
            ;;
    esac

    # Suggest restarting shell for completions
    if [[ -n "$completions_install_path" ]]; then
        echo ""
        info "Bash completions have been installed."
        info "Please restart your shell or run 'source ${completions_install_path}' for completions to take effect."
        info "(Actual sourcing mechanism might depend on your 'bash-completion' setup)."
    fi
    echo ""
    echo "You can now try running '${BINARY}'."
}

# --- Main Execution ---
main() {
    get_os
    get_arch
    check_deps # Check dependencies after determining OS
    get_latest_version
    install_hyde
}

# Run main function
main

exit 0