---
title: "Range Operator"
sidebarTitle: "Range"
weight: 40
draft: false
description: "The range operator for producing sequences of integer values."
---

# Range operator

FQL supports expressing simple numeric ranges with the ``..`` operator. This operator can be used to easily iterate over a sequence of numeric values.

The ``..`` operator will produce an array of the integer values in the defined range, with both bounding values included.

{{< editor lang="fql" >}}
return 2010..2013
{{</ editor >}}

Inside an array expansion, `.` refers to the current item. When a dot touches the range operator, the formatter adds spaces on both sides of `..` to keep the tokens separate. Here the current item is the lower bound, so the result is `[[1, 2, 3]]`.

{{< editor lang="fql" >}}
return [1][* return . .. 3]
{{</ editor >}}

## Next steps

{{< docs-related tiles="language-operators,language-control-flow-for,language-operators-array" >}}
