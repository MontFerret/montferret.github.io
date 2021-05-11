---
title: "Installation"
weight: 2
draft: false
---

# Installation

### CLI

#### From binary

You can either download latest binaries from [here](https://github.com/MontFerret/cli/releases)

or use the following shell script:

```bash
$ curl https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | LOCATION=. sh
```

#### From the source

```bash
$ go get github.com/MontFerret/cli
```

### Library

```bash
$ go get github.com/MontFerret/ferret
```

<hr />

# Environment
In order to use all Ferret features, you will need to have Chrome/Chromium either installed locally or running in Docker. For ease of use we recommend to run Chrome/Chromium inside a Docker container.

You can use any Chromium-based headless image, but we've put together an image that's ready to go:

```bash
$ docker pull montferret/chromium
$ docker run -d -p 9222:9222 montferret/chromium
```

If you'd rather see what's happening during query execution, just start launch Chrome from your host with the remote debugging port set:

```bash
$ chrome.exe --remote-debugging-port=9222
```
