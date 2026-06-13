#!/usr/bin/env bash
# =============================================================================
# SonicScrewdriver — Linux Mint ISO Builder
# =============================================================================
# Builds a customized Linux Mint ISO with uDos/Sonic pre-installed.
#
# Usage:
#   ./build-iso.sh [--input <mint-iso>] [--output <output-iso>] [--help]
#
# Requirements:
#   - Linux (or Docker on macOS)
#   - xorriso, unsquashfs, mksquashfs
#   - sudo access for chroot operations
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK_DIR="/tmp/sonic-mint-build"
OUTPUT_DIR="${SCRIPT_DIR}/output"

# Default configuration
INPUT_ISO="${SCRIPT_DIR}/linuxmint-22-cinnamon-64bit.iso"
OUTPUT_ISO="${OUTPUT_DIR}/sonic-mint-22-cinnamon-64bit.iso"
MINT_VERSION="22"
MINT_CODENAME="xia"
SONIC_VERSION="2.0.0"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info()  { echo -e "${BLUE}[INFO]${NC} $1"; }
log_ok()    { echo -e "${GREEN}[OK]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# ---------------------------------------------------------------------------
# Parse arguments
# ---------------------------------------------------------------------------
while [[ $# -gt 0 ]]; do
    case "$1" in
        --input) INPUT_ISO="$2"; shift 2 ;;
        --output) OUTPUT_ISO="$2"; shift 2 ;;
        --help)
            echo "Usage: $0 [--input <mint-iso>] [--output <output-iso>]"
            exit 0
            ;;
        *) log_error "Unknown option: $1"; exit 1 ;;
    esac
done

# ---------------------------------------------------------------------------
# Check prerequisites
# ---------------------------------------------------------------------------
check_prereqs() {
    local missing=0
    for cmd in xorriso unsquashfs mksquashfs mount umount chroot; do
        if ! command -v "$cmd" &>/dev/null; then
            log_error "Missing required tool: $cmd"
            missing=1
        fi
    done

    if [[ $missing -eq 1 ]]; then
        echo ""
        echo "Install missing tools:"
        echo "  sudo apt-get install xorriso squashfs-tools"
        exit 1
    fi
}

# ---------------------------------------------------------------------------
# Extract ISO
# ---------------------------------------------------------------------------
extract_iso() {
    log_info "Extracting ISO: ${INPUT_ISO}"

    mkdir -p "${WORK_DIR}"/{mnt,extract}

    # Mount the ISO
    sudo mount -o loop "${INPUT_ISO}" "${WORK_DIR}/mnt"

    # Copy contents
    sudo rsync -a "${WORK_DIR}/mnt/" "${WORK_DIR}/extract/"
    sudo umount "${WORK_DIR}/mnt"

    log_ok "ISO extracted to ${WORK_DIR}/extract"
}

# ---------------------------------------------------------------------------
# Customize squashfs
# ---------------------------------------------------------------------------
customize_squashfs() {
    log_info "Customizing squashfs filesystem..."

    local squashfs="${WORK_DIR}/extract/casper/filesystem.squashfs"
    local root="${WORK_DIR}/squashfs-root"

    # Unsquash the filesystem
    sudo unsquashfs -f -d "${root}" "${squashfs}"

    # Copy overlay files
    log_info "Applying overlay files..."
    sudo rsync -a "${SCRIPT_DIR}/overlay/" "${root}/"

    # Copy chroot scripts
    sudo cp "${SCRIPT_DIR}/chroot-customize.sh" "${root}/tmp/"
    sudo cp "${SCRIPT_DIR}/install-sonic.sh" "${root}/tmp/"

    # Chroot and customize
    log_info "Entering chroot to customize..."
    sudo mount --bind /dev "${root}/dev"
    sudo mount --bind /proc "${root}/proc"
    sudo mount --bind /sys "${root}/sys"
    sudo mount --bind /run "${root}/run"

    sudo chroot "${root}" /bin/bash /tmp/chroot-customize.sh

    # Cleanup
    sudo rm -f "${root}/tmp/chroot-customize.sh"
    sudo rm -f "${root}/tmp/install-sonic.sh"
    sudo rm -f "${root}/root/.bash_history"

    sudo umount "${root}/run"
    sudo umount "${root}/sys"
    sudo umount "${root}/proc"
    sudo umount "${root}/dev"

    # Repack squashfs
    log_info "Repacking squashfs..."
    sudo mksquashfs "${root}" "${WORK_DIR}/custom.squashfs" \
        -comp xz -b 1M -noappend

    # Replace the squashfs in the ISO
    sudo mv "${WORK_DIR}/custom.squashfs" "${WORK_DIR}/extract/casper/filesystem.squashfs"

    # Update filesystem size
    local fs_size
    fs_size=$(sudo du -sx --block-size=1 "${root}" | cut -f1)
    echo "${fs_size}" | sudo tee "${WORK_DIR}/extract/casper/filesystem.size" > /dev/null

    # Cleanup
    sudo rm -rf "${root}"

    log_ok "Squashfs customized and repacked"
}

