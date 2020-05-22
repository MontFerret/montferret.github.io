---
title: "io/fs"
weight: 1
draft: false
---


## READ
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/fs/read.go#L14)

Read reads from a given file.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | Path to file to read from.


**Returns** `Binary` The read file in binary format.
- - - -

## WRITE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/fs/write.go#L22)

Write writes the given data into the file.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`path` | `String` | Path to file to write into.
`data` | `Binary` | Data to write.
`params` | `Object, optional` | Additional parameters: * mode (string): * x - exclusive: returns an error if the file exist. it can be combined with other modes * a - append: will create a file if the specified file does not exist * w - write (default): will create a file if the specified file does not exist


**Returns** `None` Returns nothing.
- - - -
