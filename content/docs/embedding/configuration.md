---
title: "Configuration"
sidebarTitle: "Configuration"
weight: 50
draft: false
description: "Control the standard library, logging, encoding, concurrency, and file system access."
---

# Configuration

The engine accepts functional options that control which capabilities are available to scripts, how results are encoded, how execution is logged, and how concurrency is managed.

## Standard library

By default, the engine loads the full standard library. You can restrict it to a subset of groups or disable it entirely.

### Available groups

| Group | Contents |
|-------|----------|
| `Types` | Type checking and conversion functions |
| `Strings` | String manipulation |
| `Math` | Mathematical operations |
| `Collections` | Generic collection operations |
| `DateTime` | Date and time functions |
| `Arrays` | Array-specific operations |
| `Objects` | Object-specific operations |
| `IO` | Expands to `FS` + `NET` |
| `FS` | File system access |
| `NET` | Network operations |
| `Path` | File path manipulation |
| `Utils` | Utility functions |
| `Testing` | Test assertion functions |

### Selecting groups

{{< code lang="go" >}}
// Full standard library (default)
ferret.New(
    ferret.WithStdlib(stdlib.Full()),
)

// Everything except file system and network access
ferret.New(
    ferret.WithStdlib(stdlib.Safe()),
)

// Only specific groups
ferret.New(
    ferret.WithStdlib(stdlib.Only(stdlib.Strings, stdlib.Math, stdlib.Arrays)),
)

// Full minus specific groups
ferret.New(
    ferret.WithStdlib(stdlib.Full().Without(stdlib.IO, stdlib.Testing)),
)

// No standard library at all
ferret.New(
    ferret.WithoutStdlib(),
)
{{</ code >}}

`stdlib.Safe()` is equivalent to `stdlib.Full().Without(stdlib.IO)` — it removes file system and network access while keeping everything else.

## Logging

Logging is configured at two levels: engine-wide defaults and per-session overrides.

### Engine-level logging

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithLog(os.Stdout),
    ferret.WithLogLevel(logging.InfoLevel),
    ferret.WithLogFields(map[string]any{
        "service": "query-engine",
    }),
)
{{</ code >}}

### Session-level logging

Per-session options override the engine defaults for that execution:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithSessionLog(os.Stderr),
    ferret.WithSessionLogLevel(logging.DebugLevel),
    ferret.WithSessionLogFields(map[string]any{
        "request_id": "abc-123",
    }),
)
{{</ code >}}

### Log levels

| Level | Constant |
|-------|----------|
| Trace | `logging.TraceLevel` |
| Debug | `logging.DebugLevel` |
| Info | `logging.InfoLevel` |
| Warn | `logging.WarnLevel` |
| Error | `logging.ErrorLevel` |
| Disabled | `logging.Disabled` |

## Output encoding

Query results are encoded before being returned as an `Output`. The default encoding is JSON.

### Selecting the output format

Set the content type at the session level:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithOutputContentType("application/msgpack"),
)
{{</ code >}}

### Built-in codecs

| Content type | Format |
|-------------|--------|
| `application/json` | JSON (default) |
| `application/msgpack` | MessagePack |

### Registering a custom codec

Implement the `encoding.Codec` interface and register it on the engine:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithEncodingCodec("application/yaml", myYAMLCodec),
)
{{</ code >}}

Or replace the entire registry:

{{< code lang="go" >}}
registry := encoding.NewRegistry(jsonCodec, msgpackCodec, yamlCodec)

engine, err := ferret.New(
    ferret.WithEncodingRegistry(registry),
)
{{</ code >}}

## Concurrency control

Three options control how the engine manages concurrent execution:

### `WithMaxActiveSessions`

Limits the total number of sessions that can be running at the same time across the entire engine. When the limit is reached, `plan.NewSession` blocks until another session closes or the context is cancelled.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithMaxActiveSessions(100),
)
{{</ code >}}

Use this to cap global resource consumption — CPU, memory, network, or downstream service pressure.

### `WithMaxVMsPerPlan`

A hard cap on the total number of virtual machines a single plan can own, including both idle and active VMs. When the limit is reached and no idle VM is available, session creation fails with an error.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithMaxVMsPerPlan(16),
)
{{</ code >}}

Use this to bound the resource footprint of a single frequently-executed query.

### `WithMaxIdleVMsPerPlan`

Controls how many idle VMs each plan keeps cached for reuse after sessions close. When the idle cache is full, additional returned VMs are discarded instead of retained. The default is 8.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithMaxIdleVMsPerPlan(4),
)
{{</ code >}}

This is a memory-vs-latency trade-off: more idle VMs means faster session creation but higher steady-state memory usage.

### How they work together

- `WithMaxActiveSessions` is an engine-wide concurrency gate — it controls how many sessions run at once across all plans
- `WithMaxVMsPerPlan` is a per-plan resource cap — it bounds the VM count for one specific compiled query
- `WithMaxIdleVMsPerPlan` is a per-plan cache policy — it decides how many idle VMs stay warm

A value of 0 disables the limit for `WithMaxActiveSessions` and `WithMaxVMsPerPlan`.

## File system

Scripts that use file system functions (from the `FS` stdlib group) operate within a sandboxed file system. The host can configure the root directory and restrict access to read-only.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithFSRoot("/data/extractions"),
    ferret.WithFSReadOnly(),
)
{{</ code >}}

If no root is set, file system functions use the process working directory. `WithFSReadOnly` prevents scripts from writing files regardless of the root.

## Compiler options

The compiler can be configured through `WithCompilerOptions`:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithCompilerOptions(
        compiler.WithOptimizationLevel(compiler.O1),
    ),
)
{{</ code >}}

| Level | Constant | Description |
|-------|----------|-------------|
| None | `compiler.O0` | No optimization |
| Basic | `compiler.O1` | Basic optimizations (default) |

Debug compilation (`engine.CompileDebug`) always uses `O0` to ensure stable source-level debugging metadata.

## Next steps

{{< docs-related tiles="embedding-modules" >}}