#!/usr/bin/env bash

# Define the lock file to prevent concurrent runs
LOCK_FILE="${PWD}${1}/deploy.lock"

# Execute the deployment script with a lock to prevent conflicts
flock -n $LOCK_FILE ${PWD}${1}/scripts/deploy-if-changed.sh ${1}