---
title: "path"
weight: 1
draft: false
---


## MATCH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/match.go#L15)

Match reports whether name matches the pattern.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`pattern` | `String` | The pattern.
`name` | `String` | The name.


**Returns** `Boolean` True if the name mathes the pattern.
- - - -

## BASE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/base.go#L15)

Base returns the last component of the path. or the path itself if it does not contain any directory separators.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path.


**Returns** `String` The last component of the path.
- - - -

## CLEAN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/clean.go#L14)

Clean returns the shortest path name equivalent to path.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path.


**Returns** `String` The shortest path name equivalent to path
- - - -

## JOIN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/join.go#L14)

Join joins any number of path elements into a single path, separating them with slashes.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`elem` | `String...` `Array<String>` | The path elements


**Returns** `String` Single path from the given elements.
- - - -

## EXT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/ext.go#L14)

Ext returns the extension of the last component of path.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path.


**Returns** `String` The extension of the last component of path.
- - - -

## DIR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/dir.go#L14)

Dir returns the directory component of path.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path.


**Returns** `String` The directory component of path.
- - - -

## IS_ABS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/is_abs.go#L14)

IsAbs reports whether the path is absolute.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path.


**Returns** `Boolean` True if the path is absolute.
- - - -

## SEPARATE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/separate.go#L14)

Separate separates the path into a directory and filename component.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | The path


**Returns** `Array` First item is a directory component, and second is a filename component.
- - - -
