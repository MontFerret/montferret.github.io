---
title: "Worker"
weight: 30
draft: false
description: "Run Ferret as an HTTP service for remote FQL execution."
---

# Worker

Ferret Worker runs Ferret as a long-running HTTP service. It accepts FQL scripts over HTTP, executes them with the bundled Ferret runtime, and returns the serialized query result.

Use Worker when an application, job runner, CI system, or internal service needs to execute Ferret scripts remotely instead of starting a local runtime for each caller.

{{< docs-related tiles="tools-worker-overview,tools-worker-configuration,tools-worker-deployment,tools-worker-api" >}}