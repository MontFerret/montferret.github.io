---
title: "Array API Migration"
sidebarTitle: "Array API Migration"
weight: 75
draft: false
description: "Use the immutable arrays namespace while preserving legacy global calls during migration."
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
These functions return new arrays. This change does not add `arrays::mut::*`.

## Choose the canonical operation

| Global call | Canonical operation |
| --- | --- |
| `first`, `last`, `flatten`, `slice`, `unique`, `sorted` | Same name under `arrays::` |
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

## Go callers

The exported functions in `pkg/stdlib/arrays` adopt the canonical v2 names and
signatures. In particular, `Union` becomes distinct union, `Concat` concatenates,
and `Append` and `Remove` accept exactly two values after the context. `Slice`
rejects extra arguments beyond its optional length. Obsolete Go names and
mode-bearing wrappers are removed. Migration compatibility applies to global
FQL registrations, not the old Go functions.
