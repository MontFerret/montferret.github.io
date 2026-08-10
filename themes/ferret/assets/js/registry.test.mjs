import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import test from "node:test";

import {
    RegistryPayloadError,
    RegistryApp,
    apiArtifactURL,
    apiFunctionID,
    apiNamespaceID,
    apiReferenceAnchorIDs,
    apiSignature,
    filterModules,
    hasAPIReference,
    loadAPIReference,
    modulePath,
    modulesForOwner,
    ownerPath,
    packageInstallCommand,
    parseRegistryRoute,
    registryRouteKey,
    renderAPIReference,
    resolveArtifactURL,
    selectVersion,
    validateAPIReference,
    versionStatus
} from "./registry.mjs";

const archiveAPIReference = JSON.parse(readFileSync(new URL("./registry.api.fixture.json", import.meta.url), "utf8"));

function clone(value) {
    return structuredClone(value);
}

function versionDocument(version = archiveAPIReference.version, content = {}) {
    return {
        schemaVersion: 1,
        id: archiveAPIReference.id,
        version,
        description: "Archive tools.",
        namespace: "ARCHIVE",
        license: "Apache-2.0",
        source: { repository: "https://github.com/MontFerret/contrib" },
        package: { path: "github.com/MontFerret/contrib/modules/archive" },
        content
    };
}

class FakeRoot {
    constructor() {
        this.dataset = { registryBase: "https://registry.example/" };
        this.innerHTML = "";
        this.attributes = new Map();
    }

    setAttribute(name, value) {
        this.attributes.set(name, value);
    }

    querySelector() {
        return null;
    }
}

test("artifact URLs remain on the discovered Registry origin", () => {
    assert.equal(
        resolveArtifactURL("https://registry.example/", "./modules/index.json").href,
        "https://registry.example/modules/index.json"
    );
    assert.throws(
        () => resolveArtifactURL("https://registry.example/", "https://outside.example/modules.json"),
        RegistryPayloadError
    );
    assert.throws(
        () => resolveArtifactURL("https://registry.example/", "javascript:alert(1)"),
        RegistryPayloadError
    );
    assert.throws(
        () => resolveArtifactURL("https://registry.example/", "/modules.json#fragment"),
        RegistryPayloadError
    );
    assert.throws(
        () => resolveArtifactURL("https://registry.example/", "/modules.json?download=1"),
        RegistryPayloadError
    );
});

test("Registry routes round-trip through history paths", () => {
    assert.deepEqual(parseRegistryRoute("/registry/"), { kind: "catalog" });
    assert.deepEqual(parseRegistryRoute(ownerPath("montferret")), {
        kind: "owner",
        owner: "montferret"
    });
    assert.deepEqual(parseRegistryRoute(modulePath("montferret/html")), {
        kind: "module",
        id: "montferret/html",
        owner: "montferret",
        name: "html",
        requestedVersion: ""
    });
    assert.equal(
        parseRegistryRoute(modulePath("montferret/html", "1.2.0")).requestedVersion,
        "1.2.0"
    );
    assert.deepEqual(parseRegistryRoute("/registry/montferret/html/1.2.0/extra/"), { kind: "invalid" });
    assert.deepEqual(parseRegistryRoute("/registry/../"), { kind: "invalid" });
    assert.deepEqual(parseRegistryRoute("/registry/../html/"), { kind: "invalid" });
});

test("Registry route identity ignores fragments but preserves path and query changes", () => {
    const modulePage = new URL("https://ferret.example/registry/montferret/html/#parse");
    const anotherFragment = new URL("https://ferret.example/registry/montferret/html/#load");
    const requestedVersion = new URL("https://ferret.example/registry/montferret/html/1.2.0/#parse");
    const query = new URL("https://ferret.example/registry/montferret/html/?view=compact#parse");

    assert.equal(registryRouteKey(modulePage), registryRouteKey(anotherFragment));
    assert.notEqual(registryRouteKey(modulePage), registryRouteKey(requestedVersion));
    assert.notEqual(registryRouteKey(modulePage), registryRouteKey(query));
});

