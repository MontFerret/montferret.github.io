---
title: "Lab"
weight: 20
draft: false
description: "Test Ferret scripts with Lab, local fixtures, mock APIs, and alternate runtimes."
---

# Lab

Lab is the Ferret test runner. It discovers FQL and YAML test files, prepares parameters and local fixture services, executes each test through a Ferret runtime, and reports the result.

Use Lab when you want repeatable checks for extraction scripts, browser-backed flows, local static fixtures, mock REST APIs, or runtime compatibility.

<div class="docs-section-card-grid">
  <a class="docs-section-card" href="/docs/tools/lab/overview/">
    <strong>Overview</strong>
    <span>How Lab runs tests and where it fits with Ferret runtimes.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/installation/">
    <strong>Installation</strong>
    <span>Install Lab from releases, Docker, or source.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/writing-tests/">
    <strong>Writing Tests</strong>
    <span>Write FQL unit tests and YAML query/assert suites.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/running-tests/">
    <strong>Running Tests</strong>
    <span>Run local, HTTP, and Git-hosted test files with parameters and waits.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/static-serving/">
    <strong>Local Services</strong>
    <span>Serve static directories and mock APIs during tests or standalone.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/runners/">
    <strong>Runners</strong>
    <span>Control concurrency, retries, repeated runs, timeouts, and reporters.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/http-runtime/">
    <strong>HTTP Runtime</strong>
    <span>Run tests against a remote Ferret HTTP runtime.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/binary-runtime/">
    <strong>Binary Runtime</strong>
    <span>Run tests through an external Ferret-compatible binary.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/cdp/">
    <strong>CDP</strong>
    <span>Use browser-backed FQL with Lab and the selected runtime.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/configuration/">
    <strong>Configuration</strong>
    <span>Reference Lab flags, defaults, environment variables, and runtime params.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/lab/ci/">
    <strong>CI</strong>
    <span>Run Lab in pipelines with waits, retries, Docker, and local fixtures.</span>
  </a>
</div>
