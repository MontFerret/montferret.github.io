---
title: "Installation"
weight: 3
draft: false
---

# Installation

Even though, Ferret comes as a CLI executable, it also can be used as a library. 

### CLI

#### From binary

You can download latest binaries from [here](https://github.com/MontFerret/ferret/releases).

#### From the source

{{< highlight bash >}}
go get github.com/MontFerret/ferret
{{</ highlight >}}

### Library

{{< highlight bash >}}
go get github.com/MontFerret/ferret/pkg/compiler
{{</ highlight >}}

# Environment

In order to use all Ferret features, you will need to have Chrome either installed locally or running in Docker. For ease of use we recommend to run Chrome inside a Docker container:

{{< highlight bash >}}
docker pull alpeware/chrome-headless-stable
docker run -d -p=0.0.0.0:9222:9222 --name=chrome-headless -v /tmp/chromedata/:/data alpeware/chrome-headless-stable
{{</ highlight >}}

But if you want to see what's happening during query execution, just start your Chrome with remote debugging port:

{{< highlight bash >}}
chrome.exe --remote-debugging-port=9222
{{</ highlight >}}
