---
title: "Work with browser-driven pages"
sidebarTitle: "Browser-driven pages"
weight: 20
draft: false
description: "Use a browser to extract data from JavaScript-rendered pages."
---

# Work with browser-driven pages

Some pages load content through JavaScript — single-page applications, lazy-loaded lists, content behind client-side rendering. For these, Ferret can drive a real browser through the Chrome DevTools Protocol (CDP).

This guide covers when to use a browser, how to set one up, and how to wait for dynamic content.

## When to use a browser

Use static extraction (`WEB::HTML::OPEN` without a driver) when the data is in the initial HTML response. Use the `cdp` driver when the page:

- renders content with JavaScript
- loads data asynchronously after the initial page load
- requires user interaction (clicking, scrolling, filling forms) before content appears
- depends on browser APIs (cookies, local storage, service workers)

If you are unsure, try static extraction first. It is faster and simpler. Switch to the browser driver only if the data is missing from the static HTML.

## Set up a browser

Ferret connects to a Chrome or Chromium instance over CDP.

**Option 1: Docker** (recommended for reproducible environments)

{{< terminal command="true" >}}
docker run -d -p 9222:9222 montferret/chromium
{{< /terminal >}}

**Option 2: Managed browser** (the CLI starts and stops a browser for you)

{{< terminal command="true" >}}
ferret browser start
{{< /terminal >}}

**Option 3: Connect to an existing browser** launched with remote debugging enabled.

See [CLI Browser]({{< ref "/docs/tools/cli/browser" >}}) for full details on browser management.

## Open a page with a browser

Pass `{ driver: "cdp" }` to `WEB::HTML::OPEN`:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { 
    driver: "cdp"
})

RETURN page.title
{{< /editor >}}

Once opened, querying and element access work the same as with static pages. The difference is that the page value reflects the live browser DOM, including content added by JavaScript.

## Wait for content

JavaScript-rendered pages may not have all content immediately after the page loads. Use `WAITFOR` to pause until the data appears.

### Wait for an element

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/delayed-rendering/", {
    driver: "cdp"
})

WAITFOR EXISTS QUERY ONE '[data-testid="delayed-long"]' IN page USING css
    TIMEOUT 5s

RETURN page[~ css`[data-testid="delayed-long"]`]
{{</ editor >}}

`WAITFOR EXISTS` re-checks the expression on a polling interval until it has a value or the timeout is reached.

### Wait for a value

Use `WAITFOR VALUE` when you need the result of the check itself:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/dynamic-products/delayed/", { 
    driver: "cdp"
})

LET element = WAITFOR VALUE QUERY ONE '[data-testid="dynamic-product-card"]' IN page USING css
    TIMEOUT 10s
    EVERY 200ms

RETURN element?.textContent
{{</ editor >}}

### Wait for network idle

After a navigation or interaction, you may want to wait until the page finishes loading resources:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", {
    driver: "cdp"
})

WAITFOR EVENT "network.idle" IN page TIMEOUT 10s

RETURN page[~ css`article`]
{{</ editor >}}

### Tune the polling

`WAITFOR` supports several clauses to control timing:

{{< code lang="fql" >}}
WAITFOR EXISTS QUERY ONE ".result" IN page USING css
    TIMEOUT 10s
    EVERY 100ms, 2s
    BACKOFF EXPONENTIAL
    JITTER 0.1
{{</ code >}}

- `EVERY` — how often to re-check (with optional cap)
- `BACKOFF` — how the interval grows (`LINEAR`, `EXPONENTIAL`, or `NONE`)
- `JITTER` — randomize the interval to avoid synchronized retries
- `TIMEOUT` — maximum wait time

## Handle timeouts

When `WAITFOR` times out, it returns `false` (or `NONE` for `WAITFOR VALUE`). Add `ON TIMEOUT RETURN` for a custom fallback:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { 
    driver: "cdp"
})

LET items = WAITFOR VALUE QUERY ".product" IN page USING css
    TIMEOUT 5s
    ON TIMEOUT RETURN []

RETURN items
{{</ editor >}}

## Extract from a browser-backed page

Once the content is loaded, extraction works the same as with static pages:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/dynamic-products/basic/", {
    driver: "cdp"
})

LET grid = WAITFOR VALUE QUERY ONE "#dynamic-products-grid" IN page USING css
WAITFOR EXISTS QUERY '[data-testid="dynamic-product-card"]' IN page USING css

FOR product IN grid.children
    LET title = product[~? css`.product-title`]
    RETURN {
        title: title?.textContent,
        price: product[~? css`.product-price`],
        inStock: product[~? css`.product-in-stock`]?.attributes["data-stock-count"] != "0"
    }
{{< /editor >}}

## Navigate within a page

Use `NAVIGATE` to go to a different URL within the same browser session:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET titleBefore = page.title
LET btn = page[~? css`:nth(1, nav li)`]

WAITFOR EVENT "navigation" IN page
    TRIGGER (
        DISPATCH "click" IN btn
    )
    TIMEOUT 10s

LET titleAfter = page.title

RETURN { titleBefore, titleAfter }
{{</ editor >}}

`NAVIGATE_BACK(page)` and `NAVIGATE_FORWARD(page)` move through the browser history.

## Next steps

{{< docs-related tiles="guide-interactions,guide-pagination,language-control-flow-waitfor,tools-cli-browser" >}}
