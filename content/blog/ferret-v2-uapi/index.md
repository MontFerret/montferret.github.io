---
title: "One Interface, Many Runtimes: Introducing Ferret's Universal API"
subtitle: "Separating how Ferret runs from how applications and tools interact with it."
date: "2026-11-17"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
---

Hello, friends!

One of the ideas I've emphasized throughout Ferret's development is that **Ferret is a library first**.

You embed it into your application, configure the engine, register your own functions and modules, and decide how queries interact with the surrounding environment.

This approach has worked well. But as Ferret's ecosystem grew, it introduced a problem I hadn't needed to worry about in the early days.

We were no longer building just an execution engine.

We were building a CLI, a test runner, a debugger, IDE integrations, and eventually a way to execute queries in a runtime hosted by another application.

All of these components needed to interact with Ferret. But should every one of them depend directly on the internals of Ferret Core?

I didn't think so.

That's why we introduced the **Universal Ferret API**, or UAPI: a standalone set of interfaces and types that describes how to interact with a Ferret runtime without depending on a particular engine implementation.

It's a relatively small piece of the architecture, but it changes how the rest of the ecosystem fits together.

## The problem wasn't executing queries

Running a query through Ferret's native Go API is straightforward.

Create an engine, configure it, provide some FQL, and execute it. For applications embedding Ferret directly, that's exactly what we want.

But consider a development tool.

A test runner shouldn't need to know how Ferret allocates VM registers, constructs execution contexts, or manages internal engine components.

An editor doesn't need to know how a compiler is assembled just to compile a query and display its diagnostics.

And a debugger shouldn't have to understand the engine's private data structures to inspect a running program.

These tools need capabilities, not implementation details.

Originally, the obvious solution was to use Ferret Core's Go API everywhere. After all, it already existed and provided everything we needed.

But that creates coupling.

When the engine's internal architecture changes, tools may need to change with it. When a new execution environment appears, every consumer needs another integration path.

That becomes especially problematic when the engine isn't running inside the same process.

Imagine a CLI that can execute FQL using either an embedded Ferret engine or a runtime exposed by another application.

If the CLI has to understand two entirely different execution APIs, we've duplicated a significant part of its execution logic.

Now multiply that problem by the number of tools in the ecosystem.

The issue isn't that any individual integration is particularly difficult. It's that **the same integration problem keeps appearing in different places**.

## Separating the engine from its consumers

The solution was to establish a boundary between two responsibilities:

**Implementing a runtime** and **using a runtime**.

Ferret Core is responsible for implementing the language: parsing, compilation, optimization, VM execution, native values, and the host integrations that make queries useful.

A consumer has a different set of concerns.

It needs to submit source code, compile it, create execution sessions, provide parameters, retrieve results, and perhaps debug a running query.

Those operations don't inherently require access to Ferret Core's internal machinery.

They require a contract.

This is what the Universal API provides.