test("owner filtering matches the exact validated owner segment", () => {
    const modules = [
        { id: "montferret/html" },
        { id: "montferret/csv" },
        { id: "mont/browser" },
        { id: "MontFerret/archive" },
        { id: "montferret/../escape" }
    ];

    assert.deepEqual(modulesForOwner(modules, "montferret").map((item) => item.id), [
        "montferret/html",
        "montferret/csv"
    ]);
    assert.deepEqual(modulesForOwner(modules, "MontFerret").map((item) => item.id), ["MontFerret/archive"]);
    assert.deepEqual(modulesForOwner(modules, "missing"), []);
    assert.deepEqual(modulesForOwner(modules, "../montferret"), []);
});

test("stable latest wins and explicitly requested history is preserved", () => {
    const module = {
        latest: "1.2.0",
        versions: [
            { version: "2.0.0-rc.1" },
            { version: "1.2.0" },
            { version: "1.1.0" }
        ]
    };
    assert.equal(selectVersion(module).version, "1.2.0");
    assert.equal(selectVersion(module, "1.1.0").version, "1.1.0");
    assert.equal(selectVersion(module, "9.9.9"), null);
    assert.deepEqual(versionStatus(module, "1.2.0"), { label: "Latest stable", tone: "stable" });
    assert.deepEqual(versionStatus(module, "1.1.0", "1.1.0"), { label: "Historical release", tone: "historical" });
});

test("a catalog with no stable release selects but does not relabel its first prerelease", () => {
    const module = {
        versions: [{ version: "2.0.0-rc.2" }, { version: "2.0.0-rc.1" }]
    };
    const selected = selectVersion(module);
    assert.equal(selected.version, "2.0.0-rc.2");
    assert.deepEqual(versionStatus(module, selected.version), { label: "Latest prerelease", tone: "prerelease" });
});

test("module filtering combines text and category constraints", () => {
    const modules = [
        { id: "montferret/html", description: "Browser automation", categories: ["web"] },
        { id: "montferret/csv", description: "Read tabular files", categories: ["data"] }
    ];
    assert.deepEqual(filterModules(modules, "browser", "").map((item) => item.id), ["montferret/html"]);
    assert.deepEqual(filterModules(modules, "", "data").map((item) => item.id), ["montferret/csv"]);
    assert.deepEqual(filterModules(modules, "html", "data"), []);
});

test("install metadata uses only package.path and never reconstructs it", () => {
    assert.equal(
        packageInstallCommand({ version: "1.2.0", package: { path: "github.com/MontFerret/contrib/modules/web/html" } }),
        "go get github.com/MontFerret/contrib/modules/web/html@v1.2.0"
    );
    assert.equal(packageInstallCommand({ version: "1.2.0", source: { repository: "https://github.com/example/repo" } }), "");
});

test("optional API artifacts resolve from the selected version without assuming a filename", () => {
    const endpoint = "https://registry.example/modules/montferret/archive/versions/1.0.0/index.json";
    const withAPI = versionDocument(undefined, { api: "./references/generated-module-v1.json" });
    const withoutAPI = versionDocument(undefined, { documentationHtml: "./docs.html" });

    assert.equal(hasAPIReference(withAPI), true);
    assert.equal(
        apiArtifactURL(endpoint, withAPI).href,
        "https://registry.example/modules/montferret/archive/versions/1.0.0/references/generated-module-v1.json"
    );
    assert.equal(hasAPIReference(withoutAPI), false);
    assert.equal(apiArtifactURL(endpoint, withoutAPI), null);
    assert.throws(
        () => apiArtifactURL(endpoint, versionDocument(undefined, { api: "https://outside.example/api.json" })),
        RegistryPayloadError
    );
    assert.throws(
        () => apiArtifactURL(endpoint, versionDocument(undefined, { api: "" })),
        RegistryPayloadError
    );
});

