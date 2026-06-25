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

Runtime selection is per Lab run. The runner, sources, local services, parameters, and reporters still run in the Lab process.

## Run request

For each query, Lab sends a `POST` request to the runtime endpoint with JSON:

```json
{
  "text": "RETURN @name",
  "params": {
    "name": "Ada"
  }
}
```

The runtime response body is treated as the query output. Non-`2xx` HTTP responses fail the test with the response status.

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

Use `runtime-param=path` to override the run endpoint:

{{< terminal >}}
lab run \
  --runtime=https://ferret.example.com \
  --runtime-param='path:"/v1/run"' \
  tests/
{{< /terminal >}}

`path` overrides only query execution. `lab version --runtime=...` uses the runtime URL path and appends `info`, or requests `/info` when the URL has no path.

## Runtime params

Pass runtime-specific settings with `--runtime-param`. Values are parsed as JSON.

| Param | Type | Meaning |
| --- | --- | --- |
| `headers` | object of strings | Extra headers to add to runtime requests. |
| `path` | string | Run endpoint override. |
| `cookies` | object of strings | Cookies to add to runtime requests. |

Example:

{{< terminal >}}
lab run \
  --runtime=https://ferret.example.com \
  --runtime-param='headers:{"Authorization":"Bearer token"}' \
  --runtime-param='cookies:{"session":"abc123"}' \
  --runtime-param='path:"/v1/run"' \
  tests/
{{< /terminal >}}

Invalid runtime param types are rejected before tests run.

## Version request

Use `lab version --runtime=...` to inspect the Ferret version reported by an HTTP runtime.

{{< terminal >}}
lab version --runtime=https://ferret.example.com
{{< /terminal >}}

Lab expects the runtime info response to include `version.ferret`:

```json
{
  "version": {
    "ferret": "2.0.0"
  }
}
```

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
