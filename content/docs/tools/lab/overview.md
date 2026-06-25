---
title: "Overview"
weight: 10
draft: false
description: "Understand what Lab runs, how it prepares tests, and which commands it provides."
---

# Overview

Lab is a test runner for Ferret scripts. It loads test files, prepares parameters, starts optional local services, runs each test through a Ferret runtime, and reports pass or fail.

The runtime still owns FQL parsing, compilation, module behavior, browser behavior, and value semantics. Lab coordinates test execution around that runtime.

## Execution flow

A Lab run follows this shape:

```text
locations -> sources -> test case -> runner -> runtime -> reporter
                         |
                         +-> local services -> @lab.static / @lab.mock
```

The main packages in Lab mirror that flow:

| Area | Responsibility |
| --- | --- |
| Sources | Read test files from local paths, HTTP URLs, or Git repositories. |
| Testing | Interpret `.fql`, `.yaml`, and `.yml` test files. |
| Runner | Apply timeouts, concurrency, retries, repeated runs, and result aggregation. |
| Runtime | Execute FQL through the builtin runtime, an HTTP runtime, or a binary runtime. |
| Local services | Serve static directories and OpenAPI-backed mock APIs. |
| Reporters | Print progress and summary output. |

For local service details, see [Static File Server]({{< ref "static-serving" >}}) and [Mock API Server]({{< ref "mock-api" >}}).

## Test file types

Lab accepts these test files:

| File type | Behavior |
| --- | --- |
| `.fql` | Runs the script as a unit test. The test passes when the runtime returns no error. |
| `.fail.fql` | Runs the script as an expected-failure unit test. The test passes only when the runtime returns an error. |
| `.yaml`, `.yml` | Runs a suite with a `query` script followed by an `assert` script. |

Directories are traversed recursively. Unsupported files are ignored when Lab reads a directory.

## Commands

Lab exposes three commands:

| Command | Purpose |
| --- | --- |
| `lab run` | Run FQL and YAML test files. |
| `lab serve` | Start static and mock API services without running tests. |
| `lab version` | Print the Lab version and the selected runtime's Ferret version. |

Run tests with the builtin runtime:

{{< terminal >}}
lab run tests/
{{< /terminal >}}

Start local fixture services without a test run:

{{< terminal >}}
lab serve --static ./dist@app --mock ./users.yaml@api
{{< /terminal >}}

Inspect the runtime version:

{{< terminal >}}
lab version
{{< /terminal >}}

Use `lab version --runtime=...` when you need the version reported by a remote or binary runtime.

## Runtime boundaries

Lab passes parameters into Ferret under two groups:

| Parameter group | Source |
| --- | --- |
| User parameters | Values passed with `--param` or suite-level `params`. |
| Lab system parameters | Values Lab creates, exposed under `@lab`. |

The system namespace currently includes local service URLs under `@lab.static` and `@lab.mock`, and YAML suite query output under `@lab.data.query`.

Browser-backed FQL depends on the runtime you choose. Lab does not add a separate browser or CDP flag; configure browser access through the Ferret runtime and FQL options you are using.
