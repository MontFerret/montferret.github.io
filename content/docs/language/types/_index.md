---
title: "Values and Types"
sidebarTitle: "Values and Types"
weight: 30
draft: false
description: "An overview of FQL values, built-in types, runtime-backed values, and capabilities."
aliases:
    - /docs/fql/data-types/
---

# Values and Types

FQL works with values.

A value can be something simple, such as a number, string, boolean, or NONE. It can be structured data, such as an array or object. It can also be a runtime-backed value provided by a Ferret runtime, module, or embedding application.

This is an important part of Ferret's design: FQL is not limited to JSON-like data.

A script can work with ordinary values:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    active: true,
    roles: ["admin", "editor"]
}  

RETURN user.name
{{< /code >}}

But it can also work with values that represent external resources:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org") 

RETURN doc.title
{{< /code >}}

In this example, `doc` is not a plain object created by the script. It is a runtime-backed value returned by the `DOCUMENT` function. 
From the script author’s perspective, it can still be accessed like a normal FQL value.
The `title` value may be backed by an HTML or browser runtime. It may expose properties such as `text`, `html`, or `attributes`, depending on the module that created it.


## Value categories

FQL values can be grouped by where they come from and by what they can do.

Some values are built into the language. Other values are provided by modules, runtimes, or host applications. On top of that, values may support different capabilities, such as property access, iteration, dispatch, mutation, serialization, or cleanup.

| Concept | Description |
| --- | --- |
| Basic values | Built-in values such as `NONE`, booleans, numbers, strings, arrays, objects, and binary data. |
| Host values | Values provided by modules, runtimes, or applications that embed Ferret. |
| Capabilities | Behaviors a value may support, such as property access, iteration, dispatch, mutation, serialization, or cleanup. |

These concepts can overlap.

For example, an array is a basic value, and it also supports iteration. An object is a basic value, and it supports property access. A document returned by `DOCUMENT(...)` is a host value, and it may support property access through fields such as `title` or `url`.

A cursor returned by a database module could be a host value that supports iteration. A browser element could be a host value that supports property access and event dispatch.

In other words, **basic values** and **host values** describe what a value is, while **capabilities** describe what a value can do.

## Basic values

Basic values are the values built into FQL itself.

They include:

| Type | Example | Description |
| --- | --- | --- |
| none | NONE | Represents an absent value. |
| boolean | true, false | Represents a truth value. |
| integer | 42 | Represents a whole number. |
| float | 3.14 | Represents a floating-point number. |
| string | "hello" | Represents text. |
| array | [1, 2, 3] | Represents an ordered sequence of values. |
| object | { name: "Ada" } | Represents named fields. |
| binary | module-specific | Represents raw bytes, such as files or response bodies. |

Basic values are the values you use most often when writing ordinary FQL scripts.

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET product = {
    name: "Keyboard",
    price: 149.99,
    tags: ["hardware", "input"],
    available: true
}

RETURN product
{{< /code >}}

See Basic Values for details.

## NONE and absence

`NONE` is Ferret's equivalent of `null`, `nil`, or `None` in other languages.

It represents the absence of a value.

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN NONE
{{< /code >}}

Like `null` or `nil`, `NONE` is not the same as `false`, `0`, an empty string, an empty array, or an empty object.

That distinction is useful in extraction and automation workflows, where a missing value can mean something different from an empty value.

For example, a product may have no discount value:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN {
    name: "Keyboard",
    discount: NONE
}
{{< /code >}}

That is different from a discount value of `0`, which means the discount exists but its amount is zero.

## Structured values

Arrays and objects are the main structured values in FQL.

Arrays hold ordered values:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET names = ["Ada", "Grace", "Edsger"]  

RETURN names[0]
{{< /code >}}

Objects hold named values:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",     
    address: {
        city: "London"
    } 
}  

RETURN user.address.city
{{< /code >}}

Structured values are used heavily in FQL because scripts often extract, reshape, and return structured data.

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET products = [
    { name: "Keyboard", price: 149.99, stock: 10 },
    { name: "Mouse", price: 49.99, stock: 0 },
    { name: "Monitor", price: 299.99, stock: 5 }
]

FOR product IN products
    RETURN {
        name: product.name,
        price: product.price,
        available: product.stock > 0
    }
{{< /code >}}

See `Arrays` and `Objects` for details.

## Host values

Host values are values provided by modules, runtimes, or applications that embed Ferret.

They behave like FQL values, but they may represent things that live outside the script itself: an HTML document, a browser page, an element, a network response, a file, a cursor, a message, or an application-specific object.

For example:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org")

RETURN doc.title
{{< /code >}}

In this example, `doc` is a host value returned by the `DOCUMENT` function. The script can access `doc.title`, but `doc` does not have to be a plain object created in the script.

The module or host application decides what the value represents and what behavior it exposes.

This allows Ferret to work with external systems without forcing every external resource to become a plain object first.

## Capabilities

Different values can do different things.

Some values only carry data. For example, a number can be compared, copied, returned, or used in arithmetic, but it does not have fields and cannot be iterated.

