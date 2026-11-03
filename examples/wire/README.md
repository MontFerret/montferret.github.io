# Wire host and client

This module is the runnable example for the website's
[Remote runtimes (Wire) guide](https://ferretlang.org/docs/embedding/wire/).
The guide renders the Go files directly, so the displayed programs are the
programs tested here.

Requires Go 1.26 or later. Direct dependencies are pinned to published Ferret
Core `v2.0.0-alpha.62`, Universal API `v1.0.0-alpha.21`, Wire
`v1.0.0-alpha.4`, and gRPC `v1.84.0`. The matching website table is rendered
from `data/versions.yaml`; update it alongside this module when changing pins.
This module is independent of the website's root module and resolves published
dependencies without a local replacement. Existing ignored `.wire-dev` checkouts
are unused and may be kept.

Run returns a one-shot `api.Output`. The client collects it before closing the
output, session, plan, and runtime. Collected `*api.Content` remains usable after
cleanup; consumption may return available content alongside an error. See the
[Output reference](https://ferretlang.org/docs/embedding/go/output/) and
[migration guide](https://ferretlang.org/docs/migrations/consumable-output/).
Both Wire peers must negotiate output framing version 1; incompatible peers fail
during Connect without buffered-output fallback.

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
The client uses `client.New` with explicit `WithInsecure` and owns its channel.
The host uses `server.New` and caller-listener `Serve` so it can print its bound
address. `managed/loopback.go` shows blocking `Run` with a known loopback address.
Stop the host with Ctrl+C; it settles Wire before closing its engine. A five-second
shutdown timeout is retained while the host awaits actual settlement, then reported;
hosted cleanup can outlast that initial wait.

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

The [Wire security guide](https://ferretlang.org/docs/embedding/wire-security/)
renders the compilable `security` package's TLS, mTLS, token middleware, query,
and validation examples. The `streaming` package shows `Consume` forwarding borrowed
chunks to a writer, checking short writes and retaining cleanup errors.
Certificates and token validation are host inputs; only tests generate certificates. Security tests use real ephemeral loopback
listeners to verify trust, identities, plaintext rejection, client certificates,
startup and operation authentication, resource cleanup, startup-only cancellation,
owned versus borrowed channels, detached content after cleanup, terminal
consumption failures, and multi-message output. Authentication checks cover
`ReadExecutionOutput` as well as allocation and cleanup. TLS configurations
never skip verification.

Additional checks:

```sh
GOWORK=off GOFLAGS=-mod=readonly go test -race ./...
GOWORK=off GOFLAGS=-mod=readonly go vet ./...
go mod tidy -diff
```
