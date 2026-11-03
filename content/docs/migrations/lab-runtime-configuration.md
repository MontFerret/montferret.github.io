---
title: "Lab runtime configuration"
weight: 30
draft: false
description: "Migrate earlier Lab v2 runtime selectors, adapter settings, and FQL defaults."
---

# Migrate Lab runtime configuration

Lab's runtime target identifies the complete adapter endpoint. Adapter settings use `--runtime-option`; FQL inputs use the existing parameter facilities. This intentionally breaks the superseded v2 runtime flags.

| Previous usage | Replacement |
| --- | --- |
| `--runtime=wire --runtime-endpoint=tcp://host:port` | `--runtime=wire+tcp://host:port` (existing loopback restriction applies) |
| `--runtime-connect-timeout=10s` | `--runtime-option='connect-timeout:"10s"'` |
| `--runtime-param='headers:{...}'` | `--runtime-option='headers:{...}'` |
| `--runtime-param='path:"/v1/execute"'` | `--runtime-option='path:"/v1/execute"'` |
| FQL input through `--runtime-param` | Existing `--param`, bindings, or suite/query parameter facilities |
| `LAB_RUNTIME_PARAM` containing adapter settings | `LAB_RUNTIME_OPTION` |
| `LAB_RUNTIME_ENDPOINT` | Include the endpoint in `LAB_RUNTIME=wire+tcp://127.0.0.1:<port>` |
| `LAB_RUNTIME_CONNECT_TIMEOUT=10s` | `LAB_RUNTIME_OPTION='connect-timeout:"10s"'` |
| `--rp` | Removed; use `--runtime-option` for adapter configuration |

The separate runtime FQL defaults layer has been removed. A configuration combining runtime defaults and per-run values may need to be expressed using existing CLI, binding, or query/assertion parameter facilities; there is no guaranteed lossless one-for-one migration of that layer. No new defaults layer is introduced. External CLI binary execution remains unsupported.

Only `127.0.0.1` and ports 1–65535 are supported for Wire TCP. The new URI does not add encryption or authentication. `connect-timeout` defaults to `"5s"` and bounds connection/handshake setup; test execution timeouts remain unchanged.

HTTP accepts `headers` and existing `cookies` objects of strings and a string `path`. The path override affects execution only; version requests still use the runtime URL's info endpoint. Wire accepts only `connect-timeout`, and builtin rejects adapter options. Unsupported settings fail before connecting.

See [Lab configuration]({{< ref "/docs/tools/lab/configuration" >}}) for current flags, JSON syntax, precedence, and parameter boundaries.
