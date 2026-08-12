---
title: "Dispatch Expressions"
sidebarTitle: "Dispatch"
weight: 50
draft: false
description: "Emit an event to a host value with the dispatch expression and its arrow shorthand."
---

# Dispatch Expressions

A `dispatch` expression emits an event to a value. It is used to drive the outside world — for example, to fire a `click` or `input` event on a browser element.

{{< code lang="fql" >}}
dispatch "click" in element
{{</ code >}}

`dispatch` is performed for its effect, not its value: it runs synchronously and always evaluates to `none`.

## Payload and options

The `with` clause attaches a payload to the event, and the `options` clause carries settings that describe how the event should be emitted.

{{< code lang="fql" >}}
dispatch "input" in element with "hello"

dispatch "select" in element
    with ["1", "2"]
    options { selector: "#a", delay: 50 }
{{</ code >}}

The event name can be a literal string, a variable, or a bind parameter.

{{< code lang="fql" >}}
let eventName = "hover"

dispatch eventName in element
{{</ code >}}

## Arrow shorthand

For an event with no payload or options, the arrow operator `<-` is a concise alternative. The target is on the left and the event name on the right.

{{< code lang="fql" >}}
element <- "focus"
{{</ code >}}

The two forms are equivalent. Like the long form, the shorthand evaluates to `none`, so it can appear anywhere an expression is allowed — including inside a [match]({{< ref "match" >}}) arm or a function body.

## A host capability

`dispatch` only works when the target value can receive events — it must be a **dispatchable** value provided by a module or the host application, such as a browser element. Dispatching to a value that does not support it fails at runtime. See [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

To wait for events instead of emitting them, see [Waitfor Expressions]({{< ref "waitfor" >}}).

## Next steps

{{< docs-related tiles="language-control-flow-waitfor,language-types-capabilities,language-control-flow" >}}
