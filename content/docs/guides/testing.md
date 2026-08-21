---
title: "Test extraction scripts"
sidebarTitle: "Testing"
weight: 70
draft: false
description: "Use Lab to write and run automated tests for FQL scripts."
---

# Test extraction scripts

[Lab]({{< ref "/docs/tools/lab" >}}) runs FQL test files, applies timeouts and retries, and reports whether each file passed or failed. Assertions come from Ferret's [`t` standard-library namespace]({{< ref "/docs/standard-library/testing" >}}).

The selected Ferret runtime still owns FQL execution and module availability. Lab's builtin runtime provides Ferret core; Lab does not add runtime-specific modules such as `web::html`. To test a script that uses one, select a [binary]({{< ref "/docs/tools/lab/binary-runtime" >}}) or [HTTP]({{< ref "/docs/tools/lab/http-runtime" >}}) runtime that includes it.

## Install Lab

Download the documented Lab release for your platform. On Linux x86_64:

{{< terminal command="true" >}}
mkdir -p "$HOME/.ferret"
curl -fsSL https://github.com/MontFerret/lab/releases/download/v{{< data "versions.lab.v2" >}}/lab_linux_x86_64.tar.gz -o lab.tar.gz
tar -xzf lab.tar.gz
install -m 0755 lab "$HOME/.ferret/lab"
export PATH="$HOME/.ferret:$PATH"
{{< /terminal >}}

Verify both Lab and its selected Ferret runtime:

{{< terminal command="true" >}}
lab version
{{< /terminal >}}

See [Lab Installation]({{< ref "/docs/tools/lab/installation" >}}) for release archives, Docker, and source builds.

## Write a FQL unit test

Any `.fql` file can be a Lab unit test. Assertions return `true` when they succeed and raise an assertion error when they fail.

Create `tests/products.fql`:

{{< editor lang="fql" title="tests/products.fql" >}}
let products = [
  { name: "Mechanical Keyboard", price: 129 },
  { name: "USB-C Dock", price: 89 }
]

t::eq(length(products), 2, "expected two products")
t::not::empty(products[0].name, "the first product needs a name")

return true
{{</ editor >}}

Run the directory:

{{< terminal command="true" >}}
lab run tests/
{{< /terminal >}}

Lab traverses directories recursively and runs `.fql`, `.yaml`, and `.yml` files. Unsupported files are ignored.

### Tests pass on execution, not return values

A `.fql` test passes when the selected runtime executes it without returning an error. Lab does not inspect the returned value, so this file passes:

{{< editor lang="fql" title="tests/false.fql" >}}
return false
{{</ editor >}}

Use a `t` assertion when a false condition should fail the test:

{{< editor lang="fql" title="tests/true.fql" >}}
return t::true(false, "expected the condition to be true")
{{</ editor >}}

## Choose an assertion

The `t` namespace groups assertions by the value or relationship they check:

| Check | Assertions |
| --- | --- |
| Equality and order | `t::eq`, `t::gt`, `t::gte`, `t::lt`, `t::lte` |
| Size and content | `t::empty`, `t::len`, `t::include`, `t::match` |
| Exact values | `t::true`, `t::false`, `t::none` |
| Value types | `t::string`, `t::int`, `t::float`, `t::datetime`, `t::array`, `t::object`, `t::binary` |
| Explicit failure | `t::fail` |

Except for `t::fail`, each assertion has a negated form under `t::not`. For example:

{{< editor lang="fql" title="tests/negated-assertions.fql" >}}
let status = "ready"
let items = ["Mechanical Keyboard"]
let tags = ["stable"]

t::not::eq(status, "failed")
t::not::empty(items)
t::not::include(tags, "deprecated")
{{</ editor >}}

Assertions accept an optional message as their final argument. Use it to explain the expectation in the failure output. See [Testing]({{< ref "/docs/standard-library/testing" >}}) for every signature and accepted value type.

## Test an expected runtime error

Use a YAML suite with `expect.error` when a query must fail. Save this as `tests/rejected-input.yaml`:

```yaml
query:
  text: |
    return 1 / 0

expect:
  error:
    contains: "division by zero"
```

Lab passes the test only when the runtime returns an error containing the configured text. Matching a stable part of the message prevents an unrelated runtime failure from accidentally satisfying the test. Only `contains` is supported inside `expect.error`; unknown fields fail during suite construction instead of falling back to an unqualified error expectation. Use an empty object when any runtime error is sufficient:

