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

Use `DO WHILE` because the loop body must run at least once to process the first page. After each iteration, Ferret checks whether the next-page link still exists.

{{< editor lang="fql" >}}
// Open the first page using the browser-based CDP driver.
// Pagination requires a browser because each page is loaded by clicking a link.
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/products/", { driver: "cdp" })

// Keep selectors in variables so they can be reused throughout the query.
LET nextSelector = "a.next"
LET itemSelector = ".product-card"

// Find the "next" link, click it, and wait for the resulting navigation.
// The TRIGGER block ensures the navigation listener is active before the
// click is dispatched, avoiding race conditions.
FUNC CLICK_NEXT() {
    RETURN WAITFOR EVENT "navigation" IN page
        TRIGGER (DISPATCH "click" IN (QUERY ONE nextSelector IN page USING css))
        TIMEOUT 10s
}

// Skip navigation during the first iteration because the first page
// has already been loaded. Starting with the second iteration,
// move to the next page before extracting its products.
FUNC NEXT_PAGE(pageNum) {
    RETURN MATCH {
        WHEN pageNum > 0 => CLICK_NEXT(),
        _ => NONE
    }
}

// Process the current page, then continue while a “next” link exists.
RETURN FOR i DO WHILE QUERY EXISTS nextSelector IN page USING css
    // Prevent accidental infinite or unexpectedly long pagination.
    LIMIT 5
    
    // The first iteration processes the initial page.
    // Later iterations click "next" and wait for navigation.
    NEXT_PAGE(i)
    // Extract every product from the currently loaded page.
    FOR item IN QUERY itemSelector IN page USING css
        RETURN {
            // Optional queries return NONE when an element is missing.
            // Optional chaining also avoids accessing textContent on NONE.
            name: item[~? css`.product-title`]?.textContent,
            price: item[~? css`.product-price`]?.textContent
        }

{{</ editor >}}

Key points

* `FOR DO WHILE` processes the first page. Unlike a regular `FOR WHILE` loop, it executes its body before evaluating the continuation condition, making it ideal for pagination because the initial page is already loaded.
* The loop index distinguishes the first page. The first iteration has an index of 0. `NEXT_PAGE(i)` skips navigation on that iteration and clicks Next only on subsequent pages.
* `WAITFOR ... TRIGGER` prevents timing races. The `TRIGGER` block executes only after `WAITFOR` has started listening for the navigation event. This guarantees that even very fast navigations cannot occur before the event listener is ready.
* Navigation completes before extraction. `CLICK_NEXT()` returns only after the browser has finished navigating, so the extraction loop always runs against the newly loaded page.
* `QUERY EXISTS` controls pagination. After each iteration, Ferret checks whether a matching Next link still exists. When it no longer does, the loop terminates automatically.
* `QUERY ONE` locates the navigation control. It returns the single Next link that serves as the target for the click operation.
* `LIMIT` is a safety guard. It prevents accidental infinite or unexpectedly long pagination loops if a site behaves incorrectly.
* Optional queries tolerate missing data. The `~?` query operator returns `NONE` when an element is missing. Combined with optional chaining (`?.`), incomplete product cards can be processed without failing the query.

## Infinite scrolling

Some websites load additional content when the user scrolls to the bottom of the page. Use SCROLL_BOTTOM to trigger another content load, then process only the newly added items.

{{< editor lang="fql" >}}
// Open the page using the CDP browser driver.
// Infinite scrolling depends on browser interaction and dynamically loaded content.
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/infinite-scroll/", { driver: "cdp" })

// Number of products added after each scroll.
LET pageSize = 8

// Process the initial content, then keep scrolling until the page
// can no longer be scrolled further.
RETURN FOR i DO WHILE SCROLL_BOTTOM(page)
    // The first batch is already available.
    // For subsequent iterations, briefly wait for newly requested
    // content to be added to the page.
    LET _ = i == 0 ? NONE : WAIT(500)
    // Skip products processed during previous iterations and return
    // only the next batch of newly loaded items.
    FOR product IN QUERY `:skip(${i * pageSize}, .product-card)` IN page USING css
        RETURN {
        name: (QUERY ONE '[data-testid="product-title"]' IN product USING css)?.textContent,
        price: (QUERY ONE '[data-testid="product-price"]' IN product USING css)?.textContent
    }
{{</ editor >}}

Key points

