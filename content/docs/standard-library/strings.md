---
title: "Strings"
weight: 40
draft: false
description: "String functions in the Ferret standard library."
aliases:
  - /docs/stdlib/strings/
menuTitle: 
menu: [concat,concat_separator,contains,decode_uri_component,encode_uri_component,escape_html,find_first,find_last,fmt,from_base64,json_parse,json_stringify,left,like,lower,ltrim,md5,random_token,regex_match,regex_replace,regex_split,regex_test,right,rtrim,sha1,sha512,split,substitute,substring,to_base64,trim,unescape_html,upper,]
---



{{< header href="concat" >}}

concat

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L13)

concat concatenates one or more instances of String, or an Array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`src` | `String, repeated` `String[]`  |  | The source string / array.


**Returns** `String` A string value.
- - - -


{{< header href="concat_separator" >}}

concat_separator

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L47)

concat_separator concatenates one or more instances of String, or an Array with a given separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`separator` | `String`  |  | The separator string.
`src` | `String, repeated` `String[]`  |  | The source string / array.


**Returns** `String` Concatenated string.
- - - -


{{< header href="contains" >}}

contains

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/contains.go#L15)

contains returns a value indicating whether a specified substring occurs within a string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`returnIndex` | `Boolean`  | `False` | Values which indicates whether to return the character position of the match is returned instead of a boolean.


**Returns** `Boolean` `Int` A value indicating whether a specified substring occurs within a string.
- - - -


{{< header href="decode_uri_component" >}}

decode_uri_component

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L36)

decode_uri_component returns the decoded String of uri.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to decode.


**Returns** `String` Decoded string.
- - - -


{{< header href="encode_uri_component" >}}

encode_uri_component

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L17)

encode_uri_component returns the encoded String of uri.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to encode.


**Returns** `String` Encoded string.
- - - -


{{< header href="escape_html" >}}

escape_html

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/escape.go#L16)

"escape_html escapes special characters like \"<\" to become \"&lt;\"\. It escapes only five such characters: <, >, &, ' and \". UnescapeString(EscapeString(s))\\ == s always holds, but the converse isn't always true."


|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to escape.


**Returns** `String` Escaped string.
- - - -


{{< header href="find_first" >}}

find_first

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L17)

find_first returns the position of the first occurrence of the string search inside the string text. Positions start at 0.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`start` | `Int`  |  | Limit the search to a subset of the text, beginning at start.
`end` | `Int`  |  | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -


{{< header href="find_last" >}}

find_last

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L63)

find_last returns the position of the last occurrence of the string search inside the string text. Positions start at 0.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`src` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`start` | `Int`  |  | Limit the search to a subset of the text, beginning at start.
`end` | `Int`  |  | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -


{{< header href="fmt" >}}

fmt

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/fmt.go#L18)

fmt formats the template using these arguments.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`template` | `String`  |  | Template.
`args` | `Any, repeated`  |  | Template arguments.


**Returns** `String` String formed by template using arguments.
- - - -


{{< header href="from_base64" >}}

from_base64

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L16)

from_base64 returns the value of a base64 representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to decode.


**Returns** `String` The decoded string.
- - - -


{{< header href="json_parse" >}}

json_parse

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L15)

json_parse returns a value described by the JSON-encoded input string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to parse as json.


**Returns** `Any` Parsed value.
- - - -


{{< header href="json_stringify" >}}

json_stringify

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L36)

json_stringify returns a JSON string representation of the input value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `Any`  |  | The input value to serialize.


**Returns** `String` Json string.
- - - -


{{< header href="left" >}}

left

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L61)

left returns the leftmost characters of the string value by index.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`length` | `Int`  |  | The amount of characters to return.


**Returns** `String` The leftmost characters of the string value by index.
- - - -


{{< header href="like" >}}

like

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/like.go#L22)

like checks whether the pattern search is contained in the string text, using wildcard matching.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to search in.
`search` | `String`  |  | A search pattern that can contain the wildcard characters.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -


{{< header href="lower" >}}

lower

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L13)

lower converts strings to their lower-case counterparts. All other characters are returned unchanged.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.


**Returns** `String` This string in lower case.
- - - -


{{< header href="ltrim" >}}

ltrim

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L34)

