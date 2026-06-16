---
title: "Basic Types"
sidebarTitle: "Basic Types"
weight: 10
draft: false
description: "The basic value types available in FQL."
aliases:
    - /docs/fql/basic-types/
---

# Basic types

FQL has seven built-in value types:

| Type | Example | Description |
| --- | --- | --- |
| `NONE` | `NONE` | Represents an absent or undefined value. |
| `bool` | `true`, `false` | Represents a truth value. |
| `number` | `42`, `3.14` | Represents numeric values, both integer and floating-point. |
| `string` | `"hello"` | Represents text. |
| `array` | `[1, 2, 3]` | Represents an ordered sequence of values. |
| `object` | `{ name: "Ada" }` | Represents a set of named fields. |
| `binary` | module-specific | Represents raw bytes. |

Ferret hosts may also define additional value types — documents, elements, HTTP responses, files, handles, database connections, and other external resources. These are called host values and are covered separately in [Host Values](#).

## NONE

`NONE` is Ferret's equivalent of `null`, `nil`, or `None` in other languages.

It represents the absence of a value.

{{< editor lang="fql"  >}}
RETURN NONE
{{</ editor >}}

`NONE` is a real value in FQL. It can be assigned to variables, placed inside arrays and objects, returned from functions, and compared with other values.

{{< editor lang="fql" >}}
LET value = NONE

RETURN {
    value: value,
    isMissing: value == NONE
}
{{</ editor >}}

Unlike SQL `NULL`, comparing with `NONE` does not automatically produce another `NONE`.

{{< editor lang="fql">}}
RETURN {
    same: NONE == NONE,
    different: NONE != 1
}
{{</ editor >}}

`NONE` is useful when a value is intentionally absent:

{{< editor lang="fql">}}
LET user = {
    name: "Ada",
    email: NONE
}

RETURN user.email
{{</ editor >}}

It is useful when that absence needs to remain distinct from ordinary values such as `false`, `0`, `""`, `[]`, or `{}`.

## Booleans

Boolean values represent truth values.

FQL has two boolean values:

- `true`
- `false`

Booleans are commonly produced by comparison operators, logical operators, predicates, and conditions.

{{< editor lang="fql" >}}
LET count = 3

RETURN {
    hasItems: count > 0,
    isEmpty: count == 0
}
{{</ editor >}}

Boolean values can be combined with logical operators:

{{< editor lang="fql" >}}
LET enabled = true
LET hasAccess = false

RETURN enabled AND hasAccess
{{</ editor >}}

## Numbers

FQL uses the number type for numeric values. Numbers can be written as integers or decimals:

{{< editor lang="fql" >}}
RETURN {
    count: 10,
    price: 19.95,
    temperature: -4.5
}
{{</ editor >}}

Numbers support arithmetic operations:

{{< editor lang="fql" >}}
LET price = 20
LET tax = 1.5

RETURN price + tax
{{</ editor >}}

They can also be compared:

{{< editor lang="fql" >}}
RETURN {
    less: 1 < 2,
    equal: 10 == 10,
    greater: 5 > 3
}
{{</ editor >}}

Numbers are ordered by numeric value:
- `1 < 2`
- `10 > 2`
- `-1 < 0`

Numeric-looking strings are still strings. FQL does not treat `"10"` as the same value as `10`.

{{< editor lang="fql" >}}
RETURN {
    number: 10,
    string: "10",
    same: 10 == "10"
}
{{</ editor >}}

Use explicit conversion functions when a script needs to turn text into a number or a number into text.

## Strings

Strings represent text and can be written as quoted literals:

{{< editor lang="fql"  >}}
RETURN "Hello, Ferret"
{{</ editor >}}

Strings can be empty:

{{< editor lang="fql"  >}}
RETURN ""
{{</ editor >}}

An empty string is still a text value, not the absence of a value. It is different from `NONE`.

{{< editor lang="fql" >}}
RETURN {
    empty: "",
    missing: NONE,
    same: "" == NONE
}
{{</ editor >}}

Strings are commonly used for names, labels, URLs, selectors, extracted text, keys, and identifiers.

{{< editor lang="fql" >}}
RETURN {
    title: "Example page",
    url: "https://example.com",
    selector: ".article-title"
}
{{</ editor >}}

Strings can be compared with other strings. Comparisons use the string contents, so two strings are equal when they contain the same sequence of characters.

{{< editor lang="fql" >}}
RETURN {
    same: "ferret" == "ferret",
    different: "ferret" != "Ferret"
}
{{</ editor >}}

String comparisons are always case-sensitive.

{{< editor lang="fql" >}}
RETURN "Ferret" == "ferret"
{{</ editor >}}

## Arrays

Arrays are ordered collections of values written with square brackets:

{{< editor lang="fql" >}}
RETURN [1, 2, 3]
{{</ editor >}}

Arrays can contain any FQL value:

{{< editor lang="fql" >}}
RETURN [
    NONE,
    true,
    42,
    "text",
    [1, 2],
    { name: "Ada" }
]
{{</ editor >}}

Arrays preserve the order of their elements.

{{< editor lang="fql" >}}
RETURN ["first", "second", "third"]
{{</ editor >}}

Array items can be accessed by index using bracket notation:

{{< editor lang="fql" >}}
LET roles = ["admin", "editor", "viewer"]

RETURN roles[0]
{{</ editor >}}

An empty array represents a collection with no items:

{{< editor lang="fql" >}}
RETURN []
{{</ editor >}}

An empty array is not the same as `NONE`.

{{< editor lang="fql" >}}
RETURN {
    emptyArray: [],
    missing: NONE,
    same: [] == NONE
}
{{</ editor >}}

Arrays are commonly produced by FOR loops and collection operations.

{{< editor lang="fql" >}}
LET numbers = (
    FOR value IN [1, 2, 3]
        RETURN value * 2
)

RETURN numbers
{{</ editor >}}

Array elements can be nested:

{{< editor lang="fql" >}}
RETURN [
    ["Ada", "Lovelace"],
    ["Grace", "Hopper"]
]
{{</ editor >}}

## Objects

An object is an unordered collection of named values, written with curly braces:

{{< editor lang="fql" >}}
RETURN {
    name: "Ada",
    active: true
}
{{</ editor >}}

Each object entry has a property name and a value, which can be any FQL value.

{{< editor lang="fql" >}}
RETURN {
    name: "Ada",
    age: 36,
    email: NONE,
    tags: ["math", "computing"],
    profile: {
        active: true
    }
}
{{</ editor >}}

Object properties can be accessed with dot notation or bracket notation.

Use dot notation when the property name is known ahead of time:

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
    active: true
}

