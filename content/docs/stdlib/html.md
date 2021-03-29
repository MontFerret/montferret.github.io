---
title: "html"
weight: 1
draft: false
menuTitle: 
menu: [ATTR_GET,ATTR_QUERY,ATTR_REMOVE,ATTR_SET,BLUR,CLICK,CLICK_ALL,COOKIE_DEL,COOKIE_GET,COOKIE_SET,DOCUMENT,DOCUMENT_EXISTS,DOWNLOAD,ELEMENT,ELEMENTS,ELEMENTS_COUNT,ELEMENT_EXISTS,FOCUS,FRAMES,HOVER,INNER_HTML,INNER_HTML_ALL,INNER_HTML_SET,INNER_TEXT,INNER_TEXT_ALL,INNER_TEXT_SET,INPUT,INPUT_CLEAR,MOUSE,NAVIGATE,NAVIGATE_BACK,NAVIGATE_FORWARD,PAGINATION,PARSE,PDF,SCREENSHOT,SCROLL,SCROLL_BOTTOM,SCROLL_ELEMENT,SCROLL_TOP,SELECT,STYLE_GET,STYLE_REMOVE,STYLE_SET,WAIT_ATTR,WAIT_ATTR_ALL,WAIT_CLASS,WAIT_CLASS_ALL,WAIT_ELEMENT,WAIT_NAVIGATION,WAIT_NO_ATTR,WAIT_NO_ATTR_ALL,WAIT_NO_CLASS,WAIT_NO_CLASS_ALL,WAIT_NO_ELEMENT,WAIT_NO_STYLE,WAIT_NO_STYLE_ALL,WAIT_STYLE,WAIT_STYLE_ALL,XPATH,]
---



{{< header href="attr_get" >}}

ATTR_GET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_get.go#L14)

ATTR_GET gets single or more attribute(s) of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target node.
`attrNames` | `String, repeated`  |  | Attribute name(s).


**Returns** `Object` Key-value pairs of attribute values.
- - - -


{{< header href="attr_query" >}}

ATTR_QUERY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_query.go#L15)

ATTR_QUERY finds a single or more attribute(s) by an query selector.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target node.
`selector` | `String`  |  | Query selector.
`attrName` | `String, repeated`  |  | Attr name(s).


**Returns** `Object` Key-value pairs of attribute values.
- - - -


{{< header href="attr_remove" >}}

ATTR_REMOVE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_remove.go#L14)

ATTR_REMOVE removes single or more attribute(s) of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target node.
`attrNames` | `String, repeated`  |  | Attribute name(s).


**Returns** `None`
- - - -


{{< header href="attr_set" >}}

ATTR_SET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_set.go#L16)

ATTR_SET sets or updates a single or more attribute(s) of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target node.
`nameOrObj - Attribute name or an object representing a key` | `String` `Object`  |  | Value pair of attributes.
`value` | `String`  |  | If a second parameter is a string value, this parameter represent an attribute value.


**Returns** `None`
- - - -


{{< header href="blur" >}}

BLUR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/blur.go#L13)

BLUR Calls blur on the element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target node.
`selector` | `String`  |  | Css selector.


**Returns** `None`
- - - -


{{< header href="click" >}}

CLICK

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/click.go#L15)

CLICK dispatches click event on a given element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`cssSelectorOrClicks` | `String` `Int`  |  | Css selector or count of clicks.
`clicks` | `Int`  | `1` | Count of clicks.


**Returns** `None`
- - - -


{{< header href="click_all" >}}

CLICK_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/click_all.go#L16)

CLICK_ALL dispatches click event on all matched element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.
`clicks` | `Int`  | `1` | Optional count of clicks.


**Returns** `Boolean` True if matched at least one element.
- - - -


{{< header href="cookie_del" >}}

COOKIE_DEL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_del.go#L14)

