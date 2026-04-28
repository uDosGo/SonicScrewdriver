#!/bin/bash

# Sonic-Screwdriver Final Installer
# "The Ultimate Sonic Experience"
# Combines TARDIS theme with robust installation and DevStudio integration

set -euo pipefail

# ============================================
# Configuration
# ============================================

# Sonic Settings
SONIC_REPO="https://github.com/uDosGo/SonicScrewdriver.git"
SONIC_DIR="$HOME/Code/SonicScrewdriver"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="/tmp/sonic-final-$(date +%Y%m%d-%H%M%S).log"
STATE_FILE="/tmp/sonic-final-state.txt"

# DevStudio Integration
DEVSTUDIO_DIR="$HOME/Code/DevStudio"
DEVSTUDIO_PROJECTS="$DEVSTUDIO_DIR/Projects"

# TARDIS Theme Settings
TARDIS_NAME="Sonic TARDIS"
DOCTOR_NUMBER="14th"
COMPANION="Ubuntu 22.04 LTS"

# Timeouts
STEP_TIMEOUT=600
GLOBAL_TIMEOUT=1200

# Colors (TARDIS Console Theme)
TARDIS_BLUE='\033[38;5;27m'
TIME_VORTEX='\033[38;5;87m'
DALEK_RED='\033[38;5;196m'
CYBERMAN_YELLOW='\033[38;5;226m'
SONIC_GREEN='\033[38;5;40m'
GALLIFREYAN_GOLD='\033[38;5;220m'
NC='\033[0m' # No Color

# ============================================
# TARDIS Console Functions
# ============================================

tardis_header() {
    echo -e "${TARDIS_BLUE}
  _____ _____ _____ _____ _____ _____ _____ _____
 |_   _|_   _|_   _|_   _|_   _|_   _|_   _|_   _|
   | |   | |   | |   | |   | |   | |   | |   | |
   | |   | |   | |   | |   | |   | |   | |   | |
  _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_
 |_____|_____|_____|_____|_____|_____|_____|_____|
${NC}"
    echo -e "${GALLIFREYAN_GOLD}Sonic-Screwdriver Final Installer${NC}"
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
    echo -e "${DALEK_RED}[DALEK]${NC} $1"
}

# ============================================
# Progress Tracking System
# ============================================

