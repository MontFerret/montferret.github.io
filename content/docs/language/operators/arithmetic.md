---
title: "Arithmetic Operators"
sidebarTitle: "Arithmetic"
weight: 30
draft: false
description: "Addition, subtraction, multiplication, division, modulus, and type conversion rules for arithmetic."
---

# Arithmetic operators

Arithmetic operators work with numbers and also provide checked Duration and DateTime operations. The result type depends on the operand pair.

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

The ``+`` operator performs numeric addition when both operands are numbers. When either operand is a string, it normally performs string concatenation instead. Temporal operands take precedence: a Duration or DateTime uses the temporal rules below before string concatenation is considered.

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

## Temporal arithmetic

Duration arithmetic converts operands only where the table below says “coercible Duration.” Numeric Duration inputs are milliseconds, duration strings may be compound, and fractional nanoseconds are truncated toward zero. See [Basic Types]({{< ref "/docs/language/types/basic#durations" >}}) for the complete conversion table.

| Expression | Result |
| --- | --- |
| `Duration + coercible Duration` | Duration |
| `coercible Duration + Duration` | Duration |
| `Duration - coercible Duration` | Duration |
| `Duration * number or numeric string` | scaled Duration |
| `number or numeric string * Duration` | scaled Duration |
| `Duration / number` or numeric string | scaled Duration |
| `Duration / Duration` or duration string | Int for an exact ratio, otherwise Float |
| `DateTime + coercible Duration` | DateTime |
| `coercible Duration + DateTime` | DateTime |
| `DateTime - coercible Duration` | DateTime |
| `DateTime - DateTime` | elapsed Duration |
| `DateTime - RFC3339 string` | elapsed Duration |

{{< editor lang="fql" >}}
LET start = TO_DATETIME("2024-03-10T06:30:00Z")

RETURN {
    numericMilliseconds: 1s + 250,
    compoundString: 30m + "1h30m",
    multiplied: "2.5" * 5s,
    scaled: 5s / "2",
    ratio: 5s / "250ms",
    later: start + "90m"
}
{{</ editor >}}

For `DateTime - String`, Ferret first tries to parse the string as an RFC3339 DateTime, then as a Duration. If neither conversion succeeds, the expression raises a runtime error.

The following pairs are unsupported and raise an invalid-operation error:

- `DateTime + DateTime`
- `Duration - DateTime`
- DateTime multiplication, division, or modulus
- Duration multiplication by another Duration or a duration-form string
- reverse division such as `2 / 1s`
- Duration modulus

All Duration and DateTime arithmetic is checked for overflow. DateTime arithmetic uses canonical instants and does not retain monotonic clock metadata.

## Next steps

{{< docs-related tiles="language-operators,language-operators-range,language-operators-precedence" >}}
