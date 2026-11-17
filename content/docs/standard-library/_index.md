---
title: "Standard Library"
weight: 30
type: "docs"
draft: false
description: "Built-in Ferret functions and modules."
stdlibOrderByWeight: true
aliases:
  - /docs/stdlib/
  - /docs/standard-library/html/
  - /docs/stdlib/html/
---

# Standard Library

This reference covers the built-in functions and namespaces supplied by Ferret Core. Embedded applications can select which standard-library groups to register. A function supplied by an optional module, such as `web::html::open`, requires a runtime with that module installed. See [Go engine configuration]({{< ref "/docs/embedding/go/configuration" >}}) and [Modules]({{< ref "/docs/modules" >}}).

Registered namespace segments and host-function names are case-insensitive in FQL. This reference uses their single canonical lowercase qualified spelling.

This reference documents Ferret `{{< data "versions.core.v2" >}}`. It is generated from the published API Reference and its matching API Catalog.

{{< stdlib-children >}}

## Next steps

{{< docs-related tiles="web-extraction,embedding,tools-lab" >}}
