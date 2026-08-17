---
title: "Executing Ferret"
sidebarTitle: "Executing Ferret"
weight: 20
draft: false
description: "Run one-off programs or reuse compiled plans and sessions from JavaScript."
---

# Executing Ferret

The JavaScript API exposes three runtime objects:

| Object | Responsibility |
| --- | --- |
| `Engine` | Owns one Ferret WASM runtime and its registered JavaScript functions. |
| `Plan` | Holds a compiled FQL program that can be reused. |
| `Session` | Holds one execution context with captured parameters. |

All creation and execution operations are asynchronous.

The package exports TypeScript declarations for the runtime interfaces and their options:

| API | Input and options |
| --- | --- |
| `create(options?: CreateOptions)` | Creates an `Engine`; options cover WASM loading, host functions, and `HTTPOptions`. |
| `Engine.compile(source, options?: CompileOptions)` | Accepts `SourceInput` and an optional cancellation signal. |
| `Engine.run(source, options?: ExecutionOptions)` | Accepts `SourceInput`, `Params`, and an optional cancellation signal. |
| `Plan.createSession(options?: SessionOptions)` | Captures `Params` and accepts a signal for session creation. |
| `Plan.run(options?: ExecutionOptions)` | Runs the plan with parameters and an optional signal. |
| `Session.run(options?: SessionRunOptions)` | Runs the captured session with an optional signal. |

`SourceInput`, `Params`, `RuntimeFunction`, `HTTPOptions`, and `Version` are also exported. Cancellation options use the standard `AbortSignal` available in supported Node.js and browser environments.

## Run a one-off program

Use `engine.run()` when a program will run once:

{{< code lang="javascript" >}}
const result = await engine.run(`
    return for value in 1..3
        return value * 2
`);

console.log(result); // [2, 4, 6]
{{</ code >}}

This method compiles a temporary plan, creates a temporary session, runs it, and closes both resources before resolving.

Pass a source object when the program has a meaningful filename. The name is included in compiler diagnostics:

{{< code lang="javascript" >}}
await engine.run({
    name: "invoice-filter.fql",
    text: "return @invoices",
}, {
    params: { invoices: [] },
});
{{</ code >}}

## Compile and reuse a plan

Compile once when the same FQL source will run with different parameters:

{{< code lang="javascript" >}}
const plan = await engine.compile("return @value * 2");

try {
    console.log(plan.params); // ["value"]
    console.log(await plan.run({ params: { value: 3 } })); // 6
    console.log(await plan.run({ params: { value: 5 } })); // 10
} finally {
    await plan.close();
}
{{</ code >}}

`plan.params` is an immutable array of parameter names declared by the compiled program. `plan.run()` creates and closes a fresh session for each call.

## Reuse a session

A session captures its parameters and can run more than once:

{{< code lang="javascript" >}}
const plan = await engine.compile("return @factor * 2");
const session = await plan.createSession({
    params: { factor: 4 },
});

try {
    console.log(await session.run()); // 8
    console.log(await session.run()); // 8
} finally {
    await session.close();
    await plan.close();
}
{{</ code >}}

One session cannot execute concurrent runs. Use separate sessions when executions need to overlap.

## Cancel work

Pass an `AbortSignal` to `engine.compile()`, `engine.run()`, `plan.createSession()`, `plan.run()`, or `session.run()`:

{{< code lang="javascript" >}}
const controller = new AbortController();
const pending = engine.run("wait(10000) return true", {
    signal: controller.signal,
});

controller.abort();

try {
    await pending;
} catch (error) {
    if (error instanceof Error && error.name === "AbortError") {
        console.log("Execution cancelled");
    }
}
{{</ code >}}

Ferret execution and HTTP calls observe cancellation. See [Custom Functions]({{< ref "/docs/embedding/javascript/custom-functions" >}}) for the cancellation behavior of promises returned by JavaScript functions.

## Work with results and errors

Results are encoded by Ferret as JSON and parsed before the promise resolves. Ordinary values become JavaScript primitives, arrays, objects, or `null`. A Ferret binary result follows Ferret's JSON serialization and is returned as a base64 string.

Invalid source rejects during compilation. Runtime failures reject during execution. The current package exposes these as JavaScript `Error` objects; cancellation is distinguished by the `AbortError` name.

## Close resources

`close()` is asynchronous and idempotent on engines, plans, and sessions. The `closed` property reports local state synchronously.

- Closing a plan closes its idle sessions.
- Closing an engine closes its idle plans and sessions.
- Closing rejects without partial cleanup while compilation or session creation is pending, or while a child session is running.

Keep the engine alive while the application is using Ferret, then close it in a `finally` block during shutdown.

## Next steps

{{< docs-related tiles="embedding-javascript-parameters,embedding-javascript-custom-functions,embedding-javascript-limitations" >}}
