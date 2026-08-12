---
title: "Variables"
sidebarTitle: "Variables"
weight: 50
draft: false
description: "Declare immutable variables with let and mutable variables with var."
aliases:
    - /docs/fql/operators/high-level/let/
---

# Variables

FQL supports two kinds of variable declarations:

- `let` declares an immutable variable.
- `var` declares a mutable variable that can be reassigned later.

## let

The `let` statement assigns the result of an expression to a variable.

The variable is introduced into the scope where the `let` statement appears.

{{< code lang="fql" >}}
let variableName = expression
{{< /code >}}

After a variable is declared with `let`, it cannot be reassigned.

{{< code lang="fql" >}}
let a = [1, 2, 3]  // initial assignment

a = PUSH(a, 4)     // syntax error, unexpected identifier
let a = PUSH(a, 4) // parsing error, variable 'a' is assigned multiple times
let b = PUSH(a, 4) // allowed, result: [1, 2, 3, 4]
{{< /code >}}

`let` bindings often appear where a query needs to refer to an intermediate value, a subquery result, or another computed expression by name.

{{< editor lang="fql" >}}
let users = [
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

return for u in users {
    let numProducts = LENGTH(u.cart)

    return {
        "user": u,
        "numProducts": numProducts,
        "discount": numProducts >= 3
    }
}
{{< /editor >}}

In this example, `numProducts` stores the number of items in the user's cart. The value can then be reused in the returned object without calling `LENGTH(u.cart)` more than once.

`let` is also useful for assigning the result of a subquery to a variable.

{{< editor lang="fql" >}}
let users = [
    { "id": 1, "name": "Moises Grisewood" },
    { "id": 2, "name": "Dell Marnes" },
    { "id": 3, "name": "Tobin Bilbery" },
    { "id": 4, "name": "Lorianne Posten" },
    { "id": 5, "name": "Drucill Cryer" }
]
let friends = [
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

return for u in users {
    let friends = (
        for f in friends {
            filter u.id == f.userId
            return f
        }
    )

    return { "user": u, "friends": friends, "numFriends": LENGTH(friends) }
}
{{< /editor >}}

Here, the inner `for` query finds the friends that belong to the current user. Assigning the result to friends makes the `return` statement easier to read and allows the same value to be used more than once.

## var

The `var` statement declares a mutable variable. Unlike variables declared with `let`, variables declared with `var` can be reassigned.

The general syntax is:

{{< code lang="fql" >}}
var variableName = expression
{{< /code >}}

A `var` variable can be updated later using assignment syntax:

{{< editor lang="fql" >}}
var count = 0

count = count + 1
count = count + 1

return count
{{< /editor >}}

`var` is suitable for values that change across multiple statements, such as counters, accumulators, or flags.

{{< editor lang="fql" >}}
let numbers = [1, 2, 3, 4, 5]
var total = 0

return for n in numbers {
    total = total + n
    
    return total
}
{{< /editor >}}

In this example, `total` starts at 0 and is updated for each value in numbers.

## Destructuring declarations

`let` and `var` can bind several names from one object or array. The source expression is evaluated once.

{{< editor lang="fql" >}}
let {
    name,
    profile: { city },
    scores: [first, _, third]
} = {
    name: "Ada",
    profile: { city: "London" },
    scores: [10, 20, 30]
}

return { name, city, first, third }
{{</ editor >}}

Object entries use the property name as the binding name. Add `:` to use a different name or a nested pattern. Array entries are positional, and `_` explicitly ignores a property or position without reading it. A nested child pattern with no named bindings, such as `metadata: [_, _]` or `details: {}`, is ignored as a whole: Ferret does not retrieve or validate that child value. Patterns can be nested recursively and can end with a trailing comma. Empty patterns (`{}` and `[]`) are valid.

Every named leaf follows the declaration kind. Leaves declared by `let` are immutable. Leaves declared by `var` are independently mutable:

{{< editor lang="fql" >}}
var { count, step } = { count: 1, step: 2 }

count += step
step += 1

return { count, step }
{{</ editor >}}

Missing properties or elements bind `none`. Destructuring `none` also propagates `none` through nested patterns, while extra source values are ignored. A non-`none` value reached by the root pattern or by a child pattern needed to produce a binding must support keyed access for an object pattern or indexed access for an array pattern; otherwise execution fails with `cannot destructure <Actual> as Object` or `cannot destructure <Actual> as Array`. Explicit empty root patterns still perform this shape check; ignored child patterns do not.

Array holes are not allowed; write `_` for each ignored position. Defaults, rest or spread entries, quoted or computed object keys, and conditions are not supported in binding patterns. Destructuring is declaration syntax only: assignments and function parameters still bind one existing target or name.

## Compound assignment

FQL supports compound assignment operators that combine an arithmetic operation with assignment. They are shorthand for updating a `var` variable in place:

| Operator | Equivalent |
| --- | --- |
| `a += b` | `a = a + b` |
| `a -= b` | `a = a - b` |
| `a *= b` | `a = a * b` |
| `a /= b` | `a = a / b` |

{{< editor lang="fql" >}}
var total = 100

total += 50
total -= 10
total *= 2

return total
{{</ editor >}}

Compound assignment operators can only be used with `var` variables. Using them with `let` bindings is an error.

## Next steps

{{< docs-related tiles="language-expressions,language-operators,language-control-flow" >}}
