---
title: "Migrate Ferret"
sidebarTitle: "Migrate"
weight: 80
draft: false
description: "Choose a migration path for an existing Ferret project."
---

# Migrate Ferret

Choose the version your project uses, then follow its migration guide. For an
existing Ferret v1 project, start with **Ferret v1 → v2**.

For a Go application using the earlier buffered output API, follow [Consumable output](consumable-output/). This development API migration covers Native and UAPI handles, error timing, and ownership.

For a custom codec or direct encoding calls using the earlier byte-slice methods, follow [Encoding API](encoding-api/). This migration covers contexts, reader/writer ownership, byte helpers, and hooks.

For an earlier Lab v2 runtime configuration, follow [Lab runtime configuration](lab-runtime-configuration/).

The CLI can check FQL compatibility with `ferret migrate check` and apply
supported source changes with `ferret migrate run`. Each guide separates
automated changes from work that needs manual review.

{{< docs-related tiles="migrations-v1-to-v2,migrations-consumable-output,migrations-encoding-api,tools-cli-migrate" >}}
