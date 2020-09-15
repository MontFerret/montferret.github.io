---
title: "Ferret v0.12.0"
subtitle: "More stability"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2020-09-14"
---

Hello fellow miners,

We are happy to announce a new Ferret release - **[Ferret v0.12.0](https://github.com/MontFerret/ferret/releases/tag/v0.12.0)**.    

This release has a new useful module in the standard library - testing!   
Let's dig in. The full changelog you can find [here](https://github.com/MontFerret/ferret/blob/master/CHANGELOG.md#0120).

---


# What's added
## Assertion library
A new ``T`` namespace has been added to the standard library that provides some helpful methods for making assertions.     
The functions can be used to test UI, validate data and ensure the correctnes of a script.

{{< editor lang="fql" height="190px" readonly="true" >}}

LET doc = DOCUMENT('https://github.com/', { driver: "cdp" })

T::TRUE(ELEMENT_EXISTS(doc, "#user[login]"))
T::TRUE(ELEMENT_EXISTS(doc, "#user[email]"))
T::TRUE(ELEMENT_EXISTS(doc, "#user[password]"))

RETURN NONE

{{</ editor >}}

Read about available functions [here](https://www.montferret.dev/docs/stdlib/testing).

## New FRAMES function
A new helper function that finds HTML frames by a given property selector.

{{< editor lang="fql" height="190px" readonly="true" >}}

LET page = DOCUMENT('https://www.montferret.dev/fixtures/iframe/', {
    driver: 'cdp'
})

LET frame = FRAMES(page, "src", "https://www.montferret.dev/")

RETURN FIRST(frame).url

{{</ editor >}}

## iFrame navigation handling
Finally, we can control navigation of nested iframes! Whenever your page has a nested iframe that performs some nested navigation, you can tell Ferret to wait for its completion.

{{< editor lang="fql" height="580px" readonly="true" >}}
LET page = DOCUMENT('https://www.montferret.dev/fixtures/iframe/', {
    driver: 'cdp'
})

LET innerPage = FIRST(FRAMES(page, "src", "https://www.montferret.dev/"))

T::NOT::NONE(innerPage)

SCROLL_ELEMENT(innerPage)

CLICK(innerPage, "#navbar-burger")
WAIT_CLASS(innerPage, "#primary-nav", "is-active")

CLICK(innerPage, "#blog")

WAIT_NAVIGATION(innerPage)

WAIT_ELEMENT(innerPage, ".right-menu .mui-dropdown .login")

CLICK(ELEMENT(innerPage, ".blog-post a"))

WAIT_NAVIGATION(innerPage)


RETURN INNER_TEXT(innerPage, ".content")
{{</ editor >}}

# What's changed
## Removed property caching and tracking
This **could be** a breaking change, so, please, test your scripts. 
In previous versions, the CDP driver used aggressive caching and changes tracking on properties like ``.class``, ``.attributes`` and etc., in order to redcue amount of calls to a browser, which led to high CPU usage when a runtime held a big amount of elements in memory.   
To tackle the problem, the caching and tracking were removed and thus, every time element's property is accessed, a latest value gets retrived from a browser, without any further caching. If something needs to be cahced, it's user's responsibility to chace it in local variable.
