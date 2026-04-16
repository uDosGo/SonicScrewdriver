package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "install":
		fmt.Println("Installing...")
	case "start":
		fmt.Println("Starting...")
	case "stop":
		fmt.Println("Stopping...")
	case "list":
		fmt.Println("Listing...")
	case "remove":
		fmt.Println("Removing...")
	case "--help", "-h", "help":
		printHelp()
	case "--version", "-v", "version":
		fmt.Println("vA1.0.0")
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
