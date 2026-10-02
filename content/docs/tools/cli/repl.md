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
Welcome to Ferret REPL {{< data "versions.runtime.v2" >}}
Please use `exit` or `Ctrl-D` to exit this program.
>
```

## Evaluate expressions

Type any FQL expression at the prompt. The result prints as soon as evaluation finishes:

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

The REPL accepts the same runtime and browser flags as [`ferret run`](../run/#runtime-and-browser-flags). For example, to start a REPL with a headless browser available:

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
Welcome to Ferret REPL (Wire runtime: tcp://127.0.0.1:54321; version: v2.0.0-alpha.57)
```

The shell reuses one runtime and connection across submissions. Syntax diagnostics go to stderr and allow corrected queries. Transport loss or resource failure ends the shell; it never reconnects or falls back. The host owns browser integrations and filesystem policy, so explicit local browser or policy settings are rejected. See the [Wire runtime guide]({{< ref "wire" >}}) for setup, compatibility, and metadata behavior.

## Exit

Type `exit` or press `Ctrl-D` to leave the REPL. Ctrl+C stops the shell and interrupts pending input or execution.

## Next steps

{{< docs-related tiles="language,tools-cli-run,tools-cli-wire,tools-playground" >}}
