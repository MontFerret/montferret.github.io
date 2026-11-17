---
title: "Executing Ferret"
sidebarTitle: "Executing Ferret"
weight: 20
draft: false
description: "Run one-off programs or reuse compiled plans and sessions from JavaScript."
---

# Executing Ferret

Engines and plans can run a program in a temporary session. Create a session explicitly to reuse its parameters.

| Object | Responsibility |
| --- | --- |
| `Engine` | Owns one Ferret WASM runtime and its registered JavaScript functions and modules. |
| `Plan` | Holds a compiled FQL program that can be reused. |
| `Session` | Holds one execution context with captured parameters. |
| `Output` | Returned by `run()`. Consume it once or close it to discard unread output. |

All creation and execution operations are asynchronous.

The package exports TypeScript declarations for the runtime interfaces and their options:

| API | Input and options |
| --- | --- |
| `create(options?: CreateOptions)` | Creates an `Engine`; options cover WASM loading, host functions, modules, and `HTTPOptions`. |
| `Engine.compile(source, options?: CompileOptions)` | Accepts `SourceInput` and an optional cancellation signal. |
| `Engine.run(source, options?: ExecutionOptions)` | Accepts `SourceInput`, `Params`, and an optional cancellation signal. |
| `Plan.createSession(options?: SessionOptions)` | Captures `Params` and accepts a signal for session creation. |
| `Plan.run(options?: ExecutionOptions)` | Runs the plan with parameters and an optional signal. |
| `Session.run(options?: SessionRunOptions)` | Runs the captured session with an optional signal. |

`SourceInput`, `Params`, `RuntimeFunction`, `HTTPOptions`, and `Version` are also exported. Module integrations can use `ModuleDefinition`, `ModuleLifecycle`, `MaybePromise`, `CompileEvent`, `CompileResultEvent`, `PlanEvent`, `RunEvent`, `RunResultEvent`, and `SessionEvent`. Cancellation options use the standard `AbortSignal` available in supported Node.js and browser environments.

## Run a one-off program

Use `engine.runJSON()` when a program will run once:

{{< code lang="javascript" >}}
const result = await engine.runJSON(`
    return for value in 1..3
        return value * 2
`);

console.log(result); // [2, 4, 6]
{{</ code >}}

This method compiles a temporary plan, creates a temporary session, runs it, and closes both resources before resolving.

Pass a source object when the program has a meaningful filename. The name is included in compiler diagnostics:

{{< code lang="javascript" >}}
await engine.runJSON({
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
    console.log(await plan.runJSON({ params: { value: 3 } })); // 6
    console.log(await plan.runJSON({ params: { value: 5 } })); // 10
} finally {
    await plan.close();
}
{{</ code >}}

`plan.params` is an immutable array of parameter names declared by the compiled program. `plan.runJSON()` creates and closes a fresh session for each call.

## Reuse a session

A session captures its parameters and can run more than once:

{{< code lang="javascript" >}}
const plan = await engine.compile("return @factor * 2");
const session = await plan.createSession({
    params: { factor: 4 },
});

try {
    console.log(await session.runJSON()); // 8
    console.log(await session.runJSON()); // 8
} finally {
    await session.close();
    await plan.close();
}
{{</ code >}}

Settle the previous output before reusing a session. `runJSON()` does this for you. Create separate sessions for overlapping runs.

## Cancel work

Pass an `AbortSignal` to `engine.compile()`, `engine.runJSON()`, `plan.createSession()`, `plan.runJSON()`, or `session.runJSON()`:

{{< code lang="javascript" >}}
const controller = new AbortController();
const pending = engine.runJSON("wait(10000) return true", {
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

Use `runJSON()` for a JavaScript value. Use `run()` to inspect the encoding metadata or consume the bytes yourself. For example, collect a result before closing its engine:

{{< code lang="javascript" >}}
const output = await engine.run("return @count", { params: { count: 3 } });
try {
    const content = await output.collect();
    if (content) {
        console.log(content.metadata.contentType); // application/json
        console.log(new TextDecoder().decode(content.data)); // 3
    }
} finally {
    await output.close();
}
{{</ code >}}

`collect()` returns detached bytes, or `null` for absent content. `json()` decodes the same output, returning `undefined` for absent content. A query that returns FQL `none` encodes JSON `null`.

To forward output as it arrives, use `await output.consume(async (chunk, signal) => { … })`. Ferret waits for each callback before delivering more data. The JavaScript callback owns its bytes. A chunk can end in the middle of a JSON value or UTF-8 character.

Compilation and preparation can reject before `run()` returns a handle. Terminal execution, encoding, and cleanup failures can surface when you consume it. A rejected collection can carry partial content in `OutputError.content`. Keep the run's signal valid until consumption finishes.

## Close resources

`close()` is asynchronous and idempotent on engines, plans, and sessions. The `closed` property reports local state synchronously.

- Closing a plan closes its idle sessions.
- Closing an engine also closes unread outputs and idle plans and sessions. An active output consumption must settle before the engine can close.
- Closing rejects without partial cleanup while compilation or session creation is pending, while a child session is running, or while a child resource is already closing.
- Once close begins, the resource stays closed even if a lifecycle hook rejects. Remaining cleanup continues, errors are aggregated, and later `close()` calls do not rerun hooks.

Keep the engine alive while the application is using Ferret, then close it in a `finally` block during shutdown.

## Next steps

{{< docs-related tiles="embedding-javascript-parameters,embedding-javascript-custom-functions,embedding-javascript-modules,embedding-javascript-limitations" >}}
