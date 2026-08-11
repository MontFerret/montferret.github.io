---
title: "NONE-Coalescing Operator"
sidebarTitle: "NONE Coalescing"
weight: 25
draft: false
description: "Select a fallback only when a value is NONE with the right-associative coalescing operator."
---

# NONE-coalescing operator

The NONE-coalescing operator (`??`) selects a fallback only when its left operand evaluates to `NONE`.

{{< code lang="fql" >}}
LET user = {
    profile: {}
}

RETURN user?.profile?.displayName ?? "Anonymous"
{{</ code >}}

If `displayName` is missing, optional chaining produces `NONE` and the expression returns `"Anonymous"`.

## How it evaluates

For `left ?? right`, FQL evaluates `left` once.

- If `left` is not `NONE`, the expression returns it without evaluating `right`.
- If `left` is `NONE`, FQL evaluates and returns `right`.

The two operands may have different types. The selected operand becomes the result without conversion.

## Only NONE selects the fallback

`NULL` represents the same absent runtime value as `NONE`, so both select the fallback. Other values remain present even when a logical operator would treat them as false.

| Left operand | Result of `left ?? "fallback"` |
| --- | --- |
| `NONE` | `"fallback"` |
| `NULL` | `"fallback"` |
| `false` | `false` |
| `0` | `0` |
| `""` | `""` |
| `[]` | `[]` |
| `{}` | `{}` |

This makes `??` useful when `false`, zero, or an empty value carries meaning and must not be replaced.

## Coalescing and logical OR

Logical `OR` selects its right operand when the left operand is false according to FQL's truth-value rules. `??` selects its right operand only for `NONE`.

| Value | `value OR "fallback"` | `value ?? "fallback"` |
| --- | --- | --- |
| `NONE` | `"fallback"` | `"fallback"` |
| `false` | `"fallback"` | `false` |
| `0` | `"fallback"` | `0` |
| `""` | `"fallback"` | `""` |

Use `OR` for truth-based control flow. Use `??` for an absent-value fallback.

## Short-circuiting and errors

The fallback expression is evaluated only when the left operand is `NONE`.

{{< code lang="fql" >}}
FUNC buildDefaultName() => "Anonymous"

LET name = "Ada"

// The fallback function is not called.
RETURN name ?? buildDefaultName()
{{</ code >}}

`??` does not catch runtime errors. An error from the left operand propagates before FQL can test the value:

{{< code lang="fql" >}}
// Division by zero still fails.
RETURN (1 / 0) ?? 42
{{</ code >}}

Apply recovery first when an error should become `NONE` and then select a fallback:

{{< code lang="fql" >}}
// Postfix recovery turns the error into NONE, so this returns 42.
RETURN (1 / 0)? ?? 42
{{</ code >}}

See [Error Handling]({{% ref "../control-flow/error-handling" %}}) for explicit recovery options.

## Associativity and precedence

`??` is right-associative. A chain groups from the right:

{{< code lang="fql" >}}
LET primary = NONE
LET secondary = NONE

RETURN {
    chain: primary ?? secondary ?? "fallback",
    grouped: primary ?? (secondary ?? "fallback")
}
{{</ code >}}

Logical `OR` binds more tightly than `??`, while `??` binds more tightly than the ternary operator. Parentheses can make a different grouping explicit.

See [Operator Precedence]({{< ref "precedence" >}}) for the complete order.

## Next steps

{{< docs-related tiles="language-operators,language-operators-logical,language-operators-ternary,language-control-flow-error-handling" >}}
