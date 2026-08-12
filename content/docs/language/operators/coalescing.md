---
title: "none-Coalescing Operator"
sidebarTitle: "none Coalescing"
weight: 25
draft: false
description: "Select a fallback only when a value is none with the right-associative coalescing operator."
---

# none-coalescing operator

The none-coalescing operator (`??`) selects a fallback only when its left operand evaluates to `none`.

{{< code lang="fql" >}}
let user = {
    profile: {}
}

return user?.profile?.displayName ?? "Anonymous"
{{</ code >}}

If `displayName` is missing, optional chaining produces `none` and the expression returns `"Anonymous"`.

## How it evaluates

For `left ?? right`, FQL evaluates `left` once.

- If `left` is not `none`, the expression returns it without evaluating `right`.
- If `left` is `none`, FQL evaluates and returns `right`.

The two operands may have different types. The selected operand becomes the result without conversion.

## Only none selects the fallback

`null` represents the same absent runtime value as `none`, so both select the fallback. Other values remain present even when a logical operator would treat them as false.

| Left operand | Result of `left ?? "fallback"` |
| --- | --- |
| `none` | `"fallback"` |
| `null` | `"fallback"` |
| `false` | `false` |
| `0` | `0` |
| `""` | `""` |
| `[]` | `[]` |
| `{}` | `{}` |

This makes `??` useful when `false`, zero, or an empty value carries meaning and must not be replaced.

## Coalescing and logical or

Logical `or` selects its right operand when the left operand is false according to FQL's truth-value rules. `??` selects its right operand only for `none`.

| Value | `value or "fallback"` | `value ?? "fallback"` |
| --- | --- | --- |
| `none` | `"fallback"` | `"fallback"` |
| `false` | `"fallback"` | `false` |
| `0` | `"fallback"` | `0` |
| `""` | `"fallback"` | `""` |

Use `or` for truth-based control flow. Use `??` for an absent-value fallback.

## Short-circuiting and errors

The fallback expression is evaluated only when the left operand is `none`.

{{< code lang="fql" >}}
func buildDefaultName() => "Anonymous"

let name = "Ada"

// The fallback function is not called.
return name ?? buildDefaultName()
{{</ code >}}

`??` does not catch runtime errors. An error from the left operand propagates before FQL can test the value:

{{< code lang="fql" >}}
// Division by zero still fails.
return (1 / 0) ?? 42
{{</ code >}}

Apply recovery first when an error should become `none` and then select a fallback:

{{< code lang="fql" >}}
// Postfix recovery turns the error into NONE, so this returns 42.
return (1 / 0)? ?? 42
{{</ code >}}

See [Error Handling]({{% ref "../control-flow/error-handling" %}}) for explicit recovery options.

## Associativity and precedence

`??` is right-associative. A chain groups from the right:

{{< code lang="fql" >}}
let primary = none
let secondary = none

return {
    chain: primary ?? secondary ?? "fallback",
    grouped: primary ?? (secondary ?? "fallback")
}
{{</ code >}}

Logical `or` binds more tightly than `??`, while `??` binds more tightly than the ternary operator. Parentheses can make a different grouping explicit.

See [Operator Precedence]({{< ref "precedence" >}}) for the complete order.

## Next steps

{{< docs-related tiles="language-operators,language-operators-logical,language-operators-ternary,language-control-flow-error-handling" >}}
