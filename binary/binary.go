package binary

import (
	"fmt"
	"os"
	"os/exec"
)

const serverBinaryPath = "binary/dot-server"

// binaries maps each output path to the main package it is built from.
var binaries = map[string]string{
	serverBinaryPath: "./server/cmd",
}

func BuildBinaries() error {
	for output, pkg := range binaries {
		cmd := exec.Command("go", "build", "-trimpath", "-ldflags=-s -w", "-o", output, pkg)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", output, err)
		}
		fmt.Println("Built", output)
	}
	return nil
}
