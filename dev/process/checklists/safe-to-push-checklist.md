# Safe-to-Push Checklist

Use this checklist before pushing `sonic-screwdriver` to remote.

## 1) Confirm repository shape

- [ ] Core code is at repo root: `cmd/`, `internal/`, `library/`, `pkg/`
- [ ] Modules are under `modules/`:
  - [ ] `modules/ventoy`
  - [ ] `modules/sonic-home`
- [ ] No nested duplicate repo directory remains (single root layout)

## 2) Confirm ignored local/runtime artifacts

- [ ] Root `.gitignore` includes runtime/local patterns (`runtime/`, `state/`, `*.db`, `*.log`, `*.pid`, `*.sock`, `dist/`, `bundles/`, `*.she`, `*.sig`)
- [ ] Root `.gitignore` includes local tooling patterns (`.cursor/`, `agent-transcripts/`, `.env*`)
- [ ] `modules/sonic-home/.gitignore` excludes module runtime + bundle outputs
- [ ] `modules/ventoy/.gitignore` excludes build/tmp outputs

Quick checks:

- `git check-ignore -v .cursor/state.json`
- `git check-ignore -v modules/ventoy/build/output.iso`
- `git check-ignore -v modules/sonic-home/runtime/state.db`

## 3) Confirm workspace/build wiring

- [ ] `go.work` includes:
  - [ ] `.`
  - [ ] `./code-vault`
  - [ ] `./modules/sonic-home`
- [ ] Root `go.mod` replace path points to `./code-vault`
- [ ] `Makefile` uses root-relative paths only (no nested-prefixed paths)

## 4) Verify code health

- [ ] `go test ./...` (root) passes
- [ ] `go test ./modules/sonic-home/...` passes
- [ ] Optional: `make build` and `make test` pass

## 5) Final git review

- [ ] `git status --short` contains only intended files
- [ ] No large local artifacts appear in staged files (`*.db`, `*.log`, `*.she`, `dist/`, `bundles/`, `runtime/`)
- [ ] Docs reference current layout (`modules/ventoy`, `modules/sonic-home`, root `cmd/internal/library/pkg`)

## 6) Push readiness

- [ ] Branch name is correct
- [ ] Commit message reflects structural migration + hygiene updates
- [ ] Ready to push
