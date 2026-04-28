# Sonic-Screwdriver Quick Start Guide for Linux Ubuntu

## 🚀 Getting Started on Ubuntu

### Prerequisites
- **Ubuntu 22.04 LTS or later**
- Go 1.25+ installed
- Git
- Docker (for container operations)
- Basic development tools (make, curl, g++, etc.)

### System Requirements
- **CPU**: 2+ cores (4+ recommended)
- **RAM**: 4GB minimum (8GB recommended)
- **Disk**: 10GB free space
- **Architecture**: x86_64/amd64

### Installation

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
```

### Verify Installation

```bash
# Check sonic version
sonic --version

# Check system compatibility
sonic system check

# Show help
sonic --help
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

The repository includes Ubuntu-specific setup scripts:

```bash
# Setup git configuration
./setup-git.sh

# Setup logging
./setup-logging.sh

# Setup spine animation runtime
./setup-spine.sh

# Setup upgrade system
./setup-upgrades.sh

# Setup vendor dependencies
./setup-vendors.sh

# Install Sonic uDos integration
./install-sonic-udos.sh

# Add health check commands
./add_health_commands.sh

# Seed initial data
./seed-data.sh
```

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

# Show system information
sonic system info

# Show environment details
sonic env
```

## 🆘 Troubleshooting

```bash
# Check container logs
sonic logs game-name

# Repair unhealthy containers
sonic repair game-name

# Run health checks
sonic health --all

# Check system compatibility
sonic system check
```

## 📈 Current Version

**Sonic-Screwdriver v2.1.0** - API Central Hub with Home Assistant Integration

### What's New in v2.1.0
- Home Assistant deep integration
- Iframe embed strategy with kiosk mode
- Media player foundation
- Feed/spool system design
- Enhanced secret rotation with history
- **Linux Ubuntu specific optimizations**
- **System ID checks at startup**

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

## 🐧 Ubuntu-Specific Notes

### System ID Checks
Sonic-Screwdriver performs system ID checks at startup to ensure compatibility:

```bash
# Check system ID
cat /etc/os-release

# Sonic will verify:
# - Ubuntu version (22.04 LTS or later)
# - Architecture (x86_64/amd64)
# - Required dependencies
# - Docker availability
# - Go version compatibility
```

### Docker Configuration

Ensure Docker is properly configured:

```bash
# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Verify Docker is running
docker --version
```

### Go Environment

Set up Go environment variables:

```bash
# Add to ~/.bashrc or ~/.zshrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Apply changes
source ~/.bashrc
```

### Firewall Configuration

If using UFW, allow necessary ports:

```bash
# Allow Sonic ports (adjust as needed)
sudo ufw allow 8080/tcp
sudo ufw allow 8123/tcp

# Enable firewall
sudo ufw enable
```

## 🔄 Upgrading

```bash
# Pull latest changes
cd ~/Code/SonicScrewdriver
git pull origin main

# Rebuild and reinstall
make clean
make build
sudo make install

# Restart any running services
sonic restart
```

## 📊 System Monitoring

```bash
# Check system resources
sonic system resources

# Monitor container usage
sonic system containers

# Check disk space
sonic system disk

# View system logs
sonic system logs
```

---

*Last Updated: 2024-04-28*  
*Sonic-Screwdriver v2.1.0 for Linux Ubuntu*  
*Optimized for Ubuntu 22.04 LTS and later*
