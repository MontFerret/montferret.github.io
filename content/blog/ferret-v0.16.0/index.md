---
title: "Ferret v0.16.0"
subtitle: "Big changes"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2021-11-22"
---

Hello friends,

We hope everyone is healthy, happy, and ready for winter :)

Our team is back from the break and excited to announce a new Ferret release!

Ferret v0.16.0 is here and now it's more stable, flexible, and faster!

As always, CLI can be found [here](https://github.com/MontFerret/cli/releases/tag/v1.2.0) and the runtime itself [here](https://github.com/MontFerret/ferret/releases/tag/v0.16.0).

Let's dive in and see what we've got!


# What's new
In this release, besides internal optimizations, we've focused on new language syntax and order to 

## WAITFOR
This is the beginning of our [initiative](https://github.com/MontFerret/ferret/issues/263) to replace ``WAIT_`` functions by providing a universal interface in the runtime type system. We've started with ``WAITFOR EVENT`` statement.

The statement allows developers to declare a pause of the execution until a certain event occurs in a given source or a time is out (by default: 5000 ms).

The syntax is following:

```fql
WAITFOR EVENT event_name IN source [OPTIONS options] [FILTER expression] [TIMEOUT timeout]
```

It also can be used as an expression and assign returned event value to a variable:

```fql
LET evt = (WAITFOR EVENT event_name IN source [OPTIONS options] [FILTER expression] [TIMEOUT timeout])
```

At this moment, 3 events are supported by the ``cdp`` driver:
- navigation
- request
- response

### Navigation
``navigation`` event is intended to replace ``WAIT_NAVIGATION`` function which will be deprecated in the next release and focused to be more flexible.

```fql
WAITFOR EVENT "navigation" IN page
```

While the simple version of the statement provides extra readability, it also provides extra control over navigation using ``FILTER`` statement

```fql
WAITFOR EVENT "navigation" IN page FILTER CURREN.url LIKE "MY_DESIRABLE_URL"
```
You may have noticed ``CURRENT`` keyword inside the filter statement.
This is a pseudo-variable that provides access to the current event object. The event object will depend on the type of an event, in the case of navigation, the object has the following fields:

```fql
{
 url
 frame
 mimeType
}
```

``frame`` property allows to wait for a navigation event in one of nested frames:

```fql
WAITFOR EVENT "navigation" IN page FILTER CURREN.frame == nested_doc
```

### AJAX
``request`` and ``response`` are the events you can use to subscribe to AJAX request/response operations:

{{< editor lang="fql" height="300px" readonly="true" >}}
LET doc = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })

WAIT_ELEMENT(doc, '.chartTrack__details', 5000)
SCROLL_BOTTOM(doc)

LET evt = (WAITFOR EVENT "request" IN doc FILTER CURRENT.url LIKE "https://api-v2.soundcloud.com/charts?genre=soundcloud*")

RETURN evt.headers
{{</ editor >}}

Request object:
```fql
{
 url
 method
 headers
 body
}
```

{{< editor lang="fql" height="250px" readonly="true" >}}
LET doc = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })

WAIT_ELEMENT(doc, '.chartTrack__details', 5000)
SCROLL_BOTTOM(doc)

LET evt = (WAITFOR EVENT "response" IN doc FILTER CURRENT.url LIKE "https://api-v2.soundcloud.com/charts?genre=soundcloud*")

RETURN JSON_PARSE(evt.body)
{{</ editor >}}

Response object:
```fql
{
 url 
 statusCode
 status 
 headers 
 body 
 responseTime
}
```

They can be useful for intercepting data between a client and a server.

## Optional chaining
From the early days of Ferret, the runtime has been pretty strict on accessing properties of values that either are null or do not provide any properties by terminating an execution with a type error. But sometimes, the result of the accessing property might not be that important which makes us write long ternary operators that do some checks.

Well, now you do not have to do it anymore! This release is bringing the optional chaining aka Elvis operator to FQL language:

<div class="notification is-info">
    The optional chaining operator (?.) enables you to read the value of a property located deep within a chain of connected objects without having to check that each reference in the chain is valid.
</div>

{{< editor lang="fql" height="100px" readonly="true" >}}
LET foo = { bar: NONE }
RETURN foo?.bar?.baz
{{</ editor >}}

The ``?.`` operator is like the ``.`` chaining operator, except that instead of returning an error if a reference is null, the expression short-circuits with a return value of ``NONE``.

## Errors suppression
Errors suppression feature is somewhat similar to the optional chaining in allowing to tolerate insignificant errors during query execution and return ``NONE`` in case of occurred error:

{{< editor lang="fql" height="100px" readonly="true" >}}
LET el = ELEMENT(NONE, "#el")?

RETURN el == NONE
{{</ editor >}}

If you try to execute the same query, but without ``?`` operator, it will fail:

{{< editor lang="fql" height="100px" readonly="true" >}}
LET el = ELEMENT(NONE, "#el")

RETURN el == NONE
{{</ editor >}}

The ``?`` operator is like the ``?.`` optional chaining operator, but for function calls and inline statements:

{{< editor lang="fql" height="150px" readonly="true" >}}
LET page = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })
LET evt = (WAITFOR EVENT "navigation" IN page)?

RETURN evt == NONE ? "not navigated" : "navigated"
{{</ editor >}}

## XPath selectors
This is one of the long waiting features - universal use of XPath selectors.
With updated interface of HTML driver componenets, it is possible to pass XPath selectors to HTML API functions by using new ``X`` function:

{{< editor lang="fql" height="250px" readonly="true" >}}
LET page = DOCUMENT('https://soundcloud.com/charts/top', { driver: "cdp" })
WAIT_ELEMENT(page, X("//div[contains(@class, 'chartTrack__details')]"), 5000)

FOR el IN ELEMENTS(page, X("//div[contains(@class, 'chartTracks')]/ul/li"))
    LET details = ELEMENT(el, X('./div/div[3]'))
    RETURN {
        artist: INNER_TEXT(details, X('./div[1]')),
        title:  INNER_TEXT(details, X('./div[2]')),
    }
{{</ editor >}}

``X`` function creates an XPath type of query selector which is supported by all HTML API functions now. The function automatically resolves returned value based on its type.

## Discard results
Sometimes we need to do some prep work in a ``FOR`` loop but are not interested in the returned result. With this release we can use Go-like ``_`` variable name to discard results.

Like in Go, you cannot access this variable later:

{{< editor lang="fql" height="250px" readonly="true" >}}
LET _ = (
    FOR num IN 1..10
        PRINT(num)
        RETURN NONE
)

RETURN _
{{</ editor >}}

# What's fixed
Here is a list of important bug fixes:

- ``values.Parse`` does not parse int64
- CPU leakage
- nil pointer exception
- HTTP driver makes multiple requests #642

# Summary
That's it, folks! Full list of changes you can find [here](https://github.com/MontFerret/ferret/blob/master/CHANGELOG.md#0160).
Thanks everyone who contributed to the project!