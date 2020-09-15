

---
title: "objects"
weight: 1
draft: false
menu: [HAS,KEEP_KEYS,KEYS,MERGE,MERGE_RECURSIVE,VALUES,ZIP,]
---



{{< header >}}
HAS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/has.go#L13)

HAS returns the value stored by the given key.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`key` | `String`  |  | The key name string.


**Returns** `Boolean` True if the key exists else false.
- - - -


{{< header >}}
KEEP_KEYS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/keep_keys.go#L14)

KEEP_KEYS returns a new object with only given keys.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`obj` | `Object`  |  | Source object.
`keys` | `String, repeated`  |  | Keys that need to be kept.


**Returns** `Object` New object with only given keys.
- - - -


{{< header >}}
KEYS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/keys.go#L15)

KEYS returns string array of object's keys

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`obj` | `Object`  |  | The object whose keys you want to extract
`sort` | `Boolean`  | `False` | If sort is true, then the returned keys will be sorted.


**Returns** `String[]` Array that contains object keys.
- - - -


{{< header >}}
MERGE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/merge.go#L13)

MERGE merge the given objects into a single object.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`objects` | `Object, repeated`  |  | Objects to merge.


**Returns** `Object` Object created by merging.
- - - -


{{< header >}}
MERGE_RECURSIVE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/merge_recursive.go#L13)

MERGE_RECURSIVE recursively merge the given objects into a single object.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`objects` | `Objects, repeated`  |  | Objects to merge.


**Returns** `Object` Object created by merging.
- - - -


{{< header >}}
VALUES
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/values.go#L13)

VALUES return the attribute values of the object as an array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`object` | `Object`  |  | Target object.


**Returns** `Any[]` Values of document returned in any order.
- - - -


{{< header >}}
ZIP
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/objects/zip.go#L16)

ZIP returns an object assembled from the separate parameters keys and values. Keys and values must be arrays and have the same length.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`keys` | `String[]`  |  | An array of strings, to be used as key names in the result.
`values` | `Object[]`  |  | An array of core.value, to be used as key values.


**Returns** `Object` An object with the keys and values assembled.
- - - -
