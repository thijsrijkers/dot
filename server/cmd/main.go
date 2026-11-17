// Command dot-server runs only the SSH git server.
// Build it with `dot build`; the dot binary launches it for `dot serve`.
package main

import (
	"flag"

	"dot/server"
)

func main() {
	port := flag.String("port", "2222", "port to listen on")
	hostKey := flag.String("host-key", "host_key", "path to host key")
	authKeys := flag.String("authorized-keys", "authorized_keys", "path to authorized_keys")
	root := flag.String("root", "repos", "directory holding bare repositories")
  healthSocket := flag.String("health-socket", "dot-health.sock", "unix socket for health checks")
	flag.Parse()

	server.LaunchServer(*port, *hostKey, *authKeys, *root, *healthSocket)
}
