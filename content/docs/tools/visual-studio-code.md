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

Install [Ferret Lang from the Visual Studio Marketplace](https://marketplace.visualstudio.com/items?itemName=ferretlang.ferret-lang). The official extension is published by `ferretlang` and requires Visual Studio Code 1.95 or later.

You can also install it from Visual Studio Code:

1. Open the **Extensions** view.
2. Search for **Ferret Lang**.
3. Select the extension published by `ferretlang`.
4. Select **Install**.

Open or create a file ending in `.fql`. Visual Studio Code registers it as a Ferret document automatically.

## Bundled Ferret daemon

The extension includes the appropriate `ferretd` executable for macOS, Linux, and Windows on ARM64 and x64. You do not normally need to install the daemon separately or add it to `PATH`.

VS Code installs the package for the extension host. In Remote SSH, Dev Containers, WSL, and Codespaces, use the package for the remote host platform.

## Language intelligence

The extension provides syntax highlighting and standard editing support for `.fql` files. For saved, file-backed documents, the bundled daemon adds:

- diagnostics in the Problems view and editor
- completion suggestions
- hover information and signature help
- document symbols and outline navigation
- go to definition and find references
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

## Debugging

Open a saved `.fql` file and select the Debug button, run **Ferret: Debug Current File**, or press `F5`. The current file becomes the program, and the containing workspace folder becomes its working directory. A standalone file uses its own directory, so the normal workflow does not require `launch.json`.

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
| `cwd` | No | Working directory. When omitted, the daemon uses the program's directory; zero-configuration debugging prefers the containing workspace folder. |
| `parameters` | No | Object containing values available to FQL as bind parameters. The default is no parameters. |
| `stopOnEntry` | No | Pause before normal execution begins. The default is `false`. |

## Configuration

Normal use does not require extension settings. These public settings are intended for advanced configuration and troubleshooting:

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

## Next steps

Use the [CLI]({{< ref "/docs/tools/cli" >}}) when the same `.fql` programs need to run from a terminal, shell script, or CI job.

{{< docs-related tiles="tools-cli,tools-cli-run,tools-cli-debug" >}}
