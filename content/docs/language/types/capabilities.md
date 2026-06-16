---
title: "Capability Types"
sidebarTitle: "Capability Types"
weight: 30
draft: false
description: "Behavioral contracts that describe what operations a value supports, regardless of its concrete type."
---

# Capability types

Capability types describe behavior that a value supports.

Unlike basic types, capability types are not about a value's storage shape. They are about the operations that can be performed on the value. A capability type represents an operation-oriented contract: a value has a capability when the runtime can perform a specific class of operation on it.

This is especially important for host values, because host values may expose behavior that does not correspond to any built-in FQL type. A database connection may be queryable. A file handle may be closable. A cursor may be iterable. A browser element may accept dispatched commands.

A capability is not another primitive type like `string` or `array`. It is a behavioral contract between a value and the operations that use it.

## Capabilities vs basic types

Basic types describe what a value *is*. Capability types describe what a value can *do*.

An array is a basic value type. It is also iterable and sortable. These are separate facts: the type tells you the shape, and the capabilities tell you the operations.

A host database connection is not an array, an object, or a string. It has none of those storage shapes. However, it may be queryable if the host exposes query execution for it.

{{< code lang="fql" >}}
LET db = DB::SQLITE::OPEN("data.db")

RETURN QUERY `
  SELECT name
  FROM users
` IN db USING sql
{{</ code >}}

The `QUERY ... IN` operation does not check whether `db` is an array or an object. It checks whether `db` supports query execution. That check is a capability check.

## Capabilities vs host types

Host types often use capabilities to integrate with the language.

The runtime does not need to know every concrete host type in advance. Instead, it checks whether the value supports the capability required by an operation. This is the key mechanism that allows host values to participate in FQL expressions without being predefined by the language.

{{< code lang="fql" >}}
QUERY `
  SELECT * FROM users
` IN db USING sql
{{</ code >}}

In this example, `db` is a host value. The query expression does not require `db` to be a specific concrete database type. It requires `db` to support the queryable capability. Any host value that supports query execution can be used here, regardless of whether it wraps SQLite, PostgreSQL, an in-memory store, or something else entirely.

This is the reason capability types exist as a separate concept. They decouple language operations from concrete value implementations.

See [Host Values]({{< ref "host" >}}) for more about host values and their lifecycle.

## Common capabilities

The following capabilities describe the most common behavioral contracts that values can support.

### Queryable

A queryable value can execute query literals.

{{< code lang="fql" >}}
QUERY `SELECT * FROM users` IN db USING sql
{{</ code >}}

Query execution supports several forms:

{{< code lang="fql" >}}
// All matching values
QUERY `SELECT * FROM users` IN db USING sql

// First matching value
QUERY ONE `SELECT * FROM users WHERE id = 1` IN db USING sql

// Count of matching values
QUERY COUNT `SELECT * FROM users` IN db USING sql

// Whether any match exists
QUERY EXISTS `SELECT * FROM users WHERE active = true` IN db USING sql
{{</ code >}}

The query literal, any parameters supplied with `WITH`, and any options supplied with `OPTIONS` are all passed to the value. The value decides how to interpret and execute the query.

### Iterable

An iterable value can produce a sequence of values, which allows it to be used with `FOR ... IN`.

{{< code lang="fql" >}}
FOR item IN collection
    RETURN item
{{</ code >}}

Built-in arrays and objects are iterable. Host values such as cursors, result sets, or streams may also be iterable.

When iterating, each step produces both a value and a key. For arrays, the key is the index. For objects, the key is the property name. For host values, the key depends on the value's implementation.

### Comparable

A comparable value can define equality and ordering behavior.

Comparison is used by equality operators (`==`, `!=`), ordering operators (`<`, `>`, `<=`, `>=`), and sorting. When a value supports comparison, it controls how it is ordered relative to other values.

Built-in types define comparison behavior directly. Host values may define their own comparison rules based on identity, an internal key, a normalized representation, or another rule chosen by the runtime or module.

See [Type Ordering]({{< ref "ordering" >}}) for the full ordering model.

### Sortable

A sortable value can be sorted in place.

{{< code lang="fql" >}}
SORT values ASC
{{</ code >}}

Arrays are sortable by default. A host value that represents a collection may also support sorting if the host provides that capability.

### Observable

An observable value can produce a stream of events over time.

{{< code lang="fql" >}}
WAITFOR EVENT "load" IN page
{{</ code >}}

Observability is used for event-driven operations where the script waits for something to happen. The value defines which events it can produce and how subscriptions work.

### Dispatchable

A dispatchable value can receive and handle dispatched commands.

{{< code lang="fql" >}}
DISPATCH "click" IN button
{{</ code >}}

