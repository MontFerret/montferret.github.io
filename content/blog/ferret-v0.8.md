---
title: "Ferret v0.8"
subtitle: "More features, better API"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-07-23"
---

Hooray, **[Ferret v0.8](https://github.com/MontFerret/ferret/releases/tag/v0.8.0)** has been released!

It's been a while since the last release, but we worked hard to bring new and better Ferret. This release has many new exciting features, but unfortunately, there are also some breaking changes. 

You can find the full changelog **[here](https://github.com/MontFerret/ferret/blob/master/CHANGELOG.md#080)**.

Let's go!

# What's added
## iframe
Ferret *finally* supports ``iframe`` elements.    
When a page gets loaded, Ferret finds all available elements and provieds an access to them via the ``.frames`` property of a page object.

Here's an example of how to find a target frame:

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

Alternatively, you can filter them out by url or access to a target iframe by index, if you know it's position.

<div class="notification is-warning">
  Due to CORS security policies, you still may have issues with iframes if it points to another domain.
</div>

## Namespaces
With this release, we introduce a new language feature - **namespaces**.    

Namespaces allow library authors (and us) to isolate functions into dedicated sub sections.

Here is an example:

{{< code fql >}}
LET blob = DOWNLOAD("https://raw.githubusercontent.com/MontFerret/ferret/master/assets/logo.png")

RETURN IO::FS::WRITE("logo.png", blob)
{{</ code >}}

To namespace a function, use the new ``namespace`` method. The ``namespace`` method is chainable:

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
A good web scraping tool needs XPath support, and Ferret finally has it!   
Ferret provides simple interface to XPath engine for both drivers - CDP and HTTP.   
It automatically detects the output value type and deserializes them accordingly.    

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
This release provides a shorthand for using regexp assertions:

{{< code fql >}}
LET result = "foo" =~ "^f[o].$" // returns "true"
{{</ code >}}

{{< code fql >}}
LET result = "foo" !~ "[a-z]+bar$"  // returns "true"
{{</ code >}}

## New functions to manipulate DOM
There are some cases when you might need to change the existing DOM. To help with that, we added the ``INNER_HTML_SET`` and ``INNER_TEXT_SET`` functions.

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

Now Ferret interacts with pages in a more advanced way - your script can scrolls down or up to an element, moves the mouse, focuses and types... with random delays. Just like a real person!

## Other
There are many other many small changes here and there, like adding ``FOCUS``, ``ESCAPE_HTML``, ``UNESCAPE_HTML`` and ``DECODE_URI_COMPONENT`` functions; improving performance; and changing internal design of some parts of the system.

# What's broken
We try to maintain backwards compatibility, but some of the new features required serious design changes that lead to breaking compatibility with previous versions.  As we approach to release v1.0, the API is becoming more stable and will require fewer dramatic changes.

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

Previously, ``HTMLDocument`` contained the open page, but ``iframe`` nodes introduce the need to have multiple documents representing each node. This led to a new entity in the structure.

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
