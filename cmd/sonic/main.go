package main

import (
	"fmt"
	"log"
	"os"
	"github.com/sonic-family/sonic-screwdriver/internal/container"
	"github.com/sonic-family/sonic-screwdriver/internal/library"
	"github.com/sonic-family/sonic-screwdriver/internal/state"
)

var (
	runtime  container.Runtime
	libManager *library.Manager
	stateDB   *state.DB
)

func main() {
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
		runtime = &container.DockerRuntime{} // Fallback to mock
	} else {
		runtime = dockerRuntime
		defer dockerRuntime.Close()
	}

	// Initialize state database
	dbPath := state.GetDefaultDBPath()
	stateDB, err = state.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open state database: %v", err)
	}
	defer stateDB.Close()

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
		
		err := runtime.Start(gameName)
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
		
		err := runtime.Stop(gameName)
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
		
		err := runtime.Remove(gameName)
		if err != nil {
			log.Fatalf("Failed to remove: %v", err)
		}
		
		fmt.Printf("Removed %s\n", gameName)
	case "--help", "-h", "help":
		printHelp()
	case "--version", "-v", "version":
		fmt.Println("vA1.0.0")
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
	case "config":
		if len(os.Args) < 5 || os.Args[2] != "set" {
			fmt.Println("Usage: sonic config set <key> <value>")
			os.Exit(1)
		}
		fmt.Printf("Setting %s=%s\n", os.Args[3], os.Args[4])
		// TODO: Implement config management
	default:
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`SonicScrewdriver vA1.0.0

Usage:
  sonic install <game>     Install a game from curated library
  sonic start <game>       Start a game container
  sonic stop <game>        Stop a game container
  sonic list               List installed games
  sonic remove <game>      Remove a game

Commands:
  sonic library list       Show available curated games
  sonic logs <game>        Show container logs
  sonic config set <key> <value>

Flags:
  --help, -h               Show this help
  --version, -v            Show version`)
}
