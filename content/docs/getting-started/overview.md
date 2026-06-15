---
title: "Overview"
sidebarTitle: "Overview"
weight: 20
draft: false
description: "What Ferret is, what it is useful for, and where to start."
aliases:
    - /docs/introduction/
---

# Overview

Ferret is a programmable data extraction and automation engine for developers.
It helps you describe where data comes from, how to query it, how to wait for it, how to transform it, and how to return clean structured output.

Ferret is built around FQL, a small declarative language for targeted data workflows.
A Ferret script can query documents, work with structured values, interact with browser-backed pages, call runtime functions, and return data that another system can consume.

Ferret is commonly used for web extraction, but it is not limited to web scraping.
The same model can be used with HTML, JSON, APIs, files, browser sessions, embedded application values, and module-defined data sources.

## What Ferret is useful for

Ferret is designed for workflows where extraction and transformation logic should be explicit, repeatable, testable, and easy to embed.

Common use cases include:

- extracting structured data from websites and documents
- automating browser-driven data collection
- waiting for dynamic content, events, or changing values
- normalizing external data into predictable objects and arrays
- validating API responses, HTML pages, or browser states
- embedding user-defined extraction rules into Go applications
- evaluating filters, mappings, and lightweight rules from configuration

Ferret is not meant to replace Go, Python, JavaScript, or other general-purpose languages.
It is a focused language and runtime for describing data-oriented automation logic.

## A first look at FQL

A Ferret script usually follows a simple pattern:

1. receive or load input
2. query the input
3. transform the result
4. return structured data

For example, a script can query product cards from a document and return a normalized list of objects:

{{< editor lang="fql" >}}
LET doc = DOCUMENT("https://mockery.ferretlang.org/scenarios/ecommerce/products/")
LET products = doc[~ css`.product-card`]

FOR product IN products
    RETURN {
        name: product[~? css`.product-title`],
        price: TRIM(product[~? css`.product-price`]),
        url: product[~? css`:attr('href', a.product-link)`]
    }
{{</ editor >}}

The exact source of `doc` depends on how Ferret is being used.
It may come from a browser driver, a document loader, an embedded Go application, a test runner, or another runtime integration.

The script focuses on the extraction logic.
The host environment provides the input values, functions, modules, and runtime capabilities available to the script.

## Main pieces

Ferret is made of a few pieces that can be used together or independently.

### FQL

FQL is the Ferret query language.
It is used to describe extraction, transformation, waiting, and automation logic.

### Runtime

The runtime executes FQL scripts.
It works with regular values such as strings, numbers, arrays, and objects, as well as host-provided values such as documents, elements, browser sessions, files, or module-defined values.

### Modules

Modules extend Ferret with functions, data formats, integrations, and runtime behavior.
They let Ferret work with different environments without adding special syntax for every external system.

### CLI

The Ferret CLI runs scripts from the command line and provides development tools such as formatting and inspection.
It is the simplest way to start using Ferret locally.

### Embedding

Ferret can be embedded into Go applications.
In this mode, the host application controls which values, parameters, functions, modules, and capabilities are available to scripts.

## Ways to use Ferret

Ferret can be used in several modes:

- as a CLI tool for local scripts and extraction workflows
- as an embedded runtime inside Go applications
- as a small DSL for configuration-driven filters, mappings, and rules
- as part of a larger data pipeline, test workflow, or automation system

In all modes, the goal is the same: keep data extraction and transformation logic clear, portable, and easy to reason about.

## What Ferret is not

Ferret is not a general-purpose application language.
It is not intended to replace the languages used to build full applications.

Ferret is also not a massive crawler for downloading the internet.
Its focus is targeted, precise, repeatable extraction and automation.

## Where to go next

If you are new to Ferret, start with installation and the quick start guide.
After that, move into the language overview or the area that matches how you plan to use Ferret.

{{< docs-related tiles="getting-started-installation,getting-started-quick-start,language,web-extraction,dynamic-pages,embedding,tools-lab,tools-worker" >}}
