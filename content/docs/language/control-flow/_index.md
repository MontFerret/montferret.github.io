---
title: "Control Flow"
sidebarTitle: "Control Flow"
weight: 90
draft: false
description: "How FQL decides whether, how often, and in what order expressions are evaluated: MATCH, FOR, subqueries, QUERY, DISPATCH, and WAITFOR."
---

# Control Flow

Most FQL code runs in a straight line: every expression is evaluated once, in the order it is written. Control flow is what changes that. These constructs decide **whether** an expression is evaluated, **how often** it is evaluated, and **in what order** — or at what moment — evaluation happens.

FQL is expression-oriented, so most control-flow constructs are themselves expressions: they produce a value that can be returned, assigned, or nested inside another expression.

## Whether: choosing a value

A `MATCH` expression selects one of several branches based on a value or a set of conditions, and evaluates only the branch it chooses.

{{< code lang="fql" >}}
MATCH status (
    "active"   => "online",
    "inactive" => "offline",
    _          => "unknown",
)
{{</ code >}}

See [Match Expressions]({{< ref "match" >}}).

## How often: iterating and composing

A `FOR` loop evaluates its body once per item in a source collection. It is the primary way to repeat work in FQL.

{{< code lang="fql" >}}
FOR n IN [1, 2, 3]
    RETURN n * 2
{{</ code >}}

A subquery wraps a `FOR` block in parentheses so its result becomes a value you can assign, return, or pass to a function.

{{< code lang="fql" >}}
LET doubled = (FOR n IN [1, 2, 3] RETURN n * 2)
{{</ code >}}

See [For Loops]({{< ref "for" >}}) and [Subquery Expressions]({{< ref "subqueries" >}}).

## Delegating: running queries through host values

A `QUERY` expression hands a query — written in another dialect such as CSS — to a value that knows how to run it, and returns the results.

{{< code lang="fql" >}}
QUERY `.product .title` IN doc USING css
{{</ code >}}

See [Query Expressions]({{< ref "query" >}}).

## In what order: events and conditions

`DISPATCH` emits an event to a value, and `WAITFOR` suspends evaluation until an event arrives or a condition becomes true. Together they coordinate a query with the outside world.

{{< code lang="fql" >}}
button <- "click"

WAITFOR EVENT "navigation" IN page TIMEOUT 5s
{{</ code >}}

See [Dispatch Expressions]({{< ref "dispatch" >}}) and [Waitfor Expressions]({{< ref "waitfor" >}}).

## Language constructs and host capabilities

`MATCH`, `FOR`, and subqueries are pure language constructs. They are always available and depend only on the values you give them.

`QUERY`, `DISPATCH`, and `WAITFOR EVENT` depend on **host capabilities**. They work only when the value they act on supports the matching behavior — a queryable document, a dispatchable target, or an observable event source provided by a module or runtime. `WAITFOR`'s condition (predicate) mode is the exception: it polls an ordinary expression and needs no special capability.

When a value does not support the required capability, the query fails at runtime. See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Where to go next

{{< docs-related tiles="language-control-flow-match,language-control-flow-for,language-control-flow-subqueries,language-control-flow-query,language-control-flow-dispatch,language-control-flow-waitfor" >}}
