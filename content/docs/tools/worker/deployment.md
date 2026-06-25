---
title: "Deployment"
weight: 30
draft: false
description: "Deploy Worker with Docker images, release binaries, source builds, health checks, and reverse proxies."
---

# Deployment

Worker can run from a Docker image, a release binary, or a local source build. Docker is the usual deployment path when scripts need the CDP browser driver, because the Worker image includes a Chromium runtime.

## Run with Docker

Worker images are published to Docker Hub and GitHub Container Registry.

{{< terminal >}}
docker run --rm -p 8080:8080 montferret/worker:latest
{{< /terminal >}}

{{< terminal >}}
docker run --rm -p 8080:8080 ghcr.io/montferret/worker:latest
{{< /terminal >}}

Pin a release tag for repeatable deployments:

{{< terminal >}}
docker run --rm -p 8080:8080 montferret/worker:v2.0.0-rc.16
{{< /terminal >}}

The release image starts Chromium through the image entrypoint and then starts Worker on port `8080`.

## Configure the container

The image's default command starts Chromium and then starts Worker. Use environment variables for most container configuration:

{{< terminal >}}
docker run --rm -p 8080:8080 \
  -e LOG_LEVEL=info \
  -e BODY_LIMIT=10M \
  -e REQUEST_LIMIT=10 \
  montferret/worker:latest
{{< /terminal >}}

If you need to pass flags instead, keep the Chromium startup command:

{{< terminal >}}
docker run --rm -p 8080:8080 montferret/worker:latest \
  /bin/sh -c '/entrypoint.sh & /worker -log-level=info -body-limit=10M -request-limit=10'
{{< /terminal >}}

Mount a directory and set `FS_ROOT` when scripts need filesystem access:

{{< terminal >}}
docker run --rm -p 8080:8080 \
  -v "$PWD/data:/data:ro" \
  -e FS_ROOT=/data \
  montferret/worker:latest
{{< /terminal >}}

Use a read-only mount when scripts only need to read files.

## Run from a release binary

Download a Worker release archive from GitHub:

https://github.com/MontFerret/worker/releases

Release archives follow the GoReleaser platform naming pattern, such as `worker_linux_x86_64.tar.gz`, `worker_darwin_arm64.tar.gz`, or `worker_windows_x86_64.zip`.

After extracting the archive, place the `worker` binary on `PATH` and run:

{{< terminal >}}
worker -version
worker
{{< /terminal >}}

Local binary deployments need a reachable Chrome or Chromium process when scripts use `{ driver: "cdp" }`.

## Build from source

The Worker source is a Go module at `github.com/MontFerret/worker`. Use the Go version declared in the repository's `go.mod`.

{{< terminal >}}
git clone https://github.com/MontFerret/worker.git
cd worker
make compile
./bin/worker -version
{{< /terminal >}}

The build embeds Worker and Ferret version metadata through linker flags when built by the Makefile or release pipeline.

## Chrome dependency

Worker waits for Chrome during startup unless Chrome is disabled with `-no-chrome`. By default it checks:

```text
http://127.0.0.1:9222
```

When running outside the Docker image, start Chrome with remote debugging enabled before starting Worker:

{{< terminal >}}
google-chrome --headless --remote-debugging-port=9222
worker
{{< /terminal >}}

If Chrome runs on another host or port, set `-chrome-ip` and `-chrome-port`.

## Health checks

Use `GET /health` for load balancer and orchestration checks.

{{< terminal >}}
curl -i http://localhost:8080/health
{{< /terminal >}}

Status codes:

| Status | Meaning |
| --- | --- |
| `200` | Worker is running. If Chrome is enabled, Worker could reach Chrome's version endpoint. |
| `424` | Worker is running, but a required dependency such as Chrome is unavailable. |

The health endpoint returns an empty body and is skipped by Worker rate limiting.

## Network and security

Worker binds to `0.0.0.0` on the configured port. Put it behind your normal deployment boundary: a private network, reverse proxy, ingress controller, API gateway, or service mesh.

Worker does not include built-in authentication, authorization, TLS termination, request signing, durable queues, scheduling, or tenant isolation. Add those at the deployment layer when Worker is exposed beyond a trusted network.

Use body limits and rate limits for public or shared deployments:

{{< terminal >}}
worker \
  -log-level=info \
  -body-limit=10M \
  -request-limit=10 \
  -request-limit-time-window=60
{{< /terminal >}}

The historical baseline in the Worker README is 2 CPU and 2 GB of RAM. Browser-backed scripts, large pages, PDF generation, screenshots, and high concurrency may require more.
