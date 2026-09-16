---
title: "Migrate"
sidebarTitle: "Migrate"
weight: 77
draft: false
description: "Check FQL compatibility and migrate supported Ferret v1 behavior to Ferret v2."
---

# Migrate

The `ferret migrate` command groups the compatibility checker and the supported mechanical Ferret v1 to v2 migration. Running `ferret migrate` without a subcommand displays help and does not modify files.

For the complete project workflow, start with [Ferret v1 → v2]({{< ref "docs/migrations/v1-to-v2" >}}).

## Check compatibility

Use `check` to inspect a standalone lowercase `.fql` file or recursively scan a directory without modifying source:

{{< terminal command="true" >}}
ferret migrate check --from v1 .
{{< /terminal >}}

{{< terminal command="true" >}}
ferret migrate check scripts/query.fql
{{< /terminal >}}

The path defaults to the current directory. `--from` currently defaults to and accepts only `v1`. The check does not require a Go module, run Go tooling, resolve dependencies, or format source.

The check covers final collecting `FOR` compatibility and the same legacy stdlib
calls as `run`, including automatic replacement suggestions and calls requiring
manual review. It uses the same conservative declaration and alias guards.
A clean result covers these supported rules; it does not establish that every
v1 API or application behavior is compatible.

For example, save this legacy query as `scripts/query.fql`:

```fql
return has({ foo: "bar" }, "baz")
```

Running `ferret migrate check scripts/query.fql` reports:

```text
scripts/query.fql:1:8: Legacy stdlib call `has` should use `object::has_key`.
  help: Preview automatic replacements with `ferret migrate run --print`.

Found 1 v1 compatibility issue in 1 of 1 FQL file.
```

Replacement suggestions, manual-review findings, and malformed FQL are reported
on stderr with source locations and make the command exit nonzero after all
readable files have been checked. Filesystem, cancellation, and internal failures
stop the check immediately. The check does not attempt formatting or verify
whether a rewrite can preserve the source; use `run --print` to preview edits.

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

For loop migration, the command changes only a structurally recognized final top-level `FOR` when the program has no explicit terminal `return`. It does not independently wrap nested, assigned, expression-contained, function-contained, non-final, or already-returned loops. Files that need only formatter case or layout changes remain byte-for-byte unchanged.

Changed FQL files are rendered with the canonical formatter. A second migration leaves the explicit result unchanged.

## Migrate standard-library calls

Use a CLI release containing these stdlib migrations and a runtime release
containing the corresponding canonical namespaces.

The command recognizes legacy unqualified function calls case-insensitively and
replaces their targets with canonical lowercase names. For example, the `has`
query above becomes:

```fql
return object::has_key({ foo: "bar" }, "baz")
```

The table lists every supported automatic mapping. Unless a row describes a
rename, the function name stays the same under the listed namespace.

| Namespace | Legacy calls migrated automatically |
| --- | --- |
| `encoding::` | `json_parse`, `json_stringify`, `encode_uri_component` → `query_escape`, `decode_uri_component` → `query_unescape`, `to_base64` → `base64_encode`, `from_base64` → `base64_decode`, `escape_html` → `html_escape`, `unescape_html` → `html_unescape` |
| `crypto::` | `md5`, `sha1`, `sha512`, `random_token` |
| `path::` | `base`, `clean`, `dir`, `ext`, `is_abs`, `separate`, `match` |
| `arrays::` | `first`, `flatten`, `last`, `sorted`, `unique`, `slice`, `intersection`, `nth` → `at`, `remove_values` → `remove_any`, `minus` → `difference`, `union` → `concat`, `union_distinct` → `union`, `range` |
| `random::` | Zero-argument `rand()` → `float()` |
| `object::` | `values`, `has` → `has_key`, `zip`, `keep_keys`, `merge`, `merge_recursive` → `merge_deep`, and one-argument `keys(obj)` |
| `datetime::` | `now`, `date` → `parse`, `date_dayofweek` → `day_of_week`, `date_dayofyear` → `day_of_year`, `date_leapyear` → `is_leap_year`; `date_year`, `date_month`, `date_day`, `date_hour`, `date_minute`, `date_second`, `date_millisecond`, `date_quarter`, `date_days_in_month`, `date_format`, `date_add`, `date_subtract` lose their `date_` prefix |
| `math::` | `pi`, `abs`, `acos`, `asin`, `atan`, `atan2`, `ceil`, `cos`, `degrees`, `exp`, `exp2`, `floor`, `log`, `log2`, `log10`, `pow`, `radians`, `round`, `sin`, `sqrt`, `tan` |

Argument-aware replacements also include:

