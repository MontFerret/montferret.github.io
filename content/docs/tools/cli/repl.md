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
Welcome to Ferret REPL v{{< data "versions.cli.v2" >}}
Please use `exit` or `Ctrl-D` to exit this program.
>
```

## Evaluate expressions

Type any FQL expression at the prompt. The result prints as soon as evaluation finishes:

```
> RETURN 1 + 2
3
```

```
> RETURN { name: "Ada", active: TRUE }
{"active":true,"name":"Ada"}
```

## Multi-line input

For scripts that span multiple lines, start a block with `%`, type the lines, then close with `%` to submit:

```
> %
LET users = [
    { name: "Ada", active: TRUE },
    { name: "Grace", active: FALSE }
]

FOR user IN users
    FILTER user.active
    RETURN user.name
%
["Ada"]
```

## Pass parameters

Use the `--param` flag to pass values into REPL queries. Parameters are accessible as `@name`:

{{< terminal >}}
ferret repl --param greeting=hello
{{< /terminal >}}

```
> RETURN @greeting
"hello"
```

The `--param` flag follows the same format as the [run](../run/) command.

## Runtime and browser flags

The REPL accepts the same runtime and browser flags as [`ferret run`](../run/#runtime-and-browser-flags). For example, to start a REPL with a headless browser available:

{{< terminal >}}
ferret repl --browser-headless
{{< /terminal >}}

## Exit

Type `exit` or press `Ctrl-D` to leave the REPL.

## Next steps

{{< docs-related tiles="language,tools-cli-run,tools-playground" >}}
