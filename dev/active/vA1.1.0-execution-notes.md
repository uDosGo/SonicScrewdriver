# `vA1.1.0` Execution Notes

Purpose: daily working log for the active 2-week slice. Keep entries brief and evidence-driven.

## Slice Goals (Reference)

- P0: runtime boundary + CLI wiring
- P0: library manager + manifest validation
- P0: sqlite state bootstrap + persistence
- P0: build/test gate (`make build`, `make test`)
- P1: Ventoy patch + promotion flow draft

## Briefs and Inputs

- uHomeNest integration brief: `dev/notes/UHOMENEST-V1.1.0-INTEGRATION-BRIEF.md`
- sonic-home brief (`v1.1.1`): `dev/notes/UHOMENEST-V1.1.1-SONIC-HOME-BRIEF.md`

## Status Legend

- `[ ]` not started
- `[-]` in progress
- `[x]` done
- `[!]` blocked

## Owners

- `owner/core-runtime`
- `owner/library`
- `owner/state`
- `owner/release`
- `owner/packager`

## Day 1

Date: 2026-04-16

### Day-1 execution focus

- lock interface and schema contracts first
- ship minimum skeletons for both tracks:
  - `v1.1.0`: integration-manager baseline contracts
  - `v1.1.1`: sonic-home module scaffold + CLI entrypoints
- integrate Ventoy submodule into new dev flow

### Plan for today

- [-] runtime (`owner/core-runtime`): finalize interface boundaries
  - target files:
    - `internal/container/` (runtime interface + stub impl)
    - `cmd/sonic/main.go` (wiring points only)
  - commands:
    - `make build`
    - `make test`
- [-] library (`owner/library`): confirm index/manifests schema assumptions
  - target files:
    - `internal/library/`
    - `library/index.yaml`
  - commands:
    - `go test ./internal/library/...`
    - `sonic library list` (post-wiring smoke check)
- [-] state (`owner/state`): define initial migration and state model
  - target files:
    - `internal/state/`
  - commands:
    - `go test ./internal/state/...`
- [ ] release (`owner/release`): collect Ventoy promotion unknowns
  - target files:
    - `docs/promotion.md`
    - `ventoy/`
  - commands:
    - `rg "TODO|TBD|WIP" docs ventoy`
- [-] packager (`owner/packager`): bootstrap `v1.1.1` sonic-home track
  - target files:
    - `modules/sonic-home/cmd/pack/main.go`
    - `modules/sonic-home/cmd/install/main.go`
    - `modules/sonic-home/cmd/serve/main.go`
    - `modules/sonic-home/pkg/manifest/`
    - `modules/sonic-home/docs/BUNDLE-FORMAT.md`
  - commands:
    - `go test ./modules/sonic-home/...` (or initialize module + run)
    - `sonic-home version` (after first build)

### Minimum deliverables by end of Day 1

- [ ] runtime interface contract merged (no placeholder TODO contract gaps)
- [ ] library/state schemas written and reviewed in code
- [ ] sonic-home module tree created with compiling CLI stubs
- [ ] initial bundle manifest schema doc committed in `modules/sonic-home/docs/`
- [ ] all touched packages compile locally

### Evidence (commands/tests/docs)

- command output highlights:
  - `go test ./...` (in `modules/sonic-home/`): pass
  - `go run ./cmd/sonic-home version`: `v0.1.0-dev`
  - `go run ./cmd/pack --dry-run --source . --output /tmp/sonic-home-manifest.dryrun.json`: pass
  - `go run ./cmd/sonic-home pack --dry-run=false --output /tmp/uhome-nest-draft.she`: pass
  - `go run ./cmd/sonic-home verify /tmp/uhome-nest-draft.she`: pass
  - `go run ./cmd/sonic-home pack --dry-run=false --source . --output /tmp/uhome-step3.she`: pass
  - `go run ./cmd/sonic-home verify /tmp/uhome-step3.she`: pass (component checksum + size checks active)
