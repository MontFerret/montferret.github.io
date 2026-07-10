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
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page.title
{{< /editor >}}

The function returns an HTML page value. You can read properties like `title` directly on it.

## Query elements

Use the query expression to find elements:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

RETURN QUERY "article h2" IN page USING css
{{< /editor >}}

This returns all matching elements as an array.

To get a single element, use `QUERY ONE`:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

RETURN QUERY ONE "article h2" IN page USING css
{{< /editor >}}

More about query expressions [see the documentation]({{< ref "/docs/language/control-flow/query" >}}).

## Extract text and attributes

Once you have an element, read its properties:

{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

FOR link IN (QUERY "a" IN page USING css)
    LIMIT 5
    RETURN {
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

The `[*]` array operator lets you project fields from a list of elements without writing a `FOR` loop:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET links = QUERY "a" IN page USING css

RETURN links[*].attributes.href
{{< /editor >}}

You can also filter inline:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET links = QUERY "a" IN page USING css

RETURN links[*
    FILTER .attributes.href != NONE
    RETURN {
        text: .textContent,
        href: .attributes.href
    }
]
{{< /editor >}}

## Query nested elements

When a page has repeating structures — product cards, table rows, list items — query the container first, then query inside each one:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
LET cards = QUERY ".product-card" IN page USING css

FOR card IN cards
    RETURN {
        name: (QUERY ONE QUERY ".product-name" IN card USING css)?.textContent,
        price: (QUERY ONE QUERY ".product-price" IN card USING css)?.textContent
    }
{{< /editor >}}

> NOTE: For simple queries, you can use the shortcut query syntax. For details and limitations, see [Shortcut syntax]({{< ref "/docs/language/control-flow/query#shortcut-syntax" >}}).

The `?.` optional chaining operator returns `NONE` instead of failing when an element is not found. This keeps the script running even when some cards are missing a field.

## Handle missing elements

Not every page has the elements you expect. Use `QUERY EXISTS` to check before extracting, or `ON ERROR RETURN` to provide a fallback:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

LET title = QUERY ONE ".page-title" IN page USING css
    ON ERROR RETURN NONE

LET hasNav = QUERY EXISTS "nav" IN page USING css

RETURN {
    title: title?.textContent,
    hasNav
}
{{</ editor >}}

For more error handling patterns, see [Error handling and resilience]({{< ref "error-handling" >}}).

## Filter and sort results

Use `FILTER`, `SORT`, and `LIMIT` inside a `FOR` loop to shape the output:

{{< editor lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
LET cards = QUERY '.product-card' IN page USING css

FOR card IN cards
    LET title = (QUERY ONE '.product-title' IN card USING css)?.textContent
    LET price = (QUERY ONE '.product-price' IN card USING css)?.textContent

    FILTER title != NONE
    SORT title ASC
    LIMIT 10

    RETURN { title, price }
{{</ editor >}}

## Use parameters for reusable scripts

Save a script to a file and pass the URL as a parameter:

{{< terminal command="true" >}}
echo '
LET page = WEB::HTML::OPEN(@url)
LET headers = QUERY 'h1, h2, h3' IN page USING css
RETURN headers[*].textContent
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
LET page = WEB::HTML::OPEN(@url)
LET headers = QUERY 'h1, h2, h3' IN page USING css

RETURN headers[*].textContent
{{</ editor >}}
{{</ tab >}}
{{</ tabs >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-error-handling,language-control-flow-query,stdlib" >}}
