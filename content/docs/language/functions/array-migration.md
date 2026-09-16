---
title: "Array API Migration"
sidebarTitle: "Array API Migration"
weight: 75
draft: false
description: "Use immutable arrays, opt into explicit mutation, and preserve legacy global calls during migration."
---

The canonical v2 array library uses `arrays::`. Use a runtime containing this
cleanup before running these examples. Existing global array functions remain
available temporarily with their legacy signatures and behavior; they are
intended for removal in a later v2 minor release.

{{< code lang="fql" >}}
let values = [3, 1, 3]
return {
    appended: arrays::append(values, 2),
    concatenated: arrays::concat(values, [1, 2]),
    combined: arrays::union(values, [1, 2]),
    ordered: arrays::sorted(arrays::unique(values))
}
{{</ code >}}

The results are `[3,1,3,2]`, `[3,1,3,1,2]`, `[3,1,2]`, and `[1,3]`, respectively.
These functions return new arrays. Explicit mutation is available through
`arrays::mut::` in runtime releases containing the mutable array API.

## Choose the canonical operation

| Global call | Canonical operation |
| --- | --- |
| `first`, `last`, `flatten`, `slice`, `unique`, `sorted` | Same name under `arrays::` |
| `range(start, end[, step])` | `arrays::range(start, end[, step])` |
| `nth(xs, index)` | `arrays::at(xs, index)` |
| `append(xs, value)` or `push(xs, value)` | `arrays::append(xs, value)` |
| `union(a, b, ...)` | `arrays::concat(a, b, ...)` |
| `union_distinct(a, b, ...)` | `arrays::union(a, b, ...)` |
| `intersection(a, b, ...)` | `arrays::intersection(a, b, ...)` |
| `minus(a, b, ...)` | `arrays::difference(a, b, ...)` |
| `position(xs, value)` or `position(xs, value, false)` | `arrays::contains(xs, value)` |
| `position(xs, value, true)` | `arrays::index_of(xs, value)` |
| `remove_value(xs, value)` | `arrays::remove(xs, value)` |
| `remove_values(xs, values)` | `arrays::remove_any(xs, values)` |
| `remove_nth(xs, index)` | `arrays::remove_at(xs, index)` |
| `sorted_unique(xs)` | `arrays::sorted(arrays::unique(xs))` |
| `shift(xs)` | `arrays::slice(xs, 1)` |
| `unshift(xs, value)` | `arrays::concat([value], xs)` |

`arrays::append` adds one element: appending an array nests it. Use
`arrays::concat` to combine elements from multiple arrays. Concatenation and all
set operations require at least two array arguments.

`arrays::contains` always returns Boolean. `arrays::index_of` returns the first
matching index or `-1`. Indexes are zero-based. `arrays::at` returns `none` for
negative or missing indexes; `arrays::remove_at` returns an unchanged copy for
those indexes.

## Construct numeric ranges

In runtimes containing the math namespace migration, use
`arrays::range(start, end[, step])` to construct an array with inclusive
endpoints. The default step is positive one; pass a negative step to descend.

{{< code lang="fql" >}}
RETURN [arrays::range(1, 4), arrays::range(4, 1, -1)]
// [[1, 2, 3, 4], [4, 3, 2, 1]]
{{</ code >}}

Arguments must be finite numbers. A step pointing away from the endpoint
produces an empty array. Zero or non-advancing steps and ranges too large to
represent raise errors. The generated values retain the legacy floating-point
behavior.

The deprecated global `range` shares this implementation. It remains available
in Math-only embeddings; Arrays-only embeddings expose `arrays::range`.
There is no `math::range`.

## Compose operations without mode flags

Canonical `append` and `remove` take exactly two arguments. Canonical search
functions do not accept a Boolean return-type switch.

To deduplicate an appended result, use:

{{< code lang="fql" >}}
return arrays::unique(arrays::append([1, 1], 2))
{{</ code >}}

This returns `[1,2]`. It intentionally removes existing duplicates too. Legacy
`append(xs, value, true)` and `push(xs, value, true)` only skip adding a value
already present; they preserve existing duplicates. Keep those global calls
when that exact legacy behavior is required. A false flag is equivalent to the
canonical two-argument append.

`arrays::remove(xs, value)` removes all matches. `arrays::remove_any(xs, values)`
removes all matches to any listed value, preserving other duplicates and order.
The optional removal limit remains only on global `remove_value`: negative
means unlimited, zero removes nothing, and positive limits removals.

Legacy `unshift(xs, value, true)` corresponds to
`arrays::concat([value], arrays::remove(xs, value))`. It preserves duplicates of
other values. Global `pop` remains available and safely returns `[]` for an empty
array. To express it through canonical functions:

{{< code lang="fql" >}}
let xs = [1, 2, 3]
return length(xs) == 0 ? [] : arrays::slice(xs, 0, length(xs) - 1)
{{</ code >}}

