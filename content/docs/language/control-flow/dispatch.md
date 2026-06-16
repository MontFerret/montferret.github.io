---
title: "Dispatch Expressions"
sidebarTitle: "Dispatch"
weight: 50
draft: false
description: "Emit an event to a host value with the DISPATCH expression and its arrow shorthand."
---

# Dispatch Expressions

A `DISPATCH` expression emits an event to a value. It is used to drive the outside world — for example, to fire a `click` or `input` event on a browser element.

{{< code lang="fql" >}}
DISPATCH "click" IN element
{{</ code >}}

`DISPATCH` is performed for its effect, not its value: it runs synchronously and always evaluates to `NONE`.

## Payload and options

The `WITH` clause attaches a payload to the event, and the `OPTIONS` clause carries settings that describe how the event should be emitted.

{{< code lang="fql" >}}
DISPATCH "input" IN element WITH "hello"

DISPATCH "select" IN element
    WITH ["1", "2"]
    OPTIONS { selector: "#a", delay: 50 }
{{</ code >}}

The event name can be a literal string, a variable, or a bind parameter.

{{< code lang="fql" >}}
LET eventName = "hover"

DISPATCH eventName IN element
{{</ code >}}

## Arrow shorthand

For an event with no payload or options, the arrow operator `<-` is a concise alternative. The target is on the left and the event name on the right.

{{< code lang="fql" >}}
element <- "focus"
{{</ code >}}

The two forms are equivalent. Like the long form, the shorthand evaluates to `NONE`, so it can appear anywhere an expression is allowed — including inside a [MATCH]({{< ref "match" >}}) arm or a function body.

## A host capability

`DISPATCH` only works when the target value can receive events — it must be a **dispatchable** value provided by a module or the host application, such as a browser element. Dispatching to a value that does not support it fails at runtime. See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

To wait for events instead of emitting them, see [Waitfor Expressions]({{< ref "waitfor" >}}).
