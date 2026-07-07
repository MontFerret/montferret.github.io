---
title: "Pre-compile and distribute programs"
sidebarTitle: "Pre-compiled programs"
weight: 130
draft: false
description: "Compile FQL scripts to binary artifacts, store them, and load them at runtime without the compiler."
---

# Pre-compile and distribute programs

Ferret can compile FQL scripts into portable binary artifacts. Loading a pre-compiled artifact skips the compilation step entirely — useful for faster startup, distributing scripts without source, or caching build output in CI.

For the artifact format specification and full API, see [Programs]({{< ref "/docs/embedding/programs" >}}).

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
    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    plan, err := engine.Compile(
        context.Background(),
        source.New("extract-titles.fql", `
            LET page = WEB::HTML::OPEN(@url)
            RETURN page[~ css` + "`h1, h2`" + `][*].textContent
        `),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    data, err := plan.Marshal()
    if err != nil {
        log.Fatal(err)
    }

    if err := os.WriteFile("extract-titles.fbc", data, 0644); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("compiled %d bytes\n", len(data))
}
{{</ code >}}

`plan.Marshal()` serializes the bytecode to a self-describing binary artifact with the `.fbc` extension (Ferret Bytecode).

## Compile from the CLI

The CLI can compile scripts without writing any Go:

{{< terminal command="true" >}}
ferret build extract-titles.fql -o extract-titles.fbc
{{</ terminal >}}

Compile an entire directory:

