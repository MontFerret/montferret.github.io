---
title: "Pre-compile and distribute programs"
sidebarTitle: "Pre-compiled programs"
weight: 130
draft: false
description: "Compile FQL scripts to binary artifacts, store them, and load them without compiling the source again."
---

# Pre-compile and distribute programs

Ferret can compile FQL scripts into binary artifacts. Loading a pre-compiled artifact skips parsing and source compilation — useful for faster startup, distributing scripts without separate `.fql` files, or caching build output. The loading engine still owns a compiler and must register every module or host function used by the program.

For the artifact format specification and full API, see [Programs]({{< ref "/docs/embedding/go/programs" >}}).

## Compile from Go

Compile a query and serialize the resulting plan:

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
    engine, err := ferret.New()
    if err != nil {
        return err
    }
    defer engine.Close()

    ctx := context.Background()
    plan, err := engine.Compile(
        ctx,
        source.New("greeting.fql", `return upper(@name)`),
    )
    if err != nil {
        return err
    }
    defer plan.Close()

    data, err := plan.Marshal()
    if err != nil {
        return err
    }

    if err := os.WriteFile("greeting.fqlc", data, 0o644); err != nil {
        return err
    }

    fmt.Printf("compiled %d bytes\n", len(data))

    return nil
}
{{</ code >}}

`plan.Marshal()` serializes the bytecode to a self-describing binary artifact. These examples use the CLI's `.fqlc` extension.

## Compile from the CLI

The CLI can compile scripts without writing any Go:

{{< terminal command="true" >}}
ferret build greeting.fql -o greeting.fqlc
{{</ terminal >}}

Compile an entire directory:

{{< terminal command="true" >}}
ferret build scripts/*.fql -o dist/
{{</ terminal >}}

See [CLI Build]({{< ref "/docs/tools/cli/build" >}}) for all options.

## Choose a payload format

Artifacts support two payload formats:

| Format | Constant | Size | Use case |
| --- | --- | --- | --- |
| MessagePack | `artifact.FormatMsgPack` | Smaller | Production (default) |
| JSON | `artifact.FormatJSON` | Larger | Debugging, inspection |

To use JSON:

{{< code lang="go" >}}
import (
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

data, err := plan.Marshal(ferret.WithProgramFormat(artifact.FormatJSON))
{{</ code >}}

`ferret build` writes the default MessagePack format; it does not expose a payload-format flag. Use the Go API above when you need JSON. In either format, the first 14 bytes are the binary header and the remaining bytes are the encoded payload.

## Embed artifacts in a Go binary

Go's `embed` directive lets you compile FQL artifacts into your Go binary at build time — no external files needed at runtime:

{{< code lang="go" >}}
package main

import (
    "context"
    "embed"
    "fmt"
    "log"

    "github.com/MontFerret/ferret/v2"
)

//go:embed scripts/*.fqlc
var scripts embed.FS

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    engine, err := ferret.New()
    if err != nil {
        return err
    }
    defer engine.Close()

    data, err := scripts.ReadFile("scripts/greeting.fqlc")
    if err != nil {
        return err
    }

    plan, err := engine.Load(data)
    if err != nil {
        return err
    }
    defer plan.Close()

    ctx := context.Background()
    session, err := plan.NewSession(ctx,
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

`output.Content` is encoded output. The default JSON codec includes quotes around the returned FQL string.

Build workflow:

{{< terminal command="true" >}}
ferret build scripts/greeting.fql -o scripts/greeting.fqlc
go build -o myservice .
{{</ terminal >}}

The resulting binary contains the compiled FQL — deploy it without any external `.fql` or `.fqlc` files.

## Load and run at startup

Use `engine.Load` to create a plan from artifact bytes:

{{< code lang="go" >}}
data, err := os.ReadFile("greeting.fqlc")
if err != nil {
    log.Fatal(err)
}

plan, err := engine.Load(data)
if err != nil {
    log.Fatal(err)
}
defer plan.Close()
{{</ code >}}

### Auto-detect source vs. artifact

When your application accepts both `.fql` source and `.fqlc` artifacts, use `artifact.HasMagic` to choose the right path:

{{< code lang="go" >}}
import "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"

data, err := os.ReadFile(path)
if err != nil {
    log.Fatal(err)
}

var plan *ferret.Plan

if artifact.HasMagic(data) {
    plan, err = engine.Load(data)
} else {
    plan, err = engine.Compile(ctx, source.New(path, string(data)))
}
if err != nil {
    log.Fatal(err)
}
defer plan.Close()
{{</ code >}}

`HasMagic` checks the first 4 bytes for the `FBC2` magic number. It does not validate the full artifact — `Load` handles that.

## Handle version mismatches

Artifacts encode the bytecode instruction set version (ISA). When the runtime's ISA changes between Ferret releases, old artifacts become incompatible:

{{< code lang="go" >}}
import (
    "errors"

    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
)

plan, err := engine.Load(data)
if errors.Is(err, artifact.ErrIncompatibleISA) {
    log.Println("artifact ISA mismatch, recompiling from source")
    plan, err = engine.Compile(ctx, source.New(name, sourceText))
}

if err != nil {
    switch {
    case errors.Is(err, artifact.ErrUnsupportedSchema):
        log.Fatal("artifact schema version not supported by this runtime")
    case errors.Is(err, artifact.ErrUnknownFormat):
        log.Fatal("artifact payload format not recognized")
    default:
        log.Fatal(err)
    }
}
defer plan.Close()
{{</ code >}}

Strategy: always keep the `.fql` source alongside `.fqlc` artifacts so you can recompile when the ISA changes.

## Build a compile-and-cache workflow

Compile on first use and cache the artifact for subsequent runs:

{{< code lang="go" >}}
func loadOrCompile(ctx context.Context, engine *ferret.Engine, fqlPath string) (*ferret.Plan, error) {
    artifactPath := strings.TrimSuffix(fqlPath, ".fql") + ".fqlc"

    fqlInfo, err := os.Stat(fqlPath)
    if err != nil {
        return nil, err
    }

    // Try the cache first. A read or load failure falls through to compilation.
    if artifactInfo, statErr := os.Stat(artifactPath); statErr == nil {
        if artifactInfo.ModTime().After(fqlInfo.ModTime()) {
            data, readErr := os.ReadFile(artifactPath)
            if readErr == nil {
                if plan, loadErr := engine.Load(data); loadErr == nil {
                    return plan, nil
                }
            }
        }
    } else if !errors.Is(statErr, os.ErrNotExist) {
        return nil, statErr
    }

    src, err := os.ReadFile(fqlPath)
    if err != nil {
        return nil, err
    }

    plan, err := engine.Compile(ctx, source.New(fqlPath, string(src)))
    if err != nil {
        return nil, err
    }

    data, err := plan.Marshal()
    if err != nil {
        return nil, errors.Join(err, plan.Close())
    }
    if err := os.WriteFile(artifactPath, data, 0o644); err != nil {
        return nil, errors.Join(err, plan.Close())
    }

    return plan, nil
}
{{</ code >}}

The function checks whether the cached `.fqlc` file is newer than the `.fql` source. If it is, it loads the artifact directly. Otherwise, it compiles from source and caches the result. The caller owns the returned plan and must close it.

## Integrate into CI

Compile artifacts in CI and deploy only the generated artifacts:

```yaml
# .github/workflows/build.yml
name: Build FQL artifacts

on:
  push:
    paths:
      - 'scripts/**/*.fql'

