---
title: "Host Values"
sidebarTitle: "Host Values"
weight: 35
draft: false
description: "Implement Go types that FQL scripts can hold, access, iterate, query, and dispatch to."
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
| Function | Stateless computation: hash, format, fetch | `FN(args)` → result | No state between calls |
| Host value | Stateful resource or behavioral object | Variable with capabilities: `val.prop`, `FOR x IN val`, `QUERY ... IN val` | Host creates; runtime tracks `io.Closer` |
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

- `String()` — text representation used by `TO_STRING()` and encoding fallbacks.
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

## Minimal implementation

The smallest possible host value — a struct that wraps a string:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "hash/fnv"
    "log"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/source"
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

Register a function that returns it:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
        ns.Function().A1().Add("LABEL", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
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
LET l = LABEL("urgent")
RETURN TO_STRING(l)
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
| `KeyRemovable` | `RemoveKey(ctx, key Value) error` | `DELETE val.prop`         |
| `IndexReadable` | `At(ctx, idx Int) (Value, error)` | `val[0]`                 |
| `IndexLookup` | `LookupAt(ctx, index Int) (Value, bool, error)` | `val?[0]`                |
| `IndexWritable` | `SetAt(ctx, idx Int, value Value) error` | `val[0] = x`             |
| `IndexRemovable` | `RemoveAt(ctx, idx Int) (Value, error)` | `DELETE val[0]`          |

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
LET l = LABEL("urgent")
RETURN l.text
// "urgent"
{{</ code >}}

### Iteration and measurement

| Interface | Method | FQL syntax |
|-----------|--------|------------|
| `Iterable` | `Iterate(ctx) (Iterator, error)` | `FOR x IN val` |
| `Iterator` | `Next(ctx) (value, key Value, err error)` | (returned by `Iterate`) |
| `Measurable` | `Length(ctx) (Int, error)` | `LENGTH(val)` |
| `Containable` | `Contains(ctx, value Value) (Boolean, error)` | `x IN val` |

`Iterator.Next` must return `io.EOF` when the sequence is exhausted. If the iterator holds resources (an open cursor, for example), implement `io.Closer` on it as well.

### Query, dispatch, and observe

| Interface | Key method | FQL syntax |
|-----------|------------|------------|
| `Queryable` | `Query(ctx, Query) (List, error)` + 3 modifiers | `QUERY ... IN val USING ...` |
| `Dispatchable` | `Dispatch(ctx, DispatchEvent) error` | `DISPATCH "event" TO val` / `val <- event` |
| `Observable` | `Subscribe(ctx, Subscription) (Stream, error)` | `WAITFOR EVENT "name" IN val` |

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

### Comparison and other capabilities

| Interface | Method | Purpose |
|-----------|--------|---------|
| `Comparable` | `Compare(other Value) int` | Enables `==`, `<`, `>`, and `SORT` |
| `Sortable` | `SortAsc(ctx) error`, `SortDesc(ctx) error` | In-place sort |
| `Cloneable` | `Clone(ctx) (Cloneable, error)` | Deep copy |
| `Unwrappable` | `Unwrap() any` | Extract the inner Go value |
| `DebugInspectable` | `DebugInfo() DebugInfo` | Presentation hints for the debugger |

`Unwrappable` is useful when other Go code (custom functions, modules) needs to access the underlying Go object. The `runtime.UnwrapAs` generic helper simplifies extraction:

{{< code lang="go" >}}
if db, ok := runtime.UnwrapAs[*sql.DB](arg); ok {
    // use the database connection
}
{{</ code >}}

## Returning host values from functions

A common pattern is a namespaced function that opens a resource and returns it as a host value. The following example implements a minimal in-memory store that supports `QUERY ... IN`:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "hash/fnv"
    "io"
    "log"
    "strings"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/source"
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
        obj.Set(context.Background(), runtime.NewString("name"), runtime.NewString(r.Name))
        obj.Set(context.Background(), runtime.NewString("age"), runtime.NewInt(r.Age))
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
    ctx := context.Background()

    engine, err := ferret.New(
        ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
            db := ns.Namespace("DB")
            db.Function().A0().Add("OPEN", func(ctx context.Context) (runtime.Value, error) {
                return NewStore(), nil
            })
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    plan, err := engine.Compile(ctx, source.NewAnonymous(`
        LET db = DB::OPEN()
        RETURN QUERY "SELECT * WHERE age > 18" IN db
    `))
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    session, err := plan.NewSession(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    output, err := session.Run(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
    // [{"age":25,"name":"Carol"},{"age":30,"name":"Alice"}]
}
{{</ code >}}

## Passing host values as parameters

Instead of creating a host value inside a function, you can pass one directly as a parameter with `WithRuntimeParam` or `WithSessionRuntimeParam`:

{{< code lang="go" >}}
store := NewStore()

engine, err := ferret.New(
    ferret.WithRuntimeParam("db", store),
)
{{</ code >}}

{{< code lang="fql" >}}
RETURN QUERY "SELECT * WHERE age > 18" IN @db
{{</ code >}}

This is useful when the host application controls the resource lifecycle and wants to share a single instance across multiple query executions. The engine-level parameter is available to every session without re-creation.

For per-session values, use `WithSessionRuntimeParam`:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionRuntimeParam("db", store),
)
{{</ code >}}

## Cleanup and finalization

The runtime automatically tracks values that implement `io.Closer`:

1. During execution, any `io.Closer` value encountered in the result is added to a tracked closer set.
2. After encoding the output, the runtime calls `Close()` on every tracked closer.
3. Errors from `Close()` are propagated to the caller of `Session.Run`.

This means a host value returned by a function inside a script is encoded first, then closed. The host receives the encoded output; the value is already closed by the time `Run` returns.

For values referenced multiple times, implement the `Resource` interface to prevent double-closing:

{{< code lang="go" >}}
type Resource interface {
    io.Closer
    ResourceID() uint64
}
{{</ code >}}

`ResourceID` must return a stable, unique identifier for the resource. The runtime uses it to deduplicate closers — two references to the same resource result in a single `Close` call.

Values stored in engine-level parameters are **not** tracked by the result's closer set. The host owns their lifecycle and must close them explicitly when the engine shuts down.

## Error behavior

When a capability method returns a non-nil error, the VM raises an FQL runtime error at the point of the operation. The error message includes the value's type and the Go error text. FQL scripts can recover from these errors with `ON ERROR ...`. See [Error Handling]({{< ref "docs/language/control-flow/error-handling" >}}).

When the script attempts an operation that requires a capability the value does not implement (for example, `FOR x IN val` on a non-`Iterable` value), the VM raises a type error listing the expected capability.

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

If `*MyDB` implements `Queryable`, calls to `proxy.Query` delegate to `myDB.Query`. If it does not implement `Iterable`, `proxy.Iterate` returns a typed error. The proxy also handles `Value` methods (`String`, `Hash`, `Copy`), `Typed`, `Unwrappable`, `json.Marshaler`, and `Comparable` delegation automatically.

For Go maps and slices, use the specialized variants:

{{< code lang="go" >}}
proxyMap := sdk.NewProxyMap[string, *User](users)
proxySlice := sdk.NewProxySlice[*Record](records)
{{</ code >}}

`ProxyMap` automatically implements `KeyReadable`, `KeyWritable`, and `Iterable` for `map[string]V` types. `ProxySlice` implements `IndexReadable`, `IndexWritable`, and `Sortable`.

## Next steps

{{< docs-related tiles="language-types-host,language-types-capabilities,embedding-parameters,embedding-modules" >}}
