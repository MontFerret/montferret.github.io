---
title: "Script Structure"
weight: 30
draft: false
description: "Learn the basic structure of an FQL script."
aliases:
    - /docs/fql/syntax/
---

# Script Structure

An FQL script is a sequence of statements that runs from top to bottom and produces a final value.

Statements provide the outer structure of a script: they introduce bindings, iterate collections, and decide what value is returned. Inside those statements, expressions do most of the work. They produce the values that get assigned, returned, filtered, queried, or composed.

## Scripts

A script is a top-level sequence of statements. Statements are evaluated in order, and each statement can reference names declared by earlier ones. A name cannot be used before it is declared.

Most scripts follow a simple pattern: declare or receive input, transform data, return a result.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
let user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

let isAdmin = contains(user.roles, "admin")

return {
    name: user.name,
    isAdmin: isAdmin
}
{{< /editor >}}

## Script headers

A script can begin with one or more `use` declarations. They create local aliases for namespaces or namespaced functions and must appear before the script body. See the [`use` statement reference]({{< ref "/docs/language/script-structure/use" >}}) for the syntax and resolution rules.

## Statements

A statement describes a step in the script. Some statements create bindings; others produce results or control iteration.

`let` creates an immutable binding:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
let name = "Ada"

return name
{{< /editor >}}

`var` creates a mutable binding: one whose value can be reassigned later in the same scope:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
var total = 0

for price in [10, 20, 30] {
    total = total + price

    return total
}

return total
{{< /editor >}}

Only `var` bindings can be reassigned. `let` bindings cannot be changed after they are created, and no binding can be declared twice in the same scope. Prefer `let` unless mutation is actually needed.

More advanced scripts may also use `for`, `filter`, `collect`, `match`, `waitfor`, `dispatch`, `do while`, or function declarations. Those constructs are covered in their own pages.

## Expressions

Most of the useful work in FQL happens inside expressions. An expression is any piece of syntax that produces a value: a literal, a function call, an arithmetic combination, a field access, a query, or a nested `for`. Expressions can be assigned to bindings, passed as arguments, or returned directly.

Simple literals and arithmetic:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return (1 + 2) * 3
{{< /editor >}}

Function calls:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return upper("hello")
{{< /editor >}}

Object and array construction:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return {
    name: "Ada",
    active: true,
    score: 42,
    tags: ["admin", "editor"],
    missingValue: none
}
{{< /editor >}}

Expressions can also be composed. The output of one becomes the input of another:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
let user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

return {
    name: user.name,
    roleCount: length(user.roles)
}
{{< /editor >}}

FQL is dynamically typed. Values carry their type at runtime, and operations expect compatible types: arithmetic works on numbers, field access works on objects, collection operations expect arrays or other iterable values.

Statements describe the flow of the script. Expressions produce the values that move through that flow.

## Returning a result

Use `return` to produce a script result. A script can also finish without `return`; after its statements run, it completes successfully with `none`. Completely empty or whitespace-only source is still invalid.

The returned value can be any FQL value: `none`, a boolean, number, string, array, object, binary value, or host value.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return "Hello, world!"
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return for i in 1..10 {
    return i * i
}
{{< /editor >}}

`return for` returns the collection produced by the loop directly. Parenthesized `for` expressions remain useful when the collection must be assigned, nested, or passed to another expression.

A standalone `for` is an ordinary statement. Its body still runs, but its produced collection is discarded:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
var total = 0

for value in [1, 2, 3] {
    total = total + value
    return value
}

return total
{{< /editor >}}

Likewise, an expression used as a block-function statement is evaluated and discarded. Only an explicit `return` or an arrow body determines a function result.

## Scopes and blocks

Some statements introduce a nested scope. Names declared inside that scope are not visible outside it.

`for` is the most common block-producing statement:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
let values = (for i in 1..5 {
    let square = i * i
    return square
}
)

return values
{{< /editor >}}

`square` exists only inside the `for` block. Referencing it outside the block is an error.

Other statements have block-like shapes as well. A `match` expression describes branching logic:

{{< editor lang="fql" params=`{"status": "active"}` >}}
return match @status {
    "active" => "Account is active",
    "paused" => "Account is paused",
    _ => "Unknown status"
}
{{< /editor >}}

A `waitfor` block describes event-oriented runtime logic:

{{< code lang="fql" >}}
waitfor event network.response_received
    when event.status == 200
return event.url
{{< /code >}}

A function declaration creates a reusable local function:

{{< editor lang="fql" >}}
func fullName(user) => user.firstName + " " + user.lastName

return fullName({ firstName: "Ada", lastName: "Lovelace" })
{{< /editor >}}

Each of these has its own detailed rules and is covered in its dedicated documentation — `match` and `waitfor` in [Control Flow]({{% ref "control-flow" %}}), and function declarations in [Functions]({{% ref "functions" %}}). This section shows the structural shape only.

## Comments

FQL supports single-line comments that begin with `//` and extend to the end of the line:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
// This is a comment
return "Hello, world!"  // This is another comment
{{< /editor >}}

Multi-line comments are enclosed in `/*` and `*/`:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
/*
This is a multi-line comment.
It can span multiple lines.
*/
return "Hello, world!"
{{< /editor >}}

FQL is whitespace-insensitive. Spaces, tabs, and newlines separate tokens but do not affect semantics. Whitespace inside strings is preserved.

## Names and keywords

Names identify variables, object fields, functions, and other script-level symbols. A name must start with a letter or underscore, followed by any combination of letters, digits, and underscores:

{{< code lang="fql" >}}
let _name = "Ada"
let name2 = "Grace"
let Name = "Turing"
{{< /code >}}

Keywords are reserved words with special meaning in FQL. They are case-insensitive and conventionally written in lowercase. The full set of reserved keywords is:

{{< code lang="fql" >}}
use
as
match
when
func
for
return
query
using
waitfor
dispatch
options
timeout
every
backoff
jitter
exists
count
one
distinct
filter
sort
limit
let
var
collect
asc
desc
at
least
into
keep
with
all
any
aggregate
event
like
not
in
do
while
and
or
on
error
fail
retry
delay
delete
value
{{< /code >}}

When an object field shares its name with a keyword, quote the field name:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
return {
    "return": "This field is named 'return', which is a keyword, so it is quoted."
}
{{< /editor >}}

## Next steps

{{< docs-related tiles="language-use,language-variables,language-expressions,language-operators" >}}
