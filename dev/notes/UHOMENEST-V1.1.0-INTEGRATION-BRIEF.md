# uHomeNest — Matter + Home Assistant Integration Plan (uDOS Style)

**Document ID:** `UDN-INTEGRATION-001`  
**Status:** Active  
**Version:** 1.0.0  
**Date:** 2026-04-16  
**Related:** [uHomeNest v1.0.0 Dev Brief](./UHOMENEST-V1-DEV-BRIEF.md)

---

## Objective

Integrate **Matter** and **Home Assistant** into uHomeNest v1.0.0 using **uDOS-style clone + run** patterns — treating each integration as a **cloned component** managed by uHomeNest, not as external dependencies or embedded submodules.

**Principle:** uHomeNest owns the **orchestration**; Matter/HA repos own their **upstream code**.

---

## Integration Architecture

### Topology

```
┌─────────────────────────────────────────────────────────────┐
│                      uHomeNest Host                         │
│  (Ubuntu 22.04/24.04, ~/media/ vault, Jellyfin)            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐     ┌─────────────────────────────┐   │
│  │  uhome-api      │────▶│  uhome-integration-manager  │   │
│  │  (Go)           │     │  (systemd service)          │   │
│  └─────────────────┘     └──────────────┬──────────────┘   │
│                                          │                  │
│                    ┌─────────────────────┼──────────────┐   │
│                    ▼                     ▼              ▼   │
│         ┌──────────────────┐  ┌──────────────────┐  ┌──────┐│
│         │ matter-clone/    │  │ ha-clone/        │  │ ...  ││
│         │ (git clone)      │  │ (git clone)      │  │      ││
│         │ └── chip-tool    │  │ └── core/        │  │      ││
│         │ └── matter-server│  │ └── supervisor/  │  │      ││
│         │ └── ...          │  │ └── ...          │  │      ││
│         └────────┬─────────┘  └────────┬─────────┘  └──────┘│
│                  │                     │                    │
│                  ▼                     ▼                    │
│         ┌──────────────────────────────────────────────┐    │
│         │           Matter Fabric / HA Core            │    │
│         │         (running processes, sockets)         │    │
│         └──────────────────────────────────────────────┘    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Integration Manager (uDOS-style Service)

### Service Definition

**`uhome-integration-manager.service`**:

```ini
[Unit]
Description=uHomeNest Integration Manager (Matter + HA)
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
User=uhome
Group=uhome
WorkingDirectory=/opt/uhome-integrations
ExecStart=/opt/uhome-integrations/bin/manager

# Clone management
RuntimeDirectory=uhome-integrations
Environment=UHOME_CLONE_ROOT=/opt/uhome-integrations/clones
Environment=UHOME_CLONE_MANIFEST=/opt/uhome-integrations/manifest.json

# Matter
Environment=MATTER_CLONE_REPO=https://github.com/project-chip/connectedhomeip.git
Environment=MATTER_CLONE_BRANCH=v1.3

# Home Assistant
Environment=HA_CLONE_REPO=https://github.com/home-assistant/core.git
Environment=HA_CLONE_BRANCH=2026.4

# Resource limits
MemoryMax=2G
CPUQuota=150%

Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Manager Binary (Go)

```go
// cmd/manager/main.go
package main

import (
    "encoding/json"
    "os/exec"
    "sync"
)

type CloneManifest struct {
    Components []Component `json:"components"`
    Mu         sync.RWMutex
}

type Component struct {
    Name      string `json:"name"`
    Repo      string `json:"repo"`
    Branch    string `json:"branch"`
    Path      string `json:"path"`
    Status    string `json:"status"` // cloned, pulling, running, stopped, failed
    PID       int    `json:"pid,omitempty"`
    LastSync  string `json:"last_sync,omitempty"`
}

func main() {
    // 1. Load manifest (or create default)
    // 2. For each component: ensure cloned, git pull
    // 3. Start component services (Docker Compose or native)
    // 4. Health check loop
    // 5. Expose API on :7891 for uhome-api to query
}
```

---

## Clone Manifest Schema

**`/opt/uhome-integrations/manifest.json`**:

