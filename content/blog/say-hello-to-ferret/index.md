---
title: "Say 'Hello' to Ferret"
subtitle: "A modern web scraping tool"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-02-08"
lastmod: 2026-06-25
---

As a developer, you have probably written a web scraper at least once.

Maybe you needed product prices from a few websites. Maybe you wanted to collect images for an AI experiment. Maybe you needed public data from a website that did not provide an API. Or maybe you just needed to turn a messy web page into a clean JSON file.

At first, this usually feels simple.

You pick a language, install a library, make an HTTP request, parse some HTML, extract a few fields, and move on.

But then the simple script starts growing.

You need to handle pagination. Then retries. Then dynamic pages rendered by JavaScript. Then waiting for elements. Then browser automation. Then configuration. Then output shaping. Then another website with slightly different markup. Then another script. Then another small framework around all those scripts.

That is where Ferret comes in.

Ferret is a declarative data extraction tool built around FQL, a query language for describing what data you want and what shape you want it in.

Instead of building a new scraping project from scratch every time, Ferret lets you write compact, readable programs focused on the data itself.

## The problem

There are many good libraries and frameworks for fetching pages, parsing HTML, controlling browsers, and automating websites.

They solve important problems.

But when you build data extraction logic, there are still a few recurring issues that tend to come back again and again.

First, there is boilerplate.

Even with a good library, you still need to write glue code. You need to set up requests or browser sessions, parse documents, select elements, transform values, handle errors, and produce the final output. None of that is especially hard, but it is work you repeat every time.

Second, websites change.

Markup changes. Class names move around. Pages get redesigned. A scraper that worked yesterday might need a small update today. If the extraction logic is buried inside a deployed application, even a simple selector change can become more annoying than it should be.

Third, the web is dynamic.

A plain HTTP request is often not enough anymore. Many pages are rendered by JavaScript. Some data appears only after a delay. Some pages require user interaction: clicking a button, typing into a form, opening a menu, scrolling, or moving to the next page.

All of these problems are solvable.

But it would be nice not to solve them from scratch every time.

## The solution

Ferret was created to make this kind of work more direct.

It provides a declarative language for data extraction and automation. You describe the data you need, how to find it, and how the result should look. Ferret handles the runtime work around that description.

At its simplest, an FQL program loads a document, queries it, and returns structured data.

