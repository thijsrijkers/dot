package server

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type KeyStore struct {
	mu   sync.RWMutex
	path string
	keys map[string]map[string]ssh.PublicKey 
}

func NewKeyStore(path string) (*KeyStore, error) {
	store := &KeyStore{path: path, keys: make(map[string]map[string]ssh.PublicKey)}
	if err := store.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return store, nil
}

func (store *KeyStore) Has(user string, key ssh.PublicKey) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	_, ok := store.keys[user][ssh.FingerprintSHA256(key)]
	return ok
}

func (store *KeyStore) AddAuthorizedKey(user string, line []byte) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	key, _, _, _, err := ssh.ParseAuthorizedKey(line)
	if err != nil {
		return err
	}

	if user == "" || strings.ContainsAny(user, " \t\r\n") {
		return fmt.Errorf("invalid user name %q", user)
	}

	fingerPrint := ssh.FingerprintSHA256(key)

	if _, exists := store.keys[user][fingerPrint]; exists {
		return nil
	}

	// Write disk first, so memory never contains a non existing key
	if err := store.appendToFile(user, key); err != nil {
		return fmt.Errorf("persist key: %w", err)
	}

	if store.keys[user] == nil {
		store.keys[user] = make(map[string]ssh.PublicKey)
	}

	store.keys[user][fingerPrint] = key

	return nil
}

func (store *KeyStore) PublicKeyCallback(connection ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	if store.Has(connection.User(), key) {
		return &ssh.Permissions{
			Extensions: map[string]string{"pubkey-fp": ssh.FingerprintSHA256(key)},
		}, nil
	}
	return nil, fmt.Errorf("unknown public key for %q", connection.User())
}

func (store *KeyStore) Load() error {
	data, err := os.ReadFile(store.path)

	if err != nil {
		return err
	}

	keys := make(map[string]map[string]ssh.PublicKey)
	count := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
	
		if line == "" || line[0] == '#' {
			continue
		}

		user, rest, ok := strings.Cut(line, " ")
	
		if !ok {
			log.Printf("%s:%d: skipping line without key", store.path, lineNumber)
			continue
		}

		key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(rest)))
	
		if err != nil {
			log.Printf("%s:%d: skipping invalid key: %v", store.path, lineNumber, err)
			continue
		}
	
		if keys[user] == nil {
			keys[user] = make(map[string]ssh.PublicKey)
		}

		keys[user][ssh.FingerprintSHA256(key)] = key
		count++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	store.mu.Lock()
	store.keys = keys
	store.mu.Unlock()

	log.Printf("loaded %d keys from %s", count, store.path)

	return nil
}

func (store *KeyStore) appendToFile(user string, key ssh.PublicKey) error {
	if err := os.MkdirAll(filepath.Dir(store.path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(store.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(formatLine(user, key)); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func formatLine(user string, key ssh.PublicKey) string {
	return user + " " + string(ssh.MarshalAuthorizedKey(key))
}
