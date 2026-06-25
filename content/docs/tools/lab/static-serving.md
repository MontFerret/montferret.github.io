---
title: "Static Serving"
weight: 50
draft: false
description: "Serve local static files and mock APIs for Lab tests."
---

# Static Serving

Lab can start local HTTP services for tests. Use static services for fixture files and built frontend assets. Use mock API services for OpenAPI-described REST endpoints.

During `lab run`, service URLs are passed to FQL under `@lab.static` and `@lab.mock`.

## Serve static files during a test run

Use `run --serve` to serve a directory while tests run:

{{< terminal >}}
lab run --serve ./dist@app tests/
{{< /terminal >}}

The service URL is available as `@lab.static.app`.

{{< code lang="fql" >}}
LET doc = DOCUMENT(@lab.static.app + "/index.html")
RETURN doc.title
{{< /code >}}

Multiple directories can be served in the same run:

{{< terminal >}}
lab run \
  --serve ./dist@app \
  --serve ./fixtures@fixtures \
  tests/
{{< /terminal >}}

## Start standalone services

Use `lab serve` when you want fixture services without running tests.

{{< terminal >}}
lab serve --static ./dist@app
{{< /terminal >}}

Standalone `serve` entries must be explicit. Positional service entries are rejected; use `--static` for directories and `--mock` for mock API specs.

You can serve static directories and mock APIs together:

{{< terminal >}}
lab serve --static ./dist@app --mock ./users.yaml@api
{{< /terminal >}}

Lab prints the URL for each started service and runs until the process is cancelled.

## Entry syntax

Static and mock entries use the same binding syntax:

| Syntax | Meaning |
| --- | --- |
| `<path>` | Serve the path on a dynamic port with a default alias. |
| `<path>:<port>` | Serve the path on a fixed port. |
| `<path>@<alias>` | Serve the path with an explicit alias. |
| `<path>@<alias>:<port>` | Serve the path with an explicit alias and fixed port. |

For static directories, the default alias is the directory name. For mock API specs, the default alias is the spec filename without its extension.

Aliases must start with a letter or underscore and may contain letters, numbers, underscores, and hyphens. Use bracket access for aliases that are not valid FQL dotted-property names:

{{< code lang="fql" >}}
RETURN @lab.static["api-fixtures"] + "/users.json"
{{< /code >}}

Duplicate aliases are rejected.

## Bind and advertised hosts

By default, local services bind to `127.0.0.1` and advertise URLs with `127.0.0.1`.

Use `--serve-bind` to choose the listener host. Use `--serve-host` to choose the host placed in generated URLs.

{{< terminal >}}
lab run \
  --serve ./dist@app \
  --serve-bind 0.0.0.0 \
  --serve-host host.docker.internal \
  tests/
{{< /terminal >}}

Both values are hosts only. Do not include a port.

When `--serve-host` is set and `--serve-bind` is omitted, Lab binds to a wildcard host so remote runtimes can reach the service: `0.0.0.0` for IPv4 and hostnames, or `::` for IPv6 literals.

## Mock API services

Use `run --mock` to serve an OpenAPI-compatible mock API during a test run:

{{< terminal >}}
lab run --mock ./users.yaml@api tests/
{{< /terminal >}}

The service URL is available as `@lab.mock.api`.

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::GET(@lab.mock.api + "/users/123")
LET user = JSON_PARSE(TO_STRING(response))

RETURN user.id == "123"
{{< /code >}}

Use `serve --mock` to start the same mock API without running tests:

{{< terminal >}}
lab serve --mock ./users.yaml@api
{{< /terminal >}}

## Mock API specs

Mock specs use OpenAPI `paths` with operation-level `x-lab-mock` blocks.

```yaml
openapi: 3.1.0
info:
  title: Users API
  version: 1.0.0
paths:
  /users/{id}:
    get:
      x-lab-mock:
        status: 200
        headers:
          X-Test-Server: lab
        body:
          id: "{{ .Path.id }}"
          name: "User {{ .Path.id }}"
```

The mock server currently handles `GET`, `POST`, `PUT`, `PATCH`, and `DELETE` operations. Paths must start with `/`.

`x-lab-mock` supports:

| Field | Meaning |
| --- | --- |
| `status` | HTTP status code. Defaults to `200`. |
| `headers` | Response headers as string values. |
| `body` | Structured JSON response body. String values are rendered as templates. |
| `bodyTemplate` | Raw text response template. Mutually exclusive with `body`. |

When `body` is used, Lab writes JSON and defaults `Content-Type` to `application/json` unless the mock sets it. When `bodyTemplate` is used, Lab writes text and defaults `Content-Type` to `text/plain; charset=utf-8`.

## Template context

Mock response templates use Go template syntax. They receive this context:

| Field | Meaning |
| --- | --- |
| `.Method` | Request method. |
| `.Path` | Path parameters from routes such as `/users/{id}`. |
| `.Query` | Query parameters. |
| `.Headers` | Request headers. |
| `.Body` | Parsed JSON request body, or `nil`. |

Example:

```yaml
paths:
  /echo/{name}:
    post:
      x-lab-mock:
        status: 201
        body:
          method: "{{ .Method }}"
          name: "{{ .Path.name }}"
          active: "{{ .Body.active }}"
```

Request bodies used by templates must be JSON. If the request body is not valid JSON, the mock server returns `400 Bad Request`.
