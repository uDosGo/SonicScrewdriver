// Command sonic is the unified CLI for the SonicScrewdriver system.
// It provides container management, USB installation, secret vault access,
// device catalogue browsing, knowledge queries, and a web GUI.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/uDosGo/SonicScrewdriver/pkg/catalogue"
	"github.com/uDosGo/SonicScrewdriver/pkg/container"
	"github.com/uDosGo/SonicScrewdriver/pkg/disk"
	"github.com/uDosGo/SonicScrewdriver/pkg/gui"
	"github.com/uDosGo/SonicScrewdriver/pkg/iso"
	"github.com/uDosGo/SonicScrewdriver/pkg/knowledge"
	"github.com/uDosGo/SonicScrewdriver/pkg/library"
	"github.com/uDosGo/SonicScrewdriver/pkg/usb"
	"github.com/uDosGo/SonicScrewdriver/pkg/vault"
	"github.com/uDosGo/SonicScrewdriver/pkg/ventoy"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "container", "c":
		runContainer(args)
	case "usb", "u":
		runUSB(args)
	case "vault", "v":
		runVault(args)
	case "gui", "g":
		runGUI(args)
	case "catalogue", "cat":
		runCatalogue(args)
	case "knowledge", "k":
		runKnowledge(args)
	case "library", "lib":
		runLibrary(args)
	case "ventoy":
		runVentoy(args)
	case "help", "h", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`SonicScrewdriver — Unified System Toolkit

Usage:
  sonic <command> [arguments]

Commands:
  container, c     Manage containers (list, start, stop, restart, remove, health)
  usb, u           USB installation (list, prepare, install, full-install)
  vault, v         Secret vault (get, set, list, rotate, history)
  gui, g           Start web GUI dashboard
  catalogue, cat   Browse uCode device catalogue (list, find)
  knowledge, k     Query knowledge sources (sources, query)
  library, lib     Library index management (list, info, validate)
  ventoy           Ventoy bundle packaging (create, validate, info)
  help, h          Show this help

Examples:
  sonic container list
  sonic container start my-container
  sonic usb list
  sonic usb install --device /dev/sdb --distro ubuntu
  sonic vault set github_token ghp_xxxx
  sonic vault get github_token
  sonic gui
  sonic catalogue list
  sonic knowledge query "docker setup"
`)
}

// ─── Container ────────────────────────────────────────────────────────────────

func runContainer(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic container <list|start|stop|restart|remove|health> [name]")
		return
	}

	runtime, err := container.NewDockerRuntime()
	if err != nil {
		log.Fatalf("Failed to create runtime: %v", err)
	}

	action := args[0]
	switch action {
	case "list", "ls":
		containers, err := runtime.List()
		if err != nil {
			log.Fatalf("Failed to list containers: %v", err)
		}
		if len(containers) == 0 {
			fmt.Println("No containers found.")
			return
		}
		fmt.Println("Containers:")
		for _, name := range containers {
			health, err := runtime.CheckContainerHealth(name)
			status := "unknown"
			if err == nil && health != nil {
				status = health.Status
			}
			fmt.Printf("  %s [%s]\n", name, status)
		}

	case "start":
		if len(args) < 2 {
			log.Fatal("Usage: sonic container start <name>")
		}
		if err := runtime.Start(args[1]); err != nil {
			log.Fatalf("Failed to start container: %v", err)
		}
		fmt.Printf("Started container: %s\n", args[1])

	case "stop":
		if len(args) < 2 {
			log.Fatal("Usage: sonic container stop <name>")
		}
		if err := runtime.Stop(args[1]); err != nil {
			log.Fatalf("Failed to stop container: %v", err)
		}
		fmt.Printf("Stopped container: %s\n", args[1])

	case "restart":
		if len(args) < 2 {
			log.Fatal("Usage: sonic container restart <name>")
		}
		if err := runtime.RestartContainer(args[1]); err != nil {
			log.Fatalf("Failed to restart container: %v", err)
		}
		fmt.Printf("Restarted container: %s\n", args[1])

	case "remove", "rm":
		if len(args) < 2 {
			log.Fatal("Usage: sonic container remove <name>")
		}
		if err := runtime.Remove(args[1]); err != nil {
			log.Fatalf("Failed to remove container: %v", err)
		}
		fmt.Printf("Removed container: %s\n", args[1])

	case "health":
		healthStatuses, err := runtime.GetAllContainerHealth()
		if err != nil {
			log.Fatalf("Failed to get health: %v", err)
		}
		b, _ := json.MarshalIndent(healthStatuses, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown container action: %s\n", action)
	}
}

// ─── USB ──────────────────────────────────────────────────────────────────────

