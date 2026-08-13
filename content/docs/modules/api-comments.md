---
title: "Document module APIs"
sidebarTitle: "API comments"
weight: 25
draft: false
description: "Describe Ferret-facing module functions with Go documentation comments and structured API tags."
---

# Document module APIs

Ferret modules implemented in Go can expose APIs whose FQL signatures are not clear from their Go declarations. A variadic Go function such as:

```go
func Decode(ctx context.Context, args ...runtime.Value) (runtime.Value, error)
```

may expose this Ferret API:

```text
xml::decode(data: String | Binary) -> Object
```

Use ordinary Go documentation comments for human-readable documentation and structured tags for the Ferret-facing parameters, result, failures, and deprecation state. Barn derives the public function name and namespace from module registration rather than from the Go identifier.

## Write the documentation comment

Start with ordinary Go documentation prose. The first sentence should begin with the Go declaration name where practical. Put structured tags after all prose and examples:

```go
// Decode decodes XML content into a normalized document object.
//
// The decoder accepts UTF-8 XML input and returns a normalized Ferret object.
//
// Example:
//
//	return xml::decode("<root><item /></root>")
//
// @param data {String|Binary} XML content.
// @return {Object} Normalized XML document.
func Decode(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return decode(ctx, args)
}
```

The prose remains useful to Go tools and source readers. The structured block supplies only the Ferret API details that Go types cannot describe reliably.

## Document parameters

Use this form for each Ferret parameter:

```text
@param <name> {<type>} <description>
```

List parameters in Ferret call order:

```go
// @param path {String} Path to the archive.
// @param options {Object?} Extraction options.
```

Optionality belongs to the type. Use `?` instead of brackets around the parameter name:

```go
// @param options {Object?} Decode options.
```

Use `|` for unions, `Array<T>` for collections, and a trailing `...` when documenting a variadic Ferret value:

```go
// @param data {String|Binary} XML content.
// @param names {Array<String>} Field names.
// @param values {Any...} Values to concatenate.
```

Barn preserves the complete type expression as authored; it does not parse or
canonicalize this notation. Whether the registered function signature is
variadic still comes from the Ferret SDK registration and Go function type.

Do not use JSDoc-style parameter names or separators:

```text
@param {String} path - Path to the file.
@param [options] {Object}
@param path {String} - Path to the file.
```

## Document the result

Use singular `@return` for a meaningful result:

```text
@return {<type>} <description>
```

For example:

```go
// @return {Object} Normalized XML document.
```

A Ferret function may have at most one `@return` tag. Omit the tag when there is no useful Ferret-facing result metadata to publish.

## Document visible failures

Use `@throws` for meaningful Ferret-visible failures:

```text
@throws {<error>} <description>
```

A function may declare more than one failure:

```go
// @throws {ParseError} XML input is malformed.
// @throws {LimitError} XML input exceeds the configured limit.
```

Document failures that callers can understand or handle. Do not list every internal Go error.

## Mark deprecated functions

Keep the standard Go `Deprecated:` paragraph so Go tooling recognizes the declaration. Add `@deprecated` for generated Ferret metadata:

```text
@deprecated <description>
```

```go
// OldDecode decodes legacy XML input.
//
// Deprecated: use Decode instead.
//
// @param data {String} XML content.
// @return {Object} Normalized XML document.
// @deprecated Use Decode instead.
func OldDecode(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return decode(ctx, args)
}
```

## Keep tags ordered and focused

Use this order:

1. `@param` in Ferret call order
2. `@return`
3. `@throws`
4. `@deprecated`

Descriptions should be concise, begin with a capital letter, end with punctuation, and describe Ferret semantics instead of Go implementation details.

Prefer:

```text
@param path {String} Path to the archive.
```

Avoid:

```text
@param path {String} string passed into os.Open
```

The initial contract supports only `@param`, `@return`, `@throws`, and `@deprecated`. Keep examples, remarks, links, and longer explanations in ordinary Go documentation rather than adding new tags.

## How Barn uses the comment

Barn associates supported tags with the registered public API declaration and keeps the ordinary prose separate from structured metadata. It parses the structured block deterministically and reports malformed supported tags with the source file, declaration, and offending tag instead of guessing.

The registered function name remains authoritative. A Go function named `Decode` can therefore document `xml::decode`, and the same declaration can be registered under another namespace without changing its Go name.

Raw tags are not the final Registry presentation. The complete `Decode` comment allows the Registry to present information equivalent to:

```text
xml::decode(data: String | Binary) -> Object

Decodes XML content into a normalized document object.

Parameters
data    String | Binary    XML content.

Returns
Object    Normalized XML document.
```

Unknown ordinary `@...` text remains prose. Do not use supported tag names in prose when they could be mistaken for API metadata.

The format is intentionally small. It describes only information that cannot be recovered reliably from Go types and registration. It is not a replacement for Go documentation and does not recreate JSDoc or JavaDoc inside Go comments.

For a complete implementation walkthrough, see [Develop a Ferret module]({{< ref "/docs/guides/writing-plugins" >}}). Before tagging a release, review [Publish a module]({{< ref "/docs/modules/publish" >}}).
