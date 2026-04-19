#!/bin/bash
set -e

echo "=== Classic Modern Mint Pre-Installation ==="

# Check system requirements
check_requirements() {
    echo "Checking system requirements..."
    
    # Check disk space
    if [ $(df / --output=avail | tail -1) -lt 21474836480 ]; then
        echo "ERROR: Insufficient disk space. Minimum 20GB required."
        exit 1
    fi
    
    # Check memory
    if [ $(free -m | awk '/Mem:/{print $2}') -lt 4096 ]; then
        echo "WARNING: Low memory detected. 4GB+ recommended for best performance."
    fi
    
    echo "✓ System requirements met"
}

# Backup existing configuration
backup_config() {
    echo "Backing up existing configuration..."
    if [ -d /etc/sonic ]; then
        timestamp=$(date +%Y%m%d_%H%M%S)
        backup_dir="/var/backups/sonic_pre_install_$timestamp"
        sudo mkdir -p "$backup_dir"
        sudo cp -r /etc/sonic "$backup_dir/"
        echo "✓ Configuration backed up to $backup_dir"
    fi
}

# Install prerequisites
install_prerequisites() {
    echo "Installing prerequisites..."
    sudo apt-get update
    sudo apt-get install -y \
        curl \
        wget \
        git \
        gnupg \
        lsb-release \
        ca-certificates \
        apt-transport-https \
        software-properties-common
    echo "✓ Prerequisites installed"
}

# Main execution
main() {
    check_requirements
    backup_config
    install_prerequisites
    echo "✓ Pre-installation complete"
}

main "$@"