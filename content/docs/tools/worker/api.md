---
title: "API"
weight: 40
draft: false
description: "Use the Worker HTTP API for query execution, health checks, version info, and remote runtime clients."
---

# API

Worker exposes a small HTTP API:

| Endpoint | Purpose |
| --- | --- |
| `POST /` | Execute a FQL script. |
| `GET /health` | Check Worker and dependency health. |
| `GET /info` | Return public IP and version metadata. |

All successful query responses are the serialized Ferret result body. Worker does not return a `{ "data": ... }` wrapper.

## `POST /`

`POST /` executes a FQL script.

Request body:

```json
{
  "text": "RETURN @name",
  "params": {
    "name": "Ada"
  }
}
```

Fields:

| Field | Required | Meaning |
| --- | --- | --- |
| `text` | Yes | FQL script source. |
| `params` | No | Named values available to the script as `@name`. |

Example:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "text": "RETURN { name: @name, active: true }",
    "params": {
      "name": "Ada"
    }
  }'
{{< /terminal >}}

Example shape:

```json
{
  "active": true,
  "name": "Ada"
}
```

Status codes:

| Status | Meaning |
| --- | --- |
| `200` | The query ran and the body is the serialized result. |
| `400` | The request body could not be parsed, the query could not compile, or execution failed. |

## Query errors

Error responses use this shape:

```json
{
  "error": "run program: missing parameter",
  "details": "UnresolvedSymbol: missing parameter\n --> anonymous:1:8\n  |\n1 | RETURN @missing\n  |        ^^^^^^^^ parameter '@missing' was not provided\n"
}
```

`error` is the plain wrapped error. `details` contains formatted Ferret diagnostics when the error supports diagnostics; otherwise it repeats the plain error text.

## Browser-backed query

Use the CDP driver from FQL when the script needs Chrome:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "text": "LET page = DOCUMENT(@url, { driver: \"cdp\" }) RETURN page.title",
    "params": {
      "url": "https://example.com"
    }
  }'
{{< /terminal >}}

The Worker process must be configured with a reachable Chrome DevTools endpoint. See [Configuration]({{< ref "configuration" >}}).

## `GET /health`

`GET /health` checks whether Worker is running and, when Chrome is enabled, whether Worker can reach Chrome's version endpoint.

{{< terminal >}}
curl -i http://localhost:8080/health
{{< /terminal >}}

Responses:

| Status | Body | Meaning |
| --- | --- | --- |
| `200` | Empty | Worker is healthy. |
| `424` | Empty | Worker cannot reach a required dependency such as Chrome. |

When Worker starts with `-no-chrome`, `/health` returns `200` without checking Chrome.

## `GET /info`

`GET /info` returns public IP and version metadata.

{{< terminal >}}
curl http://localhost:8080/info
{{< /terminal >}}

Response:

```json
{
  "ip": "203.0.113.10",
  "version": {
    "worker": "2.0.0-rc.16",
    "chrome": {
      "browser": "Chrome/148.0.7778.178",
      "protocol": "1.3",
      "v8": "14.8.x",
      "webkit": "537.36"
    },
    "ferret": "2.0.0-alpha.24"
  }
}
```

Fields:

| Field | Meaning |
| --- | --- |
| `ip` | Public IP detected through `http://checkip.amazonaws.com`. |
| `version.worker` | Worker version embedded at build time. |
| `version.ferret` | Ferret runtime version embedded at build time. |
| `version.chrome` | Chrome version details from the Chrome DevTools `/json/version` endpoint. |

If Worker is built without release linker flags, `version.worker` and `version.ferret` may be empty. If Worker runs with `-no-chrome`, Chrome version fields are empty.

`/info` can return `424` when Worker cannot retrieve Chrome version data or public IP data.

## Remote runtime clients

The Ferret CLI and Lab use Worker-compatible HTTP runtimes by sending the same `POST /` payload:

```json
{
  "text": "RETURN @name",
  "params": {
    "name": "Ada"
  }
}
```

Use Worker from the CLI:

{{< terminal >}}
ferret run --runtime http://localhost:8080 script.fql
{{< /terminal >}}

Use Worker from Lab:

{{< terminal >}}
lab run --runtime=http://localhost:8080 tests/
{{< /terminal >}}

Remote runtime clients expect `/info` to expose `version.ferret` when they ask for the runtime version.

## Next steps

{{< docs-related tiles="tools-worker-configuration,tools-worker-deployment,tools-cli" >}}
