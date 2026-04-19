# `sonic-screwdriver` Dev Flow

The `dev` directory implements an opinionated, lightweight development flow for the Sonic Family.

**📚 [Complete Dev Flow Compendium](COMPLETE-DEV-FLOW-COMPENDIUM.md)** — Master document containing all briefs, workflows, and processes in one unified system.

## Core Principles

- **Active work stays visible** - current execution lives in `dev/active/`
- **Integration points are explicit** - cross-component work lives in `dev/integration/`
- **Process is documented** - workflows and checklists live in `dev/process/`
- **Old work gets composted** - completed or obsolete material moves to `dev/compost/`

## Structure

```
dev/
├── active/                  # Current execution (2-week slices)
│   ├── active-index.md       # Active request index
│   ├── execution-notes.md    # Daily execution log
│   └── completed-summary.md  # Recently completed work
│
├── integration/             # Cross-component coordination
│   ├── UHOMENEST-*.md        # uHomeNest integration briefs
│   └── VENTOY-*.md           # Ventoy integration briefs
│
├── process/                 # Workflow definitions
│   ├── checklists/           # Repeatable validation steps
│   ├── workflows/            # End-to-end process docs
│   └── templates/            # Reusable doc templates
│
├── compost/                 # Archived material
│   ├── requests/             # Completed request archives
│   ├── submissions/          # Historical submission logs
│   └── notes/                # Obsolete technical notes
│
└── README.md                # This file
```

## Workflow

### 1. Active Execution (2-week slices)
- Start each slice by updating `dev/active/active-index.md`
- Track daily progress in `dev/active/execution-notes.md`
- Move completed work to `dev/active/completed-summary.md`

### 2. Integration Coordination
- Create briefs in `dev/integration/` for cross-component work
- Use `VENTOY-` prefix for Ventoy-related integration
- Use `UHOMENEST-` prefix for uHomeNest-related integration

### 3. Process Documentation
- Add repeatable checklists to `dev/process/checklists/`
- Document end-to-end workflows in `dev/process/workflows/`
- Store reusable templates in `dev/process/templates/`

### 4. Composting
- Move obsolete material to `dev/compost/` with clear rationale
- Keep compost organized by type (requests, submissions, notes)
- Add compost date and reason in file headers

## Ventoy Integration

The Ventoy submodule (`Ventoy/`) is integrated via:
- Build scripts in `modules/ventoy/`
- Integration briefs in `dev/integration/VENTOY-*.md`
- Promotion flow documentation in `docs/promotion.md`

## Rules

1. **Keep active lanes small** - no more than 10 active requests at once
2. **Update daily** - execution notes reflect real progress
3. **Compost aggressively** - move completed work out of active lanes
4. **Document blockers explicitly** - use `[!]` marker for blocked items
5. **Align with roadmap** - active index mirrors `docs/roadmap.md` milestones

## Quick Commands

```bash
# View active work
cat dev/active/active-index.md

# Update execution notes
$EDITOR dev/active/execution-notes.md

# Compost completed work
mv dev/active/some-request.md dev/compost/requests/

# Check integration status
grep -r "TODO\|TBD\|WIP" dev/integration/
```
