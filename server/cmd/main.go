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
	flag.Parse()

	server.LaunchServer(*port, *hostKey, *authKeys)
}
