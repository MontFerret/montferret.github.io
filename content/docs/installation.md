---
title: "Installation"
weight: 2
draft: false
---

# Installation

Even though, Ferret comes as a CLI executable, it also can be used as a library. 

### CLI

#### From binary

You can download latest binaries from [here](https://github.com/MontFerret/ferret/releases).

#### From the source

```bash
$ go get github.com/MontFerret/ferret
```

### Library

```bash
$ go get github.com/MontFerret/ferret/pkg/compiler
```

<hr />

# Environment
In order to use all Ferret features, you will need to have Chrome either installed locally or running in Docker. For ease of use we recommend to run Chrome inside a Docker container:

```bash
$ docker run -d -p 9222:9222 -e CHROME_OPTS='--disable-dev-shm-usage --force-gpu-mem-available-mb --full-memory-crash-report' alpeware/chrome-headless-stable:ver-83.0.4103.61
```


But if you want to see what's happening during query execution, just start your Chrome with remote debugging port:

```bash
$ chrome.exe --remote-debugging-port=9222
```
