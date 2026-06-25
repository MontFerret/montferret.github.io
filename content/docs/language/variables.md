---
title: "Variables"
sidebarTitle: "Variables"
weight: 50
draft: false
description: "Declare immutable variables with LET and mutable variables with VAR."
aliases:
    - /docs/fql/operators/high-level/let/
---

# Variables

FQL supports two kinds of variable declarations:

- `LET` declares an immutable variable.
- `VAR` declares a mutable variable that can be reassigned later.

## LET

The `LET` statement assigns the result of an expression to a variable.

The variable is introduced into the scope where the `LET` statement appears.

{{< code lang="fql" >}}
LET variableName = expression
{{< /code >}}

After a variable is declared with `LET`, it cannot be reassigned.

{{< code lang="fql" >}}
LET a = [1, 2, 3]  // initial assignment

a = PUSH(a, 4)     // syntax error, unexpected identifier
LET a = PUSH(a, 4) // parsing error, variable 'a' is assigned multiple times
LET b = PUSH(a, 4) // allowed, result: [1, 2, 3, 4]
{{< /code >}}

`LET` bindings often appear where a query needs to refer to an intermediate value, a subquery result, or another computed expression by name.

{{< editor lang="fql" >}}
LET users = [
    {
        "id": 1,
        "firstName": "Kikelia",
        "lastName": "Coper",
        "cart": [
            { "Name": "Garlic", "Price": "$7.46" },
            { "Name": "Flower - Commercial Spider", "Price": "$6.59" }
        ]
    },
    {
        "id": 2,
        "firstName": "Toni",
        "lastName": "MacTeggart",
        "cart": [
            { "Name": "Spice - Paprika", "Price": "$6.31" },
            { "Name": "Extract - Vanilla,artificial", "Price": "$4.74" },
            { "Name": "Wine - White, Cooking", "Price": "$1.50" },
            { "Name": "Nutmeg - Ground", "Price": "$1.30" }
        ]
    },
    {
        "id": 3,
        "firstName": "Neile",
        "lastName": "Saice",
        "cart": [
            { "Name": "Mustard Prepared", "Price": "$2.28" },
            { "Name": "Flower - Commercial Bronze", "Price": "$4.80" }
        ]
    }
]

FOR u IN users
    LET numProducts = LENGTH(u.cart)

    RETURN {
        "user": u,
        "numProducts": numProducts,
        "discount": numProducts >= 3
    }
{{< /editor >}}

In this example, `numProducts` stores the number of items in the user's cart. The value can then be reused in the returned object without calling `LENGTH(u.cart)` more than once.

`LET` is also useful for assigning the result of a subquery to a variable.

{{< editor lang="fql" >}}
LET users = [
    { "id": 1, "name": "Moises Grisewood" },
    { "id": 2, "name": "Dell Marnes" },
    { "id": 3, "name": "Tobin Bilbery" },
    { "id": 4, "name": "Lorianne Posten" },
    { "id": 5, "name": "Drucill Cryer" }
]
LET friends = [
    { "id": 1, "name": "Maximo Massard", "userId": 1 },
    { "id": 2, "name": "Delainey Sancho", "userId": 1 },
    { "id": 3, "name": "Lindon Beale", "userId": 1 },
    { "id": 4, "name": "Gus Sprey", "userId": 3 },
    { "id": 5, "name": "Virgil Dallander", "userId": 3 },
    { "id": 6, "name": "Agretha Mackerel", "userId": 4 },
    { "id": 7, "name": "Christalle Aldins", "userId": 4 },
    { "id": 8, "name": "Karalynn Margery", "userId": 5 },
    { "id": 9, "name": "Rodolph Ladd", "userId": 5 },
    { "id": 10, "name": "Babette Brassill", "userId": 1 }
]

FOR u IN users
    LET friends = (
        FOR f IN friends
            FILTER u.id == f.userId
            RETURN f
    )

    RETURN { "user": u, "friends": friends, "numFriends": LENGTH(friends) }
{{< /editor >}}

Here, the inner `FOR` query finds the friends that belong to the current user. Assigning the result to friends makes the `RETURN` statement easier to read and allows the same value to be used more than once.

## VAR

The `VAR` statement declares a mutable variable. Unlike variables declared with `LET`, variables declared with `VAR` can be reassigned.

The general syntax is:

{{< code lang="fql" >}}
VAR variableName = expression
{{< /code >}}

A `VAR` variable can be updated later using assignment syntax:

{{< editor lang="fql" >}}
VAR count = 0

count = count + 1
count = count + 1

RETURN count
{{< /editor >}}

`VAR` is suitable for values that change across multiple statements, such as counters, accumulators, or flags.

{{< editor lang="fql" >}}
LET numbers = [1, 2, 3, 4, 5]
VAR total = 0

FOR n IN numbers
    total = total + n
    
    RETURN total
{{< /editor >}}

In this example, `total` starts at 0 and is updated for each value in numbers.

## Compound assignment

FQL supports compound assignment operators that combine an arithmetic operation with assignment. They are shorthand for updating a `VAR` variable in place:

| Operator | Equivalent |
| --- | --- |
| `a += b` | `a = a + b` |
| `a -= b` | `a = a - b` |
| `a *= b` | `a = a * b` |
| `a /= b` | `a = a / b` |

{{< editor lang="fql" >}}
VAR total = 100

total += 50
total -= 10
total *= 2

RETURN total
{{</ editor >}}

Compound assignment operators can only be used with `VAR` variables. Using them with `LET` bindings is an error.

## DELETE

The `DELETE` statement removes a property from an object or host value.

{{< code lang="fql" >}}
DELETE target.property
DELETE target["property"]
{{</ code >}}

Both dot notation and bracket notation are supported:

{{< editor lang="fql" >}}
VAR user = { name: "Ada", deprecated: true, role: "admin" }

DELETE user.deprecated

RETURN user
{{</ editor >}}

Deletion removes the property entirely — it is not the same as assigning `NONE`, which keeps the key present with an absent value.

`DELETE` works with any value that supports the removable capability. Built-in objects support key removal. Host values may support removal if the host runtime provides that capability.
