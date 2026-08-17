---
title: "Getting Started"
sidebarTitle: "Getting Started"
weight: 10
draft: false
description: "Install @montferret/ferret and run an FQL program from JavaScript."
---

# Getting Started

`@montferret/ferret` runs Ferret inside Node.js and browser applications. The package loads the Ferret runtime as WebAssembly and exposes an asynchronous JavaScript API.

## Requirements

Use one of these environments:

- Node.js 22 or newer
- a modern browser with WebAssembly, `fetch`, and `crypto.getRandomValues`

Building the package from source also requires Go 1.25 or newer. Applications that install the published package do not need Go.

## Install the package

{{< terminal >}}
npm install @montferret/ferret
{{< /terminal >}}

## Run a program

Import `create`, initialize an engine, and pass an FQL program to `engine.run()`:

{{< code lang="javascript" >}}
import { create } from "@montferret/ferret";

const engine = await create();

try {
    const result = await engine.run("return 1 + 1");
    console.log(result); // 2
} finally {
    await engine.close();
}
{{</ code >}}

`create()` loads and initializes the WASM runtime. There is no separate initialization call. `engine.run()` compiles and executes the program asynchronously, then returns its result as a JavaScript value.

Always close the engine when the application no longer needs it. Closing releases plans, sessions, network resources, and the WASM runtime owned by that engine.

## Use CommonJS in Node.js

The Node.js package export also supports `require`:

{{< code lang="javascript" >}}
const { create } = require("@montferret/ferret");
{{</ code >}}

The runtime API is otherwise the same.

## Handle errors

Initialization, compilation, execution, and cleanup can reject. Handle errors at the application boundary while still closing any engine that was created successfully:

{{< code lang="javascript" >}}
import { create } from "@montferret/ferret";

let engine;

try {
    engine = await create();
    const result = await engine.run("return missing_function()");
    console.log(result);
} catch (error) {
    console.error(error instanceof Error ? error.message : error);
} finally {
    await engine?.close();
}
{{</ code >}}

Compilation and execution failures reject with JavaScript errors. Cancelled operations reject with an error whose `name` is `AbortError`.

## Browser applications

Browser-aware bundlers select the browser package export. The export loads `ferret.wasm` and `wasm_exec.js` relative to the generated JavaScript entry point, so those files must be served with the bundle.

See [Runtime and WASM Limitations]({{< ref "/docs/embedding/javascript/limitations" >}}) for custom WASM loading and the runtime features currently exposed in browsers.

## Next steps

{{< docs-related tiles="embedding-javascript-executing,embedding-javascript-parameters,embedding-javascript-custom-functions,embedding-javascript-limitations" >}}
