---
title: "Language Overview"
sidebarTitle: "Overview"
weight: 10
draft: false
description: "Learn the core ideas behind FQL: structured values, expressions, queries, control flow, modules, and host capabilities."
aliases:
  - /docs/fql/introduction/
---

# Language Overview

FQL is the language used by Ferret to describe data extraction, transformation, automation, and runtime-driven workflows.

At a high level, FQL scripts take input from somewhere, shape it into useful values, and return structured results.

That input might be simple in-memory data:
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    roles: ["admin", "editor"]
}

RETURN {
    name: user.name,
    isAdmin: CONTAINS(user.roles, "admin")
}
{{< /editor >}}


Or it might come from a document, a web page, a browser session, a file, an API, or a value provided by the host application.

FQL is not meant to be a full general-purpose application language like Go, JavaScript, Python, or Java. Instead, it is designed as a small, embeddable language for describing focused pieces of data-oriented logic.

You can use FQL to:

- extract data from structured or semi-structured sources
- query documents and runtime values
- transform arrays and objects
- wait for values, elements, or events
- automate browser-driven workflows
- define reusable scripts for CLIs, test runners, workers, or embedded applications
- run a user-defined logic for a configuration-driven application

The important idea is that FQL scripts are usually not large programs. They are usually compact units of logic that describe what data you want, how to shape it, and what result should be returned.

## A small language around structured values

FQL is built around values. A script can work with simple primitive values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN 1 + 1
{{< /editor >}}

It can work with arrays:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET names = ["Ada", "Grace", "Linus"]
RETURN names
{{< /editor >}}

It can work with objects:
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    active: true
} 
RETURN user.name
{{< /editor >}}

And it can return structured data:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN {
    ok: true,
    count: 3,
    items: ["one", "two", "three"] 
}
{{< /editor >}}

This is an important part of the language model.

FQL scripts communicate through returned values rather than printed output, which allows the same script result to be serialized as JSON, consumed directly by a host application, displayed by a CLI, inspected by a test runner, or passed to another runtime component.
That makes FQL useful both as a command-line language and as an embedded language.

For example, a CLI user might run a script and receive JSON output. An application embedding Ferret might execute the same script and receive a runtime value directly in Go.

The script does not need to know who is consuming the result.

## Scripts return results

Every useful FQL script eventually produces a result.

The most direct way to do that is with `RETURN`:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN "Hello from Ferret"
{{< /editor >}}

`RETURN` can return any value:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN {
    language: "FQL",
    runtime: "Ferret",
    embedded: true 
}
{{< /editor >}}

A script can also end with a top-level `FOR` expression:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true } 
]

FOR user IN users 
    FILTER user.active     
    RETURN user.name
{{< /editor >}}

In this case, the script result is the array produced by the `FOR` expression.

Which means an FQL script normally ends in one of two ways:

* with a `RETURN` statement that returns a single value
* with a top-level `FOR` expression that returns a collection

This keeps scripts result-oriented.

An FQL script usually builds a result by naming intermediate values, applying transformations, and returning the final structure.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET price = 19.99
LET quantity = 3
RETURN {
    price: price,
    quantity: quantity,
    total: price * quantity 
}
{{< /editor >}}

This may look simple, but it has a useful consequence: scripts are easier to embed, test, and compose.

This result-oriented model makes FQL easy to reuse across different environments.
Embedded applications can consume returned values directly, test runners can assert against them, and command-line tools can serialize them for display.

The same language model works in all of those cases.

## Expressions are the main building block

FQL is expression-oriented.

That means most pieces of logic produce a value.

Arithmetic expressions produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN 10 * 2
{{< /editor >}}

Object and array expressions produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN {
    tags: ["docs", "language", "overview"] 
}
{{< /editor >}}

Conditional expressions produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET score = 87  
RETURN score >= 80 ? "passed" : "failed"
{{< /editor >}}

Function calls produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET roles = ["admin", "editor"]  
RETURN CONTAINS(roles, "admin")
{{< /editor >}}

