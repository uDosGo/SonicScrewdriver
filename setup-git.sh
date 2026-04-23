#!/bin/bash

# Configure git operations for vault and code-vault
# This script sets up git repositories for the vault and code-vault directories.

set -euo pipefail

echo "Configuring git operations for vault and code-vault..."

# Initialize git repositories
cd ~/vault && git init
cd ~/code-vault && git init

# Create .gitignore files
cat > ~/vault/.gitignore << 'EOF'
# Ignore sensitive and temporary files
@private/
trash/
compost/
*.log
*.tmp
EOF

cat > ~/code-vault/.gitignore << 'EOF'
# Ignore build artifacts and temporary files
node_modules/
dist/
*.log
*.tmp
EOF

# Configure git user
git config --global user.name "Wizard"
git config --global user.email "wizard@example.com"

echo "Git operations configured."
echo "Vault repository: ~/vault"
echo "Code-vault repository: ~/code-vault"
