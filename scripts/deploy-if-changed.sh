#!/bin/bash

echo "$(date): Fetching remote repository..."
git fetch

UPSTREAM=${1:-'@{u}'}
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

if [ $LOCAL = $REMOTE ]; then
    echo "$(date): No changes detected in git"
elif [ $LOCAL = $BASE ]; then
    BUILD_VERSION=$(git rev-parse HEAD)
    echo "$(date): Changes detected, deploying new version: $BUILD_VERSION"
    sudo ./scripts/deploy.sh
elif [ $REMOTE = $BASE ]; then
    echo "$(date): Local changes detected, stashing"
    git stash
    sudo ./scripts/deploy.sh
else
    echo "$(date): Git is diverged, this is unexpected."
fi
