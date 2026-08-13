---
title: "User-Defined Functions"
sidebarTitle: "User-Defined"
weight: 10
draft: false
description: "Declare reusable functions within a script using the func keyword."
---

# User-defined functions

A user-defined function is a reusable piece of logic declared within a script. Once declared, the function can be called like any built-in function.

{{< editor lang="fql" >}}
func double(x) => x * 2

return double(21)
{{</ editor >}}

## Declaration

A function declaration begins with the `func` keyword, followed by a name, a parameter list in parentheses, and a body.

There are two body forms: arrow and block.

### Arrow form

The arrow form uses `=>` followed by a single expression. The result of the expression is the return value.

{{< editor lang="fql" >}}
func greet(name) => concat("Hello, ", name, "!")

return greet("Ada")
{{</ editor >}}

Use the arrow form when the function body is a single expression.

### Block form

The block form encloses the function body in braces. Use it when the function needs intermediate bindings or multiple steps.

{{< editor lang="fql" >}}
func normalizePrice(input) {
    let cleaned = trim(input)
    let numeric = substitute(cleaned, "$", "")
    return to_float(numeric)
}

return normalizePrice("  $19.99  ")
{{</ editor >}}

An explicit `return` sets the block function's result. If execution reaches the closing brace without one, the function completes successfully with `none`. Arbitrary expression statements are evaluated and discarded rather than returned implicitly; use the arrow form for a single expression.

An empty block is therefore a valid effect-only function:

{{< editor lang="fql" >}}
func noop() {}

return noop()
{{</ editor >}}

### Returning a for result

A block function can return a loop directly with `return for`. The array produced by the loop becomes the function result.

{{< editor lang="fql" >}}
func doubleAll(items) {
    return for item in items {
        return item * 2
    }
}

return doubleAll([1, 2, 3])
{{</ editor >}}

This uses the same loop result as a parenthesized `for`; it does not add another wrapper. `return distinct for` applies the normal return-level deduplication directly to the loop result.

A final standalone loop is not promoted into a function result. A collecting loop still executes, but its array is discarded; a returnless braced loop executes without creating an array. In both cases the function falls through with `none`:

{{< code lang="fql" >}}
func processAll(items) {
    for item in items {
        process(item)
    }
}
{{</ code >}}

Likewise, `func value() { 42 }` evaluates `42` as an expression statement and returns `none`. Write `func value() => 42` or `func value() { return 42 }` to return the value. Use `return for` when the function should return a loop's collected array.

## Parameters

Parameters are listed inside parentheses, separated by commas.

{{< editor lang="fql" >}}
func fullName(first, last) => concat(first, " ", last)

return fullName("Ada", "Lovelace")
{{</ editor >}}

A function may have no parameters:

{{< editor lang="fql" >}}
func now() => "2024-01-01"

return now()
{{</ editor >}}

Parameters are positional. The caller must provide exactly the number of arguments the function expects.

## Capturing outer variables

A function body can read variables from the enclosing scope.

{{< editor lang="fql" >}}
let base = 10

func add(value) => base + value

return add(5)
{{</ editor >}}

If the outer variable is declared with `var`, the function can also modify it:

{{< editor lang="fql" >}}
var counter = 0

func inc() {
    counter = counter + 1
    return counter
}

return [inc(), inc(), inc()]
{{</ editor >}}

Variables declared with `let` are immutable and cannot be reassigned inside a function.

## Nesting functions

Functions can be declared inside other functions.

{{< editor lang="fql" >}}
func process(items) {
    func transform(item) => item * 2

    return (
        for item in items {
            return transform(item)
        }
    )
}

return process([1, 2, 3])
{{</ editor >}}

A nested function can access variables from all enclosing scopes, not just the immediately surrounding one.

## Using functions in loops

User-defined functions work naturally with `for` loops and other query constructs.

{{< editor lang="fql" >}}
func formatUser(user) {
    return {
        label: concat(user.name, " (", user.role, ")"),
        active: user.active
    }
}

let users = [
    { name: "Ada", role: "admin", active: true },
    { name: "Grace", role: "editor", active: false },
    { name: "Linus", role: "viewer", active: true }
]

return for user in users {
    filter user.active
    return formatUser(user)
}
{{</ editor >}}

## Function names

Function names follow the same rules as variable names: they must start with a letter or underscore, followed by any combination of letters, digits, and underscores.

Function names are case-sensitive. `add` and `Add` are different functions.

{{< editor lang="fql" >}}
func a() => 1
func A() => 2

return a() + A()
{{</ editor >}}

Built-in and host functions are documented with canonical lowercase names. User-defined function names remain case-sensitive and may use the style preferred by the script author.

## Next steps

{{< docs-related tiles="language-functions-modules,embedding-custom-functions,stdlib" >}}
