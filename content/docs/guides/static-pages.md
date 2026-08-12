---
title: "Extract data from static pages"
sidebarTitle: "Static pages"
weight: 10
draft: false
description: "Load HTML pages, query elements with CSS selectors, and shape the output."
---

# Extract data from static pages

This guide shows how to load an HTML page, find elements with CSS selectors, extract text and attributes, and return structured data.

Static extraction does not need a browser. Ferret fetches the HTML over HTTP and parses it in memory, which is fast and lightweight. Use this approach whenever the page content is present in the initial HTML response. For pages that load content through JavaScript, see [Browser-driven pages]({{< ref "browser-pages" >}}).

## Open a page

Use `WEB::HTML::OPEN` to fetch and parse an HTML page:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
return page.title
{{< /editor >}}

The function returns an HTML page value. You can read properties like `title` directly on it.

## Query elements

Use the query expression to find elements:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

return query "article h2" in page using css
{{< /editor >}}

This returns all matching elements as an array.

To get a single element, use `query one`:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

return query one "article h2" in page using css
{{< /editor >}}

More about query expressions [see the documentation]({{< ref "/docs/language/control-flow/query" >}}).

## Extract text and attributes

Once you have an element, read its properties:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

return for link in (query "a" in page using css)
    limit 5
    return {
        text: link.textContent,
        href: link.attributes.href
    }
{{< /editor >}}

Common element properties:

| Property | Description |
| --- | --- |
| `textContent` | The text content of the element |
| `innerHTML` | The inner HTML of the element |
| `attributes` | An object of attribute key-value pairs |
| `attributes.href` | A specific attribute value |

## Use array operators for compact extraction

The `[*]` array operator lets you project fields from a list of elements without writing a `for` loop:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let links = query "a" in page using css

return links[*].attributes.href
{{< /editor >}}

You can also filter inline:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
let links = query "a" in page using css

return links[*
    filter .attributes.href != none
    return {
        text: .textContent,
        href: .attributes.href
    }
]
{{< /editor >}}

## Query nested elements

When a page has repeating structures — product cards, table rows, list items — query the container first, then query inside each one:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
let cards = query ".product-card" in page using css

return for card in cards
    return {
        name: (query one query ".product-name" in card using css)?.textContent,
        price: (query one query ".product-price" in card using css)?.textContent
    }
{{< /editor >}}

> NOTE: For simple queries, you can use the shortcut query syntax. For details and limitations, see [Shortcut syntax]({{< ref "/docs/language/control-flow/query#shortcut-syntax" >}}).

The `?.` optional chaining operator returns `none` instead of failing when an element is not found. This keeps the script running even when some cards are missing a field.

## Handle missing elements

Not every page has the elements you expect. Use `query exists` to check before extracting, or `on error return` to provide a fallback:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

let title = query one ".page-title" in page using css
    on error return none

let hasNav = query exists "nav" in page using css

return {
    title: title?.textContent,
    hasNav
}
{{</ editor >}}

For more error handling patterns, see [Error handling and resilience]({{< ref "error-handling" >}}).

## Filter and sort results

Use `filter`, `sort`, and `limit` inside a `for` loop to shape the output:

{{< editor lang="fql" >}}
let page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
let cards = query '.product-card' in page using css

return for card in cards
    let title = (query one '.product-title' in card using css)?.textContent
    let price = (query one '.product-price' in card using css)?.textContent

    filter title != none
    sort title asc
    limit 10

    return { title, price }
{{</ editor >}}

## Use parameters for reusable scripts

Save a script to a file and pass the URL as a parameter:

{{< terminal command="true" >}}
echo '
let page = WEB::HTML::OPEN(@url)
let headers = query 'h1, h2, h3' in page using css
return headers[*].textContent
' > headings.fql
{{< /terminal >}}

Run it with any URL:

{{< tabs >}}
{{< tab title="Terminal" >}}

{{< terminal command="true" >}}
ferret run headings.fql --param url=https://mockery.ferretlang.org
{{< /terminal >}}
{{</ tab >}}
{{< tab title="Try in browser" >}}

{{< editor lang="fql" params=`{ "url": "https://mockery.ferretlang.org/"}` >}}
let page = WEB::HTML::OPEN(@url)
let headers = query 'h1, h2, h3' in page using css

return headers[*].textContent
{{</ editor >}}
{{</ tab >}}
{{</ tabs >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-error-handling,language-control-flow-query,stdlib" >}}
