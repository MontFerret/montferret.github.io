---
title: "Operator Precedence"
sidebarTitle: "Precedence"
weight: 60
draft: false
description: "The evaluation order of FQL operators, from lowest to highest precedence."
---

# Operator precedence

The operator precedence in FQL is similar as in other familiar languages (lowest precedence first):

- ``? :`` ternary operator
- ``||`` logical or
- ``&&`` logical and
- ``==``, ``!=`` equality and inequality
- ``IN`` in operator
- ``<``, ``<=``, ``>=``, ``>`` less than, less equal, greater equal, greater than
- ``+``, ``-`` addition, subtraction
- ``*``, ``/``, ``%`` multiplication, division, modulus
- ``!``, ``+``, ``-`` logical negation, unary plus, unary minus
- ``()`` function call
- ``.`` member access
- ``[]`` indexed value access

Operators higher in this list bind more tightly. For example, multiplication is evaluated before addition, logical AND before logical OR, and comparisons before logical operators.

{{< editor lang="fql" >}}
// Multiplication binds tighter than addition:
// interpreted as 2 + (3 * 4), not (2 + 3) * 4
RETURN 2 + 3 * 4
{{</ editor >}}

{{< editor lang="fql" >}}
// AND binds tighter than OR:
// interpreted as false || (true && true)
RETURN false || true && true
{{</ editor >}}

## Using parentheses

Parentheses ``(`` and ``)`` override the default evaluation order. Use them when the intended grouping differs from the precedence rules, or when the expression is complex enough that the precedence is not immediately obvious.

{{< editor lang="fql" >}}
RETURN (2 + 3) * 4
{{</ editor >}}

{{< editor lang="fql" height="150px" >}}
LET price = 120
LET discount = 0.1
LET tax = 0.2

// Without parentheses: discount * tax is evaluated first
// With parentheses: subtraction happens before multiplication
RETURN price * (1 - discount) * (1 + tax)
{{</ editor >}}
