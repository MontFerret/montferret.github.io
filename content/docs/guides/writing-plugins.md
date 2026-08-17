---
title: "Develop a Ferret module"
sidebarTitle: "Developing modules"
weight: 120
draft: false
description: "Build a self-contained module with namespaced functions, host values, and lifecycle hooks."
---

# Develop a Ferret module

A Ferret module bundles namespaced functions, host values, and lifecycle hooks into a single registerable Go package. This guide builds a complete module — a key-value cache that FQL scripts can create, populate, read, and iterate.

For the registration and lifecycle model, see [Modules]({{< ref "/docs/embedding/go/modules" >}}). For the underlying extension APIs, see [Custom Functions]({{< ref "/docs/embedding/go/custom-functions" >}}) and [Host Values]({{< ref "/docs/embedding/go/host-values" >}}).

## What the module will do

The `kv` module exposes three functions and a host value:

{{< code lang="fql" >}}
let cache = kv::open()

kv::set(cache, "language", "FQL")
kv::set(cache, "version", 2)

return {
    language: kv::get(cache, "language"),
    size:     length(cache),
    keys:     (for key in cache return key)
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
  --namespace kv
cd kvplugin
go mod tidy
{{</ terminal >}}

The initializer creates the module manifest, Go module, starter registration, documentation, and package directories. Replace the starter registration and add the implementation and test files as you follow the guide:

```text
kvplugin/
├── ferret.yaml
├── go.mod
├── module.go
├── module_test.go
├── options.go
├── README.md
├── core/
│   ├── cache.go
│   └── doc.go
└── lib/
    ├── doc.go
    └── lib.go
```

This structure keeps three responsibilities separate:

- the root package composes configuration, registration, and lifecycle hooks
- `core` implements the cache without knowing how FQL functions are registered
- `lib` owns the Ferret-facing function boundary

The generated `ferret.yaml` is schema-valid but contains TODO release metadata. See [Develop a module project]({{< ref "/docs/modules/develop" >}}) for the scaffold contract.

## Define configuration options

The constructor will accept functional options and validate them when Ferret registers the module. Start with a default capacity of 1,000 entries:

{{< code lang="go" title="options.go" >}}
package kvplugin

import "fmt"

const defaultMaxSize = 1000

type config struct {
    maxSize int
}

type Option func(*config) error

func WithMaxSize(maxSize int) Option {
    return func(config *config) error {
        if maxSize <= 0 {
            return fmt.Errorf("max size must be greater than zero")
        }

        config.maxSize = maxSize
        return nil
    }
}

func resolveConfig(setters []Option) (config, error) {
    config := config{maxSize: defaultMaxSize}

    for _, set := range setters {
        if set == nil {
            continue
        }

        if err := set(&config); err != nil {
            return config, err
        }
    }

    return config, nil
}
{{</ code >}}

`New` returns `module.Module`, so configuration errors surface when the host builds its Ferret engine. This keeps the constructor compatible with `ferret.WithModules(...)` while still rejecting invalid options before any functions are registered.

## Define the module with the SDK

Use `sdk.NewModule` to define the module name and registration callback. The callback resolves configuration, registers the `kv` library, and attaches lifecycle hooks:

{{< code lang="go" title="module.go" >}}
package kvplugin

import (
    "context"
    "log"
    "time"

    "github.com/acme/ferret-kvplugin/lib"

    "github.com/MontFerret/ferret/v2/pkg/module"
    "github.com/MontFerret/ferret/v2/pkg/sdk"
)

func New(setters ...Option) module.Module {
    options := append([]Option(nil), setters...)

    return sdk.NewModule("acme/kvplugin", func(bootstrap module.Bootstrap) error {
        config, err := resolveConfig(options)
        if err != nil {
            return err
        }

        namespace := bootstrap.Host().Library().Namespace("kv")
        if err := lib.RegisterLib(namespace, config.maxSize); err != nil {
            return err
        }

        registerHooks(bootstrap)
        return nil
    })
}

type runStartKey struct{}

func registerHooks(bootstrap module.Bootstrap) {
    bootstrap.Hooks().Session().BeforeRun(func(ctx context.Context) (context.Context, error) {
        return context.WithValue(ctx, runStartKey{}, time.Now()), nil
    })

    bootstrap.Hooks().Session().AfterRun(func(ctx context.Context, runErr error) error {
        start, _ := ctx.Value(runStartKey{}).(time.Time)
        log.Printf("[kv] query took %s (err=%v)", time.Since(start), runErr)
        return nil
    })
}
{{</ code >}}

`sdk.NewModule` supplies the `module.Module` implementation and adds the module name to registration errors. Keep registration itself in the callback so an invalid configuration or library definition prevents the host engine from starting with a partially registered module.

The hook functions remain normal Ferret lifecycle hooks. A context returned by `BeforeRun` flows into the VM and the matching `AfterRun` hook, which makes it suitable for request-scoped state such as tracing spans.

## Register namespaced functions

Add `kv::open`, `kv::set`, and `kv::get` with the SDK's declarative registration helpers:

{{< code lang="go" title="lib/lib.go" >}}
package lib

import (
    "context"

    "github.com/acme/ferret-kvplugin/core"

    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/sdk"
)

func RegisterLib(namespace runtime.Namespace, maxSize int) error {
    return sdk.RegisterFunctions(
        namespace,
        sdk.Func("OPEN", openWithMaxSize(maxSize)),
        sdk.Func("SET", sdk.Bind3(Set)),
        sdk.Func("GET", sdk.Bind2(Get)),
    )
}

func openWithMaxSize(maxSize int) runtime.Function0 {
    return sdk.Bind0(func(context.Context) (*core.Cache, error) {
        return core.NewCache(maxSize), nil
    })
}

func Set(
    _ context.Context,
    cache *core.Cache,
    key runtime.String,
    value runtime.Value,
) (runtime.Value, error) {
    if err := cache.SetValue(string(key), value); err != nil {
        return runtime.None, err
    }

    return runtime.None, nil
}

func Get(
    _ context.Context,
    cache *core.Cache,
    key runtime.String,
) (runtime.Value, error) {
    value, found := cache.GetValue(string(key))
    if !found {
        return runtime.None, nil
    }

    return value, nil
}
{{</ code >}}

`sdk.RegisterFunctions` validates the complete definition set before changing the namespace. A duplicate name and arity, nil handler, or invalid definition therefore cannot leave a partially registered library.

`sdk.Bind0` through `sdk.Bind4` adapt fixed-arity functions whose arguments and results already implement `runtime.Value`. Here the binders validate `*core.Cache` and `runtime.String` arguments and attach the correct argument position to type errors.

For a variadic function, validate its arity with `runtime.ValidateArgs`, then use `sdk.DecodeArg` for required arguments and `sdk.DecodeArgOr` for optional arguments. Those decoding helpers also support native Go values and structured options, but they are unnecessary for these fixed-arity handlers.

## Implement the cache host value

The cache implements several capability interfaces so FQL scripts can interact with it naturally:

| Interface | Enables |
| --- | --- |
| `runtime.Value` | Hold and pass the value in scripts |
| `runtime.Typed` | Type name in error messages |
| `runtime.KeyReadable` | `cache.key` property access |
| `runtime.Iterable` | `for key in cache` |
| `runtime.Measurable` | `length(cache)` |
| `io.Closer` | Automatic cleanup when the session ends |

{{< code lang="go" title="core/cache.go" >}}
package core

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
    mu      sync.RWMutex
    items   map[string]runtime.Value
    maxSize int
}

