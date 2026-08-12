---
title: "Waitfor Expressions"
sidebarTitle: "Waitfor"
weight: 60
draft: false
description: "Suspend evaluation until a condition becomes true or an event arrives with the waitfor expression."
---

# Waitfor Expressions

A `waitfor` expression pauses a query until something happens, then resumes and produces a value. It can wait for one condition or event, or synchronize a group with `any { ... }` or `all { ... }`.

## Waiting for a condition

In its condition (predicate) mode, `waitfor` re-checks an ordinary expression until it is satisfied or a timeout is reached. This mode is a pure language construct — it needs no special capability.

There are four forms:

| Form | Waits until | Returns |
| --- | --- | --- |
| `waitfor <expr>` | the expression is true | `true`, or `false` on timeout |
| `waitfor exists <expr>` | the expression is non-empty according to `exists` semantics | `true`, or `false` on timeout |
| `waitfor not exists <expr>` | the expression is empty or absent | `true`, or `false` on timeout |
| `waitfor value <expr>` | the expression yields anything other than `none` | that value, or `none` on timeout |

{{< editor lang="fql" >}}
return waitfor exists [1, 2, 3] timeout 100ms
{{</ editor >}}

`waitfor value` treats empty strings, arrays, and objects as valid values and returns them immediately. When you need the candidate itself but also require it to be non-empty, add a `when` condition:

{{< code lang="fql" >}}
return waitfor value loadItems()
    when LENGTH(.) > 0
    timeout 5s
{{</ code >}}

The expression is re-evaluated on each attempt, so it can reflect state that changes over time. When the wait runs out, the result reports the timeout rather than raising an error:

{{< editor lang="fql" >}}
return waitfor false timeout 50ms every 10ms
{{</ editor >}}

### Synchronizing condition groups

Add `any` or `all` before a block to evaluate several conditions as one polling operation. The mode before `any` or `all` applies to every entry:

{{< code lang="fql" >}}
let firstReady = waitfor any {
    queue.ready
    worker.ready
}
timeout 5s

let bothValues = waitfor value all {
    queue.value
    worker.value
}
timeout 5s
{{</ code >}}

The grouped forms return:

| Form | Success condition | Result |
| --- | --- | --- |
| `waitfor any { ... }` | at least one expression is true | `true` |
| `waitfor all { ... }` | every expression is true in the same polling cycle | `true` |
| `waitfor exists any/all { ... }` | one/all expressions satisfy `exists` | `true` |
| `waitfor not exists any/all { ... }` | one/all expressions satisfy `not exists` | `true` |
| `waitfor value any { ... }` | at least one expression is not `none` | the first qualifying value in declaration order |
| `waitfor value all { ... }` | every expression is not `none` in the same polling cycle | an array of values in declaration order |

Polling groups synchronize **state**. `all` does not remember an entry that passed in an earlier cycle: every entry must pass together during one cycle. `any` and `all` evaluate entries in declaration order and stop evaluating the current cycle as soon as the outcome is known.

Each entry may have its own repeated `when` conditions. Inside them, `.` is that entry's candidate:

{{< code lang="fql" >}}
return waitfor value any {
    primaryStatus()
        when .ready
        when .healthy
    fallbackStatus()
        when .ready
}
timeout 10s
every 100ms
{{</ code >}}

`timeout`, `every`, `backoff`, `jitter`, `on timeout`, and `on error` belong to the whole group. A group has one timeout and one polling schedule; entries cannot define separate policies.

### Tuning the wait

Several clauses control how the wait behaves:

- `timeout <duration>` — the maximum time to wait, provided as a value coercible to Duration.
- `every <interval>` — how often to re-check. A second coercible Duration, `every <interval>, <cap>`, caps how large the interval can grow. Without this clause, polling defaults to `100ms`.
- `backoff LINEAR | EXPONENTIAL | none` — how the interval between checks grows over time.
- `jitter <0..1>` — randomizes the interval to avoid synchronized retries.
- `when <condition>` — an additional condition that must also hold; the candidate value is available as `.`.

