---
title: "Embed Ferret in a Go application"
sidebarTitle: "Embedding the runtime"
weight: 110
draft: false
description: "Initialize the Ferret engine, compile queries, run them with parameters, and handle errors in a Go service."
aliases:
    - /docs/guides/embedding-runtime/
---

# Embed Ferret in a Go application

Embedding Ferret lets you run FQL queries inside your Go application — useful for user-defined extraction rules, dynamic data pipelines, or sandboxed scripting. This guide walks through building a Go service that accepts queries over HTTP.

For the full native API reference, see [Go embedding]({{< ref "/docs/embedding/go" >}}).

{{% notification type="warning" %}}
FQL source is code. Do not insert URLs, selectors, HTML, credentials, request bodies, or other runtime values with `fmt.Sprintf` or string concatenation. Pass them as parameters instead. See [Construct FQL safely]({{< ref "/docs/embedding/go/safe-fql-construction" >}}).
{{% /notification %}}

## Set up the project

Create a new Go module:

{{< terminal command="true" >}}
mkdir ferret-service && cd ferret-service
go mod init ferret-service
{{</ terminal >}}

Follow the [Go installation instructions]({{< ref "/docs/embedding/go/getting-started#installation" >}}) to use the consumable-output API in this guide. Then run a one-shot query to verify the setup:

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
        ferret.NewAnonymousSource(`return { message: "Ferret is running" }`),
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

    return nil
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

See [Runtime Configuration]({{< ref "/docs/embedding/go/configuration" >}}) for the full list of groups.

## Configure the HTML module

This service needs the HTML module to fetch and query pages. Install it from the `contrib` repository:

{{< terminal command="true" >}}
go get github.com/MontFerret/contrib/modules/web/html@v1.0.0-rc.29
{{</ terminal >}}

Configure the in-process memory driver for static HTML, then pass the module to `ferret.WithModules`:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

htmlmod := html.New(
    html.WithDefaultDriver(memory.New()),
)

engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
    ferret.WithModules(htmlmod),
)
{{</ code >}}

The `memory` driver fetches and parses static HTML without a browser. Ferret modules are installed, configured, and registered explicitly by the host application. See [Modules]({{< ref "/docs/embedding/go/modules" >}}) for the complete registration and lifecycle model.

### Add browser support

For JavaScript-rendered pages, add the `cdp` driver alongside the `memory` driver:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/cdp"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

htmlmod := html.New(
    html.WithDefaultDriver(memory.New()),
    html.WithDrivers(
        cdp.New(),
    ),
)

engine, err := ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
    ferret.WithModules(htmlmod),
)
{{</ code >}}

Scripts select the driver at query time with `{ driver: "cdp" }`:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
return page.title
{{</ code >}}

Without the `driver` option, the default `memory` driver is used.

## Format source without executing it

Use `pkg/formatter` to format FQL independently of an engine. Formatting checks syntax without resolving host functions or executing the program, so scripts can reference functions unavailable to the formatting application.

{{< code lang="go" >}}
import (
    "bytes"

    "github.com/MontFerret/ferret/v2/pkg/formatter"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func formatSource(text string) (string, error) {
    f, err := formatter.New()
    if err != nil {
        return "", err
    }

    var out bytes.Buffer
    if err := f.Format(&out, source.New("query.fql", text)); err != nil {
        return "", err
    }

    return out.String(), nil
}
{{</ code >}}

`formatter.WithSingleQuote(true)` prefers single quotes for ordinary string literals; the default prefers double quotes. Quote conversion preserves the exact string value. If a conversion cannot safely preserve that value, the formatter keeps the original spelling. For example, the default formatter retains `'a"b'`. Literal bodies, including backslashes and embedded line breaks, are preserved, and template chunks keep their original spelling.

`Format` writes directly to its writer. Any failure can leave partial output; use the output only after a nil error. The example buffers at the caller so it can discard incomplete output. A formatter can be reused after a failure and shared across goroutines with separate writers.

Native diagnostic spans are zero-based, half-open UTF-8 byte offsets into the original source. Native positions use one-based lines and byte columns. Convert these coordinates to an editor's coordinate system in the integration layer.

Recovered panics return a structured diagnostic with kind `diagnostics.UnexpectedError`, available through `errors.As`. Ordinary writer errors retain their identity for `errors.Is`. Internal diagnostics can contain panic and stack details; service adapters should classify them and return a generic internal-error response instead of exposing those details.

## Compile and reuse a plan

When the same query runs multiple times — with different parameters, in different goroutines, or on a schedule — compile it once and reuse the plan:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx, ferret.NewSource(
    "extract-title.fql",
    `let page = web::html::open(@url)
return page.title`,
))
if err != nil {
    return err
}
defer plan.Close()

urls := []string{
    "https://mockery.ferretlang.org",
    "https://mockery.ferretlang.org/scenarios/ecommerce/",
}

// This snippet also imports "errors".
for _, url := range urls {
    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("url", url),
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
{{</ code >}}

The plan manages an internal pool of virtual machines. Sessions borrow a VM from the pool and return it on close, so reusing a plan avoids recompilation and reduces allocation overhead.

## Build an HTTP handler

Collect before writing response headers or bytes so terminal failures can still produce an error status. For forwarding to a writer without HTTP status handling, see [Output]({{< ref "/docs/embedding/go/output#consume-borrowed-bytes" >}}).

The following handler accepts a JSON body with `query` and optional `params`, compiles and runs the query, and returns the result:

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/stdlib"

    "github.com/MontFerret/contrib/modules/web/html"
    "github.com/MontFerret/contrib/modules/web/html/drivers/memory"
)

