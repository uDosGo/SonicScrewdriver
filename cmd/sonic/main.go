package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/OkAgentDigital/universal/classicmodern"
	"github.com/OkAgentDigital/universal/container"
	"github.com/OkAgentDigital/universal/disk"
	"github.com/OkAgentDigital/universal/iso"
	"github.com/OkAgentDigital/universal/library"
	"github.com/OkAgentDigital/universal/remote"
	"github.com/OkAgentDigital/universal/secrets"
	"github.com/OkAgentDigital/universal/state"
	"github.com/OkAgentDigital/universal/usb"
	"github.com/sonic-family/sonic-screwdriver/internal/homeassistant"
	"github.com/sonic-family/sonic-screwdriver/modules/ventoy"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ============================================
// DEV MODE - DEVELOPMENT ONLY
// These variables and functions are for development only
// and can be removed for production builds.
// Tag: DEV-ONLY
// ============================================

// DevModeFlag indicates if Sonic is running in development mode
// Tag: DEV-ONLY - Remove for production
var DevModeFlag bool

// DevModeSettings holds development mode configuration
// Tag: DEV-ONLY - Remove for production
var DevModeSettings = map[string]interface{}{
	"verbose_logging":    true,
	"debug_output":       true,
	"keep_temp_files":    true,
	"archive_builds":     true,
	"auto_archive_logs":  true,
	"show_mem_stats":     false,
	"show_gc_stats":      false,
	"trace_queries":      true,
}

// Tag: DEV-ONLY - Remove for production

var (
	containerRuntime container.Runtime
	libManager        *library.Manager
	stateDB           *state.DB
	secretStore       *secrets.SecretStore
	nodeRegistry      *secrets.NodeRegistry
	proxyServer       *secrets.ProxyServer
	haIntegration     *homeassistant.HAIntegration
)

