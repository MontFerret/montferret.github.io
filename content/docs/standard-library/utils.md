---
title: "Utils"
weight: 120
draft: false
description: "Utility functions in the Ferret standard library."
aliases:
  - /docs/stdlib/utils/
menuTitle: 
menu: [print,wait,]
---



{{< header href="print" >}}

print

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/utils/log.go#L12)

print writes messages into the system log.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`message` | `Value, repeated`  |  | Print message.


**Returns** `None`
- - - -


{{< header href="wait" >}}

wait

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/utils/wait.go#L12)

wait pauses the execution for a given period.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`timeout` | `Any`  |  | Non-negative value coercible to Duration. Numbers are milliseconds.


**Returns** `None`
- - - -

## Next steps

{{< docs-related tiles="stdlib,tools-cli-debug,language-control-flow-error-handling" >}}
