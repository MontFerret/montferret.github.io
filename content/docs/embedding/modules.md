---
title: "Modules"
sidebarTitle: "Modules"
weight: 60
draft: false
description: "Install, configure, register, and manage modules in an embedded Ferret runtime."
---

# Modules

## Overview

Modules add capabilities to an embedded Ferret runtime. Official, contributed, third-party, and application-defined modules all use the same registration model:

1. Import the module as a Go package.
2. Construct and configure it in the host application.
3. Pass it to the engine with `ferret.WithModules`.
4. Let `ferret.New` register it and run its initialization hooks.
5. Close module-owned resources through the engine lifecycle.

The host application chooses which modules are available. Importing a package alone does not add its functions or runtime capabilities to Ferret.

## Installing a module

Install a module with `go get`. For example, the HTML module is published from the Ferret `contrib` repository:

{{< terminal command="true" >}}
go get github.com/MontFerret/contrib/modules/web/html
{{</ terminal >}}

Import the module and any supporting packages your configuration needs:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)
{{</ code >}}

Application-defined modules can live in the same Go module as the host application and do not need to be published separately.

## Creating and configuring a module

Construct modules before creating the engine. Constructors and functional options belong to the module package, so their configuration differs by module.

The HTML module requires a driver. This configuration uses the in-process memory driver for static HTML:

{{< code lang="go" >}}
htmlmod, err := html.New(
    html.WithDefaultDriver(memory.New()),
)
if err != nil {
    log.Fatal(err)
}
{{</ code >}}

Handle constructor errors before calling `ferret.New`. At this point, the module has not joined the engine lifecycle yet.

## Registering modules

Pass constructed modules to `ferret.WithModules` when creating the engine:

{{< code lang="go" >}}
package main

