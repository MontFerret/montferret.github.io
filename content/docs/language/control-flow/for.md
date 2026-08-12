---
title: "For Loops"
sidebarTitle: "For"
weight: 20
draft: false
description: "Iterate over collections and ranges with for, and shape results with filter, sort, limit, and collect."
---

# For Loops

A `for` loop evaluates its body once for each item produced by a source and collects the returned values into an array. It is the primary iteration construct in FQL.

{{< editor lang="fql" >}}
return for n in [1, 2, 3, 4] {
    return n * 2
}
{{</ editor >}}

The source is evaluated once. The body — including the final `return` — is evaluated once per item, and the result is always an array.

Braces are optional, but they must be paired. The documentation uses the braced form because it makes the loop boundary explicit; the unbraced form remains fully supported and has the same semantics.

### Legacy unbraced form

Existing scripts may omit both braces:

{{< code lang="fql" >}}
return for n in [1, 2, 3]
    return n * 2
{{</ code >}}

This is equivalent to the braced example above. Do not mix a single opening or closing brace with an otherwise unbraced body.

## Iterating a source

The value after `in` is the source. It can be an array, a [range]({{% ref "../operators/range" %}}), or any expression that yields a collection, including a variable, a function call, or a bind parameter.

{{< editor lang="fql" >}}
return for n in 1..5 {
    return n
}
{{</ editor >}}

The `1..5` range produces the integers from 1 to 5.

### The counter variable

A second variable after the loop variable receives the zero-based position of each item.

{{< editor lang="fql" >}}
return for value, index in ["a", "b", "c"] {
    return { index, value }
}
{{</ editor >}}

If you only need the position, ignore the value with `_`.

{{< editor lang="fql" >}}
return for _, index in ["a", "b", "c"] {
    return index
}
{{</ editor >}}

### Destructuring items

`for ... in` can destructure each item with the same recursive object and array patterns used by `let` and `var`.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", stats: { score: 3 } },
    { name: "Grace", stats: { score: 5 } }
]

return for { name, stats: { score } }, index in users {
    return { index, name, score }
}
{{</ editor >}}

The source is still evaluated once. At the start of each iteration, every named leaf is introduced as an immutable loop binding that is scoped to the loop. Object aliases and nested patterns use `:`, array entries are positional, and `_` skips a property or position without reading it. A nested child pattern with no named bindings is skipped as a whole, so Ferret neither retrieves nor validates that child value.

Missing properties or elements bind `none`, and nested patterns propagate `none`. Extra values are ignored. A non-`none` item reached by the root pattern or by a child pattern needed to produce a binding must support keyed access for `{ ... }` or indexed access for `[ ... ]`; otherwise execution fails at that pattern with `cannot destructure <Actual> as Object` or `cannot destructure <Actual> as Array`. Explicit empty root patterns still validate the item shape; ignored child patterns do not.

Array holes, defaults, rest or spread entries, quoted or computed keys, and pattern conditions are not supported. Destructuring does not apply to condition-driven `for ... while` loops.

## Shaping results

Clauses placed between the source and the `return` transform the stream of items. They take effect in the order they are written.

### filter

`filter` keeps only the items for which a condition is true.

{{< editor lang="fql" >}}
return for n in [1, 2, 3, 4, 1, 3] {
    filter n > 2
    return n
}
{{</ editor >}}

### sort

`sort` reorders the items by one or more keys. Each key may be followed by `asc` (the default) or `desc`.

{{< editor lang="fql" >}}
return for name in ["foo", "bar", "qaz", "abc"] {
    sort name
    return name
}
{{</ editor >}}

### limit

`limit count` keeps the first `count` items. `limit offset, count` skips `offset` items first, then keeps `count`.

{{< editor lang="fql" >}}
return for n in [1, 2, 3, 4, 5, 6, 7, 8] {
    limit 4, 2
    return n
}
{{</ editor >}}

### collect

`collect` groups items by one or more keys. The result has one entry per distinct group, and the original loop binding or destructured leaves are no longer in scope — only the group keys and anything you collect alongside them.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

return for u in users {
    collect dept = u.dept
    return dept
}
{{</ editor >}}

`with count into` counts the members of each group.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

return for u in users {
    collect dept = u.dept with count into total
    return { dept, total }
}
{{</ editor >}}

`aggregate` computes values across each group, such as `count`, `SUM`, `MIN`, `MAX`, or `AVERAGE`.

{{< editor lang="fql" >}}
let users = [
    { dept: "eng", age: 31 },
    { dept: "eng", age: 45 },
    { dept: "ops", age: 25 }
]

return for u in users {
    collect dept = u.dept
    aggregate headcount = COUNT(u), avgAge = AVERAGE(u.age)
    return { dept, headcount, avgAge }
}
{{</ editor >}}

## Looping on a condition

Instead of iterating a source, a `for` loop can repeat while a condition holds. `while` checks the condition before each pass; `do while` checks it after, so the body always runs at least once. An optional loop variable provides a zero-based counter.

{{< code lang="fql" >}}
for while condition {
    ...
    return value
}

for i while condition {
    ...
    return value
}

for i do while condition {
    ...
    return value
}
{{</ code >}}

Because `do while` runs the body before testing the condition, the loop below produces one item even though the condition is false from the start:

{{< editor lang="fql" >}}
return for i do while false {
    return i
}
{{</ editor >}}

{{< notification type="warning" >}}
Condition-driven loops keep running until the condition becomes false. Make sure the condition can change — for example, by mutating a var in the body — or bound the wait with a timeout-based waitfor instead.
{{</ notification >}}

## Returning and discarding loop results

Use `return for ...` to make the loop's collected array the script or block-function result. `return distinct for ...` deduplicates that array using the same semantics as other `return distinct` operands.

A standalone `for` is a statement: it executes the loop body and propagates errors or cancellation, but discards the collected array. This is useful for effect-only iteration. If the surrounding script or block function then falls through, its result is `none`.

Wrapped in parentheses, a `for` is an expression you can assign, pass, or nest. See [Subquery Expressions]({{< ref "subqueries" >}}).

## Next steps

{{< docs-related tiles="language-control-flow-subqueries,language-control-flow-error-handling,language-functions" >}}
