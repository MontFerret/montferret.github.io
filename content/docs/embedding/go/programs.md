---
title: "Programs"
sidebarTitle: "Programs"
weight: 90
draft: false
description: "Compile queries to binary artifacts, load pre-compiled programs, and work with the artifact format."
aliases:
    - /docs/embedding/programs/
---

# Programs

A compiled Ferret query can be serialized into a binary artifact and loaded later without compiling its source again. This is useful for distributing pre-compiled queries, skipping source compilation at runtime, or caching build output on disk. `Engine.Load` bypasses compilation, but the engine itself is still constructed with its compiler.

Artifacts contain bytecode and host-function signatures, not host-function or module implementations. The engine that loads an artifact must still be configured with the modules and host functions the program uses.

## Two paths to a Plan

There are two ways to get a `Plan` from the engine:

{{< code lang="go" >}}
// From FQL source — compiles and returns a plan
sourcePlan, err := engine.Compile(ctx, source.NewAnonymous(`return 1 + 1`))
if err != nil {
    log.Fatal(err)
}
defer sourcePlan.Close()

// From a pre-compiled artifact — loads and returns a plan
artifactPlan, err := engine.Load(artifactBytes)
if err != nil {
    log.Fatal(err)
}
defer artifactPlan.Close()
{{</ code >}}

Both produce the same `*Plan` that you create sessions from and run. The difference is where the bytecode comes from: the compiler or a serialized artifact.

## Serializing a program

`Plan.Marshal` uses the `artifact` package to serialize a compiled program:

{{< code lang="go" >}}
import (
    "log"
    "os"

    "github.com/MontFerret/ferret/v2/pkg/source"
)

// Compile to bytecode
plan, err := engine.Compile(ctx, source.New("query.fql", `return upper(@name)`))
if err != nil {
    log.Fatal(err)
}
defer plan.Close()

// Serialize to bytes
data, err := plan.Marshal()
if err != nil {
    log.Fatal(err)
}

// Write to file
err = os.WriteFile("query.fqlc", data, 0o644)
if err != nil {
    log.Fatal(err)
}
{{</ code >}}

The default payload format is MessagePack. To use JSON instead, set the format in the options:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

data, err := plan.Marshal(ferret.WithProgramFormat(artifact.FormatJSON))
{{</ code >}}

## Loading a program

Load the artifact back into the engine with `engine.Load`:

{{< code lang="go" >}}
data, err := os.ReadFile("query.fqlc")
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
if err != nil {
    log.Fatal(err)
}
{{</ code >}}

`Session.Run` returns encoded output. With the default JSON codec, a returned FQL string is stored in `output.Content` as a quoted JSON string.

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
var plan *ferret.Plan

if artifact.HasMagic(data) {
    plan, err = engine.Load(data)
} else {
    plan, err = engine.Compile(ctx, source.NewAnonymous(string(data)))
}
if err != nil {
    log.Fatal(err)
}
defer plan.Close()
{{</ code >}}

`HasMagic` only checks the first 4 bytes. It does not validate the full artifact — use `Load` or `Unmarshal` for that.

## Custom loaders

The engine uses a `Loader` to decode artifacts. By default it supports JSON and MessagePack payloads. You can create a custom loader with additional formats:

{{< code lang="go" >}}
loader := artifact.NewLoader(
    artifact.RegisteredFormat{ID: artifact.FormatJSON, Format: formatjson.Default},
    artifact.RegisteredFormat{ID: artifact.FormatMsgPack, Format: formatmsgpack.Default},
    artifact.RegisteredFormat{ID: 3, Format: myCustomFormat},
)

engine, err := ferret.New(
    ferret.WithProgramLoader(loader),
)
{{</ code >}}

The format ID in the `RegisteredFormat` must match the format ID written in the artifact header. When loading, the loader reads the header, selects the registered format by ID, and delegates decoding to it. A custom loader extends `Engine.Load`; `Plan.Marshal` still writes only the built-in JSON and MessagePack formats.

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
plan, err := engine.Load(data)
if errors.Is(err, artifact.ErrIncompatibleISA) {
    // artifact was compiled with a different bytecode version
}
if err != nil {
    log.Fatal(err)
}
defer plan.Close()
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
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    ctx := context.Background()

    // --- Build phase ---

    engine, err := ferret.New()
    if err != nil {
        return err
    }
    defer engine.Close()

    compiledPlan, err := engine.Compile(ctx, source.New("greeting.fql", `return upper(@name)`))
    if err != nil {
        return err
    }
    defer compiledPlan.Close()

    data, err := compiledPlan.Marshal()
    if err != nil {
        return err
    }

    if err := os.WriteFile("greeting.fqlc", data, 0o644); err != nil {
        return err
    }

    fmt.Printf("compiled %d bytes\n", len(data))

    // --- Load phase ---

    saved, err := os.ReadFile("greeting.fqlc")
    if err != nil {
        return err
    }

    loadedPlan, err := engine.Load(saved)
    if err != nil {
        return err
    }
    defer loadedPlan.Close()

    session, err := loadedPlan.NewSession(ctx,
        ferret.WithSessionParam("name", "ferret"),
    )
    if err != nil {
        return err
    }
    defer session.Close()

    output, err := session.Run(ctx)
    if err != nil {
        return err
    }

    fmt.Println(string(output.Content))
    // "FERRET"

    return nil
}
{{</ code >}}

## Next steps

{{< docs-related tiles="embedding-go-custom-functions,embedding-go-parameters,embedding-go-modules,embedding-go-configuration" >}}
