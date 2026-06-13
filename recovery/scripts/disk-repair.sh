#!/bin/bash
# =============================================================================
# SonicScrewdriver — Disk Repair & Recovery Tool
# =============================================================================
# Scans and repairs common disk issues. Part of the Sonic recovery toolkit.
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
echo "  SonicScrewdriver — Disk Repair Tool"
echo "=============================================="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    log_error "This script must be run as root (sudo)"
    exit 1
fi

# List available disks
echo "Available disks:"
lsblk -o NAME,SIZE,TYPE,FSTYPE,LABEL,MOUNTPOINT,MODEL 2>/dev/null || \
    diskutil list 2>/dev/null || \
    echo "No disk listing tool found"

echo ""
read -rp "Enter device to repair (e.g., /dev/sda): " DEVICE

if [[ ! -b "$DEVICE" ]]; then
    log_error "Not a block device: $DEVICE"
    exit 1
fi

echo ""
echo "Select repair type:"
echo "  1) Check filesystem (read-only)"
echo "  2) Repair filesystem"
echo "  3) Check for bad blocks"
echo "  4) Wipe filesystem signatures"
echo "  5) Full diagnostic"
read -rp "Choice [1-5]: " CHOICE

case "$CHOICE" in
    1)
        log_info "Checking filesystem on $DEVICE (read-only)..."
        fsck -n "$DEVICE" || true
        ;;
    2)
        log_warn "Repairing filesystem on $DEVICE..."
        fsck -y "$DEVICE" || true
        ;;
    3)
        log_info "Checking for bad blocks on $DEVICE..."
        badblocks -sv "$DEVICE" || true
        ;;
    4)
        log_warn "Wiping filesystem signatures on $DEVICE..."
        wipefs -a "$DEVICE" || true
        ;;
    5)
        log_info "Running full diagnostic on $DEVICE..."
        echo ""
        echo "--- SMART Status ---"
        smartctl -a "$DEVICE" 2>/dev/null || echo "SMART not available"
        echo ""
        echo "--- Filesystem Check ---"
        fsck -n "$DEVICE" 2>/dev/null || true
        echo ""
        echo "--- Bad Blocks ---"
        badblocks -sv "$DEVICE" 2>/dev/null || echo "badblocks not available"
        echo ""
        echo "--- Partition Table ---"
        fdisk -l "$DEVICE" 2>/dev/null || parted -l "$DEVICE" 2>/dev/null || true
        ;;
    *)
        log_error "Invalid choice"
        exit 1
        ;;
esac

log_ok "Operation complete"
