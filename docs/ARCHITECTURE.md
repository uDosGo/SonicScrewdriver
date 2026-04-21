# Sonic-Screwdriver Architecture

## 🏗️ Overview

Sonic-Screwdriver is a modular API Central Hub designed for smart home automation, game management, and system integration. The architecture follows Classic Modern Mint principles with a focus on modularity, security, and extensibility.

## 📦 Component Architecture

```
Sonic-Screwdriver
├── API Central Hub (Core)
│   ├── Secret Store
│   ├── API Proxy
│   ├── Node Registry
│   └── Authentication
├── Integrations
│   ├── Home Assistant
│   ├── Media Player
│   ├── Feed/Spool System
│   └── Remote Access
├── Container Management
│   ├── Docker Runtime
│   ├── Health Monitoring
│   └── Lifecycle Management
├── User Interface
│   ├── CLI
│   ├── TUI (Bubble Tea)
│   └── Web (Future)
└── Data Layer
    ├── SQLite State
    ├── Configuration
    └── Cache
```

## 🔐 Core Components

### 1. Secret Store

**Purpose**: Secure storage and management of sensitive data

**Features**:
- AES-256-GCM encryption
- Secret rotation with history
- Access control policies
- Atomic updates
- Backup/restore functionality

**Implementation**:
- `internal/secrets/secret_store.go`
- SQLite backend with encryption
- In-memory cache for performance

### 2. API Proxy

**Purpose**: Secure proxy for external API calls

**Features**:
- Rate limiting
- Request/response logging
- Provider-specific configurations
- Error handling and retries
- Health monitoring

**Implementation**:
- `internal/secrets/proxy.go`
- Pluggable provider system
- Circuit breaker pattern

### 3. Node Registry

**Purpose**: Management of registered nodes

**Features**:
- Node registration and authentication
- Access control
- Status monitoring
- Secret access management

**Implementation**:
- `internal/secrets/node_registry.go`
- JSON-based storage
- Token-based authentication

### 4. Container Runtime

**Purpose**: Game and service container management

**Features**:
- Docker container lifecycle
- Health monitoring
- Log management
- Resource constraints

**Implementation**:
- `internal/container/docker.go`
- Docker SDK integration
- Health check system

## 🔌 Integration Architecture

### Home Assistant Integration

**Purpose**: Deep integration with Home Assistant

**Features**:
- Iframe embed strategy
- Kiosk mode support
- API connectivity testing
- Configuration management

**Implementation**:
- `internal/homeassistant/integration.go`
- HTML iframe generation
- REST API client

### Media Player System

**Purpose**: Local media management and playback

**Features**:
- Media scanning and indexing
- Metadata extraction
- Library management
- Playback control

**Implementation**:
- SQLite media catalog
- FFmpeg integration
- Web-based player

### Feed/Spool System

**Purpose**: Content aggregation and processing

**Features**:
- Feed parsing (RSS, Atom, JSON)
- Content validation
- Spool processing pipeline
- Notification system

**Implementation**:
- Unified feed format
- Pipeline architecture
- Event-driven processing

## 🎮 Game Management

### Library System

**Purpose**: Curated game management

**Features**:
- Manifest validation
- Index management
- Installation tracking
- Version management

**Implementation**:
- YAML manifests
- SQLite index
- Docker-based execution

### Container Lifecycle

**Purpose**: Game container management

**Features**:
- Install/start/stop/remove
- Health monitoring
- Log management
- Resource isolation

**Implementation**:
- Docker containers
- State tracking
- Health checks

## 📡 Remote Access

### VNC Server

**Purpose**: Remote desktop access

**Features**:
- Password authentication
- Configurable resolution
- Session management

**Implementation**:
- `internal/remote/vnc.go`
- TigerVNC integration

### SSH Access

**Purpose**: Secure shell access

**Features**:
- Key-based authentication
- Configuration management
- Access control

**Implementation**:
- OpenSSH integration
- Configuration templates

### Samba Sharing

**Purpose**: File sharing

**Features**:
- Share management
- Access control
- Configuration

