---
title: "Parameters"
sidebarTitle: "Parameters"
weight: 40
draft: false
description: "Pass data into queries at engine and session level."
---

# Parameters

Parameters let the host application inject values into FQL queries. They are the primary way to pass dynamic data — user IDs, URLs, configuration values, thresholds — from Go into a script without string interpolation or query rewriting.

## FQL parameter syntax

In FQL, parameters are referenced with the `@` prefix:

{{< code lang="fql" >}}
LET result = @base_url + "/users/" + TO_STRING(@user_id)
RETURN result
{{</ code >}}

## Engine-level parameters

Parameters set on the engine apply as defaults to every session. Use these for values that rarely change — base URLs, API keys, environment names.

### From Go values

`WithParams` accepts a `map[string]any` and converts each value to a runtime value automatically:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithParams(map[string]any{
        "base_url":    "https://api.example.com",
        "environment": "production",
        "max_retries": 3,
    }),
)
{{</ code >}}

`WithParam` sets a single parameter:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithParam("base_url", "https://api.example.com"),
    ferret.WithParam("timeout", 30),
)
{{</ code >}}

### From runtime values

When you already have a `runtime.Value`, use the runtime variants to skip conversion:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithRuntimeParam("threshold", runtime.NewFloat(0.95)),
    ferret.WithRuntimeParams(runtime.Params{
        "tag": runtime.NewString("v2"),
    }),
)
{{</ code >}}

## Session-level parameters

Session parameters override engine defaults for a single execution. Use these for per-request values — user context, request IDs, pagination offsets.

### From Go values

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionParams(map[string]any{
        "user_id": 42,
        "page":    1,
        "limit":   25,
    }),
)
{{</ code >}}

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("user_id", 42),
)
{{</ code >}}

### From runtime values

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionRuntimeParam("score", runtime.NewFloat(0.85)),
)
{{</ code >}}

## Inspecting declared parameters

A compiled plan knows which parameters the query declares. Use `plan.Params()` to get the list:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.NewAnonymous(`
    LET url = @base_url + "/users/" + TO_STRING(@user_id)
    RETURN url
`))
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

fmt.Println(plan.Params())
// [base_url user_id]
{{</ code >}}

This is useful for validating that all required parameters are provided before creating a session.

## Supported Go types

`WithParams` and `WithParam` accept standard Go types that are converted to runtime values via `runtime.ValueOf`:

| Go type | Runtime value |
|---------|--------------|
| `nil` | `None` |
| `bool` | `Boolean` |
| `int`, `int32`, `int64` | `Int` |
| `float64` | `Float` |
| `string` | `String` |
| `time.Time` | `DateTime` |

For complex structures like arrays and objects, build them using `runtime.NewArray()` and `runtime.NewObject()`, then pass them with the runtime param variants.

## Example: parameterized query with per-session overrides

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
    engine, err := ferret.New(
        ferret.WithParam("greeting", "hello"),
        ferret.WithParam("punctuation", "!"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    ctx := context.Background()

    plan, err := engine.Compile(ctx, source.NewAnonymous(`
        RETURN CONCAT(@greeting, " ", @name, @punctuation)
    `))
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    users := []struct{ name, greeting string }{
        {"Alice", "hello"},
        {"Bob", "hey"},
        {"Carol", "hi"},
    }

    for _, u := range users {
        session, err := plan.NewSession(ctx,
            ferret.WithSessionParam("name", u.name),
            ferret.WithSessionParam("greeting", u.greeting),
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
    // "hello Alice!"
    // "hey Bob!"
    // "hi Carol!"
}
{{</ code >}}

The `greeting` parameter is set at the engine level but overridden per session for Bob and Carol. The `punctuation` parameter uses the engine default for all sessions.
