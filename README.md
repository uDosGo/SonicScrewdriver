# Sonic-Screwdriver v1.2.0

## Unified System Toolkit

Sonic-Screwdriver is a Go CLI toolkit for system administration — USB installation, Docker container management, secret storage, device cataloguing, and knowledge source querying.

## 🎯 What It Does

```
sonic container <list|start|stop|restart|remove|health>  — Docker container management
sonic usb <list|prepare|install|full-install>             — USB installer (ISO download + write)
sonic vault <get|set|list|rotate|history>                 — Encrypted secret store (via uServer)
sonic gui                                                  — Web GUI dashboard
sonic catalogue <list|find>                               — Device catalogue (scans Vaults + uCode repos)
sonic knowledge <sources|query>                           — Knowledge source querying
sonic library <list|info|validate>                        — Game library index management
sonic ventoy <create|validate|info>                       — .she bundle packaging
sonic remote <vnc|ssh|samba>                              — Remote access setup
sonic mint <check|install|apply|status|info|doctor>       — Classic Modern Mint readiness
```

## 🚀 Quick Start

```bash
# Build from source
go build -o sonic ./cmd/sonic

# View help
sonic --help
```

## 🏗️ Project Structure

```
SonicScrewdriver/
├── cmd/sonic/              # CLI entrypoint
├── pkg/
│   ├── container/          # Docker runtime wrapper
│   ├── vault/              # Secret store (wraps uServer/pkg/secrets)
│   ├── disk/               # Block device & partition management
│   ├── iso/                # ISO downloader & writer
│   ├── usb/                # USB installer
│   ├── gui/                # Web GUI (embedded HTML/JS)
│   ├── catalogue/          # Device catalogue
│   ├── knowledge/          # Knowledge source querying
│   ├── library/            # Game library index & manifest validation
│   ├── ventoy/             # .she bundle packager
│   ├── remote/             # VNC/SSH/Samba setup
│   └── classicmodern/      # Classic Modern Mint readiness
├── docs/                   # Documentation
├── test/                   # Integration tests
└── version                 # v1.1.0
```

## 🔗 Dependencies

- **uServer** (`github.com/uDosGo/uServer/pkg/secrets`) — AES-256-GCM encrypted SQLite secret store
- **Docker** — Container runtime (optional, for container commands)

## 📖 Documentation

- **[docs/legacy/](docs/legacy/)** — Archived documentation from earlier aspirational scope

## Related Repositories

- **uServer** — Backend services, secret store, API central
- **DevStudio** — Development environment configuration and tooling
