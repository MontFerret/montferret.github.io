---
title: "Waitfor Expressions"
sidebarTitle: "Waitfor"
weight: 60
draft: false
description: "Suspend evaluation until a condition becomes true or an event arrives with the WAITFOR expression."
---

# Waitfor Expressions

A `WAITFOR` expression pauses a query until something happens, then resumes and produces a value. It has two modes: waiting for a **condition** to become true, and waiting for an **event** to arrive.

## Waiting for a condition

In its condition (predicate) mode, `WAITFOR` re-checks an ordinary expression until it is satisfied or a timeout is reached. This mode is a pure language construct — it needs no special capability.

There are four forms:

| Form | Waits until | Returns |
| --- | --- | --- |
| `WAITFOR <expr>` | the expression is true | `true`, or `false` on timeout |
| `WAITFOR EXISTS <expr>` | the expression has a value | `true`, or `false` on timeout |
| `WAITFOR NOT EXISTS <expr>` | the expression is empty or absent | `true`, or `false` on timeout |
| `WAITFOR VALUE <expr>` | the expression yields a value | that value, or `NONE` on timeout |

{{< editor lang="fql" >}}
RETURN WAITFOR EXISTS [1, 2, 3] TIMEOUT 100ms
{{</ editor >}}

The expression is re-evaluated on each attempt, so it can reflect state that changes over time. When the wait runs out, the result reports the timeout rather than raising an error:

{{< editor lang="fql" >}}
RETURN WAITFOR FALSE TIMEOUT 50ms EVERY 10ms
{{</ editor >}}

### Tuning the wait

Several clauses control how the wait behaves:

- `TIMEOUT <duration>` — the maximum time to wait, written with a duration literal such as `50ms`, `0.5s`, or `5s`.
- `EVERY <interval>` — how often to re-check. A second value, `EVERY <interval>, <cap>`, caps how large the interval can grow.
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

A `WHEN` filter accepts only events that match a condition. Inside the filter, the incoming event is available as `.`. Multiple `WHEN` clauses must all pass.

{{< code lang="fql" >}}
WAITFOR EVENT "message" IN socket
    WHEN .type == "data"
    TIMEOUT 5s
{{</ code >}}

### Triggering the event

A `TRIGGER` clause runs statements *after* the subscription is set up but *before* waiting begins. This is how you cause the event you are waiting for without risking a race where it fires before you start listening.

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
