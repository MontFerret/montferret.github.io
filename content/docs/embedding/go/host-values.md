---
title: "Host Values"
sidebarTitle: "Host Values"
weight: 60
draft: false
description: "Implement Go types that FQL scripts can hold, access, iterate, query, and dispatch to."
aliases:
    - /docs/embedding/host-values/
---

# Host Values

Host values let an embedding application expose Go-managed resources and objects to FQL scripts. This page covers how to implement them in Go. For the FQL perspective — what host values look like inside a script and which operations they support — see [Host Values]({{< ref "docs/language/types/host" >}}) and [Capability Types]({{< ref "docs/language/types/capabilities" >}}).

## When to implement a host value

- **Resource handle** — the Go object owns an external resource (database connection, file handle, HTTP client) whose lifecycle outlives a single function call.
- **Lazy value** — the value materializes data only when accessed, avoiding upfront cost (cursor, paginated result set).
- **Domain object** — the value has rich behavior (property access, iteration, dispatch) that cannot be expressed as a plain array or object.
- **External client** — the value wraps a Go client library and exposes its operations through FQL capabilities (query, dispatch, subscribe).
- **Controlled capabilities** — you want to expose only certain operations (read-only properties, no iteration) rather than giving scripts full access to the underlying data.

## Choosing the right mechanism