test("Barn API Reference v1 renders functions, variadic signatures, parameters, and documentation", () => {
    const reference = validateAPIReference(
        clone(archiveAPIReference),
        archiveAPIReference.id,
        archiveAPIReference.version
    );
    const html = renderAPIReference(reference);

    assert.match(html, /id="api-namespace-named-ARCHIVE"/);
    assert.match(html, /id="api-function-named-ARCHIVE-EXTRACT"/);
    assert.match(html, /ARCHIVE::EXTRACT\(\.\.\.args\)/);
    assert.match(html, /Extract writes eligible archive entries/);
    assert.match(html, />Parameters</);
    assert.match(html, />Variadic</);
    assert.doesNotMatch(html, />Returns</);
});

test("API rendering supports global and nested namespaces, overloads, and escaped prose", () => {
    const reference = {
        schemaVersion: 1,
        id: "montferret/archive",
        version: "1.1.0",
        namespaces: [
            {
                name: "",
                functions: [{
                    name: "RUN",
                    signatures: [{ parameters: [], documentation: "Global <script> text.\n\nSecond paragraph & details." }]
                }]
            },
            {
                name: "AI::LLM",
                functions: [{
                    name: "RUN",
                    signatures: [
                        { parameters: ["input"] },
                        { parameters: ["input", "optionsValue"], documentation: "Runs with options." }
                    ]
                }]
            }
        ]
    };

    validateAPIReference(reference, reference.id, reference.version);
    const html = renderAPIReference(reference);
    assert.match(html, /id="api-function-global-RUN"/);
    assert.match(html, /id="api-function-named-AI-LLM-RUN"/);
    assert.match(html, /AI::LLM::RUN\(input, optionsValue\)/);
    assert.match(html, /Global &lt;script&gt; text/);
    assert.match(html, /Second paragraph &amp; details/);
    assert.doesNotMatch(html, /<script>/);
});

test("API anchors are deterministic and collision-safe across namespaces and case", () => {
    assert.equal(apiNamespaceID(""), "api-namespace-global");
    assert.equal(apiNamespaceID("AI::LLM"), "api-namespace-named-AI-LLM");
    assert.equal(apiFunctionID("", "RUN"), "api-function-global-RUN");
    assert.equal(apiFunctionID("AI::LLM", "RUN"), "api-function-named-AI-LLM-RUN");
    assert.notEqual(apiFunctionID("AI::LLM", "RUN"), apiFunctionID("AI::LLM", "Run"));
    assert.notEqual(apiFunctionID("", "RUN"), apiFunctionID("global", "RUN"));
    assert.equal(apiSignature("AI::LLM", "GENERATE", { parameters: ["args"], variadic: true }), "AI::LLM::GENERATE(...args)");
    assert.deepEqual(apiReferenceAnchorIDs(archiveAPIReference), [
        "api-reference",
        "api-namespace-named-ARCHIVE",
        "api-function-named-ARCHIVE-EXTRACT",
        "api-function-named-ARCHIVE-LIST",
        "api-function-named-ARCHIVE-READ"
    ]);
});

