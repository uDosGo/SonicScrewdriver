# Active Request Index

This is the compact active request index for `sonic-family`.

## Current Active Requests

- docker backend scaffold wiring
  - outcome: wire `install/start/stop/remove` commands to a deterministic runtime stub and clear Docker handoff

- library manager and manifest loading
  - outcome: parse `library/index.yaml`, validate manifests, and support `sonic library list` + `sonic install <game>` flow

- sqlite state tracking
  - outcome: persist install and runtime state with simple migration bootstrapping

- ventoy integration notes
  - outcome: define patch flow and reproducible build notes before first media-based release

## Rule

When a request is completed or absorbed into stable docs, remove it from this index instead of keeping stale trackers.
