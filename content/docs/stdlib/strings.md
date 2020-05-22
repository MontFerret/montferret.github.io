---
title: "strings"
weight: 1
draft: false
---


## LIKE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/like.go#L15)

Like checks whether the pattern search is contained in the string text, using wildcard matching.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to search in.
`search` | `String` | A search pattern that can contain the wildcard characters.
`caseInsensitive` | `Boolean` | If set to true, the matching will be case-insensitive. the default is false.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -

## TRIM
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L14)

Trim returns the string value with whitespace stripped from the start and/or end.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `String` | The string.
`chars` | `String` | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars on both sides.
- - - -

## LTRIM
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L34)

LTrim returns the string value with whitespace stripped from the start only.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `String` | The string.
`chars` | `String` | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the left-hand side.
- - - -

## RTRIM
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/trim.go#L55)

RTrim returns the string value with whitespace stripped from the end only.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `String` | The string.
`chars` | `String` | Overrides the characters that should be removed from the string. it defaults to \r\n \t.


**Returns** `String` The string without chars at the right-hand side.
- - - -

## FROM_BASE64
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L16)

FromBase64 returns the value of a base64 representation.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`base64String` | `String` | The string to decode.


**Returns** `String` The decoded string.
- - - -

## DECODE_URI_COMPONENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/decode.go#L36)

DecodeURIComponent returns the decoded String of uri.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`uri` | `String` | Uri to decode.


**Returns** `String` Decoded string.
- - - -

## CONCAT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L13)

Concat concatenates one or more instances of Read, or an Array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String...` `Array` | The source string / array.


**Returns** `String` Returns the concatenated string.
- - - -

## CONCAT_SEPARATOR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/concat.go#L47)

ConcatWithSeparator concatenates one or more instances of Read, or an Array with a given separator.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`separator` | `string` | The separator string.
`src` | `string...` `array` | The source string / array.


**Returns** `String` Returns the concatenated string.
- - - -

## LOWER
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L13)

Lower converts strings to their lower-case counterparts. All other characters are returned unchanged.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.


**Returns** `String` This string in lower case.
- - - -

## UPPER
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/case.go#L28)

Upper converts strings to their upper-case counterparts. All other characters are returned unchanged.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.


**Returns** `String` This string in upper case.
- - - -

## JSON_PARSE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L13)

JSONParse returns a FQL value described by the JSON-encoded input string.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to parse as json.


**Returns** `Value` Returns fql value
- - - -

## JSON_STRINGIFY
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/json.go#L34)

JSONStringify returns a JSON string representation of the input value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Value` | The input value to serialize.


**Returns** `String` Returns json string
- - - -

## FMT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/fmt.go#L18)

Fmt formats the template using these arguments.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`template` | `String` | Template.
`args` | `Any Values` | Template arguments.


**Returns** `String` String formed by template using arguments.
- - - -

## RANDOM_TOKEN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/random.go#L26)

RandomToken generates a pseudo-random token string with the specified length. The algorithm for token generation should be treated as opaque.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`length` | `Int` | The desired string length for the token. it must be greater than 0 and at most 65536.


**Returns** `String` A generated token consisting of lowercase letters, uppercase letters and numbers.
- - - -

## CONTAINS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/contains.go#L16)

Contains returns a value indicating whether a specified substring occurs within a string.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.
`search` | `String` | The string to seek.
`returnIndex` | `Boolean` | Values which indicates whether to return the character position of the match is returned instead of a boolean. the default is false.


**Returns** `Boolean` `Int` Returns index or a boolean value depending on returnindex.
- - - -

## ENCODE_URI_COMPONENT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L17)

EncodeURIComponent returns the encoded String of uri.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`uri` | `String` | Uri to encode.


**Returns** `String` Encoded string.
- - - -

## MD5
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L32)

Md5 calculates the MD5 checksum for text and return it in a hexadecimal string representation.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to do calculations against to.


**Returns** `String` Md5 checksum as hex string.
- - - -

## SHA1
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L48)

Sha1 calculates the SHA1 checksum for text and returns it in a hexadecimal string representation.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to do calculations against to.


**Returns** `String` Sha1 checksum as hex string.
- - - -

## SHA512
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L64)

Sha512 calculates the SHA512 checksum for text and returns it in a hexadecimal string representation.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to do calculations against to.


**Returns** `String` Sha512 checksum as hex string.
- - - -