test("generated API anchors reserve stable IDs from authored documentation", () => {
    const documentationAPIHeading = { id: "api-reference" };
    const documentationFunctionHeading = { id: "api-function-named-ARCHIVE-EXTRACT" };
    const existingDocumentationHeading = { id: "documentation-api-reference" };
    const generatedAPISection = { id: "api-reference" };
    const generatedFunctionHeading = { id: "api-function-named-ARCHIVE-EXTRACT" };
    const links = ["#api-reference", "#api-function-named-ARCHIVE-EXTRACT"].map((href) => ({
        href,
        getAttribute(name) { return name === "href" ? this.href : null; },
        setAttribute(name, value) { if (name === "href") this.href = value; }
    }));
    const documentation = {
        querySelectorAll(selector) {
            if (selector === "[id]") return [documentationAPIHeading, documentationFunctionHeading, existingDocumentationHeading];
            if (selector === "a[href]") return links;
            return [];
        }
    };
    const root = new FakeRoot();
    root.querySelector = (selector) => selector === "#registry-documentation" ? documentation : null;
    root.querySelectorAll = (selector) => selector === "[id]" ? [
        documentationAPIHeading,
        documentationFunctionHeading,
        existingDocumentationHeading,
        generatedAPISection,
        generatedFunctionHeading
    ] : [];

    new RegistryApp(root).reserveAPIAnchorIDs(archiveAPIReference);

    assert.equal(documentationAPIHeading.id, "documentation-api-reference-2");
    assert.equal(links[0].href, "#documentation-api-reference-2");
    assert.equal(documentationFunctionHeading.id, "documentation-api-function-named-ARCHIVE-EXTRACT");
    assert.equal(links[1].href, "#documentation-api-function-named-ARCHIVE-EXTRACT");
    assert.equal(generatedAPISection.id, "api-reference");
    assert.equal(generatedFunctionHeading.id, "api-function-named-ARCHIVE-EXTRACT");
});

test("malformed and mismatched API Reference artifacts are rejected", () => {
    const cases = [
        (reference) => { reference.schemaVersion = 2; },
        (reference) => { reference.id = "montferret/other"; },
        (reference) => { reference.namespaces = null; },
        (reference) => { reference.namespaces.push(clone(reference.namespaces[0])); },
        (reference) => { reference.namespaces[0].functions.push(clone(reference.namespaces[0].functions[0])); },
        (reference) => { reference.namespaces[0].functions[0].signatures[0].parameters = ["first", "rest"]; }
    ];

    for (const mutate of cases) {
        const reference = clone(archiveAPIReference);
        mutate(reference);
        assert.throws(
            () => validateAPIReference(reference, archiveAPIReference.id, archiveAPIReference.version),
            RegistryPayloadError
        );
    }
});

test("API loading returns ready, absent, and partial-failure states", async () => {
    const endpoint = "https://registry.example/modules/montferret/archive/versions/1.0.0-rc.3/index.json";
    const withAPI = versionDocument(undefined, { api: "./reference.json" });
    let requested = "";
    const ready = await loadAPIReference(endpoint, withAPI, async (url) => {
        requested = url.href;
        return clone(archiveAPIReference);
    });
    assert.equal(requested, "https://registry.example/modules/montferret/archive/versions/1.0.0-rc.3/reference.json");
    assert.equal(ready.status, "ready");
    assert.equal(ready.reference.namespaces[0].functions[0].name, "EXTRACT");

    const absent = await loadAPIReference(endpoint, versionDocument(undefined, {}), async () => {
        throw new Error("fetch must not run");
    });
    assert.deepEqual(absent, { status: "absent" });

    const failed = await loadAPIReference(endpoint, withAPI, async () => {
        throw new Error("Registry returned 500");
    });
    assert.deepEqual(failed, { status: "error" });

    const malformed = await loadAPIReference(endpoint, withAPI, async () => ({ schemaVersion: 1 }));
    assert.deepEqual(malformed, { status: "error" });
});

test("aborted API loading is not converted into a partial failure", async () => {
    const controller = new AbortController();
    const pending = loadAPIReference(
        "https://registry.example/modules/montferret/archive/versions/1.0.0-rc.3/index.json",
        versionDocument(undefined, { api: "./reference.json" }),
        async (_endpoint, signal) => new Promise((_resolve, reject) => {
            signal.addEventListener("abort", () => reject(new DOMException("Aborted", "AbortError")), { once: true });
        }),
        controller.signal
    );
    controller.abort();

    await assert.rejects(pending, (error) => error.name === "AbortError");
});

