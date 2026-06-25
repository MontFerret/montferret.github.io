---
title: "Running Tests"
weight: 40
draft: false
description: "Run Lab tests from local paths, HTTP URLs, and Git repositories."
---

# Running Tests

Use `lab run` to execute FQL unit tests and YAML suites.

{{< terminal >}}
lab run tests/
{{< /terminal >}}

If no files are provided, Lab prints command help and exits with status `1`.

## Run local files and directories

Run one file:

{{< terminal >}}
lab run tests/users.fql
{{< /terminal >}}

Run every supported test file under a directory:

{{< terminal >}}
lab run tests/
{{< /terminal >}}

You can pass multiple locations as positional arguments:

{{< terminal >}}
lab run tests/unit tests/integration/users.yaml
{{< /terminal >}}

Or repeat `--files`:

{{< terminal >}}
lab run --files tests/unit --files tests/integration/users.yaml
{{< /terminal >}}

Local directories are traversed recursively. Lab reads `.fql`, `.yaml`, and `.yml` files.

## Use file URLs and filters

Local sources can also be addressed with `file://` URLs. Add `filter` to select matching files.

{{< terminal >}}
lab run "file:///absolute/path/to/tests?filter=**/*.fql"
{{< /terminal >}}

The filter is a glob pattern. For local files, it is matched against the full file path Lab reads.

## Run HTTP sources

HTTP and HTTPS locations are fetched with `GET` and treated as a single test file.

{{< terminal >}}
lab run https://example.com/tests/users.yaml
{{< /terminal >}}

The server must return a successful `2xx` response. The URL should point directly to a supported test file.

## Run Git sources

Use `git+http://` or `git+https://` to read tests from a Git repository.

{{< terminal >}}
lab run "git+https://github.com/example/project-tests.git?filter=tests/**/*.yaml"
{{< /terminal >}}

Lab clones the repository, reads the current default HEAD, and runs supported files. Add a `filter` query parameter to limit which repository files are selected.

## Pass parameters

Use `--param` to pass user parameters into FQL. The flag can be repeated.

{{< terminal >}}
lab run tests/users.yaml \
  --param='url:"https://example.com/users"' \
  --param='limit:10' \
  --param='active:true'
{{< /terminal >}}

Parameter values are parsed as JSON. Strings must be JSON strings, so quote them inside the value as shown above.

In FQL, parameters are available by name:

```fql
RETURN {
  url: @url,
  limit: @limit,
  active: @active
}
```

## Wait for dependencies

Use `--wait` when tests depend on an HTTP service that may not be ready yet.

{{< terminal >}}
lab run tests/integration \
  --wait=http://127.0.0.1:9222/json/version \
  --wait-timeout=10 \
  --wait-attempts=12
{{< /terminal >}}

Lab waits before it creates the runtime and starts executing tests. `--wait-timeout` is the interval in seconds passed to each wait attempt, and `--wait-attempts` controls how many attempts are made.

## Choose a reporter

Lab currently provides `console` and `simple` reporters.

{{< terminal >}}
lab run --reporter=console tests/
lab run --reporter=simple tests/
{{< /terminal >}}

The `simple` reporter prints line-oriented output that is easier to read in CI logs:

```text
PASS file="tests/users.fql" duration=12ms attempts=1 times=1
DONE passed=1 failed=0 duration=12ms
```

## Exit behavior

Lab exits with status `0` when all tests pass. It exits with status `1` when any test fails, a source cannot be read, command arguments are invalid, or a wait condition times out.

If the run is interrupted, Lab cancels the context and prints `Terminated`.

## Related topics

{{< docs-related tiles="tools-lab-services,tools-lab-runners,tools-lab-runtime-remote,tools-lab-runtime-binary" >}}
