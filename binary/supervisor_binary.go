package binary

import (
	"dot/config"
	"os"
	"os/exec"
)

func StartSupervisor(config config.Config) error {
	cmd := exec.Command(supervisorBinaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

