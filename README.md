# caching

Caching infrastructure for Konflux clusters

## Development Environment

This repository includes a dev container configuration that provides a consistent development environment. To use it, you need the following on your local machine:

1. **Podman**: The dev container is based on Podman. Ensure it is installed on your system.
2. **VS Code with Dev Containers extension**: For the best experience, use Visual Studio Code with the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
3. **Running Podman Socket**: The dev container connects to your local Podman socket. Make sure the user service is running before you start the container:

    ```bash
    systemctl --user start podman.socket
    ```

    To enable the socket permanently so it starts on boot, you can run:

    ```bash
    systemctl --user enable --now podman.socket
    ```

Once these prerequisites are met, you can open this folder in VS Code and use the "Reopen in Container" command to launch the environment.
