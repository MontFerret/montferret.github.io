---
title: "Logical Operators"
sidebarTitle: "Logical"
weight: 20
draft: false
description: "Logical AND, OR, NOT operators, short-circuit evaluation, boolean conversion, and result values."
---

# Logical operators

Logical operators evaluate expressions according to their truth value. They are commonly used in filter conditions, conditional expressions, and any other context where a query needs to combine or negate conditions.

FQL supports the following symbolic logical operators:

- `&&` logical AND
- `||` logical OR
- `!` logical NOT

FQL also supports keyword-based forms of the same operators:

- `AND` logical AND
- `OR` logical OR
- `NOT` logical NOT

The keyword forms are aliases for the symbolic forms. They have the same behavior and may be used interchangeably. For example, `a && b` and `a AND b` are equivalent, as are `!value` and `NOT value`.

{{< code lang="fql" >}}
RETURN true && false
RETURN true OR false
RETURN NOT false
{{</ code >}}

## Boolean conversion

Logical operators may be applied to values that are not booleans. When a value is used in a logical operation, FQL evaluates the value according to its boolean representation.

The conversion rules are:

- `NONE` evaluates as `false`.
- Boolean values keep their original value.
- Numeric values evaluate as false when they are `0`; all other numeric values evaluate as `true`.
- Strings evaluate as `false` when they are empty; all non-empty strings evaluate as `true`.
- Arrays evaluate as `true`, regardless of whether they contain any elements.
- Objects evaluate as `true`, regardless of whether they contain any properties.
- Binary values and custom runtime-backed values evaluate as `true`.

These conversions are applied implicitly by logical operators. Passing a non-boolean value to a logical operator does not cause the query to fail.

{{< code lang="fql" >}}
RETURN !!0
RETURN !!1
RETURN !!""
RETURN !!"ferret"
RETURN !![]
RETURN !!{}
{{</ code >}}

## Logical AND

The logical `AND` operator evaluates the left-hand operand first.

If the left-hand operand evaluates as false, the operation returns the original left-hand value. In this case, the right-hand operand is not evaluated unless the query plan requires it to be evaluated earlier.

If the left-hand operand evaluates as true, the operation returns the original right-hand value.

This means that `&&` and `AND` do not always return a boolean value. They return one of their operands.

{{< code lang="fql" >}}
RETURN false && "value"
RETURN NONE && true
RETURN 0 && "fallback"
RETURN true && 23
RETURN "user" && "active"
{{</ code >}}

In the examples above, the result is the first operand that prevents the `AND` expression from continuing, or the right-hand operand when the left-hand operand is truthy.

This behavior is useful when logical expressions are used to guard access to values or to return a value only when a preceding condition is satisfied.

{{< editor lang="fql" height="150px" >}}
LET user = {
    active: true,
    name: "Ada"
}

RETURN user.active && user.name
{{</ editor >}}

The query returns "Ada" because user.active evaluates as true, so the result of the `AND` expression is the right-hand operand.

## Logical OR

The logical `OR` operator also evaluates the left-hand operand first.

If the left-hand operand evaluates as true, the operation returns the original left-hand value. In this case, the right-hand operand is not evaluated unless the query plan requires it to be evaluated earlier.

If the left-hand operand evaluates as false, the operation returns the original right-hand value.

Like logical `AND`, logical `OR` does not necessarily return a boolean value. It returns one of its operands.

{{< code lang="fql" >}}
RETURN true || "value"
RETURN 1 || 7
RETURN "ferret" OR "fallback"
RETURN NONE || "fallback"
RETURN "" || "fallback"
{{</ code >}}

In the examples above, the `OR` expression returns the first truthy operand, or the right-hand operand when the left-hand operand is falsy.

This behavior is commonly used to provide fallback values.

{{< editor lang="fql" height="150px" >}}
LET user = {
    displayName: ""
}

RETURN user.displayName || "Anonymous"
{{</ editor >}}

The query returns "Anonymous" because user.displayName is an empty string and therefore evaluates as false.

## Logical NOT

The logical `NOT` operator converts its operand to a boolean value and returns the negated boolean result.

Unlike `&&` and `||`, the `NOT` operator always returns a boolean.

{{< code lang="fql" >}}
RETURN !true
RETURN !false
RETURN !NONE
RETURN !0
RETURN !""
RETURN !"ferret"
{{</ code >}}

The keyword form `NOT` has the same behavior:

{{< code lang="fql" >}}
RETURN NOT true
RETURN NOT NONE
RETURN NOT "ferret"
{{</ code >}}

## Short-circuit evaluation

The binary logical operators use short-circuit evaluation.

The left-hand operand is evaluated first. The right-hand operand is evaluated only when it is needed to determine the result of the expression.

For logical `AND`, the right-hand operand is evaluated only if the left-hand operand evaluates as true. If the left-hand operand evaluates as false, the expression returns the left-hand value immediately.

For logical `OR`, the right-hand operand is evaluated only if the left-hand operand evaluates as false. If the left-hand operand evaluates as true, the expression returns the left-hand value immediately.

{{< editor lang="fql" >}}
LET user = {
    active: false,
    name: "Ada"
}

RETURN user.active && user.name
{{</ editor >}}

In this example, user.name does not need to determine the result of the expression because user.active already evaluates as false.

Similarly, the right-hand side of an OR expression is skipped when the left-hand side already evaluates as true:

{{< editor lang="fql" >}}
LET user = {
    active: false,
    name: "Ada"
}

RETURN user.name || "Anonymous"
{{</ editor >}}

The expression returns "Ada" without needing the fallback value to determine the result.

Subqueries are an exception to the normal short-circuit model. When an operand contains a subquery, the subquery may be evaluated before the logical operator itself as part of query execution. For this reason, logical short-circuiting should not be used to assume that a subquery operand is never evaluated.

## Result values

The result type of a logical expression depends on the operator.

The ! and NOT operators always return a boolean value.

The `&&`, `AND`, `||`, and `OR` operators return one of their operands. As a result, the returned value may have any FQL type.

Expressions that use comparison operators typically produce boolean results:

{{< code lang="fql">}}
RETURN 25 > 1 && 42 != 7
RETURN 22 IN [23, 42] || 23 NOT IN [22, 7]
RETURN 25 != 25
{{</ code >}}

However, logical operators can also return non-boolean values:

{{< code lang="fql" >}}
RETURN 1 || 7
RETURN NONE || "foo"
RETURN NONE && true
RETURN true && 23
{{</ code >}}

These expressions return `1`, `"foo"`, `NONE`, and `23`, respectively.

When a strict boolean result is required, convert the final value to a boolean by applying logical NOT twice:

{{< code lang="fql" >}}
RETURN !!"ferret"
RETURN !!NONE
{{</ code >}}

## Next steps

{{< docs-related tiles="language-operators,language-operators-comparison,language-control-flow-match" >}}
