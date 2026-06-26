---
title: "User-Defined Functions"
sidebarTitle: "User-Defined"
weight: 10
draft: false
description: "Declare reusable functions within a script using the FUNC keyword."
---

# User-defined functions

A user-defined function is a reusable piece of logic declared within a script. Once declared, the function can be called like any built-in function.

{{< editor lang="fql" >}}
FUNC double(x) => x * 2

RETURN double(21)
{{</ editor >}}

## Declaration

A function declaration begins with the `FUNC` keyword, followed by a name, a parameter list in parentheses, and a body.

There are two body forms: arrow and block.

### Arrow form

The arrow form uses `=>` followed by a single expression. The result of the expression is the return value.

{{< editor lang="fql" >}}
FUNC greet(name) => CONCAT("Hello, ", name, "!")

RETURN greet("Ada")
{{</ editor >}}

Use the arrow form when the function body is a single expression.

### Block form

Unlike many C-like languages, FQL uses parentheses for block function bodies.

{{< editor lang="fql" >}}
FUNC normalizePrice(input) (
    LET cleaned = TRIM(input)
    LET numeric = SUBSTITUTE(cleaned, "$", "")
    RETURN TO_FLOAT(numeric)
)

RETURN normalizePrice("  $19.99  ")
{{</ editor >}}

Use the block form when the function body requires intermediate variables or multiple steps.

## Parameters

Parameters are listed inside parentheses, separated by commas.

{{< editor lang="fql" >}}
FUNC fullName(first, last) => CONCAT(first, " ", last)

RETURN fullName("Ada", "Lovelace")
{{</ editor >}}

A function may have no parameters:

{{< editor lang="fql" >}}
FUNC now() => "2024-01-01"

RETURN now()
{{</ editor >}}

Parameters are positional. The caller must provide exactly the number of arguments the function expects.

## Capturing outer variables

A function body can read variables from the enclosing scope.

{{< editor lang="fql" >}}
LET base = 10

FUNC add(value) => base + value

RETURN add(5)
{{</ editor >}}

If the outer variable is declared with `VAR`, the function can also modify it:

{{< editor lang="fql" >}}
VAR counter = 0

FUNC inc() (
    counter = counter + 1
    RETURN counter
)

RETURN [inc(), inc(), inc()]
{{</ editor >}}

Variables declared with `LET` are immutable and cannot be reassigned inside a function.

## Nesting functions

Functions can be declared inside other functions.

{{< editor lang="fql" >}}
FUNC process(items) (
    FUNC transform(item) => item * 2

    RETURN (
        FOR item IN items
            RETURN transform(item)
    )
)

RETURN process([1, 2, 3])
{{</ editor >}}

A nested function can access variables from all enclosing scopes, not just the immediately surrounding one.

## Using functions in loops

User-defined functions work naturally with `FOR` loops and other query constructs.

{{< editor lang="fql" >}}
FUNC formatUser(user) (
    RETURN {
        label: CONCAT(user.name, " (", user.role, ")"),
        active: user.active
    }
)

LET users = [
    { name: "Ada", role: "admin", active: true },
    { name: "Grace", role: "editor", active: false },
    { name: "Linus", role: "viewer", active: true }
]

FOR user IN users
    FILTER user.active
    RETURN formatUser(user)
{{</ editor >}}

## Function names

Function names follow the same rules as variable names: they must start with a letter or underscore, followed by any combination of letters, digits, and underscores.

Function names are case-sensitive. `add` and `Add` are different functions.

{{< editor lang="fql" >}}
FUNC a() => 1
FUNC A() => 2

RETURN a() + A()
{{</ editor >}}

By convention, built-in functions are written in uppercase, while user-defined function names may use the style preferred by the script author.

## Next steps

{{< docs-related tiles="language-functions-modules,embedding-custom-functions,stdlib" >}}
