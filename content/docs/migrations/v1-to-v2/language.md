---
title: "FQL language migration"
sidebarTitle: "FQL language"
weight: 10
draft: false
description: "Make script results explicit when migrating FQL from Ferret v1 to v2."
---

# FQL language migration

Ferret v1 implicitly returned the value of a final top-level collecting `FOR`.
In v2, return that value explicitly with `RETURN FOR`.

Before, in v1:

{{< code lang="fql" >}}
FOR item IN 1..3
    RETURN item
{{< /code >}}

After, in v2:

{{< code lang="fql" >}}
RETURN FOR item IN 1..3 {
    RETURN item
}
{{< /code >}}

The migrated script returns `[1, 2, 3]`. The outer `RETURN` makes the loop's
collected array the script result; the inner `RETURN` supplies each item.

| v1 source | v2 action | Automated? |
| --- | --- | --- |
| Final top-level collecting `FOR` with an implicit script result | Add an explicit outer `RETURN`; format the loop | Yes, when structurally recognized and no explicit terminal return exists |
| Already explicit `RETURN FOR` | Keep the explicit result | No change needed |
| Assigned, nested, function-contained, or non-final loop | Review result ownership in context | No independent loop wrapping |
| FQL embedded in a Go string | Apply the same language review to the embedded query | No; the CLI discovers `.fql` files |

See the [CLI migration reference]({{< ref "docs/tools/cli/migrate" >}}#migrate-a-final-for)
for command behavior and source discovery.

## Review discarded loop results

A standalone collecting `FOR` in v2 executes its body but discards its collected
array. If the script falls through without a terminal `RETURN`, its result is
`none`. Parentheses around a standalone loop do not retain that result.

Use `RETURN FOR` when the array is the intended result. Use a braced loop without
an iteration `RETURN` when only side effects are intended. For example:

{{< code lang="fql" >}}
VAR total = 0
FOR value IN [1, 2, 3] {
    total += value
}
RETURN total
{{< /code >}}

This returns `6`. The loop updates `total`; the script returns it separately.

Review nested loops before changing their returns. An explicit inner
`RETURN FOR` collects an inner array, while a legacy terminal pass-through loop
can flatten results. Adding returns at every nesting level can therefore change
the output shape. The migrator wraps the recognized final top-level loop only.

Current collecting, returnless, and nested-loop behavior is documented in
[For loops]({{< ref "docs/language/control-flow/for" >}}#returning-nesting-and-discarding-loop-results).
See [Script structure]({{< ref "docs/language/script-structure" >}}) for the
current statement and return model.

## Next steps

{{< docs-related tiles="migrations-v1-to-v2-standard-library,migrations-v1-to-v2-go-embedding,migrations-v1-to-v2" >}}
