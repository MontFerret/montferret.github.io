---
title: "Install"
weight: 10
draft: false
description: "Install, update, and verify the Ferret CLI."
---

# Install

The Ferret CLI runs on Linux, macOS, and Windows. Full installation instructions, including how to embed Ferret as a Go library and set up browser support, are in the [Installation](../../getting-started/installation/) guide.

This page covers CLI-specific setup.

## From a prebuilt binary

Download the latest release for your platform from the GitHub releases page:

https://github.com/MontFerret/cli/releases

After downloading, make sure the binary is available in your `PATH`.

## Using the install script

{{< terminal >}}
curl -fsSL https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | sh
{{< /terminal >}}

The script detects your platform, downloads the matching binary, and installs it into your PATH.

## From source

Requires `Go {{< data "versions.go" >}}` or later:

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

{{< docs-related tiles="getting-started-installation,tools-cli" >}}
