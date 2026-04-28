#!/bin/bash

# Sonic-Screwdriver Ubuntu Installer - TARDIS Edition
# "The Doctor's Ultimate Installation System"
# A time-traveling, self-healing installer with Whovian flair

set -euo pipefail

# ============================================
# TARDIS Configuration
# ============================================

# TARDIS Settings
TARDIS_NAME="Sonic TARDIS"
TARDIS_VERSION="Mark VII"
DOCTOR_NUMBER="14th"
COMPANION="Ubuntu 22.04 LTS"

# Installation Settings
SONIC_REPO="https://github.com/uDosGo/SonicScrewdriver.git"
SONIC_DIR="$HOME/Code/SonicScrewdriver"
INSTALL_DIR="/usr/local/bin"
LOG_FILE="/tmp/sonic-tardis-$(date +%Y%m%d-%H%M%S).log"

# Timey-wimey settings
MAX_RETRIES=3
TIMEOUT_SECONDS=300
SPINNER_DELAY=0.1

# ============================================
# TARDIS Core Systems
# ============================================

# Colors (TARDIS Console Theme)
TARDIS_BLUE='\033[38;5;27m'
TIME_VORTEX='\033[38;5;87m'
DALEK_RED='\033[38;5;196m'
CYBERMAN_YELLOW='\033[38;5;226m'
SONIC_GREEN='\033[38;5;40m'
GALLIFREYAN_GOLD='\033[38;5;220m'
NC='\033[0m' # No Color

# TARDIS Console Spinners
SPINNER_FRAMES=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
TARDIS_SPINNER=("🚀" "🌌" "⏳" "🌀" "🔮" "🌀" "⏳" "🌌")

# ============================================
# TARDIS Console Functions
# ============================================

# Initialize TARDIS systems
init_tardis() {
    echo -e "${TARDIS_BLUE}
  _____ _____ _____ _____ _____ _____ _____ _____
 |_   _|_   _|_   _|_   _|_   _|_   _|_   _|_   _|
   | |   | |   | |   | |   | |   | |   | |   | |
   | |   | |   | |   | |   | |   | |   | |   | |
  _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_ _| |_
 |_____|_____|_____|_____|_____|_____|_____|_____|
${NC}"
    
    echo -e "${GALLIFREYAN_GOLD}TARDIS Console Initialized${NC}"
    echo -e "${TIME_VORTEX}Doctor: $DOCTOR_NUMBER | Companion: $COMPANION${NC}"
    echo -e "${TARDIS_BLUE}TARDIS: $TARDIS_NAME $TARDIS_VERSION${NC}"
    echo -e "${SONIC_GREEN}Sonic Screwdriver: ACTIVE${NC}"
    echo ""
    
    # Start TARDIS log with system information
    echo "========================================" > "$LOG_FILE"
    echo "TARDIS Log - $(date)" >> "$LOG_FILE"
    echo "========================================" >> "$LOG_FILE"
    echo "Doctor: $DOCTOR_NUMBER" >> "$LOG_FILE"
    echo "Mission: Install Sonic-Screwdriver on $COMPANION" >> "$LOG_FILE"
    echo "TARDIS: $TARDIS_NAME $TARDIS_VERSION" >> "$LOG_FILE"
    echo "Log File: $LOG_FILE" >> "$LOG_FILE"
    echo "========================================" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    
    # Log system information
    echo "=== SYSTEM INFORMATION ===" >> "$LOG_FILE"
    echo "User: $(whoami)" >> "$LOG_FILE"
    echo "Hostname: $(hostname)" >> "$LOG_FILE"
    echo "Date: $(date)" >> "$LOG_FILE"
    echo "PWD: $(pwd)" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    
    # Create restart marker file
    RESTART_FILE="/tmp/sonic-tardis-restart.marker"
    if [ -f "$RESTART_FILE" ]; then
        echo "[RESTART] TARDIS restart detected" >> "$LOG_FILE"
        rm -f "$RESTART_FILE"
    fi
}

# TARDIS Console Output
tardis_echo() {
    local message="$1"
    local color="$2"
    local log_level="$3"
    
    echo -e "${color}[TARDIS]${NC} $message" | tee -a "$LOG_FILE"
    
    if [ "$log_level" = "ERROR" ]; then
        echo -e "${DALEK_RED}[DALEK DETECTED]${NC} $message" >> "$LOG_FILE"
    fi
}

