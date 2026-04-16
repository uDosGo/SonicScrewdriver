package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	bundles := flag.String("bundles", "./bundles", "bundle directory")
	channels := flag.String("channels", "stable,beta,edge", "comma-separated channels")
	flag.Parse()

	fmt.Printf("serve stub: port=%d bundles=%s channels=%s\n", *port, *bundles, *channels)
}
