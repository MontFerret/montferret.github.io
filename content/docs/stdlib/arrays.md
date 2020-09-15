---
title: "arrays"
weight: 1
draft: false
menuTitle: 
menu: [APPEND,FIRST,FLATTEN,INTERSECTION,LAST,MINUS,NTH,OUTERSECTION,POP,POSITION,PUSH,REMOVE_NTH,REMOVE_VALUE,REMOVE_VALUES,SHIFT,SLICE,SORTED,SORTED_UNIQUE,UNION,UNION_DISTINCT,UNIQUE,UNSHIFT,]
---



{{< header href="append" >}}

APPEND

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/append.go#L15)

APPEND appends a new item to an array and returns a new array with a given element. If ``uniqueOnly`` is set to true, then will add the item only if it's unique.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.
`item` | `Any`  |  | Target value to add.


**Returns** `Any[]` New array.
- - - -


{{< header href="first" >}}

FIRST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/first.go#L13)

FIRST returns a first element from a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.


**Returns** `Any` First element in a given array.
- - - -


{{< header href="flatten" >}}

FLATTEN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/flatten.go#L18)

FLATTEN turns an array of arrays into a flat array. All array elements in array will be expanded in the result array. Non-array elements are added as they are. The function will recurse into sub-arrays up to the specified depth. Duplicates will not be removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arr` | `Any[]`  |  | Target array.
`depth` | `Int`  |  | Depth level.


**Returns** `Any[]` Flat array.
- - - -


{{< header href="intersection" >}}

INTERSECTION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/intersection.go#L15)

INTERSECTION return the intersection of all arrays specified. The result is an array of values that occur in all arguments. The element order is random. Duplicates are removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` A single array with only the elements, which exist in all provided arrays.
- - - -


{{< header href="last" >}}

LAST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/last.go#L13)

LAST returns the last element of an array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | The target array.


**Returns** `Any` Last element of an array.
- - - -


{{< header href="minus" >}}

MINUS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/minus.go#L14)

MINUS return the difference of all arrays specified. The order of the result array is undefined and should not be relied on. Duplicates will be removed.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` An array of values that occur in the first array, but not in any of the subsequent arrays.
- - - -


{{< header href="nth" >}}

NTH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/nth.go#L16)

NTH returns the element of an array at a given position. It is the same as anyArray[position] for positive positions, but does not support negative positions. If position is negative or beyond the upper bound of the array, then NONE will be returned.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | An array with elements of arbitrary type.
`index` | `Int`  |  | Position of desired element in array, positions start at 0.


**Returns** `Any` The array element at the given position.
- - - -


{{< header href="outersection" >}}

OUTERSECTION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/outersection.go#L12)

OUTERSECTION return the values that occur only once across all arrays specified. The element order is random.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Any[]` A single array with only the elements that exist only once across all provided arrays.
- - - -


{{< header href="pop" >}}

POP

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/pop.go#L13)

POP returns a new array without last element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Copy of an array without last element.
- - - -


{{< header href="position" >}}

POSITION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/position.go#L15)

POSITION returns a value indicating whether an element is contained in array. Optionally returns its position.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | The source array.
`value` | `Any`  |  | The target value.
`position` | `Boolean`  | `False` | Boolean value which indicates whether to return item's position.


**Returns** `Boolean` `Int` A value indicating whether an element is contained in array.
- - - -


{{< header href="push" >}}

PUSH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/push.go#L15)

PUSH create a new array with appended value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`value` | `Any`  |  | Target value.
`unique` | `Boolean`  | `False` | Read indicating whether to do uniqueness check.


**Returns** `Any[]` A new array with appended value.
- - - -


{{< header href="remove_nth" >}}

REMOVE_NTH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_nth.go#L14)

REMOVE_NTH returns a new array without an element by a given position.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`position` | `Int`  |  | Target element position.


**Returns** `Any[]` A new array without an element by a given position.
- - - -


{{< header href="remove_value" >}}

REMOVE_VALUE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_value.go#L16)

REMOVE_VALUE returns a new array with removed all occurrences of value in a given array. Optionally with a limit to the number of removals.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`value` | `Any`  |  | Target value.
`limit` | `Int`  |  | A limit to the number of removals.


**Returns** `Any[]` A new array with removed all occurrences of value in a given array.
- - - -


{{< header href="remove_values" >}}

REMOVE_VALUES

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_values.go#L14)

REMOVE_VALUES returns a new array with removed all occurrences of values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`values` | `Any[]`  |  | Target values.


**Returns** `Any[]` A new array with removed all occurrences of values in a given array.
- - - -


{{< header href="shift" >}}

SHIFT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/shift.go#L13)

SHIFT returns a new array without the first element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Copy of an array without the first element.
- - - -


{{< header href="slice" >}}

SLICE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/slice.go#L15)

SLICE returns a new sliced array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Source array.
`start` | `Int`  |  | Start position of extraction.
`length` | `Int`  |  | Read indicating how many elements to extract.


**Returns** `Any[]` Sliced array.
- - - -


{{< header href="sorted" >}}

SORTED

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted.go#L14)

SORTED sorts all elements in anyArray. The function will use the default comparison order for FQL value types.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Sorted array.
- - - -


{{< header href="sorted_unique" >}}

SORTED_UNIQUE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted_unique.go#L15)

SORTED_UNIQUE sorts all elements in anyArray. The function will use the default comparison order for FQL value types. Additionally, the values in the result array will be made unique

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` Sorted array.
- - - -


{{< header href="union" >}}

UNION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/union.go#L13)

UNION returns the union of all passed arrays.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | List of arrays to combine.


**Returns** `Any[]` All array elements combined in a single array, in any order.
- - - -


{{< header href="union_distinct" >}}

UNION_DISTINCT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/union_distinct.go#L13)

UNION_DISTINCT returns the union of all passed arrays with unique values.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`arrays` | `Any[], repeated`  |  | List of arrays to combine.


**Returns** `Any[]` All unique array elements combined in a single array, in any order.
- - - -


{{< header href="unique" >}}

UNIQUE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unique.go#L13)

UNIQUE returns all unique elements from a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.


**Returns** `Any[]` New array without duplicates.
- - - -


{{< header href="unshift" >}}

UNSHIFT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unshift.go#L15)

UNSHIFT prepends value to a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Any[]`  |  | Target array.
`value` | `Any`  |  | Target value to prepend.
`unique` | `Boolean`  | `False` | Optional value indicating whether a value must be unique to be prepended. default is false.


**Returns** `Any[]` New array with prepended value.
- - - -
