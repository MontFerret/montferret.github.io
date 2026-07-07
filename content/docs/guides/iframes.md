---
title: "Work with iframes"
sidebarTitle: "Iframes"
weight: 100
draft: false
description: "Access and extract data from content inside iframes."
---

# Work with iframes

Some pages embed content in iframes — login forms, embedded widgets, third-party content. This guide shows how to access iframe content with Ferret.

Iframe access requires the `cdp` driver. See [Browser-driven pages]({{< ref "browser-pages" >}}) for setup.

{{< notification type="warning" >}}
Cross-origin iframes restrict access to their content. The iframe must be on the same domain as the parent page, or the browser will block property and content access.
{{</ notification >}}

## List frames on a page

A browser-backed page exposes its frames through the `frames` property:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

FOR frame IN page.frames
    RETURN {
        url: frame.URL,
        title: frame.title
    }
{{</ code >}}

## Find a frame by URL

Use `FILTER` to locate a specific frame:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET target = (
    FOR frame IN page.frames
        FILTER CONTAINS(frame.URL, "embedded-form")
        LIMIT 1
        RETURN frame
)

LET frame = FIRST(target)
RETURN frame?.title
{{</ code >}}

## Extract content from a frame

Once you have a frame reference, query it the same way you query a page:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET target = FIRST((
    FOR frame IN page.frames
        FILTER CONTAINS(frame.URL, "embedded")
        LIMIT 1
        RETURN frame
))

LET items = target != NONE
    ? target[~ css`.content-item`]
    : []

RETURN items[*].textContent
{{</ code >}}

## Interact with elements inside a frame

Dispatch events to elements within the frame just as you would on the main page:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET frame = FIRST((
    FOR f IN page.frames
        FILTER CONTAINS(f.URL, "login")
        LIMIT 1
        RETURN f
))

LET input = QUERY ONE "input[name='email']" IN frame USING css
DISPATCH "input" IN input WITH "user@example.com"

LET submit = QUERY ONE "button[type='submit']" IN frame USING css
submit <- "click"

RETURN "submitted"
{{</ code >}}

## Handle missing frames

A frame might not be present on every page. Use error recovery or conditional checks:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET frames = (
    FOR f IN page.frames
        FILTER CONTAINS(f.URL, "target-frame")
        LIMIT 1
        RETURN f
)

RETURN LENGTH(frames) > 0
    ? FIRST(frames)[~ css`.data`][*].textContent
    : []
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-interactions,guide-error-handling,language-control-flow-for" >}}
