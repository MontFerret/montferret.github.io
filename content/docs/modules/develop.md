---
title: "Develop a module"
sidebarTitle: "Develop"
weight: 20
draft: false
description: "Initialize a Ferret module project and prepare it for implementation."
---

# Develop a module

Use `ferret mod init` to create a standalone Go project with a Ferret module manifest, starter registration code, and a small package layout. The command creates the project locally without downloading dependencies.

## Use the guided initializer

Run the command without arguments to start the interactive flow:

{{< terminal command="true" >}}
ferret mod init
{{< /terminal >}}

The initializer explains and prompts for:

| Value | Purpose |
| --- | --- |
| Module name | Canonical lowercase Registry identity in `owner/name` form |
| Go module path | Import path written to `go.mod` and used by Go consumers |
| Directory | New destination directory for the project |
| Namespace | Case-sensitive FQL namespace exposed at runtime |

It shows the resolved configuration before creating any files. Canceling the confirmation leaves the destination untouched.

## Initialize non-interactively

Provide the Registry identity and Go module path when the command cannot prompt:

{{< terminal command="true" >}}
ferret mod init acme/kvplugin \
  --go-module github.com/acme/ferret-kvplugin \
  --dir kvplugin \
  --namespace kv
{{< /terminal >}}

`--dir` and `--namespace` are optional. When omitted, both default from the `kvplugin` leaf of the module name. The destination must not already exist.

## Review the generated project

The scaffold contains:

```text
kvplugin/
├── ferret.yaml
├── go.mod
├── module.go
├── README.md
├── core/
│   └── doc.go
└── lib/
    └── doc.go
```

- `ferret.yaml` contains a valid Module Manifest v1 with TODO metadata.
- `go.mod` pins the Go and Ferret versions embedded in the CLI release.
- `module.go` returns a module named `acme/kvplugin` and provides the registration callback.
- `core` is the starting point for implementation code.
- `lib` is the starting point for Ferret-facing function bindings.
- `README.md` is the documentation Barn will snapshot from a published tag.

The generated manifest is schema-valid, but its placeholder description, license, and documentation URL are not release metadata. Replace every TODO before publishing.

## Keep distribution and runtime names separate

The generated manifest records both identities explicitly:

```yaml
name: acme/kvplugin
namespace: kv
```

`name` identifies the module in the Registry and dependency metadata. `namespace` identifies its Ferret-facing API, such as `kv::open`. Changing one does not change the other.

## Implement and test the module

Start by resolving the pinned dependencies:

{{< terminal command="true" >}}
cd kvplugin
go mod tidy
{{< /terminal >}}

Then register the module's functions, values, codecs, or lifecycle hooks in `module.go` and keep implementation details in the appropriate packages. Run the Go test suite before preparing a release:

{{< terminal command="true" >}}
go test ./...
{{< /terminal >}}

For a complete module implementation, see [Develop a Ferret module]({{< ref "/docs/guides/writing-plugins" >}}).

## Next steps

{{< docs-related tiles="guide-writing-plugins,runtime-modules-publish,embedding-go-modules" >}}
