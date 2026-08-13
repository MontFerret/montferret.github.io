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

FQL is a small, declarative-first, expression-oriented language for data automation. It is designed to be embedded in tools and applications whose host runtime determines what a script can access.

Keywords and registered host-function calls are case-insensitive, while documentation uses lowercase as the canonical spelling. Registered namespace segments and host-function names have one lowercase qualified identity; variables, user-defined functions, aliases, object properties, query dialects, strings, and comments retain their existing casing rules.

## Scripts produce values

Use `return` when an FQL script needs to produce a result:

{{< editor lang="fql" >}}
let price = 19.99
let quantity = 3

return {
    price: price,
    quantity: quantity,
    total: price * quantity
}
{{< /editor >}}

To return a loop's array, write `return for ...`. A script that reaches the end without `return` completes with `none`; standalone loops execute but discard their arrays. FQL scripts communicate through returned values rather than printed output. That makes the same script usable across different surfaces: a CLI can serialize the result as JSON, an embedded application can receive it as a Go value, and a test runner can assert against it directly. The script does not need to know who is consuming the result.

## Expressions are the main building block

FQL is expression-oriented: most pieces of logic produce a value. Arithmetic, object, and array literals, conditionals, function calls, queries, and waiting constructs all evaluate to values.

That means expressions can usually be assigned to variables, returned from scripts, passed to functions, or nested inside larger expressions.

{{< editor lang="fql" >}}
let score = 87
return score >= 80 ? "passed" : "failed"
{{< /editor >}}

`let` is how you name an intermediate value. A typical script uses `let` to break logic into readable steps - get something, transform it, return the result - rather than nesting everything into one large expression.

{{< editor lang="fql" >}}
let user = { name: "Ada", roles: ["admin", "editor"] }
let isAdmin = contains(user.roles, "admin")

return {
    name: user.name,
    isAdmin: isAdmin
}
{{< /editor >}}

## Declarative core and constrained state

FQL's default model is to describe values and transformations. Domain workflows compose through `match`, `for`, `filter`, `query`, and `waitfor`: branching, shaping collections, delegating to capable host values, and coordinating with external state all remain visible in the expression flow.

When a workflow genuinely needs state, `var` provides a mutable binding and `while` provides condition-driven iteration. These constructs are deliberately constrained; immutable `let` bindings and value-producing expressions remain the normal style.

## Collections are transformed with for

`for` is not a general-purpose loop - it is a data-shaping construct. It iterates over a collection and returns a new collection, optionally filtering and transforming along the way.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

return for user in users {
    filter user.active
    return user.name
}
{{< /editor >}}

This script starts with an array of users and returns only the names of active ones. The shape of the data changes; the logic remains declarative. `for` expressions can appear inline too, assigned via `let` and composed with the rest of the script:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/")

return for item in page[~ css`.product-card`] {
    filter item.attributes["data-in-stock"] == "true"
    return {
        title: item[~? css`.product-title`].textContent,
        url: item[~? css`a`].attributes.href
    }
}
{{< /editor >}}

Even when the source is messy, the final result can be structured and clean.

## Queries operate on capable values

FQL's query syntax is not tied to one data type or one library. Instead, querying is capability-based: a value can support one or more query dialects, and FQL can query that value using whichever dialect is appropriate.

An HTML object might support both CSS and XPath:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
let links = page[~ css`a[href]`]

return links
{{< /editor >}}

The long form of a query expression makes the structure explicit, and supports passing query-specific data or options:

{{< editor lang="fql" >}}
let db = db::sqlite::open({ memory: true })

let create = query one `
  CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
  )
` in db using sql_exec

let insert = query one `
  INSERT INTO users(name)
  VALUES (?)
` in db
    using sql_exec
    with {
        params: ["Ada"]
    }
    options { transactionLevel: "read-uncommitted" }

return db[~ sql`SELECT id, name FROM users`]
{{< /editor >}}

The meaning of `with` and `options` is defined by the value being queried and the selected dialect. The language provides a uniform syntax for query expressions; the runtime value provides the actual query behavior. This means FQL does not need built-in syntax for every possible source - HTML documents, databases, browser elements, and custom objects can all expose query capabilities while reusing the same FQL syntax.

## Waiting is explicit

Dynamic workflows often involve timing: a page may not have finished loading, an element may appear only after JavaScript runs, or a value may change in response to an event. FQL includes waiting constructs for these cases, and they are first-class expressions rather than library utilities.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/network/delayed-requests/", { driver: "cdp" })

return waitfor value page[~ css`.network-result-card p`]
    when length(.) > 0
    timeout 5s
    on timeout return false
{{< /editor >}}

This expression polls every 250ms, gives up after 5s, and returns `false` if the timeout is reached. A script that waits forever is very different from a script that waits five seconds and returns a fallback - FQL makes that difference visible in the source rather than hiding it in configuration or library code.

## Modules provide runtime behavior

FQL is intentionally small. The core language defines the syntax and execution model; modules and the host runtime supply the capabilities.

A module can provide namespaced functions:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
return yaml::decode(`
name: Ada
roles:
- admin
- editor
  `)
  {{< /editor >}}

A module can also provide value types with query support - the HTML module, for example, exposes `web::html::open`, `web::html::parse`, and the `css`, `xpath` query dialects. Other modules can provide integrations with file formats, external APIs, browser runtimes, databases, and custom application objects. The language stays consistent regardless of which modules are present; the host environment determines what is available.

This means the same script might behave differently in a CLI context than in a browser automation runtime, not because the language changes, but because the set of registered modules and host capabilities differs. This is similar to how a SQL query depends on the database engine it runs against, or how a JavaScript file behaves differently in a browser versus Node.js.

## Next steps

{{< docs-related tiles="language-structure,language-types,language-variables,language-expressions" >}}
