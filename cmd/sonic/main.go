// Command sonic is the LEGACY v1 CLI for the SonicScrewdriver system (Go).
// This is preserved for backward compatibility. The active v2 CLI is in cli/ (Python).
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
	"github.com/uDosGo/SonicScrewdriver/pkg/classicmodern"
	"github.com/uDosGo/SonicScrewdriver/pkg/container"
	"github.com/uDosGo/SonicScrewdriver/pkg/disk"
	"github.com/uDosGo/SonicScrewdriver/pkg/driver"
	"github.com/uDosGo/SonicScrewdriver/pkg/gui"
	"github.com/uDosGo/SonicScrewdriver/pkg/iso"
	"github.com/uDosGo/SonicScrewdriver/pkg/knowledge"
	"github.com/uDosGo/SonicScrewdriver/pkg/library"
	"github.com/uDosGo/SonicScrewdriver/pkg/recovery"
	"github.com/uDosGo/SonicScrewdriver/pkg/reflash"
	"github.com/uDosGo/SonicScrewdriver/pkg/remote"
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
	case "remote", "r":
		runRemote(args)
	case "mint":
		runMint(args)
	case "reflash", "rf":
		runReflash(args)
	case "driver", "dr":
		runDriver(args)
	case "recovery", "rec":
		runRecovery(args)
	case "help", "h", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`SonicScrewdriver — Unified System Toolkit

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
  remote, r        Remote access (vnc, ssh, samba, info)
  mint             Classic Modern Mint readiness (check, doctor, info, apply)
  ventoy           Ventoy bundle packaging (create, validate, info)
  reflash, rf      Device firmware flashing (list, flash, devices, download)
  driver, dr       Driver tracking and version management (list, register, check)
  recovery, rec    Device backup and restore (backup, restore, list, export)
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
  sonic reflash list
  sonic reflash flash --device /dev/ttyUSB0 --firmware firmware.bin
  sonic driver list
  sonic recovery backup --device /dev/sda --type firmware
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

// ─── Remote ───────────────────────────────────────────────────────────────────

func runRemote(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic remote <vnc|ssh|samba|info> [flags]")
		return
	}

	switch args[0] {
	case "vnc":
		port := 5901
		password := "sonic"
		geometry := "1280x720"
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--port":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &port)
					i++
				}
			case "--password":
				if i+1 < len(args) {
					password = args[i+1]
					i++
				}
			case "--geometry":
				if i+1 < len(args) {
					geometry = args[i+1]
					i++
				}
			}
		}

		server := remote.NewVNCServer(port, password, geometry)
		if err := server.SetupVNC(); err != nil {
			log.Fatalf("VNC setup failed: %v", err)
		}
		if err := server.StartVNC(); err != nil {
			log.Fatalf("VNC start failed: %v", err)
		}
		fmt.Println("VNC server started.")

	case "ssh":
		if err := remote.SetupSSH(); err != nil {
			log.Fatalf("SSH setup failed: %v", err)
		}
		fmt.Println("SSH service running.")

	case "samba":
		shareName := "shared"
		sharePath := "/home/shared"
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--name":
				if i+1 < len(args) {
					shareName = args[i+1]
					i++
				}
			case "--path":
				if i+1 < len(args) {
					sharePath = args[i+1]
					i++
				}
			}
		}
		if err := remote.SetupSamba(shareName, sharePath); err != nil {
			log.Fatalf("Samba setup failed: %v", err)
		}
		fmt.Println("Samba file sharing set up.")

	case "info":
		fmt.Println(remote.GetRemoteAccessInfo())

	default:
		fmt.Printf("Unknown remote action: %s\n", args[0])
	}
}

// ─── Mint ─────────────────────────────────────────────────────────────────────

func runMint(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic mint <check|doctor|info|apply> [component]")
		return
	}

	// Classic Modern Mint doesn't need a DB connection for basic operations
	checker := classicmodern.NewReadinessChecker(nil)

	switch args[0] {
	case "check":
		report, err := checker.GenerateInstallationReport()
		if err != nil {
			log.Fatalf("Readiness check failed: %v", err)
		}
		fmt.Println(report)

	case "doctor":
		component := "classic-modern"
		if len(args) > 1 {
			component = args[1]
		}
		check, err := checker.PerformCheck()
		if err != nil {
			log.Fatalf("Doctor check failed: %v", err)
		}
		if check.ThemeReady && check.IconsReady && check.FontsReady && check.OBFValid && check.DependenciesReady {
			fmt.Printf("✅ Component '%s' is ready.\n", component)
		} else {
			fmt.Printf("⚠️  Component '%s' has issues:\n", component)
			for _, issue := range check.Issues {
				fmt.Printf("  - %s\n", issue)
			}
		}

	case "info":
		info, err := checker.GetThemeInfo()
		if err != nil {
			log.Fatalf("Failed to get theme info: %v", err)
		}
		b, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(b))

	case "apply":
		if err := checker.ApplyTheme(); err != nil {
			log.Fatalf("Failed to apply theme: %v", err)
		}
		fmt.Println("Classic Modern theme applied.")

	default:
		fmt.Printf("Unknown mint action: %s\n", args[0])
	}
}

// ─── Reflash ──────────────────────────────────────────────────────────────────

func runReflash(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: sonic reflash <list|flash|devices|download> [flags]")
		return
	}

	flasher := reflash.New()

	switch args[0] {
	case "list", "ls":
		devices, err := flasher.ListSupportedDevices()
		if err != nil {
			log.Fatalf("Failed to list devices: %v", err)
		}
		if len(devices) == 0 {
			fmt.Println("No supported devices in library.")
			return
		}
		fmt.Println("Supported Devices:")
		for _, d := range devices {
			fmt.Printf("  %s [%s] — %s %s\n", d.Name, d.Type, d.Manufacturer, d.Model)
			fmt.Printf("    Methods: %v | Protocols: %v\n", d.Methods, d.Protocols)
			fmt.Printf("    Max Size: %s | Voltage: %s\n", d.MaxSize, d.Voltage)
		}

	case "flash":
		config := reflash.FlashConfig{}
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--device":
				if i+1 < len(args) {
					config.Device = args[i+1]
					i++
				}
			case "--firmware":
				if i+1 < len(args) {
					config.Firmware = args[i+1]
					i++
				}
			case "--method":
				if i+1 < len(args) {
					config.Method = reflash.FlashMethod(args[i+1])
					i++
				}
			case "--offset":
				if i+1 < len(args) {
					config.Offset = args[i+1]
					i++
				}
			case "--baud":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &config.BaudRate)
					i++
				}
			case "--probe":
				if i+1 < len(args) {
					config.Probe = args[i+1]
					i++
				}
			case "--verify":
				config.Verify = true
			case "--dry-run":
				config.DryRun = true
			}
		}
		if config.Device == "" {
			log.Fatal("--device is required")
		}
		if config.Firmware == "" {
			log.Fatal("--firmware is required")
		}
		result, err := flasher.Flash(config)
		if err != nil {
			log.Fatalf("Flash failed: %v", err)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))

	case "devices":
		devices, err := flasher.ListSupportedDevices()
		if err != nil {
			log.Fatalf("Failed to list devices: %v", err)
		}
		b, _ := json.MarshalIndent(devices, "", "  ")
		fmt.Println(string(b))

	case "download":
		if len(args) < 2 {
			log.Fatal("Usage: sonic reflash download <url> [dest-dir]")
		}
		destDir := "."
		if len(args) > 2 {
			destDir = args[2]
		}
		path, err := flasher.DownloadFirmware(args[1], destDir)
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		fmt.Printf("Firmware downloaded to: %s\n", path)

	default:
		fmt.Printf("Unknown reflash action: %s\n", args[0])
	}
}

// ─── Driver ───────────────────────────────────────────────────────────────────

func runDriver(args []string) {
	manager := driver.NewManager()
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load driver database: %v", err)
	}

	if len(args) < 1 {
		fmt.Println("Usage: sonic driver <list|register|check|detect|info> [name]")
		return
	}

	switch args[0] {
	case "list", "ls":
		drivers := manager.List("")
		if len(drivers) == 0 {
			// Load defaults
			for _, d := range driver.GetDefaultDrivers() {
				manager.Register(d)
			}
			drivers = manager.List("")
		}
		fmt.Println("Registered Drivers:")
		for _, d := range drivers {
			status := "❌"
			if d.Installed {
				status = "✅"
			}
			fmt.Printf("  %s %s [%s] v%s — %s\n", status, d.Name, d.Type, d.Version, d.Description)
		}

	case "register":
		if len(args) < 2 {
			log.Fatal("Usage: sonic driver register <name>")
		}
		// Look up in defaults
		for _, d := range driver.GetDefaultDrivers() {
			if d.Name == args[1] {
				if err := manager.Register(d); err != nil {
					log.Fatalf("Failed to register driver: %v", err)
				}
				fmt.Printf("Driver '%s' registered.\n", args[1])
				return
			}
		}
		log.Fatalf("Driver '%s' not found in default catalogue", args[1])

	case "check":
		if len(args) < 2 {
			log.Fatal("Usage: sonic driver check <name>")
		}
		installed, err := manager.CheckInstalled(args[1])
		if err != nil {
			log.Fatalf("Check failed: %v", err)
		}
		if installed {
			version, _ := manager.GetVersion(args[1])
			fmt.Printf("✅ Driver '%s' is installed (v%s)\n", args[1], version)
		} else {
			fmt.Printf("❌ Driver '%s' is not installed\n", args[1])
		}

	case "detect":
		found, err := manager.DetectInstalledDrivers()
		if err != nil {
			log.Fatalf("Detection failed: %v", err)
		}
		if len(found) == 0 {
			fmt.Println("No known drivers detected on system.")
			return
		}
		fmt.Println("Detected installed drivers:")
		for _, name := range found {
			fmt.Printf("  ✅ %s\n", name)
		}

	case "info":
		if len(args) < 2 {
			log.Fatal("Usage: sonic driver info <name>")
		}
		d, err := manager.Get(args[1])
		if err != nil {
			log.Fatalf("Driver not found: %v", err)
		}
		b, _ := json.MarshalIndent(d, "", "  ")
		fmt.Println(string(b))

	default:
		fmt.Printf("Unknown driver action: %s\n", args[0])
	}
}

// ─── Recovery ─────────────────────────────────────────────────────────────────

func runRecovery(args []string) {
	manager := recovery.NewManager()
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load backup index: %v", err)
	}

	if len(args) < 1 {
		fmt.Println("Usage: sonic recovery <backup|restore|list|export|import|delete> [flags]")
		return
	}

	switch args[0] {
	case "backup":
		device := ""
		backupType := recovery.BackupFirmware
		description := ""
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--device":
				if i+1 < len(args) {
					device = args[i+1]
					i++
				}
			case "--type":
				if i+1 < len(args) {
					backupType = recovery.BackupType(args[i+1])
					i++
				}
			case "--description":
				if i+1 < len(args) {
					description = args[i+1]
					i++
				}
			}
		}
		if device == "" {
			log.Fatal("--device is required")
		}
		record, err := manager.BackupDevice(device, backupType, description)
		if err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		b, _ := json.MarshalIndent(record, "", "  ")
		fmt.Println(string(b))

	case "restore":
		if len(args) < 2 {
			log.Fatal("Usage: sonic recovery restore <backup-id> [target-device]")
		}
		target := args[1]
		if len(args) > 2 {
			target = args[2]
		}
		if err := manager.RestoreDevice(args[1], target); err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
		fmt.Printf("Device restored from backup '%s'.\n", args[1])

	case "list", "ls":
		filter := ""
		if len(args) > 1 {
			filter = args[1]
		}
		records := manager.ListBackups(filter)
		if len(records) == 0 {
			fmt.Println("No backups found.")
			return
		}
		fmt.Println("Backups:")
		for _, r := range records {
			fmt.Printf("  %s — %s [%s] %s (%s)\n", r.ID, r.Device, r.Type, r.Description, r.CreatedAt)
		}

	case "export":
		if len(args) < 3 {
			log.Fatal("Usage: sonic recovery export <backup-id> <output-path>")
		}
		if err := manager.ExportBackup(args[1], args[2]); err != nil {
			log.Fatalf("Export failed: %v", err)
		}
		fmt.Printf("Backup '%s' exported to %s.\n", args[1], args[2])

	case "import":
		if len(args) < 3 {
			log.Fatal("Usage: sonic recovery import <file-path> <device-name> [backup-type]")
		}
		backupType := recovery.BackupFirmware
		if len(args) > 3 {
			backupType = recovery.BackupType(args[3])
		}
		record, err := manager.ImportBackup(args[1], args[2], backupType)
		if err != nil {
			log.Fatalf("Import failed: %v", err)
		}
		fmt.Printf("Backup imported: %s\n", record.ID)

	case "delete":
		if len(args) < 2 {
			log.Fatal("Usage: sonic recovery delete <backup-id>")
		}
		if err := manager.DeleteBackup(args[1]); err != nil {
			log.Fatalf("Delete failed: %v", err)
		}
		fmt.Printf("Backup '%s' deleted.\n", args[1])

	default:
		fmt.Printf("Unknown recovery action: %s\n", args[0])
	}
}

// Ensure we reference disk and iso for compilation
var _ = disk.DetectDevices
var _ = iso.GetDistro
var _ = strings.TrimSpace
