---
title: "Custom Functions"
sidebarTitle: "Custom Functions"
weight: 40
draft: false
description: "Register synchronous and asynchronous JavaScript functions for FQL programs."
---

# Custom Functions

Register JavaScript functions when creating an engine to make application behavior callable from FQL.

## Register functions

Pass a plain object through the `functions` option:

{{< code lang="javascript" >}}
const engine = await create({
    functions: {
        full_name: (first, last) => `${first} ${last}`,
        "inventory::total": (price, quantity) => price * quantity,
    },
});

try {
    const result = await engine.run(`
        return {
            name: full_name("Ada", "Lovelace"),
            total: inventory::total(12.5, 4),
        }
    `);
    console.log(result); // { name: "Ada Lovelace", total: 50 }
} finally {
    await engine.close();
}
{{</ code >}}

Function names are trimmed and canonicalized, including namespace segments. FQL resolution is case-insensitive. Names that normalize to the same function, such as `total` and `TOTAL`, are rejected as duplicates.

The function registry is immutable. Create another engine when an application needs a different set of functions.

## Arguments and return values

Ferret converts each argument into a JavaScript value before invoking the function. Return a supported value directly or resolve a promise with one:

{{< code lang="javascript" >}}
const engine = await create({
    functions: {
        wrap: (value) => ({ value }),
        load_status: async () => ({ status: "ready" }),
    },
});
{{</ code >}}

Functions can receive and return `null`, booleans, strings, finite numbers, arrays, plain objects, and binary values. Ferret binary arguments arrive as `Uint8Array`. The same unsupported-value rules described in [Parameters]({{< ref "/docs/embedding/javascript/parameters" >}}) apply to function results.

## Asynchronous functions

A registered function may return a promise. Ferret waits for it before continuing the FQL program:

{{< code lang="javascript" >}}
const engine = await create({
    functions: {
        lookup_user: async (id) => ({
            id,
            status: "ready",
        }),
    },
});

const user = await engine.run("return lookup_user(@id)", {
    params: { id: "user-42" },
});
{{</ code >}}

## Errors

If a synchronous function throws or an asynchronous function rejects, the current Ferret execution rejects with an error containing the JavaScript failure message:

{{< code lang="javascript" >}}
const engine = await create({
    functions: {
        fail: async () => {
            throw new Error("host lookup failed");
        },
    },
});

try {
    await engine.run("return fail()");
} catch (error) {
    console.error(error instanceof Error ? error.message : error);
    // host lookup failed
} finally {
    await engine.close();
}
{{</ code >}}

Use ordinary JavaScript error handling inside the function when the application can recover or return a domain value instead.

## Cancellation

An `AbortSignal` cancels Ferret execution, but the runtime cannot generically cancel a promise returned by application code. When cancellation arrives while a JavaScript promise is pending, the WASM runtime waits for that promise to settle and then reports the execution as aborted.

If the underlying operation supports cancellation, capture and use an application-owned signal or cancellation mechanism inside the function.

## Next steps

{{< docs-related tiles="embedding-javascript-parameters,embedding-javascript-executing,embedding-javascript-limitations,language-functions-modules" >}}