# Doctor's Announcement
doctor_says() {
    local message="$1"
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} $message" | tee -a "$LOG_FILE"
}

# Companion Response
companion_says() {
    local message="$1"
    echo -e "${CYBERMAN_YELLOW}[COMPANION]${NC} $message" | tee -a "$LOG_FILE"
}

# TARDIS System Message
tardis_system() {
    local message="$1"
    echo -e "${TARDIS_BLUE}[TARDIS]${NC} $message" | tee -a "$LOG_FILE"
}

# Sonic Screwdriver Action
sonic_action() {
    local message="$1"
    echo -e "${SONIC_GREEN}[SONIC]${NC} $message" | tee -a "$LOG_FILE"
}

# ============================================
# Time Vortex Spinners & Progress
# ============================================

# Start time vortex spinner
start_spinner() {
    local pid=$1
    local message="$2"
    local spinner_index=0
    
    # Show initial message
    echo -ne "${TIME_VORTEX}⏳${NC} $message... "
    
    # Spin in the time vortex
    while kill -0 "$pid" 2>/dev/null; do
        spinner_index=$(( (spinner_index + 1) % ${#TARDIS_SPINNER[@]} ))
        echo -ne "\b\b\b${TIME_VORTEX}${TARDIS_SPINNER[spinner_index]}${NC}"
        sleep "$SPINNER_DELAY"
    done
    
    # Clean up spinner
    echo -ne "\b\b\b"
    wait "$pid"
    local status=$?
    
    if [ $status -eq 0 ]; then
        echo -e "${SONIC_GREEN}✓${NC}"
    else
        echo -e "${DALEK_RED}✗${NC}"
    fi
    
    return $status
}

# Run command with time vortex spinner
run_with_spinner() {
    local command="$1"
    local message="$2"
    
    # Start command in background
    (eval "$command") &
    local pid=$!
    
    # Start spinner
    start_spinner "$pid" "$message"
    
    return $?
}

# Progress counter with time vortex effect
show_progress() {
    local current=$1
    local total=$2
    local message="$3"
    
    local percent=$((current * 100 / total))
    local done=$((percent / 2))
    local left=$((50 - done))
    
    # Time vortex progress bar
    echo -ne "${TIME_VORTEX}[${NC}"
    for ((i=0; i<done; i++)); do echo -ne "${TIME_VORTEX}▰${NC}"; done
    for ((i=0; i<left; i++)); do echo -ne " "; done
    echo -ne "${TIME_VORTEX}] ${percent}%${NC} - $message\r"
}

# ============================================
# TARDIS Error Handling & Self-Healing
# ============================================

# TARDIS Error Detection
tardis_error_detected() {
    local message="$1"
    echo -e "${DALEK_RED}
  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
  ███████████████████████████████████████
  ██${NC} ERROR DETECTED ${DALEK_RED}██
  ██${NC} Timey-wimey malfunction! ${DALEK_RED}██
  ███████████████████████████████████████
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀${NC}"
    echo -e "${DALEK_RED}[DALEK]${NC} $message"
    echo -e "${DALEK_RED}[DALEK]${NC} EXTERMINATE! EXTERMINATE!"
    
    # Log error with debugging info
    echo "" >> "$LOG_FILE"
    echo "===== CRITICAL ERROR =====" >> "$LOG_FILE"
    echo "Error: $message" >> "$LOG_FILE"
    echo "Timestamp: $(date)" >> "$LOG_FILE"
    echo "Log File: $LOG_FILE" >> "$LOG_FILE"
    echo "=========================" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    
    # Create restart script
    create_restart_script
}

# TARDIS Self-Healing with enhanced debugging
tardis_self_heal() {
    local attempt=$1
    local max_attempts=$2
    local command="$3"
    local message="$4"
    
    if [ $attempt -ge $max_attempts ]; then
        tardis_error_detected "Maximum regeneration limit reached!"
        echo "[REGENERATION] Failed after $max_attempts attempts" >> "$LOG_FILE"
        return 1
    fi
    
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Oh, that's not good. Let me try regenerating..."
    echo -e "${TIME_VORTEX}[TARDIS]${NC} Initiating regeneration sequence $((attempt + 1))/$max_attempts..."
    echo "[REGENERATION] Attempt $((attempt + 1))/$max_attempts for: $message" >> "$LOG_FILE"
    
    # Regeneration animation
    for i in {1..3}; do
        echo -ne "${TIME_VORTEX}🌀 Regenerating${NC}"
        for j in {1..3}; do echo -ne "."; sleep 0.5; done
        echo -ne "\r"
    done
    echo -ne "\033[K"
    
    # Try again with monitoring
    if run_with_monitoring "$command" "$message" $TIMEOUT_SECONDS; then
        echo -e "${SONIC_GREEN}[SONIC]${NC} Regeneration successful!"
        echo "[REGENERATION] Success on attempt $((attempt + 1))" >> "$LOG_FILE"
        return 0
    else
        echo "[REGENERATION] Attempt $((attempt + 1)) failed" >> "$LOG_FILE"
        tardis_self_heal $((attempt + 1)) "$max_attempts" "$command" "$message"
        return $?
    fi
}

# TARDIS Timeout Monitor with enhanced debugging
tardis_monitor() {
    local pid=$1
    local timeout=$2
    local message="$3"
    
    local start_time=$(date +%s)
    local last_check=$(date +%s)
    local check_interval=5
    
    # Log timeout monitor start
    echo "[TIMEOUT_MONITOR] Started for PID $pid, timeout: $timeout seconds, message: $message" >> "$LOG_FILE"
    
    while kill -0 "$pid" 2>/dev/null; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        # Check for timeout
        if [ $elapsed -ge $timeout ]; then
            echo -e "${DALEK_RED}
[TIMEOUT]${NC} The Doctor is taking too long! ($timeout seconds)"
            echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Oh, timey-wimey stuff got stuck in the time vortex!"
            
            # Enhanced debugging output
            echo "[TIMEOUT] Process $pid timed out after $timeout seconds" >> "$LOG_FILE"
            echo "[TIMEOUT] Last known status: $message" >> "$LOG_FILE"
            
            # Check process status
            if ps -p "$pid" -o cmd= >> "$LOG_FILE" 2>&1; then
                echo "[TIMEOUT] Process command logged" >> "$LOG_FILE"
            fi
            
            # Terminate the process
            kill "$pid" 2>/dev/null || true
            
            # Force kill if still running
            sleep 2
            if kill -0 "$pid" 2>/dev/null; then
                echo "[TIMEOUT] Process $pid still running, force killing..." >> "$LOG_FILE"
                kill -9 "$pid" 2>/dev/null || true
            fi
            
            return 1
        fi
        
        # Periodic status check (every 5 seconds)
        if [ $((current_time - last_check)) -ge $check_interval ]; then
            echo "[TIMEOUT_MONITOR] PID $pid still running, elapsed: $elapsed seconds" >> "$LOG_FILE"
            last_check=$current_time
        fi
        
        sleep 1
    done
    
    wait "$pid" 2>/dev/null || true
    local status=$?
    echo "[TIMEOUT_MONITOR] Process $pid completed with status $status" >> "$LOG_FILE"
    return $status
}

# Run command with TARDIS monitoring and enhanced error handling
run_with_monitoring() {
    local command="$1"
    local message="$2"
    local timeout=${3:-$TIMEOUT_SECONDS}
    
    echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} Entering time vortex for: $message"
    echo "[TIME_VORTEX] Starting: $command" >> "$LOG_FILE"
    
    # Start command
    (eval "$command") &
    local pid=$!
    
    # Log process start
    echo "[TIME_VORTEX] Command started with PID $pid" >> "$LOG_FILE"
    
    # Monitor in background
    tardis_monitor "$pid" "$timeout" "$message" &
    local monitor_pid=$!
    
    # Wait for command to complete
    wait "$pid" 2>/dev/null || true
    local status=$?
    
    # Stop monitor
    kill "$monitor_pid" 2>/dev/null || true
    
    # Check if process is still running (hang detection)
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "${DALEK_RED}[DALEK]${NC} Process hanging detected! Terminating..."
        echo "[HANG_DETECTED] Process $pid still running after wait, terminating..." >> "$LOG_FILE"
        kill "$pid" 2>/dev/null || true
        sleep 1
        if kill -0 "$pid" 2>/dev/null; then
            kill -9 "$pid" 2>/dev/null || true
        fi
        status=1
    fi
    
    if [ $status -eq 0 ]; then
        echo -e "${SONIC_GREEN}[SONIC]${NC} Time vortex journey complete!"
        echo "[TIME_VORTEX] Success: $message" >> "$LOG_FILE"
    else
        echo -e "${DALEK_RED}[DALEK]${NC} Time vortex disturbance detected!"
        echo "[TIME_VORTEX] Failed: $message with status $status" >> "$LOG_FILE"
        
        # Provide debugging information
        echo "" >> "$LOG_FILE"
        echo "===== DEBUGGING INFORMATION =====" >> "$LOG_FILE"
        echo "Command: $command" >> "$LOG_FILE"
        echo "Exit Status: $status" >> "$LOG_FILE"
        echo "Timeout: $timeout seconds" >> "$LOG_FILE"
        echo "=================================" >> "$LOG_FILE"
        echo "" >> "$LOG_FILE"
    fi
    
    return $status
}

# ============================================
# TARDIS System Checks
# ============================================

# Check if running as root (The Doctor doesn't like root!)
check_not_root() {
    if [ "$EUID" -eq 0 ]; then
        echo -e "${DALEK_RED}
  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
  ████████████████████████████████████████
  ██${NC} ROOT DETECTED! ${DALEK_RED}██
  ██${NC} The Doctor refuses to work as root! ${DALEK_RED}██
  ████████████████████████████████████████
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀${NC}"
        echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} I'm the Doctor, not the root user! Run this as a regular user."
        exit 1
    fi
}

# Check Ubuntu version (Companion compatibility)
check_ubuntu_version() {
    tardis_system "Checking companion compatibility..."
    
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        if [ "$ID" != "ubuntu" ]; then
            echo -e "${DALEK_RED}[DALEK]${NC} Unrecognized companion: $ID"
            echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} This TARDIS only works with Ubuntu companions!"
            exit 1
        fi
        
        local version=$(echo "$VERSION_ID" | cut -d '.' -f 1)
        local minor_version=$(echo "$VERSION_ID" | cut -d '.' -f 2)
        
        if [ "$version" -lt 22 ] || ([ "$version" -eq 22 ] && [ "$minor_version" -lt 4 ]); then
            echo -e "${DALEK_RED}[DALEK]${NC} Incompatible companion version: $VERSION_ID"
            echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Need Ubuntu 22.04 LTS or later!"
            exit 1
        fi
        
        echo -e "${SONIC_GREEN}[SONIC]${NC} Companion $VERSION_ID - compatible!"
    else
        echo -e "${DALEK_RED}[DALEK]${NC} Companion identification failed!"
        exit 1
    fi
}

