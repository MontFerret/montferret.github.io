---
title: "Type Ordering"
sidebarTitle: "Type Ordering"
weight: 50
draft: false
description: "How Ferret compares and orders values of different types."
aliases:
    - /docs/fql/type-value-order/
---

# Type ordering

FQL comparisons are deterministic. Any two values can be compared, even when they do not have the same type.

When two values are compared, FQL first looks at their types. If the types are different, the result is decided by the global type order. The actual values are only compared when both operands have the same type.

FQL uses the following type order:

{{< code >}}
none < bool < number < string < array < object
{{</ code >}}

This means that none sorts before every other value, while objects sort after every other built-in value type.

For example, a boolean value always sorts before a number, a string, an array, or an object. A string always sorts after a number, even if the string is empty or contains numeric-looking text.

{{< code lang="fql" >}}
none < false
none < 0
none < ""
none < []
none < {}

false < true
true < 0

0 < ""
0 < "0"
0 < "abc"
0 < []

"" < "abc"
"abc" < []

[] < {}
{{</ code >}}

Once the types are the same, FQL compares the values according to the rules for that type.

## Primitive values

Primitive values are ordered as follows:

- none is only equal to none.
- Booleans are ordered as false < true.
- Numbers are ordered by numeric value.
- Strings are ordered using FQL's string comparison rules.

{{< code lang="fql" >}}
none == none

false < true

1 < 2
10 > 2

"a" < "b"
{{</ code >}}

<div class="notification is-info">
  none is a regular comparable value in FQL. Comparing a value with none does not produce an unknown result.
</div>

## Arrays

Arrays are compared element by element from left to right.

At each position, FQL compares the two elements using the same comparison rules described on this page. If the elements are different, that result decides the array comparison. If the elements are equal, FQL moves to the next position.

If all compared elements are equal, the shorter array sorts first.

{{< code >}}
[] < [0]

[1] < [2]

[false] < [true]

[1, 2] < [1, 3]

[1] < [1, 0]
{{</ code >}}

Array comparison is recursive. If an element is another array or an object, that nested value is compared using the same rules.

{{< code >}}
[[1]] < [[2]]

[{ "score": 1 }] < [{ "score": 2 }]
{{</ code >}}

## Objects

Objects are compared by their attributes, not by the order in which those attributes were written.

Before two objects are compared, FQL considers their attribute names in sorted order. For each attribute name, FQL compares the corresponding values from both objects.

If one object has an attribute that the other object does not have, the missing value is treated as none for comparison purposes.

{{< code >}}
{} < { "a": 1 }

{} == { "a": none }

{ "a": 1 } < { "a": 2 }

{ "a": true } < { "a": 0 }

{ "a": { "score": 1 } } < { "a": { "score": 2 } }
{{</ code >}}

Attribute declaration order does not affect equality:

{{< code lang="fql" height="110px" >}}
{ "a": 1, "b": 2 } == { "b": 2, "a": 1 }
{{</ code >}}

If all attributes compare as equal, the objects are considered equal.