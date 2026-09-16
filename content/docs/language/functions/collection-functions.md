---
title: "Collection functions"
sidebarTitle: "Collections"
weight: 46
draft: false
description: "Count iterable values, test membership, and reverse lists or strings."
---

# Collection functions

The global `count`, `count_distinct`, `includes`, and `reverse` functions
use the capabilities supplied by each value. The iterable counting and custom
list construction described here require a runtime containing the collection
factory update.

```fql
RETURN {
    ascending: count(1..3),
    descending: count(3..1),
    distinct: count_distinct({first: 7, second: 7.0}),
    present: includes("123", 123),
    reversed: reverse([1, 2, 3])
}
// {ascending: 3, descending: 3, distinct: 1, present: true, reversed: [3, 2, 1]}
```

## Count yielded values

`count(source)` accepts an iterable, including arrays, objects, native ranges,
and read-only host sources. When the source supplies a length, count uses it
without iteration. A length failure or negative host length produces an error
without a scan fallback. Zero is a valid length. Cancellation during measurement
remains discoverable alongside an invalid-length error.
Host length computation can perform I/O and need not take constant time.

Otherwise, count traverses once. The traversal may consume a one-shot source,
perform I/O, or never finish if the source is unbounded. Counts exceeding the
runtime integer range fail instead of wrapping. Strings and non-iterable scalars
are invalid counting inputs.

`count_distinct(source)` always traverses and counts distinct yielded values.
Objects contribute values, not keys or key/value pairs. Equality follows Ferret
value semantics: equivalent integers and floats count together, hash collisions
do not merge unequal values, and Duration equality remains strict. Host equality
errors propagate.

## Test membership

`includes(container, needle)` handles strings first. It converts the needle to
text, so `includes("123", 123)` is true. The separate string `contains`
function retains its strict string arguments.

For other values, includes prefers the host's membership capability. Otherwise,
it scans iterable values using canonical equality and stops at the first match.
Objects are searched by value: `includes({key: 7}, "key")` is false, while
`includes({key: 7}, 7)` is true.

## Reverse a list or string

`reverse(list)` creates a new outer list in the source's implementation family,
preserving relevant backend configuration. It reads elements by descending
index and keeps their references shallow. It does not change the source, clone
elements, or materialize arbitrary non-list streams.

`reverse(string)` reverses Unicode code points. It does not preserve grapheme
clusters as units.

## Errors and resource ownership

Owned scans close the iterator they create exactly once, including after an
early match, cancellation, or failure. They do not explicitly close the source
or yielded values. A close failure is reported even after an otherwise successful
scan or membership match; combined errors retain their causes.

Collection functions pass the caller's context to host operations without
polling it during synchronous work. Hosts own cancellation of blocking operations.
A completed operation keeps its result; the VM observes execution cancellation
at its next cancellation boundary. Direct Go calls may complete with a canceled
context. A host scan that ignores context and never returns cannot be interrupted
by the VM.

If list reversal fails after creating a destination, it closes that incomplete
destination when closable and preserves cleanup errors with the original
failure. On success, destination ownership passes to normal result lifecycle
handling. The borrowed source and its elements are never explicitly closed.

Host authors must migrate `Spawnable.Empty` to
[the collection factory contract]({{< ref "docs/embedding/go/host-values" >}}#constructing-an-empty-collection).
Existing FQL function names and arities remain unchanged. This update does not
add `arrays::reverse` or change `is_empty`.
