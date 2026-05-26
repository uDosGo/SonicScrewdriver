# SonicScrewdriver v1.2.0 — Roadmap

> **Last Updated:** 2026-05-26
> **Status:** v1.2.0 — `remote` and `mint` commands wired, unit tests added, README synced

## ✅ Current State (v1.2.0)

SonicScrewdriver is a **Unified System Toolkit** — a Go CLI for system administration tasks.

### Implemented Commands

| Command | Status | Notes |
|---------|--------|-------|
| `sonic container` | ✅ | Docker runtime: list, start, stop, restart, remove, health |
| `sonic usb` | ✅ | USB installer: list, prepare, install, full-install |
| `sonic vault` | ✅ | Secret store (via uServer/pkg/secrets): get, set, list, rotate, history |
| `sonic gui` | ✅ | Web GUI dashboard (embedded) |
| `sonic catalogue` | ✅ | Device catalogue: list, find (scans Vaults + uCode repos) |
| `sonic knowledge` | ✅ | Knowledge sources: sources, query |
| `sonic library` | ✅ | Game library index: list, info, validate |
| `sonic ventoy` | ✅ | .she bundle: create, validate, info |
| `sonic remote` | ✅ | Remote access: vnc, ssh, samba, info |
| `sonic mint` | ✅ | Classic Modern Mint readiness: check, doctor, info, apply |

### Testing

| Area | Status | Notes |
|------|--------|-------|
| Unit tests | ✅ | `pkg/container/`, `pkg/disk/`, `pkg/iso/`, `pkg/library/`, `pkg/vault/`, `pkg/ventoy/` |
| Integration tests | ❌ | All `test/integration/` tests are `t.Skip()` stubs |

---

## 🎯 Next Steps

### Short-term (v1.3.0)

1. **Build verification / CI**
   - Ensure `go build ./cmd/sonic` passes in CI
   - Add `go vet` and `go test ./...` to `.github/workflows/ci.yaml`
   - Add linting (golangci-lint)

2. **Improve error handling**
   - Many packages use `log.Fatalf` on errors — should return errors instead
   - Add structured error types for CLI commands

3. **Documentation**
   - Add inline `sonic <command> --help` for each subcommand
   - Create a proper `docs/COMMANDS.md` that matches actual CLI

### Medium-term (v1.4.0)

4. **Integration tests**
   - Implement the integration test stubs in `test/integration/`
   - Add Docker-based test fixtures for container tests
   - Add test data for library/catalogue tests

5. **Release pipeline**
   - Tag-based release workflow (`.github/workflows/release.yaml` exists)
   - Cross-compile for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
   - Publish binaries to GitHub Releases

### Longer-term (v2.0.0)

6. **GUI improvements**
   - `pkg/gui/gui.go` serves an embedded web UI but `pkg/gui/static/` only has a `.gitkeep`
   - Build out the dashboard with real container/vault/catalogue views

7. **State management**
   - `pkg/state/db.go` uses SQLite but has no tests
   - Wire state DB into CLI for persistent install tracking

---

## 📊 Known Gaps

- `pkg/gui/gui.go` serves an embedded web UI but `pkg/gui/static/` only has a `.gitkeep`
- `pkg/state/db.go` uses SQLite but has no tests
- `pkg/vault/vault.go` depends on `uServer/pkg/secrets` via `replace` directive in `go.mod`
- Integration tests are all stubs
- No CI pipeline configured
