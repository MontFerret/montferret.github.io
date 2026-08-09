---
title: "Waitfor Expressions"
sidebarTitle: "Waitfor"
weight: 60
draft: false
description: "Suspend evaluation until a condition becomes true or an event arrives with the WAITFOR expression."
---

# Waitfor Expressions

A `WAITFOR` expression pauses a query until something happens, then resumes and produces a value. It can wait for one condition or event, or synchronize a group with `ANY { ... }` or `ALL { ... }`.

## Waiting for a condition

In its condition (predicate) mode, `WAITFOR` re-checks an ordinary expression until it is satisfied or a timeout is reached. This mode is a pure language construct — it needs no special capability.

There are four forms:

| Form | Waits until | Returns |
| --- | --- | --- |
| `WAITFOR <expr>` | the expression is true | `true`, or `false` on timeout |
| `WAITFOR EXISTS <expr>` | the expression is non-empty according to `EXISTS` semantics | `true`, or `false` on timeout |
| `WAITFOR NOT EXISTS <expr>` | the expression is empty or absent | `true`, or `false` on timeout |
| `WAITFOR VALUE <expr>` | the expression yields anything other than `NONE` | that value, or `NONE` on timeout |

{{< editor lang="fql" >}}
RETURN WAITFOR EXISTS [1, 2, 3] TIMEOUT 100ms
{{</ editor >}}

`WAITFOR VALUE` treats empty strings, arrays, and objects as valid values and returns them immediately. When you need the candidate itself but also require it to be non-empty, add a `WHEN` condition:

{{< code lang="fql" >}}
RETURN WAITFOR VALUE loadItems()
    WHEN LENGTH(.) > 0
    TIMEOUT 5s
{{</ code >}}

The expression is re-evaluated on each attempt, so it can reflect state that changes over time. When the wait runs out, the result reports the timeout rather than raising an error:

{{< editor lang="fql" >}}
RETURN WAITFOR FALSE TIMEOUT 50ms EVERY 10ms
{{</ editor >}}

### Synchronizing condition groups

Add `ANY` or `ALL` before a block to evaluate several conditions as one polling operation. The mode before `ANY` or `ALL` applies to every entry:

{{< code lang="fql" >}}
LET firstReady = WAITFOR ANY {
    queue.ready
    worker.ready
}
TIMEOUT 5s

LET bothValues = WAITFOR VALUE ALL {
    queue.value
    worker.value
}
TIMEOUT 5s
{{</ code >}}

The grouped forms return:

| Form | Success condition | Result |
| --- | --- | --- |
| `WAITFOR ANY { ... }` | at least one expression is true | `true` |
| `WAITFOR ALL { ... }` | every expression is true in the same polling cycle | `true` |
| `WAITFOR EXISTS ANY/ALL { ... }` | one/all expressions satisfy `EXISTS` | `true` |
| `WAITFOR NOT EXISTS ANY/ALL { ... }` | one/all expressions satisfy `NOT EXISTS` | `true` |
| `WAITFOR VALUE ANY { ... }` | at least one expression is not `NONE` | the first qualifying value in declaration order |
| `WAITFOR VALUE ALL { ... }` | every expression is not `NONE` in the same polling cycle | an array of values in declaration order |

Polling groups synchronize **state**. `ALL` does not remember an entry that passed in an earlier cycle: every entry must pass together during one cycle. `ANY` and `ALL` evaluate entries in declaration order and stop evaluating the current cycle as soon as the outcome is known.

Each entry may have its own repeated `WHEN` conditions. Inside them, `.` is that entry's candidate:

{{< code lang="fql" >}}
RETURN WAITFOR VALUE ANY {
    primaryStatus()
        WHEN .ready
        WHEN .healthy
    fallbackStatus()
        WHEN .ready
}
TIMEOUT 10s
EVERY 100ms
{{</ code >}}

`TIMEOUT`, `EVERY`, `BACKOFF`, `JITTER`, `ON TIMEOUT`, and `ON ERROR` belong to the whole group. A group has one timeout and one polling schedule; entries cannot define separate policies.

### Tuning the wait

Several clauses control how the wait behaves:

- `TIMEOUT <duration>` — the maximum time to wait, provided as a value coercible to Duration.
- `EVERY <interval>` — how often to re-check. A second coercible Duration, `EVERY <interval>, <cap>`, caps how large the interval can grow. Without this clause, polling defaults to `100ms`.
- `BACKOFF LINEAR | EXPONENTIAL | NONE` — how the interval between checks grows over time.
- `JITTER <0..1>` — randomizes the interval to avoid synchronized retries.
- `WHEN <condition>` — an additional condition that must also hold; the candidate value is available as `.`.