# ---------------------------------------------------------------------------
# Update ISO metadata
# ---------------------------------------------------------------------------
update_metadata() {
    log_info "Updating ISO metadata..."

    # Update .disk/info
    echo "Sonic Mint ${MINT_VERSION} \"${MINT_CODENAME}\" - Release $(date +%Y%m%d)" \
        | sudo tee "${WORK_DIR}/extract/.disk/info" > /dev/null

    # Update isolinux/txt.cfg
    local isolinux_cfg="${WORK_DIR}/extract/isolinux/txt.cfg"
    if [[ -f "${isolinux_cfg}" ]]; then
        sudo sed -i "s/default live/default sonic-live/" "${isolinux_cfg}"
        sudo sed -i "s/live-mint/sonic-mint/" "${isolinux_cfg}"
    fi

    # Update grub config
    local grub_cfg="${WORK_DIR}/extract/boot/grub/grub.cfg"
    if [[ -f "${grub_cfg}" ]]; then
        sudo sed -i "s/Linux Mint/Sonic Mint/g" "${grub_cfg}"
        sudo sed -i "s/linuxmint/sonic-mint/g" "${grub_cfg}"
    fi

    log_ok "ISO metadata updated"
}

# ---------------------------------------------------------------------------
# Build final ISO
# ---------------------------------------------------------------------------
build_iso() {
    log_info "Building final ISO: ${OUTPUT_ISO}"

    mkdir -p "${OUTPUT_DIR}"

    sudo xorriso -as mkisofs \
        -r -V "SonicMint" \
        -J -joliet-long \
        -cache-inodes \
        -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin \
        -b isolinux/isolinux.bin \
        -c isolinux/boot.cat \
        -boot-load-size 4 \
        -boot-info-table \
        -no-emul-boot \
        -eltorito-alt-boot \
        -e boot/grub/efi.img \
        -no-emul-boot \
        -isohybrid-gpt-basdat \
        -o "${OUTPUT_ISO}" \
        "${WORK_DIR}/extract"

    log_ok "ISO built: ${OUTPUT_ISO}"
}

# ---------------------------------------------------------------------------
# Cleanup
# ---------------------------------------------------------------------------
cleanup() {
    log_info "Cleaning up..."
    sudo rm -rf "${WORK_DIR}"
    log_ok "Cleanup complete"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
main() {
    echo ""
    echo "=============================================="
    echo "  SonicScrewdriver — Linux Mint ISO Builder"
    echo "  Version ${SONIC_VERSION}"
    echo "=============================================="
    echo ""

    if [[ ! -f "${INPUT_ISO}" ]]; then
        log_error "Input ISO not found: ${INPUT_ISO}"
        echo "Download Linux Mint 22 from: https://linuxmint.com/download.php"
        echo "Then run: $0 --input /path/to/linuxmint-22-cinnamon-64bit.iso"
        exit 1
    fi

    check_prereqs
    extract_iso
    customize_squashfs
    update_metadata
    build_iso
    cleanup

    echo ""
    log_ok "Sonic Mint ISO ready: ${OUTPUT_ISO}"
    echo ""
    echo "  Size: $(ls -lh "${OUTPUT_ISO}" | awk '{print $5}')"
    echo "  SHA256: $(sha256sum "${OUTPUT_ISO}" | cut -d' ' -f1)"
    echo ""
}

main "$@"
