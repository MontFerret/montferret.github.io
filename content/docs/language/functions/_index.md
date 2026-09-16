---
title: "Functions"
sidebarTitle: "Functions"
weight: 80
draft: false
description: "Call functions, pass arguments, and use results from script-defined and runtime-provided functions."
aliases:
    - /docs/fql/functions/
---

# Functions

A function call evaluates arguments and produces a value. Write the function name followed by its arguments in parentheses:

{{< editor lang="fql" >}}
let values = [1, 2, 3, 4]

return length(values)
{{</ editor >}}

This calls `length` with the array and returns `4`. Functions can be declared in the script or supplied by the runtime.

## Calling functions

Calls are expressions. You can use their results in bindings, filters, object fields, or other function calls.

{{< editor lang="fql" >}}
let message = "  hello ferret  "
let normalized = upper(trim(message))

return { text: normalized, length: length(normalized) }
{{</ editor >}}

The inner `trim` call produces the argument for `upper`. The returned object contains the normalized text and its length.

## Arguments

Arguments are expressions separated by commas. They are evaluated before the call and passed in the order written.

{{< editor lang="fql" >}}
let subtotal = 120
let taxRate = 0.08

return math::round(subtotal + subtotal * taxRate)
{{</ editor >}}

The arithmetic expression is evaluated before `math::round` receives its value. The `::` separator qualifies the function name with its namespace.

Each function defines its accepted argument counts and types. Some accept optional arguments or a variable number of arguments. For example, `concat` accepts multiple strings:

{{< editor lang="fql" >}}
return concat("ferret", "-", "lang")
{{</ editor >}}

Arguments can also come from bind parameters or host values whose types are known only during execution. Function-specific validation still applies to those values.

## Return values

A function call produces one value, which can be any FQL value, including an array, object, `none`, or host value. A function can use `none` to represent absence; the exact rule depends on the function.

{{< editor lang="fql" >}}
return arrays::first([])
{{</ editor >}}

This returns `none` because the array has no first element. A call used as an expression statement still executes, but its result is discarded. Use `return` when the value should become the script result.

## User-defined functions

Use `func` to declare a function within a script. An arrow body returns the value of its expression:

{{< editor lang="fql" >}}
func double(value) => value * 2

return double(21)
{{</ editor >}}

See [User-defined functions]({{< ref "user-defined" >}}) for block bodies, parameters, scope, and captured variables.

## Namespaced function names

Runtime-provided functions can have a global name such as `length` or a qualified name such as `arrays::first`. Namespace segments are separated by `::` and can be nested.

{{< editor lang="fql" >}}
return arrays::first([10, 20, 30])
{{</ editor >}}

Registered host-function names and namespace segments are case-insensitive. For example, `arrays::first` and `ARRAYS::FIRST` resolve to the same host function. These docs use lowercase names. User-defined functions and local aliases are case-sensitive.

The [`use` statement]({{< ref "/docs/language/script-structure/use" >}}) can create a local alias for a namespace or function. It does not load a module or make an unavailable function available.

## Where functions come from

The runtime determines which host functions a script can call:

- The [Standard Library]({{< ref "/docs/standard-library" >}}) supplies global and namespaced functions. See its reference for signatures, accepted arguments, and return values.
- [Modules]({{< ref "/docs/modules" >}}) provide additional functions and capabilities. Find their APIs and documentation in the [Registry]({{< ref "/registry" >}}).
- An embedding application can register [custom functions]({{< ref "/docs/embedding/go/custom-functions" >}}).

See [Module functions]({{< ref "modules" >}}) for qualified calls, aliases, and runtime availability.

## Errors and runtime effects

Calls can fail when a function is unavailable, an argument count or type is unsupported, or the function reports an execution error. For example, this call intentionally fails because `arrays::first` requires a list:

{{< code lang="fql" >}}
return arrays::first("not an array")
{{</ code >}}

Errors that can be detected statically are reported during compilation. Checks that depend on runtime values or registered functions happen when the runtime resolves or executes the call.

Some functions only compute values. Others mutate an input, perform I/O, or use capabilities supplied by the host. The function's provider defines those effects and requirements; consult its reference before relying on mutation or external access.

## Next steps

{{< docs-related tiles="language-functions-user-defined,language-functions-modules,stdlib,runtime-modules,registry" >}}