## Understand set ordering and symmetric difference

`unique` and `union` retain first occurrences. Intersection and difference follow
the first input's order. Equal numeric representations and nested values use
Ferret equality, with hash collisions checked for equality.

`arrays::symmetric_difference` keeps values present in an odd number of input
arrays. Duplicates within one array count once. Output follows first encounter
across inputs and retains the first representation of an equal value.

{{< code lang="fql" >}}
return {
    xor: arrays::symmetric_difference([7, 3, 7], [7, 2], [7, 2, 5]),
    legacy: outersection([7, 3, 7], [7, 2], [7, 2, 5])
}
{{</ code >}}

The XOR result is `[7,3,5]`. Legacy `outersection` returns `[3,5]`: it keeps values
present in exactly one input array, excluding `7` even though it appears in three
inputs. The two operations agree for two inputs, but migrating three or more
inputs requires choosing the intended rule. There is no canonical
`arrays::outersection` or `arrays::exclusive`.

## Slice and copy behavior

`arrays::slice(xs, start[, length])` takes a length, not an end index. Negative
starts or lengths return `[]`; an oversized length is capped at the remaining
array length. `flatten` retains its optional depth, defaulting to one; zero and
negative depths leave nested lists unexpanded.

Returned native arrays have independent top-level backing storage. Replacing,
appending, removing, or sorting elements in a slice cannot change its source or
another slice. Copies remain shallow: nested arrays, objects, and host resources
retain their identities. Host-provided lists implement their own copy contracts.

## Explicit mutation

Use `arrays::mut::` when you want to change an existing array. Immutable array
transformations and legacy globals such as `push` and `pop` continue to return
independent arrays.

{{< code lang="fql" >}}
LET values = [1, 2]
LET copied = arrays::append(values, 3)
LET changed = arrays::mut::push(values, 4)
LET removed = arrays::mut::pop(values)
RETURN {values, copied, changed, removed}
{{</ code >}}

The result is `{values: [1,2], copied: [1,2,3], changed: [1,2], removed: 4}`.
`changed` is the original array, so it also reflects the later `pop`. A `LET`
binding prevents reassignment of the binding; it does not freeze its array.

| Call under `arrays::mut::` | Result |
| --- | --- |
| `push(array, value)` | Original array with one element appended |
| `unshift(array, value)` | Original array with one element prepended |
| `set(array, index, value)` | Original array with an existing element replaced |
| `insert(array, index, value)` | Original array with an element inserted |
| `remove(array, value)` | Original array with all matching elements removed |
| `clear(array)` | Original array with all elements removed |
| `sort(array)` | Original array, stably sorted in ascending order |
| `pop(array)` | Removed last value, or `none` when empty |
| `shift(array)` | Removed first value, or `none` when empty |
| `remove_at(array, index)` | Removed value |

Indexes must be integers. `set` and mutable `remove_at` require an existing
index. `insert` accepts zero through the array length, including insertion at
the end. Negative or missing indexes produce an error; `set` never grows the
array. Immutable `arrays::remove_at` still returns an unchanged copy for a
missing index.

`remove` preserves the order of remaining elements and uses the same equality
as immutable `arrays::remove`. It accepts no removal limit. `sort` uses the
same stable ordering as `arrays::sorted` and accepts no direction argument.

{{< code lang="fql" >}}
LET original = [1, 2, 3]
LET part = arrays::slice(original, 0, 2)
LET changed = arrays::mut::set(part, 0, 9)
RETURN {original, part}
{{</ code >}}

This returns `{original: [1,2,3], part: [9,2]}`. Native slices and copies have
independent top-level storage. Copying remains shallow: nested arrays, objects,
and host resources keep their identities. Mutating a nested array can therefore
be visible through both containers.

Host-provided targets must implement the runtime list mutation contract.
Missing capabilities, nil targets, invalid indexes, and host mutation refusals
return errors. Mutable calls never silently copy a read-only target. No-op
extraction or removal does not test mutability through speculative writes.
A host, comparison, or cancellation failure can leave earlier changes applied;
there is no transaction or rollback. Mutation does not clone or close element
values.

There are no global mutable aliases or `mut::*` shortcuts. Global `push`,
`pop`, `shift`, and `unshift` keep their legacy immutable behavior.

## Go callers

The exported functions in `pkg/stdlib/arrays` adopt the canonical v2 names and
signatures. In particular, `Union` becomes distinct union, `Concat` concatenates,
and `Append` and `Remove` accept exactly two values after the context. `Slice`
rejects extra arguments beyond its optional length. Obsolete Go names and
mode-bearing wrappers are removed. Migration compatibility applies to global
FQL registrations, not the old Go functions.

Mutable Go entry points use the `Mutable` suffix, such as `arrays.PushMutable`
and `arrays.RemoveAtMutable`, and accept custom `runtime.List` implementations.
