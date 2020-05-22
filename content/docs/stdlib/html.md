---
title: "html"
weight: 1
draft: false
---


## ELEMENT_EXISTS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/element_exists.go#L13)

ELEMENT_EXISTS returns a boolean value indicating whether there is an element matched by selector.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLNode` | Parent document or element.
`selector` | `String` | Css selector.


**Returns** `Boolean` A boolean value indicating whether there is an element matched by selector.
- - - -

## BLUR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/blur.go#L13)

BLUR Calls blur on the element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`target` | `HTMLPage` `HTMLDocument` `HTMLElement` | Target node.
`selector` | `String, optional` | Optional css selector.


**Returns** `None`
- - - -

## INNER_TEXT_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_text_all.go#L15)

INNER_TEXT_ALL returns an array of inner text of matched elements.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` `HTMLElement` | Parent document or element.
`selector` | `String` | String of css selector.


**Returns** `String` An array of inner text if any element found, otherwise empty array.
- - - -

## ELEMENTS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/elements.go#L14)

ELEMENTS finds HTML elements by a given CSS selector. Returns an empty array if element not found.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLNode` | Parent document or element.
`selector` | `String` | Css selector.


**Returns** `Array` Returns an array of found html element.
- - - -

## WAIT_ELEMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_element.go#L16)

WAIT_ELEMENT waits for element to appear in the DOM. Stops the execution until it finds an element or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`false` | `HTMLDocument` | Driver htmldocument.
`selector` | `String` | Target element's selector.
`timeout` | `Int, optional` | Optional timeout. default 5000 ms.


**Returns** `None`
- - - -

## WAIT_NO_ELEMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_element.go#L25)

WAIT_NO_ELEMENT waits for element to disappear in the DOM. Stops the execution until it does not find an element or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Driver htmldocument.
`selector` | `String` | Target element's selector.
`timeout` | `Int, optional` | Optional timeout. default 5000 ms.


**Returns** `None`
- - - -

## NAVIGATE_BACK
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate_back.go#L18)

NAVIGATE_BACK navigates a given page back within its navigation history. The operation blocks the execution until the page gets loaded. If the history is empty, the function returns FALSE.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`entry` | `Int, optional` | Optional value indicating how many pages to skip. default 1.
`timeout` | `Int, optional` | Optional timeout. default is 5000.


**Returns** `Boolean` Returns true if history exists and the operation succeeded, otherwise false.
- - - -

## INPUT_CLEAR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/clear.go#L13)

INPUT_CLEAR clears a value from an underlying input element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `HTMLPage` `HTMLDocument` `HTMLElement` | Event target.
`selector` | `String, options` | Selector.


**Returns** `None`
- - - -

## SELECT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/select.go#L15)

SELECT selects a value from an underlying select element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `Open` `GetElement` | Event target.
`valueOrSelector` | `String` `Array<String>` | Selector or a an array of strings as a value.
`value` | `Array<String` | Target value. optional.


**Returns** `Array<String>` Returns an array of selected values.
- - - -

## WAIT_NAVIGATION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_navigation.go#L21)

WAIT_NAVIGATION waits for a given page to navigate to a new url. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`timeout` | `Int, optional` | Optional timeout. default 5000 ms.


**Returns** `None`
- - - -

## PDF
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/pdf.go#L42)

PDF prints a PDF of the current page.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`target` | `HTMLPage` `String` | Target page or url.
`params (Object) - Optional, An object containing the following properties : Landscape (Bool) - Paper orientation. Defaults to false. DisplayHeaderFooter (Bool) - Display header and footer. Defaults to false. PrintBackground (Bool) - Print background graphics. Defaults to false. Scale (Float64) - Scale of the webpage rendering. Defaults to 1. PaperWidth (Float64) - Paper width in inches. Defaults to 8.5 inches. PaperHeight (Float64) - Paper height in inches. Defaults to 11 inches. MarginTop (Float64) - Top margin in inches. Defaults to 1cm (~0.4 inches). MarginBottom (Float64) - Bottom margin in inches. Defaults to 1cm (~0.4 inches). MarginLeft (Float64) - Left margin in inches. Defaults to 1cm (~0.4 inches). MarginRight (Float64) - Right margin in inches. Defaults to 1cm (~0.4 inches). PageRanges (String) - Paper ranges to print, e.g., '1-5, 8, 11-13'. Defaults to the empty string, which means print all pages. IgnoreInvalidPageRanges (Bool) - to silently ignore invalid but successfully parsed page ranges, such as '3-2'. Defaults to false. HeaderTemplate (String) - HTML template for the print header. Should be valid HTML markup with following classes used to inject printing values into them: - `date`: formatted print date - `title`: document title - `url`: document location - `pageNumber`: current page number - `totalPages`: total pages in the document For example, `<span class=title></span>` would generate span containing the title. FooterTemplate (String) - HTML template for the print footer. Should use the same format as the `headerTemplate`. PreferCSSPageSize` | `Bool` | Whether or not to prefer page size as defined by css. defaults to false, in which case the content will be scaled to fit the paper size. *


**Returns** `Binary` Returns a base64 encoded string in binary format.
- - - -

## PAGINATION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/pagination.go#L16)

PAGINATION creates an iterator that goes through pages using CSS selector. The iterator starts from the current page i.e. it does not change the page on 1st iteration. That allows you to keep scraping logic inside FOR loop.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `Open` | Target document.
`selector` | `String` | Css selector for a pagination on the page.


**Returns** `None`
- - - -

## STYLE_GET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_get.go#L14)

STYLE_GET gets single or more style attribute value(s) of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`names` | `...String` | Style name(s).


**Returns** `Object` Key-value pairs of style values.
- - - -

## SCROLL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_xy.go#L16)

SCROLL scrolls by given coordinates.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Html document.
`x` | `Int` `Float` | X coordinate.
`true` | `Int` `Float` | Y coordinate.
`options` | `ScrollOptions` | Scroll options. optional.


**Returns** `None`
- - - -

## ATTR_REMOVE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_remove.go#L14)

ATTR_REMOVE removes single or more attribute(s) of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`names` | `...String` | Attribute name(s).


**Returns** `None`
- - - -

## ATTR_SET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_set.go#L15)

ATTR_SET sets or updates a single or more attribute(s) of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`nameOrObj` | `String` `Object` | Attribute name or an object representing a key-value pair of attributes.
`value` | `String` | If a second parameter is a string value, this parameter represent an attribute value.


**Returns** `None`
- - - -

## WAIT_STYLE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style.go#L12)

WAIT_STYLE

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `None`
- - - -

## WAIT_NO_STYLE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style.go#L17)

WAIT_NO_STYLE

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `None`
- - - -

## WAIT_ATTR_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr_all.go#L17)

WAIT_ATTR_ALL waits for an attribute to appear on all matched elements with a given value. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Parent document.
`selector` | `String` | String of css selector.
`class` | `String` | String of target css class.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## WAIT_NO_ATTR_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr_all.go#L27)

WAIT_NO_ATTR_ALL waits for an attribute to disappear on all matched elements by a given value. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Parent document.
`selector` | `String` | String of css selector.
`class` | `String` | String of target css class.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## COOKIE_DEL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_del.go#L14)

COOKIE_DEL gets a cookie from a given page by name.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`cookie` | `...HTTPCookie` `String` | Cookie or cookie name to delete.


**Returns** `None`
- - - -

## HOVER
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/hover.go#L15)

HOVER fetches an element with selector, scrolls it into view if needed, and then uses page.mouse to hover over the center of the element. If there's no element matching selector, the method returns an error.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLElement` | Target document or element.
`selector` | `String, options` | If document is passed, this param must represent an element selector.


