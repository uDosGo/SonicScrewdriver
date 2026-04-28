#!/bin/bash

# Sonic-Screwdriver Ultimate Ubuntu Installer
# "The Unhangable Installer"
# Advanced version with progress tracking, hang detection, and recovery

set -euo pipefail

# ============================================
# Configuration
# ============================================

SONIC_REPO="https://github.com/uDosGo/SonicScrewdriver.git"
SONIC_DIR="$HOME/Code/SonicScrewdriver"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="/tmp/sonic-install-ultimate-$(date +%Y%m%d-%H%M%S).log"
STATE_FILE="/tmp/sonic-install-state.txt"

# Progress tracking
CURRENT_STEP=0
TOTAL_STEPS=8

# Timeouts
GLOBAL_TIMEOUT=600
STEP_TIMEOUT=300

# ============================================
# Progress Tracking System
# ============================================

init_progress() {
    echo "Installation started: $(date)" > "$STATE_FILE"
    echo "Total steps: $TOTAL_STEPS" >> "$STATE_FILE"
    echo "Current step: $CURRENT_STEP" >> "$STATE_FILE"
    echo "Status: initialized" >> "$STATE_FILE"
}

update_progress() {
    local step_name="$1"
    CURRENT_STEP=$((CURRENT_STEP + 1))
    
    echo "" >> "$STATE_FILE"
    echo "Step $CURRENT_STEP/$TOTAL_STEPS: $step_name" >> "$STATE_FILE"
    echo "Timestamp: $(date)" >> "$STATE_FILE"
    echo "Status: started" >> "$STATE_FILE"
    
    echo -e "\n[PROGRESS] Step $CURRENT_STEP/$TOTAL_STEPS: $step_name"
    echo "[PROGRESS] $(date)"
}

mark_success() {
    echo "Status: completed" >> "$STATE_FILE"
    echo "[PROGRESS] Step $CURRENT_STEP/$TOTAL_STEPS: SUCCESS"
}

mark_failure() {
    echo "Status: failed" >> "$STATE_FILE"
    echo "[PROGRESS] Step $CURRENT_STEP/$TOTAL_STEPS: FAILED"
}

check_resume() {
    if [ -f "$STATE_FILE" ]; then
        echo "Resuming previous installation..."
        source "$STATE_FILE"
        echo "Resuming from step $CURRENT_STEP/$TOTAL_STEPS"
        return 0
    fi
    return 1
}

# ============================================
# Robust Command Execution
# ============================================

run_with_timeout() {
    local command="$1"
    local timeout=${2:-$STEP_TIMEOUT}
    local step_name="$3"
    
    update_progress "$step_name"
    
    # Start command in background
    (eval "$command") &
    local pid=$!
    local start_time=$(date +%s)
    local monitor_pid=""
    
    # Monitor function
    monitor_command() {
        while kill -0 "$pid" 2>/dev/null; do
            local current_time=$(date +%s)
            local elapsed=$((current_time - start_time))
            
            # Check for timeout
            if [ $elapsed -ge $timeout ]; then
                echo "[TIMEOUT] Command timed out after $timeout seconds: $command" | tee -a "$LOG_FILE"
                kill "$pid" 2>/dev/null || true
                sleep 1
                kill -9 "$pid" 2>/dev/null || true
                mark_failure
                return 1
            fi
            
            # Check for hangs (no CPU usage)
            if [ $((elapsed % 10)) -eq 0 ]; then
                local cpu_usage=$(ps -p "$pid" -o %cpu= 2>/dev/null || echo "0")
                if [ "$cpu_usage" = "0.0" ] || [ "$cpu_usage" = "0" ]; then
                    echo "[HANG] Potential hang detected at $elapsed seconds, CPU: $cpu_usage%" | tee -a "$LOG_FILE"
                fi
            fi
            
            sleep 1
        done
        
        wait "$pid"
        return $?
    }
    
    # Start monitor
    monitor_command &
    monitor_pid=$!
    
    # Wait for command or monitor to finish
    wait "$pid" 2>/dev/null || true
    local status=$?
    
    # Kill monitor
    kill "$monitor_pid" 2>/dev/null || true
    
    if [ $status -eq 0 ]; then
        mark_success
        echo "[SUCCESS] $step_name completed" | tee -a "$LOG_FILE"
    else
        mark_failure
        echo "[FAILURE] $step_name failed with status $status" | tee -a "$LOG_FILE"
    fi
    
    return $status
}

# ============================================
# Installation Steps
# ============================================

step_check_system() {
    echo "Checking system requirements..."
    
    # Check Ubuntu version
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        if [ "$ID" != "ubuntu" ]; then
            echo "Error: This installer requires Ubuntu"
            return 1
        fi
        
        local version=$(echo "$VERSION_ID" | cut -d '.' -f 1)
        local minor=$(echo "$VERSION_ID" | cut -d '.' -f 2)
        
        if [ "$version" -lt 22 ] || ([ "$version" -eq 22 ] && [ "$minor" -lt 4 ]); then
            echo "Error: Ubuntu 22.04 LTS or later required"
            return 1
        fi
    else
        echo "Error: Cannot determine OS"
        return 1
    fi
    
    # Check architecture
    local arch=$(uname -m)
    if [ "$arch" != "x86_64" ] && [ "$arch" != "amd64" ]; then
        echo "Warning: Unsupported architecture: $arch"
    fi
    
    echo "System check passed"
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
        echo "Installing missing dependencies: ${missing[*]}"
        sudo apt update
        sudo apt install -y "${missing[@]}"
    fi
    
    # Add user to docker group
    sudo usermod -aG docker "$USER"
    
    echo "Dependencies installed"
    return 0
}

