---
title: "Configuration"
weight: 20
draft: false
description: "Configure Worker flags, environment variables, limits, filesystem access, cache size, and Chrome."
---

# Configuration

Worker is configured with command-line flags and environment variables. Command-line flags take precedence over environment variables.

Environment variable names are the uppercase flag name with dashes changed to underscores. For example, `-body-limit` maps to `BODY_LIMIT`.

## Flags and environment variables

| Flag | Env var | Default | Meaning |
| --- | --- | --- | --- |
| `-port` | `PORT` | `8080` | TCP port Worker listens on. |
| `-log-level` | `LOG_LEVEL` | `debug` | Log level. Valid values follow zerolog levels such as `trace`, `debug`, `info`, `warn`, and `error`. |
| `-body-limit` | `BODY_LIMIT` | `1M` | Maximum request body size. Accepts values such as `4K`, `4KB`, `10M`, or `1G`. Empty string disables the limit. |
| `-request-limit` | `REQUEST_LIMIT` | `0` | Request rate for each IP. `0` disables rate limiting. |
| `-request-limit-time-window` | `REQUEST_LIMIT_TIME_WINDOW` | `180` | Rate-limit window in seconds. |
| `-cache-size` | `CACHE_SIZE` | `100` | Number of compiled query plans to keep in the in-memory LRU cache. |
| `-fs-root` | `FS_ROOT` | Current working directory | Filesystem root for FQL `IO::FS` functions. |
| `-chrome-ip` | `CHROME_IP` | `127.0.0.1` | Chrome remote debugging host. |
| `-chrome-port` | `CHROME_PORT` | `9222` | Chrome remote debugging port. |
| `-no-chrome` | `NO_CHROME` | `false` | Disable the CDP driver and skip Chrome checks. |
| `-version` | `VERSION` | `false` | Print Worker and Ferret version metadata, then exit. |
| `-help` | `HELP` | `false` | Print flag help, then exit. |

The current `v2.0.0-rc.16` source rejects `-cache-size=0` during startup with `must provide a positive size`, even though the generated flag help still says `0` disables caching. Use a positive cache size.

## Basic examples

Run Worker on the default port:

{{< terminal >}}
worker
{{< /terminal >}}

Run on a different port with quieter logs:

{{< terminal >}}
worker -port=3000 -log-level=info
{{< /terminal >}}

Set the same values with environment variables:

{{< terminal >}}
PORT=3000 LOG_LEVEL=info worker
{{< /terminal >}}

If both are set, the flag wins:

{{< terminal >}}
PORT=3000 worker -port=8080
{{< /terminal >}}

## Request limits

Use `-body-limit` to reject oversized JSON request bodies before Worker tries to parse them:

{{< terminal >}}
worker -body-limit=10M
{{< /terminal >}}

Use `-request-limit` to enable per-IP rate limiting. The `/health` endpoint is excluded from the rate limiter.

{{< terminal >}}
worker \
  -request-limit=10 \
  -request-limit-time-window=60
{{< /terminal >}}

This allows requests from each IP at the configured rate and uses the time window to size rate-limiter bursts.

## Filesystem access

`-fs-root` sets the root directory used by FQL filesystem functions. When the flag is omitted, Worker uses its current working directory.

{{< terminal >}}
worker -fs-root=/var/lib/ferret-worker
{{< /terminal >}}

An explicitly empty value is invalid:

{{< terminal >}}
worker -fs-root=""
{{< /terminal >}}

Keep `-fs-root` narrow in shared environments. Worker does not add authentication or per-client filesystem isolation on top of the Ferret runtime.

## Chrome and browser-backed scripts

By default, Worker expects Chrome or Chromium to expose the Chrome DevTools Protocol at `http://127.0.0.1:9222`.

Point Worker at another Chrome host:

{{< terminal >}}
worker -chrome-ip=chrome.internal -chrome-port=9222
{{< /terminal >}}

Disable Chrome when you only need non-browser scripts or the memory HTML driver:

{{< terminal >}}
worker -no-chrome
{{< /terminal >}}

With Chrome disabled, scripts that select `{ driver: "cdp" }` cannot use the browser driver.

## Query plan cache

Worker caches compiled query plans by query text. Reusing the same script text across requests avoids recompiling it each time.

{{< terminal >}}
worker -cache-size=500
{{< /terminal >}}

Parameters do not change the cache key. Put changing values in `params` instead of interpolating them into the FQL text when you want cache reuse.
