#!/bin/bash

# Ensure a directory is provided

# Change to the specified directory
cd "${PWD}${1}" || { echo "Failed to change directory to $1"; exit 1; }

echo "$(date): Fetching remote repository..."
git fetch || { echo "Failed to fetch remote repository"; exit 1; }

UPSTREAM="@{u}"
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

echo "LOCAL: $LOCAL"
echo "REMOTE: $REMOTE"
echo "BASE: $BASE"

if [ "$LOCAL" = "$REMOTE" ]; then
    echo "$(date): No changes detected in git"
elif [ "$LOCAL" = "$BASE" ]; then
    BUILD_VERSION=$(git rev-parse HEAD)
    echo "$(date): Changes detected, deploying new version: $BUILD_VERSION"
    ./scripts/deploy.sh "$1"
elif [ "$REMOTE" = "$BASE" ]; then
    echo "$(date): Local changes detected, stashing"
    git stash
    ./scripts/deploy.sh "$1"
else
    echo "$(date): Git is diverged, this is unexpected."
fi

