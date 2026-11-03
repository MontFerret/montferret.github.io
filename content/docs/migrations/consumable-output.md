---
title: "Consumable output"
sidebarTitle: "Consumable output"
weight: 20
draft: false
description: "Migrate buffered Native and UAPI output to the one-shot Output and detached Content API."
---

# Migrate to consumable output

The Core development API adopts Universal API `v1.0.0-alpha.21`: Run returns a consumable handle, and collected data lives in `Content`. Use the [development checkout instructions]({{< ref "/docs/embedding/go/getting-started#installation" >}}) for these examples. For remote runtimes, use the published dependency set in the [Wire compatibility table]({{< ref "/docs/embedding/wire#run-a-host-and-client" >}}).

This migration changes the Go output contract, error timing, and ownership. It does not change FQL semantics or require rewriting queries. The [Output reference]({{< ref "/docs/embedding/go/output" >}}) owns the complete lifecycle rules.

## Replace buffered field access

Previously, Native and UAPI returned `*Output` structs with `Content` and `ContentType` fields:

{{< code lang="go" >}}
// Before: buffered output contract.
output, err := engine.Run(ctx, source)
if err != nil {
    return err
}
fmt.Printf("%s: %s\n", output.ContentType, output.Content)
{{< /code >}}

Now collect the handle before reading bytes:

{{< code lang="go" >}}
// After: Native development API.
output, err := engine.Run(ctx, source)
if err != nil {
    return err // preparation failure
}
defer output.Close()

content, err := output.Collect(ctx)
if content != nil {
    fmt.Printf("%s: %s\n", content.Metadata.ContentType, content.Data)
}
if err != nil {
    return err // terminal failure; available content was preserved
}
{{< /code >}}

| Buffered contract | Consumable contract |
| --- | --- |
| Native Run: `(*ferret.Output, error)` | `(ferret.Output, error)` |
| UAPI Run: `(*api.Output, error)` | `(api.Output, error)` |
| `output.Content` | `content.Data`, after `Collect` |
| `output.ContentType` | `output.Metadata().ContentType` or `content.Metadata.ContentType` |
| Nil output means absent content | A usable handle may contain no content; nil collected `*Content` means absence |
| Bytes read directly, no output Close | Terminal `Collect` or `Consume`; `Close` to abandon unread output |
| Runtime failures checked from Run | Preparation checked from Run; terminal failures checked from consumption |

Root `ferret` and `pkg/encoding` alias the API's `Output`, `Content`, `Metadata`, and `Consumer` types. Output sentinels retain canonical API identity. Existing compatibility APIs and `vm.Result` remain unchanged.

## Settle every execution

For side effects, replace an ignored Run result with terminal discard consumption:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    return err
}
defer output.Close()
return output.Consume(ctx, func(context.Context, []byte) error { return nil })
{{< /code >}}

For writers, use borrowed-buffer [Consume]({{< ref "/docs/embedding/go/output#consume-borrowed-bytes" >}}). Do not retain callback bytes without copying. For HTTP handlers, collect before writing so a terminal failure can still choose an error status. Native execution and encoding are still eager; the new consumer API does not introduce incremental query execution.

Helpers that close their own runtime or session should return detached `*Content`:

{{< code lang="go" >}}
package example

import (
    "context"
    "errors"

    "github.com/MontFerret/ferret/v2"
)

func RunPlan(ctx context.Context, plan *ferret.Plan) (content *ferret.Content, err error) {
    session, err := plan.NewSession(ctx)
    if err != nil {
        return nil, err
    }
    defer func() { err = errors.Join(err, session.Close()) }()

    output, err := session.Run(ctx)
    if err != nil {
        return nil, err
    }
    defer func() { err = errors.Join(err, output.Close()) }()
    return output.Collect(ctx)
}
{{< /code >}}

Directly created sessions remain caller-owned and reusable sequentially after output settlement. Convenience Engine.Run closes its temporary session and plan before publishing the handle. Collected bytes survive all of that cleanup.

## Move runtime-error assertions

SDK tests must collect or consume before checking execution errors:

{{< code lang="go" >}}
output, err := harness.Run(t.Context(), `RETURN 1 / 0`)
if err != nil {
    t.Fatal(err) // unexpected preparation failure
}
defer output.Close()

_, err = output.Collect(t.Context())
if !errors.Is(err, runtime.ErrDivisionByZero) {
    t.Fatalf("expected division by zero, got %v", err)
}
{{< /code >}}

UAPI delegates delivery and ownership to Native and projects portable diagnostics from `Consume`, `Collect`, and `Close`. Continue using `errors.As` for diagnostics and `errors.Is` for source execution, cancellation, and cleanup causes. Debugger completion retains materialized `*Content` and can still return an event with available content alongside an error.

## Keep the Run context alive

Do not cancel a helper's Run context and return a live output for somebody else to consume. Collect inside that context's lifetime or let the caller own both context and handle. A consumption context can shorten, but cannot extend or revive, the original invocation.

Rejected consumption validates nil context and nil consumer first, then both context cancellations. It leaves the handle unclaimed and does not invoke a callback. A canceled consumption context can be replaced with a live one while the invocation is live. Once admitted, any failure finalizes the output; consumption cannot be retried.

Invocation cancellation preserves recorded operation and cleanup failures even when callers reuse the original canceled context for Collect or Consume. The same holds with a live consumption context or another canceled context with a distinct cause. Use `errors.Is` to check all joined causes. Close afterward to abandon a rejected output; repeated Close reports cleanup failures only.

Concurrent Close cancels active consumption and waits for callbacks and finalization. Borrowed bytes remain valid until the callback returns. Never call Close synchronously from that output's own consumer.

## Update serialization and encoders

Serialize detached Content explicitly. Its JSON envelope nests `metadata` and stores `data` using Go's standard byte-slice encoding, usually base64. Nil and empty slices remain distinct (`null` and `""`), while a nil content pointer means absent content. To return the query payload directly, write `content.Data` rather than serializing this envelope.

The separate [Encoding API migration]({{< ref "/docs/migrations/encoding-api" >}}) covers reader/writer method signatures, byte helpers, and changed encoder and decoder hooks. Custom encoders, including configured variants, receive a caller context and borrowed synchronous writer. They must not retain or close I/O, and must preserve write errors. Native collects the writes into its own buffer; nonempty error prefixes remain present with unknown total length, while zero-byte failures leave content absent. Successful zero-byte encoding remains present empty content. See [Value Encoders]({{< ref "/docs/embedding/go/value-encoders" >}}).

## Check integration boundaries

Update each Native or UAPI caller that decodes, prints, forwards, or ignores a Run result. Services retaining or replaying content should materialize it. Tools, SDK tests, and downstream runtime adapters need compatible contracts before adopting this Core development API. Keep transport and release pins coordinated. Wire hosts and clients must both support output framing version 1; incompatible peers fail during Connect, with no fallback to buffered inline output. Upgrade both sides together using the published dependency set in the [Wire example]({{< ref "/docs/embedding/wire" >}}).

{{< docs-related tiles="embedding-go-output,embedding-go-getting-started,embedding-go-universal-api,migrations-encoding-api,embedding-wire,migrations-v1-to-v2" >}}
