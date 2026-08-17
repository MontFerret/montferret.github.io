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

FQL has nine built-in value types:

| Type | Example | Description |
| --- | --- | --- |
| `none` | `none` | Represents an absent or undefined value. |
| `bool` | `true`, `false` | Represents a truth value. |
| `number` | `42`, `3.14` | Represents numeric values, both integer and floating-point. |
| `duration` | `250ms`, `1.5s` | Represents a signed length of time with nanosecond precision. |
| `string` | `"hello"` | Represents text. |
| `datetime` | `now()` | Represents a point in time. |
| `array` | `[1, 2, 3]` | Represents an ordered sequence of values. |
| `object` | `{ name: "Ada" }` | Represents a set of named fields. |
| `binary` | module-specific | Represents raw bytes. |

Ferret hosts may also define additional value types — documents, elements, HTTP responses, files, handles, database connections, and other external resources. These are called host values and are covered separately in [Host Values](#).

## none

`none` is Ferret's equivalent of `null`, `nil`, or `None` in other languages.

It represents the absence of a value.

{{< editor lang="fql"  >}}
return none
{{</ editor >}}

`none` is a real value in FQL. It can be assigned to variables, placed inside arrays and objects, returned from functions, and compared with other values.

{{< editor lang="fql" >}}
let value = none

return {
    value: value,
    isMissing: value == none
}
{{</ editor >}}

Unlike SQL `null`, comparing with `none` does not automatically produce another `none`.

{{< editor lang="fql">}}
return {
    same: none == none,
    different: none != 1
}
{{</ editor >}}

`none` is useful when a value is intentionally absent:

{{< editor lang="fql">}}
let user = {
    name: "Ada",
    email: none
}

return user.email
{{</ editor >}}

It is useful when that absence needs to remain distinct from ordinary values such as `false`, `0`, `""`, `[]`, or `{}`.

## Booleans

Boolean values represent truth values.

FQL has two boolean values:

- `true`
- `false`

Booleans are commonly produced by comparison operators, logical operators, predicates, and conditions.

{{< editor lang="fql" >}}
let count = 3

return {
    hasItems: count > 0,
    isEmpty: count == 0
}
{{</ editor >}}

Boolean values can be combined with logical operators:

{{< editor lang="fql" >}}
let enabled = true
let hasAccess = false

return enabled and hasAccess
{{</ editor >}}

## Numbers

FQL uses the number type for numeric values. Numbers can be written as integers or decimals:

{{< editor lang="fql" >}}
return {
    count: 10,
    price: 19.95,
    temperature: -4.5
}
{{</ editor >}}

Numbers support arithmetic operations:

{{< editor lang="fql" >}}
let price = 20
let tax = 1.5

return price + tax
{{</ editor >}}

They can also be compared:

{{< editor lang="fql" >}}
return {
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
return {
    number: 10,
    string: "10",
    same: 10 == "10"
}
{{</ editor >}}

Use explicit conversion functions when a script needs to turn text into a number or a number into text. For example, `to_number("10") - 2` returns `8`, while `"10" - 2` raises an invalid-operation error. A malformed explicit conversion raises a conversion error.

FQL internally distinguishes integers and floats. Most arithmetic and comparison operations work the same for both, but the distinction matters when type-checking functions are used. `is_int` returns true only for integer values, and `is_float` returns true only for floating-point values. Use `to_int` or `to_float` to convert between the two when needed.

## Strings

Strings represent text and can be written as quoted literals:

{{< editor lang="fql"  >}}
return "Hello, Ferret"
{{</ editor >}}

Strings can be empty:

{{< editor lang="fql"  >}}
return ""
{{</ editor >}}

An empty string is still a text value, not the absence of a value. It is different from `none`.

{{< editor lang="fql" >}}
return {
    empty: "",
    missing: none,
    same: "" == none
}
{{</ editor >}}

Strings are commonly used for names, labels, URLs, selectors, extracted text, keys, and identifiers.

{{< editor lang="fql" >}}
return {
    title: "Example page",
    url: "https://example.com",
    selector: ".article-title"
}
{{</ editor >}}

Strings can be compared with other strings. Comparisons use the string contents, so two strings are equal when they contain the same sequence of characters.

{{< editor lang="fql" >}}
return {
    same: "ferret" == "ferret",
    different: "ferret" != "Ferret"
}
{{</ editor >}}

String comparisons are always case-sensitive.

{{< editor lang="fql" >}}
return "Ferret" == "ferret"
{{</ editor >}}

When either `+` operand is a String, Ferret concatenates the String representations of both operands. This is the only arithmetic rule that stringifies otherwise unsupported operand types.

{{< editor lang="fql" >}}
return ["items: " + 3, 1s + " elapsed"]
{{</ editor >}}

### Template literals

Template literals are strings delimited by backticks that support embedded expressions. An expression inside `${...}` is evaluated and its result is inserted into the string.

{{< editor lang="fql" >}}
let name = "Ferret"
let version = 2

return `Hello from ${name} v${version}!`
{{</ editor >}}

Any FQL expression can appear inside the interpolation:

{{< editor lang="fql" >}}
let price = 19.95
let quantity = 3

return `Total: ${price * quantity}`
{{</ editor >}}

Template literals can span multiple lines. Backtick strings are useful when constructing text that includes variable data without calling `concat`.

## Arrays

Arrays are ordered collections of values written with square brackets:

{{< editor lang="fql" >}}
return [1, 2, 3]
{{</ editor >}}

Arrays can contain any FQL value:

{{< editor lang="fql" >}}
return [
    none,
    true,
    42,
    "text",
    [1, 2],
    { name: "Ada" }
]
{{</ editor >}}

Arrays preserve the order of their elements.

{{< editor lang="fql" >}}
return ["first", "second", "third"]
{{</ editor >}}

Array items can be accessed by index using bracket notation:

{{< editor lang="fql" >}}
let roles = ["admin", "editor", "viewer"]

return roles[0]
{{</ editor >}}

An empty array represents a collection with no items:

{{< editor lang="fql" >}}
return []
{{</ editor >}}

An empty array is not the same as `none`.

{{< editor lang="fql" >}}
return {
    emptyArray: [],
    missing: none,
    same: [] == none
}
{{</ editor >}}

Arrays are commonly produced by for loops and collection operations.

{{< editor lang="fql" >}}
let numbers = (
    for value in [1, 2, 3] {
        return value * 2
    }
)

return numbers
{{</ editor >}}

Array elements can be nested:

{{< editor lang="fql" >}}
return [
    ["Ada", "Lovelace"],
    ["Grace", "Hopper"]
]
{{</ editor >}}

Use `...expression` inside an array literal to copy the elements of another Array or host-provided runtime List into a new array:

{{< editor lang="fql" >}}
let middle = [2, 3]

return [1, ...middle, 4]
{{</ editor >}}

Array entries are evaluated once from left to right. Spreading `none` adds no elements. The copy is shallow: the source list is not changed, but nested values are shared. The spread operand must be an Array, another runtime List, or `none`; strings, ranges, and values that are merely iterable produce a `TypeError`.

## Objects

An object is an unordered collection of named values, written with curly braces:

{{< editor lang="fql" >}}
return {
    name: "Ada",
    active: true
}
{{</ editor >}}

Each object entry has a property name and a value, which can be any FQL value.

{{< editor lang="fql" >}}
return {
    name: "Ada",
    age: 36,
    email: none,
    tags: ["math", "computing"],
    profile: {
        active: true
    }
}
{{</ editor >}}

Object properties can be accessed with dot notation or bracket notation.

Use dot notation when the property name is known ahead of time:

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
    active: true
}

return user.name
{{</ editor >}}

Use bracket notation when the property name comes from a string literal or variable:

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
    active: true
}

