---
title: "utils"
weight: 1
draft: false
menuTitle: 
menu: [PRINT,WAIT,]
---



{{< header href="print" >}}

PRINT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/utils/log.go#L12)

PRINT writes messages into the system log.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`message` | `Value, repeated`  |  | Print message.


**Returns** `None`
- - - -


{{< header href="wait" >}}

WAIT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/utils/wait.go#L12)

WAIT pauses the execution for a given period.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`timeout` | `Int` `Float`  |  | Number value which indicates for how long to stop an execution.


**Returns** `None`
- - - -
