---
title: "Fmt"
weight: 40
draft: false
description: "Format FQL scripts to a consistent style."
---

# Fmt

The `fmt` command formats FQL scripts to a consistent style. By default, it writes language keywords and literals in their canonical lowercase form.

## Format files in place

{{< terminal >}}
ferret fmt script.fql
{{< /terminal >}}

This overwrites the file with the formatted output.

## Format multiple files

{{< terminal >}}
ferret fmt script.fql utils.fql
{{< /terminal >}}

Or use a glob:

{{< terminal >}}
ferret fmt *.fql
{{< /terminal >}}

## Read from stdin

Pipe a script through stdin to print the formatted output without writing a file:

{{< terminal >}}
cat script.fql | ferret fmt
{{< /terminal >}}

## Preview changes

Use `--dry-run` to print the formatted output to stdout instead of overwriting files:

{{< terminal >}}
ferret fmt --dry-run script.fql
{{< /terminal >}}

When formatting multiple files with `--dry-run`, each file's output is preceded by a `==> filename <==` header.

## Formatting options

| Flag | Default | Description |
| --- | --- | --- |
| `--dry-run` | `false` | Print to stdout instead of overwriting files |
| `--print-width` | `80` | Maximum line length |
| `--tab-width` | `4` | Indentation size in spaces |
| `--single-quote` | `false` | Use single quotes instead of double quotes |
| `--bracket-spacing` | `true` | Add spaces inside object brackets |
| `--case-mode` | `lower` | Keyword case: canonical `lower`, or compatibility/preference modes `upper` and `ignore` |

### Examples

Format with a wider line width:

{{< terminal >}}
ferret fmt --print-width 120 script.fql
{{< /terminal >}}

Use single quotes while keeping the canonical lowercase keyword style:

{{< terminal >}}
ferret fmt --single-quote script.fql
{{< /terminal >}}

Use uppercase output for compatibility with an existing codebase or a team preference:

{{< terminal >}}
ferret fmt --case-mode upper script.fql
{{< /terminal >}}

Use `--case-mode ignore` only when an existing workflow depends on that compatibility mode.

## Next steps

{{< docs-related tiles="tools-cli,tools-cli-check,tools-cli-run" >}}