let property = "name"

return user[property]
{{</ editor >}}

Because objects do not preserve key order, property order does not affect object equality. Two objects with the same properties and values are equal, even if those properties were written in a different order.

{{< editor lang="fql" >}}
return { a: 1, b: 2 } == { b: 2, a: 1 }
{{</ editor >}}

An empty object represents a record with no properties:

{{< editor lang="fql" >}}
return {}
{{</ editor >}}

An empty object is not the same as `none`.

{{< editor lang="fql" >}}
return {
    emptyObject: {},
    missing: none,
    same: {} == none
}
{{</ editor >}}

Use `...expression` inside an object literal to copy properties from another Object or host-provided runtime Map into a new object:

{{< editor lang="fql" >}}
let defaults = {
    theme: "light",
    pageSize: 20
}

return {
    ...defaults,
    theme: "dark"
}
{{</ editor >}}

Object entries are evaluated once from left to right, so a later property replaces an earlier property with the same name, whether either property came from a spread. Spreading `none` adds no properties. The copy is shallow and does not mutate the source map. The operand must be an Object, another runtime Map, or `none`; other concrete values produce a `TypeError`.

Spread is available only while constructing array and object literals. It does not add rest destructuring, spread call arguments, `for` expansion, deep merging, or mutation.

## DateTime values

DateTime values represent a specific point in time.

They are typically created using standard library functions such as `now()`, `date()`, or `to_datetime()`:

{{< editor lang="fql" >}}
let now = now()

return {
    current: now,
    year: date_year(now),
    month: date_month(now)
}
{{</ editor >}}

DateTime values support native instant comparison and checked arithmetic with native Duration values. Adding a Duration in either operand order produces another DateTime. Subtracting a Duration from a DateTime produces another DateTime, and subtracting two DateTime values produces the elapsed Duration between their canonical instants.

{{< editor lang="fql" >}}
let start = to_datetime("2024-01-01T00:00:00Z")
let end = start + 30d

