# Family Ecosystem Requirements

This document outlines the requirements for the family ecosystem, including sonic, udos, and related components.

## Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu/Debian recommended)
- **CPU**: Multi-core processor
- **Memory**: 4GB RAM minimum, 8GB recommended
- **Disk Space**: 10GB free space
- **Dependencies**:
  - `curl`
  - `git`
  - `xclip`
  - `docker`
  - `nodejs` (v24+)
  - `npm`

### Installation Steps

1. **Install Prerequisites**
   ```bash
   sudo apt update
   sudo apt install -y curl git xclip docker.io nodejs npm
   ```

2. **Install sonic**
   ```bash
   curl -sSL https://raw.githubusercontent.com/sonic-family/installer/main/bootstrap.sh -o /tmp/bootstrap.sh
   chmod +x /tmp/bootstrap.sh
   /tmp/bootstrap.sh
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Install udos**
   ```bash
   sonic install uDos
   sudo ln -sf ~/.local/udos/bin/udos /usr/local/bin/udos
   ```

4. **Set Up Universal Spine**
   ```bash
   mkdir -p ~/vault/{system,home,family,user,@inbox,@workspace,@toybox,@sandbox,@public,@private,binder}
   mkdir -p ~/.local/udos/{compartments,compost,legacy,trash,feeds}
   echo 'export UDOS_VAULT="$HOME/vault"' >> ~/.bashrc
   echo 'export UDOS_CODE="$HOME/code-vault"' >> ~/.bashrc
   echo 'export UDOS_STATE="$HOME/.local/udos"' >> ~/.bashrc
   source ~/.bashrc
   ```

5. **Configure Logging**
   ```bash
   mkdir -p ~/.local/share/sonic/logs
   mkdir -p ~/.local/udos/logs
   ```

6. **Set Up Local Library**
   ```bash
   mkdir -p ~/.local/share/sonic/library
   ```

7. **Initialize Git Repositories**
   ```bash
   cd ~/vault && git init
   cd ~/code-vault && git init
   ```

## Component Requirements

### sonic-screwdriver
- **Language**: Go
- **Build Tool**: `make`
- **Dependencies**: Docker, SQLite
- **Configuration**: `~/.config/sonic/config.yaml`

### uDos
- **Language**: TypeScript
- **Build Tool**: `npm`
- **Dependencies**: Node.js (v24+), Docker
- **Configuration**: `~/.config/udos/config.yaml`

### CHASIS
- **Language**: TypeScript
- **Dependencies**: Node.js (v24+), Docker
- **Configuration**: Part of uDos

### Universal Spine
- **Directories**:
  - `~/vault/`
  - `~/.local/udos/`
  - `~/code-vault/`
- **Environment Variables**:
  - `UDOS_VAULT`
  - `UDOS_CODE`
  - `UDOS_STATE`

## Development Requirements

### sonic-screwdriver
- **Go**: 1.22+
- **Tools**: `make`, `go build`
- **Testing**: Docker, `go test`

### uDos
- **Node.js**: v24+
- **Tools**: `npm`, `tsc`
- **Testing**: Docker, `npm test`

### CHASIS
- **Node.js**: v24+
- **Tools**: `npm`, `tsc`
- **Testing**: Docker, `npm test`

## Vendor Requirements

### Vendors Library
- **Directory**: `~/.local/share/sonic/vendors/`
- **Manifests**: YAML files describing vendor configurations
- **Database**: SQLite or JSON for vendor metadata

### Databases
- **Sonic**: SQLite (`~/.local/share/sonic/sonic.db`)
- **uDos**: SQLite (`~/.local/udos/udos.db`)

## Git Operations

### Vault
- **Repository**: `~/vault/`
- **Ignore**: `@private/`, `trash/`, `compost/`, `*.log`, `*.tmp`

### Code-Vault
- **Repository**: `~/code-vault/`
- **Ignore**: `node_modules/`, `dist/`, `*.log`, `*.tmp`

## Upgrades

### sonic
```bash
sonic update
```

### uDos
```bash
sonic update uDos
```

### All Extensions
```bash
sonic update --all
```

## Logging

### sonic
- **Log File**: `~/.local/share/sonic/logs/sonic.log`
- **Error Log**: `~/.local/share/sonic/logs/sonic-errors.log`
- **Audit Log**: `~/.local/share/sonic/logs/sonic-audit.log`

### uDos
- **Log File**: `~/.local/udos/logs/udos.log`

## Troubleshooting

### Common Issues
- **Docker Permission Denied**: Add user to Docker group
  ```bash
  sudo usermod -aG docker $USER
  newgrp docker
  ```
- **Node.js Version Mismatch**: Use `nvm` to manage Node.js versions
  ```bash
  nvm install 24
  nvm use 24
  ```
- **Missing Dependencies**: Install required packages
  ```bash
  sudo apt install -y curl git xclip docker.io nodejs npm
  ```

### Logs
- **Sonic Logs**: `~/.local/share/sonic/logs/`
- **uDos Logs**: `~/.local/udos/logs/`

### Configuration
- **Sonic Config**: `~/.config/sonic/config.yaml`
- **uDos Config**: `~/.config/udos/config.yaml`

## Conclusion

This document outlines the requirements and setup steps for the family ecosystem. Follow these steps to ensure a smooth installation and configuration process.
