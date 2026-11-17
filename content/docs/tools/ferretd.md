---
title: "Ferret daemon"
sidebarTitle: "Ferret daemon"
weight: 35
draft: false
description: "Use ferretd for editor language services, workspace execution, and debugging against local or Wire-hosted runtimes."
---

# Ferret daemon

`ferretd` provides developer services for FQL editors and other clients. Use it to analyze an open document, run a saved workspace script, or debug FQL from an editor. An editor integration can supply the executable. Standalone clients can download the [daemon prerelease](https://github.com/MontFerret/ferretd/releases/tag/v{{< data "versions.ferretd.v1" >}}).

| Command | Use |
| --- | --- |
| `ferretd lsp` | Start the experimental language server over stdio for diagnostics, completion, navigation, semantic highlighting, and formatting. It can analyze unsaved text. |
| `ferretd serve` | Start the local gRPC service. Clients open a workspace, create a session from a saved FQL file, run it, and read its result. |
| `ferretd dap` | Start a single-session debug adapter over stdio. The editor controls source breakpoints, stepping, variables, and evaluation. |

An editor starts each service it uses. The debug adapter runs independently of `serve`. Both can use the builtin `local` runtime or connect to a configured [Wire host]({{< ref "/docs/embedding/wire" >}}). The host supplies the functions, modules, and runtime policies used by the script. The [CLI]({{< ref "/docs/tools/cli/wire" >}}) and [Lab]({{< ref "/docs/tools/lab/wire-runtime" >}}) have their own Wire connections.

## Configure a runtime

Start the [example Wire host]({{< ref "/docs/embedding/wire" >}}#host-the-configured-engine), then save this configuration as `runtimes.json`:

```json
{
  "defaultRuntime": "local",
  "runtimes": [
    {
      "name": "application",
      "provider": "wire",
      "endpoint": "tcp://127.0.0.1:50051",
      "insecure": true
    }
  ]
}
```

The example host uses plaintext on loopback. For a TLS host, omit `insecure` and use a hostname or IP address covered by the certificate. Ferretd verifies server identity using system trust. The optional `connectTimeout` is a positive duration string and defaults to `"5s"`.

Point the workspace service and debug adapter at the same file when both need those definitions. For example, a client can launch each process with:

{{< terminal >}}
ferretd serve --runtime-config /absolute/path/to/runtimes.json
{{< /terminal >}}

{{< terminal >}}
ferretd dap --runtime-config /absolute/path/to/runtimes.json
{{< /terminal >}}

Each process reads its configuration at startup. Restart the service or adapter after editing the file. A client can list the daemon’s configured names, then select `application` when creating a workspace session. A DAP client checks the adapter’s `supportsFerretRuntimeSelection` capability and sends `runtimeName: "application"` at launch. Omitting the selector uses the configured default, which is `local` unless changed. Connection occurs when creating the session or launching the debugger. An unavailable selected host produces an error.

## Source and filesystem access

The source file and debugger locations belong to the local workspace. The selected runtime interprets filesystem operations inside FQL. With no working-directory override, local execution uses the workspace root and Wire execution uses the host's default. An explicit Wire working directory is interpreted by the host and may use a different path convention.

Ending the debug session releases its plans and connection while leaving the application host running. If the host disappears while a session is paused, the adapter reports the loss and terminates the session. Start a new session after restoring the host. See the [DAP client guide](https://github.com/MontFerret/ferretd/blob/v{{< data "versions.ferretd.v1" >}}/docs/dap.md) for launch fields and adapter capabilities.

## Connect to the workspace service

The Go client negotiates the daemon API version. API 1.2 adds runtime discovery and API 1.3 adds explicit runtime selection when creating a session. With API 1.4, a client can give a session request an acquisition ID and use it to cancel or clean up an interrupted creation, even if the reply was lost.

By default, `serve` accepts local connections over a Unix socket on macOS and Linux or a named pipe on Windows. Clients that require loopback TCP can request an ephemeral port and authenticate using a bearer token. This incoming daemon connection is separate from an outbound connection to a Wire host. See the [daemon client and transport guide](https://github.com/MontFerret/ferretd/blob/v{{< data "versions.ferretd.v1" >}}/README.md#commands) for connection details.

{{< docs-related tiles="embedding-wire,tools-cli-wire,tools-visual-studio-code" >}}
