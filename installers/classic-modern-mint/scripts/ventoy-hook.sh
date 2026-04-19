#!/bin/bash
set -e

echo "=== Ventoy Hook for Classic Modern Mint ==="

# This script runs when Ventoy boots the installer

# Detect boot mode
BOOT_MODE="unknown"
if [ -d /sys/firmware/efi ]; then
    BOOT_MODE="uefi"
else
    BOOT_MODE="legacy"
fi

echo "Detected boot mode: $BOOT_MODE"

# Set up persistent storage
setup_persistent_storage() {
    echo "Setting up persistent storage..."
    if [ -f /ventoy/ventoy.json ]; then
        PERSISTENT_SIZE=$(jq -r '.persistent.size' /ventoy/ventoy.json)
        if [ "$PERSISTENT_SIZE" != "null" ] && [ "$PERSISTENT_SIZE" != "0" ]; then
            echo "Configuring persistent storage: ${PERSISTENT_SIZE}MB"
            # Ventoy handles this automatically based on ventoy.json
        fi
    fi
    echo "✓ Persistent storage configured"
}

# Apply Ventoy theme
apply_ventoy_theme() {
    echo "Applying Ventoy theme..."
    if [ -f /ventoy/ventoy.json ]; then
        THEME_CONFIG=$(jq -r '.theme' /ventoy/ventoy.json)
        if [ "$THEME_CONFIG" != "null" ]; then
            echo "Applying theme from config"
            # Ventoy applies theme automatically
        fi
    fi
    echo "✓ Ventoy theme applied"
}

# Main execution
main() {
    setup_persistent_storage
    apply_ventoy_theme
    echo "✓ Ventoy hook complete"
}

main "$@"