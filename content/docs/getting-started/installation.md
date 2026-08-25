---
title: "Installation"
weight: 30
draft: false
description: "Choose Visual Studio Code, the Ferret CLI, or an embedding API."
aliases:
    - /docs/installation/
relatedTileOverrides:
  getting-started-quick-start:
    kicker: "Start here"
    title: "Quick Start"
    description: "Run your first Ferret script and learn the basic workflow."
---

# Installation

Choose how you will work with Ferret. Use Visual Studio Code for an editor-based workflow, the CLI for terminal and shell workflows, or embed the runtime when a Go or JavaScript application needs to execute FQL directly.

> **Alpha status**
>
> Ferret v2 is currently in alpha. The language, runtime, CLI, modules, and embedding APIs may change before beta. Pin versions for scripts, CI, and application integrations.

> **Looking for Ferret v1?**
>
> Ferret v1 remains available for existing projects, but new users should start with Ferret v2. See the [migration guide]({{< ref "/docs/tools/cli/migrate" >}}) for the supported mechanical migration steps.

## Use Visual Studio Code

Install the official [Ferret Lang extension](https://marketplace.visualstudio.com/items?itemName=ferretlang.ferret-lang), published by `ferretlang`, then open or create a `.fql` file. Use the editor actions to run or debug the current file. The extension includes the matching Ferret daemon, so normal use does not require a separate `ferretd` installation or `PATH` configuration.

See [Visual Studio Code]({{< ref "/docs/tools/visual-studio-code" >}}) for formatting, debugging, configuration, and troubleshooting.

## Use the CLI

Choose the CLI to run FQL files and expressions from a terminal, shell script, or CI job.

{{< terminal >}}
ferret version
{{< /terminal >}}

See [CLI Installation]({{< ref "/docs/tools/cli/installation" >}}) for prebuilt binaries, source installation, updates, and verification.

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

- Start with the [Quick Start]({{< ref "/docs/getting-started/quick-start" >}}) when you want to learn FQL in Visual Studio Code, the terminal, or the Playground.
- Read the [Embedding overview]({{< ref "/docs/embedding" >}}) when an application will own the runtime.
- Use [Worker]({{< ref "/docs/tools/worker" >}}) when an existing system needs Ferret through a separately deployed HTTP service.

{{< docs-related tiles="getting-started-quick-start,tools-visual-studio-code,tools-cli-installation,embedding-go,embedding-javascript,tools-worker" >}}
