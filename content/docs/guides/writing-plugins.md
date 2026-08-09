---
title: "Develop a Ferret module"
sidebarTitle: "Developing modules"
weight: 120
draft: false
description: "Build a self-contained module with namespaced functions, host values, and lifecycle hooks."
---

# Develop a Ferret module

A Ferret module bundles namespaced functions, host values, and lifecycle hooks into a single registerable Go package. This guide builds a complete module — a key-value cache that FQL scripts can create, populate, read, and iterate.

For the registration and lifecycle model, see [Modules]({{< ref "/docs/embedding/modules" >}}). For the underlying extension APIs, see [Custom Functions]({{< ref "/docs/embedding/custom-functions" >}}) and [Host Values]({{< ref "/docs/embedding/host-values" >}}).

## What the module will do

The `KV` module exposes three functions and a host value:

{{< code lang="fql" >}}
LET cache = KV::OPEN()

KV::SET(cache, "language", "FQL")
KV::SET(cache, "version", 2)

RETURN {
    language: KV::GET(cache, "language"),
    size:     LENGTH(cache),
    keys:     (FOR key IN cache RETURN key)
}
{{</ code >}}

Expected output:

```json
{"language": "FQL", "keys": ["language", "version"], "size": 2}
```

## Scaffold the module

{{< terminal command="true" >}}
ferret mod init acme/kvplugin \
  --go-module github.com/acme/ferret-kvplugin \
  --dir kvplugin \
  --namespace KV
cd kvplugin
go mod tidy
{{</ terminal >}}

The initializer creates the module manifest, Go module, starter registration, documentation, and package directories. Add `cache.go` as you follow the guide, and replace the starter `module.go` with the implementation below:

```text
kvplugin/
├── ferret.yaml
├── go.mod
├── module.go
├── cache.go
├── README.md
├── core/
│   └── doc.go
└── lib/
    └── doc.go
```

The generated `ferret.yaml` is schema-valid but contains TODO release metadata. See [Develop a module project]({{< ref "/docs/modules/develop" >}}) for the scaffold contract.

## Implement the Module interface

A module implements `Name()` and `Register(Bootstrap)`:

{{< code lang="go" title="module.go" >}}
package kvplugin

import (
    "context"
    "log"
    "time"

    "github.com/MontFerret/ferret/v2/pkg/module"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Module struct{}

func New() *Module {
    return &Module{}
}

func (m *Module) Name() string {
    return "acme/kvplugin"
}

func (m *Module) Register(boot module.Bootstrap) error {
    ns := boot.Host().Library().Namespace("KV")

    m.registerFunctions(ns)
    m.registerHooks(boot)

    return nil
}
{{</ code >}}

## Register namespaced functions

Add `KV::OPEN`, `KV::SET`, and `KV::GET` to the namespace:

{{< code lang="go" title="module.go (continued)" >}}
func (m *Module) registerFunctions(ns runtime.Namespace) {
    // KV::OPEN() — returns a new cache host value
    ns.Function().A0().Add("OPEN", func(ctx context.Context) (runtime.Value, error) {
        return NewCache(), nil
    })

    // KV::SET(cache, key, value) — stores a value
    ns.Function().A3().Add("SET", func(ctx context.Context, cacheArg, keyArg, valArg runtime.Value) (runtime.Value, error) {
        cache, err := runtime.CastArg[*Cache](cacheArg, 0)
        if err != nil {
            return nil, err
        }

        key, err := runtime.CastArg[runtime.String](keyArg, 1)
        if err != nil {
            return nil, err
        }

        cache.SetValue(string(key), valArg)

        return runtime.None, nil
    })

    // KV::GET(cache, key) — retrieves a value
    ns.Function().A2().Add("GET", func(ctx context.Context, cacheArg, keyArg runtime.Value) (runtime.Value, error) {
        cache, err := runtime.CastArg[*Cache](cacheArg, 0)
        if err != nil {
            return nil, err
        }

        key, err := runtime.CastArg[runtime.String](keyArg, 1)
        if err != nil {
            return nil, err
        }

        val, found := cache.GetValue(string(key))
        if !found {
            return runtime.None, nil
        }

        return val, nil
    })
}
{{</ code >}}

`runtime.CastArg[T]` validates the argument type and returns a clear error message if the cast fails. The second argument is the parameter index for the error message.

## Implement the cache host value

The cache implements several capability interfaces so FQL scripts can interact with it naturally:

| Interface | Enables |
| --- | --- |
| `runtime.Value` | Hold and pass the value in scripts |
| `runtime.Typed` | Type name in error messages |
| `runtime.KeyReadable` | `cache.key` property access |
| `runtime.Iterable` | `FOR key IN cache` |
| `runtime.Measurable` | `LENGTH(cache)` |
| `io.Closer` | Automatic cleanup when the session ends |

{{< code lang="go" title="cache.go" >}}
package kvplugin

import (
    "context"
    "fmt"
    "hash/fnv"
    "sort"
    "sync"

    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/sdk"
)

var CacheType = runtime.NewTypeFor[*Cache]()

type Cache struct {
    mu    sync.RWMutex
    items map[string]runtime.Value
}

func NewCache() *Cache {
    return &Cache{
        items: make(map[string]runtime.Value),
    }
}

// --- runtime.Value ---

func (c *Cache) Type() runtime.Type { return CacheType }

func (c *Cache) String() string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return fmt.Sprintf("Cache(%d)", len(c.items))
}

