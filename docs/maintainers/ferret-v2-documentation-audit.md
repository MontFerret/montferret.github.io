# Ferret v2 documentation audit

Snapshot of November 16, 2026. Release availability was checked in the distribution channel readers use. The site pin, a tool’s bundled Core, and a hosted deployment can have different versions.

## Evidence

| Project | Public default revision | Published baseline | Implementation checked |
| --- | --- | --- | --- |
| ferret | [`3dcc1849`](https://github.com/MontFerret/ferret/tree/3dcc1849b08685e9cdc71021fd2b6dd9d4003077) | v2.0.0-alpha.62 | ANTLR grammar, compiler and VM integration tests, native engine, codecs, SDK and UAPI adapter |
| api | [`ceb51144`](https://github.com/MontFerret/api/tree/ceb51144c51034ede3b766985086990b76059274) | v1.0.0-alpha.21 (Go proxy) | Runtime, plan, session, debugger and output contracts and fixtures |
| wire | [`1e92cd14`](https://github.com/MontFerret/wire/tree/1e92cd14a47a2503a7588ff5ac6f3da2617cdbe7) | v1.0.0-alpha.5; website example, CLI and Lab use alpha.4 | Client, server, output framing, transport ownership and connection-lifetime tests |
| ferretd | [`7b11ce94`](https://github.com/MontFerret/ferretd/tree/7b11ce94abb55c37a08e79b33960acc8496a8378) | v1.0.0-alpha.13 | Runtime catalog, workspace, language services, DAP and Wire integration tests |
| cli | [`9bf903c5`](https://github.com/MontFerret/cli/tree/9bf903c54bb043be4c28e9b7121f28b8129fbd57) | v2.0.0-alpha.48 | Published binary help, command definitions, debugger and runtime tests; HTML rc.27 |
| lab | [`768e1aeb`](https://github.com/MontFerret/lab/tree/768e1aeb1802dd6938088eb6e2235735375db054) | v2.0.0-alpha.38 | CLI help, runtime selection, fixture servers and runner; alpha.37 has no downloadable archives |
| editorium | [`0ac08acd`](https://github.com/MontFerret/editorium/tree/0ac08acd2cd0b873d85e7af8db3f6584e80d2cec) | vscode/v0.1.0 (cdfd1836); no JetBrains release | Six VSIX packages, Marketplace, release manifest and daemon pin; development manifests and tests |
| ferret-js | [`1c04cfb0`](https://github.com/MontFerret/ferret-js/tree/1c04cfb04d3a7840dd9b77e7dddd1e45a12080d4) | 2.0.0-alpha.12 on npm (abc24d44) | Published JS, declarations and WASM; engine, bridge, output and lifecycle tests |
| worker | [`c4bb4eb5`](https://github.com/MontFerret/worker/tree/c4bb4eb5922c142a79fcde79a9326abc717d5f37) | v2.0.0-rc.32 (fbda7ffe) | Binary, image, HTTP controllers, options, Core alpha.56 and HTML rc.23; main uses Core alpha.57 |
| contrib | [`3028694f`](https://github.com/MontFerret/contrib/tree/3028694fc88b7f0648f52e4f6ce6ebc4ab9cea26) | Independent module tags; HTML v1.0.0-rc.29 | Module manifests, driver implementations and Lab fixtures; later flat CDP session change is on main |
| barn | [`a6e0616a`](https://github.com/MontFerret/barn/tree/a6e0616acc425f42af32ed720c10b3ee5349318e) | Public Registry index at the linked main SHA; v1.4.0 release | Registry records, generated index and module API; MontFerret/registry returned HTTP 404 |
| specs | [`e40d49e0`](https://github.com/MontFerret/specs/tree/e40d49e0514def62823fcba1d884eee506066154) | v1.12.0 | Published schemas, structured API types, catalog and module contracts |
| playground | [`66674e38`](https://github.com/MontFerret/playground/tree/66674e38e939b240c2375c97bcb3dd96b169a4e8) | v1.0.0-alpha.3 | Published OpenAPI and HTTP handlers; Core alpha.61 and HTML rc.24 |
| montferret.github.io | [`5ad783a6`](https://github.com/MontFerret/montferret.github.io/tree/5ad783a6b8024a2a6b3dd00f3cf91ca8995b794f) | Audit starting revision | Mage, Hugo mounts, shortcodes, theme, release updater, example module and workflow |
| mockery | [`7c1d3dd3`](https://github.com/MontFerret/mockery/tree/7c1d3dd362bf552764de7625be0de91e9e4c4ed2) | Published site fixtures probed separately | Browser examples, dynamic products, frames, pagination and search |

The Standard Library index selects the [versioned API](https://ferretlang.org/ferret/versions/2.0.0-alpha.62/api.json) and [matching catalog](https://ferretlang.org/ferret/versions/2.0.0-alpha.62/catalog.json). The website updater and documentation now use the generator’s actual `core.v2` key. The separately pinned Wire example retains its tested dependency set. The [public Registry](https://registry.ferretlang.org/index.json) is published from Barn.

## Major corrections

- The Go install and migration paths used unreleased-checkout instructions. They now use the published Core and adapter modules. The examples settle output, propagate cleanup errors, and shut down the HTTP service before closing its resources.
- JavaScript examples treated the result of `run()` as a decoded value. The guides use `runJSON()` for values and show how to collect, consume and close an output in Node and browsers.
- The Wire, daemon, CLI and editor pages now identify which process selects the runtime and which release supports debugging. The standalone Wire example retains alpha.4 while the guide also covers the connection-lifetime interface in alpha.5. The JetBrains page records release status.
- Invalid FQL, obsolete HTML calls, and incorrect fixture selectors and HTTP paths were repaired. CLI and Lab instructions reflect the shipped formatter, runtime selection, retry behavior and available downloads.

## Page inventory

The table covers all 131 original Markdown pages and both new pages, ferretd and JetBrains. Paths are relative to `content/docs`. Status and severity describe the material found at the starting revision. Retained pages were checked against the listed source or rendered navigation. Executable coverage and exceptions are recorded below.

| Page | Subject | Source | Initial status | Severity | Action |
| --- | --- | --- | --- | --- | --- |
| `_index.md` | Docs | Website routes, release guides and navigation | accurate | low | Retain |
| `embedding/_index.md` | Embedding | Website routes, release guides and navigation | accurate | low | Retain |
| `embedding/go/_index.md` | Go | Core engine and integration tests | stale | medium | Correct installation or resource handling |
| `embedding/go/application-guide.md` | Embed Ferret in a Go application | Core engine and integration tests | stale | high | Handle service shutdown before closing engine resources |
| `embedding/go/compiler-analysis.md` | Compiler Analysis | Core pkg/compiler and analysis tests | accurate | low | Retain |
| `embedding/go/configuration.md` | Runtime Configuration | Core engine and integration tests | accurate | low | Retain |
| `embedding/go/custom-functions.md` | Custom Functions | Core engine and integration tests | stale | medium | Correct installation or resource handling |
| `embedding/go/executing.md` | Executing Ferret | Core engine and integration tests | stale | medium | Correct installation or resource handling |
| `embedding/go/getting-started.md` | Getting Started | Core engine and integration tests | stale | high | Install published Go module and preserve cleanup errors |
| `embedding/go/host-values.md` | Host Values | Core runtime, VM and materialization | stale | high | Correct VM identity and adoption of closable parameters during result encoding |
| `embedding/go/modules.md` | Modules | Core module/SDK and Contrib HTML | stale | medium | Correct installation or resource handling |
| `embedding/go/output.md` | Output | API output contract and Core engine | stale | medium | Correct installation or resource handling |
| `embedding/go/parameters.md` | Parameters | Core engine and integration tests | stale | medium | Correct installation or resource handling |
| `embedding/go/programs.md` | Programs | Core engine and bytecode/artifact | stale | medium | Correct installation or resource handling |
| `embedding/go/safe-fql-construction.md` | Construct FQL Safely | Core compiler, FQL formatter and session parameters | accurate | low | Retain |
| `embedding/go/universal-api.md` | Universal API | API and Core uapi | stale | high | Install released adapter and connect live breakpoint replacement to DAP |
| `embedding/go/value-encoders.md` | Value Encoders | Core pkg/encoding and codec tests | stale | medium | Correct installation or resource handling |
| `embedding/javascript/_index.md` | JavaScript | npm package, JS engine and bridge | accurate | low | Retain |
| `embedding/javascript/custom-functions.md` | Custom Functions | npm package, JS engine and bridge | stale | high | Correct output API and examples |
| `embedding/javascript/executing.md` | Executing Ferret | npm package, JS engine and bridge | stale | high | Correct output API and examples |
| `embedding/javascript/getting-started.md` | Getting Started | npm package, JS engine and bridge | stale | high | Correct output API and examples |
| `embedding/javascript/limitations.md` | Runtime and WASM Limitations | npm package, JS engine and bridge | stale | high | Distinguish eager WASM execution from output delivery |
| `embedding/javascript/modules.md` | Modules and Lifecycle | npm package, JS engine and bridge | stale | high | Correct output API and examples |
| `embedding/javascript/parameters.md` | Parameters | npm package, JS engine and bridge | stale | high | Correct output API and examples |
| `embedding/wire-security.md` | Wire security | Wire client/server and website example tests | accurate | low | Retain |
| `embedding/wire.md` | Remote runtimes (Wire) | Wire client/server and website example tests | incomplete | high | Document alpha.5 connection lifetime, debug support and compatible example dependencies |
| `getting-started/_index.md` | Getting Started | CLI, Editorium, Core and npm release | accurate | low | Retain |
| `getting-started/installation.md` | Installation | CLI, Editorium, Core and npm release | stale | high | Install published packages |
| `getting-started/overview.md` | Overview | CLI, Editorium, Core and npm release | accurate | low | Retain |
| `getting-started/quick-start.md` | Quick Start | CLI, Editorium, Core and npm release | stale | high | Use array membership |
| `guides/_index.md` | Guides | Website navigation and tool guides | accurate | low | Retain |
| `guides/apis.md` | Work with APIs | CLI HTTP client, Worker policy and actual hosted response | stale | high | remove failing live editor version, keep CLI workflow |
| `guides/browser-pages.md` | Work with browser-driven pages | CLI, Worker, HTML drivers rc.23/27/29 and Mockery | stale | medium | fix stock selector, executable wait example and document origin configuration |
| `guides/deploying.md` | Deploy Ferret as a service | Worker image, CLI and Lab HTTP | stale | high | Replace nonexistent HTTP routes, flags and configuration |
| `guides/error-handling.md` | Error handling and resilience | Core compiler/VM and Contrib HTML | stale | high | Repair optional access inside an array projection |
| `guides/existing-html.md` | Process existing HTML | Core and Contrib HTML memory driver | stale | medium | Update |
| `guides/iframes.md` | Work with iframes | CDP page.frames, Mockery iframe scenario | stale | high | replace fictitious targets and unsupported generalization with verified frame queries |
| `guides/interactions.md` | Interact with pages | CDP event handlers and search/load-more scenarios | stale | high | repair actions, DOM readiness, counts and links |
| `guides/pagination.md` | Pagination patterns | loop contracts and Mockery numbered links, scroll batches | stale | high | correct URL/selector errors, wait for observed count, verify results |
| `guides/precompiled-programs.md` | Pre-compile and distribute programs | CLI build and Core artifact package | stale | medium | Update |
| `guides/screenshots.md` | Screenshots and PDFs | published HTML snapshot params and real rendered image bytes | stale | high | correct JPEG default, supported crop options, Worker/CLI difference, file format |
| `guides/static-pages.md` | Extract data from static pages | memory HTML, CSS / Mockery products, CLI flags | stale | high | remove duplicate QUERY token, repair product selectors and shell quoting |
| `guides/testing.md` | Test extraction scripts | Lab runner and Core assertions | stale | high | Repair release link and retry example |
| `guides/writing-plugins.md` | Develop a Ferret module | CLI scaffold, Core SDK and sdktest | stale | medium | Update |
| `language/_index.md` | Reference navigation | Core .62 | accurate | low | Leave unchanged |
| `language/control-flow/_index.md` | Control flow and query fallback | Core .62 | stale | medium | Update |
| `language/control-flow/dispatch.md` | Event emission and shorthand | Core .62 | accurate | low | Leave unchanged |
| `language/control-flow/error-handling.md` | Recovery, retry, cancellation, optional indexing | Core .62 | stale | medium | Update |
| `language/control-flow/for.md` | Iteration, collection, scope and result ownership | Core .62 | incomplete | low | Label grammar notation |
| `language/control-flow/match.md` | Patterns, guards and branches | Core .62 | accurate | low | Leave unchanged |
| `language/control-flow/query.md` | Host queries, dialects, modifiers and arguments | Core .62 / CLI .48 | stale | high | Update |
| `language/control-flow/subqueries.md` | Nested loop expressions and result ownership | Core .62 | accurate | low | Leave unchanged |
| `language/control-flow/waitfor.md` | Polling, observable events, group subscriptions and triggers | Core .62 / CLI .48 | stale | medium | Update |
| `language/expressions.md` | Expressions, scope and returned values | Core .62 | stale | high | Update |
| `language/functions/_index.md` | Function calls and providers | Core .62 | accurate | low | Leave unchanged |
| `language/functions/modules.md` | Runtime-registered namespace and module functions | Core .62 | stale | high | Update |
| `language/functions/user-defined.md` | Named functions, capture and identifier syntax | Core .62 | stale | medium | Update |
| `language/operators/_index.md` | Operator families | Core .62 | stale | medium | Update |
| `language/operators/arithmetic.md` | Numeric, temporal and host operations | Core .62 | accurate | low | Leave unchanged |
| `language/operators/array.md` | Expansion, flattening, quantifiers and comparison | Core .62 | stale | high | Update |
| `language/operators/coalescing.md` | None coalescing and short circuit | Core .62 | accurate | low | Leave unchanged |
| `language/operators/comparison.md` | Equality, order and containment | Core .62 | accurate | low | Leave unchanged |
| `language/operators/logical.md` | Truth rules, strict negation and evaluation | Core .62 | stale | medium | Update |
| `language/operators/precedence.md` | Operator precedence and grouping | Core .62 | stale | high | Update |
| `language/operators/range.md` | Range values and iteration | Core .62 | stale | medium | Update |
| `language/operators/ternary.md` | Conditional evaluation | Core .62 | accurate | low | Leave unchanged |
| `language/overview.md` | First scripts and runtime boundaries | Core .62 / CLI .48 | stale | high | Update |
| `language/parameters.md` | Binding syntax and parameter ownership | Core .62 | stale | medium | Update |
| `language/script-structure/_index.md` | Syntax, identifiers and event expressions | Core .62 | stale | high | Update |
| `language/script-structure/use.md` | Compile-time function and namespace aliases | Core .62 | stale | medium | Update |
| `language/types/_index.md` | Basic and host value categories | Core .62 | accurate | low | Leave unchanged |
| `language/types/basic.md` | Literal, temporal, scalar and collection values | Core .62 | accurate | low | Leave unchanged |
| `language/types/capabilities.md` | Host capabilities and capability errors | Core .62 / CLI .48 | stale | high | Update |
| `language/types/host.md` | Opaque values, queries and resource ownership | Core .62 / CLI .48 | stale | high | Update |
| `language/types/ordering.md` | Comparison, recursive order and equality | Core .62 | stale | high | Update |
| `language/types/serialization.md` | Result encoding, partial errors, resource closure | Core .62 / CLI .48 | stale | high | Update |
| `language/variables.md` | Mutable/immutable bindings and lexical scope | Core .62 | accurate | low | Leave unchanged |
| `migrations/_index.md` | Migrate Ferret | Core, API and CLI migration tests | stale | medium | Update migration guidance |
| `migrations/consumable-output.md` | Consumable output | Core, API and CLI migration tests | incomplete | high | Cover released Go and JS consumers |
| `migrations/encoding-api.md` | Encoding API | Core codecs and codec tests | stale | medium | Update migration guidance |
| `migrations/lab-runtime-configuration.md` | Lab runtime configuration | Lab runner, tests and release help | stale | medium | Update migration guidance |
| `migrations/v1-to-v2/_index.md` | Ferret v1 → v2 | Core, API and CLI migration tests | accurate | low | Retain |
| `migrations/v1-to-v2/go-embedding.md` | Go embedding migration | Core, API and CLI migration tests | stale | medium | Update migration guidance |
| `migrations/v1-to-v2/language.md` | Explicit script return and migration behavior | Core .62 / CLI .48 | accurate | low | Leave unchanged |
| `migrations/v1-to-v2/standard-library.md` | Legacy-call conversion and published availability | stdlib / CLI .48 | stale | medium | Update |
| `modules/_index.md` | Modules | Contrib module tags, CLI, published Registry | incomplete | low | add first-party module discovery/version information |
| `modules/api-comments.md` | Document module APIs | Specs v1.12.0 parser/fixtures, live API artifact and Barn analyzer | stale | high | document structured union/list and opaque legacy atoms |
| `modules/develop.md` | Develop a module | published init, manifest and namespace registration | stale | medium | correct case sensitivity |
| `modules/install.md` | Install a module | released module install and tested fresh Go host | stale | low | present version pin as tested version rather than minimum |
| `modules/publish.md` | Publish a module | CLI publish flags, Barn registry release/stamping | accurate | low | retain |
| `standard-library/_index.md` | Core library, host selection and published reference | stdlib / CLI .48 | incomplete | medium | Expand |
| `tools/_index.md` | Tools | Website routes, release guides and navigation | incomplete | high | Link daemon, editors and Playground |
| `tools/cli/_index.md` | CLI | CLI commands and release help | accurate | low | Retain |
| `tools/cli/browser.md` | Browser | CLI commands and release help | stale | high | Correct released CLI behavior |
| `tools/cli/build.md` | Build | CLI commands and release help | accurate | low | Retain |
| `tools/cli/check.md` | Check | CLI commands and release help | accurate | low | Retain |
| `tools/cli/configuration.md` | Configuration | CLI commands and release help | accurate | low | Retain |
| `tools/cli/debug.md` | Debug | CLI commands and release help | stale | high | Correct released CLI behavior |
| `tools/cli/fmt.md` | Fmt | CLI commands and release help | stale | high | Use the working stdin preview mode |
| `tools/cli/inspect.md` | Inspect | CLI commands and release help | accurate | low | Retain |
| `tools/cli/installation.md` | Installation | CLI commands and release help | stale | high | Verify actual installation and version reporting |
| `tools/cli/migrate.md` | Migrate | CLI commands and release help | accurate | low | Retain |
| `tools/cli/mod.md` | Mod | CLI commands and release help | accurate | low | Retain |
| `tools/cli/overview.md` | Overview | CLI commands and release help | stale | high | Correct released CLI behavior |
| `tools/cli/repl.md` | REPL | CLI commands and release help | stale | high | Correct released CLI behavior |
| `tools/cli/run.md` | Run | CLI commands and release help | stale | high | Correct released CLI behavior |
| `tools/cli/wire.md` | Wire runtime | Wire client/server and CLI release | stale | high | Correct released CLI behavior |
| `tools/docker.md` | Docker | reserved unpublished stub | accurate as draft | low | retain draft; Docker commands belong to the tools that publish images |
| `tools/ferretd.md` | Ferret daemon | ferretd catalog, LSP/DAP and integration tests | missing | high | Create daemon guide |
| `tools/jetbrains.md` | JetBrains IDEs | plugin.xml, Gradle compatibility, release tags/Marketplace | missing | high | create concise status and source pointer |
| `tools/lab/_index.md` | Lab | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/cdp.md` | CDP | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/ci.md` | CI | Lab runner, tests and release help | stale | high | Repair download via Lab pin |
| `tools/lab/configuration.md` | Configuration | Lab runner, tests and release help | stale | high | Correct retry budget |
| `tools/lab/http-runtime.md` | HTTP Runtime | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/installation.md` | Installation | Lab runner, tests and release help | stale | high | Use published archives and tagged build |
| `tools/lab/mock-api.md` | Mock API Server | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/overview.md` | Overview | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/runners.md` | Runners | Lab runner, tests and release help | stale | high | Correct repetition and retry semantics |
| `tools/lab/running-tests.md` | Running Tests | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/static-serving.md` | Static File Server | Lab runner, tests and release help | accurate | low | Retain |
| `tools/lab/wire-runtime.md` | Wire Runtime | Wire client/server and Lab release | stale | medium | Correct Lab behavior |
| `tools/lab/writing-tests.md` | Writing Tests | Lab runner, tests and release help | accurate | low | Retain |
| `tools/playground.md` | Playground | website editor shortcode and separate Playground alpha.3 OpenAPI | missing | high | publish existing route with both workflows |
| `tools/visual-studio-code.md` | Visual Studio Code | VSIX release manifest and ferretd.json vs main | incomplete | high | qualify local versus development external Run/Debug |
| `tools/worker/_index.md` | Worker | rc.32 source/controllers | accurate | low | retain |
| `tools/worker/api.md` | API | rc.32 responses / binary version | stale | high | remove unrelated Core pin from sample response and misleading exact diagnostics |
| `tools/worker/configuration.md` | Configuration | rc.32 main/options/network and flag help | stale | high | correct -no-chrome and scopes of HTTP clients |
| `tools/worker/deployment.md` | Deployment | rc.32 release assets, Dockerfile.release, GoReleaser and main | stale | medium | correct Chrome checks |
| `tools/worker/overview.md` | Overview | rc.32 newOptions, CLI rc/debug tests | incomplete | high | document full shipped module set, HTTP vs Wire Debug |
| `web-extraction/_index.md` | Web Extraction | published child routes and docs tiles | incomplete | medium | make existing navigation landing available |
| `web-extraction/query-errors.md` | Query Errors | compiler, VM tests, CLI and runtime | stale | medium | precise error timing, missing element, link recovery |
| `web-extraction/query-results.md` | Query Results | VM and encoding, published CLI/Worker | stale | high | repair undefined names, JSON shape, output example |

Generated sections are mounted at `/docs/standard-library/`. Each row represents one category page built from the published API and catalog. All callable identities were reconciled, including their namespace and category membership.

| Page | Subject | Source | Initial status | Severity | Action |
| --- | --- | --- | --- | --- | --- |
| `arrays/` | Arrays reference | Core API and catalog | accurate | low | Regenerate |
| `collections/` | Collections reference | Core API and catalog | accurate | low | Regenerate |
| `crypto/` | Crypto reference | Core API and catalog | accurate | low | Regenerate |
| `datetime/` | Datetime reference | Core API and catalog | accurate | low | Regenerate |
| `encoding/` | Encoding reference | Core API and catalog | accurate | low | Regenerate |
| `io/` | Io reference | Core API and catalog | accurate | low | Regenerate |
| `math/` | Math reference | Core API and catalog | accurate | low | Regenerate |
| `objects/` | Objects reference | Core API and catalog | accurate | low | Regenerate |
| `path/` | Path reference | Core API and catalog | accurate | low | Regenerate |
| `random/` | Random reference | Core API and catalog | accurate | low | Regenerate |
| `strings/` | Strings reference | Core API and catalog | accurate | low | Regenerate |
| `testing/` | Testing reference | Core API and catalog | accurate | low | Regenerate |
| `types/` | Types reference | Core API and catalog | accurate | low | Regenerate |
| `utils/` | Utils reference | Core API and catalog | accurate | low | Regenerate |

| Other entry point | Subject | Source | Initial status | Severity | Action |
| --- | --- | --- | --- | --- | --- |
| `/` | Homepage, code samples and ecosystem | CLI, Contrib and homepage template | stale | high | Correct examples, runtime boundaries and destinations |
| `/try/` | Online editor | Worker endpoint and live fixture | stale | high | Replace deprecated HTML calls |
| `/registry/` and generated module routes | Published module discovery | Barn, Specs and theme tests | accurate | low | Check rendered routes and registry data |
| `/search/` | Search | Pagefind queries and Mage | accurate | low | Verify existing and new destinations |
| `/contact/` | Support and release routing | GitHub release channel | stale | medium | Link the documented v2 prerelease |
| `/sponsor/` | Project support | Site content and navigation | accurate | low | Retain |
| `/blog/` and its 28 articles | Dated announcements and examples | Article dates, existing aliases and rendered links | accurate as history | low | Preserve articles; check navigation |

## Validation and remaining limitations

- The release CLI compiled 260 documentation editor blocks and the Try page. Of 61 historical blog blocks, three use v0.16 syntax and one deliberately reassigns a v2 immutable binding. We exercised all 47 changed language examples with the release CLI or Core observable fixtures. All 217 language editor blocks and 81 static FQL snippets compiled directly. Another 47 static excerpts compiled with their omitted bindings supplied. Three intentional semantic errors produced the documented diagnostics. All six homepage examples compiled. Browser, REST, SQLite, XLSX and CSV ran with fixtures. The AI example was checked against the provider contract without an API credential.
- Core alpha.62 and API alpha.21 resolved from the Go proxy. All 17 complete v2 Go programs, eight additional package examples and the example module’s SDK tests built with `GOWORK=off go test ./...` against published modules. Executed the basic queries, codecs, host values, HTML memory driver, HTTP handlers, migration, and artifact build/load examples. Preserved historical v1 source for comparison.
- npm’s exact alpha.12 package ran the Node ESM and CommonJS examples, decoded values, output collection, cancellation and host-function errors. Two TypeScript examples type-checked against the published declarations. The browser export ran against its packaged WASM assets in headless Chromium. The rendered Node requirement and Tools navigation were checked in a browser.
- The checksummed CLI alpha.48, Lab alpha.38 and ferretd alpha.13 macOS arm64 binaries exercised one real loopback Wire host. CLI source Run, REPL and step/locals/continue worked. Lab ran a host-function assertion. A named DAP launch stopped, continued and exited, and the host accepted a fresh query afterward. Checked 32 release help surfaces, 160 tool-guide commands and 282 flag occurrences, plus a sitewide check of terminal and shell blocks. Tagged CLI, Lab and daemon integration tests and the standalone Wire host/client and TLS/mTLS tests passed.
- Ran 42 browser/extraction examples on Worker, six CLI capture examples and seven HTTP examples with a controlled server. Checked returned values, PNG/JPEG/PDF signatures, Registry discovery and installation into an empty Go module. The Try example ran on the actual hosted Worker (12 products). A local release Worker container passed `GET /health`, parameterized `POST /`, and CLI and Lab HTTP execution. The local Worker check used its Linux amd64 image.
- Ran the website’s `go test ./...`, the standalone Wire module’s `GOWORK=off GOFLAGS=-mod=readonly go test ./...`, the 24 theme tests (`npm --prefix themes/ferret test`), and `mage build`. The production build used Hugo Extended 0.159.1, Node 24, Go 1.26.0 and Mage 1.14.0. It regenerated 14 categories covering 324 functions, checked Pagefind queries, and generated Registry routes. Pagefind indexed 149 pages. The rendered crawl checked 318 HTML files, 13,314 internal links and assets, and 2,291 fragments. It found zero missing destinations or anchors. Hugo also emitted its existing `.Site.Data` deprecation warning.

### Independent review

The reviewer reconciled the before-and-after inventory, inspected the complete diff and reproduced two invalid FQL examples. The retry-delay and expression-grouping examples now each have one terminal return. Rechecking static examples also identified grammar notation that needed a text label. Wire alpha.5's logical connection notification and the distinction between VM borrowing and Native result cleanup were added after review. A probe against Core alpha.62 recorded one close for a returned parameter, one when returned twice in an array, and zero when the script returned an unrelated scalar. A direct JSON codec call also recorded zero closes. The CLI stdin formatter workaround was exercised with valid and invalid input. The reviewer rechecked the corrected workflows and found no outstanding material technical issues.

### Upstream findings and limits

- **Core alpha.62.** An unparenthesized mixture such as `false || true && true` can produce a parser ambiguity. The precedence guide uses an explicitly grouped, tested expression.
- **CLI alpha.48.** Plain stdin with `ferret fmt` panicked with exit code 2. The documented `ferret fmt --dry-run` reads stdin, writes formatted text to stdout and exits 0 on valid input.
- **Lab alpha.38.** The attempt counter spans repetitions. With `--times=2 --attempts=2`, a success followed by a failure stopped the case. With `--times=3 --attempts=2`, two successes followed by two failures and one success took five attempts and passed. The runner guide documents the effect on bounded retries.
- **Worker rc.32.** `-no-chrome` skips Chrome startup/health checks while the CDP driver remains registered. The hosted website Worker reported rc.31 and Core alpha.55 during the audit. The separate Playground deployment returned HTTP 503 for environment discovery; its tagged service API was checked independently.
- **Editors.** Confirmed published VS Code 0.1.0 and its bundled daemon alpha.4. Editorium main has external Run/Debug and a JetBrains integration. The exact JetBrains Marketplace ID returned 404 and there was no tagged package. Editor UI installation and source-only workflows were not exercised.
- Authenticated module publication, Marketplace publishing, production deployment, and external AI-provider execution were outside the executed checks.
