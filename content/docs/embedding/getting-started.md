---
title: "Getting Started"
sidebarTitle: "Getting Started"
weight: 20
draft: false
description: "Install the library, run your first query from Go, and handle the result."
---

# Getting Started

This page walks through installing Ferret as a Go dependency, running a query, and working with the result.

## Installation

Add the module to your project:

{{< terminal >}}
go get github.com/MontFerret/ferret/v2
{{</ terminal >}}

## Running a query

The simplest way to execute a query is `engine.Run`. It compiles the source, runs it in a fresh session, and returns the encoded output.

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        source.NewAnonymous(`RETURN { name: "Ferret", version: 2 }`),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
    // {"name":"Ferret","version":2}
}
{{</ code >}}

`source.NewAnonymous` wraps a query string into a `*source.Source`. For named sources use `source.New(name, text)` — the name appears in error messages and debug output.

## Compiling and reusing a plan

When the same query runs many times — with different parameters, in different goroutines, or on a schedule — compile it once and create sessions from the resulting plan:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.New("greeting", `RETURN UPPER(@name)`))
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

names := []string{"alice", "bob", "carol"}

for _, name := range names {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("name", name),
    )
    if err != nil {
        log.Fatal(err)
    }

    output, err := session.Run(ctx)
    session.Close()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
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

See [Parameters]({{< ref "parameters" >}}) for the full parameter API.

## Handling errors

Ferret returns standard Go errors. Compilation errors include source location information:

{{< code lang="go" >}}
_, err := engine.Compile(ctx, source.New("bad.fql", `RETURN @`))
if err != nil {
    fmt.Println(err)
    // compilation error with line and column
}
{{</ code >}}

Runtime errors from query execution are returned by `session.Run`:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    // handle runtime error
}
{{</ code >}}

Context cancellation and timeouts work as expected:

{{< code lang="go" >}}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

output, err := session.Run(ctx)
if err != nil {
    // may be context.DeadlineExceeded
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-modules,embedding-custom-functions,embedding-configuration,guide-writing-plugins" >}}