# Check architecture (TARDIS dimensional stability)
check_architecture() {
    tardis_system "Checking TARDIS dimensional stability..."
    
    local arch=$(uname -m)
    if [ "$arch" != "x86_64" ] && [ "$arch" != "amd64" ]; then
        echo -e "${CYBERMAN_YELLOW}[CYBERMAN]${NC} Unsupported architecture: $arch"
        echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} The TARDIS prefers x86_64/amd64 dimensions!"
    else
        echo -e "${TARDIS_BLUE}[TARDIS]${NC} Dimensional stability confirmed: $arch"
    fi
}

# ============================================
# TARDIS Installation Functions
# ============================================

# Install dependencies (TARDIS fuel)
install_dependencies() {
    tardis_system "Refueling TARDIS with essential timey-wimey stuff..."
    
    local dependencies=("git" "make" "curl" "g++" "docker" "go")
    local missing=()
    
    for dep in "${dependencies[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -eq 0 ]; then
        echo -e "${SONIC_GREEN}[SONIC]${NC} All TARDIS fuel components present!"
        return 0
    fi
    
    echo -e "${CYBERMAN_YELLOW}[CYBERMAN]${NC} Missing TARDIS components: ${missing[*]}"
    doctor_says "Don't worry, I'll sort that out!"
    
    # Install with progress
    echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} Entering installation vortex..."
    
    sudo apt update
    for ((i=0; i<=100; i+=10)); do
        show_progress "$i" "100" "Updating package lists"
        sleep 0.2
    done
    echo ""
    
    sudo apt install -y git make curl g++ docker.io golang-go libssl-dev
    
    # Add user to docker group
    sudo usermod -aG docker "$USER"
    
    echo -e "${SONIC_GREEN}[SONIC]${NC} TARDIS refueled successfully!"
    companion_says "Brilliant! What's next?"
}

