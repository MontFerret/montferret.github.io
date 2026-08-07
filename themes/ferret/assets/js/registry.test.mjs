import assert from "node:assert/strict";
import test from "node:test";

import {
    RegistryPayloadError,
    filterModules,
    modulePath,
    modulesForOwner,
    ownerPath,
    packageInstallCommand,
    parseRegistryRoute,
    resolveArtifactURL,
    selectVersion,
    versionStatus
} from "./registry.mjs";

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
