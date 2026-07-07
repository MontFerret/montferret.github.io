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

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page.title
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page.title
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

The function returns an HTML page value. You can read properties like `title` directly on it.

## Query elements

Use the query shorthand `[~ css` `` ` `` `selector` `` ` `` `]` to find elements:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page[~ css`article h2`]
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page[~ css`article h2`]
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

This returns all matching elements as an array.

To get a single element, use `[~? css` `` ` `` `selector` `` ` `` `]`:

{{< code lang="fql" >}}
LET title = page[~? css`article h2`]
{{</ code >}}

The `~?` shorthand is equivalent to `QUERY ONE`, which returns the first match or `NONE`.

You can also use the full `QUERY` syntax when you need more control:

{{< code lang="fql" >}}
LET items = QUERY "article h2" IN page USING css
LET first = QUERY ONE "article h2" IN page USING css
LET count = QUERY COUNT "article h2" IN page USING css
LET exists = QUERY EXISTS "article h2" IN page USING css
{{</ code >}}

## Extract text and attributes

Once you have an element, read its properties:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET links = page[~ css`a`]

FOR link IN links
    LIMIT 5
    RETURN {
        text: link.textContent,
        href: link.attributes.href
    }
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET links = page[~ css`a`]

FOR link IN links
    LIMIT 5
    RETURN {
        text: link.textContent,
        href: link.attributes.href
    }
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

Common element properties:

| Property | Description |
| --- | --- |
| `textContent` | The text content of the element |
| `innerHTML` | The inner HTML of the element |
| `attributes` | An object of attribute key-value pairs |
| `attributes.href` | A specific attribute value |

## Use array operators for compact extraction

The `[*]` array operator lets you project fields from a list of elements without writing a `FOR` loop:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page[~ css`a`][*].attributes.href
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page[~ css`a`][*].attributes.href
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

You can also filter inline:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
RETURN page[~ css`a`][*
    FILTER .attributes.href != NONE
    RETURN {
        text: .textContent,
        href: .attributes.href
    }
]
{{</ code >}}

## Query nested elements

When a page has repeating structures — product cards, table rows, list items — query the container first, then query inside each one:

{{< tabs >}}
{{< tab title="Terminal" >}}
{{< terminal command="true" >}}
ferret run -e '
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
LET cards = page[~ css`.product-card`]

FOR card IN cards
    RETURN {
        name: card[~? css`.product-name`]?.textContent,
        price: card[~? css`.product-price`]?.textContent
    }
'
{{< /terminal >}}
{{< /tab >}}

{{< tab title="Try in browser" >}}
{{< editor lang="fql" height="auto" copy="true" apiVersion="2" orientation="horizontal" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
LET cards = page[~ css`.product-card`]

FOR card IN cards
    RETURN {
        name: card[~? css`.product-name`]?.textContent,
        price: card[~? css`.product-price`]?.textContent
    }
{{< /editor >}}
{{< /tab >}}
{{< /tabs >}}

The `?.` optional chaining operator returns `NONE` instead of failing when an element is not found. This keeps the script running even when some cards are missing a field.

## Handle missing elements

Not every page has the elements you expect. Use `QUERY EXISTS` to check before extracting, or `ON ERROR RETURN` to provide a fallback:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")

LET title = QUERY ONE ".page-title" IN page USING css
    ON ERROR RETURN NONE

LET hasNav = QUERY EXISTS "nav" IN page USING css

RETURN {
    title: title?.textContent,
    hasNav
}
{{</ code >}}

For more error handling patterns, see [Error handling and resilience]({{< ref "error-handling" >}}).

## Filter and sort results

Use `FILTER`, `SORT`, and `LIMIT` inside a `FOR` loop to shape the output:

{{< code lang="fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/ecommerce/")
LET cards = page[~ css`.product-card`]

FOR card IN cards
    LET name = card[~? css`.product-name`]?.textContent
    LET price = card[~? css`.product-price`]?.textContent

    FILTER name != NONE
    SORT name ASC
    LIMIT 10

    RETURN { name, price }
{{</ code >}}

## Use parameters for reusable scripts

Save a script to a file and pass the URL as a parameter:

{{< terminal command="true" >}}
echo '
LET page = WEB::HTML::OPEN(@url)
RETURN page[~ css`h1, h2, h3`][*].textContent
' > headings.fql
{{< /terminal >}}

Run it with any URL:

{{< terminal command="true" >}}
ferret run headings.fql --param url=https://mockery.ferretlang.org
{{< /terminal >}}

## Next steps

{{< docs-related tiles="guide-browser-pages,guide-error-handling,language-control-flow-query,stdlib" >}}
