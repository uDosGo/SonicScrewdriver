# sonic-screwdriver

Repository for sonic-screwdriver runtime/CLI with modular components under `modules/`.

`sonic-screwdriver` is a separate project from uHomeNest. Both projects may carry
`sonic-home` and/or `sonic-express` modules to stay standards-aligned, while
implementation, runtime behavior, and release ownership remain independent per system.

## Components

| Component | Description | Version |
|-----------|-------------|---------|
| **core (root)** | Docker wrapper, container runtime, CLI/TUI/GUI | vA1.0.0 |
| **code-vault** | Shared types, protocols, API contracts | vA1.0.0 |
| **modules/sonic-express** | uDos-aligned packager and installer helper | v0.1.0-dev |
| **modules/sonic-home** | Lite packager and installer helper | v0.1.0-dev |
| **modules/ventoy** | Bootable USB tool (fork) | upstream + patches |

## Quick Start

```bash
# Build all
make build

# Install sonic
make install

# Install a game
sonic install doom
```


## Roadmap

See `docs/roadmap.md` for the current milestone plan and delivery criteria.

## License

MIT
