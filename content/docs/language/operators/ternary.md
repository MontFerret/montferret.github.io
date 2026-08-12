---
title: "Ternary Operator"
sidebarTitle: "Ternary"
weight: 50
draft: false
description: "Conditional evaluation with the ternary operator and its shortcut form."
---

# Ternary operator

The ternary operator provides compact two-way value selection. It expects a boolean condition as its first operand and returns the second operand when the condition is true, or the third operand otherwise. For structured branching, patterns, or condition chains, use [`match`]({{% ref "../control-flow/match" %}}).

{{< code lang="fql" >}}
u.age > 15 || u.active == true ? u.userId : none
{{</ code >}}

There is also a shortcut variant of the ternary operator with just two operands. This variant can be used when the expression for the boolean condition and the return value should be the same:

{{< code lang="fql" >}}
u.value ? : "value is NONE, 0 or not present"
{{</ code >}}

## Next steps

{{< docs-related tiles="language-operators,language-operators-coalescing,language-control-flow-match,language-operators-logical" >}}