```json
{
  "version": "1.0.0",
  "clone_root": "/opt/uhome-integrations/clones",
  "components": [
    {
      "name": "matter",
      "repo": "https://github.com/project-chip/connectedhomeip.git",
      "branch": "v1.3",
      "path": "clones/matter",
      "run_mode": "docker",
      "docker_compose": "clones/matter/docker-compose.yml",
      "health_check": "http://localhost:8080/health",
      "required": false,
      "env": {
        "MATTER_FABRIC_ID": "uhome.lan",
        "MATTER_COMMISSIONING_PORT": "5540"
      }
    },
    {
      "name": "home-assistant",
      "repo": "https://github.com/home-assistant/core.git",
      "branch": "2026.4",
      "path": "clones/home-assistant",
      "run_mode": "venv",
      "start_script": "clones/home-assistant/venv/bin/hass",
      "health_check": "http://localhost:8123/api/health",
      "required": false,
      "env": {
        "HA_CONFIG_DIR": "/opt/uhome-integrations/config/home-assistant"
      }
    },
    {
      "name": "matter-server",
      "repo": "https://github.com/home-assistant-libs/python-matter-server.git",
      "branch": "main",
      "path": "clones/matter-server",
      "run_mode": "docker",
      "docker_compose": "clones/matter-server/docker-compose.yml",
      "health_check": "http://localhost:5588/health",
      "required": false,
      "depends_on": ["matter"]
    }
  ]
}
```

---

## uDOS-style Clone Workflow

### Binder: `#integration/matter-ha`

```markdown
# binder: uhome-integration-matter-ha
project: uhome-nest
milestone: v1.0-integration
objective: clone, run, and orchestrate Matter + HA
status: active
priority: high

## Tasks

- [x] Design integration manager service
- [x] Create clone manifest schema
- [ ] Implement manager binary (Go)
- [ ] Add git clone/pull with retry logic
- [ ] Add Docker Compose runner
- [ ] Add venv runner for HA core
- [ ] Add health check polling
- [ ] Add API endpoint for uhome-api
- [ ] Write systemd unit
- [ ] Test on clean Ubuntu 22.04
- [ ] Document in ARCHITECTURE.md

## Clone Sources

- Matter: https://github.com/project-chip/connectedhomeip.git
- HA Core: https://github.com/home-assistant/core.git
- Python Matter Server: https://github.com/home-assistant-libs/python-matter-server.git

## Integration Points

- uhome-api → manager API (port 7891)
- manager → clone management
- clones → running services
- running services → uHomeNest browser UI (USXD tiles)

## Promotion Criteria

- `./scripts/install-integrations.sh` works idempotently
- `systemctl status uhome-integration-manager` shows active
- Matter and HA can be started/stopped via API
- Browser UI shows Matter/HA status tiles
```

---

## Clone + Run Commands (Reference)

### Manual Clone (for testing)

```bash
# As uhome user
mkdir -p /opt/uhome-integrations/clones
cd /opt/uhome-integrations/clones

# Clone Matter
git clone --depth 1 --branch v1.3 https://github.com/project-chip/connectedhomeip.git matter

# Clone HA Core
git clone --depth 1 --branch 2026.4 https://github.com/home-assistant/core.git home-assistant

# Clone Python Matter Server
git clone --depth 1 https://github.com/home-assistant-libs/python-matter-server.git matter-server
```

### Run Matter (Docker)

```bash
cd /opt/uhome-integrations/clones/matter
docker compose -f docker-compose.yml up -d
```

### Run HA Core (venv)

```bash
cd /opt/uhome-integrations/clones/home-assistant
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
hass --config /opt/uhome-integrations/config/home-assistant
```

---

## Integration API (Manager → uhome-api)

**Manager listens on `localhost:7891`**:

```go
// GET /api/integrations/status
{
  "matter": {
    "status": "running",
    "pid": 1234,
    "version": "v1.3",
    "health": "ok",
    "last_sync": "2026-04-16T10:00:00Z"
  },
  "home-assistant": {
    "status": "running",
    "pid": 5678,
    "version": "2026.4",
    "health": "ok"
  },
  "matter-server": {
    "status": "running",
    "pid": 9012,
    "health": "ok"
  }
}

// POST /api/integrations/{name}/start
// POST /api/integrations/{name}/stop
// POST /api/integrations/sync-all (git pull all clones)
```

**uhome-api consumes this** for UI status tiles:

```go
// In uhome-api handlers
func getIntegrationStatus(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get("http://localhost:7891/api/integrations/status")
    // forward to browser UI
}
```

---

## Browser UI Integration (USXD Tiles)

**Extended launcher surface** (`ui/usxd/launcher.json`):

```json
{
  "tiles": [
    {"id": "media", "label": "Media Library", "action": "navigate:/media"},
    {"id": "now-playing", "label": "Now Playing", "action": "navigate:/now-playing"},
    {
      "id": "matter",
      "label": "Matter Devices",
      "icon": "chip",
      "status_endpoint": "/api/integrations/status/matter",
      "action": "navigate:/matter/devices"
    },
    {
      "id": "home-assistant",
      "label": "Home Assistant",
      "icon": "home",
      "status_endpoint": "/api/integrations/status/home-assistant",
      "action": "navigate:/home-assistant/dashboard"
    },
    {"id": "settings", "label": "Settings", "action": "navigate:/settings"}
  ]
}
```