# Clone repository (Materializing TARDIS)
clone_repository() {
    tardis_system "Materializing TARDIS in local spacetime..."
    
    if [ -d "$SONIC_DIR" ]; then
        doctor_says "TARDIS already materialized! Updating temporal coordinates..."
        cd "$SONIC_DIR"
        run_with_spinner "git pull origin main" "Updating TARDIS coordinates"
    else
        doctor_says "Materializing TARDIS at $SONIC_DIR..."
        mkdir -p "$HOME/Code"
        run_with_spinner "git clone \"$SONIC_REPO\" \"$SONIC_DIR\"" "Materializing TARDIS"
        cd "$SONIC_DIR"
    fi
    
    echo -e "${TARDIS_BLUE}[TARDIS]${NC} Materialization complete!"
}

# Build Sonic-Screwdriver (Activating Sonic)
build_sonic() {
    tardis_system "Activating Sonic Screwdriver..."
    
    if [ ! -f "Makefile" ]; then
        echo -e "${DALEK_RED}[DALEK]${NC} TARDIS blueprint not found!"
        exit 1
    fi
    
    doctor_says "Time to build the Sonic Screwdriver!"
    echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} Sonic activation sequence initiated..."
    
    # Build with progress monitoring
    if ! run_with_monitoring "make build" "Building Sonic Screwdriver" 600; then
        echo -e "${DALEK_RED}[DALEK]${NC} Sonic activation failed!"
        return 1
    fi
    
    if [ ! -f "bin/sonic" ]; then
        echo -e "${DALEK_RED}[DALEK]${NC} Sonic Screwdriver not detected!"
        return 1
    fi
    
    echo -e "${SONIC_GREEN}[SONIC]${NC} Sonic Screwdriver activated! *zzzzzap*"
}

