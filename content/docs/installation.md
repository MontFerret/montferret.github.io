---
title: "Installation"
weight: 2
draft: false
---

# Installation

Ferret can be used either as a **command-line tool** or as a **Go library**, depending on your use case.

- Use the **CLI** to run queries locally, in scripts, or in CI
- Use the **library** to embed Ferret into your Go applications

---

## CLI
### From binary (recommended)

Download the latest prebuilt binaries from the releases page:

https://github.com/MontFerret/cli/releases

Alternatively, you can install using the provided shell script:

```bash
curl -fsSL https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | sh
```
> This script downloads the correct binary for your platform and installs it into your PATH.

If you prefer to inspect the script first:

```bash
curl -fsSL https://raw.githubusercontent.com/MontFerret/cli/master/install.sh -o install.sh
less install.sh
sh install.sh
```

### From source
If you have Go installed, you can build and install the CLI from source:

```bash
go install github.com/MontFerret/cli/ferret@latest
```

### Library
To use Ferret as a Go library, add it to your project:

```bash
go get github.com/MontFerret/ferret@latest
```

# Runtime requirements

Ferret can extract data from both static HTML and JavaScript-rendered pages.
- Static pages work out of the box
- Dynamic pages require Chrome or Chromium

## Using Docker (recommended)

For most setups, running Chrome in Docker is the easiest option.

```bash
docker pull montferret/chromium
docker run -d -p 9222:9222 montferret/chromium
```

This starts a headless Chromium instance with the remote debugging port enabled.

## Using a local Chrome installation

You can also run Chrome or Chromium directly on your machine with remote debugging enabled:

```bash
chrome --remote-debugging-port=9222
```

On Windows:
```bash
chrome.exe --remote-debugging-port=9222
```