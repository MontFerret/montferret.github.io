---
title: "Binary Runtime"
weight: 80
draft: false
description: "Run Lab tests through an external Ferret-compatible binary."
---

# Binary Runtime

Use a binary runtime when you want Lab to execute tests through an external Ferret-compatible command instead of the builtin runtime.

{{< terminal >}}
lab run --runtime=bin:/usr/local/bin/ferret tests/
{{< /terminal >}}

Lab still owns source loading, local services, retries, concurrency, and reporting. The external binary owns FQL execution.

## Execution contract

For each test query, Lab:

1. Starts the configured binary.
2. Writes the FQL source to stdin.
3. Passes parameters as `--param=name:value` arguments.
4. Reads combined stdout and stderr.
5. Treats a non-zero process exit as a test failure.

The external binary must accept FQL source from stdin and support the `--param` argument shape used by Ferret.

## Runtime path

Use `bin:/absolute/path/to/binary`:

{{< terminal >}}
lab run --runtime=bin:/opt/ferret/bin/ferret tests/
{{< /terminal >}}

To run a binary found through `PATH`, use `bin://name`:

{{< terminal >}}
lab run --runtime=bin://ferret tests/
{{< /terminal >}}

## Forward shared parameters

Runtime params other than `flags` are serialized as JSON and passed to the binary as shared `--param` values for every query.

{{< terminal >}}
lab run \
  --runtime=bin:/usr/local/bin/ferret \
  --runtime-param='baseUrl:"https://example.com"' \
  tests/
{{< /terminal >}}

If the test also receives user or Lab system parameters, Lab forwards those as additional `--param` arguments for the query.

## Pass raw binary flags

Use the special `flags` runtime param to append raw command-line flags before Lab-generated `--param` arguments.

{{< terminal >}}
lab run \
  --runtime=bin:/usr/local/bin/ferret \
  --runtime-param='flags:["--timeout=60","--verbose"]' \
  tests/
{{< /terminal >}}

`flags` must be an array of strings. Invalid values are rejected before tests run.

## Check the runtime version

Use `lab version --runtime=bin:...` to ask the external binary for its version.

{{< terminal >}}
lab version --runtime=bin:/usr/local/bin/ferret
{{< /terminal >}}

Lab runs:

```text
/usr/local/bin/ferret version
```

The returned text is displayed as the runtime version.

## Local services with binary runtimes

When you use `--serve` or `--mock`, Lab forwards the generated `@lab.static` and `@lab.mock` values as parameters to the binary runtime.

If the binary runs on the same machine as Lab, the default `127.0.0.1` service URLs usually work. If the binary is a wrapper into a container or remote environment, set `--serve-host` and `--serve-bind` so the runtime can reach the local services.
