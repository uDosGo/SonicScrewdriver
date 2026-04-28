#!/bin/bash

# Sonic-Screwdriver Ubuntu Installer
# Comprehensive installation script for Ubuntu 22.04 LTS+

set -euo pipefail

# Configuration
SONIC_REPO="https://github.com/uDosGo/SonicScrewdriver.git"
SONIC_DIR="$HOME/Code/SonicScrewdriver"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="/tmp/sonic-install-$(date +%Y%m%d-%H%M%S).log"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

# Check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        log_error "This script should not be run as root. Please run as a regular user."
        exit 1
    fi
}

# Check Ubuntu version
check_ubuntu_version() {
    log_info "Checking Ubuntu version..."
    
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        if [ "$ID" != "ubuntu" ]; then
            log_error "This script is designed for Ubuntu only. Detected: $ID"
            exit 1
        fi
        
        local version=$(echo "$VERSION_ID" | cut -d '.' -f 1)
        local minor_version=$(echo "$VERSION_ID" | cut -d '.' -f 2)
        
        if [ "$version" -lt 22 ] || ([ "$version" -eq 22 ] && [ "$minor_version" -lt 4 ]); then
            log_error "Ubuntu 22.04 LTS or later is required. Detected: $VERSION_ID"
            exit 1
        fi
        
        log_success "Ubuntu $VERSION_ID detected - compatible"
    else
        log_error "/etc/os-release not found. Are you running Ubuntu?"
        exit 1
    fi
}

# Check architecture
check_architecture() {
    log_info "Checking system architecture..."
    
    local arch=$(uname -m)
    if [ "$arch" != "x86_64" ] && [ "$arch" != "amd64" ]; then
        log_warning "Unsupported architecture: $arch. x86_64/amd64 recommended."
    else
        log_success "Architecture: $arch - compatible"
    fi
}

# Check required dependencies
check_dependencies() {
    log_info "Checking required dependencies..."
    
    local dependencies=("git" "make" "curl" "g++" "docker" "go")
    local missing=()
    
    for dep in "${dependencies[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -eq 0 ]; then
        log_success "All required dependencies are installed"
        return 0
    else
        log_warning "Missing dependencies: ${missing[*]}"
        return 1
    fi
}

# Install dependencies
install_dependencies() {
    log_info "Installing required dependencies..."
    
    sudo apt update
    sudo apt install -y git make curl g++ docker.io golang-go libssl-dev
    
    # Add user to docker group
    sudo usermod -aG docker "$USER"
    
    log_success "Dependencies installed successfully"
    log_info "Note: You may need to log out and back in for Docker group changes to take effect"
}

# Check Go version
check_go_version() {
    log_info "Checking Go version..."
    
    local go_version=$(go version 2>/dev/null | grep -o 'go[0-9.]*' | cut -d' ' -f2)
    
    if [ -z "$go_version" ]; then
        log_error "Go is not installed"
        return 1
    fi
    
    local major=$(echo "$go_version" | cut -d'.' -f1 | tr -d 'go')
    local minor=$(echo "$go_version" | cut -d'.' -f2)
    
    if [ "$major" -lt 1 ] || ([ "$major" -eq 1 ] && [ "$minor" -lt 25 ]); then
        log_warning "Go version $go_version detected. Go 1.25+ recommended."
        return 0
    else
        log_success "Go version $go_version - compatible"
        return 0
    fi
}

# Check Docker status
check_docker() {
    log_info "Checking Docker status..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        return 1
    fi
    
    if ! docker info &> /dev/null; then
        log_warning "Docker daemon is not running"
        log_info "Starting Docker daemon..."
        sudo systemctl start docker
        sudo systemctl enable docker
        
        if docker info &> /dev/null; then
            log_success "Docker daemon started successfully"
        else
            log_error "Failed to start Docker daemon"
            return 1
        fi
    else
        log_success "Docker is running"
    fi
    
    return 0
}

# Clone repository
clone_repository() {
    log_info "Cloning Sonic-Screwdriver repository..."
    
    if [ -d "$SONIC_DIR" ]; then
        log_info "Repository already exists. Pulling latest changes..."
        cd "$SONIC_DIR"
        git pull origin main
    else
        mkdir -p "$HOME/Code"
        git clone "$SONIC_REPO" "$SONIC_DIR"
        cd "$SONIC_DIR"
    fi
    
    log_success "Repository ready at $SONIC_DIR"
}

# Build Sonic-Screwdriver
build_sonic() {
    log_info "Building Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    if [ ! -f "Makefile" ]; then
        log_error "Makefile not found in $SONIC_DIR"
        exit 1
    fi
    
    if ! make build; then
        log_error "Build failed. Check the error messages above."
        exit 1
    fi
    
    if [ ! -f "bin/sonic" ]; then
        log_error "sonic binary not found after build"
        exit 1
    fi
    
    log_success "Sonic-Screwdriver built successfully"
}

# Install Sonic-Screwdriver
install_sonic() {
    log_info "Installing Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    if [ ! -f "bin/sonic" ]; then
        log_error "sonic binary not found. Run build first."
        exit 1
    fi
    
    sudo cp "bin/sonic" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/sonic"
    
    log_success "Sonic-Screwdriver installed to $INSTALL_DIR/sonic"
}

