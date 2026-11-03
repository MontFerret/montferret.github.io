---
title: "Output"
sidebarTitle: "Output"
weight: 25
draft: false
description: "Collect owned content or consume borrowed bytes, and settle output with its errors and cleanup."
---

# Output

`Engine.Run` and `Session.Run` return `(ferret.Output, error)`. Output is a caller-owned, one-shot handle: use `Collect` when you need owned bytes, or `Consume` when you want to forward or discard them. Both report terminal failures and finalize the handle.

These contracts target the [Core development checkout]({{< ref "/docs/embedding/go/getting-started#installation" >}}) and Universal API `v1.0.0-alpha.21`. See the [migration guide]({{< ref "/docs/migrations/consumable-output" >}}) for buffered-output applications. For remote runtimes, use the dependency set in the [Wire compatibility table]({{< ref "/docs/embedding/wire#run-a-host-and-client" >}}).

## Output methods

`ferret.Output` is an interface. Choose one terminal consumption method per handle:

| Method | Purpose |
| --- | --- |
| `Metadata() ferret.Metadata` | Read the immutable complete-payload descriptor |
| `Collect(ctx) (*ferret.Content, error)` | Take owned bytes and observe the terminal outcome |
| `Consume(ctx, ferret.Consumer) error` | Deliver borrowed bytes and observe the terminal outcome |
| `Close() error` | Abandon unread content or stop active delivery; report cleanup failures |

## Collect owned content

Collect before decoding, printing, or returning data from a helper that closes its execution resources:

{{< code lang="go" >}}
package example

import (
    "context"
    "errors"

    "github.com/MontFerret/ferret/v2"
)

func Query(ctx context.Context) (content *ferret.Content, err error) {
    engine, err := ferret.New()
    if err != nil {
        return nil, err
    }
    defer func() { err = errors.Join(err, engine.Close()) }()

    output, err := engine.Run(ctx, ferret.NewAnonymousSource(`RETURN 42`))
    if err != nil {
        return nil, err // preparation failed; no handle
    }
    defer func() { err = errors.Join(err, output.Close()) }()

    return output.Collect(ctx)
}
{{< /code >}}

`Collect` returns `(*ferret.Content, error)`. Inspect content independently of the error: encoding, hooks, or cleanup can fail while bytes remain available. Content belongs to the recipient and stays valid after output and parent resources close. Native collection transfers the owned encoded buffer directly; it does not copy the entire payload into a second buffer.

## Consume borrowed bytes

`Consume` invokes a synchronous `ferret.Consumer`. The callback receives a context and a borrowed, read-only byte slice. Bytes are valid only during that callback; copy any bytes you retain. Chunks arrive in order without overlap. Do not rely on chunk boundaries lining up with characters, lines, records, or values, and do not treat an empty chunk as EOF.

This writer example is separate from HTTP response handling:

{{< code lang="go" >}}
package example

import (
    "context"
    "errors"
    "io"

    "github.com/MontFerret/ferret/v2"
)

func WriteQuery(ctx context.Context, engine *ferret.Engine, dst io.Writer) (err error) {
    output, err := engine.Run(ctx, ferret.NewAnonymousSource(`RETURN 42`))
    if err != nil {
        return err
    }
    defer func() { err = errors.Join(err, output.Close()) }()

    return output.Consume(ctx, func(ctx context.Context, chunk []byte) error {
        if err := ctx.Err(); err != nil {
            return err
        }
        n, err := dst.Write(chunk)
        if err == nil && n != len(chunk) {
            return io.ErrShortWrite
        }
        return err
    })
}
{{< /code >}}

A consumer error stops delivery and is preserved with other observed failures. Panic unwinding finalizes the handle and propagates the panic. No callback remains active after `Consume` returns. Arbitrary callback code must cooperate with cancellation; the runtime cannot forcibly interrupt it.

Native execution and encoding remain **eager**. Native output delivers present content in one synchronous callback today, but the public API promises no chunk boundaries. Consumption does not imply incremental query evaluation or encoding. Writer consumption can write bytes before a terminal error is reported. An HTTP handler should [collect before writing]({{< ref "/docs/embedding/go/application-guide" >}}) so it can still choose an error status.

## Observe side-effect-only execution

Even if you only need query side effects, consume the handle to observe completion. Closing unread output is abandonment, not a successful execution check:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    return err
}
defer output.Close()

return output.Consume(ctx, func(context.Context, []byte) error {
    return nil // discard bytes; still observe terminal failures
})
{{< /code >}}

## Metadata and content presence

`output.Metadata()` returns an immutable `ferret.Metadata` descriptor. It is available before consumption, during callbacks, and after closure. It contains `ContentType`, `Length`, and `LengthKnown`. Length describes the complete encoded payload, not progress or a record count; use it only when `LengthKnown` is true and do not allocate without limits based on it.

Metadata and byte length do not establish content presence:

