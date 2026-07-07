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

Use `SCREENSHOT` to capture the visible page:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET data = SCREENSHOT(page)
RETURN data
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET data = SCREENSHOT(page)
RETURN data
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

`SCREENSHOT` returns binary data (base64-encoded PNG). To save it to a file, use `IO::FS::WRITE`:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET data = SCREENSHOT(page)
IO::FS::WRITE("screenshot.png", data)

RETURN "saved"
{{</ code >}}

## Screenshot options

Pass options to control the output:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET data = SCREENSHOT(page, {
    format: "jpeg",
    quality: 80,
    fullPage: TRUE
})

IO::FS::WRITE("full-page.jpg", data)
RETURN "saved"
{{</ code >}}

| Option | Type | Description |
| --- | --- | --- |
| `format` | string | `"png"` (default) or `"jpeg"` |
| `quality` | int | JPEG quality, 0–100 (ignored for PNG) |
| `fullPage` | bool | Capture the full scrollable page, not just the viewport |

## Generate a PDF

Use `PDF` to produce a PDF document from the page:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })
LET data = PDF(page)
IO::FS::WRITE("page.pdf", data)

RETURN "saved"
{{</ code >}}

## PDF options

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

LET data = PDF(page, {
    landscape: TRUE,
    printBackground: TRUE,
    paperWidth: 8.5,
    paperHeight: 11
})

IO::FS::WRITE("report.pdf", data)
RETURN "saved"
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
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org", { driver: "cdp" })

WAITFOR EXISTS QUERY ONE ".fully-loaded" IN page USING css
    TIMEOUT 10s

LET data = SCREENSHOT(page, { fullPage: TRUE })
IO::FS::WRITE("loaded.png", data)

RETURN "saved"
{{</ code >}}

## Capture multiple pages

Combine screenshots with pagination or navigation:

{{< code lang="fql" >}}
LET urls = [
    "https://mockery.ferretlang.org",
    "https://mockery.ferretlang.org/scenarios/ecommerce/"
]

FOR url, i IN urls
    LET page = WEB::HTML::OPEN(url, { driver: "cdp" })
    LET data = SCREENSHOT(page, { fullPage: TRUE })
    LET filename = "screenshot-" + TO_STRING(i) + ".png"
    IO::FS::WRITE(filename, data)

    RETURN { url, filename }
{{</ code >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-interactions,tools-cli-browser,stdlib" >}}