test("Registry JSON loading rejects an API redirect outside the configured origin", async () => {
    const app = new RegistryApp(new FakeRoot());
    const originalFetch = globalThis.fetch;
    globalThis.fetch = async () => ({
        ok: true,
        url: "https://outside.example/reference.json",
        json: async () => clone(archiveAPIReference)
    });

    try {
        await assert.rejects(
            app.fetchJSON(new URL("https://registry.example/reference.json")),
            RegistryPayloadError
        );
    } finally {
        globalThis.fetch = originalFetch;
    }
});

test("module rendering clears API content across version changes and preserves documentation on API failure", () => {
    const root = new FakeRoot();
    const app = new RegistryApp(root);
    const route = { id: "montferret/archive", owner: "montferret", name: "archive", requestedVersion: "" };
    const module = {
        id: route.id,
        latest: "1.0.0-rc.3",
        versions: [{ version: "1.0.0-rc.3" }, { version: "1.0.0-rc.2" }]
    };

    app.renderModuleView(
        route,
        module,
        versionDocument("1.0.0-rc.3", { api: "./reference.json" }),
        "<h1>Current documentation</h1>",
        "",
        { status: "loading" }
    );
    assert.match(root.innerHTML, /Current documentation/);
    assert.match(root.innerHTML, /id="api-reference"/);
    assert.match(root.innerHTML, /Loading API Reference/);

    app.renderModuleView(
        route,
        module,
        versionDocument("1.0.0-rc.3", { api: "./reference.json" }),
        "<h1>Current documentation</h1>",
        "",
        { status: "error" }
    );
    assert.match(root.innerHTML, /Current documentation/);
    assert.match(root.innerHTML, /API Reference unavailable/);

    app.renderModuleView(
        { ...route, requestedVersion: "1.0.0-rc.2" },
        module,
        versionDocument("1.0.0-rc.2", { documentationHtml: "./docs.html" }),
        "<h1>Previous documentation</h1>",
        "",
        { status: "absent" }
    );
    assert.match(root.innerHTML, /Previous documentation/);
    assert.doesNotMatch(root.innerHTML, /Current documentation/);
    assert.doesNotMatch(root.innerHTML, /id="api-reference"/);
});

test("a documentation failure does not prevent a successful API Reference", () => {
    const root = new FakeRoot();
    const app = new RegistryApp(root);
    const route = { id: "montferret/archive", owner: "montferret", name: "archive", requestedVersion: "" };
    const module = { id: route.id, latest: archiveAPIReference.version, versions: [{ version: archiveAPIReference.version }] };

    app.renderModuleView(
        route,
        module,
        versionDocument(undefined, { api: "./reference.json" }),
        "",
        "The documentation artifact could not be loaded for this release.",
        { status: "ready", reference: archiveAPIReference }
    );
    assert.match(root.innerHTML, /Documentation unavailable/);
    assert.match(root.innerHTML, /ARCHIVE::EXTRACT/);
});

test("different module versions render only their own API contents", () => {
    const nextReference = clone(archiveAPIReference);
    nextReference.version = "1.1.0";
    nextReference.namespaces[0].functions = [{
        name: "OPEN",
        signatures: [{ parameters: ["path"], documentation: "Open inspects an archive." }]
    }];
    validateAPIReference(nextReference, nextReference.id, nextReference.version);

    const currentHTML = renderAPIReference(archiveAPIReference);
    const nextHTML = renderAPIReference(nextReference);
    assert.match(currentHTML, /ARCHIVE::EXTRACT/);
    assert.doesNotMatch(currentHTML, /ARCHIVE::OPEN/);
    assert.match(nextHTML, /ARCHIVE::OPEN/);
    assert.doesNotMatch(nextHTML, /ARCHIVE::EXTRACT/);
});

test("a valid empty API Reference renders a useful empty state", () => {
    const reference = {
        schemaVersion: 1,
        id: "montferret/archive",
        version: "1.0.0",
        namespaces: []
    };
    validateAPIReference(reference, reference.id, reference.version);
    assert.match(renderAPIReference(reference), /does not publish any statically registered Ferret functions/);
});
