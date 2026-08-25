---
title: "Ferret Comes to Visual Studio Code"
subtitle: "Write, run, format, and debug FQL without leaving your editor."
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-08-25"
---

Hello, friends!

I'm excited to finally share something we’ve been working toward for a while: the new official Ferret extension for Visual Studio Code is now available on the Visual Studio Marketplace.

`Ferret for Visual Studio Code` brings the Ferret development experience directly into the editor, backed by Ferret Daemon (`ferretd`) and the same language tooling we’ve been building for Ferret v2.

### More than syntax highlighting

The old FQL extension provided the basics: syntax highlighting and language configuration. The new extension is a rather different beast.

It provides a proper IDE experience for Ferret, including:

- Syntax highlighting for FQL
- Language intelligence powered by the Ferret language server
- Diagnostics directly in the editor
- Document formatting, including format-on-save
- Run Current File without leaving VS Code
- Integrated debugging
- Breakpoints in `.fql` files
- Execution cancellation and output
- Support for macOS, Linux, and Windows, on both x64 and ARM64

And, importantly, you don't need to install or configure the Ferret daemon separately.

The extension ships with the appropriate `ferretd` binary for your platform and takes care of starting and managing it automatically. Open a `.fql` file and get to work.

### Run and debug Ferret like any other language

This is probably my favorite part of the release.

Ferret programs no longer have to feel like scripts you edit in one place and execute somewhere else. You can run the current file directly from the editor, see its output, set breakpoints, start a debugging session, and inspect execution using VS Code's native debugging experience.

That may sound like ordinary IDE functionality - that's precisely the point. 
Ferret is a language and it should behave like one in your editor.

### Built on `ferretd`

The extension itself deliberately doesn't try to implement Ferret's language semantics.

Language intelligence, formatting, execution, and debugging are provided through Ferret Daemon (`ferretd`). This gives us a common foundation that isn't tied to VS Code and can be reused by other editors and development tools.

VS Code is the first complete integration built on top of that architecture, but it isn't intended to be the last.

### This is v0.1.0

This is the first public release, so the extension is currently marked as Preview.

There will undoubtedly be rough edges, and the IDE experience will continue evolving alongside Ferret v2 and `ferretd`. But the important pieces are now in place: editing, language intelligence, formatting, execution, and debugging all work together as one development experience.

If you find something broken - or simply something that could be better - please open an issue. Feedback from actually using Ferret inside an IDE is going to be particularly valuable at this stage.

### Try it

The extension is available now on the Visual Studio Marketplace:

Ferret — Visual Studio Marketplace  
https://marketplace.visualstudio.com/items?itemName=ferretlang.ferret-lang

Or search for Ferret directly from the Extensions view in Visual Studio Code.

Install it, open an `.fql` file, and give it a spin.

Happy ferreting! 🦦