**Returns** `None`
- - - -

## INNER_HTML
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_html.go#L15)

INNER_HTML returns inner HTML string of a given or matched by CSS selector element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `Open` `GetElement` | Parent document or element.
`selector` | `String, optional` | String of css selector.


**Returns** `String` Inner html string if an element found, otherwise empty string.
- - - -

## INNER_HTML_SET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/set_inner_html.go#L15)

INNER_HTML_SET sets inner HTML string to a given or matched by CSS selector element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `Open` `GetElement` | Parent document or element.
`selector` | `String, optional` | String of css selector.
`innerHTML` | `String` | String of inner html.


**Returns** `None`
- - - -

## ELEMENTS_COUNT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/elements_count.go#L14)

ELEMENTS_COUNT returns a number of found HTML elements by a given CSS selector. Returns an empty array if element not found.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLNode` | Parent document or element.
`selector` | `String` | Css selector.


**Returns** `Int` A number of found html elements by a given css selector.
- - - -

## NAVIGATE_FORWARD
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate_forward.go#L18)

NAVIGATE_FORWARD navigates a given page forward within its navigation history. The operation blocks the execution until the page gets loaded. If the history is empty, the function returns FALSE.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`entry` | `Int, optional` | Optional value indicating how many pages to skip. default 1.
`timeout` | `Int, optional` | Optional timeout. default is 5000.


**Returns** `Boolean` Returns true if history exists and the operation succeeded, otherwise false.
- - - -

## MOUSE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/mouse_xy.go#L15)

MOUSE moves mouse by given coordinates.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Html document.
`x` | `Int` `Float` | X coordinate.
`true` | `Int` `Float` | Y coordinate.


**Returns** `None`
- - - -

## SCREENSHOT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/screenshot.go#L22)

SCREENSHOT takes a screenshot of a given page.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`target` | `HTMLPage` `String` | Target page or url.
`params (Object) - Optional, An object containing the following properties : x (Float|Int) - Optional, X position of the viewport. x (Float|Int) - Optional,Y position of the viewport. width (Float|Int) - Optional, Width of the viewport. height (Float|Int) - Optional, Height of the viewport. format (String) - Optional, Either "jpeg" or "png". quality` | `Int` | Optional, quality, in [0, 100], only for jpeg format.


**Returns** `Binary` Returns a base64 encoded string in binary format.
- - - -

## STYLE_REMOVE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_remove.go#L14)

STYLE_REMOVE removes single or more style attribute value(s) of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`names` | `...String` | Style name(s).


**Returns** `None`
- - - -

## WAIT_CLASS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class.go#L20)

WAIT_CLASS waits for a class to appear on a given element. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement` | Target node.
`selectorOrClass` | `String` | If document is passed, this param must represent an element selector. otherwise target class.
`classOrTimeout` | `String` `Int, optional` | If document is passed, this param must represent target class name. otherwise timeout.
`timeout` | `Int, optional` | If document is passed, this param must represent timeout. otherwise not passed.


