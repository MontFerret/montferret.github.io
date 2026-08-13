---
title: "Testing"
weight: 100
draft: false
description: "Testing assertion functions in the Ferret standard library."
aliases:
  - /docs/stdlib/testing/
menuTitle: T
menu: [array,binary,datetime,empty,eq,fail,false,float,gt,gte,include,int,len,lt,lte,match,none,object,string,true,]
---



{{< header href="array" >}}

t::array

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/array.go#L13)

array asserts that value is a array type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="binary" >}}

t::binary

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/binary.go#L13)

binary asserts that value is a binary type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="datetime" >}}

t::datetime

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/datetime.go#L13)

datetime asserts that value is a datetime type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="empty" >}}

t::empty

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/empty.go#L14)

empty asserts that the target does not contain any values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Measurable` `Binary` `Object` `Any[]` `String`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="eq" >}}

t::eq

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/equal.go#L14)

eq asserts equality of actual and expected values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="fail" >}}

t::fail

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/fail.go#L11)

fail returns an error.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="false" >}}

t::false

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/false.go#L14)

false asserts that value is false.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="float" >}}

t::float

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/float.go#L13)

float asserts that value is a float type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="gt" >}}

t::gt

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/gt.go#L14)

gt asserts that an actual value is greater than an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="gte" >}}

t::gte

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/gte.go#L14)

gte asserts that an actual value is greater than or equal to an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="include" >}}

t::include

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/include.go#L16)

include asserts that haystack includes needle.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `String` `Array` `Object` `Iterable`  |  | Haystack value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="int" >}}

t::int

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/int.go#L13)

int asserts that value is a int type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="len" >}}

t::len

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/len.go#L15)

len asserts that a measurable value has a length or size with the expected value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Measurable`  |  | Measurable value.
`length` | `Int`  |  | Target length.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="lt" >}}

t::lt

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/lt.go#L14)

lt asserts that an actual value is lesser than an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="lte" >}}

t::lte

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/lte.go#L14)

lte asserts that an actual value is lesser than or equal to an expected one.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expected` | `Any`  |  | Expected value.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="match" >}}

t::match

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/match.go#L15)

match asserts that value matches the regular expression.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Actual value.
`expression` | `String`  |  | Regular expression.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="none" >}}

t::none

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/none.go#L14)

none asserts that value is none.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="object" >}}

t::object

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/object.go#L13)

object asserts that value is a object type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="string" >}}

t::string

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/string.go#L13)

string asserts that value is a string type.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -


{{< header href="true" >}}

t::true

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/testing/true.go#L14)

true asserts that value is true.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`actual` | `Any`  |  | Value to test.
`message` | `String`  |  | Message to display on error.


**Returns** `None`
- - - -

## Next steps

{{< docs-related tiles="tools-lab,tools-lab-writing-tests,tools-lab-running-tests" >}}
