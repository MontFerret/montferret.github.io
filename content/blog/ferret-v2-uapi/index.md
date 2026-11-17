---
title: "One Interface, Many Runtimes: Introducing Ferret's Universal API"
subtitle: "Separating how Ferret runs from how applications and tools interact with it."
date: "2026-11-17"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
---

Hello, friends!

One of the ideas I've emphasized throughout Ferret's development is that [Ferret is a library first]({{< ref "blog/ferret-v2-lib-first" >}}). Applications embed the engine, configure its capabilities, register their own functions and modules, and define how queries interact with the surrounding environment.

That model worked well while Ferret was primarily an execution engine, but the surrounding ecosystem gradually became much broader. The project expanded to include a CLI, a test runner, a debugger, IDE integrations, and support for executing queries against runtimes hosted by other applications.

All of these components need access to the same core capabilities: compiling source, executing queries, managing sessions and results, and debugging running programs. Allowing each of them to depend directly on Ferret Core would tightly couple the tooling ecosystem to the engine's implementation and make support for different runtime environments increasingly difficult.

To provide a stable boundary between runtimes and their consumers, we introduced the **Universal Ferret API (UAPI)**. It is a standalone module that defines a common set of interfaces and types for interacting with Ferret runtimes independently of their underlying implementations.

This separation gives the ecosystem a consistent integration model while allowing Ferret Core, alternative runtime implementations, and the tools built around them to evolve more independently.

## The problem wasn't executing queries

Ferret's native Go API makes it straightforward to embed and execute FQL within an application. The application creates an engine, configures its environment, and executes queries using the capabilities it needs. While this level of control is valuable for embedding, it isn't necessarily appropriate for development tools that interact with the runtime.

A test runner primarily needs to execute queries and inspect their results. An editor needs compilation and diagnostics, while a debugger needs to control execution and inspect program state. None of these tools should depend on how Ferret Core initializes its compiler, manages execution contexts, or represents internal data structures.

Initially, using Ferret Core's existing Go API throughout the ecosystem seemed like a reasonable approach. However, it would also couple every consumer to the engine's implementation. Changes to internal APIs could affect multiple tools, and supporting alternative execution environments would require additional integration logic.

This becomes particularly relevant when a tool supports both local and remote runtimes. For example, the CLI may execute a query using an embedded Ferret engine or delegate execution to a runtime hosted by another application. Although these environments differ in how they operate, the CLI still performs essentially the same tasks: submitting source code, supplying parameters, and processing results.

Supporting each environment through a separate execution API would introduce unnecessary duplication, not only in the CLI but also in Lab, IDE integrations, and other tools. As the ecosystem continued to grow, maintaining these parallel integrations would become increasingly difficult.

What we needed was a common interface that allowed tools to interact with different runtime implementations without coupling their execution logic to the underlying engine.

## Separating the engine from its consumers

The Universal API addresses this by establishing a clear separation between runtime implementations and the tools that consume them. Ferret Core remains responsible for implementing the language, including parsing, compilation, optimization, VM execution, native values, and integration with host applications. Consumers, on the other hand, interact with the runtime at a higher level, using its capabilities to compile queries, create execution sessions, supply parameters, retrieve results, and debug running programs.