```yaml
expect:
  error: {}
```

The test fails if the query completes successfully. Do not define `assert` together with `expect.error`; an expected query failure produces no result for an assertion script.

The older `.fail.fql` filename convention remains supported for backward compatibility. It passes on any runtime error and fails when execution succeeds, but Lab emits a deprecation warning. Prefer `expect.error` for new negative tests; Lab does not rewrite legacy files automatically.

## Separate a query from its assertions

Use a YAML suite when the extraction and its assertions should be separate scripts. The `query` runs first; its JSON result is available to the `assert` script as `@lab.data.query.result`.

```yaml
# tests/products.yaml
query:
  text: |
    let products = [
      { name: "Mechanical Keyboard", price: 129 },
      { name: "USB-C Dock", price: 89 }
    ]

    return products[*].name

assert:
  text: |
    let names = @lab.data.query.result

    t::len(names, 2, "expected two product names")
    return t::include(names, "Mechanical Keyboard")
```

Each `query` and `assert` block must define exactly one of `text` or `ref`. Use `ref` to load FQL from another file, `params` to add user parameters for that script, and the suite's top-level `timeout` field to override the command timeout in seconds. See [Writing Tests]({{< ref "/docs/tools/lab/writing-tests" >}}) for the full suite format.

## Serve local fixtures

Stable fixtures avoid test failures caused by live APIs changing or becoming unavailable. Create this structure:

```text
tests/
├── fixtures/
│   └── products.json
└── fixture-products.fql
```

Add a small JSON response to `tests/fixtures/products.json`:

```json
[
  { "name": "Mechanical Keyboard", "price": 129 },
  { "name": "USB-C Dock", "price": 89 }
]
```

Keep the script production-oriented by reading its normal `@baseUrl` parameter:

{{< code lang="fql" title="tests/fixture-products.fql" >}}
let response = io::net::http::get(@baseUrl + "/products.json")
let products = json_parse(to_string(response))

t::len(products, 2, "expected two fixture products")
return t::gt(products[0].price, 0, "price must be positive")
{{</ code >}}

Start the fixture server for the duration of the test run, then bind the generated endpoint to `@baseUrl`:

{{< terminal command="true" >}}
lab run tests/ \
  --serve ./tests/fixtures@fixtures \
  --param-bind baseUrl=@lab.static.fixtures \
  --policy-http-allow-localhost
{{< /terminal >}}

Lab assigns a free port automatically and resolves the binding after that URL exists. The localhost policy flag permits the builtin runtime's HTTP client to fetch the local fixture. Scripts that deliberately depend on Lab can use `@lab.static.fixtures` directly. See [Static File Server]({{< ref "/docs/tools/lab/static-serving" >}}) for fixed ports, multiple directories, and runtimes outside the Lab host.

## Control the test run

Use runner options when a suite needs parallelism, repeated passes, or retries:

{{< terminal command="true" >}}
lab run tests/ \
  --concurrency=4 \
  --timeout=30 \
  --attempts=2 \
  --times=3 \
  --reporter=simple
{{< /terminal >}}

- `--concurrency=4` runs up to four test files at once.
- `--timeout=30` sets the per-test timeout to 30 seconds.
- `--attempts=2` allows at most two attempts for a failed required run.
- `--times=3` requires three successful runs of each test.
- `--reporter=simple` produces line-oriented output suited to CI logs.

See [Runners]({{< ref "/docs/tools/lab/runners" >}}) for retry and repetition behavior.

## Run tests in CI

Install Lab and use the simple reporter in GitHub Actions:

```yaml
# .github/workflows/test.yml
name: FQL tests

on:
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Lab
        run: |
          mkdir -p "$HOME/.ferret"
          curl -fsSL https://github.com/MontFerret/lab/releases/download/v{{< data "versions.lab.v2" >}}/lab_linux_x86_64.tar.gz -o lab.tar.gz
          tar -xzf lab.tar.gz
          install -m 0755 lab "$HOME/.ferret/lab"
          echo "$HOME/.ferret" >> "$GITHUB_PATH"

      - name: Run FQL tests
        run: lab run tests/ --reporter=simple
```

The release URL follows the version selected for this site. For waits, fixture services, Docker, and external runtimes, see [CI]({{< ref "/docs/tools/lab/ci" >}}).

## Next steps

{{< docs-related tiles="tools-lab,tools-lab-writing-tests,tools-lab-ci,guide-error-handling" >}}