# Install Sonic-Screwdriver (Dematerializing to system)
install_sonic() {
    tardis_system "Dematerializing Sonic Screwdriver to system pathways..."
    
    if [ ! -f "bin/sonic" ]; then
        echo -e "${DALEK_RED}[DALEK]${NC} Sonic Screwdriver not found in TARDIS!"
        return 1
    fi
    
    doctor_says "Installing Sonic Screwdriver to $INSTALL_DIR..."
    
    run_with_spinner "sudo cp \"bin/sonic\" \"$INSTALL_DIR/\"" "Dematerializing Sonic"
    run_with_spinner "sudo chmod +x \"$INSTALL_DIR/sonic\"" "Activating system interface"
    
    echo -e "${TARDIS_BLUE}[TARDIS]${NC} Sonic Screwdriver installed to $INSTALL_DIR/sonic"
}

# Setup environment (TARDIS console configuration)
setup_environment() {
    tardis_system "Configuring TARDIS console environment..."
    
    # Create .sonic directory
    mkdir -p "$HOME/.sonic"
    mkdir -p "$HOME/.sonic/logs"
    
    # Configure environment
    if ! grep -q "export PATH=\"$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc"; then
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
        echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} Adding local bin to PATH..."
    fi
    
    if ! grep -q "export GOPATH=\"$HOME/go\"" "$HOME/.bashrc"; then
        echo 'export GOPATH=$HOME/go' >> "$HOME/.bashrc"
        echo 'export PATH=$PATH:$GOPATH/bin' >> "$HOME/.bashrc"
        echo -e "${TIME_VORTEX}[TIME VORTEX]${NC} Configuring Go environment..."
    fi
    
    # Source bashrc
    source "$HOME/.bashrc"
    
    echo -e "${SONIC_GREEN}[SONIC]${NC} TARDIS console configured!"
}

