---
title: "Migrate"
sidebarTitle: "Migrate"
weight: 77
draft: false
description: "Migrate supported Ferret v1 Go imports and final-FOR FQL behavior to Ferret v2."
---

# Migrate

The `ferret migrate` command applies the supported mechanical parts of a Ferret v1 to v2 migration. Run it anywhere inside the Go module that contains the embedded application and its FQL files.

Preview the affected paths before changing the project:

{{< terminal command="true" >}}
ferret migrate --dry-run
{{< /terminal >}}

Print a unified diff for review:

{{< terminal command="true" >}}
ferret migrate --print
{{< /terminal >}}

Run without either flag to apply all planned Go and FQL replacements as one transaction:

{{< terminal command="true" >}}
ferret migrate
{{< /terminal >}}

## Migrate a final `FOR`

Ferret v1 returned the value produced by a final top-level `FOR` implicitly. Ferret v2 requires the result to be returned explicitly.

Before:

```fql
FOR item IN 1..3
    RETURN item
```

After:

```fql
return for item in 1..3 {
    return item
}
```

The command changes only a structurally recognized final top-level `FOR` when the program has no explicit terminal `return`. It does not independently wrap nested, assigned, expression-contained, function-contained, non-final, or already-returned loops. Files that need only formatter case or layout changes remain byte-for-byte unchanged.

Changed FQL files are rendered with the canonical formatter. A second migration leaves the explicit result unchanged.

## Go compatibility imports

Documented Ferret v1 imports are rewritten to their Ferret v2 compatibility packages. `go.mod` and `go.sum` are updated only when a Go import is rewritten. A project with only FQL changes does not run `go get` or change dependencies.

Generated Go files and v1 packages without a documented compatibility replacement are left unchanged and reported for manual follow-up.

## Source discovery and parse failures

The command scans lowercase `.fql` files within the containing Go module. It excludes:

- `vendor`, `testdata`, and `node_modules`
- hidden and underscore-prefixed directories
- nested Go modules

Malformed FQL is not modified. The command reports the path, first useful diagnostic, and source line, then continues planning other files. Files that can be migrated are still committed together; a commit failure rolls the transaction back.

Standalone directories without a `go.mod` are not supported.

## Scope

`ferret migrate` is a compatibility aid, not a general v1-to-v2 translator. It does not translate arbitrary v1 APIs or application logic, invent replacements for removed packages, or rewrite excluded directories.

If the project vendors dependencies, run `go mod vendor` after reviewing and applying a migration that changed Go imports.

## Flags

| Flag | Purpose |
| --- | --- |
| `--dry-run` | Show the files that would change without writing them |
| `--print` | Print only a deterministic unified diff on stdout without writing files |

The flags cannot be combined. Diagnostics and manual follow-up remain on stderr when `--print` is used.

## Next steps

{{< docs-related tiles="tools-cli-run,embedding-getting-started,tools-cli-mod" >}}
