#!/bin/bash

# Set up vendors library and databases
# This script configures the vendors library and databases for sonic and udos.

set -euo pipefail

echo "Setting up vendors library and databases..."

# Create vendors directory
mkdir -p ~/.local/share/sonic/vendors

# Create a sample vendor manifest
cat > ~/.local/share/sonic/vendors/sample-vendor.yaml << 'EOF'
name: sample-vendor
title: Sample Vendor
description: A sample vendor for testing
tags:
  - sample
  - test
container:
  image: alpine:latest
  command: ["sh", "-c", "echo 'Sample vendor running' && sleep infinity"]
EOF

# Create a sample database
cat > ~/.local/udos/vendors.db << 'EOF'
# Sample database for vendors
# This is a placeholder for a SQLite or JSON database
{
  "vendors": [
    {
      "name": "sample-vendor",
      "title": "Sample Vendor",
      "description": "A sample vendor for testing"
    }
  ]
}
EOF

echo "Vendors library and databases setup complete."
echo "Vendors directory: ~/.local/share/sonic/vendors"
echo "Sample vendor manifest created."
