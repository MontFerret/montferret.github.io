---
title: "Configuration"
weight: 80
draft: false
description: "Manage persistent CLI configuration for runtime, browser, and logging settings."
aliases:
  - /docs/tools/cli/config/
---

# Configuration

The `config` command reads and writes persistent settings that apply to all CLI commands.

## Config file

Settings are stored in `~/.ferret/config.yaml`. The file is created automatically the first time you run any `ferret` command.

## Configuration priority

When the same setting is specified in multiple places, the highest priority wins:

1. Command-line flags (highest)
2. Environment variables
3. Config file
4. Built-in defaults (lowest)

## Environment variables

Every config key maps to an environment variable with the `FERRET_` prefix. Dashes become underscores:

| Config key | Environment variable |
| --- | --- |
| `runtime` | `FERRET_RUNTIME` |
| `runtime-endpoint` | `FERRET_RUNTIME_ENDPOINT` |
| `runtime-connect-timeout` | `FERRET_RUNTIME_CONNECT_TIMEOUT` |
| `browser-address` | `FERRET_BROWSER_ADDRESS` |
| `log-level` | `FERRET_LOG_LEVEL` |

## Set a value

{{< terminal >}}
ferret config set runtime builtin
{{< /terminal >}}

{{< terminal >}}
ferret config set browser-address http://127.0.0.1:9333
{{< /terminal >}}

## Get a value

{{< terminal >}}
ferret config get runtime
{{< /terminal >}}

## List all values

{{< terminal >}}
ferret config list
{{< /terminal >}}

The `list` command (alias `ls`) prints every configurable key and its current value.

## Available config keys

| Key | Default | Description |
| --- | --- | --- |
| `log-level` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `log-output` | `stderr` | Log output destination: `stderr`, `file` |
| `log-file` | `ferret.log` | Log file path (when `log-output` is `file`) |
| `runtime` | `builtin` | `builtin`, `wire`, or a legacy Worker-compatible HTTP URL |
| `runtime-endpoint` | | Required Wire endpoint: `tcp://127.0.0.1:<port>` |
| `runtime-connect-timeout` | `5s` | Positive Wire connection and handshake timeout |
| `policy-fs-root` | Current working directory | Builtin filesystem sandbox root |
| `policy-fs-read-only` | `false` | Make the builtin filesystem sandbox read-only |
| `browser-cookies` | `false` | Persist cookies between queries |
| `browser-address` | `http://127.0.0.1:9222` | Browser remote debugging address |
| `browser-open` | `false` | Open a browser for script execution |
| `browser-headless` | `false` | Open browser in headless mode |
| `proxy` | | Proxy server address |
| `user-agent` | | Custom User-Agent header |

## Wire settings

Persist runtime selection, endpoint, and handshake timeout with `config set`:

{{< terminal >}}
ferret config set runtime wire
ferret config set runtime-endpoint tcp://127.0.0.1:54321
ferret config set runtime-connect-timeout 5s
{{< /terminal >}}

Only IPv4 loopback endpoints with a port from 1–65535 are accepted. The timeout must be strictly positive. Endpoint and timeout settings require Wire mode and are validated before connecting.

Wire mode rejects explicit builtin filesystem/HTTP policies and browser settings from flags, environment variables, or this file, including default-valued settings. Configure capabilities on the host and remove incompatible client values. Untouched defaults remain acceptable.

Use `config unset` to remove a persisted setting. When switching back to builtin mode, also remove the Wire endpoint and timeout:

{{< terminal >}}
ferret config set runtime builtin
ferret config unset runtime-endpoint
ferret config unset runtime-connect-timeout
{{< /terminal >}}

Unsetting a key does not remove an environment variable or override an explicit flag. See the [Wire runtime guide]({{< ref "wire" >}}) for connection setup and filesystem ownership.

## Global logging flags

The logging flags are available on every command as persistent flags:

{{< terminal >}}
ferret run --log-level debug script.fql
{{< /terminal >}}

{{< terminal >}}
ferret run --log-output file --log-file run.log script.fql
{{< /terminal >}}

## Version

The `version` command shows the CLI version and the runtime version:

{{< terminal >}}
ferret version
{{< /terminal >}}

```
Version:
  Self: {{< data "versions.cli.v2" >}}
  Runtime: {{< data "versions.runtime.v2" >}}
```

To read a Wire host's version:

{{< terminal >}}
ferret version --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

`Self` remains the CLI version. `Runtime` is the exact hosted Universal API version, including an empty value. Connection establishment and metadata retrieval share the Wire handshake timeout. This command accepts runtime selection flags; it has no browser or policy flags. Configured or environmental builtin settings are still rejected in Wire mode. The [Wire guide]({{< ref "wire" >}}#read-the-hosted-version) explains metadata behavior and the REPL banner.

To check a legacy HTTP runtime's version:

{{< terminal >}}
ferret version --runtime http://localhost:8080
{{< /terminal >}}

## Next steps

{{< docs-related tiles="tools-cli,tools-cli-run,tools-cli-wire,tools-cli-browser" >}}
