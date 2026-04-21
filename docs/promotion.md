# Promotion Flow

Promotion defines the validation and release gates for moving components from development to production-ready artifacts across the Sonic Family.

## Ventoy Promotion Workflow

### Overview

The Ventoy promotion workflow covers patch intake, local build, validation, and promotion to bootable media artifacts. This workflow enables creating bootable USB installers with Sonic Family tools pre-installed.

### Workflow Steps

#### 1. Patch Intake
- **Input**: Ventoy patches in `modules/ventoy/patches/`
- **Process**: Apply patches to Ventoy submodule using `git apply`
- **Validation**: Verify patch application with `modules/ventoy/build.sh --verify`
- **Output**: Patched Ventoy source in `Ventoy/`
- **Example**: `git -C Ventoy apply ../modules/ventoy/patches/*.patch`

#### 2. Local Build
- **Input**: Patched Ventoy source
- **Process**: Run `modules/ventoy/build.sh` to create bootable ISO
- **Validation**: Check build artifacts exist and are valid
- **Output**: Ventoy build artifacts in `build/ventoy/`
- **Artifacts**: `ventoy.img`, `Ventoy2Disk.sh`, and other tools

#### 3. Validation
- **Input**: Build artifacts
- **Process**: Run validation checklist (see `dev/process/checklists/ventoy-validation.md`)
- **Validation**: All checklist items pass
- **Output**: Validated build artifacts
- **Commands**: `./modules/ventoy/validate.sh`

#### 4. Promotion
- **Input**: Validated build artifacts
- **Process**: Promote to release directory with version tagging
- **Validation**: Verify promotion with `sonic ventoy verify`
- **Output**: Promoted artifacts in `release/ventoy/vX.Y.Z/`
- **Example**: `sonic ventoy promote --version 1.0.0`

#### 5. USB Creation (Optional)
- **Input**: Promoted Ventoy artifacts
- **Process**: Create bootable USB using `sonic ventoy usb`
- **Validation**: Verify USB bootability
- **Output**: Bootable USB drive with Sonic Family tools
- **Example**: `sonic ventoy usb --device /dev/sdX --bundle sonic-family-v1.0.0.she`

### Checklist

See `dev/process/checklists/ventoy-promotion.md` for the detailed promotion checklist.

### Environment Variables

- `SONIC_VENTOY_ROOT`: Path to Ventoy submodule (default: `$(pwd)/Ventoy`)
- `VENTOY_BUILD_DIR`: Build output directory (default: `$(pwd)/build/ventoy`)
- `VENTOY_RELEASE_DIR`: Release promotion directory (default: `$(pwd)/release/ventoy`)
- `VENTOY_VERSION`: Current version being promoted (e.g., `1.0.0`)

### CLI Commands

```bash
# Build Ventoy with patches
./modules/ventoy/build.sh

# Verify build artifacts
./modules/ventoy/build.sh --verify

# Run promotion checklist
sonic ventoy validate

# Promote to release with version
sonic ventoy promote --version 1.0.0

# Create bootable USB
sonic ventoy usb --device /dev/sdX --bundle sonic-family-v1.0.0.she

# Show bundle information
sonic ventoy info sonic-family-v1.0.0.she
```

### Bundle Format

Sonic Family Ventoy bundles use the `.she` (Sonic Home Edition) format:

```
sonic-family-v1.0.0.she
├── ventoy/                    # Ventoy core files
│   ├── ventoy.img             # Bootable image
│   ├── Ventoy2Disk.sh         # USB creation script
│   └── ventoy.json           # Configuration
├── sonic/                    # Sonic Family tools
│   ├── bin/                  # Binaries
│   │   ├── sonic             # Sonic Screwdriver CLI
│   │   └── udos              # uDos Connect CLI
│   └── config/               # Configuration files
├── installers/               # Installer packages
│   └── classic-modern-mint/  # Classic Modern Mint installer
└── README.md                 # Bundle documentation
```

### USB Layout

When creating bootable USB with `sonic ventoy usb`:

```
USB Drive (label: SONIC_FAMILY)
├── ventoy/                    # Ventoy boot files
│   ├── ventoy.img             # Boot image
│   ├── Ventoy2Disk.sh         # Installation script
│   └── ventoy.json           # Boot configuration
├── sonic/                    # Sonic Family tools
│   ├── bin/                  # CLI tools
│   │   ├── sonic             # Main CLI
│   │   └── udos              # uDos CLI
│   └── config/               # Configurations
├── installers/               # OS installers
│   └── classic-modern-mint/  # Linux Mint with Classic Modern theme
├── bundles/                  # Additional bundles
│   └── sonic-family-v1.0.0.she # Main bundle
└── README.txt                 # User instructions
```

### Release Process

