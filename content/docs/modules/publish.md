---
title: "Publish a module"
sidebarTitle: "Publish"
weight: 30
draft: false
description: "Prepare and submit a Ferret module release to the Registry through Barn."
---

# Publish a module

Publishing registers a tagged module release in [Barn](https://github.com/MontFerret/barn), the Git-backed source of the Ferret Registry. The CLI validates the release and prints the canonical records; publication completes when a human-reviewed Barn pull request is merged.

## Complete the module manifest

Replace the scaffold placeholders with release metadata before tagging the repository. A standalone `acme/kvplugin` module can use:

```yaml
$schema: https://schemas.ferretlang.org/module/v1.json
name: acme/kvplugin
namespace: KV
version: 1.0.0
description: Provides an in-memory key-value cache for Ferret queries.
license: Apache-2.0
documentation: https://github.com/acme/ferret-kvplugin#readme
repository:
  url: https://github.com/acme/ferret-kvplugin
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
  url: https://github.com/acme/ferret-modules
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

The repository must be public and support anonymous HTTPS Git access. Barn inspects the pushed tag without checking out or executing module code.

## Prepare the Registry records

Run the command from the module root after the tag is available remotely:

{{< terminal command="true" >}}
ferret mod publish
{{< /terminal >}}

Use `--tag` only when the release follows a different tag convention:

{{< terminal command="true" >}}
ferret mod publish --tag release-1.0.0
{{< /terminal >}}

The command:

1. validates the local `ferret.yaml`
2. checks the public Registry for the module and version
3. resolves the pushed tag through anonymous HTTPS Git
4. reads and checks the manifest, README, `go.mod`, identity, version, and commit at that tag
5. prints the deterministic Barn-relative records and pull-request guidance

It does not write to either repository, upload a package, authenticate with GitHub or another provider, or open a pull request.

## Submit the Barn pull request

For the first `acme/kvplugin` release, the output contains two records:

```text
registry/modules/acme/kvplugin/manifest.json
registry/modules/acme/kvplugin/versions/v1.0.0.json
```

Add those exact paths and contents to a Barn branch. For later releases, add only the new version record. Do not edit the module manifest or any earlier version record.

From the Barn checkout, validate the registration:

{{< terminal command="true" >}}
make check
{{< /terminal >}}

Commit only the records under `registry/modules/` and open a pull request. Do not add `publishedAt`; Barn assigns it after the version first reaches `main`. Do not create or commit `dist/`; CI generates and verifies the public Registry distribution.

After merge, Barn stamps the canonical version record and deploys only a fully stamped Registry. Published sources, version identities, and assigned timestamps are immutable.

## Next steps

{{< docs-related tiles="registry,runtime-modules-install,runtime-modules-develop" >}}
