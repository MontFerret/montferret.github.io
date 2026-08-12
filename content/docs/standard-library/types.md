---
title: "Types"
weight: 110
draft: false
description: "Type-checking and conversion functions in the Ferret standard library."
aliases:
  - /docs/stdlib/types/
menuTitle: 
menu: [IS_ARRAY,IS_BINARY,IS_BOOL,IS_DATETIME,IS_DURATION,IS_FLOAT,IS_HTML_DOCUMENT,IS_HTML_ELEMENT,IS_INT,IS_LIST,IS_MAP,IS_NAN,IS_NONE,IS_OBJECT,IS_STRING,TO_ARRAY,TO_BINARY,TO_BOOL,TO_DATETIME,TO_DURATION,TO_FLOAT,TO_INT,TO_NUMBER,TO_OBJECT,TO_STRING,TYPENAME,]
---



{{< header href="is_array" >}}

IS_ARRAY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_array.go#L13)

IS_ARRAY checks whether value is an array value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is array, otherwise false.
- - - -


{{< header href="is_binary" >}}

IS_BINARY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_binary.go#L13)

IS_BINARY checks whether value is a binary value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is binary, otherwise false.
- - - -


{{< header href="is_bool" >}}

IS_BOOL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_boolean.go#L13)

IS_BOOL checks whether value is a boolean value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is boolean, otherwise false.
- - - -


{{< header href="is_datetime" >}}

IS_DATETIME

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_date_time.go#L13)

IS_DATETIME checks whether value is a date time value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is date time, otherwise false.
- - - -


{{< header href="is_duration" >}}

IS_DURATION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_duration.go#L9)

IS_DURATION checks whether value is a native duration value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is a duration, otherwise false.
- - - -


{{< header href="is_float" >}}

IS_FLOAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_float.go#L13)

IS_FLOAT checks whether value is a float value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is float, otherwise false.
- - - -


{{< header href="is_html_document" >}}

IS_HTML_DOCUMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_html_document.go#L13)

IS_HTML_DOCUMENT checks whether value is a HTMLDocument value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is htmldocument, otherwise false.
- - - -


{{< header href="is_html_element" >}}

IS_HTML_ELEMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_html_element.go#L13)

IS_HTML_ELEMENT checks whether value is a HTMLElement value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is htmlelement, otherwise false.
- - - -


{{< header href="is_int" >}}

IS_INT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_int.go#L13)

IS_INT checks whether value is a int value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is int, otherwise false.
- - - -


{{< header href="is_list" >}}

IS_LIST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_list.go#L12)

IS_LIST checks whether value is a list value. This is an alias for IS_ARRAY.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is list, otherwise false.
- - - -


{{< header href="is_map" >}}

IS_MAP

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_map.go#L12)

IS_MAP checks whether value is a map value. This is an alias for IS_OBJECT.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is map, otherwise false.
- - - -


{{< header href="is_nan" >}}

IS_NAN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_nan.go#L13)

IS_NAN checks whether value is NaN.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is nan, otherwise false.
- - - -


{{< header href="is_none" >}}

IS_NONE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_none.go#L13)

IS_NONE checks whether value is a none value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is none, otherwise false.
- - - -


{{< header href="is_object" >}}

IS_OBJECT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_object.go#L13)

IS_OBJECT checks whether value is an object value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is object, otherwise false.
- - - -


{{< header href="is_string" >}}

IS_STRING

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_string.go#L13)

IS_STRING checks whether value is a string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is string, otherwise false.
- - - -


{{< header href="to_array" >}}

TO_ARRAY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_array.go#L15)

TO_ARRAY takes an input value of any type and convert it into an array value. None is converted to an empty array Boolean values, numbers and strings are converted to an array containing the original value as its single element Arrays keep their original value Objects / HTML nodes are converted to an array containing their attribute values as array elements.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`input` | `Any`  |  | Input value of arbitrary type.


**Returns** `Any[]` An array value.
- - - -


{{< header href="to_binary" >}}

TO_BINARY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_binary.go#L12)

