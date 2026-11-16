package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func RetrieveHostKey(path string) ssh.Signer {
	data, err := os.ReadFile(path)

	switch {
		case err == nil:
			signer, err := ssh.ParsePrivateKey(data)

			if err != nil {
				log.Fatalf("parse host key %s: %v", path, err)
			}

			return signer
		case !errors.Is(err, os.ErrNotExist):
			log.Fatalf("read host key %s: %v", path, err)
	}

	log.Printf("no host key at %s, generating a new one", path)

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("generate host key: %v", err)
	}

	block, err := ssh.MarshalPrivateKey(priv, "") 
	if err != nil {
		log.Fatalf("marshal host key: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		log.Fatalf("create key dir: %v", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		log.Fatalf("create host key %s: %v", path, err)
	}

	if err := pem.Encode(file, block); err != nil {
		file.Close()
		os.Remove(path)
		log.Fatalf("write host key: %v", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(path)
		log.Fatalf("write host key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Fatalf("host key signer: %v", err)
	}

	return signer
}
