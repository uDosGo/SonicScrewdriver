#!/bin/bash

# Install sonic and udos
# This script installs the sonic and udos binaries and sets up the environment.

set -euo pipefail

echo "Installing sonic and udos..."

# Install sonic
if [ ! -f ~/.local/bin/sonic ]; then
    echo "Installing sonic..."
    curl -sSL https://raw.githubusercontent.com/sonic-family/installer/main/bootstrap.sh -o /tmp/bootstrap.sh
    chmod +x /tmp/bootstrap.sh
    /tmp/bootstrap.sh
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
    source ~/.bashrc
fi

# Install udos
if [ ! -f ~/.local/bin/udos ]; then
    echo "Installing udos..."
    sonic install uDos
    sudo ln -sf ~/.local/udos/bin/udos /usr/local/bin/udos
fi

echo "Installation complete."
echo "sonic version: $(sonic --version)"
echo "udos version: $(udos version)"
