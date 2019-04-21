---
title: "Overview"
weight: 2
draft: false
---

# Anatomy of Ferret

Working with Ferret, you write your scrapers using its own declarative query language called FQL (Ferret Query Language).

Let's go through the examples below to get a glimpse into how a Ferret query is structured.

### Everything has a result

Every FQL query must have either a ``RETURN`` or ``FOR`` statement.

{{< highlight sql >}}
RETURN "foobar"
{{< /highlight >}}

{{< highlight sql >}}
FOR i IN 1..10
    RETURN i + 2
{{< /highlight >}}

### Functions are essential

Writing FQL scripts always involves calling functions.
One is the most important one is ``DOCUMENT`` that allows you to open an HTML page and start scraping it.

{{< highlight sql >}}
LET page = DOCUMENT("https://github.com/trending")

FOR row IN ELEMENTS(page, "ol.repo-list li")
    LET name = INNER_TEXT(row, "div:nth-child(1)")
    LET description = INNER_TEXT(row, "div:nth-child(3)")
    
    RETURN { name, description }
{{< /highlight >}}

### Dynamic pages and user events

Previous examples works with static web pages, but nowadays, more and more websites use dynamic page rendering and old plain HTTP GET request is not a solution any more.     
Therefore, there should be possible to handle this kind of web pages.    

And Ferret can handle it, the following query represents it:

{{< highlight sql >}}
LET google = DOCUMENT("https://www.google.com/", { driver: "cdp" })

INPUT(google, 'input[name="q"]', "ferret", 25)
CLICK(google, 'input[name="btnK"]')

WAIT_NAVIGATION(google)
WAIT_ELEMENT(google, '.g', 5000)

FOR result IN ELEMENTS(google, '.g')
    // filter out extra elements like videos and 'People also ask'
    FILTER TRIM(result.attributes.class) == 'g'
    RETURN {
        title: INNER_TEXT(result, 'h3'),
        description: INNER_TEXT(result, '.st'),
        url: INNER_TEXT(result, 'cite')
    }
{{< /highlight >}}