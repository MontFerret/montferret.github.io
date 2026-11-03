---
title: "REPL"
weight: 25
draft: false
description: "Launch an interactive FQL shell for experimenting with queries."
---

# REPL

The `repl` command starts an interactive FQL shell where you can type expressions and see results immediately.

{{< terminal >}}
ferret repl
{{< /terminal >}}

The shell prints a version banner and presents a `>` prompt:

```
Welcome to Ferret REPL <runtime-version>
Please use `exit` or `Ctrl-D` to exit this program.
>
```

## Evaluate expressions

Submit an FQL expression with a terminal `RETURN`. Output is written as encoded chunks arrive. Each submission finishes consumption before the shell starts the next query:

```
> return 1 + 2
3
```

```
> return { name: "Ada", active: true }
{"active":true,"name":"Ada"}
```

## Multi-line input

For scripts that span multiple lines, start a block with `%`, type the lines, then close with `%` to submit:

```
> %
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false }
]

return for user in users
    filter user.active
    return user.name
%
["Ada"]
```

## Pass parameters

Use the `--param` flag to pass values into REPL queries. Parameters are accessible as `@name`:

{{< terminal >}}
ferret repl --param greeting=hello
{{< /terminal >}}

```
> return @greeting
"hello"
```

The `--param` flag follows the same format as the [run](../run/) command.

## Runtime and browser flags

The REPL accepts the same runtime and browser flags as [`ferret run`](../run/#runtime-and-browser-flags). Scripts must [select the CDP driver]({{< ref "browser" >}}#select-cdp-in-the-script) to use Chrome or Chromium. For example, to start a REPL with a headless browser available:

{{< terminal >}}
ferret repl --browser-headless
{{< /terminal >}}

## Wire runtime

Run the shell against an application-configured host:

{{< terminal >}}
ferret repl --runtime wire --runtime-endpoint tcp://127.0.0.1:54321
{{< /terminal >}}

The banner identifies the endpoint and the exact version supplied by the host. For the development host:

```text
Welcome to Ferret REPL (Wire runtime: tcp://127.0.0.1:54321; version: v{{< data "versions.wire.core" >}})
```

The shell reuses one runtime and connection across submissions. Syntax diagnostics go to stderr and allow corrected queries. Transport loss or resource failure ends the shell; it never reconnects or falls back. The host owns browser integrations and filesystem policy, so explicit local browser or policy settings are rejected. See the [Wire runtime guide]({{< ref "wire" >}}) for setup, compatibility, and metadata behavior.

## Output and errors

The shell consumes each query's encoded output directly and preserves available bytes if a terminal error follows. Nonempty output receives a trailing newline when needed, including partial output while stdout remains writable. Empty output adds no result newline. Streaming describes encoded delivery; evaluation remains the runtime's responsibility.

Structured query diagnostics, including those reported during consumption, go to stderr and allow a corrected submission. Writer, transport, cancellation, and cleanup failures end the shell. The [HTTP Worker adapter]({{< ref "run" >}}#legacy-http-runtime) prints response bodies as raw output; HTTP statuses do not become query diagnostics.

## Exit

Type `exit` or press `Ctrl-D` to leave the REPL. Ctrl+C stops the shell and interrupts pending input or execution.

## Next steps

{{< docs-related tiles="language,tools-cli-run,tools-cli-wire,tools-playground" >}}
