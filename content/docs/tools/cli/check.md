---
title: "Check"
weight: 35
draft: false
description: "Verify FQL scripts for syntax and semantic errors without executing them."
---

# Check

The `check` command compiles FQL scripts without executing them, reporting any syntax or semantic errors.

{{< terminal >}}
ferret check script.fql
{{< /terminal >}}

If the script is valid, the command exits silently with code 0.

## Check multiple files

Pass multiple files to check them in a single run:

{{< terminal >}}
ferret check script.fql utils.fql
{{< /terminal >}}

Or use a glob:

{{< terminal >}}
ferret check *.fql
{{< /terminal >}}

When any file has errors, the command prints the diagnostics and exits with a summary:

```
2 of 5 scripts have errors
```

## Read from stdin

Pipe a script through stdin:

{{< terminal >}}
cat script.fql | ferret check
{{< /terminal >}}

## Use in CI

The `check` command returns a non-zero exit code when errors are found, making it suitable for CI pipelines:

```yaml
- name: Lint FQL scripts
  run: ferret check src/**/*.fql
```
