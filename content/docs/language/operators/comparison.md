---
title: "Comparison Operators"
sidebarTitle: "Comparison"
weight: 10
draft: false
description: "Equality, ordering, containment, pattern matching, and regular expression comparison operators."
---

# Comparison operators

Comparison operators compare two operands and return a boolean result.

FQL supports comparison between values of any type. The result of a comparison is always either true or false. Equality and inequality comparisons can be used to determine whether two values are the same or different. Ordering comparisons can be used to determine whether one value sorts before or after another value according to FQL type ordering rules.

The following comparison operators are available:

- `==` - equality
- `!=` - inequality
- `<` - less than
- `<=` - less than or equal to
- `>` - greater than
- `>=` - greater than or equal to
- `IN` - tests whether a value is contained in an array
- `NOT IN` - tests whether a value is not contained in an array
- `LIKE` - tests whether a string matches a wildcard pattern
- `NOT LIKE` - tests whether a string does not match a wildcard pattern
- `=~` - tests whether a string matches a regular expression
- `!~` - tests whether a string does not match a regular expression

Each comparison operator evaluates its operands and returns true when the comparison condition is satisfied. Otherwise, it returns false.

FQL does not implicitly convert values before comparing them. For example, the number `65` and the string `"65"` are different values. They are not converted to a common type before the comparison is evaluated.

This behavior is important when comparing values of different types. Equality comparisons only return true when the compared values are equal according to FQL value semantics. Ordering comparisons between different types are evaluated according to FQL type ordering, not by converting one operand into the type of the other operand.

{{< code lang="fql" >}}
0 == NONE
1 > 0
true != NONE
45 <= "yikes!"
65 != "65"
65 == 65
1.23 > 1.32
"abc" == "abc"
"abc" == "ABC"
{{</ code >}}

## Containment

The `IN` operator tests whether the value on the left side appears in the array on the right side.

{{< code lang="fql" >}}
1 IN [2, 3, 1.5]
"foo" IN NONE
{{</ code >}}

The `NOT IN` operator is the negated form of `IN`. It returns true when the left-hand value is not contained in the right-hand array.

{{< code lang="fql" >}}
1 NOT IN [2, 3, 1.5]
"foo" NOT IN NONE
{{</ code >}}

`IN` and `NOT IN` are only meaningful when the right-hand operand is an array. If the right-hand operand is not an array, the value cannot be found in it, and the result is false for `IN`.

## Pattern matching

The `LIKE` operator compares a string value against a wildcard pattern.

{{< code lang="fql" >}}
"foo" LIKE "f*"
"abc" LIKE "a*"
"abc" LIKE "?bc"
{{</ code >}}

The pattern is specified by the right-hand operand. The left-hand operand is the value being tested. Pattern matching performed by `LIKE` is case-sensitive.

The `NOT LIKE` operator has the same matching rules as `LIKE`, but returns the opposite result.

{{< code lang="fql" >}}
"foo" NOT LIKE "f*"
"abc" NOT LIKE "a*"
"abc" NOT LIKE "?bc"
{{</ code >}}

## Regular expressions

Regular expression comparisons use the `=~` and `!~` operators.

The `=~` operator returns true when the string on the left side matches the regular expression on the right side.

{{< code lang="fql" >}}
"foo" =~ "^f[o].$"
"foo" !~ "[a-z]+bar$"
{{</ code >}}

Pattern and regular expression operators are intended for string values. They do not perform implicit conversion from other value types to strings.

## Next steps

{{< docs-related tiles="language-operators,language-types-ordering,language-control-flow-match" >}}
