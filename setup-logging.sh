#!/bin/bash

# Set up logging and error handling for sonic and udos
# This script configures logging directories and error handling for the sonic and udos binaries.

set -euo pipefail

echo "Setting up logging and error handling..."

# Create logging directories
mkdir -p ~/.local/share/sonic/logs
mkdir -p ~/.local/udos/logs

# Configure sonic logging
cat > ~/.config/sonic/config.yaml << 'EOF'
logging:
  level: info
  file: ~/.local/share/sonic/logs/sonic.log
  max_size: 10
  max_backups: 7
  max_age: 30
EOF

# Configure udos logging
cat > ~/.config/udos/config.yaml << 'EOF'
logging:
  level: info
  file: ~/.local/udos/logs/udos.log
  max_size: 10
  max_backups: 7
  max_age: 30
EOF

echo "Logging and error handling setup complete."
echo "Sonic logs: ~/.local/share/sonic/logs/sonic.log"
echo "uDos logs: ~/.local/udos/logs/udos.log"
