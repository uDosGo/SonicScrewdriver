#!/bin/bash

# Sonic-Screwdriver Installer - TARDIS Edition for Linux
# The Doctor's Ultimate Installation System
# A time-traveling, self-healing installer with Whovian flair

set -euo pipefail

# ============================================
# TARDIS Configuration
# ============================================

TARDIS_NAME="Sonic TARDIS"
DOCTOR_NUMBER="14th"
COMPANION="Linux Ubuntu"

SONIC_REPO="https://github.com/uDosGo/SonicScrewdriver.git"
SONIC_DIR="$HOME/Code/SonicScrewdriver"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="/tmp/sonic-tardis-$(date +%Y%m%d-%H%M%S).log"
STATE_FILE="/tmp/sonic-tardis-state.txt"

# Timey-wimey settings
STEP_TIMEOUT=600
SPINNER_DELAY=0.1

# Colors (TARDIS Console Theme)
TARDIS_BLUE='\033[38;5;27m'
TIME_VORTEX='\033[38;5;87m'
DALEK_RED='\033[38;5;196m'
CYBERMAN_YELLOW='\033[38;5;226m'
SONIC_GREEN='\033[38;5;40m'
GALLIFREYAN_GOLD='\033[38;5;220m'
NC='\033[0m'

# TARDIS Console Spinners
TARDIS_SPINNER=("🚀" "🌌" "⏳" "🌀" "🔮")

# ============================================
# TARDIS Console Functions
# ============================================

tardis_header() {
    clear
    echo -e "${TARDIS_BLUE}
  _____ _____ _____ _____ _____ _____ _____ _____
 |_   _|_   _|_   _|_   _|_   _|_   _|_   _|_   _|
   | |   | |   | |   | |   | |   | |   | |   | |
   | |   | |   | |   | |   | |   | |   | |   | |
  _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_
 |_____|_____|_____|_____|_____|_____|_____|_____|
${NC}"
    echo -e "${GALLIFREYAN_GOLD}Sonic-Screwdriver Installer for Linux${NC}"
    echo -e "${TIME_VORTEX}Doctor: $DOCTOR_NUMBER | Companion: $COMPANION${NC}"
    echo -e "${TARDIS_BLUE}TARDIS: $TARDIS_NAME${NC}"
    echo -e "${SONIC_GREEN}Sonic Screwdriver: ACTIVE${NC}"
    echo ""
}

doctor_says() {
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} $1"
}

companion_says() {
    echo -e "${CYBERMAN_YELLOW}[COMPANION]${NC} $1"
}

tardis_system() {
    echo -e "${TARDIS_BLUE}[TARDIS]${NC} $1"
}

sonic_action() {
    echo -e "${SONIC_GREEN}[SONIC]${NC} $1 *zzzzzap*"
}

dalek_error() {
    echo -e "${DALEK_RED}[DALEK]${NC} EXTERMINATE! $1"
}

# ============================================
# Progress Tracking
# ============================================

init_progress() {
    echo "Installation started: $(date)" > "$STATE_FILE"
    echo "Status: initialized" >> "$STATE_FILE"
    
    echo "========================================" > "$LOG_FILE"
    echo "Sonic-Screwdriver TARDIS Installer" >> "$LOG_FILE"
    echo "Started: $(date)" >> "$LOG_FILE"
    echo "========================================" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

update_progress() {
    local step_name="$1"
    echo "Current step: $step_name" >> "$STATE_FILE"
    echo "Timestamp: $(date)" >> "$STATE_FILE"
    echo "Status: started" >> "$STATE_FILE"
    
    echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} $step_name..."
    echo "[PROGRESS] $step_name" >> "$LOG_FILE"
}

mark_success() {
    echo "Status: completed" >> "$STATE_FILE"
    echo -e "${SONIC_GREEN}[SUCCESS]${NC} Done!"
    echo "[SUCCESS] $(date)" >> "$LOG_FILE"
}

mark_failure() {
    echo "Status: failed" >> "$STATE_FILE"
    echo -e "${DALEK_RED}[FAILURE]${NC}Failed!"
    echo "[FAILURE] $(date)" >> "$LOG_FILE"
}

# ============================================
# Robust Command Execution
# ============================================

