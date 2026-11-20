package main

import (
	"dot/binary"
	"dot/config"
	"dot/supervisor"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dot <serve|stop|build> [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		cfg, err := config.LoadConfig(config.Path)
		if err != nil {
			log.Fatalf("reading %s: %v", config.Path, err)
		}

		if err := binary.StartSupervisor(cfg); err != nil {
			log.Fatalf("supervisor %v", err)
		}

	case "stop":
		cfg, err := config.LoadConfig(config.Path)
		if err != nil {
			log.Fatalf("reading %s: %v", config.Path, err)
		}
		if err := supervisor.Stop(cfg, 10*time.Second); err != nil {
			log.Fatalf("stop: %v", err)
		}
		fmt.Println("Stopped supervisor and server")

	case "build":
		if err := binary.BuildBinaries(); err != nil {
			log.Fatalf("build: %v", err)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}
