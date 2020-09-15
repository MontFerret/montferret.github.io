---
title: "path"
weight: 1
draft: false
menuTitle: PATH
menu: [BASE,CLEAN,DIR,EXT,IS_ABS,JOIN,MATCH,SEPARATE,]
---



{{< header href="base" >}}

PATH::BASE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/base.go#L14)

BASE returns the last component of the path or the path itself if it does not contain any directory separators.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The last component of the path.
- - - -


{{< header href="clean" >}}

PATH::CLEAN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/clean.go#L14)

CLEAN returns the shortest path name equivalent to path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The shortest path name equivalent to path
- - - -


{{< header href="dir" >}}

PATH::DIR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/dir.go#L14)

DIR returns the directory component of path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The directory component of path.
- - - -


{{< header href="ext" >}}

PATH::EXT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/ext.go#L14)

EXT returns the extension of the last component of path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The extension of the last component of path.
- - - -


{{< header href="is_abs" >}}

PATH::IS_ABS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/is_abs.go#L14)

IS_ABS reports whether the path is absolute.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `Boolean` True if the path is absolute.
- - - -


{{< header href="join" >}}

PATH::JOIN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/join.go#L14)

JOIN joins any number of path elements into a single path, separating them with slashes.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`elements` | `String, repeated` `String[]`  |  | The path elements


**Returns** `String` Single path from the given elements.
- - - -


{{< header href="match" >}}

PATH::MATCH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/match.go#L15)

MATCH reports whether name matches the pattern.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`pattern` | `String`  |  | The pattern.
`name` | `String`  |  | The name.


**Returns** `Boolean` True if the name matches the pattern.
- - - -


{{< header href="separate" >}}

PATH::SEPARATE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/separate.go#L14)

SEPARATE separates the path into a directory and filename component.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path


**Returns** `Any[]` First item is a directory component, and second is a filename component.
- - - -