# Setup environment
setup_environment() {
    log_info "Setting up environment..."
    
    # Create .sonic directory
    mkdir -p "$HOME/.sonic"
    mkdir -p "$HOME/.sonic/logs"
    
    # Add to bashrc if not already present
    if ! grep -q "export PATH=\"$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc"; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    fi
    
    if ! grep -q "export GOPATH=\"$HOME/go\"" "$HOME/.bashrc"; then
        echo 'export GOPATH=$HOME/go' >> "$HOME/.bashrc"
        echo 'export PATH=$PATH:$GOPATH/bin' >> "$HOME/.bashrc"
    fi
    
    # Source bashrc to apply changes
    source "$HOME/.bashrc"
    
    log_success "Environment setup complete"
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if ! command -v sonic &> /dev/null; then
        log_error "sonic command not found in PATH"
        return 1
    fi
    
    local version=$(sonic --version 2>/dev/null || echo "unknown")
    
    if [ "$version" != "unknown" ]; then
        log_success "Sonic-Screwdriver installed successfully: $version"
    else
        log_error "sonic command not working properly"
        return 1
    fi
    
    # Test system checks
    if sonic system check &> /dev/null; then
        log_success "System compatibility check passed"
    else
        log_warning "System compatibility check failed"
    fi
    
    return 0
}

# Repair installation
repair_installation() {
    log_info "Repairing installation..."
    
    # Check what's broken
    local issues=()
    
    if ! command -v sonic &> /dev/null; then
        issues+=("sonic_binary_missing")
    fi
    
    if ! command -v docker &> /dev/null; then
        issues+=("docker_not_installed")
    elif ! docker info &> /dev/null; then
        issues+=("docker_not_running")
    fi
    
    if ! command -v go &> /dev/null; then
        issues+=("go_not_installed")
    fi
    
    if [ ${#issues[@]} -eq 0 ]; then
        log_success "No issues found. Installation appears healthy."
        return 0
    fi
    
    log_info "Found issues: ${issues[*]}"
    
    # Fix issues
    for issue in "${issues[@]}"; do
        case "$issue" in
            "sonic_binary_missing")
                log_info "Reinstalling sonic binary..."
                install_sonic
                ;;
            "docker_not_installed")
                log_info "Installing Docker..."
                sudo apt install -y docker.io
                sudo usermod -aG docker "$USER"
                ;;
            "docker_not_running")
                log_info "Starting Docker..."
                sudo systemctl start docker
                sudo systemctl enable docker
                ;;
            "go_not_installed")
                log_info "Installing Go..."
                sudo apt install -y golang-go
                ;;
        esac
    done
    
    log_success "Repair complete. Please verify the installation."
}

# Setup scripts
setup_scripts() {
    log_info "Running setup scripts..."
    
    cd "$SONIC_DIR"
    
    # Run setup scripts if they exist
    local scripts=("setup-git.sh" "setup-logging.sh" "setup-vendors.sh")
    
    for script in "${scripts[@]}"; do
        if [ -f "$script" ]; then
            log_info "Running $script..."
            if ! ./"$script"; then
                log_warning "Failed to run $script"
            fi
        fi
    done
    
    log_success "Setup scripts completed"
}

# Main installation function
main_install() {
    log_info "Starting Sonic-Screwdriver installation on Ubuntu..."
    
    # Step 1: System checks
    check_root
    check_ubuntu_version
    check_architecture
    
    # Step 2: Dependency checks and installation
    if ! check_dependencies; then
        install_dependencies
    fi
    
    if ! check_go_version; then
        log_warning "Go version check failed"
    fi
    
    if ! check_docker; then
        log_error "Docker setup failed"
        exit 1
    fi
    
    # Step 3: Clone and build
    clone_repository
    build_sonic
    
    # Step 4: Install
    install_sonic
    setup_environment
    
    # Step 5: Setup scripts
    setup_scripts
    
    # Step 6: Verify
    if verify_installation; then
        log_success "Installation completed successfully!"
        echo -e "${GREEN}
========================================
Sonic-Screwdriver is ready to use!
========================================
${NC}"
        echo "Try these commands:"
        echo "  sonic --help              Show help"
        echo "  sonic --version           Show version"
        echo "  sonic system check        Check system compatibility"
        echo "  sonic tui                 Launch interactive interface"
        echo ""
        echo "Documentation:"
        echo "  $SONIC_DIR/README_LINUX.md"
        echo "  $SONIC_DIR/QUICKSTART_LINUX.md"
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

# Show help
show_help() {
    echo "Sonic-Screwdriver Ubuntu Installer"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  install       - Full installation (default)"
    echo "  check         - Check system requirements"
    echo "  dependencies  - Install dependencies only"
    echo "  build         - Build Sonic-Screwdriver only"
    echo "  install-bin   - Install binary only"
    echo "  verify        - Verify installation"
    echo "  repair        - Repair installation"
    echo "  help          - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 install           # Full installation"
    echo "  $0 check             # Check requirements"
    echo "  $0 repair            # Repair installation"
}

# Main execution
COMMAND="${1:-install}"

case "$COMMAND" in
    install|install-full)
        main_install
        ;;
    check)
        check_root
        check_ubuntu_version
        check_architecture
        check_dependencies
        check_go_version
        check_docker
        ;;
    dependencies|deps)
        check_root
        install_dependencies
        ;;
    build)
        clone_repository
        build_sonic
        ;;
    install-bin)
        build_sonic
        install_sonic
        ;;
    verify)
        verify_installation
        ;;
    repair)
        repair_installation
        verify_installation
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac
