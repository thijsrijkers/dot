# dot
A minimal, self-hostable Git server

[Documentation on roadmap steps](./docs/project_planning.md)

## Testing locally
Run all commands from the repository root:

1. Create `config.dot` in the repository root:
   ```json
   {
     "port": "8282",
     "hostKeyPath": "keys/host_key",
     "authKeysPath": "keys/authorized_keys"
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
