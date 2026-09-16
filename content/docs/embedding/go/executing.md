---
title: "Executing Ferret"
sidebarTitle: "Executing Ferret"
weight: 20
draft: false
description: "Use the Go Engine, Plan, and Session lifecycle to execute Ferret programs."
aliases:
    - /docs/embedding/overview/
---

# Executing Ferret

Embedding Ferret means importing the `github.com/MontFerret/ferret/v2` module and using its Go API to compile and execute FQL queries inside your own application.

The host application stays in control: it decides which functions are available, which parameters are passed in, which standard library groups are loaded, and which modules extend the runtime. Scripts written in FQL describe the extraction, transformation, or automation logic without needing access to the rest of the application.

## When to embed

Embedding is useful when:

- extraction or transformation rules need to be user-defined, versioned, or changed without redeploying the application
- you want a small DSL for filters, mappings, or validation checks in a configuration-driven system
- pipeline steps or automation logic should be expressed declaratively instead of hard-coded in Go
- you need to sandbox what a script can see and do by controlling the available functions and capabilities

## The execution lifecycle

Every embedded execution follows the same four-step lifecycle:

```
Engine  →  Plan  →  Session  →  Output
```

### Engine

The `Engine` is the entry point. It is created once with all configuration — functions, parameters, modules, standard library selection, logging, and concurrency limits. The engine owns the compiler and the shared host state that all queries use.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
    ferret.WithParam("api_url", "https://example.com"),
)

if err != nil {
    log.Fatal(err)
}

defer engine.Close()
{{</ code >}}

An engine is safe for concurrent use by multiple goroutines.

### Plan

A `Plan` is a compiled query. Compiling a query validates the syntax, generates bytecode, and prepares it for execution. A plan can be reused across many sessions without recompilation.

{{< code lang="go" >}}
plan, err := engine.Compile(ctx,
    ferret.NewSource("greeting.fql", `return @greeting`),
)

if err != nil {
    log.Fatal(err)
}

defer plan.Close()
{{</ code >}}

Plans are safe for concurrent use. They manage an internal pool of virtual machines that sessions borrow from.

### Session

A `Session` is a single execution of a plan. It holds the VM state, per-run parameters, and encoding configuration. Sessions are not safe for concurrent use — each goroutine should create its own session.

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("greeting", "hello"),
)

if err != nil {
    log.Fatal(err)
}

defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    log.Fatal(err)
}
{{</ code >}}

### Output

The `Output` contains the encoded result of the execution. By default, the result is encoded as JSON.

{{< code lang="go" >}}
fmt.Println(string(output.Content))
// "hello"
fmt.Println(output.ContentType)
// application/json
{{</ code >}}

## Resource management

Every level of the lifecycle — Engine, Plan, and Session — must be closed when no longer needed. Use `defer` to ensure cleanup:

{{< code lang="go" >}}
engine, err := ferret.New()
if err != nil {
    log.Fatal(err)
}
defer engine.Close()

plan, err := engine.Compile(ctx,
    ferret.NewSource("query.fql", `return 1 + 1`),
)
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

session, err := plan.NewSession(ctx)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

output, err := session.Run(ctx)
{{</ code >}}

Close directly created children before parents: session, plan, then engine. Callers own these resources; closing a parent does not close descendants. Native Ferret protects each object's state, while the host orchestrates graceful shutdown. Cancel and settle outstanding work before closing resources used by hooks or services. Cancel and wait for ordinary execution before closing its session. Debug session closure terminates and settles active commands.

Engine closure rejects new compilation and loading without waiting for admitted operations. Plan closure rejects new sessions and wakes session creation waiting for capacity, even while close hooks are running. Engine closure does not cancel those waits. Creation that has acquired capacity may finish or fail on local VM pool closure; parent closure does not revoke a successfully returned session. Callers must close any returned resources.

Close is safe to repeat concurrently and retains the completed cleanup result. A close hook must not recursively close its own object. Closing a plan closes idle VMs; borrowed VMs remain with their sessions until returned.

Debugger termination events preserve the execution cause and any retained-resource cleanup failures. Later `Close` calls retain those cleanup failures, including concurrent and repeated calls. Cancellation with successful cleanup does not itself make `Close` fail.

A direct terminal failure from retained debug execution during `Start` or a resume command settles after-run hooks immediately with the original execution error. Later `Close` releases resources without repeating those hooks or replacing the error with cancellation. A runtime-error stop also settles after-run hooks, but remains paused and inspectable until the next resume or Close.

Pass non-nil contexts to compilation, execution, session creation, and debugger commands. Already-canceled contexts fail before options or hooks run; cancellation and deadline identities remain available through `errors.Is`.

Closing a plan releases its VM pool. Closing the engine runs all registered close hooks and releases engine-scoped resources. Close hooks execute in reverse registration order (LIFO) so that resources are torn down in the correct dependency order.

## Shorthand execution

For one-shot queries that do not need plan reuse, the engine provides a `Run` method that compiles, executes, and cleans up in a single call:

{{< code lang="go" >}}
output, err := engine.Run(ctx, ferret.NewAnonymousSource(`return 1 + 1`))
if err != nil {
    log.Fatal(err)
}

fmt.Println(string(output.Content))
// 2
{{</ code >}}

`Session.Run` returns successful encoded output even when an after-run hook or result cleanup fails. `Engine.Run` preserves that output and joins execution, hook, session cleanup, and plan cleanup errors. Preserve available output before handling the error:

{{< code lang="go" >}}
output, err := engine.Run(ctx, ferret.NewAnonymousSource(`return 2`))
if output != nil {
    fmt.Println(string(output.Content))
}
if err != nil {
    return err
}
{{</ code >}}

The returned `Output` contains encoded data and has no `Close` method. Ferret releases the internal VM result after encoding; the returned content remains available after the session, plan, and engine close. Failed execution or encoding returns no output.

## Thread safety summary

| Type | Concurrent use | Notes |
|------|---------------|-------|
| `Engine` | Safe | Shared across goroutines |
| `Plan` | Safe | VM pool handles concurrent session creation |
| `Session` | Not safe | One session per goroutine |

## Next steps

{{< docs-related tiles="embedding-go-getting-started,embedding-go-programs,language" >}}
