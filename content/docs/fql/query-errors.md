---
title: "Query errors"
weight: 8
draft: false
---

# Query errors
Issuing an invalid query to Ferret will result in a parse error if the query is syntactically invalid. The Ferret compiler will detect such errors during query inspection and abort further processing. Instead, an error message are returned so that the errors can be fixed.

Under some circumstances, executing a query may also produce run-time errors that cannot be predicted from inspecting the query text alone. This is because queries may use data from collections that may also be inhomogeneous. Some examples that will cause run-time errors are:

- Division by zero: Will be triggered when an attempt is made to use the value 0 as the divisor in an arithmetic division or modulus operation
- Invalid argument types: Some function can return an error if given arguments have invalid type.
- Timeouts: Interactions with web pages may timeout due to some changes on the pages or invalid query logic or selectors.