At its center is the `api.Runtime` interface, defined in the standalone [`github.com/MontFerret/api`](https://github.com/MontFerret/api) module.

The interface exposes operations for executing source, compiling reusable plans, preparing programs for debugging, and retrieving runtime version information.

The important part is what it doesn't expose: a particular compiler implementation, VM, storage strategy, or transport.

Those decisions belong to the runtime implementation.

A consumer works with the contract, not the machinery behind it.

## One contract, several ways to execute

Here's the basic idea:

```text
                    CLI / Lab / IDE
                          |
                          v
                     api.Runtime
                          |
               +----------+----------+
               |                     |
               v                     v
         Native Adapter         Wire Client
               |                     |
               v                     v
          Ferret Core           Wire Server
                                     |
                                     v
                                Hosted Runtime
```

From the consumer's perspective, both paths expose `api.Runtime`.

One path operates locally, using Ferret Core through its Universal API adapter.

The other communicates with a runtime hosted elsewhere, using Wire.

They have different implementations, ownership arrangements, performance characteristics, and failure modes.

But the consumer doesn't need a separate execution model for each one.

This is a form of dependency inversion: tools depend on the shared contract, while implementations adapt themselves to that contract.

And because UAPI exists independently of Ferret Core, the dependency doesn't point back into the engine's implementation packages.

That distinction matters.

It gives us room to change the internals of Ferret Core without making every tool part of the same refactoring exercise.

It also means the interface isn't exclusive to Ferret Core. Another implementation can satisfy the same contracts, provided it implements the required behavior.

## Why Ferret has two APIs

This raises a reasonable question.

*Why not make the Universal API the only API?*

Because the two APIs serve different purposes.

The **Native API** is designed for applications embedding and configuring Ferret Core. It exposes engine-specific capabilities and integrates directly with Ferret's Go implementation.

The **Universal API** is designed for consumers that need to interact with a runtime without depending on how it was constructed.

Trying to combine both responsibilities into a single interface would create an awkward compromise.

Either the shared API would become increasingly coupled to Ferret Core's implementation, or the Native API would lose useful capabilities to remain sufficiently generic.

Neither outcome is desirable.

Instead, Ferret provides an adapter.

If your application owns an existing native engine, you can expose it through UAPI:

{{< code lang="go" >}}
portable := uapi.Wrap(native, "dev")
{{< /code >}}

`Wrap` doesn't create another engine. It adapts the one you already have.

The application remains responsible for the native engine's lifetime, configuration, functions, and modules. The adapter provides the common interface.

Alternatively, `uapi.New` can construct and own a native engine when that ownership model is more convenient.

Both approaches produce a runtime that implements `api.Runtime`.

The engine remains specialized. Its consumers gain a common way to interact with it.

## What using the Universal API looks like

Let's make this a little more concrete.

Suppose an application exposes a custom `REPORTS::GENERATE` function through its Ferret runtime.

A separate component wants to execute a query against that runtime.

With UAPI, its execution code can look like this:

{{< code lang="go" >}}
func runReport(
    ctx context.Context,
    runtime api.Runtime,
    period string,
) (content *api.Content, err error) {
    output, err := runtime.Run(
        ctx,
        api.NewSource(
            "report.fql",
            "RETURN REPORTS::GENERATE(@period)",
        ),
        api.WithParam("period", period),
    )
    if err != nil {
        return nil, err
    }

    defer func() {
        err = errors.Join(err, output.Close())
    }()

    return output.Collect(ctx)
}
{{< /code >}}

The function accepts an `api.Runtime`. It doesn't construct a native engine, configure the reports module, or establish a network connection.

It simply uses a runtime provided by its caller.

That runtime might be a local Ferret engine wrapped by UAPI. It might be a remote runtime obtained from Wire.

The execution code remains the same.

Of course, the supplied runtime must actually provide `REPORTS::GENERATE`. A shared interface doesn't mean every runtime exposes identical application-specific functions.

But that is precisely the separation we want.

The application defines the execution environment. The consumer provides the query.

### Reusable plans and execution sessions

Not every query should be compiled from scratch on every execution.

UAPI therefore defines more than a one-shot `Run` method.

A runtime can compile source into an `api.Plan`. The plan represents a compiled program and can be reused to create independent execution sessions.

Each `api.Session` receives its own execution configuration, including parameters and supported filesystem settings.

This lets a consumer compile a query once and execute it repeatedly with different inputs without assuming anything about the engine's internal representation of a compiled program.

Plans aren't consumed by execution, and independently created sessions have their own lifetimes.

That gives tools a useful common model for repeated execution while leaving the implementation free to manage compilation and execution resources appropriately.

## Results are part of the contract, too

One aspect of the Universal API that deserves particular attention is execution output.

Historically, Ferret's execution APIs returned buffered results. That was convenient for small queries, but not every consumer wants to accumulate an entire result in memory.

A CLI may want to print output.

An HTTP service may want to forward it.

A test runner may need to collect and inspect it.

These are different consumption patterns, but they shouldn't require different execution APIs.

The current UAPI uses a **consumable output** model.

Execution returns an `api.Output` handle. The caller can then choose how to receive the encoded result:

- `Collect` materializes available content into caller-owned memory.
- `Consume` delivers chunks to a callback, allowing output to be forwarded without retaining the complete payload.
- `Close` abandons output and releases its owned resources.

There's an important lifecycle detail here: obtaining an output handle doesn't mean execution has completed successfully.

Execution, encoding, delivery, and cleanup errors can become visible during consumption. Even a query run only for its side effects must settle its output to observe completion.

The output contract gives local and remote implementations a common way to describe result ownership, completion, cancellation, and errors.

It doesn't promise that every implementation evaluates queries incrementally or eliminates buffering internally.

It gives consumers a reliable way to receive results without forcing every runtime into the same execution strategy.

## Debugging without exposing the engine

Execution is only part of the story.

A useful language ecosystem needs tools that can inspect programs while they're running.

For Ferret, that means source-level debugging: breakpoints, stepping, stack frames, variables, and expression evaluation.

The tempting approach would be to expose the native debugger directly and let every integration adapt it to its needs.

But that would bring us back to the original problem.

An IDE shouldn't need to understand the VM's internal representation of a stack frame. Nor should a debugger client depend on the exact storage structures used for local variables.

UAPI defines portable debugger contracts instead.

A tool can request breakpoints, inspect frames, retrieve variables, evaluate expressions, and control execution through the shared debugger interface.

Ferret Core implements those operations through its native debugger. Wire can carry the same operations across a process boundary.

The interface also provides structured source locations and diagnostics, allowing tools to associate errors and execution state with the original FQL source.

This doesn't remove all the complexity of debugger development. An IDE still needs an integration layer for its user interface and debugging protocol.

But it puts an important boundary in the right place.

**Tools should understand the execution state of an FQL program, not the private implementation of Ferret's virtual machine.**

## Where Wire fits

You may have noticed that I've mentioned Wire several times.

That's not accidental.

[Wire](https://github.com/MontFerret/wire) is one of the clearest demonstrations of why UAPI exists.

But the two projects solve different problems.

UAPI defines how to interact with a runtime. Wire defines how to reach one across a process boundary.

Wire provides a Protocol Buffers and gRPC implementation of the shared runtime contracts.

On the server side, it accepts an `api.Runtime` supplied by the host application.

On the client side, it exposes an implementation of that same interface.

This means a development tool can interact with a remotely hosted runtime without learning an entirely different API for compilation, sessions, results, and debugging.

The host remains responsible for configuring its engine and deciding which functions, modules, and capabilities are available.

Wire handles communication.

The distinction is important because a remote call is not equivalent to a local one in every respect.

Network latency, disconnections, authentication, cancellation, and resource lifetimes still matter. UAPI doesn't make these differences disappear.

What it does is prevent them from forcing an entirely separate programming model onto every consumer.

I cover the remote execution side in more detail in [Introducing Ferret Wire](/blog/wire-introduction/).

## What UAPI is not

There's a danger with abstractions like this: once you introduce a common interface, it becomes tempting to turn it into a framework for everything.

That's not the goal.

UAPI isn't another Ferret engine. It doesn't compile or execute FQL by itself.

It isn't a module registry, a plugin system, or a replacement for the native embedding API.

And it doesn't define how an application should configure its runtime, enforce security policies, or expose host-specific functionality.

It's also not a promise that all implementations have identical capabilities or performance characteristics.

UAPI defines a shared set of operations and their expected behavior. An implementation still has to provide those operations correctly.

Keeping that distinction clear is important, especially as Ferret continues to evolve.

The abstraction should remain small enough to understand and useful enough that different parts of the ecosystem can rely on it.

## Building an ecosystem without coupling everything together

The immediate benefits are fairly practical.

The CLI can support different execution environments without maintaining completely separate query execution models.

Lab can run tests against different runtime adapters while preserving its own responsibilities for test discovery, scheduling, and assertions.

The daemon and IDE integrations can use common compilation, execution, and debugging contracts.

And applications embedding Ferret can expose their configured runtimes without giving external tools direct access to the engine's internal API.

But the longer-term benefit is architectural.

Ferret Core can evolve its compiler, VM, optimizations, memory management, and execution machinery independently of the tools that consume it.

Meanwhile, those tools can evolve independently of a particular runtime implementation.

That doesn't eliminate compatibility work. The Universal API itself is still in alpha, and its contracts will continue to evolve as we discover and address shortcomings.

But it gives us a much better place to manage that evolution.

Instead of allowing dependencies to spread across the entire ecosystem, we have an explicit boundary where they belong.

## A small abstraction with a larger purpose

When I started working on Ferret, I wasn't thinking about building an entire ecosystem of development tools.

I needed a language that applications could embed to query and transform data.

Over time, that grew into a compiler, a VM, a CLI, a test runner, a debugger, and integrations with increasingly different execution environments.

That growth made one thing clear: building Ferret as a library was the right starting point, but **building a library and defining a portable contract for its consumers are two different architectural decisions**.

The Universal API is the second one.

It lets Ferret Core remain an implementation of the language rather than becoming a dependency that every tool must understand intimately.

It gives us a shared foundation for execution and debugging, while leaving applications in control of their own runtime environments.

And it creates room for implementations and integrations we haven't built yet.

That's exactly what I want from an abstraction: not more layers for their own sake, but fewer assumptions about how the pieces of a system must fit together.

### Getting started

The Universal API is open source and currently available as an alpha release.

You can explore the implementation and integration guides here:

- [Universal Ferret API on GitHub](https://github.com/MontFerret/api)
- [Ferret Core](https://github.com/MontFerret/ferret)
- [Universal API embedding guide](https://ferretlang.org/docs/embedding/go/universal-api/)
- [Ferret Wire](https://github.com/MontFerret/wire)

Happy ferreting! 🦦