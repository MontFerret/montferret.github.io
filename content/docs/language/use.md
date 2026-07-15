---
title: "USE Statement"
sidebarTitle: "USE"
weight: 35
draft: false
description: "Create local aliases for namespaces and namespaced functions with the USE statement."
---

# The USE statement

`USE` creates a local alias for a fully qualified namespace or namespaced function. Use it when repeating a long qualified name would make a script harder to read.

Aliases are resolved when the script is compiled. `USE` does not import, load, or register a module.

## Alias a namespace

Use a namespace alias as the first segment of a qualified function name:

{{< code lang="fql" >}}
USE IO::NET::HTTP AS HTTP

RETURN HTTP::GET("https://api.example.com/data")
{{</ code >}}

Here, `HTTP::GET(...)` resolves to `IO::NET::HTTP::GET(...)`.

The syntax is:

{{< code lang="fql" >}}
USE target AS alias
{{</ code >}}

For a namespace alias, `target` is the fully qualified namespace and `alias` is the shorter name used by the rest of the script.

## Alias a function

`USE` can also alias a specific namespaced function. The alias is then called without a namespace prefix:

{{< code lang="fql" >}}
USE IO::NET::HTTP::GET AS GET

RETURN GET("https://api.example.com/data")
{{</ code >}}

Here, `GET(...)` resolves to `IO::NET::HTTP::GET(...)`.

## Place USE before the script body

`USE` declarations belong in the script header. Put them before variable declarations, function declarations, function calls, and terminal statements such as `RETURN` or a top-level `FOR`.

A script may declare more than one alias:

{{< code lang="fql" >}}
USE IO::FS AS FS
USE IO::NET::HTTP::GET AS GET

LET data = FS::READ("/tmp/data.json")

RETURN [data, GET("https://api.example.com/data")]
{{</ code >}}

## Resolution rules

- Aliases and qualified function names are case-sensitive. `HTTP`, `http`, and `Http` are different names.
- Reusing an alias for a different target produces a compile-time name error.
- An alias does not make its target available. The runtime still needs to register the namespace or function named by the target.
- Runtime capabilities may differ between the CLI, browser runtime, and embedding applications, so an alias can only call functions provided by the current environment.

## Next steps

{{< docs-related tiles="language-functions-modules,language-structure,embedding-modules" >}}