# Verify installation (TARDIS system diagnostic)
verify_installation() {
    tardis_system "Running TARDIS system diagnostic..."
    
    if ! command -v sonic &> /dev/null; then
        echo -e "${DALEK_RED}[DALEK]${NC} Sonic Screwdriver not found in system pathways!"
        return 1
    fi
    
    local version=$(sonic --version 2>/dev/null || echo "unknown")
    
    if [ "$version" != "unknown" ]; then
        echo -e "${SONIC_GREEN}[SONIC]${NC} System diagnostic complete: $version"
    else
        echo -e "${DALEK_RED}[DALEK]${NC} Sonic Screwdriver malfunction!"
        return 1
    fi
    
    # Test system checks
    if sonic system check &> /dev/null; then
        echo -e "${TARDIS_BLUE}[TARDIS]${NC} All systems nominal!"
    else
        echo -e "${CYBERMAN_YELLOW}[CYBERMAN]${NC} Minor system anomalies detected"
    fi
    
    return 0
}

# ============================================
# TARDIS Main Storyline
# ============================================

# TARDIS Landing Sequence
tardis_landing() {
    echo -e "${TARDIS_BLUE}
  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
  ████████████████████████████████████████
  ██${NC} TARDIS LANDING SEQUENCE INITIATED ${TARDIS_BLUE}██
  ██${NC} Destination: Ubuntu 22.04 LTS ${TARDIS_BLUE}██
  ████████████████████████████████████████
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀${NC}"
    
    # Landing animation
    for i in {1..3}; do
        echo -ne "${TIME_VORTEX}🌀 Materializing${NC}"
        for j in {1..5}; do echo -ne "."; sleep 0.3; done
        echo -ne "\r"
    done
    echo -ne "\033[K"
    
    echo -e "${SONIC_GREEN}[SONIC]${NC} *zzzzzap* Landing complete!"
    echo ""
}

# TARDIS Departure Sequence
tardis_departure() {
    echo -e "${TARDIS_BLUE}
  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
  ████████████████████████████████████████
  ██${NC} MISSION ACCOMPLISHED! ${TARDIS_BLUE}██
  ██${NC} Sonic-Screwdriver installed ${TARDIS_BLUE}██
  ████████████████████████████████████████
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀${NC}"
    
    # Departure animation
    for i in {1..3}; do
        echo -ne "${TIME_VORTEX}🌀 Dematerializing${NC}"
        for j in {1..5}; do echo -ne "."; sleep 0.3; done
        echo -ne "\r"
    done
    echo -ne "\033[K"
    
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Allons-y! The adventure continues!"
    echo ""
}

# Main TARDIS Storyline
main_storyline() {
    # Prologue
    init_tardis
    tardis_landing
    
    doctor_says "Right then! Time to install the Sonic Screwdriver on this Ubuntu system!"
    companion_says "What's the plan, Doctor?"
    doctor_says "First, we need to check the temporal compatibility!"
    
    # Act 1: System Checks
    check_not_root
    check_ubuntu_version
    check_architecture
    
    doctor_says "Excellent! This system is compatible with the TARDIS!"
    
    # Act 2: Dependency Installation
    if ! check_dependencies; then
        doctor_says "We need to refuel the TARDIS with some essential timey-wimey components!"
        install_dependencies
    fi
    
    # Act 3: TARDIS Materialization
    doctor_says "Now, let's materialize the TARDIS in this spacetime!"
    clone_repository
    
    # Act 4: Sonic Activation
    doctor_says "Time to activate the Sonic Screwdriver!"
    if ! build_sonic; then
        echo -e "${DALEK_RED}[DALEK]${NC} Dalek interference detected!"
        doctor_says "Oh no, Daleks! Let me try regenerating..."
        tardis_self_heal 1 3 "build_sonic" "Sonic activation"
    fi
    
    # Act 5: System Integration
    doctor_says "Now to integrate the Sonic with the system matrix!"
    install_sonic
    setup_environment
    
    # Act 6: Final Diagnostic
    doctor_says "Let's run a final diagnostic to make sure everything's working!"
    if verify_installation; then
        companion_says "Brilliant! It works!"
        doctor_says "Of course it works! I'm the Doctor!"
    else
        doctor_says "Hmm, something's not quite right. Let me check..."
        return 1
    fi
    
    # Epilogue
    tardis_departure
    
    echo -e "${GALLIFREYAN_GOLD}
========================================
Sonic-Screwdriver Installation Complete!
========================================
${NC}"
    echo "The Doctor has successfully installed the Sonic-Screwdriver on your system!"
    echo ""
    echo "Try these commands:"
    echo "  ${SONIC_GREEN}sonic --help${NC}              Show help"
    echo "  ${SONIC_GREEN}sonic --version${NC}           Show version"
    echo "  ${SONIC_GREEN}sonic system check${NC}        Check system compatibility"
    echo "  ${SONIC_GREEN}sonic tui${NC}                 Launch interactive interface"
    echo ""
    echo "Documentation:"
    echo "  ${TARDIS_BLUE}$SONIC_DIR/README_LINUX.md${NC}"
    echo "  ${TARDIS_BLUE}$SONIC_DIR/QUICKSTART_LINUX.md${NC}"
    echo ""
    echo "Log file: $LOG_FILE"
    
    return 0
}

