#!/bin/bash

# Seed vault and code-vault with initial data
# This script adds initial data to the vault and code-vault directories.

set -euo pipefail

echo "Seeding vault and code-vault with initial data..."

# Seed vault with example files
cat > ~/vault/@workspace/example.md << 'EOF'
# Example Document

This is an example document to seed the vault.

## Features
- Markdown support
- Tagging
- Search

## Tags
#example #seed
EOF

cat > ~/vault/binder/example.txt << 'EOF'
Example binder entry for testing.
EOF

# Seed code-vault with example files
cat > ~/code-vault/example-script.sh << 'EOF'
#!/bin/bash
# Example script for code-vault
echo "Hello, World!"
EOF

chmod +x ~/code-vault/example-script.sh

echo "Data seeding complete."
echo "Example files created in vault and code-vault."
