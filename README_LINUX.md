# Sonic-Screwdriver v2.1.0 for Linux Ubuntu

## 🚀 API Central Hub for Smart Home Automation on Ubuntu

Sonic-Screwdriver is a modular platform for managing secrets, APIs, containers, and smart home integrations with a focus on security and extensibility, optimized for Linux Ubuntu systems.

## 📖 Linux Ubuntu Documentation

### Getting Started
- **[QUICKSTART_LINUX.md](QUICKSTART_LINUX.md)** - Complete Ubuntu setup and usage guide
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

## 🎯 Key Features for Ubuntu

### v2.1.0 Highlights
- ✅ Home Assistant deep integration
- ✅ Iframe embed strategy with kiosk mode
- ✅ Enhanced secret rotation with history
- ✅ Comprehensive CLI command set
- ✅ Complete documentation overhaul
- ✅ **Linux Ubuntu specific optimizations**
- ✅ **System ID checks at startup**
- ✅ **Ubuntu 22.04 LTS+ compatibility**

### Core Capabilities
- **Secret Store**: AES-256-GCM encrypted secret management
- **API Proxy**: Secure proxy with rate limiting
- **Node Registry**: Distributed node management
- **Container Runtime**: Docker-based game management
- **TUI Interface**: Interactive terminal interface
- **CLI Commands**: Comprehensive command-line tools
- **System Monitoring**: Ubuntu-specific system checks

## 🚀 Quick Start for Ubuntu

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install required dependencies
sudo apt install -y git make curl docker.io golang-go g++ libssl-dev

# Add user to docker group (log out and back in after this)
sudo usermod -aG docker $USER

# Clone the repository
cd ~/Code
git clone https://github.com/uDosGo/SonicScrewdriver.git
cd SonicScrewdriver

# Build the application
make build

# Install to system
sudo make install

# Launch TUI
sonic tui

# Get help
sonic --help
```

See **[QUICKSTART_LINUX.md](QUICKSTART_LINUX.md)** for detailed Ubuntu setup instructions.

## 📦 Installation for Ubuntu

### Prerequisites
- **Ubuntu 22.04 LTS or later**
- Go 1.25+
- Docker 20.10+
- Git
- Make
- G++ and development tools

### Build

```bash
# Clone repository
git clone https://github.com/uDosGo/SonicScrewdriver.git
cd SonicScrewdriver

# Build from source
make build
```

### Install

```bash
# Install to system
sudo make install

# Verify installation
sonic --version

# Check system compatibility
sonic system check
```

## 🎮 Usage Examples for Ubuntu

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
sonic secret rotate openrouter_api_key --value "sk-new-..."

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

## 🔌 Ubuntu-Specific Features

### System ID Checks
Sonic-Screwdriver performs comprehensive system checks at startup:

```bash
# Check system compatibility
sonic system check

# Show system information
sonic system info

# Monitor system resources
sonic system resources
```

### Docker Integration

```bash
# Check container usage
sonic system containers

# Monitor disk space
sonic system disk

# View system logs
sonic system logs
```

### Ubuntu Configuration

```bash
# Setup Docker
sudo systemctl start docker
sudo systemctl enable docker

# Configure firewall
sudo ufw allow 8080/tcp
sudo ufw allow 8123/tcp
sudo ufw enable

# Set up Go environment
 echo 'export GOPATH=$HOME/go' >> ~/.bashrc
 echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
 source ~/.bashrc
```

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

## 🔧 Development on Ubuntu

### Setup

```bash
# Clone repository
git clone https://github.com/uDosGo/SonicScrewdriver.git
cd SonicScrewdriver

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
- Optimize for Ubuntu compatibility

## 📊 Version History

### v2.1.0 (2024-04-28)
- Home Assistant integration
- Enhanced secret management
- Documentation overhaul
- CLI improvements
- **Linux Ubuntu specific optimizations**
- **System ID checks at startup**
- **Ubuntu 22.04 LTS+ compatibility**

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

## 🐧 Ubuntu-Specific Notes

### System Requirements
- **Ubuntu 22.04 LTS or later** recommended
- **x86_64/amd64** architecture
- **4GB RAM minimum** (8GB recommended)
- **10GB free disk space**

### System ID Checks
Sonic-Screwdriver performs the following checks at startup:

1. **OS Version**: Verifies Ubuntu 22.04 LTS or later
2. **Architecture**: Confirms x86_64/amd64 compatibility
3. **Dependencies**: Checks for required tools (git, make, curl, g++, docker)
4. **Docker**: Verifies Docker installation and daemon status
5. **Go Version**: Confirms Go 1.25+ compatibility

### Troubleshooting

```bash
# Check system compatibility
sonic system check

# View system logs
sonic system logs

# Check container health
sonic health --all

# Repair unhealthy containers
sonic repair --all
```

### Performance Optimization

```bash
# Monitor system resources
sonic system resources

# Check disk space
sonic system disk

# View container usage
sonic system containers
```

---

*Sonic-Screwdriver v2.1.0 for Linux Ubuntu*  
*API Central Hub for Smart Home Automation*  
*Optimized for Ubuntu 22.04 LTS and later*  
*Built with ❤️ using Go and Classic Modern Mint principles*
