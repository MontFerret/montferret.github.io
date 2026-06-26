---
title: "Match Expressions"
sidebarTitle: "Match"
weight: 10
draft: false
description: "Select a value by matching against patterns or conditions with the MATCH expression."
---

# Match Expressions

A `MATCH` expression chooses a result from several branches. Each branch, called an **arm**, pairs a pattern or a condition with a result expression after `=>`. `MATCH` tries the arms from top to bottom, stops at the first one that applies, and produces its result.

`MATCH` is an expression: it produces a value, so it can be returned, assigned, or nested inside another expression.

{{< editor lang="fql" >}}
LET status = "active"

RETURN MATCH status (
    "active"   => "online",
    "inactive" => "offline",
    _          => "unknown",
)
{{</ editor >}}

A `MATCH` has two forms. The first matches a value against patterns; the second evaluates a list of conditions. Both forms end with a required default arm, written `_ => ...`.

## Matching a value

When `MATCH` is followed by an expression, that value — the *subject* — is compared against each arm's pattern.

{{< editor lang="fql" >}}
RETURN MATCH 2 (
    1 => "one",
    2 => "two",
    _ => "other",
)
{{</ editor >}}

The subject is evaluated once. Arms are tried in order, and only the matching arm's result is evaluated.

### Patterns

An arm's pattern can be a literal, a binding, or an object pattern.

A **literal pattern** matches a constant value — a number, string, boolean, or `NONE`.

{{< editor lang="fql" >}}
RETURN MATCH NONE (
    NONE => "missing",
    _    => "present",
)
{{</ editor >}}

A **binding pattern** is an identifier that matches any value and binds it to a name for use in that arm's result. It is most useful with a guard (see below), so the arm applies only when an extra condition holds.

{{< editor lang="fql" >}}
RETURN MATCH 5 (
    v WHEN v > 3 => v * 2,
    _            => 0,
)
{{</ editor >}}

An **object pattern** matches the shape of an object. Listed fields must be present and must themselves match; matching fields can be bound to names.

{{< editor lang="fql" >}}
RETURN MATCH { a: 1, b: 2 } (
    { a: 1, b: v } => v,
    _              => 0,
)
{{</ editor >}}

If a required field is missing, the pattern does not match and the next arm is tried.

{{< editor lang="fql" >}}
RETURN MATCH { a: 1 } (
    { a: 1, b: v } => v,
    _              => 0,
)
{{</ editor >}}

### Guards

A `WHEN` clause after a pattern adds a condition. The arm applies only when the pattern matches **and** the guard is true. A guard can use values bound by the pattern.

{{< editor lang="fql" >}}
LET n = 7

RETURN MATCH n (
    x WHEN x % 2 == 0 => "even",
    x WHEN x % 2 == 1 => "odd",
    _                 => "other",
)
{{</ editor >}}

### The default arm

Every `MATCH` ends with a default arm, written `_ => ...`. It is required, and it runs when no earlier arm applies.

## Matching conditions

The second form omits the subject. Each arm is a `WHEN` condition, and the first condition that is true wins.

{{< editor lang="fql" >}}
RETURN MATCH (
    WHEN 1 < 2 => "first",
    WHEN 2 < 3 => "second",
    _          => "default",
)
{{</ editor >}}

Even though both conditions above are true, the result is `"first"` because arms are evaluated top to bottom and matching stops at the first success.

## Using match expressions

Because `MATCH` produces a value, it can appear anywhere an expression is expected — in a `LET` declaration, as a function argument, or in the body of a loop.

{{< editor lang="fql" >}}
FOR v IN [1, 2, 3]
    RETURN MATCH v (
        1 => "one",
        2 => "two",
        _ => "other",
    )
{{</ editor >}}

## Match and the ternary operator

For a simple two-way choice on a boolean condition, the [ternary operator]({{% ref "../operators/ternary" %}}) (`? :`) is more concise. Reach for `MATCH` when you have more than two branches, when you want to match against patterns, or when you want a chain of conditions.

## Next steps

{{< docs-related tiles="language-control-flow-error-handling,language-operators-comparison,language-control-flow" >}}
