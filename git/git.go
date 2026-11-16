package git

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

var services = map[string]string{
        "git-upload-pack":  "upload-pack",  // clone / fetch / pull
        "git-receive-pack": "receive-pack", // push
}

func ParseGitCommand(command, root string) (service, repoPath string, err error) {
	verb, arg, ok := strings.Cut(strings.TrimSpace(command), " ")
	if !ok {
		return "", "", errors.New("expected: git-upload-pack|git-receive-pack '<repo>'")
	}

	service, ok = services[verb]
	if !ok {
		return "", "", fmt.Errorf("command %q is not allowed", verb)
	}

	name := strings.TrimSpace(arg)
	if len(name) >= 2 && name[0] == '\'' && name[len(name)-1] == '\'' {
		name = name[1 : len(name)-1]
	}

	if strings.ContainsAny(name, "'\\\x00") {
		return "", "", errors.New("invalid repository name")
	}

	name = strings.TrimPrefix(name, "~/")
	name = strings.TrimLeft(name, "/")

	if !filepath.IsLocal(name) {
		return "", "", errors.New("invalid repository path")
	}

	for _, candidate := range []string{name, name + ".git"} {
		full := filepath.Join(root, candidate)
		if info, err := os.Stat(full); err == nil && info.IsDir() {
			return service, full, nil
		}
	}

	return "", "", fmt.Errorf("repository %q not found", name)
}

func RunGit(channel ssh.Channel, service, repoPath string, env []string) uint32 {
	cmd := exec.Command("git", service, repoPath)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = channel
	cmd.Stderr = channel.Stderr()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("stdin pipe: %v", err)
		return 1
	}

	if err := cmd.Start(); err != nil {
		log.Printf("start git %s: %v", service, err)
		return 1
	}

	go func() {
		io.Copy(stdin, channel)
		stdin.Close()
	}()

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return uint32(exitErr.ExitCode())
		}
		log.Printf("git %s: %v", service, err)
		return 1
	}

	return 0
}

