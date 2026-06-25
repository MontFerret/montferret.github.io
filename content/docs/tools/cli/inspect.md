---
title: "Inspect"
weight: 60
draft: false
description: "Disassemble compiled FQL programs and examine bytecode, constants, functions, and source spans."
---

# Inspect

The `inspect` command compiles a FQL script and prints a disassembly of the resulting program. This is useful for understanding how the compiler translates FQL into bytecode.

## Inspect a script

{{< terminal >}}
ferret inspect script.fql
{{< /terminal >}}

Without any filter flags, the full disassembly is printed: summary, bytecode instructions, constant pool, function definitions, and source spans.

## Inspect an inline expression

Use `-e` (or `--eval`) to inspect an expression directly:

{{< terminal >}}
ferret inspect -e 'RETURN 1 + 2'
{{< /terminal >}}

## Read from stdin

{{< terminal >}}
cat script.fql | ferret inspect
{{< /terminal >}}

## Filter sections

When one or more filter flags are set, only the selected sections appear in the output. Multiple flags can be combined.

| Flag | Description |
| --- | --- |
| `--summary` | Show a high-level program summary |
| `--bytecode` | Show bytecode instructions |
| `--constants` | Show the constant pool |
| `--functions` | Show function definitions |
| `--spans` | Show debug source spans |

Show only the bytecode:

{{< terminal >}}
ferret inspect --bytecode script.fql
{{< /terminal >}}

Combine multiple sections:

{{< terminal >}}
ferret inspect --summary --bytecode script.fql
{{< /terminal >}}

When multiple sections are selected, each is preceded by a `==> SectionName <==` header.
