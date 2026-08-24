---
title: "Migrate from Ferret v1"
sidebarTitle: "Migrate from v1"
weight: 15
draft: false
description: "Move a Go application from Ferret v1 compatibility packages to the native Ferret v2 embedding API."
---

# Migrate a Go application from Ferret v1

Moving an embedded Go application to Ferret v2 has three separate parts:

1. Use `ferret migrate` for supported source and import changes.
2. Use the v2 compatibility packages temporarily where they cover your application.
3. Replace the v1 compiler, runtime, and driver composition with the native v2 Engine, Plan, Session, and module APIs.

The migration command is a compatibility aid, not a general application translator. Plan to review every change and migrate application-specific composition manually.

## Start with the mechanical migration

Run the compatibility check from the Go module root before changing files:

{{< terminal command="true" >}}
ferret migrate check --from v1 .
{{< /terminal >}}

Preview the supported changes, then apply them:

{{< terminal command="true" >}}
ferret migrate run --dry-run .
ferret migrate run --print .
ferret migrate run .
{{< /terminal >}}

The command can update supported FQL behavior, including an implicit final top-level `FOR`, and rewrite these v1 imports to their v2 compatibility equivalents:

| Ferret v1 import | Ferret v2 compatibility import |
| --- | --- |
| `github.com/MontFerret/ferret` | `github.com/MontFerret/ferret/v2/compat` |
| `github.com/MontFerret/ferret/pkg/compiler` | `github.com/MontFerret/ferret/v2/compat/compiler` |
| `github.com/MontFerret/ferret/pkg/runtime` | `github.com/MontFerret/ferret/v2/compat/runtime` |
| `github.com/MontFerret/ferret/pkg/runtime/core` | `github.com/MontFerret/ferret/v2/compat/runtime/core` |
| `github.com/MontFerret/ferret/pkg/runtime/values` | `github.com/MontFerret/ferret/v2/compat/runtime/values` |
| `github.com/MontFerret/ferret/pkg/runtime/values/types` | `github.com/MontFerret/ferret/v2/compat/runtime/values/types` |

Generated Go files and imports without a documented compatibility replacement are left unchanged and reported for manual follow-up. See the [Migrate command]({{< ref "/docs/tools/cli/migrate" >}}) for discovery rules, transaction behavior, and all available flags.

FQL migration applies to discovered lowercase `.fql` files. It does not locate FQL embedded in Go string literals, so review and format embedded queries such as the example below manually.

### Migrate driver imports manually

The command does not rewrite any import below `github.com/MontFerret/ferret/pkg/drivers`. In particular, these v1 imports have no compatibility-package replacement:

- `github.com/MontFerret/ferret/pkg/drivers`
- `github.com/MontFerret/ferret/pkg/drivers/http`
- `github.com/MontFerret/ferret/pkg/drivers/cdp`

In native v2 code, HTML support is a module from `github.com/MontFerret/contrib/modules/web/html`. Replace the v1 HTTP driver with the module's `drivers/memory` package. Replace the v1 CDP driver with the module's `drivers/cdp` package and register it with the same HTML module.

If unsupported v1 imports remain anywhere in a module-wide migration, the command retains the v1 Ferret dependency. Remove it only after all remaining imports and application logic have been migrated.

## Treat compatibility packages as temporary

The compatibility packages preserve selected v1 compiler, runtime, function, and value shapes on top of the v2 engine. They are useful for getting core-only code compiling while native composition is migrated in smaller steps.

They are not the preferred long-term embedding API. They also do not make old drivers work with v2. The compatibility compiler has no option for registering native v2 modules, so an application that calls HTML functions cannot combine its migrated compatibility compiler with the contrib HTML module. Such an application must move the engine and module composition to the native API before its HTML path can run on v2.

Use the compatibility stage to stabilize supported code, tests, and FQL changes. Then remove `/compat` imports rather than building new application features on them.

## Compare the application architectures

The following example compiles one query, caches it, registers one host function, and runs it with a URL parameter. Static HTML is the default. Setting `FERRET_CDP_ADDRESS` selects the optional CDP driver without changing the query or duplicating the application.

### Historical Ferret v1.0.0 code

{{% notification type="warning" %}}
This complete example is historical Ferret v1.0.0 code. Do not copy its APIs into a new application.
{{% /notification %}}

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "time"

    "github.com/MontFerret/ferret/pkg/compiler"
    "github.com/MontFerret/ferret/pkg/drivers"
    "github.com/MontFerret/ferret/pkg/drivers/cdp"
    httpdriver "github.com/MontFerret/ferret/pkg/drivers/http"
    "github.com/MontFerret/ferret/pkg/runtime"
    "github.com/MontFerret/ferret/pkg/runtime/core"
    "github.com/MontFerret/ferret/pkg/runtime/values"
    "github.com/MontFerret/ferret/pkg/runtime/values/types"
)

const query = `
LET page = DOCUMENT(@url, { driver: @driver })
RETURN {
    url: page.url,
    title: NORMALIZE(ELEMENT(page, "h1").innerText)
}
`