jobs:
  compile:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Ferret CLI
        run: go install github.com/MontFerret/cli/v2/ferret@v{{< data "versions.cli.v2" >}}

      - name: Compile FQL scripts
        run: |
          for f in scripts/*.fql; do
            ferret build "$f" -o "${f%.fql}.fqlc"
          done

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: fqlc-artifacts
          path: scripts/*.fqlc
```

A downstream deploy job downloads the artifacts and embeds or ships them with the application.

## Complete example

A two-phase program that compiles a directory of `.fql` files and then loads and runs them:

{{< code lang="go" >}}
package main

import (
    "context"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    if err := run(context.Background(), os.Args); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func run(ctx context.Context, args []string) (err error) {
    if len(args) < 3 {
        return fmt.Errorf("usage: %s <build|run> <dir>", args[0])
    }

    cmd, dir := args[1], args[2]

    engine, err := ferret.New()
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, engine.Close())
    }()

    switch cmd {
    case "build":
        return forEachFile(dir, "*.fql", func(path string) error {
            return compileFile(ctx, engine, path)
        })
    case "run":
        return forEachFile(dir, "*.fqlc", func(path string) error {
            return runArtifact(ctx, engine, path)
        })
    default:
        return fmt.Errorf("unknown command: %s", cmd)
    }
}

func forEachFile(dir, pattern string, visit func(string) error) error {
    files, err := filepath.Glob(filepath.Join(dir, pattern))
    if err != nil {
        return err
    }
    for _, path := range files {
        if err := visit(path); err != nil {
            return err
        }
    }
    return nil
}

func compileFile(ctx context.Context, engine *ferret.Engine, path string) (err error) {
    src, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    plan, err := engine.Compile(ctx, source.New(filepath.Base(path), string(src)))
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, plan.Close())
    }()

    data, err := plan.Marshal()
    if err != nil {
        return err
    }

    artifactPath := strings.TrimSuffix(path, ".fql") + ".fqlc"
    if err := os.WriteFile(artifactPath, data, 0o644); err != nil {
        return err
    }

    fmt.Printf("compiled %s → %s (%d bytes)\n", path, artifactPath, len(data))
    return nil
}

func runArtifact(ctx context.Context, engine *ferret.Engine, path string) (err error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    if !artifact.HasMagic(data) {
        return fmt.Errorf("%s is not a Ferret artifact", path)
    }

    plan, err := engine.Load(data)
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, plan.Close())
    }()

    session, err := plan.NewSession(ctx)
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, session.Close())
    }()

    output, err := session.Run(ctx)
    if err != nil {
        return err
    }

    fmt.Printf("--- %s ---\n%s\n", filepath.Base(path), output.Content)
    return nil
}
{{</ code >}}

Usage:

{{< terminal command="true" >}}
go run . build scripts/
go run . run scripts/
{{</ terminal >}}

## Next steps

{{< docs-related tiles="embedding-go-programs,embedding-go-executing,tools-cli-build,embedding-go-application-guide" >}}
