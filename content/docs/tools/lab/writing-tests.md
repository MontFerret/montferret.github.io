---
title: "Writing Tests"
weight: 30
draft: false
description: "Write Lab unit tests and YAML query/assert suites for Ferret scripts."
---

# Writing Tests

Lab runs two kinds of tests: FQL unit tests and YAML suites. Use a plain `.fql` file when the script itself is the test. Use a YAML suite when you want to separate the query from the assertion.

## Write a FQL unit test

A `.fql` test passes when Ferret executes the script without returning an error.

{{< code lang="fql" >}}
LET users = [
  { name: "Ada" },
  { name: "Grace" }
]

RETURN T::EQ(LENGTH(users), 2)
{{</ code >}}

Save the file as `users.fql` and run it:

{{< terminal >}}
lab run users.fql
{{< /terminal >}}

Lab does not inspect the returned value for `.fql` tests. If the script returns `false` without a runtime error, Lab still treats the unit test as passed. Use assertion helpers such as `T::EQ` when a mismatch should fail the test.

## Write an expected-failure test

A file ending in `.fail.fql` passes only when the runtime returns an error.

{{< code lang="fql" >}}
RETURN MISSING_FUNCTION()
{{</ code >}}

Save this as `invalid.fail.fql`. Lab will fail the test if the script unexpectedly succeeds.

Expected-failure tests are useful for validating syntax, runtime errors, or module behavior that should reject a script.

## Write a YAML suite

A YAML suite runs a `query` script first, stores its JSON result, and then runs an `assert` script.

```yaml
query:
  text: |
    RETURN [
      { name: "Ada" },
      { name: "Grace" }
    ]

assert:
  text: |
    LET result = @lab.data.query.result
    RETURN T::EQ(LENGTH(result), 2)
```

The suite passes when the assertion script executes successfully. Use assertion helpers or another expression that raises a runtime error when the expectation is not met. The query result is available in the assertion under `@lab.data.query.result`.

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
    FOR user IN @users
      FILTER user.active
      RETURN user.name
  params:
    users:
      - name: Ada
        active: true
      - name: Grace
        active: false

assert:
  text: |
    RETURN T::EQ(@lab.data.query.result, ["Ada"])
```

Script-level `params` are merged into the user parameters for that script. They can be used with values passed by `lab run --param`.

The assertion can also inspect the parameters used for the query:

{{< code lang="fql" >}}
RETURN T::EQ(@lab.data.query.params.users[0].name, "Ada")
{{</ code >}}

## Set a suite timeout

By default, Lab uses the runner timeout for each test. A YAML suite can override it with `timeout`, in seconds:

```yaml
timeout: 60

query:
  text: |
    RETURN DOCUMENT(@url).title
  params:
    url: "https://example.com"

assert:
  text: |
    RETURN T::NOT::EMPTY(@lab.data.query.result)
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
