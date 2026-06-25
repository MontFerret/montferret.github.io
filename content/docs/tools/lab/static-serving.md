---
title: "Static File Server"
weight: 50
draft: false
description: "Serve local static files for Lab tests or standalone fixture access."
---

# Static File Server

Lab can serve local directories over HTTP for tests. Use the static file server for fixture files, generated pages, built frontend assets, or other files that a Ferret script should read through a URL.

During `lab run`, static service URLs are passed to FQL under `@lab.static`.

## Serve static files during a test run

Use `run --serve` to serve a directory while tests run.

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

Use `lab serve --static` when you want a static fixture service without running tests.

{{< terminal >}}
lab serve --static ./dist@app
{{< /terminal >}}

Standalone `serve` entries must be explicit. Positional service entries are rejected; use `--static` for static directories.

You can serve more than one directory:

{{< terminal >}}
lab serve \
  --static ./dist@app \
  --static ./fixtures@fixtures
{{< /terminal >}}

Lab prints the URL for each started service and runs until the process is cancelled.

## Entry syntax

Static entries use this binding syntax:

| Syntax | Meaning |
| --- | --- |
| `<path>` | Serve the directory on a dynamic port with the directory name as the alias. |
| `<path>:<port>` | Serve the directory on a fixed port with the directory name as the alias. |
| `<path>@<alias>` | Serve the directory with an explicit alias. |
| `<path>@<alias>:<port>` | Serve the directory with an explicit alias and fixed port. |

Static entries must point to existing directories. The default alias is the directory name.

Aliases must start with a letter or underscore and may contain letters, numbers, underscores, and hyphens. Use bracket access for aliases that are not valid FQL dotted-property names:

{{< code lang="fql" >}}
RETURN @lab.static["api-fixtures"] + "/users.json"
{{< /code >}}

Duplicate aliases are rejected.

## Fixed ports

Use a fixed port when another process needs a stable URL.

{{< terminal >}}
lab serve --static ./dist@app:8080
{{< /terminal >}}

The fixed port is part of the service binding. The alias is still `app`, so test scripts use `@lab.static.app`.

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
