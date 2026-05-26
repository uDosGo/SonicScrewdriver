# Sonic-Screwdriver Library Format

## 📚 Overview

The Sonic-Screwdriver library system manages curated game installations using YAML manifests and a SQLite index. This document specifies the library format and structure.

## 📁 Library Structure

```
library/
├── index.yaml          # Main library index
├── games/
│   ├── game1/
│   │   ├── manifest.yaml  # Game manifest
│   │   ├── icon.png       # Game icon (optional)
│   │   └── assets/        # Additional assets (optional)
│   ├── game2/
│   │   ├── manifest.yaml
│   │   └── icon.png
│   └── ...
└── categories/         # Category definitions (optional)
    └── categories.yaml
```

## 📋 Manifest Format

### Game Manifest (`manifest.yaml`)

```yaml
# Required fields
name: "game-name"              # Unique identifier (slug format)
display_name: "Game Name"      # Human-readable name
version: "1.0.0"               # Semantic version
description: "A description"   # Game description

# Container configuration
container:
  image: "registry/image:tag"  # Docker image
  ports:                        # Port mappings
    - "8080:8080"
  volumes:                      # Volume mounts
    - "/data:/app/data"
  env:                          # Environment variables
    ENV_VAR: "value"
  command: ["/start"]          # Start command
  health_check:                 # Health check
    path: "/health"
    interval: 30
    timeout: 5

# Metadata
metadata:
  category: "adventure"        # Game category
  publisher: "Publisher Name"   # Publisher
  release_date: "2023-01-01"    # Release date
  website: "https://example.com" # Website URL
  support: "support@example.com" # Support email

# Requirements
requirements:
  min_cpu: 2                    # Minimum CPU cores
  min_ram: 4                    # Minimum RAM (GB)
  min_disk: 10                  # Minimum disk space (GB)
  gpu: false                    # GPU required
  os: ["linux", "windows"]     # Supported OS

# Installation
install:
  size: 500                     # Install size (MB)
  duration: 300                 # Estimated install time (seconds)
  steps:                        # Installation steps
    - "Download base image"
    - "Configure volumes"
    - "Set environment"

# Updates
updates:
  auto: true                    # Auto-update enabled
  channel: "stable"            # Update channel
  schedule: "daily"            # Update check schedule

# Dependencies
dependencies:
  - "dependency1"              # Required dependencies
  - "dependency2"

# Compatibility
compatibility:
  sonic_version: ">=2.0.0"      # Minimum Sonic-Screwdriver version
  docker_version: ">=20.10"      # Minimum Docker version

# Status
status: "available"            # availability status
visible: true                   # visible in library
featured: false                  # featured game

# Tags
tags:
  - "adventure"
  - "single-player"
  - "indie"

# Screenshots (optional)
screenshots:
  - "screenshot1.jpg"
  - "screenshot2.jpg"

# Videos (optional)
videos:
  - "trailer.mp4"
  - "gameplay.mp4"

# Ratings (optional)
ratings:
  esrb: "T"                     # ESRB rating
  pegi: "12"                    # PEGI rating

# Achievements (optional)
achievements:
  total: 25                     # Total achievements
  list:                         # Achievement list
    - id: "achievement1"
      name: "First Step"
      description: "Complete tutorial"

# DLC (optional)
dlc:
  - id: "dlc1"
    name: "Expansion Pack"
    price: 9.99
    description: "Additional content"

# Mods (optional)
mods:
  supported: true               # Mods supported
  workshop_id: "123456"         # Steam Workshop ID (if applicable)
```

## 📊 Index Format

### Main Index (`index.yaml`)

```yaml
# Library metadata
version: "1.1"                  # Index format version
updated: "2026-04-29"           # Last update date
maintainer: "Sonic Team"        # Maintainer

# Game index
games:
  - id: "game1"                 # Game ID (matches manifest name)
    name: "Game One"            # Display name
    version: "1.0.0"            # Current version
    status: "available"         # Status
    category: "adventure"       # Primary category
    tags: ["adventure", "indie"] # Tags
    icon: "games/game1/icon.png" # Icon path
    manifest: "games/game1/manifest.yaml" # Manifest path
    installed: false             # Installation status
    last_played: null             # Last played timestamp
    play_time: 0                 # Total play time (minutes)
    rating: null                 # User rating (1-5)
    favorite: false               # Favorite status

  - id: "game2"
    name: "Game Two"
    version: "2.1.0"
    status: "available"
    category: "puzzle"
    tags: ["puzzle", "strategy"]
    icon: "games/game2/icon.png"
    manifest: "games/game2/manifest.yaml"
    installed: true
    last_played: "2026-04-28T14:30:00Z"
    play_time: 120
    rating: 4
    favorite: true

# Categories
categories:
  - id: "adventure"
    name: "Adventure"
    description: "Adventure games"
    icon: "adventure.png"

  - id: "puzzle"
    name: "Puzzle"
    description: "Puzzle and strategy games"
    icon: "puzzle.png"

# Statistics
stats:
  total_games: 2                # Total games in library
  installed_games: 1           # Installed games
  total_play_time: 120          # Total play time (minutes)
  last_updated: "2026-04-29T10:00:00Z" # Last update
```

## 🗃️ Database Schema

### SQLite State Database