func main() {
	// Initialize DEV MODE first
	initializeDevMode()
	
	if DevModeFlag {
		log.Printf("🔧 DEV MODE ENABLED - Development settings active")
		if DevModeSettings["verbose_logging"].(bool) {
			log.Printf("   Verbose logging: ON")
		}
		if DevModeSettings["archive_builds"].(bool) {
			log.Printf("   Build archiving: ON")
		}
	}

	// Initialize library manager
	indexPath := library.GetDefaultIndexPath()
	libManager = library.NewManager(indexPath)
	if err := libManager.Load(); err != nil {
		log.Printf("Warning: Could not load library index: %v", err)
		// Continue without library for now
	}

	// Initialize Docker runtime
	dockerRuntime, err := container.NewDockerRuntime()
	if err != nil {
		log.Printf("Warning: Could not initialize Docker runtime: %v", err)
		log.Printf("Falling back to mock runtime")
		containerRuntime = &container.DockerRuntime{} // Fallback to mock
	} else {
		containerRuntime = dockerRuntime
		defer dockerRuntime.Close()
	}

	// Initialize state database
	dbPath := state.GetDefaultDBPath()
	stateDB, err = state.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open state database: %v", err)
	}
	defer stateDB.Close()

	// Initialize secret store
	masterKeyPath := getMasterKeyPath()
	secretsPath := getSecretsPath()
	
	// Generate master key if it doesn't exist
	var masterKey []byte
	if _, err := os.Stat(masterKeyPath); os.IsNotExist(err) {
		masterKey, err = secrets.GenerateMasterKey()
		if err != nil {
			log.Fatalf("Failed to generate master key: %v", err)
		}
		if err := secrets.SaveMasterKey(masterKey, masterKeyPath); err != nil {
			log.Fatalf("Failed to save master key: %v", err)
		}
		log.Printf("Generated new master key at %s", masterKeyPath)
	} else {
		masterKey, err = secrets.LoadMasterKey(masterKeyPath)
		if err != nil {
			log.Fatalf("Failed to load master key: %v", err)
		}
	}

	secretStore, err = secrets.NewSecretStore(masterKey, secretsPath)
	if err != nil {
		log.Fatalf("Failed to initialize secret store: %v", err)
	}

	// Initialize node registry
	nodeRegistryPath := getNodeRegistryPath()
	nodeRegistry, err = secrets.NewNodeRegistry(nodeRegistryPath)
	if err != nil {
		log.Fatalf("Failed to initialize node registry: %v", err)
	}

	// Initialize proxy server
	proxyServer = secrets.NewProxyServer(secretStore)

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Error: game name required")
			os.Exit(1)
		}
		gameName := os.Args[2]
		fmt.Printf("Installing %s...\n", gameName)

		// Validate game exists in library
		if libManager != nil {
			game, err := libManager.GetGame(gameName)
			if err != nil {
				log.Fatalf("Game not found in library: %v", err)
			}

			// Validate manifest
			if err := libManager.ValidateManifest(game.Path); err != nil {
				log.Fatalf("Invalid manifest: %v", err)
			}
			fmt.Printf("Validated manifest for %s\n", gameName)
		}

		// Update state - mark as installed
		if err := stateDB.SetInstalled(gameName, "1.0.0"); err != nil {
			log.Fatalf("Failed to update state: %v", err)
		}

		fmt.Printf("%s installed successfully\n", gameName)
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Error: game name required")
			os.Exit(1)
		}
		gameName := os.Args[2]

		// Update state - mark as running
		if err := stateDB.SetRunning(gameName, true); err != nil {
			log.Fatalf("Failed to update state: %v", err)
		}

		err := containerRuntime.Start(gameName)
		if err != nil {
			// Rollback state on failure
			stateDB.SetRunning(gameName, false)
			log.Fatalf("Failed to start: %v", err)
		}
		fmt.Printf("Started %s\n", gameName)
	case "stop":
		if len(os.Args) < 3 {
			fmt.Println("Error: game name required")
			os.Exit(1)
		}
		gameName := os.Args[2]

		err := containerRuntime.Stop(gameName)
		if err != nil {
			log.Fatalf("Failed to stop: %v", err)
		}

		// Update state - mark as stopped
		if err := stateDB.SetRunning(gameName, false); err != nil {
			log.Fatalf("Failed to update state: %v", err)
		}

		fmt.Printf("Stopped %s\n", gameName)
	case "list":
		installations, err := stateDB.ListInstallations()
		if err != nil {
			log.Fatalf("Failed to list installations: %v", err)
		}

		if len(installations) == 0 {
			fmt.Println("No games installed")
			return
		}

		fmt.Println("Installed games:")
		for _, inst := range installations {
			status := "installed"
			if inst.Running {
				status = "running"
			}
			fmt.Printf("  - %s (v%s) - %s\n", inst.Name, inst.Version, status)
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Error: game name required")
			os.Exit(1)
		}
		gameName := os.Args[2]

		// Update state - mark as uninstalled
		if err := stateDB.Remove(gameName); err != nil {
			log.Fatalf("Failed to update state: %v", err)
		}

		err := containerRuntime.Remove(gameName)
		if err != nil {
			log.Fatalf("Failed to remove: %v", err)
		}

		fmt.Printf("Removed %s\n", gameName)
	case "--help", "-h", "help":
		printHelp()
	case "--version", "-v", "version":
		fmt.Println("vA2.0.0")
	case "tui", "menu":
		if secretStore == nil || nodeRegistry == nil || proxyServer == nil {
			fmt.Println("Error: TUI requires initialized components")
			os.Exit(1)
		}
		tui := secrets.NewTUI(secretStore, nodeRegistry, proxyServer)
		tui.Run()
	case "library":
		if len(os.Args) < 3 || os.Args[2] != "list" {
			fmt.Println("Usage: sonic library list")
			os.Exit(1)
		}
		if libManager == nil {
			fmt.Println("Error: Library not initialized")
			os.Exit(1)
		}
		games, err := libManager.List()
		if err != nil {
			log.Fatalf("Failed to list games: %v", err)
		}
		fmt.Println("Available games:")
		for _, game := range games {
			fmt.Printf("  - %s (%s)\n", game.Name, game.Status)
		}
	case "logs":
		if len(os.Args) < 3 {
			fmt.Println("Error: game name required")
			os.Exit(1)
		}
		fmt.Printf("Showing logs for %s...\n", os.Args[2])
		// TODO: Implement log viewing
		case "health":
			if len(os.Args) < 3 {
				fmt.Println("Usage: sonic health <game>|--all")
				os.Exit(1)
			}
			if os.Args[2] == "--all" {
				healthStatuses, err := containerRuntime.GetAllContainerHealth()
				if err != nil {
					log.Fatalf("Failed to get container health: %v", err)
				}
				if len(healthStatuses) == 0 {
					fmt.Println("No Sonic containers found")
					return
				}
				fmt.Println("Container Health Status:")
				for _, status := range healthStatuses {
					healthIcon := "✅"
					if !status.Healthy {
						healthIcon = "❌"
					}
					fmt.Printf("  %s %s: %s", healthIcon, status.Name, status.Status)
					if status.Error != "" {
						fmt.Printf(" (%s)", status.Error)
					}
					fmt.Println()
				}
			} else {
				gameName := os.Args[2]
				status, err := containerRuntime.CheckContainerHealth(gameName)
				if err != nil {
					log.Fatalf("Failed to check health: %v", err)
				}
				healthIcon := "✅"
				if !status.Healthy {
					healthIcon = "❌"
				}
				fmt.Printf("%s %s Health Status:\n", healthIcon, gameName)
				fmt.Printf("  Status: %s\n", status.Status)
				fmt.Printf("  Healthy: %t\n", status.Healthy)
				if status.Error != "" {
					fmt.Printf("  Error: %s\n", status.Error)
				}
				fmt.Printf("  Last Check: %s\n", status.Timestamp.Format("2006-01-02 15:04:05"))
			}
		case "repair":
			if len(os.Args) < 3 {
				fmt.Println("Usage: sonic repair <game>|--all")
				os.Exit(1)
			}
			if os.Args[2] == "--all" {
				healthStatuses, err := containerRuntime.GetAllContainerHealth()
				if err != nil {
					log.Fatalf("Failed to get container health: %v", err)
				}
				repairedCount := 0
				for _, status := range healthStatuses {
					if !status.Healthy {
						fmt.Printf("Attempting to repair %s...\n", status.Name)
						err := containerRuntime.RestartContainer(status.Name)
						if err != nil {
							log.Printf("Failed to repair %s: %v", status.Name, err)
						} else {
							fmt.Printf("✅ Repaired %s\n", status.Name)
							repairedCount++
						}
					}
				}
				if repairedCount == 0 {
					fmt.Println("All containers are healthy, no repairs needed")
				} else {
					fmt.Printf("Repaired %d container(s)\n", repairedCount)
				}
			} else {
				gameName := os.Args[2]
				fmt.Printf("Attempting to repair %s...\n", gameName)
				status, err := containerRuntime.CheckContainerHealth(gameName)
				if err != nil {
					log.Fatalf("Failed to check health: %v", err)
				}
				if status.Healthy {
					fmt.Printf("%s is already healthy\n", gameName)
					return
				}
				err = containerRuntime.RestartContainer(gameName)
				if err != nil {
					log.Fatalf("Failed to repair: %v", err)
				}
				fmt.Printf("✅ Repaired %s\n", gameName)
			}

	case "config":
		if len(os.Args) < 5 || os.Args[2] != "set" {
			fmt.Println("Usage: sonic config set <key> <value>")
			os.Exit(1)
		}
		fmt.Printf("Setting %s=%s\n", os.Args[3], os.Args[4])
		// TODO: Implement config management
	case "secret":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic secret <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  add <name> --value <value>     - Add a new secret")
			fmt.Println("  get <name> [--cached]          - Get a secret value (with optional cache)")
			fmt.Println("  rotate <name> --value <value>  - Rotate a secret")
			fmt.Println("  history <name>                - Show secret rotation history")
			fmt.Println("  list                           - List all secrets")
			fmt.Println("  grant <name> --node <node>     - Grant access to a node")
			fmt.Println("  revoke <name> --node <node>    - Revoke access from a node")
			fmt.Println("  policy <name>                  - Show secret policy")
			fmt.Println("  backup <file>                  - Create backup")
			fmt.Println("  restore <file>                 - Restore from backup")
			fmt.Println("  export <file>                  - Export encrypted backup")
			fmt.Println("  import <file>                  - Import encrypted backup")
			os.Exit(1)
		}
		handleSecretCommand(os.Args[2:])
	case "node":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic node <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  register --master <addr> --name <name>  - Register a new node")
			fmt.Println("  list                                    - List registered nodes")
			fmt.Println("  show <name>                            - Show node details")
			fmt.Println("  revoke <name>                          - Revoke node access")
			os.Exit(1)
		}
		handleNodeCommand(os.Args[2:])
	case "proxy":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic proxy <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  status                                 - Show proxy status")
			fmt.Println("  call <provider> --data <json>         - Call proxy endpoint")
			fmt.Println("  health                                 - Check proxy health")
			os.Exit(1)
		}
		handleProxyCommand(os.Args[2:])
	case "remote":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic remote <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  vnc setup [password] [geometry]       - Setup VNC server")
			fmt.Println("  vnc start                             - Start VNC server")
			fmt.Println("  vnc stop                              - Stop VNC server")
			fmt.Println("  ssh setup                             - Setup SSH access")
			fmt.Println("  samba setup <name> <path>            - Setup Samba file sharing")
			fmt.Println("  info                                  - Show remote access info")
			os.Exit(1)
		}
		handleRemoteCommand(os.Args[2:])
	case "mint":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic mint <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  check                                - Check Classic Modern readiness")
			fmt.Println("  install                              - Install Classic Modern theme")
			fmt.Println("  apply                                 - Apply Classic Modern theme")
			fmt.Println("  status                               - Show current theme status")
			fmt.Println("  info                                 - Show theme information")
			fmt.Println("  doctor                               - Run diagnostic checks")
			os.Exit(1)
		}
		handleMintCommand(os.Args[2:])
	case "ha", "homeassistant":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic ha <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  setup <url> <token>                  - Setup Home Assistant integration")
			fmt.Println("  configure                            - Configure HA integration")
			fmt.Println("  status                               - Show HA status")
			fmt.Println("  info                                 - Show HA instance info")
			fmt.Println("  embed <output.html>                  - Generate embed HTML file")
			fmt.Println("  kiosk enable|disable                 - Enable/disable kiosk mode")
			fmt.Println("  refresh <minutes>                    - Set refresh rate")
			fmt.Println("  version                              - Get HA version")
			fmt.Println("  check                                - Check HA connection")
			os.Exit(1)
		}
		handleHACommand(os.Args[2:])
	case "ventoy":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic ventoy <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  package <installer> <output>  - Create Ventoy bundle")
			fmt.Println("  validate <bundle.she>         - Validate bundle")
			fmt.Println("  info <bundle.she>              - Show bundle info")
			os.Exit(1)
		}
		handleVentoyCommand(os.Args[2:])
	case "system":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic system <command>")
			fmt.Println("Commands:")
			fmt.Println("  check       - Check system compatibility")
			fmt.Println("  info        - Show system information")
			fmt.Println("  resources   - Show system resources")
			fmt.Println("  devmode     - Show DEV MODE status")
			os.Exit(1)
		}
		handleSystemCommand(os.Args[2:])
	case "usb":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic usb <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  list                  - List available USB devices")
			fmt.Println("  prepare <device> <distro>  - Partition and format USB for distro")
			fmt.Println("  install <device> <distro>  - Full install: download ISO + write to USB")
			fmt.Println("  write <device> <iso>       - Write existing ISO to USB")
			fmt.Println("")
			fmt.Println("Distros: ubuntu, mint, classicmodern")
			os.Exit(1)
		}
		handleUSBCommand(os.Args[2:])
	case "iso":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic iso <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  list                  - List available distros")
			fmt.Println("  download <distro>     - Download ISO to cache")
			fmt.Println("  cache                 - Show cache status")
			fmt.Println("")
			fmt.Println("Distros: ubuntu, mint, classicmodern")
			os.Exit(1)
		}
		handleISOCommand(os.Args[2:])
	case "disk":
		if len(os.Args) < 3 {
			fmt.Println("Usage: sonic disk <command> [args]")
			fmt.Println("Commands:")
			fmt.Println("  list                  - List all block devices")
			fmt.Println("  info <device>         - Show device details")
			fmt.Println("  wipe <device>         - Wipe partition table (DESTRUCTIVE)")
			os.Exit(1)
		}
		handleDiskCommand(os.Args[2:])
	default:
		printHelp()
	}
}

func handleVentoyCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ventoy command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "package":
		if len(args) < 3 {
			fmt.Println("Usage: sonic ventoy package <installer-dir> <output.she>")
			os.Exit(1)
		}
		packageVentoyBundle(args[1], args[2])
	case "validate":
		if len(args) < 2 {
			fmt.Println("Usage: sonic ventoy validate <bundle.she>")
			os.Exit(1)
		}
		validateVentoyBundle(args[1])
	case "info":
		if len(args) < 2 {
			fmt.Println("Usage: sonic ventoy info <bundle.she>")
			os.Exit(1)
		}
		showBundleInfo(args[1])
	default:
		fmt.Printf("Error: unknown ventoy command: %s\n", command)
		os.Exit(1)
	}
}

func packageVentoyBundle(sourceDir, outputPath string) {
	fmt.Printf("Packaging Ventoy bundle from %s to %s\n", sourceDir, outputPath)
	
	// Create packager
	packager := ventoy.NewPackager(sourceDir, filepath.Dir(outputPath), libManager)
	
	// Create bundle
	bundlePath, err := packager.CreateBundle(filepath.Base(outputPath))
	if err != nil {
		log.Fatalf("Failed to create bundle: %v", err)
	}
	
	fmt.Printf("✓ Bundle created: %s\n", bundlePath)
	
	// Validate bundle
	err = packager.ValidateBundle(bundlePath)
	if err != nil {
		log.Fatalf("Bundle validation failed: %v", err)
	}
	
	fmt.Printf("✓ Bundle validated successfully\n")
	
	// Show bundle info
	info, err := ventoy.GetBundleInfo(bundlePath)
	if err != nil {
		log.Fatalf("Failed to get bundle info: %v", err)
	}
	
	fmt.Printf("Bundle Info:\n")
	for key, value := range info {
		fmt.Printf("  %s: %s\n", key, value)
	}
}

func validateVentoyBundle(bundlePath string) {
	fmt.Printf("Validating Ventoy bundle: %s\n", bundlePath)
	
	packager := ventoy.NewPackager("", "", libManager)
	err := packager.ValidateBundle(bundlePath)
	if err != nil {
		log.Fatalf("Bundle validation failed: %v", err)
	}
	
	fmt.Printf("✓ Bundle is valid\n")
}

// ============================================
// DEV MODE FUNCTIONS - DEVELOPMENT ONLY
// Tag: DEV-ONLY - Remove for production
// ============================================

// initializeDevMode checks environment and sets DevModeFlag
// DEV MODE is automatically enabled when:
// - SONIC_DEV_MODE=true environment variable is set
// - DEV_MODE=true environment variable is set
// - DEVSTUDIO_DEV_MODE=true environment variable is set
// - Running from VibeCli environment (VIBE_HOME or VIBE_SESSION_ID set)
// Tag: DEV-ONLY - Remove for production
func initializeDevMode() {
	// Check explicit environment variables
	if os.Getenv("SONIC_DEV_MODE") == "true" {
		DevModeFlag = true
		return
	}
	if os.Getenv("DEV_MODE") == "true" {
		DevModeFlag = true
		return
	}
	if os.Getenv("DEVSTUDIO_DEV_MODE") == "true" {
		DevModeFlag = true
		return
	}
	
	// Check for VibeCli environment (auto-enable DEV MODE)
	if os.Getenv("VIBE_HOME") != "" || os.Getenv("VIBE_SESSION_ID") != "" {
		DevModeFlag = true
		os.Setenv("SONIC_DEV_MODE", "true")
		os.Setenv("DEV_MODE", "true")
		os.Setenv("DEVSTUDIO_DEV_MODE", "true")
		log.Printf("🤖 VibeCli environment detected - DEV MODE auto-enabled")
		return
	}
	
	// Check if we're running in a development directory
	wd, _ := os.Getwd()
	if strings.Contains(wd, "Code/SonicScrewdriver") || 
	   strings.Contains(wd, "Code/DevStudio") {
		// Check if there's a .devmode file in DevStudio
		devModeFile := filepath.Join(os.Getenv("HOME"), "Code", "DevStudio", ".devmode")
		if data, err := os.ReadFile(devModeFile); err == nil {
			if strings.TrimSpace(string(data)) == "true" {
				DevModeFlag = true
				return
			}
		}
	}
	
	// Default: DEV MODE off in production
	DevModeFlag = false
}

// isDevMode returns true if running in development mode
// Tag: DEV-ONLY - Remove for production
func isDevMode() bool {
	return DevModeFlag
}

// getDevModeSetting returns a dev mode setting value
// Tag: DEV-ONLY - Remove for production
func getDevModeSetting(key string) bool {
	if !DevModeFlag {
		return false
	}
	if val, ok := DevModeSettings[key].(bool); ok {
		return val
	}
	return false
}

// Tag: DEV-ONLY - Remove for production
// ============================================

func getMasterKeyPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./master.key"
	}
	return filepath.Join(homeDir, ".sonic", "master.key")
}

func getSecretsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./secrets.enc"
	}
	return filepath.Join(homeDir, ".sonic", "secrets.enc")
}

func getNodeRegistryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./nodes.json"
	}
	return filepath.Join(homeDir, ".sonic", "nodes.json")
}

func handleRemoteCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: remote command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "vnc":
		if len(args) < 2 {
			fmt.Println("Usage: sonic remote vnc <subcommand> [args]")
			os.Exit(1)
		}
		handleVNCCommand(args[1:])
	case "ssh":
		if len(args) < 2 {
			fmt.Println("Usage: sonic remote ssh <subcommand>")
			os.Exit(1)
		}
		handleSSHCommand(args[1:])
	case "samba":
		if len(args) < 2 {
			fmt.Println("Usage: sonic remote samba <subcommand> [args]")
			os.Exit(1)
		}
		handleSambaCommand(args[1:])
	case "info":
		showRemoteAccessInfo()
	default:
		fmt.Printf("Error: unknown remote command: %s\n", command)
		os.Exit(1)
	}
}

func handleVNCCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: vnc subcommand required")
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "setup":
		password := "password"
		geometry := "1920x1080"
		if len(args) > 1 {
			password = args[1]
		}
		if len(args) > 2 {
			geometry = args[2]
		}
		setupVNC(password, geometry)
	case "start":
		startVNC()
	case "stop":
		stopVNC()
	default:
		fmt.Printf("Error: unknown vnc subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleSSHCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ssh subcommand required")
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "setup":
		setupSSH()
	default:
		fmt.Printf("Error: unknown ssh subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleSambaCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: samba subcommand required")
		os.Exit(1)
	}

	subcommand := args[0]
	switch subcommand {
	case "setup":
		if len(args) < 3 {
			fmt.Println("Usage: sonic remote samba setup <name> <path>")
			os.Exit(1)
		}
		setupSamba(args[1], args[2])
	default:
		fmt.Printf("Error: unknown samba subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func showBundleInfo(bundlePath string) {
	fmt.Printf("Showing bundle info: %s\n", bundlePath)
	
	info, err := ventoy.GetBundleInfo(bundlePath)
	if err != nil {
		log.Fatalf("Failed to get bundle info: %v", err)
	}
	
	fmt.Printf("Bundle Information:\n")
	for key, value := range info {
		fmt.Printf("  %s: %s\n", key, value)
	}
}

func handleSecretCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: secret command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "add":
		if len(args) < 4 || args[1] != "--value" {
			fmt.Println("Usage: sonic secret add <name> --value <value>")
			os.Exit(1)
		}
		addSecret(args[1], args[3])
	case "get":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret get <name>")
			os.Exit(1)
		}
		getSecretWithCache(args[1], false)
	case "list":
		listSecrets()
	case "grant":
		if len(args) < 4 || args[2] != "--node" {
			fmt.Println("Usage: sonic secret grant <name> --node <node>")
			os.Exit(1)
		}
		grantSecretAccess(args[1], args[3])
	case "revoke":
		if len(args) < 4 || args[2] != "--node" {
			fmt.Println("Usage: sonic secret revoke <name> --node <node>")
			os.Exit(1)
		}
		revokeSecretAccess(args[1], args[3])
	case "policy":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret policy <name>")
			os.Exit(1)
		}
		showSecretPolicy(args[1])
	case "backup":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret backup <file>")
			os.Exit(1)
		}
		backupSecrets(args[1])
	case "restore":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret restore <file>")
			os.Exit(1)
		}
		restoreSecrets(args[1])
	case "export":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret export <file>")
			os.Exit(1)
		}
		exportSecrets(args[1])
	case "import":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret import <file>")
			os.Exit(1)
		}
		importSecrets(args[1])
	case "rotate":
		if len(args) < 4 || args[2] != "--value" {
			fmt.Println("Usage: sonic secret rotate <name> --value <value>")
			os.Exit(1)
		}
		rotateSecret(args[1], args[3])
	case "history":
		if len(args) < 2 {
			fmt.Println("Usage: sonic secret history <name>")
			os.Exit(1)
		}
		showSecretHistory(args[1])
	default:
		fmt.Printf("Error: unknown secret command: %s\n", command)
		os.Exit(1)
	}
}

func handleNodeCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: node command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "register":
		if len(args) < 5 || args[1] != "--master" || args[3] != "--name" {
			fmt.Println("Usage: sonic node register --master <addr> --name <name>")
			os.Exit(1)
		}
		registerNode(args[2], args[4])
	case "list":
		listNodes()
	case "show":
		if len(args) < 2 {
			fmt.Println("Usage: sonic node show <name>")
			os.Exit(1)
		}
		showNode(args[1])
	case "revoke":
		if len(args) < 2 {
			fmt.Println("Usage: sonic node revoke <name>")
			os.Exit(1)
		}
		revokeNode(args[1])
	default:
		fmt.Printf("Error: unknown node command: %s\n", command)
		os.Exit(1)
	}
}

func handleProxyCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: proxy command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "status":
		proxyStatus()
	case "call":
		if len(args) < 4 || args[2] != "--data" {
			fmt.Println("Usage: sonic proxy call <provider> --data <json>")
			os.Exit(1)
		}
		proxyCall(args[1], args[3])
	case "health":
		proxyHealth()
	default:
		fmt.Printf("Error: unknown proxy command: %s\n", command)
		os.Exit(1)
	}
}

func addSecret(name, value string) {
	fmt.Printf("Adding secret %s...\n", name)
	
	// Add the secret
	if err := secretStore.AddSecret(name, value); err != nil {
		log.Fatalf("Failed to add secret: %v", err)
	}
	
	// Check if this is an API key and test it automatically
	if name == "openrouter_api_key" || name == "deepseek_api_key" || name == "gemini_api_key" || name == "github_token" {
		provider := ""
		if name == "openrouter_api_key" {
			provider = "openrouter"
		} else if name == "deepseek_api_key" {
			provider = "deepseek"
		} else if name == "gemini_api_key" {
			provider = "gemini"
		} else if name == "github_token" {
			provider = "github"
		}
		
		if provider != "" {
			fmt.Printf("Testing %s API key...\n", provider)
			err := proxyServer.TestAPIKey(provider, value)
			if err != nil {
				fmt.Printf("⚠️  API key test failed: %v\n", err)
				fmt.Printf("The key has been stored but may not be valid.\n")
			} else {
				fmt.Printf("✅ API key is valid\n")
			}
		}
	}
	
	fmt.Printf("✓ Secret %s added\n", name)
}

func getSecret(name string) {
	getSecretWithCache(name, false)
}

func getSecretWithCache(name string, allowCached bool) {
	fmt.Printf("Getting secret %s...\n", name)
	
	value, fromCache, err := secretStore.GetSecretWithCache(name, allowCached)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}
	
	if fromCache {
		fmt.Printf("⚠️  Using cached value (master may be unreachable)\n")
	}
	
	fmt.Printf("Secret value: %s\n", value)
}

func listSecrets() {
	fmt.Println("Listing secrets...")
	
	secrets, err := secretStore.ListSecrets()
	if err != nil {
		log.Fatalf("Failed to list secrets: %v", err)
	}
	
	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return
	}
	
	fmt.Println("Secrets:")
	for _, secret := range secrets {
		policy, _ := secretStore.GetPolicy(secret)
		allowedNodes := "all"
		if len(policy.AllowedNodes) > 0 {
			allowedNodes = "specific nodes"
		}
		fmt.Printf("  - %s (allowed: %s)\n", secret, allowedNodes)
	}
}

func grantSecretAccess(name, node string) {
	fmt.Printf("Granting access to secret %s for node %s...\n", name, node)
	
	// Get current policy
	policy, err := secretStore.GetPolicy(name)
	if err != nil {
		// If no policy exists, create a new one
		policy = secrets.SecretPolicy{
			AllowedNodes: []string{},
			AllowedRoles:  []string{},
			RateLimit:     "60/min",
		}
	}
	
	// Add node to allowed nodes
	for _, allowedNode := range policy.AllowedNodes {
		if allowedNode == node {
			fmt.Printf("✓ Access already granted\n")
			return
		}
	}
	
	policy.AllowedNodes = append(policy.AllowedNodes, node)
	
	if err := secretStore.SetPolicy(name, policy); err != nil {
		log.Fatalf("Failed to grant access: %v", err)
	}
	
	// Also update node registry
	if err := nodeRegistry.GrantSecretAccess(node, name); err != nil {
		log.Printf("Warning: Failed to update node registry: %v", err)
	}
	
	fmt.Printf("✓ Access granted\n")
}

func revokeSecretAccess(name, node string) {
	fmt.Printf("Revoking access to secret %s from node %s...\n", name, node)
	
	policy, err := secretStore.GetPolicy(name)
	if err != nil {
		log.Fatalf("Failed to get policy: %v", err)
	}
	
	// Remove node from allowed nodes
	newAllowedNodes := []string{}
	for _, allowedNode := range policy.AllowedNodes {
		if allowedNode != node {
			newAllowedNodes = append(newAllowedNodes, allowedNode)
		}
	}
	
	policy.AllowedNodes = newAllowedNodes
	
	if err := secretStore.SetPolicy(name, policy); err != nil {
		log.Fatalf("Failed to revoke access: %v", err)
	}
	
	// Also update node registry
	if err := nodeRegistry.RevokeSecretAccess(node, name); err != nil {
		log.Printf("Warning: Failed to update node registry: %v", err)
	}
	
	fmt.Printf("✓ Access revoked\n")
}

func showSecretPolicy(name string) {
	fmt.Printf("Showing policy for secret %s...\n", name)
	
	policy, err := secretStore.GetPolicy(name)
	if err != nil {
		fmt.Println("Policy: No policy set")
		return
	}
	
	fmt.Printf("Policy for %s:\n", name)
	fmt.Printf("  Allowed Nodes: %v\n", policy.AllowedNodes)
	fmt.Printf("  Allowed Roles: %v\n", policy.AllowedRoles)
	fmt.Printf("  Rate Limit: %s\n", policy.RateLimit)
}

func backupSecrets(filePath string) {
	fmt.Printf("Creating backup to %s...\n", filePath)
	
	if err := secretStore.Backup(filePath); err != nil {
		log.Fatalf("Failed to create backup: %v", err)
	}
	
	fmt.Printf("✓ Backup created successfully\n")
}

func restoreSecrets(filePath string) {
	fmt.Printf("Restoring from backup %s...\n", filePath)
	
	if err := secretStore.Restore(filePath); err != nil {
		log.Fatalf("Failed to restore backup: %v", err)
	}
	
	fmt.Printf("✓ Backup restored successfully\n")
}

func exportSecrets(filePath string) {
	fmt.Printf("Exporting encrypted backup to %s...\n", filePath)
	
	if err := secretStore.ExportBackup(filePath); err != nil {
		log.Fatalf("Failed to export backup: %v", err)
	}
	
	fmt.Printf("✓ Encrypted backup exported successfully\n")
}

func importSecrets(filePath string) {
	fmt.Printf("Importing encrypted backup from %s...\n", filePath)
	
	if err := secretStore.ImportBackup(filePath); err != nil {
		log.Fatalf("Failed to import backup: %v", err)
	}
	
	fmt.Printf("✓ Encrypted backup imported successfully\n")
}

