---
title: "IO / FS"
sidebarTitle: "FS"
weight: 10
draft: false
description: "File system functions in the Ferret standard library."
aliases:
  - /docs/stdlib/io-fs/
menuTitle: io::fs
menu: [read,write,]
---



{{< header href="read" >}}

io::fs::read

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/fs/read.go#L14)

read reads from a given file.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | Path to file to read from.


**Returns** `Binary` File content in binary format.
- - - -


{{< header href="write" >}}

io::fs::write

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/io/fs/write.go#L20)

write writes the given data into the file.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`path` | `String`  |  | File path to write into.
`data` | `Binary`  |  | Data to write.
`params` | `Object`  |  | Additional parameters:
`params.mode` | `String`  |  | Write (default): will create a file if the specified file does not exist


**Returns** `None`
- - - -

## Next steps

{{< docs-related tiles="stdlib,embedding,tools-cli-run" >}}
