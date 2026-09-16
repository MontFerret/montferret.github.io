---
title: "Ferret v1 → v2"
sidebarTitle: "Ferret v1 → v2"
weight: 10
draft: false
description: "Check an existing Ferret project, apply supported migrations, and finish the manual move to v2."
---

# Migrate from Ferret v1 to v2

Start with a compatibility check, preview the automated changes, and then review
the language, standard-library, and embedding changes that affect your project.
Use a CLI release containing the migration rules you need and a runtime
containing the corresponding v2 APIs. Older releases can support fewer rules.

## 1. Check the project

Run the read-only FQL check from the project root:

{{< terminal command="true" >}}
ferret migrate check --from v1 .
{{< /terminal >}}

The check reports supported compatibility findings, including calls that need
manual review. A clean result covers only the implemented rules; it does not
establish that all v1 code or application behavior is compatible with v2.
The check examines FQL files, not Go application architecture or queries stored
in Go string literals.

## 2. Preview automated changes

Show the paths that would change, then inspect the diff:

{{< terminal command="true" >}}
ferret migrate run --dry-run .
ferret migrate run --print .
{{< /terminal >}}

These are separate preview modes and do not write files. Check and run have
different directory discovery rules; see the [CLI command reference]({{< ref "docs/tools/cli/migrate" >}})
before choosing the migration boundary.

## 3. Apply supported migrations

After reviewing the preview, apply the supported changes:

{{< terminal command="true" >}}
ferret migrate run .
{{< /terminal >}}

The CLI can make a final top-level collecting `FOR` result explicit, rewrite
supported legacy standard-library calls, and move supported Go imports to v2
compatibility packages. It leaves calls requiring manual review unchanged while
other safe calls can migrate. Review its diagnostics as well as the resulting diff.

## 4. Review FQL language changes

Follow [FQL language migration]({{< ref "docs/migrations/v1-to-v2/language" >}})
to check script results and loop-result ownership. Review embedded queries
manually: the CLI discovers lowercase `.fql` files, not FQL inside Go strings.

## 5. Migrate standard-library usage

Use the [standard-library migration guide]({{< ref "docs/migrations/v1-to-v2/standard-library" >}})
to review names, argument modes, and changed semantics. An automatic rename
does not eliminate runtime behavior changes such as duplicate-key handling.

## 6. Migrate Go embedding code

For embedded Go applications, continue with [Go embedding migration]({{< ref "docs/migrations/v1-to-v2/go-embedding" >}}).
Compatibility imports are a temporary stage. Engine composition, modules,
drivers, execution, and resource ownership require application-level work.

## 7. Format, test, and finish

Format the remaining source with the [FQL formatter]({{< ref "docs/tools/cli/fmt" >}}).
Changed FQL files are already formatted by the migrator; unchanged files may
still need formatting. Rerun the compatibility check and resolve its remaining
findings. Test representative inputs and compare the results with the v1 project.

Remove deprecated calls and compatibility imports after replacing their behavior
and validating the application. For Go applications, also test cancellation and
cleanup, and remove the v1 dependency only after all v1 imports are gone.

## Next steps

{{< docs-related tiles="migrations-v1-to-v2-language,migrations-v1-to-v2-standard-library,migrations-v1-to-v2-go-embedding,tools-cli-migrate" >}}
