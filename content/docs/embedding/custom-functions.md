---
title: "Custom Functions"
sidebarTitle: "Custom Functions"
weight: 30
draft: false
description: "Register host-defined functions and namespaces that scripts can call."
---

# Custom Functions

Custom functions let the host application expose Go logic to FQL scripts. Functions can be registered individually or organized into namespaces.

## Registering functions

The most direct way to add functions is `WithFunctionsRegistrar`. The callback receives a `runtime.Namespace` where you register functions by arity:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
        ns.Function().A0().Add("NOW_UNIX", func(ctx context.Context) (runtime.Value, error) {
            return runtime.NewInt(int(time.Now().Unix())), nil
        })

        ns.Function().A1().Add("DOUBLE", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
            n, err := runtime.CastArg[runtime.Int](arg, 0)
            if err != nil {
                return nil, err
            }

            return runtime.NewInt(int(n * 2)), nil
        })
    }),
)
{{</ code >}}

Scripts can then call these functions directly:

{{< code lang="fql" >}}
RETURN { timestamp: NOW_UNIX(), doubled: DOUBLE(21) }
{{</ code >}}

## Function signatures

Ferret provides typed function signatures for each arity:

| Builder method | Go signature | FQL call |
|---------------|-------------|----------|
| `A0()` | `func(ctx) (Value, error)` | `FN()` |
| `A1()` | `func(ctx, arg) (Value, error)` | `FN(x)` |
| `A2()` | `func(ctx, arg1, arg2) (Value, error)` | `FN(x, y)` |
| `A3()` | `func(ctx, arg1, arg2, arg3) (Value, error)` | `FN(x, y, z)` |
| `A4()` | `func(ctx, arg1, arg2, arg3, arg4) (Value, error)` | `FN(a, b, c, d)` |
| `Var()` | `func(ctx, args ...Value) (Value, error)` | `FN(a, b, ...)` |

Fixed-arity functions (`A0` through `A4`) automatically validate the argument count. Variadic functions (`Var`) receive a slice and must validate the count themselves.

## Namespaces

Use `WithNamespace` with a library builder to create a named group of functions:

{{< code lang="go" >}}
lib := runtime.NewLibrary()
ns := lib.Namespace("CRYPTO")

ns.Function().A1().Add("MD5", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
    input := arg.String()
    hash := md5.Sum([]byte(input))
    return runtime.NewString(hex.EncodeToString(hash[:])), nil
})

ns.Function().A1().Add("SHA256", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
    input := arg.String()
    hash := sha256.Sum256([]byte(input))
    return runtime.NewString(hex.EncodeToString(hash[:])), nil
})

engine, err := ferret.New(
    ferret.WithNamespace(lib),
)
{{</ code >}}

Scripts call these as:

{{< code lang="fql" >}}
RETURN {
    md5: CRYPTO::MD5("hello"),
    sha: CRYPTO::SHA256("hello")
}
{{</ code >}}

Namespaces can be nested. A namespace created with `ns.Namespace("SUB")` produces functions accessible as `CRYPTO::SUB::FUNCTION_NAME`.

## Argument validation

For variadic functions, use the validation helpers from the `runtime` package to check argument count and types:

{{< code lang="go" >}}
fns.Var().Add("CONCAT_WITH", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
    // At least 2 arguments: separator + one or more values
    if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
        return nil, err
    }

    separator, err := runtime.CastArg[runtime.String](args[0], 0)
    if err != nil {
        return nil, err
    }

    parts := make([]string, 0, len(args)-1)
    for _, arg := range args[1:] {
        parts = append(parts, arg.String())
    }

    return runtime.NewString(strings.Join(parts, string(separator))), nil
})
{{</ code >}}

### Type casting helpers

The `runtime` package provides generic casting functions that return typed values with clear error messages:

| Helper | Purpose |
|--------|---------|
| `CastArg[T](arg, index)` | Cast a single argument to type `T` |
| `CastArgAt[T](args, index)` | Cast argument at position in a slice |
| `CastArgs[T](args)` | Cast all arguments to the same type |
| `CastArgs2[T1, T2](a, b)` | Cast two arguments to different types |
| `CastArgs3[T1, T2, T3](a, b, c)` | Cast three arguments to different types |

### Arity and type validation

| Helper | Purpose |
|--------|---------|
| `ValidateArgs(args, min, max)` | Check argument count is within range |
| `ValidateArgsType(args, types...)` | Check each argument matches expected type |
| `ValidateArgType(arg, pos, types...)` | Check a single argument matches one of the expected types |

## Merging pre-built function sets

If you have a `*runtime.Functions` object built separately, merge it into the engine with `WithFunctions`:

{{< code lang="go" >}}
builder := runtime.NewFunctionsBuilder()
builder.A1().Add("REVERSE", reverseFunc)
builder.A2().Add("REPEAT", repeatFunc)

funcs, err := builder.Build()
if err != nil {
    log.Fatal(err)
}

engine, err := ferret.New(
    ferret.WithFunctions(funcs),
)
{{</ code >}}

## Complete example

A custom namespace that provides string transformation functions:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    lib := runtime.NewLibrary()
    ns := lib.Namespace("TEXT")

    ns.Function().A1().Add("TITLE_CASE", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
        return runtime.NewString(strings.Title(arg.String())), nil
    })

    ns.Function().A2().Add("WRAP", func(ctx context.Context, text, wrapper runtime.Value) (runtime.Value, error) {
        w := wrapper.String()
        return runtime.NewString(w + text.String() + w), nil
    })

    ns.Function().Var().Add("JOIN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
        if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
            return nil, err
        }

        parts := make([]string, len(args))
        for i, arg := range args {
            parts[i] = arg.String()
        }

        return runtime.NewString(strings.Join(parts, "")), nil
    })

    engine, err := ferret.New(
        ferret.WithNamespace(lib),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        source.NewAnonymous(`
            RETURN {
                titled: TEXT::TITLE_CASE("hello world"),
                wrapped: TEXT::WRAP("content", "**"),
                joined: TEXT::JOIN("a", "b", "c")
            }
        `),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
    // {"joined":"abc","titled":"Hello World","wrapped":"**content**"}
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-modules,embedding-host-values,embedding-configuration,guide-writing-plugins" >}}
