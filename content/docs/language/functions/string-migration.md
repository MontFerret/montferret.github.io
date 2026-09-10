---
title: "String API Migration"
sidebarTitle: "String API Migration"
weight: 70
draft: false
description: "Migrate late-alpha string, path, encoding, and crypto calls to the streamlined Ferret v2 API."
---

The Ferret v2 string cleanup introduces strict text arguments, explicit regex
results, and separate namespaces for path, encoding, and crypto. These changes target
the upcoming v2 runtime build; use a runtime containing the cleanup before
running the examples below. Removed names have no compatibility aliases.

{{< code lang="fql" >}}
let tags = ["docs", "fql", "runtime"]
return {
    joined: join(tags, ", "),
    digest: crypto::sha256("Ferret"),
    encoded: encoding::base64_encode("é")
}
{{</ code >}}

## Rename calls

| Former call | Current call |
| --- | --- |
| `concat_separator(separator, values)` | `join(values, separator)` |
| `substitute(text, search, replacement[, limit])` | `replace(text, search, replacement[, limit])` |
| `substitute(text, search)` | `replace(text, search, "")` |
| `contains(text, search, true)` | `find_first(text, search)` for an index |
| `regex_match(text, pattern)` | `regex_find(text, pattern)` or `regex_find_all(text, pattern)` |
| Regex Boolean case options | Inline pattern flags such as `(?i)` |
| `json_parse`, `json_stringify` | `encoding::json_parse`, `encoding::json_stringify` |
| `encode_uri_component`, `decode_uri_component` | `encoding::query_escape`, `encoding::query_unescape` |
| `to_base64`, `from_base64` | `encoding::base64_encode`, `encoding::base64_decode` |
| `escape_html`, `unescape_html` | `encoding::html_escape`, `encoding::html_unescape` |
| `md5`, `sha1`, `sha512`, `random_token` | The same names under `crypto::` |
| `base`, `clean`, `dir`, `ext`, `is_abs`, `separate`, `match` | The same names under `path::` |
| Global path `join(parts...)` | `path::join(parts...)` |

Actual string operations remain global. `starts_with`, `ends_with`, `repeat`,
and `crypto::sha256` are new. Identifiers can contain underscores after digits,
including `base64_encode`.

All eight path functions move into `path::` without global aliases. Existing path
behavior and argument counts are preserved. Global `join(values, separator)`
continues to join strings; use `path::join` to assemble paths.

{{< code lang="fql" >}}
return path::base(path::join("reports", "result.json"))
// "result.json"
{{</ code >}}

## Validate text and bounds

Text, patterns, separators, cutsets, and `join` elements require String values.
Convert other values explicitly with `to_string` when that is intended. Supplied
bounds, lengths, and limits require Int values. Omitted options use defaults;
wrong types return errors.

`left`, `right`, `substring`, `find_first`, and `find_last` count Unicode runes,
not bytes or grapheme clusters. Left/right lengths must be non-negative and
clamp to the string length. Negative or out-of-range substring offsets produce
an empty string; its optional length is non-negative and clamps safely.
Find bounds clamp to the string length, with an exclusive end. Reversed bounds
produce `-1`; an empty search returns the effective start for `find_first` or
the effective end for `find_last`.

All trim variants default to Unicode whitespace, including tabs and newlines.
Explicit cutsets are literal rune sets; an empty cutset leaves text unchanged.

`split` and `regex_split` limits are non-negative maximum result counts. Zero
returns `[]`; positive limits retain the unsplit remainder. `replace` limits
count replacements, with zero leaving text unchanged. Omit a limit for unlimited
operation; explicit `-1` is invalid. Literal replacement does not expand dollar
expressions. An empty search matches at rune boundaries. `repeat` rejects
negative counts and counts outside the host integer range. Results must not
exceed 64 MiB (67,108,864 UTF-8 bytes), including when the count is one. Oversized
results are rejected before allocation with an argument error that `ON ERROR`
can catch. Zero count or empty text returns an empty string for any valid count.

`join` preserves order and empty elements, inserts exactly one separator between
neighbors, and returns an empty string for an empty list. Nested lists and None
elements are errors. `concat` retains its conversion of Any values to text;
a single List argument contributes its elements, and None contributes no text.

## Format values

`fmt` accepts either automatic `{}` placeholders or zero-based decimal `{n}`
indexes. Do not mix the two modes. Repeated indexes are allowed, but every
supplied value must be referenced. `{{` and `}}` produce literal braces.
Malformed placeholders, missing values, overflowing indexes, and unused values
return errors. Substitution values use their runtime string representations.

{{< code lang="fql" >}}
return fmt("{{{1}}} {0} {1}", "é", 42) // "{42} é 42"
{{</ code >}}

## Match patterns

`contains`, `starts_with`, and `ends_with` always return Boolean. `like` matches
the whole string with glob syntax: `*`, `?`, character classes, and alternatives.
Percent and underscore are literal characters. Its optional case-insensitive
argument must be Boolean.

Regex functions use Go regular expressions and inline flags such as `(?i)`.
`regex_test` returns Boolean. `regex_find` returns None or an object containing
`match`, `groups`, and `named`; `regex_find_all` returns an array of these objects.

{{< code lang="fql" >}}
return regex_find("é😀", "(?P<letter>é)(😀)?")
// {match: "é😀", groups: ["é", "😀"], named: {letter: "é"}}
{{</ code >}}

`groups` excludes the full match and preserves capture order. Unmatched groups
are empty strings, including in `named`. `named` maps capture names to strings.
No captures yields an empty array and object.

`regex_find` and `regex_find_all` reject duplicate non-empty capture names. For
example, `(?P<value>a)(?P<value>b)` returns an invalid-argument error on the
pattern (argument 2), identifying `value` as the duplicate. Rename the captures
or make a group unnamed. Names are compared exactly: `value` and `Value` are
distinct. Duplicates are rejected before matching, including across alternatives,
in optional groups, and when the text would produce no match. These errors can
be caught with `ON ERROR`.

`regex_test`, `regex_replace`, and `regex_split` retain Go's duplicate-name
behavior because they do not return named capture objects.
All-match operations use Go's non-overlapping and zero-width match rules.
`regex_replace` expands `$1`, `${name}`, and `$$` using Go replacement semantics.
Invalid expressions return argument errors. Boolean regex options and the old
reserved `regex_split` argument are removed.

## Encode data and generate tokens

Query escaping uses query-form rules: spaces become `+`, and a literal plus
becomes `%2B`. Base64 encoding accepts String UTF-8 bytes or Binary and emits
standard padded Base64. Decoding returns Binary, including arbitrary non-text
bytes. HTML helpers escape and decode HTML character references. JSON parsing
requires String; JSON stringify accepts Any and retains existing codec behavior.

Hash functions accept String or Binary bytes and return lowercase hexadecimal
digests. MD5 and SHA-1 remain available for legacy digest interoperability;
they are unsuitable for security-sensitive collision resistance.
`crypto::random_token` accepts lengths from 1 through 65536 inclusive, uses
cryptographically secure randomness, and samples lowercase letters, uppercase
letters, and digits uniformly. Invalid lengths and entropy-source failures
return errors.

Go embedders can select `stdlib.Strings`, `stdlib.Encoding`, and `stdlib.Crypto`
independently. Both `stdlib.Full()` and `stdlib.Safe()` include all three.
