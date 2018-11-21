package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marabunta/client"
)

var version string

func main() {
	var (
		v        = flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))
		id       = flag.String("id", "", "ant ID")
		host     = flag.String("host", "", "Connect to `host` default (marabunta)")
		port     = flag.Int("port", 1415, "Port number to use for the connection")
		certFile = flag.String("cert", "server.crt", "TLS cert")
	)

	flag.Parse()
	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if _, err := os.Stat(*certFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Cannot read the server certificate file: %s, use -h for more info.\n", *certFile)
		os.Exit(1)
	}

	ant, err := client.New(*id, *host, *port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Run it
	if err := ant.Run(*certFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
