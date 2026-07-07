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

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
RETURN page.title
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
RETURN page.title
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

Once opened, querying and element access work the same as with static pages. The difference is that the page value reflects the live browser DOM, including content added by JavaScript.

## Wait for content

JavaScript-rendered pages may not have all content immediately after the page loads. Use `WAITFOR` to pause until the data appears.

### Wait for an element

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

WAITFOR EXISTS QUERY ONE ".dynamic-content" IN page USING css
    TIMEOUT 5s

RETURN page[~ css`.dynamic-content`]
{{</ code >}}

`WAITFOR EXISTS` re-checks the expression on a polling interval until it has a value or the timeout is reached.

### Wait for a value

Use `WAITFOR VALUE` when you need the result of the check itself:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET element = WAITFOR VALUE QUERY ONE ".lazy-item" IN page USING css
    TIMEOUT 10s
    EVERY 200ms

RETURN element?.textContent
{{</ code >}}

### Wait for network idle

After a navigation or interaction, you may want to wait until the page finishes loading resources:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

WAITFOR EVENT "idle" IN page TIMEOUT 10s

RETURN page[~ css`article`]
{{</ code >}}

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

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET items = WAITFOR VALUE QUERY ".product" IN page USING css
    TIMEOUT 5s
    ON TIMEOUT RETURN []

RETURN items
{{</ code >}}

## Extract from a browser-backed page

Once the content is loaded, extraction works the same as with static pages:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET articles = page[~ css`article`]

FOR article IN articles
    LET title = article[~? css`h2`]
    RETURN {
        title: title?.textContent,
        html: article.innerHTML
    }
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET articles = page[~ css`article`]

FOR article IN articles
    LET title = article[~? css`h2`]
    RETURN {
        title: title?.textContent,
        html: article.innerHTML
    }
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

## Navigate within a page

Use `NAVIGATE` to go to a different URL within the same browser session:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET titleBefore = page.title

NAVIGATE(page, "https://mockery.ferretlang.org/scenarios/ecommerce/")
WAIT_NAVIGATION(page)

LET titleAfter = page.title

RETURN { titleBefore, titleAfter }
{{</ code >}}

`NAVIGATE_BACK(page)` and `NAVIGATE_FORWARD(page)` move through the browser history.

## Next steps

{{< docs-related tiles="guide-interactions,guide-pagination,language-control-flow-waitfor,tools-cli-browser" >}}