func runUSB(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic usb <list|prepare|install|full-install> [flags]")
		return
	}

	action := args[0]
	switch action {
	case "list", "ls":
		devices, err := usb.ListUSBDevices()
		if err != nil {
			log.Fatalf("Failed to list USB devices: %v", err)
		}
		if len(devices) == 0 {
			fmt.Println("No USB devices found.")
		}

	case "prepare":
		config := parseUSBConfig(args[1:])
		if config.Device == "" {
			log.Fatal("--device is required")
		}
		if config.DistroName == "" {
			log.Fatal("--distro is required")
		}
		layout, err := usb.GetLayout(config.DistroName)
		if err != nil {
			log.Fatalf("Unknown distro: %v", err)
		}
		config.Layout = layout
		if err := usb.PrepareDisk(config); err != nil {
			log.Fatalf("Disk preparation failed: %v", err)
		}

	case "install":
		config := parseUSBConfig(args[1:])
		if config.Device == "" {
			log.Fatal("--device is required")
		}
		if config.DistroName == "" {
			log.Fatal("--distro is required")
		}
		result, err := usb.InstallDistro(config)
		if err != nil {
			log.Fatalf("Installation failed: %v", err)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))

	case "full-install", "full":
		config := parseUSBConfig(args[1:])
		if config.Device == "" {
			log.Fatal("--device is required")
		}
		if config.DistroName == "" {
			log.Fatal("--distro is required")
		}
		layout, err := usb.GetLayout(config.DistroName)
		if err != nil {
			log.Fatalf("Unknown distro: %v", err)
		}
		config.Layout = layout
		result, err := usb.FullInstall(config)
		if err != nil {
			log.Fatalf("Full install failed: %v", err)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown USB action: %s\n", action)
	}
}

func parseUSBConfig(args []string) usb.InstallConfig {
	config := usb.InstallConfig{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--device":
			if i+1 < len(args) {
				config.Device = args[i+1]
				i++
			}
		case "--distro":
			if i+1 < len(args) {
				config.DistroName = args[i+1]
				i++
			}
		case "--dry-run":
			config.DryRun = true
		case "--force":
			config.Force = true
		}
	}
	return config
}

// ─── Vault ────────────────────────────────────────────────────────────────────

func runVault(args []string) {
	v, err := vault.Open()
	if err != nil {
		log.Fatalf("Failed to open vault: %v", err)
	}

	if len(args) < 1 {
		fmt.Println("Usage: sonic vault <get|set|list|rotate|history> [key] [value]")
		return
	}

	switch args[0] {
	case "get":
		if len(args) < 2 {
			log.Fatal("Usage: sonic vault get <key>")
		}
		value, err := v.Get(args[1])
		if err != nil {
			log.Fatalf("Failed to get secret: %v", err)
		}
		fmt.Println(value)

	case "set":
		if len(args) < 3 {
			log.Fatal("Usage: sonic vault set <key> <value>")
		}
		if err := v.Set(args[1], args[2]); err != nil {
			log.Fatalf("Failed to set secret: %v", err)
		}
		fmt.Printf("Secret '%s' set successfully\n", args[1])

	case "list", "ls":
		secrets, err := v.List()
		if err != nil {
			log.Fatalf("Failed to list secrets: %v", err)
		}
		if len(secrets) == 0 {
			fmt.Println("No secrets stored.")
			return
		}
		fmt.Println("Secrets:")
		for _, s := range secrets {
			fmt.Printf("  %s\n", s)
		}

	case "rotate":
		if len(args) < 3 {
			log.Fatal("Usage: sonic vault rotate <key> <new-value>")
		}
		if err := v.Rotate(args[1], args[2]); err != nil {
			log.Fatalf("Failed to rotate secret: %v", err)
		}
		fmt.Printf("Secret '%s' rotated successfully\n", args[1])

	case "history":
		if len(args) < 2 {
			log.Fatal("Usage: sonic vault history <key>")
		}
		history, err := v.History(args[1])
		if err != nil {
			log.Fatalf("Failed to get history: %v", err)
		}
		b, _ := json.MarshalIndent(history, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown vault action: %s\n", args[0])
	}
}

// ─── GUI ──────────────────────────────────────────────────────────────────────

func runGUI(args []string) {
	port := 8080
	for i := 0; i < len(args); i++ {
		if args[i] == "--port" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &port)
			i++
		}
	}

	runtime, err := container.NewDockerRuntime()
	if err != nil {
		log.Fatalf("Failed to create runtime: %v", err)
	}

	server := gui.NewServer(runtime, port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start GUI: %v", err)
	}

	fmt.Printf("SonicScrewdriver GUI running at http://localhost:%d\n", server.Port())
	fmt.Println("Press Ctrl+C to stop.")

	// Block forever
	select {}
}

// ─── Catalogue ────────────────────────────────────────────────────────────────

