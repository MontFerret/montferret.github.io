---
title: "strings"
weight: 1
draft: false
menuTitle: 
menu: [CONCAT,CONCAT_SEPARATOR,CONTAINS,DECODE_URI_COMPONENT,ENCODE_URI_COMPONENT,ESCAPE_HTML,FIND_FIRST,FIND_LAST,FMT,FROM_BASE64,JSON_PARSE,JSON_STRINGIFY,LEFT,LIKE,LOWER,LTRIM,MD5,RANDOM_TOKEN,REGEX_MATCH,REGEX_REPLACE,REGEX_SPLIT,REGEX_TEST,RIGHT,RTRIM,SHA1,SHA512,SPLIT,SUBSTITUTE,SUBSTRING,TO_BASE64,TRIM,UNESCAPE_HTML,UPPER,]
---



{{< header >}}

CONCAT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L13)

CONCAT concatenates one or more instances of String, or an Array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`src` | `String, repeated` `String[]`  |  | The source string / array.


**Returns** `String` A string value.
- - - -


{{< header >}}

CONCAT_SEPARATOR

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L47)

CONCAT_SEPARATOR concatenates one or more instances of String, or an Array with a given separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`separator` | `String`  |  | The separator string.
`src` | `String, repeated` `String[]`  |  | The source string / array.


**Returns** `String` Concatenated string.
- - - -


{{< header >}}

CONTAINS

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/contains.go#L15)

CONTAINS returns a value indicating whether a specified substring occurs within a string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`returnIndex` | `Boolean`  | `False` | Values which indicates whether to return the character position of the match is returned instead of a boolean.


**Returns** `Boolean` `Int` A value indicating whether a specified substring occurs within a string.
- - - -


{{< header >}}

DECODE_URI_COMPONENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L36)

DECODE_URI_COMPONENT returns the decoded String of uri.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to decode.


**Returns** `String` Decoded string.
- - - -


{{< header >}}

ENCODE_URI_COMPONENT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L17)

ENCODE_URI_COMPONENT returns the encoded String of uri.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to encode.


**Returns** `String` Encoded string.
- - - -


{{< header >}}

ESCAPE_HTML

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/escape.go#L16)

ESCAPE_HTML escapes special characters like "<" to become "&lt;". It escapes only five such characters: <, >, &, ' and ". UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to escape.


**Returns** `String` Escaped string.
- - - -


{{< header >}}

FIND_FIRST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L17)

FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`start` | `Int`  |  | Limit the search to a subset of the text, beginning at start.
`end` | `Int`  |  | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -


{{< header >}}

FIND_LAST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L63)

FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`src` | `String`  |  | The source string.
`search` | `String`  |  | The string to seek.
`start` | `Int`  |  | Limit the search to a subset of the text, beginning at start.
`end` | `Int`  |  | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -


{{< header >}}

FMT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/fmt.go#L18)

FMT formats the template using these arguments.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`template` | `String`  |  | Template.
`args` | `Any, repeated`  |  | Template arguments.


**Returns** `String` String formed by template using arguments.
- - - -


{{< header >}}

FROM_BASE64

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L16)

FROM_BASE64 returns the value of a base64 representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to decode.


**Returns** `String` The decoded string.
- - - -


{{< header >}}

JSON_PARSE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L13)

JSON_PARSE returns a value described by the JSON-encoded input string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to parse as json.


**Returns** `Any` Parsed value.
- - - -


{{< header >}}

JSON_STRINGIFY

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L34)

JSON_STRINGIFY returns a JSON string representation of the input value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `Any`  |  | The input value to serialize.


**Returns** `String` Json string.
- - - -


{{< header >}}

LEFT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L61)

LEFT returns the leftmost characters of the string value by index.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`length` | `Int`  |  | The amount of characters to return.


**Returns** `String` The leftmost characters of the string value by index.
- - - -


{{< header >}}

LIKE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/like.go#L15)

LIKE checks whether the pattern search is contained in the string text, using wildcard matching.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to search in.
`search` | `String`  |  | A search pattern that can contain the wildcard characters.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -


{{< header >}}

LOWER

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L13)

LOWER converts strings to their lower-case counterparts. All other characters are returned unchanged.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.


**Returns** `String` This string in lower case.
- - - -


{{< header >}}

LTRIM

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L34)

LTRIM returns the string value with whitespace stripped from the start only.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the left-hand side.
- - - -


{{< header >}}

MD5

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L32)

