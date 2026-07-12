---
title: "SonicScrewdriver Dev Plan (uCore Integration + Tasker Flow)"
status: active
last_updated: 2026-07-09T16:00:00+08:00
category: planning
tags: [sonic, ucore, tasker, usx, mcp, snackbar]
description: "Execution plan for integrating SonicScrewdriver with uCore developer tools, tasker dev-flow, USX style kit, skills, and snackbar/MCP surfaces. Re-baselined 2026-07-09 after repo assessment found 0/14 tasks started."
---

# SonicScrewdriver Dev Plan (2026-07)

This is the active engineering plan for migrating SonicScrewdriver toward
shared uCore developer tooling and service conventions while preserving
hardware-safe, local-first USB workflows.

## Goals

1. Adopt uCore-style development workflow with tasker as the canonical task lane.
2. Refactor operator-facing UI/document surfaces toward USX token conventions.
3. Integrate uCore Skills and MCP-compatible tool/service surfaces.
4. Align runtime and CLI eventing with snackbar + spool observability.

## Source of Truth

- Task flow model: `../../uCore/.tasker.dev-flow.yaml`
- Dev tools posture: `../../uCore/docs/DEVMODE_CODE_ANALYSIS_SKILLS.md`
- MCP conventions: `../../uCore/docs/MCP_SETUP.md` and `../../uCore/docs/mcp-policy.md`
- Snackbar and spool model: `../../uCore/docs/SNACKS_SYSTEM_SPEC.md`, `../../uCore/docs/SPOOL_SPEC.md`, and `../../uCore/docs/FEED_SYSTEM_SPEC.md`

## Assessment (2026-07-09)

Repo assessed against Dev Plan deliverables. All 5 workstreams were at 0% implementation:
no event schema, no tests, no MCP server/skills directory, no snackbar envelopes,
no spool adapters, no USX token application, and stale Go v1 CI. Sprint
`sprint.2026-07-08.sonic-ucore-integration` superseded by
`sprint.2026-07-09.sonic-ucore-integration-v2` with realistic 3-wave sequencing
(Foundation → Integration → Stabilize).

## Workstream A - uCore Developer Tools Integration

### Scope

- Standardize SonicScrewdriver dev checks and diagnostics around uCore-style
  categories (analysis, maintenance, workflow, orchestration).
- Add shared telemetry wrappers for CLI, bootloader build, and USB lifecycle
  operations where safe.

### Deliverables

1. Dev-tools parity map: Sonic commands and equivalent uCore capability classes.
2. Event adapter for spool/feed write shape:
   `timestamp`, `module`, `level`, `message`, `tags`.
3. Contract drift checks in CI for lifecycle state names and event schema.

### Exit Criteria

- Sonic command and build events can be consumed by uCore spool tooling.
- One complete USB build/test workflow is represented with uCore-style event logs.

## Workstream B - Adopt Tasker Dev Flow

### Scope

- Add Sonic-local `.tasker.dev-flow.yaml` aligned to uCore schema shape.
- Move roadmap execution tracking into lane-based task IDs and sprint metadata.

### Deliverables

1. Initial lanes:
   - `cli-runtime`
   - `bootloader-ux`
   - `skills-mcp`
   - `runtime-snackbar`
   - `maintenance`
2. Task UID conventions: `task.sonic.<lane>.<nnn>`.
3. Sprint block with `start`, `end`, `status`, and completion counters.

### Exit Criteria

- Active planning changes are updated in `.tasker.dev-flow.yaml` first.
- Major roadmap milestones map to task IDs with status transitions.

## Workstream C - USX Style Kit and Operator Surface Refactor

### Scope

- Align docs/operator-facing visual surfaces with shared USX token and style
  conventions where frontend surfaces exist.
- Preserve CLI-first functionality while improving consistency for onboarding,
  status output, and docs UI surfaces.

### Deliverables

1. Tokenized style baseline for docs/operator pages.
2. Shared visual primitives for status, warnings, and action/result output blocks.
3. Migration guide for any future Vue/web control surface tied to Sonic workflows.

### Exit Criteria

- Operator-facing surfaces rely on shared token conventions.
- No critical operator route depends on one-off style systems.

## Workstream D - uCore Skills + MCP Integration

### Scope

- Define Sonic MCP tool surfaces that align with uCore skill posture:
  diagnostics, build orchestration, device inventory, and repair helpers.
- Keep hardware-sensitive actions gated and explicit.

### Deliverables

1. Skill catalog: Sonic-native vs uCore-compatible shared skill patterns.
2. MCP manifest update for Sonic tools with safety policy notes.
3. Structured skill execution telemetry bridged to spool events.

### Exit Criteria

- Sonic skills are discoverable through MCP with stable naming and response shape.
- Failed tool runs emit structured diagnostics suitable for snackbar/spool.

## Workstream E - Snackbar + MCP Runtime Operations

### Scope

- Implement snackbar-compatible response envelopes for CLI/runtime-facing actions.
- Ensure bootloader/USB operations remain auditable with local-safe defaults.

### Deliverables

1. Event taxonomy for `usb`, `bootloader`, `security`, `mint`, and `device` lanes.
2. Shared response envelope shape for success/warn/error states.
3. Operator diagnostics summary for recent event stream and failures.

### Exit Criteria

- Operator actions emit consistent, actionable status feedback.
- MCP and runtime failures are traceable through spool logs and status endpoints.

## Phase Plan (Reset 2026-07-09)

### Wave 1: Foundation (Days 1–7)

- `task.sonic.maintenance.001` — Create dev-tools parity table against uCore developer tools.
- `task.sonic.cli-runtime.001` — Add end-to-end tests for core CLI command groups.
- `task.sonic.cli-runtime.002` — Standardize error envelopes for CLI runtime failures.
- `task.sonic.bootloader-ux.001` — Define bootloader status taxonomy for diagnostics.
- `task.sonic.skills-mcp.001` — Draft Sonic skill catalog with uCore compatibility labels.

### Wave 2: Integration (Days 8–14)

- `task.sonic.runtime-snackbar.001` — Define snackbar response envelope for runtime and mint workflows.
- `task.sonic.runtime-snackbar.002` — Emit structured spool/feed events for USB lifecycle changes.
- `task.sonic.cli-runtime.003` — Emit structured spool events for major command lifecycle stages.
- `task.sonic.bootloader-ux.002` — Add consistent result blocks for install/check flows.
- `task.sonic.skills-mcp.002` — Add MCP manifest entries for Sonic tool surfaces.

### Wave 3: Stabilize (Days 15–21)

- `task.sonic.maintenance.002` — Add CI check for shared event schema drift (replace stale Go CI).
- `task.sonic.runtime-snackbar.003` — Add operator diagnostics summary for recent event stream.
- `task.sonic.bootloader-ux.003` — Add arm64 and BIOS/UEFI lifecycle state checks.
- `task.sonic.skills-mcp.003` — Bridge skill execution telemetry into spool events.

## Risks and Mitigations

- Risk: Hardware workflows are disrupted by tooling refactor.
  Mitigation: Keep execution paths feature-flagged and preserve existing CLI contracts.
- Risk: MCP/skills drift across repos.
  Mitigation: Add naming/schema lint and manifest checks in CI.
- Risk: Event noise reduces operator signal.
  Mitigation: Introduce event severity + scoped tags per subsystem.

## Definition of Done

1. Sonic planning and sprint tracking run through tasker dev-flow.
2. Shared uCore developer-tool patterns are documented and operational.
3. Operator/document surfaces are aligned to USX style-kit conventions.
4. Skills and MCP endpoints are integrated with structured snackbar/spool
   observability.