func (c *Cache) Hash() uint64 {
    h := fnv.New64a()
    h.Write([]byte("kv-cache"))
    return h.Sum64()
}

func (c *Cache) Copy() runtime.Value {
    return c
}

// --- Cache operations ---

func (c *Cache) SetValue(key string, val runtime.Value) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}

func (c *Cache) GetValue(key string) (runtime.Value, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

// --- KeyReadable: cache.key ---

func (c *Cache) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    val, ok := c.items[key.String()]
    if !ok {
        return runtime.None, nil
    }

    return val, nil
}

// --- Measurable: LENGTH(cache) ---

func (c *Cache) Length(_ context.Context) (runtime.Int, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return runtime.Int(len(c.items)), nil
}

// --- Iterable: FOR key IN cache ---

func (c *Cache) Iterate(_ context.Context) (runtime.Iterator, error) {
    c.mu.RLock()
    keys := make([]string, 0, len(c.items))
    for k := range c.items {
        keys = append(keys, k)
    }
    c.mu.RUnlock()

    sort.Strings(keys)

    return sdk.NewSliceIterator(keys), nil
}

// --- io.Closer ---

func (c *Cache) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items = nil
    return nil
}
{{</ code >}}

## Add lifecycle hooks

Hooks let the module react to engine events. Add timing and cleanup:

{{< code lang="go" title="module.go (continued)" >}}
type ctxStartKey struct{}

func (m *Module) registerHooks(boot module.Bootstrap) {
    boot.Hooks().Session().BeforeRun(func(ctx context.Context) (context.Context, error) {
        return context.WithValue(ctx, ctxStartKey{}, time.Now()), nil
    })

    boot.Hooks().Session().AfterRun(func(ctx context.Context, runErr error) error {
        start, _ := ctx.Value(ctxStartKey{}).(time.Time)
        log.Printf("[kv] query took %s (err=%v)", time.Since(start), runErr)
        return nil
    })
}
{{</ code >}}

`BeforeRun` hooks can return a modified context that flows through to `AfterRun` and into the VM. This is useful for injecting request-scoped values like tracing spans.

## Register and use the module

In the host application:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"

    kvplugin "github.com/acme/ferret-kvplugin"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    engine, err := ferret.New(
        ferret.WithModules(kvplugin.New()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        source.NewAnonymous(`
            LET cache = KV::OPEN()

            KV::SET(cache, "language", "FQL")
            KV::SET(cache, "version", 2)

            RETURN {
                language: KV::GET(cache, "language"),
                size:     LENGTH(cache),
                keys:     (FOR key IN cache RETURN key)
            }
        `),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
}
{{</ code >}}

## Accept configuration

Use the functional options pattern to make the module configurable:

{{< code lang="go" >}}
type Option func(*Module)

func WithMaxSize(n int) Option {
    return func(m *Module) {
        m.maxSize = n
    }
}

type Module struct {
    maxSize int
}

func New(opts ...Option) *Module {
    m := &Module{
        maxSize: 1000,
    }
    for _, opt := range opts {
        opt(m)
    }
    return m
}
{{</ code >}}

Then in the host:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithModules(kvplugin.New(
        kvplugin.WithMaxSize(500),
    )),
)
{{</ code >}}

## Test the module

Write a Go test that creates an engine with the module, runs a query, and asserts on the output:

{{< code lang="go" title="module_test.go" >}}
package kvplugin_test

import (
    "context"
    "encoding/json"
    "testing"

    kvplugin "github.com/acme/ferret-kvplugin"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func TestCache(t *testing.T) {
    engine, err := ferret.New(
        ferret.WithModules(kvplugin.New()),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        source.NewAnonymous(`
            LET c = KV::OPEN()
            KV::SET(c, "a", 1)
            KV::SET(c, "b", 2)
            RETURN {
                size: LENGTH(c),
                a: KV::GET(c, "a"),
                keys: (FOR k IN c RETURN k)
            }
        `),
    )
    if err != nil {
        t.Fatal(err)
    }

    var result struct {
        Size int      `json:"size"`
        A    int      `json:"a"`
        Keys []string `json:"keys"`
    }

    if err := json.Unmarshal(output.Content, &result); err != nil {
        t.Fatal(err)
    }

    if result.Size != 2 {
        t.Errorf("expected size 2, got %d", result.Size)
    }

    if result.A != 1 {
        t.Errorf("expected a=1, got %d", result.A)
    }

    if len(result.Keys) != 2 {
        t.Errorf("expected 2 keys, got %d", len(result.Keys))
    }
}
{{</ code >}}

{{< terminal command="true" >}}
go test ./...
{{</ terminal >}}

## Next steps

{{< docs-related tiles="runtime-modules-develop,runtime-modules-publish,embedding-modules,embedding-custom-functions" >}}