COOKIE_DEL gets a cookie from a given page by name.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`cookiesOrNames` | `HTTPCookie, repeated` `String, repeated`  |  | Cookie or cookie name to delete.


**Returns** `None`
- - - -


{{< header href="cookie_get" >}}

COOKIE_GET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_get.go#L15)

COOKIE_GET gets a cookie from a given page by name.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`name` | `String`  |  | Cookie or cookie name to delete.


**Returns** `HTTPCookie` Cookie if found, otherwise none.
- - - -


{{< header href="cookie_set" >}}

COOKIE_SET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_set.go#L13)

COOKIE_SET sets cookies to a given page

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`cookies` | `HTTPCookie, repeated`  |  | Target cookies.


**Returns** `None`
- - - -


{{< header href="document" >}}

DOCUMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/document.go#L53)

DOCUMENT opens an HTML page by a given url. By default, loads a page by http call - resulted page does not support any interactions.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`params` | `Object`  |  | An object containing the following properties :
`params.driver` | `String`  |  | Driver name to use.
`params.timeout` | `Int`  | `60000` | Page load timeout.
`params.userAgent` | `String`  |  | Custom user agent.
`params.keepCookies` | `Boolean`  | `False` | Boolean value indicating whether to use cookies from previous sessions i.e. not to open a page in the incognito mode.
`params.cookies` | `Object[]` `Object`  |  | Set of http cookies to use during page loading.
`params.cookies.*.name` | `String`  |  | Cookie name.
`params.cookies.*.value` | `String`  |  | Cookie value.
`params.cookies.*.path` | `String`  |  | Cookie path.
`params.cookies.*.domain` | `String`  |  | Cookie domain.
`params.cookies.` | `Int`  | `*.maxAge` | Cookie max age.
`params.cookies.` | `String` `DateTime`  | `*.expires` | Cookie expiration date time.
`params.cookies.` | `String`  | `*.sameSite` | Origin policy.
`params.cookies.` | `Boolean`  | `*.httpOnly=false` | Cookie cannot be accessed through client side script.
`params.cookies.` | `Boolean`  | `*.secure=false` | Cookie sent to the server only with an encrypted request over the https protocol.
`params.headers` | `Object`  |  | Set of http headers to use during page loading.
`params.ignore` | `Object`  |  | Set of parameters to ignore some page functionality or behavior.
`params.ignore.resources` | `Object[]`  |  | Collection of rules to ignore resources during page load and navigation.
`params.ignore.resources.` | `String`  | `*.url` | > exactly one) are allowed. escape character is backslash. omitting is equivalent to "*".
`params.ignore.resources.` | `String`  | `*.type` | Resource type. if set, requests for matching resource types will be blocked.
`params.ignore.statusCodes` | `Object[]`  |  | Collection of rules to ignore certain http codes that can cause failures.
`params.ignore.statusCodes.` | `String`  | `*.url` | > exactly one) are allowed. escape character is backslash. omitting is equivalent to "*".
`params.ignore.statusCodes.` | `Int`  | `*.code` | Http code to ignore.
`params.viewport` | `Object`  |  | Viewport params.
`params.viewport.height` | `Int`  |  | Viewport height.
`params.viewport.width` | `Int`  |  | Viewport width.
`params.viewport.scaleFactor` | `Float`  |  | Viewport scale factor.
`params.viewport.mobile` | `Boolean`  |  | Value that indicates whether to emulate mobile device.
`params.viewport.landscape` | `Boolean`  |  | Value that indicates whether to render a page in landscape position.


**Returns** `HTMLPage` Loaded html page.
- - - -


{{< header href="document_exists" >}}

DOCUMENT_EXISTS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/document_exists.go#L16)

DOCUMENT_EXISTS returns a boolean value indicating whether a web page exists by a given url.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`url` | `String`  |  | Target url.
`options` | `Object`  |  | Request options.
`options.headers` | `Object`  |  | Request headers.


**Returns** `Boolean` A boolean value indicating whether a web page exists by a given url.
- - - -


