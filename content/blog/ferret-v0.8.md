---
title: "Ferret v0.8"
subtitle: "More features, better API"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-07-23"
---

Hooray, **[Ferret v0.8](https://github.com/MontFerret/ferret/releases/tag/v0.8.0)** has been released!

It's been a while since the last release, but we worked hard to bring new and better Ferret.

This release has many new exciting features, but unfortunately, also some breaking changes. 

The full changelog you can find **[here](https://github.com/MontFerret/ferret/blob/master/CHANGELOG.md#080)**.

Let's go!

# What's added
## iframe
Finally, Ferret has a support of ``iframe`` elements.    
When a page gets loaded, Ferret finds all available elements and provieds an access to them via ``.frames`` property of a page object.    

In order to find a target frame, you may use the following snippet:

{{< code fql >}}
LET page = DOCUMENT("https://www.bbc.com/", {
    driver: "cdp"
})

LET frame = (
    FOR frame IN page.frames
        FILTER frame.name == "smphtml5iframeplayer"
        RETURN frame
)

RETURN INNER_TEXT(FIRST(frame))
{{</ code >}}

Alternatevley, you can filter them out by url or access to a target iframe by index if you know it's position.

<div class="notification is-warning">
  You still may have issues with iframes pointing to another domain due to CORS security policy.
</div>

## Namespaces
With this release, we introduce a new language feature - namespaces.    

Namespaces will allow library authors (and us) to isoalate their functions into a dedicated sub sections.

Here is an example of possible use:

{{< code fql >}}
LET blob = DOWNLOAD("https://raw.githubusercontent.com/MontFerret/ferret/master/assets/logo.png")

RETURN IO::FS::WRITE("logo.png", blob)
{{</ code >}}

In order to use namespaces, you need to use new ``namespace`` method which can be chained:

{{< code go >}}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func main() {
    c := compiler.New()
    c.Namespace("IO").Namespace("FS").RegisterFunction("Read", fs.Read)
}
{{</ code >}}

<div class="notification is-info">
	In future releases we will put HTML related functions into HTML:: namespaces.
</div>

## XPath
Good web scraping tool cannot be without XPath support and Ferret finaly has it!   
Ferret provides simple interface to XPath engine for both drivers - CDP and HTTP.   
It automatically detects a type of output values and deserializes them accordingly.    

{{< code fql >}}
RETURN XPATH(page, "//div[contains(@class, 'form-group')]")
{{</ code >}}

{{< code fql >}}
RETURN XPATH(page, "count(//div)")
{{</ code >}}

These two queries will return 2 different types:    

1. Returns an array of serialized elements (their inner HTML)
2. Returns a number indicating how many "div" elements are on the page.

## Regular expression operator
This is a shorthand for using regexp assertions:

{{< code fql >}}
LET result = "foo" =~ "^f[o].$" // returns "true"
{{</ code >}}

{{< code fql >}}
LET result = "foo" !~ "[a-z]+bar$"  // returns "true"
{{</ code >}}

## New functions to manipulate DOM
There are some cases when you might need to do some changes in the existing DOM. That's why we added new ``INNER_HTML_SET`` and ``INNER_TEXT_SET`` functions to help you do that.

{{< code fql >}}
// Using document and selector
INNER_HTML_SET(doc, "body", "<span>Hello</span>")
INNER_TEXT_SET(doc, "body", "Hello")

// Or an element directly
INNER_HTML_SET(doc.body, "<span>Hello</span>")
INNER_TEXT_SET(doc.body, "Hello")
{{</ code >}}

## Viewport settings
In this release, you can override default values of a viewport in headless mode.

{{< code fql >}}
LET doc = DOCUMENT(@url, {
    driver: 'cdp',
    viewport: {
        width: 1920,
        height: 1080
    }
})
{{</ code >}}

## Better emulation of user interaction
This is a big change in how Ferret handles page interactions.     

Now Ferret performs it in a more advanced way - scrolls down or up to an element, moves the mouse, focuses and types... with a random delay. As a real person!

## Other
There are many other many small changes here and there like adding ``FOCUS``, ``ESCAPE_HTML``, ``UNESCAPE_HTML`` and ``DECODE_URI_COMPONENT`` functions, improving perfomance and changing internal design of some parts of the system.

# What's broken
We are trying to make things compatible between versions, but some features required serious design changes that lead to breaking the compatibility. 

The good news is that as we aproach to release v1.0, the API gets more stable and require less dramatic changes.

<div class="notification is-info">
	Most of the breaking changes will affect only embedded solutions, use of HTML drivers in particular. No changes in the syntax, so no scripts need to change.
</div>

## Virtual DOM structure
Work on ``iframe`` support required us to redesign the structure of the virtual DOM by introducing top level entity called ``HTMLPage``:

{{< code go>}}
type HTMLPage interface {
	core.Value
	core.Iterable
	core.Getter
	core.Setter
	collections.Measurable
	io.Closer

	IsClosed() values.Boolean
	GetURL() values.String
	GetMainFrame() HTMLDocument
	GetFrames(ctx context.Context) (*values.Array, error)
	GetFrame(ctx context.Context, idx values.Int) (core.Value, error)
	GetCookies(ctx context.Context) (*values.Array, error)
	SetCookies(ctx context.Context, cookies ...HTTPCookie) error
	DeleteCookies(ctx context.Context, cookies ...HTTPCookie) error
	PrintToPDF(ctx context.Context, params PDFParams) (values.Binary, error)
	CaptureScreenshot(ctx context.Context, params ScreenshotParams) (values.Binary, error)
	WaitForNavigation(ctx context.Context) error
	Navigate(ctx context.Context, url values.String) error
	NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error)
	NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error)
}
{{</ code >}}

Previously, the role of open page was played by ``HTMLDocument``, but the need of having multiple documents representing ``iframe`` nodes led us to bring a new entity to the structure.

## Driver API

Because of the changes in Virtual DOM structure, the driver API has been changed as well in order to be reasonable.

``Driver.LoadDocument`` and ``LoadDocumentParams`` are renamed to ``Driver.Open`` and ``Params``.

{{< code go>}}
type Driver interface {
	io.Closer
	Name() string
	Open(ctx context.Context, params Params) (HTMLPage, error)
}
{{</ code >}}

## Other
In the context of API stabilization and consistency, there are some other minor changes in vDOM elements like extra returned value (usually an error) or ``Get`` prefix in some methods.