package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marabunta/ant"
)

var version string

func main() {
	parser := &ant.Parse{}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = parser.Usage(fs)

	cfg, err := parser.ParseArgs(fs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if (fs.Lookup("v")).Value.(flag.Getter).Get().(bool) {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	a, err := ant.New(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.Fatal(a.Start())
}
