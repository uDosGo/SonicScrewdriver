# Ventoy Integration Brief

## Overview

Ventoy submodule integration for Sonic Family bootable USB workflows.

## Current State

- **Submodule**: `Ventoy/` (https://github.com/fredporter/Ventoy)
- **Module**: `modules/ventoy/` (build scripts and integration glue)
- **Docs**: `docs/promotion.md` (promotion flow - needs expansion)

## Integration Points

### 1. Build Integration
- **Location**: `modules/ventoy/build.sh`
- **Status**: Scaffold complete, patch integration pending
- **Next Steps**: Define patch application workflow

### 2. Promotion Flow
- **Location**: `docs/promotion.md`
- **Status**: Placeholder exists, needs detailed workflow
- **Next Steps**: Document patch intake → build → validation → promotion

### 3. CLI Integration
- **Location**: `cmd/sonic/main.go` (future)
- **Status**: Not yet implemented
- **Next Steps**: Add `sonic ventoy` subcommands

## Open Questions

1. **Patch Management**: How are Ventoy patches tracked and applied?
2. **Build Artifacts**: Where are built Ventoy images stored?
3. **Validation**: What validation steps are required before promotion?
4. **CLI Surface**: What Ventoy-specific commands are needed in `sonic` CLI?

## Action Items

- [ ] Define patch application workflow in `modules/ventoy/`
- [ ] Expand `docs/promotion.md` with concrete promotion steps
- [ ] Create Ventoy build validation checklist in `dev/process/checklists/`
- [ ] Design CLI integration for Ventoy operations
- [ ] Document environment variables for Ventoy path discovery

## References

- Ventoy upstream: https://github.com/ventoy/Ventoy
- Family alignment: `Ventoy/DOC/uDos-family-alignment.md`
- Ventoy roadmap: `Ventoy/DOC/ROADMAP.md`