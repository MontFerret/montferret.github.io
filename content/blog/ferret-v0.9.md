---
title: "Ferret v0.9"
subtitle: "Moving forward"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-09-08"
---

Hello fellow miners, **[Ferret v0.9](https://github.com/MontFerret/ferret/releases/tag/v0.9.0)** has been released!

In this release, we mostly focused on bug fixes and filling gaps in user interaction functionality.

Let's see what we've got.

# What's added
## Clearing input values
Since early versions, Ferret has ``INPUT`` HTML function that allowed to type any value to input HTML element.
This works by appending the function parameter to any pre-existing `<input>` value. While it works fine in many cases, there are some scenarios when a target input field has a preset value which needs to be erased. With this release, Ferret has a new method ``INPUT_CLEAR`` that erases any data from a target input.

{{< code fql >}}
LET doc = DOCUMENT(@url, {
    driver: 'cdp'
})

RETURN INPUT_CLEAR(doc, '#my-input')
{{</ code >}}

## Double and more click
Single, double, or more. Now it's possible to specify the number of clicks to ``CLICK`` function.

{{< code fql >}}
LET doc = DOCUMENT(@url, {
    driver: 'cdp'
})

RETURN CLICK(doc, '#like', 2)
{{</ code >}}

## Unfocus
``BLUR`` function removes focus from an active input field.

{{< code fql >}}
LET doc = DOCUMENT(@url, {
    driver: 'cdp'
})

FOCUS(doc, '#input')
BLUR(doc, '#input')

RETURN NONE
{{</ code >}}

## Default headers and cookies
With this release, you can set default headers and/or cookies to HTML drivers.
One of the possible use cases is to create drivers with different custom names and pre-defined headers/cookies:

{{< code go>}}
cdp.NewDriver(
    cdp.WithCustomName("cdp_headers"),
    cdp.WithHeader("Single_header", []string{"single_header_value"}),
    cdp.WithHeaders(drivers.HTTPHeaders{
        "Multi_set_header":  []string{"multi_set_header_value"},
        "Multi_set_header2": []string{"multi_set_header2_value"},
    }),
)

cdp.NewDriver(
    cdp.WithCustomName("cdp_cookies"),
    cdp.WithCookie(drivers.HTTPCookie{
        Name:     "single_cookie",
        Value:    "single_cookie_value",
        Path:     "/",
        MaxAge:   0,
        Secure:   false,
        HTTPOnly: false,
        SameSite: 0,
    }),
    cdp.WithCookies([]drivers.HTTPCookie{
        {
            Name:     "multi_set_cookie",
            Value:    "multi_set_cookie_value",
            Path:     "/",
            MaxAge:   0,
            Secure:   false,
            HTTPOnly: false,
            SameSite: 0,
        },
        {
            Name:     "multi_set_cookie2",
            Value:    "multi_set_cookie2_value",
            Path:     "/",
            MaxAge:   0,
            Secure:   false,
            HTTPOnly: false,
            SameSite: 0,
        },
    }),
)
{{</ code >}}

## Params in dot notation
In previous versions of Ferret, it was possible to pass complex types as parameters, but it was impossible to read data from it using dot notation. Now it is possible. But even more - you can use params as dot notation segments:

{{< code fql >}}
RETURN @one.@two.@three
{{</ code >}}

## String literals
In addition to defining string literals with quotes (', ") and backtick (\`), you can now define string literals with the forwardtick (´).

# What's fixed
## Open tabs on error
Ferret didn't close Chrome/Chromium tab if an error occurs during a page load.

## CLICK
``CLICK`` could not be used with both an element and selector.

## Dot notation after a function call
Parser could not properly handle scenarios where dot notation was used right after a function call:

{{< code fql >}}
LET items = [{name: "foo"}]
RETURN FIRST(items).name
{{</ code >}}

# What's changed
## Internal cleanup and refactoring
We've made some internal cleanups and minor refactoring.    
There might be minor breaking changes for those who use HTML drivers directly. 

#### 2nd returned value
The following methods have an error as asecond returned value now: 

- ``HTMLNode.GetChildNodes``
- ``HTMLNode.GetChildNode``
- ``HTMLNode.QuerySelector``
- ``HTMLNode.QuerySelectorAll``
- ``HTMLNode.CountBySelector``
- ``HTMLElement.GetValue``
- ``HTMLElement.GetAttributes``
- ``HTMLElement.GetAttribute``

#### Re-arranged methods
The following methods have been moved from ``HTMLDocument`` to ``HTMLElement`` interface:

- ``SelectBySelector``
- ``FocusBySelector``
- ``HoverBySelector``

#### Removed methods
The following method(s) have been removed from ``HTMLDocument`` interface:

- ``MoveMouseBySelector``
