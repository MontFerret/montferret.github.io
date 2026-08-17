---
title: "Modules and Lifecycle"
sidebarTitle: "Modules and Lifecycle"
weight: 50
draft: false
description: "Define JavaScript-backed Ferret modules with functions and asynchronous lifecycle hooks."
---

# Modules and Lifecycle

JavaScript modules group host functions with callbacks for engine initialization, compilation, execution, and cleanup. Use them when an integration needs to acquire resources, observe work, or release state along with exposing functions to FQL.

## Define a module

Use `defineModule()` to validate a definition while preserving its inferred TypeScript type:

{{< code lang="javascript" >}}
import { create, defineModule } from "@montferret/ferret";

const events = [];

const telemetry = defineModule({
    name: "telemetry",
    functions: {
        track: (event) => {
            events.push(event);
            return event;
        },
    },
    lifecycle: {
        async onInit() {
            await Promise.resolve();
            events.push({ type: "engine:init" });
        },
        async onClose() {
            events.push({ type: "engine:close" });
        },
    },
});

const engine = await create({
    modules: [telemetry],
});

try {
    await engine.run(`return track({ type: "query:run" })`);
} finally {
    await engine.close();
}
{{</ code >}}

The module name identifies the registration; it does not prefix the module's FQL functions. The example therefore calls `track()`, not `telemetry::track()`.

Every lifecycle callback may return immediately or return a promise. Ferret waits for the promise before continuing the operation.

## Define modules in TypeScript

The Node.js and browser entry points export `ModuleDefinition`, `ModuleLifecycle`, `MaybePromise`, and all lifecycle event types. A definition can also be declared directly:

{{< code lang="typescript" >}}
import type { ModuleDefinition } from "@montferret/ferret";

const audit: ModuleDefinition = {
    name: "audit",
    lifecycle: {
        beforeRun(event) {
            console.log("Ferret execution started", event);
        },
    },
};
{{</ code >}}

`defineModule()` returns the same object that it receives. It provides validation and preserves the specific inferred type rather than wrapping the module in a runtime class.

## Registration and immutability

`create({ modules })` requires an array. Each module definition, its `functions` object, and its `lifecycle` object must be a plain JavaScript object. Module names must contain at least one non-whitespace character, function values must be callable, and lifecycle objects may contain only the documented callback names.

Module names are case-sensitive stable identifiers. Two modules named `audit` are rejected, while `audit` and `Audit` are distinct registrations as long as their functions do not conflict.

Engine creation revalidates every definition and snapshots its functions and lifecycle callbacks. Mutating the original objects afterward does not change the engine. Create another engine when the application needs a different module configuration.

Function names continue to use Ferret's case-insensitive canonical resolution. Engine creation rejects conflicts between two modules, a module and the `functions` shorthand, or a user-defined function and the standard library instead of choosing an implementation. Module names do not create an FQL namespace automatically.

The existing `functions` option remains available for integrations that only need host functions. It can be combined with modules and is registered before the modules:

{{< code lang="javascript" >}}
const engine = await create({
    functions: {
        app_name: () => "inventory",
    },
    modules: [telemetry, audit],
});
{{</ code >}}

## Lifecycle callbacks

The `ModuleLifecycle` interface exposes callbacks from Ferret Core's engine, plan, and session lifecycle:

| Callback | When it runs | Event |
| --- | --- | --- |
| `onInit` | During `create()`, after module registration | None |
| `onClose` | During engine cleanup | None |
| `beforeCompile` | Before a source is compiled | `CompileEvent` |
| `afterCompile` | After a compilation attempt | `CompileResultEvent` |
| `onPlanClose` | When a plan closes | `PlanEvent` |
| `beforeRun` | Before a session executes | `RunEvent` |
| `afterRun` | After an execution attempt | `RunResultEvent` |
| `onSessionClose` | When a session closes | `SessionEvent` |

Each callback returns `MaybePromise<void>`, which means either `void` or a promise-like value that resolves to `void`.

### Compile events

Compile callbacks receive an immutable snapshot of the normalized source:

{{< code lang="typescript" >}}
const compilerLog = defineModule({
    name: "compiler-log",
    lifecycle: {
        beforeCompile(event) {
            console.log(event.source.name, event.source.text);
        },
        afterCompile(event) {
            if (event.error !== undefined) {
                console.error("Compilation failed", event.error);
            }
        },
    },
});
{{</ code >}}

`CompileEvent.source` always contains `name` and `text`. `CompileResultEvent` adds `error?: unknown`; a Ferret compiler failure is represented as a JavaScript `Error`.

### Run and cleanup events

`RunResultEvent` adds `error?: unknown`, with Ferret VM failures represented as JavaScript `Error` objects. `PlanEvent`, `RunEvent`, and `SessionEvent` are immutable empty objects today because Ferret Core does not provide stable per-object metadata for those hooks. Their interfaces remain extensible for future stable fields.

Lifecycle events do not expose Go contexts, pointers, registrars, bridge identifiers, or handles. Cancellation remains available through the `AbortSignal` passed to compile, run, and session APIs.

## Ordering and errors

Lifecycle ordering follows Ferret Core:

| Hook group | Order | Failure behavior |
| --- | --- | --- |
| Initialization and before hooks | Module registration order (FIFO) | Stop at the first failure |
| After and close hooks | Reverse registration order (LIFO) | Continue and aggregate failures |

If a before hook fails, its operation does not begin and the corresponding after hooks do not run. When compilation or execution begins and then fails, result hooks still run and receive the underlying compiler or VM error.

A callback that throws synchronously or returns a rejected promise rejects the operation associated with that hook. Error messages include the module and lifecycle phase while retaining the original thrown or rejected message. An `onInit` failure rejects `create()` and triggers cleanup of the registered engine hooks.

## Closing resources

Lifecycle-aware cleanup preserves the resource rules described in [Executing Ferret]({{< ref "/docs/embedding/javascript/executing" >}}):

- A close preflight rejects without cleanup while compilation or session creation is pending, or while a child session is running or closing.
- Once close begins, the resource becomes closed even if a lifecycle callback rejects. Remaining child cleanup and close callbacks continue, and the close promise rejects with the aggregated errors.
- Calling `close()` again does not rerun lifecycle callbacks.
- Operations that re-enter a resource while it is initializing or closing reject instead of racing its lifecycle. Do not call operations on that same resource from its close callback.

Always await `close()` and handle cleanup failures at the application boundary.

## Next steps

{{< docs-related tiles="embedding-javascript-custom-functions,embedding-javascript-executing,embedding-javascript-limitations,language-functions-modules" >}}
