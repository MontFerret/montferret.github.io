---
title: "Installation"
weight: 30
draft: false
description: "Install Ferret as a command-line tool, embed it as a Go library, and configure browser support for dynamic pages."
aliases:
    - /docs/installation/
relatedTileOverrides:
  getting-started-quick-start:
    kicker: "Start here"
    title: "Quick Start"
    description: "Run your first Ferret script and learn the basic workflow."
---

# Installation

Ferret can be used in two main ways:

- as a command-line tool for running Ferret scripts locally, in shell scripts, or in CI
- as a Go library for embedding Ferret into your own applications

The official CLI includes the standard modules needed for common web extraction workflows. When embedding Ferret as a library, HTML and browser support are added explicitly by registering the relevant modules.

> **Alpha status**
>
> Ferret v2 is currently in alpha. It is ready for experimentation, feedback, prototypes, internal tools, and early integration work, but the language, runtime, CLI, modules, and embedding APIs may still change before beta.
>
> Pin CLI and Go module versions when using Ferret in scripts, CI, or embedded applications.

> **Looking for Ferret v1?**
>
> Ferret v1 remains available for existing projects, but new users should start with Ferret v2. See the migration guide for differences between v1 and v2.

## Install the CLI

The Ferret CLI is the easiest way to run Ferret from your terminal.

### From a prebuilt binary

Download the latest release for your platform from the GitHub releases page:

https://github.com/MontFerret/cli/releases

After downloading the binary, make sure it is available in your `PATH`.

You can verify the installation with:

{{< terminal >}}
ferret version
{{< /terminal >}}

### Using the install script

You can also install the CLI with the provided shell script:

{{< terminal >}}
curl -fsSL https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | sh
{{< /terminal >}}

The script detects your platform, downloads the matching binary, and installs it into your PATH.

If you prefer to inspect the script before running it:

{{< terminal >}}
curl -fsSL https://raw.githubusercontent.com/MontFerret/cli/master/install.sh -o install.sh
less install.sh
sh install.sh
{{< /terminal >}}

### From source

If you already have Go installed, you can build and install a specific Ferret CLI version from source.

This requires `Go {{< data "versions.go" >}}` or later.

{{< terminal >}}
go install github.com/MontFerret/cli/v2/ferret@v{{< data "versions.cli.v2" >}}
{{< /terminal >}}

During the alpha stage, prefer installing a specific tagged version instead of using `@latest`.

Verify the installation with:

{{< terminal >}}
ferret version
{{< /terminal >}}

## Add Ferret to a Go project

Ferret can also be embedded into Go applications.

Add the module to your project with:

{{< terminal >}}
go get github.com/MontFerret/ferret/v2@v{{< data "versions.runtime.v2" >}}
{{< /terminal >}}

This is useful when you want to use Ferret as a data extraction engine inside your own services, workers, tools, or automation pipelines.

## Browser and HTML support

Ferret’s core runtime is intentionally small. HTML querying, browser automation, and related web capabilities are provided by modules.

The official CLI distribution includes the standard web/HTML modules, so common web extraction workflows work out of the box when using the CLI.

If you embed Ferret as a Go library, you decide which modules to register in your runtime. This lets applications keep their Ferret environment small and capability-oriented, while still enabling HTML, browser, or other integrations when needed.

| Use case | Available in core? | Available in CLI? | Notes |
| --- | --- | --- | --- |
| General Ferret language runtime | Yes | Yes | Expressions, control flow, values, functions, and execution |
| Static HTML querying | No | Yes | Provided by the HTML/web module bundled with the CLI |
| Browser automation | No | Yes | Requires the browser module and a Chrome/Chromium runtime |
| JavaScript-rendered pages | No | Yes | Requires Chrome or Chromium through CDP |
| Custom application integrations | Via modules | Depends on distribution | Register only the capabilities your application needs |

## Browser runtime requirements

Ferret does not require a browser for every workflow.

You only need Chrome or Chromium when using browser-backed features, such as querying JavaScript-rendered pages, waiting for page state, or dispatching browser events.

For browser-based workflows, Ferret connects to Chrome or Chromium through the Chrome DevTools Protocol, usually on port `9222`.

## Run Chromium with Docker

For most setups, running Chromium in Docker is the easiest option.

{{< terminal >}}
docker pull montferret/chromium
docker run -d -p 9222:9222 montferret/chromium
{{< /terminal >}}

This starts a headless Chromium instance with the remote debugging port enabled.

You can check that the browser is running with:

{{< terminal >}}
curl http://127.0.0.1:9222/json/version
{{< /terminal >}}

If the command returns browser metadata, Chromium is ready to use.

## Use a local Chrome or Chromium installation

You can also run Chrome or Chromium directly on your machine with remote debugging enabled.

On macOS or Linux:

{{< terminal >}}
chrome --remote-debugging-port=9222
{{< /terminal >}}

Depending on your system, the executable may also be named google-chrome, chromium, or chromium-browser.

On Windows:

{{< terminal >}}
powershell chrome.exe --remote-debugging-port=9222
{{< /terminal >}}

Once Chrome is running, verify that the debugging endpoint is available:

{{< terminal >}}
curl http://127.0.0.1:9222/json/version
{{< /terminal >}}

## Next steps

After installing Ferret, choose where you want to go next.

{{< docs-related tiles="getting-started-quick-start,tools-cli,embedding,tools-lab" >}}
