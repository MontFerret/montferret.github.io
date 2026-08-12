---
title: "Query Expressions"
sidebarTitle: "Query"
weight: 40
draft: false
description: "Delegate a query to a host value using the query ... in ... expression."
---

# Query Expressions

A `query` expression delegates work to a value that knows how to run queries. Instead of expressing the logic in FQL, you provide a query payload and pass it to a host value, such as an HTML document, database connection, API client, or JSON document.

{{< code lang="fql" >}}
query `.product .title` in doc
{{</ code >}}

This reads as: run the query `` `.product .title` `` against `doc`.

Some host values support more than one query dialect. In that case, the value may provide a default dialect, so `using` is not always required. For example, an HTML document may default to CSS selectors.

When you need to choose a specific dialect, add `using`:

{{< code lang="fql" >}}
query `//article/h1` in doc using xpath
{{</ code >}}

This reads as: run the query `` `//article/h1` `` against `doc`, using the `xpath` dialect.

The query expression is often written with backticks so selectors and other syntax do not need escaping. It can also be a regular string, variable, bind parameter, member expression, indexed expression, or function call. Compound expressions must be wrapped in parentheses so the boundary before `in` is explicit.

{{< code lang="fql" >}}
query config.selector in page
query (".item-" + category) in page
{{</ code >}}

## Anatomy

{{< code lang="fql" >}}
query [modifier] <expression> in <source> [using <dialect>] [with <params>] [options <options>]
{{</ code >}}

- **expression** — the query to run. Literals and atomic expressions can be used directly; compound expressions such as concatenation, comparisons, logical operations, ternaries, or `in` predicates must be parenthesized.
- **source** — an expression that yields the value the query runs against. The value must support querying.
- **using** — optionally selects a query dialect, such as `css`, `xpath`, or `sql`. If omitted, the source value chooses its default dialect when one is available.
- **with** — an optional value passed to the query as parameters.
- **options** — an optional value carrying execution settings, such as a timeout.

If `using` is omitted and the source value does not provide a default dialect, the query fails at runtime.

## Result modifiers

By default, a `query` returns a list of every match. A modifier after `query` changes the shape of the result.

| Form | Returns |
| --- | --- |
| `query ...` | a list of all matches |
| `query one ...` | the first match, or `none` if there are none |
| `query count ...` | the number of matches |
| `query exists ...` | `true` if there is at least one match, otherwise `false` |

{{< code lang="fql" >}}
let total = query count `.item` in doc
let hasNext = query exists `.pagination .next` in doc
let title = query one `.product .title` in doc

return { total, hasNext, title }
{{</ code >}}

## Shortcut syntax

The two most common query forms have shorter equivalents:

| Shortcut                             | Equivalent form                                  | Returns |
|--------------------------------------|--------------------------------------------------| --- |
| `<source>[~ <dialect>'<expression>']` | `query <expression> in <source> using <dialect>` | a list of all matches |
| `<source>[~? <dialect>'<expression>']`          | `query one <expression> in <source> using <dialect>`            | the first match, or `none` |

{{< code lang="fql" >}}
let products = doc[~ css`.product`]
let title = doc[~? css`.product .title`]
{{</ code >}}

These are equivalent to:

{{< code lang="fql" >}}
let products = query `.product` in doc
let title = query one `.product .title` in doc
{{</ code >}}

The shortcut provides only the simple query form: a query expression, a source value, and the source value’s dialect. It does not support `with`, or `options`. Use the full `query` form when you need to pass parameters, configure execution options, or make the query behavior explicit.

`query count` and `query exists` do not have shortcut forms at this time.

## Parameters and options

`with` supplies parameters to the query, and `options` carries execution settings. Both are evaluated once.

{{< code lang="fql" >}}
query `SELECT name, price FROM products WHERE category = $c` in db
with { c: "laptops" }
options { timeout: 5000 }
{{</ code >}}

For host values with multiple dialects, `using` can still be used together with `with` and `options`:

{{< code lang="fql" >}}
query `SELECT name, price FROM products WHERE category = $c` in db using sql
with { c: "laptops" }
options { timeout: 5000 }
{{</ code >}}

## A host capability

`query` only works when the source value supports querying — it must be a **queryable** value. Queryable values may include HTML documents, database connections, API clients, JSON documents, or other values provided by a module or the host application.

The source value decides which dialects are available. It may also define a default dialect. If the source does not support querying, the requested dialect is unavailable, or no dialect is provided and the source has no default, the query fails at runtime.

You can recover from such failures with a recovery clause.

{{< code lang="fql" >}}
let rows = query `.row` in doc on error return []
{{</ code >}}

For more on values that expose querying and other behaviors, see [Value Capabilities]({{% ref "../types/capabilities" %}}) and [Host Values]({{% ref "../types/host" %}}).

## Next steps

{{< docs-related tiles="web-extraction,language-types-capabilities,language-control-flow" >}}