| Mechanism | Best for | Script sees | Lifecycle |
|-----------|----------|-------------|-----------|
| Parameter | Static data: strings, numbers, config, JSON-like structures | `@name` — a plain value | Host sets before execution |
| Function | Stateless computation: hash, format, fetch | `fn(args)` → result | No state between calls |
| Host value | Stateful resource or behavioral object | Variable with capabilities: `val.prop`, `for x in val`, `query ... in val` | See [cleanup and finalization](#cleanup-and-finalization) |
| Module | Self-contained extension bundling functions + hooks | Namespaced functions and values | Module `Register` + engine lifecycle hooks |

When in doubt, walk through these questions:

| Question | If yes | If no |
|----------|--------|-------|
| Does the Go object own a resource that must be closed? | Host value (implement `io.Closer`) | ↓ |
| Does the script need to access properties, iterate, query, or dispatch? | Host value (implement capability interfaces) | ↓ |
| Is the data static and known before execution? | Parameter (`WithRuntimeParam`) | ↓ |
| Is it a one-shot computation with no state? | Function | ↓ |
| Does it bundle multiple functions + hooks? | Module | Function or host value |

## The Value interface

Every host value must satisfy `runtime.Value`:

{{< code lang="go" >}}
type Value interface {
    fmt.Stringer   // String() string
    Hashable       // Hash() uint64
    Copy() Value
}
{{</ code >}}

- `String()` — text representation used by `to_string()` and encoding fallbacks.
- `Hash()` — identity hash used by the VM for deduplication and map keys. Use `fnv.New64a` from the standard library.
- `Copy()` — shallow copy. For host values that wrap a pointer, returning the same pointer is usually correct.

Optionally implement `Typed` so the runtime can report your type in error messages and enable generic argument casting with `CastArg`:

{{< code lang="go" >}}
type Typed interface {
    Type() Type
}
{{</ code >}}

Create a type with `runtime.NewTypeFor`:

{{< code lang="go" >}}
var MyValueType = runtime.NewTypeFor[*MyValue]()
{{</ code >}}

## Copying a capability value

Use `runtime.Copy` when the caller needs the copied value to preserve a concrete Go type or a capability interface:

{{< code lang="go" >}}
func copyList(list runtime.List) (runtime.List, error) {
    next, err := runtime.Copy(list)
    if err != nil {
        return nil, err
    }

    return next, nil
}
{{</ code >}}

The helper calls `Value.Copy()` once and checks that its result implements the static type of the argument (`runtime.List` here). A different concrete implementation is accepted if it still implements that interface. An incompatible result produces an error wrapping `runtime.ErrInvalidType`, including the source and returned types, instead of a type-assertion panic.

Copying remains shallow. The host implementation must provide the storage independence required by its capability; this check does not validate that independence or change resource ownership. It does not recover panics inside a host's `Copy` method. The existing `Value.Copy() Value` method stays unchanged.

## Constructing an empty collection

`runtime.List`
embeds `Factory[List]` and `runtime.Map` embeds `Factory[Map]`:

{{< code lang="go" >}}
type Factory[T runtime.Collection] interface {
    New(context.Context) (T, error)
}
{{</ code >}}

This intentionally replaces `Spawnable[T].Empty`. Rename custom collection
methods to `New` and update calls; no compatibility alias is provided.

`New` creates an empty, independently mutable collection of the receiver's
implementation family. Preserve configuration such as backend, encoding, limits,
and placement policy without copying elements or changing the source. Collections
may share a backend service, but must have distinct logical mutable contents.
A wrapper must override an inherited factory if it would lose that wrapper's
implementation or configuration.

Unlike `Copy` and `Clone`, `New` contains no source elements. Unlike `Clear`,
it leaves the existing collection intact. Construction can perform I/O and fail;
blocking implementations must honor the context. Ownership transfers only on
success. If construction fails, the factory releases partially acquired resources
and joins cleanup failures with the construction error.

`collections::reverse(list)` uses this factory and appends existing element references
in reverse index order. Failed reversal closes the incomplete destination when
closable, preserving both operation and cleanup errors. A successful destination
passes to normal result lifecycle handling. The source and its elements remain
borrowed. Immutable object merges also construct through the first source map's
`New`, retaining their existing clone/copy behavior for inserted values. After
successful construction, a failed or canceled immutable merge closes its incomplete
closable destination and joins cleanup errors with the primary error. Success
transfers the open result; source maps and mutable merge targets remain borrowed.

Read-only values need only `Value` and `Iterable` for `count` and
`count_distinct`; they need not implement List or Map. `count` additionally
uses `Measurable` when present, without creating an iterator or retrying a failed
length operation. Negative host lengths fail without traversal; zero is valid.

## Minimal implementation

The smallest possible host value — a struct that wraps a string:

{{< code lang="go" >}}
package example

import (
    "hash/fnv"

    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Label struct {
    text string
}

var LabelType = runtime.NewTypeFor[*Label]()

func (l *Label) Type() runtime.Type { return LabelType }
func (l *Label) String() string     { return l.text }

func (l *Label) Hash() uint64 {
    h := fnv.New64a()
    h.Write([]byte(l.text))
    return h.Sum64()
}

func (l *Label) Copy() runtime.Value {
    return &Label{text: l.text}
}
{{</ code >}}

Register a function that returns it. This snippet also imports `context` and the root `ferret` package:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
        ns.Function().A1().Add("label", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
            text, err := runtime.CastArg[runtime.String](arg, 0)
            if err != nil {
                return nil, err
            }

            return &Label{text: string(text)}, nil
        })
    }),
)
{{</ code >}}

{{< code lang="fql" >}}
let l = label("urgent")
return to_string(l)
// "urgent"
{{</ code >}}

## Adding capabilities

A host value without capability interfaces is opaque — scripts can hold it and pass it around, but cannot access properties, iterate over it, or query it. Implement capability interfaces to enable FQL operations on the value.

### Property and index access

| Interface | Method | FQL syntax               |
|-----------|--------|--------------------------|
| `KeyReadable` | `Get(ctx, key Value) (Value, error)` | `val.prop`, `val["key"]` |
| `KeyLookup` | `Lookup(ctx, key Value) (Value, bool, error)` | `val?.prop`              |
| `KeyWritable` | `Set(ctx, key, value Value) error` | `val.prop = x`           |
| `KeyRemovable` | `RemoveKey(ctx, key Value) error` | `delete val.prop`         |
| `IndexReadable` | `At(ctx, idx Int) (Value, error)` | `val[0]`                 |
| `IndexLookup` | `LookupAt(ctx, index Int) (Value, bool, error)` | `val?[0]`                |
| `IndexWritable` | `SetAt(ctx, idx Int, value Value) error` | `val[0] = x`             |
| `IndexRemovable` | `RemoveAt(ctx, idx Int) (Value, error)` | `delete val[0]`          |

