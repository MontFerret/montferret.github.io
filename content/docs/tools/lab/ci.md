---
title: "CI"
weight: 110
draft: false
description: "Run Lab in CI with stable output, waits, retries, Docker, and local fixtures."
---

# CI

Lab fits CI best when test dependencies are explicit and output is easy to scan. Use `--reporter=simple`, wait for external services, and keep local fixture URLs reachable from the selected runtime.

## Basic CI command

{{< terminal >}}
lab run tests/ --reporter=simple
{{< /terminal >}}

The process exits with status `0` when all tests pass and status `1` when any test fails.

## Wait for services

Use `--wait` for HTTP services that may start in parallel with the test job.

{{< terminal >}}
lab run tests/integration \
  --wait=http://127.0.0.1:9222/json/version \
  --wait-timeout=10 \
  --wait-attempts=12 \
  --reporter=simple
{{< /terminal >}}

Waits happen before the runtime is created and before test files execute.

## Retry unstable dependencies

Use `--attempts` for tests that can fail because a dependency is not fully settled even after a health check.

{{< terminal >}}
lab run tests/integration \
  --attempts=3 \
  --timeout=60 \
  --reporter=simple
{{< /terminal >}}

`--attempts=3` allows up to three attempts total for a failing test.

## Run in GitHub Actions

```yaml
name: Lab

on:
  pull_request:

jobs:
  lab:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Lab
        run: |
          mkdir -p "$HOME/.ferret"
          curl -fsSL https://raw.githubusercontent.com/MontFerret/lab/main/install.sh | sh
          echo "$HOME/.ferret" >> "$GITHUB_PATH"

      - name: Run Lab tests
        run: |
          lab run tests/ \
            --reporter=simple \
            --attempts=3 \
            --timeout=60
```

Use a pinned `VERSION` in the install step when the pipeline must stay on a specific Lab release.

## Run with Docker

The Docker image can run tests from a mounted directory:

{{< terminal >}}
docker run --rm -v "$PWD/tests:/test" montferret/lab:latest
{{< /terminal >}}

To pass explicit flags, mount the workspace and call `run`:

{{< terminal >}}
docker run --rm -v "$PWD:/workspace" montferret/lab:latest \
  run /workspace/tests \
  --reporter=simple \
  --attempts=3
{{< /terminal >}}

The image entrypoint starts the bundled browser entrypoint before invoking Lab for Lab commands.

## Serve local fixtures in CI

Use `--serve` for static fixtures and `--mock` for OpenAPI mock APIs:

{{< terminal >}}
lab run tests/e2e \
  --serve ./dist@app \
  --mock ./users.yaml@api \
  --reporter=simple
{{< /terminal >}}

Tests can read `@lab.static.app` and `@lab.mock.api`.

For the full local service syntax, see [Static File Server]({{< ref "static-serving" >}}) and [Mock API Server]({{< ref "mock-api" >}}).

If the selected Ferret runtime runs outside the Lab process, make the fixture services reachable from that runtime:

{{< terminal >}}
lab run tests/e2e \
  --runtime=https://ferret.example.com \
  --serve ./dist@app \
  --mock ./users.yaml@api \
  --serve-bind 0.0.0.0 \
  --serve-host host.docker.internal \
  --reporter=simple
{{< /terminal >}}

`--serve-host` controls the URLs passed to FQL. `--serve-bind` controls where Lab listens.

## Keep output deterministic

For CI logs:

- Prefer `--reporter=simple`.
- Keep `--concurrency` low when tests share external state.
- Use `--times` for repeatability checks.
- Use explicit aliases for `--serve` and `--mock` so FQL reads stable parameter names.

Example:

{{< terminal >}}
lab run tests/smoke \
  --concurrency=2 \
  --times=3 \
  --serve ./fixtures@fixtures \
  --reporter=simple
{{< /terminal >}}
