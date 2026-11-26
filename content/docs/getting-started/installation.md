---
title: "Installation"
weight: 30
draft: false
description: "Set up Visual Studio Code, a supported JetBrains IDE, the Ferret CLI, or an embedding API."
aliases:
    - /docs/installation/
relatedTileOverrides:
  getting-started-quick-start:
    kicker: "Start here"
    title: "Quick Start"
    description: "Run your first Ferret script and learn the basic workflow."
---

# Installation

Choose how you will work with Ferret. Use an IDE to write, run, and debug FQL, the CLI for terminal and shell workflows, or an embedding API to execute FQL from a Go or JavaScript application.

> **Alpha status**
>
> Ferret v2 is currently in alpha. The language, runtime, CLI, modules, and embedding APIs may change before beta. Pin versions for scripts, CI, and application integrations.

> **Looking for Ferret v1?**
>
> Ferret v1 remains available for existing projects, but new users should start with Ferret v2. See the [v1 → v2 migration overview]({{< ref "/docs/migrations/v1-to-v2" >}}) for automated changes and manual follow-up.

## Use your IDE

- [Visual Studio Code installation]({{< ref "/docs/tools/visual-studio-code#installation" >}}) for the official Ferret Lang extension.
- [JetBrains IDEs installation]({{< ref "/docs/tools/jetbrains#installation" >}}) to install the Ferret plugin from JetBrains Marketplace in IntelliJ IDEA or GoLand.

Both integrations include the Ferret daemon. After installing, save a `.fql` file and use your IDE's Run or Debug action.

## Use the CLI

Choose the CLI to run FQL files and expressions from a terminal, shell script, or CI job.

{{< terminal >}}
ferret version
{{< /terminal >}}

See [CLI Installation]({{< ref "/docs/tools/cli/installation" >}}) for prebuilt binaries, source installation, updates, and verification.

## Embed Ferret in Go

Choose native Go embedding when the host needs full control over runtime configuration, modules, host values, codecs, and sandboxed services.

Install the published module from [Go Embedding: Getting Started]({{< ref "/docs/embedding/go/getting-started#installation" >}}). The guide runs a query and shows how to consume and close its output.

## Embed Ferret in JavaScript

Choose JavaScript embedding to run Ferret from Node.js or a modern browser through `@montferret/ferret`.

{{< terminal >}}
npm install --save-exact @montferret/ferret@{{< data "versions.core_js.v2" >}}
{{< /terminal >}}

Continue with [JavaScript Embedding: Getting Started]({{< ref "/docs/embedding/javascript/getting-started" >}}).

## What to choose next

- Start with the [Quick Start]({{< ref "/docs/getting-started/quick-start" >}}) to try FQL in your IDE, the terminal, or the online editor.
- Read the [Embedding overview]({{< ref "/docs/embedding" >}}) when an application will own the runtime.
- Use [Worker]({{< ref "/docs/tools/worker" >}}) when an existing system needs Ferret through a separately deployed HTTP service.

{{< docs-related tiles="getting-started-quick-start,tools-visual-studio-code,tools-jetbrains,tools-cli-installation,embedding-go,embedding-javascript,tools-worker" >}}
