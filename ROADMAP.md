# SonicScrewdriver v1.1.0 — Roadmap

> **Last Updated:** 2026-05-26
> **Status:** v1.1.0 — Scope trimmed, docs archived, code consolidated

## ✅ Current State (v1.1.0)

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

### Existing Packages (Not Wired to CLI)

| Package | Status | Notes |
|---------|--------|-------|
| `pkg/remote/` | ✅ | VNC/SSH/Samba setup functions — not exposed in CLI |
| `pkg/classicmodern/` | ✅ | Classic Modern Mint readiness checker — not exposed in CLI |
| `pkg/state/` | ✅ | SQLite state DB — used internally by other packages |
| `pkg/disk/` | ✅ | Block device/partition management — used by USB installer |
| `pkg/iso/` | ✅ | ISO downloader/verifier — used by USB installer |

### Testing

| Area | Status | Notes |
|------|--------|-------|
| Unit tests | ⚠️ Minimal | Only `pkg/container/health_simple_test.go` has real tests |
| Integration tests | ❌ | All `test/integration/` tests are `t.Skip()` stubs |

---

## 🎯 Next Steps

### Short-term (v1.2.0)

1. **Wire `remote` and `mint` commands into CLI**
   - `pkg/remote/vnc.go` and `pkg/classicmodern/readiness.go` exist but aren't exposed in `cmd/sonic/main.go`
   - Add `sonic remote <vnc|ssh|samba>` and `sonic mint <check|install|apply|status|info|doctor>` subcommands

2. **Add real unit tests**
   - `pkg/vault/vault.go` — test get/set/rotate/history
   - `pkg/library/` — test index loading, manifest validation
   - `pkg/ventoy/packager.go` — test bundle create/validate
   - `pkg/iso/downloader.go` — test SHA256 verification
   - `pkg/disk/disk.go` — test device listing parsing

3. **Fix README to match actual CLI**
   - Remove `remote` and `mint` from README until they're wired into the CLI
   - Or wire them in as part of this release

### Medium-term (v1.3.0)

4. **Build verification / CI**
   - Ensure `go build ./cmd/sonic` passes in CI
   - Add `go vet` and `go test ./...` to `.github/workflows/ci.yaml`
   - Add linting (golangci-lint)

5. **Improve error handling**
   - Many packages use `log.Fatalf` on errors — should return errors instead
   - Add structured error types for CLI commands

6. **Documentation**
   - Add inline `sonic <command> --help` for each subcommand
   - Create a proper `docs/COMMANDS.md` that matches actual CLI

### Longer-term (v2.0.0)

7. **Integration tests**
   - Implement the integration test stubs in `test/integration/`
   - Add Docker-based test fixtures for container tests
   - Add test data for library/catalogue tests

8. **Release pipeline**
   - Tag-based release workflow (`.github/workflows/release.yaml` exists)
   - Cross-compile for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
   - Publish binaries to GitHub Releases

---

## 📊 Known Gaps

- `pkg/remote/` and `pkg/classicmodern/` are orphaned — code exists but no CLI entry
- `pkg/state/db.go` uses SQLite but has no tests
- `pkg/gui/gui.go` serves an embedded web UI but `pkg/gui/static/` only has a `.gitkeep`
- `pkg/library/` has schema validation but no actual game library data
- `pkg/vault/vault.go` depends on `uServer/pkg/secrets` via `replace` directive in `go.mod`
