---
title: "Ferret v0.17.0"
subtitle: "Iterative improvements"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2023-02-07"
---
Hello Ferret users!

We are happy to announce the release of Ferret v0.17.0.

This release focuses on housekeeping and tooling improvements, tightening up the build and modernizing the dependency stack.

# Changes

- Updated project dependencies (#769)
- Dropped support for pre-generics versions of Go
- Updated ANTLR and removed the legacy SDK
- Added staticcheck and goimports to the toolchain
- Updated build steps to reflect the new setup
- Fixed unit tests for HTTP helper functions

While there are no user-facing features in this release, 0.17.0 lays important groundwork for future development by improving correctness, maintainability, and build consistency.
