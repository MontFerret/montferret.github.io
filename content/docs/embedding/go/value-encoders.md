---
title: "Value Encoders"
sidebarTitle: "Value Encoders"
weight: 80
draft: false
description: "Understand how output values are encoded, use built-in codecs, and implement custom encoders."
aliases:
    - /docs/embedding/value-encoders/
---

# Value Encoders

When a session finishes, Ferret encodes the result value into bytes using a codec selected by content type. `Run` returns a consumable `Output` handle. Collect it to read the encoded bytes and their MIME type:

These examples use the [documented Core checkout]({{< ref "/docs/embedding/go/getting-started#installation" >}}). For an existing byte-slice codec, follow the [Encoding API migration]({{< ref "/docs/migrations/encoding-api" >}}).

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    log.Fatal(err)
}
defer output.Close()

content, err := output.Collect(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println(content.Metadata.ContentType) // "application/json"
fmt.Println(string(content.Data))
{{</ code >}}

## Built-in codecs

Two codecs are registered by default:

| Content type | Format | Package |
|-------------|--------|---------|
| `application/json` | JSON (default) | `pkg/encoding/json` |
| `application/vnd.msgpack` | MessagePack | `pkg/encoding/msgpack` |

MessagePack decodes signed integers and unsigned integers up to `math.MaxInt64`
as `runtime.Int`, preserving their exact values on both 32-bit and 64-bit
systems. Unsigned integers above that limit return an overflow error. Values
encoded as MessagePack floats decode as `runtime.Float`, even when they have
no fractional part.

## Selecting the output format

Set the content type when creating a session:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithOutputContentType("application/vnd.msgpack"),
)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

output, err := session.Run(ctx)
if err != nil {
    log.Fatal(err)
}
defer output.Close()

