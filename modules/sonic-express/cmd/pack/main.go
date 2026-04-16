package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sonic-family/sonic-express/pkg/bundle"
	"github.com/sonic-family/sonic-express/pkg/manifest"
)

func main() {
	source := flag.String("source", ".", "source directory to inspect")
	output := flag.String("output", "manifest.dryrun.json", "output manifest path")
	version := flag.String("version", "0.0.0-dev", "bundle version")
	channel := flag.String("channel", "edge", "release channel")
	archCSV := flag.String("arch", "amd64", "comma-separated architectures")
	dryRun := flag.Bool("dry-run", true, "generate manifest only")
	flag.Parse()

	input := manifest.BuildInput{
		SourcePath:   *source,
		Version:      *version,
		Channel:      *channel,
		Architecture: splitCSV(*archCSV),
	}

	m, err := manifest.BuildDryRun(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build manifest: %v\n", err)
		os.Exit(1)
	}

	if !*dryRun {
		if err := bundle.BuildDraftBundle(*output, m); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build draft bundle: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("draft bundle written: %s\n", *output)
		return
	}

	if err := manifest.Write(*output, m); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write dry-run manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("dry-run manifest written: %s\n", *output)
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
