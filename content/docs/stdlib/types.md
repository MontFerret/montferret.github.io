---
title: "types"
weight: 1
draft: false
---


## IS_BINARY
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_binary.go#L13)

IsBinary checks whether value is a binary value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is binary, otherwise false.
- - - -

## TO_ARRAY
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_array.go#L15)

toArray takes an input value of any type and convert it into an array value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Array` None is converted to an empty array boolean values, numbers and strings are converted to an array containing the original value as its single element arrays keep their original value objects / html nodes are converted to an array containing their attribute values as array elements.
- - - -

## IS_DATETIME
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_date_time.go#L13)

IsDateTime checks whether value is a date time value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is date time, otherwise false.
- - - -

## IS_BOOL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_boolean.go#L13)

IsBool checks whether value is a boolean value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is boolean, otherwise false.
- - - -

## IS_FLOAT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_float.go#L13)

IsFloat checks whether value is a float value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is float, otherwise false.
- - - -

## TO_INT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_int.go#L20)

ToInt takes an input value of any type and convert it into an integer value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Int` None and false are converted to the value 0 true is converted to 1 numbers keep their original value strings are converted to their numeric equivalent if the string contains a valid representation of a number. string values that do not contain any valid representation of a number will be converted to the number 0. an empty array is converted to 0, an array with one member is converted into the result of to_number() for its sole member. an array with two or more members is converted to the number 0. an object / html node is converted to the number 0.
- - - -

## IS_STRING
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_string.go#L13)

IsString checks whether value is a string value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is string, otherwise false.
- - - -

## IS_INT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_int.go#L13)

IsInt checks whether value is a int value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is int, otherwise false.
- - - -

## IS_NONE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_none.go#L13)

IsNone checks whether value is a none value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is none, otherwise false.
- - - -

## IS_NAN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_nan.go#L13)

IsNaN checks whether value is NaN.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is nan, otherwise false.
- - - -

## IS_HTML_DOCUMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_html_document.go#L13)

IsHTMLDocument checks whether value is a HTMLDocument value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is htmldocument, otherwise false.
- - - -

## TO_FLOAT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_float.go#L20)

ToFloat takes an input value of any type and convert it into a float value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Float` None and false are converted to the value 0 true is converted to 1 numbers keep their original value strings are converted to their numeric equivalent if the string contains a valid representation of a number. string values that do not contain any valid representation of a number will be converted to the number 0. an empty array is converted to 0, an array with one member is converted into the result of to_number() for its sole member. an array with two or more members is converted to the number 0. an object / html node is converted to the number 0.
- - - -

## TO_BOOL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_boolean.go#L18)

ToBool takes an input value of any type and converts it into the appropriate boolean value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` None is converted to false numbers are converted to true, except for 0, which is converted to false strings are converted to true if they are non-empty, and to false otherwise dates are converted to true if they are not zero, and to false otherwise arrays are always converted to true (even if empty) objects / htmlnodes / binary are always converted to true
- - - -

## IS_HTML_ELEMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_html_element.go#L13)

IsHTMLElement checks whether value is a HTMLElement value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is htmlelement, otherwise false.
- - - -

## TO_STRING
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_string.go#L12)

ToString takes an input value of any type and convert it into a string value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `String` String representation of a given value.
- - - -

## TO_DATETIME
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/to_date_time.go#L12)

ToDateTime takes an input value of any type and converts it into the appropriate date time value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `DateTime` Parsed date time.
- - - -

## IS_OBJECT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_object.go#L13)

IsObject checks whether value is an object value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is object, otherwise false.
- - - -

## IS_ARRAY
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/is_array.go#L13)

IsArray checks whether value is an array value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns true if value is array, otherwise false.
- - - -

## TYPENAME
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/types/type_name.go#L12)

TypeName returns the data type name of value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | Input value of arbitrary type.


**Returns** `Boolean` Returns string representation of a type.
- - - -
