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
LET page = DOCUMENT("https://www.w3schools.com/html/html_iframe.asp", {
    driver: "cdp"
})

LET content = (
    FOR f IN page.frames
        FILTER f.URL == "https://www.w3schools.com/html/default.asp"
            RETURN f.head.innerHTML
)

RETURN FIRST(content)
{{< /editor >}}