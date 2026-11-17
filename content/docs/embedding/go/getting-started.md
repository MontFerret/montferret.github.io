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

From your Go application module, install Ferret Core v{{< data "versions.core.v2" >}}. This release requires Go {{< data "versions.go" >}} or later. For a new application, run `go mod init example.com/myapp` first.

{{< terminal >}}
go get github.com/MontFerret/ferret/v2@v{{< data "versions.core.v2" >}}
{{</ terminal >}}

For an existing integration, follow the [consumable output migration]({{< ref "/docs/migrations/consumable-output" >}}) and the [Encoding API migration]({{< ref "/docs/migrations/encoding-api" >}}).

## Running a query

The simplest way to execute a query is `engine.Run`. It compiles the source, runs it in a fresh session, and returns a one-shot output handle. Collect it to obtain owned encoded bytes.

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    engine, err := ferret.New()
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, engine.Close()) }()

    output, err := engine.Run(
        context.Background(),
        ferret.NewAnonymousSource(`return { name: "Ferret", version: 2 }`),
    )
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, output.Close()) }()

    content, err := output.Collect(context.Background())
    if err != nil {
        return err
    }

    fmt.Println(string(content.Data))
    // {"name":"Ferret","version":2}

    return nil
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
    return err
}
defer plan.Close()

names := []string{"alice", "bob", "carol"}

// This snippet also imports "errors".
for _, name := range names {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("name", name),
    )
    if err != nil {
        return err
    }

    output, err := session.Run(ctx)
    if err != nil {
        return errors.Join(err, session.Close())
    }

    content, collectErr := output.Collect(ctx)
    closeErr := errors.Join(output.Close(), session.Close())
    if err := errors.Join(collectErr, closeErr); err != nil {
        return err
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
// [name] for greeting.fql
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
