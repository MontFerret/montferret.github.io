---
title: "Runners"
weight: 60
draft: false
description: "Control Lab concurrency, retries, repeated runs, test timeouts, and reporter output."
---

# Runners

The Lab runner coordinates source files, test cases, runtime execution, and reporting. Its main controls are concurrency, attempts, repeated successful runs, intervals, and per-test timeouts.

## Run tests concurrently

Use `--concurrency` to control how many test files Lab can execute at once.

{{< terminal >}}
lab run --concurrency=4 tests/
{{< /terminal >}}

The default is `1`. If `0` is provided, Lab falls back to `1`.

Each test receives a cloned parameter set, so Lab system parameters and suite data are isolated per test file.

## Set the test timeout

Use `--timeout` to set the per-test timeout in seconds.

{{< terminal >}}
lab run --timeout=60 tests/
{{< /terminal >}}

The default is `30` seconds. YAML suites can override this with a top-level `timeout` field.

## Retry failed tests

Use `--attempts` to set the maximum number of attempts for a test.

{{< terminal >}}
lab run --attempts=3 tests/
{{< /terminal >}}

`--attempts=3` means Lab may run a failing test up to three times total. A successful attempt stops the retry loop for the current required run.

The default is `1`, which means no retry after the first failure.

## Repeat successful runs

Use `--times` to require each test to pass multiple times.

{{< terminal >}}
lab run --times=5 tests/
{{< /terminal >}}

The default is `1`. Lab increments the run count only after a successful execution. If a run fails, Lab retries according to `--attempts`; if attempts are exhausted, the test fails.

Use `--times-interval` to wait between repeated runs or retry attempts:

{{< terminal >}}
lab run --times=5 --times-interval=2 tests/
{{< /terminal >}}

The interval is measured in seconds.

## Pick a reporter

Use `--reporter` to choose output formatting.

| Reporter | Use when |
| --- | --- |
| `console` | You want human-readable structured console output. |
| `simple` | You want line-oriented output for scripts and CI logs. |

Example:

{{< terminal >}}
lab run --reporter=simple tests/
{{< /terminal >}}

The simple reporter emits progress lines and a final summary:

```text
PASS file="tests/users.fql" duration=12ms attempts=1 times=1
FAIL file="tests/api.yaml" duration=30s attempts=3 times=1 error="operation timed out"
DONE passed=1 failed=1 duration=30.012s
```

## Result fields

Each reported test result includes:

| Field | Meaning |
| --- | --- |
| `file` | Source file name. |
| `duration` | Average duration for successful runs. If no run succeeds, it reflects accumulated failed attempt time. |
| `attempts` | Number of attempts made. |
| `times` | Number of successful required runs completed. |
| `error` | Runtime, source, timeout, or assertion error for failed tests. |

The summary includes total passed tests, failed tests, and total runner duration.

## Exit status

Reporters return an error when the summary contains failures. The Lab process exits with status `1` for failed test runs and status `0` when all tests pass.

## Next steps

{{< docs-related tiles="tools-lab-running-tests,tools-lab-configuration,tools-lab" >}}
