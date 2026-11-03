---
title: "HTTP Runtime"
weight: 70
draft: false
description: "Run Lab tests against a Ferret-compatible HTTP runtime."
---

# HTTP Runtime

By default, Lab runs tests with its builtin Ferret runtime. Use an HTTP runtime when you want Lab to send each query to a remote Ferret-compatible service.

{{< terminal >}}
lab run --runtime=http://localhost:8080 tests/
{{< /terminal >}}

Lab uses the Universal Ferret API's one-shot `api.Runtime.Run` for HTTP execution, including repetitions and retries. Its case executor remembers that HTTP compilation is unsupported and falls back without issuing a compilation request. HTTP keeps the request protocol described below. Runtime selection is per Lab run. The runner, sources, local services, parameters, and reporters still run in the Lab process.

## Run request

For each query, Lab sends a `POST` request to the runtime endpoint with JSON:

```json
{
  "text": "return @name",
  "params": {
    "name": "Ada"
  }
}
```

A successful request returns a one-shot Universal API output as soon as the response headers arrive. Its metadata records the response content type and the length when known. Go callers can stream the body with `Consume` or use `Collect` to retain the bytes received, including a prefix delivered before an error. Lab checks completion before passing a test. YAML suites decode successful query output as JSON, with an empty response yielding a null result. Non-`2xx` HTTP responses fail the test with the response status.

By default, Lab sends these request headers:

| Header | Value |
| --- | --- |
| `Content-Type` | `application/json` |
| `Accept` | `*/*` |
| `Accept-Charset` | `utf-8` |
| `Accept-Encoding` | `gzip, deflate` |
| `Cache-Control` | `no-cache` |

## Endpoint paths

If the runtime URL includes a path, Lab uses that path for `run` requests:

{{< terminal >}}
lab run --runtime=https://ferret.example.com/v1/run tests/
{{< /terminal >}}

Use `runtime-option=path` to override the run endpoint:

{{< terminal >}}
lab run \
  --runtime=https://ferret.example.com \
  --runtime-option='path:"/v1/run"' \
  tests/
{{< /terminal >}}

`path` overrides only query execution. `lab version --runtime=...` uses the runtime URL path and appends `info`, or requests `/info` when the URL has no path.

## Runtime options

Pass adapter-specific settings with `--runtime-option`. Values are parsed as JSON.

| Option | Type | Meaning |
| --- | --- | --- |
| `headers` | object of strings | Extra headers to add to runtime requests. |
| `path` | string | Run endpoint override. |
| `cookies` | object of strings | Cookies to add to runtime requests. |

Example:

{{< terminal >}}
lab run \
  --runtime=https://ferret.example.com \
  --runtime-option='headers:{"Authorization":"Bearer token"}' \
  --runtime-option='cookies:{"session":"abc123"}' \
  --runtime-option='path:"/v1/run"' \
  tests/
{{< /terminal >}}

Unknown options, settings for another adapter, and invalid types or values are rejected before external setup. FQL values use `--param`, bindings, or YAML query/assertion parameters. Ordinary FQL names `headers` and `path` do not configure the HTTP adapter; adapter options never become FQL inputs. See [Configuration](../configuration/) for quoting and precedence rules.

## Version request

Use `lab version --runtime=...` to inspect the Ferret version reported by an HTTP runtime.

{{< terminal >}}
lab version \
  --runtime=https://ferret.example.com/api \
  --runtime-option='headers:{"Authorization":"Bearer token123"}'
{{< /terminal >}}

Lab expects the runtime info response to include `version.ferret`:

```json
{
  "version": {
    "ferret": "2.0.0"
  }
}
```

The version value is opaque and printed unchanged.

## Capabilities and ownership

The HTTP protocol supports one-shot execution and version reporting. Reusable compilation and debugging, per-session filesystem roots, and output-codec selection are unsupported by this adapter. Go callers receive errors matching `errors.ErrUnsupported` without an HTTP request. A canceled context takes precedence.

The output owns the response body until consumption completes, the caller closes it, or the invocation is canceled. The version request reads and closes its body before returning. Both paths preserve cleanup errors. The adapter borrows the shared HTTP client, so its own close operation leaves that client usable. Filesystem and outbound HTTP execution policies belong to the host; Lab rejects local `--policy-fs-*` and `--policy-http-*` settings for HTTP runtimes.

## Local services with remote runtimes

When tests use `--serve` or `--mock`, Lab starts those services locally and passes URLs through `@lab.static` and `@lab.mock`.

If the HTTP runtime runs outside the Lab host, `127.0.0.1` may point to the runtime host instead of the Lab host. Use `--serve-host` to advertise a host the runtime can reach, and `--serve-bind` when Lab must listen on a non-loopback interface.

For local service entry syntax, see [Static File Server]({{< ref "static-serving" >}}) and [Mock API Server]({{< ref "mock-api" >}}).

{{< terminal >}}
lab run \
  --runtime=https://ferret.example.com \
  --serve ./dist@app \
  --mock ./users.yaml@api \
  --serve-bind 0.0.0.0 \
  --serve-host host.docker.internal \
  tests/
{{< /terminal >}}

See [Migrate Lab runtime configuration]({{< ref "/docs/migrations/lab-runtime-configuration" >}}) when upgrading an earlier v2 configuration.

## Next steps

{{< docs-related tiles="tools-lab-wire-runtime,tools-lab-running-tests,tools-worker" >}}
