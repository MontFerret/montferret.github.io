---
title: "CLI"
weight: 10
draft: false
description: "Use the Ferret command-line interface to run, debug, build, and extend FQL applications."
aliases:
  - /docs/cli/
---

# CLI

The Ferret CLI is the primary tool for working with FQL scripts. Use it to execute queries, start an interactive shell, debug scripts, format source code, compile to bytecode, manage modules and browser instances, and configure persistent settings.

All commands accept `--log-level`, `--log-output`, and `--log-file` flags for controlling log output. Settings can be persisted through the [`config`](configuration/) command or environment variables with the `FERRET_` prefix. See [Configuration](configuration/) for details.

{{< docs-related tiles="tools-cli-installation,tools-cli-run,tools-cli-repl,tools-cli-debug,tools-cli-check,tools-cli-fmt,tools-cli-build,tools-cli-inspect,tools-cli-browser,tools-cli-mod,tools-cli-migrate,tools-cli-configuration" >}}
