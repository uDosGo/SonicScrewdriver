# Sonic Roadmap

## Current Baseline

- Release marker: `vA1.0.0-scaffold`
- Status: scaffold complete, runtime and persistence still stubbed

## Delivery Cadence

- Keep active work visible in `dev/requests/active-index.md`
- Capture short implementation notes in `dev/notes/`
- Move shipped work into `dev/submissions/completed-summary.md`
- Update this roadmap when milestone scope changes (not after every commit)

## Milestone: `vA1.1.0` (Runtime + State Foundation)

- implement container runtime boundary in `internal/container/`
- implement library index parsing and manager actions in `internal/library/`
- implement sqlite state layer in `internal/state/`
- wire CLI commands in `cmd/sonic/main.go` to concrete services
- complete curated manifests and validation checks in `library/`
- define Ventoy patch + build promotion flow in `modules/ventoy/` and `docs/promotion.md`

## Exit Criteria for `vA1.1.0`

- `sonic library list` returns curated index entries from YAML
- `sonic install/start/stop/remove <game>` invoke runtime and state layers (stubbed Docker calls acceptable)
- state database initializes with first migration and records install/runtime status
- `make build` and `make test` pass on clean checkout

## Milestone: `vA1.2.0` (Container Integration + Reliability)

### Scope

- replace runtime stubs with real Docker/OCI operations and deterministic error handling
- add health/status lifecycle checks (`created`, `running`, `stopped`, `failed`) in state + CLI output
- add basic log surfaces (`sonic logs <game>`) and runtime diagnostics for failed starts
- enforce manifest validation gates (required fields, image format, launch command constraints)
- publish first reproducible local integration script in `scripts/` for install -> start -> stop -> remove

### Exit Criteria

- `sonic start <game>` starts a real container and transitions persisted runtime state
- `sonic status <game>` returns a stable summary from state + runtime reconciliation
- failed launches produce actionable messages and non-zero exit codes
- integration scenario passes in CI on a Docker-enabled runner

## Milestone: `vA1.3.0` (Distribution + UX Hardening)

### Scope

- finalize `modules/ventoy/` patch workflow and versioned build artifacts
- document operator flow for building, validating, and promoting bootable media
- add first recovery/repair commands for inconsistent local state
- harden CLI UX (consistent output shape, machine-readable mode, and clearer error taxonomy)
- define compatibility matrix for host OS/runtime assumptions in docs

### Exit Criteria

- `docs/promotion.md` describes a reproducible promotion path from source to media artifact
- release notes template includes runtime changes, manifest changes, and upgrade notes
- boot media smoke test checklist is documented and executable by a new contributor

## Milestone: `vA1.4.0` (Polish + Extensibility)

### Scope

- add plugin-style game/provider abstraction points without breaking curated defaults
- add migration tooling for state schema upgrades between minor versions
- improve test depth (state migration tests, runtime integration tests, manifest fuzz/edge cases)
- define stable extension contract for future GUI/TUI feature parity

### Exit Criteria

- extension points are documented and covered with contract tests
- migration path from `vA1.1.x` databases is tested and automated
- release gate includes unit, integration, and migration suites with clear pass/fail criteria

## Cross-Cutting Risks

- Docker/runtime differences across environments can cause nondeterministic behavior
- manifest quality drift can break install/start flows unless validation remains strict
- state/runtime divergence must be reconciled explicitly to avoid stale UX
- Ventoy patch maintenance can become a release bottleneck without repeatable scripts

## Operating Rule

Prefer small, shippable milestones. If a scope item cannot be validated by a concrete command or test, split it before starting implementation.