content, err := output.Collect(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println(content.Metadata.ContentType)
// application/vnd.msgpack
{{</ code >}}

If no content type is set, the session defaults to `application/json`.

## The Codec interface

A codec combines encoding and decoding behind a content type:

{{< code lang="go" >}}
type Codec interface {
    ContentType() string
    Encoder
    Decoder
}

type Encoder interface {
    Encode(ctx context.Context, dst io.Writer, value runtime.Value) error
    EncodeWith() EncoderConfigurer
}

type Decoder interface {
    Decode(ctx context.Context, src io.Reader) (runtime.Value, error)
    DecodeWith() DecoderConfigurer
}
{{</ code >}}

### Writing and reading values

The following snippets use `encoding` from `pkg/encoding` and `json` from `pkg/encoding/json`. Pass your context and destination directly:

{{< code lang="go" >}}
if err := json.Default.Encode(ctx, writer, value); err != nil {
    return err
}

decoded, err := json.Default.Decode(ctx, reader)
if err != nil {
    return err
}
fmt.Println(decoded)
{{</ code >}}

The codec borrows I/O synchronously until the method returns. It must not retain,
close, or use the supplied reader or writer afterward. Scratch buffers may be
reused after each `Write` returns. Success completes all codec-owned writes and
flushing; callers flush, finalize, and close their own wrappers.

Failure may leave a written prefix or consume part or all of the input.
Non-nil I/O errors retain their identity. An incomplete write without an error
returns `io.ErrShortWrite` and stops without retrying. Decoding returns a fully
materialized value, may read ahead, and rejects malformed, truncated, or multiple
documents.

Use the byte helpers when your caller needs bytes:

{{< code lang="go" >}}
data, err := encoding.EncodeBytes(ctx, json.Default, value)
if err != nil {
    return err
}

decoded, err := encoding.DecodeBytes(ctx, json.Default, data)
if err != nil {
    return err
}
fmt.Println(decoded)
{{</ code >}}

`EncodeBytes` collects into a `bytes.Buffer` and returns the written prefix even
on error. `DecodeBytes` uses `bytes.NewReader`. Check encoding success before
publishing buffered bytes when publication must be atomic.

Native encoding remains eager and byte-backed. Successful zero-byte encoding
produces present empty `Content`; a failed zero-byte encoding leaves content
absent. A nonempty prefix remains available with `LengthKnown=false`.
Collected bytes survive encoder scratch reuse and result/session cleanup.
See [Output]({{< ref "/docs/embedding/go/output" >}}).

### Context and buffering

Caller contexts reach hooks and runtime capabilities. Built-in codecs first
reject nil contexts and nil or typed-nil readers/writers with errors wrapping
`runtime.ErrInvalidArgument`, before resource discovery, hooks, or I/O.

After validation, encoding completes stored-root resource discovery before
checking cancellation. An already canceled context, or cancellation during
discovery, returns before hooks, serialization, or writes. Decoding checks
cancellation immediately after validation, before hooks or reads. During active
work, both codecs check cancellation at bounded I/O boundaries, every 256
traversal values/tokens, and between large payload chunks. Arbitrary blocking I/O
and context-free custom marshalers cannot be forcibly interrupted.

JSON writes progressively, including escaped strings and Base64 payloads in
bounded chunks. Custom marshalers and generic Go-value fallbacks may buffer.
MessagePack uses the supplied I/O directly, but generic iterables still
materialize to determine container lengths. Both decoders materialize their
result; this interface makes no constant-memory guarantee. Query execution,
record iteration, and output lifecycle behavior are unchanged.

Encoding borrows the source value. A distinct acquired iterator closes once after
traversal and child encoding, including on failure or cancellation. An iterator
that is the source itself, or shares its `runtime.Resource` identity, stays with
the source's owner. Non-comparable resource aliases need `runtime.Resource` for
stable identity. Direct codec callers own cleanup of their source and child
resources. Native materialization adopts resources into its VM result.

Both built-in JSON and MessagePack adopt successful iterator yields, host
`List.At` results, and map callback values and keys before cancellation checks,
key/separator writes, or user hooks can abandon them. They also discover resources
already stored in concrete built-in arrays and objects at the root and these
acquisition boundaries, and around hooks or container failures. This iterative,
cycle-safe discovery completes even when the context is canceled. It never
reads lazy iterables, invokes host collection capabilities, marshals values, or
runs user encoding hooks. Repeated containers are rescanned so later mutations
do not hide new resources.

The codecs report discoveries to the VM result without closing source or child
resources. Comparable identities and `runtime.Resource` aliases use the result's
existing deduplication; dynamically non-comparable aliases need
`runtime.Resource` for a safe identity. Configured and pointer instances of the
built-in codecs share this behavior. Custom codecs and decorators retain the
existing ownership pre-hook; the built-in path is selected by Go type, not
content type. The original contexts still reach hooks and host capabilities.

MessagePack encodes its collected items only after traversal succeeds, then
closes a distinct iterator. A traversal failure leaves yielded resources with
the VM result for cleanup even though child encoding and its hooks never run.
Traversal errors joined with `io.EOF` retain their other causes. Iterator
cleanup errors join traversal or encoding errors before the parent's
post-hooks. Yielded-resource cleanup errors belong to result cleanup and remain
separate from encoding errors.

## Encoder and decoder hooks

Configurers build encoders and decoders with hooks attached:

{{< code lang="go" >}}
encoder := codec.EncodeWith().
    PreHook(func(ctx context.Context, value runtime.Value) error {
        return nil
    }).
    PostHook(func(ctx context.Context, value runtime.Value, err error) error {
        return nil
    }).
    Encoder()
{{</ code >}}

| Hook | Signature | When it runs |
|------|-----------|-------------|
| `PreEncoderHook` | `func(context.Context, runtime.Value) error` | Before each visited value |
| `PostEncoderHook` | `func(context.Context, runtime.Value, error) error` | After each value, including on encoding failure |
| `PreDecoderHook` | `func(context.Context) error` | Once before decoding, without consuming input |
| `PostDecoderHook` | `func(context.Context, runtime.Value, error) error` | Once after decoding |

Encoder hooks run recursively in depth-first order. Decoder post-hooks receive
the decoded value on success and `runtime.None` on decoding failure. Hooks of the
same kind run in registration order and stop at the first hook error. A pre-hook
error skips that operation and its post-hooks; ancestor encoder post-hooks still
receive descendant errors. Operation and hook errors preserve their causes.

Decoder hooks work with the operation and decoded value. For example, reject a
successfully decoded `none` without hiding a decoding failure:

{{< code lang="go" >}}
decoder := codec.DecodeWith().
    PreHook(func(ctx context.Context) error {
        return ctx.Err()
    }).
    PostHook(func(_ context.Context, value runtime.Value, decodeErr error) error {
        if decodeErr == nil && value == runtime.None {
            return errors.New("document must not be none")
        }

        return nil
    }).
    Decoder()

decoded, err := decoder.Decode(ctx, reader)
if err != nil {
    return err
}
fmt.Println(decoded)
{{</ code >}}

Built-in JSON and MessagePack configurers publish independent snapshots.
Adding hooks to a configurer or deriving another encoder does not change an
encoder or decoder already returned by it. Copies of a built-in codec retain
their published configuration. Configure before sharing; configurer mutation
and any mutable state captured by hooks require caller synchronization.

## The Registry

The `encoding.Registry` stores codecs by MIME-normalized content type:

{{< code lang="go" >}}
// Create a registry seeded with codecs
registry := encoding.NewRegistry(jsonCodec, msgpackCodec)

// Or start empty and register individually
empty := encoding.NewEmptyRegistry()
if err := empty.Register(myCodec); err != nil {
    return err
}
{{</ code >}}

### Registry methods

| Method | Returns | Purpose |
|--------|---------|---------|
| `Register(codec)` | `error` | Store a codec by its content type |
| `Codec(contentType)` | `Codec, error` | Look up a full codec |
| `Encoder(contentType)` | `Encoder, error` | Look up an encoder |
| `Decoder(contentType)` | `Decoder, error` | Look up a decoder |
| `Clone()` | `*Registry` | Copy registrations while retaining the codec values |

Content types are normalized using MIME media type parsing, so `application/json` and `application/json; charset=utf-8` resolve to the same codec. Registering the same normalized type replaces the earlier entry.

`NewRegistry` ignores nil seed codecs and panics if a seed has an invalid content type. Use `NewEmptyRegistry` and check each `Register` error for runtime-supplied configuration. A clone has its own registration map; it does not deep-copy custom codec state.

## Registering codecs on the engine

Add or override a single codec:

{{< code lang="go" >}}
engine, err := ferret.New(
    ferret.WithEncodingCodec("text/csv", myCSVCodec),
)
{{</ code >}}

Replace the entire registry:

{{< code lang="go" >}}
registry := encoding.NewRegistry(jsonCodec, myCSVCodec)

engine, err := ferret.New(
    ferret.WithEncodingRegistry(registry),
)
{{</ code >}}

When you use `WithEncodingRegistry`, only the codecs in the provided registry are available. The default JSON and MessagePack codecs are not included unless you add them yourself.

## Complete example

This plain-text codec writes directly, preserves short-write errors, and runs
configured hooks. Its decoder buffers the input until EOF because plain text
has no document delimiter. Applications needing interruption during large
custom operations should add bounded cancellation checks between chunks.

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "log"
    "reflect"
    "slices"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/encoding"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
)

