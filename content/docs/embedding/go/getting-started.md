---
title: "Getting Started"
sidebarTitle: "Getting Started"
weight: 10
draft: false
description: "Install the library, run your first query from Go, and handle the result."
aliases:
    - /docs/embedding/getting-started/
---

# Getting Started

This page walks through installing Ferret as a Go dependency, running a query, and working with the result.

## Installation

These Go embedding pages target the **Core development API**, using Universal API `v1.0.0-alpha.21` and consumable output. The website's Standard Library version listing retains its separate release pin. The [Wire example's compatibility table]({{< ref "docs/embedding/wire#run-a-host-and-client" >}}) lists its published dependency set for remote runtimes.

The documented checkout requires **Go 1.26 or later**. It includes both the reader/writer Encoding API and consumable Output API. Use this checkout as a local module replacement. From the directory containing your Go application (with its own `go.mod`), clone Core into a sibling directory and select the revision documented here:

{{< terminal >}}
git clone https://github.com/MontFerret/ferret.git ../ferret-core
git -C ../ferret-core checkout 3dcc1849b08685e9cdc71021fd2b6dd9d4003077
go mod edit -replace=github.com/MontFerret/ferret/v2=../ferret-core
go get github.com/MontFerret/ferret/v2
{{</ terminal >}}

After adding the Go source below, run `go mod tidy`. Initialize a new application's module first with `go mod init example.com/myapp`. The replacement resolves Core locally at the exact revision above; it does not depend on the version selected by an unpinned `go get`. Remove it when upgrading to a compatible published release.

For an existing integration, follow the [consumable output migration]({{< ref "/docs/migrations/consumable-output" >}}) and the [Encoding API migration]({{< ref "/docs/migrations/encoding-api" >}}).

## Running a query

The simplest way to execute a query is `engine.Run`. It compiles the source, runs it in a fresh session, and returns a one-shot output handle. Collect it to obtain owned encoded bytes.

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
)

func main() {
    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        ferret.NewAnonymousSource(`return { name: "Ferret", version: 2 }`),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer output.Close()

    content, err := output.Collect(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(content.Data))
    // {"name":"Ferret","version":2}
}
{{</ code >}}

`ferret.NewAnonymousSource` is convenient for one-shot queries. Use `ferret.NewSource(name, content)` when the source name should appear in compilation errors and debugger output. Both constructors return `ferret.Source` values.

Native applications normally use the root `ferret` package for sources, engine and session options, log levels, module contracts, parameter values, program formats, and output. Import a package under `pkg` only when using the specialized API owned by that package.

## Compiling and reusing a plan

When the same query runs many times — with different parameters, in different goroutines, or on a schedule — compile it once and create sessions from the resulting plan:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx,
    ferret.NewSource("greeting.fql", `return upper(@name)`),
)
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

names := []string{"alice", "bob", "carol"}

// This snippet also imports "errors".
for _, name := range names {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("name", name),
    )
    if err != nil {
        log.Fatal(err)
    }

    output, err := session.Run(ctx)
    if err != nil {
        session.Close()
        log.Fatal(err)
    }

    content, collectErr := output.Collect(ctx)
    output.Close()
    closeErr := session.Close()
    if err := errors.Join(collectErr, closeErr); err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(content.Data))
}
// "ALICE"
// "BOB"
// "CAROL"
{{</ code >}}

The plan manages an internal pool of virtual machines. Sessions borrow a VM from the pool and return it on close, so creating many sessions from the same plan is efficient.

## Passing parameters

Parameters let the host application inject values into a query at runtime. In FQL, parameters are referenced with the `@` prefix.

Engine-level parameters apply to every session:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithParam("base_url", "https://api.example.com"),
)
{{</ code >}}

Session-level parameters override engine defaults for a single execution:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("user_id", 42),
    ferret.WithSessionParam("base_url", "https://staging.example.com"),
)
{{</ code >}}

You can inspect which parameters a compiled query declares:

{{< code lang="go" >}}
params := plan.Params()
fmt.Println(params)
// [base_url user_id]
{{</ code >}}

See [Parameters]({{< ref "/docs/embedding/go/parameters" >}}) for the full parameter API.

## Handling errors

Ferret returns standard Go errors. Compilation errors include source location information:

{{< code lang="go" >}}
_, err := engine.Compile(ctx,
    ferret.NewSource("bad.fql", `return @`),
)
if err != nil {
    fmt.Println(err)
    // compilation error with line and column
}
{{</ code >}}

Preparation failures are returned by `Run`. Once execution enters the VM, execution, encoding, hooks, and output cleanup failures are reported by `Collect` or `Consume`:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    return err // preparation failed
}
defer output.Close()

content, err := output.Collect(ctx)
if err != nil {
    // Inspect available content independently, then handle the terminal error.
    return err
}
_ = content
{{</ code >}}

Context cancellation and timeouts work as expected:

{{< code lang="go" >}}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

output, err := session.Run(ctx)
if err != nil {
    return err // cancellation may prevent VM admission
}
defer output.Close()

_, err = output.Collect(ctx)
if err != nil {
    // errors.Is(err, context.DeadlineExceeded) detects an expired deadline.
    return err
}
{{</ code >}}

Keep the original Run context alive until output consumption completes. A new consumption context cannot revive a canceled invocation. See [Output]({{< ref "/docs/embedding/go/output" >}}) for borrowing, cancellation, and closure.

The VM observes cancellation at structural execution boundaries rather than polling every native operation. Blocking host functions, iterators, queries, streams, and other context-aware capabilities receive this same context and must observe it while they retain control. Cancellation and deadline errors propagate to the caller and cannot be suppressed by FQL error recovery.

## Next steps

{{< docs-related tiles="embedding-go-output,embedding-go-executing,embedding-go-modules,embedding-go-custom-functions,embedding-go-configuration" >}}
