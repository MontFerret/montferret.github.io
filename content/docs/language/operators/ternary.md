---
title: "Ternary Operator"
sidebarTitle: "Ternary"
weight: 50
draft: false
description: "Conditional evaluation with the ternary operator and its shortcut form."
---

# Ternary operator

FQL supports a ternary operator that can be used for conditional evaluation. The ternary operator expects a boolean condition as its first operand, and it returns the result of the second operand if the condition evaluates to true, and the third operand otherwise.

{{< code lang="fql" >}}
u.age > 15 || u.active == true ? u.userId : NONE
{{</ code >}}

There is also a shortcut variant of the ternary operator with just two operands. This variant can be used when the expression for the boolean condition and the return value should be the same:

{{< code lang="fql" >}}
u.value ? : "value is NONE, 0 or not present"
{{</ code >}}

## Next steps

{{< docs-related tiles="language-operators,language-control-flow-match,language-operators-logical" >}}
