---
title: "objects"
weight: 1
draft: false
---


## KEYS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/keys.go#L15)

Keys returns string array of object's keys

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`obj` | `Object` | The object whose keys you want to extract
`sort` | `Boolean, optional` | If sort is true, then the returned keys will be sorted.


**Returns** `Array of String` Array that contains object keys.
- - - -

## HAS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/has.go#L13)

Has returns the value stored by the given key.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`key` | `String` | The key name string.


**Returns** `Boolean` True if the key exists else false.
- - - -

## VALUES
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/values.go#L13)

Values return the attribute values of the object as an array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`obj` | `Object` | An object.


**Returns** `Array of Value` The values of document returned in any order.
- - - -

## KEEP_KEYS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/keep_keys.go#L14)

KeepKeys returns a new object with only given keys.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `Object` | Source object.
`keys` | `Array Of String OR Strings` | Keys that need to be keeped.


**Returns** `Object` New object with only given keys.
- - - -

## MERGE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/merge.go#L13)

Merge merge the given objects into a single object.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`objs` | `Array Of Object OR Objects` | Objects to merge.


**Returns** `Object` Object created by merging.
- - - -

## MERGE_RECURSIVE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/merge_recursive.go#L13)

MergeRecursive recursively merge the given objects into a single object.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`objs` | `Objects` | Objects to merge.


**Returns** `Object` Object created by merging.
- - - -

## ZIP
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/zip.go#L16)

Zip returns an object assembled from the separate parameters keys and values. Keys and values must be arrays and have the same length.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`keys` | `Array of Strings` | An array of strings, to be used as key names in the result.
`values` | `Array of Objects` | An array of core.value, to be used as key values.


**Returns** `Object` An object with the keys and values assembled.
- - - -