return {
    start: start,
    end: end,
    elapsed: end - start
}
{{</ editor >}}

Use `is_datetime` to check whether a value is a DateTime. `to_datetime` accepts an existing DateTime or an RFC3339 string. It also accepts an Int or finite Float Unix epoch offset when the unit is explicit: `to_datetime(value, "s")`, `to_datetime(value, "ms")`, `to_datetime(value, "us")`, or `to_datetime(value, "ns")`. Ferret does not infer an epoch unit from numeric magnitude, and numeric strings are not treated as epoch values.

See [the DateTime Standard Library category]({{% ref "docs/standard-library/datetime" %}}) for the full list of available operations.

## Binary values

Binary values represent raw bytes.

They are used for data that should be handled as bytes instead of ordinary text, such as downloaded files, encoded payloads, images, archives, or data exchanged with runtime modules.

Binary values are part of the FQL value model, but they are usually returned by functions, modules, or runtime operations. For example, an HTTP, file, or encoding module may return binary data when the result should be treated as bytes.

{{< editor lang="fql" >}}
let file = io::net::http::get("https://avatars.githubusercontent.com/u/39228646?s=200&v=4")

return file
{{</ editor >}}

<div class="notification is-info">
  Ferret’s default serializer encodes binary values as Base64 strings, so byte-oriented data can be represented safely in text-based output formats.
</div>  

## Durations

Durations are native FQL values backed by signed nanoseconds. A duration literal is a number immediately followed by a unit suffix. Suffixes are case-insensitive:

| Suffix | Unit |
| --- | --- |
| `ms` | milliseconds |
| `s` | seconds |
| `m` | minutes |
| `h` | hours |
| `d` | days (exactly 24 hours) |

{{< code lang="fql" >}}
100ms       // 100 milliseconds
5s          // 5 seconds
1.5m        // 90 seconds
2H          // suffixes are case-insensitive
1d          // exactly 24 hours
2.5e-1s     // 250 milliseconds
{{</ code >}}

Duration literals work anywhere an ordinary expression is accepted. Compound source literals such as `1h30m` are not supported; compose literals with arithmetic instead. Duration strings do support compound forms, so `to_duration("1h30m")` produces the same value as `1h + 30m`.

{{< code lang="fql" >}}
let interval = 1h + 30m

return {
    doubled: interval * 2,
    ratio: 1s / 250ms,
    isDuration: is_duration(interval),
    type: typename(interval)
}
{{</ code >}}

`to_duration` and scheduling expressions use the broad Duration conversion rules:

| Source value | Duration result |
| --- | --- |
| Duration | unchanged |
| Int or Float | milliseconds; fractional nanoseconds are truncated toward zero |
| Duration string | parsed using duration syntax, including compound forms such as `"1h30m"` |
| `none` or `false` | `0ms` |
| `true` | `1ms` |
| empty list | `0ms` |
| singleton list | recursively converts its only element |
| list with multiple elements | runtime error |
| object or unsupported value | runtime error |

Arithmetic and comparison operators do not apply this conversion implicitly. Duration addition and subtraction require two native Durations. Multiplication accepts a native `Int` or `Float` in either operand order. Division accepts a native number for scaling or another Duration for a ratio. An exact Duration ratio produces an integer; a fractional ratio produces a float.

Duration equality is also strict: `1s == 1000ms` is true, while `1s == 1000` and `1s == "1s"` are false. Relational comparison between Duration and a non-Duration value raises an invalid-operation error. Use `to_duration(value)` explicitly when conversion is intended.

String-triggered `+` remains concatenation. For example, `1s + "1s"` returns `"1s1s"`, while `1s + to_duration("1s")` returns `2s`.

Conversion, parsing, scaling, and division truncate fractional nanoseconds toward zero. Values outside the signed Duration range raise a range error instead of wrapping. Use `to_string` when text is required.

Equivalent values have the same normalized string form. For example, `5000ms` renders as `5s`, and days may render as hours. A zero Duration is false when evaluated by binary `and` or `or`; any non-zero Duration is true. Unary `!` and `not` accept only Boolean values. Negative durations are valid values and arithmetic results, but scheduling operations reject them.

## Type checks

Scripts often receive values in different shapes. Use type-related functions or predicates when logic needs to branch based on the kind of value being handled.

{{< editor lang="fql" >}}
let value = "42"

return is_string(value) ? "text" : "not text"
{{</ editor >}}

Type checks are useful when working with external data, optional fields, runtime-backed values, or module results.

{{< editor lang="fql" >}}
let user = {
    name: "Ada",
}
let value = user.email

return value == none ? "missing" : value
{{</ editor >}}

See [the Types Standard Library category]({{% ref "docs/standard-library/types" %}}) for the full list of available type-checking functions.

## Next steps

{{< docs-related tiles="language-types-host,language-types-capabilities,language-types-ordering,language-types-serialization" >}}
