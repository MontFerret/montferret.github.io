---
title: "Subquery Expressions"
sidebarTitle: "Subqueries"
weight: 30
draft: false
description: "Use a parenthesized for block as a value to compose and nest FQL transformations."
---

# Subquery Expressions

A subquery is a collecting query block wrapped in parentheses and used as a value. Most often it is a [`for`]({{< ref "for" >}}) loop whose array result is assigned, returned, or passed to another expression. Because a subquery is a value, the loop must have its own `return`.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

let activeUsers = (
    for u in users {
        filter u.active
        return u.name
    }
)

return activeUsers
{{</ editor >}}

The parentheses are required when a `for` is embedded in another expression. Without them, a top-level `for` is a statement and any collected array is discarded.

At a statement or direct-return boundary, the formatter removes unnecessary loop grouping: `(for ... )` becomes `for ...`, and `return (for ... )` becomes `return for ...`. It preserves the parentheses in `let result = (for ... )`, function arguments, operator operands, and member sources because those positions do not accept an embedded bare `for`.

Parentheses do not turn a returnless loop into a value. This is intentionally invalid:

{{< code lang="fql" >}}
let result = (for item in items {
    process(item)
})
{{</ code >}}

The compiler reports `A FOR loop used as an expression must return a value.` Add a loop-owned `return`, or move the returnless loop into statement position.

## Composing transformations

Because a subquery is just a value, it can be used anywhere a value is expected — including as an argument to a function.

{{< editor lang="fql" >}}
return length(
    (for n in 1..10 { filter n % 2 == 0 return n })
)
{{</ editor >}}

A subquery can also be indexed like any other array.

{{< editor lang="fql" >}}
return (for n in 1..5 { return n * n })[2]
{{</ editor >}}

## Nesting

The `return` of one loop can be another subquery, which produces nested arrays.

{{< editor lang="fql" >}}
return (
    for i in 1..3 {
        return (
            for j in 1..3 {
                return i * j
            }
        )
    }
)
{{</ editor >}}

Each inner subquery is evaluated once per iteration of the outer loop.

The direct spelling `return for ...` expresses the same ownership without the extra parentheses. Use it when the outer loop should collect each inner array. A final bare nested collecting loop is flattened only as a compatibility exception; explicit `return for` is clearer when nested result shape matters.

## Subqueries and query expressions

A subquery composes FQL transformations: it runs a `for` block and hands you the result. This is different from a [Query Expression]({{< ref "query" >}}), which delegates a query to a host value such as an HTML document. They share the word "query" but solve different problems — use a subquery to shape data with FQL, and `query` to extract data through a host capability.

## Next steps

{{< docs-related tiles="language-control-flow-query,language-control-flow-for,language-control-flow" >}}
