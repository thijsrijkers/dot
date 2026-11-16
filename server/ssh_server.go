package server

import (
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)


func LaunchServer(port string, hostKeyPath string, authorizationKeyStorePath string) {
	store, err := NewKeyStore(authorizationKeyStorePath)

	if err != nil {
		log.Printf("Failed to create key store: %v", err)
		return
	}

	store.Load()

	config := &ssh.ServerConfig{
		PublicKeyCallback: store.PublicKeyCallback, 
	}

	config.AddHostKey(RetrieveHostKey(hostKeyPath))

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Println("Dot server listening on port: "+port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		// TODO: handle connection. Closed for now so sockets don't leak.
		connection.Close()
	}

}

