---
title: "Visual Studio Code"
weight: 40
draft: false
description: "Write, format, run, and debug FQL with the official Visual Studio Code integration."
---

# Visual Studio Code

The official Ferret integration for Visual Studio Code provides language intelligence, formatting, execution, and debugging for `.fql` files. It is backed by the Ferret daemon and includes the matching daemon executable for supported platforms.

Use Visual Studio Code when you want to write and run Ferret programs from an editor. The [CLI]({{< ref "/docs/tools/cli" >}}) remains available for terminal, shell, and CI workflows.

## Installation

Ferret Lang v0.2.0 requires Visual Studio Code 1.95 or later.

1. From the [Ferret Lang v0.2.0 release](https://github.com/MontFerret/editorium/releases/tag/vscode/v0.2.0), download the `.vsix` for your extension host's operating system and architecture.
2. In the Command Palette, run **Extensions: Install from VSIX…** and select the downloaded file.
3. Open or create a `.fql` file. Visual Studio Code recognizes it as a Ferret document.

The official [Visual Studio Marketplace listing](https://marketplace.visualstudio.com/items?itemName=ferretlang.ferret-lang) is published by `ferretlang`. Use the release VSIX to install this exact version.

## Bundled Ferret daemon

The extension bundles `ferretd` 1.0.0-alpha.13 for macOS, Linux, and Windows on ARM64 and x64.

VS Code installs the package for the extension host. In Remote SSH, Dev Containers, WSL, and Codespaces, use the package for the remote host platform.

## Language intelligence

The extension provides syntax highlighting and standard editing support for `.fql` files. For saved, file-backed documents, the bundled daemon adds:

- diagnostics in the Problems view and editor
- completion suggestions
- hover information and signature help
- document symbols and outline navigation
- go to definition and find references within the current document
- semantic highlighting
- document formatting

Untitled Ferret documents still receive syntax highlighting, bracket handling, indentation, and folding. Save the document as a `.fql` file before using daemon-backed language features, execution, or debugging.

## Formatting

Use **Format Document** or **Format Document With...** and select **Ferret Lang**. Formatting uses the Ferret formatter exposed by the running language server and includes unsaved editor changes.

To select Ferret as the default formatter and format on save, add these language-specific settings:

```json
{
  "[ferret]": {
    "editor.defaultFormatter": "ferretlang.ferret-lang",
    "editor.formatOnSave": true
  }
}
```

## Running programs

Open a saved `.fql` file and select the Run button in the editor title bar, or run **Ferret: Run Current File** from the Command Palette. If the document has unsaved changes, select **Save and Run** so the daemon executes the saved version.

Results, elapsed time, compilation errors, and runtime errors appear in **Ferret Execution** under **View > Output**. Run **Ferret: Show Output** to reveal that channel later.

While the current file is running, the editor action changes to Cancel. Select it or run **Ferret: Cancel Execution** to stop that execution.

### Current-file settings

**Run Current File** and **Debug Current File** use these settings for the source file. Parameter names omit the `@` prefix. For example, `{"limit": 10}` binds `@limit` to `10`.

| Setting | Default | Purpose |
| --- | --- | --- |
| `ferret.run.parameters` | `{}` | JSON object containing FQL bind parameter values. |
| `ferret.run.workingDirectory` | Empty | Runtime filesystem directory. A blank setting uses the selected runtime's default. Local relative paths resolve from the compilation workspace. Wire hosts receive the directory as entered. |

## Debugging

Open a saved `.fql` file and select the Debug button, run **Ferret: Debug Current File**, or press `F5` for the current file. The containing workspace folder is the compilation root. A standalone file uses its own directory. You can debug the current file without creating `launch.json`.

Set breakpoints in the editor gutter before or during a session. Use the standard Visual Studio Code controls to continue, pause, step over, step into, step out, or stop execution. The Run and Debug views provide the call stack, variables, watches, and debug console while the program is paused.

### Configure launch.json

Create `.vscode/launch.json` when a program needs explicit paths, parameters, or entry behavior:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "ferret",
      "request": "launch",
      "name": "Debug API scraper",
      "program": "${workspaceFolder}/scripts/scrape.fql",
      "cwd": "${workspaceFolder}",
      "runtimeName": "local",
      "parameters": {
        "baseUrl": "https://example.com",
        "limit": 10
      },
      "stopOnEntry": true
    }
  ]
}
```

The Ferret launch properties are:

| Property | Required | Description |
| --- | --- | --- |
| `program` | Yes | Path to the `.fql` program. Zero-configuration debugging uses the active file. |
| `cwd` | No | Local compilation root containing the program. An explicit launch defaults to the program's directory. Current-file debugging prefers the containing workspace folder. |
| `runtimeName` | No | Omit to use the runtime picker selection. Use `""` to follow the daemon default or specify an exact, case-sensitive runtime name such as `local`. |
| `workingDirectory` | No | Runtime filesystem directory. A blank or omitted value uses the runtime default. Local relative paths resolve from `cwd`. For Wire, the host interprets the path after VS Code resolves any launch variables. |
| `parameters` | No | Object containing values available to FQL as bind parameters. The default is no parameters. |
| `stopOnEntry` | No | Pause before normal execution begins. The default is `false`. |

## Configure runtimes

Set `ferret.runtime.configPath` to a [daemon runtime configuration file]({{< ref "/docs/tools/ferretd#configure-a-runtime" >}}) on the extension host. Use a literal path. Relative paths resolve from the workspace folder in a window with exactly one file-backed folder. Use an absolute path in other windows. A blank setting uses the built-in configuration.

Run **Ferret: Select Runtime** from the Command Palette. The catalog shows names, labels, and providers, including Wire hosts that are offline. Select **Daemon default** to follow the configured default, `local` for built-in execution, or a named Wire runtime. The selection is saved for future Run and Debug launches in that workspace. An explicit `runtimeName` in a debug configuration overrides it.

After the configuration is first used, run **Ferret: Reload Runtimes** when you edit the file or change or clear the path. Finish or stop active Run and Debug sessions and wait for them to close before reloading. Use **Ferret: Retry Runtime Discovery** if catalog discovery failed and the daemon is still running. If daemon startup failed, correct the reported error and reload.

### Run and debug on a Wire host

1. Start the [example Wire host]({{< ref "/docs/embedding/wire#host-the-configured-engine" >}}) on the extension host and point `ferret.runtime.configPath` at the daemon guide's `runtimes.json`.
2. Run **Ferret: Select Runtime** and choose `application`. Set `ferret.run.parameters` to `{"name": "Ada"}` and leave `ferret.run.workingDirectory` blank to use the host's default.
3. Save the Wire guide's `RETURN app::greet(@name)` query in a local `.fql` file. Run the current file to get `"Hello, Ada!"`. Set a breakpoint in the file and use **Ferret: Debug Current File** to debug it on the same runtime.

The application supplies the functions, modules, and filesystem policies. The source and breakpoints stay on the extension host. Stopping an editor launch releases its query and connection while the application stays running.

The bundled daemon supports both local and Wire launches. A custom `ferretd` needs API 1.3 or later in major 1 for runtime discovery and API 1.4 or later in major 1 for Wire Run and Debug. Configured or explicitly selected Debug also requires an adapter that advertises runtime selection.

## Configuration

These settings control the daemon and language services.

| Setting | Default | Purpose |
| --- | --- | --- |
| `ferret.server.path` | Empty | Authoritative override for the bundled `ferretd` used by language features, execution, and new debug sessions. An invalid override fails instead of falling back. |
| `ferret.server.args` | `[]` | Advanced arguments appended after the required `lsp` command. These affect only the language server. |
| `ferret.trace.server` | `off` | LSP tracing in the **Ferret** output channel: `off`, `messages`, or `verbose`. |

Changing `ferret.server.path` restarts the language server and execution daemon, invalidating active executions. Existing debug sessions keep their current adapter process. Changing `ferret.server.args` restarts only the language server.

## Troubleshooting

If language features stop responding, run **Ferret: Restart Language Server**. This restarts language features without interrupting the execution daemon or active executions.

Check the two output channels for different information:

- **Ferret** contains daemon selection, language-server lifecycle messages, server errors, and protocol traces.
- **Ferret Execution** contains program results, elapsed time, and execution failures.

Set `ferret.trace.server` to `messages` for protocol-level troubleshooting. Use `verbose` only when payload details are necessary, and review verbose traces for source code or other sensitive data before sharing them.

If `ferret.server.path` is set, confirm that it points to an executable available in the extension host and that the executable can start with the `lsp` command. Clear the setting to return to the bundled daemon. Syntax highlighting remains available while the language server is unavailable.

If a Wire runtime is missing or unavailable, check the runtime name and configuration file, then reload runtimes. If the name appears but a launch cannot connect, check that the application is running and reachable from the extension host.

## Next steps

Use the [CLI]({{< ref "/docs/tools/cli" >}}) when the same `.fql` programs need to run from a terminal, shell script, or CI job.

{{< docs-related tiles="tools-ferretd,embedding-wire,tools-cli,tools-cli-run,tools-cli-debug" >}}
