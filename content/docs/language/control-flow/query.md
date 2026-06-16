---
title: "Query Expressions"
sidebarTitle: "Query"
weight: 40
draft: false
description: "Delegate a query to a host value using the QUERY ... IN ... USING ... expression."
---

# Query Expressions

A `QUERY` expression delegates work to a value that knows how to run queries. Instead of expressing the logic in FQL, you write a query in another dialect — such as CSS for an HTML document — and `QUERY` passes it to the value and returns the results.

{{< code lang="fql" >}}
QUERY `.product .title` IN doc USING css
{{</ code >}}

This reads as: run the query `` `.product .title` `` against `doc`, using the `css` dialect.

The query string is often written with backticks so selectors and other syntax do not need escaping. It can also be a regular string, a variable, or a bind parameter.

## Anatomy

{{< code lang="fql" >}}
QUERY [modifier] <payload> IN <source> USING <dialect> [WITH <params>] [OPTIONS <options>]
{{</ code >}}

- **payload** — the query to run, as a string, variable, or parameter.
- **source** — an expression that yields the value the query runs against. The value must support querying.
- **dialect** — an identifier naming the query language, such as `css`, `xpath`, or `sql`. The available dialects depend on the modules and runtime in use.
- **WITH** — an optional value passed to the query as parameters.
- **OPTIONS** — an optional value carrying execution settings, such as a timeout.

## Result modifiers

By default, a `QUERY` returns a list of every match. A modifier after `QUERY` changes the shape of the result.

| Form | Returns |
| --- | --- |
| `QUERY ...` | a list of all matches |
| `QUERY ONE ...` | the first match, or `NONE` if there are none |
| `QUERY COUNT ...` | the number of matches |
| `QUERY EXISTS ...` | `true` if there is at least one match, otherwise `false` |

{{< code lang="fql" >}}
LET total = QUERY COUNT `.item` IN doc USING css
LET hasNext = QUERY EXISTS `.pagination .next` IN doc USING css
LET title = QUERY ONE `.product .title` IN doc USING css

RETURN { total, hasNext, title }
{{</ code >}}

## Parameters and options

`WITH` supplies parameters to the query, and `OPTIONS` carries execution settings. Both are evaluated once.

{{< code lang="fql" >}}
QUERY `SELECT name, price FROM products WHERE category = $c` IN db USING sql
    WITH { c: "laptops" }
    OPTIONS { timeout: 5000 }
{{</ code >}}

## A host capability

`QUERY` only works when the source value supports querying — it must be a **queryable** value, such as an HTML document, a database connection, or a JSON document provided by a module or the host application. If the source does not support querying, or the requested dialect is unavailable, the query fails at runtime.

You can recover from such failures with a recovery clause.

{{< code lang="fql" >}}
LET rows = QUERY `.row` IN doc USING css ON ERROR RETURN []
{{</ code >}}

For more on values that expose querying and other behaviors, see [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}). For querying documents in practice, see [Web Extraction]({{% ref "../../web-extraction" %}}).
