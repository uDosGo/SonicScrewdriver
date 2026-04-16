# Sonic Roadmap

## Current Baseline

- Release marker: `vA1.0.0-scaffold`
- Status: scaffold complete, runtime and persistence still stubbed

## Next Steps

- implement container runtime boundary in `sonic-screwdriver/internal/container/`
- implement library index parsing and manager actions in `sonic-screwdriver/internal/library/`
- implement sqlite state layer in `sonic-screwdriver/internal/state/`
- wire CLI commands in `sonic-screwdriver/cmd/sonic/main.go` to concrete services
- complete curated manifests and validation checks in `sonic-screwdriver/library/`
- define Ventoy patch + build promotion flow in `ventoy/` and `docs/promotion.md`

## Process Lanes

Track active planning in `dev/requests/active-index.md`, keep concise technical notes in `dev/notes/`, and summarize finished work in `dev/submissions/completed-summary.md`.

## Exit Criteria for `vA1.1.0`

- `sonic library list` returns curated index entries from YAML
- `sonic install/start/stop/remove <game>` invoke runtime and state layers (stubbed Docker calls acceptable)
- state database initializes with first migration and records install/runtime status
- `make build` and `make test` pass on clean checkout
