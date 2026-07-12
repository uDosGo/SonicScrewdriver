# SonicScrewdriver v2.1.0

## Universal USB Bootloader & System Toolkit

SonicScrewdriver v2 is a **Python CLI + C bootloader** for creating
bootable USB drives, managing device firmware, enrolling security keys,
and deploying Linux Mint systems — with uCore-compatible spool
observability, MCP tool surfaces, and USX design token alignment.

> ✅ **Status:** Sprint `sprint.2026-07-09.sonic-ucore-integration-v2`
> complete (14/14). uCore dev-tools parity, snackbar/spool telemetry,
> MCP manifest, bootloader lifecycle taxonomy, and CI schema checks are
> all operational.

## 🎯 What It Does (v2 Python CLI)

```
sonic usb create <device>              — Triple-partition USB (ESP + ext4 + exFAT)
sonic security enroll <type>           — FIDO2/U2F, GPG, or SSH key enrollment
sonic mint build <input> <output>      — Linux Mint ISO customization
sonic mesh init|connect|discover       — Peer-to-peer mesh networking
sonic chasis add|launch|list           — Game library management
sonic bootloader install|build|test    — SonicScrewloader management
sonic device scan|identify|lookup      — Device library management
sonic device add|remove|repurpose      — CRUD + router repurposing
sonic diagnostics summary|health       — Spool event stream + system health
```

## 🚀 Quick Start

```bash
# Install Python CLI (with dev dependencies for tests)
cd cli && pip install -e ".[dev]"

# Run tests (11 passing)
python -m pytest tests/ -v

# View help
sonic --help

# Run diagnostics
sonic diagnostics summary
sonic diagnostics health
```

## 🏗️ Project Structure

```
SonicScrewdriver/
├── bootloader/             # C/asm bootloader (Teletext, UEFI, BIOS, ARM64)
│   ├── src/                # teletext.c, framebuffer.c, detect.c, menu.c, chainload.c
│   ├── include/            # teletext.h, detect.h, menu.h (lifecycle enum)
│   ├── config/menus/       # YAML menu configs (mac.yaml, pc.yaml, bios.yaml)
│   └── Makefile            # UEFI (gnu-efi) + BIOS targets
├── cli/                    # Python CLI (v2)
│   ├── sonic/
│   │   ├── commands/       # usb, security, mint, mesh, chasis, bootloader, device, diagnostics
│   │   ├── lib/            # envelope, spool, usx_theme, mcp_bridge
│   │   └── data/devices/   # YAML device entries (PCs, routers, ESP32)
│   ├── tests/              # pytest suites (11 tests, all passing)
│   └── setup.py            # v2.1.0
├── mcp/                    # MCP tool manifest
│   └── sonic-mcp-manifest.json  # 32 tool definitions with safety labels
├── mint/                   # Linux Mint ISO build system
│   ├── build-iso.sh        # Full ISO build pipeline
│   ├── chroot-customize.sh # Hostname, locale, user, desktop config
│   ├── install-sonic.sh    # Sonic CLI + dependencies
│   └── overlay/            # Plymouth theme, skel, backgrounds
├── recovery/               # Recovery tools
│   └── scripts/            # disk-repair, password-reset, data-recovery, memtest
├── docs/                   # Documentation
│   ├── SONICSCREWDRIVER_DEV_PLAN.md    # Active dev plan
│   ├── dev-tools-parity.md             # Sonic ↔ uCore capability map
│   ├── bootloader-status-taxonomy.md   # 17 lifecycle states
│   ├── sonic-skill-catalog.md          # 40 MCP tools
│   ├── usx-sonic-tokens.md             # USX token reference
│   └── usx-migration-guide.md          # Vue/web surface guide
├── .github/workflows/ci.yaml  # CI: pytest + schema validation + spool lint
├── .tasker.dev-flow.yaml      # Canonical tasker sprint tracker
└── version                     # v2.1.0
```

## 📖 Documentation

- **[docs/USB-CREATION.md](docs/USB-CREATION.md)** — Full USB creation guide
- **[docs/SONICSCREWDRIVER_DEV_PLAN.md](docs/SONICSCREWDRIVER_DEV_PLAN.md)** — Active dev plan (uCore integration, USX, skills/MCP, snackbar)
- **[.tasker.dev-flow.yaml](.tasker.dev-flow.yaml)** — Canonical tasker lane + sprint execution tracker
- **[docs/dev-tools-parity.md](docs/dev-tools-parity.md)** — Sonic ↔ uCore capability mapping
- **[docs/sonic-skill-catalog.md](docs/sonic-skill-catalog.md)** — 40 MCP tools with safety classifications
- **[docs/bootloader-status-taxonomy.md](docs/bootloader-status-taxonomy.md)** — 17 lifecycle states
- **[docs/usx-sonic-tokens.md](docs/usx-sonic-tokens.md)** — USX design token reference for CLI/docs
- **[docs/usx-migration-guide.md](docs/usx-migration-guide.md)** — Future Vue/web surface migration guide
- **[mcp/sonic-mcp-manifest.json](mcp/sonic-mcp-manifest.json)** — MCP tool manifest (32 tools)

## USX Design Alignment

SonicScrewdriver CLI output follows [USX token conventions](docs/usx-sonic-tokens.md):
- Status icons and colors map to `--usx-color-success` / `--usx-color-warning` / `--usx-color-error`
- Result blocks use `rich.panel.Panel.fit()` with semantic border colors
- Snackbar response envelopes (`SnackbarResponse`) align with uCore snackbar status shapes
- Future web surfaces should import `@udos/usx-tokens` (see [migration guide](docs/usx-migration-guide.md))

## Related Repositories

- **uCore** — Backend services, MCP, skills, developer tooling, USX tokens
- **uServer** — Backend services, secret store, API central