// VNC server instance
var vncServer *remote.VNCServer

func setupVNC(password, geometry string) {
	fmt.Printf("Setting up VNC server with password and geometry %s...\n", geometry)
	
	vncServer = remote.NewVNCServer(5901, password, geometry)
	
	if err := vncServer.SetupVNC(); err != nil {
		log.Fatalf("Failed to setup VNC: %v", err)
	}
	
	fmt.Println("✓ VNC server setup completed")
	fmt.Println("Run 'sonic remote vnc start' to start the server")
}

func startVNC() {
	if vncServer == nil {
		fmt.Println("Error: VNC server not set up. Run 'sonic remote vnc setup' first.")
		os.Exit(1)
	}
	
	fmt.Println("Starting VNC server...")
	
	if err := vncServer.StartVNC(); err != nil {
		log.Fatalf("Failed to start VNC: %v", err)
	}
	
	fmt.Println("✓ VNC server started")
	fmt.Printf("Connect using: vncviewer %s:1\n", remote.GetLocalIP())
}

func stopVNC() {
	if vncServer == nil {
		fmt.Println("Error: VNC server not set up.")
		os.Exit(1)
	}
	
	fmt.Println("Stopping VNC server...")
	
	if err := vncServer.StopVNC(); err != nil {
		log.Fatalf("Failed to stop VNC: %v", err)
	}
	
	fmt.Println("✓ VNC server stopped")
}

func setupSSH() {
	fmt.Println("Setting up SSH for remote access...")
	
	if err := remote.SetupSSH(); err != nil {
		log.Fatalf("Failed to setup SSH: %v", err)
	}
	
	fmt.Println("✓ SSH setup completed")
	fmt.Printf("Connect using: ssh %s@%s\n", remote.GetCurrentUser(), remote.GetLocalIP())
}

func setupSamba(shareName, sharePath string) {
	fmt.Printf("Setting up Samba file sharing for %s at %s...\n", shareName, sharePath)
	
	if err := remote.SetupSamba(shareName, sharePath); err != nil {
		log.Fatalf("Failed to setup Samba: %v", err)
	}
	
	fmt.Println("✓ Samba setup completed")
	fmt.Printf("Access shared files at: smb://%s/%s\n", remote.GetLocalIP(), shareName)
}

func showRemoteAccessInfo() {
	fmt.Println("🌐 Remote Access Information")
	fmt.Println("============================")
	fmt.Println(remote.GetRemoteAccessInfo())
}

func rotateSecret(name, newValue string) {
	fmt.Printf("Rotating secret %s...\n", name)
	
	if err := secretStore.RotateSecret(name, newValue); err != nil {
		log.Fatalf("Failed to rotate secret: %v", err)
	}
	
	fmt.Printf("✓ Secret %s rotated successfully\n", name)
	
	// Show history
	history, err := secretStore.GetSecretHistory(name)
	if err == nil && len(history) > 0 {
		fmt.Println("\nRotation History:")
		for i, entry := range history {
			fmt.Printf("  %d. %s (on %s)\n", i+1, entry["action"], entry["date"])
		}
	}
}

func showSecretHistory(name string) {
	fmt.Printf("Showing rotation history for secret %s...\n", name)
	
	history, err := secretStore.GetSecretHistory(name)
	if err != nil {
		if err.Error() == "no history available" {
			fmt.Println("No rotation history available for this secret")
			return
		}
		log.Fatalf("Failed to get secret history: %v", err)
	}
	
	if len(history) == 0 {
		fmt.Println("No rotation history available")
		return
	}
	
	fmt.Printf("\nRotation History for %s:\n", name)
	fmt.Println("┌─────────┬──────────────────────────────────────────────┐")
	fmt.Println("│ #      │ Date                │ Action          │ Value (partial) │")
	fmt.Println("├─────────┼──────────────────────────────────────────────┤")
	
	for i, entry := range history {
		valuePreview := entry["value"]
		if len(valuePreview) > 10 {
			valuePreview = valuePreview[:10] + "..."
		}
		fmt.Printf("│ %-7d │ %-20s │ %-14s │ %-15s │\n", 
			i+1, entry["date"], entry["action"], valuePreview)
	}
	
	fmt.Println("└─────────┴──────────────────────────────────────────────┘")
}

func handleMintCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: mint command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "check":
		checkClassicModernReadiness()
	case "install":
		installClassicModern()
	case "apply":
		applyClassicModern()
	case "status":
		checkClassicModernStatus()
	case "info":
		showClassicModernInfo()
	case "doctor":
		runClassicModernDoctor()
	default:
		fmt.Printf("Error: unknown mint command: %s\n", command)
		os.Exit(1)
	}
}

func checkClassicModernReadiness() {
	fmt.Println("Checking Classic Modern Mint installation readiness...")
	
	// Initialize database connection
	dbPath := state.GetDefaultDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	checker := classicmodern.NewReadinessChecker(db)
	report, err := checker.GenerateInstallationReport()
	if err != nil {
		log.Fatalf("Failed to generate readiness report: %v", err)
	}
	
	fmt.Println(report)
}

func installClassicModern() {
	fmt.Println("Installing Classic Modern theme...")
	fmt.Println("✓ Theme installation would be handled here")
	fmt.Println("Note: Actual installation requires theme files to be present")
}

func applyClassicModern() {
	fmt.Println("Applying Classic Modern theme...")
	
	// Initialize database connection
	dbPath := state.GetDefaultDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	checker := classicmodern.NewReadinessChecker(db)
	
	// Check if theme is ready
	check, err := checker.PerformCheck()
	if err != nil {
		log.Fatalf("Failed to check readiness: %v", err)
	}
	
	if !check.ThemeReady {
		log.Fatal("Theme files are not ready. Run 'sonic mint check' first.")
	}
	
	// Apply theme
	if err := checker.ApplyTheme(); err != nil {
		log.Fatalf("Failed to apply theme: %v", err)
	}
	
	fmt.Println("✓ Classic Modern theme applied successfully")
	fmt.Println("Note: You may need to log out and back in for full effect")
}

func checkClassicModernStatus() {
	fmt.Println("Checking Classic Modern theme status...")
	
	// Initialize database connection
	dbPath := state.GetDefaultDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	checker := classicmodern.NewReadinessChecker(db)
	
	// Check if theme is installed
	installed, err := checker.CheckThemeInstallation()
	if err != nil {
		log.Fatalf("Failed to check theme installation: %v", err)
	}
	
	if installed {
		fmt.Println("✅ Classic Modern theme is currently active")
	} else {
		fmt.Println("❌ Classic Modern theme is not active")
		fmt.Println("Run 'sonic mint apply' to activate it")
	}
}

func showClassicModernInfo() {
	fmt.Println("Showing Classic Modern theme information...")
	
	// Initialize database connection
	dbPath := state.GetDefaultDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	checker := classicmodern.NewReadinessChecker(db)
	
	info, err := checker.GetThemeInfo()
	if err != nil {
		log.Fatalf("Failed to get theme info: %v", err)
	}
	
	fmt.Printf("Classic Modern Theme Information:\n")
	fmt.Printf("  Installed: %v\n", info["installed"])
	if info["installed"] == true {
		fmt.Printf("  Location: %s\n", info["location"])
		fmt.Printf("  Name: %s\n", info["name"])
		fmt.Printf("  Version: %s\n", info["version"])
		fmt.Printf("  Description: %s\n", info["description"])
		fmt.Printf("  OBF Present: %v\n", info["obf_present"])
	}
}