type QueryRequest struct {
    Query  string         `json:"query"`
    Params map[string]any `json:"params,omitempty"`
}

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    htmlmod := html.New(
        html.WithDefaultDriver(memory.New()),
    )

    engine, err := ferret.New(
        ferret.WithStdlib(stdlib.Safe()),
        ferret.WithModules(htmlmod),
    )
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, engine.Close()) }()

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

        plan, err := engine.Compile(ctx,
            ferret.NewSource("api-query.fql", req.Query),
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
            return
        }
        defer plan.Close()

        session, err := plan.NewSession(ctx,
            ferret.WithSessionParams(req.Params),
        )
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

        defer output.Close()
        content, err := output.Collect(ctx)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if content == nil {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        w.Header().Set("Content-Type", content.Metadata.ContentType)
        if _, err := w.Write(content.Data); err != nil {
            log.Printf("write response: %v", err)
        }
    })

    log.Println("listening on :8080")
    return http.ListenAndServe(":8080", nil)
}
{{</ code >}}

Test it:

{{< terminal command="true" >}}
curl -s -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"query": "return upper(@name)", "params": {"name": "ferret"}}'
{{</ terminal >}}

{{< notification type="info" >}}
For production use, consider caching compiled plans by query hash to avoid recompilation on repeated queries.
{{</ notification >}}

## Handle errors

Compilation and runtime errors are separate concerns:

{{< code lang="go" >}}
plan, err := engine.Compile(ctx,
    ferret.NewSource("user-script.fql", query),
)
if err != nil {
    // Compilation error — syntax or semantic issue.
    // The error includes source location (line and column).
    // Safe to return to the user.
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
    return
}
{{</ code >}}

Handle preparation errors from `session.Run`, then runtime and encoding errors from consumption:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    // Preparation failed, possibly before VM entry due to cancellation.
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
defer output.Close()

content, err := output.Collect(ctx)
if err != nil {
    // Execution, encoding, hooks, cancellation, or cleanup failed.
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
_ = content // write only after the terminal error check
{{</ code >}}

Runtime error identities remain available through diagnostic wrapping. Use the public runtime sentinels with `errors.Is` when the host needs to distinguish zero-divisor failures:

{{< code lang="go" >}}
switch {
case errors.Is(err, runtime.ErrDivisionByZero):
    // Int, Float, or Duration division by zero
case errors.Is(err, runtime.ErrModuloByZero):
    // Int or Float modulo by zero
case errors.Is(err, runtime.ErrInvalidOperation):
    // another unsupported runtime operation
}
{{</ code >}}

`runtime.ErrDivisionByZero` and `runtime.ErrModuloByZero` are the canonical identities. The existing `vm.ErrDivisionByZero` and `vm.ErrModuloByZero` names remain aliases for source compatibility, but embedding code should prefer the runtime package.

Always set a timeout on the context to prevent runaway queries:

{{< code lang="go" >}}
ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
defer cancel()
{{</ code >}}

The VM checks cancellation at structural execution boundaries. Native work between those boundaries completes atomically. Any blocking or remote host operation receives the execution context and is responsible for observing cancellation while it retains control. A returned `context.Canceled` or `context.DeadlineExceeded`, including a wrapped error, bypasses FQL recovery. Cancellation before VM admission can fail `Session.Run`; after admission it is reported by output consumption, with recorded failures preserved.

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

See [Runtime Configuration]({{< ref "/docs/embedding/go/configuration" >}}) for details on how these interact.

## Shut down cleanly

Stop accepting new requests and drain active handlers before closing the engine. Each handler owns its output, session, and plan. The complete example gives each query a deadline and waits for HTTP shutdown before releasing the engine.

## Complete example

A full runnable service combining everything above:

{{< code lang="go" >}}
package main

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/MontFerret/ferret/v2"
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
    if err := serve(); err != nil {
        log.Fatal(err)
    }
}

func serve() (err error) {
    htmlmod := html.New(
        html.WithDefaultDriver(memory.New()),
        html.WithDrivers(
            cdp.New(),
        ),
    )

    engine, err := ferret.New(
        ferret.WithStdlib(stdlib.Safe()),
        ferret.WithModules(htmlmod),
        ferret.WithMaxActiveSessions(50),
        ferret.WithMaxVMsPerPlan(8),
        ferret.WithMaxIdleVMsPerPlan(4),
    )
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, engine.Close()) }()

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

        plan, err := engine.Compile(ctx,
            ferret.NewSource("api-query.fql", req.Query),
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnprocessableEntity)
            return
        }
        defer plan.Close()

        session, err := plan.NewSession(ctx,
            ferret.WithSessionParams(req.Params),
        )
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

        defer output.Close()
        content, err := output.Collect(ctx)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if content == nil {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        w.Header().Set("Content-Type", content.Metadata.ContentType)
        if _, err := w.Write(content.Data); err != nil {
            log.Printf("write response: %v", err)
        }
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    shutdown, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    served := make(chan error, 1)
    go func() { served <- srv.ListenAndServe() }()

    select {
    case <-shutdown.Done():
    case err = <-served:
        if errors.Is(err, http.ErrServerClosed) {
            err = nil
        }
    }
    return errors.Join(err, srv.Shutdown(context.Background()))
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-go-modules,embedding-go-configuration,guide-writing-plugins,guide-precompiled-programs" >}}
