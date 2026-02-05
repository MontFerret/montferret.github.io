---
title: "Ferret v0.18.1"
subtitle: "v0.18.1"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2025-05-07"
---
Hello Ferret users!

This is a small cleanup-focused release that modernizes the codebase and keeps things aligned with current Go best practices.

# What’s changed
- Simplified formatting code by using fmt.Printf instead of nested fmt.Println(fmt.Sprintf(...)) (#789)
- Updated ANTLR and other dependencies, and upgraded the Go version (#796)
- Removed deprecated io/ioutil usage as part of an internal refactor (#792)

No functional changes in behavior, but the internals are now cleaner, more idiomatic, and future-proofed. 

As always, thanks to everyone contributing and reporting issues.