**Returns** `None`
- - - -

## WAIT_NO_CLASS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class.go#L33)

WAIT_NO_CLASS waits for a class to disappear on a given element. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement` | Target node.
`selectorOrClass` | `String` | If document is passed, this param must represent an element selector. otherwise target class.
`classOrTimeout` | `String` `Int, optional` | If document is passed, this param must represent target class name. otherwise timeout.
`timeout` | `Int, optional` | If document is passed, this param must represent timeout. otherwise not passed.


**Returns** `None`
- - - -

## SCROLL_TOP
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_top.go#L13)

SCROLL_TOP scrolls the document's window to its top.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Target document.
`options` | `ScrollOptions` | Scroll options. optional.


**Returns** `None`
- - - -

## INNER_TEXT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_text.go#L15)

INNER_TEXT returns inner text string of a given or matched by CSS selector element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` `HTMLElement` | Parent document or element.
`selector` | `String, optional` | String of css selector.


**Returns** `String` Inner text if an element found, otherwise empty string.
- - - -

## ATTR_GET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/attr_get.go#L14)

ATTR_GET gets single or more attribute(s) of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`names` | `...String` | Attribute name(s).


**Returns** `Object` Key-value pairs of attribute values.
- - - -

## FOCUS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/focus.go#L13)

