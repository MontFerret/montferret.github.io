---
title: "Test extraction scripts"
sidebarTitle: "Testing"
weight: 70
draft: false
description: "Use Lab to write and run automated tests for FQL scripts."
---

# Test extraction scripts

As extraction scripts grow, testing helps catch regressions when pages change or scripts are refactored. [Lab]({{< ref "/docs/tools/lab" >}}) is Ferret's test runner — it executes FQL test files and reports pass/fail results.

This guide walks through the testing workflow. For the complete reference, see the [Lab documentation]({{< ref "/docs/tools/lab" >}}).

## Install Lab

{{< terminal command="true" >}}
go install github.com/nicktomlin/ferret-lab@latest
{{< /terminal >}}

Or pull the Docker image:

{{< terminal command="true" >}}
docker pull montferret/lab
{{< /terminal >}}

See [Lab Installation]({{< ref "/docs/tools/lab/installation" >}}) for all options.

## Write a test file

A test file is a regular `.fql` script that uses assertion functions from the `TESTING` namespace. Name the file with a `.test.fql` extension:

```
tests/
  headings.test.fql
```

{{< code lang="fql" title="tests/headings.test.fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET headings = page[~ css`h1, h2`]

TESTING::NOT_EMPTY(headings, "page should have headings")
TESTING::GT(LENGTH(headings), 0, "at least one heading expected")

RETURN TRUE
{{</ code >}}

## Run the tests

{{< terminal command="true" >}}
lab --dir ./tests
{{< /terminal >}}

Lab discovers all `.test.fql` files in the directory and runs them. A test passes when it completes without error; it fails when an assertion fails or a runtime error is raised.

## Write expected-failure tests

Name a file with `.fail.fql` to indicate that it *should* fail:

{{< code lang="fql" title="tests/missing-element.fail.fql" >}}
LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org")
LET el = QUERY ONE ".does-not-exist" IN page USING css
RETURN el.textContent
{{</ code >}}

This test passes only if the script produces an error.

## Write YAML test suites

For testing many queries with structured assertions, use `.yaml` test files:

```yaml
# tests/api.test.yaml
tests:
  - name: "posts endpoint returns data"
    query: |
      LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts")
      LET posts = JSON_PARSE(TO_STRING(response))
      RETURN LENGTH(posts)
    assert:
      gt: 0

  - name: "single post has title"
    query: |
      LET response = IO::NET::HTTP::GET("https://jsonplaceholder.typicode.com/posts/1")
      LET post = JSON_PARSE(TO_STRING(response))
      RETURN post.title
    assert:
      not_empty: true
```

## Use fixtures for reproducible tests

Testing against live websites is fragile — the page may change at any time. Lab can serve static HTML fixtures locally so tests run against stable content.

Create a fixtures directory:

```
tests/
  fixtures/
    products.html
  products.test.fql
```

Start the fixture server and run tests:

{{< terminal command="true" >}}
lab --dir ./tests --static ./tests/fixtures --static-port 8080
{{< /terminal >}}

Reference the local server in your test:

{{< code lang="fql" title="tests/products.test.fql" >}}
LET page = WEB::HTML::OPEN("http://localhost:8080/products.html")
LET items = page[~ css`.product-card`]

TESTING::NOT_EMPTY(items, "should find product cards")
TESTING::EQ(LENGTH(items), 3, "expected 3 products")

RETURN TRUE
{{</ code >}}

See [Static File Server]({{< ref "/docs/tools/lab/static-serving" >}}) for details.

## Run tests in CI

Lab can run in a CI pipeline. A common setup:

```yaml
# .github/workflows/test.yml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run FQL tests
        run: |
          lab --dir ./tests --static ./tests/fixtures --static-port 8080
```

For browser-backed tests, add a Chromium container:

```yaml
services:
  chromium:
    image: montferret/chromium
    ports:
      - 9222:9222
```

See [CI]({{< ref "/docs/tools/lab/ci" >}}) for the full reference.

## Control test execution

Lab supports several options for test execution:

{{< terminal command="true" >}}
lab --dir ./tests --concurrency 4 --timeout 30s --retries 2
{{< /terminal >}}

- `--concurrency` — run tests in parallel
- `--timeout` — maximum time per test
- `--retries` — retry failed tests

See [Runners]({{< ref "/docs/tools/lab/runners" >}}) for all options.

## Next steps

{{< docs-related tiles="tools-lab,tools-lab-writing-tests,tools-lab-ci,guide-error-handling" >}}
