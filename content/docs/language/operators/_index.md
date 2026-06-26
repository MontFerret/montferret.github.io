---
title: "Operators"
sidebarTitle: "Operators"
weight: 70
draft: false
description: "Comparison, logical, arithmetic, ternary, range, array, and precedence operators."
aliases:
  - /docs/fql/operators/
---

# Operators

FQL provides a set of operators for comparing values, combining conditions, performing arithmetic, working with arrays, and controlling evaluation order. Operators are used throughout the language in expressions, filters, projections, and conditions.

## Comparison

Comparison operators compare two operands and return a boolean result. They include equality (`==`, `!=`), ordering (`<`, `<=`, `>`, `>=`), containment (`IN`, `NOT IN`), pattern matching (`LIKE`, `NOT LIKE`), and regular expression matching (`=~`, `!~`).

{{< code lang="fql" >}}
65 == 65
"abc" != "ABC"
1 IN [1, 2, 3]
"foo" LIKE "f*"
{{</ code >}}

See [Comparison Operators]({{< ref "comparison" >}}).

## Logical

Logical operators evaluate expressions according to their truth value. FQL supports `&&` / `AND`, `||` / `OR`, and `!` / `NOT`. The binary operators use short-circuit evaluation and return one of their operands rather than always returning a boolean.

{{< code lang="fql" >}}
RETURN true && "value"
RETURN NONE || "fallback"
RETURN NOT false
{{</ code >}}

See [Logical Operators]({{< ref "logical" >}}).

## Arithmetic

Arithmetic operators perform operations on numeric operands: addition (`+`), subtraction (`-`), multiplication (`*`), division (`/`), and modulus (`%`). Unary plus and minus are also supported.

{{< code lang="fql" >}}
1 + 1
33 - 99
12.4 * 4.5
{{</ code >}}

See [Arithmetic Operators]({{< ref "arithmetic" >}}).

## Range

The range operator (`..`) produces an array of integer values between two bounds, inclusive.

{{< code lang="fql" >}}
2010..2013
{{</ code >}}

See [Range Operator]({{< ref "range" >}}).

## Ternary

The ternary operator provides conditional evaluation. It returns one of two values depending on a boolean condition.

{{< code lang="fql" >}}
u.age > 15 ? u.userId : NONE
{{</ code >}}

See [Ternary Operator]({{< ref "ternary" >}}).

## Array

Array operators work with arrays and nested array structures. They include indexed access (`[]`), expansion (`[*]`), flattening (`[**]`), inline filtering and projection, the question mark operator (`[?]`), and array comparison operators (`ANY`, `ALL`, `NONE`).

{{< code lang="fql" >}}
users[*].name
values[**]
values[* FILTER . > 2 LIMIT 3 RETURN . * 10]
tags ANY == "fql"
{{</ code >}}

See [Array Operators]({{< ref "array" >}}).

## DELETE

The `DELETE` statement removes a property from an object or host value.

{{< code lang="fql" >}}
DELETE target.property
DELETE target["property"]
{{</ code >}}

Both dot notation and bracket notation are supported:

{{< editor lang="fql" >}}
VAR user = { name: "Ada", deprecated: true, role: "admin" }

DELETE user.deprecated

RETURN user
{{</ editor >}}

Deletion removes the property entirely — it is not the same as assigning `NONE`, which keeps the key present with an absent value.

`DELETE` works with any value that supports the removable capability. Built-in objects support key removal. Host values may support removal if the host runtime provides that capability.

## Precedence

Operator precedence determines the order in which operators are evaluated. Parentheses can override the default evaluation order.

See [Operator Precedence]({{< ref "precedence" >}}).

## Where to go next

{{< docs-related tiles="language-operators-comparison,language-operators-logical,language-operators-arithmetic,language-operators-range,language-operators-ternary,language-operators-array,language-operators-precedence" >}}
