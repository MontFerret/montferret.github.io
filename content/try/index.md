---
title: "Try it!"
slug: "/try/"
type: "repl"
draft: false
---

{{< editor id="replEditorNext" sharable="true" apiVersion="2" orientation="horizontal" >}}
// Open the the product listing page using a browser-based driver (CDP)
// This allows Ferret to execute JavaScript and work with dynamic content
let doc = DOCUMENT('https://mockery.ferretlang.org/scenarios/dynamic-products/basic/', {
    driver: 'cdp'
})

// Wait until at least one product card is present on the page
// This is important because the page loads content asynchronously
// WAITFOR VALUE accepts any non-NONE candidate; WHEN keeps polling until this list is non-empty
let products = waitfor value doc[~ css`.product-card`]
    when LENGTH(.) > 0
    timeout 5000

// Iterate over each product card and extract useful data
return for product in products
    return {
        brand: product[~? css`.product-brand`].textContent,
        title: product[~? css`.product-title`].textContent,
        price: TO_FLOAT(SUBSTITUTE(product[~? css`.product-price`], '$', '')) on error return 0
    }
{{</ editor >}}
