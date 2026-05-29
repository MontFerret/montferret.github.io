---
title: "Syntax"
weight: 20
draft: false
description: "Learn the basic structure of an FQL script."
aliases:
  - /docs/fql/syntax/
---

# Script Structure

An FQL script is a sequence of statements that produce a final value.

Most scripts follow a simple shape:

1. declare or receive input
2. transform data
3. return a result

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

The final value may be produced by a top-level `RETURN`, or by another terminal statement such as a top-level `FOR` expression-style block.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
FOR role IN ["admin", "editor", "viewer"]
    RETURN UPPER(role)
{{< /editor >}}

In both cases, the script produces a value that can be returned to the CLI, a host application, a test runner, or another runtime.

## Statements

A script is made of statements.

Common statements include:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET name = "Ada"

RETURN name
{{< /editor >}}

`LET` creates an immutable binding.

`RETURN` produces the final result of the script.

When a value needs to change over time, FQL also supports `VAR`:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
VAR count = 0
count = count + 1

RETURN count
{{< /editor >}}

Most scripts should prefer `LET` unless mutation is actually needed.

More advanced scripts may also use statements such as `FOR`, `FILTER`, `COLLECT`, `WAITFOR`, `DISPATCH`, or function declarations, depending on what the script needs to do.

## Expressions
Most of the useful work in FQL happens inside expressions. Expressions can be simple values, function calls, or complex combinations of operators and nested expressions.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN 1 + 2
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN UPPER("hello")
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN (1 + 2) * 3
{{< /editor >}}

Expressions can be assigned to variables and used in subsequent statements:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET user = {
    firstName: "Ada",
    lastName: "Lovelace"
}
LET fullName = user.firstName + " " + user.lastName

RETURN fullName
{{< /editor >}}

They can also be nested inside larger expressions:

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

Expressions are the building blocks of FQL scripts. Statements describe the flow of the script, while expressions produce the values that move through that flow.

## Producing a result
Every script that is meant to produce a result needs a terminal statement.

The most common terminal statement is `RETURN`.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN "Hello, world!"
{{< /editor >}}

The returned value can be any FQL value:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN 42
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN ["a", "b", "c"]
{{< /editor >}}

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN {
    ok: true,
    items: []
}
{{< /editor >}}

A `FOR` loop can also be terminal when the script produces a collection.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
FOR i IN 1..10
    RETURN i * i
{{< /editor >}}

In this form, the loop itself produces the script result.
This is useful when the script is primarily a projection, extraction, or transformation over a collection.

You can also declare intermediate values before and inside the loop to support more complex logic:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET products = [
    { name: "Widget", price: 9.99, active: true },
    { name: "Gadget", price: 19.99, active: false },
    { name: "Doohickey", price: 14.99, active: true }
]
LET discount = 0.1

FOR product IN products
    LET discountedPrice = product.price * (1 - discount)
    FILTER product.active AND discountedPrice < 15
    LET productName = product.name + " (discounted)"
    RETURN {
        name: productName,
        price: discountedPrice
    }
{{< /editor >}}

## Top-level flow

FQL scripts are usually written from top to bottom. The flow of execution follows the order of statements, with each statement able to reference values declared in previous statements.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET price = 19.99
LET tax = price * 0.08

RETURN price + tax
{{< /editor >}}

A statement cannot use a value before it has been declared:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET total = price + tax
LET price = 19.99
LET tax = price * 0.08

RETURN total
{{< /editor >}}

## Bindings

Bindings created by `LET` and `VAR` are immutable and mutable, respectively. They can be used in subsequent statements but cannot be re-declared.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET name = "Ada"
LET name = "Grace"  // Error: name is already declared

RETURN name
{{< /editor >}}

A binding lets you give a meaningful name to an intermediate value. This can make your script easier to read and maintain, especially when the value is used multiple times.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET price = 19.99
LET tax = price * 0.08
LET total = price + tax
LET message = "The total price is $" + ROUND(total)

RETURN message
{{< /editor >}}

## Parameters

Parameters are a special kind of binding that receive their value from outside the script. They are variables prefixed with `@`.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN "Hello, " + @name + "!"
{{< /editor >}}

For example, a host application or CLI command can pass name at runtime instead of hardcoding it in the script.

Parameters are often used as the starting point of a script:

{{< code lang="fql" >}}
FOR product IN @products
    FILTER product.inStock
    RETURN {
        title: product.title,
        price: product.price
    }
{{< /code >}}

## Structured results

FQL scripts can return any valid FQL value, including complex objects and collections.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN {
    user: {
        name: "Ada",
        roles: ["admin", "editor"]
    },
    stats: {
        posts: 42,
        followers: 1000
    },
    recentActivity: [
        { type: "post", title: "My latest post" },
        { type: "comment", content: "Great article!" }
    ],
    createdAt: NOW()
}
{{< /editor >}}

This is especially useful for extraction and automation workflows, where the script should return predictable output.

## Blocks

Some statements, such as `FOR`, `MATCH`, and function declarations, create a new block scope. This means that variables declared inside the block are not accessible outside of it.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
LET res = (FOR i IN 1..5
    LET square = i * i
    RETURN square
)
RETURN [res, square]  // Error: square is not defined outside the block
{{< /editor >}}

Inside the loop, FILTER controls which items are included, and RETURN describes the value produced for each included item.

Blocks let you describe a local flow of work without leaving the declarative style of the language.

The exact block structure depends on the statement being used. For example, FOR, MATCH, WAITFOR, and function declarations each have their own shape.

## Whitespace
FQL is whitespace-insensitive, which means that you can use spaces, tabs, and newlines to format your code in a way that is most readable to you.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN 1 +     2
{{< /editor >}}

Spaces, tabs, and line breaks separate tokens. Whitespace inside strings is preserved.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN "Hello,     world!"
{{< /editor >}}

## Comments

FQL supports single-line comments that start with `//` and continue to the end of the line.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
// This is a comment
RETURN "Hello, world!"  // This is another comment
{{< /editor >}}

It also supports multi-line comments enclosed in `/*` and `*/`.

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
/*
This is a multi-line comment.
It can span multiple lines.
*/
RETURN "Hello, world!"
{{< /editor >}}

## Names and keywords

Names identify variables, object fields, functions, and other script-level symbols. They must start with a letter or underscore, followed by letters, digits, or underscores.

{{< code lang="fql"  >}}
LET _name = "Ada"
LET name2 = "Grace"
LET Name = "Turing"
{{< /code >}}

Keywords are words with special meaning in FQL.

{{< code lang="fql"  >}}
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

Keywords are case-sensitive and are conventionally written in uppercase.

When an object field has the same name as a keyword, quote the field name:

{{< editor lang="fql" apiVersion="2" orientation="horizontal" >}}
RETURN {
    "return": "This field is named 'return', which is a keyword, so it is quoted."
}
{{< /editor >}}