{{< header href="download" >}}

DOWNLOAD

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/download.go#L15)

DOWNLOAD downloads a resource from the given GetURL.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`url` | `String`  |  | Url to download.


**Returns** `Binary` A base64 encoded string in binary format.
- - - -


{{< header href="element" >}}

ELEMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/element.go#L16)

ELEMENT finds an element by a given CSS selector. Returns NONE if element not found.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `HTMLElement` A matched html element
- - - -


{{< header href="elements" >}}

ELEMENTS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/elements.go#L14)

ELEMENTS finds HTML elements by a given CSS selector. Returns an empty array if element not found.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `HTMLElement[]` An array of matched html elements.
- - - -


{{< header href="elements_count" >}}

ELEMENTS_COUNT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/elements_count.go#L14)

ELEMENTS_COUNT returns a number of found HTML elements by a given CSS selector. Returns an empty array if element not found.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `Int` A number of matched html elements by a given css selector.
- - - -


{{< header href="element_exists" >}}

ELEMENT_EXISTS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/element_exists.go#L13)

ELEMENT_EXISTS returns a boolean value indicating whether there is an element matched by selector.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `Boolean` A boolean value indicating whether there is an element matched by selector.
- - - -


{{< header href="focus" >}}

FOCUS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/focus.go#L13)

FOCUS Sets focus on the element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `None`
- - - -


{{< header href="frames" >}}

FRAMES

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/find_frames.go#L15)

FRAMES finds HTML frames by a given property selector. Returns an empty array if frames not found.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Html page.
`property` | `String`  |  | Property selector.
`value` | `Any`  |  | Property value.


**Returns** `HTMLDocument[]` Returns an array of found html frames.
- - - -


{{< header href="hover" >}}

HOVER

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/hover.go#L15)

HOVER fetches an element with selector, scrolls it into view if needed, and then uses page.mouse to hover over the center of the element. If there's no element matching selector, the method returns an error.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | If document is passed, this param must represent an element selector.


**Returns** `None`
- - - -


{{< header href="inner_html" >}}

INNER_HTML

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_html.go#L15)

INNER_HTML returns inner HTML string of a given or matched by CSS selector element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.


**Returns** `String` Inner html string if a matched element, otherwise empty string.
- - - -


{{< header href="inner_html_all" >}}

INNER_HTML_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_html_all.go#L15)

INNER_HTML_ALL returns an array of inner HTML strings of matched elements.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.


**Returns** `String[]` An array of inner html strings if all matched elements, otherwise empty array.
- - - -


{{< header href="inner_html_set" >}}

INNER_HTML_SET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/set_inner_html.go#L15)

INNER_HTML_SET sets inner HTML string to a given or matched by CSS selector element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`htmlOrSelector` | `String`  |  | Html or css selector.
`html` | `String`  |  | String of inner html.


**Returns** `None`
- - - -


{{< header href="inner_text" >}}

INNER_TEXT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_text.go#L15)

INNER_TEXT returns inner text string of a given or matched by CSS selector element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.


**Returns** `String` Inner text if a matched element, otherwise empty string.
- - - -


{{< header href="inner_text_all" >}}

INNER_TEXT_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_text_all.go#L15)

INNER_TEXT_ALL returns an array of inner text of matched elements.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.


**Returns** `String[]` An array of inner text if all matched elements, otherwise empty array.
- - - -


{{< header href="inner_text_set" >}}

INNER_TEXT_SET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/set_inner_text.go#L15)

INNER_TEXT_SET sets inner text string to a given or matched by CSS selector element

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`textOrCssSelector` | `String`  |  | String of css selector.
`text` | `String`  |  | String of inner text.


**Returns** `None`
- - - -


{{< header href="input" >}}

INPUT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/input.go#L16)

INPUT types a value to an underlying input element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`valueOrSelector` | `String`  |  | Css selector or a value.
`value` | `String`  |  | Target value.
`delay` | `Int`  |  | Target value.


