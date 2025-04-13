#!/bin/bash

# Get the first argument as the VERSION
VERSION="$1"

# Check if VERSION matches the format [0-9].[0-9].[0-9]
if [ -z "$VERSION" ]; then
    echo "Error: VERSION number is required"
    echo "Usage: $0 version_number"
    exit 1
fi

# Package name and output name
PACKAGE_NAME="cryptid"
OUTPUT_NAME="${PACKAGE_NAME}-${VERSION}"

# Platforms to build for
PLATFORMS=(
    "darwin/amd64" "darwin/arm64"
    # "linux/amd64" "linux/arm64"
    # "windows/amd64" "windows/arm64"
)

# Ensure the dist directory exists
mkdir -p dist

build() {
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r -a platform_split <<<"$platform"
        local GOOS="${platform_split[0]}"
        local GOARCH="${platform_split[1]}"

        local binary_dir="dist/binaries/$GOOS-$GOARCH"
        local binary="$PACKAGE_NAME"
        local tar_name="dist/${OUTPUT_NAME}-${GOOS}-${GOARCH}.tar.gz"

        if [[ "$GOOS" == "windows" ]]; then
            binary+=".exe"
        fi

        env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$binary_dir/$binary" "backend/cmd/cli/main.go"

        tar -czf "$tar_name" -C "$binary_dir" "$binary"
    done

    rm -rf dist/binaries
}

build "$@"    
