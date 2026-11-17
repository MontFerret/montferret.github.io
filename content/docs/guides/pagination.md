---
title: "Pagination patterns"
sidebarTitle: "Pagination"
weight: 40
draft: false
description: "Collect data across multiple pages using common pagination patterns."
---

# Pagination patterns

Most websites split data across multiple pages. This guide shows how to handle the most common pagination patterns: clicking a "next" button, iterating numbered pages by URL, and collecting results across all pages.

Clicking a link or scrolling to load more content requires the `cdp` driver. Numbered pages with server-rendered HTML can be fetched using the memory driver. See [Browser-driven pages]({{< ref "browser-pages" >}}) for setup.

## Click a "next" button

The most common pagination pattern is to click a Next link until it disappears.

Use `do while` because the loop body must run at least once to process the first page. After each iteration, Ferret checks whether the next-page link still exists.

{{< editor lang="fql" >}}
// Open the first page using the browser-based CDP driver.
// Pagination requires a browser because each page is loaded by clicking a link.
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/", { driver: "cdp" })

// Keep selectors in variables so they can be reused throughout the query.
let nextSelector = "a.next[href]"
let itemSelector = ".product-card"

// Find the "next" link, click it, and wait for the resulting navigation.
// The TRIGGER block ensures the navigation listener is active before the
// click is dispatched, avoiding race conditions.
func CLICK_NEXT() {
    return waitfor event "navigation" in page
        trigger (dispatch "click" in (query one nextSelector in page using css))
        timeout 10s on timeout fail
}

// Skip navigation during the first iteration because the first page
// has already been loaded. Starting with the second iteration,
// move to the next page before extracting its products.
func NEXT_PAGE(pageNum) {
    return match {
        when pageNum > 0 => CLICK_NEXT(),
        _ => none
    }
}

// Process the current page, then continue while a “next” link exists.
return for i do while query exists nextSelector in page using css
    // Prevent accidental infinite or unexpectedly long pagination.
    limit 5
    
    // The first iteration processes the initial page.
    // Later iterations click "next" and wait for navigation.
    NEXT_PAGE(i)
    // Extract every product from the currently loaded page.
    for item in query itemSelector in page using css
        return {
            // Optional queries return NONE when an element is missing.
            // Optional chaining also avoids accessing textContent on NONE.
            name: item[~? css`.product-title`]?.textContent,
            price: item[~? css`.product-price`]?.textContent
        }

{{</ editor >}}

Key points

* `for do while` processes the first page. Unlike a regular `for while` loop, it executes its body before evaluating the continuation condition, making it ideal for pagination because the initial page is already loaded.
* The loop index distinguishes the first page. The first iteration has an index of 0. `NEXT_PAGE(i)` skips navigation on that iteration and clicks Next only on subsequent pages.
* `waitfor ... trigger` prevents timing races. The `trigger` block executes only after `waitfor` has started listening for the navigation event. This guarantees that even very fast navigations cannot occur before the event listener is ready.
* Navigation completes before extraction. `CLICK_NEXT()` returns only after the browser has finished navigating, so the extraction loop always runs against the newly loaded page.
* `query exists` controls pagination. After each iteration, Ferret checks whether a matching Next link still exists. When it no longer does, the loop terminates automatically.
* `query one` locates the navigation control. It returns the single Next link that serves as the target for the click operation.
* `limit` is a safety guard. It prevents accidental infinite or unexpectedly long pagination loops if a site behaves incorrectly.
* Optional queries tolerate missing data. The `~?` query operator returns `none` when an element is missing. Combined with optional chaining (`?.`), incomplete product cards can be processed without failing the query.

## Infinite scrolling

On Mockery’s infinite-scroll page the first batch appears automatically. Scroll after that batch renders, then wait for the count to increase before reading the appended cards:

