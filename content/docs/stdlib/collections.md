

---
title: "collections"
weight: 1
draft: false
menu: [INCLUDES,LENGTH,REVERSE,]
---



{{< header >}}
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


{{< header >}}
LENGTH
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/collections/length.go#L14)

LENGTH returns the length of a measurable value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Measurable`  |  | The value to measure.


**Returns** `Int` The length of the value.
- - - -


{{< header >}}
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
