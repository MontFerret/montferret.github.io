---
title: "Browser"
weight: 70
draft: false
description: "Start and stop managed browser instances for scripts that use browser-backed features."
---

# Browser

The `browser` command manages Chrome or Chromium instances for scripts that need a browser through the Chrome DevTools Protocol.

Most users do not need this command directly. When you pass `--browser-open` or `--browser-headless` to `ferret run`, the CLI starts and stops a browser automatically. The `browser` command is useful when you want a long-lived browser that multiple script runs can share.

## Open a browser

{{< terminal >}}
ferret browser open
{{< /terminal >}}

This starts a browser and prints the process ID. The CLI waits until the browser exits.

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-d`, `--detach` | `false` | Start in background and print process ID |
| `--headless` | `false` | Start in headless mode |
| `-p`, `--port` | `9222` | Remote debugging port |
| `--user-dir` | `.ferret-browser` | Browser user data directory |

To start a headless browser in the background:

{{< terminal >}}
ferret browser open --detach --headless
{{< /terminal >}}

Then run scripts against it:

{{< terminal >}}
ferret run --browser-address http://127.0.0.1:9222 script.fql
{{< /terminal >}}

## Close a browser

Stop a browser by process ID:

{{< terminal >}}
ferret browser close 12345
{{< /terminal >}}

Or close the browser associated with the current configuration:

{{< terminal >}}
ferret browser close
{{< /terminal >}}

## Next steps

{{< docs-related tiles="tools-cli-run,web-extraction,tools-lab" >}}