func runClassicModernDoctor() {
	fmt.Println("Running Classic Modern diagnostic checks...")
	
	// Initialize database connection
	dbPath := state.GetDefaultDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	checker := classicmodern.NewReadinessChecker(db)
	
	// Perform comprehensive check
	check, err := checker.PerformCheck()
	if err != nil {
		log.Fatalf("Failed to perform diagnostic check: %v", err)
	}
	
	fmt.Println("🏥 Classic Modern Diagnostic Report")
	fmt.Println("===================================")
	
	// Export to JSON for detailed analysis
	exportPath := "~/classic-modern-diagnostic.json"
	expandedPath := expandPath(exportPath)
	if err := checker.ExportCheckToJSON(expandedPath); err != nil {
		log.Printf("Warning: Failed to export diagnostic report: %v", err)
	} else {
		fmt.Printf("✓ Detailed report exported to: %s\n", expandedPath)
	}
	
	// Show summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Theme:        %s\n", statusEmoji(check.ThemeReady))
	fmt.Printf("  Icons:        %s\n", statusEmoji(check.IconsReady))
	fmt.Printf("  Fonts:        %s\n", statusEmoji(check.FontsReady))
	fmt.Printf("  OBF:          %s\n", statusEmoji(check.OBFValid))
	fmt.Printf("  Dependencies: %s\n", statusEmoji(check.DependenciesReady))
	
	if len(check.Issues) > 0 {
		fmt.Printf("\n❌ Issues found: %d\n", len(check.Issues))
		for _, issue := range check.Issues {
			fmt.Printf("  - %s\n", issue)
		}
	} else {
		fmt.Printf("\n✅ All systems healthy!\n")
	}
}

func statusEmoji(ready bool) string {
	if ready {
		return "✅"
	}
	return "❌"
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func registerNode(masterAddr, nodeName string) {
	fmt.Printf("Registering node %s with master %s...\n", nodeName, masterAddr)
	
	node, err := nodeRegistry.RegisterNode(nodeName, masterAddr)
	if err != nil {
		log.Fatalf("Failed to register node: %v", err)
	}
	
	fmt.Printf("✓ Node %s registered (ID: %s)\n", node.Name, node.ID)
}

func listNodes() {
	fmt.Println("Listing registered nodes...")
	
	nodes, err := nodeRegistry.ListNodes()
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}
	
	if len(nodes) == 0 {
		fmt.Println("No nodes registered")
		return
	}
	
	fmt.Printf("%-15s %-10s %-20s %-10s\n", "ID", "NAME", "STATUS", "LAST SEEN")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")
	
	for _, node := range nodes {
		fmt.Printf("%-15s %-10s %-20s %-10s\n", 
			node.ID, node.Name, node.Status, node.LastSeen.Format("2006-01-02 15:04:05"))
	}
}

func showNode(name string) {
	fmt.Printf("Showing details for node %s...\n", name)
	
	node, err := nodeRegistry.GetNode(name)
	if err != nil {
		log.Fatalf("Failed to get node: %v", err)
	}
	
	fmt.Printf("Node: %s\n", node.Name)
	fmt.Printf("  ID: %s\n", node.ID)
	fmt.Printf("  Status: %s\n", node.Status)
	fmt.Printf("  Last seen: %s\n", node.LastSeen.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Allowed secrets: %v\n", node.AllowedSecrets)
}

func revokeNode(name string) {
	fmt.Printf("Revoking node %s...\n", name)
	
	if err := nodeRegistry.RevokeNode(name); err != nil {
		log.Fatalf("Failed to revoke node: %v", err)
	}
	
	fmt.Printf("✓ Node %s revoked\n", name)
}

func proxyStatus() {
	fmt.Println("Proxy Status:")
	
	status := proxyServer.GetStatus()
	
	fmt.Println("┌─────────────┬──────────┬────────────┬─────────────┬──────────────┐")
	fmt.Println("│ Provider    │ Calls    │ Errors     │ Rate Limit  │ Status       │")
	fmt.Println("├─────────────┼──────────┼────────────┼─────────────┼──────────────┤")
	
	providers := []string{"openrouter", "deepseek", "gemini", "github"}
	for _, provider := range providers {
		if providerStatus, ok := status[provider]; ok {
			calls := providerStatus["calls"].(int)
			rateLimit := providerStatus["rate_limit"].(string)
			statusText := providerStatus["status"].(string)
			
			// Simple status emoji
			emoji := "🟢"
			if statusText != "healthy" {
				emoji = "🔴"
			}
			
			fmt.Printf("│ %-10s │ %-7d │ 0 (0%%)     │ %-10s │ %s %-11s │\n", 
				provider, calls, rateLimit, emoji, statusText)
		}
	}
	
	fmt.Println("└─────────────┴──────────┴────────────┴─────────────┴──────────────┘")
}

func proxyCall(provider, data string) {
	fmt.Printf("Calling proxy %s with data: %s\n", provider, data)
	
	// Parse the data as JSON
	var requestData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &requestData); err != nil {
		log.Fatalf("Failed to parse request data: %v", err)
	}
	
	request := secrets.ProxyRequest{
		Provider: provider,
		Method:   "POST",
		Path:     "/chat/completions",
		Headers:  map[string]string{},
		Body:     requestData,
	}
	
	response, err := proxyServer.HandleRequest(request)
	if err != nil {
		log.Fatalf("Failed to handle proxy request: %v", err)
	}
	
	fmt.Printf("✓ Proxy call completed\n")
	fmt.Printf("Status: %d\n", response.Status)
	fmt.Printf("Response: %+v\n", response.Body)
}

func proxyHealth() {
	fmt.Println("Checking proxy health...")
	
	health := proxyServer.GetHealth()
	
	for provider, status := range health {
		emoji := "✅"
		if status != "healthy" {
			emoji = "❌"
		}
		fmt.Printf("%s %s: %s\n", emoji, provider, status)
	}
	
	fmt.Println("\nProxy health check completed!")
}

func handleHACommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: ha command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "setup":
		if len(args) < 3 {
			fmt.Println("Usage: sonic ha setup <url> <token>")
			os.Exit(1)
		}
		setupHA(args[1], args[2])
	case "configure":
		configureHA()
	case "status":
		showHAStatus()
	case "info":
		showHAInfo()
	case "embed":
		if len(args) < 2 {
			fmt.Println("Usage: sonic ha embed <output.html>")
			os.Exit(1)
		}
		generateEmbedFile(args[1])
	case "kiosk":
		if len(args) < 2 {
			fmt.Println("Usage: sonic ha kiosk enable|disable")
			os.Exit(1)
		}
		configureKioskMode(args[1])
	case "refresh":
		if len(args) < 2 {
			fmt.Println("Usage: sonic ha refresh <minutes>")
			os.Exit(1)
		}
		configureRefreshRate(args[1])
	case "version":
		showHAVersion()
	case "check":
		checkHAConnection()
	default:
		fmt.Printf("Error: unknown ha command: %s\n", command)
		os.Exit(1)
	}
}

func setupHA(baseURL, apiToken string) {
	fmt.Printf("Setting up Home Assistant integration with %s...\n", baseURL)
	
	haIntegration = homeassistant.NewHAIntegration(baseURL, apiToken)
	
	if err := haIntegration.Configure(); err != nil {
		log.Fatalf("Failed to configure HA integration: %v", err)
	}
	
	fmt.Printf("✓ Home Assistant integration configured\n")
	
	// Save configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not get home directory: %v", err)
		return
	}
	
	configPath := filepath.Join(homeDir, ".sonic", "ha-config.json")
	if err := haIntegration.SaveConfig(configPath); err != nil {
		log.Printf("Warning: Failed to save HA config: %v", err)
	}
}

func configureHA() {
	fmt.Println("Configuring Home Assistant integration...")
	
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	if err := haIntegration.Configure(); err != nil {
		log.Fatalf("Failed to configure HA integration: %v", err)
	}
	
	fmt.Printf("✓ Home Assistant integration configured\n")
}

func showHAStatus() {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	config := haIntegration.GetStatus()
	fmt.Printf("Home Assistant Status:\n")
	fmt.Printf("  Base URL: %s\n", config.BaseURL)
	fmt.Printf("  Status: %s\n", config.Status)
	fmt.Printf("  Last Checked: %s\n", config.LastChecked)
	fmt.Printf("  Kiosk Mode: %t\n", config.KioskMode)
	fmt.Printf("  Refresh Rate: %d minutes\n", config.RefreshRate)
}

