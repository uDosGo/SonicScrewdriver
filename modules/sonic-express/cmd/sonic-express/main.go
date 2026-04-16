package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sonic-family/sonic-express/pkg/bundle"
	"github.com/sonic-family/sonic-express/pkg/manifest"
)

const version = "v0.1.0-dev"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Println(version)
	case "pack":
		runPack(os.Args[2:])
	case "verify":
		runVerify(os.Args[2:])
	case "install":
		fmt.Println("install command is scaffolded in cmd/install")
	case "serve":
		fmt.Println("serve command is scaffolded in cmd/serve")
	case "help", "--help", "-h":
		printHelp()
	default:
		printHelp()
		os.Exit(2)
	}
}

func runPack(args []string) {
	fs := flag.NewFlagSet("pack", flag.ExitOnError)
	source := fs.String("source", ".", "source directory")
	output := fs.String("output", "manifest.dryrun.json", "manifest output")
	versionFlag := fs.String("version", "0.0.0-dev", "bundle version")
	channel := fs.String("channel", "edge", "release channel")
	dryRun := fs.Bool("dry-run", true, "generate dry-run manifest")

	_ = fs.Parse(args)
	m, err := manifest.BuildDryRun(manifest.BuildInput{
		SourcePath: *source,
		Version:    *versionFlag,
		Channel:    *channel,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pack failed: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		if err := manifest.Write(*output, m); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write manifest: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("pack dry-run complete: %s\n", *output)
		return
	}

	if err := bundle.BuildDraftBundle(*output, m); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write draft bundle: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("pack draft bundle complete: %s\n", *output)
}

func runVerify(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: sonic-express verify <bundle.she>")
		os.Exit(2)
	}
	if err := bundle.VerifyDraftBundle(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "verify failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("verify ok: %s\n", args[0])
}

func printHelp() {
	fmt.Println(`sonic-express v0.1.0-dev

Usage:
  sonic-express version
  sonic-express pack --dry-run --source . --output manifest.dryrun.json
  sonic-express pack --dry-run=false --output uhome-nest-draft.she
  sonic-express verify ./uhome-nest-draft.she
  sonic-express install [bundle.she] [--usb]
  sonic-express serve --port 8080 --bundles ./bundles`)
}
