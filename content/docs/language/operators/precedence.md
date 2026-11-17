---
title: "Operator Precedence"
sidebarTitle: "Precedence"
weight: 60
draft: false
description: "The evaluation order of FQL operators, from lowest to highest precedence."
---

# Operator precedence

The operator precedence in FQL is listed below from lowest to highest. This order follows the FQL grammar; do not assume unary and comparison operators have the same relative order as another language.

- ``? :`` ternary operator
- ``??`` none-coalescing operator
- ``||`` logical or
- ``&&`` logical and
- ``!``, ``+``, ``-`` logical negation, unary plus, unary minus
- ``like``, ``not like`` pattern matching
- ``in``, ``not in`` containment
- ``all``, ``any``, ``none`` array comparison operators
- ``==``, ``!=``, ``<``, ``<=``, ``>=``, ``>`` comparison operators
- ``=~``, ``!~`` regular expression matching
- ``+``, ``-`` addition, subtraction
- ``*``, ``/``, ``%`` multiplication, division, modulus
- ``..`` range and primary expressions
- ``()`` function call
- ``?.`` optional chaining
- ``.`` member access
- ``[]`` indexed value access

Operators nearer the bottom bind more tightly. Assignment is a statement and cannot be nested in an expression. For example, multiplication is evaluated before addition, logical and before logical or, and comparisons before unary and logical operators.

{{< editor lang="fql" >}}
// Multiplication binds tighter than addition:
// interpreted as 2 + (3 * 4), not (2 + 3) * 4
return 2 + 3 * 4
{{</ editor >}}

{{< editor lang="fql" >}}
// The result is true: (false && true) || true.
return false && true || true
{{</ editor >}}

{{< code lang="fql" >}}
let cached = none
let fetched = none
let active = true
let nickname = none

// OR binds tighter than NONE coalescing:
// interpreted as (cached OR fetched) ?? "fallback"
// NONE coalescing binds tighter than the ternary operator:
// interpreted as active ? (nickname ?? "Anonymous") : "Inactive"
return {
    cached: cached or fetched ?? "fallback",
    status: active ? nickname ?? "Anonymous" : "Inactive"
}
{{</ code >}}

When `or` precedes `and`, write the grouping explicitly. For example, `false || (true && true)` compiles in this release. The compiler can report an ambiguity for the unparenthesized form.

## Using parentheses

Parentheses ``(`` and ``)`` group expressions. Use them when the intended grouping differs from the precedence rules. Most binary operators associate to the left; `??` associates to the right, so `a ?? b ?? c` means `a ?? (b ?? c)`.

{{< editor lang="fql" >}}
return (2 + 3) * 4
{{</ editor >}}

{{< editor lang="fql" height="150px" >}}
let price = 120
let discount = 0.1
let tax = 0.2

// Calculate the discount and tax factors before multiplying the price.
return price * (1 - discount) * (1 + tax)
{{</ editor >}}

The formatter removes grouping that does not affect syntax or semantics. For example, it writes `1 + (2 * 3)` as `1 + 2 * 3`, but keeps `(1 + 2) * 3`, `a - (b - c)`, and `(a ?? b) ?? c`. It also keeps parentheses required by grammar boundaries, recovery-tail ownership, comments, or lexical safety such as `-(-value)`.

## Next steps

{{< docs-related tiles="language-operators,language-operators-coalescing,language-expressions,language-operators-arithmetic" >}}