**Returns** `Boolean` Returns true if an element was found.
- - - -


{{< header href="input_clear" >}}

INPUT_CLEAR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/clear.go#L13)

INPUT_CLEAR clears a value from an underlying input element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector.


**Returns** `None`
- - - -


{{< header href="mouse" >}}

MOUSE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/mouse_xy.go#L15)

MOUSE moves mouse by given coordinates.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`document` | `HTMLDocument`  |  | Html document.
`x` | `Int` `Float`  |  | X coordinate.
`true` | `Int` `Float`  |  | Y coordinate.


**Returns** `None`
- - - -


{{< header href="navigate" >}}

NAVIGATE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate.go#L17)

NAVIGATE navigates a given page to a new resource. The operation blocks the execution until the page gets loaded. Which means there is no need in WAIT_NAVIGATION function.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`url` | `String`  |  | Target url to navigate.
`timeout` | `Int`  | `5000` | Navigation timeout.


**Returns** `None`
- - - -


{{< header href="navigate_back" >}}

NAVIGATE_BACK

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate_back.go#L18)

NAVIGATE_BACK navigates a given page back within its navigation history. The operation blocks the execution until the page gets loaded. If the history is empty, the function returns FALSE.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`entry` | `Int`  | `1` | An integer value indicating how many pages to skip.
`timeout` | `Int`  | `5000` | Navigation timeout.


**Returns** `Boolean` True if history exists and the operation succeeded, otherwise false.
- - - -


{{< header href="navigate_forward" >}}

NAVIGATE_FORWARD

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate_forward.go#L18)

NAVIGATE_FORWARD navigates a given page forward within its navigation history. The operation blocks the execution until the page gets loaded. If the history is empty, the function returns FALSE.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`entry` | `Int`  | `1` | An integer value indicating how many pages to skip.
`timeout` | `Int`  | `5000` | Navigation timeout.


**Returns** `Boolean` True if history exists and the operation succeeded, otherwise false.
- - - -


{{< header href="pagination" >}}

PAGINATION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/pagination.go#L16)

PAGINATION creates an iterator that goes through pages using CSS selector. The iterator starts from the current page i.e. it does not change the page on 1st iteration. That allows you to keep scraping logic inside FOR loop.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Css selector for a pagination on the page.


**Returns** `None`
- - - -


{{< header href="parse" >}}

PARSE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/parse.go#L31)

PARSE loads an HTML page from a given string or byte array

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`html` | `String`  |  | Html string to parse.
`params` | `Object`  |  | An object containing the following properties:
`params.driver` | `String`  |  | Name of a driver to parse with.
`params.keepCookies` | `Boolean`  | `False` | Boolean value indicating whether to use cookies from previous sessions i.e. not to open a page in the incognito mode.
`params.cookies` | `HTTPCookies`  |  | Set of http cookies to use during page loading.
`params.headers` | `HTTPHeaders`  |  | Set of http headers to use during page loading.
`params.viewport` | `Object`  |  | Viewport params.
`params.viewport.height` | `Int`  |  | Viewport height.
`params.viewport.width` | `Int`  |  | Viewport width.
`params.viewport.scaleFactor` | `Float`  |  | Viewport scale factor.
`params.viewport.mobile` | `Boolean`  |  | Value that indicates whether to emulate mobile device.
`params.viewport.landscape` | `Boolean`  |  | Value that indicates whether to render a page in landscape position.


**Returns** `HTMLPage` Returns parsed and loaded html page.
- - - -


{{< header href="pdf" >}}

PDF

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/pdf.go#L42)

