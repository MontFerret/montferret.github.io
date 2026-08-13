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
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

return for frame in page.frames
    return {
        url: frame.URL,
        title: frame.title
    }
{{</ code >}}

## Find a frame by URL

Use `filter` to locate a specific frame:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let target = (
    for frame in page.frames
        filter contains(frame.URL, "embedded-form")
        limit 1
        return frame
)

let frame = first(target)
return frame?.title
{{</ code >}}

## Extract content from a frame

Once you have a frame reference, query it the same way you query a page:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let target = first((
    for frame in page.frames
        filter contains(frame.URL, "embedded")
        limit 1
        return frame
))

let items = target != none
    ? target[~ css`.content-item`]
    : []

return items[*].textContent
{{</ code >}}

## Interact with elements inside a frame

Dispatch events to elements within the frame just as you would on the main page:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let frame = first((
    for f in page.frames
        filter contains(f.URL, "login")
        limit 1
        return f
))

let input = query one "input[name='email']" in frame using css
dispatch "input" in input with "user@example.com"

let submit = query one "button[type='submit']" in frame using css
submit <- "click"

return "submitted"
{{</ code >}}

## Handle missing frames

A frame might not be present on every page. Use error recovery or conditional checks:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let frames = (
    for f in page.frames
        filter contains(f.URL, "target-frame")
        limit 1
        return f
)

return length(frames) > 0
    ? first(frames)[~ css`.data`][*].textContent
    : []
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-interactions,guide-error-handling,language-control-flow-for" >}}
