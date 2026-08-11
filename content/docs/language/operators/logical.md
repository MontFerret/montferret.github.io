---
title: "Logical Operators"
sidebarTitle: "Logical"
weight: 20
draft: false
description: "Logical AND, OR, and strict Boolean negation, including short-circuit evaluation and result values."
---

# Logical operators

FQL supports symbolic and keyword forms of the logical operators:

- `&&` or `AND`
- `||` or `OR`
- `!` or `NOT`

The paired forms have the same behavior.

{{< code lang="fql" >}}
RETURN true && false
RETURN true OR false
RETURN NOT false
{{</ code >}}

## Truth values for AND and OR

Binary `AND` and `OR` evaluate each operand according to its Boolean representation:

- `NONE` is false.
- Booleans keep their value.
- Numeric zero is false; other numbers are true.
- A zero Duration or DateTime is false; other temporal values are true.
- An empty String is false; other strings are true.
- Arrays and objects are true, even when empty.
- Binary and host values are true.

The conversion is used only to decide control flow. `AND` and `OR` return one of their original operands rather than an automatically converted Boolean.

## Logical AND

`AND` evaluates the left operand first. If it is false, the expression returns that operand without evaluating the right operand. Otherwise, it evaluates and returns the right operand.

{{< editor lang="fql" height="150px" >}}
LET user = {
    active: true,
    name: "Ada"
}

RETURN user.active && user.name
{{</ editor >}}

This returns `"Ada"` because the left operand is true.

{{< code lang="fql" >}}
false && "value"  // false
NONE AND true      // NONE
0 && "fallback"   // 0
true && 23         // 23
{{</ code >}}

## Logical OR

`OR` evaluates the left operand first. If it is true, the expression returns that operand without evaluating the right operand. Otherwise, it evaluates and returns the right operand.

{{< editor lang="fql" height="150px" >}}
LET user = {
    displayName: ""
}

RETURN user.displayName || "Anonymous"
{{</ editor >}}

{{< code lang="fql" >}}
true || "value"       // true
1 OR 7                 // 1
NONE || "fallback"    // "fallback"
"" || "fallback"      // "fallback"
{{</ code >}}

{{< notification type="info" >}}
Use <code>??</code> when only <code>NONE</code> should select the fallback. Unlike <code>OR</code>, NONE coalescing preserves <code>false</code>, zero, and empty strings.
{{</ notification >}}

See [NONE-Coalescing Operator]({{< ref "coalescing" >}}) for the fallback semantics and examples.

## Logical NOT

Unary `!` and `NOT` accept only Boolean operands and always return a Boolean.

{{< editor lang="fql" >}}
RETURN {
    notTrue: !true,
    notFalse: NOT false
}
{{</ editor >}}

Other types produce an operator-oriented runtime error:

{{< code lang="text" >}}
operator '!' cannot be applied to String
operator '!' cannot be applied to Int
{{</ code >}}

Use `TO_BOOL(value)` when explicit Boolean conversion is intended.

{{< notification type="info" >}}
Double negation is no longer an implicit conversion mechanism. Replace expressions such as <code>!!value</code> with <code>TO_BOOL(value)</code>.
{{</ notification >}}

## Short-circuit evaluation

The right side of `AND` is evaluated only when the left side is true. The right side of `OR` is evaluated only when the left side is false.

{{< editor lang="fql" >}}
LET user = {
    active: false,
    name: "Ada"
}

RETURN user.active && user.name
{{</ editor >}}

Subqueries are an exception to this source-level model. Query planning may evaluate a subquery operand before the logical operator, so short-circuiting should not be used to suppress a subquery's execution.

## Result values

`!` and `NOT` always return a Boolean. `AND` and `OR` return an operand, so their result may have any FQL type.

{{< code lang="fql" >}}
RETURN 25 > 1 && 42 != 7
RETURN 22 IN [23, 42] || 23 NOT IN [22, 7]
RETURN NONE || "fallback"
RETURN true && 23
{{</ code >}}

When a strict Boolean result is required from a binary logical expression, pass the result to `TO_BOOL`.

## Next steps

{{< docs-related tiles="language-operators,language-operators-coalescing,language-operators-comparison,language-control-flow-match" >}}