{{< code lang="fql" >}}
waitfor value loadStatus()
    timeout 10s
    every 100ms, 1s
    backoff EXPONENTIAL
    jitter 0.2
{{</ code >}}

`timeout`, `every`, and its cap accept ordinary expressions, so values can be stored or computed:

{{< code lang="fql" >}}
let base = 50ms

return waitfor value loadStatus()
    timeout base * 20
    every base, base * 4
{{</ code >}}

`timeout`, `every`, and its cap use the canonical Duration conversion rules. Numbers are milliseconds, duration strings may be compound, and singleton lists are converted recursively:

{{< code lang="fql" >}}
let timeout = "1s500ms"

return waitfor false
    timeout timeout
    every 25, [100]
{{</ code >}}

All scheduling results must be non-negative. Conversion failures, overflow, and negative values raise runtime errors.

### Recovering from a timeout

By default a timed-out wait returns `false` (or `none` for the `value` form). A recovery clause lets you choose a different result.

{{< editor lang="fql" >}}
return waitfor value none timeout 30ms every 5ms on timeout return "gave up"
{{</ editor >}}

Use `on error` to handle a failure raised while evaluating the condition.

## Waiting for an event

In event mode, `waitfor` subscribes to an event source and waits for a matching event, which it returns as a value.

{{< code lang="fql" >}}
let event = waitfor event "navigation" in page timeout 5s
{{</ code >}}

Event timeouts use the same Duration conversion and non-negative scheduling policy as condition waits.

A `when` filter accepts only events that match a condition. Inside the filter, the incoming event is available as `.`. Multiple `when` clauses must all pass.

{{< code lang="fql" >}}
waitfor event "message" in socket
    when .type == "data"
    timeout 5s
{{</ code >}}

### Synchronizing event occurrences

Event groups use the existing event entry syntax and may subscribe to different sources:

{{< code lang="fql" >}}
let event = waitfor event any {
    "navigation" in page
    "download" in browser
    "disconnect" in socket
}
timeout 10s
{{</ code >}}

`event any` returns the first qualifying event that occurs and closes the other subscriptions. Concurrent events have no declaration-order tie-break. If one source ends without a qualifying event, the other arms keep waiting; if every source ends without a match, the result is `none`, as with a singular event wait.

`event all` waits until every subscription has produced one qualifying event. Each arm remains satisfied after its event occurs, and the result is an array in declaration order even when the events arrive in another order:

{{< code lang="fql" >}}
let events = waitfor event all {
    "domcontentloaded" in page
    "networkidle" in page
}
timeout 10s
{{</ code >}}

Every `event all` arm must produce an event that passes its `when` filters. If an unmatched source ends, the wait fails immediately because the group can no longer be satisfied. A source ending after its arm has matched does not affect the result.

Event groups synchronize **occurrences**, unlike polling groups, which synchronize state in one cycle. Event filters are per entry:

{{< code lang="fql" >}}
return waitfor event any {
    "response" in page
        when .status >= 500
    "dialog" in page
        when .type == "confirm"
}
timeout 10s
{{</ code >}}

All subscriptions are established concurrently. A timeout, cancellation, setup failure, trigger failure, stream error, or completed wait closes every remaining subscription. `timeout`, `trigger`, `on timeout`, and `on error` apply once to the whole group.

### Triggering the event

A `trigger` clause runs statements *after* the subscription is set up but *before* waiting begins. For an event group, every subscription is established before the one shared trigger runs. This is how you cause the event you are waiting for without risking a race where it fires before you start listening.

{{< code lang="fql" >}}
waitfor event "response" in page
    trigger ( button <- "click" )
    timeout 10s
{{</ code >}}

The trigger body can dispatch events and run other statements. See [Dispatch Expressions]({{< ref "dispatch" >}}).

## A host capability

Event mode only works when the source is an **observable** value — one that produces a stream of events, provided by a module or the host application, such as a browser page. Condition mode has no such requirement. See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Next steps

{{< docs-related tiles="language-control-flow-dispatch,language-types-capabilities,language-control-flow" >}}
