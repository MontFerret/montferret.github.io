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
RETURN 2010..2013
{{</ editor >}}

## Next steps

{{< docs-related tiles="language-operators,language-control-flow-for,language-operators-array" >}}