FOCUS Sets focus on the element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`target` | `HTMLPage` `HTMLDocument` `HTMLElement` | Target node.
`selector` | `String, optional` | Optional css selector.


**Returns** `None`
- - - -

## INNER_HTML_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/get_inner_html_all.go#L15)

INNER_HTML_ALL returns an array of inner HTML strings of matched elements.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` `HTMLElement` | Parent document or element.
`selector` | `String` | String of css selector.


**Returns** `String` An array of inner html strings if any element found, otherwise empty array.
- - - -

## SCROLL_BOTTOM
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_bottom.go#L13)

SCROLL_BOTTOM scrolls the document's window to its bottom.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Target document.
`options` | `ScrollOptions` | Scroll options. optional.


**Returns** `None`
- - - -

## COOKIE_SET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_set.go#L13)

COOKIE_SET sets cookies to a given page

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`cookie...` | `HTTPCookie` | Target cookies.


**Returns** `None`
- - - -

## WAIT_CLASS_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class_all.go#L17)

WAIT_CLASS_ALL waits for a class to appear on all matched elements. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Parent document.
`selector` | `String` | String of css selector.
`class` | `String` | String of target css class.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## WAIT_NO_CLASS_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_class_all.go#L27)

WAIT_NO_CLASS_ALL waits for a class to disappear on all matched elements. Stops the execution until the navigation ends or operation times out.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `HTMLDocument` | Parent document.
`selector` | `String` | String of css selector.
`class` | `String` | String of target css class.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## WAIT_STYLE_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style_all.go#L12)

WAIT_STYLE_ALL

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `None`
- - - -

## WAIT_NO_STYLE_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_style_all.go#L17)

WAIT_NO_STYLE_ALL

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `None`
- - - -

## COOKIE_GET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/cookie_get.go#L14)

COOKIE_GET gets a cookie from a given page by name.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`name` | `String` | Cookie or cookie name to delete.


**Returns** `None`
- - - -

## XPATH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/xpath.go#L14)

XPATH evaluates the XPath expression.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `HTMLPage` `HTMLDocument` `HTMLElement` | Target html object.
`expression` | `String` | Xpath expression.


**Returns** `Value` Returns result of a given xpath expression.
- - - -

## PARSE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/parse.go#L26)

PARSE loads an HTML page from a given string or byte array

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - Optional, an object containing the following properties : driver (String) - Optional, driver name. keepCookies (Boolean) - Optional, boolean value indicating whether to use cookies from previous sessions. i.e. not to open a page in the Incognito mode. cookies (HTTPCookies) - Optional, set of HTTP cookies. headers (HTTPHeaders) - Optional, HTTP headers. viewport` | `Viewport` | Optional, viewport params.


**Returns** `HTMLPage` Returns parsed and loaded html page.
- - - -

## INNER_TEXT_SET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/set_inner_text.go#L15)

INNER_TEXT_SET sets inner text string to a given or matched by CSS selector element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`doc` | `Open` `GetElement` | Parent document or element.
`selector` | `String, optional` | String of css selector.
`innerText` | `String` | String of inner text.


**Returns** `None`
- - - -

## ELEMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/element.go#L16)

