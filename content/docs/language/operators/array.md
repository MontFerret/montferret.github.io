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
- `[* filter condition]` array expansion with filtering
- `[* return projection]` array expansion with projection
- `[** limit n]` array flattening with limit
- `[? condition]` array search operator

## Array expansion

Array expansion applies an expression to each item in an array and returns the produced values as a new array.

The syntax for it is the following:

{{< code lang="fql" >}}
array[*].expression
{{</ code >}}

The `[*]` operator is used to indicate that the array should be expanded. It can be used with any array expression, including variables, function calls, and nested arrays.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", active: true },
    { name: "Grace", active: false },
    { name: "Linus", active: true }
]

return users[*].name
{{</ editor >}}

If an expanded item does not contain the requested property, the produced value is `none`.

{{< editor lang="fql" >}}
let users = [
    { name: "Ada", email: "ada@example.com" },
    { name: "Grace" }
]

return users[*].email
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
let values = [[1, 2], [3, 4], [5]]

return values[**]
{{</ editor >}}

The `[***]` operator removes two levels of nesting:

{{< editor lang="fql" >}}
let values = [
    [[1, 2]],
    [[3, 4]],
    [[5]]
]

return values[***]
{{</ editor >}}

Flattening only affects arrays at the level where the operator is applied. It does not recursively flatten every nested array unless enough `*` characters are used.

For example:

{{< editor lang="fql" >}}
let values = [
    [1, [2, 3]],
    [4, [5]]
]

return values[**]
{{</ editor >}}

Only one level was removed. The nested arrays `[2, 3]` and `[5]` remain because they are one level deeper.

To remove two levels, use `[***]`:

{{< editor lang="fql" >}}
let values = [
    [1, [2, 3]],
    [4, [5]]
]

return values[***]
{{</ editor >}}

Flattening is most useful when a query produces an array for each input value, but the final result should be a single array.

The following query returns the friend names for each user:

