---
title: "Mock API Server"
weight: 55
draft: false
description: "Serve OpenAPI-described mock APIs for Lab tests or standalone fixture access."
---

# Mock API Server

Lab can serve OpenAPI-described mock APIs over HTTP. Use mock APIs when a test needs stable REST responses without depending on an external service.

During `lab run`, mock service URLs are passed to FQL under `@lab.mock`.

## Serve a mock API during a test run

Use `run --mock` to serve an OpenAPI-compatible mock API while tests run.

{{< terminal >}}
lab run --mock ./users.yaml@api tests/
{{< /terminal >}}

The service URL is available as `@lab.mock.api`.

{{< code lang="fql" >}}
LET response = IO::NET::HTTP::GET(@lab.mock.api + "/users/123")
LET user = JSON_PARSE(TO_STRING(response))

RETURN user.id == "123"
{{< /code >}}

Multiple mock APIs can be served in the same run:

{{< terminal >}}
lab run \
  --mock ./users.yaml@users \
  --mock ./billing.yaml@billing \
  tests/
{{< /terminal >}}

## Start standalone mock APIs

Use `lab serve --mock` when you want a mock API service without running tests.

{{< terminal >}}
lab serve --mock ./users.yaml@api
{{< /terminal >}}

Standalone `serve` entries must be explicit. Positional service entries are rejected; use `--mock` for mock API specs.

Lab prints the URL for each started service and runs until the process is cancelled.

## Entry syntax

Mock entries use this binding syntax:

| Syntax | Meaning |
| --- | --- |
| `<path>` | Serve the spec on a dynamic port with the filename, without extension, as the alias. |
| `<path>:<port>` | Serve the spec on a fixed port with the filename, without extension, as the alias. |
| `<path>@<alias>` | Serve the spec with an explicit alias. |
| `<path>@<alias>:<port>` | Serve the spec with an explicit alias and fixed port. |

Mock entries must point to existing files. The default alias is the spec filename without its extension.

Aliases must start with a letter or underscore and may contain letters, numbers, underscores, and hyphens. Use bracket access for aliases that are not valid FQL dotted-property names:

{{< code lang="fql" >}}
RETURN @lab.mock["users-api"] + "/users/123"
{{< /code >}}

Duplicate aliases are rejected.

## Fixed ports

Use a fixed port when another process needs a stable URL.

{{< terminal >}}
lab serve --mock ./users.yaml@api:8081
{{< /terminal >}}

The fixed port is part of the service binding. The alias is still `api`, so test scripts use `@lab.mock.api`.

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

## Bind and advertised hosts

By default, mock services bind to `127.0.0.1` and advertise URLs with `127.0.0.1`.

Use `--serve-bind` to choose the listener host. Use `--serve-host` to choose the host placed in generated URLs.

{{< terminal >}}
lab run \
  --mock ./users.yaml@api \
  --serve-bind 0.0.0.0 \
  --serve-host host.docker.internal \
  tests/
{{< /terminal >}}

Both values are hosts only. Do not include a port.

When `--serve-host` is set and `--serve-bind` is omitted, Lab binds to a wildcard host so remote runtimes can reach the service: `0.0.0.0` for IPv4 and hostnames, or `::` for IPv6 literals.
