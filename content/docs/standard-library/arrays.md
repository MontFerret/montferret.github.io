---
title: "Arrays"
weight: 10
draft: false
description: "Array functions in the Ferret standard library."
aliases:
  - /docs/stdlib/arrays/
menuTitle: 
menu: [append,first,flatten,intersection,last,minus,nth,outersection,pop,position,push,remove_nth,remove_value,remove_values,shift,slice,sorted,sorted_unique,union,union_distinct,unique,unshift,]
---



{{< header href="append" >}}

append

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/append.go#L15)

append appends a new item to an array and returns a new array with a given element. If ``uniqueOnly`` is set to true, then will add the item only if it's unique.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.
`item` | `Any`  |  | Target value to add.


**Returns** `Any[]` New array.
- - - -


{{< header href="first" >}}

first

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/first.go#L13)

first returns a first element from a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.


**Returns** `Any` First element in a given array.
- - - -


{{< header href="flatten" >}}

flatten

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/flatten.go#L18)

flatten turns an array of arrays into a flat array. All array elements in array will be expanded in the result array. Non-array elements are added as they are. The function will recurse into sub-arrays up to the specified depth. Duplicates will not be removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.
`depth` | `Int`  |  | Depth level.


**Returns** `Any[]` Flat array.
- - - -


{{< header href="intersection" >}}

intersection

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/intersection.go#L15)

intersection return the intersection of all arrays specified. The result is an array of values that occur in all arguments. The element order is random. Duplicates are removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` A single array with only the elements, which exist in all provided arrays.
- - - -


{{< header href="last" >}}

last

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/last.go#L13)

last returns the last element of an array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | The target array.


**Returns** `Any` Last element of an array.
- - - -


{{< header href="minus" >}}

minus

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/minus.go#L14)

minus return the difference of all arrays specified. The order of the result array is undefined and should not be relied on. Duplicates will be removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` An array of values that occur in the first array, but not in any of the subsequent arrays.
- - - -


{{< header href="nth" >}}

nth

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/nth.go#L16)

nth returns the element of an array at a given position. It is the same as anyArray[position] for positive positions, but does not support negative positions. If position is negative or beyond the upper bound of the array, then `none` will be returned.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | An array with elements of arbitrary type.
`index` | `Int`  |  | Position of desired element in array, positions start at 0.


**Returns** `Any` The array element at the given position.
- - - -


{{< header href="outersection" >}}

outersection

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/outersection.go#L12)

outersection return the values that occur only once across all arrays specified. The element order is random.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` A single array with only the elements that exist only once across all provided arrays.
- - - -


{{< header href="pop" >}}

pop

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/pop.go#L13)

pop returns a new array without last element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Copy of an array without last element.
- - - -


{{< header href="position" >}}

position

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/position.go#L15)

position returns a value indicating whether an element is contained in array. Optionally returns its position.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | The source array.
`value` | `Any`  |  | The target value.
`position` | `Boolean`  | `False` | Boolean value which indicates whether to return item's position.


**Returns** `Boolean` `Int` A value indicating whether an element is contained in array.
- - - -


{{< header href="push" >}}

push

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/push.go#L15)

push create a new array with appended value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`value` | `Any`  |  | Target value.
`unique` | `Boolean`  | `False` | Read indicating whether to do uniqueness check.


**Returns** `Any[]` A new array with appended value.
- - - -


{{< header href="remove_nth" >}}

remove_nth

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_nth.go#L14)

remove_nth returns a new array without an element by a given position.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`position` | `Int`  |  | Target element position.


**Returns** `Any[]` A new array without an element by a given position.
- - - -


{{< header href="remove_value" >}}

remove_value

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_value.go#L16)

remove_value returns a new array with removed all occurrences of value in a given array. Optionally with a limit to the number of removals.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`value` | `Any`  |  | Target value.
`limit` | `Int`  |  | A limit to the number of removals.


**Returns** `Any[]` A new array with removed all occurrences of value in a given array.
- - - -


{{< header href="remove_values" >}}

remove_values

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_values.go#L14)

remove_values returns a new array with removed all occurrences of values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`values` | `Any[]`  |  | Target values.


**Returns** `Any[]` A new array with removed all occurrences of values in a given array.
- - - -


{{< header href="shift" >}}

shift

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/shift.go#L13)

shift returns a new array without the first element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Copy of an array without the first element.
- - - -


{{< header href="slice" >}}

slice

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/slice.go#L15)

slice returns a new sliced array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`start` | `Int`  |  | Start position of extraction.
`length` | `Int`  |  | Read indicating how many elements to extract.


**Returns** `Any[]` Sliced array.
- - - -


{{< header href="sorted" >}}

sorted

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted.go#L14)

sorted sorts all elements in anyArray. The function will use the default comparison order for FQL value types.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Sorted array.
- - - -


{{< header href="sorted_unique" >}}

sorted_unique

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted_unique.go#L15)

sorted_unique sorts all elements in anyArray. The function will use the default comparison order for FQL value types. Additionally, the values in the result array will be made unique

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Sorted array.
- - - -


{{< header href="union" >}}

union

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/union.go#L13)

union returns the union of all passed arrays.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | List of arrays to combine.


**Returns** `Any[]` All array elements combined in a single array, in any order.
- - - -


{{< header href="union_distinct" >}}

union_distinct

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/union_distinct.go#L13)

union_distinct returns the union of all passed arrays with unique values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | List of arrays to combine.


**Returns** `Any[]` All unique array elements combined in a single array, in any order.
- - - -


{{< header href="unique" >}}

unique

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unique.go#L13)

unique returns all unique elements from a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` New array without duplicates.
- - - -


{{< header href="unshift" >}}

unshift

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unshift.go#L15)

unshift prepends value to a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.
`value` | `Any`  |  | Target value to prepend.
`unique` | `Boolean`  | `False` | Optional value indicating whether a value must be unique to be prepended. default is false.


**Returns** `Any[]` New array with prepended value.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-operators-array,language-types-basic" >}}
