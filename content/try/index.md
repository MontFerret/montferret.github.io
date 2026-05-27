---
title: "Try it!"
slug: "/try/"
type: "repl"
draft: false
---

<p style="margin-bottom: 2em;">
    <a class="button is-dark" href="/try/next/">Try Ferret v2</a>
</p>

{{< editor id="replEditor" sharable="true" >}}
// Open the the product listing page using a browser-based driver (CDP)
// This allows Ferret to execute JavaScript and work with dynamic content
LET doc = DOCUMENT('https://mockery.montferret.dev/scenarios/dynamic-products/basic/', {
    driver: 'cdp'
})

// Wait until at least one product card is present on the page
// This is important because the page loads content asynchronously
WAIT_ELEMENT(doc, '.product-card', 5000)

// Select all product cards from the page
LET products = ELEMENTS(doc, '.product-card')

// Iterate over each product card and extract useful data
FOR product IN products
    RETURN {
        brand: TRIM(INNER_TEXT(product, '.product-brand')),
        title: TRIM(INNER_TEXT(product, '.product-title')),
        price: TO_FLOAT(SUBSTITUTE(INNER_TEXT(product, '.product-price'), '$', ''))
    }
{{</ editor >}}
