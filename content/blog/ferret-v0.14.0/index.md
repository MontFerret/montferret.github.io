---
title: "Ferret v0.14.0"
subtitle: "Happy new year!"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2021-03-06"
---

Hello friends,

Belated happy new 2021 year! Hope things will be better and we can finally get back to normal life.

Meanwhile, in between of work, personal lives and anxiety, we've managed to bring you some great features with a new release of Ferret - **[Ferret v0.14.0](https://github.com/MontFerret/ferret/releases/tag/v0.14.0)**.  

This release contains some syntax updates, DOM API fixes and extra flexibility. Let's dive in!

# What's added

## Support of History API

Before v0.14.0, Page object from the CDP driver, would only return the url that was set during **full** page load or redirect but ignored all **in page** navigation using History API (e.g. using ``react-router``). 

Starting this relase, the behavior has changed. Now the Page object always returns the url which is in the url-bar of your browser. If you need previous behavior, access url property of the root document.

By other words, ``page.url`` and ``document.url`` may point to different locations if History API is used on your page.

{{< editor lang="fql" height="300px" readonly="true" >}}
LET page = DOCUMENT("https://soundcloud.com", { driver: "cdp"})
LET doc = page.mainFrame

WAIT_ELEMENT(doc, ".trendingTracks")
SCROLL_ELEMENT(doc, ".trendingTracks")
WAIT_ELEMENT(doc, ".trendingTracks .badgeList__item")

LET song = ELEMENT(doc, ".trendingTracks .badgeList__item")
CLICK(song)

WAIT_ELEMENT(doc, ".l-listen-hero")

RETURN {
    page: page.url,
    doc: doc.url
}
{{</ editor >}}

## Support of custom HTTP transport in the HTTP driver

Now it's possible to provide your custom HTTP transport to the underlying HTTP client ([pester](https://github.com/sethgrid/pester)) in HTTP driver:

{{< code lang="golang" height="250px" readonly="true" >}}
import (
    h "net/http"
    "github.com/MontFerret/ferret/pkg/drivers/http"
)

func main() {
    httpDriver := http.NewDriver(
		http.WithCustomTransport(&h.Transport{}),
	)
}

{{</ code >}}

## Added LIKE operator
It's been a part of ArrangoDB Query Language since the beginning, but was not ported to FQL in its early times.
And now, ``LIKE`` operator has finally landed!

{{< editor lang="fql" height="300px" readonly="false" >}}
LET values = (
    FOR str IN ["foo", "bar", "qaz"]
				FILTER str LIKE "*a*"
				RETURN str 
)

RETURN FIRST(values) NOT LIKE "b**" ? 'failure' : 'success' 
{{</ editor >}}

As you might noticed, FQL's implementation has some deviations from its ArrangoDB counterpart: syntax of pattern matching. Ferret is using [stadard Unix wildcards](http://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm).

## Support of ignoring page resources
Now you can ignore all the "noise" loaded with your web pages and speed up your scraping by disabling particular files from loading:

{{< editor lang="fql" height="300px" readonly="true" >}}
LET p = DOCUMENT("https://www.gettyimages.com/", {
    driver: "cdp",
    ignore: {
        resources: [
            {
                url: "*",
                type: "image"
            }
        ]
    }
})

RETURN NONE

{{</ editor >}}

You can either use just a type of the resource (stylesheet, script, image, font and etc) or add an url pattern to do it more selectevley.

## Support of handling non-200 HTTP status codes in the HTTP driver
This possibility added to the HTTP driver that allows you to handle situation when a target website responds with HTTP code other than 200 but still has content.

{{< editor lang="fql" height="300px" readonly="false" >}}
LET p = DOCUMENT("https://www.gettyimages.com/", {
    ignore: {
        statusCodes: [
            {
                url: "*",
                code: 418
            }
        ]
    }
})

RETURN NONE

{{</ editor >}}

As with ignoring resources, you can either use just a code or add an url pattern to do it more selectevley.

## DOCUMENT_EXISTS function

Now you can check if a webpage exists before navigating to it.

{{< editor lang="fql" height="80px" readonly="false" >}}
RETURN DOCUMENT_EXISTS("www.asdadasda.sds")
{{</ editor >}}

# What's fixed

Ok, now, let's take a look what bugs have been fixed!

- RAND(0,100) always same result
- Element.children always returns empty array
- Passing parameters with a nested nil structure leads to panic

# Summary
Thanks everyone who contributed to the project by either providing new feature or fixing bugs or just asking questions!