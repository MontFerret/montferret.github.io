---
title: "Playground"
weight: 50
draft: false
description: "Try FQL in the website editor and find the standalone Playground service."
---

# Playground

Open [Try Ferret]({{< ref "/try" >}}), edit the query, and select **Run**. Try `RETURN 1 + 2` to get `3`. The v2 editor sends source and bind parameters to a hosted [Worker]({{< ref "/docs/tools/worker" >}}) and displays the JSON response. Its modules, browser access and filesystem belong to that host. Use your own runtime for private inputs.

The separate [Playground service](https://github.com/MontFerret/playground) exposes an HTTP API for environment discovery, analysis, formatting, compilation and execution. Its [published API contract](https://github.com/MontFerret/playground/blob/v1.0.0-alpha.3/api/openapi.yaml) describes the endpoints and response shapes for applications integrating that service.
