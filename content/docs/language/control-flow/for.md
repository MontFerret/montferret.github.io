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
FOR n IN [1, 2, 3, 4]
    RETURN n * 2
{{</ editor >}}

The source is evaluated once. The body — including the final `RETURN` — is evaluated once per item, and the result is always an array.

## Iterating a source

The value after `IN` is the source. It can be an array, a [range]({{% ref "../operators/range" %}}), or any expression that yields a collection, including a variable, a function call, or a bind parameter.

{{< editor lang="fql" >}}
FOR n IN 1..5
    RETURN n
{{</ editor >}}

The `1..5` range produces the integers from 1 to 5.

### The counter variable

A second variable after the loop variable receives the zero-based position of each item.

{{< editor lang="fql" >}}
FOR value, index IN ["a", "b", "c"]
    RETURN { index, value }
{{</ editor >}}

If you only need the position, ignore the value with `_`.

{{< editor lang="fql" >}}
FOR _, index IN ["a", "b", "c"]
    RETURN index
{{</ editor >}}

## Shaping results

Clauses placed between the source and the `RETURN` transform the stream of items. They take effect in the order they are written.

### FILTER

`FILTER` keeps only the items for which a condition is true.

{{< editor lang="fql" >}}
FOR n IN [1, 2, 3, 4, 1, 3]
    FILTER n > 2
    RETURN n
{{</ editor >}}

### SORT

`SORT` reorders the items by one or more keys. Each key may be followed by `ASC` (the default) or `DESC`.

{{< editor lang="fql" >}}
FOR name IN ["foo", "bar", "qaz", "abc"]
    SORT name
    RETURN name
{{</ editor >}}

### LIMIT

`LIMIT count` keeps the first `count` items. `LIMIT offset, count` skips `offset` items first, then keeps `count`.

{{< editor lang="fql" >}}
FOR n IN [1, 2, 3, 4, 5, 6, 7, 8]
    LIMIT 4, 2
    RETURN n
{{</ editor >}}

### COLLECT

`COLLECT` groups items by one or more keys. The result has one entry per distinct group, and the original loop variable is no longer in scope — only the group keys and anything you collect alongside them.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

FOR u IN users
    COLLECT dept = u.dept
    RETURN dept
{{</ editor >}}

`WITH COUNT INTO` counts the members of each group.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", dept: "eng" },
    { name: "Grace", dept: "eng" },
    { name: "Linus", dept: "ops" }
]

FOR u IN users
    COLLECT dept = u.dept WITH COUNT INTO total
    RETURN { dept, total }
{{</ editor >}}

`AGGREGATE` computes values across each group, such as `COUNT`, `SUM`, `MIN`, `MAX`, or `AVERAGE`.

{{< editor lang="fql" >}}
LET users = [
    { dept: "eng", age: 31 },
    { dept: "eng", age: 45 },
    { dept: "ops", age: 25 }
]

FOR u IN users
    COLLECT dept = u.dept
    AGGREGATE headcount = COUNT(u), avgAge = AVERAGE(u.age)
    RETURN { dept, headcount, avgAge }
{{</ editor >}}

## Looping on a condition

Instead of iterating a source, a `FOR` loop can repeat while a condition holds. `WHILE` checks the condition before each pass; `DO WHILE` checks it after, so the body always runs at least once. An optional loop variable provides a zero-based counter.

{{< code lang="fql" >}}
FOR WHILE condition
    ...
    RETURN value

FOR i WHILE condition
    ...
    RETURN value

FOR i DO WHILE condition
    ...
    RETURN value
{{</ code >}}

Because `DO WHILE` runs the body before testing the condition, the loop below produces one item even though the condition is false from the start:

{{< editor lang="fql" >}}
FOR i DO WHILE false
    RETURN i
{{</ editor >}}

{{< notification type="warning" >}}
Condition-driven loops keep running until the condition becomes false. Make sure the condition can change — for example, by mutating a VAR in the body — or bound the wait with a TIMEOUT-based WAITFOR instead.
{{</ notification >}}

## Loops as values

A `FOR` loop written on its own is the query's output. Wrapped in parentheses, it becomes a value you can assign or nest. See [Subquery Expressions]({{< ref "subqueries" >}}).

## Next steps

{{< docs-related tiles="language-control-flow-subqueries,language-control-flow-error-handling,language-functions" >}}