type plainTextCodec struct {
    encPre  []encoding.PreEncoderHook
    encPost []encoding.PostEncoderHook
    decPre  []encoding.PreDecoderHook
    decPost []encoding.PostDecoderHook
}

func (plainTextCodec) ContentType() string { return "text/plain" }

func checkContext(ctx context.Context, input any) error {
    if ctx == nil || input == nil {
        return fmt.Errorf("%w: nil codec context or I/O", runtime.ErrInvalidArgument)
    }

    // An io.Reader or io.Writer interface can hold a nil pointer or other nil value.
    reflected := reflect.ValueOf(input)
    switch reflected.Kind() {
    case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
        if reflected.IsNil() {
            return fmt.Errorf("%w: nil codec I/O", runtime.ErrInvalidArgument)
        }
    }

    return ctx.Err()
}

func (c plainTextCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
    if err := checkContext(ctx, dst); err != nil {
        return err
    }
    if value == nil {
        value = runtime.None
    }

    for _, hook := range c.encPre {
        if err := hook(ctx, value); err != nil {
            return err
        }
    }

    text := value.String()
    err := ctx.Err()
    if err == nil {
        var n int
        n, err = io.WriteString(dst, text)
        if err == nil && n != len(text) {
            err = io.ErrShortWrite
        }
        if err == nil {
            err = ctx.Err()
        }
    }
    for _, hook := range c.encPost {
        if hookErr := hook(ctx, value, err); hookErr != nil {
            return errors.Join(err, hookErr)
        }
    }
    return err
}