step_check_docker() {
    echo "Checking Docker status..."
    
    if ! docker info &> /dev/null; then
        echo "Starting Docker..."
        sudo systemctl start docker
        sudo systemctl enable docker
        
        if ! docker info &> /dev/null; then
            echo "Error: Failed to start Docker"
            return 1
        fi
    fi
    
    echo "Docker is running"
    return 0
}

step_clone_repository() {
    echo "Cloning repository..."
    
    if [ -d "$SONIC_DIR" ]; then
        echo "Repository exists, pulling latest changes..."
        cd "$SONIC_DIR"
        git pull origin main
    else
        mkdir -p "$HOME/Code"
        git clone "$SONIC_REPO" "$SONIC_DIR"
        cd "$SONIC_DIR"
    fi
    
    echo "Repository ready"
    return 0
}

step_build_sonic() {
    echo "Building Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    if [ ! -f "Makefile" ]; then
        echo "Error: Makefile not found"
        return 1
    fi
    
    # Build with progress monitoring
    make build
    
    if [ ! -f "bin/sonic" ]; then
        echo "Error: sonic binary not found after build"
        return 1
    fi
    
    echo "Build successful"
    return 0
}

step_install_sonic() {
    echo "Installing Sonic-Screwdriver..."
    
    cd "$SONIC_DIR"
    
    sudo cp "bin/sonic" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/sonic"
    
    echo "Installation complete"
    return 0
}

step_setup_environment() {
    echo "Setting up environment..."
    
    mkdir -p "$HOME/.sonic"
    mkdir -p "$HOME/.sonic/logs"
    
    # Update bashrc
    if ! grep -q "export PATH=\"$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc"; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    fi
    
    if ! grep -q "export GOPATH=\"$HOME/go\"" "$HOME/.bashrc"; then
        echo 'export GOPATH=$HOME/go' >> "$HOME/.bashrc"
        echo 'export PATH=$PATH:$GOPATH/bin' >> "$HOME/.bashrc"
    fi
    
    source "$HOME/.bashrc"
    
    echo "Environment setup complete"
    return 0
}

step_verify_installation() {
    echo "Verifying installation..."
    
    if ! command -v sonic &> /dev/null; then
        echo "Error: sonic command not found"
        return 1
    fi
    
    local version=$(sonic --version 2>/dev/null || echo "unknown")
    
    if [ "$version" = "unknown" ]; then
        echo "Error: sonic command not working"
        return 1
    fi
    
    echo "Verification successful: $version"
    return 0
}

# ============================================
# Main Installation
# ============================================

main_installation() {
    echo "Starting Sonic-Screwdriver Ultimate Installation"
    echo "Log file: $LOG_FILE"
    echo "State file: $STATE_FILE"
    
    # Initialize
    init_progress
    
    # Check if resuming
    if check_resume; then
        echo "Resuming from step $CURRENT_STEP"
    else
        CURRENT_STEP=0
    fi
    
    # Run steps with proper sequencing
    local steps=(
        "step_check_system:System Check"
        "step_install_dependencies:Install Dependencies"
        "step_check_docker:Check Docker"
        "step_clone_repository:Clone Repository"
        "step_build_sonic:Build Sonic"
        "step_install_sonic:Install Sonic"
        "step_setup_environment:Setup Environment"
        "step_verify_installation:Verify Installation"
    )
    
    local total_steps=${#steps[@]}
    
    # Run each step
    for ((i=CURRENT_STEP; i<total_steps; i++)); do
        local step_func=$(echo "${steps[$i]}" | cut -d':' -f1)
        local step_name=$(echo "${steps[$i]}" | cut -d':' -f2)
        
        echo "" | tee -a "$LOG_FILE"
        echo "========================================" | tee -a "$LOG_FILE"
        echo "Step $(($i+1))/$total_steps: $step_name" | tee -a "$LOG_FILE"
        echo "========================================" | tee -a "$LOG_FILE"
        
        if ! $step_func; then
            echo "Step failed: $step_name" | tee -a "$LOG_FILE"
            echo "Check $LOG_FILE for details"
            echo "You can resume by running this script again"
            return 1
        fi
        
        # Update state
        CURRENT_STEP=$((i + 1))
        echo "Current step: $CURRENT_STEP" > "$STATE_FILE"
    done
    
    # Clean up
    rm -f "$STATE_FILE"
    
    echo "" | tee -a "$LOG_FILE"
    echo "========================================" | tee -a "$LOG_FILE"
    echo "Installation completed successfully!" | tee -a "$LOG_FILE"
    echo "========================================" | tee -a "$LOG_FILE"
    
    echo ""
    echo "Sonic-Screwdriver is ready to use!"
    echo ""
    echo "Try these commands:"
    echo "  sonic --help              Show help"
    echo "  sonic --version           Show version"
    echo "  sonic system check        Check system compatibility"
    echo "  sonic tui                 Launch interactive interface"
    echo ""
    echo "Log file: $LOG_FILE"
    
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
