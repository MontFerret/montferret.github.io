---
title: "Debug"
weight: 30
draft: false
description: "Debug FQL scripts interactively with breakpoints, stepping, and variable inspection."
---

# Debug

The `debug` command starts an interactive debugger for a FQL script:

{{< terminal >}}
ferret debug script.fql
{{< /terminal >}}

The debugger compiles the script, pauses at the first executable location, and presents a `(fdb)` prompt:

```
Ferret debugger started.
Paused at script.fql:1:1
1 | LET x = 42
    ~~~~~~~~~~
Type "help" for available commands.
(fdb)
```

## Requirements

Pass one local FQL source file. The builtin runtime is the default. A Wire host can compile and debug the same file with its own functions and modules. Compiled artifacts, stdin, inline evaluation, and Worker HTTP endpoints are unsupported.

## Debug on a Wire host

Start the [example Wire host]({{< ref "/docs/embedding/wire" >}}#host-the-configured-engine). Save the following as `greeting.fql`:

{{< code lang="fql" >}}
LET greeting = app::greet(@name)
RETURN greeting
{{< /code >}}

From the directory containing the file, connect to the host and pass the script parameter:

{{< terminal >}}
ferret debug --runtime wire --runtime-endpoint tcp://127.0.0.1:50051 --param name=Ada greeting.fql
{{< /terminal >}}

At the initial stop, type `next` and then `locals` to inspect `greeting`. Type `continue` to finish. The CLI reads the local source and sends it to the host. You can step through FQL around a host function. The host retains its native code, filesystem policy, and browser configuration. Exiting the debugger releases the session and leaves the host running.

## Set breakpoints

Use the `break` command (alias `b`) to set a breakpoint at a source location:

```
(fdb) break 5
Breakpoint 1 set at script.fql:5:1 (requested script.fql:5, next-file).
```

Locations can include a column and a filename:

| Format | Example |
| --- | --- |
| Line | `break 5` |
| Line and column | `break 5:10` |
| File and line | `break utils.fql:12` |
| File, line, and column | `break utils.fql:12:4` |

### Binding modes

By default, a breakpoint binds to the next executable location in the file (`--next`). Two other modes are available:

| Flag | Behavior |
| --- | --- |
| `--next` | Bind to the next executable location in the file (default) |
| `--exact` | Bind only at the exact location; fail if it is not executable |
| `--in-function` | Bind to the next executable location within the enclosing function |

```
(fdb) break --exact 10
(fdb) break --in-function 15
```

## Manage breakpoints

List all breakpoints with `breakpoints` (aliases `bp`, `bl`):

```
(fdb) breakpoints
ID  Requested     Bound           Mode       State
1   script.fql:5  script.fql:5:1  next-file  bound
```

Delete a breakpoint by ID with `delete` (alias `d`):

```
(fdb) delete 1
Breakpoint 1 deleted.
```

## Execution control

| Command | Alias | Description |
| --- | --- | --- |
| `continue` | `c` | Resume execution until the next breakpoint or completion |
| `step` | `s` | Step into the next source location |
| `next` | `n` | Step over the current source location |
| `out` | | Step out of the current frame |
| `pause` | | Queue a pause after the next resume on the builtin runtime. Wire accepts pause only while running and rejects it at a stopped prompt. |

Pressing Enter with no input repeats the last `continue`, `step`, `next`, or `out` command.

## Inspect state

| Command | Alias | Description |
| --- | --- | --- |
| `where` | `w`, `bt` | Show the current stack trace |
| `locals` | `l` | Show local variables and parameters in the current frame |
| `print <expr>` | `p`, `e`, `eval` | Evaluate a safe expression (no calls, queries, or mutation) |

```
(fdb) locals
Locals:
  name = "Ada"
  x = 42
```

```
(fdb) print x + 1
43
```

## Exit

Type `quit` (alias `q`) or press `Ctrl-D` to end the debug session.

## Command reference

| Command | Aliases | Description |
| --- | --- | --- |
| `break <location>` | `b` | Set breakpoint at next executable location in file |
| `break --exact <location>` | | Set breakpoint only at exact location |
| `break --next <location>` | | Set breakpoint at next executable location in file |
| `break --in-function <location>` | | Set breakpoint at next executable location in function |
| `breakpoints` | `bp`, `bl` | List all breakpoints |
| `delete <id>` | `d` | Delete a breakpoint |
| `continue` | `c` | Resume execution |
| `step` | `s` | Step into next source location |
| `next` | `n` | Step over current source location |
| `out` | | Step out of current frame |
| `pause` | | Request a pause (builtin runtime at a stopped prompt) |
| `where` | `w`, `bt` | Show stack trace |
| `locals` | `l` | Show local variables |
| `print <expr>` | `p`, `e`, `eval` | Evaluate a safe expression |
| `quit` | `q` | Stop debugging and exit |

## Runtime and browser flags

Use the [runtime flags](../run/#runtime-and-browser-flags) to select Wire, set its endpoint and connection timeout, and pass parameters. Browser and filesystem/HTTP policy flags configure the builtin runtime. Explicit local policies or browser settings are rejected in Wire mode.

## Next steps

{{< docs-related tiles="tools-cli-inspect,tools-cli-run,language-control-flow-error-handling" >}}
