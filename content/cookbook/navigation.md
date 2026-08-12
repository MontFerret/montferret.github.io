---
title: "Navigation"
weight: 4
draft: false
---

### Navigate by url

{{< editor height="250px" >}}
let page = DOCUMENT("https://github.com/", { driver: "cdp" })
let header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE(page, "https://github.com/features", 10000)

let header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

return [header1, header2]
{{< /editor >}}

### Navigate back

{{< editor height="250px" >}}
let page = DOCUMENT("https://github.com/", { driver: "cdp" })

NAVIGATE(page, "https://github.com/features", 10000)

let header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE_BACK(page)

let header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

return [header1, header2]
{{< /editor >}}


### Navigate forward

{{< editor height="250px" >}}
let page = DOCUMENT("https://github.com/", { driver: "cdp" })

NAVIGATE(page, "https://github.com/features", 10000)

let header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE_BACK(page)
NAVIGATE_FORWARD(page)

let header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

return [header1, header2]
{{< /editor >}}