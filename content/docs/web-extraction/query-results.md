---
title: "Query Results"
weight: 10
draft: false
description: "Understand the shape and handling of Ferret query output."
aliases:
  - /docs/fql/query-results/
---

# Query results

`RETURN` produces the value of its expression. For example, `RETURN 1` produces the number `1`. A `FOR` expression collects the values returned by its iterations into an array:

{{< editor lang="fql" >}}
RETURN FOR link IN [{ href: "/docs/" }, { href: "/try/" }]
    RETURN link.href
{{< /editor >}}

The result is `["/docs/", "/try/"]`. A loop with no matching iterations produces `[]`.

Ferret values are encoded when the host delivers output. The CLI’s default JSON encoding represents objects and arrays directly and binary values as base64 strings. See [Output]({{< ref "/docs/embedding/go/output" >}}) for consumption and codec options in Go.
