---
title: "Language Overview"
sidebarTitle: "Overview"
weight: 20
draft: false
description: "The core ideas behind FQL: how scripts produce values, how expressions compose, how collections transform, and how the runtime shapes what a script can do."
aliases:
    - /docs/fql/introduction/
---

# Language Overview

FQL is a small, expression-oriented language for extracting, transforming, and returning structured data. It's designed to be embedded in tools and runtimes that determine what a script can access.

## Scripts produce values

Every useful FQL script ends by producing a result. The most direct way is with `RETURN`:

{{< editor lang="fql" >}}
LET price = 19.99
LET quantity = 3

RETURN {
    price: price,
    quantity: quantity,
    total: price * quantity
}
{{< /editor >}}

A script can also end with a top-level `FOR` expression, in which case the result is the array that `FOR` produces. Either way, FQL scripts communicate through returned values rather than printed output. That makes the same script usable across different surfaces: a CLI can serialize the result as JSON, an embedded application can receive it as a Go value, and a test runner can assert against it directly. The script does not need to know who is consuming the result.

## Expressions are the main building block

FQL is expression-oriented: most pieces of logic produce a value. Arithmetic, object, and array literals, conditionals, function calls, queries, and waiting constructs all evaluate to values that can be assigned, returned, or composed with other expressions.

{{< editor lang="fql" >}}
LET score = 87
RETURN score >= 80 ? "passed" : "failed"
{{< /editor >}}

`LET` is how you name an intermediate value. A typical script uses `LET` to break logic into readable steps - get something, transform it, return the result - rather than nesting everything into one large expression.

{{< editor lang="fql" >}}
LET user = { name: "Ada", roles: ["admin", "editor"] }
LET isAdmin = CONTAINS(user.roles, "admin")

RETURN {
    name: user.name,
    isAdmin: isAdmin
}
{{< /editor >}}

## Collections are transformed with FOR

`FOR` is not a general-purpose loop - it is a data-shaping construct. It iterates over a collection and returns a new collection, optionally filtering and transforming along the way.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

FOR user IN users
    FILTER user.active
    RETURN user.name
{{< /editor >}}

This script starts with an array of users and returns only the names of active ones. The shape of the data changes; the logic remains declarative. `FOR` expressions can appear inline too, assigned via `LET` and composed with the rest of the script:

{{< editor lang="fql" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org/scenarios/ecommerce/products/")

FOR item IN doc[~ css`.product-card`]
    FILTER item.attributes["data-in-stock"] == "true"
    RETURN {
        title: item[~? css`.product-title`].textContent,
        url: item[~? css`a`].attributes.href
    }
{{< /editor >}}

Even when the source is messy, the final result can be structured and clean.

## Queries operate on capable values

FQL's query syntax is not tied to one data type or one library. Instead, querying is capability-based: a value can support one or more query dialects, and FQL can query that value using whichever dialect is appropriate.

An HTML document might support both CSS and XPath:

{{< editor lang="fql" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
LET links = doc[~ css`a[href]`]

RETURN links
{{< /editor >}}

The long form of a query expression makes the structure explicit, and supports passing query-specific data or options:

{{< editor lang="fql" >}}
LET db = DB::SQLITE::OPEN({ memory: true })

LET create = QUERY ONE `
  CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
  )
` IN db USING sql_exec

LET insert = QUERY ONE `
  INSERT INTO users(name)
  VALUES (?)
`   IN db 
    USING sql_exec 
    WITH {
        params: ["Ada"]
    }
    OPTIONS { transactionLevel: "read-uncommitted" }

RETURN db[~ sql`SELECT id, name FROM users`]
{{< /editor >}}

The meaning of `WITH` and `OPTIONS` is defined by the value being queried and the selected dialect. The language provides a uniform syntax for query expressions; the runtime value provides the actual query behavior. This means FQL does not need built-in syntax for every possible source - HTML documents, databases, browser elements, and custom objects can all expose query capabilities while reusing the same FQL syntax.

## Waiting is explicit

Dynamic workflows often involve timing: a page may not have finished loading, an element may appear only after JavaScript runs, or a value may change in response to an event. FQL includes waiting constructs for these cases, and they are first-class expressions rather than library utilities.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org/scenarios/network/delayed-requests/", { driver: "cdp" })
RETURN WAITFOR VALUE doc[~ css`.network-result-card p`]
    TIMEOUT 5s
    EVERY 250ms
    ON TIMEOUT RETURN false
{{< /editor >}}

This expression polls every 250ms, gives up after 5s, and returns `false` if the timeout is reached. A script that waits forever is very different from a script that waits five seconds and returns a fallback - FQL makes that difference visible in the source rather than hiding it in configuration or library code.

## Modules provide runtime behavior

FQL is intentionally small. The core language defines the syntax and execution model; modules and the host runtime supply the capabilities.

A module can provide namespaced functions:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN YAML::DECODE(`
name: Ada
roles:
- admin
- editor
  `)
  {{< /editor >}}

A module can also provide value types with query support - the HTML module, for example, exposes `DOCUMENT`, `PARSE` functions, and the `css`, `xpath` query dialects. Other modules can provide integrations with file formats, external APIs, browser runtimes, databases, and custom application objects. The language stays consistent regardless of which modules are present; the host environment determines what is available.

This means the same script might behave differently in a CLI context than in a browser automation runtime, not because the language changes, but because the set of registered modules and host capabilities differs. This is similar to how a SQL query depends on the database engine it runs against, or how a JavaScript file behaves differently in a browser versus Node.js.

## Where to go next

{{< docs-related tiles="language-structure,language-types,language-variables,language-expressions" >}}