{{< editor lang="fql" readonly="true" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/products/")

FOR card IN page[~ css`.product-card`]
    RETURN {
        title: card[~? css`.product-title`].textContent,
        price: card[~? css`.product-price`].textContent,
        url: card[~? css`a`].attributes.href
    }
{{</ editor >}}

The important part is not the specific website or selector syntax. The important part is the shape of the program.

Load a document. Find matching elements. Iterate over them. Return clean structured data.

That is the core idea behind Ferret.

## Focus on the data, not the glue

Ferret is designed so that extraction logic can stay close to the shape of the data you want.

You do not need to start with a full application just to describe a small data extraction task. You do not need to wrap every script in request setup, parsing setup, browser setup, and output conversion.

FQL gives you language-level constructs for the common parts of data extraction:

{{< editor lang="fql" readonly="true" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/")

FOR article IN page[~ css`article`]
    LET title = article[~? css`h2`].textContent
    LET summary = article[~? css`p`].textContent
    RETURN {
        title,
        summary
    }
{{</ editor >}}

The result is a collection of objects, not a pile of intermediate parsing code.

That also makes scripts easier to read later. When extraction rules need to change, you can update the FQL program itself instead of digging through application logic that mixes infrastructure, navigation, parsing, and transformation in one place.

## Dynamic pages

Of course, not every page is available as static HTML.

Many modern websites render content in the browser. A document may load first, then fetch data, then render the part you care about a second later. In those cases, the extraction script needs to wait for the page to reach the right state before reading from it.

Ferret supports this kind of workflow too.

{{< editor lang="fql" readonly="true" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/dynamic-products/basic/", {
    driver: "cdp"
})

LET cards = WAITFOR VALUE page[~ css`.product-card`]
    TIMEOUT 5s

FOR card IN cards
    RETURN {
        title: card[~? css`.product-title`].textContent,
        status: card[~? css`.product-stock`].textContent
    }
{{</ editor >}}

The query still describes the data, but it also gives Ferret enough information to deal with the dynamic nature of the page.

Use a browser-backed driver. Wait for the content to appear. Then extract the result.

The script does not need to change into a completely different kind of program just because the page uses JavaScript.

## Page interactions

Sometimes waiting is not enough.

Some pages require interaction before the data appears. A user might need to type a search query, click a button, open a filter, or move through pagination.

Ferret treats those actions as part of the program.

{{< editor lang="fql" params=`{"q": "ferret"}` >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/forms/", {
    driver: "cdp"
})

DISPATCH "input" IN page[~? css`#query`] WITH {
    value: @q
}

DISPATCH "click" IN page[~? css`button[type="submit"]`]

LET result = WAITFOR VALUE page[~ css`#form-result`] TIMEOUT 5s

RETURN result
{{</ editor >}}

This is where Ferret starts to feel less like a traditional scraper and more like a small language for data-oriented browser automation.

The program describes the same steps a person would take:

Open the page. Enter a search term. Submit the form. Wait for the results. Extract the data.

The difference is that the final output is structured and repeatable.

## More than just selectors

A useful extraction script rarely stops at selecting elements.

You usually need to clean text, normalize values, handle missing fields, filter results, limit output, combine values, or reshape the final result.

FQL includes common language features for that work: variables, expressions, loops, conditionals, objects, arrays, functions, and parameters.

{{< code lang="fql" >}}
FOR product IN page[~ css`.product`]
    LET title = TRIM(product[~? css`.title`].textContent)
    LET priceText = TRIM(product[~? css`.price`].textContent)
    
    RETURN {
        title,
        price: priceText,
        available: product[~? css`.sold-out`] == NONE
    }
{{</ code >}}

The goal is not to hide programming completely.

The goal is to give data extraction its own language, so the logic can stay concise, readable, and focused.

## Embeddable by design

Ferret can be used from the command line, but it was designed to be embedded into applications as well.

That matters because extraction logic often lives inside a larger system.

A backend service might need to run FQL scripts. A data pipeline might need to execute them on a schedule. A testing tool might use them to verify what appears in a browser. An application might expose its own host values and let FQL query them.

Ferret is not only a standalone tool. It is also a runtime that can be integrated into other software.

That is one of the reasons FQL exists as a dedicated language instead of just another library API. Scripts can be stored, updated, passed around, and executed by a host application without turning every extraction rule into application code.

## Extensible runtime

Ferret is also extensible.

Applications can provide their own functions, modules, and host values. That means FQL programs are not limited to built-in document querying.

A host application can decide what values exist, how they behave, and what it means to query them.

For example, one value might represent an HTML document. Another might represent an API client. Another might represent a database connection. Another might represent something specific to the application embedding Ferret.

The language stays the same, while the runtime can grow around the needs of the host environment.

## Who is Ferret for?

Ferret is useful for anyone who needs to extract or automate data from places where a small script is not quite enough, but a full custom scraping system feels like too much.

That includes developers who need to collect data from websites, data engineers building repeatable extraction tasks, QA engineers testing browser-rendered content, researchers collecting public datasets, and teams that want extraction logic to be easier to read and update.

It is especially useful when the extraction rules matter as much as the application around them.

If the logic changes often, it helps to keep that logic in a small declarative program. If the same task needs to run in different environments, it helps to have a runtime. If the target is dynamic or interactive, it helps to have browser automation available without turning the whole script into low-level control flow.

## What Ferret is not

Ferret is not a magic solution for every data problem. It will not make unstable websites stable, remove the need to understand the page or system you are querying, or turn messy data into clean data without clear rules.

Like any extraction tool, it also needs to be used responsibly. Not every website should be automated, and not every piece of data should be collected. Rate limits, terms of service, robots.txt, privacy, and basic respect for other people’s systems still matter. Ferret does not remove those responsibilities; it simply gives you a better language and runtime for the work when that work is appropriate.

## What next?

Ferret is an open source project, and there is still a lot of work ahead. The current focus is on making the language, runtime, tools, and documentation more stable, more consistent, and easier to use. That includes clearer examples, stronger CLI support, better browser-oriented flows, improved embedding documentation, more modules, and eventually better editor tooling.

The long-term goal is to make structured data extraction easier to write, easier to read, and easier to maintain. Ferret started as a way to avoid writing the same scraping glue code over and over again, and that motivation is still very much alive. More broadly, though, Ferret is about extracting data with intent: you describe what you want, and Ferret gives you a language and runtime designed to help you get it.