# ============================================
# TARDIS Console Interface
# ============================================

# Show TARDIS help
tardis_help() {
    echo -e "${TARDIS_BLUE}
  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
  ████████████████████████████████████████
  ██${NC} TARDIS CONSOLE INTERFACE ${TARDIS_BLUE}██
  ██${NC} Doctor Who Themed Installer ${TARDIS_BLUE}██
  ████████████████████████████████████████
  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀${NC}"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  ${GALLIFREYAN_GOLD}install${NC}       - Full installation with Doctor Who storyline (default)"
    echo "  ${GALLIFREYAN_GOLD}quick${NC}         - Quick installation (no storyline)"
    echo "  ${GALLIFREYAN_GOLD}check${NC}         - Check system requirements"
    echo "  ${GALLIFREYAN_GOLD}repair${NC}        - Repair installation"
    echo "  ${GALLIFREYAN_GOLD}help${NC}          - Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 install           # Full Doctor Who experience"
    echo "  $0 quick             # Quick installation"
    echo "  $0 repair            # Repair installation"
}

# Quick installation (without storyline)
quick_install() {
    echo -e "${TARDIS_BLUE}[TARDIS]${NC} Initiating quick installation sequence..."
    
    check_not_root
    check_ubuntu_version
    check_architecture
    
    if ! check_dependencies; then
        install_dependencies
    fi
    
    clone_repository
    build_sonic
    install_sonic
    setup_environment
    
    if verify_installation; then
        echo -e "${SONIC_GREEN}[SONIC]${NC} Installation complete!"
    else
        echo -e "${DALEK_RED}[DALEK]${NC} Installation failed!"
        return 1
    fi
}

# ============================================
# Main Execution
# ============================================

COMMAND="${1:-install}"

case "$COMMAND" in
    install|full)
        main_storyline
        ;;
    quick)
        quick_install
        ;;
    check)
        check_not_root
        check_ubuntu_version
        check_architecture
        check_dependencies
        ;;
    repair)
        doctor_says "Repairing TARDIS systems..."
        verify_installation || true
        # Add repair logic here
        ;;
    help|--help|-h)
        tardis_help
        ;;
    *)
        echo "Unknown command: $COMMAND"
        tardis_help
        exit 1
        ;;
esac

# Create restart script for manual recovery
create_restart_script() {
    local restart_script="/tmp/sonic-tardis-restart.sh"
    
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Creating restart script at $restart_script"
    
    cat > "$restart_script" << 'INNER_EOF'
#!/bin/bash
# TARDIS Restart Script
# Manual recovery for Sonic-Screwdriver installation

echo "Restarting TARDIS installation..."

# Mark as restart
RESTART_FILE="/tmp/sonic-tardis-restart.marker"
touch "$RESTART_FILE"

# Find the installer
INSTALLER="$(find /tmp -name "install-ubuntu-tardis.sh" 2>/dev/null | head -1)"

if [ -z "$INSTALLER" ]; then
    echo "Error: Could not find TARDIS installer"
    exit 1
fi

# Restart with quick mode
bash "$INSTALLER" quick

# Clean up
rm -f "$RESTART_FILE"
rm -f "$0"
INNER_EOF
    
    chmod +x "$restart_script"
    
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Restart script created: $restart_script"
    echo -e "${GALLIFREYAN_GOLD}[DOCTOR]${NC} Run: bash $restart_script"
    echo "[RECOVERY] Restart script created at $restart_script" >> "$LOG_FILE"
}

# Add error logging to existing error detection
