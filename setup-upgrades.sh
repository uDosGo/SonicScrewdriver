#!/bin/bash

# Set up upgrade and local library management for sonic and udos
# This script configures the upgrade process and local library for sonic and udos.

set -euo pipefail

echo "Setting up upgrade and local library management..."

# Create local library directory
mkdir -p ~/.local/share/sonic/library

# Configure sonic to use local library
cat > ~/.config/sonic/library.yaml << 'EOF'
local_library_path: ~/.local/share/sonic/library
auto_update: true
EOF

# Create a sample game manifest for testing
cat > ~/.local/share/sonic/library/test-game.yaml << 'EOF'
name: test-game
title: Test Game
description: A test game for verifying the installation
tags:
  - test
  - example
container:
  image: alpine:latest
  command: ["sh", "-c", "echo 'Test game running' && sleep infinity"]
EOF

echo "Upgrade and local library management setup complete."
echo "Local library: ~/.local/share/sonic/library"
echo "Test game manifest created."
