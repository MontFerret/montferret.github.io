---
title: "Overview"
sidebarTitle: "Overview"
weight: 10
draft: false
description: "What Ferret is and how its main pieces fit together."
aliases:
  - /docs/introduction/
---

# What is Ferret?
Ferret is a portable runtime and declarative query language for extracting structured data from web pages and documents. It lets you describe the data you need while Ferret handles the lower-level work of loading pages, querying HTML, waiting for browser state, and returning structured results.

Use Ferret when you need repeatable extraction for UI testing, machine learning pipelines, analytics, monitoring, or internal automation.

<div class="docs-card-grid docs-card-grid-compact">
  <a class="docs-card" href="/docs/getting-started/installation/">
    <span class="docs-card-kicker">1</span>
    <strong>Install Ferret</strong>
    <span>Set up the command-line tool or Go library.</span>
  </a>
  <a class="docs-card" href="/docs/getting-started/first-script/">
    <span class="docs-card-kicker">2</span>
    <strong>Run a first script</strong>
    <span>Query static HTML or a JavaScript-rendered page.</span>
  </a>
  <a class="docs-card" href="/docs/web-extraction/">
    <span class="docs-card-kicker">3</span>
    <strong>Query HTML</strong>
    <span>Understand documents, results, and extraction errors.</span>
  </a>
  <a class="docs-card" href="/docs/tools/lab/">
    <span class="docs-card-kicker">4</span>
    <strong>Test with Lab</strong>
    <span>Move repeatable Ferret scripts into test suites.</span>
  </a>
</div>

## How It Works

<img src="/img/design.png"  />

Ferret consists of the following main parts:

- FQL parser, compiler, and runtime
- Standard library and module functions
- User function registry for application-specific behavior
- HTML drivers and types for static documents and Chrome DevTools Protocol
- Ecosystem tools including CLI, Lab, and Worker

For static page parsing, no external browser is required. For JavaScript-rendered pages, Ferret can connect to Chrome or Chromium through the Chrome DevTools Protocol.

## Where To Go Next

- Read [Language](/docs/language/) when you want the FQL concepts behind a query.
- Read [Web Extraction](/docs/web-extraction/) when you want to work with documents and HTML results.
- Read [Tools](/docs/tools/) when you need CLI commands, Lab, Worker, Docker images, or editor integrations.
- Read [Reference](/docs/reference/) when you need detailed statement and API entries.