{{< editor lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/infinite-scroll/", { driver: "cdp" })
let firstBatch = waitfor value (query count ".product-card" in page using css)
    when . > 0 timeout 5s on timeout fail

web::html::scroll_bottom(page)
let cards = waitfor value (query ".product-card" in page using css)
    when length(.) > firstBatch timeout 5s on timeout fail

return for card, index in cards
    filter index >= firstBatch
    return card[~? css`[data-testid="product-title"]`]?.textContent
{{< /editor >}}

For a longer collection, repeat until the page signals completion. Limit the number of scrolls and give each wait a deadline.

## Iterate numbered pages by URL

When pages are addressable by URL (e.g., `?page=1`, `?page=2`), use a `for` loop with a range:

Annotated query

{{< editor lang="fql" >}}
// Base URL of the paginated resource.
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products"
// Total number of pages to visit.
let totalPages = 5
// Iterate over every page number in the range.
return for pageNum in 1..totalPages
    // Build the URL for the current page.
    // The first page uses the base URL, while subsequent pages
    // follow the site's numbered URL pattern.
    let url = pageNum == 1 ? baseURL: `${baseURL}/page/${pageNum}`
    // Fetch the page over HTTP.
    // Because each page has its own URL, no browser automation is required.
    let page = web::html::open(url)
    // Extract every product from the current page.
    for item in page[~ css`.product-card`]
        return {
            page: pageNum,
            // Optional queries return NONE when an element is missing.
            // Optional chaining avoids accessing properties on NONE.
            name: item[~? css`.product-title`]?.textContent,
            price: item[~? css`.product-price`]?.textContent
    }
{{</ editor >}}
Key points

* A range drives the pagination. `1..totalPages` generates the sequence of page numbers to visit, making it easy to iterate over numbered URLs.
* Each URL is built independently. The current page number is used to construct the appropriate URL for each request. The first page often uses a different URL pattern than subsequent pages.
* Each page is fetched directly. `web::html::open()` loads every page independently over HTTP, so no browser session or navigation is required.
* Pages are processed independently. Each iteration loads, extracts, and returns data from one page before moving on to the next.
* Browser automation is unnecessary. This approach is faster and more resource-efficient because it avoids rendering pages or simulating user interactions.
* Best suited for static websites. Use this pattern when the desired content is present in the initial HTML returned by the server. If the content is loaded dynamically with JavaScript, use browser automation instead.
* Optional queries tolerate missing data. The `~?` query operator returns `none` when an element is missing. Combined with optional chaining (`?.`), incomplete product cards can be processed without failing the query.

### Detect the last page

If you do not know the total number of pages, open the first page to find out:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products/"
let firstPage = web::html::open(baseURL)

let totalPages = query count '[data-testid="page-link"]' in firstPage using css

return for pageNum in 1..totalPages
    let page = web::html::open(pageNum == 1 ? baseURL : `${baseURL}page/${pageNum}/`)
    for item in page[~ css`.product-card`]
        return item[~? css`.product-title`]?.textContent
{{</ code >}}

## Collect results into a flat array

When each page returns an array of items, the outer loop produces an array of arrays. Use `arrays::flatten` to merge them:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products/"

let result = (
    for pageNum in 1..3
        let page = web::html::open(pageNum == 1 ? baseURL : `${baseURL}page/${pageNum}/`)

        let items = (
            for item in page[~ css`.product-card`]
                return {
                    name: item[~? css`.product-title`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        return items
)

return arrays::flatten(result)
{{</ code >}}

Alternatively, use the `[**]` array contraction operator to flatten inline:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products/"
let result = (
    for pageNum in 1..3
        let page = web::html::open(pageNum == 1 ? baseURL : `${baseURL}page/${pageNum}/`)
        return page[~ css`.product-card`][*
            return {
                name: .[~? css`.product-title`]?.textContent,
                price: .[~? css`.product-price`]?.textContent
            }
        ]
)

return result[**]
{{</ code >}}

## Add error recovery

Pagination scripts are long-running and may encounter network errors or missing elements. Wrap page loads and interactions with error recovery:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products/"

let result = (
    for pageNum in 1..10
        let page = web::html::open(pageNum == 1 ? baseURL : `${baseURL}page/${pageNum}/`)
            on error retry 2 delay 1s backoff EXPONENTIAL
            or return none

        filter page != none

        let items = (
            for item in page[~ css`.product-card`]
                return {
                    page: pageNum,
                    name: item[~? css`.product-title`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        return items
)

return arrays::flatten(result)
{{</ code >}}

See [Error handling and resilience]({{< ref "error-handling" >}}) for more patterns.

## Next steps

{{< docs-related tiles="guide-error-handling,guide-interactions,language-control-flow-for,language-control-flow-error-handling" >}}
