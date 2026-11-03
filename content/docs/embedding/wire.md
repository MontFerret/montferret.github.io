---
title: "Remote runtimes (Wire)"
weight: 30
draft: false
description: "Use Wire to compile, execute, and debug scripts against an application's configured remote runtime."
---

# Remote runtimes (Wire)

Your application embeds Ferret and registers an `app::greet` function. A separate client needs to run a script that calls it, using the application's existing runtime setup:

{{< code lang="fql" name="greeting" >}}
RETURN app::greet(@name)
{{< /code >}}

Wire lets the client compile and execute this script against the hosted runtime. With `name` set to `"Ada"`, it returns the JSON string `"Hello, Ada!"`. The custom function is registered on the host; the client sends source and parameters.

Use Wire when another process needs your application's functions, modules, services, or runtime configuration without rebuilding that setup locally.

## How the pieces fit together

**Ferret Core** provides the native engine and its embedding configuration. **The Universal API** provides shared runtime, plan, session, and debugger contracts. **Wire** exposes a hosted implementation remotely; its Go client returns an `api.Runtime`.

This guide and its runnable example use published releases with Universal API `v{{< data "versions.wire.api" >}}` consumable output. A successful `Run` returns a one-shot handle. Collect or consume it to observe completion; close it if you abandon it. See [Output]({{< ref "docs/embedding/go/output" >}}) for the shared contract.

Runtime construction and configuration remain the host's responsibility. The [Universal API guide]({{< ref "docs/embedding/go/universal-api" >}}) explains the adapter; the example below pins a compatible published dependency set.

The current implementation uses protobuf and gRPC, with a Go server and client. This guide uses those implementations.

## Run a host and client

