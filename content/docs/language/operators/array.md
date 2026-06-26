---
title: "Array Operators"
sidebarTitle: "Array"
weight: 70
draft: false
description: "Array expansion, flattening, inline expressions, the question mark operator, and array comparison operators."
---

# Array operators

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

## Array expansion

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

## Array flattening

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

## Inline expressions

Array expansion and array flattening operators can include inline expressions.

Inline expressions make it possible to filter array elements, limit the number of elements included in the result, and project each element into a different value.

The following operations are supported inside an array operator:

- `FILTER` selects which elements are included.
- `LIMIT` restricts how many elements are included.
- `RETURN` defines the value produced for each included element.

Inline expressions can be used with array expansion and flattening operators:

{{< code lang="fql" >}}
array[* FILTER condition LIMIT count RETURN expression]
array[** FILTER condition LIMIT count RETURN expression]
array[*** FILTER condition LIMIT count RETURN expression]
{{</ code >}}

When multiple inline operations are used together, they must appear in this order:

{{< code lang="fql" >}}
FILTER
LIMIT
RETURN
{{</ code >}}

Inside each expression, `.` refers to the element currently being processed. The condition can refer to `.` as well as variables from the surrounding scope.

Each operation can appear at most once. Sorting is not supported by this shorthand form yet.

### Inline filter

`FILTER` includes only the elements that satisfy a condition.

{{< editor lang="fql" >}}
LET users = [
    {
        name: "Ada",
        age: 37,
        friends: [
            { name: "Grace", age: 41 },
            { name: "Linus", age: 31 },
            { name: "Barbara", age: 45 }
        ]
    },
    {
        name: "Alan",
        age: 48,
        friends: [
            { name: "Edsger", age: 50 },
            { name: "Donald", age: 39 }
        ]
    }
]

RETURN users[* RETURN {
    name: .name,
    olderFriends: .friends[* FILTER .age > 40].name
}]
{{</ editor >}}

The `FILTER` condition can refer to `.`, functions, operators, and variables from the outer scope.

{{< editor lang="fql" >}}
LET minAge = 40

LET users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace", age: 41 },
            { name: "Linus", age: 31 },
            { name: "Barbara", age: 45 }
        ]
    }
]

RETURN users[* RETURN {
    name: .name,
    friends: .friends[* FILTER .age >= minAge].name
}]
{{</ editor >}}

Inside nested inline expressions, `.` always refers to the current element of the innermost array operator.

Inline expressions are a shorthand form and do not support local variable declarations. For transformations that require named intermediate values, multiple scopes, or more complex control over the current element, use a regular `FOR` loop.

### Inline limit

`LIMIT` restricts the number of elements included in the array result.

{{< editor lang="fql" >}}
LET users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace" },
            { name: "Linus" },
            { name: "Barbara" }
        ]
    },
    {
        name: "Alan",
        friends: [
            { name: "Edsger" },
            { name: "Donald" }
        ]
    }
]

RETURN users[* RETURN {
    name: .name,
    friends: .friends[* LIMIT 1].name
}]
{{</ editor >}}

When `FILTER` and `LIMIT` are used together, `FILTER` must appear first. The limit is then applied to the filtered elements.

{{< editor lang="fql" >}}
LET values = [1, 2, 3, 4, 5, 6]

RETURN values[* FILTER . > 2 LIMIT 2]
{{</ editor >}}

`LIMIT` also supports an offset form:

{{< code lang="fql" >}}
LIMIT offset, count
{{</ code >}}

The first number specifies how many matching elements to skip. The second number specifies the maximum number of elements to include.

{{< editor lang="fql" >}}
LET users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace" },
            { name: "Linus" },
            { name: "Barbara" }
        ]
    },
    {
        name: "Alan",
        friends: [
            { name: "Edsger" },
            { name: "Donald" },
            { name: "Frances" }
        ]
    }
]

RETURN users[* RETURN {
    name: .name,
    friends: .friends[* LIMIT 1, 2].name
}]
{{</ editor >}}

This form skips the first element and includes at most two elements after it.

### Inline projection

`RETURN` defines the value produced for each element.

Without an inline `RETURN`, the array operator returns the selected elements themselves. With an inline `RETURN`, each selected element is replaced by the value produced by the return expression.

{{< editor lang="fql" >}}
LET friends = [
    { name: "Grace", age: 41 },
    { name: "Linus", age: 31 },
    { name: "Barbara", age: 45 }
]

RETURN friends[* RETURN .name]
{{</ editor >}}

The projection can produce any FQL value, including an object, array, string, number, or computed expression.

{{< editor lang="fql" >}}
LET friends = [
{ name: "Grace", age: 41 },
{ name: "Linus", age: 31 },
{ name: "Barbara", age: 45 }
]

