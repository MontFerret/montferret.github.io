---
title: "Ferret v0.12.0"
subtitle: "More stability"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2020-09-03"
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

{{< editor lang="fql" height="190px" >}}

LET doc = DOCUMENT('https://github.com/', { driver: "cdp" })

T::TRUE(ELEMENT_EXISTS(doc, "#user[login]"))
T::TRUE(ELEMENT_EXISTS(doc, "#user[email]"))
T::TRUE(ELEMENT_EXISTS(doc, "#user[password]"))

RETURN data

{{</ editor >}}

Read about available functions [here](https://www.montferret.dev/docs/stdlib/testing).

## iFrame navigation handling
Finally, we can control navigation of nested iframes! Whenever your page has a nested iframe that performs some nested navigation, you can tell Ferret to wait for its completion.

# What's changed
## Removed property caching and tracking
This **could be** a breaking change, so, please, test your scripts. 
In previous versions, the CDP driver used aggressive caching and changes tracking on properties like ``.class``, ``.attributes`` and etc., in order to redcue amount of calls to a browser, which led to high CPU usage when a runtime held a big amount of elements in memory.   
To tackle the problem, the caching and tracking were removed and thus, every time element's property is accessed, a latest value gets retrived from a browser, without any further caching. If something needs to be cahced, it's user's responsibility to chace it in local variable.
