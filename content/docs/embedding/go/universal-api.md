---
title: "Universal API"
weight: 150
draft: false
description: "Expose a Native Ferret engine through the portable Universal API."
---

# Use the Universal API

Use `github.com/MontFerret/ferret/v2/uapi`, Ferret's official adapter, when
your integration accepts `github.com/MontFerret/api` interfaces. For Native
embedding, use the root `github.com/MontFerret/ferret/v2` package and
`ferret.New`. Use `uapi.New` to create and own a Native engine:

{{< code lang="go" >}}
package example

import (
    "context"

    "github.com/MontFerret/api"
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/uapi"
)

func Run(ctx context.Context) (*api.Output, error) {
    portable, err := uapi.New(ferret.WithParam("value", 41))
    if err != nil {
        return nil, err
    }
    defer portable.Close()

    return portable.Run(ctx, api.NewAnonymousSource("RETURN @value + 1"))
}
{{< /code >}}

The runtime owns the engine created by `New`; closing the runtime closes that
engine. Pass Native options to configure modules, host services, codecs, and
defaults. Construction failures return a nil runtime and a projected Native
error; Native handles rollback. Nil options, including `uapi.New(nil)`,
are skipped by Native. Call `uapi.New()` to use Native defaults.

Use `Wrap` when the caller supplies and owns the Native engine:

{{< code lang="go" >}}
native, err := ferret.New()
if err != nil {
    return err
}
defer native.Close()

var portable api.Runtime = uapi.Wrap(native)
{{< /code >}}

`Wrap` acquires no resources and leaves engine cleanup to its caller.
`uapi.Wrap(nil)` panics. Both constructors return `*uapi.Runtime`,
which implements `api.Runtime`. Both constructors live in `uapi`; the root
`ferret` package exposes the Native API.

The adapter translates portable source, options, and diagnostics while Native
Ferret performs compilation, execution, debugging, and resource cleanup. Native
and Universal option types remain separate.

## Reuse a plan

`portable.Compile(ctx, source, options...)` returns an `api.Plan`.
Create independent sessions with `plan.NewSession(ctx, options...)`, run them
with `session.Run(ctx)`, and close each session before closing the plan.
Ordinary sessions also support sequential runs when their environment remains
unchanged.

For debugging, use `portable.CompileDebug`, then `plan.NewDebugSession`.
The returned `api/debugger.Session` supports entry, breakpoints, stepping,
frames, variables, evaluation, pause, and termination through Native debugging.

## Configure portable execution

| Option | Behavior |
| --- | --- |
| `api.WithOptimizationLevel` | Selects None, Basic, or Full for one compilation. Aggressive and unknown levels fail. |
| `api.WithParam`, `api.WithParams` | Supply host parameters using Native value conversion. Later settings override earlier keys; maps merge. |
| `api.WithOutputContentType` | Selects a codec registered on the Native engine. JSON is the default. |
| `api.WithFSRoot` | Selects a session-owned filesystem root while retaining the engine's read-only policy. |

Omitted optimization inherits the engine default. Plan setters queue supported
levels; Native validates them for the compilation mode. Debug compilation accepts
omission or None. Basic/Full fails during Native option application, even if
followed by None. Session options apply to both ordinary and debug sessions, and
to the convenience `Runtime.Run` call.

Non-nil option callbacks run once, in order, independently of the operation
context. Session setters queue Native options and return nil; Native owns
parameter conversion and validation. If callbacks return errors, they are joined
and delegation stops. Cancellation is included only if a callback returns it.
Otherwise Native receives the original caller context and owns context validation
and cancellation. Native applies queued options and joins option-application
errors before acquiring session resources. `Runtime.Run` may compile before
session validation.

The portable contract allows runtime-specific validation at the point of use.
Native checks output codec availability during result encoding, after the query
has run. An unavailable codec therefore returns an operation error even though
callbacks and session creation succeeded and query side effects may have occurred.
The adapter preserves Native result and session cleanup on this path.

Parameter maps are converted after all portable callbacks finish. Changes to
those maps during translation are therefore visible to conversion. Native
rejects blank output content types and filesystem roots, and rejects an empty
name or nil value for a single parameter. In parameter maps, nil becomes `none`;
nil and empty maps are no-ops when applied. Runtime-specific extension options
must reject incompatible option targets.

Configure Native modules, host services, codecs, and other engine-specific
features through `New` options or before wrapping an existing engine. See [configuration](../configuration/) and
[host values](../host-values/).

## Own the lifetime

For runtimes created by `New`, `Runtime.Close` delegates to the owned Native
engine. Native releases its resources and rejects subsequent runtime operations.
For `Wrap`, Close returns nil and leaves the adapter, engine, and directly created
plans usable. Multiple adapters may borrow one engine independently. The engine's
owner remains responsible for closing it; external closure is visible to all
borrowers.

Plan, session, and debug-session close calls delegate to Native. Plan close
rejects new sessions and wakes capacity waiters without waiting for outstanding
constructors. Use caller contexts to cancel work, settle it, then close sessions
before plans and their owning runtime or borrowed Native engine. Native engine
close does not wait for operations already started or close descendants. Parent
close does not implicitly cancel caller-owned work. Directly created sessions
and debug sessions retain their own lifecycle after plan closure.
The portable contract requires cleanup of owned resources and rejection of new
descendants after plan closure; it does not require waiting for constructors or
coordinating descendant cleanup. Callers coordinate cleanup when descendants use
parent-owned resources.
`Runtime.Run` delegates temporary session and plan cleanup to Native.

Settle ordinary Run calls before closing their session. Debug Close terminates
and waits for active commands. Repeated and concurrent Close calls preserve the
Native cleanup result and diagnostic information; projected error wrappers need
not have identical pointers. Close hooks must not recursively close their own
Native object. Context arguments must be non-nil.

## Preserve output and diagnostics

`Runtime.Run` and `Session.Run` return `(*api.Output, error)`. A nil output
means no output was produced. A non-nil output represents a produced value,
including empty output, and may accompany an error from a hook or cleanup step.
A zero-valued `Output` is not an absence sentinel. Inspect output independently
of the error:

{{< code lang="go" >}}
output, err := portable.Run(ctx, api.NewAnonymousSource("RETURN 42"))
if output != nil {
    fmt.Printf("%s: %s\n", output.ContentType, output.Content)
}
if err != nil {
    return err
}
{{< /code >}}

Output bytes belong to the caller, survive cleanup, and need no Close call.
Debugger commands can likewise return a completion event and output alongside
an error.

Use `errors.As` to obtain `github.com/MontFerret/api/diagnostics.Diagnostics`.
The projection preserves source identity and byte coordinates, while
`errors.Is` and `errors.As` still reach Native errors and cancellation causes.
See [executing queries](../executing/) for Native output and lifecycle details.
