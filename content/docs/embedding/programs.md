---
title: "Programs"
sidebarTitle: "Programs"
weight: 56
draft: false
description: "Compile queries to portable binary artifacts, load pre-compiled programs, and work with the artifact format."
---

# Programs

A compiled Ferret query can be serialized into a binary artifact and loaded later without the compiler. This is useful for distributing pre-compiled queries, skipping compilation at runtime, or caching build output on disk.

## Two paths to a Plan

There are two ways to get a `Plan` from the engine:

{{< code lang="go" >}}
// From FQL source — compiles and returns a plan
plan, err := engine.Compile(ctx, source.NewAnonymous(`RETURN 1 + 1`))

// From a pre-compiled artifact — loads and returns a plan
plan, err := engine.Load(artifactBytes)
{{</ code >}}

Both produce the same `*Plan` that you create sessions from and run. The difference is where the bytecode comes from: the compiler or a serialized artifact.

## Serializing a program

The `artifact` package provides `Marshal` to serialize a compiled program:

{{< code lang="go" >}}
import (
    "os"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/artifact"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

// Compile to bytecode
plan, err := engine.Compile(ctx, source.New("query.fql", `RETURN UPPER(@name)`))
if err != nil {
    log.Fatal(err)
}

// Serialize to bytes
data, err := plan.Marshal()
if err != nil {
    log.Fatal(err)
}

// Write to file
err = os.WriteFile("query.fbc", data, 0644)
if err != nil {
    log.Fatal(err)
}
{{</ code >}}

The default payload format is MessagePack. To use JSON instead, set the format in the options:

{{< code lang="go" >}}
data, err := plan.Marshal(ferret.WithProgramFormat(artifact.FormatJSON))
{{</ code >}}

## Loading a program

Load the artifact back into the engine with `engine.Load`:

{{< code lang="go" >}}
data, err := os.ReadFile("query.fbc")
if err != nil {
    log.Fatal(err)
}

plan, err := engine.Load(data)
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

session, err := plan.NewSession(ctx,
    ferret.WithSessionParam("name", "ferret"),
)
if err != nil {
    log.Fatal(err)
}
defer session.Close()

output, err := session.Run(ctx)
{{</ code >}}

For lower-level access, `ferret.UnmarshalProgram` returns a `*bytecode.Program` without wrapping it in a plan:

{{< code lang="go" >}}
program, err := ferret.UnmarshalProgram(data)
{{</ code >}}

## Artifact format

An artifact is a self-describing binary format with a 14-byte header followed by the encoded payload:

| Bytes | Field | Type | Description |
|-------|-------|------|-------------|
| 0–3 | Magic | `[4]byte` | `FBC2` — identifies the file as a Ferret artifact |
| 4 | Format | `uint8` | Payload format ID (1 = JSON, 2 = MessagePack) |
| 5 | Schema | `uint8` | Schema version (currently 1) |
| 6–7 | ISA | `uint16 LE` | Bytecode instruction set version |
| 8–9 | Flags | `uint16 LE` | Reserved (must be 0 in schema v1) |
| 10–13 | Length | `uint32 LE` | Payload length in bytes |

All multi-byte fields are little-endian. The header is followed by exactly `Length` bytes of payload data in the format specified by the Format field.

## Payload formats

| ID | Constant | Format | Use case |
|----|----------|--------|----------|
| 1 | `artifact.FormatJSON` | JSON | Human-readable, debugging |
| 2 | `artifact.FormatMsgPack` | MessagePack | Compact, production (default) |

Both formats implement the `format.Format` interface:

{{< code lang="go" >}}
type Format interface {
    Name() string
    Marshal(*bytecode.Program) ([]byte, error)
    Unmarshal([]byte) (*bytecode.Program, error)
}
{{</ code >}}

## Detecting artifacts

Use `HasMagic` to quickly check whether a byte slice looks like a Ferret artifact:

{{< code lang="go" >}}
if artifact.HasMagic(data) {
    plan, err := engine.Load(data)
    // ...
} else {
    plan, err := engine.Compile(ctx, source.NewAnonymous(string(data)))
    // ...
}
{{</ code >}}

`HasMagic` only checks the first 4 bytes. It does not validate the full artifact — use `Load` or `Unmarshal` for that.

## Custom loaders

The engine uses a `Loader` to decode artifacts. By default it supports JSON and MessagePack payloads. You can create a custom loader with additional formats:

{{< code lang="go" >}}
loader := artifact.NewLoader(
    artifact.RegisteredFormat{ID: 1, Format: formatjson.Default},
    artifact.RegisteredFormat{ID: 2, Format: formatmsgpack.Default},
    artifact.RegisteredFormat{ID: 3, Format: myCustomFormat},
)

engine, err := ferret.New(
    ferret.WithProgramLoader(loader),
)
{{</ code >}}

The format ID in the `RegisteredFormat` must match the format ID written in the artifact header. When loading, the loader reads the header, selects the registered format by ID, and delegates decoding to it.

## Error handling

The `artifact` package defines sentinel errors for each validation failure:

| Error | Cause |
|-------|-------|
| `ErrInvalidMagic` | First 4 bytes are not `FBC2` |
| `ErrUnsupportedSchema` | Schema version is not supported by this loader |
| `ErrIncompatibleISA` | Bytecode ISA version does not match the runtime |
| `ErrUnknownFormat` | Payload format ID is not registered |
| `ErrInvalidHeader` | Header is malformed or payload length does not match |
| `ErrInvalidPayload` | Payload format decoder failed |
| `ErrInvalidArtifact` | General loader or artifact state error |

Use `errors.Is` to check for specific failures:

{{< code lang="go" >}}
_, err := engine.Load(data)
if errors.Is(err, artifact.ErrIncompatibleISA) {
    // artifact was compiled with a different bytecode version
}
{{</ code >}}

## Complete example

Compile a query, save it to disk, then load and run it in a separate step:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
    "github.com/MontFerret/ferret/v2/pkg/compiler"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    ctx := context.Background()

    // --- Build phase ---

    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    plan, err := engine.Compile(source.New("greeting.fql", `RETURN UPPER(@name)`))
    if err != nil {
        log.Fatal(err)
    }

    data, err := plan.Marshal()
    if err != nil {
        log.Fatal(err)
    }

    if err := os.WriteFile("greeting.fbc", data, 0644); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("compiled %d bytes\n", len(data))

    // --- Load phase ---

    saved, err := os.ReadFile("greeting.fbc")
    if err != nil {
        log.Fatal(err)
    }

    plan, err := engine.Load(saved)
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    session, err := plan.NewSession(ctx,
        ferret.WithSessionParam("name", "ferret"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    output, err := session.Run(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
    // "FERRET"
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-value-encoders,embedding-configuration" >}}
