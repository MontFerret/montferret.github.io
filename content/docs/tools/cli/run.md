---
title: "Run"
weight: 20
draft: false
description: "Execute FQL scripts from files, inline expressions, stdin, or compiled artifacts."
---

# Run

The `run` command executes a FQL script and prints the encoded result to stdout. It is also available as `ferret exec`. Builtin execution is the default; use [Wire]({{< ref "wire" >}}) to execute source against an application-configured host.

Errors and source diagnostics go to stderr. Available output is preserved even when execution also returns an error; successfully produced empty output prints nothing.

## Run a script file

{{< terminal >}}
ferret run script.fql
{{< /terminal >}}

## Evaluate an inline expression

Use `--eval` (or `-e`) to run a FQL expression directly:

{{< terminal >}}
ferret run -e 'return 1 + 2'
{{< /terminal >}}

The `--eval` flag cannot be combined with a file argument.

## Read from stdin

Pipe a script through stdin:

{{< terminal >}}
cat script.fql | ferret run
{{< /terminal >}}

## Run a compiled artifact

Compiled `.fqlc` artifacts produced by [`ferret build`](../build/) can be executed directly:

{{< terminal >}}
ferret run script.fqlc
{{< /terminal >}}

The CLI auto-detects whether the input is FQL source or a compiled artifact. Artifacts require the builtin runtime and cannot be used with a remote runtime.

## Pass parameters

Use `--param` (or `-p`) to pass named values into a script. Parameters are accessible as `@name` in FQL:

{{< terminal >}}
ferret run script.fql --param url=https://example.com --param limit=50
{{< /terminal >}}

{{< code lang="fql" >}}
LET page = web::html::open(@url)

RETURN FOR item IN page[~ css`.result`]
    LIMIT @limit
    RETURN item.textContent
{{< /code >}}

Values are parsed as JSON when possible, otherwise treated as strings:

| Example | Parsed as |
| --- | --- |
| `--param name=Steve` | String `"Steve"` |
| `--param age=42` | Number `42` |
| `--param active=true` | Boolean `true` |
| `--param value=null` | Null (`none` in FQL) |
| `--param tags='["admin","editor"]'` | Array `["admin", "editor"]` |
| `--param user='{"name":"Ada"}'` | Object `{"name": "Ada"}` |
| `--param code='"123"'` | String `"123"` (quoted to prevent number parsing) |

The `--param` flag can be repeated to pass multiple parameters.

## Runtime and browser flags

These flags control where and how scripts execute. They are shared by the `run`, [`repl`](../repl/), and [`debug`](../debug/) commands; debugging requires builtin mode. The `version` command accepts only the runtime selection flags from this table.

| Flag | Default | Description |
| --- | --- | --- |
| `-r`, `--runtime` | `builtin` | `builtin`, `wire`, or a legacy Worker-compatible HTTP URL |
| `--runtime-endpoint` | | Wire endpoint: `tcp://127.0.0.1:<port>` |
| `--runtime-connect-timeout` | `5s` | Positive Wire connection and handshake timeout |
| `--policy-fs-root` | Current working directory | Builtin filesystem sandbox root |
| `--policy-fs-read-only` | `false` | Make the builtin filesystem sandbox read-only |
| `--proxy` | | Proxy server address |
| `--user-agent` | | Custom User-Agent header |
| `-d`, `--browser-address` | `http://127.0.0.1:9222` | Browser remote debugging address |
| `-B`, `--browser-open` | `false` | Open a visible browser for script execution |
| `-b`, `--browser-headless` | `false` | Open a headless browser for script execution |
| `-c`, `--browser-cookies` | `false` | Persist cookies between queries |

Builtin filesystem and HTTP execution policies are shared by `run`, `repl`, and `debug`. See [Filesystem policy]({{< ref "configuration" >}}#filesystem-policy) and [HTTP policy]({{< ref "configuration" >}}#http-policy) for their flags, config keys, environment variables, and defaults. They do not configure remote hosts.

### Browser modes

Scripts that do not use browser-backed features run without a browser. The HTML module defaults to its memory driver; [select `{ driver: "cdp" }` in the script]({{< ref "browser" >}}#select-cdp-in-the-script) when it needs Chrome or Chromium. Browser flags do not change the selected HTML driver.

For a CDP script, choose one of:

**Automatic browser** — the CLI starts and stops a browser for each run:

{{< terminal >}}
ferret run --browser-headless script.fql
{{< /terminal >}}

**Existing browser** — connect to a browser that is already running:

{{< terminal >}}
ferret run --browser-address http://127.0.0.1:9222 script.fql
{{< /terminal >}}

**Managed browser** — start a long-lived browser with [`ferret browser open`](../browser/) and reuse it across runs.

### Remote runtime

Send source to an application-configured Wire host:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 script.fql
{{< /terminal >}}

The host compiles and executes the query using its functions, modules, policies, and filesystem. Explicit builtin policies and browser settings are rejected in Wire mode, including configured or environmental values. See the [Wire runtime guide]({{< ref "wire" >}}) for the development host, endpoint validation, ownership, and diagnostics.

### Legacy HTTP runtime

To execute against a Worker-compatible HTTP runtime:

{{< terminal >}}
ferret run --runtime http://localhost:8080 script.fql
{{< /terminal >}}

The Worker adapter sends source text and parameters with `POST /` and reads the hosted Ferret version with `GET /info`. Execution responses stream unchanged regardless of HTTP status: an error status alone does not fail the command, and the body is not interpreted as structured query diagnostics. Request, body-read, writer, cancellation, and cleanup failures still fail the command.

Remote runtimes cannot run compiled artifacts or CLI debug sessions. Worker does not expose reusable Universal API compilation, output codec selection, or per-session filesystem-root selection.

All runtime and browser flags can also be set persistently through the [`config`](../configuration/) command or environment variables.

## Output delivery

The CLI consumes encoded output directly, writing each chunk completely before accepting the next. Slow stdout applies backpressure. Bytes are not decoded, re-encoded, or collected into a complete response by the command.

Available bytes remain on stdout if execution or delivery later fails. A successful empty result prints nothing. Terminal execution, read, write, cancellation, and cleanup failures cause a failing exit. HTTP response delivery uses a bounded buffer. Streaming encoded output does not imply incremental FQL evaluation.

## Go execution helpers

The CLI's Go helpers `runtime.Run`, `RunSource`, `RunArtifact`, `Builtin.RunArtifact`, and `run.Execute` return `(api.Output, error)`. A returned handle is one-shot: `Consume` processes borrowed, read-only chunks, while `Collect` returns detached `Content.Data`. Copy callback bytes you retain. Non-nil collected content may accompany an error and contain partial output.

`runtime.Run`, `RunArtifact`, and `run.Execute` transfer their invocation-owned runtime and transport resources to the output. `RunSource` borrows its supplied `Resources`. `Builtin.RunArtifact` borrows its receiver and transfers its temporary session and plan to the output. Consume or close output before releasing borrowed parents; callers remain responsible for closing those parents.

Keep the invocation context alive until consumption or closure finishes. `Consume` and `Collect` observe terminal execution, delivery, and cleanup failures. `Close` abandons unread output without certifying success. See [Output]({{< ref "/docs/embedding/go/output" >}}) for the shared lifecycle contract.

## Next steps

{{< docs-related tiles="tools-cli-wire,tools-cli-debug,tools-cli-check,tools-lab" >}}