- files touched:
  - `go.work`
  - `modules/sonic-home/go.mod`
  - `modules/sonic-home/cmd/pack/main.go`
  - `modules/sonic-home/cmd/install/main.go`
  - `modules/sonic-home/cmd/serve/main.go`
  - `modules/sonic-home/cmd/sonic-home/main.go`
  - `modules/sonic-home/pkg/manifest/manifest.go`
  - `modules/sonic-home/pkg/manifest/manifest_test.go`
  - `modules/sonic-home/docs/BUNDLE-FORMAT.md`
  - `modules/sonic-home/README.md`
- blockers:
  - payload currently packages source tree files into `payload/base/` (docker/venv artifact flows pending)
  - cryptographic signing, delta updates, and USB auto-install are tracked as follow-up tasks

### End-of-day status

- runtime: unchanged today
- library: unchanged today
- state: unchanged today
- release: Ventoy integration brief created, promotion flow documented
- packager: step-3 baseline complete; ready for step 4 (sign/verify keys + docker/venv component modes)
- dev-flow: new structure implemented, Ventoy integrated into workflows

## Day 2

Date:

### Plan for today

- [ ] runtime: implement deterministic stub path for install/start/stop/remove
- [ ] library: implement index parsing + manifest loading
- [ ] state: migration bootstrap + open/create DB path
- [ ] release: draft structure for `docs/promotion.md`

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 3

Date:

### Plan for today

- [ ] runtime + CLI wiring: route commands through service layer
- [ ] library validation: required fields and unknown game handling
- [ ] state writes: install/start/stop/remove transitions persisted
- [ ] release notes: patch intake and build checkpoints drafted

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 4

Date:

### Plan for today

- [ ] integration pass across runtime/library/state flows
- [ ] add or fix unit tests for manifest + state behavior
- [ ] verify command failure modes return actionable non-zero exits
- [ ] release doc: add validation + promotion checkpoints

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 5

Date:

### Plan for today

- [ ] stabilize P0 implementation details
- [ ] close high-priority test failures
- [ ] sync roadmap/request docs to implementation reality
- [ ] list unresolved risks before week close

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 6

Date:

### Plan for today

- [ ] begin week-2 integration and cleanup
- [ ] tighten manifest validation edge cases
- [ ] verify migration behavior on fresh and existing DBs
- [ ] refine `docs/promotion.md` first full draft

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 7

Date:

### Plan for today

- [ ] end-to-end command flow verification
- [ ] reconcile runtime state with persisted state reporting
- [ ] eliminate critical TODO markers in touched P0 paths
- [ ] document remaining open questions for release

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 8

Date:

### Plan for today

- [ ] run targeted regression checks for install/start/stop/remove
- [ ] ensure `sonic library list` behavior is stable and deterministic
- [ ] prepare near-final P0 completion notes
- [ ] mark P1 spill items explicitly if not complete

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 9

Date:

### Plan for today

- [ ] run clean-checkout gate checks (`make build`, `make test`)
- [ ] address final P0 defects discovered in gate run
- [ ] finalize docs alignment (`docs/roadmap.md`, active index, promotion draft)
- [ ] draft candidate completion notes for submission summary

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Day 10

Date:

### Plan for today

- [ ] confirm all P0 definitions of done are satisfied
- [ ] decide P1 completion vs spill with explicit rationale
- [ ] prepare merge-ready checklist and residual risk notes
- [ ] update `dev/submissions/completed-summary.md` when slice is complete

### Evidence (commands/tests/docs)

- command output highlights:
- files touched:
- blockers:

### End-of-day status

- runtime:
- library:
- state:
- release:

## Slice Close Checklist

- [ ] all P0 definitions of done met
- [ ] `make build` passed on clean checkout
- [ ] `make test` passed on clean checkout
- [ ] no critical TODOs remain in touched P0 paths
- [ ] active index updated (completed or remaining items)
- [ ] completed summary updated with meaningful outcomes
- [ ] unresolved risks and follow-up tasks documented
