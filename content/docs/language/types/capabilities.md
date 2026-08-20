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
let db = db::sqlite::open("data.db")

return query `
  SELECT name
  FROM users
` in db using sql
{{</ code >}}

The `query ... in` operation does not check whether `db` is an array or an object. It checks whether `db` supports query execution. That check is a capability check.

## Capabilities vs host types

Host types often use capabilities to integrate with the language.

The runtime does not need to know every concrete host type in advance. Instead, it checks whether the value supports the capability required by an operation. This is the key mechanism that allows host values to participate in FQL expressions without being predefined by the language.

{{< code lang="fql" >}}
query `
  SELECT * FROM users
` in db using sql
{{</ code >}}

In this example, `db` is a host value. The query expression does not require `db` to be a specific concrete database type. It requires `db` to support the queryable capability. Any host value that supports query execution can be used here, regardless of whether it wraps SQLite, PostgreSQL, an in-memory store, or something else entirely.

This is the reason capability types exist as a separate concept. They decouple language operations from concrete value implementations.

See [Host Values]({{< ref "host" >}}) for more about host values and their lifecycle.

## Common capabilities

The following capabilities describe the most common behavioral contracts that values can support.

### Queryable

A queryable value can execute query literals.

{{< code lang="fql" >}}
query `SELECT * FROM users` in db using sql
{{</ code >}}

Query execution supports several forms:

{{< code lang="fql" >}}
// All matching values
query `SELECT * FROM users` in db using sql

// First matching value
query one `SELECT * FROM users WHERE id = 1` in db using sql

// Count of matching values
query count `SELECT * FROM users` in db using sql

// Whether any match exists
query exists `SELECT * FROM users WHERE active = true` in db using sql
{{</ code >}}

The query literal, any parameters supplied with `with`, and any options supplied with `options` are all passed to the value. The value decides how to interpret and execute the query.

### Iterable

An iterable value can produce a sequence of values, which allows it to be used with `for ... in`.

{{< code lang="fql" >}}
return for item in collection {
    return item
}
{{</ code >}}

Built-in arrays and objects are iterable. Host values such as cursors, result sets, or streams may also be iterable.

When iterating, each step produces both a value and a key. For arrays, the key is the index. For objects, the key is the property name. For host values, the key depends on the value's implementation.

### Equatable

An equatable host value defines equality for `==`, `!=`, membership, grouping, set operations, and deduplication. Equality is independent from relational comparison: a value may support equality without supporting `<`, `>`, `<=`, or `>=`.

Equal values must produce equal hashes. A hash only selects possible matches; Ferret still verifies equality before treating two values as the same.

### Comparable

A comparable host value defines relational comparison for `<`, `>`, `<=`, `>=`, and sorting within a compatible comparison domain.

Built-in types define comparison behavior directly. Host values may compare by identity, an internal key, a normalized representation, or another stable rule chosen by the runtime or module. If a host value implements both Equatable and Comparable, the two contracts must agree about which values are equal.

See [Type Ordering]({{< ref "ordering" >}}) for the full ordering model.

### Arithmetic

Host values can define binary arithmetic with independent capabilities for addition, subtraction, multiplication, division, and modulus. Supporting one operation does not imply support for another.

Each capability explicitly handles both operand positions. For example, `host + value` uses the host's `Add` method, while `value + host` uses its `RightAdd` method. The right-hand method receives the original left operand; Ferret never silently reverses subtraction, division, or modulus.

Native arithmetic runs first, including String-triggered concatenation. Host capability negotiation is used only when the native operand pair is otherwise unsupported. A host implementation can decline one operand arrangement so Ferret can try the other value's right-hand method; genuine host errors stop evaluation immediately.

These capabilities apply only to binary operators. Unary arithmetic and increment or decrement remain defined by native runtime semantics.

### Sortable

A sortable value can be sorted in place.

{{< code lang="fql" >}}
sort values asc
{{</ code >}}

Arrays are sortable by default. A host value that represents a collection may also support sorting if the host provides that capability.

### Observable

An observable value can produce a stream of events over time.

{{< code lang="fql" >}}
waitfor event "load" in page
{{</ code >}}

Observability is used for event-driven operations where the script waits for something to happen. The value defines which events it can produce and how subscriptions work.

### Dispatchable

A dispatchable value can receive and handle dispatched commands.

