# dot
A minimal, self-hostable Git server

[Documentation on roadmap steps](./docs/project_planning.md)

## Testing locally
Run all commands from the repository root:

1. Create `config.dot` in the repository root:
   ```json
   {
     "port": "8282",
     "hostKeyPath": "path/to/host_key_storage",
     "authKeysPath": "path/to/authorized_keys_storage",
     "root": "path/to/repository_storage",
     "healthSocket": "path/to/health.sock"
   }
   ```
   A host key is generated at `hostKeyPath` on first start if none exists.

2. Commands to use with Dot:
   ```sh
   go run . build
   go run . serve
   ```

3. From another terminal, check that the port is open:
   ```sh
   nc -z localhost 8282 && echo up
   ```

4. Create a bare repository and clone it (the SSH user must have a key in `authorized_keys`):
   ```sh
   git init --bare repos/demo.git
   git clone ssh://<user>@localhost:8282/demo.git
   ```
   Only `git-upload-pack` (clone/fetch) and `git-receive-pack` (push) are accepted, and only for repos inside `root`.
