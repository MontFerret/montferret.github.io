---
title: "Installation"
weight: 10
draft: false
description: "Install, update, and verify the Ferret CLI."
aliases:
  - /docs/tools/cli/install/
---

# Installation

The Ferret CLI runs on Linux, macOS, and Windows. This page covers CLI-specific setup. To compare the CLI with Go and JavaScript embedding, see the [Installation chooser]({{< ref "/docs/getting-started/installation" >}}).

## Prebuilt binary

Download `v{{< data "versions.cli.v2" >}}` for your platform from the [CLI release page](https://github.com/MontFerret/cli/releases/tag/v{{< data "versions.cli.v2" >}}).

Release archives use these names:

| Platform | Archive |
| --- | --- |
| Linux x86_64 | `cli_linux_x86_64.tar.gz` |
| Linux arm64 | `cli_linux_arm64.tar.gz` |
| macOS x86_64 | `cli_darwin_x86_64.tar.gz` |
| macOS arm64 | `cli_darwin_arm64.tar.gz` |
| Windows x86_64 | `cli_windows_x86_64.zip` |
| Windows arm64 | `cli_windows_arm64.zip` |

Extract the archive and place the `ferret` binary in a directory on your `PATH`. For a user-local installation on Linux or macOS:

{{< terminal >}}
mkdir -p "$HOME/.ferret"
mv ferret "$HOME/.ferret/"
export PATH="$HOME/.ferret:$PATH"
{{< /terminal >}}

On Windows, extract `ferret.exe` from the zip file and place it in a directory included in `PATH`.

## Build from source

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

These docs track the Ferret v2 prerelease channel. To update within that channel, install the newer tagged archive or rerun `go install` with the newer v2 tag.

The `ferret update self` command checks GitHub's latest non-prerelease release, so it does not advance between v2 prereleases. It also cannot install the Windows zip release. Do not use it as the v2 prerelease update path.

## Uninstall

Remove the `ferret` binary from the directory where you installed it. For the user-local path shown above:

{{< terminal >}}
rm "$HOME/.ferret/ferret"
{{< /terminal >}}

The CLI configuration remains at `$HOME/.ferret/config.yaml`. Remove that file separately only when you also want to discard the saved configuration. Do not remove the shared `$HOME/.ferret` directory if it contains Lab or other files.

## Next steps

{{< docs-related tiles="tools-cli-run,tools-cli-mod,tools-cli-repl,getting-started-quick-start" >}}
