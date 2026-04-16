package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	usb := flag.Bool("usb", false, "auto-detect install bundle from USB")
	channel := flag.String("channel", "", "install channel")
	from := flag.String("from", "", "update server URL")
	dryRun := flag.Bool("dry-run", false, "print install plan only")
	flag.Parse()

	target := ""
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	if target == "" && !*usb && *from == "" {
		fmt.Fprintln(os.Stderr, "usage: install [bundle.she] [--usb] [--channel <name> --from <url>]")
		os.Exit(2)
	}

	mode := "bundle"
	switch {
	case *usb:
		mode = "usb"
	case *from != "":
		mode = "channel"
	}

	fmt.Printf("install stub: mode=%s target=%s channel=%s from=%s dry_run=%t\n", mode, target, *channel, *from, *dryRun)
}
