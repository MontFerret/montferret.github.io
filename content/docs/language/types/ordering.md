---
title: "Type Ordering"
sidebarTitle: "Type Ordering"
weight: 40
draft: false
description: "How Ferret compares native values, including cross-type order and strict Duration behavior."
aliases:
    - /docs/fql/type-value-order/
---

# Type ordering

FQL defines a relational order for the scalar and collection types listed below. It is used by `<`, `<=`, `>`, and `>=`, sorting, and recursive array and object comparison.

When two values in this sequence have different types, their type decides the result:

{{< code >}}
none < bool < number < string < datetime < binary < array < object
{{</ code >}}

Values are compared within their type only after the types match. `Int` and `Float` share the numeric comparison domain and compare by numeric value.

{{< code lang="fql" >}}
none < false
false < true
true < 0
1 < 2.5
0 < "0"
"abc" < []
[] < {}
{{</ code >}}

## Duration values

Duration is intentionally separate from the cross-type chain. Two Durations compare by signed nanosecond value, and equivalent units compare equal.

{{< editor lang="fql" >}}
return {
    less: 500ms < 1s,
    equal: 1000ms == 1s
}
{{</ editor >}}

A native Duration and any non-Duration value are unequal. Relational comparison between them is invalid in either operand order.

{{< code lang="fql" >}}
1s == 1000                 // false
1s != "1s"                 // true
1s < 1000                  // runtime error
"1s" >= 1s                // runtime error
{{</ code >}}

Use `to_duration` explicitly when the other value should be interpreted as a Duration.

Sorting uses the same relational contract, so a collection that mixes Duration with another type cannot be sorted without first normalizing its values.

{{< code lang="fql" >}}
return arrays::sorted([1s, "2s"]) // runtime error
{{</ code >}}

## Equality, membership, and uniqueness

Equality does not use the cross-type order to make different types equal. Membership, `match`, grouping, set operations, and deduplication all verify canonical equality recursively.

Equivalent native Durations collapse to the first representative. This example retains three values, the Duration `1s`, the string `"1s"`, and the integer `1000`. The Duration and the string have the same representation when serialized to JSON.

{{< code lang="fql" >}}
return distinct [1s, 1000ms, "1s", 1000]
{{</ code >}}

Hashes may select candidate buckets internally, but equality determines whether values are actually the same.

## Primitive values

Within a built-in type:

- `none` is equal only to `none`.
- Booleans use `false < true`.
- Numbers compare by numeric value across `Int` and `Float`.
- Durations compare by signed nanoseconds, only with Duration.
- Strings compare lexically and case-sensitively.
- DateTime values compare by canonical instant.
- Binary values compare by length first, then by byte contents at equal lengths.

DateTime comparison does not parse strings or convert numeric epoch values. Use `to_datetime` before comparison when conversion is intended.

## Arrays

The shorter array comes first. Arrays of equal length compare element by element from left to right until two elements differ.

{{< code lang="fql" >}}
[] < [0]
[1] < [2]
[1, 2] < [1, 3]
[9] < [0, 0]
{{</ code >}}

Array comparison is recursive. A nested Duration/non-Duration pair remains unequal for equality and invalid for relational comparison.

{{< code lang="fql" >}}
[1s] == ["1s"] // false
[1s] < ["1s"]  // runtime error
{{</ code >}}

## Objects

Objects compare by number of properties first. For equal sizes, Ferret walks the keys in sorted order. At the first pair of different keys, the lexically larger key sorts first. Matching keys compare their values recursively. Equality requires exactly the same keys and equal values, including properties explicitly set to `none`.

{{< code lang="fql" >}}
{} < { "a": 1 }
{} != { "a": none }
{ "b": 1 } < { "a": 1 }
{ "a": 1 } < { "a": 2 }
{ "a": true } < { "a": 0 }
{{</ code >}}

Declaration order does not affect equality:

{{< editor lang="fql" height="110px" >}}
return { "a": 1, "b": 2 } == { "b": 2, "a": 1 }
{{</ editor >}}

## Ranges

Two ranges compare by their integer endpoints, first by start and then by end. They are equal when both endpoints match.

## Host values

Host values participate only through compatible comparison capabilities supplied by the embedding application. Equality and relational comparison are separate capabilities; an opaque host value does not become comparable through its String representation or a generic numeric conversion.

See [Capability Types]({{< ref "capabilities" >}}) and [Host Values]({{< ref "host" >}}) for the host contracts.

## Next steps

{{< docs-related tiles="language-operators-comparison,language-types-basic,language-types-capabilities" >}}
