---
title: "For Loops"
sidebarTitle: "For"
weight: 20
draft: false
description: "Iterate over collections and ranges with FOR, and shape results with FILTER, SORT, LIMIT, and COLLECT."
---

# For Loops

A `FOR` loop evaluates its body once for each item produced by a source and collects the returned values into an array. It is the primary iteration construct in FQL.

{{< editor lang="fql" >}}
RETURN FOR n IN [1, 2, 3, 4] {
    RETURN n * 2
}
{{</ editor >}}

The source is evaluated once. The body — including the final `RETURN` — is evaluated once per item, and the result is always an array.

Braces are optional, but they must be paired. The documentation uses the braced form because it makes the loop boundary explicit; the unbraced form remains fully supported and has the same semantics.

### Legacy unbraced form

Existing scripts may omit both braces:

{{< code lang="fql" >}}
RETURN FOR n IN [1, 2, 3]
    RETURN n * 2
{{</ code >}}

This is equivalent to the braced example above. Do not mix a single opening or closing brace with an otherwise unbraced body.

## Iterating a source

The value after `IN` is the source. It can be an array, a [range]({{% ref "../operators/range" %}}), or any expression that yields a collection, including a variable, a function call, or a bind parameter.

{{< editor lang="fql" >}}
RETURN FOR n IN 1..5 {
    RETURN n
}
{{</ editor >}}

The `1..5` range produces the integers from 1 to 5.

### The counter variable

A second variable after the loop variable receives the zero-based position of each item.

{{< editor lang="fql" >}}
RETURN FOR value, index IN ["a", "b", "c"] {
    RETURN { index, value }
}
{{</ editor >}}

If you only need the position, ignore the value with `_`.

{{< editor lang="fql" >}}
RETURN FOR _, index IN ["a", "b", "c"] {
    RETURN index
}
{{</ editor >}}

### Destructuring items

`FOR ... IN` can destructure each item with the same recursive object and array patterns used by `LET` and `VAR`.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", stats: { score: 3 } },
    { name: "Grace", stats: { score: 5 } }
]

RETURN FOR { name, stats: { score } }, index IN users {
    RETURN { index, name, score }
}
{{</ editor >}}

The source is still evaluated once. At the start of each iteration, every named leaf is introduced as an immutable loop binding that is scoped to the loop. Object aliases and nested patterns use `:`, array entries are positional, and `_` skips a property or position without reading it. A nested child pattern with no named bindings is skipped as a whole, so Ferret neither retrieves nor validates that child value.

Missing properties or elements bind `NONE`, and nested patterns propagate `NONE`. Extra values are ignored. A non-`NONE` item reached by the root pattern or by a child pattern needed to produce a binding must support keyed access for `{ ... }` or indexed access for `[ ... ]`; otherwise execution fails at that pattern with `cannot destructure <Actual> as Object` or `cannot destructure <Actual> as Array`. Explicit empty root patterns still validate the item shape; ignored child patterns do not.

Array holes, defaults, rest or spread entries, quoted or computed keys, and pattern conditions are not supported. Destructuring does not apply to condition-driven `FOR ... WHILE` loops.

## Shaping results

Clauses placed between the source and the `RETURN` transform the stream of items. They take effect in the order they are written.

### FILTER

`FILTER` keeps only the items for which a condition is true.

{{< editor lang="fql" >}}
RETURN FOR n IN [1, 2, 3, 4, 1, 3] {
    FILTER n > 2
    RETURN n
}
{{</ editor >}}

### SORT

`SORT` reorders the items by one or more keys. Each key may be followed by `ASC` (the default) or `DESC`.

{{< editor lang="fql" >}}
RETURN FOR name IN ["foo", "bar", "qaz", "abc"] {
    SORT name
    RETURN name
}
{{</ editor >}}

### LIMIT

`LIMIT count` keeps the first `count` items. `LIMIT offset, count` skips `offset` items first, then keeps `count`.

{{< editor lang="fql" >}}
RETURN FOR n IN [1, 2, 3, 4, 5, 6, 7, 8] {
    LIMIT 4, 2
    RETURN n
}
{{</ editor >}}

### COLLECT

`COLLECT` groups items by one or more keys. The result has one entry per distinct group, and the original loop binding or destructured leaves are no longer in scope — only the group keys and anything you collect alongside them.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

RETURN FOR u IN users {
    COLLECT dept = u.dept
    RETURN dept
}
{{</ editor >}}

`WITH COUNT INTO` counts the members of each group.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

RETURN FOR u IN users {
    COLLECT dept = u.dept WITH COUNT INTO total
    RETURN { dept, total }
}
{{</ editor >}}

`AGGREGATE` computes values across each group, such as `COUNT`, `SUM`, `MIN`, `MAX`, or `AVERAGE`.

{{< editor lang="fql" >}}
LET users = [
    { dept: "eng", age: 31 },
    { dept: "eng", age: 45 },
    { dept: "ops", age: 25 }
]

RETURN FOR u IN users {
    COLLECT dept = u.dept
    AGGREGATE headcount = COUNT(u), avgAge = AVERAGE(u.age)
    RETURN { dept, headcount, avgAge }
}
{{</ editor >}}

## Looping on a condition

Instead of iterating a source, a `FOR` loop can repeat while a condition holds. `WHILE` checks the condition before each pass; `DO WHILE` checks it after, so the body always runs at least once. An optional loop variable provides a zero-based counter.

{{< code lang="fql" >}}
FOR WHILE condition {
    ...
    RETURN value
}

FOR i WHILE condition {
    ...
    RETURN value
}

FOR i DO WHILE condition {
    ...
    RETURN value
}
{{</ code >}}

Because `DO WHILE` runs the body before testing the condition, the loop below produces one item even though the condition is false from the start:

{{< editor lang="fql" >}}
RETURN FOR i DO WHILE false {
    RETURN i
}
{{</ editor >}}

{{< notification type="warning" >}}
Condition-driven loops keep running until the condition becomes false. Make sure the condition can change — for example, by mutating a VAR in the body — or bound the wait with a TIMEOUT-based WAITFOR instead.
{{</ notification >}}

## Returning and discarding loop results

Use `RETURN FOR ...` to make the loop's collected array the script or block-function result. `RETURN DISTINCT FOR ...` deduplicates that array using the same semantics as other `RETURN DISTINCT` operands.

A standalone `FOR` is a statement: it executes the loop body and propagates errors or cancellation, but discards the collected array. This is useful for effect-only iteration. If the surrounding script or block function then falls through, its result is `NONE`.

Wrapped in parentheses, a `FOR` is an expression you can assign, pass, or nest. See [Subquery Expressions]({{< ref "subqueries" >}}).

## Next steps

{{< docs-related tiles="language-control-flow-subqueries,language-control-flow-error-handling,language-functions" >}}
