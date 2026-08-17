---
title: "Install"
weight: 10
draft: false
description: "Install, update, and verify the Ferret CLI."
---

# Install

The Ferret CLI runs on Linux, macOS, and Windows. This page covers CLI-specific setup. To compare the CLI with Go and JavaScript embedding, see the [Installation chooser]({{< ref "/docs/getting-started/installation" >}}).

## From a prebuilt binary

Download `v{{< data "versions.cli.v2" >}}` for your platform from the [CLI release page](https://github.com/MontFerret/cli/releases/tag/v{{< data "versions.cli.v2" >}}).

After downloading, make sure the binary is available in your `PATH`.

## From source

Requires `Go {{< data "versions.cli.go" >}}` or later:

{{< terminal >}}
go install github.com/MontFerret/cli/v2/ferret@v{{< data "versions.cli.v2" >}}
{{< /terminal >}}

## Verify installation

{{< terminal >}}
ferret version
{{< /terminal >}}

This prints the CLI version and the runtime version.

## Update

The CLI can update itself to the latest release:

{{< terminal >}}
ferret update self
{{< /terminal >}}

This downloads the latest binary for your platform and replaces the current installation.

## Uninstall

Remove the `ferret` binary from your PATH and, optionally, delete the configuration directory:

{{< terminal >}}
rm -rf ~/.ferret
{{< /terminal >}}

## Next steps

{{< docs-related tiles="tools-cli-run,tools-cli-mod,tools-cli-repl,getting-started-quick-start" >}}