PDF prints a PDF of the current page.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`target` | `HTMLPage` `String`  |  | Target page or url.
`params` | `Object`  |  | An object containing the following properties:
`params.landscape` | `Bool`  | `False` | Paper orientation.
`params.displayHeaderFooter` | `Bool`  | `False` | Display header and footer.
`params.printBackground` | `Bool`  | `False` | Print background graphics.
`params.scale` | `Float`  | `1` | Scale of the webpage rendering.
`params.paperWidth` | `Float`  | `22` | Paper width in inches.
`params.paperHeight` | `Float`  | `28` | Paper height in inches.
`params.marginTo` | `Float`  | `1` | Top margin in inches.
`params.marginBottom` | `Float`  | `1` | Bottom margin in inches.
`params.marginLeft` | `Float`  | `1` | Left margin in inches.
`params.marginRight` | `Float`  | `1` | Right margin in inches.
`params.pageRanges` | `String`  |  | 13'.
`params.ignoreInvalidPageRanges` | `Bool`  | `False` | 2'.
`params.headerTemplate` | `String`  |  | `totalpages`: total pages in the document for example, `<span class=title></span>` would generate span containing the title.
`params.footerTemplate` | `String`  |  | Html template for the print footer. should use the same format as the `headertemplate`.
`params.preferCSSPageSize` | `Bool`  | `False` | Whether or not to prefer page size as defined by css. defaults to false, in which case the content will be scaled to fit the paper size. *


**Returns** `Binary` Pdf document in binary format.
- - - -


{{< header href="screenshot" >}}

SCREENSHOT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/screenshot.go#L22)

SCREENSHOT takes a screenshot of a given page.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`target` | `HTMLPage` `String`  |  | Target page or url.
`params` | `Object`  |  | An object containing the following properties :
`params.x` | `Float` `Int`  | `0` | X position of the viewport.
`params.y` | `Float` `Int`  | `0` | Y position of the viewport.
`params.width` | `Float` `Int`  |  | Width of the viewport.
`params.height` | `Float` `Int`  |  | Height of the viewport.
`params.format` | `String`  | `"jpeg"` | Either "jpeg" or "png".
`params.quality` | `Int`  | `100` | Quality, in [0, 100], only for jpeg format.


**Returns** `Binary` Screenshot in binary format.
- - - -


{{< header href="scroll" >}}

SCROLL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_xy.go#L19)

SCROLL scrolls by given coordinates.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`document` | `HTMLDocument`  |  | Html document.
`x` | `Int` `Float`  |  | X coordinate.
`true` | `Int` `Float`  |  | Y coordinate.
`params` | `Object`  |  | Scroll params.
`params.behavior` | `String`  | `"instant"` | Scroll behavior
`params.block` | `String`  | `"center"` | Scroll vertical alignment.
`params.inline` | `String`  | `"center"` | Scroll horizontal alignment.


**Returns** `None`
- - - -


{{< header href="scroll_bottom" >}}

SCROLL_BOTTOM

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_bottom.go#L18)

SCROLL_BOTTOM scrolls the document's window to its bottom.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`document` | `HTMLDocument`  |  | Html document.
`x` | `Int` `Float`  |  | X coordinate.
`true` | `Int` `Float`  |  | Y coordinate.
`params` | `Object`  |  | Scroll params.
`params.behavior` | `String`  | `"instant"` | Scroll behavior
`params.block` | `String`  | `"center"` | Scroll vertical alignment.
`params.inline` | `String`  | `"center"` | Scroll horizontal alignment.


**Returns** `None`
- - - -


{{< header href="scroll_element" >}}

SCROLL_ELEMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_element.go#L20)

SCROLL_ELEMENT scrolls an element on.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | If document is passed, this param must represent an element selector.
`params` | `Object`  |  | Scroll params.
`params.behavior` | `String`  | `"instant"` | Scroll behavior
`params.block` | `String`  | `"center"` | Scroll vertical alignment.
`params.inline` | `String`  | `"center"` | Scroll horizontal alignment.


**Returns** `None`
- - - -


{{< header href="scroll_top" >}}

SCROLL_TOP

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_top.go#L18)

