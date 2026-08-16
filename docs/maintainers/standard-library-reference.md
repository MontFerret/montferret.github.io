# Standard Library reference generation

The website generates its Standard Library reference from the versioned Ferret
Core API published at `https://ferretlang.org/ferret/index.json`.
Catalog-bearing publications contain sibling `api.json` and `catalog.json`
artifacts with the canonical identity `montferret/core`. Their reusable API
contracts are owned by `github.com/MontFerret/specs`.

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
fails the command. Only an HTTP `404` for `catalog.json` invokes the legacy flat
renderer used by immutable API-only releases. Other catalog failures never
fall back to stale or incomplete content.

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

`api.json` remains the authority for function identities, signatures, and
descriptions. `catalog.json` supplies category grouping, order, titles, and real
namespace roots. Categories are presentation concepts; they do not create
callable namespaces such as `math::abs`.

## Version history

Specs must publish the API Catalog package before Ferret can publish a release
that uses it. Ferret then publishes both sibling artifacts before this site's
exact `runtime.v2` pin is advanced. Existing API-only alpha releases remain
immutable and use the flat fallback. Historical and version-aware Standard
Library documentation remains deferred until after the final 2.0.0 release.