{{< code lang="fql" >}}
WAITFOR VALUE loadStatus()
    TIMEOUT 10s
    EVERY 100ms, 1s
    BACKOFF EXPONENTIAL
    JITTER 0.2
{{</ code >}}

`TIMEOUT`, `EVERY`, and its cap accept ordinary expressions, so values can be stored or computed:

{{< code lang="fql" >}}
LET base = 50ms

RETURN WAITFOR VALUE loadStatus()
    TIMEOUT base * 20
    EVERY base, base * 4
{{</ code >}}

`TIMEOUT`, `EVERY`, and its cap use the canonical Duration conversion rules. Numbers are milliseconds, duration strings may be compound, and singleton lists are converted recursively:

{{< code lang="fql" >}}
LET timeout = "1s500ms"

RETURN WAITFOR FALSE
    TIMEOUT timeout
    EVERY 25, [100]
{{</ code >}}

All scheduling results must be non-negative. Conversion failures, overflow, and negative values raise runtime errors.

### Recovering from a timeout

By default a timed-out wait returns `false` (or `NONE` for the `VALUE` form). A recovery clause lets you choose a different result.

{{< editor lang="fql" >}}
RETURN WAITFOR VALUE NONE TIMEOUT 30ms EVERY 5ms ON TIMEOUT RETURN "gave up"
{{</ editor >}}

Use `ON ERROR` to handle a failure raised while evaluating the condition.

## Waiting for an event

In event mode, `WAITFOR` subscribes to an event source and waits for a matching event, which it returns as a value.

{{< code lang="fql" >}}
LET event = WAITFOR EVENT "navigation" IN page TIMEOUT 5s
{{</ code >}}

Event timeouts use the same Duration conversion and non-negative scheduling policy as condition waits.

A `WHEN` filter accepts only events that match a condition. Inside the filter, the incoming event is available as `.`. Multiple `WHEN` clauses must all pass.

{{< code lang="fql" >}}
WAITFOR EVENT "message" IN socket
    WHEN .type == "data"
    TIMEOUT 5s
{{</ code >}}

### Synchronizing event occurrences

Event groups use the existing event entry syntax and may subscribe to different sources:

{{< code lang="fql" >}}
LET event = WAITFOR EVENT ANY {
    "navigation" IN page
    "download" IN browser
    "disconnect" IN socket
}
TIMEOUT 10s
{{</ code >}}

`EVENT ANY` returns the first qualifying event that occurs and closes the other subscriptions. Concurrent events have no declaration-order tie-break.

`EVENT ALL` waits until every subscription has produced one qualifying event. Each arm remains satisfied after its event occurs, and the result is an array in declaration order even when the events arrive in another order:

{{< code lang="fql" >}}
LET events = WAITFOR EVENT ALL {
    "domcontentloaded" IN page
    "networkidle" IN page
}
TIMEOUT 10s
{{</ code >}}

Event groups synchronize **occurrences**, unlike polling groups, which synchronize state in one cycle. Event filters are per entry:

{{< code lang="fql" >}}
RETURN WAITFOR EVENT ANY {
    "response" IN page
        WHEN .status >= 500
    "dialog" IN page
        WHEN .type == "confirm"
}
TIMEOUT 10s
{{</ code >}}

All subscriptions are established concurrently. A timeout, cancellation, setup failure, trigger failure, stream error, or completed wait closes every remaining subscription. `TIMEOUT`, `TRIGGER`, `ON TIMEOUT`, and `ON ERROR` apply once to the whole group.

### Triggering the event

A `TRIGGER` clause runs statements *after* the subscription is set up but *before* waiting begins. For an event group, every subscription is established before the one shared trigger runs. This is how you cause the event you are waiting for without risking a race where it fires before you start listening.

{{< code lang="fql" >}}
WAITFOR EVENT "response" IN page
    TRIGGER ( button <- "click" )
    TIMEOUT 10s
{{</ code >}}

The trigger body can dispatch events and run other statements. See [Dispatch Expressions]({{< ref "dispatch" >}}).

## A host capability

Event mode only works when the source is an **observable** value — one that produces a stream of events, provided by a module or the host application, such as a browser page. Condition mode has no such requirement. See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Next steps

{{< docs-related tiles="language-control-flow-dispatch,language-types-capabilities,language-control-flow" >}}