SCROLL_TOP scrolls the document's window to its top.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`document` | `HTMLDocument`  |  | Html document.
`x` | `Int` `Float`  |  | X coordinate.
`true` | `Int` `Float`  |  | Y coordinate.
`params` | `Object`  |  | Scroll params.
`params.behavior` | `String`  | `"instant"` | Scroll behavior
`params.block` | `String`  | `"center"` | Scroll vertical alignment.
`params.inline` | `String`  | `"center"` | Scroll horizontal alignment.


**Returns** `None`
- - - -


{{< header href="select" >}}

SELECT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/select.go#L15)

SELECT selects a value from an underlying select element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`element` | `HTMLElement`  |  | Target html element.
`valueOrSelector` | `String` `String[]`  |  | Selector or a an array of strings as a value.
`value` | `String[]`  |  | Target value. optional.


**Returns** `String[]` Array of selected values.
- - - -


{{< header href="style_get" >}}

STYLE_GET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_get.go#L14)

STYLE_GET gets single or more style attribute value(s) of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`element` | `HTMLElement`  |  | Target html element.
`names` | `String, repeated`  |  | Style name(s).


**Returns** `Object` Collection of key-value pairs of style values.
- - - -


{{< header href="style_remove" >}}

STYLE_REMOVE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_remove.go#L14)

STYLE_REMOVE removes single or more style attribute value(s) of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`element` | `HTMLElement`  |  | Target html element.
`names` | `String, repeated`  |  | Style name(s).


**Returns** `None`
- - - -


{{< header href="style_set" >}}

STYLE_SET

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_set.go#L15)

STYLE_SET sets or updates a single or more style attribute value of a given element.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`element` | `HTMLElement`  |  | Target html element.
`nameOrObj - Style name or an object representing a key` | `String` `Object`  |  | Value pair of attributes.
`value` | `String`  |  | If a second parameter is a string value, this parameter represent a style value.


**Returns** `None`
- - - -


{{< header href="wait_attr" >}}

WAIT_ATTR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr.go#L17)

WAIT_ATTR waits until a target attribute's value appears

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`attrNameOrSelector` | `String`  |  | String of an attr name or css selector.
`attrValueOrAttrName` | `String` `Any`  |  | Attr value or name.
`attrValueOrTimeout` | `Any` `Int`  |  | Attr value or a timeout.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_attr_all" >}}

WAIT_ATTR_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr_all.go#L17)

WAIT_ATTR_ALL waits for an attribute to appear on all matched elements with a given value. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.
`class` | `String`  |  | String of target css class.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_class" >}}

WAIT_CLASS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class.go#L17)

WAIT_CLASS waits for a class to appear on a given element. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selectorOrClass` | `String`  |  | If document is passed, this param must represent an element selector. otherwise target class.
`classOrTimeout` | `String` `Int`  |  | If document is passed, this param must represent target class name. otherwise timeout.
`timeout` | `Int`  |  | If document is passed, this param must represent timeout. otherwise not passed.


**Returns** `None`
- - - -


{{< header href="wait_class_all" >}}

WAIT_CLASS_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class_all.go#L17)

WAIT_CLASS_ALL waits for a class to appear on all matched elements. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.
`class` | `String`  |  | String of target css class.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_element" >}}

WAIT_ELEMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_element.go#L16)

WAIT_ELEMENT waits for element to appear in the DOM. Stops the execution until it finds an element or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Target element's selector.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_navigation" >}}

WAIT_NAVIGATION

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_navigation.go#L26)

WAIT_NAVIGATION waits for a given page to navigate to a new url. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`page` | `HTMLPage`  |  | Target page.
`timeout` | `Int`  | `5000` | Navigation timeout.
`params` | `Object`  | `None` | Navigation parameters.
`params.timeout` | `Int`  | `5000` | Navigation timeout.
`params.target` | `String`  |  | Navigation target url.
`params.frame` | `HTMLDocument`  |  | Navigation frame.


**Returns** `None`
- - - -


{{< header href="wait_no_attr" >}}

WAIT_NO_ATTR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr.go#L27)

WAIT_NO_ATTR waits until a target attribute's value disappears

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`attrNameOrSelector` | `String`  |  | String of an attr name or css selector.
`attrValueOrAttrName` | `String` `Any`  |  | Attr value or name.
`attrValueOrTimeout` | `Any` `Int`  |  | Attr value or wait timeout.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_no_attr_all" >}}

WAIT_NO_ATTR_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr_all.go#L27)

WAIT_NO_ATTR_ALL waits for an attribute to disappear on all matched elements by a given value. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.
`class` | `String`  |  | String of target css class.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_no_class" >}}

WAIT_NO_CLASS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class.go#L27)

WAIT_NO_CLASS waits for a class to disappear on a given element. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selectorOrClass` | `String`  |  | If document is passed, this param must represent an element selector. otherwise target class.
`classOrTimeout` | `String` `Int`  |  | If document is passed, this param must represent target class name. otherwise timeout.
`timeout` | `Int`  |  | If document is passed, this param must represent timeout. otherwise not passed.


