---
title: "CDP"
weight: 90
draft: false
description: "Use browser-backed FQL in Lab without treating CDP as a Lab-owned flag."
---

# CDP

Lab can run tests that use browser-backed FQL, but Lab does not own browser execution or Chrome DevTools Protocol behavior. That belongs to the selected Ferret runtime.

In current Lab v2, there is no `--cdp` flag and no `LAB_CDP` environment variable. Configure browser access through the Ferret runtime and FQL options you are using.

## Browser-backed tests

A browser-backed test is still a normal Lab test file:

{{< code lang="fql" >}}
LET doc = DOCUMENT("https://example.com", {
  driver: "cdp"
})

RETURN doc.title
{{< /code >}}

Run it with Lab:

{{< terminal >}}
lab run browser.fql
{{< /terminal >}}

If the runtime expects an existing browser endpoint, start that browser separately and make sure the runtime can reach it.

## Wait for a browser endpoint

Use `--wait` when a browser endpoint is started outside Lab and may not be ready yet.

{{< terminal >}}
lab run browser-tests/ \
  --wait=http://127.0.0.1:9222/json/version \
  --wait-timeout=5 \
  --wait-attempts=12
{{< /terminal >}}

The wait happens before tests execute. It does not configure the runtime; it only prevents tests from starting before the dependency responds.

## Builtin runtime

With the builtin runtime, Ferret runs inside the Lab process. Browser endpoints such as `127.0.0.1:9222` are resolved from the same host where Lab runs.

Use FQL and Ferret runtime configuration for browser driver details. Lab passes user parameters and `@lab` system parameters into the runtime, but it does not rewrite browser options.

## HTTP and binary runtimes

When you use `--runtime=http(s)://...` or `--runtime=bin:...`, browser connectivity is evaluated from that runtime's environment.

For example, if Lab runs on your laptop but the HTTP runtime runs in a container, `127.0.0.1:9222` inside the query refers to the runtime container, not your laptop. Expose or advertise the browser endpoint in a way that the runtime can reach.

The same rule applies to Lab local services. If a remote runtime must fetch `@lab.static` or `@lab.mock` URLs, set `--serve-bind` and `--serve-host` so those URLs are reachable from the runtime.

## CI pattern

A common CI shape is:

1. Start the browser or runtime service.
2. Use `lab run --wait=...` to wait for its health endpoint.
3. Run tests with `--reporter=simple`.

{{< terminal >}}
lab run browser-tests/ \
  --wait=http://127.0.0.1:9222/json/version \
  --reporter=simple
{{< /terminal >}}
