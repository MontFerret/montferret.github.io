---
title: "Remote runtimes (Wire)"
weight: 30
draft: false
description: "Use Wire to compile, execute, and debug scripts against an application's configured remote runtime."
---

# Remote runtimes (Wire)

Your application embeds Ferret and registers an `app::greet` function. A separate client needs to run a script that calls it, using the application's existing runtime setup:

{{< code lang="fql" name="greeting" >}}
return app::greet(@name)
{{< /code >}}

Wire lets the client compile and execute this script against the hosted runtime. With `name` set to `"Ada"`, it returns the JSON string `"Hello, Ada!"`. The custom function is registered on the host; the client sends source and parameters.

Use Wire when another process needs your application's functions, modules, services, or runtime configuration without rebuilding that setup locally.

## How the pieces fit together

**Ferret Core** provides the native engine and its embedding configuration. **The Universal API** provides shared runtime, plan, session, and debugger contracts. **Wire** exposes a hosted implementation remotely; its Go client returns an `api.Runtime`.

The native adapter and Wire client let consumers use the same contracts. Runtime construction and configuration remain the host's responsibility. See the [Universal API guide]({{< ref "docs/embedding/go/universal-api" >}}) for the shared programming model.

The current implementation uses protobuf and gRPC, with a Go server and client. This guide uses those implementations.

## Run a host and client

The example contains two complete programs in a [separate Go module](https://github.com/MontFerret/montferret.github.io/tree/main/examples/wire). It requires Go 1.25 or later and pins these compatible versions:

| Dependency | Version |
| --- | --- |
| Ferret Core | `v2.0.0-alpha.57` |
| Universal API | `v1.0.0-alpha.20` |
| Wire | `v1.0.0-alpha.2` |
| gRPC | `v1.84.0` |

Download the example and its dependencies:

{{< terminal >}}
git clone https://github.com/MontFerret/montferret.github.io.git
cd montferret.github.io/examples/wire
go mod download
{{< /terminal >}}

The module is independent of the website build. You can also copy the two programs below into `host/main.go` and `client/main.go` in a new directory, then initialize the same dependencies:

{{< terminal >}}
go mod init example.com/wire-demo
go get github.com/MontFerret/ferret/v2@v2.0.0-alpha.57 github.com/MontFerret/api@v1.0.0-alpha.20 github.com/MontFerret/wire@v1.0.0-alpha.2 google.golang.org/grpc@v1.84.0
{{< /terminal >}}

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

Save this program as `client/main.go`. It constructs its transport connection and follows the normal Universal API flow: `Compile → NewSession → Run`.

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

Stop the host with Ctrl+C. It uses a fresh five-second shutdown context to settle Wire-owned work before closing the native engine. For a different local port, start `go run ./host -listen 127.0.0.1:50052` and use `go run ./client -address 127.0.0.1:50052`. The `-listen`, `-address`, `-name`, and `-timeout` flags belong to these example programs.

## Connect with the CLI

A Wire-capable Ferret CLI can execute source, open a REPL, and report the host's version against this example:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:50051 --eval 'RETURN app::greet(@name)' --param name=Ada
ferret repl --runtime wire --runtime-endpoint tcp://127.0.0.1:50051
ferret version --runtime wire --runtime-endpoint tcp://127.0.0.1:50051
{{< /terminal >}}

The CLI uses the host's functions, modules, policies, and filesystem configuration. It supports source execution through Wire; artifacts and `ferret debug` require the builtin runtime. The [CLI Wire guide]({{< ref "docs/tools/cli/wire" >}}) covers its loopback-only endpoint syntax, timeout, compatibility, diagnostics, and ownership.

## What crosses the remote boundary

Registered functions and modules, host services, and runtime configuration stay on the host. Wire does not construct that runtime or transfer its configuration to the client.

Portable interfaces do not make every Go value transportable. Wire parameters support null, booleans, signed integers, unsigned integers that fit in `int64`, finite floats, strings, bytes, `[]any`, and `map[string]any`. Arbitrary Go functions, native host objects, datetime, duration, and regexp values cannot be sent as parameters. See the [client's parameter rules](https://github.com/MontFerret/wire/blob/v1.0.0-alpha.2/docs/client.md#options-and-parameters) for conversion details and limits.

Execution returns an `api.Output` with a content type and encoded bytes. The host selects an available encoder, and Wire preserves its output without decoding it. A client may decode those bytes according to the content type. Returning encoded data does not transfer a live host object or its capabilities.

## Wire or Worker?

Use [Worker]({{< ref "docs/tools/worker/overview" >}}) when you want to deploy the HTTP execution service with its supported runtime configuration and service controls.

Use Wire when your application already owns a configured runtime and needs to expose compilation, execution, and debugging through the Universal API. Wire is an integration library and protocol around that runtime. Your application owns the endpoint and deployment concerns.

## Debug a hosted script

Wire supports the Universal API debugging path: compile with `remote.CompileDebug`, then create a debugger with `plan.NewDebugSession`. The returned debugger session supports entry stops, breakpoints, stepping, inspection, and evaluation through the hosted implementation. Close the debugger and plan when finished.

See the [Wire client debugger documentation](https://github.com/MontFerret/wire/blob/v1.0.0-alpha.2/docs/client.md#debugger) and [Universal API guide]({{< ref "docs/embedding/go/universal-api" >}}) for the complete API. Wire's debugging capability does not imply that a particular editor or CLI can connect to its endpoint; that requires a tool-specific integration.

## Own resources and secure the endpoint

The host owns its runtime and endpoint. Wire borrows the runtime and does not close it. The client owns its transport connection; closing the remote runtime leaves that connection open.

Close every plan, normal session, and debug session you create. Settle work first, then close sessions before plans, the remote runtime, and the transport. Ordinary parent closure is not a substitute for closing caller-owned descendants. The example preserves cleanup errors with `errors.Join`.

Deployment security is an integration responsibility. Authenticate callers and authorize the compilation, execution, and debugging capabilities they can use. Protect transport traffic with an appropriate authenticated deployment boundary, or use a private socket with restrictive directory and socket permissions. Restrict the host functions, modules, filesystem, and network services exposed to scripts according to your application's needs. A loopback address alone is not authentication or isolation between local users.

Wire is **pre-stable**. Pin a compatible set of Core, Universal API, and Wire versions, and verify changes before upgrading. The `ferret.wire.v1` protobuf package name does not guarantee a stable public release.

## Next steps

- [Wire repository](https://github.com/MontFerret/wire)
- [Client documentation for the example's Wire version](https://github.com/MontFerret/wire/blob/v1.0.0-alpha.2/docs/client.md)
- [Protocol reference for the example's Wire version](https://github.com/MontFerret/wire/blob/v1.0.0-alpha.2/docs/protocol.md)

{{< docs-related tiles="embedding-go,embedding-go-custom-functions,tools-cli-wire,tools-worker" >}}
