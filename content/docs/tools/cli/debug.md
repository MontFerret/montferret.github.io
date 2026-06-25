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
Paused at script.fql:1
1 | LET x = 42
Type "help" for available commands.
(fdb)
```

## Requirements

The debugger requires the builtin runtime and a source file. It does not support compiled artifacts, stdin, inline evaluation (`--eval`), or remote runtimes.

## Set breakpoints

Use the `break` command (alias `b`) to set a breakpoint at a source location:

```
(fdb) break 5
Breakpoint 1 set at script.fql:5 (next-file).
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
ID  Requested      Bound           Mode       State
1   script.fql:5   script.fql:5    next-file  bound
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
| `pause` | | Request a pause after the next resume |

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
  x = 42
  name = "Ada"
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
| `pause` | | Request pause after next resume |
| `where` | `w`, `bt` | Show stack trace |
| `locals` | `l` | Show local variables |
| `print <expr>` | `p`, `e`, `eval` | Evaluate a safe expression |
| `quit` | `q` | Stop debugging and exit |

## Runtime and browser flags

The debug command accepts the same runtime and browser flags as [`ferret run`](../run/#runtime-and-browser-flags). The `--runtime` flag is ignored since debugging requires the builtin runtime.
