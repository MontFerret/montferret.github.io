---
title: "Interact with pages"
sidebarTitle: "Page interactions"
weight: 30
draft: false
description: "Click elements, fill forms, and trigger browser events."
---

# Interact with pages

Many extraction tasks require interaction first — submitting a search form, clicking a filter, dismissing a popup, or logging in. This guide shows how to dispatch events to browser elements and wait for the results.

Page interactions require the `cdp` driver. See [Browser-driven pages]({{< ref "browser-pages" >}}) for setup.

## Click an element

Use `DISPATCH "click"` or the arrow shorthand `<-` to click:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET button = QUERY ONE "button.load-more" IN page USING css

DISPATCH "click" IN button

// Arrow shorthand — equivalent to the above
button <- "click"
{{</ code >}}

After a click, the page may change. Use `WAITFOR` to wait for the result before extracting data:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET button = QUERY ONE "button.load-more" IN page USING css

button <- "click"

WAITFOR EXISTS QUERY ONE ".results" IN page USING css
    TIMEOUT 5s

RETURN page[~ css`.results .item`]
{{</ code >}}

## Fill a text input

Use `DISPATCH "input"` with a `WITH` payload to type into an input field:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET input = QUERY ONE "input[name='search']" IN page USING css

DISPATCH "input" IN input WITH "ferret"
{{</ code >}}

You can also use the `INPUT` stdlib function:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
INPUT(page, "input[name='search']", "ferret")
{{</ code >}}

## Select from a dropdown

Use `DISPATCH "select"` with an array of values:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET dropdown = QUERY ONE "select[name='category']" IN page USING css

DISPATCH "select" IN dropdown WITH ["electronics"]
{{</ code >}}

Or with the `SELECT` stdlib function:

{{< code lang="fql" >}}
SELECT(page, "select[name='category']", ["electronics"])
{{</ code >}}

## Submit a form

A typical form interaction combines filling inputs with a click or form submission:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

INPUT(page, "input[name='search']", "ferret")

LET submit = QUERY ONE "button[type='submit']" IN page USING css
submit <- "click"

WAITFOR EXISTS QUERY ONE ".search-results" IN page USING css
    TIMEOUT 10s

RETURN page[~ css`.search-results .item`][*].textContent
{{</ code >}}

## Wait for the result of an interaction

The `WAITFOR EVENT ... TRIGGER` pattern is the safest way to combine an interaction with waiting for its result. It subscribes to the event *before* triggering the action, avoiding a race condition where the event fires before listening begins:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET button = QUERY ONE "button.submit" IN page USING css

WAITFOR EVENT "navigation" IN page
    TRIGGER ( button <- "click" )
    TIMEOUT 10s

RETURN page[~ css`.result`]
{{</ code >}}

This reads as: start listening for a `navigation` event, then click the button, then wait until the event arrives or the timeout is reached.

## Focus and hover

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET input = QUERY ONE "input[name='email']" IN page USING css
input <- "focus"

LET menu = QUERY ONE ".dropdown-trigger" IN page USING css
DISPATCH "hover" IN menu
{{</ code >}}

## Scroll the page

Use `SCROLL_BOTTOM` or `SCROLL_TOP` to scroll the page, or `SCROLL_ELEMENT` for a specific element:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

SCROLL_BOTTOM(page)
SCROLL_TOP(page)

LET container = QUERY ONE ".scroll-container" IN page USING css
SCROLL_ELEMENT(container, { y: 500 })
{{</ code >}}

## Multi-step interaction

Complex workflows chain several interactions together. Each step waits for the previous one to complete before proceeding:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

// Step 1: fill and submit a search form
INPUT(page, "input[name='search']", "ferret")

LET submit = QUERY ONE "button[type='submit']" IN page USING css
submit <- "click"

WAITFOR EXISTS QUERY ONE ".search-results" IN page USING css
    TIMEOUT 10s

// Step 2: click the first result
LET firstResult = QUERY ONE ".search-results a" IN page USING css
firstResult <- "click"

WAITFOR EVENT "navigation" IN page TIMEOUT 10s

// Step 3: extract data from the detail page
RETURN {
    title: page.title,
    content: QUERY ONE ".content" IN page USING css
}
{{</ code >}}

## Error recovery for interactions

Interactions can fail — an element might not be clickable, or the page might not respond. Attach `ON ERROR RETURN` to handle failures gracefully:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET button = QUERY ONE ".optional-popup-close" IN page USING css
    ON ERROR RETURN NONE

LET dismissed = button != NONE
    ? (DISPATCH "click" IN button ON ERROR RETURN NONE)
    : NONE

RETURN page[~ css`article`]
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-pagination,guide-error-handling,language-control-flow-dispatch,language-control-flow-waitfor" >}}
