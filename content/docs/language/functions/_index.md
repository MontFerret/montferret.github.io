---
title: "Functions"
sidebarTitle: "Functions"
weight: 80
draft: false
description: "Call built-in, module-provided, and runtime-provided functions from FQL expressions."
aliases:
    - /docs/fql/functions/
---

# Functions

Functions are expressions that evaluate input values and produce a result.

A function call consists of a function name followed by a parenthesized argument list. The arguments are expressions, so they may be literals, variables, object or array values, arithmetic expressions, subqueries, or other function calls.

{{< code lang="fql">}}
function_name(argument1, argument2, ...)
{{</ code >}}

For example:

{{< editor lang="fql" >}}
let values = [1, 2, 3, 4]

return length(values)
{{</ editor >}}

Host-function names are case-insensitive in FQL and have one canonical lowercase spelling. For example, `length`, `LENGTH`, and `LeNgTh` resolve to the same host function. User-defined functions retain their own case-sensitive rules.

## Calling functions

A function call may be used anywhere an expression is accepted.

{{< editor lang="fql" >}}
let name = "Ada Lovelace"

return {
    original: name,
    upper: upper(name),
    length: length(name)
}
{{</ editor >}}

Arguments are evaluated before the function is called. The resulting values are then passed to the function in the order in which the arguments appear.

{{< editor lang="fql" >}}
let subtotal = 120
let taxRate = 0.08

return round(subtotal + subtotal * taxRate)
{{</ editor >}}

A function may also be called with the result of another function.

{{< editor lang="fql" >}}
let message = "  hello ferret  "

return upper(trim(message))
{{</ editor >}}

Nested calls are ordinary expressions. They should be used only where the resulting expression remains clear enough for the query being written.

## Function arguments

Function arguments are separated by commas.

{{< editor lang="fql" >}}
return concat("ferret", "-", "lang")
{{</ editor >}}

The number and type of accepted arguments depend on the function. Some functions require a fixed number of arguments. Others accept optional arguments or a variable number of arguments.

For example, a function may require a string:

{{< editor lang="fql" >}}
return lower("FERRET")
{{</ editor >}}

Another function may accept an array:

{{< editor lang="fql" >}}
return first([10, 20, 30])
{{</ editor >}}

If a function receives an unsupported number of arguments or a value of an unsupported type, the query fails with an error.

## Return values

Every function call produces a value.

The returned value may be any FQL value: `none`, a `boolean`, `number`, `string`, `array`, `object`, `binary` value, or host value.

{{< editor lang="fql" >}}
let numbers = [4, 8, 15, 16, 23, 42]

return {
    count: length(numbers),
    first: first(numbers),
    hasValue: 15 in numbers
}
{{</ editor >}}

A function may return `none` when no value is available, when a lookup fails, or when the function is defined to represent absence that way. The exact behavior depends on the function.

{{< editor lang="fql" >}}
let values = []

return first(values)
{{</ editor >}}

## Built-in functions

FQL provides built-in functions for common operations on values such as `strings`, `arrays`, `objects`, `numbers`, and `types`.

Built-in functions are available at the top level, without a namespace prefix.

{{< editor lang="fql" >}}
let tags = ["docs", "fql", "runtime"]

return {
    count: length(tags),
    first: first(tags),
    joined: concat_separator(", ", tags)
}
{{</ editor >}}

The available built-in functions are documented in [the standard library reference]({{% ref "../../standard-library" %}}).

## Function calls in expressions

Because function calls are expressions, they can be combined with other expression forms.

They can be used in let declarations:

{{< editor lang="fql" >}}
let normalized = lower(trim("  Ferret  "))

return normalized
{{</ editor >}}

They can be used in object projections:

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
    email: "ADA@example.com"
}

return {
    name: user.name,
    email: lower(user.email)
}
{{</ editor >}}

They can be used in filters:

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

return for user in users {
    filter user.active == true and length(user.name) >= 4
    return user.name
}
{{</ editor >}}

They can also be used in array inline expressions:

{{< editor lang="fql" >}}
let names = [" Ada ", " Grace ", " Linus "]

return names[* return trim(.)]
{{</ editor >}}

Inside an inline array expression, the current element is accessed with `.`.

## Dynamic values

Some functions operate on values whose exact type is known only at runtime. This commonly happens when values come from bind parameters, host applications, modules, external data, or runtime capabilities.

{{< editor lang="fql" >}}
return length(@items)
{{</ editor >}}

The query may compile even when the final value of `@items` is not known yet. If the provided value is not supported by the function at runtime, the query fails with an error.

This allows the same query text to be reused with different input values while still preserving function-specific validation when the query is executed.

## Errors

A function call can fail for several reasons:

- the function name does not exist;
- the function is not available in the current scope;
- the function receives too few or too many arguments;
- an argument has an unsupported type;
- a runtime capability required by the function is unavailable;
- the function itself reports an error while executing.

For example, a function that expects an array may reject a string:

{{< editor lang="fql" >}}
return first("not an array")
{{</ editor >}}

Errors are reported by the compiler when they can be detected statically. Errors that depend on runtime values are reported during execution.

## Function purity and runtime effects

Many functions are pure value operations. They compute a result from their arguments and do not modify external state.

Other functions may interact with runtime systems. A function provided by a module may read from an external source, perform I/O, query a database, or return a capability-backed value.

The behavior of such functions is defined by the module that provides them. When a function depends on runtime capabilities, the same query may behave differently depending on the runtime configuration, available permissions, and provided host environment.

## Next steps

{{< docs-related tiles="language-functions-user-defined,language-functions-modules,stdlib,embedding-custom-functions" >}}