These interactions are defined through a set of interfaces in the standalone [`github.com/MontFerret/api`](https://github.com/MontFerret/api) module. The central `api.Runtime` interface provides operations for executing source code, compiling reusable plans, preparing programs for debugging, and retrieving runtime version information. Additional interfaces describe compiled plans, execution sessions, and debugging operations, giving consumers a consistent way to interact with the runtime throughout the execution lifecycle.

Importantly, these contracts make no assumptions about the underlying compiler, virtual machine, or communication transport. Each runtime implementation remains free to manage its internal architecture and resources, provided it satisfies the shared contracts. This allows Ferret Core to evolve independently while giving development tools a stable integration boundary that can accommodate different runtime implementations.

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

From a consumer's perspective, both local and remote execution are represented through the same `api.Runtime` contract. A local implementation may adapt Ferret Core directly through UAPI, while a remote implementation may communicate with a runtime exposed through [Wire]({{< ref "docs/embedding/wire" >}}). The underlying behavior differs in important ways, including ownership, latency, resource management, and failure handling, but those differences do not require consumers to adopt separate execution models.

This structure follows a dependency-inversion approach: tools depend on the shared API, while runtime implementations adapt themselves to it. Because the Universal API is defined independently of Ferret Core, consumers are not coupled to the engine's implementation packages or internal architecture.

That separation allows Ferret Core to evolve without forcing corresponding changes throughout the tooling ecosystem. It also keeps the contract open to other runtime implementations, provided they satisfy the same behavioral requirements.

## Why Ferret has two APIs

The `Native API` and the `Universal API` are intentionally separate because they serve different roles within the architecture.

The Native API is designed for applications that embed and configure Ferret Core directly. It exposes engine-specific capabilities and gives the host access to the parts of the runtime that are necessary for configuring modules, functions, services, codecs, and other implementation-specific behavior.

The Universal API operates at a different level. It defines the contracts used by tools and other consumers that need to execute, compile, or debug FQL without depending on how the underlying runtime was constructed.

Combining these responsibilities into a single API would weaken that separation. A sufficiently rich shared interface would gradually become coupled to Ferret Core's implementation details, while a sufficiently generic Native API would have to omit capabilities that are useful to embedding applications. Keeping the two APIs distinct allows each one to remain focused on its intended use case.

Ferret connects them through an adapter. An application that already owns a native engine can expose it through UAPI without creating another runtime:

{{< code lang="go" >}}
portable := uapi.Wrap(native, "dev")
{{< /code >}}

`uapi.Wrap` adapts the existing engine to the `api.Runtime` contract while leaving ownership with the application. The host continues to manage the engine's lifetime, configuration, registered functions, modules, and other native capabilities.

For cases where the Universal API runtime should own the underlying engine, `uapi.New` constructs and manages one directly.

Both approaches expose the same `api.Runtime` interface to consumers while preserving the distinction between the specialized Native API and the portable contracts defined by UAPI.

## What using the Universal API looks like

Consider an application that exposes a custom `reports::generate` function through its Ferret runtime. A separate component can execute a query against that runtime without needing to know how the engine was configured or where it is running.

With UAPI, the execution code depends only on `api.Runtime`:

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
            "return reports::generate(@period)",
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

The function receives a runtime from its caller and uses the shared execution contract directly. It does not construct a native engine, configure the reports module, or establish a transport connection. Those responsibilities remain outside the execution logic.

The supplied runtime may be a local Ferret Core instance exposed through the UAPI adapter, or a remote runtime provided by Wire. In either case, the consumer uses the same `api.Runtime` interface and follows the same execution flow.

This does not imply that all runtimes provide identical application-specific capabilities. The runtime used in this example must still expose `reports::generate`, along with any other functions or modules required by the query. UAPI standardizes how consumers interact with the runtime; it does not standardize the runtime's domain-specific configuration.

This keeps the responsibilities clearly separated: the host defines the execution environment and its available capabilities, while the consumer operates against that environment through a common API.

### Reusable plans and execution sessions

For workloads that execute the same query multiple times, recompiling the source for every invocation would be unnecessary. UAPI therefore exposes compiled plans and execution sessions in addition to the convenience `Run` method.

A runtime can compile source into an `api.Plan`, which represents a reusable compiled program. The plan can then create independent `api.Session` instances, each with its own execution configuration, such as parameters and supported filesystem settings.

This allows consumers to separate compilation from execution and reuse the same compiled plan across multiple invocations with different inputs. At the same time, they remain independent of the runtime's internal representation of compiled code or session state.

Plans are not consumed when executed, and sessions manage their own lifecycles. The result is a common model for repeated execution that gives consumers predictable semantics while leaving each runtime implementation free to manage compilation and execution resources internally.

## Results are part of the contract

Execution output is another important part of the Universal API contract. Earlier Ferret APIs returned fully buffered results, which works well for many cases but is less suitable when consumers need to process or forward large outputs without materializing them entirely in memory.

Different consumers also have different requirements. A CLI may stream encoded output directly to the terminal, an HTTP service may forward it to another connection, while a test runner may need to collect the complete result for inspection and assertions. These use cases should not require separate execution models.

UAPI addresses this through a consumable output abstraction. Execution returns an `api.Output` handle, and the caller decides how to process it:

- `Collect` materializes the available result as caller-owned content.
- `Consume` delivers encoded data incrementally to a callback, allowing it to be forwarded or processed without retaining the complete payload.
- `Close` abandons the remaining output and releases resources owned by the handle.

The output handle also defines the lifecycle of an execution result. Successfully obtaining an `api.Output` means that execution has been admitted and a usable handle has been created; it does not mean the query has already completed successfully. Errors that occur during execution, encoding, delivery, or output cleanup may therefore surface when the output is consumed or collected. Even executions performed primarily for side effects must settle their output in order to observe their final outcome.

By including these semantics in the shared contract, UAPI gives both local and remote runtime implementations a consistent model for ownership, cancellation, completion, and error reporting. It does not require every implementation to evaluate queries incrementally or avoid buffering internally; those remain implementation details. What it standardizes is how consumers receive and manage execution results.

## Debugging without exposing the engine

The same separation applies to debugging. Ferret's tooling needs source-level capabilities such as breakpoints, stepping, stack inspection, variable access, and expression evaluation, but exposing the native debugger directly would couple every integration to the VM's internal representation of execution state.

UAPI defines portable debugger contracts for these operations instead. Tools can manage breakpoints, inspect frames and variables, evaluate expressions, and control execution through shared interfaces without depending on how Ferret Core represents stack frames, locals, or paused execution internally.

Ferret Core implements these contracts using its native debugger, while Wire can expose the same debugging operations across a process boundary. Structured source locations and diagnostics are part of the shared model as well, allowing clients to associate runtime state and errors with the corresponding FQL source.

This does not eliminate the need for tool-specific integration. An IDE still has to map these capabilities into its own debugging protocol and user interface. What UAPI provides is a stable boundary between those tools and the runtime implementation, so debugger clients operate on the observable state of an FQL program rather than the private structures of Ferret's virtual machine.

## Where Wire fits

Wire builds directly on the contracts defined by UAPI, but the two projects address different parts of the architecture. UAPI defines how consumers interact with a runtime, while Wire provides a transport for accessing that runtime across a process boundary.

Wire implements the shared runtime contracts using Protocol Buffers and gRPC. On the server side, the host supplies an `api.Runtime` that represents its configured execution environment. On the client side, Wire exposes a remote implementation of the same interface. As a result, tools can use the same APIs for compilation, reusable plans, sessions, output handling, and debugging regardless of whether the runtime is local or hosted by another process.

The host remains responsible for constructing and configuring the runtime, including its modules, functions, services, and other capabilities. Wire does not reproduce or manage that configuration; it provides the communication layer through which consumers access it.

Remote execution still introduces concerns that do not exist, or have different characteristics, in a local runtime. Network latency, connection failures, authentication, cancellation, and distributed resource lifecycles all remain relevant. UAPI does not attempt to abstract those differences away completely. Instead, it provides a common interaction model so that supporting a remote runtime does not require consumers to adopt a separate execution API.

The remote execution model and Wire's architecture are covered in more detail in [Introducing Ferret Wire](/blog/wire-introduction/).

## What UAPI is not

UAPI is intentionally limited in scope. It defines the contracts through which consumers interact with a Ferret runtime, but it is not itself a runtime implementation and does not compile or execute FQL on its own.

It also does not replace the Native API or take responsibility for runtime configuration. Module registration, host-specific functions, services, security policies, codecs, and other implementation-specific capabilities remain under the control of the runtime and its host application.

Similarly, UAPI does not require every implementation to provide identical performance characteristics or internal behavior. Local and remote runtimes may differ significantly in how they manage resources, perform I/O, or execute operations. What they share is a common set of observable contracts and lifecycle semantics.

Keeping this boundary narrow is important for the long-term design of the API. UAPI should provide enough structure for tools and runtime implementations to interoperate consistently without becoming a general-purpose framework for configuring or implementing Ferret runtimes.

## Building an ecosystem without coupling everything together

The practical benefit of UAPI is that different parts of the Ferret ecosystem can share the same execution model without depending directly on Ferret Core.

The CLI can support multiple runtime environments through the same set of contracts. Lab can execute tests against different runtime adapters while keeping test discovery, scheduling, retries, and assertions within its own domain. Daemons and IDE integrations can rely on common APIs for compilation, execution, and debugging, while applications that embed Ferret can expose their configured runtimes without exposing the engine's internal interfaces.

The broader benefit is architectural independence. Ferret Core can continue to evolve its compiler, virtual machine, optimization pipeline, memory management, and execution internals without requiring every consumer to evolve in lockstep. At the same time, tools can develop new capabilities without being tied to a specific runtime implementation.

UAPI does not eliminate compatibility concerns, and its contracts are still evolving during the alpha stage. What it provides is a defined boundary where those compatibility decisions can be managed explicitly, rather than allowing implementation dependencies to spread across the ecosystem.

## A small abstraction with a larger purpose

Ferret began as an embeddable language for querying and transforming data, but the surrounding ecosystem has grown to include a compiler, virtual machine, CLI, test runner, debugger, and integrations with multiple execution environments.

That growth introduced a distinction that was less important in the early stages of the project: building Ferret as a reusable library and defining a portable interface for its consumers are separate architectural concerns.

The Universal API addresses the latter. It allows Ferret Core to remain focused on implementing the language and its runtime, while tools depend on a shared set of execution and debugging contracts rather than on the engine's internal APIs.

This separation gives applications control over how their runtimes are configured and hosted, while providing consumers with a consistent way to interact with them. It also leaves room for new runtime implementations, transports, and tooling integrations without requiring them to reproduce the assumptions of the current engine.

The purpose of UAPI is therefore not to introduce abstraction for its own sake, but to reduce the number of assumptions that different parts of the ecosystem make about one another. That gives Ferret a cleaner foundation for evolving its runtime and tooling independently.

### Getting started

The Universal API is open source and currently available as an alpha release.

You can explore the implementation and integration guides here:

- [Universal API on GitHub](https://github.com/MontFerret/api)
- [Ferret Core](https://github.com/MontFerret/ferret)
- [Universal API embedding guide]({{<ref "docs/embedding/go/universal-api/" >}})
- [Wire](https://github.com/MontFerret/wire)
- [Wire embedding guide](https://ferretlang.org/docs/embedding/wire/)

Happy ferreting! 🦦