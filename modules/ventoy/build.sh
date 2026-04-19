#!/usr/bin/env sh
set -eu

# Ventoy build script for Sonic Family installers
# Supports creating bootable USB images with Classic Modern Mint

COMMAND="${1:-help}"
INSTALLER="${2:-classic-modern-mint}"
OUTPUT="${3:-/tmp}"

show_help() {
    echo "Ventoy Build Script for Sonic Family"
    echo ""
    echo "Usage: $0 <command> [installer] [output]"
    echo ""
    echo "Commands:"
    echo "  package       - Create installer package (default: classic-modern-mint)"
    echo "  usb           - Create bootable USB image"
    echo "  test          - Test in QEMU"
    echo "  clean         - Clean build artifacts"
    echo "  help          - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 package classic-modern-mint /tmp"
    echo "  $0 usb /dev/sdX classic-modern-mint.she"
    echo "  $0 test classic-modern-mint.she"
}

package_installer() {
    echo "=== Packaging $INSTALLER installer ==="
    
    SOURCE_DIR="installers/$INSTALLER"
    if [ ! -d "$SOURCE_DIR" ]; then
        echo "ERROR: Installer source not found: $SOURCE_DIR"
        exit 1
    fi
    
    echo "Creating bundle..."
    go run cmd/sonic/main.go ventoy package "$SOURCE_DIR" "$OUTPUT/$INSTALLER.she"
    
    if [ $? -eq 0 ]; then
        echo "✓ Bundle created: $OUTPUT/$INSTALLER.she"
    else
        echo "ERROR: Failed to create bundle"
        exit 1
    fi
}

create_usb() {
    USB_DEVICE="${2:-}"
    BUNDLE="${3:-$OUTPUT/$INSTALLER.she}"
    
    if [ -z "$USB_DEVICE" ]; then
        echo "ERROR: USB device not specified"
        echo "Usage: $0 usb /dev/sdX [bundle.she]"
        exit 1
    fi
    
    if [ ! -f "$BUNDLE" ]; then
        echo "ERROR: Bundle not found: $BUNDLE"
        exit 1
    fi
    
    echo "=== Creating bootable USB on $USB_DEVICE ==="
    echo "WARNING: This will erase all data on $USB_DEVICE"
    
    # Check if Ventoy is available
    if ! command -v ventoy &> /dev/null; then
        echo "ERROR: Ventoy not found. Please install Ventoy first."
        echo "Download: https://www.ventoy.net"
        exit 1
    fi
    
    # Install Ventoy on USB device
    echo "Installing Ventoy..."
    ventoy -i "$USB_DEVICE"
    
    if [ $? -ne 0 ]; then
        echo "ERROR: Ventoy installation failed"
        exit 1
    fi
    
    # Copy bundle to USB
    echo "Copying bundle to USB..."
    cp "$BUNDLE" "/Volumes/VENTOY/"
    cp "$SOURCE_DIR/config/ventoy.json" "/Volumes/VENTOY/"
    
    # Copy Ventoy theme
    if [ -d "$SOURCE_DIR/ventoy/themes" ]; then
        cp -r "$SOURCE_DIR/ventoy/themes" "/Volumes/VENTOY/"
    fi
    
    echo "✓ Bootable USB created successfully"
    echo "You can now boot from $USB_DEVICE"
}

test_qemu() {
    BUNDLE="${2:-$OUTPUT/$INSTALLER.she}"
    
    if [ ! -f "$BUNDLE" ]; then
        echo "ERROR: Bundle not found: $BUNDLE"
        exit 1
    fi
    
    echo "=== Testing in QEMU ==="
    
    if ! command -v qemu-system-x86_64 &> /dev/null; then
        echo "ERROR: QEMU not found. Please install QEMU first."
        exit 1
    fi
    
    # Create temporary disk image
    DISK_IMG="/tmp/ventoy_test.img"
    if [ ! -f "$DISK_IMG" ]; then
        echo "Creating test disk image (this may take a while)..."
        qemu-img create -f raw "$DISK_IMG" 4G
        ventoy -i "$DISK_IMG"
        # Mount and copy files (simplified for example)
        echo "Test setup complete. Run manually:"
        echo "qemu-system-x86_64 -m 4G -boot d -cdrom $DISK_IMG"
    else
        echo "Using existing test image: $DISK_IMG"
        echo "Run: qemu-system-x86_64 -m 4G -boot d -cdrom $DISK_IMG"
    fi
}

clean_build() {
    echo "=== Cleaning build artifacts ==="
    rm -f "$OUTPUT/$INSTALLER.she"
    rm -f /tmp/ventoy_test.img
    echo "✓ Cleanup complete"
}

# Main execution
case "$COMMAND" in
    package|pkg)
        package_installer
        ;;
    usb)
        create_usb "$@"
        ;;
    test)
        test_qemu "$@"
        ;;
    clean)
        clean_build
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "ERROR: Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac
