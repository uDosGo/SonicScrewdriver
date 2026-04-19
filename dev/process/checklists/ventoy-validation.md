# Ventoy Validation Checklist

## Pre-Build Validation

- [ ] Ventoy submodule is initialized and up to date
  ```bash
  git submodule update --init --recursive
  ```

- [ ] Required patches exist in `modules/ventoy/patches/`
  ```bash
  ls -la modules/ventoy/patches/
  ```

- [ ] Build environment variables are set correctly
  ```bash
  echo "SONIC_VENTOY_ROOT=$SONIC_VENTOY_ROOT"
  echo "VENTOY_BUILD_DIR=$VENTOY_BUILD_DIR"
  ```

## Build Validation

- [ ] Build script executes without errors
  ```bash
  ./modules/ventoy/build.sh
  ```

- [ ] Build artifacts are created in expected location
  ```bash
  ls -la $VENTOY_BUILD_DIR
  ```

- [ ] Build artifacts have correct permissions
  ```bash
  stat $VENTOY_BUILD_DIR/*
  ```

- [ ] Build artifacts are not empty
  ```bash
  find $VENTOY_BUILD_DIR -type f -size 0
  ```

## Post-Build Validation

- [ ] Ventoy version information is correct
  ```bash
  strings $VENTOY_BUILD_DIR/ventoy.img | grep VERSION
  ```

- [ ] Required binaries are present
  ```bash
  file $VENTOY_BUILD_DIR/Ventoy2Disk.sh
  file $VENTOY_BUILD_DIR/ventoy.img
  ```

- [ ] Checksums match expected values
  ```bash
  sha256sum $VENTOY_BUILD_DIR/*
  ```

## Functional Validation

- [ ] Ventoy image can be mounted
  ```bash
  sudo mount -o loop $VENTOY_BUILD_DIR/ventoy.img /mnt/test
  sudo umount /mnt/test
  ```

- [ ] Ventoy tools are executable
  ```bash
  chmod +x $VENTOY_BUILD_DIR/Ventoy2Disk.sh
  $VENTOY_BUILD_DIR/Ventoy2Disk.sh -v
  ```

- [ ] Patch verification passes
  ```bash
  ./modules/ventoy/build.sh --verify
  ```

## Promotion Validation

- [ ] Release directory exists
  ```bash
  mkdir -p $VENTOY_RELEASE_DIR
  ```

- [ ] Artifacts can be copied to release directory
  ```bash
  cp -v $VENTOY_BUILD_DIR/* $VENTOY_RELEASE_DIR/
  ```

- [ ] Release artifacts have correct permissions
  ```bash
  chmod 644 $VENTOY_RELEASE_DIR/*.img
  chmod 755 $VENTOY_RELEASE_DIR/*.sh
  ```

- [ ] Release manifest can be generated
  ```bash
  ls -la $VENTOY_RELEASE_DIR > $VENTOY_RELEASE_DIR/MANIFEST.txt
  ```

## Validation Commands

```bash
# Run full validation
./modules/ventoy/validate.sh

# Check specific validation step
./modules/ventoy/validate.sh --build
./modules/ventoy/validate.sh --functional
./modules/ventoy/validate.sh --promotion
```

## Troubleshooting

### Common Issues

1. **Missing patches**: Ensure all required patches are in `modules/ventoy/patches/`
2. **Build failures**: Check build script output for specific errors
3. **Permission issues**: Verify build directory permissions
4. **Missing dependencies**: Ensure all build tools are installed

### Debugging Commands

```bash
# Verbose build
bash -x ./modules/ventoy/build.sh

# Check Ventoy submodule status
git submodule status Ventoy

# List all patches
find modules/ventoy/patches -type f

# Check environment
env | grep VENTOY
```