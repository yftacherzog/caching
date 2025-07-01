#!/bin/bash
set -e

# --- Failsafe and OS-specific Setup ---
OS=$(uname -s)
SOCKET_PATH=""

if [ "$OS" == "Linux" ]; then
    SOCKET_PATH="${XDG_RUNTIME_DIR}/podman/podman.sock"
elif [ "$OS" == "Darwin" ]; then
    echo "Error: macOS is not supported by this dev container configuration." >&2
    echo "This setup depends on Linux-specific features like host networking and user namespace mapping (--userns=keep-id)." >&2
    echo "It also assumes the standard Linux rootless Podman socket path (\${XDG_RUNTIME_DIR}/podman/podman.sock)." >&2
    echo "To add macOS support, this script would need to be modified to handle the Podman Machine VM." >&2
    exit 1
else
    echo "Error: Unsupported operating system '$OS'." >&2
    exit 1
fi

if [ -z "$SOCKET_PATH" ] || [ ! -S "$SOCKET_PATH" ]; then
    echo "Error: Podman socket not found or not running at the expected path: ${SOCKET_PATH}" >&2
    echo "For Fedora/RHEL, please ensure the Podman service is running." >&2
    echo "You can typically start it with: systemctl --user start podman.socket" >&2
    exit 1
fi
# --- End Failsafe ---

# This script is called by `initializeCommand` in devcontainer.json.
# It populates a .devcontainer.env file with host-specific information,
# which is then used by the container.

ENV_FILE=".devcontainer/devcontainer.env"

# Clear the file to start fresh
> "$ENV_FILE"

# Write git config to the env file so it's available inside the container
echo "GIT_AUTHOR_NAME=$(git config --get user.name)" >> "$ENV_FILE"
echo "GIT_AUTHOR_EMAIL=$(git config --get user.email)" >> "$ENV_FILE"
echo "GIT_COMMITTER_NAME=$(git config --get user.name)" >> "$ENV_FILE"
echo "GIT_COMMITTER_EMAIL=$(git config --get user.email)" >> "$ENV_FILE" 