1. **Prepare Release**:
   ```bash
   # Update version in source
   sed -i 's/vA1.0.0/vA1.1.0/g' version
   
   # Build all components
   make build
   
   # Create Ventoy bundle
   sonic ventoy package --installer classic-modern-mint --output sonic-family-v1.1.0.she
   ```

2. **Validate Release**:
   ```bash
   # Verify bundle
   sonic ventoy verify sonic-family-v1.1.0.she
   
   # Test health monitoring
   sonic health --all
   
   # Test repair functionality
   sonic repair --all
   ```

3. **Promote Release**:
   ```bash
   # Promote to release directory
   sonic ventoy promote --version 1.1.0 --bundle sonic-family-v1.1.0.she
   
   # Create checksums
   sha256sum release/ventoy/v1.1.0/* > release/ventoy/v1.1.0/CHECKSUMS.txt
   
   # Sign release (optional)
   gpg --detach-sign --armor release/ventoy/v1.1.0/CHECKSUMS.txt
   ```

4. **Publish Release**:
   ```bash
   # Update CHANGELOG
   # Tag release in git
   git tag vA1.1.0
   git push origin vA1.1.0
   
   # Create GitHub release
   gh release create vA1.1.0 --notes-file CHANGELOG.md
   
   # Upload artifacts
   gh release upload vA1.1.0 release/ventoy/v1.1.0/*
   ```

### Rollback Procedure

If issues are found after release:

1. **Identify Issue**: Check logs and health status
   ```bash
   sonic health --all
   journalctl -u sonic-daemon
   ```

2. **Attempt Repair**:
   ```bash
   sonic repair --all
   ```

3. **Rollback to Previous Version**:
   ```bash
   # Download previous release
   wget https://github.com/sonic-family/sonic-screwdriver/releases/download/vA1.0.0/sonic-family-v1.0.0.she
   
   # Reinstall
   sonic ventoy promote --version 1.0.0 --bundle sonic-family-v1.0.0.she --force
   ```

4. **Notify Users**:
   - Create GitHub issue for tracking
   - Update documentation with workaround
   - Announce rollback in release notes

### Success Criteria

- ✅ Ventoy bundle creates successfully
- ✅ All health checks pass
- ✅ USB creation works on multiple devices
- ✅ Boot process completes without errors
- ✅ Sonic Family tools are accessible from USB
- ✅ Rollback procedure is tested and documented

### Troubleshooting

**Issue: USB not booting**
- Verify SHA256 checksums of download
- Try different USB port
- Check BIOS boot order
- Test on different machine

**Issue: Health checks failing**
```bash
# Check Docker daemon
systemctl status docker

# Restart Sonic daemon
systemctl restart sonic-daemon

# Check logs
journalctl -u sonic-daemon -f
```

**Issue: Promotion validation fails**
- Verify all required files are present
- Check file permissions
- Validate checksums
- Review validation checklist

### Best Practices

1. **Test Before Release**: Always test USB creation on multiple machines
2. **Document Changes**: Update CHANGELOG.md with all changes
3. **Version Consistency**: Ensure version numbers match across all files
4. **Backup**: Keep previous release available for rollback
5. **Sign Releases**: Use GPG signing for security
6. **Automate**: Use CI/CD for consistent releases

## Sonic Home/Express Promotion

### Overview

Sonic Home and Sonic Express modules follow similar promotion workflows tailored to their specific artifact types.

### Workflow Steps

#### 1. Module Build
- **Input**: Module source in `modules/sonic-home/` or `modules/sonic-express/`
- **Process**: Run module-specific build script
- **Validation**: Check build artifacts
- **Output**: Module build artifacts

#### 2. Validation
- **Input**: Build artifacts
- **Process**: Run module validation checklist
- **Validation**: All checklist items pass
- **Output**: Validated build artifacts

#### 3. Promotion
- **Input**: Validated build artifacts
- **Process**: Promote to release directory
- **Validation**: Verify promotion
- **Output**: Promoted artifacts in `release/sonic-home/` or `release/sonic-express/`

### Checklists

- Sonic Home: `dev/process/checklists/sonic-home-promotion.md`
- Sonic Express: `dev/process/checklists/sonic-express-promotion.md`

## Cross-Component Promotion

For full Sonic Family releases, coordinate promotion across all components:

1. Promote Ventoy artifacts
2. Promote Sonic Home artifacts  
3. Promote Sonic Express artifacts
4. Create unified release bundle
5. Run cross-component validation

## Rules

1. **No partial promotions**: All validation steps must pass before promotion
2. **Artifacts are immutable**: Once promoted, artifacts should not be modified
3. **Document everything**: All promotion steps and validations must be documented
4. **Automate where possible**: Prefer scripted validation over manual checks

## Open Questions

- How do we handle version synchronization across components?
- What's the rollback procedure for failed promotions?
- How do we validate cross-component compatibility?
- What signing/verification mechanisms are required?