func NewCache(maxSize int) *Cache {
    return &Cache{
        items:   make(map[string]runtime.Value),
        maxSize: maxSize,
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
    _, _ = h.Write([]byte("kv-cache"))
    return h.Sum64()
}

func (c *Cache) Copy() runtime.Value {
    return c
}

// --- Cache operations ---

func (c *Cache) SetValue(key string, value runtime.Value) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if _, exists := c.items[key]; !exists && len(c.items) >= c.maxSize {
        return fmt.Errorf("cache has reached its limit of %d entries", c.maxSize)
    }

    c.items[key] = value
    return nil
}

func (c *Cache) GetValue(key string) (runtime.Value, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, found := c.items[key]
    return value, found
}

// --- KeyReadable: cache.key ---

func (c *Cache) Get(_ context.Context, key runtime.Value) (runtime.Value, error) {
    value, found := c.GetValue(key.String())
    if !found {
        return runtime.None, nil
    }

    return value, nil
}

// --- Measurable: length(cache) ---

func (c *Cache) Length(_ context.Context) (runtime.Int, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return runtime.Int(len(c.items)), nil
}

// --- Iterable: FOR key IN cache ---

func (c *Cache) Iterate(_ context.Context) (runtime.Iterator, error) {
    c.mu.RLock()
    keys := make([]string, 0, len(c.items))
    for key := range c.items {
        keys = append(keys, key)
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

The limit applies only when adding a new key. Updating an existing key is allowed when the cache is full, which makes capacity predictable without preventing normal replacement.

## Register and use the module

The host imports only the root package. It does not need to know about `core` or `lib`:

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
        ferret.WithModules(kvplugin.New(
            kvplugin.WithMaxSize(500),
        )),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    output, err := engine.Run(
        context.Background(),
        source.NewAnonymous(`
            let cache = kv::open()

            kv::set(cache, "language", "FQL")
            kv::set(cache, "version", 2)

            return {
                language: kv::get(cache, "language"),
                size:     length(cache),
                keys:     (for key in cache return key)
            }
        `),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
}
{{</ code >}}

## Test through the SDK harness

Use `sdktest` for black-box tests that exercise registration, compilation, runtime argument checks, lifecycle hooks, and result encoding together:

{{< code lang="go" title="module_test.go" >}}
package kvplugin_test

import (
    "encoding/json"
    "testing"

    kvplugin "github.com/acme/ferret-kvplugin"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/sdk/sdktest"
)

func newHarness(t *testing.T, options ...kvplugin.Option) *sdktest.Harness {
    t.Helper()
    return sdktest.New(t, ferret.WithModules(kvplugin.New(options...)))
}

func TestCache(t *testing.T) {
    harness := newHarness(t)

    output, err := harness.Run(t.Context(), `
        let cache = kv::open()
        kv::set(cache, "a", 1)
        kv::set(cache, "b", 2)
        return {
            size: length(cache),
            a: kv::get(cache, "a"),
            keys: (for key in cache return key)
        }
    `)
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

    if result.Size != 2 || result.A != 1 || len(result.Keys) != 2 {
        t.Fatalf("unexpected result: %+v", result)
    }
}

func TestCacheRejectsWrongArgumentType(t *testing.T) {
    harness := newHarness(t)

    if _, err := harness.Run(t.Context(), `return kv::get("not a cache", "key")`); err == nil {
        t.Fatal("expected a cache argument error")
    }
}

func TestCacheCapacity(t *testing.T) {
    harness := newHarness(t, kvplugin.WithMaxSize(1))

    output, err := harness.Run(t.Context(), `
        let cache = kv::open()
        kv::set(cache, "key", 1)
        kv::set(cache, "key", 2)
        return kv::get(cache, "key")
    `)
    if err != nil {
        t.Fatal(err)
    }
    if string(output.Content) != "2" {
        t.Fatalf("expected overwrite result 2, got %s", output.Content)
    }

    if _, err := harness.Run(t.Context(), `
        let cache = kv::open()
        kv::set(cache, "first", 1)
        kv::set(cache, "second", 2)
        return true
    `); err == nil {
        t.Fatal("expected the second key to exceed capacity")
    }
}

func TestWithMaxSizeRejectsNonPositiveValue(t *testing.T) {
    engine, err := ferret.New(
        ferret.WithModules(kvplugin.New(kvplugin.WithMaxSize(0))),
    )
    if engine != nil {
        _ = engine.Close()
    }
    if err == nil {
        t.Fatal("expected module configuration to fail")
    }
}
{{</ code >}}

Run the complete module test suite:

{{< terminal command="true" >}}
go test ./...
{{</ terminal >}}

## Next steps

{{< docs-related tiles="runtime-modules-develop,runtime-modules-publish,embedding-go-modules,embedding-go-custom-functions" >}}
