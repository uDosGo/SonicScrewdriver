# Sonic-Screwdriver v2.1.0

## TARDIS Console: API Central Hub for Smart Home Automation

Sonic-Screwdriver is a modular platform for managing secrets, APIs, containers, and smart home integrations with a focus on security and extensibility.

*TARDIS Narrator applied: LIGHT intensity* - Themed introduction, technical content unchanged.

## 📖 Documentation

### Getting Started
- **[QUICKSTART.md](QUICKSTART.md)** - Complete setup and usage guide
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - System architecture overview

### Core Features
- **[docs/CLI_COMMANDS.md](docs/CLI_COMMANDS.md)** - Comprehensive CLI reference
- **[docs/LIBRARY_FORMAT.md](docs/LIBRARY_FORMAT.md)** - Game library format specification
- **[docs/SECRET_ROTATION_GUIDE.md](docs/SECRET_ROTATION_GUIDE.md)** - Secret management guide

### Integrations
- **Home Assistant**: Built-in integration with iframe embed support
- **Media Player**: Local media management system
- **Feed/Spool**: Content aggregation and processing
- **Remote Access**: VNC, SSH, and Samba support

### Development
- **[docs/DEVLOG.md](docs/DEVLOG.md)** - Current development status
- **[docs/ROADMAP.md](docs/ROADMAP.md)** - Future plans and roadmap
- **[CHANGELOG.md](CHANGELOG.md)** - Version history

## 🎯 Key Features

### v2.1.0 Highlights
- ✅ Home Assistant deep integration
- ✅ Iframe embed strategy with kiosk mode
- ✅ Enhanced secret rotation with history
- ✅ Comprehensive CLI command set
- ✅ Complete documentation overhaul

### Core Capabilities
- **Secret Store**: AES-256-GCM encrypted secret management
- **API Proxy**: Secure proxy with rate limiting
- **Node Registry**: Distributed node management
- **Container Runtime**: Docker-based game management
- **TUI Interface**: Interactive terminal interface
- **CLI Commands**: Comprehensive command-line tools

## 🚀 Quick Start

```bash
# Build from source
go build -o sonic ./cmd/sonic

# Or use the installer
./install.sh

# Check system
sonic system check

# View help
sonic --help
```

## Narrator Reference

This project uses TARDIS Narrator (LIGHT intensity) for themed documentation. Technical instructions remain unchanged for clarity.

| Themed Term | Meaning |
|-------------|---------|
| TARDIS Console | Sonic-Screwdriver (the main application) |

*Designed by Wizard. Built by Vibe.*
