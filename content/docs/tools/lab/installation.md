---
title: "Installation"
weight: 20
draft: false
description: "Install Lab from a release archive, run it with Docker, or build it from source."
---

# Installation

Lab is distributed as release archives, Docker images, and Go source. After installing it, run `lab version` to verify both Lab and the selected Ferret runtime.

## Prebuilt binary

Download `v{{< data "versions.lab.v2" >}}` for your platform from the [Lab release page](https://github.com/MontFerret/lab/releases/tag/v{{< data "versions.lab.v2" >}}).

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
export PATH="$HOME/.ferret:$PATH"
{{< /terminal >}}

On Windows, extract `lab.exe` from the zip file and place it in a directory included in `PATH`.

## Docker

Lab release images are published to Docker Hub and GitHub Container Registry. Use the documented version tag so the container runs the same Lab version as these pages.

{{< terminal >}}
docker run --rm -v "$PWD/tests:/test" montferret/lab:{{< data "versions.lab.v2" >}}
{{< /terminal >}}

The equivalent GitHub Container Registry image is `ghcr.io/montferret/lab:{{< data "versions.lab.v2" >}}`.

The container's default command runs Lab against files mounted at `/test`. You can also pass an explicit Lab command:

{{< terminal >}}
docker run --rm -v "$PWD:/workspace" montferret/lab:{{< data "versions.lab.v2" >}} \
  run --reporter=simple /workspace/tests
{{< /terminal >}}

The container entrypoint treats `run`, `serve`, `version`, `help`, and leading flag arguments as Lab invocations. Other commands are executed directly in the container.

## Build from source

Lab v2 is a Go module at `github.com/MontFerret/lab/v2`. Building `v{{< data "versions.lab.v2" >}}` requires Go `{{< data "versions.lab.go" >}}` or later.

{{< terminal >}}
git clone --branch v{{< data "versions.lab.v2" >}} --depth 1 https://github.com/MontFerret/lab.git
cd lab
make build
{{< /terminal >}}

The build writes the binary to `./bin/lab`.

For a narrower local compile, run:

{{< terminal >}}
go build -o ./bin/lab ./main.go
{{< /terminal >}}

## Verify installation

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

## Update

For a prebuilt installation, download the newer release archive and replace the existing `lab` binary. For Docker, change the image tag to the newer release and pull or run that image. Source builds should check out the newer release tag before rebuilding.

## Uninstall

Remove the `lab` binary from the directory where you installed it. For the user-local path shown above:

{{< terminal >}}
rm "$HOME/.ferret/lab"
{{< /terminal >}}

Do not remove the shared `$HOME/.ferret` directory if it contains the CLI configuration or other files. On Windows, remove `lab.exe` from its installation directory.

## Next steps

{{< docs-related tiles="tools-lab-writing-tests,tools-lab-running-tests,tools-lab-configuration" >}}