type (
    App struct {
        program    *runtime.Program
        httpDriver *httpdriver.Driver
        cdpDriver  *cdp.Driver
        driverName string
    }

    Result struct {
        URL   string `json:"url"`
        Title string `json:"title"`
    }
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    app, err := NewApp(os.Getenv("FERRET_CDP_ADDRESS"))
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, app.Close())
    }()

    runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    data, err := app.Run(runCtx, "https://example.com")
    if err != nil {
        return err
    }

    var result Result
    if err := json.Unmarshal(data, &result); err != nil {
        return err
    }

    fmt.Printf("%s: %s\n", result.Title, result.URL)
    return nil
}

func NewApp(cdpAddress string) (*App, error) {
    comp := compiler.New()
    if err := comp.RegisterFunction("NORMALIZE", normalize); err != nil {
        return nil, err
    }

    program, err := comp.Compile(query)
    if err != nil {
        return nil, err
    }

    app := &App{
        program:    program,
        httpDriver: httpdriver.NewDriver(),
        driverName: httpdriver.DriverName,
    }

    if cdpAddress != "" {
        app.cdpDriver = cdp.NewDriver(cdp.WithAddress(cdpAddress))
        app.driverName = cdp.DriverName
    }

    return app, nil
}

func (app *App) Run(ctx context.Context, url string) ([]byte, error) {
    ctx = drivers.WithContext(ctx, app.httpDriver, drivers.AsDefault())
    if app.cdpDriver != nil {
        ctx = drivers.WithContext(ctx, app.cdpDriver)
    }

    return app.program.Run(ctx,
        runtime.WithParam("url", url),
        runtime.WithParam("driver", app.driverName),
    )
}

func (app *App) Close() error {
    var cdpErr error
    if app.cdpDriver != nil {
        cdpErr = app.cdpDriver.Close()
    }

    return errors.Join(cdpErr, app.httpDriver.Close())
}

func normalize(_ context.Context, args ...core.Value) (core.Value, error) {
    if err := core.ValidateArgs(args, 1, 1); err != nil {
        return values.None, err
    }
    if err := core.ValidateType(args[0], types.String); err != nil {
        return values.None, err
    }

    return values.NewString(strings.TrimSpace(args[0].String())), nil
}
{{</ code >}}

This application caches a `runtime.Program`, attaches drivers to each execution context, and passes values through `runtime.WithParam`. `Program.Run` returns raw JSON bytes. The application also owns and closes both driver instances.

### Native Ferret v2 code

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "time"

    htmlmodule "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/cdp"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/encoding"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

const query = `let page = document(@url, { driver: @driver })
return { url: page.url, title: normalize(element(page, "h1").innerText) }`

type (
    App struct {
        engine     *ferret.Engine
        plan       *ferret.Plan
        driverName string
    }

    Result struct {
        URL   string `json:"url"`
        Title string `json:"title"`
    }
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    app, err := NewApp(ctx, os.Getenv("FERRET_CDP_ADDRESS"))
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, app.Close())
    }()

    runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    output, err := app.Run(runCtx, "https://example.com")
    if err != nil {
        return err
    }

    var result Result
    if err := json.Unmarshal(output.Content, &result); err != nil {
        return err
    }

    fmt.Printf("%s: %s (%s)\n", result.Title, result.URL, output.ContentType)
    return nil
}

func NewApp(ctx context.Context, cdpAddress string) (*App, error) {
    htmlOptions := []htmlmodule.Option{
        htmlmodule.WithDefaultDriver(memory.New()),
    }
    driverName := memory.DriverName

    if cdpAddress != "" {
        htmlOptions = append(htmlOptions,
            htmlmodule.WithDrivers(cdp.New(cdp.WithAddress(cdpAddress))),
        )
        driverName = cdp.DriverName
    }

    htmlModule, err := htmlmodule.New(htmlOptions...)
    if err != nil {
        return nil, err
    }

    engine, err := ferret.New(
        ferret.WithModules(htmlModule),
        ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
            ns.Function().A1().Add("normalize", normalize)
        }),
    )
    if err != nil {
        return nil, err
    }

    plan, err := engine.Compile(ctx, source.New("extract-title.fql", query))
    if err != nil {
        return nil, errors.Join(err, engine.Close())
    }

    return &App{
        engine:     engine,
        plan:       plan,
        driverName: driverName,
    }, nil
}

func (app *App) Run(ctx context.Context, url string) (*encoding.Output, error) {
    session, err := app.plan.NewSession(ctx,
        ferret.WithSessionParam("url", url),
        ferret.WithSessionParam("driver", app.driverName),
    )
    if err != nil {
        return nil, err
    }

    output, runErr := session.Run(ctx)
    closeErr := session.Close()

    return output, errors.Join(runErr, closeErr)
}

func (app *App) Close() error {
    planErr := app.plan.Close()
    engineErr := app.engine.Close()

    return errors.Join(planErr, engineErr)
}

