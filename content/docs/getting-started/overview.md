---
title: "Overview"
sidebarTitle: "Overview"
weight: 20
draft: false
description: "What Ferret is and how its main pieces fit together."
aliases:
    - /docs/introduction/
---

# Overview

Ferret is a declarative-first, expression-oriented embedded language and runtime for data automation. FQL scripts query and transform ordinary data, application-defined host values, documents, APIs, and browser-backed resources while the embedding application controls the available capabilities.

Ferret is especially useful when data logic should be portable, reviewable, and change independently from the host application. At a high level, a script describes:

- where data comes from
- how to query it
- how to wait for it
- how to transform it
- how to return clean structured output

Ferret can be used as a command-line tool, embedded into Go applications, extended with modules, or used as a small domain-specific language for configuration-driven systems.

## What Ferret is
Ferret is built around FQL, a declarative-first language designed for data extraction, transformation, and automation workflows. Its core is expression-oriented: `match`, `for`, `filter`, `query`, and `waitfor` compose domain logic and produce values. `var` and `while` add constrained mutable state for workflows that need it.

Instead of writing large amounts of host-language glue code, Ferret lets you describe the shape of the data you want and the operations needed to get it.
A Ferret program can query a document, interact with a browser-backed value, transform arrays and objects, call host-provided functions, and return structured data that can be consumed by another system.
Ferret is not limited to HTML scraping.

The web is one important use case, but the language is designed around a broader idea: values can support different operations depending on what they represent.

Structured data can be filtered, mapped, and transformed. A document can be queried. A browser-backed element can be queried and can receive dispatched events.

The runtime uses these distinctions — called capabilities — to decide what a script can do with a given value. The same language can work across different inputs as long as those inputs provide the capabilities the script needs.

## What you can build with Ferret
Ferret can be used for many kinds of targeted data workflows:
* extracting structured data from websites, documents, and APIs
* automating browser-driven workflows, including dynamic content and event-based waits
* normalizing and transforming external data into predictable structures
* embedding user-defined extraction logic into Go applications
* evaluating filters, mappings, and expressions in configuration-driven systems
* testing and validating APIs, HTML pages, and browser-driven interfaces
* and more!

Ferret can power scraping and data collection workflows, including workflows that collect unstructured data. Its main focus, however, is not raw scale for its own sake. Ferret is designed to make extraction logic explicit, repeatable, testable, and easy to embed into developer workflows.

## A first look at FQL
A Ferret script usually follows a simple pattern:
1. load or receive some input
2. query the input
3. transform the result
4. return structured data

For example, a script might query product cards from a document and return a normalized list of objects:

{{< editor lang="fql" readonly="true" params="false" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
let products = page[~ css`.product-card`]

return for product in products
    return {
        name: product[~ css`.product-title`],
        price: product[~ css`.product-price`],
        url: product[~ css`:attr("href", .product-link)`]
    }
{{</ editor >}}

The exact source of `page` depends on how Ferret is being used. It may come from a browser driver, a document loader, an embedded Go application, a test runner, or another runtime integration.

The important idea is that the script focuses on the extraction logic, while the host environment provides the values, functions, modules, and capabilities available at runtime.

## The core mental model

Ferret has a few core concepts that appear throughout the documentation.

### Scripts

Ferret programs are written in FQL.

A script describes how to query, transform, automate, and return data. Scripts can be run from the CLI, executed by a test runner, or embedded inside another application.

### Values

Ferret works with runtime values such as strings, numbers, booleans, arrays, objects, documents, elements, and module-defined values.

Some values are simple data. Others may expose behavior through capabilities.

For example, a plain object can be transformed. A document can be queried. A browser-backed element may support both querying and dispatching events.

### Capabilities

Capabilities are one of the central ideas in Ferret.

Instead of hard-coding every possible operation into the language, Ferret lets runtime values expose specific capabilities.

For example:

* a queryable value can be queried
* a dispatchable value can receive events or actions
* a readable value can provide data
* a module-defined value can expose domain-specific behavior

This keeps the core language small while allowing Ferret to support different data sources, document types, protocols, and runtime integrations.

### Modules and drivers

Ferret can be extended through modules and drivers.

Modules can add functions, data formats, protocols, integrations, or new runtime behavior. Drivers can provide capabilities for specific environments, such as HTML documents or browser-controlled pages.

This means Ferret’s language does not need special syntax for every external system. Instead, external systems can be exposed through values, functions, and capabilities.

## Ways to use Ferret

Ferret can be used in several different modes.

### As a CLI tool

Ferret can be used from the command line to run scripts, format code, inspect programs, and work with local extraction workflows.

This is the simplest way to start using Ferret.

### As an embedded runtime

Ferret can be embedded into Go applications.

In this mode, the host application provides input values, parameters, functions, modules, and runtime capabilities. Ferret provides the execution engine and the language used to describe the logic.

This is useful when extraction or transformation logic needs to be configurable, versioned, or provided outside the main application code.

### As an expression engine

Ferret can also be used as a small DSL inside configuration-driven applications. The host application evaluates Ferret expressions or scripts at runtime instead of hard-coding every filter, mapping, or transformation in Go.

This is useful when extraction rules, pipeline steps, validation checks, or automation logic need to be user-defined, versioned separately, or changed without redeploying the application.

The host application remains in control: it decides which functions are available, which values are passed into the script, which modules are loaded, and which capabilities the script can use.

For example, a billing service can pass invoice data directly from Go through the `@invoices` parameter. The FQL script selects open invoices and assigns each one to a review queue without loading a web page:

{{< code lang="fql" >}}
return for invoice in @invoices {
    filter invoice.status == "open"
    return {
        id: invoice.id,
        queue: match {
            when invoice.total >= 10000 => "review",
            _                           => "standard",
        }
    }
}
{{</ code >}}

The host owns the invoice records and execution boundary; the embedded script owns the configurable data decision.

### As part of a larger workflow

Ferret can also be used as one piece of a larger system.

For example, Ferret can extract and shape data, while another system handles storage, analytics, machine learning, reporting, or orchestration.

Ferret is designed to complement general-purpose languages and data tools, not replace them.

## What Ferret is not

Ferret is not a general-purpose programming language replacement.

It is not intended to replace Go, Python, JavaScript, or other languages used to build full applications.

Ferret is also not a massive web crawler for downloading the internet. Its focus is targeted, precise, repeatable extraction and automation.

Ferret is not limited to web scraping either. HTML and browser automation are important parts of the ecosystem, but Ferret’s core model is broader: querying, transforming, and automating capable values through a small declarative-first, expression-oriented language.

## The Ferret ecosystem

Ferret is more than a single executable.

The ecosystem includes:

- the Ferret language, runtime, and standard library
- the Ferret CLI for running, formatting, and debugging scripts
- optional modules and drivers for additional functions, data formats, and integrations
- embedding APIs for Go applications that need to control what scripts can see and do
- Lab, a test runner for Ferret scripts
- Mockery, a safe fake website used in examples, demos, and driver testing

A Ferret daemon is also in development to provide editor integration through the Language Server Protocol, including syntax highlighting, autocomplete, and debugging support.

These pieces are designed to work together while keeping the core language and runtime small.

## Next steps

If you are new to Ferret, start with the basics and then move into the areas that match how you plan to use it.

{{< docs-related tiles="getting-started-installation,getting-started-quick-start,language,web-extraction,embedding,tools-lab,tools-worker" >}}
