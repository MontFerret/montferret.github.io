---
title: "Object functions and migration"
sidebarTitle: "Object migration"
weight: 45
draft: false
description: "Use immutable object functions and migrate legacy global calls."
---

# Object functions and migration

The new v2 alpha object library uses the `object::` namespace. Use a runtime
release containing this namespace when adopting these calls; earlier alpha
releases expose the legacy global functions.

```fql
let base = {name: "Ferret", internal: true}
let clean = object::omit_keys(base, "internal")
let result = object::merge(clean, {version: 2})
return object::entries(result)
```

Object transformations return independent containers. Nested values are deep
cloned when supported; other host values follow their shallow-copy contract.
The source object is unchanged. There is no `object::mut` API yet.

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

Key filters accept variadic String keys or one list of String keys. At least
one key argument is required. Missing or repeated keys are harmless.
`object::keep_keys(value, [])` returns an empty map;
`object::omit_keys(value, [])` returns an independent copy.

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

This intentionally changes the old global `ZIP`, which kept the first value
for duplicate keys. Review any code relying on that behavior during migration.

## Migrate global calls

The CLI's object source migration recognizes legacy calls case-insensitively:

| Legacy call | Replacement |
| --- | --- |
| `KEYS(value)` | `object::keys(value)` |
| `VALUES(value)` | `object::values(value)` |
| `HAS(value, key)` | `object::has_key(value, key)` |
| `KEEP_KEYS(value, keys...)` | `object::keep_keys(value, keys...)` |
| `MERGE(values...)` | `object::merge(values...)` |
| `MERGE_RECURSIVE(values...)` | `object::merge_deep(values...)` |
| `ZIP(keys, values)` | `object::zip(keys, values)` |

`KEYS(value, true)` becomes `sorted(object::keys(value))`;
`KEYS(value, false)` becomes `object::keys(value)`.
The object expression is evaluated once. Dynamic sorting expressions and
shadowed sorting functions require manual migration; the CLI reports the
source location and leaves the file unchanged. It preserves user-defined
functions and qualified calls.

Preview changes with `ferret migrate run --print path/to/query.fql`.
See [Migrate]({{< ref "/docs/tools/cli/migrate" >}}) for checking, applying,
and reviewing source migrations. Use a CLI release containing this migration
together with the corresponding object-capable runtime release.
