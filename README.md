# Sonic Family

Monorepo for SonicScrewdriver container runtime, Code Vault shared contracts, and Ventoy bootable USB fork.

## Components

| Component | Description | Version |
|-----------|-------------|---------|
| **code-vault** | Shared types, protocols, API contracts | vA1.0.0 |
| **sonic-screwdriver** | Docker wrapper, container runtime, CLI/TUI/GUI | vA1.0.0 |
| **ventoy** | Bootable USB tool (fork) | upstream + patches |

## Quick Start

```bash
# Build all
make build

# Install sonic
cd sonic-screwdriver && make install

# Install a game
sonic install doom
```

## License

MIT
