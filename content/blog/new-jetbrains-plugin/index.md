---
title: "Ferret Comes to JetBrains IDEs"
subtitle: "Native FQL development in IntelliJ IDEA and GoLand."
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-11-26"
---

Hello, friends!

A few months ago, when I introduced the official Visual Studio Code extension, I mentioned that it wouldn't be the last editor integration for Ferret.

Today, I'm happy to share that the first release of the official Ferret plugin for JetBrains IDEs is available!

The plugin brings FQL development directly into IntelliJ IDEA and GoLand, with support for:

- Language intelligence and diagnostics
- Code completion and formatting
- Native Run and Debug configurations
- Breakpoints and stepping
- Variable inspection, expression evaluation, and watches

Just like the Visual Studio Code extension, it comes with a bundled Ferret daemon (`ferretd`), so there's no need to install it separately.

## Beyond local execution

One feature I'm particularly excited about is the ability to run and debug FQL programs against runtimes hosted by other applications through Wire.

If you're embedding Ferret in your own application, you can now develop and debug FQL programs directly from your IDE while using the functions, modules, and capabilities provided by your application.

This functionality is also available in the latest release of the Visual Studio Code extension.

It's nice to see the ideas behind Ferret's Universal API and Wire coming together in something that developers can actually use.

## Give it a try

This is the initial v0.1.0 release, with support for IntelliJ IDEA and GoLand 2026.2 and 2026.3.

You can download the plugin from [GitHub](https://github.com/MontFerret/editorium/releases/tag/jetbrains/v0.1.0) and find installation and usage instructions in the [documentation]({{< ref "/docs/tools/jetbrains" >}}).

As always, feedback, suggestions, and bug reports are very welcome.

Happy ferreting! 🦦
