---
title: "Installation"
weight: 20
draft: false
description: "Install Lab from release archives, the install script, Docker, or source."
---

# Installation

Lab is distributed as release archives, a Docker image, and Go source. After installing it, run `lab version` to verify both Lab and the selected Ferret runtime.

## Install from a release archive

Download an archive from the [Lab releases](https://github.com/MontFerret/lab/releases) page for your platform.

Release archives follow this naming pattern:

| Platform | Archive |
| --- | --- |
| Linux x86_64 | `lab_linux_x86_64.tar.gz` |
| Linux arm64 | `lab_linux_arm64.tar.gz` |
| macOS x86_64 | `lab_darwin_x86_64.tar.gz` |
| macOS arm64 | `lab_darwin_arm64.tar.gz` |
| Windows x86_64 | `lab_windows_x86_64.zip` |
| Windows arm64 | `lab_windows_arm64.zip` |

Extract the archive and place the `lab` binary in a directory on your `PATH`.

{{< terminal >}}
mkdir -p "$HOME/.ferret"
tar -xzf lab_linux_x86_64.tar.gz
mv lab "$HOME/.ferret/"
export PATH="$PATH:$HOME/.ferret"
{{< /terminal >}}

On Windows, extract `lab.exe` from the zip file and place it in a directory included in `PATH`.

## Use the install script

The install script detects the current OS and architecture, downloads the matching release archive, and copies the binary into the install directory.

The default install directory is `$HOME/.ferret`. The directory must already exist and be writable.

{{< terminal >}}
mkdir -p "$HOME/.ferret"
curl -fsSL https://raw.githubusercontent.com/MontFerret/lab/main/install.sh | sh
{{< /terminal >}}

Use `LOCATION` to install somewhere else:

{{< terminal >}}
mkdir -p "$HOME/bin"
curl -fsSL https://raw.githubusercontent.com/MontFerret/lab/main/install.sh -o install.sh
LOCATION="$HOME/bin" sh install.sh
{{< /terminal >}}

Use `VERSION` to install a specific release tag:

{{< terminal >}}
VERSION=v2.0.0-rc.18 sh install.sh
{{< /terminal >}}

The script updates a shell profile when it can find one. If the binary is not found after installation, add the install directory to `PATH` manually.

## Run with Docker

Lab release images are published as `montferret/lab` and `ghcr.io/montferret/lab`.

{{< terminal >}}
docker run --rm -v "$PWD/tests:/test" montferret/lab:latest
{{< /terminal >}}

The container's default command runs Lab against files mounted at `/test`. You can also pass an explicit Lab command:

{{< terminal >}}
docker run --rm -v "$PWD:/workspace" montferret/lab:latest \
  run --reporter=simple /workspace/tests
{{< /terminal >}}

The container entrypoint treats `run`, `serve`, `version`, `help`, and leading flag arguments as Lab invocations. Other commands are executed directly in the container.

## Build from source

Lab v2 is a Go module at `github.com/MontFerret/lab/v2`. The current source declares Go `1.25.6`, and CI uses Go `>=1.25`.

{{< terminal >}}
git clone https://github.com/MontFerret/lab.git
cd lab
make build
{{< /terminal >}}

The build writes the binary to `./bin/lab`.

For a narrower local compile, run:

{{< terminal >}}
go build -o ./bin/lab ./main.go
{{< /terminal >}}

## Verify the install

Run:

{{< terminal >}}
lab version
{{< /terminal >}}

The output includes the Lab binary version and the Ferret version reported by the selected runtime.

To verify a remote or binary runtime:

{{< terminal >}}
lab version --runtime=http://localhost:8080
lab version --runtime=bin:/usr/local/bin/ferret
{{< /terminal >}}