**Returns** `None`
- - - -


{{< header href="wait_no_class_all" >}}

WAIT_NO_CLASS_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class_all.go#L27)

WAIT_NO_CLASS_ALL waits for a class to disappear on all matched elements. Stops the execution until the navigation ends or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | String of css selector.
`class` | `String`  |  | String of target css class.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_no_element" >}}

WAIT_NO_ELEMENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_element.go#L25)

WAIT_NO_ELEMENT waits for element to disappear in the DOM. Stops the execution until it does not find an element or operation times out.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`selector` | `String`  |  | Target element's selector.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_no_style" >}}

WAIT_NO_STYLE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style.go#L27)

WAIT_NO_STYLE waits until a target style value disappears

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`styleNameOrSelector` | `String`  |  | Style name or css selector.
`valueOrStyleName` | `String` `Any`  |  | Style value or name.
`valueOrTimeout` | `Any` `Int`  |  | Style value or wait timeout.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_no_style_all" >}}

WAIT_NO_STYLE_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style_all.go#L27)

WAIT_NO_STYLE_ALL waits until a target style value disappears on all matched elements with a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`styleNameOrSelector` | `String`  |  | Style name or css selector.
`valueOrStyleName` | `String` `Any`  |  | Style value or name.
`valueOrTimeout` | `Any` `Int`  |  | Style value or wait timeout.
`timeout` | `Int`  | `5000` | Timeout.


**Returns** `None`
- - - -


{{< header href="wait_style" >}}

WAIT_STYLE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style.go#L17)

WAIT_STYLE waits until a target style value appears

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`styleNameOrSelector` | `String`  |  | Style name or css selector.
`valueOrStyleName` | `String` `Any`  |  | Style value or name.
`valueOrTimeout` | `Any` `Int`  |  | Style value or wait timeout.
`timeout` | `Int`  | `5000` | Wait timeout.


**Returns** `None`
- - - -


{{< header href="wait_style_all" >}}

WAIT_STYLE_ALL

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style_all.go#L17)

WAIT_STYLE_ALL waits until a target style value appears on all matched elements with a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`styleNameOrSelector` | `String`  |  | Style name or css selector.
`valueOrStyleName` | `String` `Any`  |  | Style value or name.
`valueOrTimeout` | `Any` `Int`  |  | Style value or wait timeout.
`timeout` | `Int`  | `5000` | Timeout.


**Returns** `None`
- - - -


{{< header href="xpath" >}}

XPATH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/xpath.go#L14)

XPATH evaluates the XPath expression.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement`  |  | Target html node.
`expression` | `String`  |  | Xpath expression.


**Returns** `Any` Returns result of a given xpath expression.
- - - -