Query expressions produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
RETURN doc[~? css`.product-card`]
{{< /editor >}}

Waiting expressions produce values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/delayed-requests/", { driver: "cdp" })
RETURN WAITFOR EXISTS doc[~ css`.network-result-card`]
    TIMEOUT 5s
    EVERY 250ms
    ON TIMEOUT RETURN false
{{< /editor >}}

This is one of the core ideas behind FQL.

Although FQL includes statements, the language is primarily value-oriented. Scripts are usually written by defining intermediate values, transforming data step by step, and returning the final structured result.

That keeps FQL compact and makes scripts easier to reason about.

## LET gives names to values

Use `LET` to bind a value to a name:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET name = "Ada"  
RETURN name
{{< /editor >}}

A `LET` binding is useful when a value needs a name, when an expression becomes too large, or when you want to make a script easier to read.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET user = {
    name: "Ada",
    roles: ["admin", "editor"] }  
LET isAdmin = CONTAINS(user.roles, "admin")  

RETURN {
    name: user.name,
    isAdmin: isAdmin 
}
{{< /editor >}}

In most scripts, `LET` is the main way to break logic into readable steps.

This is especially helpful in data extraction workflows, where a script often has a sequence like:

1. get a document or page
2. query some values from it
3. transform those values
4. return a clean result

For example:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
LET cards = doc[~ css`.product-card`]  
LET names = (
    FOR card IN cards 
        RETURN card[~? css`.product-title`].textContent
)  

RETURN names
{{< /editor >}}

The exact query syntax is covered later in the docs. For now, the important point is the shape of the script: bind intermediate values, transform them, return the result.

## Data flow is explicit

FQL is designed to make data flow visible.

A typical script reads from top to bottom:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true } 
]  

FOR user IN users 
    FILTER user.active     
    RETURN user.name
{{< /editor >}}

This script starts with an array of users, filters only active users, and returns their names.

The loop does not exist just to repeat commands. It describes a transformation from one collection into another collection.

That is the general idea behind FOR in FQL: it is a data-shaping construct.

Later sections cover loop clauses in more detail, including filtering, sorting, collecting, and returning transformed values. The overview-level idea is enough for now:

> FOR expressions are how FQL transforms collections.

This style is useful for extraction scripts because raw data is rarely in the exact shape you need.

You might extract a list of elements from a page, filter out irrelevant ones, map each element into an object, and return a clean array.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
FOR item IN doc[~ css`.product-card`]
    FILTER item.attributes["data-in-stock"] == "true"  
    RETURN {
        title: item[~? css`.product-title`].textContent,         
        url: item[~?  css`a`].attributes.href
}
{{< /editor >}}

Even when the source is messy, the final result can be structured.

## Queries are capability-based

Querying is one of the most important ideas in FQL.

In many languages, query syntax is tied to one specific data type or one specific library. FQL takes a different approach: querying is based on capabilities.

A value can support one or more query dialects. If it does, FQL can query that value using the appropriate dialect.

For example, an HTML document may support CSS queries:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
LET links = doc[~ css`a[href]`]
RETURN links
{{< /editor >}}

or XPATH queries:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
LET links = doc[~ xpath`//a[@href]`]
RETURN links
{{< /editor >}}

In both cases, FQL queries the same target value, doc, but uses a different dialect to interpret the query expression.
The target value is the same, but the dialect changes how the query text is interpreted. The `css` dialect treats `a[href]` as a CSS selector, while the `xpath` dialect treats `//a[@href]` as an XPath expression.

The long form is more explicit:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/ecommerce/products/")
RETURN QUERY `a[href]` IN doc USING css
{{< /editor >}}

The long form is the full query expression. Use it when the query needs to stand on its own, when readability matters, or when you need to pass additional data to the underlying query implementation.

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}

RETURN QUERY `SELECT * FROM users WHERE id = :id` 
    IN db 
    USING sql 
    WITH {
        params: {
            id: @id
        }
    }
{{< /code >}}

