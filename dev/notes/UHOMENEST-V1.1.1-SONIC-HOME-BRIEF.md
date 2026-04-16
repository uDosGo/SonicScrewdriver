# Sonic-Home — Lite Packager Module for uHomeNest

**Document ID:** `UDN-SONIC-001`  
**Status:** Active  
**Version:** 1.0.0  
**Date:** 2026-04-16  
**Related:** [uHomeNest v1.0.0 Dev Brief](./UHOMENEST-V1-DEV-BRIEF.md), [Matter+HA Integration Plan](./UDN-INTEGRATION-001.md)

---

## Objective

Create **`sonic-home`** — a **lite packager, handler, distribution, and installer helper** module for uHomeNest v1.0.0, based on **uDOS `sonic-express`** patterns.

**Purpose:** Provide a **minimal, fast, dependency-light** pathway to:
- Package uHomeNest components into distributable bundles
- Handle installation on target hardware (single command)
- Distribute updates via USB, network, or local archive
- Link to future **Sonic-family** tooling without requiring full Sonic stack

**Project boundary note:** uHomeNest and sonic-screwdriver are separate systems. Both can
host a `sonic-home` module with shared standards and compatible interfaces, but each keeps
independent runtime policy, implementation details, and release cadence.

**Principle:** `sonic-home` is the **on-ramp** to Sonic-family for uHomeNest; it works standalone but can be upgraded to full Sonic when available.

---

## Architecture

### Sonic Family Layering

```
┌─────────────────────────────────────────────────────────────────┐
│                      Sonic-Family (future)                      │
│  Full installer, Ventoy, dual-boot, recovery, fleet management │
│  Repository: sonic-screwdriver, sonic-studio, sonic-orchestra   │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ optional upgrade path
                              │
┌─────────────────────────────────────────────────────────────────┐
│                    sonic-home (v1.0)                    │
│  Lite packager, handler, installer helper for uHomeNest        │
│  • Bundle generation (.she, .bundle)                           │
│  • Single-command install                                       │
│  • USB auto-detection                                           │
│  • Update channels (stable, beta, edge)                         │
│  • Local distribution over LAN                                  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ consumes
                              │
┌─────────────────────────────────────────────────────────────────┐
│                         uHomeNest v1.0                          │
│  Core: Jellyfin + ~/media/ + API + UI + Matter/HA integrations │
└─────────────────────────────────────────────────────────────────┘
```

### Module Placement

```
uHomeNest/
├── sonic-home/              # New module
│   ├── cmd/
│   │   ├── pack/                    # Bundle packager
│   │   ├── install/                 # Installer handler
│   │   └── serve/                   # Local distribution server
│   ├── pkg/
│   │   ├── bundle/                  # .she bundle format
│   │   ├── manifest/                # Manifest schema
│   │   ├── handler/                 # Install handlers
│   │   └── transport/               # USB, HTTP, file transport
│   ├── templates/
│   │   ├── Dockerfile.install       # Install-time container
│   │   └── preseed.cfg              # Auto-install answers
│   ├── scripts/
│   │   ├── build-bundle.sh
│   │   ├── install-from-usb.sh
│   │   └── create-update-channel.sh
│   └── docs/
│       ├── BUNDLE-FORMAT.md
│       ├── INSTALLER-API.md
│       └── UPGRADE-TO-SONIC.md
├── server/                          # Existing (unmodified)
├── ui/                              # Existing (unmodified)
└── media-vault/                     # Existing (unmodified)
```

---

## Bundle Format (.she — Sonic Home)

### Design Principles

- **Self-contained**: One file, no external dependencies at install time
- **Verifiable**: Signed manifest, checksums
- **Incremental**: Delta updates between versions
- **Portable**: Works on Ubuntu 20.04–24.04, Debian 12, Raspberry Pi OS

### Bundle Structure

```
uhome-nest-v1.0.0.she
├── header.sheh                     # 512-byte header (magic, version, size)
├── manifest.json                   # Bundle metadata
├── signature.sig                   # Ed25519 signature
├── payload/
│   ├── base/                       # Base system files
│   │   ├── usr/local/bin/uhome-api
│   │   ├── etc/systemd/uhome*.service
│   │   └── opt/uhome-integrations/manifest.json
│   ├── docker/                     # Docker images (layer tarballs)
│   │   ├── jellyfin.tar
│   │   ├── matter.tar
│   │   └── matter-server.tar
│   ├── venv/                       # Python virtualenv (HA Core)
│   │   └── home-assistant.tar.gz
│   ├── assets/                     # UI static files
│   │   └── ui.tar.gz
│   └── scripts/
│       ├── pre-install.sh
│       ├── post-install.sh
│       └── healthcheck.sh
└── updates/
    └── v1.0.0-to-v1.1.0.delta       # Optional delta patch
```