RETURN friends[* RETURN {
    label: CONCAT(.name, " is ", .age),
    adult: .age >= 18
}]
{{</ editor >}}

RETURN can be combined with FILTER and LIMIT.

{{< editor lang="fql" >}}
LET friends = [
    { name: "Grace", age: 41 },
    { name: "Linus", age: 31 },
    { name: "Barbara", age: 45 },
    { name: "Donald", age: 39 }
]

RETURN friends[*
    FILTER .age >= 40
    LIMIT 2
    RETURN CONCAT(.name, " is ", .age)
]
{{</ editor >}}

When all three operations are used together, the array is processed in the same order as the syntax: elements are filtered first, the limit is applied next, and the projection is evaluated last.

## Question mark operator

The question mark operator checks whether an array contains elements that satisfy a condition.

Unlike inline `FILTER`, which returns matching elements, the question mark operator returns a boolean value. It is used when a query needs to test an array rather than return a transformed array.

{{< editor lang="fql" >}}
LET values = [1, 2, 3, 4]

RETURN values[? 2 FILTER . % 2 == 0]
{{</ editor >}}

The value after `?` is a quantifier. It defines how many elements must satisfy the condition.

The quantifier is optional. When it is omitted, `ANY` is used.

The following quantifiers are supported:

- An integer number, such as `2`, requires exactly that many matching elements.
- A range, such as `2..4`, requires the number of matching elements to be within the range.
- `NONE` requires no matching elements.
- `ANY` requires at least one matching element.
- `ALL` requires every element to match.
- `AT LEAST n` requires at least `n` matching elements.

The quantifier is followed by `FILTER` when a condition is specified.

{{< editor lang="fql" >}}
LET values = [1, 2, 3, 4, 5, 6]

RETURN {
    exactlyThree: values[? 3 FILTER . % 2 == 0],
    betweenTwoAndFour: values[? 2..4 FILTER . > 2],
    noneNegative: values[? NONE FILTER . < 0],
    anyEven: values[? ANY FILTER . % 2 == 0],
    allPositive: values[? ALL FILTER . > 0],
    atLeastTwoLarge: values[? AT LEAST 2 FILTER . > 4]
}
{{</ editor >}}

Inside the question mark operator, `.` refers to the array element being tested by the `FILTER` condition.

{{< editor lang="fql" >}}
LET minAge = 40

LET users = [
{
    name: "Ada",
    friends: [
            { name: "Grace", age: 41 },
            { name: "Linus", age: 31 }
        ]
    },
    {
    name: "Alan",
    friends: [
            { name: "Edsger", age: 50 },
            { name: "Donald", age: 39 }
        ]
    }
]

RETURN users[* RETURN {
    name: .name,
    hasOlderFriend: .friends[? ANY FILTER .age >= minAge]
}]
{{</ editor >}}

When the operator is used without a quantifier and without FILTER, it checks whether the value is a non-empty array.

{{< editor lang="fql" >}}
RETURN {
    empty: [][?],
    nonEmpty: [1, 2, 3][?],
    notArray: NONE[?]
}
{{</ editor >}}

Conceptually, the question mark operator is equivalent to filtering an array and then checking the number of matching elements.

| Question mark expression | Equivalent length check |
| --- | -- |
| `array[? n FILTER condition]` | `LENGTH(array[* FILTER condition]) == n` |
| `array[? min..max FILTER condition]` | the number of matching elements is between min and max |
| `array[? NONE FILTER condition]` | `LENGTH(array[* FILTER condition]) == 0` |
| `array[? ANY FILTER condition]` | `LENGTH(array[* FILTER condition]) > 0` |
| `array[? ALL FILTER condition]` | `LENGTH(array[* FILTER condition]) == LENGTH(array)` |
| `array[? AT LEAST n FILTER condition]` | `LENGTH(array[* FILTER condition]) >= n` |
| `array[?]` | the value is an array with at least one element |

The question mark operator is especially useful for nested search, where a document, object, or result value contains arrays that need to be tested without expanding them into the final output.

## Array comparison operators

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

The comparison operator can be any [supported comparison operator]({{< ref "comparison" >}}).

### ANY

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


### ALL

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

### NONE

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

### Comparing arrays with arrays

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

### Empty arrays

Array comparison operators have well-defined behavior for empty arrays.

{{< editor lang="fql" >}}
LET values = []

RETURN {
    any: values ANY == 1,
    all: values ALL == 1,
    none: values NONE == 1
}
{{</ editor >}}

### Using array comparisons in filters

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

### Difference from direct array comparison

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

Use direct array comparison when the array as a whole is the value being compared. Use array comparison operators when the condition should be applied to the array's elements.

## Next steps

{{< docs-related tiles="language-operators,language-types-basic,language-operators-comparison" >}}
