---
title: "types"
weight: 1
draft: false
menuTitle: 
menu: [IS_ARRAY,IS_BINARY,IS_BOOL,IS_DATETIME,IS_FLOAT,IS_HTML_DOCUMENT,IS_HTML_ELEMENT,IS_INT,IS_NAN,IS_NONE,IS_OBJECT,IS_STRING,TO_ARRAY,TO_BOOL,TO_DATETIME,TO_FLOAT,TO_INT,TO_STRING,TYPENAME,]
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
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_date_time.go#L12)

TO_DATETIME takes an input value of any type and converts it into the appropriate date time value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `DateTime` Parsed date time.
- - - -


{{< header href="to_float" >}}

TO_FLOAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_float.go#L20)

TO_FLOAT takes an input value of any type and convert it into a float value. None and false are converted to the value 0 true is converted to 1 Numbers keep their original value Strings are converted to their numeric equivalent if the string contains a valid representation of a number. String values that do not contain any valid representation of a number will be converted to the number 0. An empty array is converted to 0, an array with one member is converted into the result of TO_NUMBER() for its sole member. An array with two or more members is converted to the number 0. An object / HTML node is converted to the number 0.

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

TO_INT takes an input value of any type and convert it into an integer value. None and false are converted to the value 0 true is converted to 1 Numbers keep their original value Strings are converted to their numeric equivalent if the string contains a valid representation of a number. String values that do not contain any valid representation of a number will be converted to the number 0. An empty array is converted to 0, an array with one member is converted into the result of TO_NUMBER() for its sole member. An array with two or more members is converted to the number 0. An object / HTML node is converted to the number 0.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Int` An integer value.
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
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/type_name.go#L12)

TYPENAME returns the data type name of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Any`  |  | Input value of arbitrary type.


**Returns** `Boolean` Returns string representation of a type.
- - - -