run_safe() {
    local step_name="$1"
    local timeout=${2:-$STEP_TIMEOUT}
    shift 2
    local command="$@"
    
    update_progress "$step_name"
    echo "[CMD] $command" >> "$LOG_FILE"
    
    # Start command
    (eval "$command") &
    local pid=$!
    local start_time=$(date +%s)
    
    # Monitor loop
    while kill -0 "$pid" 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        # Show spinner
        local spinner_index=$(( elapsed % ${#TARDIS_SPINNER[@]} ))
        echo -ne "\r${TIME_VORTEX}${TARDIS_SPINNER[spinner_index]}${NC} $step_name... $elapsed seconds "
        
        # Timeout check
        if [ $elapsed -ge $timeout ]; then
            echo -ne "\r${DALEK_RED}[TIMEOUT]${NC} $timeout seconds!\n"
            echo "[TIMEOUT] Command timed out after $timeout seconds: $command" >> "$LOG_FILE"
            kill "$pid" 2>/dev/null || true
            sleep 1
            kill -9 "$pid" 2>/dev/null || true
            mark_failure
            return 1
        fi
        
        sleep $SPINNER_DELAY
    done
    
    # Wait for completion
    wait "$pid"
    local status=$?
    
    echo -ne "\r\033[K"
    
    if [ $status -eq 0 ]; then
        echo -e "${SONIC_GREEN}[SUCCESS]${NC} $step_name completed"
        mark_success
        return 0
    else
        echo -e "${DALEK_RED}[FAILURE]${NC} $step_name failed with status $status"
        echo "[FAILURE] Status: $status" >> "$LOG_FILE"
        mark_failure
        return $status
    fi
}

# ============================================
# Self-Healing
# ============================================

self_heal() {
    local attempt=$1
    local max_attempts=$2
    local step_name="$3"
    shift 3
    local command="$@"
    
    if [ $attempt -ge $max_attempts ]; then
        dalek_error "Maximum regeneration limit reached!"
        return 1
    fi
    
    doctor_says "Regenerating... (attempt $((attempt+1))/$max_attempts)"
    
    # Regeneration animation
    for i in {1..3}; do
        echo -ne "${TIME_VORTEX}🌀 ${NC}"
        for j in {1..3}; do echo -ne "."; sleep 0.3; done
        echo -ne "\r\033[K"
    done
    
    if eval "$command"; then
        sonic_action "Regeneration successful!"
        return 0
    else
        self_heal $((attempt+1)) $max_attempts "$step_name" "$command"
        return $?
    fi
}

# ============================================
# Installation Steps
# ============================================

step_system_check() {
    tardis_system "Checking system compatibility..."
    
    # Check for root
    if [ "$EUID" -eq 0 ]; then
        dalek_error "Do not run as root!"
        return 1
    fi
    
    # Check Ubuntu version
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        if [ "$ID" != "ubuntu" ]; then
            dalek_error "This installer requires Ubuntu!"
            return 1
        fi
        
        local version=$(echo "$VERSION_ID" | cut -d '.' -f 1)
        local minor=$(echo "$VERSION_ID" | cut -d '.' -f 2)
        
        if [ "$version" -lt 22 ] || ([ "$version" -eq 22 ] && [ "$minor" -lt 4 ]); then
            dalek_error "Ubuntu 22.04 LTS or later required!"
            return 1
        fi
    else
        dalek_error "Cannot determine OS!"
        return 1
    fi
    
    # Check architecture
    local arch=$(uname -m)
    if [ "$arch" != "x86_64" ] && [ "$arch" != "amd64" ]; then
        echo -e "${CYBERMAN_YELLOW}[WARNING]${NC} Unsupported architecture: $arch"
    fi
    
    sonic_action "System compatible"
    return 0
}

step_install_dependencies() {
    tardis_system "Installing timey-wimey dependencies..."
    
    local dependencies=("git" "make" "curl" "g++" "docker.io" "golang-go" "libssl-dev")
    local missing=()
    
    for dep in "${dependencies[@]}"; do
        if ! dpkg -l "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        doctor_says "Installing ${#missing[@]} missing dependencies..."
        sudo apt update
        sudo apt install -y "${missing[@]}"
    else
        sonic_action "All dependencies installed"
    fi
    
    # Add user to docker group
    sudo usermod -aG docker "$USER"
    
    # Check Go
    if ! command -v go &> /dev/null; then
        dalek_error "Go not installed!"
        return 1
    fi
    
    sonic_action "Dependencies ready"
    return 0
}

step_check_docker() {
    tardis_system "Checking Docker status..."
    
    if ! docker info &> /dev/null; then
        doctor_says "Starting Docker daemon..."
        sudo systemctl start docker
        sudo systemctl enable docker
        
        if ! docker info &> /dev/null; then
            dalek_error "Failed to start Docker!"
            return 1
        fi
    fi
    
    sonic_action "Docker online"
    return 0
}

step_clone_repository() {
    tardis_system "Materializing TARDIS..."
    
    if [ -d "$SONIC_DIR" ]; then
        doctor_says "TARDIS already materialized, updating..."
        cd "$SONIC_DIR"
        
        # Stash local changes if any
        if ! git diff --quiet || ! git diff --cached --quiet; then
            git stash push --include-untracked --message "TARDIS stash" 2>/dev/null || true
        fi
        
        git pull origin main 2>/dev/null || true
        
        # Apply stashed changes
        if git stash list | grep -q "TARDIS stash"; then
            git stash pop 2>/dev/null || true
        fi
    else
        doctor_says "Materializing TARDIS..."
        mkdir -p "$HOME/Code"
        git clone "$SONIC_REPO" "$SONIC_DIR"
        cd "$SONIC_DIR"
    fi
    
    sonic_action "TARDIS materialized"
    return 0
}

step_build_sonic() {
    tardis_system "Building Sonic Screwdriver..."
    
    cd "$SONIC_DIR"
    
    if [ ! -f "Makefile" ]; then
        dalek_error "Makefile not found!"
        return 1
    fi
    
    # Check if binary exists and is recent
    if [ -f "bin/sonic" ]; then
        doctor_says "Binary exists, checking if rebuild needed..."
        local makefile_time=$(stat -c %Y Makefile 2>/dev/null || date +%s)
        local binary_time=$(stat -c %Y bin/sonic 2>/dev/null || date +%s)
        
        if [ "$binary_time" -gt "$makefile_time" ]; then
            sonic_action "Using existing binary"
            return 0
        fi
    fi
    
    doctor_says "Building with make..."
    
    if ! make build; then
        dalek_error "Build failed!"
        return 1
    fi
    
    if [ ! -f "bin/sonic" ]; then
        dalek_error "Sonic binary not found!"
        return 1
    fi
    
    sonic_action "Sonic Screwdriver built"
    return 0
}

step_install_sonic() {
    tardis_system "Installing Sonic Screwdriver..."
    
    cd "$SONIC_DIR"
    
    # Install to user local first
    mkdir -p "$HOME/.local/bin"
    cp "bin/sonic" "$HOME/.local/bin/"
    chmod +x "$HOME/.local/bin/sonic"
    
    # Also install to system
    sudo cp "bin/sonic" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/sonic"
    
    sonic_action "Sonic Screwdriver installed"
    return 0
}

step_setup_environment() {
    tardis_system "Configuring TARDIS console..."
    
    # Create directories
    mkdir -p "$HOME/.sonic"
    mkdir -p "$HOME/.sonic/logs"
    mkdir -p "$HOME/.local/go"
    
    # Update bashrc
    if ! grep -q "export PATH=\"$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc"; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    fi
    
    # Setup correct Go workspace in .local
    if ! grep -q "export GOPATH=\"$HOME/.local/go\"" "$HOME/.bashrc"; then
        echo 'export GOPATH="$HOME/.local/go"' >> "$HOME/.bashrc"
        echo 'export PATH="$PATH:$GOPATH/bin"' >> "$HOME/.bashrc"
    fi
    
    # Check for old ~/go
    if [ -d "$HOME/go" ]; then
        echo -e "${CYBERMAN_YELLOW}[WARNING]${NC} Old Go workspace found at ~/go"
        echo "The correct location is ~/.local/go"
    fi
    
    # Source bashrc
    source "$HOME/.bashrc" 2>/dev/null || true
    
    sonic_action "TARDIS console configured"
    return 0
}

step_devstudio_integration() {
    tardis_system "Integrating with DevStudio..."
    
    # Check if DevStudio exists
    if [ ! -d "$HOME/Code/DevStudio" ]; then
        echo -e "${CYBERMAN_YELLOW}[INFO]${NC} DevStudio not found, skipping integration"
        return 0
    fi
    
    # Create project structure
    mkdir -p "$HOME/Code/DevStudio/Projects/SonicScrewdriver/config"
    
    # Create config
    cat > "$HOME/Code/DevStudio/Projects/SonicScrewdriver/config/config.yaml" << EOF
sonic:
  source_dir: "~/Code/SonicScrewdriver"
  build_dir: "/tmp/sonic-builds"
  install_path: "/usr/local/bin"
  go_path: "~/.local/go"

lechat:
  api_key: ""
  api_url: "https://api.lechat.pro"
  enabled: false

development:
  debug_mode: false
  log_level: "info"
  go_workspace: "~/.local/go"

devstudio:
  projects_dir: "~/Code/DevStudio/Projects"
  installed: true
  runtime_rules:
    go_workspace: "~/.local/go"
EOF
    
    # Create README
    cat > "$HOME/Code/DevStudio/Projects/SonicScrewdriver/README.md" << EOF
# Sonic-Screwdriver DevStudio Project

Sonic-Screwdriver integration with DevStudio.

## Configuration
`config/config.yaml` contains project settings.

## Usage
```bash
cd ~/Code/DevStudio
./install.sh
```
EOF
    
    sonic_action "DevStudio integration complete"
    return 0
}

step_verify_installation() {
    tardis_system "Running final diagnostic..."
    
    if ! command -v sonic &> /dev/null; then
        dalek_error "Sonic command not found!"
        return 1
    fi
    
    local version=$(sonic --version 2>/dev/null || echo "unknown")
    
    if [ "$version" = "unknown" ]; then
        dalek_error "Sonic command not working!"
        return 1
    fi
    
    # Test system check
    if sonic system check &> /dev/null; then
        echo "System compatibility: OK"
    else
        echo "System compatibility: WARNING"
    fi
    
    sonic_action "All systems nominal! Version: $version"
    return 0
}

# ============================================
# Main Storyline
# ============================================

main_install() {
    tardis_header
    doctor_says "Allons-y! Let's install the Sonic-Screwdriver!"
    companion_says "What's the plan, Doctor?"
    
    init_progress
    
    local steps=(
        "step_system_check:System Check"
        "step_install_dependencies:Install Dependencies"
        "step_check_docker:Check Docker"
        "step_clone_repository:Clone Repository"
        "step_build_sonic:Build Sonic"
        "step_install_sonic:Install Sonic"
        "step_setup_environment:Setup Environment"
        "step_devstudio_integration:DevStudio Integration"
        "step_verify_installation:Verify Installation"
    )
    
    local total_steps=${#steps[@]}
    
    for ((i=0; i<total_steps; i++)); do
        local step_func=$(echo "${steps[$i]}" | cut -d':' -f1)
        local step_name=$(echo "${steps[$i]}" | cut -d':' -f2)
        
        echo ""
        echo "========================================"
        echo "Step $(($i+1))/$total_steps: $step_name"
        echo "========================================"
        
        if ! $step_func; then
            dalek_error "Mission failed at step: $step_name"
            echo ""
            echo "Check $LOG_FILE for details"
            echo "You can resume by running this script again"
            return 1
        fi
    done
    
    # Clean up
    rm -f "$STATE_FILE"
    
    echo ""
    echo "========================================"
    sonic_action "Installation complete!"
    echo "========================================"
    
    doctor_says "There you go! Sonic-Screwdriver is ready!"
    
    echo ""
    echo "Try these commands:"
    echo "  sonic --help"
    echo "  sonic --version"
    echo "  sonic system check"
    echo ""
    echo "Log: $LOG_FILE"
    echo ""
    
    return 0
}

# ============================================
# Main Execution
# ============================================

COMMAND="${1:-install}"

case "$COMMAND" in
    install|run)
        main_install
        ;;
    quick)
        # Quick mode with less theme
        echo "Quick installation mode..."
        STEP_TIMEOUT=300
        SPINNER_DELAY=0.5
        main_install
        ;;
    help)
        echo "Sonic-Screwdriver TARDIS Installer for Linux"
        echo ""
        echo "Usage: $0 [install|quick|help]"
        echo ""
        echo "Commands:"
        echo "  install   - Full TARDIS experience (default)"
        echo "  quick     - Fast installation with minimal theme"
        echo "  help      - Show this help"
        ;;
    *)
        echo "Unknown command: $COMMAND"
        echo "Usage: $0 [install|quick|help]"
        exit 1
        ;;
esac
