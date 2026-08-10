---
title: "Introducing the Ferret Module Registry"
subtitle: "Discover, install, and publish Ferret modules"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-08-10"
---

Ferret now has a module Registry.

Until now, modules could be published as Go packages, but there was no shared place to discover them, compare versions, check compatibility, or find their documentation. You needed to know that a module existed and where its source lived before you could use it.

The [Ferret Registry]({{< ref "/registry" >}}) gives modules a common home. It is the public catalog that people browse and that the Ferret CLI reads when searching for or installing modules.

[Barn](https://github.com/MontFerret/barn) is the Git-backed source repository and publication pipeline behind the Registry. Contributors submit registration records to Barn for validation and human review; Barn then generates the public Registry from the accepted records.

The Registry is available now as part of the Ferret v2 alpha. The workflow is intentionally practical: find a module, install a compatible release into a Go application, or publish your own module through a Barn pull request.

## One place for modules

A Registry entry brings together the information needed to understand and install a module:

- its canonical `owner/name` identity
- available versions and their Ferret compatibility
- the FQL namespace exposed at runtime
- the public source repository and exact release commit
- the installable Go package path
- versioned module documentation

These identifiers are related, but they are not interchangeable. For the Archive module, they look like this:

| Concept | Value |
| --- | --- |
| Registry identity | `montferret/archive` |
| FQL namespace | `ARCHIVE` |
| Go package path | `github.com/MontFerret/contrib/modules/archive` |
| Runtime registration | `ferret.WithModules(archive.New())` |

The Registry identity is how users search for and install the module. The namespace is what FQL scripts call. The package path is what a Go application compiles. Keeping these identities explicit lets a module use the names that make sense at each boundary without relying on naming conventions to guess the others.

## Find a module

The Registry is the easiest way to browse the catalog from the web. Each module page shows its description, namespace, versions, compatibility, source, and published documentation.

The CLI provides the same discovery workflow from a terminal. Search by identity or description:

{{< terminal command="true" >}}
ferret mod search archive
{{< /terminal >}}

Then inspect a module before adding it to an application:

{{< terminal command="true" >}}
ferret mod info montferret/archive
{{< /terminal >}}

`mod info` shows the available and selected versions, runtime namespace, Ferret compatibility, repository, commit, and documentation URL. This makes the Registry useful before installation, not only as a package lookup service.

## Install a module into an application

Run `ferret mod install` from the Go application that hosts Ferret:

{{< terminal command="true" >}}
ferret mod install montferret/archive
{{< /terminal >}}

Without an explicit version, the CLI selects the newest registered release compatible with the Ferret version already used by the project. It does not silently move the application to a different Ferret release to make a module fit.

You can also request a specific release:

{{< terminal command="true" >}}
ferret mod install montferret/archive@1.0.0-rc.3
{{< /terminal >}}

An exact release still has to declare compatibility with the project's Ferret version.

The installer resolves the Registry record, adds the published Go package, updates the application’s ferret.New(...) composition, and verifies that the affected package builds before applying the staged changes. If the application is missing a safe Ferret dependency or composition helper, interactive use can show and request approval for those prerequisite changes.

There is one important boundary: this command installs a module into an existing Go host application. It does not extend the standalone Ferret CLI runtime. The module is compiled into the application and runs with that process's permissions; the official CLI continues to ship with its own selected module set.

For project prerequisites, non-interactive setup, and configuration after installation, see [Install a module]({{< ref "/docs/modules/install" >}}).

## Publish a module

Publishing starts from a tagged Go module in a public Git repository. The Registry does not accept uploaded package archives. Instead, Ferret resolves the public tag, inspects the files at its exact commit, and prepares the immutable Barn records that identify the release.

The example below publishes `ziflex/kvplugin` at version `1.0.0`.

### Complete the release metadata

Keep `ferret.yaml`, `go.mod`, and `README.md` together at the module root. A complete standalone manifest looks like this:

```yaml
$schema: https://schemas.ferretlang.org/module/v1.json
name: ziflex/kvplugin
namespace: KV
version: 1.0.0
description: Provides an in-memory key-value cache for Ferret queries.
license: Apache-2.0
documentation: https://github.com/ziflex/ferret-kvplugin#readme
repository:
  url: https://github.com/ziflex/ferret-kvplugin
compatibility:
  ferret: ">=2.0.0-alpha.44 <3.0.0"
keywords:
  - cache
categories:
  - data
```

Module Manifest v1 is strict. Unknown fields, duplicate keys, malformed versions, mixed-case Registry identities, invalid URLs, and invalid SPDX license expressions are rejected.

The adjacent `go.mod` supplies the installable package path. The `README.md` supplies the versioned documentation Barn will snapshot and publish as `Markdown` and sanitized `HTML`. A module `README.md` typically includes `Overview`, `Installation`, `Quick Start`, and `API Reference` sections.

`compatibility.ferret` is optional in the schema, but Registry releases should declare it. The installer needs a usable compatibility range to select a release for a host application.

### Tag and push the release

Commit the manifest, module implementation, tests, `go.mod`, and `README.md` before creating the release tag. For a standalone module, the default tag is `v<version>`:

{{< terminal command="true" >}}
git tag v1.0.0
git push origin v1.0.0
{{< /terminal >}}

For a monorepo module, the default tag is `<repository.directory>/v<version>`. Pass `--tag` to each publication command when a repository uses a different convention. The complete [publication guide]({{< ref "/docs/modules/publish" >}}) covers both layouts.

### Validate the release

After the tag is public, validate the complete release without authenticating or changing anything on GitHub:

{{< terminal command="true" >}}
ferret mod publish --dry-run
{{< /terminal >}}

The command checks the local manifest against the current Registry, resolves the public tag to its exact commit, and verifies the README.md, go.mod, module identity, version, and package path. Nothing is submitted, and no module code is checked out or executed.

### Authenticate with GitHub

Before submitting, set `GH_TOKEN` or `GITHUB_TOKEN`. If neither is set, Ferret reads the current `github.com` token by running:

{{< terminal command="true" >}}
gh auth token --hostname github.com
{{< /terminal >}}

Authenticate the GitHub CLI when needed:

{{< terminal command="true" >}}
gh auth login --hostname github.com
{{< /terminal >}}

The credential must be able to write to your personal Barn fork and open a pull request against `MontFerret/barn`.

### Submit the publication

Run the default command from the module root:

{{< terminal command="true" >}}
ferret mod publish
{{< /terminal >}}

After validation, the CLI uses the GitHub API to create or reuse your personal Barn fork, create a focused publication branch, and open a pull request against `MontFerret/barn`. A module's first release prepares its module and version records; later releases prepare only the new version record. You do not need a local Barn checkout or knowledge of its record layout.

To inspect the records without submitting them, print a deterministic, versioned JSON document:

{{< terminal command="true" >}}
ferret mod publish --print
{{< /terminal >}}

`--dry-run` and `--print` cannot be combined, and neither mode submits records. Retries are safe: an already-published version exits successfully, while an exact open pull request or publication branch is reused. Ferret refuses to overwrite an immutable Registry record or a divergent branch. If a stale publication branch in your fork contains different changes, delete it before retrying.

## Publication model

The Registry's publication process is deliberately Git-backed and review-based.

Specs owns the strict portable contracts for module manifests and Registry records. Barn owns repository inspection, cross-document checks, publication history, and generation of the public catalog. The CLI consumes the resulting versioned artifacts and applies compatibility rules when installing a module.

Before submitting, authenticate with GitHub using `GH_TOKEN`, `GITHUB_TOKEN`, or the GitHub CLI. See the publication guide for authentication and permission details.

## Available during the v2 alpha

The Registry and `ferret mod` workflow are available now with the Ferret v2 alpha CLI. They are useful today, but the wider v2 ecosystem is still moving toward beta, so feedback on module authoring, discovery, compatibility, and publication is especially useful.

You can:

- [browse the Registry]({{< ref "/registry" >}})
- [install a module]({{< ref "/docs/modules/install" >}})
- [initialize a module project]({{< ref "/docs/modules/develop" >}})
- [build a complete module with the Ferret SDK]({{< ref "/docs/guides/writing-plugins" >}})
- [publish a module through Barn]({{< ref "/docs/modules/publish" >}})
- [review the `ferret mod` command reference]({{< ref "/docs/tools/cli/mod" >}})

If you have a module that would be useful outside your own application, the Registry now gives it a path from a public Git tag to something other Ferret users can find, inspect, and install.
