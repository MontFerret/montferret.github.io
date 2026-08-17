---
title: "Install a module"
sidebarTitle: "Install"
weight: 10
draft: false
description: "Install a registered Ferret module into a Go application."
---

# Install a module

Use `ferret mod install` to add a registered module to an existing Go application. The installer resolves a compatible release, adds its Go package, and registers it in the application's Ferret engine composition.

This workflow requires Ferret CLI `v{{< data "versions.cli.v2" >}}` or later and a Go project with an existing `go.mod`.

## Find a module

Search the public Registry by module identity or description:

{{< terminal command="true" >}}
ferret mod search archive
{{< /terminal >}}

Inspect a module before installing it:

{{< terminal command="true" >}}
ferret mod info montferret/archive
{{< /terminal >}}

You can also browse module versions and documentation in the [Ferret Registry]({{< ref "/registry" >}}).

## Install a compatible release

Run the command from the Go project that will host Ferret:

{{< terminal command="true" >}}
ferret mod install montferret/archive
{{< /terminal >}}

Without a version, the installer selects the newest registered release compatible with the Ferret version already selected by the project.

To request an exact release, append its strict semantic version without a leading `v`:

{{< terminal command="true" >}}
ferret mod install montferret/archive@1.0.0-rc.3
{{< /terminal >}}

An exact release still has to support the project's Ferret version. The command fails instead of changing the project to a different Ferret release.

## Approve missing project setup

The application must have:

- `github.com/MontFerret/ferret/v2` in its Go module graph
- one unambiguous `ferret.New(...)` composition where modules can be registered

When either prerequisite is missing, an interactive install shows every proposed setup change and asks for approval. For an empty project or a project with one Go package, it can add the Ferret dependency and create an exported composition helper in `ferret.go`:

```go
func NewFerret(options ...ferret.Option) (*ferret.Engine, error) {
    return ferret.New(options...)
}
```

Projects with multiple possible packages must add their composition manually so the CLI does not guess which package owns the engine.

In CI or another non-interactive environment, approve the same safe prerequisites explicitly:

{{< terminal command="true" >}}
ferret mod install --yes montferret/archive
{{< /terminal >}}

Without `--yes`, a non-interactive command exits with the equivalent manual setup steps instead of reading from standard input.

## What the installer changes

The installer:

1. resolves the Registry release and its published Go package path
2. stages changes to `go.mod`, `go.sum`, and the owning composition file
3. adds a normal Go import and `ferret.WithModules(module.New())`
4. asks the Go toolchain for the exact published package version
5. builds the owning package with the staged changes
6. applies all validated files together

It rejects replacements for the selected Registry package, a dependency resolution that changes the project's Ferret version, and a release whose origin conflicts with the Registry commit when Go reports that origin. A failed validation leaves the original project files in place.

Run the same command again to confirm the exact dependency and registration are already present; the installer reports the module as installed without adding a duplicate.

## Configure the module

Installation uses the module package's zero-argument `New()` constructor. If the module exposes functional options, edit the generated registration after installation:

```go
ferret.WithModules(
    archive.New(
        archive.WithMaxEntrySize(32 << 20),
    ),
)
```

Use options documented by the module. This example limits the materialized size of each archive entry to 32 MiB.

## Next steps

{{< docs-related tiles="embedding-go-modules,runtime-modules-develop,registry" >}}
