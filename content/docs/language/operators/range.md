---
title: "Range Operator"
sidebarTitle: "Range"
weight: 40
draft: false
description: "The range operator for producing sequences of integer values."
---

# Range operator

FQL supports expressing simple numeric ranges with the ``..`` operator. This operator can be used to easily iterate over a sequence of numeric values.

The ``..`` operator produces an iterable range that includes both integer bounds. Use it in a loop, or return it directly to serialize its values as a JSON array. A range can also count down, as in `3..1`.

{{< editor lang="fql" >}}
return 2010..2013
{{</ editor >}}

Inside an array expansion, `.` refers to the current item. When a dot touches the range operator, the formatter adds spaces on both sides of `..` to keep the tokens separate. Here the current item is the lower bound, so the result is `[[1, 2, 3]]`.

{{< editor lang="fql" >}}
return [1][* return . .. 3]
{{</ editor >}}

## Next steps

{{< docs-related tiles="language-operators,language-control-flow-for,language-operators-array" >}}
