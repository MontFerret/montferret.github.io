---
title: "Migrate"
sidebarTitle: "Migrate"
weight: 77
draft: false
description: "Check FQL compatibility and migrate supported Ferret v1 behavior to Ferret v2."
---

# Migrate

The `ferret migrate` command groups the compatibility checker and the supported mechanical Ferret v1 to v2 migration. Running `ferret migrate` without a subcommand displays help and does not modify files.

## Check compatibility

Use `check` to inspect a standalone lowercase `.fql` file or recursively scan a directory without modifying source:

{{< terminal command="true" >}}
ferret migrate check --from v1 .
{{< /terminal >}}

{{< terminal command="true" >}}
ferret migrate check scripts/query.fql
{{< /terminal >}}

The path defaults to the current directory. `--from` currently defaults to and accepts only `v1`. The check does not require a Go module, run Go tooling, resolve dependencies, or format source.

Compatibility findings and malformed FQL are reported with source locations and make the command exit nonzero after all readable files have been checked. Filesystem, cancellation, and internal failures stop the check immediately.

Directory checks include lowercase `.fql` files in `testdata`, hidden and underscore-prefixed directories, and nested Go modules. They skip `.git`, `.hg`, `.svn`, `vendor`, and `node_modules`, and do not follow directory symlinks.

## Run a migration

Use `run` with a standalone lowercase `.fql` file or a project directory. The path defaults to the current directory.

Preview the affected paths before changing the project:

{{< terminal command="true" >}}
ferret migrate run --dry-run
{{< /terminal >}}

Print a unified diff for review:

{{< terminal command="true" >}}
ferret migrate run --print path/to/project
{{< /terminal >}}

Run without either flag to apply all planned Go and FQL replacements as one transaction:

{{< terminal command="true" >}}
ferret migrate run path/to/project
{{< /terminal >}}

A standalone file migration changes only that file. A selected directory is always the migration boundary, including when the path is omitted and defaults to `.`. A directory with no eligible Go source is migrated as an FQL-only project even when it is inside a Go module. If selected Go source exists without a containing `go.mod`, the command fails before applying any FQL changes.

When selected Go source belongs to a containing module, `run` uses that module for Go metadata and dependency updates without scanning source outside the selected directory.

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

Documented Ferret v1 imports are rewritten to their Ferret v2 compatibility packages. `go.mod` and `go.sum` are updated only when a Go import is rewritten. FQL-only targets do not require the Go toolchain and do not change Go dependencies.

Generated Go files and v1 packages without a documented compatibility replacement are left unchanged and reported for manual follow-up.

When the selected directory covers only part of a Go module, the existing Ferret v1 dependency is retained because source outside the migration boundary is not inspected. Run the migration from the module root when the whole module is ready to remove that dependency.

## Source discovery and parse failures

Directory migration scans lowercase `.fql` files within the selected directory. Unlike the broader read-only check, `run` excludes the following descendants:

- `vendor`, `testdata`, and `node_modules`
- hidden and underscore-prefixed directories
- nested Go modules

Malformed FQL is not modified. The command reports the path, first useful diagnostic, and source line, then continues planning other files. Files that can be migrated are still committed together; a commit failure rolls the transaction back.

The selected directory itself is scanned even when its name would be excluded as a descendant, such as `.tmp` or `testdata`. Explicit symlink targets are rejected, and directory symlinks are not followed. A standalone `.fql` target is migrated directly even when it is located under a directory that recursive migration would exclude.

## Scope

`ferret migrate run` is a compatibility aid, not a general v1-to-v2 translator. It does not translate arbitrary v1 APIs or application logic, invent replacements for removed packages, or rewrite source in excluded descendant directories.

If the project vendors dependencies, run `go mod vendor` after reviewing and applying a migration that changed Go imports.

## Commands and flags

| Command or flag | Purpose |
| --- | --- |
| `check [path]` | Check FQL source without modifying files |
| `check --from v1` | Select the source Ferret version; currently only `v1` is supported |
| `run [path]` | Apply the supported migration to a file or selected directory |
| `run --dry-run` | Show the files that would change without writing them |
| `run --print` | Print only a deterministic unified diff on stdout without writing files |

The `run` flags cannot be combined. Diagnostics and manual follow-up remain on stderr when `--print` is used.

## Next steps

{{< docs-related tiles="tools-cli-run,embedding-go-getting-started,tools-cli-mod" >}}
