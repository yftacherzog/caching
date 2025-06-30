#!/bin/bash
set -e

# This script is called by `initializeCommand` in devcontainer.json.
# It populates a .devcontainer.env file with the host's git user info,
# which is then used to set environment variables inside the container.

ENV_FILE=".devcontainer/devcontainer.env"

# Clear the file to start fresh
> "$ENV_FILE"

# Write git config values to the env file
echo "GIT_AUTHOR_NAME=$(git config --get user.name)" >> "$ENV_FILE"
echo "GIT_AUTHOR_EMAIL=$(git config --get user.email)" >> "$ENV_FILE"
echo "GIT_COMMITTER_NAME=$(git config --get user.name)" >> "$ENV_FILE"
echo "GIT_COMMITTER_EMAIL=$(git config --get user.email)" >> "$ENV_FILE" 