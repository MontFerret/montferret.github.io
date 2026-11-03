---
title: "Introducing Ferret Wire: Your Runtime, Connected"
subtitle: "Compile, run, and debug FQL against the Ferret runtime embedded in your application."
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-10-02"
---

Hello, friends!

In an earlier post, I talked about why Ferret was built as a library first. Your application embeds the engine, decides which modules and custom functions are available, and defines the environment in which queries run. The CLI is one way to use Ferret, but it is not the center of the architecture.

That approach creates an interesting question:

*How do development tools work with the runtime inside your application, rather than a separate runtime that merely resembles it?*

Today, I want to introduce `Ferret Wire`: a protocol and Go SDK that lets another process compile, execute, and debug queries using an application's configured Ferret runtime. The current implementation uses `Protocol Buffers` and `gRPC`.

The idea is simple: bring the tools to your runtime, instead of rebuilding your runtime inside every tool.

## The missing piece is your application

Imagine an inventory service that embeds Ferret and registers an `inventory` namespace. A query might look like this:

{{< code lang="fql" >}}
RETURN inventory::search({
    category: @category,
    available: true
})
{{< /code >}}

`inventory::search` is not a built-in Ferret function. In this example, it is something your application provides, backed by its own services and business logic. This ability to expose application-specific functionality is part of Ferret's library-first design.

Now imagine opening that query in a separate development tool.

The tool might understand FQL perfectly. It might have the same version of Ferret. But that does not mean its runtime knows what `inventory::search` does or even that the function exists.

You could reproduce the application's setup inside the tool. Register the same functions, configure the same dependencies, and keep everything synchronized as the application evolves.

But that is another environment to maintain. And when you are investigating why a query behaves a certain way, another environment is not necessarily what you need.

Wire provides a different route: send the query to the application's exposed runtime. Compilation and execution happen there, using the runtime the host has already configured. The client does not need to recreate that configuration.

## One API, on either side of the connection

Wire builds on the Universal Ferret API, which defines shared contracts for runtimes, compiled plans, execution sessions, and source-level debugging.

These contracts describe how a consumer uses a runtime without prescribing how the host constructs it. A consumer can compile a query into a reusable plan, create a session with parameters, and execute it through the same interfaces.

On the server side, Wire accepts an `api.Runtime`. On the client side, `client.New(ctx, target, options...)` returns a remote implementation of that same interface. There is no separate public Wire execution model that every consumer has to learn.

For native Ferret, the ferret/uapi adapter connects the engine to those shared contracts. Wire then carries operations across the process boundary.

That separation is useful for tool authors. Query-running code can depend on `api.Runtime`, while connection setup remains a separate concern.

A remote call is still a remote call, of course. Latency, cancellation, and connection failures matter. Sharing an interface does not make those differences disappear; it keeps them from requiring an entirely different way to describe and execute a query.

The examples in this introduction follow the published releases pinned in the [runnable Wire guide]({{< ref "docs/embedding/wire" >}}): Wire `v{{< data "versions.wire.v1" >}}`, Core `v{{< data "versions.wire.core" >}}`, and Universal API `v{{< data "versions.wire.api" >}}`. Execution returns [consumable output]({{< ref "docs/embedding/go/output" >}}); the client collects it before releasing its runtime.

### Connecting to our inventory service

Let's return to the inventory query:

{{< code lang="fql" >}}
RETURN inventory::search({
    category: @category,
    available: true
})
{{< /code >}}

In our example, the inventory service already embeds Ferret and registers `inventory::search`, backed by the application's inventory data. We now want to execute that query from a separate development tool.

There are two sides to that connection: the application exposes its existing runtime, and the tool connects to it.

### First, expose the application's runtime

We do not need to build another engine or register the inventory function again. We take the engine the application already uses and adapt it to the Universal API with `uapi.Wrap`. This adapter borrows the engine without taking ownership of it.

The following function accepts that configured engine and exposes it through Wire on `127.0.0.1:50051`. The `coreVersion` argument identifies the embedded Ferret Core version. `server.New` constructs the server; blocking `Run` opens the TCP listener and manages shutdown.

{{< code lang="go" >}}
import (
    "context"
    "errors"

    "github.com/MontFerret/api"
    "github.com/MontFerret/ferret/v2/pkg/engine"
    "github.com/MontFerret/ferret/v2/uapi"
    "github.com/MontFerret/wire/server"
)

func serveInventory(
    ctx context.Context,
    eng *engine.Engine,
    coreVersion api.Version,
) (err error) {
    // eng already contains the application's inventory::search function.
    hostRuntime := uapi.Wrap(eng, coreVersion)

    srv, err := server.New(hostRuntime)
    if err != nil {
        return err
    }

    defer func() {
        // The serving context may already be canceled.
        err = errors.Join(err, srv.Shutdown(context.Background()))
    }()

    return srv.Run(ctx, "127.0.0.1:50051")
}
{{< /code >}}

The inventory service runs this function as part of its application lifecycle. Canceling `ctx` initiates shutdown with `Run`'s default 30-second cleanup budget. If that wait expires, the deferred call still waits for retained cleanup to settle before the caller can release the engine. The application still owns its engine: stopping the Wire endpoint does not close it.

*This example uses unencrypted localhost TCP for development. It does not authenticate connecting processes. The host must choose appropriate transport protection and access controls before exposing its runtime; Wire does not supply that security policy.* 