Other values expose additional behavior to FQL. These behaviors are called **capabilities**.

A capability describes what operations a value supports. Conceptually, it is like an interface between a value and the Ferret runtime. When FQL code performs an operation, the runtime checks whether the value supports the capability required by that operation.

For example, `FOR ... IN` requires the value to be iterable. Dot access requires the value to expose properties. Dispatch requires the value to accept events or commands.

This is useful because FQL does not need a separate built-in type for every possible external resource. Instead, modules and host applications can provide values that support the behaviors that make sense for that resource.

| Capability | Enables |
| --- | --- |
| Property access | Reading fields with `.name` or `[name]`. |
| Iteration | Using the value in `FOR ... IN`. |
| Mutation | Updating the value or its fields, if supported. |
| Dispatch | Sending an event or command to the value. |
| Serialization | Returning the value as output. |
| Cleanup | Releasing host-side resources after execution. |

Basic values can have capabilities too.

For example, arrays support iteration:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
FOR item IN [1, 2, 3]

RETURN item * 2
{{< /code >}}

Objects support property access:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    active: true
}

RETURN user.name
{{< /code >}}

Host values can expose capabilities as well.

For example, a document value may expose fields such as `title` or `url`:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org")

RETURN doc.title
{{< /code >}}

A cursor returned by a database module may support iteration:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET rows = QUERY "SELECT id, name FROM users" IN db USING sql

FOR row IN rows
    RETURN row.name
{{< /code >}}

A browser element may support dispatch, allowing the script to send an event or command to it:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
DISPATCH "click" IN button
{{< /code >}}

The exact capabilities depend on the value and on the module, runtime, or host application that provides it.

This means two values may both be host values but support different operations. A document may support property access. A cursor may support iteration. A browser element may support property access and dispatch. A binary response body may be serializable but not iterable.

When an operation is used on a value that does not support it, Ferret reports an error unless the script handles that error explicitly.

See [Value Capabilities](../value-capabilities/) for details.

## Host values

Host values are values provided by modules, runtimes, or applications that embed Ferret.

They are useful when a script needs to work with something that is not just plain FQL data. A host value may represent an external resource, a handle to something managed outside the script, or a domain-specific object provided by the application.

Examples include:

- an HTML document
- a browser page
- a DOM element
- an HTTP response
- a file or blob
- a database record
- a database cursor
- a queue message
- an application-specific domain object

From the script author's perspective, a host value is still a value. It can be assigned to a variable, passed to a function, returned from a function, or used in expressions.

For example, DOCUMENT(...) may return a document value:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org")

RETURN doc.title
{{< /code >}}

In this example, doc is not a plain object created by the script. It is a value provided by a module or runtime. The script can still access doc.title because the document value exposes that property to FQL.

Host values can support different capabilities depending on what they represent.

A database cursor may support iteration:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
FOR row IN DB::QUERY("SELECT id, name FROM users")
RETURN row.name
{{< /code >}}

An HTTP response may support property access:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET response = HTTP::GET("https://example.com/api/products")

RETURN response.status
{{< /code >}}

A browser element may support dispatch, allowing the script to send an event or command to it.

The exact behavior depends on the module, runtime, or host application that provides the value.

Host values are one of Ferret's main extension points. They allow the same language to work across static documents, browser automation, APIs, tests, embedded applications, and custom runtimes without adding a separate built-in type for every external system.

See Host Values for details.

## Type checks

FQL values have runtime types.

In most scripts, you do not need to declare types explicitly. Values are created by literals, expressions, functions, modules, or the host application.

fql LET name = "Ada" LET count = 10 LET active = true

Functions and operators may expect certain kinds of values. If a value has the wrong type for an operation, Ferret reports an error unless the script handles it.

fql RETURN "price: " + 10

Depending on the operator and runtime rules, this may either be supported through conversion or rejected as a type error.

When writing reusable scripts, it is often better to make assumptions visible:

fql LET price = product.price  RETURN price == NONE ? NONE : price * 1.2

## Values as output

A script returns a value.

That value may be a primitive value:

fql RETURN 42

An object:

fql RETURN {     name: "Ada",     active: true }

An array:

fql RETURN [     { name: "Ada" },     { name: "Grace" } ]

Or a value provided by a module or runtime.

When a value is returned from a script, Ferret needs to serialize it into a result format. Basic values can usually be serialized directly. Runtime-defined values decide how they should appear in output, or whether they can be returned directly at all.

For extraction workflows, scripts usually return plain arrays and objects because they are easy to consume by other tools.

fql FOR product IN products     RETURN {         name: product.name,         price: product.price,         url: product.url     }

Runtime-backed values are often most useful inside the script, while plain values are usually preferred at the boundary.

## Where to go next

Use the following pages depending on what you want to learn:

- Basic Values explains NONE, booleans, numbers, strings, arrays, objects, and binary values.
- Arrays and Objects explains structured data, nested access, dynamic property access, and object shorthand.
- Value Capabilities explains property access, iteration, dispatch, mutation, serialization, and cleanup.
- Runtime-defined Values explains values provided by runtimes, modules, and embedded applications.