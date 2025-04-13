#!/bin/bash

# Check if old version argument is provided 
if [ -z "$1" ]; then
    echo "Error: Old version number is required"
    echo "Usage: $0 old_version new_version"
    exit 1
fi

# Check if new version argument is provided
if [ -z "$2" ]; then
    echo "Error: New version number is required"
    echo "Usage: $0 old_version new_version"
    exit 1
fi
# Validate version format 
if ! [[ $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Invalid version format. Must be in format x.x.x"
    exit 1
fi

if ! [[ $2 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Invalid version format. Must be in format x.x.x"
    exit 1
fi

# Store the version
OLD_VERSION=$1
NEW_VERSION=$2

# show me the files that will be updated
FILES_TO_UPDATE=(
	"backend/cmd/cli/main.go"
	"Makefile"
)

for FILE in "${FILES_TO_UPDATE[@]}"; do
    gsed -i "s/$OLD_VERSION/$NEW_VERSION/g" "$FILE"
done

printf "Version updated from %s to %s in all files\n" "$OLD_VERSION" "$NEW_VERSION"