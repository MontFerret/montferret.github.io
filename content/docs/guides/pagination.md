---
title: "Pagination patterns"
sidebarTitle: "Pagination"
weight: 40
draft: false
description: "Collect data across multiple pages using common pagination patterns."
---

# Pagination patterns

Most websites split data across multiple pages. This guide shows how to handle the most common pagination patterns: clicking a "next" button, iterating numbered pages by URL, and collecting results across all pages.

All pagination examples use the `cdp` driver because page navigation and clicking require a browser. See [Browser-driven pages]({{< ref "browser-pages" >}}) for setup.

## Click a "next" button

The most common pagination pattern is to click a Next link until it disappears.

Use `do while` because the loop body must run at least once to process the first page. After each iteration, Ferret checks whether the next-page link still exists.

{{< editor lang="fql" >}}
// Open the first page using the browser-based CDP driver.
// Pagination requires a browser because each page is loaded by clicking a link.
let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/products/", { driver: "cdp" })

// Keep selectors in variables so they can be reused throughout the query.
let nextSelector = "a.next"
let itemSelector = ".product-card"

// Find the "next" link, click it, and wait for the resulting navigation.
// The TRIGGER block ensures the navigation listener is active before the
// click is dispatched, avoiding race conditions.
func CLICK_NEXT() {
    return waitfor event "navigation" in page
        trigger (dispatch "click" in (query one nextSelector in page using css))
        timeout 10s
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

Some websites load additional content when the user scrolls to the bottom of the page. Use `web::html::scroll_bottom` to trigger another content load, then process only the newly added items.

{{< editor lang="fql" >}}
// Open the page using the CDP browser driver.
// Infinite scrolling depends on browser interaction and dynamically loaded content.
let page = web::html::open("https://mockery.ferretlang.org/scenarios/infinite-scroll/", { driver: "cdp" })

// Number of products added after each scroll.
let pageSize = 8

// Process the initial content, then keep scrolling until the page
// can no longer be scrolled further.
return for i do while web::html::scroll_bottom(page)
    // The first batch is already available.
    // For subsequent iterations, briefly wait for newly requested
    // content to be added to the page.
    let _ = i == 0 ? none : wait(500)
    // Skip products processed during previous iterations and return
    // only the next batch of newly loaded items.
    for product in query `:skip(${i * pageSize}, .product-card)` in page using css
        return {
        name: (query one '[data-testid="product-title"]' in product using css)?.textContent,
        price: (query one '[data-testid="product-price"]' in product using css)?.textContent
    }
{{</ editor >}}

Key points

* Infinite scrolling requires a browser. The page uses JavaScript and user interaction to load additional content, so the query opens it with the CDP driver.
* `for do while` processes the initial batch. The loop body runs before its condition is evaluated, allowing the products already present on the first page to be extracted before the first scroll.
* `web::html::scroll_bottom` drives the loop. After each iteration, Ferret scrolls to the bottom of the page. The loop continues while further scrolling is possible.
* The loop index identifies each batch. On the first iteration, i is 0, so the query processes the products initially present on the page. Later values correspond to batches loaded by subsequent scrolls.
* The first iteration does not wait. The initial products are already available, so `wait` is used only after a scroll has triggered another content load.
* The temporary wait avoids reading the page too early. It gives the website time to fetch and render newly loaded products before extraction begins. This is a temporary workaround until `page.network.status` can be used to wait for the page to become idle.
* `:skip()` prevents duplicate results. Each iteration skips the products returned by earlier iterations and selects only the newly loaded batch.
* `pageSize` must match the website’s behavior. The skip offset assumes that every scroll loads eight products. If the site returns variable-sized batches, track the number of already processed elements instead.
* `query one` extracts fields from each product. Each title and price lookup is scoped to the current product card rather than the entire page.
* This pattern needs a stopping condition. The loop ends when `web::html::scroll_bottom(page)` indicates that the page can no longer scroll further. A `limit` can also be added as a safety guard for sites that scroll indefinitely.

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
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="
let firstPage = web::html::open(baseURL + "1")

let lastPageLink = query one ".pagination a:last-child" in firstPage using css
let totalPages = to_int(lastPageLink?.textContent) on error return 1

return for pageNum in 1..totalPages
    let page = web::html::open(baseURL + to_string(pageNum))
    for item in page[~ css`.product-card`]
        return item[~? css`.product-name`]?.textContent
{{</ code >}}

## Collect results into a flat array

When each page returns an array of items, the outer loop produces an array of arrays. Use `flatten` to merge them:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="

let result = (
    for pageNum in 1..3
        let page = web::html::open(baseURL + to_string(pageNum))

        let items = (
            for item in page[~ css`.product-card`]
                return {
                    name: item[~? css`.product-name`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        return items
)

return flatten(result)
{{</ code >}}

Alternatively, use the `[**]` array contraction operator to flatten inline:

{{< code lang="fql" >}}
let result = (
    for pageNum in 1..3
        let page = web::html::open("https://mockery.ferretlang.org/scenarios/ecommerce/?page=" + to_string(pageNum))
        return page[~ css`.product-card`][*
            return {
                name: .[~? css`.product-name`]?.textContent,
                price: .[~? css`.product-price`]?.textContent
            }
        ]
)

return result[**]
{{</ code >}}

## Add error recovery

Pagination scripts are long-running and may encounter network errors or missing elements. Wrap page loads and interactions with error recovery:

{{< code lang="fql" >}}
let baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="

let result = (
    for pageNum in 1..10
        let page = web::html::open(baseURL + to_string(pageNum))
            on error retry 2 delay 1s backoff EXPONENTIAL
            or return none

        filter page != none

        let items = (
            for item in page[~ css`.product-card`]
                return {
                    page: pageNum,
                    name: item[~? css`.product-name`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        return items
)

return flatten(result)
{{</ code >}}

See [Error handling and resilience]({{< ref "error-handling" >}}) for more patterns.

## Next steps

{{< docs-related tiles="guide-error-handling,guide-interactions,language-control-flow-for,language-control-flow-error-handling" >}}
