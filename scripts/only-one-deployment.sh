#!/usr/bin/env bash

# Define the lock file to prevent concurrent runs
LOCK_FILE="${PWD}/deploy.lock"

# Execute the deployment script with a lock to prevent conflicts
flock -n $LOCK_FILE sudo ./scripts/deploy-if-changed.sh