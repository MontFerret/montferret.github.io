---
title: "Configuration"
weight: 100
draft: false
description: "Reference Lab command flags, defaults, environment variables, and runtime params."
---

# Configuration

Lab configuration is primarily command-line flags. Most flags also have `LAB_*` environment variable bindings.

The current Lab v2 command surface does not include `--cdp` or `LAB_CDP`. Browser/CDP behavior belongs to the selected Ferret runtime, not to Lab's CLI flags.

## `lab run`

`lab run` executes FQL unit tests and YAML suites.

| Flag | Alias | Env var | Default | Meaning |
| --- | --- | --- | --- | --- |
| `--files` | `-f` | `LAB_FILES` | | Location of test files to run. Can be repeated. |
| `--timeout` | `-t` | `LAB_TIMEOUT` | `30` | Per-test timeout in seconds. |
| `--reporter` | | `LAB_REPORTER` | `console` | Reporter: `console` or `simple`. |
| `--runtime` | `-r` | `LAB_RUNTIME` | | Runtime selector. Empty means builtin. Use `http(s)://...` or `bin:...` for adapters. |
| `--runtime-param` | `--rp` | `LAB_RUNTIME_PARAM` | | Runtime adapter parameter. Can be repeated. |
| `--concurrency` | `-c` | `LAB_CONCURRENCY` | `1` | Number of tests to run at once. |
| `--times` | | `LAB_TIMES` | `1` | Required successful runs per test. |
| `--attempts` | `-a` | `LAB_ATTEMPTS` | `1` | Maximum attempts per required run. |
| `--times-interval` | | `LAB_TIMES_INTERVAL` | `0` | Seconds to wait between repeated runs or retry attempts. |
| `--serve` | | `LAB_SERVE` | | Static directory mapping for test execution. Can be repeated. |
| `--mock` | | `LAB_MOCK` | | OpenAPI mock API spec mapping for test execution. Can be repeated. |
| `--serve-bind` | | `LAB_SERVE_BIND` | | Host to bind local services to. Must not include a port. |
| `--serve-host` | | `LAB_SERVE_HOST` | | Host to advertise in local service URLs. Must not include a port. |
| `--param` | `-p` | `LAB_PARAM` | | User parameter for FQL. Can be repeated. |
| `--wait` | `-w` | `LAB_WAIT` | | HTTP resource to wait for before tests run. Can be repeated. |
| `--wait-timeout` | `--wt` | `LAB_WAIT_TIMEOUT` | `5` | Wait interval in seconds. |
| `--wait-attempts` | | `LAB_WAIT_ATTEMPTS` | `5` | Number of wait attempts. |

Example:

{{< terminal >}}
lab run tests/ \
  --timeout=60 \
  --concurrency=4 \
  --attempts=3 \
  --reporter=simple \
  --param='baseUrl:"https://example.com"'
{{< /terminal >}}

## `lab serve`

`lab serve` starts static and mock API services without running tests.

| Flag | Env var | Default | Meaning |
| --- | --- | --- | --- |
| `--static` | `LAB_STATIC` | | Static directory mapping. Can be repeated. |
| `--mock` | `LAB_MOCK` | | OpenAPI mock API spec mapping. Can be repeated. |
| `--serve-bind` | `LAB_SERVE_BIND` | | Host to bind local services to. Must not include a port. |
| `--serve-host` | `LAB_SERVE_HOST` | | Host to advertise in local service URLs. Must not include a port. |

Example:

{{< terminal >}}
lab serve --static ./dist@app --mock ./users.yaml@api
{{< /terminal >}}

Standalone service entries must use `--static` or `--mock`.

For service-specific examples, see [Static File Server]({{< ref "static-serving" >}}) and [Mock API Server]({{< ref "mock-api" >}}).

## `lab version`

`lab version` prints the Lab version and the selected runtime version.

| Flag | Alias | Env var | Default | Meaning |
| --- | --- | --- | --- | --- |
| `--runtime` | `-r` | `LAB_RUNTIME` | | Runtime selector to inspect. Empty means builtin. |

Examples:

{{< terminal >}}
lab version
lab version --runtime=http://localhost:8080
lab version --runtime=bin:/usr/local/bin/ferret
{{< /terminal >}}

## Parameter syntax

`--param` and `--runtime-param` use `name:value` syntax. The value must be valid JSON.

{{< terminal >}}
lab run tests/ \
  --param='name:"Ada"' \
  --param='limit:10' \
  --param='active:true' \
  --param='tags:["admin","editor"]'
{{< /terminal >}}

Strings must be JSON strings. Use shell quoting to keep the inner quotes intact.

## Runtime param keys

HTTP runtimes support:

| Key | Type | Meaning |
| --- | --- | --- |
| `headers` | object of strings | Extra request headers. |
| `path` | string | Run endpoint override. |
| `cookies` | object of strings | Cookies sent to the runtime. |

Binary runtimes support:

| Key | Type | Meaning |
| --- | --- | --- |
| `flags` | array of strings | Raw command-line flags passed to the binary before generated `--param` values. |
| Any other key | JSON value | Shared parameter forwarded as `--param=<key>:<json>`. |

## Local service entry syntax

`--serve`, `--static`, and `--mock` entries use:

| Syntax | Meaning |
| --- | --- |
| `<path>` | Dynamic port and default alias. |
| `<path>:<port>` | Fixed port and default alias. |
| `<path>@<alias>` | Explicit alias and dynamic port. |
| `<path>@<alias>:<port>` | Explicit alias and fixed port. |

Static entries must point to existing directories. Mock entries must point to existing files.

## Next steps

{{< docs-related tiles="tools-lab-running-tests,tools-lab-writing-tests,tools-lab" >}}
