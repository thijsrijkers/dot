package main

import (
	"context"
	"dot/config"
	"dot/supervisor"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig(config.Path)
	if err != nil {
		log.Fatalf("reading %s: %v", config.Path, err)
	}

	if err := supervisor.WritePid(cfg); err != nil {
		log.Fatalf("writing pid file: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	err = supervisor.Run(ctx, cfg, 5*time.Second)
	stop()
	supervisor.RemovePid(cfg)

	if err != nil {
		log.Fatalf("supervisor: %v", err)
	}
}
