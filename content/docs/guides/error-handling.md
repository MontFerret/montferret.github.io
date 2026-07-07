---
title: "Error handling and resilience"
sidebarTitle: "Error handling"
weight: 50
draft: false
description: "Handle missing elements, network failures, and timeouts in extraction scripts."
---

# Error handling and resilience

Real-world pages are unpredictable. Elements may be missing, networks may be slow, and page structures may change. This guide shows practical patterns for writing extraction scripts that handle these failures gracefully.

For the full syntax reference, see [Error Handling]({{< ref "/docs/language/control-flow/error-handling" >}}).

## Return a fallback value

The most common pattern: attach `ON ERROR RETURN` to provide a default when an expression fails.

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET title = QUERY ONE ".nonexistent" IN page USING css ON ERROR RETURN NONE
RETURN title?.textContent
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET title = QUERY ONE ".nonexistent" IN page USING css ON ERROR RETURN NONE
RETURN title?.textContent
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

Common fallback values:

| Fallback | When to use |
| --- | --- |
| `NONE` | Single optional value |
| `[]` | Expected list that may be empty |
| `{}` | Expected object with no data |
| `"unknown"` | Display-safe placeholder |

## Optional chaining

The `?.` operator accesses properties on values that may be `NONE`. Instead of failing, it returns `NONE`:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET element = QUERY ONE ".maybe-missing" IN page USING css ON ERROR RETURN NONE

// Safe — returns NONE if element is NONE
LET text = element?.textContent
LET href = element?.attributes?.href

RETURN { text, href }
{{</ code >}}

Without `?.`, accessing a property on `NONE` is a runtime error.

## Retry on failure

Use `ON ERROR RETRY` when a transient failure may resolve on a second attempt:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    ON ERROR RETRY 3
{{</ code >}}

### Add a delay between retries

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    ON ERROR RETRY 3 DELAY 500ms
{{</ code >}}

### Use exponential backoff

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    ON ERROR RETRY 3 DELAY 200ms BACKOFF EXPONENTIAL
{{</ code >}}

The delay doubles on each retry: 200ms, 400ms, 800ms.

### Fall back after all retries fail

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    ON ERROR RETRY 3 DELAY 500ms BACKOFF EXPONENTIAL
    OR RETURN NONE

RETURN page != NONE ? page.title : "page unavailable"
{{</ code >}}

## Handle timeouts

`WAITFOR` expressions support `ON TIMEOUT RETURN` for a timeout-specific fallback. This is separate from `ON ERROR`:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET result = WAITFOR VALUE QUERY ONE ".slow-loading" IN page USING css
    TIMEOUT 5s
    ON TIMEOUT RETURN NONE
    ON ERROR RETURN NONE

RETURN result?.textContent
{{</ code >}}

## Extract with fallback selectors

When a site changes its layout, old selectors may stop working. Try multiple selectors and use the first one that matches:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

LET title = QUERY ONE ".new-title" IN page USING css ON ERROR RETURN NONE
LET titleFallback = title == NONE
    ? (QUERY ONE ".old-title" IN page USING css ON ERROR RETURN NONE)
    : title

RETURN titleFallback?.textContent
{{</ code >}}

Or with a function for reuse:

{{< code lang="fql" >}}
FUNC queryFirst(page, selectors) {
    FOR selector IN selectors
        LET result = QUERY ONE selector IN page USING css ON ERROR RETURN NONE
        FILTER result != NONE
        LIMIT 1
        RETURN result
}

LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET title = FIRST(queryFirst(page, [".new-title", ".old-title", "h1"]))

RETURN title?.textContent
{{</ code >}}

## Protect a loop from individual failures

When iterating over items, a single failure in one item should not stop the entire extraction. Use `ON ERROR RETURN` on individual operations, or wrap the whole loop body:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET items = page[~ css`article`]

FOR item IN items
    LET title = item[~? css`h2`]?.textContent
    LET link = item[~? css`a`]?.attributes?.href
    LET description = item[~? css`p`]?.textContent

    RETURN {
        title: title != NONE ? title : "untitled",
        link,
        description
    }
{{</ code >}}

To catch and skip items that fail entirely, wrap the loop body in a grouped expression:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET items = page[~ css`article`]

FOR item IN items
    LET result = ({
        title: item[~? css`h2`]?.textContent,
        link: item[~? css`a`]?.attributes?.href
    }) ON ERROR RETURN NONE

    FILTER result != NONE
    RETURN result
{{</ code >}}

## Combine error and timeout handling

A common pattern for browser-backed extraction: retry on error, provide a timeout fallback, and use optional chaining throughout:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
    ON ERROR RETRY 2 DELAY 1s BACKOFF EXPONENTIAL
    OR RETURN NONE

LET loaded = page != NONE
    ? (WAITFOR EXISTS QUERY ONE ".content" IN page USING css
        TIMEOUT 10s
        ON TIMEOUT RETURN FALSE)
    : FALSE

RETURN loaded ? {
    title: page?.title,
    items: page[~ css`.content .item`][*
        RETURN {
            name: .?.textContent
        }
    ]
} : { error: "page unavailable" }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-static-pages,guide-browser-pages,language-control-flow-error-handling,language-control-flow-waitfor" >}}
