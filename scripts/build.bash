#!/bin/bash

set -e

# Get the first argument as the version
version="$1"

# Check if version matches the format v[0-9].[0-9].[0-9]
if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid version number. Please use the format v1.0.0"
    exit 1
fi

# Package name and output name
package_name="cryptid"
output_name="${package_name}-${version}"

# Platforms to build for
platforms=(
    "darwin/amd64" "darwin/arm64"
    "linux/amd64" "linux/arm64"
    "freebsd/amd64" "freebsd/arm64"
    "windows/amd64" "windows/arm64"
)

# Ensure the dist directory exists
mkdir -p dist

build() {
    local platform="$1"
    local platform_split=(${platform//\// })
    local GOOS="${platform_split[0]}"
    local GOARCH="${platform_split[1]}"
    
    # echo "Building for $GOOS-$GOARCH"
    
    local binary_dir="dist/binaries/$GOOS-$GOARCH"
    local binary="$package_name"
    local tar_name="dist/${output_name}-${GOOS}-${GOARCH}.tar.gz"
    
    if [[ "$GOOS" == "windows" ]]; then
        binary+=".exe"
    fi
    
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -o "$binary_dir/$binary" "cmd/$package_name.go"
    chmod +x "$binary_dir/$binary"
    
    tar -czf "$tar_name" -C "$binary_dir" "$binary"
}

export -f build
export output_name
export package_name

# Build in parallel
printf "%s\n" "${platforms[@]}" | xargs -n 1 -P 4 -I {} bash -c 'build "{}"' && rm -rf dist/binaries

echo "Build completed successfully."