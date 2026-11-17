---
title: "Deploy Ferret as a service"
sidebarTitle: "Deploying"
weight: 80
draft: false
description: "Run Ferret Worker as a long-running HTTP service for remote script execution."
---

# Deploy Ferret as a service

[Ferret Worker]({{< ref "/docs/tools/worker" >}}) is an HTTP service that accepts FQL scripts and returns results. This guide walks through deploying Worker with Docker for remote or shared use.

## Start Worker with Docker

The release image includes Chromium. Start the service on your machine:

{{< terminal command="true" >}}
docker run --rm -d --name ferret-worker -p 127.0.0.1:8080:8080 montferret/worker:{{< data "versions.worker.v2" >}}
{{< /terminal >}}

Once Worker is ready, its health check returns HTTP 200 with an empty body:

{{< terminal command="true" >}}
curl --fail --silent --show-error http://localhost:8080/health
{{< /terminal >}}

## Execute a query

Send FQL source in `text` and bind values in `params`. A successful response contains the encoded FQL result.

{{< terminal command="true" >}}
curl --fail --silent --show-error http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"text":"return @count + 1","params":{"count":2}}'
{{< /terminal >}}

## Use the CLI or Lab

With a local `script.fql` or a Lab suite under `tests/`, point the client at Worker:

{{< terminal command="true" >}}
ferret run --runtime http://localhost:8080 script.fql
lab run --runtime=http://localhost:8080 tests/
{{< /terminal >}}

Worker supplies the modules, browser connection, filesystem, and outbound HTTP policy for those scripts. For browser setup, container configuration, and exposing the service through a reverse proxy, see [Worker deployment]({{< ref "/docs/tools/worker/deployment" >}}). The [HTTP API]({{< ref "/docs/tools/worker/api" >}}) describes query errors and health checks.

When you finish, run `docker stop ferret-worker`.

## Next steps

{{< docs-related tiles="tools-worker,tools-worker-configuration,tools-worker-deployment,tools-worker-api" >}}