func showHAInfo() {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	info, err := haIntegration.GetHAInfo()
	if err != nil {
		log.Fatalf("Failed to get HA info: %v", err)
	}
	
	fmt.Printf("Home Assistant Information:\n")
	for key, value := range info {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

func generateEmbedFile(outputPath string) {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	fmt.Printf("Generating embed HTML file at %s...\n", outputPath)
	
	if err := haIntegration.GenerateEmbedFile(outputPath); err != nil {
		log.Fatalf("Failed to generate embed file: %v", err)
	}
	
	fmt.Printf("✓ Embed HTML file generated\n")
}

func configureKioskMode(mode string) {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	if mode != "enable" && mode != "disable" {
		fmt.Println("Error: mode must be 'enable' or 'disable'")
		os.Exit(1)
	}
	
	enabled := mode == "enable"
	haIntegration.SetKioskMode(enabled)
	
	fmt.Printf("✓ Kiosk mode %s\n", mode + "d")
}

func configureRefreshRate(rateStr string) {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	// Simple parsing - in production you'd want better error handling
	var rate int
	fmt.Sscanf(rateStr, "%d", &rate)
	
	if rate <= 0 {
		fmt.Println("Error: refresh rate must be positive")
		os.Exit(1)
	}
	
	haIntegration.SetRefreshRate(rate)
	
	fmt.Printf("✓ Refresh rate set to %d minutes\n", rate)
}

func showHAVersion() {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	version, err := haIntegration.GetHAVersion()
	if err != nil {
		log.Fatalf("Failed to get HA version: %v", err)
	}
	
	fmt.Printf("Home Assistant Version: %s\n", version)
}

func checkHAConnection() {
	if haIntegration == nil {
		fmt.Println("Error: HA integration not set up. Run 'sonic ha setup' first.")
		os.Exit(1)
	}
	
	fmt.Println("Checking Home Assistant connection...")
	
	status, err := haIntegration.CheckHAStatus()
	if err != nil {
		log.Fatalf("Failed to check HA status: %v", err)
	}
	
	if status {
		fmt.Println("✅ Home Assistant is accessible")
	} else {
		fmt.Println("❌ Home Assistant is not accessible")
	}
}

func printHelp() {
	fmt.Println(`SonicScrewdriver vA2.0.0 - API Central Hub

Usage:
  sonic install <game>     Install a game from curated library
  sonic start <game>       Start a game container
  sonic stop <game>        Stop a game container
  sonic list               List installed games
  sonic remove <game>      Remove a game

Commands:
  sonic tui                Launch interactive TUI menu
  sonic menu               Launch interactive menu
  sonic library list       Show available curated games
  sonic logs <game>        Show container logs
  sonic health <game>|--all  Check container health status
  sonic repair <game>|--all  Repair unhealthy containers
  sonic config set <key> <value>
  
  # Secret Management (API Central Hub)
  sonic secret add <name> --value <value>  Add a secret
  sonic secret get <name>                  Get a secret
  sonic secret list                        List secrets
  sonic secret grant <name> --node <node> Grant access
  sonic secret revoke <name> --node <node> Revoke access
  sonic secret policy <name>               Show policy
  
  # Node Management
  sonic node register --master <addr> --name <name>  Register node
  sonic node list                                    List nodes
  sonic node show <name>                            Show node
  sonic node revoke <name>                          Revoke node
  
  # API Proxy
  sonic proxy status                                Show proxy status
  sonic proxy call <provider> --data <json>         Call proxy
  sonic proxy health                                Check health
  
  # Remote Access
  sonic remote vnc setup [password] [geometry]     Setup VNC server
  sonic remote vnc start                           Start VNC server
  sonic remote vnc stop                            Stop VNC server
  sonic remote ssh setup                           Setup SSH access
  sonic remote samba setup <name> <path>           Setup Samba sharing
  sonic remote info                                Show remote access info
  
  # Home Assistant
  sonic ha setup <url> <token>                  Setup Home Assistant integration
  sonic ha configure                            Configure HA integration
  sonic ha status                               Show HA status
  sonic ha info                                 Show HA instance info
  sonic ha embed <output.html>                  Generate embed HTML file
  sonic ha kiosk enable|disable                 Enable/disable kiosk mode
  sonic ha refresh <minutes>                    Set refresh rate
  sonic ha version                              Get HA version
  sonic ha check                                Check HA connection
  
  # Classic Modern Mint
  sonic mint check                                Check installation readiness
  sonic mint install                              Install Classic Modern theme
  sonic mint apply                                 Apply Classic Modern theme
  sonic mint status                               Show theme status
  sonic mint info                                 Show theme information
  sonic mint doctor                               Run diagnostic checks
  
  # Ventoy
  sonic ventoy package    Create Ventoy installer bundle
  sonic ventoy validate    Validate Ventoy bundle
  sonic ventoy info        Show bundle information
  
  # System
  sonic system check       Check system compatibility
  sonic system info        Show system information
  sonic system resources   Show system resources
  sonic system devmode     Show DEV MODE status

Flags:
  --help, -h               Show this help
  --version, -v            Show version`)
}

func handleSystemCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: system command required")
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "check":
		performSystemChecks()
	case "info":
		showSystemInfo()
	case "resources":
		showSystemResources()
	case "devmode":
		showDevModeStatus()
	default:
		fmt.Printf("Error: unknown system command: %s\n", command)
		os.Exit(1)
	}
}

// showDevModeStatus displays the current DEV MODE status
// Tag: DEV-ONLY - Remove for production
func showDevModeStatus() {
	fmt.Println("🔧 DEV MODE Status")
	fmt.Println("===================")
	if DevModeFlag {
		fmt.Println("Status:  ✅ ENABLED")
		fmt.Println("Source:  Environment variable or VibeCli auto-detection")
		fmt.Println("")
		fmt.Println("Active Settings:")
		for key, val := range DevModeSettings {
			if b, ok := val.(bool); ok && b {
				fmt.Printf("  ✓ %s\n", key)
			}
		}
	} else {
		fmt.Println("Status:  ❌ DISABLED")
		fmt.Println("")
		fmt.Println("To enable DEV MODE:")
		fmt.Println("  - Set environment variable: export SONIC_DEV_MODE=true")
		fmt.Println("  - Or run from VibeCli (auto-enabled)")
		fmt.Println("  - Or create ~/.devmode file with 'true' in DevStudio")
	}
	fmt.Println("")
	
	// Show relevant environment variables
	fmt.Println("Environment:")
	if val := os.Getenv("SONIC_DEV_MODE"); val != "" {
		fmt.Printf("  SONIC_DEV_MODE=%s\n", val)
	}
	if val := os.Getenv("DEV_MODE"); val != "" {
		fmt.Printf("  DEV_MODE=%s\n", val)
	}
	if val := os.Getenv("DEVSTUDIO_DEV_MODE"); val != "" {
		fmt.Printf("  DEVSTUDIO_DEV_MODE=%s\n", val)
	}
	if val := os.Getenv("VIBE_HOME"); val != "" {
		fmt.Printf("  VIBE_HOME=%s\n", val)
	}
	if val := os.Getenv("VIBE_SESSION_ID"); val != "" {
		fmt.Printf("  VIBE_SESSION_ID=%s\n", val)
	}
}

func performSystemChecks() {
	fmt.Println("🔍 Performing system compatibility checks...")
	
	// Check OS
	osInfo, err := getOSInfo()
	if err != nil {
		log.Printf("Warning: Could not determine OS: %v", err)
		return
	}
	
	fmt.Printf("  OS: %s %s\n", osInfo.Name, osInfo.Version)
	
	// Check architecture
	arch := "unknown"
	if runtime.GOOS == "linux" {
		out, _ := exec.Command("uname", "-m").Output()
		arch = strings.TrimSpace(string(out))
	}
	fmt.Printf("  Architecture: %s\n", arch)
	
	// Check required dependencies
	checkDependencies()
	
	// Check Docker
	checkDocker()
	
	// Check Go version
	checkGoVersion()
	
	fmt.Println("✅ System checks complete")
}

func getOSInfo() (*SystemInfo, error) {
	// Read /etc/os-release for Linux
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("could not read OS info: %v", err)
	}
	
	info := &SystemInfo{}
	lines := strings.Split(string(data), "\n")
	
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			info.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			info.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}
	
	if info.Name == "" || info.Version == "" {
		return nil, fmt.Errorf("could not parse OS info")
	}
	
	return info, nil
}

func checkDependencies() {
	fmt.Println("  Checking dependencies...")
	
	dependencies := []string{"git", "make", "curl", "g++", "docker"}
	missing := []string{}
	
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			missing = append(missing, dep)
		}
	}
	
	if len(missing) > 0 {
		log.Printf("  ⚠️  Missing dependencies: %v", missing)
	} else {
		fmt.Println("  ✓ All dependencies found")
	}
}

