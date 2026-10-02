---
title: "Wire runtime"
weight: 27
draft: false
description: "Run FQL source and an interactive shell against an application's configured Wire runtime."
---

# Use a Wire runtime

Use Wire when a separate application owns the Ferret engine, custom functions, modules, and runtime configuration that your script needs. The CLI sends source and parameters to that host:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 --eval 'RETURN 42'
{{< /terminal >}}

Replace the endpoint with the one printed by your host. Wire source execution is available in `run`, its `exec` alias, and `repl`. Builtin execution remains the default. Worker-compatible HTTP URLs in `--runtime` retain their existing behavior.

These commands require a CLI build with Wire support. Check `ferret run --help` for `wire` in the runtime flag. The host dependency versions below do not identify a CLI release.

## Start a development host

The [CLI development example](https://github.com/MontFerret/cli/tree/master/examples/wire-host) exposes a configured native engine through the Universal API and Wire. It adds a host-only `DEMO::IDENTITY(value)` function and prints an ephemeral IPv4 loopback endpoint.

In terminal 1, from a CLI source checkout containing this integration, build the CLI and start the host:

{{< terminal >}}
GOFLAGS=-mod=mod make compile
go run -mod=mod ./examples/wire-host
{{< /terminal >}}

The commands require Go {{< data "versions.cli.go" >}} or later. Module mode bypasses the checkout's vendor directory. Example endpoint:

```text
tcp://127.0.0.1:54321
```

In terminal 2, from the same CLI checkout, run the built binary using the endpoint printed by terminal 1:

{{< terminal >}}
./bin/ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 --eval 'RETURN DEMO::IDENTITY(@value)' --param value=42
{{< /terminal >}}

It prints `42`. The function runs on the host; the CLI does not register it locally. Stop the example host with Ctrl+C when finished. Its shutdown is bounded and closes the host-owned engine after Wire work settles.

For a host you build into your own application, see [Remote runtimes (Wire)]({{< ref "docs/embedding/wire" >}}).

### Compatible host versions

This integration uses the following compatible dependency set:

| Dependency | Version |
| --- | --- |
| Ferret Core | `v2.0.0-alpha.57` |
| Universal API | `v1.0.0-alpha.20` |
| Wire | `v1.0.0-alpha.2` |
| gRPC | `v1.84.0` |

The host must implement the alpha.2 handshake, including hosted runtime version metadata. Alpha.1 hosts are incompatible. Pin a compatible set when building a host; Wire is pre-stable.

## Run files, stdin, and parameters

After placing a Wire-capable `ferret` binary on your `PATH`, use the same source forms as builtin execution:

{{< terminal >}}
ferret exec --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 script.fql
cat script.fql | ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

Input files are read by the CLI. It sends their source name and contents to the host, where compilation and execution happen. The CLI creates no local engine for Wire execution. Compiled artifacts and `ferret debug` require the builtin runtime.

Parameters use the same [JSON-aware parsing]({{< ref "run" >}}#pass-parameters) in builtin and Wire modes. Null, arrays, objects, and explicitly quoted strings retain their values:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 --eval 'RETURN DEMO::IDENTITY(@value)' --param value='{"items":[42,null,"123"]}'
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 --eval 'RETURN @code' --param code='"123"'
{{< /terminal >}}

## Use the interactive shell

{{< terminal >}}
ferret repl --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

For the example host, the banner and a submission look like:

```text
Welcome to Ferret REPL (Wire runtime: tcp://127.0.0.1:54321; version: v2.0.0-alpha.57)
Please use `exit` or `Ctrl-D` to exit this program.
> RETURN DEMO::IDENTITY(42)
42
```

The REPL opens one runtime and transport pair and reuses it for every submission. Query diagnostics go to stderr; correct the query and submit again. Transport loss or resource failure ends the shell with an error. It does not reconnect, retry a query, or fall back to builtin execution.

The existing [editing and multiline conventions]({{< ref "repl" >}}#multi-line-input) apply. Use `exit` or Ctrl+D to exit; Ctrl+C stops the shell, including pending input or execution.

## Read the hosted version

{{< terminal >}}
ferret version --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

The command preserves the normal output format:

```text
Version:
  Self: <cli-version>
  Runtime: v2.0.0-alpha.57
```

`Self` identifies the CLI. `Runtime` and the Wire REPL banner use the exact version supplied by the hosted Universal API runtime. The value is opaque and is preserved without parsing or normalization, including an empty value. The development example derives it from its Ferret dependency's Go build information. Builtin execution uses the CLI's injected runtime version, or `unknown` in a build without version injection.

Wire captures metadata during the initial handshake. A failed or stalled metadata request fails connection construction within the connection timeout, before the REPL reads input. Later version reads use the captured value and make no extra RPC. The `version` command accepts runtime selection flags without adding browser or policy flags.

## Configure the connection

| Flag and config key | Environment variable | Default |
| --- | --- | --- |
| `runtime` | `FERRET_RUNTIME` | `builtin`; select `wire` explicitly |
| `runtime-endpoint` | `FERRET_RUNTIME_ENDPOINT` | Required in Wire mode; no default port |
| `runtime-connect-timeout` | `FERRET_RUNTIME_CONNECT_TIMEOUT` | `5s`; must be strictly positive |

Only `tcp://127.0.0.1:<port>` is accepted, with a port from 1–65535. Other hosts, schemes, userinfo, paths, queries, and fragments are rejected before connecting. Endpoint and connection-timeout settings require Wire mode.

The timeout bounds connection establishment and the Wire handshake, including hosted metadata retrieval. It does not set a query execution timeout. Command cancellation interrupts connection setup and active work.

These settings are shared by `run`, `exec`, `repl`, and `version`. [Configuration precedence]({{< ref "configuration" >}}#configuration-priority) remains flags, environment variables, config file, then defaults. To persist a Wire connection:

{{< terminal >}}
ferret config set runtime wire
ferret config set runtime-endpoint tcp://127.0.0.1:54321
ferret config set runtime-connect-timeout 5s
{{< /terminal >}}

To restore builtin selection, unset the Wire-specific settings as well:

{{< terminal >}}
ferret config set runtime builtin
ferret config unset runtime-endpoint
ferret config unset runtime-connect-timeout
{{< /terminal >}}

`config unset` removes persisted values; environment variables still take precedence. The example host selects a new ephemeral port each time it starts, so update any saved endpoint after restarting it.

## Host policies and filesystem ownership

The host owns its functions, modules, filesystem, HTTP policy, and browser integrations. FQL filesystem operations run against the host's configured filesystem; the CLI does not send its working directory as a filesystem-root override. The development host uses its own working directory as its filesystem root. Reading a local input script is separate from filesystem access inside that script.

Wire mode rejects explicitly supplied builtin filesystem/HTTP policies and browser settings, including values in config files or environment variables. Even an explicit default-valued flag such as `--browser-headless=false` counts as configuration. Untouched defaults remain acceptable. Configure those capabilities on the host, and remove incompatible client settings before connecting.

Each command owns one logical runtime and transport. Active calls settle before cleanup closes the runtime, then the transport. Closing the CLI connection leaves the host engine usable. After cancellation, Wire uses bounded detached cleanup; cleanup failures remain part of the reported error.

This CLI transport uses **plaintext loopback for trusted local development**, without authentication or TLS. Local access does not authenticate the client or host. HTTP headers and policy flags are not Wire credentials.

## Output and diagnostics

`run` prints the host's encoded result bytes to stdout without decoding or re-encoding them and adds no connection banner. Available output is written even when execution also returns an error. Successfully produced empty output emits nothing. Errors and source diagnostics go to stderr and cause a failing command exit.

To inspect an intentionally invalid query, run:

{{< terminal >}}
ferret run --runtime wire --runtime-endpoint tcp://127.0.0.1:54321 --eval 'LET arr = []

RETURN'
{{< /terminal >}}

The host reports the source location, `RETURN` line, insertion caret after the keyword, and missing-value hint. `RETURN )` highlights the actual invalid token. In the REPL, these query errors allow a corrected submission; connection failures require checking the host and starting a new command.

## Next steps

{{< docs-related tiles="tools-cli-run,tools-cli-repl,embedding-wire,tools-worker" >}}
