---
title: "Module Functions"
sidebarTitle: "Modules"
weight: 20
draft: false
description: "Call namespaced functions and create aliases with the USE statement."
---

# Module functions

In addition to built-in functions, which are available at the top level, Ferret supports functions organized in namespaces. These namespaced functions are provided by the host application, modules, or extensions registered at engine startup. They are always available by their full qualified name — no import step is needed.

## Namespaced calls

A namespaced function is called by prefixing the function name with the namespace and `::`.

{{< code lang="fql" >}}
IO::FS::READ("/tmp/data.json")
{{</ code >}}

Namespaces can be nested. Each `::` separator introduces another level.

{{< code lang="fql" >}}
IO::NET::HTTP::GET("https://api.example.com/data")
{{</ code >}}

The namespace makes it clear where the function comes from and avoids name conflicts between different providers.

## The USE statement

The `USE` statement creates a compile-time alias for a namespace or a specific namespaced function. It does not import or load anything — the functions are already registered and available. `USE` provides a shorter way to refer to them.

The syntax is:

{{< code lang="fql" >}}
USE target AS alias
{{</ code >}}

### Aliasing a namespace

When the target is a namespace, the alias can be used as a prefix in place of the original namespace.

{{< code lang="fql" >}}
USE IO::FS AS fs

RETURN fs::READ("/tmp/data.json")
{{</ code >}}

Here, `fs::READ(...)` resolves to `IO::FS::READ(...)` at compile time.

### Aliasing a function

When the target includes a function name, the alias can be called directly without any namespace prefix.

{{< code lang="fql" >}}
USE IO::FS::READ AS read

RETURN read("/tmp/data.json")
{{</ code >}}

Here, `read(...)` resolves to `IO::FS::READ(...)`.

## Where functions come from

The set of available namespaced functions depends on the host environment. Functions are registered at engine startup by:

- **The standard library** — built-in functions like `LENGTH`, `CONCAT`, and `FLOOR` are registered at the top level (no namespace needed). Some groups, such as IO, register under a namespace.
- **Modules** — external modules passed to the engine at startup can register functions under any namespace.
- **The host application** — an application embedding Ferret can register its own functions and namespaces directly.

The language does not distinguish between these sources. A namespaced function call works the same way regardless of where the function was registered.

## Runtime-backed functions

Some namespaced functions interact with external systems. They may read files, issue HTTP requests, query databases, control browsers, or work with binary data.

The behavior of such functions depends on the runtime configuration. The same query may behave differently in a CLI context than in a browser automation runtime — not because the language changes, but because the set of available functions and runtime capabilities differs.

## Next steps

{{< docs-related tiles="embedding-modules,stdlib,language-functions" >}}