Adding property access to the `Label` type:

{{< code lang="go" >}}
func (l *Label) Get(_ context.Context, key runtime.Value) (runtime.Value, error) {
    switch key.String() {
    case "text":
        return runtime.NewString(l.text), nil
    case "length":
        return runtime.NewInt(len(l.text)), nil
    default:
        return runtime.None, nil
    }
}
{{</ code >}}

{{< code lang="fql" >}}
let l = label("urgent")
return l.text
// "urgent"
{{</ code >}}

### Iteration and measurement

| Interface | Method | FQL syntax |
|-----------|--------|------------|
| `Iterable` | `Iterate(ctx) (Iterator, error)` | `for x in val` |
| `Iterator` | `Next(ctx) (value, key Value, err error)` | (returned by `Iterate`) |
| `Measurable` | `Length(ctx) (Int, error)` | `length(val)` |
| `Containable` | `Contains(ctx, value Value) (Boolean, error)` | `x in val` |

`Iterator.Next` must return `io.EOF` when the sequence is exhausted. If the iterator holds resources (an open cursor, for example), implement `io.Closer` on it as well.

### Bulk snapshots for spread

Array spread accepts values that implement `runtime.List`, and object spread accepts values that implement `runtime.Map`. Ferret can enumerate those interfaces directly, so host values do not need another capability for correctness.

Remote or paginated collections can optionally provide a native bulk representation:

{{< code lang="go" >}}
type ListSnapshotter interface {
    Snapshot(context.Context) (*Array, error)
}

type MapSnapshotter interface {
    Snapshot(context.Context) (*Object, error)
}
{{</ code >}}

When a valid List or Map implements the matching snapshot interface, spread calls `Snapshot` once and copies from the returned native Array or Object. A successful call must return a non-nil value. If it returns an error, spread propagates that error without retrying through element-by-element or property-by-property enumeration.

The snapshot interfaces are optional optimizations. Implementing `ListSnapshotter` without `runtime.List`, or `MapSnapshotter` without `runtime.Map`, does not make a value spreadable.

`Snapshot` is a bulk read of the collection's current contents. It does not promise transactional or atomic point-in-time consistency for a remote source unless the host implementation provides that guarantee. It also does not transfer ownership and is separate from `vm.Materialize`, which performs terminal conversion of a VM result with explicit lifecycle semantics.

### Query, dispatch, and observe

| Interface | Key method | FQL syntax |
|-----------|------------|------------|
| `Queryable` | `Query(ctx, Query) (List, error)` + 3 modifiers | `query ... in val using ...` |
| `Dispatchable` | `Dispatch(ctx, DispatchEvent) error` | `dispatch "event" TO val` / `val <- event` |
| `Observable` | `Subscribe(ctx, Subscription) (Stream, error)` | `waitfor event "name" in val` |

The `Queryable` interface has four methods: `Query`, `QueryOne`, `QueryCount`, and `QueryExists`. If your implementation only needs the list-returning `Query`, delegate the other three to the built-in helpers:

{{< code lang="go" >}}
func (s *Store) QueryOne(ctx context.Context, q runtime.Query) (runtime.Value, error) {
    return runtime.DefaultQueryOne(ctx, q, s.Query)
}

func (s *Store) QueryCount(ctx context.Context, q runtime.Query) (runtime.Int, error) {
    return runtime.DefaultQueryCount(ctx, q, s.Query)
}

func (s *Store) QueryExists(ctx context.Context, q runtime.Query) (runtime.Boolean, error) {
    return runtime.DefaultQueryExists(ctx, q, s.Query)
}
{{</ code >}}

The `Query` struct carries the expression, dialect, parameters, and options from the FQL statement:

{{< code lang="go" >}}
type Query struct {
    Expression String // the query text
    Kind       String // dialect from USING (empty = default)
    Params     Value  // input from WITH
    Options    Value  // execution policy from OPTIONS
}
{{</ code >}}

### Arithmetic operators