* Infinite scrolling requires a browser. The page uses JavaScript and user interaction to load additional content, so the query opens it with the CDP driver.
* `FOR DO WHILE` processes the initial batch. The loop body runs before its condition is evaluated, allowing the products already present on the first page to be extracted before the first scroll.
* `SCROLL_BOTTOM` drives the loop. After each iteration, Ferret scrolls to the bottom of the page. The loop continues while further scrolling is possible.
* The loop index identifies each batch. On the first iteration, i is 0, so the query processes the products initially present on the page. Later values correspond to batches loaded by subsequent scrolls.
* The first iteration does not wait. The initial products are already available, so `WAIT` is used only after a scroll has triggered another content load.
* The temporary wait avoids reading the page too early. It gives the website time to fetch and render newly loaded products before extraction begins. This is a temporary workaround until `page.network.status` can be used to wait for the page to become idle.
* `:skip()` prevents duplicate results. Each iteration skips the products returned by earlier iterations and selects only the newly loaded batch.
* `pageSize` must match the website’s behavior. The skip offset assumes that every scroll loads eight products. If the site returns variable-sized batches, track the number of already processed elements instead.
* `QUERY ONE` extracts fields from each product. Each title and price lookup is scoped to the current product card rather than the entire page.
* This pattern needs a stopping condition. The loop ends when `SCROLL_BOTTOM(page)` indicates that the page can no longer scroll further. A `LIMIT` can also be added as a safety guard for sites that scroll indefinitely.

## Iterate numbered pages by URL

When pages are addressable by URL (e.g., `?page=1`, `?page=2`), use a `FOR` loop with a range:

Annotated query

{{< editor lang="fql" >}}
// Base URL of the paginated resource.
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/products"
// Total number of pages to visit.
LET totalPages = 5
// Iterate over every page number in the range.
RETURN FOR pageNum IN 1..totalPages
    // Build the URL for the current page.
    // The first page uses the base URL, while subsequent pages
    // follow the site's numbered URL pattern.
    LET url = pageNum == 1 ? baseURL: `${baseURL}/page/${pageNum}`
    // Fetch the page over HTTP.
    // Because each page has its own URL, no browser automation is required.
    LET page = WEB::HTML::OPEN(url)
    // Extract every product from the current page.
    FOR item IN page[~ css`.product-card`]
        RETURN {
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
* Each page is fetched directly. `WEB::HTML::OPEN()` loads every page independently over HTTP, so no browser session or navigation is required.
* Pages are processed independently. Each iteration loads, extracts, and returns data from one page before moving on to the next.
* Browser automation is unnecessary. This approach is faster and more resource-efficient because it avoids rendering pages or simulating user interactions.
* Best suited for static websites. Use this pattern when the desired content is present in the initial HTML returned by the server. If the content is loaded dynamically with JavaScript, use browser automation instead.
* Optional queries tolerate missing data. The `~?` query operator returns `NONE` when an element is missing. Combined with optional chaining (`?.`), incomplete product cards can be processed without failing the query.

### Detect the last page

If you do not know the total number of pages, open the first page to find out:

{{< code lang="fql" >}}
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="
LET firstPage = WEB::HTML::OPEN(baseURL + "1")

LET lastPageLink = QUERY ONE ".pagination a:last-child" IN firstPage USING css
LET totalPages = TO_INT(lastPageLink?.textContent) ON ERROR RETURN 1

RETURN FOR pageNum IN 1..totalPages
    LET page = WEB::HTML::OPEN(baseURL + TO_STRING(pageNum))
    FOR item IN page[~ css`.product-card`]
        RETURN item[~? css`.product-name`]?.textContent
{{</ code >}}

## Collect results into a flat array

When each page returns an array of items, the outer loop produces an array of arrays. Use `FLATTEN` to merge them:

{{< code lang="fql" >}}
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="

LET result = (
    FOR pageNum IN 1..3
        LET page = WEB::HTML::OPEN(baseURL + TO_STRING(pageNum))

        LET items = (
            FOR item IN page[~ css`.product-card`]
                RETURN {
                    name: item[~? css`.product-name`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        RETURN items
)

RETURN FLATTEN(result)
{{</ code >}}

Alternatively, use the `[**]` array contraction operator to flatten inline:

{{< code lang="fql" >}}
LET result = (
    FOR pageNum IN 1..3
        LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/?page=" + TO_STRING(pageNum))
        RETURN page[~ css`.product-card`][*
            RETURN {
                name: .[~? css`.product-name`]?.textContent,
                price: .[~? css`.product-price`]?.textContent
            }
        ]
)

RETURN result[**]
{{</ code >}}

## Add error recovery

Pagination scripts are long-running and may encounter network errors or missing elements. Wrap page loads and interactions with error recovery:

{{< code lang="fql" >}}
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="

LET result = (
    FOR pageNum IN 1..10
        LET page = WEB::HTML::OPEN(baseURL + TO_STRING(pageNum))
            ON ERROR RETRY 2 DELAY 1s BACKOFF EXPONENTIAL
            OR RETURN NONE

        FILTER page != NONE

        LET items = (
            FOR item IN page[~ css`.product-card`]
                RETURN {
                    page: pageNum,
                    name: item[~? css`.product-name`]?.textContent,
                    price: item[~? css`.product-price`]?.textContent
                }
        )

        RETURN items
)

RETURN FLATTEN(result)
{{</ code >}}

See [Error handling and resilience]({{< ref "error-handling" >}}) for more patterns.

## Next steps

{{< docs-related tiles="guide-error-handling,guide-interactions,language-control-flow-for,language-control-flow-error-handling" >}}