func checkDocker() {
	fmt.Println("  Checking Docker...")
	
	// Check if Docker is installed
	if _, err := exec.LookPath("docker"); err != nil {
		log.Println("  ⚠️  Docker not found")
		return
	}
	
	// Check if Docker daemon is running
	if _, err := exec.Command("docker", "info").Output(); err != nil {
		log.Println("  ⚠️  Docker daemon not running")
		return
	}
	
	fmt.Println("  ✓ Docker is running")
}

func checkGoVersion() {
	fmt.Println("  Checking Go version...")
	
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("  ⚠️  Go not found: %v", err)
		return
	}
	
	versionStr := string(output)
	fmt.Printf("  ✓ Go: %s", strings.TrimSpace(versionStr))
}

func showSystemInfo() {
	fmt.Println("📋 System Information")
	fmt.Println("======================")
	
	osInfo, err := getOSInfo()
	if err != nil {
		log.Printf("Warning: Could not get OS info: %v\n", err)
		return
	}
	
	fmt.Printf("OS: %s %s\n", osInfo.Name, osInfo.Version)
	
	arch := "unknown"
	if runtime.GOOS == "linux" {
		cmd := exec.Command("uname", "-m")
		output, _ := cmd.Output()
		arch = strings.TrimSpace(string(output))
	}
	fmt.Printf("Architecture: %s\n", arch)
	
	// Show Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err == nil {
		fmt.Printf("Go: %s", strings.TrimSpace(string(output)))
	}
	
	// Show Docker version
	if _, err := exec.LookPath("docker"); err == nil {
		cmd := exec.Command("docker", "--version")
		output, err := cmd.Output()
		if err == nil {
			fmt.Printf("Docker: %s", strings.TrimSpace(string(output)))
		}
	}
}

func showSystemResources() {
	fmt.Println("💻 System Resources")
	fmt.Println("===================")
	
	// Get CPU info
	fmt.Printf("CPU: %d cores\n", runtime.NumCPU())
	
	// Get memory info
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/meminfo")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "MemTotal:") {
					var total uint64
					fmt.Sscanf(line, "MemTotal: %d kB", &total)
					fmt.Printf("Memory: %.2f GB total\n", float64(total)/1024/1024)
					break
				}
			}
		}
	}
}

type SystemInfo struct {
	Name    string
	Version string
}

// ============================================
// USB Installer Commands
// ============================================

func handleUSBCommand(args []string) {
if len(args) == 0 {
fmt.Println("Error: usb command required")
os.Exit(1)
}

switch args[0] {
case "list":
usb.ListUSBDevices()
case "prepare":
if len(args) < 3 {
fmt.Println("Usage: sonic usb prepare <device> <distro>")
fmt.Println("Distros: ubuntu, mint, classicmodern")
os.Exit(1)
}
device := args[1]
distroName := args[2]
layout, err := usb.GetLayout(distroName)
if err != nil {
log.Fatalf("Error: %v", err)
}
if !usb.Confirm(device) {
fmt.Println("Aborted")
os.Exit(0)
}
config := usb.InstallConfig{
Device: device,
DistroName: distroName,
Layout: layout,
}
if err := usb.PrepareDisk(config); err != nil {
log.Fatalf("Prepare failed: %v", err)
}
case "install":
if len(args) < 3 {
fmt.Println("Usage: sonic usb install <device> <distro>")
fmt.Println("Distros: ubuntu, mint, classicmodern")
os.Exit(1)
}
device := args[1]
distroName := args[2]
layout, err := usb.GetLayout(distroName)
if err != nil {
log.Fatalf("Error: %v", err)
}
if !usb.Confirm(device) {
fmt.Println("Aborted")
os.Exit(0)
}
config := usb.InstallConfig{
Device: device,
DistroName: distroName,
Layout: layout,
}
result, err := usb.FullInstall(config)
if err != nil {
log.Fatalf("Install failed: %v", err)
}
fmt.Printf("\n✅ %s installed to %s\n", result.Distro, result.Device)
case "write":
if len(args) < 3 {
fmt.Println("Usage: sonic usb write <device> <iso-path>")
os.Exit(1)
}
device := args[1]
isoPath := args[2]
if !usb.Confirm(device) {
fmt.Println("Aborted")
os.Exit(0)
}
if err := iso.WriteISOToDisk(isoPath, device); err != nil {
log.Fatalf("Write failed: %v", err)
}
default:
fmt.Printf("Unknown usb command: %s\n", args[0])
os.Exit(1)
}
}

// ============================================
// ISO Download Commands
// ============================================

func handleISOCommand(args []string) {
if len(args) == 0 {
fmt.Println("Error: iso command required")
os.Exit(1)
}

switch args[0] {
case "list":
fmt.Println("\nAvailable Distros:")
for _, d := range iso.ListDistros() {
fmt.Printf("  %-20s v%-8s %s (%s)\n", d.Name, d.Version, d.Size, d.Arch)
}
case "download":
if len(args) < 2 {
fmt.Println("Usage: sonic iso download <distro>")
os.Exit(1)
}
distro, err := iso.GetDistro(args[1])
if err != nil {
log.Fatalf("Error: %v", err)
}
progressCh := make(chan iso.DownloadStatus, 10)
go func() {
for status := range progressCh {
if status.Complete {
break
}
fmt.Printf("\r  Downloading: %.1f%%", status.Progress)
}
}()
isoPath, err := iso.Download(distro, progressCh)
close(progressCh)
if err != nil {
log.Fatalf("Download failed: %v", err)
}
fmt.Printf("\n✅ Downloaded: %s\n", isoPath)
case "cache":
cacheDir := iso.GetCacheDir()
fmt.Printf("ISO Cache: %s\n", cacheDir)
if entries, err := os.ReadDir(cacheDir); err == nil {
fmt.Printf("Cached ISOs: %d\n", len(entries))
for _, e := range entries {
info, _ := e.Info()
fmt.Printf("  %s (%s)\n", e.Name(), formatBytes(uint64(info.Size())))
}
}
default:
fmt.Printf("Unknown iso command: %s\n", args[0])
os.Exit(1)
}
}

// ============================================
// Disk Operations Commands
// ============================================

func handleDiskCommand(args []string) {
if len(args) == 0 {
fmt.Println("Error: disk command required")
os.Exit(1)
}

switch args[0] {
case "list":
devices, err := disk.DetectDevices(false)
if err != nil {
log.Fatalf("Error: %v", err)
}
fmt.Println("\nBlock Devices:")
for _, d := range devices {
removable := ""
if d.Removable {
removable = " [USB]"
}
fmt.Printf("  %-12s %-8s %s%s\n", d.Path, d.Size, d.Model, removable)
}
case "info":
if len(args) < 2 {
fmt.Println("Usage: sonic disk info <device>")
os.Exit(1)
}
devices, err := disk.DetectDevices(false)
if err != nil {
log.Fatalf("Error: %v", err)
}
for _, d := range devices {
if d.Path == args[1] || d.Name == args[1] {
fmt.Printf("\nDevice: %s\n", d.Path)
fmt.Printf("  Size:      %s\n", d.Size)
fmt.Printf("  Model:     %s\n", d.Model)
fmt.Printf("  Removable: %v\n", d.Removable)
return
}
}
fmt.Printf("Device not found: %s\n", args[1])
case "wipe":
if len(args) < 2 {
fmt.Println("Usage: sonic disk wipe <device>")
os.Exit(1)
}
device := args[1]
if !usb.Confirm(device) {
fmt.Println("Aborted")
os.Exit(0)
}
if err := disk.WipeDevice(device); err != nil {
log.Fatalf("Wipe failed: %v", err)
}
fmt.Println("✅ Disk wiped")
default:
fmt.Printf("Unknown disk command: %s\n", args[0])
os.Exit(1)
}
}

func formatBytes(bytes uint64) string {
const unit = 1024
if bytes < unit {
return fmt.Sprintf("%dB", bytes)
}
div, exp := uint64(unit), 0
for n := bytes / unit; n >= unit; n /= unit {
div *= unit
exp++
}
return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
