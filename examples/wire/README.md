# Wire host and client

This module is the runnable example for the website's
[Remote runtimes (Wire) guide](https://ferretlang.org/docs/embedding/wire/).
The guide renders the Go files directly, so the displayed programs are the
programs tested here.

Requires Go 1.25 or later. Direct dependencies are pinned to Ferret Core
`v2.0.0-alpha.57`, Universal API `v1.0.0-alpha.20`, Wire `v1.0.0-alpha.2`,
and gRPC `v1.84.0`. This module is independent of the website's root module.

From this directory:

```sh
go mod download
go run ./host
```

In another terminal, from this directory:

```sh
go run ./client
```

Expected output:

```text
Content-Type: application/json
"Hello, Ada!"
```

The host alone registers `app::greet`; the client sends a named FQL source and
the `name` parameter. Try `go run ./client -name Grace` to change the parameter.
Stop the host with Ctrl+C; it settles Wire before closing its engine.

The plaintext loopback endpoint is for local development only. Loopback is not
authentication. Deployment security and runtime capabilities belong to the host.

The host's `-listen` flag defaults to `127.0.0.1:50051`; use `127.0.0.1:0`
for an automatically assigned local port. The client accepts `-address`, `-name`,
and `-timeout` (default `10s`). These are example-program flags.

Run the integration check:

```sh
GOWORK=off GOFLAGS=-mod=readonly go test ./...
```

It builds and runs the exact host and client as separate processes over a real
loopback TCP connection, checks two parameter values and an expired context,
verifies signal-driven host shutdown, and checks an unavailable endpoint.