| Legacy call | Replacement |
| --- | --- |
| `position(a, v)` or `position(a, v, false)` | `arrays::contains(a, v)` |
| `position(a, v, true)` | `arrays::index_of(a, v)` |
| `append(a, v)` / `push(a, v)`, optionally with `false` | `arrays::append(a, v)` |
| `remove_value(a, v)`, optionally with a negative integer literal | `arrays::remove(a, v)` |
| `sorted_unique(a)` | `arrays::sorted(arrays::unique(a))` |
| `shift(a)` | `arrays::slice(a, 1)` |
| `outersection(a, b)` | `arrays::symmetric_difference(a, b)` |
| `keys(o, false)` | `object::keys(o)` |
| `keys(o, true)` | `arrays::sorted(object::keys(o))` |
| `date_diff(a, b, unit, true)` | `datetime::diff(a, b, unit)` |

Boolean modes and negative integer limits may be parenthesized; arbitrary
constant expressions are not evaluated. Retained arguments preserve their
evaluation order and count. Array and object replacements use immutable
operations, never `arrays::mut::` or `object::mut::`.

Only parsed call targets are matched. Strings, comments, object keys, and
variable names are not treated as calls. Argument expressions, their order,
and error operators retain their meaning. Nested calls are considered
independently, including eligible calls inside already-qualified calls.
Already-qualified targets themselves are preserved.

Changed files use the canonical formatter, so whitespace and layout can change.
Files with no automatic edits retain their original bytes. Rerunning a successful
migration produces no further edits; unresolved manual findings remain.

### Migrate object functions

Object replacements use immutable `object::` operations, never `object::mut::`.
`keys(value)` and `keys(value, false)` become `object::keys(value)`.
`keys(value, true)` becomes `arrays::sorted(object::keys(value))`. Dynamic modes
and unsupported arities remain unchanged for manual review.

In v1, `ZIP` kept the first value for duplicate keys. The canonical `object::zip`
and its deprecated v2 global alias both use the last value. Review code relying
on the v1 behavior when upgrading the runtime, even before rewriting the call.

See [Object migration]({{< ref "/docs/migrations/v1-to-v2/standard-library" >}}#objects)
for compatibility changes and the [Objects reference]({{< ref "docs/standard-library/objects" >}})
for the current API.

### Calls requiring manual review

Both `check` and `run` explain why a supported call needs review:

- `join` is ambiguous between legacy path joining and modern global string joining.
- Dynamic modes in `position`, `keys`, `append`, and `push` need manual review.
- Unique `append`/`push` suppresses only the incoming duplicate; applying
  `arrays::unique` would also remove existing duplicates.
- Zero, positive, and dynamic `remove_value` limits have no canonical limit mode.
- `outersection` with three or more arrays uses exactly-one-input semantics;
  canonical symmetric difference uses odd-number-of-inputs semantics.
- `pop` needs an evaluation-count-preserving replacement; `unshift` must retain
  argument evaluation order; `remove_nth` retains host-list removal semantics.
- Parameterized `rand` uses historical rounded/floored calculations rather than
  canonical continuous or integer bounds.
- Unsupported arities of argument-aware rules require manual review.
- `date_compare` has component-range semantics that differ from `datetime::same`;
  `date_diff` without a literal `true` floating flag may truncate toward zero,
  whereas `datetime::diff` always returns a Float.
- `average`, `sum`, `min`, `max`, `median`, `percentile`, `stddev_population`,
  `stddev_sample`, `variance_population`, and `variance_sample` have permissive
  legacy behavior that differs from strict canonical math.
- A matching function declaration or function alias anywhere in the file may
  change call resolution. These checks are deliberately conservative, including
  case variants and declarations in nested scopes. A namespace alias blocks a
  replacement when it would redirect any introduced namespace, including both
  `arrays` and `object` for sorted keys.

Each call receives at most one semantic finding, with collision explanations
taking precedence. `check` reports the original path, line, and column; `run`
manual actions retain the original path and line.

A manual finding preserves that call and does not prevent other safe calls in
the same file from migrating. Manual findings alone do not make `run` fail.

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

If analysis succeeds but applying edits, formatting, or validating the rewritten
source fails, the entire file remains unchanged. This includes cases where the
formatter cannot preserve comments. Previously discovered semantic manual
actions remain in the report alongside the file-level failure, including in
`--dry-run` and `--print` modes. Other files can still migrate.

The selected directory itself is scanned even when its name would be excluded as a descendant, such as `.tmp` or `testdata`. Explicit symlink targets are rejected, and directory symlinks are not followed. A standalone `.fql` target is migrated directly even when it is located under a directory that recursive migration would exclude.

## Scope

`ferret migrate run` is a compatibility aid, not a general v1-to-v2 translator. It does not translate arbitrary v1 APIs or application logic, invent replacements for removed packages, or rewrite source in excluded descendant directories.

For an embedded Go application, continue with [Go embedding migration]({{< ref "/docs/migrations/v1-to-v2/go-embedding" >}}). The guide covers the temporary compatibility packages and the manual move from v1 compiler, runtime, and driver composition to the native Engine, Plan, Session, and module APIs.

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

{{< docs-related tiles="migrations-v1-to-v2,migrations-v1-to-v2-standard-library,migrations-v1-to-v2-go-embedding,tools-cli-run,embedding-go-getting-started,tools-cli-mod" >}}
