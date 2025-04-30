#!/usr/bin/env bash

# Package name
PACKAGE_NAME="hyde"

PLATFORMS=(
    "darwin/amd64" "darwin/arm64"
    "linux/amd64" "linux/arm64"
)

# Ensure the dist directory exists
echo "Ensuring 'dist' directory exists..."
mkdir -p dist

build() {
    echo "Starting build process..."
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r -a platform_split <<<"$platform"
        local GOOS="${platform_split[0]}"
        local GOARCH="${platform_split[1]}"

        # Define the output binary path and name
        local output_binary="dist/${PACKAGE_NAME}-${GOOS}-${GOARCH}"

        echo "Building for ${GOOS}/${GOARCH} -> ${output_binary}"

        # Build the binary directly into the dist directory

        env GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w" -o "$output_binary" "backend/cmd/cli/main.go"

        # Check if the build was successful
        if [ $? -ne 0 ]; then
            echo "Error: Build failed for ${GOOS}/${GOARCH}"
            exit 1
        else
            echo "Successfully built ${output_binary}"
        fi
        echo # Add a newline for readability
    done

    echo "Build process completed."
    echo "Binaries are located in the 'dist' directory:"
    ls -l dist/
}

# Run the build function
build "$@"

exit 0