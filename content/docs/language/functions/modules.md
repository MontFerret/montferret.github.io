---
title: "Module Functions"
sidebarTitle: "Modules"
weight: 20
draft: false
description: "Call namespaced functions and create aliases with the use statement."
---

# Module functions

Namespaces organize related runtime-provided functions. The Standard Library, modules, and embedding applications can register namespaced functions. Once a provider is registered, its functions can be called by their qualified names without an import statement.

## Namespaced calls

A namespaced function is called by prefixing the function name with the namespace and `::`.

{{< code lang="fql" >}}
return io::fs::read("/tmp/data.json")
{{</ code >}}

Namespaces can be nested. Each `::` separator introduces another level.

{{< code lang="fql" >}}
return io::net::http::get("https://api.example.com/data")
{{</ code >}}

The namespace groups related functions and helps avoid name conflicts between providers. These examples require a runtime with the relevant I/O functions and access to the file or endpoint.

## The use statement

When a script repeatedly calls functions from the same namespace, `use` can create a shorter local alias. Put the declaration at the start of the script, before the script body.

{{< code lang="fql" >}}
use io::fs as fs

return fs::read("/tmp/data.json")
{{</ code >}}

`use` is a compile-time alias; it does not load a module or change which functions the runtime provides. See the [`use` statement reference]({{< ref "/docs/language/script-structure/use" >}}) for function aliases, placement rules, and name-resolution behavior.

## Where functions come from

The set of available namespaced functions depends on the host environment. Functions are registered at engine startup by:

- **The [Standard Library]({{< ref "/docs/standard-library" >}})** supplies both global functions, such as `length`, and namespaced functions, such as `math::floor`.
- **[Modules]({{< ref "/docs/modules" >}})** register additional functions. The [Registry]({{< ref "/registry" >}}) links to module APIs and documentation.
- **The host application** can register its own functions and namespaces directly when embedding Ferret.

The language does not distinguish between these sources. A namespaced function call works the same way regardless of where the function was registered.

## Runtime-backed functions

Some namespaced functions interact with external systems. They may read files, issue HTTP requests, query databases, control browsers, or work with binary data.

The behavior of such functions depends on the runtime configuration. The same query may behave differently in a CLI context than in a browser automation runtime — not because the language changes, but because the set of available functions and runtime capabilities differs.

## Next steps

{{< docs-related tiles="language-use,runtime-modules,registry,embedding-go-modules,stdlib,language-functions" >}}
