---
title: "Screenshots and PDFs"
sidebarTitle: "Screenshots & PDFs"
weight: 90
draft: false
description: "Capture screenshots and generate PDFs from web pages."
---

# Screenshots and PDFs

Ferret can capture screenshots and generate PDF documents from browser-backed pages. Both require the `cdp` driver.

## Take a screenshot

Use `web::html::screenshot` to capture the visible page:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
let data = web::html::screenshot(page)
return data
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
let data = web::html::screenshot(page)
return data
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

`web::html::screenshot` returns binary data (base64-encoded PNG). To save it to a file, use `io::fs::write`:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
let data = web::html::screenshot(page)
io::fs::write("screenshot.png", data)

return "saved"
{{</ code >}}

## Screenshot options

Pass options to control the output:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let data = web::html::screenshot(page, {
    format: "jpeg",
    quality: 80,
    fullPage: true
})

io::fs::write("full-page.jpg", data)
return "saved"
{{</ code >}}

| Option | Type | Description |
| --- | --- | --- |
| `format` | string | `"png"` (default) or `"jpeg"` |
| `quality` | int | JPEG quality, 0–100 (ignored for PNG) |
| `fullPage` | bool | Capture the full scrollable page, not just the viewport |

## Generate a PDF

Use `web::html::pdf` to produce a PDF document from the page:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })
let data = web::html::pdf(page)
io::fs::write("page.pdf", data)

return "saved"
{{</ code >}}

## PDF options

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

let data = web::html::pdf(page, {
    landscape: true,
    printBackground: true,
    paperWidth: 8.5,
    paperHeight: 11
})

io::fs::write("report.pdf", data)
return "saved"
{{</ code >}}

| Option | Type | Description |
| --- | --- | --- |
| `landscape` | bool | Landscape orientation |
| `printBackground` | bool | Include background graphics |
| `paperWidth` | float | Page width in inches |
| `paperHeight` | float | Page height in inches |
| `marginTop`, `marginBottom`, `marginLeft`, `marginRight` | float | Margins in inches |

## Wait before capturing

Dynamic pages may need time to render. Wait for the content to stabilize before capturing:

{{< code lang="fql" >}}
let page = web::html::open("https://mockery.ferretlang.org", { driver: "cdp" })

waitfor exists query one ".fully-loaded" in page using css
    timeout 10s

let data = web::html::screenshot(page, { fullPage: true })
io::fs::write("loaded.png", data)

return "saved"
{{</ code >}}

## Capture multiple pages

Combine screenshots with pagination or navigation:

{{< code lang="fql" >}}
let urls = [
    "https://mockery.ferretlang.org",
    "https://mockery.ferretlang.org/scenarios/ecommerce/"
]

return for url, i in urls
    let page = web::html::open(url, { driver: "cdp" })
    let data = web::html::screenshot(page, { fullPage: true })
    let filename = "screenshot-" + to_string(i) + ".png"
    io::fs::write(filename, data)

    return { url, filename }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-interactions,tools-cli-browser,stdlib" >}}
