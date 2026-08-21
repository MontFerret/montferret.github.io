---
title: "Writing Tests"
weight: 30
draft: false
description: "Write Lab unit tests and YAML assertion or expected-error suites for Ferret scripts."
---

# Writing Tests

Lab runs two kinds of tests: FQL unit tests and YAML suites. Use a plain `.fql` file when the script itself is the test. Use a YAML suite when you want to separate the query from an assertion or require a runtime error.

## Write a FQL unit test

A `.fql` test passes when Ferret executes the script without returning an error.

{{< code lang="fql" >}}
let users = [
  { name: "Ada" },
  { name: "Grace" }
]

return t::eq(length(users), 2)
{{</ code >}}

Save the file as `users.fql` and run it:

{{< terminal >}}
lab run users.fql
{{< /terminal >}}

Lab does not inspect the returned value for `.fql` tests. If the script returns `false` without a runtime error, Lab still treats the unit test as passed. Use assertion helpers such as `t::eq` when a mismatch should fail the test.

## Write a YAML suite

A YAML suite runs a `query` script first, stores its JSON result, and then runs an `assert` script.

```yaml
query:
  text: |
    return [
      { name: "Ada" },
      { name: "Grace" }
    ]

assert:
  text: |
    let result = @lab.data.query.result
    return t::eq(length(result), 2)
```

The suite passes when the assertion script executes successfully. Use assertion helpers or another expression that raises a runtime error when the expectation is not met. The query result is available in the assertion under `@lab.data.query.result`.

## Expect a runtime error

Use a YAML suite with `expect.error` when the query must fail. An empty error object accepts any error returned by the runtime:

```yaml
query:
  text: |
    return 1 / 0

expect:
  error: {}
```

The test fails if the query completes successfully. Add `contains` when the failure reason has a stable message:

```yaml
query:
  text: |
    return 1 / 0

expect:
  error:
    contains: "division by zero"
```

Lab performs a substring match against the runtime error message. Matching a stable part of the message prevents an unrelated runtime failure from accidentally passing the test. Only `contains` is supported inside `expect.error`; unknown fields fail during suite construction instead of falling back to an unqualified error expectation. Do not define `assert` together with `expect.error`; an expected query failure produces no result for an assertion script.

The older `.fail.fql` filename convention remains supported for backward compatibility. It passes when the runtime returns any error and fails when execution succeeds, but Lab emits a deprecation warning. Prefer `expect.error` for new negative tests because it makes the expectation explicit and can verify the failure reason.

## Use inline scripts or refs

Each `query` and `assert` block must use exactly one of:

| Field | Meaning |
| --- | --- |
| `text` | Inline FQL source in the YAML file. |
| `ref` | A source reference to another test file. |

Example with referenced scripts:

```yaml
query:
  ref: ./queries/users.fql

assert:
  ref: ./assertions/users-count.fql
```

References are resolved by the source that loaded the suite. For local files, relative refs are resolved from the suite's directory.

## Pass suite parameters

The `params` field adds user parameters for that script. Parameters are available in FQL as `@name`.

```yaml
query:
  text: |
    return for user in @users
      filter user.active
      return user.name
  params:
    users:
      - name: Ada
        active: true
      - name: Grace
        active: false

assert:
  text: |
    return t::eq(@lab.data.query.result, ["Ada"])
```

Script-level `params` are merged into the user parameters for that script. They can be used with values passed by `lab run --param`. The `--param-bind` option can adapt a value that already exists during Lab setup into an ordinary user parameter path before scripts run.

The assertion can also inspect the parameters used for the query:

{{< code lang="fql" >}}
return t::eq(@lab.data.query.params.users[0].name, "Ada")
{{</ code >}}

## Set a suite timeout

By default, Lab uses the runner timeout for each test. A YAML suite can override it with `timeout`, in seconds:

```yaml
timeout: 60

query:
  text: |
    return web::html::open(@url).title
  params:
    url: "https://example.com"

assert:
  text: |
    return t::not::empty(@lab.data.query.result)
```

Timeouts are enforced around the suite run. Use this for tests that legitimately need more time than the command default.

## Use Lab system parameters

Lab reserves the `@lab` parameter namespace for values it creates during a test run.

| Parameter | Available when |
| --- | --- |
| `@lab.static.<alias>` | A static service is started with `--serve`. |
| `@lab.mock.<alias>` | A mock API service is started with `--mock`. |
| `@lab.data.query.result` | A YAML suite assertion runs after the query. |
| `@lab.data.query.params` | A YAML suite assertion runs after the query. |

Keep user parameters outside `@lab`; Lab writes that namespace for each test.
Use direct `@lab.*` access when a script intentionally depends on Lab, or use `--param-bind` to keep the script's parameter names environment-neutral.

## Next steps

{{< docs-related tiles="tools-lab-running-tests,tools-lab-configuration,stdlib" >}}
