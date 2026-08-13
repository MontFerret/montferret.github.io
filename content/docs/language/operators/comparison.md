---
title: "Comparison Operators"
sidebarTitle: "Comparison"
weight: 10
draft: false
description: "Equality, relational comparison, containment, pattern matching, and regular expression operators."
---

# Comparison operators

Comparison operators test two values. A valid comparison returns a Boolean; a relational operator raises a runtime error when the operand pair has no defined comparison.

FQL provides:

- `==` and `!=` for equality
- `<`, `<=`, `>`, and `>=` for relational comparison
- `in` and `not in` for containment
- `like` and `not like` for wildcard matching
- `=~` and `!~` for regular expression matching

## Equality

Equality and inequality are valid across all operand types. Incompatible values are unequal: `==` returns false and `!=` returns true. FQL does not implicitly convert strings, numbers, DateTime values, or Duration values before testing equality.

{{< editor lang="fql" >}}
return {
    mixedNumeric: 1 == 1.0,
    numericString: 65 == "65",
    noneAndZero: none == 0,
    sameText: "abc" == "abc",
    incompatible: 1s != "1s"
}
{{</ editor >}}

`Int` and `Float` are the numeric exception: they compare by numeric value. Arrays and objects compare recursively, so nested values follow these same rules. Object property declaration order does not affect equality.

`match`, membership, grouping, set operations, and deduplication use the same canonical equality semantics.

## Relational comparison

Values of the same built-in type compare according to that type. Mixed `Int`/`Float` pairs compare numerically. Non-Duration built-in types use FQL's deterministic cross-type order instead of implicit conversion.

{{< code lang="fql" >}}
1 < 2.5
"abc" < "abd"
false < 0
to_datetime("2026-08-02T12:00:00Z") > "not-a-date"
{{</ code >}}

See [Type Ordering]({{< ref "../types/ordering" >}}) for structural and cross-type details.

### Duration comparisons

Native Duration values compare only with other native Duration values. Equivalent units normalize to the same value.

{{< editor lang="fql" >}}
return {
    equivalent: 1s == 1000ms,
    greater: 5s > 4999ms,
    stringIsDifferent: 1s == "1s",
    numberIsDifferent: 1s != 1000,
    explicitConversion: 1s == to_duration("1s")
}
{{</ editor >}}

Equality remains non-failing for incompatible types. Relational Duration/non-Duration pairs have no defined result and report the actual source operator with operands in source order:

{{< code lang="text" >}}
operator '>' cannot be applied to String and Duration
operator '<=' cannot be applied to Duration and String
{{</ code >}}

The same rule applies recursively. For example, `[1s] == ["1s"]` is false, while `[1s] < ["1s"]` raises an invalid-operation error. Element-wise `any`, `all`, and `none` comparisons also use this strict behavior.

{{< notification type="info" >}}
Duration comparison no longer converts numbers or strings implicitly. Use <code>to_duration(value)</code> before comparison when conversion is intended.
{{</ notification >}}

DateTime values compare by canonical instant when both operands are DateTime. A string is never parsed implicitly; mixed DateTime/non-Duration relational comparisons continue to use the normal cross-type order.

## Containment

`in` tests the left value against the container on the right:

- an Array or other runtime List tests whether an equal element exists;
- an Object or other runtime Map tests whether an equal value exists;
- a String tests whether it contains the left operand's String representation;
- any other right operand returns false.

{{< editor lang="fql" >}}
return {
    arrayValue: 1 in [2, 3, 1],
    objectValue: "Ada" in { name: "Ada" },
    substring: "err" in "runtime error",
    unsupported: "foo" in none
}
{{</ editor >}}

Object containment checks values, not property names. Use the appropriate property-access or existence operation when testing a key.

`not in` negates the `in` result. Consequently, it returns true for an unsupported right operand.

## Pattern matching

`like` compares a String against a case-sensitive wildcard pattern. `*` matches any sequence and `?` matches one character. `not like` returns the opposite result.

{{< code lang="fql" >}}
"foo" like "f*"
"abc" like "?bc"
"abc" not like "a*"
{{</ code >}}

## Regular expressions

`=~` returns true when the String on the left matches the regular expression on the right. `!~` returns the opposite result.

{{< code lang="fql" >}}
"foo" =~ "^f[o].$"
"foo" !~ "[a-z]+bar$"
{{</ code >}}

Pattern and regular expression operators require String operands; they do not implicitly stringify other values.

## Next steps

{{< docs-related tiles="language-operators,language-types-ordering,language-control-flow-match" >}}
