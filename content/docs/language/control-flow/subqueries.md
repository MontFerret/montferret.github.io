---
title: "Subquery Expressions"
sidebarTitle: "Subqueries"
weight: 30
draft: false
description: "Use a parenthesized FOR block as a value to compose and nest FQL transformations."
---

# Subquery Expressions

A subquery is a query block wrapped in parentheses and used as a value. Most often it is a [`FOR`]({{< ref "for" >}}) loop whose result — always an array — is assigned, returned, or passed to another expression.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

LET activeUsers = (
    FOR u IN users
        FILTER u.active
        RETURN u.name
)

RETURN activeUsers
{{</ editor >}}

The parentheses are required. A `FOR` loop written without them is the output of the query itself, not a value you can place inside another expression.

## Composing transformations

Because a subquery is just a value, it can be used anywhere a value is expected — including as an argument to a function.

{{< editor lang="fql" >}}
RETURN LENGTH(
    (FOR n IN 1..10 FILTER n % 2 == 0 RETURN n)
)
{{</ editor >}}

A subquery can also be indexed like any other array.

{{< editor lang="fql" >}}
RETURN (FOR n IN 1..5 RETURN n * n)[2]
{{</ editor >}}

## Nesting

The `RETURN` of one loop can be another subquery, which produces nested arrays.

{{< editor lang="fql" >}}
RETURN (
    FOR i IN 1..3
        RETURN (
            FOR j IN 1..3
                RETURN i * j
        )
)
{{</ editor >}}

Each inner subquery is evaluated once per iteration of the outer loop.

## Subqueries and query expressions

A subquery composes FQL transformations: it runs a `FOR` block and hands you the result. This is different from a [Query Expression]({{< ref "query" >}}), which delegates a query to a host value such as an HTML document. They share the word "query" but solve different problems — use a subquery to shape data with FQL, and `QUERY` to extract data through a host capability.

## Next steps

{{< docs-related tiles="language-control-flow-query,language-control-flow-for,language-control-flow" >}}
