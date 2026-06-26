---
title: "Config"
weight: 80
draft: false
description: "Manage persistent CLI configuration for runtime, browser, and logging settings."
---

# Config

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
| `runtime` | `builtin` | Runtime type: `builtin` or a remote URL |
| `runtime-fs-root` | Current working directory | File system root for the runtime |
| `browser-cookies` | `false` | Persist cookies between queries |
| `browser-address` | `http://127.0.0.1:9222` | Browser remote debugging address |
| `browser-open` | `false` | Open a browser for script execution |
| `browser-headless` | `false` | Open browser in headless mode |
| `proxy` | | Proxy server address |
| `user-agent` | | Custom User-Agent header |

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
  Self: 2.0.0-alpha.26
  Runtime: 2.0.0-alpha.26
```

To check the version of a remote runtime:

{{< terminal >}}
ferret version --runtime http://localhost:8080
{{< /terminal >}}

## Next steps

{{< docs-related tiles="tools-cli,tools-cli-run,tools-cli-browser" >}}
