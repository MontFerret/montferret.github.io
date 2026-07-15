---
title: "Embed Ferret in a Go application"
sidebarTitle: "Embedding the runtime"
weight: 110
draft: false
description: "Initialize the Ferret engine, compile queries, run them with parameters, and handle errors in a Go service."
---

# Embed Ferret in a Go application

Embedding Ferret lets you run FQL queries inside your Go application — useful for user-defined extraction rules, dynamic data pipelines, or sandboxed scripting. This guide walks through building a Go service that accepts queries over HTTP.

For the full API reference, see the [Embedding]({{< ref "/docs/embedding" >}}) section.

## Set up the project

Create a new Go module and add the Ferret dependency:

{{< terminal command="true" >}}
mkdir ferret-service && cd ferret-service
go mod init ferret-service
go get github.com/MontFerret/ferret/v2
{{</ terminal >}}

Run a one-shot query to verify the setup:

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
        source.NewAnonymous(`RETURN { message: "Ferret is running" }`),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
}
{{</ code >}}

## Choose a standard library profile

The full standard library includes file system and network access. For a service that runs user-submitted queries, restrict it:

{{< code lang="go" >}}
import "github.com/MontFerret/ferret/v2/pkg/stdlib"
{{</ code >}}

| Profile | What it does | When to use |
| --- | --- | --- |
| `stdlib.Full()` | All function groups | Trusted scripts, internal pipelines |
| `stdlib.Safe()` | Full minus IO (file system + network) | User-submitted queries |
| `stdlib.Only(groups...)` | Only the listed groups | Minimal sandboxes |

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
)
{{</ code >}}

See [Configuration]({{< ref "/docs/embedding/configuration" >}}) for the full list of groups.

## Configure the HTML module

This service needs the HTML module to fetch and query pages. Install it from the `contrib` repository:

{{< terminal command="true" >}}
go get github.com/MontFerret/contrib/modules/web/html
{{</ terminal >}}

Configure the in-process memory driver for static HTML, then pass the module to `ferret.WithModules`:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

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
{{</ code >}}

The `memory` driver fetches and parses static HTML without a browser. Ferret modules are installed, configured, and registered explicitly by the host application. See [Modules]({{< ref "/docs/embedding/modules" >}}) for the complete registration and lifecycle model.

### Add browser support

For JavaScript-rendered pages, add the `cdp` driver alongside the `memory` driver:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/cdp"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

htmlmod, err := html.New(
    html.WithDefaultDriver(memory.New()),
    html.WithDrivers(
        cdp.New(),
    ),
)
if err != nil {
    log.Fatal(err)
}

engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
    ferret.WithModules(htmlmod),
)
{{</ code >}}

Scripts select the driver at query time with `{ driver: "cdp" }`:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
RETURN page.title
{{</ code >}}

Without the `driver` option, the default `memory` driver is used.

## Compile and reuse a plan

When the same query runs multiple times — with different parameters, in different goroutines, or on a schedule — compile it once and reuse the plan:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.New("extract-title", `
    LET page = WEB::HTML::OPEN(@url)
    RETURN page.title
`))
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

urls := []string{
    "https://mockery.ferretlang.org",
    "https://mockery.ferretlang.org/scenarios/ecommerce/",
}

for _, url := range urls {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("url", url),
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
{{</ code >}}

The plan manages an internal pool of virtual machines. Sessions borrow a VM from the pool and return it on close, so reusing a plan avoids recompilation and reduces allocation overhead.

## Build an HTTP handler

The following handler accepts a JSON body with `query` and optional `params`, compiles and runs the query, and returns the result:

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
    "github.com/MontFerret/ferret/v2/pkg/stdlib"

    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

type QueryRequest struct {
    Query  string         `json:"query"`
    Params map[string]any `json:"params,omitempty"`
}

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

    http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req QueryRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request body", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
        defer cancel()

        plan, err := engine.Compile(ctx, source.New("api-query", req.Query))
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
            return
        }
        defer plan.Close()

        sessionOpts := make([]ferret.Option, 0, len(req.Params))
        for k, v := range req.Params {
            sessionOpts = append(sessionOpts, ferret.WithSessionParam(k, v))
        }

        session, err := plan.NewSession(ctx, sessionOpts...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer session.Close()

        output, err := session.Run(ctx)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(output.Content)
    })

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
{{</ code >}}

