---
title: "Operators"
weight: 60
draft: false
description: "Comparison, logical, arithmetic, ternary, range, and precedence operators."
aliases:
  - /docs/fql/operators/
---

# Operators
## Comparison operators

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

Regular expression comparisons use the `=~` and `!~` operators.

The `=~` operator returns true when the string on the left side matches the regular expression on the right side.

{{< code lang="fql" >}}
"foo" =~ "^f[o].$"
"foo" !~ "[a-z]+bar$"
{{</ code >}}

Pattern and regular expression operators are intended for string values. They do not perform implicit conversion from other value types to strings.

## Logical operators

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

### Boolean conversion

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

### Logical AND

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

### Logical OR

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

### Logical NOT

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

### Short-circuit evaluation

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

### Result values

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

## Arithmetic operators

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

The arithmetic operators accept operands of any type. Passing non-numeric values to an arithmetic operator will cast the operands to numbers:

- ``NONE`` will be converted to ``0``
- ``false`` will be converted to ``0``, ``true`` will be converted to ``1``
- a valid numeric value remains unchanged, but ``NaN`` and Infinity will be converted to ``0``
- string values are converted to a number if they contain a valid string representation of a number. Any whitespace at the start or the end of the string is ignored. Strings with any other contents are converted to the number ``0``
- an empty array is converted to ``0``, an array with one member is converted to the numeric representation of its sole member. Arrays with more members are converted to the number ``0``.
- objects, binary and custom types are converted to the number ``0``.

Here are a few examples:

{{< editor lang="fql" >}}
RETURN [
    1 + "a",
    1 + "99",
    1 + NONE,
    NONE + 1,
    3 + [ ],
    24 + [ 2 ],
    24 + [ 2, 4 ],
    25 - NONE,
    17 - true,
    23 * { },
    5 * [ 7 ],
    24 / "12",
    0 / 1
]
{{</ editor >}}

## Ternary operator
FQL also supports a ternary operator that can be used for conditional evaluation. The ternary operator expects a boolean condition as its first operand, and it returns the result of the second operand if the condition evaluates to true, and the third operand otherwise.

{{< code lang="fql" >}}
u.age > 15 || u.active == true ? u.userId : null
{{</ code >}}

There is also a shortcut variant of the ternary operator with just two operands. This variant can be used when the expression for the boolean condition and the return value should be the same:

{{< code lang="fql" >}}
u.value ? : 'value is null, 0 or not present'
{{</ code >}}

## Range operator
FQL supports expressing simple numeric ranges with the ``..`` operator. This operator can be used to easily iterate over a sequence of numeric values.

The ``..`` operator will produce an array of the integer values in the defined range, with both bounding values included.

{{< editor lang="fql" >}}
RETURN 2010..2013
{{</ editor >}}

## Operator precedence

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

The parentheses ``(`` and ``)`` can be used to enforce a different operator evaluation order.


## Array operators
FQL provides several operators for working with arrays and nested array structures.

Array operators can access individual items, expand arrays into derived values, flatten nested arrays, or search through nested array structures.

The most common forms are:
- `[]` indexed value access operator
- `[*]` array expansion operator
- `[**]` array flattening operator
- `[* FILTER condition]` array expansion with filtering
- `[* RETURN projection]` array expansion with projection
- `[** LIMIT n]` array flattening with limit
- `[? condition]` array search operator

### Array expansion

Array expansion applies an expression to each item in an array and returns the produced values as a new array.

The syntax for it is the following:

{{< code lang="fql" >}}
array[*].expression
{{</ code >}}

The `[*]` operator is used to indicate that the array should be expanded. It can be used with any array expression, including variables, function calls, and nested arrays.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

RETURN users[*].name
{{</ editor >}}

If an expanded item does not contain the requested property, the produced value is `NONE`.

{{< editor lang="fql" >}}
LET users = [
    { name: "Ada", email: "ada@example.com" },
    { name: "Grace" }
]

