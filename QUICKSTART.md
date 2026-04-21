# Sonic-Screwdriver Quick Start Guide

## 🚀 Getting Started

### Prerequisites
- Go 1.25+ installed
- Git
- Docker (for container operations)
- Basic development tools (make, curl, etc.)

### Installation

```bash
# Clone the repository
cd /home/wizard/code-vault
git clone https://github.com/fredporter/sonic-screwdriver.git
cd sonic-screwdriver

# Build the application
go build -o sonic ./cmd/sonic

# Install to system (optional)
sudo cp sonic /usr/local/bin/
```

### Basic Usage

```bash
# Show help
sonic --help

# Launch interactive TUI
sonic tui

# List available commands
sonic menu
```

## 📦 Key Components

### 1. API Central Hub
```bash
# Add a secret
sonic secret add my_api_key --value "your-key-here"

# Get a secret
sonic secret get my_api_key

# Rotate a secret
sonic secret rotate my_api_key --value "new-key-here"

# List all secrets
sonic secret list
```

### 2. Home Assistant Integration (v2.1.0)
```bash
# Setup HA integration
sonic ha setup "http://your-ha-url:8123" "your-long-lived-token"

# Generate embed HTML
sonic ha embed /path/to/output.html

# Check HA status
sonic ha status

# Enable kiosk mode
sonic ha kiosk enable
```

### 3. Node Management
```bash
# Register a node
sonic node register --master "master-address" --name "node-name"

# List registered nodes
sonic node list

# Grant secret access
sonic secret grant my_api_key --node "node-name"
```

### 4. Remote Access
```bash
# Setup VNC server
sonic remote vnc setup "password" "1920x1080"
sonic remote vnc start

# Setup SSH access
sonic remote ssh setup

# Setup Samba sharing
sonic remote samba setup "share-name" "/path/to/share"
```

## 🎮 Game Management

```bash
# Install a game from library
sonic install game-name

# Start a game
sonic start game-name

# Stop a game
sonic stop game-name

# List installed games
sonic list

# Check game health
sonic health game-name
```

## 🔧 Setup Scripts

The repository includes several setup scripts:

- `setup-git.sh` - Configure git settings
- `setup-logging.sh` - Setup logging configuration
- `setup-spine.sh` - Setup spine animation runtime
- `setup-upgrades.sh` - Configure upgrade system
- `setup-vendors.sh` - Setup vendor dependencies
- `install-sonic-udos.sh` - Install Sonic uDos integration
- `add_health_commands.sh` - Add health check commands
- `seed-data.sh` - Seed initial data

Run them as needed:
```bash
./setup-git.sh
./setup-vendors.sh
# etc.
```

## 📚 Documentation

- **Full Documentation**: See `docs/` directory
- **Secret Rotation Guide**: `docs/SECRET_ROTATION_GUIDE.md`
- **Media Catalog Schema**: `docs/MEDIA_CATALOG_SCHEMA.md`
- **Feed/Spool Spec**: `docs/FEED_SPOOL_SPEC.md`

## 🐳 Docker Operations

```bash
# Start a container
sonic start game-name

# Stop a container
sonic stop game-name

# Check container health
sonic health game-name

# Repair containers
sonic repair --all
```

## 🔐 Security Features

```bash
# Backup secrets
sonic secret backup secrets-backup.json

# Restore secrets
sonic secret restore secrets-backup.json

# Export encrypted backup
sonic secret export encrypted-backup.enc

# Import encrypted backup
sonic secret import encrypted-backup.enc
```

## 🎨 Classic Modern Mint

```bash
# Check Classic Modern readiness
sonic mint check

# Install Classic Modern theme
sonic mint install

# Apply Classic Modern theme
sonic mint apply

# Run diagnostic checks
sonic mint doctor
```

## 🌐 Ventoy Integration

```bash
# Create Ventoy bundle
sonic ventoy package /path/to/installer /output/path/bundle.she

# Validate bundle
sonic ventoy validate bundle.she

# Show bundle info
sonic ventoy info bundle.she
```

## 📊 Version Information

```bash
# Show version
sonic --version
# or
sonic -v
```

## 🆘 Troubleshooting

```bash
# Check container logs
sonic logs game-name

# Repair unhealthy containers
sonic repair game-name

# Run health checks
sonic health --all
```

## 📈 Current Version

**Sonic-Screwdriver v2.1.0** - API Central Hub with Home Assistant Integration

### What's New in v2.1.0
- Home Assistant deep integration
- Iframe embed strategy with kiosk mode
- Media player foundation
- Feed/spool system design
- Enhanced secret rotation with history

## 🔗 Related Repositories

- **uDev Framework**: Core framework
- **uDos Connect**: Connectivity layer
- **uHome Nest**: Home automation surface

## 📝 License

MIT License - See LICENSE file for details.

## 🤝 Support

For issues, questions, or contributions:
- Check the GitHub issues page
- Review the documentation
- Contact the development team

---

*Last Updated: 2026-04-29*
*Sonic-Screwdriver v2.1.0*