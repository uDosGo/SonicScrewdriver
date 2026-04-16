# Sonic Home Bundle Format (`.she`)

Status: draft bootstrap for `v1.1.1`.

## Goals

- single-file distribution artifact for uHomeNest
- manifest-first validation and reproducible metadata
- easy upgrade path to future full Sonic tooling

## File Layout (Draft)

```text
<bundle>.she
├── header.sheh
├── manifest.json
├── signature.sig
├── payload/
│   ├── base/
│   ├── docker/
│   ├── venv/
│   ├── assets/
│   └── scripts/
└── updates/
```

## Draft Bundle + Verify Support (Implemented)

Current bootstrap supports dry-run manifest generation and draft archive output:

- command:
  - `go run ./cmd/pack --dry-run --source . --output manifest.dryrun.json`
  - `go run ./cmd/pack --dry-run=false --output uhome-nest-draft.she`
  - `go run ./cmd/verify ./uhome-nest-draft.she`
- output:
  - draft `manifest.json` metadata for source path, version, channel, and architecture
  - tar-based `.she` draft archive with required baseline entries
  - structural verification for required entries and manifest parse
- non-goal in this step:
  - payload population from real artifacts, cryptographic signatures, delta generation

## Baseline Manifest Fields

- `schema_version`
- `bundle_id`
- `version`
- `release_channel`
- `build_date`
- `architecture`
- `source_path`
- `source_digest`
- `components`
- `scripts`
- `signatures`
- `update_info`

## Follow-Up Tasks

- signing:
  - add Ed25519 key loading and detached signature generation
  - add `verify` command for signature + checksum checks
- delta updates:
  - add `pack --delta-from --delta-to` flow
  - evaluate `bsdiff` and `xdelta` compatibility
- USB auto-install:
  - define portable USB layout and auto-detection scripts
  - add installer mode parity for `--usb` and non-interactive preseed

