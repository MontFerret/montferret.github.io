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

The most common pattern: attach `on error return` to provide a default when an expression fails.

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let title = query one ".nonexistent" in page using css on error return none
return title?.textContent
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let title = query one ".nonexistent" in page using css on error return none
return title?.textContent
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

Common fallback values:

| Fallback | When to use |
| --- | --- |
| `none` | Single optional value |
| `[]` | Expected list that may be empty |
| `{}` | Expected object with no data |
| `"unknown"` | Display-safe placeholder |

## Optional chaining

The `?.` operator accesses properties on values that may be `none`. Instead of failing, it returns `none`:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let element = query one ".maybe-missing" in page using css on error return none

// Safe — returns NONE if element is NONE
let text = element?.textContent
let href = element?.attributes?.href

return { text, href }
{{</ code >}}

Without `?.`, accessing a property on `none` is a runtime error.

## Provide a value when data is absent

Combine optional chaining with `??` when missing data should become a concrete value:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let element = query one ".maybe-missing" in page using css on error return none

return element?.textContent ?? "Unknown"
{{</ code >}}

This keeps `false`, zero, and empty strings unchanged. `??` is not an error handler: the query uses `on error return none` first, optional chaining safely reads the member, and then `??` replaces only the resulting `none`.

See [none-Coalescing Operator]({{< ref "/docs/language/operators/coalescing" >}}) for details.

## Retry on failure

Use `on error retry` when a transient failure may resolve on a second attempt:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    on error retry 3
{{</ code >}}

### Add a delay between retries

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    on error retry 3 delay 500ms
{{</ code >}}

### Use exponential backoff

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    on error retry 3 delay 200ms backoff EXPONENTIAL
{{</ code >}}

The delay doubles on each retry: 200ms, 400ms, 800ms.

### Fall back after all retries fail

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
    on error retry 3 delay 500ms backoff EXPONENTIAL
    or return none

return page != none ? page.title : "page unavailable"
{{</ code >}}

## Handle timeouts

`waitfor` expressions support `on timeout return` for a timeout-specific fallback. This is separate from `on error`:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

let result = waitfor value query one ".slow-loading" in page using css
    timeout 5s
    on timeout return none
    on error return none

return result?.textContent
{{</ code >}}

## Extract with fallback selectors

When a site changes its layout, old selectors may stop working. Try multiple selectors and use the first one that matches:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

let title = query one ".new-title" in page using css on error return none
let titleFallback = title == none
    ? (query one ".old-title" in page using css on error return none)
    : title

return titleFallback?.textContent
{{</ code >}}

Or with a function for reuse:

{{< code lang="fql" >}}
func queryFirst(page, selectors) {
    return for selector in selectors
        let result = query one selector in page using css on error return none
        filter result != none
        limit 1
        return result
}

let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let title = FIRST(queryFirst(page, [".new-title", ".old-title", "h1"]))

return title?.textContent
{{</ code >}}

## Protect a loop from individual failures

When iterating over items, a single failure in one item should not stop the entire extraction. Use `on error return` on individual operations, or wrap the whole loop body:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let items = page[~ css`article`]

return for item in items
    let title = item[~? css`h2`]?.textContent
    let link = item[~? css`a`]?.attributes?.href
    let description = item[~? css`p`]?.textContent

    return {
        title: title != none ? title : "untitled",
        link,
        description
    }
{{</ code >}}

To catch and skip items that fail entirely, wrap the loop body in a grouped expression:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let items = page[~ css`article`]

return for item in items
    let result = ({
        title: item[~? css`h2`]?.textContent,
        link: item[~? css`a`]?.attributes?.href
    }) on error return none

    filter result != none
    return result
{{</ code >}}

## Combine error and timeout handling

A common pattern for browser-backed extraction: retry on error, provide a timeout fallback, and use optional chaining throughout:

{{< code lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
    on error retry 2 delay 1s backoff EXPONENTIAL
    or return none

let loaded = page != none
    ? (waitfor exists query one ".content" in page using css
        timeout 10s
        on timeout return false)
    : false

return loaded ? {
    title: page?.title,
    items: page[~ css`.content .item`][*
        return {
            name: .?.textContent
        }
    ]
} : { error: "page unavailable" }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-static-pages,guide-browser-pages,language-control-flow-error-handling,language-control-flow-waitfor" >}}
