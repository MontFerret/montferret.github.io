---
title: "Serialization"
sidebarTitle: "Serialization"
weight: 50
draft: false
description: "How Ferret converts runtime values into external representations when returning results to the host."
---

# Serialization

Serialization is the process of converting Ferret runtime values into an external representation that can be returned to the host application, CLI, API client, or embedding environment.

Ferret programs operate on runtime values. When a program finishes, the final result is passed back to the host. Before that result can cross the host boundary, it is serialized into the encoding format selected by the host runtime. The default encoding is JSON, but embedders may choose a different representation.

While the program is running, this is an FQL object containing FQL values. After execution completes, the runtime serializes the result — for example, as a JSON object:

{{< editor lang="fql" >}}
RETURN {
  name: "Ada",
  active: true,
  tags: ["math", "logic"]
}
{{</ editor >}}

## When serialization happens

Serialization happens after program execution completes, before the result is returned to the host.

During execution, values remain Ferret runtime values. Arrays, objects, host values, functions, and other runtime constructs are evaluated according to FQL semantics. Only the final result is serialized.

Temporary values that are created and discarded during execution are not serialized. A value only needs to be serializable if it becomes part of the returned result.

{{< editor lang="fql" >}}
LET items = [1, 2, 3, 4, 5]

LET doubled = (
    FOR item IN items
        RETURN item * 2
)

RETURN doubled[0]
{{</ editor >}}

Here, `items` and `doubled` are intermediate values. Only the final result — a single number — is serialized.

## Serializable values

Basic FQL values have well-defined serialized representations. The exact encoding depends on the output format selected by the host, but the following table shows the typical mapping for JSON:

| FQL value | JSON representation |
| --- | --- |
| `NONE` | `null` |
| `true`, `false` | `true`, `false` |
| integer | number |
| float | number |
| string | string |
| array | array |
| object | object |
| binary | Base64-encoded string |
| date/time | RFC 3339 string |

See [Basic Types]({{< ref "basic" >}}) for details on each value type.

## Arrays and objects

Arrays and objects are serialized recursively. Each element or field value must itself be serializable.

{{< editor lang="fql" >}}
RETURN {
  users: [
    { name: "Ada", scores: [95, 87] },
    { name: "Grace", scores: [91, 88] }
  ]
}
{{</ editor >}}

The runtime serializes the outer object, then each array element, then each nested object and its fields, until every value has been converted.

Array order is preserved. Elements appear in the serialized output in the same order they have in the runtime value.

Object field order is not guaranteed. FQL objects are unordered collections of named fields, and the serialized representation may list fields in any order. Scripts should not rely on a specific field order in the output unless the selected encoding explicitly preserves it.

Object keys are always strings in the serialized representation.

## NONE

`NONE` represents the absence of a value in FQL. When serialized to JSON, `NONE` is encoded as `null`.

{{< editor lang="fql" >}}
RETURN {
  name: "Ada",
  middleName: NONE
}
{{</ editor >}}

Fields with a `NONE` value are included in the serialized output. A field set to `NONE` is not the same as a missing field — it is a field that is explicitly present with an absent value.

## Binary values

Binary values represent raw bytes. Because many output formats are text-based, binary values are encoded in a format-appropriate way.

{{< editor lang="fql" >}}
RETURN IO::NET::HTTP::GET("https://mockery.ferretlang.org/")
{{</ editor >}}

In JSON, binary values are encoded as Base64 strings. Other encodings may use a native binary representation.

## Date and time values

Date and time values are serialized as strings. In JSON, date/time values are encoded using the RFC 3339 format.

{{< editor lang="fql" >}}
RETURN NOW()
{{</ editor >}}

## Host values

Host values represent values provided by the embedding runtime. They may wrap external resources such as database connections, browser pages, files, HTTP clients, streams, or other host-managed objects.

Host values are not serialized by inspecting their internal structure. Instead, their serialized form is determined by the host runtime.

