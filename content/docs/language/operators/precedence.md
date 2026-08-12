---
title: "Operator Precedence"
sidebarTitle: "Precedence"
weight: 60
draft: false
description: "The evaluation order of FQL operators, from lowest to highest precedence."
---

# Operator precedence

The operator precedence in FQL is similar as in other familiar languages (lowest precedence first):

- ``=``, ``+=``, ``-=``, ``*=``, ``/=`` assignment and compound assignment
- ``? :`` ternary operator
- ``??`` none-coalescing operator
- ``||`` logical or
- ``&&`` logical and
- ``==``, ``!=`` equality and inequality
- ``in``, ``not in`` containment
- ``like``, ``not like`` pattern matching
- ``=~``, ``!~`` regular expression matching
- ``<``, ``<=``, ``>=``, ``>`` less than, less equal, greater equal, greater than
- ``..`` range
- ``+``, ``-`` addition, subtraction
- ``*``, ``/``, ``%`` multiplication, division, modulus
- ``!``, ``+``, ``-`` logical negation, unary plus, unary minus
- ``()`` function call
- ``?.`` optional chaining
- ``.`` member access
- ``[]`` indexed value access

Operators higher in this list bind more tightly. For example, multiplication is evaluated before addition, logical and before logical or, and comparisons before logical operators.

{{< editor lang="fql" >}}
// Multiplication binds tighter than addition:
// interpreted as 2 + (3 * 4), not (2 + 3) * 4
return 2 + 3 * 4
{{</ editor >}}

{{< editor lang="fql" >}}
// AND binds tighter than OR:
// interpreted as false || (true && true)
return false || true && true
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

## Using parentheses

Parentheses ``(`` and ``)`` override the default evaluation order. Use them when the intended grouping differs from the precedence rules, or when the expression is complex enough that the precedence is not immediately obvious.

{{< editor lang="fql" >}}
return (2 + 3) * 4
{{</ editor >}}

{{< editor lang="fql" height="150px" >}}
let price = 120
let discount = 0.1
let tax = 0.2

// Without parentheses: discount * tax is evaluated first
// With parentheses: subtraction happens before multiplication
return price * (1 - discount) * (1 + tax)
{{</ editor >}}

## Next steps

{{< docs-related tiles="language-operators,language-operators-coalescing,language-expressions,language-operators-arithmetic" >}}
