---
title: "Values and Types"
sidebarTitle: "Values and Types"
weight: 30
draft: false
description: "An overview of FQL values, basic types, host values, and value capabilities."
aliases:
- /docs/fql/data-types/
---

# Values and Types

FQL is built around values. Every expression produces one, every variable holds one, and every script returns one. A value can be as simple as a number or a string, as structured as an array or object, or as rich as a host value provided by a Ferret runtime, module, or embedding application.

This last category is an important part of Ferret's design. FQL is not limited to JSON-like data. A script can work with plain data - constructing objects, iterating arrays, computing expressions - but it can also work with values that represent live resources outside the script itself:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org")

RETURN doc.title
{{< /code >}}

Here, `doc` is not a plain object constructed by the script. It is a host value returned by `DOCUMENT(...)`, and it may represent an HTML document, a remote resource, or something else entirely depending on the runtime. From the script author's perspective, it is still a value: it can be assigned to a variable, passed to functions, and accessed through the operations it supports.

## Value categories

FQL values fall into three broad categories.

| Category | Description |
| --- | --- |
| Basic values | Values built into FQL itself: `NONE`, booleans, numbers, strings, arrays, objects, and binary values. |
| Host values | Values provided by a runtime, module, or embedding application, representing external resources such as documents, elements, responses, files, cursors, or application-specific objects. |
| Capabilities | Behaviors a value supports, such as property access, iteration, dispatch, mutation, serialization, or cleanup. |

The distinction between the first two categories and the third is worth noting. Basic values and host values describe what a value *is*. Capabilities describe what a value can *do*. An array is a basic value that supports iteration; an object is a basic value that supports property access; a document returned by `DOCUMENT(...)` is a host value that may expose properties like `title`, `url`, or `body`. The category a value belongs to does not determine its capabilities — that is defined by the value itself.

## Basic values

Basic values are the ordinary, built-in values of the language. They include the absence marker `NONE`, the booleans `true` and `false`, numeric values like `42` or `3.14`, text strings, ordered arrays, named-field objects, and raw binary data. Most scripts work primarily with basic values, either constructing them directly or extracting them from host values:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
name: "Ada",
active: true,
roles: ["admin", "editor"]
}

RETURN user.name
{{< /code >}}

| Type | Example | Description |
| --- | --- | --- |
| `NONE` | `NONE` | Represents an absent or undefined value. |
| `bool` | `true`, `false` | Represents a truth value. |
| `number` | `42`, `3.14` | Represents numeric values, both integer and floating-point. |
| `string` | `"hello"` | Represents text. |
| `array` | `[1, 2, 3]` | Represents an ordered sequence of values. |
| `object` | `{ name: "Ada" }` | Represents a set of named fields. |
| `binary` | module-specific | Represents raw bytes. |

Basic values are covered in detail in [Basic Types](#).

## Host values

Host values are values defined outside the core language. They are provided by Ferret runtimes, modules, or applications that embed Ferret, and they may represent anything managed by the host environment — an HTML document, a browser element, an HTTP response, a file, a database cursor, or an application-specific object.

The script excerpt above, using `DOCUMENT(...)`, illustrates this. The `doc` value may expose fields like `title` or `url` to FQL, but it is not a plain object: its structure and behavior are defined by the module or runtime that created it, not by the language itself.

This extensibility is central to Ferret's design. Rather than adding a separate built-in type for every external system — browsers, APIs, databases, files, test harnesses — Ferret allows host environments to introduce their own value types that integrate naturally with the rest of the language. A host value participates in the same expressions, assignments, and function calls as any basic value.

Host values are covered in detail in [Host Values](#).

## Capabilities

Not every value supports every operation, and capabilities are how FQL describes this. A capability is a behavior that a value exposes: property access, iteration, dispatch, mutation, serialization, or cleanup. Whether a given value supports a given capability depends on the value itself, not on its category.

Arrays, for instance, support iteration. The `FOR` statement relies on this:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
FOR item IN [1, 2, 3]
RETURN item * 2
{{< /code >}}

Objects support property access, which is what makes dot notation meaningful:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = { name: "Ada", active: true }

RETURN user.name
{{< /code >}}

A browser element provided by a host runtime may support dispatch, allowing commands to be sent to it:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
DISPATCH "click" IN button
{{< /code >}}

Capabilities are particularly important for host values because different host values may support different operations. A document may expose properties; a cursor may be iterable; a browser element may accept dispatched commands. Understanding which capabilities a value supports tells you what you can do with it in a script.

Capabilities are covered in detail in [Value Capabilities](#).

## Returning values

Every FQL script returns a value. That value may be a primitive, a structured object, or a collection built from a query:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
FOR product IN products
    RETURN {
        name: product.name,
        price: product.price,
        available: product.stock > 0
    }
{{< /code >}}

For extraction workflows, scripts typically return plain arrays and objects, since these are straightforward for other tools to consume. Host values are most useful inside the script itself — when a host value is returned, it is up to the host to determine whether it can be serialized and what its output representation should be.

## Where to go next

TBD