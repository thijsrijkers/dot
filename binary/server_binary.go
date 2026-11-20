package binary

import (
	"dot/config"
	"os"
	"os/exec"
)

func StartServer(config config.Config) (*exec.Cmd, error) {
	cmd := exec.Command(serverBinaryPath,
		"-port", config.Port,
		"-host-key", config.HostKeyPath,
		"-authorized-keys", config.AuthKeysPath,
		"-root", config.Root,
    "-health-socket", config.HealthSocket,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
  
	if err := cmd.Start(); err != nil {
  	return nil, err
  }
  
	return cmd, nil
}

func LaunchServer(config config.Config) error {
	cmd, err := StartServer(config)
  
	if err != nil {
  	return err
  }
  
	return cmd.Wait()
}

