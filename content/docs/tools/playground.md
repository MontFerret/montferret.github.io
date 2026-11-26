---
title: "Playground"
weight: 50
draft: false
description: "Run, format, and inspect FQL in the hosted Ferret Playground."
---

# Playground

Open [Try Ferret]({{< ref "/try" >}}), select **Basic evaluation** under **Examples**, and select **Run**. The query `RETURN 1 + 2` produces `3` in the **Result** tab. You can also run from the editor with Ctrl+Enter or Cmd+Enter.

The editor checks your source as you type. Select a diagnostic with a position to jump to it. Use **Format** to apply the runtime's formatter or **Bytecode** to inspect the compiled query. Both operations work on source alone.

Open **Params** to provide bind values. The editor reads valid JSON such as `false`, `null`, or `{"tags":["sale"]}` as typed values. It uses plain text when a value is invalid JSON. Enter `"false"` to pass it as a string.

**Share** copies a link containing the current source and parameter values. Opening a link fills the editor. Select **Run** when you're ready to execute it.

The hosted runtime uses Ferret's safe standard library and provides `DOCUMENT(url)` for supported pages on Mockery. Try **Mockery product extraction** to see how the browser can wait for JavaScript-rendered content. The toolbar reports the runtime version. The [Playground API contract](https://github.com/MontFerret/playground/blob/v1.0.0-alpha.5/api/openapi.yaml) describes discovery, tooling, and execution for applications that use the service.