Test it:

{{< terminal command="true" >}}
curl -s -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"query": "RETURN UPPER(@name)", "params": {"name": "ferret"}}'
{{</ terminal >}}

{{< notification type="info" >}}
For production use, consider caching compiled plans by query hash to avoid recompilation on repeated queries.
{{</ notification >}}

## Handle errors

Compilation and runtime errors are separate concerns:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, source.New("user-script.fql", query))
if err != nil {
    // Compilation error — syntax or semantic issue.
    // The error includes source location (line and column).
    // Safe to return to the user.
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
    return
}
{{</ code >}}

Runtime errors come from `session.Run`:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    // Runtime error — execution failure, context timeout, etc.
    // May be context.DeadlineExceeded if the timeout was hit.
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
{{</ code >}}

Always set a timeout on the context to prevent runaway queries:

{{< code lang="go" >}}
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()
{{</ code >}}

## Control concurrency

Three options control how the engine manages concurrent execution:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
    ferret.WithModules(htmlmod),
    ferret.WithMaxActiveSessions(50),
    ferret.WithMaxVMsPerPlan(8),
    ferret.WithMaxIdleVMsPerPlan(4),
)
{{</ code >}}

| Option | Effect |
| --- | --- |
| `WithMaxActiveSessions` | Engine-wide cap on concurrent sessions. `NewSession` blocks when full. |
| `WithMaxVMsPerPlan` | Per-plan cap on total VMs (idle + active). Bounds resource use for hot queries. |
| `WithMaxIdleVMsPerPlan` | Per-plan idle VM cache size. Trade memory for faster session creation. |

See [Configuration]({{< ref "/docs/embedding/configuration" >}}) for details on how these interact.

## Shut down cleanly

Handle OS signals and close the engine before exiting:

{{< code lang="go" >}}
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/stdlib"

    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
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
        ferret.WithMaxActiveSessions(50),
    )
    if err != nil {
        log.Fatal(err)
    }

    srv := &http.Server{Addr: ":8080"}

    // Register handlers using engine ...

    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("shutting down")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    srv.Shutdown(ctx)
    engine.Close()
}
{{</ code >}}

Close active sessions and plans before closing the engine. `engine.Close()` runs engine-level cleanup hooks in reverse registration order (LIFO). Always close the HTTP server first so no new requests arrive while the engine is shutting down.

## Complete example

A full runnable service combining everything above:

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
    "github.com/MontFerret/ferret/v2/pkg/stdlib"

    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/cdp"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

type QueryRequest struct {
    Query  string         `json:"query"`
    Params map[string]any `json:"params,omitempty"`
}

func main() {
    htmlmod, err := html.New(
        html.WithDefaultDriver(memory.New()),
        html.WithDrivers(
            cdp.New(),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    engine, err := ferret.New(
        ferret.WithStdlib(stdlib.Safe()),
        ferret.WithModules(htmlmod),
        ferret.WithMaxActiveSessions(50),
        ferret.WithMaxVMsPerPlan(8),
        ferret.WithMaxIdleVMsPerPlan(4),
    )
    if err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()

    mux.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req QueryRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request body", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
        defer cancel()

        plan, err := engine.Compile(ctx, source.New("api-query", req.Query))
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
            return
        }
        defer plan.Close()

        sessionOpts := make([]ferret.Option, 0, len(req.Params))
        for k, v := range req.Params {
            sessionOpts = append(sessionOpts, ferret.WithSessionParam(k, v))
        }

        session, err := plan.NewSession(ctx, sessionOpts...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer session.Close()

        output, err := session.Run(ctx)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(output.Content)
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    go func() {
        log.Println("listening on :8080")
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("shutting down")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    srv.Shutdown(ctx)
    engine.Close()
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-modules,embedding-configuration,guide-writing-plugins,guide-precompiled-programs" >}}
