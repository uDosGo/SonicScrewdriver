#!/bin/bash
# =============================================================================
# SonicScrewdriver — Chroot Customization Script
# =============================================================================
# Runs inside the Linux Mint squashfs chroot to apply customizations.
# Called by build-iso.sh during the ISO build process.
# =============================================================================

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------
HOSTNAME="sonic-mint"
USERNAME="sonic"
USER_PASSWORD="sonic"
TIMEZONE="Australia/Brisbane"
LOCALE="en_AU.UTF-8"
KEYBOARD_LAYOUT="us"

# ---------------------------------------------------------------------------
# Mount special filesystems
# ---------------------------------------------------------------------------
mount -t proc none /proc
mount -t sysfs none /sys
mount -t devtmpfs none /dev

# ---------------------------------------------------------------------------
# Set hostname
# ---------------------------------------------------------------------------
echo "${HOSTNAME}" > /etc/hostname
sed -i "s/127.0.1.1.*/127.0.1.1\t${HOSTNAME}/" /etc/hosts

# ---------------------------------------------------------------------------
# Set locale
# ---------------------------------------------------------------------------
echo "${LOCALE} UTF-8" > /etc/locale.gen
locale-gen
update-locale LANG="${LOCALE}" LANGUAGE="${LOCALE%%.*}"

# ---------------------------------------------------------------------------
# Set timezone
# ---------------------------------------------------------------------------
echo "${TIMEZONE}" > /etc/timezone
ln -sf "/usr/share/zoneinfo/${TIMEZONE}" /etc/localtime
dpkg-reconfigure -f noninteractive tzdata

# ---------------------------------------------------------------------------
# Set keyboard layout
# ---------------------------------------------------------------------------
cat > /etc/default/keyboard << EOF
XKBMODEL="pc105"
XKBLAYOUT="${KEYBOARD_LAYOUT}"
XKBVARIANT=""
XKBOPTIONS=""
BACKSPACE="guess"
EOF

# ---------------------------------------------------------------------------
# Create user
# ---------------------------------------------------------------------------
useradd -m -G sudo,adm,cdrom,plugdev -s /bin/bash "${USERNAME}"
echo "${USERNAME}:${USER_PASSWORD}" | chpasswd

# Configure sudo (no password for sudo group)
sed -i 's/%sudo.*/%sudo ALL=(ALL:ALL) NOPASSWD:ALL/' /etc/sudoers

# ---------------------------------------------------------------------------
# Install Sonic tools
# ---------------------------------------------------------------------------
if [[ -f /tmp/install-sonic.sh ]]; then
    bash /tmp/install-sonic.sh
fi

# ---------------------------------------------------------------------------
# Configure desktop
# ---------------------------------------------------------------------------
# Set default background
SUDO_USER="${USERNAME}" dbus-launch --exit-with-session gsettings set \
    org.cinnamon.desktop.background picture-uri \
    "file:///usr/share/backgrounds/sonic/sonic-wallpaper.png" 2>/dev/null || true

# Add Sonic launcher to desktop
cat > "/home/${USERNAME}/Desktop/sonic.desktop" << EOF
[Desktop Entry]
Name=SonicScrewdriver
Comment=Universal USB Bootloader & System Toolkit
Exec=sonic
Terminal=true
Type=Application
Icon=utilities-terminal
Categories=System;
EOF

chown "${USERNAME}:${USERNAME}" "/home/${USERNAME}/Desktop/sonic.desktop"
chmod +x "/home/${USERNAME}/Desktop/sonic.desktop"

# ---------------------------------------------------------------------------
# Cleanup
# ---------------------------------------------------------------------------
apt-get clean
rm -rf /var/lib/apt/lists/*
rm -rf /tmp/*

# Unmount special filesystems
umount /proc
umount /sys
umount /dev

exit 0
