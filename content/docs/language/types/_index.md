---
title: "Values and Types"
sidebarTitle: "Values and Types"
weight: 40
draft: false
description: "An overview of FQL values, basic types, host values, and value capabilities."
aliases:
    - /docs/fql/data-types/
---

# Values and Types

FQL is built around values. Every expression produces one, every variable holds one, and every script returns one. A value can be as simple as a number or a string, as structured as an array or object, or as rich as a host value provided by a Ferret runtime, module, or embedding application.

## Value model

To understand FQL values, it helps to separate three related ideas.

| Idea | Description |
| --- | --- |
| Basic values | Values built into FQL itself: `NONE`, booleans, numbers, strings, arrays, objects, and binary values. |
| Host values | Values provided by a runtime, module, or embedding application, such as documents, elements, responses, files, cursors, or application-specific objects. |
| Capabilities | Behaviors a value supports, such as property access, iteration, querying, dispatching actions, or serialization. |

Basic values and host values describe what a value *is*. Capabilities describe what a value can *do*.

For example, an array is a basic value that supports iteration. An object is a basic value that supports property access. A document returned by a runtime may be a host value that supports property access, querying, or other runtime-defined behavior.

The category a value belongs to does not fully determine what operations are available. The value itself decides which capabilities it exposes.

## Basic values

Basic values are the ordinary, built-in values of the language: the absence marker `NONE`, the booleans `true` and `false`, numbers, strings, arrays, objects, and raw binary data. Most scripts work primarily with basic values, either constructing them directly or extracting them from host values.

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
    active: true,
    roles: ["admin", "editor"]
}

RETURN user.name
{{< /editor >}}

Each type has specific rules for comparison, equality, and how it participates in expressions. These are documented in [Basic Types]({{< ref "basic" >}}).

## Host values

Host values are values defined outside the core language. They are provided by Ferret runtimes, modules, or applications that embed Ferret, and they may represent anything managed by the host environment — an HTML document, a browser element, an HTTP response, a file, a database cursor, or an application-specific object.

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

RETURN {
    title: page.title,
    url: page.url
}
{{< /editor >}}

Here, `page` is not a plain object constructed by the script. It is a host value returned by `WEB::HTML::OPEN(...)`, and it may represent an HTML page, a remote resource, or something else entirely depending on the runtime. From the script author's perspective, it is still a value: it can be assigned to a variable, passed to functions, and accessed through the operations it supports — such as the `title` and `url` properties read above.

This extensibility is central to Ferret's design. Rather than adding a separate built-in type for every external system — browsers, APIs, databases, files, test harnesses — Ferret allows host environments to introduce their own value types that integrate naturally with the rest of the language. A host value participates in the same expressions, assignments, and function calls as any basic value.

Host values are explained in [Host Values]({{< ref "host" >}}).

## Capabilities

Not every value supports every operation, and capabilities are how FQL describes this. A capability is a behavior that a value exposes: property access, iteration, querying, dispatching actions, or serialization. Whether a given value supports a given capability depends on the value itself, not on its category.

Arrays support iteration, which is what the `FOR` statement relies on. Objects support property access, which is what makes dot notation meaningful. A browser element provided by a host runtime may support dispatch, allowing commands to be sent to it:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
DISPATCH "click" IN button
{{< /code >}}

Capabilities are particularly important for host values because different host values may support different operations. A document may expose properties; a cursor may be iterable; a browser element may accept dispatched commands. Understanding which capabilities a value supports tells you what you can do with it in a script.

Capability-specific behavior is documented in [Capabilities]({{< ref "capabilities" >}}).

## Values at the boundary

Inside a script, values may be basic values or host values. At the boundary of a script, plain arrays, objects, strings, numbers, booleans, binary values, and `NONE` are the most portable results:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
FOR product IN products
    RETURN {
        name: product.name,
        price: product.price,
        available: product.stock > 0
    }
{{< /code >}}

Host values are usually most useful while the script is running. When a script returns a host value, the embedding application or runtime decides whether it can be serialized and what representation to use.

## Next steps

{{< docs-related tiles="language-types-basic,language-types-host,language-types-capabilities,language-types-ordering,language-types-serialization,language-control-flow" >}}