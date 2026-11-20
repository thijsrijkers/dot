package supervisor

import (
	"context"
	"dot/binary"
	"dot/config"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"syscall"
	"time"
)

type process struct {
	cmd  *exec.Cmd
  done chan struct{} 
}

func newUnixClient(socket string) *http.Client {
	return &http.Client{
  	Timeout: 2 * time.Second,
    Transport: &http.Transport{
      DisableKeepAlives: true, 
      DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
        var dialer net.Dialer
        return dialer.DialContext(ctx, "unix", socket)
      },
  	},
	}
}

func Run(ctx context.Context, cfg config.Config, interval time.Duration) error {
	process, err := start(cfg)
  if err != nil {
    return err
  }

  defer func() { process.stop() }()

  client := newUnixClient(cfg.HealthSocket)
  ticker := time.NewTicker(interval)
  defer ticker.Stop()

  for {
    select {
      case <-ctx.Done():
        return nil
      case <-ticker.C:
        if err := check(client); err == nil {
          continue
        } else {
          log.Printf("supervisor: health check failed: %v; restarting server", err)
        }

        process.stop()
        if process, err = start(cfg); err != nil {
          return fmt.Errorf("restarting server: %w", err)
        }
    }
  }
}

func start(cfg config.Config) (*process, error) {
  cmd, err := binary.StartServer(cfg)

  if err != nil {
    return nil, err
  }

  process := &process{cmd: cmd, done: make(chan struct{})}
  go func() {
    err := cmd.Wait()
    log.Printf("supervisor: server (pid %d) exited: %v", cmd.Process.Pid, err)
    close(process.done)
  }()

  log.Printf("supervisor: server started (pid %d)", cmd.Process.Pid)

  return process, nil
}

func (process *process) stop() {
	select {
  	case <-process.done:
      return
    default:
  }
  
	_ = process.cmd.Process.Signal(syscall.SIGTERM)
  select {
  	case <-process.done:
    case <-time.After(5 * time.Second):
    	_ = process.cmd.Process.Kill()
      <-process.done
  }
}

func check(client *http.Client) error {
	response, err := client.Get("http://unix/health")
  
	if err != nil {
  	return err
  }

	defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
  	return fmt.Errorf("status %d", response.StatusCode)
  }
  
	return nil
}