### Manifest Schema

```json
{
  "schema_version": "1.0.0",
  "bundle_id": "com.uhome.nest",
  "version": "1.0.0",
  "release_channel": "stable",
  "build_date": "2026-04-16T10:00:00Z",
  "architecture": ["amd64", "arm64"],
  "os_release": ["ubuntu-22.04", "ubuntu-24.04", "debian-12"],
  "dependencies": {
    "docker": ">=24.0.0",
    "systemd": ">=249",
    "python": ">=3.10"
  },
  "components": [
    {
      "name": "uhome-api",
      "type": "binary",
      "source": "payload/base/usr/local/bin/uhome-api",
      "checksum": "sha256:abc123...",
      "size_bytes": 15728640
    },
    {
      "name": "jellyfin",
      "type": "docker",
      "source": "payload/docker/jellyfin.tar",
      "image": "jellyfin/jellyfin:10.9.0",
      "checksum": "sha256:def456..."
    },
    {
      "name": "matter",
      "type": "docker",
      "source": "payload/docker/matter.tar",
      "image": "connectedhomeip/chip-tool:v1.3",
      "checksum": "sha256:ghi789..."
    },
    {
      "name": "home-assistant",
      "type": "venv",
      "source": "payload/venv/home-assistant.tar.gz",
      "python_version": "3.11",
      "checksum": "sha256:jkl012..."
    }
  ],
  "scripts": {
    "pre_install": "payload/scripts/pre-install.sh",
    "post_install": "payload/scripts/post-install.sh",
    "healthcheck": "payload/scripts/healthcheck.sh"
  },
  "signatures": [
    {
      "key_id": "uhome-release-2026",
      "algorithm": "ed25519",
      "signature": "base64encoded..."
    }
  ],
  "update_info": {
    "previous_version": "0.9.0",
    "delta_available": false,
    "changelog_url": "https://uhome.local/changelog"
  }
}
```

---

## Packager CLI (`sonic-home pack`)

### Command Interface

```bash
# Build bundle from current uHomeNest checkout
sonic-home pack \
  --source /opt/uhome-nest \
  --output /dist/uhome-nest-v1.0.0.she \
  --version 1.0.0 \
  --channel stable \
  --arch amd64,arm64

# Build with custom components
sonic-home pack \
  --source . \
  --include-components uhome-api,jellyfin,matter \
  --exclude-components home-assistant \
  --output ./uhome-nest-lite.she

# Create delta update between versions
sonic-home pack \
  --delta-from v1.0.0 \
  --delta-to v1.1.0 \
  --output v1.0.0-to-v1.1.0.delta

# Sign existing bundle
sonic-home sign \
  --bundle uhome-nest-v1.0.0.she \
  --key /opt/uhome/keys/release-key.priv \
  --output uhome-nest-v1.0.0.signed.she
```

### Packager Implementation (Go)

```go
// cmd/pack/main.go
package main

import (
    "archive/tar"
    "compress/gzip"
    "crypto/ed25519"
    "encoding/json"
    "io"
    "os"
    "path/filepath"
)

type Packager struct {
    SourceDir   string
    OutputPath  string
    Version     string
    Channel     string
    Architectures []string
}

func (p *Packager) Build() error {
    // 1. Scan source directory for components
    // 2. Generate manifest with checksums
    // 3. Build Docker images if needed (docker save)
    // 4. Create Python venv tarball
    // 5. Write header + manifest + payload
    // 6. Sign with Ed25519 key
    return nil
}

func (p *Packager) CreateDelta(oldVersion, newVersion string) error {
    // 1. Load old and new bundles
    // 2. Compute binary diff (bsdiff/xdelta)
    // 3. Generate delta manifest
    // 4. Output delta bundle
    return nil
}
```

---

## Installer Handler (`sonic-home install`)

### Command Interface

```bash
# Install from local bundle file
sonic-home install /path/to/uhome-nest-v1.0.0.she

# Install from USB (auto-detect)
sonic-home install --usb

# Install from network channel
sonic-home install --channel stable --from https://updates.uhome.local

# Install with preseed answers (unattended)
sonic-home install --preseed install-config.yaml

# Verify bundle before install
sonic-home verify /path/to/bundle.she

# Dry run (what would be installed)
sonic-home install --dry-run /path/to/bundle.she
```

