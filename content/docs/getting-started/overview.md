---
title: "Overview"
sidebarTitle: "Overview"
weight: 10
draft: false
description: "What Ferret is and how its main pieces fit together."
aliases:
  - /docs/introduction/
---

# Overview
Ferret is a programmable data extraction and automation engine for developers.
It provides a small declarative language, an embeddable runtime, and an extensible execution model for querying, transforming, and automating data from websites, documents, APIs, and application-defined sources.
Ferret is especially useful when data lives behind messy HTML, browser interactions, inconsistent external systems, or configuration-driven workflows that need more structure than ad-hoc scripts.
At a high level, Ferret helps you describe:
- where data comes from
- how to query it
- how to wait for it
- how to transform it
- how to return clean structured output
  Ferret can be used as a command-line tool, embedded into Go applications, extended with modules, or used as a small domain-specific language for configuration-driven systems.

## What Ferret is
Ferret is built around FQL, a declarative language designed for data extraction, transformation, and automation workflows.
Instead of writing large amounts of imperative glue code, Ferret lets you describe the shape of the data you want and the operations needed to get it.
A Ferret program can query a document, interact with a browser-backed value, transform arrays and objects, call host-provided functions, and return structured data that can be consumed by another system.
Ferret is not limited to HTML scraping. 

The web is one important use case, but the language is designed around a broader idea: values expose capabilities, and the runtime uses those capabilities to decide what operations are available.
That means the same language can work with different kinds of inputs, as long as those inputs provide the capabilities the script needs.

## What you can build with Ferret
Ferret can be used for many kinds of targeted data workflows:
- extracting structured data from websites
- querying HTML, JSON, XML, CSV, and other document-like data
- automating browser-driven workflows
- waiting for dynamic content, events, or changing values
- normalizing external data into predictable structures
- embedding user-defined extraction rules into Go applications
- evaluating filters, mappings, and expressions in configuration-driven systems
- testing API responses, HTML pages, and browser-driven user interfaces
- validating that external systems return expected structured data
- writing repeatable workflow checks for dynamic pages or application states
- and more!

Ferret can power scraping and data collection workflows, including workflows that collect unstructured data. Its main focus, however, is not raw scale for its own sake. Ferret is designed to make extraction logic explicit, repeatable, testable, and easy to embed into developer workflows.

## A first look at FQL
A Ferret script usually follows a simple pattern:
1. load or receive some input
2. query the input
3. transform the result
4. return structured data
   For example, a script might query product cards from a document and return a normalized list of objects:

{{< code lang="fql" height="160px" >}}
LET products = doc[~ css`.product-card`]
FOR product IN products
    RETURN {
        name: product[~ css`.product-title`],
        price: product[~ css`.price`],
        url: product[~ css`a`]
    }
{{</ code >}}

The exact source of doc depends on how Ferret is being used. It may come from a browser driver, a document loader, an embedded Go application, a test runner, or another runtime integration.

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

Ferret can also be used as a small DSL for configuration-driven applications.

Instead of hard-coding every filter, rule, mapping, or transformation in Go, an application can evaluate Ferret expressions or scripts at runtime.

This is useful for systems that need user-defined or configuration-defined behavior, such as:

* filtering records
* mapping external data into internal structures
* defining extraction rules
* evaluating lightweight business rules
* transforming API responses
* configuring pipeline steps
* describing application-specific automation logic

In this mode, the host application remains in control. It decides which functions are available, which values are passed into the script, which modules are loaded, and which capabilities the script can use.

### As part of a larger workflow

Ferret can also be used as one piece of a larger system.

For example, Ferret can extract and shape data, while another system handles storage, analytics, machine learning, reporting, or orchestration.

Ferret is designed to complement general-purpose languages and data tools, not replace them.

## What Ferret is not

Ferret is not a general-purpose programming language replacement.

It is not intended to replace Go, Python, JavaScript, or other languages used to build full applications.

Ferret is also not a massive web crawler for downloading the internet. Its focus is targeted, precise, repeatable extraction and automation.

Ferret is not limited to web scraping either. HTML and browser automation are important parts of the ecosystem, but Ferret’s core model is broader: querying, transforming, and automating capable values through a small declarative language.

## The Ferret ecosystem

Ferret is more than a single executable.

The ecosystem includes:

* the Ferret language and runtime
* the Ferret CLI
* standard library functions
* optional modules and drivers
* embedding APIs for Go applications
* Lab, a test runner for Ferret scripts
* Mockery, a safe fake website for demos, examples, and driver testing

These pieces are designed to work together while keeping the core language and runtime small.

## Where to go next

If you are new to Ferret, start with the basics and then move into the areas that match how you plan to use it.

<div class="docs-card-grid">
  <a class="docs-card" href="{{< ref "installation" >}}">
    <span class="docs-card-kicker">Setup</span>
    <strong>Install Ferret</strong>
    <span>Set up the CLI or add Ferret as a Go library.</span>
  </a>
  <a class="docs-card" href="{{< ref "quick-start" >}}">
    <span class="docs-card-kicker">Try</span>
    <strong>Run your first script</strong>
    <span>Execute a Ferret query against static or dynamic pages.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">Learn</span>
    <strong>Learn the language basics</strong>
    <span>Understand variables, expressions, arrays, objects, control flow, and functions.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">Extract</span>
    <strong>Query documents</strong>
    <span>Use query expressions and dialects such as CSS to extract structured data.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">Browse</span>
    <strong>Work with dynamic pages</strong>
    <span>Wait for elements, dispatch events, and automate browser interactions.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">Embed</span>
    <strong>Embed Ferret in Go</strong>
    <span>Use Ferret as an embedded runtime or expression engine.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">Test</span>
    <strong>Test scripts with Lab</strong>
    <span>Run Ferret scripts against controlled pages and scenarios.</span>
  </a>
</div>