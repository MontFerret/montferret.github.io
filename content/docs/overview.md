---
title: "Overview"
weight: 1
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

Previous examples work with

{{< highlight sql >}}
LET doc = DOCUMENT('https://soundcloud.com/charts/top', true)

WAIT_ELEMENT(doc, '.chartTrack__details', 5000)

LET tracks = ELEMENTS(doc, '.chartTrack__details')

FOR track IN tracks
    RETURN {
        artist: TRIM(INNER_TEXT(track, '.chartTrack__username')),
        track: TRIM(INNER_TEXT(track, '.chartTrack__title'))
    }
{{< /highlight >}}