### Now, run the query from another process

With the inventory service running, our development tool connects to the same address: `127.0.0.1:50051`.

This time, we do not create a Ferret engine. `client.New` creates and owns a channel and returns a remote implementation of `api.Runtime`. We explicitly select plaintext for this loopback example; the default is verified TLS with no plaintext fallback. For a caller-owned gRPC channel, `client.From(ctx, conn)` borrows it instead.

{{< code lang="go" >}}
import (
    "context"
    "errors"

    "github.com/MontFerret/api"
    "github.com/MontFerret/wire/client"
)

func searchInventory(
    ctx context.Context,
    category string,
) (content *api.Content, err error) {
    remote, err := client.New(ctx, "127.0.0.1:50051", client.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer func() {
        err = errors.Join(err, remote.Close())
    }()

    output, err := remote.Run(
        ctx,
        api.NewSource("inventory.fql", `
            RETURN inventory::search({
                category: @category,
                available: true
            })
        `),
        api.WithParam("category", category),
        api.WithOutputContentType("application/json"),
    )
    if err != nil {
        return nil, err
    }
    defer func() { err = errors.Join(err, output.Close()) }()

    return output.Collect(ctx)
}
{{< /code >}}

Calling `searchInventory(ctx, "camping")` sends our original query and its category parameter to the inventory service. Wire invokes the hosted runtime's `Run` operation. The client collects its encoded output, settles cleanup, and returns detached content. Available content may accompany a terminal error; callers should inspect both. Compilation and execution happen against the engine we exposed in the first example—the one where `inventory::search` is registered.

Notice what is absent from the client: no inventory function registration, no inventory database configuration, and no duplicate Ferret engine. Those belong to the application. The tool supplies the query and parameters; the application supplies the execution environment.

There is one distinction worth making: this example already knows the name and calling convention of `inventory::search`. Wire does not *yet* expose a function catalog for tools to discover automatically. Executing a known host function works; discovering available functions and their signatures is a separate, future capability. 

The query stays the same. What changes is that our tool can now run it where its application-specific functionality actually lives.

## Debugging belongs with execution

Running a query against the right runtime is useful. Understanding what happened inside that query is just as important.

Suppose the inventory query returns something unexpected. Reproducing the problem in a different runtime might help but it might also remove the conditions that caused the problem.

Wire carries source-level debugging operations across the same boundary. A client can set breakpoints, continue execution, step into or over FQL calls, inspect stack frames and local variables, and evaluate expressions in a stopped frame. These operations use the Universal API's debugger contracts.

Diagnostics travel with structure, too. Rather than reducing a compiler error to an arbitrary string, Wire preserves the runtime's typed diagnostics, including source information, annotations, ranges, hints, and notes. That gives a consuming tool the information it needs to present a useful error at the relevant location.

This is debugging of the FQL program, not a replacement for debugging the host application's Go code. Likewise, the runtime protocol is one part of a tooling integration: adapting it to an editor's user interface or debugging protocol remains a separate responsibility.

## Your application remains in control

Wire follows the same ownership principle as Ferret itself.

It does not construct a replacement engine or take over the application. The host supplies the runtime and chooses an endpoint: `Run` opens a TCP listener, or `Serve` accepts one created by the host. Wire manages its own server-side resources without assuming ownership of the configured runtime. Shutting down Wire does not shut down the engine that the application supplied.

“Remote” also does not have to mean another machine. A development tool and an application can communicate through a private Unix socket on the same computer. The host still chooses and secures that endpoint.

That security responsibility matters. Wire does not provide an authentication, authorization, or TLS policy on the application's behalf. The host must decide who can connect and which runtime capabilities those clients should be allowed to exercise. The [Wire security guide]({{< ref "docs/embedding/wire-security" >}}) shows paired TLS, mTLS, and token-middleware configurations.

If a query can invoke a powerful host function, exposing query execution also exposes a route to that function. A Wire endpoint should therefore be treated as access to an application's execution environment—not merely a convenient place to upload query text.

## An early foundation, with a practical purpose

Wire is currently available as an alpha release. Its API and protocol are still evolving, so early adopters should pin compatible versions rather than assume a stable cross-version contract.

The implementation already has integration coverage that crosses the complete path: a native Ferret engine, the Universal API adapter, a Wire server, real protobuf/gRPC serialization, and the public Wire client. That coverage includes execution, reusable sessions, resource lifecycle, compiler diagnostics, and debugging.

The next question is how to make that foundation useful in everyday workflows. A command-line tool, an interactive environment, or an editor integration should be able to work with an application-provided runtime without each one inventing its own remote execution protocol.

That is the direction Wire is intended to support—not by becoming a larger application framework, but by giving those tools a shared connection point.

The [runnable Wire guide]({{< ref "docs/embedding/wire" >}}) contains the host and client programs. The [Wire repository](https://github.com/MontFerret/wire/tree/v{{< data "versions.wire.v1" >}}) contains protocol documentation and current limitations for anyone interested in trying it.

For me, the motivation comes back to the original library-first idea: embedding a language should not mean giving up useful development tools, and using those tools should not require rebuilding the application around them.

Your application owns the runtime. Wire gives your tools a way to reach it.

Happy ferreting! 🦦