A query may also accept implementation-specific options:

{{< code lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN QUERY `SELECT * FROM users WHERE id = :id` 
    IN db
    USING sql
    WITH {
        params: {
            id: @id
        }
    }
    OPTIONS {
        transactionLevel: "read-uncommitted"
    }
{{< /code >}}

The meaning of `WITH` and `OPTIONS` is defined by the value being queried and the selected dialect. FQL provides a common syntax for query expressions; the runtime value provides the actual query behavior.

This is the core of Ferret’s capability-based query model. The language does not need built-in syntax for every possible source or query system. HTML documents, browser elements, databases, custom objects, or module-defined values can expose their own query capabilities while still using the same FQL syntax.

The available dialects and options depend on the runtime and modules in use, but the language model stays the same: select a value, choose a dialect, provide a query, and optionally pass query-specific arguments or options.

This keeps FQL small while allowing the ecosystem to grow.

## Waiting is part of the language

Many extraction and automation workflows are not instantaneous.

Static values can usually be queried immediately, but dynamic workflows often involve timing. Pages may need to finish loading, elements may appear after JavaScript runs, values may change in response to events, and browser workflows may need to wait for specific conditions before continuing.

FQL includes waiting constructs for these cases.

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/delayed-requests/", { driver: "cdp" })
RETURN WAITFOR VALUE doc[~? css`.network-result-card p`]
    TIMEOUT 5s
    EVERY 250ms
    ON TIMEOUT RETURN false
{{< /editor >}}

This expression waits for `.network-result-card` to exist. It checks every 250ms, gives up after 5s, and returns false if the timeout is reached.

The exact syntax is covered later, but the design goal is worth understanding early:

> Waiting is not an afterthought hidden in library code. Waiting is part of the script.

That matters because timeout behavior is part of the logic.

A script that waits forever is very different from a script that waits for five seconds and returns a fallback value.

By making waiting explicit, FQL makes dynamic workflows easier to read and safer to run.

A script can say:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/delayed-requests/", { driver: "cdp" })

RETURN WAITFOR EXISTS doc[~ css`.foobar`]
    TIMEOUT 1s
    EVERY 500ms
    ON TIMEOUT FAIL
{{< /editor >}}

Or it can choose a softer failure mode:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev/scenarios/network/delayed-requests/", { driver: "cdp" })
RETURN WAITFOR EXISTS doc[~ css`.foobar`]
    TIMEOUT 1s
    EVERY 500ms
    ON TIMEOUT RETURN false
{{< /editor >}}

Those are different policies, and the script makes that difference visible.

## Control flow stays data-oriented

FQL includes control flow, but its goal is not to turn scripts into large imperative programs.

Control flow exists to keep data logic readable when expressions become too complex.

For small choices, a ternary expression is often enough:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET count = 5
RETURN count > 0 ? "found" : "empty"
{{< /editor >}}

For larger branching logic, `MATCH` can make intent clearer:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET status = "pending"
RETURN MATCH status (
    "ready" => "continue",     
    "pending" => "wait",     
    _ => "stop"
)
{{< /editor >}}

The goal is not to make FQL feel like a traditional statement-heavy language.

The goal is to preserve the declarative, result-oriented feel of scripts while still providing enough structure for real-world cases.

This is especially important in extraction and automation workflows, where branching is common:

- a page may have one of several layouts
- a value may be missing
- a response may represent success or failure
- a workflow may need different behavior depending on state

FQL should let you express those cases clearly without forcing everything into nested ternary expressions or host-side code.

## Modules extend the runtime

FQL itself is intentionally small.

The core language provides the syntax and execution model, but many useful capabilities come from modules and the host runtime.

A module can provide functions:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
RETURN YAML::DECODE(`
name: Ada
roles:
  - admin
  - editor
`)
{{< /editor >}}

A module can provide runtime values:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = PARSE("<h1>Hello</h1>")
RETURN doc
{{< /editor >}}

> HTML module is registered without a namespace prefix for legacy reasons, but it is still a module. The `DOCUMENT` function and the `css` query dialect are provided by that module.

A module can provide query behavior:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")

RETURN WEB::ARTICLE::EXTRACT(doc)
{{< /editor >}}

Modules can provide integrations with external systems, file formats, browser runtimes, APIs, custom application objects, and other domain-specific capabilities.
This separation keeps the core language focused while allowing runtimes and modules to add behavior for specific use cases.
Ferret does not need to put every possible feature into the language core. Instead, the language can stay focused while modules add domain-specific capabilities.

For example, HTML support does not need to be built into the core language. It can be provided by a module, just as browser automation can be provided by a runtime, host applications can expose their own functions and values, and test runners can add test-specific parameters and helpers.

The language stays consistent from the script’s point of view, while the host environment determines which capabilities are available.

## Host capabilities shape what a script can do

Because FQL is embeddable, the available capabilities are defined by the environment in which a script runs.

A CLI environment may provide filesystem access, standard modules, and command-line parameters. A browser-oriented runtime may expose document and element values. A test runner may add fixture URLs, assertions, or static test assets. An application embedding Ferret may register its own domain objects and functions.

The language remains consistent across these environments, but each host controls which capabilities are available to the script.

It allows Ferret to be used as:

- a command-line scripting tool
- an embedded expression and automation engine
- a data extraction runtime
- a browser automation layer
- a test runner language
- a configuration-driven workflow language

FQL source code is portable, but it is always evaluated in a specific runtime environment. That environment determines which modules, functions, value types, query dialects, and host capabilities are available to the script.

In other words, a script is defined by both its source code and the capabilities of the runtime that executes it. This is similar to how a SQL query depends on the database engine it runs against, or how JavaScript code depends on whether it runs in a browser, Node.js, or another host environment.

## Error and timeout behavior can be part of expressions

Real-world extraction workflows often involve partial or uncertain conditions: a selector may not match, a page may load slowly, a value may be missing, a network request may fail, or a runtime capability may reject an operation. 
FQL is designed to keep the handling of those cases close to the operation itself, so the script makes its failure and fallback behavior explicit.

{{< editor lang="fql" height="128px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN QUERY ONE "#price" IN doc USING css
    ON ERROR RETURN NONE
{{</ editor >}}

{{< editor lang="fql" height="120px" apiVersion="2" >}}
LET parsed = TO_FLOAT("not a number") ON ERROR RETURN 42
RETURN parsed
{{</ editor >}}

Timeout behavior can be expressed in a similar way where supported:

{{< editor lang="fql" height="146px" apiVersion="2" >}}
LET doc = DOCUMENT("https://mockery.montferret.dev")
RETURN WAITFOR VALUE doc[~ css`.loaded`]
    TIMEOUT 5s
    EVERY 250ms
    ON TIMEOUT RETURN NONE
{{</ editor >}}

The script author decides how each failure mode should be treated. 
In automation and extraction workflows, a timeout or missing value may be an error, or it may be an expected outcome that should produce `NONE`, `false`, an empty array, or another fallback value. 
FQL keeps that behavior explicit and close to the operation that needs it.

## FQL can be used from different tools

FQL is the language, but Ferret is an ecosystem.

You can run FQL through the CLI:

{{< terminal command="true" >}}
ferret run script.fql
{{< /terminal >}}

You can pass a short expression directly:

{{< terminal command="true" >}}
ferret run -e 'RETURN 1 + 1'
{{< /terminal >}}

FQL is not limited to the command line. Ferret can be embedded into Go applications, where scripts run as part of application code and return values directly to the host. Tools such as Lab can use FQL scripts as tests against static pages, browser runtimes, APIs, or other controlled scenarios. Specialized runtimes and workers can execute the same language in environments designed for automation, scheduling, or distributed execution.

The key idea is that FQL is portable across interfaces. The same language can be used interactively, in automation, in CI, in tests, or inside another application.

## What to learn next

TODO