ELEMENT finds an element by a given CSS selector. Returns NONE if element not found.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLElement` | Parent document or element.
`selector` | `String` | Css selector.


**Returns** `HTMLElement` `None` Returns an htmlelement if found, otherwise none.
- - - -

## CLICK
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/click.go#L15)

CLICK dispatches click event on a given element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `Open` `GetElement` | Event source.
`selectorOrCount` | `String` `Int, optional` | Optional selector or count of clicks.
`count` | `Int, optional` | Optional count of clicks.


**Returns** `None`
- - - -

## CLICK_ALL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/click_all.go#L16)

CLICK_ALL dispatches click event on all matched element

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `Open` | Open.
`selector` | `String` | Selector.
`count` | `Int, optional` | Optional count of clicks.


**Returns** `Boolean` Returns true if matched at least one element.
- - - -

## DOWNLOAD
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/download.go#L15)

Download downloads a resource from the given GetURL.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`GetURL` | `String` | Geturl to download.


**Returns** `Binary` Returns a base64 encoded string in binary format.
- - - -

## INPUT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/input.go#L16)

INPUT types a value to an underlying input element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`source` | `HTMLPage` `HTMLDocument` `HTMLElement` | Event target.
`valueOrSelector` | `String` | Selector or a value.
`value` | `String` | Target value.
`delay` | `Int, optional` | Target value.


**Returns** `Boolean` Returns true if an element was found.
- - - -

## NAVIGATE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/navigate.go#L17)

NAVIGATE navigates a given page to a new resource. The operation blocks the execution until the page gets loaded. Which means there is no need in WAIT_NAVIGATION function.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`page` | `HTMLPage` | Target page.
`url` | `String` | Target url to navigate.
`timeout` | `Int, optional` | Optional timeout. default is 5000.


**Returns** `None`
- - - -

## STYLE_SET
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/style_set.go#L15)

STYLE_SET sets or updates a single or more style attribute value of a given element.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`el` | `HTMLElement` | Target element.
`nameOrObj` | `String` `Object` | Style name or an object representing a key-value pair of attributes.
`value` | `String` | If a second parameter is a string value, this parameter represent a style value.


**Returns** `None`
- - - -

## WAIT_ATTR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr.go#L17)

WAIT_ATTR waits until a target attribute's value appears

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement` | Parent document.
`attrNameOrSelector` | `String` | String of an attr name or css selector.
`attrValueOrAttrName` | `String` `Any` | Attr value or name.
`attrValueOrTimeout` | `Any` `Int, optional` | Attr value or an optional timeout.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## WAIT_NO_ATTR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/wait_attr.go#L27)

WAIT_NO_ATTR waits until a target attribute's value disappears

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`node` | `HTMLPage` `HTMLDocument` `HTMLElement` | Parent document.
`attrNameOrSelector` | `String` | String of an attr name or css selector.
`attrValueOrAttrName` | `String` `Any` | Attr value or name.
`attrValueOrTimeout` | `Any` `Int, optional` | Attr value or an optional timeout.
`timeout` | `Int, optional` | Optional timeout.


**Returns** `None`
- - - -

## DOCUMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/document.go#L32)

DOCUMENT opens an HTML page by a given url. By default, loads a page by http call - resulted page does not support any interactions.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`params (Object) - Optional, An object containing the following properties : driver (String) - Optional, driver name. timeout (Int) - Optional, timeout. userAgent (String) - Optional, user agent. keepCookies (Boolean) - Optional, boolean value indicating whether to use cookies from previous sessions. i.e. not to open a page in the Incognito mode. cookies (HTTPCookies) - Optional, set of HTTP cookies. headers (HTTPHeaders) - Optional, HTTP headers. viewport` | `Viewport` | Optional, viewport params.


**Returns** `HTMLPage` Returns loaded html page.
- - - -

## SCROLL_ELEMENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/html/scroll_element.go#L17)

SCROLL_ELEMENT scrolls an element on.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`docOrEl` | `HTMLDocument` `HTMLElement` | Target document or element.
`selector` | `String` | If document is passed, this param must represent an element selector.
`options` | `ScrollOptions` | Scroll options. optional.


**Returns** `None`
- - - -
