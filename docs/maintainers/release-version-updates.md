# Release version updates

The website owns the versions used for current tool download and release links.
Published releases notify this repository, and the `Update release version`
workflow proposes the corresponding `data/versions.yaml` change in a pull
request against `main`.

The workflow must be present on the repository's default branch before GitHub
will accept repository or manual dispatches for it.

## Release notification payload

Release repositories send a `repository_dispatch` event named
`ferret-release-published` to `MontFerret/montferret.github.io`:

```json
{
  "event_type": "ferret-release-published",
  "client_payload": {
    "repository": "MontFerret/cli",
    "tag": "v2.0.0-alpha.41",
    "release_url": "https://github.com/MontFerret/cli/releases/tag/v2.0.0-alpha.41",
    "commit_sha": "0123456789abcdef0123456789abcdef01234567"
  }
}
```

All four payload fields are required and treated as untrusted input. The source
repository must match the updater's exact allowlist, the tag must be strict
SemVer with at most one leading `v`, the release URL must match that repository
and tag, and the commit SHA must be a 40- or 64-character hexadecimal value.
The event cannot choose the target repository, branch, version file, YAML key,
or token scope.

## Update the file locally

Run the same updater from the repository root when checking a release by hand:

```sh
go run ./cmd/versions update \
  --repository MontFerret/cli \
  --tag v2.0.0-alpha.41 \
  --format json
```

The command defaults to `data/versions.yaml`. JSON output contains the resolved
product and major together with an `updated`, `unchanged`, or `stale` status.
Diagnostics and stale-release warnings are written to stderr.

## Retry an update manually

Open **Actions**, select **Update release version**, and choose **Run workflow**
on `main`. Enter the same repository, tag, release URL, and commit SHA that the
release notification would contain.

The updater is idempotent. A notification matching the current version exits
successfully without a branch. An older release is a successful stale no-op.
If the stable automation branch already proposes a newer version, that pending
version is also preserved. SemVer versions with equal precedence but different
build metadata fail for maintainer review because their release order cannot be
inferred.

## GitHub App configuration

The workflow reads these repository or organization secrets:

- `FERRET_RELEASE_APP_CLIENT_ID`
- `FERRET_RELEASE_APP_PRIVATE_KEY`

The Ferret Release App must be installed on
`MontFerret/montferret.github.io` with **Contents: write** and **Pull requests:
write**. The workflow requests a token for only this repository and those two
permissions. It uses the App token for the automation branch and pull request,
so the normal website pull-request checks run. Pull requests are never merged
automatically.

## Add a tracked tool or version channel

The authoritative repository mapping is in `internal/versions`. To add a tool:

1. Add its exact `owner/repository` mapping to the package.
2. Add the product and supported major key to `data/versions.yaml` with a
   normalized, double-quoted SemVer value.
3. Add mapping, channel, validation, and preservation tests.
4. Configure the release repository to send the documented dispatch after a
   release is published.

The SemVer major selects `product.v<major>`. The updater only modifies a key
that already exists, so a new major must be added deliberately to
`data/versions.yaml`; otherwise the notification fails without changing the
file. Top-level and channel keys must occur exactly once. The global `go` value
and product-specific `go` values are unrelated to releases and are never
updated by this command.
