---
title: "Publish a module"
sidebarTitle: "Publish"
weight: 30
draft: false
description: "Prepare and submit a Ferret module release to the Registry through Barn."
---

# Publish a module

Publishing registers a tagged module release in [Barn](https://github.com/MontFerret/barn), the Git-backed source of the Ferret Registry. The CLI validates the public release, prepares its immutable records, and uses the GitHub API to open the Barn pull request. Publication completes after that pull request passes review and is merged.

## Complete the module manifest

Replace the scaffold placeholders with release metadata before tagging the repository. A standalone `ziflex/kvplugin` module can use:

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

Module Manifest v1 is closed: unknown fields, duplicate keys, malformed versions, mixed-case Registry identities, non-HTTPS URLs, and invalid SPDX license expressions are rejected. The `name` must match the module being registered, while `namespace` remains an independent runtime identity.

`compatibility.ferret` is optional in the schema but should be present for Registry releases. `ferret mod install` needs a valid compatibility range to choose a release for a host application.

For a monorepo module, add its normalized repository-relative directory:

```yaml
repository:
  url: https://github.com/ziflex/ferret-modules
  directory: modules/kvplugin
```

Keep `ferret.yaml`, `go.mod`, and `README.md` together at that module root. Barn reads the installable package path from the adjacent `go.mod` directive and checks that it is compatible with the published version.

## Write release documentation

Barn snapshots `README.md` from the exact tagged commit and publishes it as Markdown and sanitized HTML. The README should contain these sections:

- Overview
- Installation
- Quick Start
- API Reference

The manifest's `documentation` URL is the canonical base used to resolve relative links in that README. Keep the manifest description to a concise, single-line summary and put detailed usage in the documentation.

## Tag and push the release

Commit the completed manifest, `go.mod`, implementation, tests, and README before creating the tag. The manifest version must equal the release version.

For a standalone module, the default tag is `v<version>`:

{{< terminal command="true" >}}
git tag v1.0.0
git push origin v1.0.0
{{< /terminal >}}

For the monorepo example above, the default is `modules/kvplugin/v1.0.0`:

{{< terminal command="true" >}}
git tag modules/kvplugin/v1.0.0
git push origin modules/kvplugin/v1.0.0
{{< /terminal >}}

The repository must be public and support anonymous HTTPS Git access. Publication inspects the pushed tag without checking out or executing module code.

For a non-standard release tag, pass the same `--tag` value whenever you validate, submit, or print the release records. For example:

{{< terminal command="true" >}}
ferret mod publish --tag release-1.0.0 --dry-run
{{< /terminal >}}

## Validate without submitting

Run a dry run from the module root after the tag is available remotely:

{{< terminal command="true" >}}
ferret mod publish --dry-run
{{< /terminal >}}

The command:

1. validates the local `ferret.yaml`
2. checks the public Registry for the module and version
3. resolves the pushed tag and pins its exact commit through anonymous HTTPS Git
4. reads and checks `ferret.yaml`, `README.md`, and `go.mod` at that commit
5. verifies the identity, version, package path, documentation, and immutable Barn records needed for the release

`--dry-run` performs the complete release validation without resolving a GitHub credential or changing anything on GitHub. It still reads the public Registry and the module's public Git repository.

## Authenticate with GitHub

The submitting command resolves a GitHub token in this order:

1. `GH_TOKEN`
2. `GITHUB_TOKEN`
3. `gh auth token --hostname github.com`

Set one of the environment variables before publishing, or authenticate the GitHub CLI when needed:

{{< terminal command="true" >}}
gh auth login --hostname github.com
{{< /terminal >}}

The credential must be able to write to your personal Barn fork and open a pull request against the public `MontFerret/barn` repository. Ferret does not print the token or add it to the prepared records.

## Submit the release

Run the default command after validation and authentication:

{{< terminal command="true" >}}
ferret mod publish
{{< /terminal >}}

After preparing the release, the CLI:

1. creates or reuses your personal Barn fork
2. creates a focused publication branch containing only the required source records
3. opens a pull request against `MontFerret/barn`
4. prints the new or existing pull request URL

The first `ziflex/kvplugin` release prepares its module record and its `1.0.0` version record. Later releases prepare only the new version record. You do not need a local Barn checkout, a manual `make check`, or knowledge of the Registry record layout.

The CLI does not upload module code. Barn continues to reference the exact public tag and pinned commit, and the pull request remains subject to Barn's review and CI validation.

## Inspect the prepared records

Use `--print` to emit the deterministic Barn-relative records as a versioned JSON document:

{{< terminal command="true" >}}
ferret mod publish --print
{{< /terminal >}}

This mode is intended for inspection or unusual manual automation. It does not resolve a GitHub credential or submit the records. `--dry-run` and `--print` cannot be combined.

## Retry safely

Publication is safe to retry:

- an already-published version exits successfully without changing GitHub
- an exact open pull request is returned instead of creating a duplicate
- an exact publication branch is reused

Ferret refuses to overwrite an immutable Registry record, reuse a pull request with different records, or replace a publication branch whose base or contents differ. If a stale publication branch in your personal fork conflicts with the prepared release, delete that branch before retrying.

## After merge

Barn assigns `publishedAt` after the version first reaches `main`; contributors and the CLI do not supply that value. Barn then generates and verifies `dist/` in CI and deploys only a fully stamped Registry. Published sources, version identities, records, and assigned timestamps are immutable.

## Next steps

{{< docs-related tiles="registry,runtime-modules-install,runtime-modules-develop" >}}
