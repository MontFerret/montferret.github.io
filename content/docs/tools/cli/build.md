---
title: "Build"
weight: 50
draft: false
description: "Compile FQL scripts into bytecode artifacts for faster execution."
---

# Build

The `build` command compiles FQL scripts into bytecode artifacts. Compiled artifacts skip the parsing and compilation steps at execution time.

## Compile a script

{{< terminal >}}
ferret build script.fql
{{< /terminal >}}

This produces `script.fqlc` alongside the source file.

## Compile multiple scripts

{{< terminal >}}
ferret build script.fql utils.fql
{{< /terminal >}}

Or use a glob:

{{< terminal >}}
ferret build *.fql
{{< /terminal >}}

Each input file produces a corresponding `.fqlc` artifact in the same directory.

## Custom output path

Use `-o` (or `--output`) to control where artifacts are written.

For a single input, the output can be a file or directory:

{{< terminal >}}
ferret build script.fql -o dist/script.fqlc
{{< /terminal >}}

For multiple inputs, the output must be a directory:

{{< terminal >}}
ferret build *.fql -o dist/
{{< /terminal >}}

The output directory is created automatically if it does not exist.

## Run a compiled artifact

Pass the `.fqlc` file to `ferret run`:

{{< terminal >}}
ferret run script.fqlc
{{< /terminal >}}

Compiled artifacts require the builtin runtime. They cannot be used with a remote runtime.

## Inspect compiled output

Use [`ferret inspect`](../inspect/) to disassemble a compiled program and examine its bytecode, constants, and function definitions.
