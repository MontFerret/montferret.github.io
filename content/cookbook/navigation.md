---
title: "Navigation"
weight: 4
draft: false
---

### Navigate by url

{{< editor height="250px" >}}
LET page = DOCUMENT("https://github.com/", { driver: "cdp" })
LET header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE(page, "https://github.com/features", 10000)

LET header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

RETURN [header1, header2]
{{< /editor >}}

### Navigate back

{{< editor height="250px" >}}
LET page = DOCUMENT("https://github.com/", { driver: "cdp" })

NAVIGATE(page, "https://github.com/features", 10000)

LET header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE_BACK(page)

LET header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

RETURN [header1, header2]
{{< /editor >}}


### Navigate forward

{{< editor height="250px" >}}
LET page = DOCUMENT("https://github.com/", { driver: "cdp" })

NAVIGATE(page, "https://github.com/features", 10000)

LET header1 = ELEMENT(page, 'meta[name="description"]').attributes.content

NAVIGATE_BACK(page)
NAVIGATE_FORWARD(page)

LET header2 = ELEMENT(page, 'meta[name="description"]').attributes.content

RETURN [header1, header2]
{{< /editor >}}