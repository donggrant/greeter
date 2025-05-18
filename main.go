package main

import (
	"flag"
	"os"
)

func main() {
	// Define command-line flags
	serverMode := flag.Bool("server", false, "Run in server mode")
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	// Set port environment variable if specified
	if *port != "8080" {
		os.Setenv("PORT", *port)
	}

	if *serverMode {
		// Run in server mode
		RunServer()
	} else {
		// Run in CLI mode
		RunCLI()
	}
}