RETURN users[*].email
{{</ editor >}}

Array expansion preserves the number and order of items in the input array and produces a new array with the same number of items.

### Array flattening

The array flattening operator expands an array and removes one or more levels of nested arrays from the result.

The syntax is based on the array expansion operator, but uses multiple asterisks:

- `[*]` expands an array.
- `[**]` expands an array and flattens one nested level.
- `[***]` expands an array and flattens two nested levels.
- Additional `*` characters flatten additional levels.

In other words, the number of extra `*` characters determines how many array nesting levels are removed.

For example, `[**]` removes one level of nesting:

{{< editor lang="fql" >}}
LET values = [[1, 2], [3, 4], [5]]

RETURN values[**]
{{</ editor >}}

The `[***]` operator removes two levels of nesting:

{{< editor lang="fql" >}}
LET values = [
    [[1, 2]],
    [[3, 4]],
    [[5]]
]

RETURN values[***]
{{</ editor >}}

Flattening only affects arrays at the level where the operator is applied. It does not recursively flatten every nested array unless enough `*` characters are used.

For example:

{{< editor lang="fql" >}}
LET values = [
    [1, [2, 3]],
    [4, [5]]
]

RETURN values[**]
{{</ editor >}}

Only one level was removed. The nested arrays `[2, 3]` and `[5]` remain because they are one level deeper.

To remove two levels, use `[***]`:

{{< editor lang="fql" >}}
LET values = [
    [1, [2, 3]],
    [4, [5]]
]

RETURN values[***]
{{</ editor >}}

Flattening is most useful when a query produces an array for each input value, but the final result should be a single array.

The following query returns the friend names for each user:

{{< editor lang="fql" >}}
LET users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace" },
            { name: "Linus" }
        ]
    },
    {
        name: "Alan",
        friends: [
            { name: "Edsger" },
            { name: "Barbara" }
        ]
    }
]

RETURN users[*].friends[*].name
{{</ editor >}}

The result is a nested array, with one array of friend names for each user.

To return a single array of friend names, apply `[**]` to the nested result:

{{< editor lang="fql" >}}
LET users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace" },
            { name: "Linus" }
        ]
    },
    {
        name: "Alan",
        friends: [
            { name: "Edsger" },
            { name: "Barbara" }
        ]
    }
]

RETURN users[*].friends[*].name[**]
{{</ editor >}}

The flattening operator does not remove duplicate values. If the same value appears multiple times, it appears multiple times in the flattened result.

{{< editor lang="fql" >}}
LET groups = [
    ["admin", "editor"],
    ["editor", "viewer"]
]

RETURN groups[**]
{{</ editor >}}

The query returns:

Use `DISTINCT` operator when duplicate values need to be removed.

### Array comparison operators
Array comparison operators evaluate a comparison against the elements of an array.

They specify how the results of those element-wise comparisons are combined. Depending on the operator, the expression returns true when at least one element matches, when every element matches, or when no elements match.

FQL provides three array comparison operators:

* `ANY`
* `ALL`
* `NONE`

The general form is:

- `<array> ANY <comparison value>`
- `<array> ALL <comparison value>`
- `<array> NONE <comparison value>`

The comparison operator can be any [supported comparison operator]({{% ref "#comparison-operators" %}}).

#### ANY

`ANY` returns true when at least one element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If one or more elements produce true, the whole `ANY` expression returns true. If no elements satisfy the comparison, the expression returns false.

{{< editor lang="fql" >}}
LET tags = ["docs", "fql", "arrays"]

RETURN tags ANY == "fql"
{{</ editor >}}

This expression returns true because the array contains "fql". The other elements do not need to match; a single matching element is enough.

`ANY` is useful when the presence of at least one matching value is enough to satisfy a condition.

{{< editor lang="fql" >}}
LET prices = [12, 18, 31, 9]

RETURN prices ANY >= 30
{{</ editor >}}


#### ALL

