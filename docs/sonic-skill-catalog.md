# SonicScrewdriver — Skill Catalog (MCP Tool Inventory)

**Created:** 2026-07-09
**Status:** Foundation deliverable — `task.sonic.skills-mcp.001`
**Purpose:** Inventory all Sonic capabilities as MCP tools with uCore domain
compatibility labels and safety classifications.

---

## Safety Classifications

| Label | Meaning | MCP Behavior |
|---|---|---|
| **safe** | Read-only, no side effects. No confirmation required. | Auto-execute |
| **caution** | Reads system state, may prompt user. | Ask once per session |
| **gated** | Writes to hardware/filesystem. Destructive potential. | Require explicit confirmation every invocation |

---

## Tool Catalog

### usb — USB Creation & Management

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.usb.list` | List connected USB block devices | analysis | safe | Read-only device enumeration |
| `sonic.usb.info` | Show partition layout and metadata for a device | analysis | safe | Read-only partition inspection |
| `sonic.usb.create` | Create triple-partition Sonic USB (ESP + ext4 + exFAT) | workflow, orchestration | gated | Destructive — wipes target device |
| `sonic.usb.destroy` | Wipe Sonic partitions and restore single volume | maintenance | gated | Destructive — wipes target device |

### security — Security Device Enrollment

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.security.enroll` | Enroll FIDO2/GPG/SSH/TPM security device | workflow | gated | Writes security key material |
| `sonic.security.list` | List enrolled security devices | analysis | safe | Read-only device inventory |
| `sonic.security.revoke` | Revoke an enrolled security device | maintenance | gated | Irreversible revocation |
| `sonic.security.generate-ssh` | Generate Sonic SSH key pair (ed25519/rsa/ecdsa) | workflow | caution | Creates key files on disk |
| `sonic.security.generate-gpg` | Generate Sonic GPG key pair | workflow | caution | Creates key files on disk |

### mint — Linux Mint ISO Customization

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.mint.build` | Customize Linux Mint ISO (hostname, user, branding, packages) | orchestration | caution | Reads/writes ISO files on disk |
| `sonic.mint.verify` | Verify ISO checksum integrity | analysis | safe | Read-only checksum validation |
| `sonic.mint.preseed` | Extract preseed configuration from customized ISO | analysis | safe | Read-only config extraction |

### bootloader — Bootloader Management

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.bootloader.install` | Install SonicScrewloader to USB ESP | workflow | gated | Writes to device ESP |
| `sonic.bootloader.status` | Check bootloader installation health | analysis | safe | Read-only status check |
| `sonic.bootloader.remove` | Remove SonicScrewloader from device | maintenance | gated | Destructive — removes boot entry |
| `sonic.bootloader.build` | Cross-compile bootloader (UEFI x86_64/aarch64, BIOS) | orchestration | caution | Requires gnu-efi toolchain |
| `sonic.bootloader.test` | Test bootloader in QEMU emulator | analysis | safe | Read-only — emulated execution |

### device — Device Library

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.device.scan` | Scan for connected devices (USB, PCI, Bluetooth) | analysis | safe | Read-only hardware enumeration |
| `sonic.device.identify` | Identify hardware capabilities of a device | analysis | safe | Read-only capability detection |
| `sonic.device.lookup` | Query local device database | analysis | safe | Read-only DB query |
| `sonic.device.add` | Add custom device entry to local database | workflow | safe | Local DB write, no hardware impact |
| `sonic.device.remove` | Remove device entry from local database | maintenance | safe | Local DB write, no hardware impact |
| `sonic.device.repurpose` | Repurpose router/device (beacon, OpenWrt, mesh) | orchestration | gated | May flash firmware — destructive |
| `sonic.device.flash` | Flash firmware to ESP32/router via serial | workflow | gated | Destructive — writes firmware |
| `sonic.device.ewaste` | Classify device for e-waste disposal | analysis | safe | Read-only classification guidance |
| `sonic.device.submit` | Submit device to global Sonic registry | workflow | caution | Network POST to remote API |

### mesh — Mesh Networking

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.mesh.init` | Initialize mesh node on this machine | orchestration | caution | Starts background daemon |
| `sonic.mesh.discover` | Discover nearby Sonic mesh nodes | analysis | safe | Read-only network scan |
| `sonic.mesh.connect` | Connect to a specific mesh node | workflow | safe | Network connect, no local writes |
| `sonic.mesh.status` | Show mesh network status and connected peers | analysis | safe | Read-only daemon state |
| `sonic.mesh.disconnect` | Disconnect from a mesh node | maintenance | safe | Graceful disconnect |
| `sonic.mesh.stop` | Stop the mesh networking daemon | maintenance | safe | Process termination |

### chasis — CHASIS Game Library

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.chasis.add` | Add EFI game binary to CHASIS library | workflow | safe | Local catalog entry |
| `sonic.chasis.remove` | Remove game from CHASIS library | maintenance | safe | Local catalog removal |
| `sonic.chasis.list` | List games in CHASIS library | analysis | safe | Read-only catalog |
| `sonic.chasis.install` | Install game EFI binary to USB drive | workflow | caution | Writes to USB |
| `sonic.chasis.launch` | Launch a CHASIS game EFI binary | workflow | caution | Executes binary |
| `sonic.chasis.export` | Export CHASIS library as tar.gz archive | workflow | safe | Creates archive on disk |
| `sonic.chasis.import` | Import CHASIS library from tar.gz archive | workflow | safe | Extracts archive on disk |

### diagnostics — Cross-Cutting (New)

| Tool Name | Description | uCore Domain | Safety | Notes |
|---|---|---|---|---|
| `sonic.diagnostics.summary` | Show recent spool event summary | analysis | safe | Read-only — reads spool journal |
| `sonic.diagnostics.health` | Full system health check (USB, bootloader, device DB) | analysis | safe | Read-only multi-subsystem scan |

---

## Summary

| Safety | Count | % |
|---|---|---|
| safe | 24 | 60% |
| caution | 9 | 23% |
| gated | 7 | 18% |
| **Total** | **40** | |

| uCore Domain | Count |
|---|---|
| analysis | 18 |
| workflow | 13 |
| maintenance | 6 |
| orchestration | 4 |

---

## MCP Naming Convention

All Sonic MCP tools follow the pattern:

```
sonic.<command-group>.<action>
```

Examples:
- `sonic.usb.create`
- `sonic.security.enroll`
- `sonic.mint.build`
- `sonic.bootloader.install`
- `sonic.device.scan`
- `sonic.mesh.init`
- `sonic.chasis.add`
- `sonic.diagnostics.summary`