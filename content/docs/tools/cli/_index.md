---
title: "CLI"
weight: 10
draft: false
description: "Use the Ferret command-line interface to run, debug, format, build, and inspect FQL scripts."
aliases:
  - /docs/cli/
---

# CLI

The Ferret CLI is the primary tool for working with FQL scripts. Use it to execute queries, start an interactive shell, debug scripts, format source code, compile to bytecode, manage browser instances, and configure persistent settings.

All commands accept `--log-level`, `--log-output`, and `--log-file` flags for controlling log output. Settings can be persisted through the [config](config/) command or environment variables with the `FERRET_` prefix. See [Config](config/) for details.

<div class="docs-section-card-grid">
  <a class="docs-section-card" href="/docs/tools/cli/install/">
    <strong>Install</strong>
    <span>Install, update, and verify the Ferret CLI.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/run/">
    <strong>Run</strong>
    <span>Execute FQL scripts from files, inline expressions, stdin, or compiled artifacts.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/repl/">
    <strong>REPL</strong>
    <span>Launch an interactive FQL shell for experimenting with queries.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/debug/">
    <strong>Debug</strong>
    <span>Debug FQL scripts interactively with breakpoints, stepping, and variable inspection.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/check/">
    <strong>Check</strong>
    <span>Verify FQL scripts for syntax and semantic errors without executing them.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/fmt/">
    <strong>Fmt</strong>
    <span>Format FQL scripts to a consistent style.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/build/">
    <strong>Build</strong>
    <span>Compile FQL scripts into bytecode artifacts for faster execution.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/inspect/">
    <strong>Inspect</strong>
    <span>Disassemble compiled programs and examine bytecode, constants, and source spans.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/browser/">
    <strong>Browser</strong>
    <span>Start and stop managed browser instances for browser-backed scripts.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/cli/config/">
    <strong>Config</strong>
    <span>Manage persistent CLI configuration for runtime, browser, and logging settings.</span>
  </a>
</div>
