---
title: "Configuration"
sidebarTitle: "Configuration"
weight: 50
draft: false
description: "Control the standard library, logging, encoding, concurrency, and sandboxed host services."
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

Query results are encoded before being returned as an `Output`. The default encoding is JSON; MessagePack is also built in. Set the content type at the session level with `WithOutputContentType`:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithOutputContentType("application/vnd.msgpack"),
)
{{</ code >}}

You can register custom codecs on the engine with `WithEncodingCodec`. See [Value Encoders]({{< ref "value-encoders" >}}) for the codec interfaces, hooks, registry, and a complete custom codec example.

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

## Sandboxed components

Sandboxed components control how scripts access host resources outside the Ferret runtime. Configure each component independently according to what the embedding application should allow. The engine adds the configured components to the `context.Context` passed to functions during execution.

### File system

Scripts that use file system functions (from the `FS` stdlib group) operate within a sandboxed file system. The host can configure the root directory and restrict access to read-only.

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithFSRoot("/data/extractions"),
    ferret.WithFSReadOnly(),
)
{{</ code >}}

If `WithFSRoot` is not set, file system access is disabled and `FS` functions return a root-not-configured error. `WithFSReadOnly` prevents scripts from writing within the configured root.

#### Accessing the file system from a function

Retrieve the configured file system from the function context. Prefer the narrowest helper for the operation the function performs:

{{< code lang="go" title="read_file.go" >}}
package files

import (
    "context"

    ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

func ReadFile(ctx context.Context, pathArg runtime.Value) (runtime.Value, error) {
    path, err := runtime.CastArg[runtime.String](pathArg, 0)
    if err != nil {
        return runtime.None, err
    }

    reader, err := ferretfs.ReaderFrom(ctx)
    if err != nil {
        return runtime.None, err
    }

    data, err := reader.ReadFile(path.String())
    if err != nil {
        return runtime.None, err
    }

    return runtime.NewBinary(data), nil
}
{{</ code >}}

The `fs` package also provides `WriterFrom`, `DirectoriesFrom`, and `RemoverFrom`. Use `FileSystemFrom` only when a function needs several capabilities.

For public modules, these context helpers are the recommended way to access files. They keep the module inside the root and read-only policy selected by the embedding host. Calling `os.ReadFile`, `os.WriteFile`, or constructing a separate file system would bypass that configuration.

### HTTP client

Scripts that use network functions (from the `NET` stdlib group) make outbound requests through the engine's network service. Supply a policy-configured HTTP client to restrict which destinations scripts can reach and how much data they can send or receive.

{{< code lang="go" >}}
httpClient := ferrethttp.New(
    ferrethttp.WithAllowedSchemes("https"),
    ferrethttp.WithAllowedHosts("api.example.com"),
    ferrethttp.WithTimeout(10*time.Second),
    ferrethttp.WithMaxRequestSize(1<<20),  // 1 MiB
    ferrethttp.WithMaxResponseSize(4<<20), // 4 MiB
)

engine, err := ferret.New(
    ferret.WithNetwork(
        ferretnet.New(
            ferretnet.WithHTTPClient(httpClient),
        ),
    ),
)
{{</ code >}}

The HTTP client supports these policy controls:

| Control | Options |
|---------|---------|
| URL schemes and destinations | `WithAllowedSchemes`, `WithAllowedHosts`, `WithBlockedHosts` |
| Local, private, and link-local addresses | `WithAllowLocalhost`, `WithAllowPrivateNetworks`, `WithAllowLinkLocal` |
| Redirects | `WithFollowRedirects`, `WithMaxRedirects` |
| Request headers | `WithDefaultHeader`, `WithDefaultHeaders`, `WithBlockedRequestHeaders` |
| Time and payload limits | `WithTimeout`, `WithMaxRequestSize`, `WithMaxResponseSize` |

By default, the client allows HTTP and HTTPS and follows up to 10 redirects, but denies localhost, private networks, carrier-grade NAT, and link-local destinations. Unspecified, multicast, reserved, and invalid destinations are always denied. This includes cloud metadata endpoints on link-local addresses. Timeout and payload-size limits are disabled until configured.

{{< notification type="warning" >}}
The secure destination defaults are intentionally backward-incompatible. Applications that previously used the built-in client to reach development servers, containers, cluster-local services, or private APIs must opt in to the required address classes or inject a custom HTTP client.
{{</ notification >}}

The built-in client resolves hostnames and checks every returned address before the initial request and each followed redirect. It also checks the concrete numeric address immediately before connecting, so a DNS change between validation and connection cannot redirect the dial to a forbidden address. Redirects are checked before their requests are sent. An allowed-host list restricts names but never overrides these address-class checks.

For production, use the narrowest practical host allowlist, as in the example above. If an application deliberately needs internal destinations, enable each class separately:

{{< code lang="go" >}}
internalHTTPClient := ferrethttp.New(
    ferrethttp.WithAllowedHosts("localhost", "api.internal.example"),
    ferrethttp.WithAllowLocalhost(true),
    ferrethttp.WithAllowPrivateNetworks(true),
    ferrethttp.WithAllowLinkLocal(true),
)
{{</ code >}}

`WithAllowLocalhost` enables only localhost names and loopback addresses. `WithAllowPrivateNetworks` enables RFC 1918, IPv6 unique-local, and carrier-grade NAT addresses; it does not enable link-local addresses. Enable `WithAllowLinkLocal` only when link-local access, including access to potential cloud metadata services, is explicitly required.

The built-in client does not inherit `HTTP_PROXY`, `HTTPS_PROXY`, or `NO_PROXY`. A proxy may resolve or connect to a different destination than the client validated. Applications that require a proxy or an exceptional destination class can inject their own `ferrethttp.Client` with `ferretnet.WithHTTPClient`; that client is responsible for equivalent destination and redirect controls.

HTTP policy is defense in depth, not a complete sandbox for arbitrary untrusted FQL. Combine it with production allowlists, execution limits, restricted module sets, and infrastructure-level egress controls.

#### Accessing HTTP from a function

Retrieve the configured HTTP client from the function context and pass the same context to the request:

{{< code lang="go" title="fetch_status.go" >}}
package request

import (
    "context"

    ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
    ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

func FetchStatus(ctx context.Context, urlArg runtime.Value) (runtime.Value, error) {
    url, err := runtime.CastArg[runtime.String](urlArg, 0)
    if err != nil {
        return runtime.None, err
    }

    client, err := ferretnet.HTTPClientFrom(ctx)
    if err != nil {
        return runtime.None, err
    }

    response, err := client.Do(ctx, &ferrethttp.Request{
        Method: "GET",
        URL:    url.String(),
    })
    if err != nil {
        return runtime.None, err
    }
    if response == nil {
        return runtime.None, runtime.Error(runtime.ErrUnexpected, "HTTP response is nil")
    }

    return runtime.NewInt(response.StatusCode), nil
}
{{</ code >}}

For public modules, use `HTTPClientFrom` instead of `net/http` or a separately constructed client. Requests then follow the host's destination, redirect, header, timeout, and payload-size policies, while the execution context continues to carry cancellation and deadlines. Use `NetworkFrom` only when a function needs the complete network service.

Both examples have valid Ferret function signatures and can be registered with `sdk.Func`. See [Writing plugins]({{< ref "/docs/guides/writing-plugins" >}}) for the complete registration pattern.

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

{{< docs-related tiles="embedding-getting-started,embedding-modules,embedding-value-encoders,embedding-programs" >}}
