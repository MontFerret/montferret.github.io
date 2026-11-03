---
title: "Wire Runtime"
weight: 75
draft: false
description: "Run Lab tests through the hosted Universal Ferret API over gRPC."
---

# Wire Runtime

Use Wire when a separately configured Ferret host should execute your tests through the Universal Ferret API over gRPC.

{{< terminal >}}
lab run --runtime=wire+tcp://127.0.0.1:50051 tests/
{{< /terminal >}}

Use Lab `v{{< data "versions.lab.v2" >}}` with a host such as the [published Wire example]({{< ref "/docs/embedding/wire" >}}), pinned to Wire `v{{< data "versions.wire.v1" >}}`, Universal API `v{{< data "versions.wire.api" >}}`, and Core `v{{< data "versions.wire.core" >}}`. Lab requires output framing version **1** and checks compatibility during the connection handshake. Incompatible peers fail at startup without a buffered-output fallback.

Lab still owns test discovery, suites, parameters, scheduling, local services, and reporting. Wire is an additional integration adapter; the [HTTP/Worker adapter](../http-runtime/) remains available.

## Configure the connection

The runtime target contains the complete endpoint. Only `wire+tcp://127.0.0.1:<port>` with a port from 1 to 65535 is accepted. Hostnames, IPv6, other addresses, credentials, paths, queries, and fragments are rejected. Lab translates the public URI to its existing internal TCP endpoint. The plaintext loopback transport is intended for trusted local development, including a host exposed through a local port mapping. The new URI scheme adds no encryption or authentication.

{{< terminal >}}
lab run \
  --runtime=wire+tcp://127.0.0.1:50051 \
  --runtime-option='connect-timeout:"10s"' \
  --param='value:"shared"' \
  tests/

LAB_RUNTIME=wire+tcp://127.0.0.1:50051 \
LAB_RUNTIME_OPTION='connect-timeout:"10s"' \
lab run tests/
{{< /terminal >}}

`connect-timeout` is the only Wire adapter option. It must be a positive JSON duration string and defaults to `"5s"`. It bounds transport establishment and the Connect handshake; cancellation or an earlier caller deadline also stops setup. The per-test `--timeout` remains separate and bounds execution. Unknown options, HTTP settings such as `headers` and `path`, and arbitrary FQL names fail before connecting. See [Configuration](../configuration/) for value syntax, environment precedence, and repeated keys.

## Pass parameters

FQL inputs use the existing `--param`, `--param-bind`, and YAML query/assertion facilities. Query and assertion parameter overlays remain independent. Adapter settings never become FQL inputs; ordinary parameters named `headers`, `cookies`, `path`, and `flags` remain valid.

Nested YAML objects with string keys are converted to Wire's portable object shape without changing binary or numeric values. Objects with non-string keys are rejected.

Lab preserves the source name, exact FQL text, user parameters, and materialized `@lab` values. It uses Wire's one-shot `api.Runtime.Run` by default. With repetitions or retries enabled, the case executor compiles scripts lazily, reuses successful plans within the case, and opens a fresh session for every invocation. Changed referenced script names or text select a different plan. Every successful invocation returns a one-shot output handle. Lab consumes it within the test deadline and checks the terminal error before marking an attempt complete. Suite assertions receive the query result only after the entire query output succeeds. Wire message limits do not cap aggregate output; Lab materializes test results, so choose queries and host policies that bound result sizes.

Go callers use the canonical `api.Runtime` contract. Caller session options run once in order and retain callback errors. Parameter setters preserve nested YAML conversion; other supported UAPI settings are forwarded. Compilation methods use the host compiler. A thin plan wrapper applies the same YAML conversion, validation, and option ordering to ordinary and debug sessions as one-shot execution. Metadata and closure retain the underlying plan contract. Lab closes each invocation's session after its output settles and releases cached plans before reporting the case; direct Go callers own their plans and sessions. Cleanup failures remain observable and cannot satisfy expected-error tests.

## Inspect the hosted version

{{< terminal >}}
lab version \
  --runtime=wire+tcp://127.0.0.1:50051 \
  --runtime-option='connect-timeout:"10s"'
{{< /terminal >}}

`Runtime:` is the exact version supplied by the hosted `api.Runtime.Version`, rather than the Wire package or host application version. Wire captures it during Connect; version reads do not require another RPC.

## Configure Ferret at the host

The hosted runtime owns Ferret configuration, modules, filesystem access, and outbound HTTP execution policy. Lab rejects local `--policy-fs-*` and `--policy-http-*` options for Wire. Configure those capabilities at the host.

Filesystem paths and browser endpoints used by FQL resolve in the host's environment. The Wire endpoint's loopback address identifies the client transport; it does not change where the host executes FQL.

Local fixture and mock URLs still arrive under `@lab.static` and `@lab.mock`. If the host runs in another environment, use `--serve-bind` and `--serve-host` so it can reach those services. See [static serving](../static-serving/).

## Connection lifetime

Lab creates and owns the gRPC client transport. Wire borrows it and owns its logical connection. After executions settle, Lab closes the logical runtime before the transport and retains cleanup failures. Failed setup also closes the transport.

Each adapter uses one physical connection. Transport loss fails the adapter; start a new Lab invocation to reconnect. Lab does not automatically redial or switch runtimes.

See [Migrate Lab runtime configuration]({{< ref "/docs/migrations/lab-runtime-configuration" >}}) when upgrading an earlier v2 configuration.

## Next steps

{{< docs-related tiles="embedding-wire,tools-lab-configuration,tools-lab-http-runtime,tools-lab-running-tests" >}}
