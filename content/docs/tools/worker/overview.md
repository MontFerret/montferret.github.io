---
title: "Overview"
weight: 10
draft: false
description: "Understand what Worker runs, where it fits, and which Ferret capabilities it exposes."
---

# Overview

Worker is an HTTP runtime for Ferret. It receives a FQL script, compiles or reuses a cached plan, runs the script with optional parameters, and writes the serialized Ferret result back to the client.

If your application already owns a configured Ferret engine and needs to expose it to another process, see [Wire]({{< ref "docs/embedding/wire" >}}). Wire is an integration library around that runtime; Worker is a deployable HTTP execution service.

These pages document Worker `v{{< data "versions.worker.v2" >}}`.

## Run a query

Start Worker, then send a JSON request to `POST /`:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"text":"return 1 + 2"}'
{{< /terminal >}}

The response body is the query result:

```json
3
```

Worker does not wrap successful results in a service envelope. If the script returns an object, array, string, number, boolean, or `none`, the response body is that serialized value.

## Pass parameters

Use the `params` object to pass values that the script reads with `@name`:

{{< terminal >}}
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "text": "let doc = web::html::open(@url) return doc.title",
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
| JWT, OAuth2 | Work with tokens and OAuth2 flows. |
| Archive, PDF, XLSX | Work with archives, documents, and spreadsheets. |
| REST, LLM | Make REST requests and interact with configured model providers. |

The Worker executable registers these modules at startup. Network operations still depend on the host’s network policy and any credentials the operation requires.

The HTML module registers two drivers:

| Driver | Behavior |
| --- | --- |
| `memory` | The default driver. Fetches and parses documents without a browser. |
| `cdp` | Connects to Chrome or Chromium over the Chrome DevTools Protocol for browser-backed scripts. |

Use `driver: "cdp"` when a script needs browser rendering or browser interaction:

{{< editor lang="fql" >}}
let page = web::html::open("https://ferretlang.org", { driver: "cdp" })
return page.title
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

Worker accepts source queries over HTTP. The CLI runs compiled artifacts locally. For remote interactive debugging, use a [Wire runtime]({{< ref "/docs/tools/cli/wire" >}}) that exposes the debugger.

## Next steps

{{< docs-related tiles="tools-worker-api,tools-worker-configuration,tools-worker-deployment" >}}
