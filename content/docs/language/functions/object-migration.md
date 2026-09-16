---
title: "Object functions and migration"
sidebarTitle: "Object migration"
weight: 45
draft: false
description: "Use immutable object functions, opt into mutation, and migrate legacy global calls."
---

# Object functions and migration

The canonical v2 object library uses the `object::` namespace. In runtime
releases containing the compatibility layer, seven legacy global names remain
available temporarily as deprecated aliases. They use the canonical signatures
and behavior; use `object::` for new code.

```fql
let base = {name: "Ferret", internal: true}
let clean = object::omit_keys(base, "internal")
let result = object::merge(clean, {version: 2})
return object::entries(result)
```

Immutable object transformations return independent containers. Nested values are deep
cloned when supported; other host values follow their shallow-copy contract.
The source object is unchanged. Explicit mutation is available through
`object::mut::` in runtime releases containing the mutable object API.

## Functions

| Call | Behavior |
| --- | --- |
| `object::keys(value)` | Return keys in unspecified order |
| `object::values(value)` | Return copied values in unspecified order |
| `object::entries(value)` | Return associated `["key", value]` pairs in unspecified order |
| `object::has_key(value, key)` | Check whether a String key exists, including a present `none` value |
| `object::keep_keys(value, keys...)` | Retain only the requested String keys |
| `object::omit_keys(value, keys...)` | Remove the requested String keys |
| `object::merge(values...)` | Merge maps, with later values replacing earlier ones |
| `object::merge_deep(values...)` | Recursively merge conflicting maps |
| `object::zip(keys, values)` | Construct an object from equal-length parallel lists |
| `object::from_entries(entries)` | Construct an object from an iterable of two-item pairs |

Both merge functions accept variadic maps or a single list of maps. An empty
list returns an empty object; calling either function without arguments fails.
Deep merge replaces arrays, scalars, and `none` rather than combining them.

