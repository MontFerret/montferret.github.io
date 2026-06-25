---
title: "Worker"
weight: 30
draft: false
description: "Run Ferret as an HTTP service for remote FQL execution."
---

# Worker

Ferret Worker runs Ferret as a long-running HTTP service. It accepts FQL scripts over HTTP, executes them with the bundled Ferret runtime, and returns the serialized query result.

Use Worker when an application, job runner, CI system, or internal service needs to execute Ferret scripts remotely instead of starting a local runtime for each caller.

<div class="docs-section-card-grid">
  <a class="docs-section-card" href="/docs/tools/worker/overview/">
    <strong>Overview</strong>
    <span>Learn where Worker fits, how execution works, and what runtime capabilities it provides.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/worker/configuration/">
    <strong>Configuration</strong>
    <span>Set ports, logging, limits, filesystem access, caching, and Chrome connection options.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/worker/deployment/">
    <strong>Deployment</strong>
    <span>Run Worker with Docker, release binaries, source builds, health checks, and reverse proxies.</span>
  </a>
  <a class="docs-section-card" href="/docs/tools/worker/api/">
    <strong>API</strong>
    <span>Use the HTTP query endpoint, health endpoint, info endpoint, and remote runtime contract.</span>
  </a>
</div>
