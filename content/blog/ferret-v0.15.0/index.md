---
title: "Ferret v0.15.0"
subtitle: "Moving forward"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2021-04-25"
---

Hello friends,

Ferret v0.15.0 is out. This release is more focused on bug fixes and API improvements.

Since the runtime and CLI have been splitted into 2 different repositories, you can find CLI binaries [here](https://github.com/MontFerret/cli/releases/tag/v1.1.0) which was upgraded to v1.1.0.


# What's added

## Support of document charset in HTTP driver
If your target page does not contain information if its content encoding and it is other than ``utf-8``, you can help ``http`` driver to properly convert the page content into ``utf-8``.

{{< code lang="fql" height="300px" readonly="true" >}}
LET page = DOCUMENT("www.non-utf8-page.co", {
    driver: "http",
    charset: "windows-1251"
})
{{</ code >}}

## Possibility to send keyboard events

Have you ever scraped a page with search input field that could be tirggered by pressing "Enter" key on your keyboard and there was no way of doing this using Ferret? Such a bumer!

Well, now you can do that with new ``PRESS`` and ``PRESS_SELECTOR`` functions!

{{< editor lang="fql" height="300px" readonly="true" >}}
LET page = DOCUMENT("https://soundcloud.com/", {
    driver: "cdp"
})

WAIT_ELEMENT(page, '.headerSearch__input')
LET search = ELEMENT(page, '.headerSearch__input')

INPUT(search, "Singing ferret")
WAIT(200)
PRESS(search, "Enter")

WAIT_ELEMENT(page, ".lazyLoadingList__list")

FOR track IN ELEMENTS(page, ".searchItem")
    LIMIT 3
    RETURN {
        artist: TRIM(INNER_TEXT(track, '.soundTitle__usernameText')),
        track: TRIM(INNER_TEXT(track, '.soundTitle__title'))
    }
{{</ editor >}}

# What's fixed

Ok, now, let's take a look what bugs have been fixed!

- Passing headers and cookies to HTTP driver
- Reading property of anyonymous object
- Clearing input text containing special characteers

# Summary
As always, thanks everyone who contributed to the project!