The example contains two complete programs in a [separate Go module](https://github.com/MontFerret/montferret.github.io/tree/main/examples/wire). It requires Go {{< data "versions.wire.go" >}} or later and pins these compatible versions:

| Dependency | Version |
| --- | --- |
| Ferret Core | `v{{< data "versions.wire.core" >}}` |
| Universal API | `v{{< data "versions.wire.api" >}}` |
| Wire | [`v{{< data "versions.wire.v1" >}}`](https://github.com/MontFerret/wire/releases/tag/v{{< data "versions.wire.v1" >}}) |
| gRPC | `v{{< data "versions.wire.grpc" >}}` |

Download the example and its dependencies:

{{< terminal >}}
git clone https://github.com/MontFerret/montferret.github.io.git
cd montferret.github.io/examples/wire
go mod download
{{< /terminal >}}

The module resolves published dependencies directly, with no local checkout replacement. Keep these pins together when upgrading. Both Wire peers must negotiate output framing version **1** during `Connect`; incompatible peers fail at startup, with no buffered-output fallback. Upgrade older hosts and clients together using the [consumable output migration]({{< ref "docs/migrations/consumable-output" >}}).

### Host the configured engine

Save this program as `host/main.go`. The host registers a typed, one-argument function, constructs the native engine, wraps it through the official adapter, and serves it through Wire.

{{< code lang="go" name="host/main" file="examples/wire/host/main.go" >}}{{< /code >}}

`uapi.Wrap` borrows the engine. The supplied version identifies Ferret Core, independently of the host application and Wire. The host remains responsible for shutting down Wire and closing the engine.

Start the host in the example directory:

{{< terminal >}}
go run ./host
{{< /terminal >}}

It prints:

```text
Listening on 127.0.0.1:50051
```

This example uses **local-development-only plaintext transport** on loopback. Loopback binding limits network exposure; it does not authenticate clients. Configure deployment security before exposing an endpoint beyond this example.

### Execute from another process

Save this program as `client/main.go`. It uses `client.New(ctx, target, client.WithInsecure())` to create and own its channel, then follows the normal Universal API flow: `Compile → NewSession → Run → Collect`. The output is settled before the session, plan, and remote runtime are closed.

{{< code lang="go" name="client/main" file="examples/wire/client/main.go" >}}{{< /code >}}

In another terminal, from the same example directory, run:

{{< terminal >}}
go run ./client
{{< /terminal >}}

The output is:

```text
Content-Type: application/json
"Hello, Ada!"
```

The client has no native engine or function registration. `app::greet` runs on the host with the supplied parameter. To change it, run `go run ./client -name Grace`.

Both programs handle Ctrl+C and termination signals. The client bounds its operations with a ten-second context; use `-timeout` to change that bound. Cancellation ends work but does not replace explicit resource cleanup. Execution and cleanup errors are reported with a nonzero exit status, and available output is inspected independently of an error.

Stop the host with Ctrl+C. It first waits with a fresh five-second shutdown context. If that wait expires, it retains the timeout error and waits for cleanup to settle before closing the native engine and reporting the result. Uncooperative hosted cleanup can therefore keep the process alive beyond five seconds. For a different local port, start `go run ./host -listen 127.0.0.1:50052` and use `go run ./client -address 127.0.0.1:50052`. The `-listen`, `-address`, `-name`, and `-timeout` flags belong to these example programs.

## Choose startup and shutdown ownership

`server.New(hosted, options...)` configures Wire without listening or taking ownership of the supplied `api.Runtime`. A server accepts one serving lifetime. Both startup paths use the same transport credentials, interceptors, limits, and shutdown lifecycle.

For a known TCP address, `Server.Run` blocks through serving and managed shutdown. This compiled example uses the host's existing runtime:

{{< code lang="go" name="managed/loopback" file="examples/wire/managed/loopback.go" >}}{{< /code >}}

Canceling or reaching the deadline of the serving context triggers shutdown. `Run` computes a separate cleanup deadline: **30 seconds by default**, or a positive `server.WithShutdownTimeout(duration)` override. The budget starts when shutdown begins. A deadline supplied to `Shutdown` can shorten the effective deadline; later calls cannot extend it. `Shutdown(context.Background())` cannot disable an accepted managed budget.

At the effective deadline, Wire forces transport shutdown independently of blocked hosted cleanup. `Run` returns an error matching `context.DeadlineExceeded`, joined with known independent failures. This does **not** mean hosted cleanup has finished. A subsequent `Shutdown` waits on the retained lifecycle and returns its actual settlement result without repeating cleanup. The example waits for that settlement before its caller may close the borrowed runtime. Successful cancellation or explicit shutdown returns `nil` only after serving and cleanup both settle successfully.

Use `Serve(ctx, listener)` for a caller-created listener, including Unix sockets or custom transports. The greeting host uses it to print the actual address when `-listen 127.0.0.1:0` chooses an ephemeral port. Once accepted, [gRPC closes the listener when serving returns](https://pkg.go.dev/google.golang.org/grpc@v1.84.0#Server.Serve); rejected startup leaves it untouched. `Serve` has no managed default timeout. Its caller coordinates shutdown and awaits `Shutdown` before releasing the runtime.

## Connect with the CLI

Ferret CLI `v{{< data "versions.cli.v2" >}}` can execute source, open a REPL, and report the host's version against this example. The example's published dependency set uses the consumable-output framing required by current CLI `main`; older buffered-output hosts are incompatible. See the [CLI Wire guide]({{< ref "docs/tools/cli/wire" >}}) for the current compatibility requirements and separate CLI development host.

These commands use this example's `app::greet` function, not the CLI development host's `DEMO::IDENTITY`:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:50051 --eval 'RETURN app::greet(@name)' --param name=Ada
ferret repl --runtime wire --runtime-endpoint tcp://127.0.0.1:50051
ferret version --runtime wire --runtime-endpoint tcp://127.0.0.1:50051
{{< /terminal >}}

The CLI uses the host's functions, modules, policies, and filesystem configuration. It supports source execution through Wire; artifacts and `ferret debug` require the builtin runtime. The [CLI Wire guide]({{< ref "docs/tools/cli/wire" >}}) covers its loopback-only endpoint syntax, timeout, compatibility, diagnostics, and ownership.

## What crosses the remote boundary

Registered functions and modules, host services, and runtime configuration stay on the host. Wire does not construct that runtime or transfer its configuration to the client.

Portable interfaces do not make every Go value transportable. Wire parameters support null, booleans, signed integers, unsigned integers that fit in `int64`, finite floats, strings, bytes, `[]any`, and `map[string]any`. Arbitrary Go functions, native host objects, datetime, duration, and regexp values cannot be sent as parameters. See the [client's parameter rules](https://github.com/MontFerret/wire/blob/v{{< data "versions.wire.v1" >}}/docs/client.md#options-and-parameters) for conversion details and limits.

The host selects an available encoder, and Wire preserves its bytes without decoding them. `output.Metadata()` describes the complete representation and remains readable after closure. `output.Collect(ctx)` returns detached `*api.Content`: inspect `content.Metadata.ContentType` and `content.Data`. Nil content means absent output; a non-nil content value may contain zero bytes. Returning encoded data does not transfer a live host object or its capabilities.

## Consume output and observe completion

A `Run` error is an admission or preparation failure and returns no usable handle. After successful preparation, execution, encoding, delivery, and output-owned cleanup failures are reported by `Collect` or `Consume`. Collection may return received content alongside an error; inspect both. Keep the invocation context alive until output is settled. A consumption context may shorten that lifetime, but cannot extend it.

Use `Collect` for small results you need to retain. To forward bytes without materializing the entire response, use this compiled example:

{{< code lang="go" name="streaming/query" file="examples/wire/streaming/query.go" >}}{{< /code >}}

`Consume` supplies borrowed chunks synchronously with backpressure. Chunk boundaries may split encoded values or UTF-8 characters; copy bytes you retain after the callback. The writer must return promptly and honor any cancellation policy it needs. Callbacks must not call `Close` on their own output because closure waits for the callback. For side-effect-only queries, consume with a callback that discards bytes so terminal failures are still observed.

After consumption starts, both terminal operations finalize output-owned resources, including on failure. `Close` abandons unread output without draining it and does not certify successful execution. Close each output before reusing its durable session, and retain cleanup errors.

Wire's default inbound and outbound limits are **4 MiB per gRPC message**, including framing. Large results are split across messages; Wire imposes **no aggregate output-size cap**. `Collect` materializes received bytes, and debugger output is materialized on both sides. Hosts and consumers must budget for total result sizes. Chunked delivery does not imply that the hosted engine evaluates or encodes queries incrementally.

## Wire or Worker?

Use [Worker]({{< ref "docs/tools/worker/overview" >}}) when you want to deploy the HTTP execution service with its supported runtime configuration and service controls.

Use Wire when your application already owns a configured runtime and needs to expose compilation, execution, and debugging through the Universal API. Wire is an integration library and protocol around that runtime. Your application owns the endpoint and deployment concerns.

## Debug a hosted script

Wire supports the Universal API debugging path: compile with `remote.CompileDebug`, then create a debugger with `plan.NewDebugSession`. The returned debugger session supports entry stops, breakpoints, stepping, inspection, and evaluation through the hosted implementation. Close the debugger and plan when finished.

See the [Wire client debugger documentation](https://github.com/MontFerret/wire/blob/v{{< data "versions.wire.v1" >}}/docs/client.md#debugger) and [Universal API guide]({{< ref "docs/embedding/go/universal-api" >}}) for the Wire and Native adapter contracts. Wire's debugging capability does not imply that a particular editor or CLI can connect to its endpoint; that requires a tool-specific integration.

## Own resources and secure the endpoint

The host owns its runtime and endpoint. Wire borrows the runtime and never closes it. Settle every output and close every plan, normal session, and debug session you create. Ordinary parent closure does not replace closing caller-owned descendants; the examples preserve cleanup errors with `errors.Join`.

`client.New(ctx, target, options...)` creates and owns a gRPC channel. Its default is TLS with system trust and normal server identity verification. `client.WithInsecure()` explicitly selects plaintext for this loopback demo. There is **no plaintext fallback** if TLS fails. Constructor cancellation controls startup only; after construction, use each operation's context to cancel work.

Closing the remote runtime gates new root operations. Existing plans and admitted work may retain the logical connection, so the owned channel closes only after its last descendant and in-flight work finish. Settle outputs, then close sessions, plans, and the remote runtime. With `client.From(ctx, conn)`, Wire always borrows the caller's transport, including on failed startup: the caller closes that channel after logical resources settle.

Configure TLS, mTLS, or host-supplied token middleware using the [Wire security examples]({{< ref "docs/embedding/wire-security" >}}). Server options and `Run` options collect validation failures before construction or listening; the security guide shows `errors.As` inspection and collection diagnostics. Loopback alone is not authentication or isolation between local users. Authorization and runtime capabilities remain the host's responsibility.

Wire is **pre-stable**. Pin a compatible set of Core, Universal API, and Wire versions, and verify changes before upgrading. The `ferret.wire.v1` protobuf package name does not guarantee a stable public release.

## Next steps

- [Wire repository](https://github.com/MontFerret/wire)
- [Client documentation for the example's Wire version](https://github.com/MontFerret/wire/blob/v{{< data "versions.wire.v1" >}}/docs/client.md)
- [Protocol reference for the example's Wire version](https://github.com/MontFerret/wire/blob/v{{< data "versions.wire.v1" >}}/docs/protocol.md)

{{< docs-related tiles="embedding-wire-security,embedding-go,embedding-go-custom-functions,tools-cli-wire,tools-worker" >}}