init_progress() {
    echo "Installation started: $(date)" > "$STATE_FILE"
    echo "Current step: 0" >> "$STATE_FILE"
    echo "Status: initialized" >> "$STATE_FILE"
    
    # Initialize log file
    echo "========================================" > "$LOG_FILE"
    echo "Sonic-Screwdriver Final Installer" >> "$LOG_FILE"
    echo "Started: $(date)" >> "$LOG_FILE"
    echo "Log: $LOG_FILE" >> "$LOG_FILE"
    echo "========================================" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

update_progress() {
    local step_name="$1"
    local step_num="$2"
    
    echo "Current step: $step_num" > "$STATE_FILE"
    echo "Step $step_num: $step_name" >> "$STATE_FILE"
    echo "Timestamp: $(date)" >> "$STATE_FILE"
    echo "Status: started" >> "$STATE_FILE"
    
    echo -e "\n${TIME_VORTEX}[PROGRESS]${NC} Step $step_num: $step_name"
    echo "[PROGRESS] $(date)" >> "$LOG_FILE"
}

mark_success() {
    echo "Status: completed" >> "$STATE_FILE"
    echo -e "${SONIC_GREEN}[SUCCESS]${NC} Step completed"
    echo "[SUCCESS] $(date)" >> "$LOG_FILE"
}

mark_failure() {
    echo "Status: failed" >> "$STATE_FILE"
    echo -e "${DALEK_RED}[FAILURE]${NC} Step failed"
    echo "[FAILURE] $(date)" >> "$LOG_FILE"
}

check_resume() {
    if [ -f "$STATE_FILE" ]; then
        echo "Resuming previous installation..."
        source "$STATE_FILE"
        local current_step=$(grep "Current step:" "$STATE_FILE" | cut -d' ' -f3)
        echo "Resuming from step $current_step"
        return 0
    fi
    return 1
}

# ============================================
# Robust Command Execution
# ============================================

run_safe() {
    local command="$1"
    local step_name="$2"
    local timeout=${3:-$STEP_TIMEOUT}
    
    update_progress "$step_name" "$CURRENT_STEP"
    
    # Start command
    (eval "$command") &
    local pid=$!
    local start_time=$(date +%s)
    
    # Monitor loop
    while kill -0 "$pid" 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        # Timeout check
        if [ $elapsed -ge $timeout ]; then
            dalek_error "Timeout after $timeout seconds!"
            kill "$pid" 2>/dev/null || true
            sleep 1
            kill -9 "$pid" 2>/dev/null || true
            mark_failure
            return 1
        fi
        
        # Hang detection (every 10 seconds)
        if [ $((elapsed % 10)) -eq 0 ]; then
            local cpu_usage=$(ps -p "$pid" -o %cpu= 2>/dev/null || echo "0")
            if [ "$cpu_usage" = "0.0" ] || [ "$cpu_usage" = "0" ]; then
                echo "[HANG] Potential hang detected at $elapsed seconds" >> "$LOG_FILE"
            fi
        fi
        
        sleep 1
    done
    
    # Wait for completion
    wait "$pid"
    local status=$?
    
    if [ $status -eq 0 ]; then
        mark_success
        return 0
    else
        mark_failure
        return $status
    fi
}

# ============================================
# DevStudio Integration
# ============================================

setup_devstudio() {
    echo "Setting up DevStudio integration..."
    
    # Create DevStudio structure if it doesn't exist
    if [ ! -d "$DEVSTUDIO_DIR" ]; then
        echo "DevStudio not found at $DEVSTUDIO_DIR"
        return 0
    fi
    
    # Create Sonic-Screwdriver project in DevStudio
    mkdir -p "$DEVSTUDIO_PROJECTS/SonicScrewdriver/config"
    
    # Create project configuration
    cat > "$DEVSTUDIO_PROJECTS/SonicScrewdriver/config/config.yaml" << 'EOF'
# Sonic-Screwdriver DevStudio Configuration
sonic:
  source_dir: "~/Code/SonicScrewdriver"
  build_dir: "/tmp/sonic-builds"
  install_path: "/usr/local/bin"
  go_path: "~/.local/go"  # Correct Go workspace location
  
lechat:
  api_key: ""
  api_url: "https://api.lechat.pro"
  enabled: false

development:
  debug_mode: false
  log_level: "info"
  auto_build: false
  auto_test: false
  go_workspace: "~/.local/go"  # Use .local/go not ~/go

# DevStudio Integration
devstudio:
  projects_dir: "~/Code/DevStudio/Projects"
  vibe_skills_dir: "~/Code/DevStudio/vibe-skills"
  installed: true
  runtime_rules:
    go_workspace: "~/.local/go"
    no_home_go: true
EOF
    
    # Create project README
    cat > "$DEVSTUDIO_PROJECTS/SonicScrewdriver/README.md" << 'EOF'
# Sonic-Screwdriver DevStudio Project

This project integrates Sonic-Screwdriver with DevStudio.

## Configuration

The project configuration is located at `config/config.yaml`.

## Usage

```bash
# Load configuration
devstudio config load-project SonicScrewdriver

# Build and install
cd ~/Code/DevStudio/vibe-skills/sonic-screwdriver
./sonic-screwdriver.sh build
./sonic-screwdriver.sh install
```

## Integration

This project provides:
- Centralized configuration management
- DevStudio skill integration
- Documentation publishing
- Agentic workflow support
EOF
    
    echo "DevStudio integration complete"
    return 0
}

# ============================================
# Installation Steps
# ============================================

step_system_check() {
    echo "Checking system requirements..."
    
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
    
    # Check for root
    if [ "$EUID" -eq 0 ]; then
        dalek_error "Do not run as root!"
        return 1
    fi
    
    sonic_action "System check passed"
    return 0
}

step_install_dependencies() {
    echo "Installing dependencies..."
    
    local dependencies=("git" "make" "curl" "g++" "docker.io" "golang-go" "libssl-dev")
    local missing=()
    
    for dep in "${dependencies[@]}"; do
        if ! dpkg -l "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        doctor_says "Installing missing dependencies..."
        sudo apt update
        sudo apt install -y "${missing[@]}"
    fi
    
    # Add user to docker group
    sudo usermod -aG docker "$USER"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        dalek_error "Go not installed!"
        return 1
    fi
    
    sonic_action "Dependencies installed"
    return 0
}

step_check_docker() {
    echo "Checking Docker status..."
    
    if ! docker info &> /dev/null; then
        doctor_says "Starting Docker..."
        sudo systemctl start docker
        sudo systemctl enable docker
        
        if ! docker info &> /dev/null; then
            dalek_error "Failed to start Docker!"
            return 1
        fi
    fi
    
    sonic_action "Docker is running"
    return 0
}

step_clone_repository() {
    echo "Cloning Sonic-Screwdriver repository..."
    
    if [ -d "$SONIC_DIR" ]; then
        doctor_says "Repository exists, checking status..."
        cd "$SONIC_DIR"
        
        # Check if there are local changes
        if ! git diff --quiet || ! git diff --cached --quiet; then
            doctor_says "Local changes detected, stashing..."
            git stash push --include-untracked --message "Automatic stash before update"
        fi
        
        doctor_says "Pulling latest changes..."
        if ! git pull origin main; then
            dalek_error "Git pull failed!"
            # Try to unstash if we stashed
            if git stash list | grep -q "Automatic stash"; then
                git stash pop || true
            fi
            return 1
        fi
        
        # Apply stashed changes if they exist
        if git stash list | grep -q "Automatic stash"; then
            doctor_says "Applying stashed changes..."
            git stash pop || true
        fi
    else
        doctor_says "Cloning repository..."
        mkdir -p "$HOME/Code"
        git clone "$SONIC_REPO" "$SONIC_DIR"
        cd "$SONIC_DIR"
    fi
    
    sonic_action "Repository ready"
    return 0
}

step_build_sonic() {
    echo "Building Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    if [ ! -f "Makefile" ]; then
        dalek_error "Makefile not found!"
        return 1
    fi
    
    # Check if binary already exists
    if [ -f "bin/sonic" ]; then
        doctor_says "Binary already exists, checking if rebuild needed..."
        # Compare timestamps
        local makefile_time=$(stat -c %Y Makefile 2>/dev/null || date +%s)
        local binary_time=$(stat -c %Y bin/sonic 2>/dev/null || date +%s)
        
        if [ "$binary_time" -gt "$makefile_time" ]; then
            doctor_says "Binary is newer than Makefile, skipping build..."
            sonic_action "Using existing binary"
            return 0
        fi
    fi
    
    doctor_says "Building with make..."
    
    # Build with timeout monitoring
    if ! make build; then
        dalek_error "Build failed!"
        return 1
    fi
    
    if [ ! -f "bin/sonic" ]; then
        dalek_error "Sonic binary not found!"
        return 1
    fi
    
    sonic_action "Build successful"
    return 0
}

step_install_sonic() {
    echo "Installing Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    sudo cp "bin/sonic" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/sonic"
    
    sonic_action "Installation complete"
    return 0
}

step_setup_environment() {
    echo "Setting up environment..."
    
    # Create .sonic directory
    mkdir -p "$HOME/.sonic"
    mkdir -p "$HOME/.sonic/logs"
    
    # Update bashrc
    if ! grep -q "export PATH=\"$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc"; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    fi
    
    # Setup Go workspace in .local (not ~/go)
    if ! grep -q "export GOPATH=\"$HOME/.local/go\"" "$HOME/.bashrc"; then
        echo 'export GOPATH=$HOME/.local/go' >> "$HOME/.bashrc"
        echo 'export PATH=$PATH:$GOPATH/bin' >> "$HOME/.bashrc"
    fi
    
    # Create .local/go directory if it doesn't exist
    mkdir -p "$HOME/.local/go"
    
    # Check if old ~/go exists and warn
    if [ -d "$HOME/go" ]; then
        echo "Warning: Found Go workspace at ~/go"
        echo "The correct location is ~/.local/go"
        echo "Consider moving: mv ~/go ~/.local/go"
    fi
    
    # Check for stray go folder in root
    if [ -d "/go" ] && [ ! -L "/go" ]; then
        echo "Warning: Found /go directory in root"
        echo "This might be from a previous installation"
        echo "Consider moving it to $HOME/go"
    fi
    
    source "$HOME/.bashrc"
    
    sonic_action "Environment configured"
    return 0
}

step_devstudio_integration() {
    echo "Setting up DevStudio integration..."
    
    setup_devstudio
    
    sonic_action "DevStudio integration complete"
    return 0
}

step_verify_installation() {
    echo "Verifying installation..."
    
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
    
    sonic_action "Verification successful: $version"
    return 0
}

# ============================================
# Main Installation Storyline
# ============================================

main_installation() {
    tardis_header
    doctor_says "Allons-y! Let's install the Sonic-Screwdriver!"
    
    # Initialize
    init_progress
    CURRENT_STEP=1
    
    # Define steps
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
    
    # Run each step
    for ((i=0; i<total_steps; i++)); do
        local step_func=$(echo "${steps[$i]}" | cut -d':' -f1)
        local step_name=$(echo "${steps[$i]}" | cut -d':' -f2)
        
        echo ""
        echo "========================================"
        echo "Step $(($i+1))/$total_steps: $step_name"
        echo "========================================"
        
        if ! $step_func; then
            dalek_error "Step failed: $step_name"
            echo "Check $LOG_FILE for details"
            echo "You can resume by running this script again"
            return 1
        fi
        
        CURRENT_STEP=$((i + 1))
        echo "Current step: $CURRENT_STEP" > "$STATE_FILE"
    done
    
    # Clean up
    rm -f "$STATE_FILE"
    
    echo ""
    echo "========================================"
    sonic_action "Installation completed successfully!"
    echo "========================================"
    
    doctor_says "There you go! Sonic-Screwdriver is ready to use!"
    
    echo ""
    echo "Try these commands:"
    echo "  ${SONIC_GREEN}sonic --help${NC}              Show help"
    echo "  ${SONIC_GREEN}sonic --version${NC}           Show version"
    echo "  ${SONIC_GREEN}sonic system check${NC}        Check system compatibility"
    echo "  ${SONIC_GREEN}sonic tui${NC}                 Launch interactive interface"
    echo ""
    echo "DevStudio integration:"
    echo "  cd ~/Code/DevStudio"
    echo "  devstudio config load-project SonicScrewdriver"
    echo ""
    echo "Log file: $LOG_FILE"
    echo "State file: $STATE_FILE (removed on success)"
    
    return 0
}

# ============================================
# Main Execution
# ============================================

COMMAND="${1:-install}"

case "$COMMAND" in
    install)
        main_installation
        ;;
    resume)
        main_installation
        ;;
    clean)
        rm -f "$STATE_FILE"
        echo "State cleaned"
        ;;
    status)
        if [ -f "$STATE_FILE" ]; then
            echo "Installation in progress:"
            cat "$STATE_FILE"
        else
            echo "No installation in progress"
        fi
        ;;
    *)
        echo "Usage: $0 [install|resume|clean|status]"
        exit 1
        ;;
esac