func runCatalogue(args []string) {
	cat := catalogue.New("")

	if len(args) < 1 {
		fmt.Println("Usage: sonic catalogue <list|find> [name]")
		return
	}

	switch args[0] {
	case "list", "ls":
		devices, err := cat.ListDevices()
		if err != nil {
			log.Fatalf("Failed to list devices: %v", err)
		}
		if len(devices) == 0 {
			fmt.Println("No devices found in catalogue.")
			return
		}
		fmt.Println("Device Catalogue:")
		for _, d := range devices {
			fmt.Printf("  %s [%s] (%s)\n", d.Name, d.Type, d.Repo)
		}

	case "find":
		if len(args) < 2 {
			log.Fatal("Usage: sonic catalogue find <name>")
		}
		device, err := cat.FindDevice(args[1])
		if err != nil {
			log.Fatalf("Device not found: %v", err)
		}
		b, _ := json.MarshalIndent(device, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown catalogue action: %s\n", args[0])
	}
}

// ─── Knowledge ────────────────────────────────────────────────────────────────

func runKnowledge(args []string) {
	k := knowledge.New("")

	if len(args) < 1 {
		fmt.Println("Usage: sonic knowledge <sources|query> [term]")
		return
	}

	switch args[0] {
	case "sources", "ls":
		sources, err := k.ListSources()
		if err != nil {
			log.Fatalf("Failed to list sources: %v", err)
		}
		if len(sources) == 0 {
			fmt.Println("No knowledge sources found.")
			return
		}
		fmt.Println("Knowledge Sources:")
		for _, s := range sources {
			fmt.Printf("  %s [%s] — %s\n", s.Name, s.Type, s.Description)
		}

	case "query", "q":
		if len(args) < 2 {
			log.Fatal("Usage: sonic knowledge query <term>")
		}
		results, err := k.Query(args[1])
		if err != nil {
			log.Fatalf("Query failed: %v", err)
		}
		if len(results) == 0 {
			fmt.Printf("No results found for '%s'.\n", args[1])
			return
		}
		fmt.Printf("Results for '%s':\n", args[1])
		for _, r := range results {
			fmt.Println("  " + r)
		}

	default:
		fmt.Printf("Unknown knowledge action: %s\n", args[0])
	}
}

// ─── Library ──────────────────────────────────────────────────────────────────

func runLibrary(args []string) {
	manager := library.NewManager(library.GetDefaultIndexPath())

	if len(args) < 1 {
		fmt.Println("Usage: sonic library <list|info|validate> [name]")
		return
	}

	switch args[0] {
	case "list", "ls":
		if err := manager.Load(); err != nil {
			log.Fatalf("Failed to load library: %v", err)
		}
		entries, err := manager.List()
		if err != nil {
			log.Fatalf("Failed to list library: %v", err)
		}
		if len(entries) == 0 {
			fmt.Println("Library is empty.")
			return
		}
		fmt.Println("Library Index:")
		for _, e := range entries {
			fmt.Printf("  %s — %s\n", e.Name, e.Status)
		}

	case "info":
		if len(args) < 2 {
			log.Fatal("Usage: sonic library info <name>")
		}
		if err := manager.Load(); err != nil {
			log.Fatalf("Failed to load library: %v", err)
		}
		entry, err := manager.GetGame(args[1])
		if err != nil {
			log.Fatalf("Failed to get info: %v", err)
		}
		b, _ := json.MarshalIndent(entry, "", "  ")
		fmt.Println(string(b))

	case "validate":
		if len(args) < 2 {
			log.Fatal("Usage: sonic library validate <manifest-path>")
		}
		if err := manager.ValidateManifest(args[1]); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		fmt.Println("Manifest is valid.")

	default:
		fmt.Printf("Unknown library action: %s\n", args[0])
	}
}

// ─── Ventoy ───────────────────────────────────────────────────────────────────

func runVentoy(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic ventoy <create|validate|info> [flags]")
		return
	}

	switch args[0] {
	case "create":
		source := "."
		output := "."
		name := "bundle"
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--source":
				if i+1 < len(args) {
					source = args[i+1]
					i++
				}
			case "--output":
				if i+1 < len(args) {
					output = args[i+1]
					i++
				}
			case "--name":
				if i+1 < len(args) {
					name = args[i+1]
					i++
				}
			}
		}

		manager := library.NewManager("")
		packager := ventoy.NewPackager(source, output, manager)
		bundlePath, err := packager.CreateBundle(name)
		if err != nil {
			log.Fatalf("Failed to create bundle: %v", err)
		}
		fmt.Printf("Bundle created: %s\n", bundlePath)

	case "validate":
		if len(args) < 2 {
			log.Fatal("Usage: sonic ventoy validate <bundle-path>")
		}
		manager := library.NewManager("")
		packager := ventoy.NewPackager("", "", manager)
		if err := packager.ValidateBundle(args[1]); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		fmt.Println("Bundle is valid.")

	case "info":
		if len(args) < 2 {
			log.Fatal("Usage: sonic ventoy info <bundle-path>")
		}
		info, err := ventoy.GetBundleInfo(args[1])
		if err != nil {
			log.Fatalf("Failed to get bundle info: %v", err)
		}
		b, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown ventoy action: %s\n", args[0])
	}
}

// Ensure we reference disk and iso for compilation
var _ = disk.DetectDevices
var _ = iso.GetDistro
var _ = strings.TrimSpace
