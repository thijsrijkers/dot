1. Project setup [✓]

Write it in Go so it builds into one static binary.
Use one binary, dot, with subcommands: serve, supervise, dashboard, user, repo.
Split the code into modules: server, supervisor, status, dashboard, admin.

2. SSH git server (serve) [✓]

Allow only SSH, with public-key login and no passwords.
Check keys against a dot-managed authorized_keys file.
Allow only pull (git-upload-pack) and push (git-receive-pack).
Clean every repo path so users can’t reach anything outside the repo root.
Use a persistent host key so clients don’t get “host key changed” warnings.

3. Health API [✓]

Serve health over a local unix socket, so nothing new is exposed on the network.

4. Supervisor (supervise)

Start the server as a child process.
Check health every few seconds via the socket and a test connection to the SSH port.
Restart the server if it crashes or fails several checks in a row, with increasing delays between restarts.
Expose its own state (PID, restart count, last check) on a second socket.

5. Terminal dashboard (dashboard)

Build a full-screen terminal UI that refreshes every second.
Panels: server status and uptime, supervisor health and restarts, active users, recent pushes and pulls, errors.
Use keyboard navigation, with q to quit.

6. Admin commands

Add and remove users together with their SSH public keys.
Create, list and delete bare repositories.

7. Configuration

Use one config file for the port, repo directory, key paths and socket paths.
Add per-repo read/write permissions later.

8. Deployment

Build the binary and install it on the server.
Create a dedicated system user and directories for config, repos and runtime files.
Generate the host key.
Run supervise as a systemd service, so systemd restarts the supervisor and the supervisor restarts the server.
Set the firewall to allow only the SSH git port.

9. Usage

Admins add users and repos, then watch the server with the dashboard.
Git users clone and push over SSH as usual.

10. Testing and hardening

Write integration tests that run a real clone and push against the binary.
Test rejected keys, path tricks and blocked commands.
Add log rotation, rate limiting on failed logins, metrics and backups.
