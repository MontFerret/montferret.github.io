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

## List frames on a page

Open the [Mockery iframe scenario](https://mockery.ferretlang.org/scenarios/iframes/) with the browser driver. The page’s `frames` collection exposes documents, including the main document. Use their URLs to identify them; frame ordering can vary.

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/iframes/", { driver: "cdp" })
return for frame in page.frames
    return { url: frame.URL, title: frame.title }
{{< /code >}}

## Extract content from a frame

Frame A contains a product card. Find its document by URL and run the CSS query against that document. Polling allows time for the browser to load it:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org/scenarios/iframes/", { driver: "cdp" })
let frame = waitfor value arrays::first((
    for candidate in page.frames
        filter contains(candidate.URL, "/iframes/frame-a/")
        return candidate
)) timeout 5s on timeout fail

return frame[~? css`[data-testid="product-title"]`]?.textContent
{{< /code >}}

These examples use same-origin frames. If a frame is optional, use `ON TIMEOUT RETURN NONE` and check the value before querying its document.

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-interactions,guide-error-handling,language-control-flow-for" >}}
