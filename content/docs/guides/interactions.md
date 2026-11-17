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

Use `dispatch "click"` or the arrow shorthand `<-` to click:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/dynamic-products/load-more/", { driver: "cdp" })
let button = query one '[data-testid="load-more-products"]' in page using css
waitfor exists query one ".product-card" in page using css timeout 5s on timeout fail
let before = query count ".product-card" in page using css

dispatch "click" in button
return waitfor value (query count ".product-card" in page using css)
    when . > before
    timeout 5s on timeout fail
{{< /editor >}}

The query waits for the card count to increase after the click. To use the arrow form, replace the dispatch statement with `button <- "click"`.

## Fill a text input

Use `dispatch "input"` with a `with` payload to type into an input field:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/search/", {
    driver: "cdp"
})
let input = query one "#search-query" in page using css

waitfor (query exists '#search-status[data-state="ready"]' in page using css) timeout 5s on timeout fail
dispatch "input" in input with { value: "camera" }
dispatch "click" in (query one '[data-testid="search-submit"]' in page using css)
waitfor (query exists '#search-status[data-state="ready"]' in page using css) timeout 5s on timeout fail
return query count ".product-card" in page using css
{{< /editor >}}

## Select from a dropdown

Use `dispatch "select"` with the option value:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/search/", {
    driver: "cdp"
})
let categories = query one "#search-category" in page using css

waitfor (query exists '#search-status[data-state="ready"]' in page using css) timeout 5s on timeout fail
dispatch "select" in categories with { value: "laptops" }
dispatch "click" in (query one '[data-testid="search-submit"]' in page using css)
waitfor (query exists '#search-status[data-state="ready"]' in page using css) timeout 5s on timeout fail
return query count ".product-card" in page using css
{{< /editor >}}

## Submit a form

A typical form interaction combines filling inputs with a click or form submission:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/forms/", {
    driver: "cdp"
})

dispatch "input" in (query one "#query" in page using css) with { value: "ferret" }
dispatch "click" in (query one "#search-form button[type='submit']" in page using css)

let result = waitfor value (query one "#form-result" in page using css)
    when .textContent != ""
    timeout 3s
    on timeout fail

return result.textContent
{{</ editor >}}

## Wait for the result of an interaction

The `waitfor event ... trigger` pattern is the safest way to combine an interaction with waiting for its result. It subscribes to the event *before* triggering the action, avoiding a race condition where the event fires before listening begins:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/", { driver: "cdp" })
let next = query one '[data-testid="page-next"]' in page using css

waitfor event "navigation" in page
    trigger dispatch "click" in next
    timeout 10s

return query ".product-grid" in page using css
{{</ editor >}}

This reads as: start listening for a `navigation` event, then click the button, then wait until the event arrives or the timeout is reached.

## Focus

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/forms/", { driver: "cdp" })

let input = query one '#query' in page using css
dispatch "focus" in input
let test = query one '[data-testid="form-event-status"]' in page using css

return test.textContent
{{</ editor >}}

## Hover

Use `mouseover` and `mouseout` to trigger hover effects. Some pages may require a `mousemove` event instead:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/mouse/", { driver: "cdp" })

let target = query one '#mouse-hover-target' in page using css
dispatch "mouseover" in target
dispatch "mouseout" in target

let test = query one '#mouse-hover-status' in page using css

return test.textContent
{{</ editor >}}

Or, helper functions `HOVER` and `UNHOVER` can be used:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/mouse/", { driver: "cdp" })
let target = query one '#mouse-hover-target' in page using css
web::html::hover(target)
let test = query one '#mouse-hover-status' in page using css

return test.textContent
{{</ editor >}}

## Scroll the page

Use `web::html::scroll_bottom` or `web::html::scroll_top` to scroll the page, or `web::html::scroll_element` for a specific element.
The [pagination guide]({{< ref "pagination" >}}#infinite-scrolling) shows how to wait for the next batch of cards after a scroll.

## Multi-step interaction

Complex workflows chain several interactions together. Each step waits for the previous one to complete before proceeding:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/search/", { driver: "cdp" })

func PARSE_PRICE(product) {
    let priceNode = query one ".product-price" in product using css
    let priceText = priceNode.attributes["data-price"]
    let price = to_float(replace(priceText, "$", ""))
    return price
}

// Step 1a: fill and submit a search form
dispatch "input" in (query one "#search-query" in page using css) with { value: "laptop" }

// Step 1b: submit the form
dispatch "click" in (query one '[data-testid="search-submit"]' in page using css)

// Step 2a: wait for the results to load by checking the results status
waitfor (query exists '#search-status[data-state="ready"]' in page using css)
    timeout 10s on timeout fail

// Step 2b: collect product cards
let products = query ".product-card" in page using css

// Step 3: extract data from each product card
return for product in products
    return {
        brand: query one ".product-brand" in product using css,
        name: query one ".product-title" in product using css,
        price: PARSE_PRICE(product),
    }
{{</ editor >}}

## Error recovery for interactions

Interactions can fail — an element might not be clickable, or the page might not respond. Attach `on error return` to handle failures gracefully:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let button = query one ".optional-popup-close" in page using css
    on error return none

let dismissed = button != none
    ? (dispatch "click" in button on error return none)
    : none

return query "article" in page using css
{{</ editor >}}

## Next steps

{{< docs-related tiles="guide-pagination,guide-error-handling,language-control-flow-dispatch,language-control-flow-waitfor" >}}
