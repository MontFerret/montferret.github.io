---
title: "Collections"
weight: 20
draft: false
description: "Collection functions in the Ferret standard library."
aliases:
  - /docs/stdlib/collections/
menuTitle: 
menu: [COUNT,COUNT_DISTINCT,INCLUDES,LENGTH,REVERSE,]
---



{{< header href="includes" >}}

INCLUDES

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/collections/include.go#L14)

INCLUDES checks whether a container includes a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`haystack` | `String` `Any[]` `Object` `Iterable`  |  | The value container.
`needle` | `Any`  |  | The target value to assert.


**Returns** `Boolean` A boolean value that indicates whether a container contains a given value.
- - - -


{{< header href="count" >}}

COUNT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/collections/count.go#L9)

COUNT returns the number of elements in a collection.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`collection` | `Collection`  |  | The collection to count.


**Returns** `Int` The number of elements in the collection.
- - - -


{{< header href="count_distinct" >}}

COUNT_DISTINCT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/collections/count_distinct.go#L11)

COUNT_DISTINCT computes the number of distinct elements in the given collection and returns the count as an integer.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`collection` | `Collection`  |  | The collection to count distinct elements in.


**Returns** `Int` The number of distinct elements in the collection.
- - - -


{{< header href="length" >}}

LENGTH

{{</ header >}}

LENGTH returns the number of elements or characters in a value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `String` `Any[]` `Object` `Binary`  |  | The value to measure.


**Returns** `Int` The length of the value — number of characters for strings, number of elements for arrays and objects, number of bytes for binary values.
- - - -


{{< header href="reverse" >}}

REVERSE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/collections/reverse.go#L13)

REVERSE returns the reverse of a given string or array value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `String` `Any[]`  |  | The string or array to reverse.


**Returns** `String` `Any[]` A reversed version of a given value.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-operators-array,language-control-flow-for" >}}
