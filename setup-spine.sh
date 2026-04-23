#!/bin/bash

# Setup Universal Spine for uDos
# This script initializes the directory structure and environment variables required for uDos.

set -euo pipefail

echo "Setting up Universal Spine for uDos..."

# Create vault directories
mkdir -p ~/vault/{system,home,family,user,@inbox,@workspace,@toybox,@sandbox,@public,@private,binder}

# Create uDos state directories
mkdir -p ~/.local/udos/{compartments,compost,legacy,trash,feeds}

# Set environment variables
echo 'export UDOS_VAULT="$HOME/vault"' >> ~/.bashrc
echo 'export UDOS_CODE="$HOME/code-vault"' >> ~/.bashrc
echo 'export UDOS_STATE="$HOME/.local/udos"' >> ~/.bashrc

# Reload shell configuration
source ~/.bashrc

echo "Universal Spine setup complete."
echo "Vault directories created at: $UDOS_VAULT"
echo "State directories created at: $UDOS_STATE"
