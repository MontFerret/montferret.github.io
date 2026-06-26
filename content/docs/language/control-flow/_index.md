---
title: "Control Flow"
sidebarTitle: "Control Flow"
weight: 90
draft: false
description: "How FQL controls evaluation through branching, iteration, subqueries, queries, dispatch, waiting, and error recovery."
---

# Control Flow

Control flow describes how an FQL script determines which expressions are evaluated, how many times they are evaluated, and when evaluation proceeds.

Most FQL code is evaluated in the order it is written. Control-flow constructs change that linear execution model by introducing branching, iteration, delegation to host values, or suspension until a condition is met.

FQL is expression-oriented, so most control-flow constructs produce values. Their results can be returned, assigned to variables, passed to functions, or composed with other expressions.

## Branching with `MATCH`

A `MATCH` expression selects one branch from a set of alternatives. Only the selected branch is evaluated.

{{< code lang="fql" >}}
MATCH status (
    "active"   => "online",
    "inactive" => "offline",
    _          => "unknown",
)
{{</ code >}}

Use `MATCH` when a value should be derived from a known set of cases or conditions.

See [Match Expressions]({{< ref "match" >}}).

## Iteration with `FOR`

A `FOR` expression evaluates its body once for each item in a source collection. It is the primary construct for iterating over arrays, query results, and other iterable values.

{{< code lang="fql" >}}
FOR n IN [1, 2, 3]
    RETURN n * 2
{{</ code >}}

The result of a `FOR` expression is the collection of values returned by its body.

See [For Loops]({{< ref "for" >}}).

## Composition with subqueries

A subquery wraps a `FOR` expression in parentheses, allowing its result to be used as a regular value.

{{< code lang="fql" >}}
LET doubled = (FOR n IN [1, 2, 3] RETURN n * 2)
{{</ code >}}

Subqueries are useful when an intermediate collection needs to be assigned, returned, passed to a function, or embedded inside another expression.

See [Subquery Expressions]({{< ref "subqueries" >}}).

## Delegation with `QUERY`

A `QUERY` expression delegates query execution to a value that supports a query capability. The query itself may use another language or selector syntax, such as CSS.

{{< code lang="fql" >}}
QUERY `.product .title` IN doc USING css
{{</ code >}}

The runtime does not interpret the query payload directly. Instead, it passes the query to the target value and returns the result produced by that value.

See [Query Expressions]({{< ref "query" >}}).

## Coordination with `DISPATCH` and `WAITFOR`

`DISPATCH` sends an event or action to a target value. `WAITFOR` suspends evaluation until an event is observed, a value becomes available, or a condition becomes true.

{{< code lang="fql" >}}
button <- "click"

WAITFOR EVENT "navigation" IN page TIMEOUT 5s
{{</ code >}}

These constructs are commonly used when FQL interacts with external systems, such as browser pages, dynamic documents, or event-driven runtimes.

See [Dispatch Expressions]({{< ref "dispatch" >}}) and [Waitfor Expressions]({{< ref "waitfor" >}}).

## Recovering from errors

A recovery tail attaches to an expression and tells the runtime what to do when it fails — return a fallback value, retry the operation, or propagate the error explicitly. The optional chaining operator (`?.`) handles the common case of accessing a member on a value that may be `NONE`.

{{< code lang="fql" >}}
LET title = QUERY ONE `.title` IN doc ON ERROR RETURN "untitled"

LET name = user?.name
{{</ code >}}

Recovery tails work on function calls, `QUERY`, `DISPATCH`, `WAITFOR`, and grouped expressions. `WAITFOR` additionally supports `ON TIMEOUT` for timeout-specific recovery.

See [Error Handling]({{< ref "error-handling" >}}).

## Language constructs and host capabilities

Some control-flow constructs are part of the language itself. `MATCH`, `FOR`, and subqueries are always available and operate on ordinary FQL values.

Other constructs depend on host capabilities. `QUERY`, `DISPATCH`, and `WAITFOR EVENT` require the target value to provide the corresponding runtime behavior. For example, a document may support querying, a browser element may support dispatching events, and a page may support observing navigation or network events.

`WAITFOR` condition mode is different: it evaluates an ordinary expression repeatedly until the condition succeeds or the timeout is reached. It does not require the target value to provide an event capability.

If a value does not support the capability required by a control-flow construct, evaluation fails at runtime.

See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Next steps

{{< docs-related tiles="language-control-flow-match,language-control-flow-for,language-control-flow-subqueries,language-control-flow-query,language-control-flow-dispatch,language-control-flow-waitfor,language-control-flow-error-handling" >}}