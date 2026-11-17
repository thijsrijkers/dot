package server

import (
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func HealthCheck(socket string) error {
	if err := os.MkdirAll(filepath.Dir(socket), 0o750); err != nil {
  	return err
  }
  if err := os.Remove(socket); err != nil && !errors.Is(err, os.ErrNotExist) {
    return err
  }

  listener, err := net.Listen("unix", socket)
  if err != nil {
    return err
  }

  if err := os.Chmod(socket, 0o660); err != nil {
    listener.Close()
    return err
  }
        
  return http.Serve(listener, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte("ok\n"))
  }))
}
