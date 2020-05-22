---
title: "arrays"
weight: 1
draft: false
---


## FIRST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/first.go#L13)

First returns a first element from a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arr` | `Array` | Target array.


**Returns** `Value` First element in a given array.
- - - -

## REMOVE_VALUE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_value.go#L16)

RemoveValue returns a new array with removed all occurrences of value in a given array. Optionally with a limit to the number of removals.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Source array.
`value` | `Value` | Target value.
`limit` | `Int, optional` | A limit to the number of removals.


**Returns** `Array` A new array with removed all occurrences of value in a given array.
- - - -

## INTERSECTION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/intersection.go#L15)

Intersection return the intersection of all arrays specified. The result is an array of values that occur in all arguments.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arrays` | `Array, repeated` | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Array` A single array with only the elements, which exist in all provided arrays. the element order is random. duplicates are removed.
- - - -

## FLATTEN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/flatten.go#L18)

Flatten turn an array of arrays into a flat array. All array elements in array will be expanded in the result array. Non-array elements are added as they are. The function will recurse into sub-arrays up to the specified depth. Duplicates will not be removed.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arr` | `Array` | Target array.
`depth` | `Int, optional` | Depth level.


**Returns** `Array` Flat array.
- - - -

## SORTED_UNIQUE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted_unique.go#L15)

SortedUnique sorts all elements in anyArray. The function will use the default comparison order for FQL value types. Additionally, the values in the result array will be made unique

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.


**Returns** `Array` Sorted array.
- - - -

## SHIFT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/shift.go#L13)

Shift returns a new array without the first element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.


**Returns** `Array` Copy of an array without the first element.
- - - -

## UNSHIFT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unshift.go#L16)

Unshift prepends value to a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.
`value` | `Value` | Target value to prepend.
`unique` | `Boolean, optional` | Optional value indicating whether a value must be unique to be prepended. default is false.


**Returns** `Array` New array with prepended value.
- - - -

## POSITION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/position.go#L14)

Position returns a value indicating whether an element is contained in array. Optionally returns its position.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | The source array.
`value` | `Value` | The target value.
`returnIndex` | `Boolean, optional` | Read which indicates whether to return item's position.


**Returns** `None`
- - - -

## POP
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/pop.go#L13)

Pop returns a new array without last element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.


**Returns** `Array` Copy of an array without last element.
- - - -

## PUSH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/push.go#L15)

Push create a new array with appended value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Source array.
`value` | `Value` | Target value.
`unique` | `Boolean, optional` | Read indicating whether to do uniqueness check.


**Returns** `Array` A new array with appended value.
- - - -

## SLICE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/slice.go#L15)

Slice returns a new sliced array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Source array.
`start` | `Int` | Start position of extraction.
`length` | `Int, optional` | Read indicating how many elements to extract.


**Returns** `Array` Sliced array.
- - - -

## UNIQUE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/unique.go#L13)

Unique returns all unique elements from a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.


**Returns** `Array` New array without duplicates.
- - - -

## OUTERSECTION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/outersection.go#L12)

Outersection return the values that occur only once across all arrays specified.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arrays` | `Array, repeated` | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Array` A single array with only the elements that exist only once across all provided arrays. the element order is random.
- - - -

## SORTED
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/sorted.go#L14)

Sorted sorts all elements in anyArray. The function will use the default comparison order for FQL value types.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Target array.


**Returns** `Array` Sorted array.
- - - -

## REMOVE_VALUES
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_values.go#L14)

RemoveValues returns a new array with removed all occurrences of values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Source array.
`values` | `Array` | Target values.


**Returns** `Array` A new array with removed all occurrences of values in a given array.
- - - -

## MINUS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/minus.go#L14)

Minus return the difference of all arrays specified.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arrays` | `Array, repeated` | An arbitrary number of arrays as multiple arguments (at least 2).


**Returns** `Array` An array of values that occur in the first array, but not in any of the subsequent arrays. the order of the result array is undefined and should not be relied on. duplicates will be removed.
- - - -

## LAST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/last.go#L13)

Last returns the last element of an array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | The target array.


**Returns** `Value` Last element of an array.
- - - -

## REMOVE_NTH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/remove_nth.go#L14)

RemoveNth returns a new array without an element by a given position.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Source array.
`position` | `Int` | Target element position.


**Returns** `Array` A new array without an element by a given position.
- - - -

## NTH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/nth.go#L16)

Nth returns the element of an array at a given position. It is the same as anyArray[position] for positive positions, but does not support negative positions.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | An array with elements of arbitrary type.
`index` | `Int` | Position of desired element in array, positions start at 0.


**Returns** `Value` The array element at the given position. if position is negative or beyond the upper bound of the array, then none will be returned.
- - - -

## APPEND
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/append.go#L15)

Append appends a new item to an array and returns a new array with a given element. If ``uniqueOnly`` is set to true, then will add the item only if it's unique.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arr` | `Array` | Target array.
`item` | `Value` | Target value to add.


**Returns** `Array` New array.
- - - -

## UNION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/arrays/union.go#L13)

Union returns the union of all passed arrays.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`arrays` | `Array, repeated` | List of arrays to combine.


**Returns** `Array` All array elements combined in a single array, in any order.
- - - -