Host values can participate in binary arithmetic by implementing one primitive capability per operator. Each capability has an ordinary method for the host value on the left and an explicit right-hand method for the host value on the right:

| Interface | Left-hand method | Right-hand method | FQL operator |
|-----------|------------------|-------------------|--------------|
| `Addable` | `Add(ctx, right)` | `RightAdd(ctx, left)` | `+` |
| `Subtractable` | `Subtract(ctx, right)` | `RightSubtract(ctx, left)` | `-` |
| `Multipliable` | `Multiply(ctx, right)` | `RightMultiply(ctx, left)` | `*` |
| `Dividable` | `Divide(ctx, right)` | `RightDivide(ctx, left)` | `/` |
| `Modulable` | `Mod(ctx, right)` | `RightMod(ctx, left)` | `%` |

For `host + value`, Ferret calls `host.Add(ctx, value)`. For `value + host`, it calls `host.RightAdd(ctx, value)`. Right-hand methods receive the original left operand; Ferret does not reverse non-commutative operations.

This `Score` value supports addition in both positions. Its ordinary `runtime.Value` methods are omitted here for brevity:

{{< code lang="go" >}}
var _ runtime.Addable = (*Score)(nil)

func (s *Score) Add(ctx context.Context, right runtime.Value) (runtime.Value, error) {
    if err := ctx.Err(); err != nil {
        return runtime.None, err
    }

    amount, ok := right.(runtime.Int)
    if !ok {
        return runtime.None, runtime.ErrUnsupportedOperands
    }

    return &Score{Value: s.Value + int64(amount)}, nil
}

func (s *Score) RightAdd(ctx context.Context, left runtime.Value) (runtime.Value, error) {
    if err := ctx.Err(); err != nil {
        return runtime.None, err
    }

    amount, ok := left.(runtime.Int)
    if !ok {
        return runtime.None, runtime.ErrUnsupportedOperands
    }

    return &Score{Value: int64(amount) + s.Value}, nil
}
{{</ code >}}

Return `runtime.ErrUnsupportedOperands` when the capability exists but does not accept a particular operand arrangement. Ferret then tries the other operand's right-hand method. Return any other error for a genuine execution failure; Ferret propagates it immediately without trying another implementation. If both operands decline, the expression produces the normal invalid-operation error.

The primitive capabilities are independent. Implementing `Addable` does not imply subtraction, and implementing `Multipliable` does not imply division or modulus. `Additive`, `Multiplicative`, and `Arithmetic` combine primitive interfaces for compile-time assertions, but runtime dispatch never checks these compound interfaces.

Native arithmetic remains authoritative. In particular, if either `+` operand is a native String, Ferret performs String concatenation before checking `Addable`. Unary `+`, unary `-`, `++`, and `--` do not use these binary capabilities.

### Comparison and other capabilities

| Interface | Method | Purpose |
|-----------|--------|---------|
| `Equatable` | `Equal(ctx, other Value) (bool, error)` | Enables host-defined `==`, `!=`, membership, grouping, and deduplication |
| `Comparable` | `Compare(ctx, other Value) (Ordering, error)` | Enables relational comparison and sorting |
| `Sortable` | `SortAsc(ctx) error`, `SortDesc(ctx) error` | In-place sort |
| `Cloneable` | `Clone(ctx) (Cloneable, error)` | Deep copy |
| `Unwrappable` | `Unwrap() any` | Extract the inner Go value |
| `DebugInspectable` | `DebugInfo() DebugInfo` | Presentation hints for the debugger |

Use `DebugInfo` to give a host value a semantic type or a short display. For
example, an HTML node that already holds its tag and class can return
`runtime.DebugInfo{TypeName: "HTML::Node", Display: "<div.product-card>"}`.
Each field is optional. An opaque value with no hints appears as `<Resource>`.
An `sdk.HostValue` with an explicitly declared Ferret type can show that name
while keeping its wrapped target private.

The debugger calls the hook once per description. Return metadata already in
memory. Reading children, running a query, or formatting the entire resource
inside the hook makes listing locals pay that cost. Implementing `List`,
`Map`, or `Queryable` alone grants no debugger access to children.

