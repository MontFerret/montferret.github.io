# Standard Library reference generation

The website generates its Standard Library reference from the versioned Ferret
Core API published at `https://ferretlang.org/ferret/index.json`. The published
API has the canonical identity `montferret/core` and follows the API contracts
owned by `github.com/MontFerret/specs`.

## Select the documented Ferret version

`data/versions.yaml` is the version source for the website. The generator reads
`runtime.v2`, finds that exact version in the published index, and follows its
`href`. It does not select `latest` or the greatest prerelease.

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

Generation requires network access to the published index and selected API. A
missing version, network failure, malformed document, unsupported schema,
unexpected API identity, or version mismatch fails the command. The build never
falls back to stale or incomplete API content.

Requests use a 30-second client timeout and accept at most 4 MiB per document.
The selected artifact URL is resolved from the index entry, and redirects must
remain HTTP(S) requests on the index origin.

## Generated files

Generated Hugo pages live under `.generated/content/docs/standard-library` and
are mounted into Hugo's content tree. The directory is ignored by Git and is
replaced atomically on each successful run. Do not edit files in it.

The handwritten introduction remains at
`content/docs/standard-library/_index.md`. It is composed with the generated
sections through Hugo's page tree, without mixing curated prose into generated
files.

## Version history

Phase 1 intentionally publishes one current, unversioned reference during the
Ferret v2 alpha. Historical and version-aware Standard Library documentation is
deferred until after the final 2.0.0 release.