TO_BINARY takes an input value of any type and converts it into a binary value. The value is first converted to its string representation, then the string bytes are used as the binary content.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Binary` A binary value.
- - - -


{{< header href="to_bool" >}}

TO_BOOL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_boolean.go#L18)

TO_BOOL takes an input value of any type and converts it into the appropriate boolean value. None is converted to false Numbers are converted to true, except for 0, which is converted to false Strings are converted to true if they are non-empty, and to false otherwise Dates are converted to true if they are not zero, and to false otherwise Arrays are always converted to true (even if empty) Objects / HtmlNodes / Binary are always converted to true

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` The appropriate boolean value.
- - - -


{{< header href="to_datetime" >}}

TO_DATETIME

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_date_time.go#L20)

TO_DATETIME returns an existing DateTime or parses an RFC3339 string. With an explicit unit, it also converts an Int or finite Float offset from the Unix epoch. Supported units are `s`, `ms`, `us`, and `ns`; aliases include `sec`, `second`, `seconds`, `millisecond`, `milliseconds`, `µs`, `μs`, `microsecond`, `microseconds`, `nanosecond`, and `nanoseconds`. Unit matching is case-insensitive.

Numeric values always require a unit; Ferret never guesses from their magnitude. Negative epochs and fractional values are supported, and fractions smaller than one nanosecond are truncated toward zero. Numeric strings are not epoch values. Unknown units, non-finite numbers, and values outside the DateTime range raise runtime errors. A unit supplied with a DateTime or RFC3339 string is also an error.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `DateTime, String, Int, Float`  |  | Existing DateTime, RFC3339 string, or numeric Unix epoch offset.
`unit` | `String`  |  | Optional epoch unit for numeric values: `s`, `ms`, `us`, or `ns`.


**Returns** `DateTime` Parsed date time.
- - - -

{{< editor lang="fql" >}}
return {
    parsed: TO_DATETIME("2026-08-02T12:00:00Z"),
    seconds: TO_DATETIME(1690992000, "s"),
    milliseconds: TO_DATETIME(1690992000000, "ms"),
    fractional: TO_DATETIME(1.5, "s"),
    beforeEpoch: TO_DATETIME(-1, "s")
}
{{</ editor >}}


{{< header href="to_duration" >}}

TO_DURATION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_duration.go#L13)

TO_DURATION converts a value to a native Duration. Numeric values are milliseconds, duration strings may use compound syntax, and fractional nanoseconds are truncated toward zero.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Value supported by the canonical Duration conversion rules.


**Returns** `Duration` Converted duration value.
- - - -


{{< header href="to_float" >}}

TO_FLOAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_float.go#L20)

TO_FLOAT converts a supported value to a Float. Int values are converted to Float, numeric strings are parsed, Booleans become `0` or `1`, DateTime values become Unix seconds, and list values are converted recursively and summed. Malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Float` A float value.
- - - -


{{< header href="to_int" >}}

TO_INT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_int.go#L20)

TO_INT converts a supported value to an Int. Float values are truncated toward zero, integer strings are parsed, Booleans become `0` or `1`, DateTime values become Unix seconds, and list values are converted recursively and summed. Malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Int` An integer value.
- - - -


{{< header href="to_number" >}}

TO_NUMBER

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_number.go#L20)

TO_NUMBER converts a supported value to a native number. Existing Int and Float values keep their type; other supported values use Float conversion. Numeric strings are parsed explicitly, while malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Value to convert to a number.


**Returns** `Int, Float` Converted numeric value.

{{< editor lang="fql" >}}
return {
    integer: TO_NUMBER(10),
    parsed: TO_NUMBER("10"),
    fractional: TO_NUMBER("2.5"),
    arithmetic: TO_NUMBER("10") - 2
}
{{</ editor >}}

- - - -


{{< header href="to_object" >}}

TO_OBJECT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_object.go#L12)

TO_OBJECT converts the given value to an object. The conversion rules depend on the input type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Object` The object representation of the given value.
- - - -


{{< header href="to_string" >}}

TO_STRING

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_string.go#L12)

TO_STRING takes an input value of any type and convert it into a string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `String` String representation of a given value.
- - - -


{{< header href="typename" >}}

TYPENAME

{{</ header >}}

TYPENAME returns the data type name of a value as a string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `String` Returns string representation of a type (e.g. `"Int"`, `"Float"`, `"Duration"`, `"String"`, `"Array"`, `"Object"`, `"DateTime"`, `"Boolean"`, `"Binary"`, `"None"`).
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-types-basic,language-types-host" >}}
