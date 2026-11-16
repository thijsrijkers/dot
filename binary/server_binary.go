package binary

import (
	"dot/config"
	"os"
	"os/exec"
)

func LaunchServer(config config.Config) error {
	cmd := exec.Command(serverBinaryPath,
		"-port", config.Port,
		"-host-key", config.HostKeyPath,
		"-authorized-keys", config.AuthKeysPath,
		"-root", config.Root,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
