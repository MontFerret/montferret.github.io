# Standard Library reference generation

The website generates its Standard Library reference from the versioned Ferret
Core API published at `https://ferretlang.org/ferret/index.json`. Selected
publications must contain sibling `api.json` and `catalog.json` artifacts with
the canonical identity `montferret/core`. Their reusable API contracts are
owned by `github.com/MontFerret/specs`.

## Select the documented Ferret version

`data/versions.yaml` is the version source for the website. The generator reads
`runtime.v2`, finds that exact version in the published index, and follows its
`href` to `api.json`. The matching `catalog.json` URL is derived from that
authoritative artifact directory. It does not select `latest` or the greatest
prerelease.

Changing `runtime.v2` replaces the current unversioned reference in place. Alpha
releases do not create retained documentation trees or a version picker.

## Generate locally

Run the focused target from the repository root:

```sh
mage generateStdlib
```

`mage generate` is an alias. `mage build`, `mage serve`, and `mage serveSearch`
all run generation before Hugo, so contributors do not copy or download an API
artifact manually.

Generation requires network access to the published index and selected API
artifacts. A missing version, network failure, malformed document, unsupported
schema, unexpected identity or version, or API/catalog membership mismatch
fails the command. The catalog is required: a `404` or any other catalog
failure does not fall back to flat, stale, or incomplete content. The prior
generated tree remains unchanged after any failure.

Requests use a 30-second client timeout and accept at most 4 MiB per document.
The selected API URL is resolved from the index entry, the catalog is resolved
as its sibling, and redirects must remain HTTP(S) requests on the index origin.

## Generated files

Generated Hugo pages live under `.generated/content/docs/standard-library` and
are mounted into Hugo's content tree. The directory is ignored by Git and is
replaced atomically on each successful run. Do not edit files in it.

The handwritten introduction remains at
`content/docs/standard-library/_index.md`. It is composed with the generated
sections through Hugo's page tree, without mixing curated prose into generated
files.

`api.json` remains the authority for function identities, signatures,
parameters, return values, descriptions, failures, and deprecation metadata.
`catalog.json` supplies category grouping, order, titles, and descriptions. Its
structured members can identify global functions with an empty namespace or
namespaced functions with their canonical API namespace. Categories are
presentation concepts; they do not create callable namespaces such as
`math::abs` or alter callable names.

Each catalog category produces one page beneath `/docs/standard-library/`.
Functions render as full API sections on that page, and the right-side "On this
page" menu links to identity-qualified anchors such as `#global-abs`,
`#io-fs-read`, and `#t-not-eq`. The menu follows the shared documentation
layout and is hidden at narrow widths. The generator does not create
`/functions/` or namespace-segment function pages. Former published I/O and
Testing namespace section URLs continue to redirect to their catalog
categories.

The importable Mage helpers live in `tools/stdlibdocs` and
`tools/registryroutes`. They intentionally have no command wrappers.

## Version history

Specs must publish the API Catalog package before Ferret can publish a release
that uses it. Ferret then publishes both sibling artifacts before this site's
exact `runtime.v2` pin is advanced. Existing API-only alpha releases remain
immutable publication state, but the website cannot select them. Historical
and version-aware Standard Library documentation remains deferred until after
the final 2.0.0 release.
