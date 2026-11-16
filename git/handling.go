package git

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func HandleConnection(connection net.Conn, config *ssh.ServerConfig, root string) {
	defer connection.Close()

	sshConn, channels, requests, err := ssh.NewServerConn(connection, config)
	if err != nil {
		log.Printf("handshake from %s failed: %v", connection.RemoteAddr(), err)
		return
	}
	defer sshConn.Close()
	go ssh.DiscardRequests(requests)

	for newChannel := range channels {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "only sessions are allowed")
			continue
		}
		channel, channelRequests, err := newChannel.Accept()
		if err != nil {
			log.Printf("accept channel: %v", err)
			continue
		}
		go handleSession(sshConn.User(), channel, channelRequests, root)
	}
}

func handleSession(user string, channel ssh.Channel, requests <-chan *ssh.Request, root string) {
	defer channel.Close()
	var env []string

	for req := range requests {
		switch req.Type {
		//this is setting the env. Git uses it for one thing: telling the server which version of the git wire protocol it wants to speak.
		case "env":
			var keyValue struct{ Name, Value string }
			if ssh.Unmarshal(req.Payload, &keyValue) == nil && keyValue.Name == "GIT_PROTOCOL" {
				env = append(env, "GIT_PROTOCOL="+keyValue.Value)
				req.Reply(true, nil)
			} else {
				req.Reply(false, nil)
			}

		//Command execution
		case "exec":
			var payload struct{ Command string }
			if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
				req.Reply(false, nil)
				return
			}

			service, repoPath, err := ParseGitCommand(payload.Command, root)
			if err != nil {
				log.Printf("%s: rejected %q: %v", user, payload.Command, err)
				req.Reply(true, nil)
				fmt.Fprintf(channel.Stderr(), "dot: %v\n", err)
				sendExitStatus(channel, 1)
				return
			}

			req.Reply(true, nil)
			log.Printf("%s: %s %s", user, service, repoPath)
			sendExitStatus(channel, RunGit(channel, service, repoPath, env))

			return

		default:
			req.Reply(false, nil)
		}
	}
}

func sendExitStatus(channel ssh.Channel, status uint32) {
	channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{status}))
}
