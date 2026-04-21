#!/bin/bash
set -e

echo "=== Classic Modern Mint Post-Installation ==="

# Apply Classic Modern theme
apply_theme() {
    echo "Applying Classic Modern theme..."
    
    # Install theme files
    sudo mkdir -p /usr/share/themes/ClassicModern
    sudo cp -r assets/themes/ClassicModern/* /usr/share/themes/ClassicModern/
    
    # Install icons
    sudo mkdir -p /usr/share/icons/ClassicModern-Icons
    sudo cp -r assets/icons/* /usr/share/icons/ClassicModern-Icons/
    
    # Set theme
    gsettings set org.cinnamon.desktop.interface gtk-theme "ClassicModern"
    gsettings set org.cinnamon.desktop.interface icon-theme "ClassicModern-Icons"
    gsettings set org.cinnamon.desktop.interface cursor-theme "DMZ-White"
    gsettings set org.cinnamon.desktop.interface font-name "Roboto 11"
    
    # Set wallpaper
    gsettings set org.cinnamon.desktop.background picture-uri "file:///usr/share/backgrounds/classic-modern-default.jpg"
    sudo cp assets/wallpapers/classic-modern-default.jpg /usr/share/backgrounds/
    
    echo "✓ Classic Modern theme applied"
}

# Install Sonic Family tools
install_sonic_tools() {
    echo "Installing Sonic Family tools..."
    
    # Install Sonic Screwdriver
    sudo cp bin/sonic /usr/local/bin/
    sudo chmod +x /usr/local/bin/sonic
    
    # Create config directory
    sudo mkdir -p /etc/sonic
    sudo cp config/sonic-config.json /etc/sonic/config.json
    
    # Set up Go environment in uDos vendor
    sudo mkdir -p /home/$USER/uDos/vendor/go
    echo "export GOPATH=/home/$USER/uDos/vendor/go" | sudo tee -a /etc/profile.d/sonic-go.sh
    echo "export PATH=$PATH:$GOPATH/bin" | sudo tee -a /etc/profile.d/sonic-go.sh
    
    # Set up systemd service
    sudo cp scripts/sonic.service /etc/systemd/system/
    sudo systemctl enable sonic.service
    
    echo "✓ Sonic Family tools installed"
}

# Configure Ventoy integration
configure_ventoy() {
    echo "Configuring Ventoy integration..."
    
    # Copy Ventoy config
    sudo mkdir -p /etc/ventoy
    sudo cp config/ventoy.json /etc/ventoy/config.json
    
    # Install Ventoy themes
    sudo mkdir -p /usr/share/ventoy/themes
    sudo cp -r ventoy/themes/* /usr/share/ventoy/themes/
    
    echo "✓ Ventoy integration configured"
}

# Clean up
cleanup() {
    echo "Cleaning up..."
    sudo apt-get autoremove -y
    sudo apt-get clean
    echo "✓ Cleanup complete"
}

# Main execution
main() {
    apply_theme
    install_sonic_tools
    configure_ventoy
    cleanup
    echo "✓ Post-installation complete"
    echo "Please reboot to apply all changes."
}

main "$@"