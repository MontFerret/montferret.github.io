---
title: "Expressions"
sidebarTitle: "Expressions"
weight: 110
draft: false
description: "Use expressions to produce values in FQL queries."
aliases:
  - /docs/fql/expressions/
---

# Expressions

An expression is a piece of FQL code that produces a value.

Expressions appear throughout a query: on the right side of `LET` and `VAR`, inside `RETURN`, in function arguments, conditions, filters, object fields, array items, and other places where a value is expected.

{{< editor lang="fql" >}}
LET name = "Ada"
LET active = true
LET roles = ["admin", "editor"]

RETURN {
    name: name,
    active: active,
    roleCount: LENGTH(roles)
}
{{</ editor >}}

In this example, `"Ada"`, `true`, `["admin", "editor"]`, `name`, `active`, and `LENGTH(roles)` are all expressions.

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
LET users = [
    { name: "Ada", age: 36, active: true },
    { name: "Grace", age: 42, active: false },
    { name: "Linus", age: 31, active: true }
]

FOR user IN users
    FILTER user.active && user.age >= 35
    RETURN {
        name: user.name,
        label: CONCAT(user.name, " is active")
    }
{{</ editor >}}

The query uses expressions in several different positions:

- `users` is the input expression for the `FOR` loop.
- `user.active && user.age >= 35` is the filter expression.
- `user.name` is an object field value expression.
- `CONCAT(user.name, " is active")` is a function call expression.

## Literal expressions

A literal expression writes a value directly in the query.

{{< editor lang="fql" >}}
RETURN {
    none: NONE,
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
LET name = "Ada"
LET greeting = CONCAT("Hello, ", name)

RETURN greeting
{{</ editor >}}

The expression name evaluates to the value assigned by the earlier `LET` statement.

Variables are resolved from the current scope. A variable can only be referenced after it has been declared in a scope where it is visible.

## Property access

Property access reads a field from an object or runtime-backed value.

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
    profile: {
        city: "London"
    }
}

RETURN user.profile.city
{{</ editor >}}

Property access can be chained when nested values are being read.

For objects, the accessed property is matched by name. For runtime-backed values, the behavior depends on the value and the runtime that provides it.

## Indexed access

Indexed access reads an item from a value by position or key.

{{< editor lang="fql" >}}
LET users = ["Ada", "Grace", "Linus"]

RETURN users[0]
{{</ editor >}}

Indexes are expressions too:

{{< editor lang="fql" >}}
LET users = ["Ada", "Grace", "Linus"]
LET index = 1

RETURN users[index]
{{</ editor >}}

Indexed access is commonly used with arrays. Host values may also support indexed access if the runtime defines that behavior.

## Operators

Operators combine expressions into larger expressions.

{{< editor lang="fql">}}
LET price = 100
LET quantity = 3
LET total = price * quantity

RETURN total >= 250
{{</ editor >}}

The expression price `*` quantity produces a number.

The expression `total >= 250` produces a boolean.

Operators are available for common operations such as arithmetic, comparison, logical checks, and value composition. 
See the [Operators section]({{% ref "operators" %}}) for the full list of supported operators and precedence rules.

## Function calls

A function call is an expression that invokes a function and produces its result.

{{< editor lang="fql" >}}
LET name = "Ada Lovelace"

RETURN LOWER(name)
{{</ editor >}}

Function arguments are expressions:

{{< editor lang="fql" >}}
LET firstName = "Ada"
LET lastName = "Lovelace"

RETURN CONCAT(UPPER(firstName), " ", UPPER(lastName))
{{</ editor >}}

The inner calls to `UPPER` are evaluated and passed as arguments to `CONCAT`.

Functions declared in the current script are called the same way:

{{< editor lang="fql" >}}
FUNC DISPLAY_NAME(firstName, lastName) (
    RETURN CONCAT(firstName, " ", lastName)
)

RETURN DISPLAY_NAME("Ada", "Lovelace")
{{</ editor >}}

Functions can come from the standard library, modules, the runtime environment, or declarations in the current script.

For details on function declarations and calls, see the [Functions section]({{% ref "functions" %}}).

## Array expressions

An array expression creates an array value.

{{< editor lang="fql" >}}
LET first = "Ada"
LET second = "Grace"

RETURN [first, second, "Linus"]
{{</ editor >}}

Each array item is an expression. Items are evaluated in order and stored as array values.

Array expressions can contain any value type, including nested arrays and objects.

{{< editor lang="fql" >}}
RETURN [
    { name: "Ada", roles: ["admin", "editor"] },
    { name: "Grace", roles: ["viewer"] }
]
{{</ editor >}}

## Object expressions

An object expression creates an object value.

{{< editor lang="fql" >}}
LET name = "Ada"
LET active = true

RETURN {
    name: name,
    active: active,
    label: CONCAT(name, " is active")
}
{{</ editor >}}

Object field values are expressions.

Object field names are not evaluated unless the syntax explicitly supports a dynamic field name. A regular field name such as name or active is used as the property name in the resulting object.

## Conditional expressions

A conditional expression selects a value based on a condition.

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
    active: true
}

RETURN user.active ? "active" : "inactive"
{{</ editor >}}

The condition is evaluated first. If it is true, the first branch is used. Otherwise, the second branch is used.

Both branches are expressions, so they can contain literals, variables, function calls, object expressions, array expressions, or other expressions.

## Subquery expressions

Some query constructs can produce values and be used as expressions.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

LET activeUsers = (
    FOR user IN users
        FILTER user.active
        RETURN user.name
)

RETURN activeUsers
{{</ editor >}}

Nested query blocks can also be used this way:

{{< editor lang="fql" >}}
LET products = (
    FOR i IN 1..5
        FOR x IN 1..5
            RETURN i * x
)

RETURN products
{{</ editor >}}

A subquery expression evaluates a query block and uses its result as a value. This allows a query result to be assigned to a variable, returned as part of another value, or passed to a function.

## Expressions and statements

Expressions produce values.

Statements describe query structure, control flow, or variable declarations.

For example, `LET` is a statement. The code on the right side of `=` is an expression:

{{< code lang="fql" >}}
LET total = price * quantity
{{</ code >}}

`RETURN` is also a statement. The value after `RETURN` is an expression:

{{< editor lang="fql" >}}
RETURN total >= 250
{{</ editor >}}

Only the expression parts of a statement can be nested inside other expressions.

## Evaluation

Expressions are evaluated when the surrounding statement or expression is evaluated.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false }
]

FOR user IN users
    RETURN {
        name: user.name,
        active: user.active
    }
{{</ editor >}}

In this query, the object expression after `RETURN` is evaluated once for each item produced by the loop.

Expression evaluation follows the structure of the query. Nested expressions are evaluated as needed to produce the value of the outer expression.

## Result values

Every expression produces a value. The result can be a basic value, a structured value, or a runtime-backed value provided by the environment.

{{< editor lang="fql" >}}
LET response = DOCUMENT("https://mockery.ferretlang.org")

RETURN response
{{</ editor >}}

In this example, the exact kind of value returned by `DOCUMENT` depends on the runtime and available modules.
