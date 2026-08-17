---
title: "Compiler Analysis"
sidebarTitle: "Compiler Analysis"
weight: 100
draft: false
description: "Inspect compiler-resolved symbols, references, calls, diagnostics, and type facts without executing a query."
aliases:
    - /docs/embedding/compiler-analysis/
---

# Compiler Analysis

Use `compiler.Analyze` when an embedding application needs semantic information about FQL source without executing it. The result is a protocol-independent snapshot suitable for editor integrations, source navigation, validation UIs, and other developer tools.

Analysis uses the same parser, name resolution, scopes, captures, call resolution, diagnostics, and type facts as normal compilation. It does not maintain a separate language model.

## Analyze a source

Create a compiler and pass it a `source.Source`:

{{< code lang="go" >}}
package main

import (
    "fmt"

    "github.com/MontFerret/ferret/v2/pkg/compiler"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    src := source.New("query.fql", `
LET products = @products
FUNC available(product) => product.stock > 0
RETURN FOR product IN products
    FILTER available(product)
    RETURN product.name
`)

    analysis, err := compiler.New().Analyze(src)
    if analysis == nil {
        panic("analysis did not produce a snapshot")
    }

    for _, symbol := range analysis.Symbols() {
        fmt.Printf("%s: %v\n", symbol.Name, symbol.Kind)
    }

    for _, diagnostic := range analysis.Diagnostics() {
        fmt.Println(diagnostic.Message)
    }

    if err != nil {
        // The snapshot remains useful when source diagnostics are present.
        fmt.Println("source needs attention:", err)
    }
}
{{</ code >}}

`Analyze` is safe to call concurrently on the same `Compiler`. It always uses unoptimized, full-source analysis, regardless of the compiler's optimization or debug options.

## Work with partial results

Source diagnostics do not discard semantic work that already succeeded:

{{< code lang="go" >}}
analysis, err := c.Analyze(src)
renderDiagnostics(analysis.Diagnostics())
updateOutline(analysis.Symbols())

if err != nil {
    // The source has syntax or semantic diagnostics.
}
{{</ code >}}

For empty or syntactically invalid input, the snapshot contains cloned diagnostics but does not guarantee semantic entities. After parsing succeeds, symbols, references, calls, and type facts established before or alongside semantic errors remain available.

The returned `error` follows the compiler's diagnostic error behavior. Use `Analysis.Diagnostics()` when you need the structured diagnostics and their spans.

## Syntax tokens

`Analysis.SyntaxTokens()` returns the non-whitespace lexer tokens captured by the same frontend pass that produced the semantic snapshot. The result includes hidden-channel line and block comments, omits whitespace and newline trivia, and remains available in partial analyses with syntax errors.

Each `SyntaxToken` contains a parser-independent `SyntaxTokenKind`, an optional canonical `SyntaxWord` identity, and a source span. Kinds distinguish identifiers, namespace segments, keywords, strings, numbers, durations, comments, operators, punctuation, and otherwise unknown input. Word identities distinguish case-insensitive FQL words without exposing generated lexer token numbers; non-word tokens use `SyntaxWordUnknown`. Template string text and delimiters are string tokens; embedded-expression delimiters remain punctuation so consumers can combine syntax and semantic classifications without inspecting generated parser types.

{{< code lang="go" >}}
for _, token := range analysis.SyntaxTokens() {
    fmt.Println(token.Kind, token.Word, token.Span.Start, token.Span.End)
}
{{</ code >}}

`compiler.SyntaxWords()` returns a deterministic defensive copy of the canonical word metadata. Each entry includes the typed identity, uppercase canonical spelling, and a category that distinguishes keywords, word operators, literals, and contextual parser words. This lets completion and other language tooling obtain spelling and classification from Ferret while choosing which categories make sense in a particular editor context.

Token spans use the same zero-based, half-open UTF-8 byte offsets as the other analysis APIs. No ANTLR or generated-parser type is exposed, and callers should depend only on the stable compiler token kinds and word identities rather than the grammar's internal token numbers.

## Symbols and references

Every source-visible declaration receives a `SymbolID`. Distinct shadowed declarations receive distinct IDs even when their names are the same.

`Symbol` reports:

- the source name and `SymbolKind`;
- whether the binding is mutable;
- the compiler-established `ValueType`;
- the full declaration span and exact name-selection span; and
- whether a source declaration exists.

