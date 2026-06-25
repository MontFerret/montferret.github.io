---
title: "Error Handling"
sidebarTitle: "Error Handling"
weight: 65
draft: false
description: "Recover from runtime failures with recovery tails and optional chaining."
---

# Error Handling

FQL does not use try/catch blocks. Instead, a **recovery tail** attaches directly to the expression that might fail. The tail tells the runtime what to do when a failure occurs — return a fallback value, retry the operation, or propagate the error.

For member access on values that may be `NONE`, FQL provides the **optional chaining operator** (`?.`), which returns `NONE` instead of failing.

## Returning a fallback value

The most common recovery is `ON ERROR RETURN`, which catches a runtime failure and produces a fallback value instead.

{{< code lang="fql" >}}
RETURN getData() ON ERROR RETURN NONE
{{</ code >}}

The fallback expression is only evaluated when the guarded expression fails. Any expression can serve as the fallback — a literal, a variable, a function call, or a collection.

{{< code lang="fql" >}}
LET rows = QUERY `.items` IN doc ON ERROR RETURN []
{{</ code >}}

{{< code lang="fql" >}}
DISPATCH "click" IN element ON ERROR RETURN NONE
{{</ code >}}

## Propagating errors

By default, a runtime failure propagates up and halts evaluation. `ON ERROR FAIL` makes that behavior explicit. It is useful when a recovery tail is required for readability or when paired with a separate `ON TIMEOUT` clause.

{{< code lang="fql" >}}
LET value = WAITFOR VALUE check()
    TIMEOUT 5s
    ON TIMEOUT RETURN NONE
    ON ERROR FAIL
{{</ code >}}

## Retrying on failure

`ON ERROR RETRY` re-executes the guarded expression a given number of times before giving up.

{{< code lang="fql" >}}
RETURN fetchData() ON ERROR RETRY 3
{{</ code >}}

The retry count is the number of *additional* attempts after the first failure. If all retries are exhausted and the expression still fails, the final error propagates.

### Delay and backoff

A `DELAY` clause adds a pause between retries. An optional `BACKOFF` strategy controls how the delay grows.

{{< code lang="fql" >}}
RETURN fetchData() ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL
{{</ code >}}

| Strategy | Behavior |
| --- | --- |
| `CONSTANT` | every retry waits the same duration |
| `LINEAR` | the delay grows by a fixed increment each retry |
| `EXPONENTIAL` | the delay doubles each retry |

`BACKOFF` requires `DELAY`. Without `BACKOFF`, the delay is constant.

### Fallback after retries

When all retries are exhausted, `OR RETURN` provides a fallback value instead of propagating the final error. `OR FAIL` makes propagation explicit.

{{< code lang="fql" >}}
RETURN fetchData() ON ERROR RETRY 3 DELAY 100ms BACKOFF EXPONENTIAL OR RETURN "unavailable"
{{</ code >}}

{{< code lang="fql" >}}
RETURN fetchData() ON ERROR RETRY 2 OR FAIL
{{</ code >}}

## Handling timeouts

`ON TIMEOUT` handles timeout failures separately from other errors. It is only valid on `WAITFOR` expressions that include a `TIMEOUT` clause.

{{< code lang="fql" >}}
LET result = WAITFOR VALUE loadStatus()
    TIMEOUT 10s
    ON TIMEOUT RETURN "timed out"
{{</ code >}}

`ON ERROR` and `ON TIMEOUT` are independent — they can appear together on the same expression, each with its own action.

{{< code lang="fql" >}}
LET token = WAITFOR VALUE authenticate()
    TIMEOUT 5s
    ON TIMEOUT RETURN "timeout"
    ON ERROR RETRY 2 DELAY 100ms OR RETURN "error"
{{</ code >}}

A timeout is not retried by `ON ERROR RETRY`. The two handlers apply to different failure kinds.

## Grouped expressions

Any expression can be wrapped in parentheses to attach a recovery tail. This is how you add error recovery to constructs that do not accept recovery tails directly, such as `FOR` loops.

{{< code lang="fql" >}}
LET results = (FOR item IN items
    RETURN process(item)
) ON ERROR RETURN []
{{</ code >}}

Retry works on grouped expressions too. When a grouped `FOR` is retried, the loop restarts from the beginning — partial results from a failed attempt are discarded.

{{< code lang="fql" >}}
LET results = (FOR item IN items
    RETURN process(item)
) ON ERROR RETRY 1 OR RETURN []
{{</ code >}}

## Optional chaining

The optional chaining operator `?.` accesses a member on a value that may be `NONE`. Instead of failing, it produces `NONE`.

{{< editor lang="fql" >}}
LET obj = NONE

RETURN obj?.name
{{</ editor >}}

It works with computed property names as well.

{{< editor lang="fql" >}}
LET obj = NONE
LET key = "name"

RETURN obj?.[key]
{{</ editor >}}

Without `?.`, accessing a member on `NONE` is a runtime error.

Optional chaining applies only to member access. `func?()` and `arr?[0]` are not supported — use `ON ERROR RETURN` or a grouped expression instead.

## Where recovery applies

| Construct | Recovery tails | Optional chaining |
| --- | --- | --- |
| Function calls | `func() ON ERROR ...` | — |
| Member access | `obj.prop ON ERROR ...` | `obj?.prop` |
| `QUERY` | `QUERY ... ON ERROR ...` | — |
| `DISPATCH` | `DISPATCH ... ON ERROR ...` | — |
| `WAITFOR` | `ON ERROR ...`, `ON TIMEOUT ...` | — |
| Grouped `(...)` | `(...) ON ERROR ...` | — |

Each expression may define `ON ERROR` at most once. `WAITFOR` expressions may additionally define `ON TIMEOUT` at most once. `RETRY` and its `OR` fallback are only available under `ON ERROR`.
