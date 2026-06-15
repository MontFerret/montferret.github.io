---
title: "Arithmetic Operators"
sidebarTitle: "Arithmetic"
weight: 30
draft: false
description: "Addition, subtraction, multiplication, division, modulus, and type conversion rules for arithmetic."
---

# Arithmetic operators

Arithmetic operators perform an arithmetic operation on two numeric operands. The result of an arithmetic operation is again a numeric value.

FQL supports the following arithmetic operators:

- ``+`` addition
- ``-`` subtraction
- ``*`` multiplication
- ``/`` division
- ``%`` modulus

Unary plus and unary minus are supported as well:

{{< editor lang="fql" >}}
LET x = -5
LET y = 1
RETURN [-x, +y]
{{</ editor >}}

For exponentiation, there is a numeric function `POW()`. The syntax base ** exp is not supported.

Some example arithmetic operations:

{{< code lang="fql" >}}
1 + 1
33 - 99
12.4 * 4.5
13.0 / 0.1
23 % 7
-15
+9.99
{{</ code >}}

## Type conversion

Arithmetic operators accept operands of any type. The conversion rules depend on the operator.

### Addition

The ``+`` operator performs numeric addition when both operands are numbers. When either operand is a string, it performs string concatenation instead. All other types are converted to their string representation before concatenation.

{{< editor lang="fql" >}}
RETURN [
    1 + 2,
    1 + "99",
    1 + "a",
    1 + NONE,
    NONE + 1,
    3 + [ ],
    24 + [ 2 ],
    "hello" + " " + "world"
]
{{</ editor >}}

### Subtraction, multiplication, division, and modulus

The ``-``, ``*``, ``/``, and ``%`` operators always convert their operands to numbers. The conversion rules are:

- ``NONE`` is converted to ``0``.
- ``false`` is converted to ``0``, ``true`` is converted to ``1``.
- A valid numeric value remains unchanged.
- String values are converted to a number if they contain a valid numeric representation. Strings with non-numeric contents are converted to ``0``.
- An empty array is converted to ``0``. A non-empty array is converted by summing the numeric values of all its elements.
- Objects, binary, and custom types are converted to ``0``.

{{< editor lang="fql" >}}
RETURN [
    25 - NONE,
    17 - true,
    23 * { },
    5 * [ 7 ],
    10 - [ 2, 3 ],
    24 / "12",
    0 / 1
]
{{</ editor >}}