{{< code lang="fql" >}}
dispatch "click" in button
{{</ code >}}

Dispatch is effectful: it causes a side effect on the target value without producing a return value. Browser elements, UI controls, and other interactive host values commonly support dispatch.

### Measurable

A measurable value has a defined length.

{{< code lang="fql" >}}
length(elements)
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
let first = items[0]
let last = items[length(items) - 1]
{{</ code >}}

Key-based access uses dot notation or bracket notation with a string key:

{{< code lang="fql" >}}
let name = user.name
let name = user["name"]
{{</ code >}}

Safe navigation returns `none` instead of raising an error when the target value is `none`:

{{< code lang="fql" >}}
let city = user?.address?.city
let first = items?[0]
{{</ code >}}

Built-in arrays are readable by index. Built-in objects are readable by key. Host values may support either or both forms of access.

### Writable

A writable value supports member assignment — setting an element at a position or under a key.

{{< code lang="fql" >}}
var items = [1, 2, 3]
items[0] = 10

var user = { name: "Ada" }
user.name = "Grace"
user["active"] = true
{{</ code >}}

Built-in arrays support index assignment. Built-in objects support key assignment. Host values may support either or both forms.

### Removable

A removable value supports member deletion through the `delete` statement.

{{< code lang="fql" >}}
var user = { name: "Ada", deprecated: true }
delete user.deprecated
delete user["deprecated"]
{{</ code >}}

Deletion removes the member entirely — it is not the same as assigning `none`, which keeps the key present with an absent value.

Built-in objects support key removal. Host values may support removal if the host provides that capability.

## How operations use capabilities

When an operation requires a capability, the runtime checks whether the value supports it. If it does, the operation proceeds. If it does not, the operation fails with a runtime error.

This check happens at runtime, not at parse time. FQL does not statically verify that a variable holds a value with a particular capability. The check occurs when the operation is actually executed.

| Operation | Required capability |
| --- | --- |
| `query ... in value` | Queryable |
| `for item in value` | Iterable |
| `sort value` | Sortable |
| `waitfor event ... in value` | Observable |
| `dispatch ... in value` | Dispatchable |
| `value == other` / `value != other` | Equatable |
| `value < other` and other relational comparisons | Comparable |
| `value + other` | Addable |
| `value - other` | Subtractable |
| `value * other` | Multipliable |
| `value / other` | Dividable |
| `value % other` | Modulable |
| `value[index]` | Readable (by index) |
| `value.key` / `value["key"]` | Readable (by key) |
| `value[index] = x` | Writable (by index) |
| `value.key = x` | Writable (by key) |
| `delete value.key` | Removable |

A value may support multiple capabilities simultaneously. An array is iterable, sortable, measurable, equatable, comparable, and readable by index. An object is iterable, measurable, equatable, comparable, readable by key, writable by key, and removable by key. A host cursor might be iterable and closable but not queryable or sortable.

## Runtime-defined behavior

Capability behavior is runtime-defined.

For built-in values, Ferret defines the capability behavior directly. For host values, the embedding runtime defines which capabilities are supported and how they behave. Two different host values may support the same capability but implement it differently.

A host value may support equality, relational comparison, querying, iteration, cleanup, or serialization differently from another host value. As long as the behavioral contract is satisfied, the runtime and the language operations work correctly.

The exact implementation mechanism depends on the host runtime. In the Go runtime, capabilities are represented by interfaces implemented by runtime values.

Context-aware capabilities receive the execution context unchanged. Ferret does not poll cancellation inside a capability method; an implementation that may block or perform remote work must observe cancellation while it retains control.

## Error behavior

Using an operation with a value that does not support the required capability results in a runtime error.

{{< editor lang="fql" >}}
return query `SELECT * FROM users` in "not a database" using sql
{{</ editor >}}

This fails because a string does not support query execution.

{{< editor lang="fql" >}}
return for item in 42 {
    return item
}
{{</ editor >}}

This fails because a number is not iterable.

{{< editor lang="fql" >}}
return dispatch "click" in "not an element"
{{</ editor >}}

This fails because a string does not support dispatch.

Capability errors are reported at the point where the operation is attempted. The error identifies the value, its type, and the capability that was expected.

## Next steps

{{< docs-related tiles="language-types-host,language-types-ordering,language-types-serialization,embedding" >}}
