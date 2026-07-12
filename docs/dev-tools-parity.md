# SonicScrewdriver ↔ uCore Developer Tools Parity Table

**Created:** 2026-07-09
**Status:** Foundation deliverable — `task.sonic.maintenance.001`
**Purpose:** Map every Sonic command to uCore capability classes for integration planning.

---

## uCore Capability Classes

uCore developer tools are organized into four capability domains:

| Domain | Description |
|---|---|
| **analysis** | Inspection, diagnostics, scanning, auditing |
| **maintenance** | Repair, cleanup, archiving, technical debt |
| **workflow** | Lifecycle operations, orchestration, build pipelines |
| **orchestration** | Multi-step automation, agent coordination, MCP dispatch |

---

## Command Parity Map

### `sonic usb` — USB Creation & Management

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic usb create` | workflow, orchestration | Plate scaffolding (`plate_refresh/`) | Multi-step partition + format + bootloader install pipeline |
| `sonic usb list` | analysis | Device scanner (`device_scanner.py`) | Read-only, safe |
| `sonic usb info` | analysis | Spool reader (`spool_reader.py`) | Partition metadata inspection |
| `sonic usb destroy` | maintenance | DESTROY/REBUILD protocol | Destructive with confirmation gate |

### `sonic security` — Security Device Enrollment

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic security enroll` | workflow | Key management (`keygen.py`) | FIDO2/GPG/SSH enrollment |
| `sonic security list` | analysis | Catalog service (`catalog/`) | Read-only device inventory |
| `sonic security revoke` | maintenance | Revocation pipeline | Destructive with confirmation |
| `sonic security generate` | workflow | Key generation | SSH/GPG key creation |
| `sonic security gpg-init` | workflow | Key generation | GPG initialization |

### `sonic mint` — Linux Mint ISO Customization

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic mint build` | orchestration | Plate refresh (`plate_refresh/refresh.py`) | Multi-stage extract→customize→repack pipeline |
| `sonic mint verify` | analysis | Checksum validation | Read-only |
| `sonic mint preseed` | analysis | Config extraction | Read-only metadata export |

### `sonic bootloader` — Bootloader Management

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic bootloader install` | workflow | Package installation (`package_manager/`) | Writes to ESP, requires device path |
| `sonic bootloader status` | analysis | Health check (`GET /api/health/full`) | Read-only status check |
| `sonic bootloader remove` | maintenance | Package removal | Destructive with confirmation |
| `sonic bootloader build` | orchestration | Build pipeline | Cross-compile UEFI + BIOS |
| `sonic bootloader test` | analysis | QEMU test harness | Read-only (emulated) |

### `sonic device` — Device Library

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic device scan` | analysis | Device scanner | Read-only hardware enumeration |
| `sonic device identify` | analysis | Hardware fingerprinting | Read-only capability detection |
| `sonic device lookup` | analysis | Catalog search (`/api/catalog/search`) | Read-only database query |
| `sonic device add` | workflow | Catalog entry creation | Write operation (local DB) |
| `sonic device remove` | maintenance | Catalog entry removal | Destructive with confirmation |
| `sonic device repurpose` | orchestration | Device transformation pipeline | Router→beacon, router→OpenWrt |
| `sonic device flash` | workflow | Firmware deployment | Requires serial port, destructive |
| `sonic device ewaste` | analysis | Classification pipeline | Read-only with guidance output |
| `sonic device submit` | workflow | Registry submission | Network write, POST to registry API |

### `sonic mesh` — Mesh Networking

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic mesh init` | orchestration | Service daemon start | Background process + mDNS discovery |
| `sonic mesh discover` | analysis | Service discovery | Read-only network scan |
| `sonic mesh connect` | workflow | Peer connection | Network I/O |
| `sonic mesh status` | analysis | Health check | Read-only daemon state |
| `sonic mesh disconnect` | maintenance | Connection teardown | Graceful disconnect |
| `sonic mesh stop` | maintenance | Service daemon stop | Process termination |

### `sonic chasis` — CHASIS Game Library

| Command | uCore Domain | uCore Equivalent | Notes |
|---|---|---|---|
| `sonic chasis add` | workflow | Catalog entry creation | EFI binary ingestion |
| `sonic chasis remove` | maintenance | Catalog entry removal | Destructive |
| `sonic chasis list` | analysis | Catalog listing | Read-only |
| `sonic chasis install` | workflow | Package installation | Copies EFI to USB |
| `sonic chasis launch` | workflow | Process launcher | Runs game binary |
| `sonic chasis export` | workflow | Archive creation | Portable tar.gz export |
| `sonic chasis import-lib` | workflow | Archive ingestion | tar.gz import |

---

## Summary Statistics

| Domain | Command Count |
|---|---|
| analysis | 16 |
| workflow | 14 |
| maintenance | 7 |
| orchestration | 4 |
| **Total** | **41 command/option combinations** |

---

## Event Adapter Shape

All Sonic commands should emit events in the uCore-compatible shape:

```python
{
    "timestamp": "2026-07-09T22:00:00+08:00",
    "module": "sonic.usb",
    "level": "INFO",
    "message": "USB create started for /dev/sdb",
    "tags": ["usb", "create", "lifecycle"]
}
```

Mapped from uCore `SpoolEntry` dataclass:
- `timestamp` → ISO 8601
- `level` → INFO, WARNING, ERROR, DEBUG, CRITICAL
- `module` → `sonic.<command_group>` (e.g., `sonic.usb`, `sonic.mint`)
- `message` → Human-readable event description
- `tags` → Lowercase, kebab-case labels for filtering