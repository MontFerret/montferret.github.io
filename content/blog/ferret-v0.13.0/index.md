---
title: "Ferret v0.13.0"
subtitle: "Syntax change"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2020-11-24"
---

Hello friends,

We are excited to announce a new Ferret release - **[Ferret v0.13.0](https://github.com/MontFerret/ferret/releases/tag/v0.13.0)**.  

This release brings new syntax, DOM API updates and some bug fixes. Let's dive in!

# What's added

## While loop
Finally! While loops have landed FQL!    

Ferret query language inherently lacked while loops due to its origin from ArangoDB query language. But over time, it's become obvious that this kind of looping is essential for data scraping.

While loops in FQL are implemented as an extension of ``FOR IN`` loop and come in two modes: ``while-do`` and ``do-while``:

```
FOR counter [DO] WHILE condition
    RETURN value
```

``while-do`` is a traditional while loop where a condition is checked before every iteration. It might be not very useful in web scraping (e.g. pagination), thus we have implemented ``do-while`` where condition checks start to occur on 2nd iteration i.e. there is always at least one iteration.

Since data in FQL is immutable, there are no direct mutations in code which while loops can depend on. Rather while loops depend on checks of external resources like existance of HTML elements or database records.

Ok, enough of theory! Let's take a look at the query that gets a list of API managment solutions from the GitHub Marketplace:

{{< editor lang="fql" height="340px" readonly="true" >}}
LET doc = DOCUMENT("https://github.com/marketplace/category/api-management", { driver: "cdp"})
LET nextSelector = ".paginate-container .BtnGroup a:nth-child(2)"
LET elementsSelector = '[data-hydro-click]'

WAIT_ELEMENT(doc, elementsSelector)

FOR i DO WHILE ELEMENT_EXISTS(doc, nextSelector)
    LIMIT 2
	LET wait = i > 0 ? CLICK(doc, nextSelector) : false
	LET nav = wait ? WAIT(2000) && WAIT_ELEMENT(doc, elementsSelector) : false

	FOR el IN ELEMENTS(doc, elementsSelector)
        FILTER ELEMENT_EXISTS(el, '.h4')
		RETURN {
            url: el.attributes.href,
            name: INNER_TEXT(el, '.h4')
        }
{{</ editor >}}

In this while loop, we are checking if the next button exists (external mutation) but always execute the loop body at least once just in case there is only one page.

The ``i`` variable represents a loop counter that gives information which iteration we are in.
Being an extension of ``FOR-IN`` loop gives us all statements that can be used with regular ``FOR-IN`` loop like ``FILTER``, ``LIMIT`` and etc.

## Computed styles

Until this version, CDP driver returned attribute based styles whenver you accessed them via ``Element.style`` or ``Element.attributes.style`` or ``STYLE_GET(Element)``.

Starting with this release, ``Element.style`` returns an object containing computed styles while ``Element.attributes.style`` returns an object containg styles defined in ``style`` attribute of an element. ``STYLE_GET(Element)`` returns computed styles as well.

{{< editor lang="fql" height="200px" readonly="true" >}}
LET doc = DOCUMENT("https://github.com/", { driver: "cdp"})
LET el = ELEMENT(doc, "#user_email")

RETURN {
    computed: el.style,
    attribute: el.attributes.style
}
{{</ editor >}}

## Element's parent and siblings

With this release, it's possible to direclty access parents and siblings of elements:

{{< editor lang="fql" height="200px" readonly="true" >}}
LET doc = DOCUMENT("https://getbootstrap.com/docs/4.5/components/list-group/", { driver: "cdp"})
LET list = ELEMENT(doc, 'ul.list-group')
LET items = ELEMENTS(list, 'li')
LET firstItem = items[0]
LET secondItem = items[1]

RETURN {
    parent: secondItem.parentElement.innerText,
    first: [secondItem.previousElementSibling, firstItem],
    second: [firstItem.nextElementSibling, secondItem]
}
{{</ editor >}}

# What's fixed
## HTML escaping

Long lasting issue of HTML escaping has been finally fixed. 

The problem was that Go's JSON encoder from the standard library does escaping of HTML strings.