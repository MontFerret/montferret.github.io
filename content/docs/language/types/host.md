---
title: "Host Values"
sidebarTitle: "Host Values"
weight: 30
draft: false
description: "How FQL represents runtime-provided values such as database connections, browser pages, files, clients, and other host-managed resources."
---

# Host Values

Host values are values provided by a runtime, module, or embedded application.

Unlike basic FQL values such as numbers, strings, arrays, and objects, a host value is backed by implementation-specific data outside the FQL value model. FQL can hold the value, pass it to functions, return it from expressions, and use it with operations supported by its capabilities.

Common examples include database connections, browser pages, HTML documents, HTTP clients, files, streams, and other resources managed outside FQL.

## What is a host value?

A host value represents something owned by the host environment rather than by FQL itself.

For example, a database connection value may refer to an open connection managed by a Go module. A browser page value may refer to a live browser tab managed by the browser runtime. FQL does not inspect or copy those values directly. Instead, it works with them through functions, operators, and runtime capabilities.

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")

RETURN QUERY `
  SELECT name
  FROM users
` IN db
{{</ code >}}

In this example, `db` is a host value. It is not an object containing database fields. It is a runtime-provided value that can be used by operations which understand database connections.

## Using host values

Host values can participate in normal expression flow.

They can be assigned to variables:

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")
{{</ code >}}

They can be passed to functions:

{{< code lang="fql" >}}
RETURN DB::TABLES(db)
{{</ code >}}

They can be grouped with related data while a query is running:

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")

LET context = {
  source: db,
  table: "users"
}

RETURN QUERY `
  SELECT *
  FROM users
` IN context.source
{{</ code >}}

However, the internal data of a host value is opaque to FQL. A host value does not automatically expose object fields, array elements, or scalar operations unless the runtime explicitly supports them.

## Opaque by default

Host values are opaque by default.

FQL can hold and pass a host value, but it cannot assume anything about the value's internal implementation. The supported operations depend on the value's runtime type and capabilities.

Some host values expose capabilities that make them usable with familiar FQL syntax. For example, an HTML node may support property access, allowing fields such as `text`, `attrs`, or other module-defined properties to be read using normal property syntax.

Other host values may expose no properties at all and can only be used through functions or specific language operations.

{{< code lang="fql" >}}
LET doc = HTML::PARSE("<a href='https://example.com'>Example</a>")

RETURN doc.root.text
{{</ code >}}

In this example, property access works because the HTML host value explicitly supports it. FQL is not accessing the host value's internal implementation directly; it is using a capability exposed by that value.

## Capabilities

Some host values expose capabilities.

A capability describes an operation that a value supports. For example, a database connection may support query execution, which allows it to be used with `QUERY ... IN`.

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")

RETURN QUERY `
  SELECT *
  FROM users
` IN db
{{</ code >}}

The `QUERY ... IN db` expression works because `db` provides the capability required by query execution.

FQL does not need to know how the database connection is implemented internally. It only needs to know that the value supports the required operation.

See [Capability Types]({{< ref "capabilities" >}}) for the full capability model.

## Lifetime and ownership

Some host values represent external resources with a lifetime, such as open files, network connections, database handles, browser sessions, cursors, or runtime clients.

Those resources are owned by the host environment, module, or runtime that created them. FQL can receive and use host values during program execution, but resource cleanup is handled by the runtime according to the capabilities exposed by those values.

At the end of program execution, before the final result is returned to the host, Ferret finalizes the result. During finalization, resource-backed values may be closed and the result is converted into the output format selected by the host.

Host values that support the close capability are closed automatically during this step. This is similar to values that implement `io.Closer` in Go: the value may hold an external resource, and the runtime is responsible for releasing that resource when the program result is finalized.

After cleanup, the final result is serialized using the encoding configured by the host. By default, Ferret serializes results as JSON.

This means FQL code can work with host values while the program is running, but the value returned to the outside world must still be converted into the host-selected output format.

## Result serialization

Host values are regular FQL values during execution, but they are not always serializable as data.

When a program returns a result, Ferret converts that result into the output format selected by the host. The default output format is JSON, but embedded hosts may choose a different encoding.

{{< code lang="fql" >}}
RETURN DB::SQLITE::OPEN("data.db")
{{</ code >}}

This expression returns a host value during execution. When the program result is finalized, the host decides how that value should be represented in the selected output format.

A host value may be serialized using a descriptive representation, converted using a host-defined rule, or rejected if it cannot be represented safely in the selected format.

See [Serialization]({{< ref "serialization" >}}) for more information about how FQL values are converted to external formats.

## Equality and ordering

Host values may define their own comparison behavior.

FQL does not compare host values by inspecting their internal implementation. Instead, comparison is delegated to the host value when it supports comparison.

A comparable host value defines how it is ordered relative to other values, including values of other types. That ordering may be based on identity, an internal key, a normalized representation, or another rule chosen by the runtime or module.

If a host value does not support comparison, it should not be used with ordering operations such as sorting or range comparison.

This keeps host values opaque while still allowing modules and runtimes to make specific host values participate in FQL's ordering model.

See [Type Ordering]({{< ref "ordering" >}}) and [Capability Types]({{< ref "capabilities" >}}) for more information about value ordering.
