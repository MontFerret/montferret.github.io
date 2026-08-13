---
title: "Expressions"
sidebarTitle: "Expressions"
weight: 60
draft: false
description: "Use expressions to produce values in FQL queries."
aliases:
  - /docs/fql/expressions/
---

# Expressions

An expression is a piece of FQL code that produces a value.

Expressions appear throughout a query: on the right side of `let` and `var`, inside `return`, in function arguments, conditions, filters, object fields, array items, and other places where a value is expected.

{{< editor lang="fql" >}}
let name = "Ada"
let active = true
let roles = ["admin", "editor"]

return {
    name: name,
    active: active,
    roleCount: length(roles)
}
{{</ editor >}}

In this example, `"Ada"`, `true`, `["admin", "editor"]`, `name`, `active`, and `length(roles)` are all expressions.

## Where expressions are used

Expressions can be used anywhere FQL expects a value.

Common examples include:

- variable assignments
- return values
- function arguments
- array items
- object field values
- filter conditions
- loop inputs
- conditional branches

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", age: 36, active: true },
    { name: "Grace", age: 42, active: false },
    { name: "Linus", age: 31, active: true }
]

return for user in users {
    filter user.active && user.age >= 35
    return {
        name: user.name,
        label: concat(user.name, " is active")
    }
}
{{</ editor >}}

The query uses expressions in several different positions:

- `users` is the input expression for the `for` loop.
- `user.active && user.age >= 35` is the filter expression.
- `user.name` is an object field value expression.
- `concat(user.name, " is active")` is a function call expression.

## Literal expressions

A literal expression writes a value directly in the query.

{{< editor lang="fql" >}}
return {
    none: none,
    boolean: true,
    number: 42,
    string: "hello",
    array: [1, 2, 3],
    object: { name: "Ada" }
}
{{</ editor >}}

Literals are the most direct way to create basic values.

For the complete list of built-in value types, see the [Values and Types section]({{% ref "types/_index.md" %}}).

## Variable references

A variable reference is an expression that reads the value of a variable.

{{< editor lang="fql" >}}
let name = "Ada"
let greeting = concat("Hello, ", name)

return greeting
{{</ editor >}}

The expression name evaluates to the value assigned by the earlier `let` statement.

Variables are resolved from the current scope. A variable can only be referenced after it has been declared in a scope where it is visible.

## Property access

Property access reads a field from an object or runtime-backed value.

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
    profile: {
        city: "London"
    }
}

return user.profile.city
{{</ editor >}}

Property access can be chained when nested values are being read.

For objects, the accessed property is matched by name. For runtime-backed values, the behavior depends on the value and the runtime that provides it.

## Indexed access

Indexed access reads an item from a value by position or key.

{{< editor lang="fql" >}}
let users = ["Ada", "Grace", "Linus"]

return users[0]
{{</ editor >}}

Indexes are expressions too:

{{< editor lang="fql" >}}
let users = ["Ada", "Grace", "Linus"]
let index = 1

return users[index]
{{</ editor >}}

Indexed access is commonly used with arrays. Host values may also support indexed access if the runtime defines that behavior.

## Operators

Operators combine expressions into larger expressions.

{{< editor lang="fql">}}
let price = 100
let quantity = 3

return price * quantity >= 250
{{</ editor >}}

The expression `price * quantity` produces a number, and comparing it with `>= 250` produces a boolean.

See the [Operators section]({{% ref "operators" %}}) for the full list of supported operators and precedence rules.

## Function calls

A function call is an expression that invokes a function and produces its result.

{{< editor lang="fql" >}}
let firstName = "Ada"
let lastName = "Lovelace"

return concat(upper(firstName), " ", upper(lastName))
{{</ editor >}}

Function arguments are expressions too. The inner calls to `upper` are evaluated and passed as arguments to `concat`.

For details on function declarations and calls, see the [Functions section]({{% ref "functions" %}}).

## Collection expressions

A collection expression creates an array or object value.

{{< editor lang="fql" >}}
let first = "Ada"
let second = "Grace"

return [first, second, "Linus"]
{{</ editor >}}

{{< editor lang="fql" >}}
let name = "Ada"
let active = true

return {
    name: name,
    active: active,
    label: concat(name, " is active")
}
{{</ editor >}}

Array items and object field values are expressions. They are evaluated in order and stored in the resulting collection.

Collections can contain any value type, including nested arrays and objects.

{{< editor lang="fql" >}}
return [
    { name: "Ada", roles: ["admin", "editor"] },
    { name: "Grace", roles: ["viewer"] }
]
{{</ editor >}}

Object field names become property names in the resulting object.

## Conditional expressions

A conditional expression selects a value based on a condition, using the ternary operator.

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
    active: true
}

return user.active ? "active" : "inactive"
{{</ editor >}}

The condition is evaluated first. If it is true, the first branch is used. Otherwise, the second branch is used. Both branches are expressions.

See the [Ternary Operator section]({{% ref "operators/ternary" %}}) for the full syntax, including the shortcut form.

## Subquery expressions

Some query constructs can produce values and be used as expressions.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

let activeUsers = (
    for user in users {
        filter user.active
        return user.name
    }
)

return activeUsers
{{</ editor >}}

Nested query blocks can also be used this way:

{{< editor lang="fql" >}}
let products = (
    for i in 1..5 {
        for x in 1..5 {
            return i * x
        }
    }
)

return products
{{</ editor >}}

A subquery expression evaluates a query block and uses its result as a value. This allows a query result to be assigned to a variable, returned as part of another value, or passed to a function.

See [Subquery Expressions]({{% ref "control-flow/subqueries" %}}) for more, including nesting and indexing.

## Expressions and statements

Expressions produce values.

Statements describe query structure, control flow, or variable declarations.

For example, `let` is a statement. The code on the right side of `=` is an expression:

{{< code lang="fql" >}}
let total = price * quantity
{{</ code >}}

`return` is also a statement. The value after `return` is an expression:

{{< editor lang="fql" >}}
return total >= 250
{{</ editor >}}

Only the expression parts of a statement can be nested inside other expressions.

## Evaluation

Expressions are evaluated when the surrounding statement or expression is evaluated.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false }
]

return for user in users {
    return {
        name: user.name,
        active: user.active
    }
}
{{</ editor >}}

In this query, the object expression after `return` is evaluated once for each item produced by the loop.

Expression evaluation follows the structure of the query. Nested expressions are evaluated as needed to produce the value of the outer expression.

## Next steps

{{< docs-related tiles="language-operators,language-control-flow,language-functions" >}}
