# Sonic-Screwdriver v2.1.0

## 🚀 API Central Hub for Smart Home Automation

Sonic-Screwdriver is a modular platform for managing secrets, APIs, containers, and smart home integrations with a focus on security and extensibility.

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

# Install
sudo cp sonic /usr/local/bin/

# Launch TUI
sonic tui

# Get help
sonic --help
```

See **[QUICKSTART.md](QUICKSTART.md)** for detailed setup instructions.

## 📦 Installation

### Prerequisites
- Go 1.25+
- Docker 20.10+
- Git
- Make

### Build

```bash
git clone https://github.com/fredporter/sonic-screwdriver.git
cd sonic-screwdriver
go build -o sonic ./cmd/sonic
```

### Install

```bash
sudo cp sonic /usr/local/bin/
sonic --version
```

## 🎮 Usage Examples

### Home Assistant Integration

```bash
# Setup HA integration
sonic ha setup "http://ha.local:8123" "your-long-lived-token"

# Generate embed HTML
sonic ha embed /var/www/ha-embed.html

# Enable kiosk mode
sonic ha kiosk enable
sonic ha refresh 60
```

### Secret Management

```bash
# Add API key
sonic secret add openrouter_api_key --value "sk-..."

# Rotate secret
sonic secret rotate openrouter_api_key --value "sk-new..."

# View history
sonic secret history openrouter_api_key
```

### Game Management

```bash
# Install and run game
sonic install my-game
sonic start my-game
sonic health my-game
```

## 🔌 Integrations

### Home Assistant
- Iframe embed strategy
- Kiosk mode with auto-refresh
- API connectivity testing
- Configuration management

### Media Player
- Media scanning and indexing
- Metadata extraction
- Library management
- Playback control

### Feed/Spool System
- Feed parsing (RSS, Atom, JSON)
- Content validation
- Spool processing pipeline
- Notification system

## 📁 Project Structure

```
sonic-screwdriver/
├── cmd/                # CLI commands
├── internal/           # Core modules
│   ├── container/      # Docker runtime
│   ├── homeassistant/  # HA integration
│   ├── library/        # Game library
│   ├── remote/         # Remote access
│   ├── secrets/        # Secret store
│   └── state/          # State management
├── docs/               # Documentation
├── dev/                # Development files
├── modules/            # Integration modules
├── scripts/            # Utility scripts
└── library/            # Game library
```

## 🔧 Development

### Setup

```bash
# Clone repository
git clone https://github.com/fredporter/sonic-screwdriver.git
cd sonic-screwdriver

# Install dependencies
go mod download

# Build
make build

# Test
make test
```

### Documentation

```bash
# Update documentation
# Follow patterns in existing docs

# Add new documentation
# Create comprehensive guides

# Compost old documentation
mv old-doc.md dev/compost/old-docs/
```

## 🤝 Contributing

See **[CONTRIBUTING.md](CONTRIBUTING.md)** for contribution guidelines.

### Development Workflow

1. Fork the repository
2. Create feature branch
3. Implement changes
4. Write tests
5. Update documentation
6. Submit pull request

### Code Standards

- Follow Go best practices
- Write comprehensive tests
- Document all public APIs
- Maintain consistent style

## 📊 Version History

### v2.1.0 (2026-04-29)
- Home Assistant integration
- Enhanced secret management
- Documentation overhaul
- CLI improvements

### v2.0.0 (2026-04-22)
- API Central Hub foundation
- Secret store with encryption
- Node registry
- Container runtime
- TUI interface

### v1.1.0 (2026-04-15)
- Runtime foundation
- Library management
- State persistence
- CLI wiring

See **[CHANGELOG.md](CHANGELOG.md)** for complete history.

## 📚 Related Projects

- **uDev Framework**: Core framework
- **uDos Connect**: Connectivity layer
- **uHome Nest**: Home automation surface

## 🔐 License

MIT License - See **[LICENSE](LICENSE)** for details.

## 🤝 Support

- **Issues**: GitHub Issues
- **Discussion**: GitHub Discussions
- **Documentation**: This repository
- **Contact**: project maintainers

---

*Sonic-Screwdriver v2.1.0*
*API Central Hub for Smart Home Automation*
*Built with ❤️ using Go and Classic Modern Mint principles*