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

When a script repeatedly calls functions from the same namespace, `USE` can create a shorter local alias. Put the declaration at the start of the script, before the script body.

{{< code lang="fql" >}}
USE IO::FS AS FS

RETURN FS::READ("/tmp/data.json")
{{</ code >}}

`USE` is a compile-time alias; it does not load a module or change which functions the runtime provides. See the [`USE` statement reference]({{< ref "/docs/language/script-structure/use" >}}) for function aliases, placement rules, and name-resolution behavior.

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

{{< docs-related tiles="language-use,embedding-modules,stdlib,language-functions" >}}