Dispatch is effectful: it causes a side effect on the target value without producing a return value. Browser elements, UI controls, and other interactive host values commonly support dispatch.

### Measurable

A measurable value has a defined length.

{{< code lang="fql" >}}
LENGTH(elements)
{{</ code >}}

This capability is used when the runtime needs to know the size of a value, such as for length checks, emptiness tests, or size-based operations. Arrays, objects, strings, and host collections may all be measurable.

### Closable

A closable value owns an external resource that should be released after execution.

Database connections, file handles, browser sessions, network clients, and other resource-backed values may support close. The runtime tracks closable values and releases their resources during program finalization, before the result is returned to the host.

Host values that support close do not need to be closed manually in FQL code. The runtime handles cleanup automatically.

See [Host Values]({{< ref "host" >}}) for more about host values and their lifecycle.

### Readable

A readable value supports member access — retrieving an element by position or by key.

Index-based access uses bracket notation with an integer position:

{{< code lang="fql" >}}
LET first = items[0]
LET last = items[LENGTH(items) - 1]
{{</ code >}}

Key-based access uses dot notation or bracket notation with a string key:

{{< code lang="fql" >}}
LET name = user.name
LET name = user["name"]
{{</ code >}}

Safe navigation returns `NONE` instead of raising an error when the target value is `NONE`:

{{< code lang="fql" >}}
LET city = user?.address?.city
LET first = items?[0]
{{</ code >}}

Built-in arrays are readable by index. Built-in objects are readable by key. Host values may support either or both forms of access.

### Writable

A writable value supports member assignment — setting an element at a position or under a key.

{{< code lang="fql" >}}
VAR items = [1, 2, 3]
items[0] = 10

VAR user = { name: "Ada" }
user.name = "Grace"
user["active"] = true
{{</ code >}}

Built-in arrays support index assignment. Built-in objects support key assignment. Host values may support either or both forms.

### Removable

A removable value supports member deletion through the `DELETE` statement.

{{< code lang="fql" >}}
VAR user = { name: "Ada", deprecated: true }
DELETE user.deprecated
DELETE user["deprecated"]
{{</ code >}}

Deletion removes the member entirely — it is not the same as assigning `NONE`, which keeps the key present with an absent value.

Built-in objects support key removal. Host values may support removal if the host provides that capability.

## How operations use capabilities

When an operation requires a capability, the runtime checks whether the value supports it. If it does, the operation proceeds. If it does not, the operation fails with a runtime error.

This check happens at runtime, not at parse time. FQL does not statically verify that a variable holds a value with a particular capability. The check occurs when the operation is actually executed.

| Operation | Required capability |
| --- | --- |
| `QUERY ... IN value` | Queryable |
| `FOR item IN value` | Iterable |
| `SORT value` | Sortable |
| `WAITFOR EVENT ... IN value` | Observable |
| `DISPATCH ... IN value` | Dispatchable |
| `value == other` | Comparable (when custom ordering is needed) |
| `value[index]` | Readable (by index) |
| `value.key` / `value["key"]` | Readable (by key) |
| `value[index] = x` | Writable (by index) |
| `value.key = x` | Writable (by key) |
| `DELETE value.key` | Removable |

A value may support multiple capabilities simultaneously. An array is iterable, sortable, measurable, comparable, and readable by index. An object is iterable, measurable, comparable, readable by key, writable by key, and removable by key. A host cursor might be iterable and closable but not queryable or sortable.

## Runtime-defined behavior

Capability behavior is runtime-defined.

For built-in values, Ferret defines the capability behavior directly. For host values, the embedding runtime defines which capabilities are supported and how they behave. Two different host values may support the same capability but implement it differently.

A host value may support equality, ordering, querying, iteration, cleanup, or serialization differently from another host value. As long as the behavioral contract is satisfied, the runtime and the language operations work correctly.

The exact implementation mechanism depends on the host runtime. In the Go runtime, capabilities are represented by interfaces implemented by runtime values.

## Error behavior

Using an operation with a value that does not support the required capability results in a runtime error.

{{< editor lang="fql" >}}
RETURN QUERY `SELECT * FROM users` IN "not a database" USING sql
{{</ code >}}

This fails because a string does not support query execution.

{{< editor lang="fql" >}}
FOR item IN 42
    RETURN item
{{</ code >}}

This fails because a number is not iterable.

{{< editor lang="fql" >}}
RETURN DISPATCH "click" IN "not an element"
{{</ code >}}

This fails because a string does not support dispatch.

Capability errors are reported at the point where the operation is attempted. The error identifies the value, its type, and the capability that was expected.

## Where to go next

{{< docs-related tiles="language-types-host,language-types-ordering,language-types-serialization,embedding" >}}
