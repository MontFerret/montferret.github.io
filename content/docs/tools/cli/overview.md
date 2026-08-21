---
title: "Overview"
weight: 5
draft: false
description: "Understand the Ferret CLI workflow, command groups, runtime selection, and persistent configuration."
---

# Overview

The Ferret CLI provides the local command-line workflow around FQL. It can run source scripts, open an interactive shell, check and format code, debug execution, build and inspect bytecode, manage browser processes, and work with Ferret modules.

Run an expression without creating a file:

{{< terminal >}}
ferret run --eval 'return 1 + 2'
{{< /terminal >}}

The command prints the serialized result:

```text
3
```

## Typical workflow

1. Follow [Installation]({{< ref "installation" >}}) and verify the CLI with `ferret version`.
2. Use [Run]({{< ref "run" >}}) for scripts and inline expressions, or [REPL]({{< ref "repl" >}}) while exploring FQL interactively.
3. Run [Check]({{< ref "check" >}}) before execution and [Fmt]({{< ref "fmt" >}}) when source should use the canonical format.
4. Use [Debug]({{< ref "debug" >}}) for breakpoints, stepping, and variable inspection.
5. Use [Build]({{< ref "build" >}}) to create a bytecode artifact and [Inspect]({{< ref "inspect" >}}) to examine compiled output.

## Command groups

| Area | Commands |
| --- | --- |
| Execute FQL | [Run]({{< ref "run" >}}), [REPL]({{< ref "repl" >}}) |
| Validate and format | [Check]({{< ref "check" >}}), [Fmt]({{< ref "fmt" >}}) |
| Debug and inspect | [Debug]({{< ref "debug" >}}), [Build]({{< ref "build" >}}), [Inspect]({{< ref "inspect" >}}) |
| Browser processes | [Browser]({{< ref "browser" >}}) |
| Modules and compatibility | [Mod]({{< ref "mod" >}}), [Migrate]({{< ref "migrate" >}}) |
| Persistent settings | [Configuration]({{< ref "configuration" >}}) |

Run `ferret <command> --help` for the flags accepted by a specific command.

## Runtime and browser boundaries

The CLI uses its builtin Ferret runtime by default. The runtime owns FQL parsing, compilation, module behavior, parameter semantics, execution, and result serialization. The CLI owns command parsing, source selection, runtime selection, persistent settings, logging, and local browser-process management.

[Run]({{< ref "run" >}}) can send source queries to a Worker-compatible HTTP runtime with `--runtime`. Remote runtimes execute source queries; compiled artifacts and interactive debugging require the builtin local runtime.

Browser-backed FQL also depends on the selected runtime and its registered modules. The [Browser]({{< ref "browser" >}}) commands start and stop managed Chrome or Chromium processes, while runtime and FQL options select how a script uses the browser.

## Configuration and logging

All commands accept `--log-level`, `--log-output`, and `--log-file`. The [`config`]({{< ref "configuration" >}}) command persists supported settings in `$HOME/.ferret/config.yaml`.

Settings can also come from environment variables with the `FERRET_` prefix. Command-line flags take priority over environment variables, which take priority over the configuration file and builtin defaults.

## Next steps

{{< docs-related tiles="tools-cli-installation,tools-cli-run,tools-cli-repl,tools-cli-migrate" >}}
