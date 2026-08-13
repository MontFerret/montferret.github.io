---
title: "Path"
weight: 70
draft: false
description: "Path functions in the Ferret standard library."
aliases:
  - /docs/stdlib/path/
menuTitle: PATH
menu: [base,clean,dir,ext,is_abs,join,match,separate,]
---



{{< header href="base" >}}

path::base

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/base.go#L14)

base returns the last component of the path or the path itself if it does not contain any directory separators.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The last component of the path.
- - - -


{{< header href="clean" >}}

path::clean

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/clean.go#L14)

clean returns the shortest path name equivalent to path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The shortest path name equivalent to path
- - - -


{{< header href="dir" >}}

path::dir

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/dir.go#L14)

dir returns the directory component of path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The directory component of path.
- - - -


{{< header href="ext" >}}

path::ext

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/ext.go#L14)

ext returns the extension of the last component of path.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `String` The extension of the last component of path.
- - - -


{{< header href="is_abs" >}}

path::is_abs

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/is_abs.go#L14)

is_abs reports whether the path is absolute.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path.


**Returns** `Boolean` True if the path is absolute.
- - - -


{{< header href="join" >}}

path::join

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/join.go#L14)

join joins any number of path elements into a single path, separating them with slashes.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`elements` | `String, repeated` `String[]`  |  | The path elements


**Returns** `String` Single path from the given elements.
- - - -


{{< header href="match" >}}

path::match

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/match.go#L15)

match reports whether name matches the pattern.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`pattern` | `String`  |  | The pattern.
`name` | `String`  |  | The name.


**Returns** `Boolean` True if the name matches the pattern.
- - - -


{{< header href="separate" >}}

path::separate

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/path/separate.go#L14)

separate separates the path into a directory and filename component.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | The path


**Returns** `Any[]` First item is a directory component, and second is a filename component.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-types-basic,embedding-custom-functions" >}}