A host value may serialize to an object, string, number, array, `NONE`, or another supported representation. The encoding depends on what the host value exposes. For example, a host value may implement a standard marshaling interface that defines how it is converted to the output format.

Some host values may not be serializable at all. If a non-serializable host value is part of the returned result, serialization produces an error.

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")

// This may fail — a connection handle may not be serializable
RETURN db
{{</ code >}}

Instead, return data extracted from the host value:

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN({ memory: true })

RETURN QUERY `
  SELECT id, name
  FROM users
` IN db USING sql
{{</ code >}}

The first example attempts to serialize a database connection handle, which may not have a meaningful external representation. The second example returns query results — basic values that serialize naturally.

See [Host Values]({{< ref "host" >}}) and [Capability Types]({{< ref "capabilities" >}}) for more about host values and their behavioral contracts.

## Iterable values

Some runtime values are iterable but are not arrays. For example, host values that represent cursors, result sets, or lazy sequences may produce values on demand rather than holding them all in memory.

{{< editor lang="fql" >}}
LET page = WEB::HTML::WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
LET elements = QUERY `.product-card` IN page USING css

RETURN elements
{{</ editor >}}

When an iterable value is part of the returned result, the runtime materializes it into a list before serialization. The iterable is fully consumed, and the resulting list is serialized as an array.

This means iterable values are serializable as long as the values they produce are themselves serializable.

## Cleanup before return

Before the serialized result is returned, the runtime gives resource-backed values a chance to release external resources.

Host values that support the close capability are closed automatically at the end of program execution. This applies to values such as open database connections, file handles, network clients, browser sessions, and other host-owned resources.

Cleanup happens after execution completes and alongside serialization, before the final result is returned to the host. Errors during cleanup are reported alongside the output.

FQL code does not need to close these resources manually. The runtime tracks closable values and handles cleanup automatically.

See [Capability Types]({{< ref "capabilities" >}}) for more about the close capability.

## Encoding formats

Serialization is controlled by the host runtime.

Ferret includes built-in support for JSON and MessagePack encodings. The CLI serializes results as JSON by default. An embedded Go application may convert results into Go values, choose MessagePack for a binary format, or register a custom encoding.

FQL defines the value semantics. The host defines how those values are encoded outside the runtime.

| Encoding | Content type | Notes |
| --- | --- | --- |
| JSON | `application/json` | Default. Text-based, widely supported. |
| MessagePack | `application/vnd.msgpack` | Binary format, compact. |
| Custom | host-defined | Hosts may register additional encodings. |

## Non-serializable values

Some runtime values do not have a portable serialized representation. Returning them as part of the final result produces a serialization error.

Examples include:

- **Functions** — function values are runtime constructs and cannot be represented in an external format.
- **Non-serializable host values** — host values that do not define an external representation, such as raw connection handles or internal runtime objects.

{{< code lang="fql" >}}
LET format = (val) => val + "!"

// This fails — functions are not serializable
RETURN format
{{</ code >}}

If a non-serializable value is nested inside an otherwise serializable structure, the entire serialization fails:

{{< code lang="fql" >}}
LET format = (val) => val + "!"

// This also fails — the function is nested inside the returned object
RETURN {
  name: "Ada",
  formatter: format
}
{{</ code >}}

## Serialization errors

Serialization may fail if the returned value contains something that cannot be represented by the selected encoding.

Common causes include:

- Returning a host value that does not support serialization, such as a database connection or file handle.
- Returning a function value.
- Returning a value that contains a non-serializable value nested within it.

When serialization fails, the error identifies what could not be serialized. The fix is usually to return extracted data instead of the raw resource:

{{< code lang="fql" >}}
// Instead of returning the connection
LET db = DB::SQLITE::OPEN("data.db")

// Return the data
LET users = (
    QUERY `SELECT id, name FROM users` IN db USING sql
)

RETURN users
{{</ code >}}

## Where to go next

{{< docs-related tiles="language-types-host,language-types-capabilities,language-types-ordering,embedding" >}}
