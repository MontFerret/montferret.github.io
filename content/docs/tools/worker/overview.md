---
title: "Overview"
weight: 10
draft: false
description: "Understand what Worker runs, where it fits, and which Ferret capabilities it exposes."
---

# Overview

Worker is an HTTP runtime for Ferret. It receives a FQL script, compiles or reuses a cached plan, runs the script with optional parameters, and writes the serialized Ferret result back to the client.

The current Worker source documented here is `main` at `v2.0.0-rc.16`.

## Run a query

Start Worker, then send a JSON request to `POST /`:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"text":"RETURN 1 + 2"}'
{{< /terminal >}}

The response body is the query result:

```json
3
```

Worker does not wrap successful results in a service envelope. If the script returns an object, array, string, number, boolean, or `NONE`, the response body is that serialized value.

## Pass parameters

Use the `params` object to pass values that the script reads with `@name`:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "text": "LET doc = DOCUMENT(@url) RETURN doc.title",
    "params": {
      "url": "https://example.com"
    }
  }'
{{< /terminal >}}

Missing parameters are reported as query execution errors with diagnostic details.

## Runtime capabilities

Worker creates a Ferret engine with the standard Ferret runtime plus these contributed modules:

| Module area | Purpose |
| --- | --- |
| CSV, TOML, XML, YAML | Parse and work with common data formats. |
| Web HTML | Open, parse, query, and interact with HTML documents. |
| Web article, robots, sitemap | Extract article content and inspect web metadata. |
| SQLite | Use the SQLite module in memory-only mode. |
| JWT | Work with JSON Web Tokens. |

The HTML module registers two drivers:

| Driver | Behavior |
| --- | --- |
| `memory` | The default driver. Fetches and parses documents without a browser. |
| `cdp` | Connects to Chrome or Chromium over the Chrome DevTools Protocol for browser-backed scripts. |

Use `driver: "cdp"` when a script needs browser rendering or browser interaction:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://ferretlang.org", { driver: "cdp" })
RETURN page.title
{{< /editor >}}

## Runtime boundaries

Worker owns service concerns around HTTP, rate limiting, request body limits, query-plan caching, filesystem root selection, and Chrome connection settings.

The Ferret runtime still owns FQL parsing, compilation, module behavior, parameter semantics, browser-driver behavior, timeout behavior inside scripts, and result serialization.

Worker does not provide authentication, authorization, durable job queues, persistence, scheduling, or multi-tenant isolation by itself. Add those concerns at the deployment layer, usually with a reverse proxy, service mesh, job runner, or application-specific wrapper.

## Use as a remote runtime

The Ferret CLI and Lab can call a Worker-compatible HTTP runtime. For example:

{{< terminal >}}
ferret run --runtime http://localhost:8080 script.fql
lab run --runtime=http://localhost:8080 tests/
{{< /terminal >}}

Remote runtimes execute source queries. Compiled artifacts and interactive debugging require the builtin local runtime.
