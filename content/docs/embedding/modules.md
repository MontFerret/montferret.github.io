---
title: "Modules & Hooks"
sidebarTitle: "Modules & Hooks"
weight: 60
draft: false
description: "Write self-contained extensions and react to engine lifecycle events."
---

# Modules & Hooks

Modules are self-contained extensions that bundle functions, parameters, codecs, and lifecycle hooks into a single registerable unit. Hooks let both modules and the host application react to events in the engine lifecycle.

## The Module interface

A module implements two methods:

{{< code lang="go" >}}
type Module interface {
    Name() string
    Register(Bootstrap) error
}
{{</ code >}}

`Name` returns an identifier for logging and diagnostics. `Register` receives a `Bootstrap` context that provides access to the engine's registries and hook system.

## Bootstrap context

The `Bootstrap` passed to `Register` exposes two surfaces:

### Host context

`bootstrap.Host()` returns a `HostContext` with access to:

| Method | Purpose |
|--------|---------|
| `Library()` | Register functions and namespaces |
| `Params()` | Set default parameters |
| `Encoding()` | Register encoding codecs |
| `Logger()` | Access the engine logger |
| `FileSystem()` | Access the engine file system |

### Hook registrar

`bootstrap.Hooks()` returns a `HookRegistrar` for subscribing to lifecycle events (see [Lifecycle hooks](#lifecycle-hooks) below).

## Registering modules

Pass modules to the engine with `WithModules`:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithModules(
        &metricsModule{},
        &cacheModule{},
    ),
)
{{</ code >}}

Modules are registered in order. Each module's `Register` is called during engine construction, before the engine is returned to the caller.

## Writing a module

A module that registers a `METRICS` namespace and logs query execution times:

{{< code lang="go" >}}
package metrics

import (
    "context"
    "log"
    "time"

    "github.com/MontFerret/ferret/v2/pkg/module"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Module struct {
    prefix string
}

func New(prefix string) *Module {
    return &Module{prefix: prefix}
}

func (m *Module) Name() string {
    return "metrics"
}

func (m *Module) Register(boot module.Bootstrap) error {
    // Register functions
    ns := boot.Host().Library().Namespace("METRICS")

    ns.Function().A1().Add("TAG", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
        return runtime.NewString(m.prefix + "." + arg.String()), nil
    })

    // Register hooks
    boot.Hooks().BeforeRun(func(ctx context.Context) (context.Context, error) {
        return context.WithValue(ctx, ctxStartKey{}, time.Now()), nil
    })

    boot.Hooks().AfterRun(func(ctx context.Context, runErr error) error {
        start, _ := ctx.Value(ctxStartKey{}).(time.Time)
        elapsed := time.Since(start)
        log.Printf("[%s] query took %s (err=%v)", m.prefix, elapsed, runErr)
        return nil
    })

    return nil
}

type ctxStartKey struct{}
{{</ code >}}

Use it in the host application:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithModules(metrics.New("myapp")),
)
{{</ code >}}

## Lifecycle hooks

Hooks let you react to events at different stages of the engine lifecycle. They can be registered through modules or directly on the engine with option functions.

### Engine hooks

| Hook | When it runs | Signature |
|------|-------------|-----------|
| `WithEngineInitHook` | After engine construction | `func() error` |
| `WithEngineCloseHook` | When `engine.Close()` is called | `func() error` |

Init hooks run in registration order (FIFO) and stop on the first error. Close hooks run in reverse registration order (LIFO) so that dependencies are torn down correctly.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithEngineInitHook(func() error {
        log.Println("engine initialized")
        return nil
    }),
    ferret.WithEngineCloseHook(func() error {
        log.Println("engine closing")
        return nil
    }),
)
{{</ code >}}

### Compilation hooks

| Hook | When it runs | Signature |
|------|-------------|-----------|
| `WithBeforeCompileHook` | Before each `Compile` call | `func(ctx context.Context) error` |
| `WithAfterCompileHook` | After each `Compile` call | `func(ctx context.Context, compileErr error) error` |

Before-compile hooks run FIFO and stop on the first error. After-compile hooks run LIFO and always execute, even if compilation failed — the compile error is passed as the second argument.

### Execution hooks

| Hook | When it runs | Signature |
|------|-------------|-----------|
| `WithBeforeRunHook` | Before `session.Run` | `func(ctx context.Context) (context.Context, error)` |
| `WithAfterRunHook` | After `session.Run` | `func(ctx context.Context, runErr error) error` |

Before-run hooks can return a modified context that is passed to subsequent hooks and the VM. This is useful for injecting request-scoped values like tracing spans or deadlines.

After-run hooks always execute, receiving the run error if one occurred.

### Cleanup hooks

| Hook | When it runs | Signature |
|------|-------------|-----------|
| `WithPlanCloseHook` | When `plan.Close()` is called | `func() error` |
| `WithSessionCloseHook` | When `session.Close()` is called | `func() error` |

Both run in LIFO order.

## Hook execution order

All "before" hooks run in registration order (FIFO) — first registered, first executed. If any before hook returns an error, the remaining hooks are skipped and the operation fails.

All "after" and "close" hooks run in reverse registration order (LIFO) — last registered, first executed. This ensures that resources are cleaned up in the correct dependency order. After and close hooks continue executing even if an earlier hook fails; errors are aggregated.

| Hook type | Order | Stops on error |
|-----------|-------|---------------|
| Init, Before | FIFO | Yes |
| After, Close | LIFO | No (aggregated) |

## Hooks via options vs. modules

The same hooks are available through both paths:

- **Engine options** (`WithBeforeRunHook`, etc.) — for hooks that belong to the host application
- **Module Bootstrap** (`boot.Hooks().BeforeRun(...)`) — for hooks that belong to a self-contained extension

Both register into the same hook chain. Module hooks are registered during engine construction in the order modules are listed.

## Next steps

{{< docs-related tiles="embedding-custom-functions,embedding-host-values,embedding-value-encoders" >}}