| Collection result | Meaning |
| --- | --- |
| `content == nil` | No content is available. Native consumption invokes no callback. |
| `content != nil`, `len(content.Data) == 0` | Content is present and empty. Native consumption invokes an empty callback. |
| `content != nil`, nonempty data | Encoded content is available, possibly alongside an error. |

A non-nil `&ferret.Content{}` is present content, including when `Data` is nil. Successful native encoding records the actual codec's media type and a known length equal to the encoded byte count. Successful encoding with zero writes is present-empty content.

An encoder may write a prefix before returning an error. Native preserves a nonempty prefix as content, with the selected media type and unknown total length. A zero-byte encoding failure leaves content absent. Later hook or cleanup failures do not erase an established descriptor. For portable implementations, partial collected data can be a prefix; metadata continues to describe the original complete payload.

## Serialize materialized content

Serialize `Content` explicitly, rather than a live output handle:

{{< code lang="go" >}}
content, err := output.Collect(ctx)
if err != nil {
    return err
}
encoded, err := json.Marshal(content)
{{< /code >}}

For JSON-encoded FQL `RETURN 42`, this produces:

```json
{"metadata":{"contentType":"application/json","length":2,"lengthKnown":true},"data":"NDI="}
```

The standard Go byte-slice representation encodes `Data` as base64. It does not embed those bytes as raw JSON. Absent `*Content` serializes as `null`. In present content, nil `Data` serializes as `null` and an empty non-nil slice as `""`; neither makes the containing content absent.

To send the query's encoded payload itself, write `content.Data` and use `content.Metadata.ContentType`. That is different from serializing the content envelope.

## Error timing

Preparation errors are returned immediately by `Run`: compilation, options, session construction, before-run hooks, and context validation before VM entry. Once admitted to the VM, Native returns a usable handle and a nil Run error, including eager VM, encoding, after-run-hook, and result-cleanup failures. `Collect` and `Consume` report those recorded failures with delivery, context, and cleanup errors.

Native closes VM results after materialization. Convenience `Engine.Run` also closes its temporary session and plan before returning, recording their cleanup failures. An unread handle retains bytes, metadata, its invocation context, and recorded outcomes rather than VM resources or execution permits. Caller-created sessions remain caller-owned and can be reused sequentially after output settlement.

Ordinary after-run hooks receive the primary execution error and run before encoding. Debugger completion keeps its separate ordering and retains materialized `*Content`, including available content alongside an error; it does not publish a consumable handle.

Use `errors.Is` and `errors.As` to traverse joined causes. Through [UAPI]({{< ref "/docs/embedding/go/universal-api" >}}), portable diagnostics and source coordinates remain available from `Consume`, `Collect`, and `Close` errors.

## Context lifetime and rejected consumption

Keep the **original Run context** alive until output settlement. A context substituted by an execution hook affects execution only. The callback's effective context observes both the original invocation and consumption cancellations and their earliest deadline. A later context cannot extend or revive the invocation.

Admission checks a non-nil context, then a non-nil consumer for `Consume`, then both context cancellations before claiming the handle. Rejected calls return nil collected content or invoke no consumer, leave the handle unclaimed, and do not disturb active consumption. If only the consumption context is canceled, retry with a live context while the invocation remains live.

If the invocation is canceled, the admission error preserves its recorded execution, encoding, hook, and cleanup failures. This holds whether consumption uses the original context, a separate live context, or another canceled context. Distinct cancellation causes remain traversable with `errors.Is`, along with `context.Canceled` or `context.DeadlineExceeded`. Rejection does not erase the recorded outcome or settle the handle; close it when abandoning it.

After successful admission, any failure or cancellation finalizes the handle. Consumption cannot resume or repeat. Valid competing or reentrant consumption matches `ferret.ErrOutputInUse` while delivery is active; once stopping or finalization begins it matches `ferret.ErrOutputClosed`. Invalid or canceled calls retain their validation precedence even then.

## Close and settlement

`Close` abandons unread output without draining it. It returns only the recorded cleanup outcome, usually nil; it does not report execution failures or certify successful execution. Consume or collect to observe those failures. Automatic finalization also leaves cleanup errors available to later `Close` calls.

During active consumption, `Close` requests cancellation before waiting for callbacks and finalization. Borrowed bytes stay valid until their callback returns. Concurrent and repeated closers converge on the same finalization and cleanup outcome. If explicit closure wins before the consumption outcome is committed, that operation also matches `ferret.ErrOutputClosed`. Closure after commitment only waits; automatic finalization does not relabel ordinary errors as explicit closure.

**Do not call synchronous `Close` from the output's own consumer.** It waits for that callback and would deadlock. Metadata reads and rejected reentrant consumption are supported from the callback. Output closure does not close caller-owned sessions or borrowed parents, and it does not modify bytes transferred by collection. Metadata remains readable after closure.

## Next steps

{{< docs-related tiles="embedding-go-executing,embedding-go-value-encoders,embedding-go-universal-api,migrations-consumable-output" >}}
