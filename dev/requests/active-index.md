# Active Request Index

This is the compact active request index for `sonic-screwdriver`.

## 2-Week Execution Slice (Aligned to `vA1.1.0`)

Timebox: next 10 working days.  
Owner labels:
- `owner/core-runtime`: container + CLI runtime wiring
- `owner/library`: manifest/index parsing and install inputs
- `owner/state`: sqlite schema + persistence lifecycle
- `owner/release`: Ventoy/promotion docs and release flow notes
- `owner/packager`: sonic-home bundling/install flow

### P0 (Must Ship This Slice)

- runtime boundary + command wiring
  - owner: `owner/core-runtime`
  - scope: implement runtime interface in `internal/container/` and wire `install/start/stop/remove` in `cmd/sonic/main.go`
  - definition of done:
    - `sonic install <game>` resolves manifest and calls runtime path
    - `sonic start <game>` and `sonic stop <game>` transition runtime state through service layer
    - `sonic remove <game>` removes tracked install/runtime metadata
    - command failures return non-zero exit code with actionable message

- library manager + manifest validation
  - owner: `owner/library`
  - scope: parse `library/index.yaml`, load game manifests, validate required manifest fields
  - definition of done:
    - `sonic library list` returns curated entries from YAML index
    - `sonic install <game>` rejects unknown game IDs and malformed manifests
    - manifest validation checks are covered by unit tests

- sqlite state bootstrap + persistence
  - owner: `owner/state`
  - scope: initialize sqlite DB, apply first migration, persist install/runtime records
  - definition of done:
    - first run initializes database automatically
    - install/start/stop/remove operations write deterministic state transitions
    - migration boot path is exercised in tests

- build and test gate for slice
  - owner: `owner/core-runtime`
  - scope: enforce clean pass for build + tests before closing slice
  - definition of done:
    - `make build` passes on clean checkout
    - `make test` passes on clean checkout

### P1 (Should Progress In Parallel, Can Spill)

- Ventoy patch + promotion flow draft
  - owner: `owner/release`
  - scope: define `modules/ventoy/` patch handling expectations and first draft of `docs/promotion.md`
  - definition of done:
    - doc covers patch intake, local build, validation, and promotion checkpoints
    - known open questions are listed explicitly (do not hide unresolved steps)

- sonic-home bootstrap (`v1.1.1` track)
  - owner: `owner/packager`
  - scope: scaffold `modules/sonic-home/` module with `pack/install/serve` CLI stubs, baseline `.she` manifest schema, and starter docs
  - definition of done:
    - `modules/sonic-home/` tree exists with compiling Go entrypoints
    - initial manifest schema documented in `modules/sonic-home/docs/BUNDLE-FORMAT.md`
    - `pack` command supports dry-run manifest generation path (no full payload required yet)
    - follow-up tasks explicitly listed for signing, delta, and USB auto-install

## Tracking Rhythm

- day 1-2: finalize interfaces and schema choices
- day 3-6: implement runtime/library/state paths in parallel
- day 7-8: integrate CLI command flows and close failing tests
- day 9-10: stabilization, docs sync, and `vA1.1.0` gate check

## Exit Check For This Slice

- all P0 definitions of done met
- no critical TODOs remain in touched runtime/library/state paths
- roadmap/docs remain aligned with implemented behavior

## Rule

When a request is completed or absorbed into stable docs, remove it from this index instead of keeping stale trackers.
