#!/usr/bin/env bash

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

IS_GSED_INSTALLED=false
IS_SED_INSTALLED=false
# make sure at least one of gsed or sed is installed
if ! command -v gsed &> /dev/null; then
    IS_GSED_INSTALLED=false
fi

if ! command -v sed &> /dev/null; then
    IS_SED_INSTALLED=false
fi

if ! $IS_GSED_INSTALLED && ! $IS_SED_INSTALLED; then
    echo "Error: gsed or sed is not installed"
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
    if $IS_GSED_INSTALLED; then 
        gsed -i "s/$OLD_VERSION/$NEW_VERSION/g" "$FILE"
    else
        sed -i "s/$OLD_VERSION/$NEW_VERSION/g" "$FILE"
    fi
done

printf "Version updated from %s to %s in all files\n" "$OLD_VERSION" "$NEW_VERSION"