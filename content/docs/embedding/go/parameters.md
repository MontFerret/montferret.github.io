---
title: "Parameters"
sidebarTitle: "Parameters"
weight: 30
draft: false
description: "Pass data into queries at engine and session level."
aliases:
    - /docs/embedding/parameters/
---

# Parameters

Parameters let the host application inject values into FQL queries. They are the primary way to pass dynamic data — user IDs, URLs, configuration values, thresholds — from Go into a script without string interpolation or query rewriting.

{{% notification type="warning" %}}
FQL source is code. Do not insert runtime data with `fmt.Sprintf` or string concatenation. See [Construct FQL safely]({{< ref "/docs/embedding/go/safe-fql-construction" >}}) for URL, selector, HTML, and source-composition examples.
{{% /notification %}}

## FQL parameter syntax

In FQL, parameters are referenced with the `@` prefix:

{{< code lang="fql" >}}
let result = @base_url + "/users/" + to_string(@user_id)
return result
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

## Inspecting referenced parameters

A compiled plan records the parameters referenced by the query. Use `plan.Params()` to get their names in first-seen order:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.NewAnonymous(`
    let url = @base_url + "/users/" + to_string(@user_id)
    return url
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

`WithParams`, `WithParam`, and their session equivalents convert Go values through `runtime.ValueOf`:

| Go input | Runtime value |
|---------|--------------|
| `nil` | `None` |
| An existing `runtime.Value` | The same value, without conversion |
| `bool` | `Boolean` |
| `string` | `String` |
| `int`, `int8`, `int16`, `int32`, `int64` | `Int` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `Int`, when the value fits in `int64` |
| `float32`, `float64` | `Float` |
| `time.Time` | `DateTime` |
| `time.Duration` | `Duration` |
| `[]byte` | `Binary` |
| Slices and arrays | `Array`, with elements converted recursively |
| Maps | `Object`, with values converted recursively and keys represented as strings |
| Structs | `Object`, with exported fields converted recursively |
| Pointers | The pointed-to value converted recursively; a nil pointer becomes `None` |

Struct field names are used exactly as declared in Go. Unexported fields are skipped, and `runtime.ValueOf` does not read struct tags or flatten embedded fields. The scalar cases above apply to the concrete built-in types; a defined scalar type must implement `runtime.Value` or be converted to a supported Go type first. Unsupported values, nested unsupported values, and unsigned integers larger than `math.MaxInt64` return an error.

`runtime.ValueOf(nil)` returns `None`, so nil entries work in the map-based `WithParams` and `WithSessionParams` options. The single-value `WithParam` and `WithSessionParam` options reject a nil `any`; pass `runtime.None` through the corresponding runtime-value option when you need an explicit `None`.

This conversion produces in-memory `runtime.Value` instances for execution. It is separate from result encoding: `Session.Run` returns a `*ferret.Output`, whose `Content` contains encoded bytes. The default output codec is JSON.

## Example: parameterized query with per-session overrides

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    engine, err := ferret.New(
        ferret.WithParam("greeting", "hello"),
        ferret.WithParam("punctuation", "!"),
    )
    if err != nil {
        return err
    }
    defer engine.Close()

    ctx := context.Background()

    plan, err := engine.Compile(ctx, source.NewAnonymous(`
        return concat(@greeting, " ", @name, @punctuation)
    `))
    if err != nil {
        return err
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
            return err
        }

        output, runErr := session.Run(ctx)
        closeErr := session.Close()
        if err := errors.Join(runErr, closeErr); err != nil {
            return err
        }

        fmt.Println(string(output.Content))
    }
    // "hello Alice!"
    // "hey Bob!"
    // "hi Carol!"

    return nil
}
{{</ code >}}

The `greeting` parameter is set at the engine level but overridden per session for Bob and Carol. The `punctuation` parameter uses the engine default for all sessions.

## Next steps

{{< docs-related tiles="language-parameters,embedding-go-configuration,embedding-go-modules" >}}
