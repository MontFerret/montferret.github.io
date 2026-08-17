---
title: "use Statement"
sidebarTitle: "Use"
weight: 10
draft: false
description: "Create local aliases for namespaces and namespaced functions with the use statement."
---

# The use statement

`use` creates a local alias for a fully qualified namespace or namespaced function. Use it when repeating a long qualified name would make a script harder to read.

Aliases are resolved when the script is compiled. `use` does not import, load, or register a module.

## Alias a namespace

Use a namespace alias as the first segment of a qualified function name:

{{< code lang="fql" >}}
use io::net::http as http

return http::get("https://api.example.com/data")
{{</ code >}}

Here, `http::get(...)` resolves to `io::net::http::get(...)`.

The syntax is:

{{< code lang="fql" >}}
use target as alias
{{</ code >}}

For a namespace alias, `target` is the fully qualified namespace and `alias` is the shorter name used by the rest of the script.

## Alias a function

`use` can also alias a specific namespaced function. The alias is then called without a namespace prefix:

{{< code lang="fql" >}}
use io::net::http::get as get

return get("https://api.example.com/data")
{{</ code >}}

Here, `get(...)` resolves to `io::net::http::get(...)`.

## Place use before the script body

`use` declarations belong in the script header. Put them before variable declarations, function declarations, function calls, standalone loops, and `return` statements.

A script may declare more than one alias:

{{< code lang="fql" >}}
use io::fs as fs
use io::net::http::get as get

let data = fs::read("/tmp/data.json")

return [data, get("https://api.example.com/data")]
{{</ code >}}

## Resolution rules

- Alias names remain case-sensitive. Registered namespace segments and host-function names are case-insensitive and are documented in canonical lowercase.
- Reusing an alias for a different target produces a compile-time name error.
- An alias does not make its target available. The runtime still needs to register the namespace or function named by the target.
- Runtime capabilities may differ between the CLI, browser runtime, and embedding applications, so an alias can only call functions provided by the current environment.

## Next steps

{{< docs-related tiles="language-functions-modules,language-structure,embedding-go-modules" >}}
