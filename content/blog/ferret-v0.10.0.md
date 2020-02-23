---
title: "Ferret v0.10.0"
subtitle: "Comeback"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2020-02-21"
---

Hello fellow miners, **[Ferret v0.10.0](https://github.com/MontFerret/ferret/releases/tag/v0.10.0)** has been **finally** released!

It's been long and busy 6 months since the last release and many of you might thought that the project got abandoned.    
We assure you it's not! But sometimes life happens and slows you down.

But we are back and ready to rock!

This release has many great features and improvements thank to all our contributors.    
Let's take a look at the most important ones. The full changelog you can find [here](https://github.com/MontFerret/ferret/blob/master/CHANGELOG.md#0100).

# What's added
## I/O functions
You asked - we added! Finally, I/O functions have arrived.
In the beginning of the project, we were concerned about security threats that access to filesystem poses, but seeing how demanded this functionality is we agreed to add it.    
Now, Ferret has new I/O namespaces, with the following functions:

- ``IO::FS::READ``
- ``IO::FS::WRITE``
- ``IO::NET::HTTP:GET``
- ``IO::NET::HTTP:POST``
- ``IO::NET::HTTP:PUT``
- ``IO::NET::HTTP:DELETE``
- ``IO::NET::HTTP:REQUEST``

Enjoy!

## Loading HTML string into memory
Another popular feature request - possibility to load a prefetched HTML string into Ferret.    
With this release, you can finally do it with ``HTML_PARSE`` function.

{{< code fql >}}
LET file = IO::FS::READ(@myfile)

LET doc = HTML_PARSE(TO_STRING(file), {
    driver: 'cdp' // or 'http'
})

RETURN INNER_TEXT(doc, 'title')
{{</ code >}}

## Parameter value availability check
In this release, we do pre-runtime check whether all values are provided for parameters used by a script.    
It was very frustrating when your script was failing in the middle of its execution just because you forgot to provide a value for one of the paramaters.

# What's changed
## Case insensitivity
Finally, FQL keywords are case insensitive!

That's how Google Search query looks like in lower case now:

{{< code fql >}}
let google = document("https://www.google.com/", {
    driver: "cdp",
})

input(google, 'input[name="q"]', "ferret")
click(google, 'input[name="btnK"]')

wait_navigation(google)

for result in elements(google, '.g')
    filter trim(result.attributes.class) == 'g'
    return {
        title: inner_text(result, 'h3'),
        description: inner_text(result, '.st'),
        url: inner_text(result, 'cite')
    }
{{</ code >}}

## Improved CDP driver performance
CDP driver performance has been drastically improved that brings CPU usage down from 60% for 3 pages to <1%.   

# What's fixed
## Redirects
It was a long lasting problem that Ferret could not correctly handle redirects that could occur during a page navigation that was leading to timeouts or corrupt page state.    

Now it has finally been fixed!   
```WAIT_NAVIGATION``` function has 2nd optional parameter that can be an object with the following fields:

{{< code go >}}
type Parameters struct {
    TargetURL string
    Timeout time.Duration
}
{{</ code >}}

``TargetURL`` is a regexp string that can be used to give Ferret a hint what the destination url is:

{{< code fql >}}
LET doc = DOCUMENT("http://waos.ovh/redirect.html", {
    driver: 'cdp'
})

CLICK(doc, '.click')

WAIT_NAVIGATION(doc, { targetURL: "redirect2\.html" })

RETURN ELEMENT(doc, '.title')
{{</ code >}}