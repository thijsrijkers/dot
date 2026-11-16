package main

import (
	"dot/binary"
	"dot/config"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dot <serve|build> [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		cfg, err := config.LoadConfig(config.Path)
		if err != nil {
			log.Fatalf("reading %s: %v", config.Path, err)
		}
		if err := binary.LaunchServer(cfg); err != nil {
			log.Fatalf("server: %v", err)
		}

	case "build":
		if err := binary.BuildBinaries(); err != nil {
			log.Fatalf("build: %v", err)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}
