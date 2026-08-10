---
title: "Modules"
weight: 60
draft: false
description: "Discover, install, develop, and publish Ferret modules."
---

# Modules

Modules add functions, host values, data formats, and integrations to a Ferret runtime. Public modules are distributed as Go packages and can be discovered through the [Ferret Registry]({{< ref "/registry" >}}).

The CLI covers the complete module workflow:

- install a registered module into a Go application
- initialize a module project for development
- publish a tagged release through Barn

## How module identities fit together

A module has several related identifiers. They are not interchangeable.

| Concept | Example | Used by |
| --- | --- | --- |
| Registry identity | `ziflex/kvplugin` | `ferret.yaml`, the Registry, and `ferret mod` commands |
| Runtime namespace | `KV` | FQL functions such as `KV::GET` |
| Go package path | `github.com/ziflex/ferret-kvplugin` | Go imports and dependency resolution |
| Runtime registration | `ferret.WithModules(kvplugin.New())` | The host application's Ferret engine |

The Registry identity is a canonical lowercase `owner/name` value. The runtime namespace is an independent, case-sensitive FQL namespace. Barn obtains the installable Go package path from the `go.mod` beside the module manifest instead of deriving it from either identity.

## Modules run in a host application

`ferret mod install` installs a module into an existing Go application. The command adds the Go dependency and registers the module in that application's `ferret.New(...)` composition.

It does not extend or modify the standalone Ferret CLI runtime. Modules installed this way are compiled into the host application and run with that process's permissions. The official CLI includes its own selected module set.

## Choose a workflow

{{< docs-related tiles="runtime-modules-install,runtime-modules-develop,runtime-modules-publish,tools-cli-mod" >}}