**Status badges** show green/yellow/red based on integration health.

---

## Installation Script Additions

**`scripts/install-integrations.sh`**:

```bash
#!/bin/bash
# uHomeNest Integration Installer (Matter + HA)

set -e

UHOME_USER=${SUDO_USER:-$USER}
CLONE_ROOT="/opt/uhome-integrations"

echo "📦 Installing uHomeNest Integrations (Matter + HA)"

# Create directories
sudo mkdir -p $CLONE_ROOT/{clones,config,bin}
sudo chown -R $UHOME_USER:$UHOME_USER $CLONE_ROOT

# Install integration manager binary
sudo cp ./bin/integration-manager $CLONE_ROOT/bin/
sudo chmod +x $CLONE_ROOT/bin/integration-manager

# Create default manifest if not exists
if [ ! -f "$CLONE_ROOT/manifest.json" ]; then
    sudo cp ./config/manifest.default.json $CLONE_ROOT/manifest.json
    sudo chown $UHOME_USER:$UHOME_USER $CLONE_ROOT/manifest.json
fi

# Install systemd service
sudo cp ./systemd/uhome-integration-manager.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable uhome-integration-manager
sudo systemctl start uhome-integration-manager

# Verify
sleep 3
if systemctl is-active --quiet uhome-integration-manager; then
    echo "✅ Integration manager running"
else
    echo "❌ Integration manager failed to start"
    journalctl -u uhome-integration-manager -n 20
    exit 1
fi

# Initial clone of all components
echo "🔄 Cloning integration components (first sync)..."
curl -X POST http://localhost:7891/api/integrations/sync-all

echo "✅ Integrations installed"
echo "   Manager API: http://localhost:7891"
echo "   Matter:      http://localhost:8080 (if running)"
echo "   HA Core:     http://localhost:8123 (if running)"
```

---

## Manifest Default Configuration

**`config/manifest.default.json`**:

```json
{
  "version": "1.0.0",
  "clone_root": "/opt/uhome-integrations/clones",
  "auto_start": ["matter-server", "home-assistant"],
  "components": [
    {
      "name": "matter",
      "enabled": true,
      "repo": "https://github.com/project-chip/connectedhomeip.git",
      "branch": "v1.3",
      "run_mode": "docker",
      "compose_file": "docker-compose.yml",
      "ports": [8080, 5540]
    },
    {
      "name": "home-assistant",
      "enabled": true,
      "repo": "https://github.com/home-assistant/core.git",
      "branch": "2026.4",
      "run_mode": "venv",
      "config_dir": "/opt/uhome-integrations/config/home-assistant",
      "ports": [8123]
    },
    {
      "name": "matter-server",
      "enabled": true,
      "repo": "https://github.com/home-assistant-libs/python-matter-server.git",
      "branch": "main",
      "run_mode": "docker",
      "compose_file": "docker-compose.yml",
      "ports": [5588],
      "depends_on": ["matter"]
    }
  ]
}
```

---

## uHomeNest v1.0.0 Integration Boundary

| Layer | Owns | Does Not Own |
|-------|------|--------------|
| **uhome-api** | Integration status API, UI tiles, start/stop commands | Clone management, git operations, process supervision |
| **integration-manager** | Clone lifecycle, git sync, process start/stop, health checks | Media scanning, Jellyfin, ~/media/ vault |
| **Matter/HA repos** | Their own code, upstream updates | uHomeNest integration logic |

---

## Integration Workflow (Operator View)

```bash
# Fresh install
./scripts/install.sh                    # Base uHomeNest
./scripts/install-integrations.sh       # Matter + HA clones

# Check status
curl http://localhost:7891/api/integrations/status

# Start/stop integrations
curl -X POST http://localhost:7891/api/integrations/matter-server/start
curl -X POST http://localhost:7891/api/integrations/home-assistant/stop

# Sync all clones (git pull)
curl -X POST http://localhost:7891/api/integrations/sync-all

# Browser UI
# http://localhost:7890/ → shows Matter/HA status tiles
```

---

## Testing Matrix

| Scenario | Expected |
|----------|----------|
| Fresh Ubuntu, no Docker | Manager installs Docker, clones, runs |
| Git clone fails | Retry 3x, logs error, continues other components |
| HA venv fails | Fallback to Docker mode (if available) |
| Matter port conflict | Manager detects, tries next port |
| Integration crashes | Manager restarts (configurable backoff) |
| Upstream repo updates | `sync-all` does `git pull --rebase` |

---

## Related Binders

- `#integration/clone-manager` — Manager service implementation
- `#integration/matter-runtime` — Matter Docker/venv runner
- `#integration/ha-runtime` — HA Core runner
- `#integration/ui-tiles` — USXD status integration

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-04-16 | Initial uDOS-style integration plan |
