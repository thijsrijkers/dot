package server

import (
	"dot/git"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)


func LaunchServer(port string, hostKeyPath string, authorizationKeyStorePath string, root string, healthSocket string) {
	if root == "" {
		log.Fatal("Root is not set")
	}

	store, err := NewKeyStore(authorizationKeyStorePath)

	if err != nil {
		log.Printf("Failed to create key store: %v", err)
		return
	}

	log.Println("Loading store values...")
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

	go func() {
    if err := HealthCheck(healthSocket); err != nil {
      log.Fatalf("health socket: %v", err)
  	}
  }()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		git.HandleConnection(connection, config, root)
	}
}

