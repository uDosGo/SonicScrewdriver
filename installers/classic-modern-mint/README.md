# Classic Modern Mint Installer

Sonic Family bootable USB installer with Classic Modern style modifications.

## Package Structure

```
installers/classic-modern-mint/
├── iso/                  # ISO image output
├── config/               # Configuration files
│   ├── ventoy.json       # Ventoy configuration
│   └── sonic-config.json # Sonic installer config
├── assets/               # Style assets
│   ├── wallpapers/       # Classic Modern wallpapers
│   ├── icons/           # Custom icons
│   └── themes/           # GTK/Cinnamon themes
├── scripts/              # Installation scripts
│   ├── pre-install.sh    # Pre-installation hooks
│   ├── post-install.sh   # Post-installation setup
│   └── ventoy-hook.sh   # Ventoy-specific hooks
├── ventoy/               # Ventoy integration
│   ├── plugins/         # Ventoy plugins
│   └── themes/          # Ventoy boot themes
└── README.md             # This file
```

## Build Process

1. **Package Creation**: `sonic package --installer classic-modern-mint`
2. **USB Preparation**: `sonic ventoy prepare /dev/sdX`
3. **Bundle Copy**: `sonic ventoy copy classic-modern-mint.she`
4. **Validation**: `sonic ventoy verify`

## Classic Modern Style Modifications

- **Cinnamon Theme**: Classic Modern Mint theme
- **Wallpapers**: Curated Classic Modern collection
- **Icons**: Custom icon set
- **Boot Theme**: Classic Modern Ventoy theme
- **Fonts**: Optimized for readability

## Ventoy Integration

The installer uses Ventoy for multi-boot USB support with:
- Persistent storage support
- UEFI/Legacy boot compatibility
- Custom boot menu branding
- Automatic driver injection

## Usage

```bash
# Build the installer package
./modules/ventoy/build.sh --installer classic-modern-mint

# Create bootable USB
sonic ventoy create --usb /dev/sdX --installer classic-modern-mint.she

# Test in QEMU
sonic ventoy test --qemu
```

## Configuration

Edit `config/ventoy.json` for Ventoy-specific settings and `config/sonic-config.json` for installer behavior.

## Requirements

- Ventoy 1.0.93+
- Sonic Screwdriver vA1.1.0+
- USB drive (8GB+ recommended)
- Linux/macOS host for building

## Troubleshooting

See `dev/process/checklists/ventoy-validation.md` for common issues.