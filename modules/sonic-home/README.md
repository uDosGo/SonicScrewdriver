# sonic-home

Lite packager and installer helper module for Sonic-family workflows.

This module is standards-aligned with uHomeNest's `sonic-home` effort, but remains
project-local to `sonic-screwdriver` (separate codebase, release cycle, and runtime boundaries).

## Commands (Bootstrap)

- `go run ./cmd/sonic-home version`
- `go run ./cmd/sonic-home pack --dry-run --source . --output manifest.dryrun.json`
- `go run ./cmd/sonic-home pack --dry-run=false --output uhome-nest-draft.she`
- `go run ./cmd/sonic-home verify ./uhome-nest-draft.she`
- `go run ./cmd/install --usb`
- `go run ./cmd/serve --port 8080`

## Scope in this bootstrap

- compiling CLI stubs for `pack`, `install`, and `serve`
- manifest dry-run generation path
- draft `.she` tar layout writer and structural verifier
- starter bundle format documentation
