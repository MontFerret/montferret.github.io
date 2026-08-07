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
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/dynamic-products/basic/", {
    driver: "cdp"
})
LET button = QUERY ONE '[data-testid="page-next"]' IN page USING css

DISPATCH "click" IN button

// Arrow shorthand — equivalent to the above
button <- "click"
{{</ code >}}

After a click, the page may change. Use `WAITFOR` to wait for the result before extracting data:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/dynamic-products/basic/", {
    driver: "cdp"
})
LET button = QUERY ONE '[data-testid="page-next"]' IN page USING css

DISPATCH "click" IN button

WAITFOR EXISTS (QUERY ONE '[data-testid="dynamic-product-card"]' IN page USING css)
    TIMEOUT 5s

RETURN QUERY '[data-testid="dynamic-product-card"]' IN page USING css
{{</ editor >}}

## Fill a text input

Use `DISPATCH "input"` with a `WITH` payload to type into an input field:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/search/", {
    driver: "cdp"
})
LET input = QUERY ONE "#search-query" IN page USING css

LET beforeCount = QUERY COUNT '.product-card' IN page USING css
DISPATCH "input" IN input WITH { value: "camera" }
LET afterCount = QUERY COUNT '.product-card' IN page USING css

RETURN {
    beforeCount,
    afterCount
}
{{</ editor >}}

## Select from a dropdown

Use `DISPATCH "select"` with an array of values:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/search/", {
    driver: "cdp"
})
LET categories = QUERY ONE "#search-category" IN page USING css

LET beforeCount = QUERY COUNT '.product-card' IN page USING css
DISPATCH "select" IN categories WITH { value: "laptops" }
LET afterCount = QUERY COUNT '.product-card' IN page USING css

RETURN {
    beforeCount,
    afterCount
}
{{</ editor >}}

## Submit a form

A typical form interaction combines filling inputs with a click or form submission:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/forms/", {
    driver: "cdp"
})

DISPATCH "input" IN (QUERY ONE "#query" IN page USING css) WITH { value: "ferret" }
DISPATCH "click" IN (QUERY ONE "#search-form button[type='submit']" IN page USING css)

LET result = WAITFOR VALUE (QUERY ONE "#form-result" IN page USING css)
    WHEN .textContent != ""
    TIMEOUT 3s
    ON TIMEOUT FAIL

RETURN result.textContent
{{</ editor >}}

## Wait for the result of an interaction

The `WAITFOR EVENT ... TRIGGER` pattern is the safest way to combine an interaction with waiting for its result. It subscribes to the event *before* triggering the action, avoiding a race condition where the event fires before listening begins:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/products/", { driver: "cdp" })
LET next = QUERY ONE '[data-testid="page-next"]' IN page USING css

WAITFOR EVENT "navigation" IN page
    TRIGGER DISPATCH "click" IN next
    TIMEOUT 10s

RETURN QUERY ".product-grid" IN page USING css
{{</ editor >}}

This reads as: start listening for a `navigation` event, then click the button, then wait until the event arrives or the timeout is reached.

## Focus

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/forms/", { driver: "cdp" })

LET input = QUERY ONE '#query' IN page USING css
DISPATCH "focus" IN input
LET test = QUERY ONE '[data-testid="form-event-status"]' IN page USING css

RETURN test.textContent
{{</ editor >}}

## Hover

Use `mouseover` and `mouseout` to trigger hover effects. Some pages may require a `mousemove` event instead:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/mouse/", { driver: "cdp" })

LET target = QUERY ONE '#mouse-hover-target' IN page USING css
DISPATCH "mouseover" IN target
DISPATCH "mouseout" IN target

LET test = QUERY ONE '#mouse-hover-status' IN page USING css

RETURN test.textContent
{{</ editor >}}

Or, helper functions `HOVER` and `UNHOVER` can be used:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/mouse/", { driver: "cdp" })
LET target = QUERY ONE '#mouse-hover-target' IN page USING css
WEB::HTML::HOVER(target)
LET test = QUERY ONE '#mouse-hover-status' IN page USING css

RETURN test.textContent
{{</ editor >}}

## Scroll the page

Use `SCROLL_BOTTOM` or `SCROLL_TOP` to scroll the page, or `SCROLL_ELEMENT` for a specific element.
Here is a simple example of scrolling to the bottom of a page to load more products:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/infinite-scroll/", { driver: "cdp" })
LET pageSize = 8

FOR i WHILE SCROLL_BOTTOM(page)
    WAIT(500)
    FOR product IN QUERY `:skip(${i * pageSize}, .product-card)` IN page USING css
        RETURN {
            name: QUERY ONE '[data-testid="product-title"]' IN product USING css,
            price: QUERY ONE '[data-testid="product-price"]' IN product USING css,
        }
{{</ editor >}}

## Multi-step interaction

Complex workflows chain several interactions together. Each step waits for the previous one to complete before proceeding:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/search/", { driver: "cdp" })

FUNC PARSE_PRICE(product) {
    LET priceNode = QUERY ONE ".product-price" IN product USING css
    LET priceText = priceNode.attributes["data-price"]
    LET price = TO_FLOAT(SUBSTITUTE(priceText, "$", ""))
    RETURN price
}

// Step 1a: fill and submit a search form
DISPATCH "input" IN (QUERY ONE "#search-query" IN page USING css) WITH { value: "laptop" }

// Step 1b: submit the form
DISPATCH "click" IN (QUERY ONE '[data-testid="search-submit"]' IN page USING css)

// Step 2a: wait for the results to load by checking for the loader to disappear
WAITFOR EXISTS QUERY ONE "#search-loader" IN page USING css
    WHEN TO_BOOL(.attributes.disabled) == true
    TIMEOUT 10s

// Step 2b: collect product cards
LET products = QUERY ".product-card" IN page USING css

// Step 3: extract data from each product card
FOR product IN products
    RETURN {
        brand: QUERY ONE ".product-brand" IN product USING css,
        name: QUERY ONE ".product-title" IN product USING css,
        price: PARSE_PRICE(product),
    }
{{</ editor >}}

## Error recovery for interactions

Interactions can fail — an element might not be clickable, or the page might not respond. Attach `ON ERROR RETURN` to handle failures gracefully:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET button = QUERY ONE ".optional-popup-close" IN page USING css
    ON ERROR RETURN NONE

LET dismissed = button != NONE
    ? (DISPATCH "click" IN button ON ERROR RETURN NONE)
    : NONE

RETURN QUERY "article" IN page USING css
{{</ editor >}}

## Next steps

{{< docs-related tiles="guide-pagination,guide-error-handling,language-control-flow-dispatch,language-control-flow-waitfor" >}}