### Installer Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    sonic-home install               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Verify bundle signature (if present)                    │
│  2. Check system requirements (Docker, systemd, space)      │
│  3. Run pre-install.sh (backup existing config)             │
│  4. Extract payload to /opt/uhome-nest/                     │
│  5. Load Docker images (docker load < jellyfin.tar)         │
│  6. Extract venv to /opt/uhome-integrations/venv/           │
│  7. Install systemd units                                   │
│  8. Run post-install.sh (enable services)                   │
│  9. Run healthcheck                                         │
│  10. Register installation (write /etc/uhome/install.json)  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Preseed Configuration

```yaml
# install-config.yaml
uhome:
  version: 1.0.0
  install_path: /opt/uhome-nest
  media_vault_path: /home/uhome/media
  integrations:
    matter: true
    home-assistant: true
  network:
    api_port: 7890
    matter_port: 5540
    ha_port: 8123
  auto_start: true
  backup_existing: true
  backup_path: /opt/uhome-backups/preinstall
```

### Installer API (for remote orchestration)

```go
// POST /api/sonic/install
// Accepts bundle URL or upload

type InstallRequest struct {
    BundleURL   string            `json:"bundle_url,omitempty"`
    BundleData  []byte            `json:"bundle_data,omitempty"`
    Preseed     map[string]interface{} `json:"preseed"`
    DryRun      bool              `json:"dry_run"`
}

type InstallResponse struct {
    Status      string   `json:"status"` // started, complete, failed
    InstallID   string   `json:"install_id"`
    LogURL      string   `json:"log_url,omitempty"`
    Components  []string `json:"components_installed"`
    HealthCheck string   `json:"health_check"`
}
```

---

## Distribution Server (`sonic-home serve`)

### Local Update Channel

```bash
# Start distribution server on LAN
sonic-home serve \
  --port 8080 \
  --bundles /opt/uhome-dist/bundles \
  --channels stable,beta,edge

# Serve with metadata
sonic-home serve \
  --manifest /opt/uhome-dist/channels.json \
  --tls \
  --cert /etc/uhome/ssl/cert.pem
```

### Channel Manifest

```json
{
  "server_version": "1.0.0",
  "channels": {
    "stable": {
      "current": "1.0.0",
      "bundle_url": "https://updates.uhome.local/bundles/uhome-nest-v1.0.0.she",
      "signature_url": "https://updates.uhome.local/bundles/uhome-nest-v1.0.0.sig",
      "manifest_url": "https://updates.uhome.local/manifests/v1.0.0.json",
      "release_date": "2026-04-16",
      "changelog": "Initial stable release"
    },
    "beta": {
      "current": "1.1.0-beta.1",
      "bundle_url": "https://updates.uhome.local/bundles/uhome-nest-v1.1.0-beta.1.she",
      "release_date": "2026-04-20"
    },
    "edge": {
      "current": "1.1.0-dev",
      "bundle_url": "https://updates.uhome.local/bundles/uhome-nest-latest-edge.she"
    }
  },
  "update_check_interval_hours": 24,
  "public_keys": [
    {
      "key_id": "uhome-release-2026",
      "key_data": "base64ed25519publickey..."
    }
  ]
}
```

---

## USB Auto-Install

### USB Layout

```
USB Drive (label: UHOME_INSTALL)
├── sonic-home/          # Installer itself
│   └── sonic-home.bin   # Static binary
├── bundles/
│   └── uhome-nest-v1.0.0.she    # Bundle to install
├── preseed.yaml                 # Optional auto-answers
├── auto-install.sh              # Detects and runs
└── README.txt                   # User instructions
```

### Auto-Detection Script

```bash
#!/bin/bash
# auto-install.sh - runs from USB on insertion

# Detect if we're on a uHomeNest-capable system
if [ -f /etc/os-release ] && grep -q "Ubuntu\|Debian" /etc/os-release; then
    echo "🔍 uHomeNest installation media detected"
    
    # Check if already installed
    if systemctl is-active --quiet uhome-api; then
        echo "⚠️  uHomeNest already installed. Use --force to reinstall."
        exit 0
    fi
    
    # Run installer
    /media/$(whoami)/UHOME_INSTALL/sonic-home/sonic-home.bin \
        install /media/$(whoami)/UHOME_INSTALL/bundles/uhome-nest-v1.0.0.she \
        --preseed /media/$(whoami)/UHOME_INSTALL/preseed.yaml \
        --non-interactive
fi
```

### udev Rule (optional)

```bash
# /etc/udev/rules.d/99-uhome-install.rules
ACTION=="add", SUBSYSTEM=="block", ENV{ID_FS_LABEL}=="UHOME_INSTALL", \
    RUN+="/usr/bin/screen -S uhome-install -dm /media/auto-install.sh"
```

---