```sql
-- Games table
CREATE TABLE games (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    status TEXT NOT NULL,
    category TEXT,
    icon_path TEXT,
    manifest_path TEXT NOT NULL,
    installed BOOLEAN DEFAULT FALSE,
    install_date TEXT,
    last_played TEXT,
    play_time INTEGER DEFAULT 0,
    rating INTEGER,
    favorite BOOLEAN DEFAULT FALSE,
    last_updated TEXT NOT NULL
);

-- Game metadata table
CREATE TABLE game_metadata (
    game_id TEXT PRIMARY KEY,
    publisher TEXT,
    release_date TEXT,
    website TEXT,
    support_email TEXT,
    description TEXT,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Game requirements table
CREATE TABLE game_requirements (
    game_id TEXT PRIMARY KEY,
    min_cpu INTEGER,
    min_ram INTEGER,
    min_disk INTEGER,
    gpu_required BOOLEAN,
    supported_os TEXT,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Game tags table
CREATE TABLE game_tags (
    game_id TEXT,
    tag TEXT,
    PRIMARY KEY (game_id, tag),
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Categories table
CREATE TABLE categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    icon_path TEXT
);

-- Installations table
CREATE TABLE installations (
    game_id TEXT PRIMARY KEY,
    install_date TEXT NOT NULL,
    last_update TEXT,
    status TEXT NOT NULL,
    error_message TEXT,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Play sessions table
CREATE TABLE play_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id TEXT NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT,
    duration INTEGER,
    FOREIGN KEY (game_id) REFERENCES games(id)
);
```

## 📦 Library Management

### Adding Games

1. **Create game directory**: `library/games/<game-id>/`
2. **Create manifest**: `manifest.yaml` with required fields
3. **Add to index**: Update `index.yaml` with game entry
4. **Validate**: Run `sonic library validate`

### Updating Games

1. **Update manifest**: Modify `manifest.yaml`
2. **Update index**: Modify `index.yaml`
3. **Increment version**: Update version number
4. **Validate**: Run `sonic library validate`

### Removing Games

1. **Remove from index**: Update `index.yaml`
2. **Remove directory**: Delete `library/games/<game-id>/`
3. **Clean database**: Run `sonic library cleanup`

## 🔍 Validation

### Manifest Validation

The system validates:
- Required fields present
- Valid YAML syntax
- Semantic version format
- Container configuration valid
- File paths exist
- Icon dimensions correct

### Index Validation

The system validates:
- No duplicate IDs
- All referenced manifests exist
- Valid status values
- Consistent versioning

## 🎮 Game Lifecycle

### Installation

```bash
# Install from library
sonic install <game-id>

# Installation process
1. Validate manifest
2. Pull Docker image
3. Create volumes
4. Set environment
5. Start container
6. Update state
```

### Starting

```bash
# Start game
sonic start <game-id>

# Start process
1. Check installed
2. Validate container
3. Start container
4. Update status
5. Record session
```

### Stopping

```bash
# Stop game
sonic stop <game-id>

# Stop process
1. Check running
2. Stop container
3. Update status
4. End session
```

### Removal

```bash
# Remove game
sonic remove <game-id>

# Removal process
1. Stop container
2. Remove container
3. Clean volumes
4. Update state
5. Remove from index
```

## 📊 Library Operations

### List Games

```bash
# List all games
sonic library list

# List installed games
sonic library list --installed

# List by category
sonic library list --category adventure
```

### Search Games

```bash
# Search by name
sonic library search "adventure"

# Search by tag
sonic library search --tag "indie"
```

### Game Information

```bash
# Show game info
sonic library info <game-id>

# Show detailed info
sonic library info <game-id> --detailed
```

### Validate Library

```bash
# Validate all manifests
sonic library validate

# Validate specific game
sonic library validate <game-id>
```

### Update Library

```bash
# Update from remote
sonic library update

# Update specific game
sonic library update <game-id>
```

## 🎨 Icons and Assets

### Icon Requirements
- **Format**: PNG
- **Size**: 256x256 pixels
- **Color**: 32-bit RGBA
- **Naming**: `icon.png`

### Screenshot Requirements
- **Format**: JPG or PNG
- **Size**: 1280x720 pixels recommended
- **Color**: 24-bit RGB or 32-bit RGBA
- **Naming**: `screenshot1.jpg`, `screenshot2.jpg`, etc.

### Video Requirements
- **Format**: MP4
- **Codec**: H.264
- **Resolution**: 1280x720 or 1920x1080
- **Bitrate**: 5 Mbps recommended
- **Naming**: `trailer.mp4`, `gameplay.mp4`, etc.

## 📚 Best Practices

### Manifest Design
- Use semantic versioning
- Provide clear descriptions
- Specify accurate requirements
- Include comprehensive metadata
- Use appropriate tags and categories

### Library Organization
- Group games by category
- Use consistent naming conventions
- Include icons for all games
- Keep manifests up to date
- Validate regularly

### Performance
- Optimize container configurations
- Specify resource limits
- Use efficient health checks
- Minimize startup time
- Optimize asset sizes

## 🔧 Troubleshooting

### Common Issues

**Manifest validation failed**
- Check required fields
- Validate YAML syntax
- Verify file paths
- Check version format

**Game won't install**
- Verify Docker image exists
- Check port conflicts
- Validate volume permissions
- Review container logs

**Game won't start**
- Check container status
- Verify dependencies
- Review resource requirements
- Check health status

**Library validation errors**
- Check for duplicate IDs
- Verify all manifests exist
- Validate index structure
- Review category definitions

## 📈 Future Enhancements

### Planned Features
- Automatic library updates
- Remote library sync
- User ratings and reviews
- Game recommendations
- Achievement tracking
- Cloud save support

### Architecture Improvements
- Plugin system for game types
- Enhanced validation rules
- Better error reporting
- Performance optimization

---

*Last Updated: 2026-04-29*
*Sonic-Screwdriver v2.1.0*
*Library Format Version: 1.1*