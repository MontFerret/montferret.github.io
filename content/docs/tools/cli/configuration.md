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

Settings are stored in `~/.ferret/config.yaml`, inside the resolved user home directory on every supported platform. The file is created automatically the first time you run any `ferret` command.

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

`config set` validates the complete persisted filesystem and HTTP policy before writing. Invalid values or conflicting policy controls leave the config file unchanged.

## Get a value

{{< terminal >}}
ferret config get runtime
{{< /terminal >}}

## Unset a value

Remove a persisted value and restore its implicit default:

{{< terminal >}}
ferret config unset policy-http-allowed-hosts
{{< /terminal >}}

This leaves environment variables and explicit command-line flags unchanged.

## List all values

{{< terminal >}}
ferret config list
{{< /terminal >}}

The `list` command (alias `ls`) prints every configurable key and its current value.

## Available config keys

| Key | Default | Description |
| --- | --- | --- |
| `log-level` | `info` | Logging level: `trace`, `debug`, `info`, `warn`, `error`, `fatal` |
| `log-output` | `stderr` | Log output destination: `stderr`, `file`, `none` |
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

## Filesystem policy

Filesystem and HTTP execution policies apply to the builtin runtime on `run`, `repl`, and `debug`. Supplying them with a remote runtime is a configuration error. The `version` command has no execution-policy flags.

The builtin filesystem is writable and rooted at the CLI's current working directory by default. Choose a narrower relative or absolute root, and optionally make it read-only:

{{< terminal >}}
ferret run --policy-fs-root ./fixtures --policy-fs-read-only script.fql
{{< /terminal >}}

The config keys are `policy-fs-root` and `policy-fs-read-only`; their environment variables are `FERRET_POLICY_FS_ROOT` and `FERRET_POLICY_FS_READ_ONLY`.

## HTTP policy

Ferret's HTTP client blocks localhost, loopback, private-network, and link-local access by default. These controls apply to integrations using that client, such as `io::net::http` and `net::rest`. They do not govern requests made by Chrome itself. The `proxy` and `user-agent` options configure HTML/browser drivers separately.

For a script that intentionally uses a local development service:

{{< terminal >}}
ferret run --policy-http-allow-localhost --policy-http-default-headers='{"X-Trace":"local"}' script.fql
{{< /terminal >}}

Each key below is also a `--flag` on `run`, `repl`, and `debug`:

| Flag and config key | Environment variable | Default | Behavior |
| --- | --- | --- | --- |
| `policy-http-allowed-schemes` | `FERRET_POLICY_HTTP_ALLOWED_SCHEMES` | `http,https` | Allowed URL schemes |
| `policy-http-allowed-methods` | `FERRET_POLICY_HTTP_ALLOWED_METHODS` | `GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS` | Allowed HTTP methods |
| `policy-http-allowed-hosts` | `FERRET_POLICY_HTTP_ALLOWED_HOSTS` | unrestricted | Exact allowed hosts or `host:port` values |
| `policy-http-blocked-hosts` | `FERRET_POLICY_HTTP_BLOCKED_HOSTS` | none | Exact blocked hosts or `host:port` values |
| `policy-http-allow-localhost` | `FERRET_POLICY_HTTP_ALLOW_LOCALHOST` | `false` | Allow localhost and loopback addresses |
| `policy-http-allow-private-networks` | `FERRET_POLICY_HTTP_ALLOW_PRIVATE_NETWORKS` | `false` | Allow private-network addresses |
| `policy-http-allow-link-local` | `FERRET_POLICY_HTTP_ALLOW_LINK_LOCAL` | `false` | Allow link-local addresses |
| `policy-http-default-headers` | `FERRET_POLICY_HTTP_DEFAULT_HEADERS` | none | Default request headers as a JSON string map |
| `policy-http-blocked-request-headers` | `FERRET_POLICY_HTTP_BLOCKED_REQUEST_HEADERS` | none | Block requests containing these header names |
| `policy-http-timeout` | `FERRET_POLICY_HTTP_TIMEOUT` | `30s` | Overall HTTP timeout |
| `policy-http-no-timeout` | `FERRET_POLICY_HTTP_NO_TIMEOUT` | `false` | Explicitly disable the overall timeout |
| `policy-http-max-request-size` | `FERRET_POLICY_HTTP_MAX_REQUEST_SIZE` | `16777216` | Maximum request-body size in bytes |
| `policy-http-unlimited-request-size` | `FERRET_POLICY_HTTP_UNLIMITED_REQUEST_SIZE` | `false` | Explicitly disable the request-body limit |
| `policy-http-max-response-size` | `FERRET_POLICY_HTTP_MAX_RESPONSE_SIZE` | `16777216` | Maximum response-body size in bytes |
| `policy-http-unlimited-response-size` | `FERRET_POLICY_HTTP_UNLIMITED_RESPONSE_SIZE` | `false` | Explicitly disable the response-body limit |
| `policy-http-max-response-header-size` | `FERRET_POLICY_HTTP_MAX_RESPONSE_HEADER_SIZE` | `1048576` | Maximum response-header size in bytes |
| `policy-http-follow-redirects` | `FERRET_POLICY_HTTP_FOLLOW_REDIRECTS` | `true` | Follow HTTP redirects |
| `policy-http-max-redirects` | `FERRET_POLICY_HTTP_MAX_REDIRECTS` | `10` | Maximum redirects to follow |

Only explicitly supplied flags, environment variables, and config values override Ferret's defaults. Numeric zero retains the Ferret default. Use `no-timeout`, `unlimited-request-size`, or `unlimited-response-size` to disable those limits. Enabling one conflicts with explicitly setting its corresponding timeout or size limit.

List flags accept repeated or comma-separated values. For `config set`, use a comma-separated string. Default headers require a JSON object with string values:

{{< terminal >}}
ferret config set policy-http-allowed-hosts api.example.com,cdn.example.com
ferret config set policy-http-default-headers '{"X-Trace":"local"}'
{{< /terminal >}}

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
  Self: <cli-version>
  Runtime: <runtime-version>
```

To read a Wire host's version:

{{< terminal >}}
ferret version --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

`Self` remains the CLI version. `Runtime` is the exact hosted runtime version supplied through the Universal API, including an empty value. Connection establishment and metadata retrieval share the Wire handshake timeout. This command accepts runtime selection flags; it has no browser or policy flags. Configured or environmental builtin settings are still rejected in Wire mode. The [Wire guide]({{< ref "wire" >}}#read-the-hosted-version) explains metadata behavior and the REPL banner.

To check a legacy HTTP runtime's version:

{{< terminal >}}
ferret version --runtime http://localhost:8080
{{< /terminal >}}

## Next steps

{{< docs-related tiles="tools-cli,tools-cli-run,tools-cli-wire,tools-cli-browser" >}}