For host maps in runtimes containing the collection factory API, both immutable
merges construct their destination through the first source's `New(ctx)`.
The empty destination retains that map's implementation family and backend
configuration; incoming values still follow the clone/copy rules above. After
construction succeeds, a population or cancellation failure closes the incomplete
destination when it is closable, preserving both the primary and cleanup errors.
Successful results stay open until normal result cleanup. Source maps and mutable
merge targets remain borrowed; this cleanup does not provide rollback. The factory
remains responsible for cleanup when construction itself fails.
Host authors should follow the [collection factory migration]({{< ref "docs/embedding/go/host-values" >}}#constructing-an-empty-collection).


Key filters accept variadic String keys or one list of String keys. At least
one key argument is required. Missing or repeated keys are harmless.
`object::keep_keys(value, [])` returns an empty map;
`object::omit_keys(value, [])` returns an independent copy.

## Explicit mutation

Standard library operations are immutable by default. A `::mut::` subnamespace
explicitly opts into mutation of an existing value:

```fql
LET original = {name: "Ferret", internal: true}
LET clean = object::omit_keys(original, "internal")
// original still includes internal.
LET alias = original
LET result = object::mut::omit_keys(original, "internal")
RETURN [original, alias, result, clean]
// All four values are {name: "Ferret"}.
```

`object::merge(a, b)` returns an independent value without modifying `a`.
`object::mut::merge(a, b)` modifies `a` and returns that same target. Top-level
aliases observe the changes, and the returned target can be composed with
another mutable operation.

| Call | Behavior |
| --- | --- |
| `object::mut::merge(target, sources...)` | Merge into the target; later sources win |
| `object::mut::merge_deep(target, sources...)` | Merge recursively, copying conflicting nested branches before modification |
| `object::mut::keep_keys(target, keys...)` | Remove every unselected key from the target |
| `object::mut::omit_keys(target, keys...)` | Remove selected keys from the target |

Mutable merges accept variadic source maps or one list after the target.
`object::mut::merge(target)` and `object::mut::merge(target, [])` return the
unchanged target; the same applies to `merge_deep`. The target itself must be a
map, not a list of maps.

Mutable key filters still require a key argument. An explicit empty list clears
the target for `keep_keys` and leaves it unchanged for `omit_keys`. Missing and
repeated keys are harmless. Retained values are not copied.

Deep mutation preserves the target's top-level identity. It copies a conflicting
nested branch, merges into the copy, and replaces the target's branch after
success. Source objects and external aliases to the original nested branch stay
unchanged; untouched branches retain their identity. Incoming values use the
same clone/copy contracts as immutable merges, including host-value guarantees.

Host-backed targets must implement the runtime map mutation contract. Missing
capabilities and rejected writes fail the call; there is no fallback to mutating
a copy of the target. No-op calls do not issue writes to probe a host's policy.

Mutation is not transactional. An error can leave earlier updates applied,
including a host write that changes state before reporting failure. A failed
nested merge does not replace its original branch. Iteration, cloning, mutation,
and cleanup errors propagate to the caller.

## Entries and constructors

```fql
let value = {name: "Ferret", version: 2}
return object::from_entries(object::entries(value)) == value
```

Entries preserve the association between each key and its value. Separate calls
to `keys` and `values` do not promise matching positions.
Every entry key is a String; only the value is cloned or copied. A host runtime
map exposing a non-string key causes `entries` to fail with a type error.

`from_entries` accepts any runtime iterable and consumes it once. Each entry
must be an indexed, measurable pair of exactly two values, with a String key
first. Invalid pairs and host iteration/access errors fail the call.

Both constructors use **last-key-wins**:

```fql
return [
    object::zip(["a", "a"], [1, 2]),
    object::from_entries([["a", 1], ["a", 2]])
]
// Both objects are {a: 2}.
```

This intentionally changes v1 ZIP, which kept the first value for duplicate
keys. The deprecated global alias also uses last-key-wins. Review any code
relying on the v1 behavior, even before rewriting its function name.

## Migrate global calls

The seven global aliases below are deprecated in generated Core API metadata,
with a message naming the canonical replacement. This metadata does not produce
compiler or runtime deprecation warnings. The compatibility layer is temporary;
no removal release is specified.

There are no global aliases for the new `object::entries`,
`object::from_entries`, or `object::omit_keys` APIs, or for mutable operations.
Mutation remains explicit under `object::mut::`.

The CLI migration remains available to rewrite legacy calls case-insensitively:

| Legacy call | Replacement |
| --- | --- |
| `KEYS(value)` | `object::keys(value)` |
| `VALUES(value)` | `object::values(value)` |
| `HAS(value, key)` | `object::has_key(value, key)` |
| `KEEP_KEYS(value, keys...)` | `object::keep_keys(value, keys...)` |
| `MERGE(values...)` | `object::merge(values...)` |
| `MERGE_RECURSIVE(values...)` | `object::merge_deep(values...)` |
| `ZIP(keys, values)` | `object::zip(keys, values)` |

The compatibility alias for keys accepts one argument, and the CLI rewrites
only `KEYS(value)`. Every other arity, including `KEYS(value, true)` and
`KEYS(value, false)`, remains unchanged and receives a manual-review explanation.

Calls potentially resolved by a local function declaration or explicit function
alias are also preserved for manual review. These guards are case-insensitive
and include forward and nested declarations. Namespace aliases block only
replacements whose canonical targets they would redirect. Already-qualified
calls are preserved.

Other safe calls in the same file can still migrate. If rewriting or formatting
fails, the whole file remains unchanged, and previously discovered manual
actions are reported alongside the file-level failure.

Preview changes with `ferret migrate run --print path/to/query.fql`.
See [Migrate]({{< ref "/docs/tools/cli/migrate" >}}) for checking, applying,
and reviewing source migrations. Use a CLI release containing this migration
together with the corresponding object-capable runtime release.
