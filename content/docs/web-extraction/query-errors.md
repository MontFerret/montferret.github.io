---
title: "Query Errors"
weight: 20
draft: false
description: "Understand Ferret query error output."
aliases:
  - /docs/fql/query-errors/
---

# Query errors

Compilation checks syntax and bindings before execution. Use [`ferret check`]({{< ref "/docs/tools/cli/check" >}}) to locate a malformed expression or an undeclared variable.

Errors that depend on input or a host operation can surface during execution. A divisor can be zero, a function can receive an unsupported type, or a browser operation can fail. For a missing element, `QUERY ONE` returns `NONE`; optional member access such as `title?.textContent` lets the script return that absence.

Use the [error handling reference]({{< ref "/docs/language/control-flow/error-handling" >}}) for recovery and [WAITFOR]({{< ref "/docs/language/control-flow/waitfor" >}}) to choose how a wait should handle its timeout.
