---
title: "Value Encoders"
sidebarTitle: "Value Encoders"
weight: 55
draft: false
description: "Understand how output values are encoded, use built-in codecs, and implement custom encoders."
---

# Value Encoders

When a session finishes, Ferret encodes the result value into bytes using a codec selected by content type. The `Output` struct carries the encoded bytes and their MIME type:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println(output.ContentType) // "application/json"
fmt.Println(string(output.Content))
{{</ code >}}

## Built-in codecs

Two codecs are registered by default:

| Content type | Format | Package |
|-------------|--------|---------|
| `application/json` | JSON (default) | `pkg/encoding/json` |
| `application/vnd.msgpack` | MessagePack | `pkg/encoding/msgpack` |

## Selecting the output format

Set the content type when creating a session:

{{< code lang="go" >}}
session, err := plan.NewSession(ctx,
    ferret.WithOutputContentType("application/vnd.msgpack"),
)
if err != nil {
    log.Fatal(err)
}

output, err := session.Run(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println(output.ContentType)
// application/vnd.msgpack
{{</ code >}}

If no content type is set, the session defaults to `application/json`.

## The Codec interface

A codec combines encoding and decoding behind a single content type:

{{< code lang="go" >}}
type Codec interface {
    ContentType() string
    Encoder
    Decoder
}
{{</ code >}}

### Encoder

{{< code lang="go" >}}
type Encoder interface {
    Encode(value runtime.Value) ([]byte, error)
    EncodeWith() EncoderConfigurer
}
{{</ code >}}

`Encode` converts a runtime value to bytes. `EncodeWith` returns a configurer for attaching hooks (see below).

### Decoder

{{< code lang="go" >}}
type Decoder interface {
    Decode(data []byte) (runtime.Value, error)
    DecodeWith() DecoderConfigurer
}
{{</ code >}}

`Decode` converts bytes back into a runtime value. `DecodeWith` returns a configurer for attaching hooks.

## Encoder and decoder hooks

Hooks let you intercept encoding and decoding without replacing the codec. The configurer chain builds a new encoder or decoder with hooks attached:

{{< code lang="go" >}}
encoder := codec.EncodeWith().
    PreHook(func(value runtime.Value) error {
        // runs before encoding
        return nil
    }).
    PostHook(func(value runtime.Value, err error) error {
        // runs after encoding
        return nil
    }).
    Encoder()
{{</ code >}}

### Hook types

| Hook | Signature | When it runs |
|------|-----------|-------------|
| `PreEncoderHook` | `func(value runtime.Value) error` | Before encoding a value |
| `PostEncoderHook` | `func(value runtime.Value, err error) error` | After encoding; receives the encode error |
| `PreDecoderHook` | `func(data []byte) error` | Before decoding bytes |
| `PostDecoderHook` | `func(data []byte, err error) error` | After decoding; receives the decode error |

Multiple hooks of the same type run in registration order. If any hook returns an error, processing stops.

## The Registry

The `encoding.Registry` stores codecs by MIME-normalized content type:

{{< code lang="go" >}}
// Create a registry seeded with codecs
registry := encoding.NewRegistry(jsonCodec, msgpackCodec)

// Or start empty and register individually
registry := encoding.NewEmptyRegistry()
registry.Register(myCodec)
{{</ code >}}

### Registry methods

| Method | Returns | Purpose |
|--------|---------|---------|
| `Register(codec)` | `error` | Store a codec by its content type |
| `Codec(contentType)` | `Codec, error` | Look up a full codec |
| `Encoder(contentType)` | `Encoder, error` | Look up an encoder |
| `Decoder(contentType)` | `Decoder, error` | Look up a decoder |
| `Clone()` | `*Registry` | Create an independent copy |

Content types are normalized using MIME media type parsing, so `application/json` and `application/json; charset=utf-8` resolve to the same codec.

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

A custom codec that encodes values as plain text:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/encoding"
    "github.com/MontFerret/ferret/v2/pkg/runtime"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

type plainTextCodec struct{}

func (plainTextCodec) ContentType() string {
    return "text/plain"
}

func (plainTextCodec) Encode(value runtime.Value) ([]byte, error) {
    return []byte(value.String()), nil
}

func (c plainTextCodec) EncodeWith() encoding.EncoderConfigurer {
    return &plainTextEncoderConfigurer{codec: c}
}

func (plainTextCodec) Decode(data []byte) (runtime.Value, error) {
    return runtime.NewString(string(data)), nil
}

func (c plainTextCodec) DecodeWith() encoding.DecoderConfigurer {
    return &plainTextDecoderConfigurer{codec: c}
}

type plainTextEncoderConfigurer struct {
    codec plainTextCodec
    pre   []encoding.PreEncoderHook
    post  []encoding.PostEncoderHook
}

func (e *plainTextEncoderConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
    if hook != nil {
        e.pre = append(e.pre, hook)
    }
    return e
}

func (e *plainTextEncoderConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
    if hook != nil {
        e.post = append(e.post, hook)
    }
    return e
}

func (e *plainTextEncoderConfigurer) Encoder() encoding.Encoder {
    return e.codec
}

type plainTextDecoderConfigurer struct {
    codec plainTextCodec
    pre   []encoding.PreDecoderHook
    post  []encoding.PostDecoderHook
}

func (d *plainTextDecoderConfigurer) PreHook(hook encoding.PreDecoderHook) encoding.DecoderConfigurer {
    if hook != nil {
        d.pre = append(d.pre, hook)
    }
    return d
}

func (d *plainTextDecoderConfigurer) PostHook(hook encoding.PostDecoderHook) encoding.DecoderConfigurer {
    if hook != nil {
        d.post = append(d.post, hook)
    }
    return d
}

func (d *plainTextDecoderConfigurer) Decoder() encoding.Decoder {
    return d.codec
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
        source.NewAnonymous(`RETURN "Hello, Ferret!"`),
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

    fmt.Println(output.ContentType)
    // text/plain
    fmt.Println(string(output.Content))
    // Hello, Ferret!
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-programs,embedding-modules" >}}
