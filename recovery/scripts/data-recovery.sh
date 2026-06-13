#!/bin/bash
# =============================================================================
# SonicScrewdriver — Data Recovery Tool
# =============================================================================
# Recovers data from damaged or inaccessible drives.
# Part of the Sonic recovery toolkit.
# =============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

echo ""
echo "=============================================="
echo "  SonicScrewdriver — Data Recovery Tool"
echo "=============================================="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root (sudo)"
    exit 1
fi

# List available devices
echo "Available devices:"
lsblk -o NAME,SIZE,TYPE,FSTYPE,LABEL,MOUNTPOINT,MODEL 2>/dev/null || \
    diskutil list 2>/dev/null

echo ""
read -rp "Enter source device (e.g., /dev/sdb): " SOURCE_DEVICE
read -rp "Enter destination directory (e.g., /mnt/recovered): " DEST_DIR

if [[ ! -b "$SOURCE_DEVICE" ]]; then
    log_error "Not a block device: $SOURCE_DEVICE"
    exit 1
fi

mkdir -p "$DEST_DIR"

echo ""
echo "Select recovery method:"
echo "  1) ddrescue (best for failing drives)"
echo "  2) rsync (best for accessible filesystems)"
echo "  3) photorec (file carving)"
echo "  4) testdisk (partition recovery)"
read -rp "Choice [1-4]: " CHOICE

case "$CHOICE" in
    1)
        log_info "Using ddrescue to clone $SOURCE_DEVICE..."
        IMAGE_FILE="${DEST_DIR}/disk-image.dd"
        MAP_FILE="${DEST_DIR}/disk-image.map"
        ddrescue -d -r3 "$SOURCE_DEVICE" "$IMAGE_FILE" "$MAP_FILE" || true
        log_ok "Disk image saved to $IMAGE_FILE"
        ;;
    2)
        log_info "Mounting $SOURCE_DEVICE and using rsync..."
        MOUNT_POINT="/mnt/sonic-source"
        mkdir -p "$MOUNT_POINT"
        mount "$SOURCE_DEVICE" "$MOUNT_POINT" 2>/dev/null || \
            mount -o ro "$SOURCE_DEVICE" "$MOUNT_POINT" || {
            log_error "Failed to mount $SOURCE_DEVICE"
            exit 1
        }
        rsync -avh --progress "$MOUNT_POINT/" "$DEST_DIR/"
        umount "$MOUNT_POINT"
        log_ok "Data recovered to $DEST_DIR"
        ;;
    3)
        log_info "Using photorec for file carving on $SOURCE_DEVICE..."
        photorec /d "$DEST_DIR" /log "$SOURCE_DEVICE" || true
        log_ok "File carving complete. Results in $DEST_DIR"
        ;;
    4)
        log_info "Using testdisk for partition recovery on $SOURCE_DEVICE..."
        testdisk "$SOURCE_DEVICE" || true
        ;;
    *)
        log_error "Invalid choice"
        exit 1
        ;;
esac

log_ok "Recovery operation complete"