import (
    "log"

    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func main() {
    htmlmod, err := html.New(
        html.WithDefaultDriver(memory.New()),
    )
    if err != nil {
        log.Fatal(err)
    }

    engine, err := ferret.New(
        ferret.WithStdlib(stdlib.Safe()),
        ferret.WithModules(htmlmod),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    // Compile and run FQL programs with engine.
}
{{</ code >}}

`ferret.New` calls the module's `Register` method while assembling the engine. The engine is returned only after every module has registered, the host registries have been built, and all engine initialization hooks have completed.

## Registering multiple modules

`ferret.WithModules` accepts any number of modules:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithModules(
        htmlmod,
        sqlitemod,
        appmod,
    ),
)
{{</ code >}}

Modules register sequentially in the order passed. Multiple `WithModules` options append to the same ordered list.

Modules may add different functions to the same namespace. A shared namespace is not itself an error. If two modules register the same fully qualified function name, however, engine construction fails instead of choosing one implementation or overriding the earlier registration.

## Registration errors

Module construction and engine registration report errors at different points:

| Failure | Result |
| --- | --- |
| A module constructor returns an error | Handle the error before calling `ferret.New` |
| A nil module is passed to `WithModules` | Engine option validation fails |
| A module's `Register` method returns an error | Registration stops; later modules are not registered |
| Host registries cannot be built, such as after a duplicate function registration | Engine construction fails after module registration |
| An engine initialization hook returns an error | Initialization stops and engine construction fails |

In each engine-construction failure that occurs after bootstrap begins, Ferret runs the engine close hooks registered so far in reverse registration order. Close-hook errors are aggregated with the construction error, so cleanup failures are not discarded.

Always check the error returned by `ferret.New`. When construction fails, it returns no usable engine.

## Module lifecycle

A module participates in the engine lifecycle through registration and hooks:

| Stage | Behavior |
| --- | --- |
| Construction | The host application creates and configures the module |
| Registration | `ferret.New` calls each module's `Register` method in order |
| Host build | Ferret finalizes functions, parameters, codecs, and host services |
| Initialization | Engine initialization hooks run in registration order and stop on the first error |
| Shutdown | `engine.Close()` runs engine close hooks in reverse registration order and aggregates their errors |

The `module.Module` interface does not define a `Close` method. A module that owns engine-scoped resources registers an engine close hook and releases them there. Plan- and session-scoped resources belong in their corresponding close hooks and are released when the host application closes those plans and sessions.

## Writing a custom module

Applications can implement the module contract to expose their own functions, namespaces, host values, codecs, parameters, and lifecycle behavior:

{{< code lang="go" >}}
type Module interface {
    Name() string
    Register(Bootstrap) error
}
{{</ code >}}

`Register` receives a `module.Bootstrap` with access to host registries and the engine, plan, and session hook registrars. Keep registration focused on wiring the module into the host; construct and validate module-specific configuration before registration when practical.

For a complete implementation walkthrough, see [Writing plugins]({{< ref "/docs/guides/writing-plugins" >}}).

## Lifecycle hooks

Lifecycle hooks let modules and host applications react to engine, compilation, execution, and cleanup events. Host applications can register them through `ferret` option functions. Modules register the same hooks through `module.Bootstrap`.

### Engine hooks

| Hook | When it runs | Signature |
| --- | --- | --- |
| `WithEngineInitHook` | During `ferret.New`, after registration and host construction | `func() error` |
| `WithEngineCloseHook` | When construction fails or `engine.Close()` is called | `func() error` |

Initialization hooks run in registration order (FIFO) and stop on the first error. Close hooks run in reverse registration order (LIFO), continue after errors, and aggregate those errors.

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

Modules use the engine hook registrar:

{{< code lang="go" >}}
func (m *Module) Register(boot module.Bootstrap) error {
    boot.Hooks().Engine().OnInit(func() error {
        return m.initialize()
    })
    boot.Hooks().Engine().OnClose(func() error {
        return m.close()
    })

    return nil
}
{{</ code >}}

### Compilation hooks

| Hook | When it runs | Signature |
| --- | --- | --- |
| `WithBeforeCompileHook` | Before compilation begins | `func(ctx context.Context) error` |
| `WithAfterCompileHook` | After a compilation attempt | `func(ctx context.Context, compileErr error) error` |

Before-compile hooks run FIFO and stop on the first error. Once compilation is attempted, after-compile hooks run LIFO and receive the compilation error, if any. They continue after hook errors and aggregate them.

Modules register these hooks with `boot.Hooks().Plan().BeforeCompile(...)` and `boot.Hooks().Plan().AfterCompile(...)`.

### Execution hooks

| Hook | When it runs | Signature |
| --- | --- | --- |
| `WithBeforeRunHook` | Before `session.Run` begins | `func(ctx context.Context) (context.Context, error)` |
| `WithAfterRunHook` | After a run attempt | `func(ctx context.Context, runErr error) error` |

Before-run hooks run FIFO and can return a derived context for subsequent hooks and VM execution. They stop on the first error. Once execution is attempted, after-run hooks run LIFO, receive the run error, and aggregate hook errors.

Modules register these hooks with `boot.Hooks().Session().BeforeRun(...)` and `boot.Hooks().Session().AfterRun(...)`.

### Cleanup hooks

| Hook | When it runs | Module registrar |
| --- | --- | --- |
| `WithPlanCloseHook` | When `plan.Close()` is called | `boot.Hooks().Plan().OnClose(...)` |
| `WithSessionCloseHook` | When `session.Close()` is called | `boot.Hooks().Session().OnClose(...)` |

Plan and session close hooks run in LIFO order, continue after errors, and aggregate those errors.

### Hook execution order

| Hook type | Order | Stops on error |
| --- | --- | --- |
| Init, Before | FIFO | Yes |
| After, Close | LIFO | No; errors are aggregated |

Hooks registered through engine options and modules join the same hook chains. Engine-option hooks are registered while options are processed; module hooks are added later as modules register in the order passed to `WithModules`.

## Next steps

{{< docs-related tiles="embedding-getting-started,embedding-custom-functions,embedding-configuration,guide-writing-plugins" >}}
