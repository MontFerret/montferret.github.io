---
title: "Arithmetic Operators"
sidebarTitle: "Arithmetic"
weight: 30
draft: false
description: "Addition, subtraction, multiplication, division, modulus, and the operand types accepted by arithmetic operators."
---

# Arithmetic operators

FQL arithmetic is defined directly over native numbers and temporal values. Operators do not implicitly route arbitrary values through `TO_NUMBER`, `TO_DURATION`, or `TO_DATETIME`.

FQL supports:

- `+` addition or String-triggered concatenation
- `-` subtraction
- `*` multiplication
- `/` division
- `%` modulus
- unary `+` and `-`

For exponentiation, use `POW()`. The syntax `base ** exponent` is not supported.

## Native numeric arithmetic

Numeric arithmetic accepts only `Int` and `Float` values.

{{< editor lang="fql" >}}
return {
    integer: 1 + 2,
    mixed: 1 + 2.5,
    product: 2.5 * 3,
    exactDivision: 6 / 3,
    fractionalDivision: 5 / 2,
    floatRemainder: 5.5 % 2
}
{{</ editor >}}

`Int + Int`, `Int - Int`, `Int * Int`, and `Int % Int` return an `Int`. Exact division of two integers also returns an `Int`; division with a remainder returns a `Float`. Any numeric pair containing a `Float` returns a `Float`, and floating-point modulus uses the ordinary remainder operation.

Arithmetic checks integer overflow, division and modulus by zero, and non-finite floating-point results. These conditions raise runtime errors instead of wrapping or producing `NaN` or infinity.

Zero divisors produce the specific diagnostics `division by zero` and `modulo by zero`.

Numeric-looking strings are not numbers. Convert them explicitly:

{{< editor lang="fql" >}}
let input = "10"

return TO_NUMBER(input) - 2
{{</ editor >}}

{{< notification type="info" >}}
Implicit numeric-string, Boolean, <code>none</code>, collection, Binary, and opaque host-value arithmetic is no longer supported. Use an explicit <code>TO_*</code> function when conversion is intended.
{{</ notification >}}

## String concatenation

If either `+` operand is an actual String, Ferret concatenates both operands' String representations. The String may appear on either side and may be combined with any runtime value.

{{< editor lang="fql" >}}
return [
    "a" + 1,
    1 + "a",
    "enabled=" + true,
    none + " value",
    [1, 2] + " items",
    1s + " elapsed"
]
{{</ editor >}}

Expressions are evaluated from left to right. A later String does not rescue an earlier invalid pair:

{{< code lang="fql" >}}
true + 1 + " items" // fails at true + 1
{{</ code >}}

## Temporal arithmetic

DateTime and Duration arithmetic accepts only the native operand pairs below. Convert text or other values before applying an operator.

| Expression | Result |
| --- | --- |
| `Duration + Duration` | Duration |
| `Duration - Duration` | Duration |
| `Duration * Number` | scaled Duration |
| `Number * Duration` | scaled Duration |
| `Duration / Number` | scaled Duration |
| `Duration / Duration` | Int for an exact ratio, otherwise Float |
| `DateTime + Duration` | DateTime |
| `Duration + DateTime` | DateTime |
| `DateTime - Duration` | DateTime |
| `DateTime - DateTime` | elapsed Duration |
| unary `+Duration` or `-Duration` | Duration |

{{< editor lang="fql" >}}
let start = TO_DATETIME("2024-03-10T06:30:00Z")

return {
    combined: 1s + 250ms,
    multiplied: 2.5 * 5s,
    scaled: 5s / 2,
    ratio: 5s / 250ms,
    later: start + TO_DURATION("90m"),
    elapsed: (start + 90m) - start
}
{{</ editor >}}

All Duration scaling truncates fractional nanoseconds toward zero. Duration and DateTime range overflow raises an error.

String concatenation still takes precedence for `+`. For example, `1s + "1s"` returns the String `"1s1s"`; use `1s + TO_DURATION("1s")` for temporal addition.

Unsupported temporal pairs include `DateTime + DateTime`, `Duration - DateTime`, reverse division such as `2 / 1s`, Duration modulus, and DateTime multiplication, division, or modulus.

## Unary operators

Unary `+` and `-` accept only `Int`, `Float`, and `Duration`. Other operand types raise an invalid-operation error.

{{< editor lang="fql" >}}
let x = -5
let interval = 500ms

return [-x, +interval, -interval]
{{</ editor >}}

## Invalid operations

An unsupported operand pair reports the source operator and preserves operand order:

{{< code lang="text" >}}
operator '-' cannot be applied to String and Int
operator '+' cannot be applied to Boolean and Int
operator '/' cannot be applied to Duration and String
{{</ code >}}

Equality remains valid across incompatible types; these errors apply only to operators that require a supported arithmetic pair.

## Next steps

{{< docs-related tiles="language-operators,language-operators-comparison,language-types-basic" >}}