Kinds distinguish ordinary bindings, function parameters, user-defined functions, bind parameters, loop bindings, `MATCH` bindings, `COLLECT` bindings, and namespace aliases.

A distinct bind parameter such as `@products` is represented by one analysis-local symbol. Its `HasDeclaration` field is false because the parameter has no FQL declaration. Each occurrence is still returned as a resolved reference to that symbol.

Use the accessors according to the task:

{{< code lang="go" >}}
symbols := analysis.Symbols()

symbol, ok := analysis.Symbol(symbolID)
references := analysis.References()
referencesToSymbol := analysis.ReferencesTo(symbolID)

symbolAtCursor, ok := analysis.SymbolAt(byteOffset)
visible := analysis.VisibleSymbols(byteOffset)
{{</ code >}}

`VisibleSymbols` applies lexical shadowing and the language's activation rules. User-defined functions are visible throughout their owning scope. Variables appear after a successful declaration, function parameters inside their bodies, loop bindings after the input expression, `MATCH` bindings in their arm, and `COLLECT` bindings after their clause. Namespace aliases are source-wide, and bind parameters become source-wide once known to the analysis.

Compiler-only pseudo-bindings and forwarding bindings used to implement captures are not exposed. A captured reference resolves to the original lexical declaration's `SymbolID`, including transitive captures through nested functions.

## Calls

`Analysis.Calls()` returns resolved user-defined, compiler-builtin, and host calls. Each `Call` includes:

- the full call, callee, and individual argument spans;
- a resolved display name;
- a normalized runtime identity for builtin and host calls;
- the result type when the compiler knows it; and
- a target `SymbolID` for a user-defined function.

Source calls to user-defined functions also appear in `References()`. Host and registered calls do not receive fabricated symbols, so their target ID is zero.

{{< code lang="go" >}}
call, ok := analysis.CallAt(byteOffset)
if ok {
    fmt.Println(call.Name, call.Identity, call.ArgumentSpans)
}
{{</ code >}}

Analysis classifies a host call from compiler-visible syntax and `USE` aliases. It does not validate the call against an engine or module function registry. A host call may therefore be absent at runtime even though it appears in the semantic snapshot.

## Type facts

`TypeFacts()` reports only facts established by the current compiler. It does not run a separate general-purpose type inference pass.

The public `ValueType` set covers unknown, dynamic (`ValueTypeAny`), none, integer, float, duration, boolean, string, array, object, list, and map. Host and user-defined function results remain dynamic when the compiler cannot prove a narrower result. `ValueTypeUnknown` is a meaningful result: it means no existing compiler fact established the type.

Use `TypeAt` to find the narrowest expression fact at a byte offset:

{{< code lang="go" >}}
fact, ok := analysis.TypeAt(byteOffset)
if ok {
    fmt.Println(fact.Type)
}
{{</ code >}}

## Offsets and ID lifetime

All positions and spans are zero-based UTF-8 byte offsets with half-open ranges: the start byte is included and the end byte is excluded. Convert editor line and character positions before calling the offset-based queries.

`SymbolID` zero is invalid. Non-zero IDs are deterministic only within one `Analysis` result. Do not persist them as cross-analysis, cross-version, or cross-document identities. Re-analyze the source and rebuild any external index from the new snapshot.

## Snapshot ownership

An `Analysis` is immutable after `Analyze` returns. Its accessors return defensive copies, including syntax tokens, nested argument spans, and diagnostic spans, so callers can sort or annotate returned values without modifying the stored snapshot.

The API has no `context.Context` parameter because compiler analysis does not currently expose meaningful cancellation. Applications that need scheduling control should manage analysis calls at their own worker or request boundary.

## What this API does not provide

Compiler analysis is a single-source semantic snapshot. It does not provide:

- an LSP server or protocol types;
- incremental parsing or incremental analysis;
- project-wide indexing or persistent symbol IDs;
- completion, hover, rename, or navigation handlers; or
- validation against a host function registry.

Those features can consume this snapshot, but their protocol, caching, scheduling, and project ownership remain in the embedding tool.

## Next steps

{{< docs-related tiles="embedding-go-executing,embedding-go-configuration,embedding-go-programs,embedding-go-custom-functions" >}}