{{< editor lang="fql" >}}
let users = [
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

return users[*].friends[*].name
{{</ editor >}}

The result is a nested array, with one array of friend names for each user.

To return a single array of friend names, apply `[**]` to the nested result:

{{< editor lang="fql" >}}
let users = [
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

return users[*].friends[*].name[**]
{{</ editor >}}

The flattening operator does not remove duplicate values. If the same value appears multiple times, it appears multiple times in the flattened result.

{{< editor lang="fql" >}}
let groups = [
    ["admin", "editor"],
    ["editor", "viewer"]
]

return groups[**]
{{</ editor >}}

The query returns:

Use `distinct` operator when duplicate values need to be removed.

## Inline expressions

Array expansion and array flattening operators can include inline expressions.

Inline expressions make it possible to filter array elements, limit the number of elements included in the result, and project each element into a different value.

The following operations are supported inside an array operator:

- `filter` selects which elements are included.
- `limit` restricts how many elements are included.
- `return` defines the value produced for each included element.

Inline expressions can be used with array expansion and flattening operators:

{{< code lang="fql" >}}
array[* filter condition limit count return expression]
array[** filter condition limit count return expression]
array[*** filter condition limit count return expression]
{{</ code >}}

When multiple inline operations are used together, they must appear in this order:

{{< code lang="fql" >}}
filter
limit
return
{{</ code >}}

Inside each expression, `.` refers to the element currently being processed. The condition can refer to `.` as well as variables from the surrounding scope.

Each operation can appear at most once. Sorting is not supported by this shorthand form yet.

### Inline filter

`filter` includes only the elements that satisfy a condition.

{{< editor lang="fql" >}}
let users = [
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

return users[* return {
    name: .name,
    olderFriends: .friends[* filter .age > 40].name
}]
{{</ editor >}}

The `filter` condition can refer to `.`, functions, operators, and variables from the outer scope.

{{< editor lang="fql" >}}
let minAge = 40

let users = [
    {
        name: "Ada",
        friends: [
            { name: "Grace", age: 41 },
            { name: "Linus", age: 31 },
            { name: "Barbara", age: 45 }
        ]
    }
]

return users[* return {
    name: .name,
    friends: .friends[* filter .age >= minAge].name
}]
{{</ editor >}}

Inside nested inline expressions, `.` always refers to the current element of the innermost array operator.

Inline expressions are a shorthand form and do not support local variable declarations. For transformations that require named intermediate values, multiple scopes, or more complex control over the current element, use a regular `for` loop.

### Inline limit

`limit` restricts the number of elements included in the array result.

{{< editor lang="fql" >}}
let users = [
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

return users[* return {
    name: .name,
    friends: .friends[* limit 1].name
}]
{{</ editor >}}

When `filter` and `limit` are used together, `filter` must appear first. The limit is then applied to the filtered elements.

{{< editor lang="fql" >}}
let values = [1, 2, 3, 4, 5, 6]

return values[* filter . > 2 limit 2]
{{</ editor >}}

`limit` also supports an offset form:

{{< code lang="fql" >}}
limit offset, count
{{</ code >}}

The first number specifies how many matching elements to skip. The second number specifies the maximum number of elements to include.

{{< editor lang="fql" >}}
let users = [
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

return users[* return {
    name: .name,
    friends: .friends[* limit 1, 2].name
}]
{{</ editor >}}

This form skips the first element and includes at most two elements after it.

### Inline projection

`return` defines the value produced for each element.

Without an inline `return`, the array operator returns the selected elements themselves. With an inline `return`, each selected element is replaced by the value produced by the return expression.

{{< editor lang="fql" >}}
let friends = [
    { name: "Grace", age: 41 },
    { name: "Linus", age: 31 },
    { name: "Barbara", age: 45 }
]

return friends[* return .name]
{{</ editor >}}

The projection can produce any FQL value, including an object, array, string, number, or computed expression.

{{< editor lang="fql" >}}
let friends = [
{ name: "Grace", age: 41 },
{ name: "Linus", age: 31 },
{ name: "Barbara", age: 45 }
]

return friends[* return {
    label: concat(.name, " is ", .age),
    adult: .age >= 18
}]
{{</ editor >}}

return can be combined with filter and limit.

{{< editor lang="fql" >}}
let friends = [
    { name: "Grace", age: 41 },
    { name: "Linus", age: 31 },
    { name: "Barbara", age: 45 },
    { name: "Donald", age: 39 }
]

return friends[*
    filter .age >= 40
    limit 2
    return concat(.name, " is ", .age)
]
{{</ editor >}}

When all three operations are used together, the array is processed in the same order as the syntax: elements are filtered first, the limit is applied next, and the projection is evaluated last.

## Question mark operator

The question mark operator checks whether an array contains elements that satisfy a condition.

Unlike inline `filter`, which returns matching elements, the question mark operator returns a boolean value. It is used when a query needs to test an array rather than return a transformed array.

{{< editor lang="fql" >}}
let values = [1, 2, 3, 4]

return values[? 2 filter . % 2 == 0]
{{</ editor >}}

The value after `?` is a quantifier. It defines how many elements must satisfy the condition.

The quantifier is optional. When it is omitted, `any` is used.

The following quantifiers are supported:

- An integer number, such as `2`, requires exactly that many matching elements.
- A range, such as `2..4`, requires the number of matching elements to be within the range.
- `none` requires no matching elements.
- `any` requires at least one matching element.
- `all` requires every element to match.
- `at least n` requires at least `n` matching elements.

The quantifier is followed by `filter` when a condition is specified.

{{< editor lang="fql" >}}
let values = [1, 2, 3, 4, 5, 6]

return {
    exactlyThree: values[? 3 filter . % 2 == 0],
    betweenTwoAndFour: values[? 2..4 filter . > 2],
    noneNegative: values[? none filter . < 0],
    anyEven: values[? any filter . % 2 == 0],
    allPositive: values[? all filter . > 0],
    atLeastTwoLarge: values[? at least 2 filter . > 4]
}
{{</ editor >}}

Inside the question mark operator, `.` refers to the array element being tested by the `filter` condition.

{{< editor lang="fql" >}}
let minAge = 40

let users = [
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

return users[* return {
    name: .name,
    hasOlderFriend: .friends[? any filter .age >= minAge]
}]
{{</ editor >}}

When the operator is used without a quantifier and without filter, it checks whether the value is a non-empty array.

{{< editor lang="fql" >}}
return {
    empty: [][?],
    nonEmpty: [1, 2, 3][?],
    notArray: none[?]
}
{{</ editor >}}

Conceptually, the question mark operator is equivalent to filtering an array and then checking the number of matching elements.

| Question mark expression | Equivalent length check |
| --- | -- |
| `array[? n filter condition]` | `length(array[* filter condition]) == n` |
| `array[? min..max filter condition]` | the number of matching elements is between min and max |
| `array[? none filter condition]` | `length(array[* filter condition]) == 0` |
| `array[? any filter condition]` | `length(array[* filter condition]) > 0` |
| `array[? all filter condition]` | `length(array[* filter condition]) == length(array)` |
| `array[? at least n filter condition]` | `length(array[* filter condition]) >= n` |
| `array[?]` | the value is an array with at least one element |

The question mark operator is especially useful for nested search, where a document, object, or result value contains arrays that need to be tested without expanding them into the final output.

## Array comparison operators

Array comparison operators evaluate a comparison against the elements of an array.

They specify how the results of those element-wise comparisons are combined. Depending on the operator, the expression returns true when at least one element matches, when every element matches, or when no elements match.

FQL provides three array comparison operators:

* `any`
* `all`
* `none`

The general form is:

- `<array> any <comparison value>`
- `<array> all <comparison value>`
- `<array> none <comparison value>`

The comparison operator can be any [supported comparison operator]({{< ref "comparison" >}}).

Each element uses the same strict comparison rules as a scalar expression. In particular, Duration values compare only with Duration values: incompatible equality simply does not match, while an incompatible relational comparison raises an error.

{{< code lang="fql" >}}
[1s, 2s] any == "2s"               // false
[1s, 2s] any == to_duration("2s")  // true
[1s, 2s] all > 999                  // runtime error
{{</ code >}}

### any

`any` returns true when at least one element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If one or more elements produce true, the whole `any` expression returns true. If no elements satisfy the comparison, the expression returns false.

{{< editor lang="fql" >}}
let tags = ["docs", "fql", "arrays"]

return tags any == "fql"
{{</ editor >}}

This expression returns true because the array contains "fql". The other elements do not need to match; a single matching element is enough.

`any` is useful when the presence of at least one matching value is enough to satisfy a condition.

{{< editor lang="fql" >}}
let prices = [12, 18, 31, 9]

return prices any >= 30
{{</ editor >}}


### all

`all` returns true only when every element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If all elements produce true, the whole `all` expression returns true. If at least one element does not satisfy the comparison, the expression returns false.

{{< editor lang="fql" >}}
let scores = [82, 91, 74, 88]

return scores all >= 70
{{</ editor >}}

In this example, the expression returns true because every value in scores is greater than or equal to 70.

{{< editor lang="fql" >}}
let scores = [82, 91, 64, 88]

return scores all >= 70
{{</ editor >}}

Use `all` when the expression should return true only if the comparison holds for the complete array.

### none

`none` returns true when no element in the array satisfies the comparison.

The comparison is evaluated for each element of the array. If every element produces false, the whole `none` expression returns true. If at least one element satisfies the comparison, the expression returns false.

{{< editor lang="fql" >}}
let blocked = ["spam", "phishing", "malware"]
let labels = ["docs", "release", "fql"]

return labels none in blocked
{{</ editor >}}

In this example, the expression returns true because none of the values in labels are present in blocked.

`none` is the inverse form of any: it requires the comparison to fail for the complete array. Use none when the expression should return true only if the array contains no matching values.

{{< editor lang="fql" >}}
let values = [3, 7, 11]

return values none == 0
{{</ editor >}}

### Comparing arrays with arrays

Array comparison operators can also be used when the right side of the comparison is another array.

{{< editor lang="fql" >}}
let userRoles = ["reader", "editor"]
let requiredRoles = ["admin", "owner"]

return userRoles none in requiredRoles
{{</ editor >}}

This returns true because none of the user roles are present in `requiredRoles`.

The comparison is still applied element by element. In this example, each value from `userRoles` is checked against `requiredRoles` using `in`.

Another common case is checking whether at least one value from one array exists in another array:

{{< editor lang="fql" >}}
let selected = ["fql", "arrays"]
let supported = ["values", "arrays", "operators"]

return selected any in supported
{{</ editor >}}

This returns true because `"arrays"` exists in both arrays.

To require every value from one array to exist in another array, use `all`:

{{< editor lang="fql" >}}
let selected = ["arrays", "operators"]
let supported = ["values", "arrays", "operators"]

return selected all in supported
{{</ editor >}}

### Empty arrays

Array comparison operators have well-defined behavior for empty arrays.

{{< editor lang="fql" >}}
let values = []

return {
    any: values any == 1,
    all: values all == 1,
    none: values none == 1
}
{{</ editor >}}

### Using array comparisons in filters

Array comparison operators are often used inside `filter` statements to keep or reject records based on array contents.

{{< editor lang="fql" >}}
let products = [
    { name: "Basic plan", features: ["docs", "search"] },
    { name: "Team plan", features: ["docs", "search", "sharing"] },
    { name: "Enterprise plan", features: ["docs", "search", "sharing", "sso"] }
]

return for product in products {
    filter product.features any == "sso"
    return product.name
}
{{</ editor >}}

The same idea can be used to exclude records:

{{< editor lang="fql" >}}
let products = [
    { name: "Basic plan", tags: ["public"] },
    { name: "Internal preview", tags: ["internal", "beta"] },
    { name: "Stable release", tags: ["public", "stable"] }
]

return for product in products {
    filter product.tags none == "internal"
    return product.name
}
{{</ editor >}}

### Difference from direct array comparison

Array comparison operators are different from comparing an array value directly.

{{< editor lang="fql" >}}
let values = [1, 2, 3]

return values == 2
{{</ editor >}}

This compares the array itself with the number `2`, so the result is `false`.

To compare the elements inside the array, use an array comparison operator:

{{< editor lang="fql" >}}
let values = [1, 2, 3]

return values any == 2
{{</ editor >}}

Use direct array comparison when the array as a whole is the value being compared. Use array comparison operators when the condition should be applied to the array's elements.

## Next steps

{{< docs-related tiles="language-operators,language-types-basic,language-operators-comparison" >}}