ltrim returns the string value with whitespace stripped from the start only.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the left-hand side.
- - - -


{{< header href="md5" >}}

md5

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L32)

md5 calculates the md5 checksum for text and return it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Md5 checksum as hex string.
- - - -


{{< header href="random_token" >}}

random_token

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/random.go#L26)

random_token generates a pseudo-random token string with the specified length. The algorithm for token generation should be treated as opaque.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`len` | `Int`  |  | The desired string length for the token. it must be greater than 0 and at most 65536.


**Returns** `String` A generated token consisting of lowercase letters, uppercase letters and numbers.
- - - -


{{< header href="regex_match" >}}

regex_match

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L16)

regex_match returns the matches in the given string text, using the regex.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to search in.
`expression` | `String`  |  | A regular expression to use for matching the text.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.


**Returns** `Any[]` An array of strings containing the matches.
- - - -


{{< header href="regex_replace" >}}

regex_replace

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L133)

regex_replace replace every substring matched with the regexp with a given string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`expression` | `String`  |  | A regular expression search pattern.
`replacement` | `String`  |  | The string to replace the search pattern with
`caseInsensitive` | `Boolean`  | `False` | Insensitive.


**Returns** `String` Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
- - - -


{{< header href="regex_split" >}}

regex_split

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L58)

regex_split splits the given string text into a list of strings, using the separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`expression` | `String`  |  | A regular expression to use for splitting the text.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.
`limit` | `Int`  |  | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `Any[]` An array of strings splitted by the expression.
- - - -


{{< header href="regex_test" >}}

regex_test

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L100)

regex_test test whether the regexp has at least one match in the given text.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to test.
`expression` | `String`  |  | A regular expression to use for splitting the text.
`caseInsensitive` | `Boolean`  | `False` | Insensitive.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -


{{< header href="right" >}}

right

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L88)

right returns the rightmost characters of the string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`length` | `Int`  |  | The amount of characters to return.


**Returns** `String` The rightmost characters of the string value.
- - - -


{{< header href="rtrim" >}}

rtrim

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L55)

rtrim returns the string value with whitespace stripped from the end only.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the right-hand side.
- - - -


{{< header href="sha1" >}}

sha1

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L48)

sha1 calculates the sha1 checksum for text and returns it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Sha1 checksum as hex string.
- - - -


{{< header href="sha512" >}}

sha512

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L64)

sha512 calculates the sha512 checksum for text and returns it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Sha512 checksum as hex string.
- - - -


{{< header href="split" >}}

split

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/split.go#L16)

split splits the given string value into a list of strings, using the separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`separator` | `String`  |  | The separator.
`limit` | `Int`  |  | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `String[]` Array of strings.
- - - -


{{< header href="substitute" >}}

substitute

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substitute.go#L17)

substitute replaces search values in the string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to modify
`search` | `String`  |  | The string representing a search pattern
`replace` | `String`  |  | The string representing a replace value
`limit` | `Int`  |  | The cap the number of replacements to this value.


**Returns** `String` Returns a string with replace substring.
- - - -


{{< header href="substring" >}}

substring

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L15)

substring returns a substring of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`offset` | `Int`  |  | Start at offset, offsets start at position 0.
`length` | `Int`  |  | At most length characters, omit to get the substring from offset to the end of the string.


**Returns** `String` A substring of value.
- - - -


{{< header href="to_base64" >}}

to_base64

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L80)

to_base64 returns the base64 representation of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to encode.


**Returns** `String` A base64 representation of the string.
- - - -


{{< header href="trim" >}}

trim

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L14)

trim returns the string value with whitespace stripped from the start and/or end.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars on both sides.
- - - -


{{< header href="unescape_html" >}}

unescape_html

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/unescape.go#L17)

unescape_html unescapes entities like "&lt;" to become "<". It unescapes a larger range of entities than EscapeString escapes. For example, "&aacute;" unescapes to "á", as does "&#225;" and "&#xE1;". UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to escape.


**Returns** `String` Escaped string.
- - - -


{{< header href="upper" >}}

upper

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L28)

upper converts strings to their upper-case counterparts. All other characters are returned unchanged.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.


**Returns** `String` This string in upper case.
- - - -

## Next steps

{{< docs-related tiles="stdlib,language-operators-comparison,language-types-basic" >}}