MD5 calculates the MD5 checksum for text and return it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Md5 checksum as hex string.
- - - -


{{< header >}}

RANDOM_TOKEN

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/random.go#L26)

RANDOM_TOKEN generates a pseudo-random token string with the specified length. The algorithm for token generation should be treated as opaque.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`len` | `Int`  |  | The desired string length for the token. it must be greater than 0 and at most 65536.


**Returns** `String` A generated token consisting of lowercase letters, uppercase letters and numbers.
- - - -


{{< header >}}

REGEX_MATCH

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L16)

REGEX_MATCH returns the matches in the given string text, using the regex.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to search in.
`expression` | `String`  |  | A regular expression to use for matching the text.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.


**Returns** `Any[]` An array of strings containing the matches.
- - - -


{{< header >}}

REGEX_REPLACE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L133)

REGEX_REPLACE replace every substring matched with the regexp with a given string.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`expression` | `String`  |  | A regular expression search pattern.
`replacement` | `String`  |  | The string to replace the search pattern with
`caseInsensitive` | `Boolean`  | `False` | Insensitive.


**Returns** `String` Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
- - - -


{{< header >}}

REGEX_SPLIT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L58)

REGEX_SPLIT splits the given string text into a list of strings, using the separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`expression` | `String`  |  | A regular expression to use for splitting the text.
`caseInsensitive - If set to true, the matching will be case` | `Boolean`  |  | Insensitive. the default is false.
`limit` | `Int`  |  | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `Any[]` An array of strings splitted by the expression.
- - - -


{{< header >}}

REGEX_TEST

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L100)

REGEX_TEST test whether the regexp has at least one match in the given text.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to test.
`expression` | `String`  |  | A regular expression to use for splitting the text.
`caseInsensitive` | `Boolean`  | `False` | Insensitive.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -


{{< header >}}

RIGHT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L88)

RIGHT returns the rightmost characters of the string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`length` | `Int`  |  | The amount of characters to return.


**Returns** `String` The rightmost characters of the string value.
- - - -


{{< header >}}

RTRIM

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L55)

RTRIM returns the string value with whitespace stripped from the end only.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the right-hand side.
- - - -


{{< header >}}

SHA1

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L48)

SHA1 calculates the SHA1 checksum for text and returns it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Sha1 checksum as hex string.
- - - -


{{< header >}}

SHA512

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L64)

SHA512 calculates the SHA512 checksum for text and returns it in a hexadecimal string representation.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to do calculations against to.


**Returns** `String` Sha512 checksum as hex string.
- - - -


{{< header >}}

SPLIT

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/split.go#L16)

SPLIT splits the given string value into a list of strings, using the separator.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to split.
`separator` | `String`  |  | The separator.
`limit` | `Int`  |  | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `String[]` Array of strings.
- - - -


{{< header >}}

SUBSTITUTE

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substitute.go#L17)

SUBSTITUTE replaces search values in the string value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to modify
`search` | `String`  |  | The string representing a search pattern
`replace` | `String`  |  | The string representing a replace value
`limit` | `Int`  |  | The cap the number of replacements to this value.


**Returns** `String` Returns a string with replace substring.
- - - -


{{< header >}}

SUBSTRING

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L15)

SUBSTRING returns a substring of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.
`offset` | `Int`  |  | Start at offset, offsets start at position 0.
`length` | `Int`  |  | At most length characters, omit to get the substring from offset to the end of the string.


**Returns** `String` A substring of value.
- - - -


{{< header >}}

TO_BASE64

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L80)

TO_BASE64 returns the base64 representation of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string to encode.


**Returns** `String` A base64 representation of the string.
- - - -


{{< header >}}

TRIM

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L14)

TRIM returns the string value with whitespace stripped from the start and/or end.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The string.
`chars` | `String`  |  | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars on both sides.
- - - -


{{< header >}}

UNESCAPE_HTML

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/unescape.go#L17)

UNESCAPE_HTML unescapes entities like "&lt;" to become "<". It unescapes a larger range of entities than EscapeString escapes. For example, "&aacute;" unescapes to "á", as does "&#225;" and "&#xE1;". UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`uri` | `String`  |  | Uri to escape.


**Returns** `String` Escaped string.
- - - -


{{< header >}}

UPPER

{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L28)

UPPER converts strings to their upper-case counterparts. All other characters are returned unchanged.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`str` | `String`  |  | The source string.


**Returns** `String` This string in upper case.
- - - -