func normalize(_ context.Context, arg runtime.Value) (runtime.Value, error) {
    value, err := runtime.CastArg[runtime.String](arg, 0)
    if err != nil {
        return nil, err
    }

    return runtime.NewString(strings.TrimSpace(value.String())), nil
}
{{</ code >}}

The memory driver is always registered as the default. When `FERRET_CDP_ADDRESS` is set, the application registers CDP as an additional named driver and passes `cdp` as the per-session `driver` parameter. The URL remains data rather than being interpolated into FQL source.

For a one-off host function, `WithFunctionsRegistrar` keeps registration close to engine construction. If the application exposes a related set of functions, values, services, or lifecycle hooks, implement an application [module]({{< ref "/docs/embedding/go/modules" >}}) and register it with `WithModules` instead.

## Map v1 concepts to native v2

| Ferret v1 | Native Ferret v2 | Migration note |
| --- | --- | --- |
| `compiler.New()` | `ferret.New(...)` | The Engine owns compiler and shared host configuration. Construction can fail, so handle its error. |
| `drivers.WithContext` | `html.New(...)` plus `ferret.WithModules(...)` | Register capabilities once during Engine construction instead of attaching drivers to every caller context. |
| `pkg/drivers/http` | `modules/web/html/drivers/memory` | Use the memory driver for static HTTP loading, parsing, and in-memory DOM operations. Its registered name is `memory`, not `http`. |
| `pkg/drivers/cdp` | `modules/web/html/drivers/cdp` | Register CDP as the default or as an additional named HTML driver. It connects to an existing browser endpoint. |
| `compiler.RegisterFunction` | `WithFunctionsRegistrar`, `WithFunctions`, or an application module | Host functions are finalized when the Engine is built. Do not mutate the function set after construction. |
| `runtime.WithParam` | `ferret.WithParam` or `ferret.WithSessionParam` | Engine parameters are shared defaults. Session parameters are per-execution values and override engine defaults. |
| `compiler.Compile` returning `*runtime.Program` | `engine.Compile` returning `*ferret.Plan` | A Plan is the reusable compiled query and owns its VM pool. |
| `program.Run(ctx)` | `plan.NewSession(ctx)` then `session.Run(ctx)` | Create a separate Session for each concurrent execution and close it after use. |
| Cached `*runtime.Program` | Cached `*ferret.Plan` | Plans are safe for concurrent use; Sessions are not. Create one Session per goroutine or request. |
| Raw JSON `[]byte` | `*encoding.Output` | Read encoded bytes from `Output.Content` and the selected MIME type from `Output.ContentType`. JSON remains the default codec. |
| Manual driver cleanup | Session, Plan, Engine, and module lifecycle | Close directly created Sessions, then Plans, then the Engine. Modules use lifecycle hooks for resources they own, and runtime-owned HTML values are closed as execution results are materialized. Resources retained by the caller remain the caller's responsibility. |
| Driver registration and cancellation in one context | Engine composition plus caller execution context | Pass the caller context to Compile, NewSession, and Run. Modules and host functions receive a derived context that preserves cancellation and deadlines. |

See [Parameters]({{< ref "/docs/embedding/go/parameters" >}}) for conversion and override rules, [Custom Functions]({{< ref "/docs/embedding/go/custom-functions" >}}) for typed function registration, and [Executing Ferret]({{< ref "/docs/embedding/go/executing" >}}) for the complete lifecycle and concurrency contract.

## Preserve cancellation and ownership

Registration no longer travels through the execution context, but cancellation still does. Use the request, job, or process context supplied by the caller. Pass it to `engine.Compile`, `plan.NewSession`, and `session.Run` rather than replacing it with `context.Background()` inside application helpers.

`Plan.NewSession` observes cancellation while waiting for execution capacity. `Session.Run` passes the execution context to the VM, module hooks, HTML operations, and host functions. Blocking host code must observe that context while it retains control. Cancellation and deadline errors propagate to the caller.

The host owns every Engine, Plan, and Session it constructs directly. Close children before parents:

1. Close each Session to return its VM and run session hooks.
2. Close the reusable Plan to release its VM pool and run plan hooks.
3. Close the Engine to run module and engine close hooks and release engine-scoped services.

`Engine.Run` is different: it owns and closes the temporary Plan and Session it creates for a one-shot execution. Use the explicit Plan and Session path when caching compiled work.

## Finish the migration

Before removing the v1 dependency:

- replace every `/compat` import with its native v2 package or API
- replace every v1 driver import and context registration
- run the same Plan through multiple Sessions with representative parameters
- test caller cancellation, deadlines, and cleanup on both success and failure
- check `Output.ContentType` before assuming a result encoding other than the JSON default

Native composition is complete when the application constructs one Engine with its functions and modules, caches reusable Plans, creates Sessions for individual executions, and no longer imports Ferret v1 or v2 compatibility packages.

## Next steps

{{< docs-related tiles="tools-cli-migrate,embedding-go-modules,embedding-go-parameters,embedding-go-custom-functions,embedding-go-executing" >}}
