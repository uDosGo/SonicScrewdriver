# Sonic-Screwdriver CLI Commands

## 📖 Command Reference

### Basic Commands

```bash
# Show help
sonic --help
sonic -h

# Show version
sonic --version
sonic -v

# Launch interactive TUI
sonic tui
sonic menu
```

### Game Management

```bash
# Install a game
sonic install <game-name>

# Start a game
sonic start <game-name>

# Stop a game
sonic stop <game-name>

# List installed games
sonic list

# Remove a game
sonic remove <game-name>

# Show game logs
sonic logs <game-name>

# Check game health
sonic health <game-name>
sonic health --all

# Repair game containers
sonic repair <game-name>
sonic repair --all
```

### Secret Management

```bash
# Add a secret
sonic secret add <name> --value <value>

# Get a secret
sonic secret get <name>

# Rotate a secret
sonic secret rotate <name> --value <new-value>

# Show secret history
sonic secret history <name>

# List all secrets
sonic secret list

# Grant secret access
sonic secret grant <name> --node <node-name>

# Revoke secret access
sonic secret revoke <name> --node <node-name>

# Show secret policy
sonic secret policy <name>

# Backup secrets
sonic secret backup <file>

# Restore secrets
sonic secret restore <file>

# Export encrypted backup
sonic secret export <file>

# Import encrypted backup
sonic secret import <file>
```

### Node Management

```bash
# Register a node
sonic node register --master <address> --name <node-name>

# List registered nodes
sonic node list

# Show node details
sonic node show <name>

# Revoke node access
sonic node revoke <name>
```

### API Proxy

```bash
# Show proxy status
sonic proxy status

# Make proxy call
sonic proxy call <provider> --data <json-data>

# Check proxy health
sonic proxy health
```

### Home Assistant Integration (v2.1.0)

```bash
# Setup HA integration
sonic ha setup <url> <token>

# Configure HA integration
sonic ha configure

# Show HA status
sonic ha status

# Show HA instance info
sonic ha info

# Generate embed HTML
sonic ha embed <output.html>

# Enable/disable kiosk mode
sonic ha kiosk enable
sonic ha kiosk disable

# Set refresh rate (minutes)
sonic ha refresh <minutes>

# Get HA version
sonic ha version

# Check HA connection
sonic ha check
```

### Remote Access

```bash
# Setup VNC server
sonic remote vnc setup [password] [geometry]
sonic remote vnc start
sonic remote vnc stop

# Setup SSH access
sonic remote ssh setup

# Setup Samba sharing
sonic remote samba setup <name> <path>

# Show remote access info
sonic remote info
```

### Classic Modern Mint

```bash
# Check Classic Modern readiness
sonic mint check

# Install Classic Modern theme
sonic mint install

# Apply Classic Modern theme
sonic mint apply

# Show theme status
sonic mint status

# Show theme information
sonic mint info

# Run diagnostic checks
sonic mint doctor
```

### Ventoy Integration

```bash
# Create Ventoy bundle
sonic ventoy package <installer-dir> <output.she>

# Validate Ventoy bundle
sonic ventoy validate <bundle.she>

# Show bundle information
sonic ventoy info <bundle.she>
```

### Library Management

```bash
# List available games
sonic library list
```

### Configuration

```bash
# Set configuration value
sonic config set <key> <value>
```

## 🎯 Command Categories

### Core Commands
- `install`, `start`, `stop`, `list`, `remove` - Game lifecycle
- `tui`, `menu` - Interactive interfaces
- `logs`, `health`, `repair` - Monitoring and maintenance

### Integration Commands
- `ha` - Home Assistant integration
- `remote` - Remote access management
- `mint` - Classic Modern theme management
- `ventoy` - Ventoy bundle management

### Management Commands
- `secret` - Secret store management
- `node` - Node registry management
- `proxy` - API proxy management
- `library` - Game library management
- `config` - Configuration management

## 📋 Command Structure

```
sonic <category> <command> [options] [arguments]
```

### Categories
- `game` - Game management (default)
- `secret` - Secret management
- `node` - Node management
- `proxy` - API proxy
- `ha` - Home Assistant
- `remote` - Remote access
- `mint` - Classic Modern
- `ventoy` - Ventoy integration
- `library` - Game library
- `config` - Configuration

### Common Options
- `--help`, `-h` - Show help
- `--version`, `-v` - Show version
- `--verbose` - Verbose output (where supported)
- `--dry-run` - Test without making changes (where supported)

## 💡 Usage Tips

### Secret Management
```bash
# Add API key with automatic testing
sonic secret add openrouter_api_key --value "sk-..."

# Rotate secret and show history
sonic secret rotate openrouter_api_key --value "sk-new..."
sonic secret history openrouter_api_key
```

### Home Assistant
```bash
# Setup and generate embed
sonic ha setup "http://ha.local:8123" "eyJhbGciOiJIUz..."
sonic ha embed /var/www/ha-embed.html
sonic ha kiosk enable
sonic ha refresh 60
```

### Game Management
```bash
# Full lifecycle
sonic install my-game
sonic start my-game
sonic health my-game
sonic stop my-game
sonic remove my-game
```

### Remote Access
```bash
# Setup complete remote access
sonic remote vnc setup "password" "1920x1080"
sonic remote vnc start
sonic remote ssh setup
sonic remote samba setup "shared" "/home/user/shared"
```

## 🔧 Advanced Usage

### Batch Operations
```bash
# Start multiple games
for game in game1 game2 game3; do
    sonic start "$game"
done

# Health check all games
sonic health --all
```

### Scripting
```bash
# Check if game is running
if sonic health "my-game" | grep -q "healthy"; then
    echo "Game is healthy"
else
    echo "Game needs attention"
fi
```

### Configuration
```bash
# Setup environment
sonic config set docker.socket "/var/run/docker.sock"
sonic config set ha.url "http://ha.local:8123"
```

## 📚 Help System

```bash
# General help
sonic --help

# Category-specific help
sonic secret --help
sonic ha --help
sonic remote --help

# Command-specific help
sonic secret add --help
sonic ha setup --help
```

## 🎨 Classic Modern Integration

```bash
# Complete theme setup
sonic mint check
sonic mint install
sonic mint apply
sonic mint doctor
```

## 🌐 Network Commands

```bash
# Remote management
ssh user@remote-host sonic health --all

# API access
curl -X POST http://localhost:8080/api/secret/get -d '{"name":"key"}'
```

## 📈 Monitoring

```bash
# System health
sonic health --all

# Proxy status
sonic proxy status

# Node status
sonic node list
```

---

*Last Updated: 2026-04-29*
*Sonic-Screwdriver v2.1.0*
*Command Reference Version: 2.1*