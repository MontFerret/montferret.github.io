---
title: "Types"
weight: 110
draft: false
description: "Type-checking and conversion functions in the Ferret standard library."
aliases:
  - /docs/stdlib/types/
menuTitle: 
menu: [is_array,is_binary,is_bool,is_datetime,is_duration,is_float,IS_HTML_DOCUMENT,IS_HTML_ELEMENT,is_int,is_list,is_map,is_nan,is_none,is_object,is_string,to_array,to_binary,to_bool,to_datetime,to_duration,to_float,to_int,to_number,to_object,to_string,typename,]
---



{{< header href="is_array" >}}

is_array

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_array.go#L13)

is_array checks whether value is an array value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is array, otherwise false.
- - - -


{{< header href="is_binary" >}}

is_binary

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_binary.go#L13)

is_binary checks whether value is a binary value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is binary, otherwise false.
- - - -


{{< header href="is_bool" >}}

is_bool

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_boolean.go#L13)

is_bool checks whether value is a boolean value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is boolean, otherwise false.
- - - -


{{< header href="is_datetime" >}}

is_datetime

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_date_time.go#L13)

is_datetime checks whether value is a date time value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is date time, otherwise false.
- - - -


{{< header href="is_duration" >}}

is_duration

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_duration.go#L9)

is_duration checks whether value is a native duration value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is a duration, otherwise false.
- - - -


{{< header href="is_float" >}}

is_float

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_float.go#L13)

is_float checks whether value is a float value.

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

is_int

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_int.go#L13)

is_int checks whether value is a int value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is int, otherwise false.
- - - -


{{< header href="is_list" >}}

is_list

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_list.go#L12)

is_list checks whether value is a list value. This is an alias for is_array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is list, otherwise false.
- - - -


{{< header href="is_map" >}}

is_map

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_map.go#L12)

is_map checks whether value is a map value. This is an alias for is_object.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is map, otherwise false.
- - - -


{{< header href="is_nan" >}}

is_nan

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_nan.go#L13)

is_nan checks whether value is NaN.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is nan, otherwise false.
- - - -


{{< header href="is_none" >}}

is_none

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_none.go#L13)

is_none checks whether value is a none value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is none, otherwise false.
- - - -


{{< header href="is_object" >}}

is_object

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_object.go#L13)

is_object checks whether value is an object value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is object, otherwise false.
- - - -


{{< header href="is_string" >}}

is_string

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_string.go#L13)

is_string checks whether value is a string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is string, otherwise false.
- - - -


{{< header href="to_array" >}}

to_array

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_array.go#L15)

to_array takes an input value of any type and convert it into an array value. None is converted to an empty array Boolean values, numbers and strings are converted to an array containing the original value as its single element Arrays keep their original value Objects / HTML nodes are converted to an array containing their attribute values as array elements.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`input` | `Any`  |  | Input value of arbitrary type.


**Returns** `Any[]` An array value.
- - - -


{{< header href="to_binary" >}}

to_binary

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_binary.go#L12)

to_binary takes an input value of any type and converts it into a binary value. The value is first converted to its string representation, then the string bytes are used as the binary content.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Binary` A binary value.
- - - -


{{< header href="to_bool" >}}

to_bool

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_boolean.go#L18)

to_bool takes an input value of any type and converts it into the appropriate boolean value. None is converted to false Numbers are converted to true, except for 0, which is converted to false Strings are converted to true if they are non-empty, and to false otherwise Dates are converted to true if they are not zero, and to false otherwise Arrays are always converted to true (even if empty) Objects / HtmlNodes / Binary are always converted to true

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` The appropriate boolean value.
- - - -


{{< header href="to_datetime" >}}

to_datetime

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_date_time.go#L20)

to_datetime returns an existing DateTime or parses an RFC3339 string. With an explicit unit, it also converts an Int or finite Float offset from the Unix epoch. Supported units are `s`, `ms`, `us`, and `ns`; aliases include `sec`, `second`, `seconds`, `millisecond`, `milliseconds`, `µs`, `μs`, `microsecond`, `microseconds`, `nanosecond`, and `nanoseconds`. Unit matching is case-insensitive.

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
    parsed: to_datetime("2026-08-02T12:00:00Z"),
    seconds: to_datetime(1690992000, "s"),
    milliseconds: to_datetime(1690992000000, "ms"),
    fractional: to_datetime(1.5, "s"),
    beforeEpoch: to_datetime(-1, "s")
}
{{</ editor >}}


{{< header href="to_duration" >}}

to_duration

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_duration.go#L13)

to_duration converts a value to a native Duration. Numeric values are milliseconds, duration strings may use compound syntax, and fractional nanoseconds are truncated toward zero.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Value supported by the canonical Duration conversion rules.


**Returns** `Duration` Converted duration value.
- - - -


{{< header href="to_float" >}}

to_float

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_float.go#L20)

to_float converts a supported value to a Float. Int values are converted to Float, numeric strings are parsed, Booleans become `0` or `1`, DateTime values become Unix seconds, and list values are converted recursively and summed. Malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Float` A float value.
- - - -


{{< header href="to_int" >}}

to_int

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_int.go#L20)

to_int converts a supported value to an Int. Float values are truncated toward zero, integer strings are parsed, Booleans become `0` or `1`, DateTime values become Unix seconds, and list values are converted recursively and summed. Malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Int` An integer value.
- - - -


{{< header href="to_number" >}}

to_number

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_number.go#L20)

to_number converts a supported value to a native number. Existing Int and Float values keep their type; other supported values use Float conversion. Numeric strings are parsed explicitly, while malformed strings and unsupported types raise conversion errors.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Value to convert to a number.


**Returns** `Int, Float` Converted numeric value.

{{< editor lang="fql" >}}
return {
    integer: to_number(10),
    parsed: to_number("10"),
    fractional: to_number("2.5"),
    arithmetic: to_number("10") - 2
}
{{</ editor >}}

- - - -


{{< header href="to_object" >}}

to_object

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_object.go#L12)

to_object converts the given value to an object. The conversion rules depend on the input type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Object` The object representation of the given value.
- - - -


{{< header href="to_string" >}}

to_string

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_string.go#L12)

to_string takes an input value of any type and convert it into a string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `String` String representation of a given value.
- - - -


{{< header href="typename" >}}

typename

{{</ header >}}

typename returns the data type name of a value as a string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `String` Returns string representation of a type (e.g. `"Int"`, `"Float"`, `"Duration"`, `"String"`, `"Array"`, `"Object"`, `"DateTime"`, `"Boolean"`, `"Binary"`, `"None"`).
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-types-basic,language-types-host" >}}
