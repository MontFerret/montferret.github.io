---
title: "Mod"
sidebarTitle: "Mod"
weight: 75
draft: false
description: "Discover, install, initialize, and publish Ferret modules from the CLI."
---

# Mod

The `ferret mod` command groups Registry discovery, project installation, module initialization, and publication.

```text
ferret mod search [query]
ferret mod info <name>
ferret mod install <module>[@version]
ferret mod init [name]
ferret mod publish [--tag <tag>] [--dry-run | --print]
```

## Search and inspect

Search module identities and descriptions:

{{< terminal command="true" >}}
ferret mod search sqlite
{{< /terminal >}}

Show one module's available versions, selected version, namespace, source, compatibility, and documentation:

{{< terminal command="true" >}}
ferret mod info montferret/sqlite
{{< /terminal >}}

## Install

Install the newest compatible release or request an exact version:

{{< terminal command="true" >}}
ferret mod install montferret/archive
ferret mod install montferret/archive@1.0.0-rc.3
{{< /terminal >}}

| Flag | Purpose |
| --- | --- |
| `-y`, `--yes` | Add safe missing Ferret and composition prerequisites without prompting |

The command installs into the current Go application, not the Ferret CLI runtime. See [Install a module]({{< ref "/docs/modules/install" >}}) for project discovery, compatibility, composition changes, and validation behavior.

## Initialize

Start the guided flow:

{{< terminal command="true" >}}
ferret mod init
{{< /terminal >}}

For automation, provide at least the Registry identity and Go module path:

{{< terminal command="true" >}}
ferret mod init acme/kvplugin --go-module github.com/acme/ferret-kvplugin
{{< /terminal >}}

| Flag | Purpose |
| --- | --- |
| `--go-module` | Go import path written to the generated `go.mod` |
| `--dir` | Destination directory; defaults to the module-name leaf |
| `--namespace` | Runtime namespace; defaults to the module-name leaf |

See [Develop a module]({{< ref "/docs/modules/develop" >}}) for the generated project and next steps.

## Publish

Validate and submit the current tagged release through Barn:

{{< terminal command="true" >}}
ferret mod publish
{{< /terminal >}}

| Flag | Purpose |
| --- | --- |
| `--tag` | Override the default standalone or monorepo release tag |
| `--dry-run` | Validate and prepare the release without authenticating or submitting to GitHub |
| `--print` | Print the prepared Barn records as versioned JSON without submitting |

Without a non-submitting flag, the command creates or reuses a personal Barn fork and focused publication branch, then opens or reuses a pull request against `MontFerret/barn`. It reads `GH_TOKEN`, then `GITHUB_TOKEN`, and otherwise uses the current `gh` CLI credential. `--dry-run` and `--print` cannot be combined; neither mode resolves a GitHub credential or submits records.

The command does not upload module code. Publication remains tied to the public tag and pinned commit, and an already-published version is a successful no-op. See [Publish a module]({{< ref "/docs/modules/publish" >}}) for authentication, validation, and retry behavior.

## Next steps

{{< docs-related tiles="runtime-modules-install,runtime-modules-develop,runtime-modules-publish,registry" >}}
