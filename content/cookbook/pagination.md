---
title: "Pagination"
weight: 2
draft: false
---

There are several way how to implement pagination.

### WHILE loop
Since v0.13.0, pagination can be implemented with ``for-while`` loop.

{{< notification type="info">}}
For paginations, it's recommended to use DO-WHILE variation of the loop, in order to process at least the first page.
{{</ notification >}}

{{< editor height="300px" >}}
LET doc = DOCUMENT("https://github.com/MontFerret/ferret/stargazers", { driver: "cdp" })

LET nextSelector = ".paginate-container .BtnGroup a:nth-child(2)"
LET elementsSelector = '.follow-list li'

FOR i DO WHILE ELEMENT_EXISTS(doc, nextSelector)
    LIMIT 3
	LET wait = i > 0 ? CLICK(doc, nextSelector) : false
	LET nav = wait ? WAIT_NAVIGATION(doc) : false
	
	FOR el IN ELEMENTS(doc, elementsSelector)
		FILTER ELEMENT_EXISTS(el, ".octicon-organization")

		RETURN {
			name: INNER_TEXT(el, ".follow-list-name"),
			company: INNER_TEXT(el, ".follow-list-info span")
		}
{{</ editor >}}

### Controlled
You can also use ``for-in`` loop with specified range of iterations that can be either fixed or extrapolated from a target page:

{{< editor height="600px" >}}
LET baseURL = 'https://www.amazon.com/'
LET amazon = DOCUMENT(baseURL, { driver: "cdp" })

INPUT(amazon, '#twotabsearchtextbox', "ferret")
CLICK(amazon, '.nav-search-submit input[type="submit"]')
WAIT_NAVIGATION(amazon)

LET resultListSelector = 'div.s-result-list'
LET resultItemSelector = '[data-component-type="s-search-result"]'
LET nextBtnSelector = 'ul.a-pagination .a-last a'
LET priceWholeSelector = '.a-price-whole'
LET priceFracSelector = '.a-price-fraction'
LET pagers = ELEMENTS(amazon, 'ul.a-pagination li.a-disabled')
LET pages = LENGTH(pagers) > 0 ? TO_INT(INNER_TEXT(LAST(pagers))) : 0

LET result = (
    FOR pageNum IN 1..pages

        LET clicked = pageNum == 1 ? false : CLICK(amazon, nextBtnSelector)
        LET wait = clicked ? WAIT_NAVIGATION(amazon, 10000) : false
        LET waitSelector = wait ? WAIT_ELEMENT(amazon, resultListSelector) : false

        LET items = (
            FOR el IN ELEMENTS(amazon, resultItemSelector)
                LET hasPrice = ELEMENT_EXISTS(el, priceWholeSelector)
                LET priceWholeTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceWholeSelector), "[0-9]+")) : "0"
                LET priceFracTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceFracSelector), "[0-9]+")) : "00"
		        LET price = TO_FLOAT(priceWholeTxt + "." + priceFracTxt)
		        LET anchor = ELEMENT(el, "a")

                RETURN {
                    url: baseURL + anchor.attributes.href,
                    title: INNER_TEXT(el, 'h2'),
                    price
                }
        )

        RETURN items
)

RETURN FLATTEN(result)
{{</ editor >}}

### Uncontrolled
In turn, in uncontrolled pagination, we use a helper function [PAGINATION](/docs/stdlib/html/#pagination). The functions accepts an HTML element and a CSS selector for "Next" button. Once the the given selector returns empty result, iteration ends.    

{{< notification type="info">}}
Iteration always starts with a current page.
{{</ notification >}}

{{< editor height="600px" >}}
LET baseURL = 'https://www.amazon.com/'
LET amazon = DOCUMENT(baseURL, { driver: "cdp" })

INPUT(amazon, '#twotabsearchtextbox', "ferret")
CLICK(amazon, '.nav-search-submit input[type="submit"]')
WAIT_NAVIGATION(amazon)

LET resultListSelector = '#s-results-list-atf'
LET resultItemSelector = '[data-component-type="s-search-result"]'
LET nextBtnSelector = 'ul.a-pagination .a-last a'
LET priceWholeSelector = '.a-price-whole'
LET priceFracSelector = '.a-price-fraction'

LET result = (
    FOR pageNum IN PAGINATION(amazon, nextBtnSelector)
        LIMIT 3

        LET wait = pageNum > 0 ? WAIT_NAVIGATION(amazon, 20000) : false
        LET waitSelector = wait ? WAIT_ELEMENT(amazon, resultListSelector) : false

        LET items = (
            FOR el IN ELEMENTS(amazon, resultItemSelector)
                LET hasPrice = ELEMENT_EXISTS(el, priceWholeSelector)
                LET priceWholeTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceWholeSelector), "[0-9]+")) : "0"
                LET priceFracTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceFracSelector), "[0-9]+")) : "00"
		        LET price = TO_FLOAT(priceWholeTxt + "." + priceFracTxt)
		        LET anchor = ELEMENT(el, "a")

                RETURN {
                    url: baseURL + anchor.attributes.href,
                    title: INNER_TEXT(el, 'h2'),
                    price
                }
        )

        RETURN items
)

RETURN FLATTEN(result)
{{</ editor >}}