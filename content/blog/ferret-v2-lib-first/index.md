---
title: "Why Ferret Was Built as a Library First"
subtitle: "Build the library first. Everything else follows."
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-08-05"
---

When people first discover Ferret, they usually encounter the command-line interface.

They write a query, run it against a website or an API, inspect the results, and move on, and for many use cases, that’s all they need.

But in reality, the CLI was never the center of the project. Ferret was built around a much simpler idea: build the library first, then build applications around it.

One engineering principle that has influenced almost every project I’ve built over the past fifteen years is this: whenever possible, I start by building a library.

I do this because libraries naturally encourage clear boundaries and reusable abstractions. They typically solve one problem well while leaving the surrounding architecture to the application.

Applications come and go, interfaces evolve, and new deployment models appear, but a good library remains useful regardless of how it’s used.

That philosophy has always felt very much in the spirit of Unix: build focused components that compose well instead of monolithic tools that try to own everything.

Ferret follows exactly the same principle. The runtime is the product, while the CLI is simply one application built on top of it, no different from a backend service, a worker, a desktop application, or a test runner.

## We already trust specialized engines

Let's think about databases for a moment.

Very few applications implement their own SQL parser, query planner, optimizer, and execution engine. We use a database because querying relational data is a solved problem.

Our application simply asks a question.

```sql
SELECT *
FROM inventory
WHERE category = 'Books'
  AND available = true;
```

And the database decides how to answer it.

Yet when it comes to extracting data from websites, APIs, or documents, many applications still embed that entire execution engine directly into their own codebase.

They contain hundreds — or sometimes thousands — of lines of code responsible for:

* making HTTP requests
* traversing HTML
* parsing JSON
* handling pagination
* retrying failed requests
* filtering and transforming data
* coordinating browser automation

Eventually, that extraction logic becomes tightly intertwined with the application’s business logic, so even simple changes require touching application code.

## A query engine instead of glue code

Ferret takes a different approach.

Instead of embedding all that extraction logic into your application, you embed a query engine. Your application asks the questions, and Ferret figures out how to answer them.

By keeping business logic in your application and extraction logic in queries, the two can evolve independently. That’s a surprisingly small architectural change with a much bigger impact than it might seem.

Instead of hundreds of lines of imperative code describing how to retrieve data, your application executes a query describing what data it needs.

## Your application remains in control

One consequence of building Ferret as a library is that it never tries to become your framework. 

Your application continues to own its architecture, networking, configuration, logging, authentication, and everything else that makes it your application. 
Ferret simply does the one thing it was designed to do: evaluate queries.

In practice, that means your application decides:

* which modules are available
* what network access is permitted
* whether browser automation is enabled
* what files can be accessed
* which custom functions are exposed
* which policies queries must follow

The runtime operates entirely within the boundaries defined by the host.

                 Your Application
        ┌─────────────────────────┐
        │ Business Logic          │
        │ HTTP Server             │
        │ Authentication          │
        │ Logging                 │
        │                         │
        │   ┌─────────────────┐   │
        │   │ Ferret Runtime  │   │
        │   └─────────────────┘   │
        └────────────┬────────────┘
                     │
        Browser • APIs • DB • PDFs

## The CLI is just another application

The runtime operates entirely within the boundaries defined by the host.

This is also why I think of the CLI as just another application rather than _the_ application. It creates an engine, registers the modules it needs, executes a query, and prints the result — exactly what a backend service, a desktop application, a testing framework, or an internal automation platform might do.

{{< code lang="go" >}}
engine := ferret.New()

result, err := engine.Run(ctx, query, ferret.Params{
    "url": url,
})

return result
{{</ code >}}

The CLI exists because it’s a convenient way to execute and experiment with queries, not because Ferret revolves around it. 

In many ways, it’s simply a reference application that demonstrates how the runtime can be embedded.

## Extending the language with your own domain

Because the runtime lives inside your application, it can naturally expose concepts that are unique to your domain.

Suppose you’re building an inventory system. 
Your application could register an `INVENTORY` namespace that exposes products, warehouses, suppliers, or purchase orders directly to queries.

{{< code lang="fql" >}}
LET products = INVENTORY::SEARCH({
    category: "Books",
    available: true
})

RETURN products
{{</ code >}}

`INVENTORY` isn’t built into Ferret at all; it’s simply a namespace provided by your application. 
To someone writing queries, though, that distinction is almost invisible - it feels like a natural part of the language. 

The same approach works equally well for CRMs, ERPs, internal APIs, cloud platforms, AI services, or any other domain your application exposes.

## One runtime, many kinds of data

Databases specialize in querying relational data. Ferret applies the same idea to heterogeneous data, allowing a single query to work with web pages, REST APIs, SQL databases, spreadsheets, PDF documents, AI models, browser sessions, or objects provided by your own application.

From the runtime’s perspective, they’re simply values that expose capabilities. Queries don’t need to care where those values came from; they simply describe the transformations to perform.

## Separating responsibilities

One of the nicest side effects of this design is the clear separation of responsibilities. 

Business logic remains in your application, while extraction logic lives in Ferret, allowing each to evolve independently. 
The application no longer needs to know about CSS selectors, XPath expressions, pagination, or browser events, just as the query doesn’t need to care about dependency injection, HTTP routing, authentication middleware, or deployment infrastructure. 

Each focuses on solving its own problem, making the overall system easier to understand, easier to test, and easier to evolve over time.

## Looking back

Looking back, I realize I was never trying to build another scraping tool. What I really wanted was a reusable query engine that applications could embed in the same way they embed a database, a JavaScript engine, or a template engine.

The command-line interface came naturally because it’s a convenient way to experiment with queries, automate repetitive tasks, and integrate with scripts. But it’s only one consumer of the runtime, not the foundation it was built on.

That idea has shaped Ferret from the very beginning, and it continues to influence how I think about new features and the overall architecture of the project.