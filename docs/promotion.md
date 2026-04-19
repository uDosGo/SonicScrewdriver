# Promotion Flow

Promotion defines the validation and release gates for moving components from development to production-ready artifacts across the Sonic Family.

## Ventoy Promotion Workflow

### Overview

The Ventoy promotion workflow covers patch intake, local build, validation, and promotion to bootable media artifacts.

### Workflow Steps

#### 1. Patch Intake
- **Input**: Ventoy patches in `modules/ventoy/patches/`
- **Process**: Apply patches to Ventoy submodule
- **Validation**: Verify patch application with `modules/ventoy/build.sh --verify`
- **Output**: Patched Ventoy source in `Ventoy/`

#### 2. Local Build
- **Input**: Patched Ventoy source
- **Process**: Run `modules/ventoy/build.sh`
- **Validation**: Check build artifacts exist and are valid
- **Output**: Ventoy build artifacts in `build/ventoy/`

#### 3. Validation
- **Input**: Build artifacts
- **Process**: Run validation checklist (see `dev/process/checklists/ventoy-validation.md`)
- **Validation**: All checklist items pass
- **Output**: Validated build artifacts

#### 4. Promotion
- **Input**: Validated build artifacts
- **Process**: Promote to release directory
- **Validation**: Verify promotion with `sonic ventoy verify`
- **Output**: Promoted artifacts in `release/ventoy/`

### Checklist

See `dev/process/checklists/ventoy-promotion.md` for the detailed promotion checklist.

### Environment Variables

- `SONIC_VENTOY_ROOT`: Path to Ventoy submodule (default: `$(pwd)/Ventoy`)
- `VENTOY_BUILD_DIR`: Build output directory (default: `$(pwd)/build/ventoy`)
- `VENTOY_RELEASE_DIR`: Release promotion directory (default: `$(pwd)/release/ventoy`)

### CLI Commands

```bash
# Build Ventoy with patches
./modules/ventoy/build.sh

# Verify build artifacts
./modules/ventoy/build.sh --verify

# Run promotion checklist
sonic ventoy validate

# Promote to release
sonic ventoy promote
```

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
