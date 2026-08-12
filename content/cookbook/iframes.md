---
title: "iframes"
weight: 3
draft: false
---

In order to find a particular iframe, you need to find it by its name or src.

{{< notification type="warning">}}
Beware, a target iframe must be in the same domain, otherwise its properties and content will be unavailable.
{{</ notification >}}

{{< editor height="250px" >}}
let page = DOCUMENT("https://www.w3schools.com/html/html_iframe.asp", {
    driver: "cdp"
})

let content = (
    for f in page.frames
        filter f.URL == "https://www.w3schools.com/html/default.asp"
            return f.head.innerHTML
)

return FIRST(content)
{{< /editor >}}