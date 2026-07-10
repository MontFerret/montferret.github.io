---
title: "Query Expressions"
sidebarTitle: "Query"
weight: 40
draft: false
description: "Delegate a query to a host value using the QUERY ... IN ... expression."
---

# Query Expressions

A `QUERY` expression delegates work to a value that knows how to run queries. Instead of expressing the logic in FQL, you provide a query payload and pass it to a host value, such as an HTML document, database connection, API client, or JSON document.

{{< code lang="fql" >}}
QUERY `.product .title` IN doc
{{</ code >}}

This reads as: run the query `` `.product .title` `` against `doc`.

Some host values support more than one query dialect. In that case, the value may provide a default dialect, so `USING` is not always required. For example, an HTML document may default to CSS selectors.

When you need to choose a specific dialect, add `USING`:

{{< code lang="fql" >}}
QUERY `//article/h1` IN doc USING xpath
{{</ code >}}

This reads as: run the query `` `//article/h1` `` against `doc`, using the `xpath` dialect.

The query expression is often written with backticks so selectors and other syntax do not need escaping. It can also be a regular string, a variable, or a bind parameter.

## Anatomy

{{< code lang="fql" >}}
QUERY [modifier] <expression> IN <source> [USING <dialect>] [WITH <params>] [OPTIONS <options>]
{{</ code >}}

- **expression** — the query to run, as a string, variable, or parameter.
- **source** — an expression that yields the value the query runs against. The value must support querying.
- **USING** — optionally selects a query dialect, such as `css`, `xpath`, or `sql`. If omitted, the source value chooses its default dialect when one is available.
- **WITH** — an optional value passed to the query as parameters.
- **OPTIONS** — an optional value carrying execution settings, such as a timeout.

If `USING` is omitted and the source value does not provide a default dialect, the query fails at runtime.

## Result modifiers

By default, a `QUERY` returns a list of every match. A modifier after `QUERY` changes the shape of the result.

| Form | Returns |
| --- | --- |
| `QUERY ...` | a list of all matches |
| `QUERY ONE ...` | the first match, or `NONE` if there are none |
| `QUERY COUNT ...` | the number of matches |
| `QUERY EXISTS ...` | `true` if there is at least one match, otherwise `false` |

{{< code lang="fql" >}}
LET total = QUERY COUNT `.item` IN doc
LET hasNext = QUERY EXISTS `.pagination .next` IN doc
LET title = QUERY ONE `.product .title` IN doc

RETURN { total, hasNext, title }
{{</ code >}}

## Shortcut syntax

The two most common query forms have shorter equivalents:

| Shortcut                             | Equivalent form                                  | Returns |
|--------------------------------------|--------------------------------------------------| --- |
| `<source>[~ <dialect>'<expression>']` | `QUERY <expression> IN <source> USING <dialect>` | a list of all matches |
| `<source>[~? <dialect>'<expression>']`          | `QUERY ONE <expression> IN <source> USING <dialect>`            | the first match, or `NONE` |

{{< code lang="fql" >}}
LET products = doc[~ css`.product`]
LET title = doc[~? css`.product .title`]
{{</ code >}}

These are equivalent to:

{{< code lang="fql" >}}
LET products = QUERY `.product` IN doc
LET title = QUERY ONE `.product .title` IN doc
{{</ code >}}

The shortcut provides only the simple query form: a query expression, a source value, and the source value’s dialect. It does not support `WITH`, or `OPTIONS`. Use the full `QUERY` form when you need to pass parameters, configure execution options, or make the query behavior explicit.

`QUERY COUNT` and `QUERY EXISTS` do not have shortcut forms at this time.

## Parameters and options

`WITH` supplies parameters to the query, and `OPTIONS` carries execution settings. Both are evaluated once.

{{< code lang="fql" >}}
QUERY `SELECT name, price FROM products WHERE category = $c` IN db
WITH { c: "laptops" }
OPTIONS { timeout: 5000 }
{{</ code >}}

For host values with multiple dialects, `USING` can still be used together with `WITH` and `OPTIONS`:

{{< code lang="fql" >}}
QUERY `SELECT name, price FROM products WHERE category = $c` IN db USING sql
WITH { c: "laptops" }
OPTIONS { timeout: 5000 }
{{</ code >}}

## A host capability

`QUERY` only works when the source value supports querying — it must be a **queryable** value. Queryable values may include HTML documents, database connections, API clients, JSON documents, or other values provided by a module or the host application.

The source value decides which dialects are available. It may also define a default dialect. If the source does not support querying, the requested dialect is unavailable, or no dialect is provided and the source has no default, the query fails at runtime.

You can recover from such failures with a recovery clause.

{{< code lang="fql" >}}
LET rows = QUERY `.row` IN doc ON ERROR RETURN []
{{</ code >}}

For more on values that expose querying and other behaviors, see [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Next steps

{{< docs-related tiles="web-extraction,language-types-capabilities,language-control-flow" >}}