`ALL` returns true only when every element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If all elements produce true, the whole `ALL` expression returns true. If at least one element does not satisfy the comparison, the expression returns false.

{{< editor lang="fql" >}}
LET scores = [82, 91, 74, 88]

RETURN scores ALL >= 70
{{</ editor >}}

In this example, the expression returns true because every value in scores is greater than or equal to 70.

{{< editor lang="fql" >}}
LET scores = [82, 91, 64, 88]

RETURN scores ALL >= 70
{{</ editor >}}

Use `ALL` when the expression should return true only if the comparison holds for the complete array.

#### NONE

`NONE` returns true when no element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If every element produces false, the whole `NONE` expression returns true. If at least one element satisfies the comparison, the expression returns false.

{{< editor lang="fql" >}}
LET blocked = ["spam", "phishing", "malware"]
LET labels = ["docs", "release", "fql"]

RETURN labels NONE IN blocked
{{</ editor >}}

In this example, the expression returns true because none of the values in labels are present in blocked.

`NONE` is the inverse form of ANY: it requires the comparison to fail for the complete array. Use NONE when the expression should return true only if the array contains no matching values.

{{< editor lang="fql" >}}
LET values = [3, 7, 11]

RETURN values NONE == 0
{{</ editor >}}

#### Comparing arrays with arrays

Array comparison operators can also be used when the right side of the comparison is another array.

{{< editor lang="fql" >}}
LET userRoles = ["reader", "editor"]
LET requiredRoles = ["admin", "owner"]

RETURN userRoles NONE IN requiredRoles
{{</ editor >}}

This returns true because none of the user roles are present in `requiredRoles`.

The comparison is still applied element by element. In this example, each value from `userRoles` is checked against `requiredRoles` using `IN`.

Another common case is checking whether at least one value from one array exists in another array:

{{< editor lang="fql" >}}
LET selected = ["fql", "arrays"]
LET supported = ["values", "arrays", "operators"]

RETURN selected ANY IN supported
{{</ editor >}}

This returns true because `"arrays"` exists in both arrays.

To require every value from one array to exist in another array, use `ALL`:

{{< editor lang="fql" >}}
LET selected = ["arrays", "operators"]
LET supported = ["values", "arrays", "operators"]

RETURN selected ALL IN supported
{{</ editor >}}

#### Empty arrays

Array comparison operators have well-defined behavior for empty arrays.

{{< editor lang="fql" >}}
LET values = []

RETURN {
    any: values ANY == 1,
    all: values ALL == 1,
    none: values NONE == 1
}
{{</ editor >}}

#### Using array comparisons in filters

Array comparison operators are often used inside `FILTER` statements to keep or reject records based on array contents.

{{< editor lang="fql" >}}
LET products = [
    { name: "Basic plan", features: ["docs", "search"] },
    { name: "Team plan", features: ["docs", "search", "sharing"] },
    { name: "Enterprise plan", features: ["docs", "search", "sharing", "sso"] }
]

FOR product IN products
    FILTER product.features ANY == "sso"
    RETURN product.name
{{</ editor >}}

The same idea can be used to exclude records:

{{< editor lang="fql" >}}
LET products = [
    { name: "Basic plan", tags: ["public"] },
    { name: "Internal preview", tags: ["internal", "beta"] },
    { name: "Stable release", tags: ["public", "stable"] }
]

FOR product IN products
    FILTER product.tags NONE == "internal"
    RETURN product.name
{{</ editor >}}

#### Difference from direct array comparison

Array comparison operators are different from comparing an array value directly.

{{< editor lang="fql" >}}
LET values = [1, 2, 3]

RETURN values == 2
{{</ editor >}}

This compares the array itself with the number `2`, so the result is `false`.

To compare the elements inside the array, use an array comparison operator:

{{< editor lang="fql" >}}
LET values = [1, 2, 3]

RETURN values ANY == 2
{{</ editor >}}

Use direct array comparison when the array as a whole is the value being compared. Use array comparison operators when the condition should be applied to the array’s elements.