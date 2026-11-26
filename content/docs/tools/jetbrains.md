---
title: "JetBrains IDEs"
weight: 45
draft: false
description: "Write, format, run, and debug FQL with the official Ferret plugin for IntelliJ IDEA and GoLand."
---

# JetBrains IDEs

The official Ferret plugin for IntelliJ IDEA and GoLand provides language intelligence, formatting, and native Run and Debug actions for `.fql` files. It uses the bundled Ferret daemon to run queries locally or against an application-hosted runtime.

## Installation

The v0.1.0 plugin supports IntelliJ IDEA and GoLand on the **2026.2** and **2026.3** release lines.

1. From the [JetBrains v0.1.0 release](https://github.com/MontFerret/editorium/releases/tag/jetbrains/v0.1.0), download `ferret-jetbrains-0.1.0-signed.zip` and its `.sha256` file. Keep both files in the same directory.
2. Verify the ZIP's SHA-256 against the checksum file. On macOS or Linux, run `shasum -a 256 -c ferret-jetbrains-0.1.0-signed.zip.sha256`. In PowerShell, use `Get-FileHash ferret-jetbrains-0.1.0-signed.zip -Algorithm SHA256` and compare the hash with the published value.
3. In **Settings → Plugins**, open the gear menu and choose **Install Plugin from Disk…**. Select the ZIP and restart the IDE if prompted.

## Bundled Ferret daemon

The universal plugin archive bundles `ferretd` 1.0.0-alpha.13 for macOS, Linux, and Windows on ARM64 and x64. The plugin selects the daemon for the IDE host automatically.

## Language intelligence

Open a local `.fql` file. The IDE recognizes it as a Ferret file and starts language services when needed. The bundled daemon provides diagnostics, completion, hover information, and definition navigation within the current document.

## Formatting

With the `.fql` file open and the selection cleared, choose **Code → Reformat Code**. The action uses Ferret's formatter on the whole document. Fix syntax errors before reformatting.

## Running programs

Open a local `.fql` file and choose **Run** from its context menu or select **Current File** in the Run/Debug configuration selector. For a first run, save `RETURN 42` in the file. The **Run** tool window shows the selected runtime, the result, and any errors.

### Create a Run/Debug configuration

Choose **Run → Edit Configurations… → Add New Configuration → Ferret**. Use a saved configuration to reuse the same source, runtime, parameters, and working directory for Run and Debug.

| Field | What to enter |
| --- | --- |
| **Runtime** | **Daemon default** or the exact runtime name. Choose `local` to use the built-in runtime even when the configured default is Wire. |
| **Source file** | A local `.fql` file inside the project. You can use an absolute path or a path relative to the project directory. |
| **Working directory** | Leave blank to use the runtime's default. For local execution, a relative path resolves from the project directory. For Wire, the application host receives the path as entered. |
| **Parameters (JSON object)** | A JSON object such as `{"limit": 10}`. Keys omit the `@` prefix. A blank field supplies an empty object. |

With that parameter value, save the following as your source and run the configuration to get `11`.

```fql
LET total = @limit + 1
RETURN total
```

## Debugging

Set a line breakpoint in the gutter and choose **Debug** for the current file or a saved Ferret configuration. For the example above, break on `RETURN total`. The **Debug** tool window shows the call stack. Select a frame to inspect **Locals** and **Parameters**, then use **Evaluate Expression** or add a watch such as `total + 1`.

Use **Resume**, **Pause**, **Step Over**, **Step Into**, **Step Out**, and **Stop** from the standard debugger controls. Debug uses the saved source at launch. Start a new session to debug subsequent source changes.

Evaluation supports reading bindings and parameters, member and index access, and scalar operators. Function calls and mutation are unavailable in the evaluator. Large collections can appear as summaries with no expandable children. To watch a nested value, enter its full member or index expression.

The plugin supports ordinary line breakpoints. Conditional breakpoints, logpoints, attaching to existing processes, and interactive console input are unavailable.

## Configure runtimes

Open **Settings → Languages & Frameworks → Ferret** and set **Runtime configuration file** to a [daemon runtime configuration]({{< ref "/docs/tools/ferretd#configure-a-runtime" >}}). Relative paths resolve from the project directory. A blank setting uses the built-in configuration.

Select **Apply**, then expand **Runtime catalog** or open the **Runtime** picker in a Ferret Run/Debug configuration. Discovery shows each configured name, label, and provider, even if its Wire host is offline. The picker follows the daemon default or saves a case-sensitive name for that configuration.

After the configuration is first used, choose **Reload runtimes** when you edit the file or change or clear the path. Finish or stop active Run and Debug sessions and wait for them to close before reloading. Use **Retry discovery** if the catalog request failed and the daemon is still running. If daemon startup failed, correct the reported error and reload.

### Run and debug on a Wire host

1. Start the [example Wire host]({{< ref "/docs/embedding/wire#host-the-configured-engine" >}}) and select the daemon guide's `runtimes.json` in the Ferret settings.
2. Save the Wire guide's `RETURN app::greet(@name)` query in a local `.fql` file. Create a **Ferret** configuration with **Runtime** set to `application` and **Parameters (JSON object)** set to `{"name": "Ada"}`. Leave **Working directory** blank to use the host's default.
3. Run the configuration to get `"Hello, Ada!"`. Add a line breakpoint in the local file and Debug the same configuration.

Each launch uses the application's functions, modules, and filesystem policies. The debugger uses your local source for breakpoints and stepping. Stopping the launch releases its query and connection. The application stays running.

## Troubleshooting

- If the plugin fails to load, check the IDE version and enable **Ferret** under **Settings → Plugins → Installed**. For a daemon startup error, check the IDE log and reinstall the signed archive if the bundled executable is missing.
- For a Run or Debug failure, check the corresponding tool window. Confirm the source file is saved, the parameters are a JSON object, and the selected runtime is in the catalog.
- If a Wire launch cannot connect, check the endpoint and confirm the host is running. After an endpoint change, update the configuration file, reload runtimes, and start a new launch.

## Next steps

{{< docs-related tiles="tools-ferretd,embedding-wire,tools-cli" >}}
