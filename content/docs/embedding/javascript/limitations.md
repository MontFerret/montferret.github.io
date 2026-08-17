---
title: "Runtime and WASM Limitations"
sidebarTitle: "Runtime and WASM Limitations"
weight: 60
draft: false
description: "Understand the current environment, loading, networking, and API boundaries of @montferret/ferret."
---

# Runtime and WASM Limitations

`@montferret/ferret` runs the Ferret v2 runtime through its WebAssembly build. The language and virtual machine are Ferret; the boundaries below come from the host environment or from features the current JavaScript API does not expose.

## Runtime initialization

Each `create()` call loads and instantiates a separate WASM runtime. Initialization is asynchronous and has a meaningful startup and memory cost, so applications should normally reuse an engine instead of creating one for every query.

Close each engine during application shutdown. An engine owns its plans, sessions, HTTP transport, and WASM runtime.

## Browser assets

The browser export loads `ferret.wasm` and `wasm_exec.js` relative to the package entry point. Both assets must remain available after bundling and deployment.

Override the WASM source when the application serves it from another location:

{{< code lang="javascript" >}}
const engine = await create({
    wasm: new URL("/assets/ferret.wasm", location.href),
});
{{</ code >}}

The `wasm` option accepts a URL, a Node.js file path, an `ArrayBuffer`, a `Uint8Array`, or a precompiled `WebAssembly.Module`. It does not replace `wasm_exec.js`; the matching Go runtime support file must still be served with the package entry point.

## HTTP behavior

The package supplies platform-specific HTTP transports while Ferret retains request policy, size limits, redirect checks, and cancellation.

| Environment | Current behavior |
| --- | --- |
| Node.js | Uses per-engine `node:http` and `node:https` connection pools. DNS results are policy-checked and the selected address is pinned for the request. Allowed redirects are checked before being followed. |
| Browser | Allows same-origin requests only. Redirects are rejected because the browser does not expose enough DNS and redirect information for Ferret to apply the same policy safely. Requests use same-origin credentials. |

Localhost and loopback destinations are denied by default. Trusted applications can enable loopback access when creating an engine:

{{< code lang="javascript" >}}
const engine = await create({
    http: { allowLocalhost: true },
});
{{</ code >}}

This option enables loopback only. Private and link-local networks remain denied.

## Filesystem access

The full Ferret standard library is registered, but the JavaScript API does not expose a filesystem-root option. FQL filesystem functions therefore fail because no filesystem root is configured.

This is a current API limitation, not a claim that WebAssembly can never access host files.

## Modules and host integrations

The JavaScript API supports declarative modules made from JavaScript functions and lifecycle callbacks. These modules participate in Ferret Core's engine, plan, and session hooks, but they do not load arbitrary compiled Go modules or expose Go registration machinery.

The package does not expose APIs for:

- registering Go modules or Registry packages
- implementing Ferret host values or capability interfaces in JavaScript
- registering custom output codecs
- selecting standard-library groups
- configuring logging, filesystem roots, or VM/session limits
- loading or serializing precompiled Ferret program artifacts
- compiler-analysis snapshots

Use native [Go embedding]({{< ref "/docs/embedding/go" >}}) when an application needs those integration surfaces.

## Browser automation

Running Ferret inside a web browser does not automatically provide Ferret's browser-automation modules. JavaScript-defined modules cannot supply the Go host values and CDP driver required for browser-backed pages, elements, screenshots, and dispatch capabilities, so those features remain unavailable through this package alone.

Use a Ferret distribution that registers the appropriate browser modules when an FQL program needs those capabilities.

## Values and output

The JavaScript boundary accepts JSON-compatible values plus `Uint8Array`. It does not accept arbitrary class instances, cyclic object graphs, non-finite numbers, or JavaScript-specific values without a Ferret representation.

Results are JSON-decoded JavaScript values. The current API does not expose MessagePack or custom output codecs, and binary results use Ferret's JSON base64 representation.

## Execution concurrency

Multiple engines can exist independently, and plans can create separate sessions. A single session cannot run concurrently, and closing an engine or plan rejects while a child execution or resource creation is active.

## Next steps

{{< docs-related tiles="embedding-javascript-getting-started,embedding-javascript-executing,embedding-javascript-custom-functions,embedding-javascript-modules" >}}