When you pass a pointer to a native runtime scalar, the debugger displays the
pointed-to value. For example, a pointer to `runtime.Int(42)` displays `42`,
and a nil `*runtime.Int` displays `<Int>`. A pointer to `runtime.Binary`
displays its byte length. DateTime displays preserve the numeric timezone
offset, for example `2026-02-03T04:05:06.123456789-03:30`.

Descriptions are shallow. A positive reference lets you request a page from
the paused session. For example, to read elements 1000 through 1019 of a
buffered array, use the reference from its local or evaluation result.

{{< code lang="go" >}}
children, err := session.Variables(ctx, value.Reference, ferret.DebugVariablePage{
    Filter: ferret.DebugVariableFilterIndexed,
    Start:  1000,
    Count:  20,
})
{{</ code >}}

Each child includes its absolute index or exact property name. The default
page size is 100 and the maximum is 1,000. `IndexedCount` and `NamedCount` are
optional exact counts. A nil count is unknown. A zero reference leaves the
value available as a summary. Repeated reads of a binding or child reuse its
reference during the stop. After resume or close, retained page data remains
readable and its expansion references expire.

To expose safe children of a host value, implement
`runtime.DebugChildrenProvider`. Return exact known counts and readable views
from `DebugChildrenInfo()`. Implement `DebugChildAt(ctx, kind, position)` using
stable in-memory storage, observe cancellation, and return `found=false` at
the end. Keep reachable values and borrowed array backing slices stable
throughout the pause. A deferred query or cursor can use presentation alone.

Evaluating a member expression such as `items[0]` reads the selected child.
Use `ferret.WithDebugFormat` to set `MaxBytes`, the byte budget for each type and display.
The budget includes string quotes, escapes, and the `...` truncation marker.
A truncated field with only one or two bytes available uses one or two dots.

`Equatable` and `Comparable` are independent. `Equal` must be reflexive, symmetric, and transitive within the value's compatible comparison domain, and must return false without an error for an incompatible value. `Compare` must provide a reflexive, antisymmetric, and transitive order within that domain and return an error compatible with `runtime.ErrInvalidOperation` for an incompatible value. If both interfaces are implemented, `Equal` must return true exactly when `Compare` returns `runtime.Equal`.

Semantically equal values must return the same `Hash()`. A hash is only a candidate selector; Ferret verifies `Equal` before grouping or deduplicating host values.

Both methods receive the execution context unchanged and may return operational errors. Any negative or positive `Ordering` is normalized to `runtime.Less` or `runtime.Greater` by runtime dispatch.

`Unwrappable` is useful when other Go code (custom functions, modules) needs to access the underlying Go object. The `runtime.UnwrapAs` generic helper simplifies extraction:

{{< code lang="go" >}}
if db, ok := runtime.UnwrapAs[*sql.DB](arg); ok {
    // use the database connection
}
{{</ code >}}

Native `runtime.Int` unwraps to Go `int64` on every architecture. Use `runtime.UnwrapAs[int64]` when extracting it. To bind an Int into a Go `int`, use `sdk.Decode`, which checks whether the value fits the destination type before assigning it.

For code upgrading from native `int` assertions, see [Update integer unwrapping]({{< ref "/docs/migrations/v1-to-v2/go-embedding#update-integer-unwrapping" >}}).

When a Go API requires a native integer, use `runtime.ToNativeInt(value)`. It returns the exact Go `int` and `true` when the value fits, or `0, false` otherwise. Negative values are accepted; apply any index or size restrictions before calling the destination API.

For optional allocation sizing, use `runtime.CapacityHint(length, multiplier, extra)`. It checks `length*multiplier+extra` and returns zero when the result cannot fit, length or extra is negative, or multiplier is nonpositive. Zero means to omit the hint. This helper performs no allocation, does not validate mandatory sizes, and does not guarantee allocation success.

## Returning host values from functions

A common pattern is a namespaced function that opens a resource and returns it as a host value. The following example implements a minimal in-memory store that supports `query ... in`:

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "fmt"
    "hash/fnv"
    "log"
    "strings"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Record struct {
    Name string
    Age  int
}

