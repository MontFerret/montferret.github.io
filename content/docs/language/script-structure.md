---
title: "Script structure"
weight: 30
draft: false
description: "Learn the basic structure of an FQL script."
aliases:
    - /docs/fql/syntax/
---

## Script structure

An FQL script is a sequence of statements that runs from top to bottom and produces a final value.

Statements provide the outer structure of a script: they introduce bindings, iterate collections, and decide what value is returned. Inside those statements, expressions do most of the work. They produce the values that get assigned, returned, filtered, queried, or composed.

## Scripts

A script is a top-level sequence of statements. Statements are evaluated in order, and each statement can reference names declared by earlier ones. A name cannot be used before it is declared.

Most scripts follow a simple pattern: declare or receive input, transform data, return a result.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

LET isAdmin = CONTAINS(user.roles, "admin")

RETURN {
    name: user.name,
    isAdmin: isAdmin
}
{{< /editor >}}

## Script headers

A script can begin with one or more `USE` declarations. They create local aliases for namespaces or namespaced functions and must appear before the script body. See the [`USE` statement reference]({{< ref "/docs/language/use" >}}) for the syntax and resolution rules.

## Statements

A statement describes a step in the script. Some statements create bindings; others produce results or control iteration.

`LET` creates an immutable binding:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET name = "Ada"

RETURN name
{{< /editor >}}

`VAR` creates a mutable binding: one whose value can be reassigned later in the same scope:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
VAR total = 0

FOR price IN [10, 20, 30]
    total = total + price
    
    RETURN total
{{< /editor >}}

Only `VAR` bindings can be reassigned. `LET` bindings cannot be changed after they are created, and no binding can be declared twice in the same scope. Prefer `LET` unless mutation is actually needed.

More advanced scripts may also use `FOR`, `FILTER`, `COLLECT`, `MATCH`, `WAITFOR`, `DISPATCH`, `DO WHILE`, or function declarations. Those constructs are covered in their own pages.

## Expressions

Most of the useful work in FQL happens inside expressions. An expression is any piece of syntax that produces a value: a literal, a function call, an arithmetic combination, a field access, a query, or a nested `FOR`. Expressions can be assigned to bindings, passed as arguments, or returned directly.

Simple literals and arithmetic:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN (1 + 2) * 3
{{< /editor >}}

Function calls:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN UPPER("hello")
{{< /editor >}}

Object and array construction:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN {
    name: "Ada",
    active: true,
    score: 42,
    tags: ["admin", "editor"],
    missingValue: NONE
}
{{< /editor >}}

Expressions can also be composed. The output of one becomes the input of another:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

RETURN {
    name: user.name,
    roleCount: LENGTH(user.roles)
}
{{< /editor >}}

FQL is dynamically typed. Values carry their type at runtime, and operations expect compatible types: arithmetic works on numbers, field access works on objects, collection operations expect arrays or other iterable values.

Statements describe the flow of the script. Expressions produce the values that move through that flow.

## Returning a result

Every FQL script must end in a terminal statement.

A terminal statement produces the script result. FQL currently has two terminal forms:

- `RETURN`, which returns a single value
- a top-level `FOR` statement, which iterates over a collection and returns the produced collection

The returned value can be any FQL value: `NONE`, a boolean, number, string, array, object, binary value, or host value.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN "Hello, world!"
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
FOR i IN 1..10
    RETURN i * i
{{< /editor >}}

The loop body is the same in both cases. The difference is what happens to the produced collection: a top-level `FOR` makes it the script result; a parenthesized `FOR` stores it for further use.

## Scopes and blocks

Some statements introduce a nested scope. Names declared inside that scope are not visible outside it.

`FOR` is the most common block-producing statement:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET values = (FOR i IN 1..5
    LET square = i * i
    RETURN square
)

RETURN values
{{< /editor >}}

`square` exists only inside the `FOR` block. Referencing it outside the block is an error.

Other statements have block-like shapes as well. A `MATCH` expression describes branching logic:

{{< editor lang="fql" params=`{"status": "active"}` >}}
RETURN MATCH @status (
    "active" => "Account is active",
    "paused" => "Account is paused",
    _ => "Unknown status"
)
{{< /editor >}}

A `WAITFOR` block describes event-oriented runtime logic:

{{< code lang="fql" >}}
WAITFOR EVENT network.response_received
    WHEN event.status == 200
RETURN event.url
{{< /code >}}

A function declaration creates a reusable local function:

{{< editor lang="fql" >}}
FUNC fullName(user) => user.firstName + " " + user.lastName

RETURN fullName({ firstName: "Ada", lastName: "Lovelace" })
{{< /editor >}}

Each of these has its own detailed rules and is covered in its dedicated documentation — `MATCH` and `WAITFOR` in [Control Flow]({{% ref "control-flow" %}}), and function declarations in [Functions]({{% ref "functions" %}}). This section shows the structural shape only.

## Comments

FQL supports single-line comments that begin with `//` and extend to the end of the line:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
// This is a comment
RETURN "Hello, world!"  // This is another comment
{{< /editor >}}

Multi-line comments are enclosed in `/*` and `*/`:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
/*
This is a multi-line comment.
It can span multiple lines.
*/
RETURN "Hello, world!"
{{< /editor >}}

FQL is whitespace-insensitive. Spaces, tabs, and newlines separate tokens but do not affect semantics. Whitespace inside strings is preserved.

## Names and keywords

Names identify variables, object fields, functions, and other script-level symbols. A name must start with a letter or underscore, followed by any combination of letters, digits, and underscores:

{{< code lang="fql" >}}
LET _name = "Ada"
LET name2 = "Grace"
LET Name = "Turing"
{{< /code >}}

Keywords are reserved words with special meaning in FQL. They are case-sensitive and conventionally written in uppercase. The full set of reserved keywords is:

{{< code lang="fql" >}}
USE
AS
MATCH
WHEN
FUNC
FOR
RETURN
QUERY
USING
WAITFOR
DISPATCH
OPTIONS
TIMEOUT
EVERY
BACKOFF
JITTER
EXISTS
COUNT
ONE
DISTINCT
FILTER
SORT
LIMIT
LET
VAR
COLLECT
ASC
DESC
AT
LEAST
INTO
KEEP
WITH
ALL
ANY
AGGREGATE
EVENT
LIKE
NOT
IN
DO
WHILE
AND
OR
ON
ERROR
FAIL
RETRY
DELAY
DELETE
VALUE
{{< /code >}}

When an object field shares its name with a keyword, quote the field name:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN {
    "return": "This field is named 'return', which is a keyword, so it is quoted."
}
{{< /editor >}}

## Next steps

{{< docs-related tiles="language-use,language-variables,language-expressions,language-operators" >}}
