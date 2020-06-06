---
title: "Query results"
weight: 7
draft: false
---

# Query results

Unlike AQL, the result of an FQL query is an not always array of values. The individual values can be returned, not wrapped by an array.

For example, when the ``RETURN`` statement is used as the last query statement, a values gets returned as it is:

{{< code lang="fql" height="90px" >}}
RETURN 1
{{</ code >}}

{{< code lang="fql" height="90px" >}}
1
{{</ code >}}

{{< code lang="fql" height="90px" >}}
RETURN { foo: "bar" }
{{</ code >}}

{{< code lang="fql" height="90px" >}}
{ "foo": "bar" }
{{</ code >}}

However, when returning data from an iteration, the result values will be always an array:

{{< code lang="fql" height="100px" >}}
FOR u IN elements
    RETURN i.href
{{</ code >}}

{{< code lang="fql" height="90px" >}}
{ "foo": "bar" }
{{</ code >}}

## Result type

The result data type is in JSON format. 
All binary data gets encoded into base64 strings.