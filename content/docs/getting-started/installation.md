---
title: "Installation"
weight: 30
draft: false
description: "Choose the Ferret CLI or embed Ferret in a Go or JavaScript application."
aliases:
    - /docs/installation/
relatedTileOverrides:
  getting-started-quick-start:
    kicker: "Start here"
    title: "Quick Start"
    description: "Run your first Ferret script and learn the basic workflow."
---

# Installation

Choose how your application will run Ferret. Use the CLI for local scripts and shell workflows, or embed the runtime when a Go or JavaScript application needs to execute FQL directly.

> **Alpha status**
>
> Ferret v2 is currently in alpha. The language, runtime, CLI, modules, and embedding APIs may change before beta. Pin versions for scripts, CI, and application integrations.

> **Looking for Ferret v1?**
>
> Ferret v1 remains available for existing projects, but new users should start with Ferret v2. See the [migration guide]({{< ref "/docs/tools/cli/migrate" >}}) for the supported mechanical migration steps.

## Use the CLI

Choose the CLI to run FQL files and expressions from a terminal, shell script, or CI job.

{{< terminal >}}
ferret version
{{< /terminal >}}

See [Install the CLI]({{< ref "/docs/tools/cli/install" >}}) for prebuilt binaries, source installation, updates, and verification.

## Embed Ferret in Go

Choose native Go embedding when the host needs full control over runtime configuration, modules, host values, codecs, and sandboxed services.

{{< terminal >}}
go get github.com/MontFerret/ferret/v2@v{{< data "versions.runtime.v2" >}}
{{< /terminal >}}

Continue with [Go Embedding: Getting Started]({{< ref "/docs/embedding/go/getting-started" >}}).

## Embed Ferret in JavaScript

Choose JavaScript embedding to run Ferret from Node.js or a modern browser through `@montferret/ferret`.

{{< terminal >}}
npm install @montferret/ferret
{{< /terminal >}}

Continue with [JavaScript Embedding: Getting Started]({{< ref "/docs/embedding/javascript/getting-started" >}}).

## What to choose next

- Start with the [Quick Start]({{< ref "/docs/getting-started/quick-start" >}}) when you want to learn FQL from the terminal or playground.
- Read the [Embedding overview]({{< ref "/docs/embedding" >}}) when an application will own the runtime.
- Use [Worker]({{< ref "/docs/tools/worker" >}}) when an existing system needs Ferret through a separately deployed HTTP service.

{{< docs-related tiles="getting-started-quick-start,tools-cli-install,embedding-go,embedding-javascript,tools-worker" >}}