type Store struct {
    records []Record
}

var StoreType = runtime.NewTypeFor[*Store]()

func NewStore() *Store {
    return &Store{
        records: []Record{
            {"Alice", 30},
            {"Bob", 17},
            {"Carol", 25},
        },
    }
}

func (s *Store) Type() runtime.Type { return StoreType }
func (s *Store) String() string     { return "Store" }

func (s *Store) Hash() uint64 {
    h := fnv.New64a()
    h.Write([]byte("store"))
    return h.Sum64()
}

func (s *Store) Copy() runtime.Value { return s }
func (s *Store) Close() error        { return nil }

func (s *Store) Query(_ context.Context, q runtime.Query) (runtime.List, error) {
    out := runtime.NewArray(len(s.records))

    for _, r := range s.records {
        if strings.Contains(strings.ToLower(q.Expression.String()), "age > 18") && r.Age <= 18 {
            continue
        }

        obj := runtime.NewObject()
        obj.SetString("name", runtime.NewString(r.Name))
        obj.SetString("age", runtime.NewInt(r.Age))
        out.Append(context.Background(), obj)
    }

    return out, nil
}

func (s *Store) QueryOne(ctx context.Context, q runtime.Query) (runtime.Value, error) {
    return runtime.DefaultQueryOne(ctx, q, s.Query)
}

func (s *Store) QueryCount(ctx context.Context, q runtime.Query) (runtime.Int, error) {
    return runtime.DefaultQueryCount(ctx, q, s.Query)
}

func (s *Store) QueryExists(ctx context.Context, q runtime.Query) (runtime.Boolean, error) {
    return runtime.DefaultQueryExists(ctx, q, s.Query)
}

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() (err error) {
    ctx := context.Background()

    engine, err := ferret.New(
        ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
            db := ns.Namespace("db")
            db.Function().A0().Add("open", func(ctx context.Context) (runtime.Value, error) {
                return NewStore(), nil
            })
        }),
    )
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, engine.Close()) }()

    plan, err := engine.Compile(ctx, ferret.NewSource("host-values.fql", `
        let db = db::open()
        return query "SELECT * WHERE age > 18" in db
    `))
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, plan.Close()) }()

    session, err := plan.NewSession(ctx)
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, session.Close()) }()

    output, err := session.Run(ctx)
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, output.Close()) }()

    content, err := output.Collect(ctx)
    if err != nil {
        return err
    }

    fmt.Println(string(content.Data))
    // [{"name":"Alice","age":30},{"name":"Carol","age":25}]

    return nil
}
{{</ code >}}

## Passing host values as parameters

Instead of creating a host value inside a function, you can pass one directly as a parameter with `WithRuntimeParam` or `WithSessionRuntimeParam`. A type that implements `runtime.Value` also satisfies the root `ferret.Value` alias accepted by these options:

{{< code lang="go" >}}
store := NewStore()

engine, err := ferret.New(
    ferret.WithRuntimeParam("db", store),
)
{{</ code >}}

{{< code lang="fql" >}}
return query "SELECT * WHERE age > 18" in @db
{{</ code >}}

This is useful when the host application controls the resource lifecycle and wants to share a single instance across multiple query executions. The engine-level parameter is available to every session without re-creation.

For per-session values, use `WithSessionRuntimeParam`:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionRuntimeParam("db", store),
)
{{</ code >}}

## Cleanup and finalization

For a VM-owned value that implements `io.Closer`, Core tracks a comparable Go value (including a pointer) by its identity. A non-comparable closer needs to implement `runtime.Resource` to participate in VM ownership tracking.

1. During execution, the VM tracks resources it owns. It transfers escaping resources to the result. The encoder can also discover nested resources during materialization.
2. After encoding the output, the runtime calls `Close()` on every tracked closer.
3. Result cleanup errors are retained by the returned `Output` and reported by `Collect` or `Consume`. `Output.Close` also retains the cleanup outcome; it does not report execution or encoding failures.

