---
title: "Encoding API"
sidebarTitle: "Encoding API"
weight: 25
draft: false
description: "Migrate byte-slice codecs to context-aware reader/writer methods and value-aware hooks."
---

# Migrate the Encoding API

Ferret codecs accept a caller context and borrowed `io.Writer` or `io.Reader`.
Use the [documented Core checkout]({{< ref "/docs/embedding/go/getting-started#installation" >}})
for these examples. The [Value Encoders guide]({{< ref "/docs/embedding/go/value-encoders" >}})
describes the current contract and provides a complete custom codec.

This changes Go codec implementations and direct callers. FQL queries do not
need rewriting. SDK host-value codecs and bytecode serialization have separate
contracts. Applications using the earlier buffered Run result should also
follow the [consumable output migration]({{< ref "/docs/migrations/consumable-output" >}}).

## Update direct codec calls

| Earlier method | Current method |
| --- | --- |
| `Encode(value) ([]byte, error)` | `Encode(ctx, dst, value) error` |
| `Decode(data) (runtime.Value, error)` | `Decode(ctx, src) (runtime.Value, error)` |
| Collect bytes from `Encode` | `encoding.EncodeBytes(ctx, encoder, value)` |
| Pass bytes to `Decode` | `encoding.DecodeBytes(ctx, decoder, data)` |

Pass a writer or reader when the application already has one:

{{< code lang="go" >}}
if err := codec.Encode(ctx, writer, value); err != nil {
    return err
}

decoded, err := codec.Decode(ctx, reader)
if err != nil {
    return err
}
fmt.Println(decoded)
{{< /code >}}

Use the byte helpers when callers need an owned byte slice:

{{< code lang="go" >}}
package example

import (
    "context"

    "github.com/MontFerret/ferret/v2/pkg/encoding"
    "github.com/MontFerret/ferret/v2/pkg/encoding/json"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

func EncodeJSON(ctx context.Context, value runtime.Value) ([]byte, error) {
    return encoding.EncodeBytes(ctx, json.Default, value)
}

func DecodeJSON(ctx context.Context, data []byte) (runtime.Value, error) {
    return encoding.DecodeBytes(ctx, json.Default, data)
}
{{< /code >}}

`EncodeBytes` returns any written prefix alongside an error. Do not publish it
as a complete document unless encoding succeeded. `DecodeBytes` wraps the input
in a reader. There is no adapter interface for the old byte-slice methods;
custom codecs must implement the new methods, including configured variants.

## Update custom codecs and I/O ownership

The interfaces now require:

{{< code lang="go" >}}
Encode(ctx context.Context, dst io.Writer, value runtime.Value) error
Decode(ctx context.Context, src io.Reader) (runtime.Value, error)
{{< /code >}}

Borrow I/O synchronously until the call returns. Never close or retain it or use
it from background work after returning. Finish codec-owned writes and flushing
before returning success. The caller still flushes, finalizes, and closes its
own wrappers, such as a `bufio.Writer`.

Preserve non-nil reader and writer errors so `errors.Is` can find them. If a
write returns fewer bytes than requested with no error, return
`io.ErrShortWrite` and stop without retrying. Failed encoding may leave a prefix;
failed decoding may have consumed input. Buffer before publishing when the
destination requires a complete document.

Reject nil contexts and nil I/O with an error wrapping
`runtime.ErrInvalidArgument`, including interfaces holding typed-nil I/O.
Propagate the caller context to hooks and host capabilities. Custom traversal
and large operations need bounded cancellation checks. Context cannot forcibly
interrupt an arbitrary blocking reader, writer, or context-free marshaler.

## Update encoder and decoder hooks

| Hook | Earlier arguments | Current arguments |
| --- | --- | --- |
| `PreEncoderHook` | `runtime.Value` | `context.Context, runtime.Value` |
| `PostEncoderHook` | `runtime.Value, error` | `context.Context, runtime.Value, error` |
| `PreDecoderHook` | `[]byte` | `context.Context` |
| `PostDecoderHook` | `[]byte, error` | `context.Context, runtime.Value, error` |

Encoder hooks still run recursively for visited values. Decoder hooks run once
per operation: the pre-hook runs before reads, and the post-hook receives the
materialized value on decoding success or `runtime.None` on decoding failure.
Decoder hooks no longer receive the source bytes. If an application needs those
bytes, capture them at its own reader boundary.

Hooks of the same kind run in registration order and stop at the first hook
failure. A pre-hook failure skips the operation and its post-hooks; ancestor
encoder post-hooks still receive descendant failures. Preserve both operation
and hook causes when joining errors. Custom encoders must retain configured
hooks; Native uses an encoding pre-hook for ownership tracking with custom
codecs and decorators.

Built-in configurers publish independent snapshots. Adding hooks later does
not mutate an encoder or decoder already returned. Synchronize mutable state
captured by hooks when sharing configured codecs.

## Preserve output behavior

Reader/writer codecs do not make Native query execution lazy. Native still
encodes eagerly into an owned buffer before returning its Output handle.
Successful zero-byte encoding produces present, empty Content. A zero-byte
encoding failure leaves content absent; a nonempty error prefix remains
available with `LengthKnown=false`. Consume or collect to observe the encoding
and cleanup outcome, and inspect content independently of the terminal error.

Both built-in decoders return fully materialized values and may read ahead.
JSON can encode generic iterables progressively; MessagePack materializes them
to determine container lengths. Reader/writer signatures do not promise
constant memory usage or incremental query evaluation.

Direct codec callers retain ownership of source and child resources. Native
result processing adopts resources discovered during built-in encoding and
closes its result after materialization. Keep the original Run context alive
until output settlement. See [Output]({{< ref "/docs/embedding/go/output" >}})
for content ownership, terminal errors, and closure.

{{< docs-related tiles="embedding-go-value-encoders,embedding-go-output,migrations-consumable-output,embedding-go-getting-started" >}}
