#!/bin/bash
# =============================================================================
# SonicScrewdriver — Password Reset Tool
# =============================================================================
# Resets forgotten passwords on Linux systems by mounting the root partition
# and chrooting in. Part of the Sonic recovery toolkit.
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
echo "  SonicScrewdriver — Password Reset Tool"
echo "=============================================="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root (sudo)"
    exit 1
fi

# List available partitions
echo "Available partitions:"
lsblk -o NAME,SIZE,TYPE,FSTYPE,LABEL,MOUNTPOINT 2>/dev/null || \
    diskutil list 2>/dev/null

echo ""
read -rp "Enter root partition (e.g., /dev/sda2): " PARTITION

if [[ ! -b "$PARTITION" ]]; then
    log_error "Not a block device: $PARTITION"
    exit 1
fi

# Mount the partition
MOUNT_POINT="/mnt/sonic-recovery"
mkdir -p "$MOUNT_POINT"

log_info "Mounting $PARTITION to $MOUNT_POINT..."
mount "$PARTITION" "$MOUNT_POINT" || {
    log_error "Failed to mount $PARTITION"
    exit 1
}

# Check for Linux filesystem
if [[ ! -f "${MOUNT_POINT}/etc/passwd" ]]; then
    log_error "No Linux installation found on $PARTITION"
    umount "$MOUNT_POINT"
    exit 1
fi

echo ""
log_info "Linux installation detected!"
echo ""

# List users
echo "Available users:"
awk -F: '$3 >= 1000 && $3 < 65534 {print $1}' "${MOUNT_POINT}/etc/passwd"
echo ""

read -rp "Enter username to reset password for: " USERNAME

if ! grep -q "^${USERNAME}:" "${MOUNT_POINT}/etc/passwd"; then
    log_error "User '$USERNAME' not found"
    umount "$MOUNT_POINT"
    exit 1
fi

# Chroot and reset password
log_info "Resetting password for $USERNAME..."

mount --bind /dev "${MOUNT_POINT}/dev"
mount --bind /proc "${MOUNT_POINT}/proc"
mount --bind /sys "${MOUNT_POINT}/sys"

chroot "$MOUNT_POINT" /bin/bash -c "passwd $USERNAME"

umount "${MOUNT_POINT}/sys"
umount "${MOUNT_POINT}/proc"
umount "${MOUNT_POINT}/dev"
umount "$MOUNT_POINT"

log_ok "Password reset complete for $USERNAME"