When a script returns a tracked resource created by a host function, Native encodes the value and closes it before `Run` returns. The caller receives a handle to the encoded content. Consume or collect the handle to observe execution, encoding, and cleanup failures. Direct calls to a codec borrow values and leave source and child-resource cleanup to the caller. See [Value Encoders]({{< ref "/docs/embedding/go/value-encoders#context-and-buffering" >}}).

Result materialization can also adopt a non-comparable closer. Without a stable identity, repeated discoveries may close it more than once. Implement `runtime.Resource` for such values or when different Go values represent the same underlying resource.

{{< code lang="go" >}}
type Resource interface {
    io.Closer
    ResourceID() uint64
}
{{</ code >}}

`ResourceID` must return a stable, unique identifier for the resource. The runtime uses it to deduplicate closers — two references to the same resource result in a single `Close` call.

A parameter is borrowed during evaluation. If a script returns `@db` or a container holding it, Native encoding adopts the closer and result cleanup closes it before `Run` returns. To share a connection across sessions, return the queried records and keep the handle out of the result. The host closes the shared connection after its last use.

## Context and cancellation

The VM observes query cancellation at structural execution boundaries, such as loop backedges, calls, debugger observations, and top-level completion. Native arithmetic, comparison, property access, and collection work between those boundaries is atomic from the VM's perspective.

Context-aware host functions and capabilities receive the original execution context. An implementation that waits, performs network or database I/O, advances a remote iterator, or otherwise retains control must observe that context itself:

{{< code lang="go" >}}
func (s *Store) Query(ctx context.Context, q runtime.Query) (runtime.List, error) {
    select {
    case result := <-s.query(q):
        return result, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
{{</ code >}}

Return cancellation with `%w` when adding context. Errors compatible with `context.Canceled` or `context.DeadlineExceeded` terminate execution directly and cannot be caught by `on error`, retry recovery, or another protected FQL operation.

`io.Closer` has no context parameter. Its `Close` call is treated as one atomic operation.

## Error behavior

When a capability method returns a non-cancellation error, the VM raises an FQL runtime error at the point of the operation. The error message includes the value's type and the Go error text. FQL scripts can recover from eligible errors with `on error ...`. Cancellation and deadline errors always propagate to the embedding caller. See [Error Handling]({{< ref "docs/language/control-flow/error-handling" >}}).

When the script attempts an operation that requires a capability the value does not implement (for example, `for x in val` on a non-`Iterable` value), the VM raises a type error listing the expected capability.

Use the runtime error helpers to return well-formatted errors from your capability methods:

{{< code lang="go" >}}
runtime.Error(runtime.ErrNotFound, "record does not exist")
runtime.TypeError(runtime.TypeOf(val), runtime.TypeIterable)
{{</ code >}}

## The Proxy shortcut

For cases where a Go type already implements some capability interfaces, `sdk.Proxy[T]` avoids the boilerplate of implementing `Value` manually. The proxy wraps any Go value and delegates capability calls via type assertion:

{{< code lang="go" >}}
import "github.com/MontFerret/ferret/v2/pkg/sdk"

proxy := sdk.NewProxy[*MyDB](myDB)
{{</ code >}}

If `*MyDB` implements `Queryable`, calls to `proxy.Query` delegate to `myDB.Query`. If it does not implement `Iterable`, `proxy.Iterate` returns a typed error. The proxy also handles `Value` methods (`String`, `Hash`, `Copy`), `Typed`, `Unwrappable`, `json.Marshaler`, and equality delegation automatically.

For Go maps and slices, use the specialized variants:

{{< code lang="go" >}}
proxyMap := sdk.NewProxyMap[string, *User](users)
proxySlice := sdk.NewProxySlice[*Record](records)
{{</ code >}}

`ProxyMap` automatically implements `KeyReadable`, `KeyWritable`, and `Iterable` for `map[string]V` types. `ProxySlice` implements `IndexReadable`, `IndexWritable`, and `Sortable`.

## Next steps

{{< docs-related tiles="language-types-host,language-types-capabilities,embedding-go-parameters,embedding-go-modules" >}}
