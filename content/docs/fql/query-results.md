---
title: "Query results"
weight: 7
draft: false
---

# Query results

Unlike AQL, The result of an FQL query is an not always array of values. The individual values can be returned, not wrapped by an array.

For example, when the ``RETURN`` statement is used as the last query statment, a values gets returned as it is:

{{< code fql >}}
RETURN 1
{{</ code >}}

{{< highlight bash >}}
1
{{</ code >}}

{{< code fql >}}
RETURN { foo: "bar" }
{{</ code >}}

{{< highlight bash >}}
{ "foo": "bar" }
{{</ code >}}

However, when returning data from an iteration, the result values will be always an array:

{{< code fql >}}
FOR u IN elements
    RETURN i.href
{{</ code >}}

{{< code bash >}}
{ "foo": "bar" }
{{</ code >}}

## Result type

The result data type is in JSON format. 
All binary data gets endcoded into base64 strings.