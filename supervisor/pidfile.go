package supervisor

import (
	"dot/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func PidFile(cfg config.Config) string {
	return filepath.Join(filepath.Dir(cfg.HealthSocket), "supervisor.pid")
}

func WritePid(cfg config.Config) error {
	path := PidFile(cfg)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strconv.Itoa(os.Getpid())+"\n"), 0o640)
}

func RemovePid(cfg config.Config) {
	_ = os.Remove(PidFile(cfg))
}

func Stop(cfg config.Config, timeout time.Duration) error {
	path := PidFile(cfg)
	data, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		return errors.New("supervisor is not running")
	}

	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return fmt.Errorf("invalid pid file %s: %w", path, err)
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			_ = os.Remove(path) 
			return errors.New("supervisor is not running (removed stale pid file)")
		}
		return err
	}

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("supervisor (pid %d) did not exit within %s", pid, timeout)
}