**Implementation**:
- Samba configuration
- Share management

## 🎨 Classic Modern Mint

### Theme System

**Purpose**: Classic Modern theme integration

**Features**:
- Readiness checking
- Theme installation
- Theme application
- Diagnostic tools

**Implementation**:
- `internal/classicmodern/`
- SQLite state tracking

## 📁 Data Architecture

### State Management

**Purpose**: Persistent state storage

**Features**:
- Installation tracking
- Runtime status
- Configuration
- History

**Implementation**:
- SQLite database
- Migration system
- Transaction support

### Configuration

**Purpose**: Application configuration

**Features**:
- JSON configuration
- Environment variables
- Command-line overrides

**Implementation**:
- Config structs
- Validation
- Merge logic

### Cache Layer

**Purpose**: Performance optimization

**Features**:
- In-memory caching
- Cache invalidation
- Offline mode support

**Implementation**:
- Memory cache
- Fallback logic

## 🔒 Security Architecture

### Encryption

**Purpose**: Data protection

**Features**:
- AES-256-GCM for secrets
- Master key encryption
- Secure key storage

**Implementation**:
- Crypto packages
- Key management

### Authentication

**Purpose**: Access control

**Features**:
- Token-based auth
- Node authentication
- Role-based access

**Implementation**:
- JWT tokens
- Access control lists

### Audit Logging

**Purpose**: Activity tracking

**Features**:
- Operation logging
- Error tracking
- Audit trails

**Implementation**:
- Log files
- Structured logging

## 🧪 Testing Architecture

### Unit Tests

**Purpose**: Component testing

**Features**:
- Isolated testing
- Mock dependencies
- High coverage

**Implementation**:
- Go testing package
- Testify assertions

### Integration Tests

**Purpose**: System testing

**Features**:
- End-to-end testing
- Real dependencies
- Scenario testing

**Implementation**:
- Test containers
- Real databases

### Health Checks

**Purpose**: Runtime monitoring

**Features**:
- Container health
- Service monitoring
- Status reporting

**Implementation**:
- Health endpoints
- Status commands

## 📚 Design Principles

### Modularity
- Clear component boundaries
- Minimal dependencies
- Pluggable architecture

### Security
- Defense in depth
- Least privilege
- Secure by default

### Reliability
- Error handling
- Retry logic
- Graceful degradation

### Maintainability
- Clear documentation
- Consistent style
- Comprehensive tests

### Extensibility
- Plugin architecture
- Configuration-driven
- Open interfaces

## 🔧 Technical Stack

### Languages
- **Primary**: Go (Golang)
- **Secondary**: JavaScript, Python
- **Configuration**: YAML, JSON

### Databases
- **Primary**: SQLite
- **Secondary**: (Future) PostgreSQL

### Technologies
- Docker for containerization
- Bubble Tea for TUI
- REST and WebSocket APIs
- AES-256-GCM for encryption

### Tools
- Go modules for dependency management
- Make for build automation
- Git for version control
- SQLite for data storage

## 🗺️ Deployment Architecture

### Local Deployment
```
User → CLI/TUI → Sonic-Screwdriver → Docker → Services
```

### Network Deployment
```
Client → API → Sonic-Screwdriver → Docker → Services
                        ↓
                   SQLite State
```

### Home Assistant Integration
```
Sonic-Screwdriver → HA API → Home Assistant
     ↓
Iframe Embed → Browser
```

## 📈 Performance Considerations

### Optimization Areas
- Secret cache for frequent access
- Batch operations for bulk actions
- Connection pooling for APIs
- Lazy loading for large datasets

### Scalability
- Modular design for horizontal scaling
- Stateless components where possible
- Resource constraints for containers

## 🔮 Future Architecture Evolution

### v2.2.0 Plans
- Microservices architecture
- Plugin system
- Enhanced API layer
- Improved monitoring

### v3.0.0 Vision
- Cloud-native deployment
- Distributed architecture
- Advanced AI integration
- Enterprise features

---

*Last Updated: 2026-04-29*
*Sonic-Screwdriver v2.1.0*
*Architecture Version: 2.1*