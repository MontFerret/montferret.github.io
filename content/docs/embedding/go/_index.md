---
title: "Go"
weight: 10
draft: false
description: "Embed the native Ferret runtime into a Go application."
---

# Embed Ferret in Go

Use the native Go library when your application needs full control over Ferret's engine, modules, host values, codecs, sandboxed services, and execution lifecycle.

Start with a minimal query, then move into the API area that matches your integration.

For integrations using portable runtime interfaces, see the [Universal API adapter](universal-api/).

For an existing v1 application, follow [Go embedding migration]({{< ref "docs/migrations/v1-to-v2/go-embedding" >}}).

{{< docs-related tiles="embedding-go-getting-started,embedding-go-executing,embedding-go-parameters,embedding-go-safe-fql-construction,embedding-go-custom-functions,embedding-go-configuration,embedding-go-host-values,embedding-go-modules,embedding-go-value-encoders,embedding-go-programs,embedding-go-compiler-analysis,embedding-go-application-guide" >}}