## TO_BASE64
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/encode.go#L80)

ToBase64 returns the base64 representation of value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `string` | The string to encode.


**Returns** `String` A base64 representation of the string.
- - - -

## FIND_FIRST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L18)

FindFirst returns the position of the first occurrence of the string search inside the string text. Positions start at 0.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.
`search` | `String` | The string to seek.
`start` | `Int, optional` | Limit the search to a subset of the text, beginning at start.
`end` | `Int, optional` | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -

## FIND_LAST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/find.go#L65)

FindLast returns the position of the last occurrence of the string search inside the string text. Positions start at 0.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.
`search` | `String` | The string to seek.
`start` | `Int, optional` | Limit the search to a subset of the text, beginning at start.
`end` | `Int, optional` | Limit the search to a subset of the text, ending at end


**Returns** `Int` The character position of the match. if search is not contained in text, -1 is returned. if search is empty, start is returned.
- - - -

## SUBSTRING
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L15)

Substring returns a substring of value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `String` | The source string.
`offset` | `Int` | Start at offset, offsets start at position 0.
`length` | `Int, optional` | At most length characters, omit to get the substring from offset to the end of the string. optional.


**Returns** `String` A substring of value.
- - - -

## LEFT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L61)

Left returns the leftmost characters of the string value by index.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.
`length` | `Int` | The amount of characters to return.


**Returns** `String` Returns the left substring.
- - - -

## RIGHT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substr.go#L88)

Right returns the rightmost characters of the string value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`src` | `String` | The source string.
`length` | `Int` | The amount of characters to return.


**Returns** `String` Returns the right substring.
- - - -

## UNESCAPE_HTML
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/unescape.go#L17)

UnescapeHTML unescapes entities like "&lt;" to become "<". It unescapes a larger range of entities than EscapeString escapes. For example, "&aacute;" unescapes to "á", as does "&#225;" and "&#xE1;". UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`string` | `String` | Uri to unescape.


**Returns** `String` Escaped string.
- - - -

## SPLIT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/split.go#L16)

Split splits the given string value into a list of strings, using the separator.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to split.
`separator` | `String` | The separator.
`limit` | `Int` | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `Array<String>` Array of strings.
- - - -

## REGEX_MATCH
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L16)

RegexMatch returns the matches in the given string text, using the regex.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to search in.
`regex` | `String` | A regular expression to use for matching the text.
`caseInsensitive` | `Boolean` | If set to true, the matching will be case-insensitive. the default is false.


**Returns** `Array` An array of strings containing the matches.
- - - -

## REGEX_SPLIT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L58)

RegexSplit splits the given string text into a list of strings, using the separator.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to split.
`regex` | `String` | A regular expression to use for splitting the text.
`caseInsensitive` | `Boolean` | If set to true, the matching will be case-insensitive. the default is false.
`limit` | `Int` | Limit the number of split values in the result. if no limit is given, the number of splits returned is not bounded.


**Returns** `Array` An array of strings splited by the expression.
- - - -

## REGEX_TEST
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L100)

RegexTest test whether the regexp has at least one match in the given text.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to split.
`regex` | `String` | A regular expression to use for splitting the text.
`caseInsensitive` | `Boolean` | If set to true, the matching will be case-insensitive. the default is false.


**Returns** `Boolean` Returns true if the pattern is contained in text, and false otherwise.
- - - -

## REGEX_REPLACE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/regex.go#L133)

RegexReplace replace every substring matched with the regexp with a given string.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to split.
`regex` | `String` | A regular expression search pattern.
`replacement` | `String` | The string to replace the search pattern with
`caseInsensitive` | `Boolean` | If set to true, the matching will be case-insensitive. the default is false.


**Returns** `String` Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
- - - -

## ESCAPE_HTML
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/escape.go#L16)

EscapeHTML escapes special characters like "<" to become "&lt;". It escapes only five such characters: <, >, &, ' and ". UnescapeString(EscapeString(s)) == s always holds, but the converse isn't always true.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`string` | `String` | Uri to escape.


**Returns** `String` Escaped string.
- - - -

## SUBSTITUTE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/strings/substitute.go#L17)

Substitute replaces search values in the string value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`text` | `String` | The string to modify
`search` | `String` | The string representing a search pattern
`replace` | `String` | The string representing a replace value
`limit` | `Int` | The cap the number of replacements to this value.


**Returns** `String` Returns a string with replace substring.
- - - -