{{< terminal command="true" >}}
for f in scripts/*.fql; do
    ferret build "$f" -o "${f%.fql}.fbc"
done
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
import "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"

data, err := plan.Marshal(ferret.WithProgramFormat(artifact.FormatJSON))
{{</ code >}}

JSON artifacts are human-readable — useful for inspecting the compiled output or debugging issues:

{{< terminal command="true" >}}
ferret build script.fql -o script.fbc --format json
cat script.fbc | python3 -c "import sys; data=sys.stdin.buffer.read(); print(data[14:].decode())" | jq .
{{</ terminal >}}

The first 14 bytes are the binary header; the rest is the JSON payload.

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

//go:embed scripts/*.fbc
var scripts embed.FS

func main() {
    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    data, err := scripts.ReadFile("scripts/extract-titles.fbc")
    if err != nil {
        log.Fatal(err)
    }

    plan, err := engine.Load(data)
    if err != nil {
        log.Fatal(err)
    }
    defer plan.Close()

    session, err := plan.NewSession(context.Background(),
        ferret.WithSessionParam("url", "https://mockery.ferretlang.org"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    output, err := session.Run(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(output.Content))
}
{{</ code >}}

Build workflow:

{{< terminal command="true" >}}
ferret build scripts/extract-titles.fql -o scripts/extract-titles.fbc
go build -o myservice .
{{</ terminal >}}

The resulting binary contains the compiled FQL — deploy it without any `.fql` or `.fbc` files.

## Load and run at startup

Use `engine.Load` to create a plan from artifact bytes:

{{< code lang="go" >}}
data, err := os.ReadFile("extract-titles.fbc")
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

When your application accepts both `.fql` source and `.fbc` artifacts, use `artifact.HasMagic` to choose the right path:

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
if err != nil {
    if errors.Is(err, artifact.ErrIncompatibleISA) {
        // Artifact was compiled with a different bytecode version.
        // Recompile from source.
        log.Println("artifact ISA mismatch, recompiling from source")
        plan, err = engine.Compile(ctx, source.New(name, sourceText))
    } else if errors.Is(err, artifact.ErrUnsupportedSchema) {
        log.Fatal("artifact schema version not supported by this runtime")
    } else if errors.Is(err, artifact.ErrUnknownFormat) {
        log.Fatal("artifact payload format not recognized")
    } else {
        log.Fatal(err)
    }
}
{{</ code >}}

Strategy: always keep the `.fql` source alongside `.fbc` artifacts so you can recompile when the ISA changes.

## Build a compile-and-cache workflow

Compile on first use and cache the artifact for subsequent runs:

{{< code lang="go" >}}
func loadOrCompile(engine *ferret.Engine, fqlPath string) (*ferret.Plan, error) {
    ctx := context.Background()
    fbcPath := strings.TrimSuffix(fqlPath, ".fql") + ".fbc"

    fqlInfo, err := os.Stat(fqlPath)
    if err != nil {
        return nil, err
    }

    // Try loading cached artifact
    if fbcInfo, err := os.Stat(fbcPath); err == nil {
        if fbcInfo.ModTime().After(fqlInfo.ModTime()) {
            data, err := os.ReadFile(fbcPath)
            if err == nil {
                plan, err := engine.Load(data)
                if err == nil {
                    return plan, nil
                }
                // Fall through to recompile on load error
            }
        }
    }

    // Compile from source
    src, err := os.ReadFile(fqlPath)
    if err != nil {
        return nil, err
    }

    plan, err := engine.Compile(ctx, source.New(fqlPath, string(src)))
    if err != nil {
        return nil, err
    }

    // Cache the artifact
    data, err := plan.Marshal()
    if err == nil {
        os.WriteFile(fbcPath, data, 0644)
    }

    return plan, nil
}
{{</ code >}}

The function checks whether the cached `.fbc` file is newer than the `.fql` source. If it is, it loads the artifact directly. Otherwise, it compiles from source and caches the result.

## Integrate into CI

Compile artifacts in CI and deploy only the binaries:

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
        run: go install github.com/MontFerret/ferret/v2/cmd/ferret@latest

      - name: Compile FQL scripts
        run: |
          for f in scripts/*.fql; do
            ferret build "$f" -o "${f%.fql}.fbc"
          done

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: fbc-artifacts
          path: scripts/*.fbc
```

A downstream deploy job downloads the artifacts and embeds or ships them with the application.

## Complete example

A two-phase program that compiles a directory of `.fql` files and then loads and runs them:

{{< code lang="go" >}}
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"

    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
    "github.com/MontFerret/ferret/v2/pkg/source"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Fprintf(os.Stderr, "usage: %s <build|run> <dir>\n", os.Args[0])
        os.Exit(1)
    }

    cmd, dir := os.Args[1], os.Args[2]

    engine, err := ferret.New()
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    ctx := context.Background()

    switch cmd {
    case "build":
        files, _ := filepath.Glob(filepath.Join(dir, "*.fql"))
        for _, f := range files {
            src, err := os.ReadFile(f)
            if err != nil {
                log.Fatal(err)
            }

            plan, err := engine.Compile(ctx, source.New(filepath.Base(f), string(src)))
            if err != nil {
                log.Fatalf("compile %s: %v", f, err)
            }

            data, err := plan.Marshal()
            plan.Close()
            if err != nil {
                log.Fatal(err)
            }

            out := strings.TrimSuffix(f, ".fql") + ".fbc"
            if err := os.WriteFile(out, data, 0644); err != nil {
                log.Fatal(err)
            }

            fmt.Printf("compiled %s → %s (%d bytes)\n", f, out, len(data))
        }

    case "run":
        files, _ := filepath.Glob(filepath.Join(dir, "*.fbc"))
        for _, f := range files {
            data, err := os.ReadFile(f)
            if err != nil {
                log.Fatal(err)
            }

            if !artifact.HasMagic(data) {
                log.Fatalf("%s is not a valid artifact", f)
            }

            plan, err := engine.Load(data)
            if err != nil {
                log.Fatalf("load %s: %v", f, err)
            }

            session, err := plan.NewSession(ctx)
            if err != nil {
                plan.Close()
                log.Fatal(err)
            }

            output, err := session.Run(ctx)
            session.Close()
            plan.Close()

            if err != nil {
                log.Fatalf("run %s: %v", f, err)
            }

            fmt.Printf("--- %s ---\n%s\n", filepath.Base(f), output.Content)
        }

    default:
        log.Fatalf("unknown command: %s", cmd)
    }
}
{{</ code >}}

Usage:

{{< terminal command="true" >}}
go run . build scripts/
go run . run scripts/
{{</ terminal >}}

## Next steps

{{< docs-related tiles="embedding-programs,embedding-overview,tools-cli-build,guide-embedding-runtime" >}}
