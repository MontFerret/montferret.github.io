---
title: "Try Ferret v2"
url: "/try/next/"
type: "repl"
draft: false
---

{{< editor id="replEditorNext" sharable="true" apiVersion="2" height="384px" >}}
// Open the the product listing page using a browser-based driver (CDP)
// This allows Ferret to execute JavaScript and work with dynamic content
LET doc = DOCUMENT('https://mockery.ferretlang.org/scenarios/dynamic-products/basic/', {
    driver: 'cdp'
})

// Wait until at least one product card is present on the page
// This is important because the page loads content asynchronously
// WAITFOR VALUE returns the value of the expression once it evaluates to a non-empty value or until the timeout is reached
LET products = WAITFOR VALUE doc[~ css`.product-card`]
    WHEN LENGTH(.) > 0
    TIMEOUT 5000

// Iterate over each product card and extract useful data
FOR product IN products
    RETURN {
        brand: product[~? css`.product-brand`].textContent,
        title: product[~? css`.product-title`].textContent,
        price: TO_FLOAT(SUBSTITUTE(product[~? css`.product-price`], '$', '')) ON ERROR RETURN 0
    }
{{</ editor >}}
