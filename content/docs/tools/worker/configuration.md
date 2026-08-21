---
title: "Configuration"
weight: 30
draft: false
description: "Configure Worker flags, environment variables, limits, filesystem access, outbound HTTP, cache size, and Chrome."
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
| `-policy-fs-root` | `POLICY_FS_ROOT` | Current working directory | Filesystem root for FQL `io::fs` functions. |
| `-policy-http-allowed-hosts` | `POLICY_HTTP_ALLOWED_HOSTS` | Empty | Comma-separated exact hosts or `host:port` values allowed for outbound Ferret HTTP requests. |
| `-http-allow-all-hosts` | `HTTP_ALLOW_ALL_HOSTS` | `false` | Allow outbound Ferret HTTP requests to any host while retaining the other HTTP policy checks. |
| `-policy-http-blocked-hosts` | `POLICY_HTTP_BLOCKED_HOSTS` | Empty | Comma-separated exact hosts or `host:port` values blocked for outbound Ferret HTTP requests. |
| `-policy-http-timeout` | `POLICY_HTTP_TIMEOUT` | `10s` | Overall timeout for outbound Ferret HTTP requests. |
| `-policy-http-max-request-size` | `POLICY_HTTP_MAX_REQUEST_SIZE` | `1048576` | Maximum outbound HTTP request body size in bytes. `0` disables the limit. |
| `-policy-http-max-response-size` | `POLICY_HTTP_MAX_RESPONSE_SIZE` | `10485760` | Maximum outbound HTTP response body size in bytes. `0` disables the limit. |
| `-policy-http-max-redirects` | `POLICY_HTTP_MAX_REDIRECTS` | `3` | Maximum redirects to follow. `0` uses the Go standard library default. |
| `-policy-http-follow-redirects` | `POLICY_HTTP_FOLLOW_REDIRECTS` | `true` | Follow outbound HTTP redirects. |
| `-policy-http-allow-localhost` | `POLICY_HTTP_ALLOW_LOCALHOST` | `false` | Allow localhost and loopback literal addresses. |
| `-policy-http-allow-private-networks` | `POLICY_HTTP_ALLOW_PRIVATE_NETWORKS` | `false` | Allow private-network literal IP addresses. |
| `-policy-http-blocked-request-headers` | `POLICY_HTTP_BLOCKED_REQUEST_HEADERS` | `Authorization,Cookie,Proxy-Authorization` | Comma-separated request headers removed from outbound Ferret HTTP requests. |
| `-chrome-ip` | `CHROME_IP` | `127.0.0.1` | Chrome remote debugging host. |
| `-chrome-port` | `CHROME_PORT` | `9222` | Chrome remote debugging port. |
| `-no-chrome` | `NO_CHROME` | `false` | Disable the CDP driver and skip Chrome checks. |
| `-version` | `VERSION` | `false` | Print Worker and Ferret version metadata, then exit. |
| `-help` | `HELP` | `false` | Print flag help, then exit. |

The current `v{{< data "versions.worker.v2" >}}` source rejects `-cache-size=0` during startup with `must provide a positive size`, even though the generated flag help still says `0` disables caching. Use a positive cache size.

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

`-policy-fs-root` sets the root directory used by FQL filesystem functions. When the flag is omitted, Worker uses its current working directory.

{{< terminal >}}
worker -policy-fs-root=/var/lib/ferret-worker
{{< /terminal >}}

An explicitly empty value is invalid:

{{< terminal >}}
worker -policy-fs-root=""
{{< /terminal >}}

Keep `-policy-fs-root` narrow in shared environments. Worker does not add authentication or per-client filesystem isolation on top of the Ferret runtime.

## Outbound HTTP access

Worker disables Ferret's policy-backed outbound HTTP client by default. Configure an allowlist when scripts need to call specific services:

{{< terminal >}}
worker \
  -policy-http-allowed-hosts=api.example.com,uploads.example.com:8443 \
  -policy-http-timeout=5s \
  -policy-http-max-response-size=5242880
{{< /terminal >}}

Allowed and blocked host entries are exact host or `host:port` values. Separate multiple values with commas in both flags and environment variables.

Use `-http-allow-all-hosts` only when an explicit allowlist is not practical:

{{< terminal >}}
worker -http-allow-all-hosts
{{< /terminal >}}

`-http-allow-all-hosts` cannot be combined with `-policy-http-allowed-hosts`. Blocked hosts, timeout, request and response size limits, redirect controls, and blocked request headers still apply. Localhost, loopback literals, and private-network literal IP addresses remain blocked unless their corresponding `-policy-http-allow-*` flag is enabled.

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

## Next steps

{{< docs-related tiles="tools-worker-deployment,tools-worker-api,tools-cli-configuration" >}}