RETURN user.name
{{</ editor >}}

Use bracket notation when the property name comes from a string literal or variable:

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
    active: true
}

LET property = "name"

RETURN user[property]
{{</ editor >}}

Because objects do not preserve key order, property order does not affect object equality. Two objects with the same properties and values are equal, even if those properties were written in a different order.

{{< editor lang="fql" >}}
RETURN { a: 1, b: 2 } == { b: 2, a: 1 }
{{</ editor >}}

An empty object represents a record with no properties:

{{< editor lang="fql" >}}
RETURN {}
{{</ editor >}}

An empty object is not the same as `NONE`.

{{< editor lang="fql" >}}
RETURN {
    emptyObject: {},
    missing: NONE,
    same: {} == NONE
}
{{</ editor >}}

## Binary values

Binary values represent raw bytes.

They are used for data that should be handled as bytes instead of ordinary text, such as downloaded files, encoded payloads, images, archives, or data exchanged with runtime modules.

Binary values are part of the FQL value model, but they are usually returned by functions, modules, or runtime operations. For example, an HTTP, file, or encoding module may return binary data when the result should be treated as bytes.

{{< editor lang="fql" >}}
LET file = IO::NET::HTTP::GET("https://avatars.githubusercontent.com/u/39228646?s=200&v=4")

RETURN file
{{</ editor >}}

<div class="notification is-info">
  Ferret’s default serializer encodes binary values as Base64 strings, so byte-oriented data can be represented safely in text-based output formats.
</div>  

## Type checks

Scripts often receive values in different shapes. Use type-related functions or predicates when logic needs to branch based on the kind of value being handled.

{{< editor lang="fql" >}}
LET value = "42"

RETURN IS_STRING(value) ? "text" : "not text"
{{</ editor >}}

Type checks are useful when working with external data, optional fields, runtime-backed values, or module results.

{{< editor lang="fql" >}}
LET user = {
    name: "Ada",
}
LET value = user.email

RETURN value == NONE ? "missing" : value
{{</ editor >}}

See [the Standard Library section]({{% ref "docs/standard-library/types" %}}) for the full list of available type-checking functions.

## Where to go next

{{< docs-related tiles="language-types-host,language-types-capabilities,language-types-ordering,language-types-serialization" >}}
