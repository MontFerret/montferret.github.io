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

The most common pattern: click a "next" link until it disappears.

Use `FOR DO WHILE` so the body runs at least once (for the first page), and check whether the "next" button still exists after each iteration:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/", {
    driver: "cdp"
})

LET nextSelector = "a.next-page"
LET itemSelector = ".product-card"

FOR i DO WHILE QUERY EXISTS nextSelector IN page USING css
    LIMIT 5

    // On subsequent pages, click "next" and wait for navigation
    LET _ = i > 0 ? (
        LET btn = QUERY ONE nextSelector IN page USING css
        btn <- "click"
        WAITFOR EVENT "navigation" IN page TIMEOUT 10s
    ) : NONE

    FOR item IN QUERY itemSelector IN page USING css
        RETURN {
            name: item[~? css`.product-name`]?.textContent,
            price: item[~? css`.product-price`]?.textContent
        }
{{</ code >}}

Key points:
- `DO WHILE` ensures the first page is always processed.
- `LIMIT 5` caps the number of pages (remove for all pages, but add a reasonable limit to avoid infinite loops).
- The conditional `i > 0 ?` skips clicking on the first iteration.

## Iterate numbered pages by URL

When pages are addressable by URL (e.g., `?page=1`, `?page=2`), use a `FOR` loop with a range:

{{< code lang="fql" >}}
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="
LET totalPages = 5

FOR pageNum IN 1..totalPages
    LET page = WEB::HTML::OPEN(baseURL + TO_STRING(pageNum))

    FOR item IN page[~ css`.product-card`]
        RETURN {
            page: pageNum,
            name: item[~? css`.product-name`]?.textContent,
            price: item[~? css`.product-price`]?.textContent
        }
{{</ code >}}

This approach works without a browser since each page is fetched independently via HTTP. Use static extraction when the page content is in the initial HTML.

### Detect the last page

If you do not know the total number of pages, open the first page to find out:

{{< code lang="fql" >}}
LET baseURL = "https://mockery.ferretlang.org/scenarios/ecommerce/?page="
LET firstPage = WEB::HTML::OPEN(baseURL + "1")

LET lastPageLink = QUERY ONE ".pagination a:last-child" IN firstPage USING css
LET totalPages = TO_INT(lastPageLink?.textContent) ON ERROR RETURN 1

FOR pageNum IN 1..totalPages
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
