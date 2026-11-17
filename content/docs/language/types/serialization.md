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
return {
  name: "Ada",
  active: true,
  tags: ["math", "logic"]
}
{{</ editor >}}

## When serialization happens

Serialization happens after program execution completes, before the result is returned to the host.

During execution, values remain Ferret runtime values. Arrays, objects, host values, and other runtime constructs are evaluated according to FQL semantics. Only the final result is serialized.

Temporary values that are created and discarded during execution are not serialized. A value only needs to be serializable if it becomes part of the returned result.

{{< editor lang="fql" >}}
let items = [1, 2, 3, 4, 5]

let doubled = (
    for item in items {
        return item * 2
    }
)

return doubled[0]
{{</ editor >}}

Here, `items` and `doubled` are intermediate values. Only the final result — a single number — is serialized.

## Serializable values

Basic FQL values have well-defined serialized representations. The exact encoding depends on the output format selected by the host, but the following table shows the typical mapping for JSON:

| FQL value | JSON representation |
| --- | --- |
| `none` | `null` |
| `true`, `false` | `true`, `false` |
| integer | number |
| float | number |
| duration | normalized duration string |
| string | string |
| array | array |
| object | object |
| binary | Base64-encoded string |
| date/time | RFC 3339 string |

When decoding JSON, integer-form numbers within the signed 64-bit range become `Int` values on both 32-bit and 64-bit hosts, preserving the exact value. Numbers with a fractional or exponent part, and integers outside that range, use the codec's Float conversion. This also applies to numbers nested in arrays and objects.

MessagePack preserves signed 64-bit integer values. Collection lengths must fit both the format's unsigned 32-bit header and the encoder's native integer range; oversized lengths return a range error.

See [Basic Types]({{< ref "basic" >}}) for details on each value type.

## Arrays and objects

Arrays and objects are serialized recursively. Each element or field value must itself be serializable.

{{< editor lang="fql" >}}
return {
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

## none

`none` represents the absence of a value in FQL. When serialized to JSON, `none` is encoded as `null`.

{{< editor lang="fql" >}}
return {
  name: "Ada",
  middleName: none
}
{{</ editor >}}

Fields with a `none` value are included in the serialized output. A field set to `none` is not the same as a missing field — it is a field that is explicitly present with an absent value.

## Binary values

Binary values represent raw bytes. Because many output formats are text-based, binary values are encoded in a format-appropriate way.

{{< editor lang="fql" >}}
return to_binary("Hi")
{{</ editor >}}

The JSON result is `"SGk="`, the Base64 encoding of the two bytes in `Hi`. Other encodings may use a native binary representation.

## Date and time values

Date and time values are serialized as strings. In JSON, date/time values are encoded using the RFC 3339 format.

{{< editor lang="fql" >}}
return datetime::now()
{{</ editor >}}

## Duration values

Duration values are serialized as normalized strings in JSON and MessagePack, including when nested in arrays or objects.

{{< editor lang="fql" >}}
return {
  timeout: 5000ms,
  intervals: [100ms, 1s]
}
{{</ editor >}}

The result uses strings such as `"5s"`, `"100ms"`, and `"1s"`. Ordinary JSON or MessagePack deserialization produces FQL strings; it does not reconstruct duration values automatically.

## Host values

Host values represent values provided by the embedding runtime. They may wrap external resources such as database connections, browser pages, files, HTTP clients, streams, or other host-managed objects.

The selected codec uses the capabilities exposed by the value. For example, a host type can implement a marshaling interface to supply its own JSON representation. A SQLite connection from the optional SQLite module encodes as the display string `"<db.sqlite.connection>"`.

{{< code lang="fql" >}}
let db = db::sqlite::open({ memory: true })
return db
{{</ code >}}

To return records from the database, query the connection. This example returns a row with id `1` and name `Ada`.

{{< code lang="fql" >}}
let db = db::sqlite::open({ memory: true })
return query `SELECT 1 AS id, 'Ada' AS name` in db using sql
{{</ code >}}

See [Host Values]({{< ref "host" >}}) and [Capability Types]({{< ref "capabilities" >}}) for more about host values and their behavioral contracts.

## Iterable values

Some runtime values are iterable but are not arrays. For example, host values that represent cursors, result sets, or lazy sequences may produce values on demand rather than holding them all in memory.

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
let elements = query `.product-card` in page using css

return elements
{{</ editor >}}

For a value encoded through its iterable capability, the JSON encoder walks the items and writes a JSON array. An error while iterating or encoding an item is reported during output consumption.

## Cleanup before return

Native Ferret encodes the result and then closes the VM result, releasing the resources owned by that execution. This includes a closable host parameter when the encoder encounters it in the returned value. The caller observes encoding and cleanup failures when consuming the output. See [Go output]({{< ref "/docs/embedding/go/output" >}}).

## Encoding formats

Serialization is controlled by the host runtime.

Ferret includes built-in support for JSON and MessagePack encodings. The CLI serializes results as JSON by default. An embedded Go application may choose MessagePack or register a custom encoding. The JavaScript package uses JSON. Call `runJSON()` for a decoded value or use the output returned by `run()` for encoded bytes.

FQL defines the value semantics. The host defines how those values are encoded outside the runtime.

| Encoding | Content type | Notes |
| --- | --- | --- |
| JSON | `application/json` | Default. Text-based, widely supported. |
| MessagePack | `application/vnd.msgpack` | Binary format, compact. |
| Custom | host-defined | Hosts may register additional encodings. |

## Non-serializable values

A host type may reject serialization for a particular format. Return the fields the caller needs when the host type has no useful external representation.

## Serialization errors

Encoding can fail when a host value's marshaler returns an error or an iterable fails while producing its items. The same error can occur inside an array or object, because nested values are encoded recursively.

## Next steps

{{< docs-related tiles="language-types-host,language-types-capabilities,language-types-ordering,embedding" >}}
