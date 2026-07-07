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

The simplest deployment uses the Worker Docker image:

{{< terminal command="true" >}}
docker run -d -p 8080:8080 montferret/worker
{{< /terminal >}}

Worker is now listening on port 8080. Test it with a curl request:

{{< terminal command="true" >}}
curl -s -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{"query": "RETURN 1 + 1"}'
{{< /terminal >}}

## Add browser support

For scripts that need the `cdp` driver, run Worker alongside a Chromium container:

```yaml
# docker-compose.yml
services:
  worker:
    image: montferret/worker
    ports:
      - "8080:8080"
    environment:
      - FERRET_WORKER_CHROME_URL=ws://chromium:9222
    depends_on:
      - chromium

  chromium:
    image: montferret/chromium
```

{{< terminal command="true" >}}
docker compose up -d
{{< /terminal >}}

Now Worker can execute browser-backed scripts:

{{< terminal command="true" >}}
curl -s -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "LET p = WEB::HTML::OPEN(\"https://mockery.ferretlang.org\", { driver: \"cdp\" })\nRETURN p.title"
  }'
{{< /terminal >}}

## Configure Worker

Worker is configured through environment variables or command-line flags. Common settings:

| Variable | Default | Description |
| --- | --- | --- |
| `FERRET_WORKER_PORT` | `8080` | HTTP listen port |
| `FERRET_WORKER_CHROME_URL` | — | WebSocket URL to a Chrome instance |
| `FERRET_WORKER_LOG_LEVEL` | `info` | Log verbosity |
| `FERRET_WORKER_QUERY_TIMEOUT` | `30s` | Maximum query execution time |
| `FERRET_WORKER_RATE_LIMIT` | — | Maximum requests per second |
| `FERRET_WORKER_FS_ENABLED` | `false` | Allow file system access |

See [Worker Configuration]({{< ref "/docs/tools/worker/configuration" >}}) for the full list.

## Pass parameters

Send parameters alongside the query:

{{< terminal command="true" >}}
curl -s -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "LET p = WEB::HTML::OPEN(@url)\nRETURN p.title",
    "params": { "url": "https://mockery.ferretlang.org" }
  }'
{{< /terminal >}}

## Health checks

Worker exposes a health endpoint:

{{< terminal command="true" >}}
curl -s http://localhost:8080/api/health
{{< /terminal >}}

Use this in Docker health checks or load balancer probes:

```yaml
services:
  worker:
    image: montferret/worker
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 5s
      retries: 3
```

## Run behind a reverse proxy

For production deployments, place Worker behind a reverse proxy (nginx, Caddy, Traefik) to handle TLS, rate limiting, and authentication.

Example nginx configuration:

```nginx
upstream worker {
    server 127.0.0.1:8080;
}

server {
    listen 443 ssl;
    server_name ferret.example.com;

    location /api/ {
        proxy_pass http://worker;
        proxy_read_timeout 60s;
    }
}
```

## Use Worker as a remote runtime

The Ferret CLI and Lab can execute scripts against a remote Worker instance:

{{< terminal command="true" >}}
ferret run --runtime http://localhost:8080 script.fql
{{< /terminal >}}

Lab supports this with the HTTP runtime:

{{< terminal command="true" >}}
lab --dir ./tests --runtime http://localhost:8080
{{< /terminal >}}

See [HTTP Runtime]({{< ref "/docs/tools/lab/http-runtime" >}}) for details.

## Next steps

{{< docs-related tiles="tools-worker,tools-worker-configuration,tools-worker-deployment,tools-worker-api" >}}
