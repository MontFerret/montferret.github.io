---
title: "testing"
weight: 1
draft: false
menuTitle: T
menu: [ARRAY,BINARY,DATETIME,EMPTY,EQUAL,FAIL,FALSE,FLOAT,GT,GTE,INCLUDE,INT,LEN,LT,LTE,MATCH,NONE,OBJECT,STRING,TRUE,]
---



{{< header href="array" >}}

T::ARRAY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/array.go#L13)

ARRAY asserts that value is a array type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="binary" >}}

T::BINARY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/binary.go#L13)

BINARY asserts that value is a binary type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="datetime" >}}

T::DATETIME

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/datetime.go#L13)

DATETIME asserts that value is a datetime type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="empty" >}}

T::EMPTY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/empty.go#L14)

EMPTY asserts that the target does not contain any values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Measurable` `Binary` `Object` `Any[]` `String`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="equal" >}}

T::EQUAL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/equal.go#L14)

EQUAL asserts equality of actual and expected values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="fail" >}}

T::FAIL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/fail.go#L11)

FAIL returns an error.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="false" >}}

T::FALSE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/false.go#L14)

FALSE asserts that value is false.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="float" >}}

T::FLOAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/float.go#L13)

FLOAT asserts that value is a float type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="gt" >}}

T::GT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/gt.go#L14)

GT asserts that an actual value is greater than an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="gte" >}}

T::GTE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/gte.go#L14)

GTE asserts that an actual value is greater than or equal to an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="include" >}}

T::INCLUDE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/include.go#L16)

INCLUDE asserts that haystack includes needle.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `String` `Array` `Object` `Iterable`  |  | Haystack value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="int" >}}

T::INT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/int.go#L13)

INT asserts that value is a int type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="len" >}}

T::LEN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/len.go#L15)

LEN asserts that a measurable value has a length or size with the expected value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Measurable`  |  | Measurable value.
`length` | `Int`  |  | Target length.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="lt" >}}

T::LT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/lt.go#L14)

LT asserts that an actual value is lesser than an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="lte" >}}

T::LTE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/lte.go#L14)

LTE asserts that an actual value is lesser than or equal to an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="match" >}}

T::MATCH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/match.go#L15)

MATCH asserts that value matches the regular expression.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expression` | `String`  |  | Regular expression.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="none" >}}

T::NONE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/none.go#L14)

NONE asserts that value is none.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="object" >}}

T::OBJECT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/object.go#L13)

OBJECT asserts that value is a object type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="string" >}}

T::STRING

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/string.go#L13)

STRING asserts that value is a string type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="true" >}}

T::TRUE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/true.go#L14)

TRUE asserts that value is true.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -
