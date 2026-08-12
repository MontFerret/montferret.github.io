---
title: "Subquery Expressions"
sidebarTitle: "Subqueries"
weight: 30
draft: false
description: "Use a parenthesized for block as a value to compose and nest FQL transformations."
---

# Subquery Expressions

A subquery is a query block wrapped in parentheses and used as a value. Most often it is a [`for`]({{< ref "for" >}}) loop whose result — always an array — is assigned, returned, or passed to another expression.

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

The parentheses are required. A `for` loop written without them is the output of the query itself, not a value you can place inside another expression.

## Composing transformations

Because a subquery is just a value, it can be used anywhere a value is expected — including as an argument to a function.

{{< editor lang="fql" >}}
return LENGTH(
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

## Subqueries and query expressions

A subquery composes FQL transformations: it runs a `for` block and hands you the result. This is different from a [Query Expression]({{< ref "query" >}}), which delegates a query to a host value such as an HTML document. They share the word "query" but solve different problems — use a subquery to shape data with FQL, and `query` to extract data through a host capability.

## Next steps

{{< docs-related tiles="language-control-flow-query,language-control-flow-for,language-control-flow" >}}