func (c plainTextCodec) Decode(ctx context.Context, src io.Reader) (runtime.Value, error) {
    if err := checkContext(ctx, src); err != nil {
        return runtime.None, err
    }
    for _, hook := range c.decPre {
        if err := hook(ctx); err != nil {
            return runtime.None, err
        }
    }
    // Plain text has no document boundary: this example buffers until EOF.
    var data []byte
    err := ctx.Err()
    if err == nil {
        data, err = io.ReadAll(src)
    }
    if err == nil {
        err = ctx.Err()
    }
    var value runtime.Value = runtime.None
    if err == nil {
        value = runtime.NewString(string(data))
    }
    for _, hook := range c.decPost {
        if hookErr := hook(ctx, value, err); hookErr != nil {
            return runtime.None, errors.Join(err, hookErr)
        }
    }
    return value, err
}

func (c plainTextCodec) EncodeWith() encoding.EncoderConfigurer {
    c.encPre = slices.Clone(c.encPre)
    c.encPost = slices.Clone(c.encPost)
    return &plainTextEncoderConfigurer{codec: c}
}

type plainTextEncoderConfigurer struct { codec plainTextCodec }

func (e *plainTextEncoderConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
    if hook != nil { e.codec.encPre = append(e.codec.encPre, hook) }
    return e
}
func (e *plainTextEncoderConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
    if hook != nil { e.codec.encPost = append(e.codec.encPost, hook) }
    return e
}
func (e *plainTextEncoderConfigurer) Encoder() encoding.Encoder {
    c := e.codec
    c.encPre = slices.Clone(c.encPre)
    c.encPost = slices.Clone(c.encPost)
    return c
}

func (c plainTextCodec) DecodeWith() encoding.DecoderConfigurer {
    c.decPre = slices.Clone(c.decPre)
    c.decPost = slices.Clone(c.decPost)
    return &plainTextDecoderConfigurer{codec: c}
}

type plainTextDecoderConfigurer struct { codec plainTextCodec }

func (d *plainTextDecoderConfigurer) PreHook(hook encoding.PreDecoderHook) encoding.DecoderConfigurer {
    if hook != nil { d.codec.decPre = append(d.codec.decPre, hook) }
    return d
}
func (d *plainTextDecoderConfigurer) PostHook(hook encoding.PostDecoderHook) encoding.DecoderConfigurer {
    if hook != nil { d.codec.decPost = append(d.codec.decPost, hook) }
    return d
}
func (d *plainTextDecoderConfigurer) Decoder() encoding.Decoder {
    c := d.codec
    c.decPre = slices.Clone(c.decPre)
    c.decPost = slices.Clone(c.decPost)
    return c
}

func main() {
    engine, err := ferret.New(
        ferret.WithEncodingCodec("text/plain", plainTextCodec{}),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    plan, err := engine.Compile(
        context.Background(),
        ferret.NewAnonymousSource(`return "Hello, Ferret!"`),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    session, err := plan.NewSession(
        context.Background(),
        ferret.WithOutputContentType("text/plain"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    output, err := session.Run(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    defer output.Close()

    content, err := output.Collect(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(content.Metadata.ContentType)
    // text/plain
    fmt.Println(string(content.Data))
    // Hello, Ferret!
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-go-output,migrations-encoding-api,embedding-go-programs,embedding-go-modules" >}}
