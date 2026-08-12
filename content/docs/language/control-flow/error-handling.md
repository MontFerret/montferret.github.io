---
title: "Error Handling"
sidebarTitle: "Error Handling"
weight: 65
draft: false
description: "Recover from runtime failures with recovery tails and optional chaining."
---

# Error Handling

FQL does not use try/catch blocks. Instead, a **recovery tail** attaches directly to the expression that might fail. The tail tells the runtime what to do when a failure occurs — return a fallback value, retry the operation, or propagate the error.

For member access on values that may be `none`, FQL provides the **optional chaining operator** (`?.`), which returns `none` instead of failing.

## Returning a fallback value

The most common recovery is `on error return`, which catches a runtime failure and produces a fallback value instead.

{{< code lang="fql" >}}
return getData() on error return none
{{</ code >}}

The fallback expression is only evaluated when the guarded expression fails. Any expression can serve as the fallback — a literal, a variable, a function call, or a collection.

{{< code lang="fql" >}}
let rows = query `.items` in doc on error return []
{{</ code >}}

{{< code lang="fql" >}}
dispatch "click" in element on error return none
{{</ code >}}

## Propagating errors

By default, a runtime failure propagates up and halts evaluation. `on error fail` makes that behavior explicit. It is useful when a recovery tail is required for readability or when paired with a separate `on timeout` clause.

{{< code lang="fql" >}}
let value = waitfor value check()
    timeout 5s
    on timeout return none
    on error fail
{{</ code >}}

## Retrying on failure

`on error retry` re-executes the guarded expression a given number of times before giving up.

{{< code lang="fql" >}}
return fetchData() on error retry 3
{{</ code >}}

The retry count is the number of *additional* attempts after the first failure. If all retries are exhausted and the expression still fails, the final error propagates.

### Delay and backoff

A `delay` clause adds a pause between retries. An optional `backoff` strategy controls how the delay grows.

{{< code lang="fql" >}}
return fetchData() on error retry 3 delay 100ms backoff EXPONENTIAL
{{</ code >}}

| Strategy | Behavior |
| --- | --- |
| `constant` | every retry waits the same duration |
| `linear` | the delay grows by a fixed increment each retry |
| `exponential` | the delay doubles each retry |

`backoff` requires `delay`. Without `backoff`, the delay is constant.

`delay` accepts any value supported by the canonical Duration conversion. Numbers are milliseconds, duration strings may be compound, and singleton lists are converted recursively. The converted delay must be non-negative; conversion failures, overflow, and negative values raise runtime errors.

Because the unparenthesized `or` token begins the retry fallback, wrap a logical `or` used inside the delay expression in parentheses:

{{< code lang="fql" >}}
let base = 100
let preferredDelay = none

return fetchData()
    on error retry 3 delay base * 2 or return none

return fetchData()
    on error retry 3 delay (preferredDelay or base) or return none
{{</ code >}}

### Fallback after retries

When all retries are exhausted, `or return` provides a fallback value instead of propagating the final error. `or fail` makes propagation explicit.

{{< code lang="fql" >}}
return fetchData() on error retry 3 delay 100ms backoff EXPONENTIAL or return "unavailable"
{{</ code >}}

{{< code lang="fql" >}}
return fetchData() on error retry 2 or fail
{{</ code >}}

## Handling timeouts

`on timeout` handles timeout failures separately from other errors. It is only valid on `waitfor` expressions that include a `timeout` clause.

{{< code lang="fql" >}}
let result = waitfor value loadStatus()
    timeout 10s
    on timeout return "timed out"
{{</ code >}}

`on error` and `on timeout` are independent — they can appear together on the same expression, each with its own action.

{{< code lang="fql" >}}
let token = waitfor value authenticate()
    timeout 5s
    on timeout return "timeout"
    on error retry 2 delay 100ms or return "error"
{{</ code >}}

A timeout is not retried by `on error retry`. The two handlers apply to different failure kinds.

## Grouped expressions

Any expression can be wrapped in parentheses to attach a recovery tail. This is how you add error recovery to constructs that do not accept recovery tails directly, such as `for` loops.

{{< code lang="fql" >}}
let results = (for item in items {
    return process(item)
}
) on error return []
{{</ code >}}

Retry works on grouped expressions too. When a grouped `for` is retried, the loop restarts from the beginning — partial results from a failed attempt are discarded.

{{< code lang="fql" >}}
let results = (for item in items {
    return process(item)
}
) on error retry 1 or return []
{{</ code >}}

## Optional chaining

The optional chaining operator `?.` accesses a member on a value that may be `none`. Instead of failing, it produces `none`.

{{< editor lang="fql" >}}
let obj = none

return obj?.name
{{</ editor >}}

It works with computed property names as well.

{{< editor lang="fql" >}}
let obj = none
let key = "name"

return obj?.[key]
{{</ editor >}}

Without `?.`, accessing a member on `none` is a runtime error.

Use `??` after optional chaining when the final result should have a concrete fallback:

{{< code lang="fql" >}}
let obj = none

return obj?.name ?? "Unknown"
{{</ code >}}

`??` selects a fallback only after an expression produces `none`; it does not catch runtime errors. Apply `on error`, postfix recovery after the expression (for example, `func()?`), or optional chaining first when a failure should become `none`.

See [none-Coalescing Operator]({{% ref "../operators/coalescing" %}}) for the complete behavior.

Optional chaining applies only to member access. `func?()` and `arr?[0]` are not supported — use `on error return` or a grouped expression instead.

## Where recovery applies

| Construct | Recovery tails | Optional chaining |
| --- | --- | --- |
| Function calls | `func() on error ...` | — |
| Member access | `obj.prop on error ...` | `obj?.prop` |
| `query` | `query ... on error ...` | — |
| `dispatch` | `dispatch ... on error ...` | — |
| `waitfor` | `on error ...`, `on timeout ...` | — |
| Grouped `(...)` | `(...) on error ...` | — |

Each expression may define `on error` at most once. `waitfor` expressions may additionally define `on timeout` at most once. `retry` and its `or` fallback are only available under `on error`.

## Next steps

{{< docs-related tiles="language-control-flow,tools-lab,stdlib" >}}