## Integration with uHomeNest (Future Sonic-Family)

### Upgrade Path

```bash
# Current: sonic-home (lite)
sonic-home install bundle.she

# Future: full Sonic-family
sonic upgrade --from home --to full
```

### Compatibility Layer

```go
// pkg/sonic/compat.go
// Detect if full Sonic is available

type SonicProvider interface {
    IsAvailable() bool
    Install(bundlePath string) error
    CreateUSB(device string) error
    ManageDualBoot(config DualBootConfig) error
}

func GetSonic() SonicProvider {
    if _, err := exec.LookPath("sonic"); err == nil {
        return &FullSonic{}  // sonic-screwdriver available
    }
    return &HomeSonic{} // fallback to lite
}
```

### Future Sonic-Family Integration Points

| Sonic Component | Integration with sonic-home |
|----------------|-------------------------------------|
| `sonic-screwdriver` | Consumes `.she` bundles, provides Ventoy/USB creation |
| `sonic-studio` | Bundle authoring GUI, manifest editor |
| `sonic-orchestra` | Fleet management, update channels, multi-node |
| `sonic-recovery` | Bundle-based recovery, rollback |

---

## Binder: `#sonic/home`

```markdown
# binder: sonic-home
project: uhome-nest
milestone: v1.0-sonic
objective: lite packager, installer, distribution helper
status: active
priority: high

## Tasks

- [ ] Design .she bundle format spec
- [ ] Implement packager CLI (Go)
- [ ] Bundle header + manifest generation
- [ ] Ed25519 signing/verification
- [ ] Docker image extraction (docker save/load)
- [ ] Python venv tarball generation
- [ ] Installer handler with preseed support
- [ ] USB auto-detection and installation
- [ ] Local distribution server (HTTP)
- [ ] Update channel manifest
- [ ] Delta update generation (bsdiff)
- [ ] Integration with uhome-api (status endpoints)
- [ ] Documentation (BUNDLE-FORMAT.md, INSTALLER-API.md)
- [ ] Upgrade guide to full Sonic-family

## Dependencies

- Go 1.21+
- Docker (for bundle creation, optional at runtime)
- bsdiff (for delta generation)

## Outputs

- Binary: `sonic-home` (static, ~15MB)
- Bundle format: `.she` (Sonic Home)
- Default bundles for uHomeNest v1.0.0

## Testing

- [ ] Create bundle from clean uHomeNest checkout
- [ ] Install bundle on fresh Ubuntu 22.04 VM
- [ ] Verify all services start correctly
- [ ] USB auto-install on Raspberry Pi 4
- [ ] Delta update v1.0.0 → v1.1.0
- [ ] Upgrade to mock Sonic-family (future compatibility)
```

---

## Installation (for uHomeNest Operators)

```bash
# Install sonic-home as optional module
cd /opt/uhome-nest
./scripts/install-sonic-express.sh

# Or download pre-built binary
curl -L https://github.com/uhome-project/sonic-home/releases/download/v1.0.0/sonic-home \
    -o /usr/local/bin/sonic-home
chmod +x /usr/local/bin/sonic-home

# Test
sonic-home version
# sonic-home v1.0.0 (build 2026-04-16, ed25519)

# Create your first bundle
sonic-home pack --source . --output uhome-backup.she
```

---

## Non-Goals (v1.0)

- Full Ventoy/USB creation (future: sonic-screwdriver)
- Dual-boot management (future: sonic-screwdriver)
- Fleet orchestration (future: sonic-orchestra)
- GUI installer (CLI only for v1)
- Windows installer (Linux only)
- Encrypted bundles (v2 feature)
- P2P distribution (v2 feature)

---

## Success Criteria

- [ ] `sonic-home pack` creates valid `.she` bundle from source
- [ ] Bundle size < 2GB for full uHomeNest (with Docker images)
- [ ] `sonic-home install` completes on clean Ubuntu in < 10 minutes
- [ ] USB auto-install works without network (offline install)
- [ ] Update channel serves bundles, client can check/install
- [ ] Bundle signature prevents tampering
- [ ] Healthcheck passes after installation
- [ ] Upgrade to full Sonic-family possible without reinstall

---

## Related Documents

- [uHomeNest v1.0.0 Dev Brief](./UHOMENEST-V1-DEV-BRIEF.md)
- [Matter + HA Integration Plan](./UDN-INTEGRATION-001.md)
- [Sonic-Family Architecture (future)](https://github.com/fredporter/sonic-screwdriver/docs/ARCHITECTURE.md)
- [Bundle Format Specification](./docs/BUNDLE-FORMAT.md)

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-04-16 | Initial sonic-home brief |
