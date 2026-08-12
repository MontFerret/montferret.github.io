---
title: "Pagination"
weight: 2
draft: false
---

There are several way how to implement pagination.

### while loop
Since v0.13.0, pagination can be implemented with ``for-while`` loop.

{{< notification type="info">}}
For paginations, it's recommended to use do-while variation of the loop, in order to process at least the first page.
{{</ notification >}}

{{< editor height="300px" >}}
let doc = DOCUMENT("https://github.com/MontFerret/ferret/stargazers", { driver: "cdp" })

let nextSelector = ".paginate-container .BtnGroup a:nth-child(2)"
let elementsSelector = '.follow-list li'

return for i do while ELEMENT_EXISTS(doc, nextSelector)
    limit 3
	let wait = i > 0 ? CLICK(doc, nextSelector) : false
	let nav = wait ? WAIT_NAVIGATION(doc) : false
	
	return for el in ELEMENTS(doc, elementsSelector)
		filter ELEMENT_EXISTS(el, ".octicon-organization")

		return {
			name: INNER_TEXT(el, ".follow-list-name"),
			company: INNER_TEXT(el, ".follow-list-info span")
		}
{{</ editor >}}

### Controlled
You can also use ``for-in`` loop with specified range of iterations that can be either fixed or extrapolated from a target page:

{{< editor height="600px" >}}
let baseURL = 'https://www.amazon.com/'
let amazon = DOCUMENT(baseURL, { driver: "cdp" })

INPUT(amazon, '#twotabsearchtextbox', "ferret")
CLICK(amazon, '.nav-search-submit input[type="submit"]')
WAIT_NAVIGATION(amazon)

let resultListSelector = 'div.s-result-list'
let resultItemSelector = '[data-component-type="s-search-result"]'
let nextBtnSelector = 'ul.a-pagination .a-last a'
let priceWholeSelector = '.a-price-whole'
let priceFracSelector = '.a-price-fraction'
let pagers = ELEMENTS(amazon, 'ul.a-pagination li.a-disabled')
let pages = LENGTH(pagers) > 0 ? TO_INT(INNER_TEXT(LAST(pagers))) : 0

let result = (
    for pageNum in 1..pages

        let clicked = pageNum == 1 ? false : CLICK(amazon, nextBtnSelector)
        let wait = clicked ? WAIT_NAVIGATION(amazon, 10000) : false
        let waitSelector = wait ? WAIT_ELEMENT(amazon, resultListSelector) : false

        let items = (
            for el in ELEMENTS(amazon, resultItemSelector)
                let hasPrice = ELEMENT_EXISTS(el, priceWholeSelector)
                let priceWholeTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceWholeSelector), "[0-9]+")) : "0"
                let priceFracTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceFracSelector), "[0-9]+")) : "00"
		        let price = TO_FLOAT(priceWholeTxt + "." + priceFracTxt)
		        let anchor = ELEMENT(el, "a")

                return {
                    url: baseURL + anchor.attributes.href,
                    title: INNER_TEXT(el, 'h2'),
                    price
                }
        )

        return items
)

return FLATTEN(result)
{{</ editor >}}

### Uncontrolled
In turn, in uncontrolled pagination, we use a helper function [PAGINATION](/docs/standard-library/html/#pagination). The functions accepts an HTML element and a CSS selector for "Next" button. Once the the given selector returns empty result, iteration ends.

{{< notification type="info">}}
Iteration always starts with a current page.
{{</ notification >}}

{{< editor height="600px" >}}
let baseURL = 'https://www.amazon.com/'
let amazon = DOCUMENT(baseURL, { driver: "cdp" })

INPUT(amazon, '#twotabsearchtextbox', "ferret")
CLICK(amazon, '.nav-search-submit input[type="submit"]')
WAIT_NAVIGATION(amazon)

let resultListSelector = '#s-results-list-atf'
let resultItemSelector = '[data-component-type="s-search-result"]'
let nextBtnSelector = 'ul.a-pagination .a-last a'
let priceWholeSelector = '.a-price-whole'
let priceFracSelector = '.a-price-fraction'

let result = (
    for pageNum in PAGINATION(amazon, nextBtnSelector)
        limit 3

        let wait = pageNum > 0 ? WAIT_NAVIGATION(amazon, 20000) : false
        let waitSelector = wait ? WAIT_ELEMENT(amazon, resultListSelector) : false

        let items = (
            for el in ELEMENTS(amazon, resultItemSelector)
                let hasPrice = ELEMENT_EXISTS(el, priceWholeSelector)
                let priceWholeTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceWholeSelector), "[0-9]+")) : "0"
                let priceFracTxt = hasPrice ? FIRST(REGEX_MATCH(INNER_TEXT(el, priceFracSelector), "[0-9]+")) : "00"
		        let price = TO_FLOAT(priceWholeTxt + "." + priceFracTxt)
		        let anchor = ELEMENT(el, "a")

                return {
                    url: baseURL + anchor.attributes.href,
                    title: INNER_TEXT(el, 'h2'),
                    price
                }
        )

        return items
)

return FLATTEN(result)
{{</ editor >}}
