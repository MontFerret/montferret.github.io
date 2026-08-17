---
title: "Embedding"
weight: 50
draft: false
description: "Run Ferret inside Go and JavaScript applications while the host controls runtime capabilities."
---

# Embedding

Ferret is designed to run inside host applications. The host owns the runtime lifecycle and decides which parameters, functions, modules, and capabilities an FQL program can use.

Ferret currently supports two embedding environments:

- **Go embedding** uses the native `github.com/MontFerret/ferret/v2` library. Go applications can configure the complete runtime, including modules, host values, codecs, filesystem access, and execution limits.
- **JavaScript embedding** uses `@montferret/ferret`. The package compiles the Ferret runtime to WebAssembly and exposes it through a JavaScript API for Node.js and modern browsers.

The JavaScript package is not a separate implementation of Ferret. It runs the same Ferret compiler, virtual machine, and standard library through the WASM build.

Choose the host environment in which your application runs. FQL syntax and language behavior stay the same, while the available host integrations and configuration APIs depend on that environment.

{{< docs-related tiles="embedding-go,embedding-javascript,language,language-parameters" >}}
