---
title: "Query results"
weight: 6
draft: false
---

# Query results

Unlike AQL, The result of an FQL query is an not always array of values. The individual values can be returned, not wrapped by an array.

For example, when the ``RETURN`` statement is used as the last query statment, a values gets returned as it is:

{{< highlight javascript >}}
RETURN 1
{{</ highlight >}}

{{< highlight bash >}}
1
{{</ highlight >}}

{{< highlight javascript >}}
RETURN { foo: "bar" }
{{</ highlight >}}

{{< highlight bash >}}
{ "foo": "bar" }
{{</ highlight >}}

However, when returning data from an iteration, the result values will be always an array:

{{< highlight javascript >}}
FOR u IN elements
    RETURN i.href
{{</ highlight >}}

{{< highlight bash >}}
{ "foo": "bar" }
{{</ highlight >}}

## Result type

The result data type is in JSON format. 
All binary data gets endcoded into base64 strings.