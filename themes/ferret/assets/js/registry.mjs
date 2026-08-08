const SCHEMA_VERSION = 1;
const ROUTE_SEGMENT = /^[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?$/;

export class RegistryPayloadError extends Error {}

export function resolveArtifactURL(baseURL, href) {
    if (typeof href !== "string" || href.length === 0) {
        throw new RegistryPayloadError("An artifact link is missing.");
    }

    const base = new URL(baseURL);
    const resolved = new URL(href, base);
    if (!/^https?:$/.test(resolved.protocol) || resolved.origin !== base.origin || resolved.username || resolved.password || resolved.search || resolved.hash) {
        throw new RegistryPayloadError("An artifact link points outside the Registry.");
    }

    return resolved;
}

export function parseRegistryRoute(pathname) {
    const parts = pathname.split("/").filter(Boolean);
    if (parts[0] !== "registry") {
        return { kind: "invalid" };
    }
    if (parts.length === 1) {
        return { kind: "catalog" };
    }
    if (parts.length === 2) {
        return ROUTE_SEGMENT.test(parts[1])
            ? { kind: "owner", owner: parts[1] }
            : { kind: "invalid" };
    }
    if ((parts.length !== 3 && parts.length !== 4) || parts.slice(1).some((part) => !ROUTE_SEGMENT.test(part))) {
        return { kind: "invalid" };
    }

    return {
        kind: "module",
        id: `${parts[1]}/${parts[2]}`,
        owner: parts[1],
        name: parts[2],
        requestedVersion: parts[3] || ""
    };
}

export function selectVersion(moduleDocument, requestedVersion = "") {
    const versions = Array.isArray(moduleDocument?.versions) ? moduleDocument.versions : [];
    if (requestedVersion) {
        return versions.find((entry) => entry.version === requestedVersion) || null;
    }
    if (moduleDocument?.latest) {
        return versions.find((entry) => entry.version === moduleDocument.latest) || null;
    }
    return versions[0] || null;
}

export function versionStatus(moduleDocument, selectedVersion, requestedVersion = "") {
    if (moduleDocument.latest && selectedVersion === moduleDocument.latest) {
        return { label: "Latest stable", tone: "stable" };
    }
    if (!moduleDocument.latest && moduleDocument.versions?.[0]?.version === selectedVersion) {
        return { label: "Latest prerelease", tone: "prerelease" };
    }
    if (requestedVersion) {
        return { label: "Historical release", tone: "historical" };
    }
    return { label: "Prerelease", tone: "prerelease" };
}

export function filterModules(modules, query, categoryID = "") {
    const normalized = query.trim().toLowerCase();
    return modules.filter((module) => {
        const searchable = `${module.id} ${module.description || ""}`.toLowerCase();
        const matchesQuery = !normalized || searchable.includes(normalized);
        const matchesCategory = !categoryID || module.categories?.includes(categoryID);
        return matchesQuery && matchesCategory;
    });
}

export function packageInstallCommand(versionDocument) {
    const packagePath = versionDocument?.package?.path;
    if (typeof packagePath !== "string" || packagePath.trim() === "" || typeof versionDocument?.version !== "string") {
        return "";
    }
    return `go get ${packagePath}@v${versionDocument.version}`;
}

function assertV1(document, description) {
    if (!document || typeof document !== "object" || document.schemaVersion !== SCHEMA_VERSION) {
        throw new RegistryPayloadError(`${description} is not a Registry v1 document.`);
    }
    return document;
}

function escapeHTML(value) {
    return String(value ?? "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#39;");
}

function externalHTTPSURL(value) {
    try {
        const parsed = new URL(value);
        return parsed.protocol === "https:" && !parsed.username && !parsed.password ? parsed.href : "";
    } catch {
        return "";
    }
}

export function modulePath(id, version = "") {
    return `/registry/${id}${version ? `/${version}` : ""}/`;
}

export function ownerPath(owner) {
    return `/registry/${owner}/`;
}

export function registryRouteKey(location) {
    return `${location.pathname}${location.search}`;
}

export function modulesForOwner(modules, owner) {
    if (!ROUTE_SEGMENT.test(owner)) return [];
    return modules.filter((module) => validModuleID(module?.id) && module.id.split("/")[0] === owner);
}

function fragmentID(hash) {
    if (typeof hash !== "string" || !hash.startsWith("#") || hash.length === 1) {
        return "";
    }

    try {
        return decodeURIComponent(hash.slice(1));
    } catch {
        return "";
    }
}

class RegistryApp {
    constructor(root) {
        this.root = root;
        this.baseURL = new URL(root.dataset.registryBase);
        this.discovery = null;
        this.moduleCatalog = null;
        this.categoryCatalog = null;
        this.renderController = null;
        this.renderedRouteKey = "";
    }

    start() {
        document.addEventListener("click", (event) => {
            const link = event.target.closest?.("a[href]");
            if (!link || event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) {
                return;
            }

            const destination = new URL(link.href, window.location.href);
            if (destination.origin !== window.location.origin || !destination.pathname.startsWith("/registry/")) {
                return;
            }
            if (destination.pathname === window.location.pathname && destination.search === window.location.search && destination.hash) {
                return;
            }

            event.preventDefault();
            window.history.pushState({}, "", `${destination.pathname}${destination.search}${destination.hash}`);
            this.render();
        });
        window.addEventListener("popstate", () => {
            if (registryRouteKey(window.location) === this.renderedRouteKey) {
                this.scrollToFragment();
                return;
            }

            this.render();
        });
        this.render();
    }

    async fetchJSON(endpoint, signal) {
        const response = await fetch(endpoint, {
            headers: { Accept: "application/json" },
            signal
        });
        if (!response.ok) {
            throw new Error(`Registry returned ${response.status}.`);
        }
        if (new URL(response.url).origin !== this.baseURL.origin) {
            throw new RegistryPayloadError("A Registry request redirected outside the Registry.");
        }
        return response.json();
    }

    async discover(signal) {
        if (this.discovery) {
            return this.discovery;
        }
        const root = assertV1(await this.fetchJSON(this.baseURL, signal), "The Registry root");
        if (!root.artifacts?.modules) {
            throw new RegistryPayloadError("The Registry root does not publish a module catalog.");
        }
        this.discovery = {
            modules: resolveArtifactURL(this.baseURL, root.artifacts.modules),
            categories: root.artifacts.categories ? resolveArtifactURL(this.baseURL, root.artifacts.categories) : null
        };
        return this.discovery;
    }

    async catalogs(signal, includeCategories = true) {
        const discovery = await this.discover(signal);
        if (!this.moduleCatalog) {
            this.moduleCatalog = assertV1(await this.fetchJSON(discovery.modules, signal), "The module catalog");
            if (!Array.isArray(this.moduleCatalog.modules)) {
                throw new RegistryPayloadError("The module catalog has no module list.");
            }
        }
        if (includeCategories && !this.categoryCatalog && discovery.categories) {
            try {
                const categories = assertV1(await this.fetchJSON(discovery.categories, signal), "The category catalog");
                this.categoryCatalog = Array.isArray(categories.categories) ? categories : { schemaVersion: 1, categories: [] };
            } catch (error) {
                if (error.name === "AbortError") throw error;
                this.categoryCatalog = { schemaVersion: 1, categories: [], failed: true };
            }
        }
        return {
            modules: this.moduleCatalog,
            categories: includeCategories ? this.categoryCatalog || { schemaVersion: 1, categories: [] } : { schemaVersion: 1, categories: [] },
            moduleCatalogURL: discovery.modules,
            categoryCatalogURL: discovery.categories
        };
    }

    async loadModuleRecords(entries, moduleCatalogURL, signal) {
        const results = await Promise.allSettled(entries.map(async (entry) => {
            if (!validModuleID(entry?.id) || !entry?.href) throw new RegistryPayloadError("A module catalog entry is incomplete.");
            const endpoint = resolveArtifactURL(moduleCatalogURL, entry.href);
            const module = assertV1(await this.fetchJSON(endpoint, signal), `Module ${entry.id}`);
            if (module.id !== entry.id) throw new RegistryPayloadError(`Module ${entry.id} has mismatched metadata.`);
            return { ...module, categories: [] };
        }));
        if (signal.aborted) throw new DOMException("Aborted", "AbortError");

        return {
            modules: results.filter((result) => result.status === "fulfilled").map((result) => result.value),
            failed: results.filter((result) => result.status === "rejected").length
        };
    }

    setBusy(busy) {
        this.root.setAttribute("aria-busy", String(busy));
    }

    loading(title, message) {
        this.setBusy(true);
        this.root.innerHTML = `
            <section class="registry-state" role="status">
                <span class="registry-spinner" aria-hidden="true"></span>
                <div><h2>${escapeHTML(title)}</h2><p>${escapeHTML(message)}</p></div>
            </section>`;
    }

    unavailable(title, message, retry = true) {
        this.setBusy(false);
        this.root.innerHTML = `
            <section class="registry-state registry-state-warning" role="alert">
                <p class="registry-state-label">Registry unavailable</p>
                <h2>${escapeHTML(title)}</h2>
                <p>${escapeHTML(message)}</p>
                ${retry ? '<button class="button registry-retry" type="button">Try again</button>' : '<a class="button" href="/registry/">Back to Registry</a>'}
            </section>`;
        this.root.querySelector(".registry-retry")?.addEventListener("click", () => {
            this.discovery = null;
            this.moduleCatalog = null;
            this.categoryCatalog = null;
            this.render();
        });
    }

    async render() {
        this.renderController?.abort();
        this.renderController = new AbortController();
        const { signal } = this.renderController;
        this.renderedRouteKey = registryRouteKey(window.location);
        const route = parseRegistryRoute(window.location.pathname);

        try {
            if (route.kind === "catalog") {
                await this.renderCatalog(signal);
            } else if (route.kind === "owner") {
                await this.renderOwner(route, signal);
            } else if (route.kind === "module") {
                await this.renderModule(route, signal);
            } else {
                this.unavailable("This Registry route is not valid", "Use the Registry catalog to choose a published module.", false);
            }
        } catch (error) {
            if (error.name === "AbortError") return;
            this.unavailable("The Registry could not be loaded", error instanceof RegistryPayloadError ? error.message : "The service may be temporarily unavailable. Try again shortly.");
        }
    }

    async renderCatalog(signal) {
        this.loading("Loading modules", "Following the Registry's published catalog links…");
        const { modules: catalog, categories, moduleCatalogURL, categoryCatalogURL } = await this.catalogs(signal);
        const { modules: loadedModules, failed: failedModules } = await this.loadModuleRecords(catalog.modules, moduleCatalogURL, signal);
        if (loadedModules.length === 0 && catalog.modules.length > 0) {
            throw new Error("No module records could be loaded.");
        }

        const moduleByID = new Map(loadedModules.map((module) => [module.id, module]));
        const categoryResults = await Promise.allSettled(categories.categories.map(async (entry) => {
            const endpoint = resolveArtifactURL(categoryCatalogURL, entry.href);
            const document = assertV1(await this.fetchJSON(endpoint, signal), `Category ${entry.id}`);
            for (const moduleEntry of document.modules || []) {
                moduleByID.get(moduleEntry.id)?.categories.push(entry.id);
            }
        }));
        if (signal.aborted) throw new DOMException("Aborted", "AbortError");
        const failedCategories = (categories.failed ? 1 : 0) + categoryResults.filter((result) => result.status === "rejected").length;

        this.renderCatalogView(loadedModules, categories.categories, failedModules, failedCategories);
    }

    renderCatalogView(modules, categories, failedModules, failedCategories) {
        this.setBusy(false);
        this.root.innerHTML = `
            <section class="registry-catalog" aria-labelledby="registry-catalog-title">
                <div class="registry-catalog-heading">
                    <div><p class="registry-section-label">Published modules</p><h2 id="registry-catalog-title">Browse the catalog</h2></div>
                    <p id="registry-result-count" class="registry-result-count" role="status"></p>
                </div>
                ${(failedModules || failedCategories) ? `<div class="registry-notice" role="status">Some Registry metadata could not be loaded. Available modules are shown below.</div>` : ""}
                <div class="registry-filters">
                    <label class="registry-field">
                        <span>Search modules</span>
                        <input id="registry-search" class="input" type="search" placeholder="Name or description" autocomplete="off">
                    </label>
                    <label class="registry-field">
                        <span>Category</span>
                        <select id="registry-category" class="select-input">
                            <option value="">All categories</option>
                            ${categories.map((category) => `<option value="${escapeHTML(category.id)}">${escapeHTML(category.name)} (${escapeHTML(category.count)})</option>`).join("")}
                        </select>
                    </label>
                </div>
                <div id="registry-results" class="registry-grid"></div>
            </section>`;

        const search = this.root.querySelector("#registry-search");
        const category = this.root.querySelector("#registry-category");
        const results = this.root.querySelector("#registry-results");
        const count = this.root.querySelector("#registry-result-count");
        const update = () => {
            const filtered = filterModules(modules, search.value, category.value);
            count.textContent = `${filtered.length} ${filtered.length === 1 ? "module" : "modules"}`;
            if (filtered.length === 0) {
                results.innerHTML = `<section class="registry-state registry-state-inline"><h3>No modules found</h3><p>Try another search or category.</p></section>`;
                return;
            }
            results.innerHTML = this.renderModuleCards(filtered);
        };
        search.addEventListener("input", update);
        category.addEventListener("change", update);
        update();
    }

    async renderOwner(route, signal) {
        this.loading("Loading owner", `Resolving modules published by ${route.owner}…`);
        const { modules: catalog, moduleCatalogURL } = await this.catalogs(signal, false);
        const entries = modulesForOwner(catalog.modules, route.owner);
        if (entries.length === 0) {
            this.unavailable("Owner not found", "This owner does not have any modules in the current Registry catalog.", false);
            return;
        }

        const { modules, failed } = await this.loadModuleRecords(entries, moduleCatalogURL, signal);
        if (modules.length === 0) {
            throw new Error("No module records could be loaded.");
        }

        this.renderOwnerView(route.owner, modules, failed);
    }

    renderOwnerView(owner, modules, failedModules) {
        this.setBusy(false);
        this.root.innerHTML = `
            <section class="registry-catalog" aria-labelledby="registry-owner-title">
                <nav class="registry-breadcrumbs" aria-label="Breadcrumb">
                    <a href="/registry/">Registry</a><span aria-hidden="true">/</span><span aria-current="page">${escapeHTML(owner)}</span>
                </nav>
                <div class="registry-catalog-heading">
                    <div><p class="registry-section-label">Module owner</p><h2 id="registry-owner-title">Modules by ${escapeHTML(owner)}</h2></div>
                    <p class="registry-result-count" role="status">${modules.length} ${modules.length === 1 ? "module" : "modules"}</p>
                </div>
                ${failedModules ? `<div class="registry-notice" role="status">Some module metadata could not be loaded. Available modules are shown below.</div>` : ""}
                <div class="registry-grid">${this.renderModuleCards(modules)}</div>
            </section>`;
    }

    renderModuleCards(modules) {
        return modules.map((module) => {
            const selected = selectVersion(module);
            const status = selected ? versionStatus(module, selected.version) : null;
            return `<a class="registry-card" href="${escapeHTML(modulePath(module.id))}">
                <span class="registry-card-arrow" aria-hidden="true">→</span>
                <h3>${escapeHTML(module.id)}</h3>
                <p>${escapeHTML(module.description || "No description is available for this module.")}</p>
                <span class="registry-card-meta">${selected ? `<span>${escapeHTML(selected.version)}</span><span class="registry-status registry-status-${status.tone}">${escapeHTML(status.label)}</span>` : "No releases"}</span>
            </a>`;
        }).join("");
    }

    async renderModule(route, signal) {
        this.loading("Loading module", `Resolving ${route.id}${route.requestedVersion ? ` at ${route.requestedVersion}` : ""}…`);
        const { modules: catalog, moduleCatalogURL } = await this.catalogs(signal);
        const entry = catalog.modules.find((candidate) => candidate.id === route.id);
        if (!entry) {
            this.unavailable("Module not found", "This module is not present in the current Registry catalog. The route may be stale.", false);
            return;
        }

        const moduleEndpoint = resolveArtifactURL(moduleCatalogURL, entry.href);
        const module = assertV1(await this.fetchJSON(moduleEndpoint, signal), `Module ${route.id}`);
        if (module.id !== route.id || !Array.isArray(module.versions)) {
            throw new RegistryPayloadError("The module document does not match the requested module.");
        }
        if (module.versions.some((entry) => !ROUTE_SEGMENT.test(entry?.version || "") || !entry?.href)) {
            throw new RegistryPayloadError("The module document contains an invalid version entry.");
        }

        const selected = selectVersion(module, route.requestedVersion);
        if (!selected) {
            const message = route.requestedVersion
                ? `Version ${route.requestedVersion} is not published for ${route.id}. The route may be stale.`
                : `${route.id} has no published versions.`;
            this.unavailable("Version not found", message, false);
            return;
        }

        const versionEndpoint = resolveArtifactURL(moduleEndpoint, selected.href);
        const version = assertV1(await this.fetchJSON(versionEndpoint, signal), `${route.id}@${selected.version}`);
        if (version.id !== route.id || version.version !== selected.version) {
            throw new RegistryPayloadError("The version document does not match the requested release.");
        }

        let documentation = "";
        let documentationError = "";
        if (version.content?.documentationHtml) {
            try {
                const docsEndpoint = resolveArtifactURL(versionEndpoint, version.content.documentationHtml);
                const response = await fetch(docsEndpoint, { headers: { Accept: "text/html" }, signal });
                if (!response.ok) throw new Error(`Registry returned ${response.status}.`);
                if (new URL(response.url).origin !== this.baseURL.origin) throw new RegistryPayloadError("The documentation request redirected outside the Registry.");
                documentation = await response.text();
            } catch (error) {
                if (error.name === "AbortError") throw error;
                documentationError = "The documentation artifact could not be loaded for this release.";
            }
        } else {
            documentationError = "This release does not publish browser-ready documentation.";
        }

        this.renderModuleView(route, module, version, documentation, documentationError);
    }

    renderModuleView(route, module, version, documentation, documentationError) {
        this.setBusy(false);
        const status = versionStatus(module, version.version, route.requestedVersion);
        const installCommand = packageInstallCommand(version);
        const repository = externalHTTPSURL(version.source?.repository);
        const links = Object.entries(version.links || {}).filter(([, href]) => externalHTTPSURL(href));
        const homepage = links.find(([name]) => name.toLowerCase() === "homepage");
        const otherLinks = links.filter(([name]) => name.toLowerCase() !== "homepage");

        this.root.innerHTML = `
            <article class="registry-module">
                <nav class="registry-breadcrumbs" aria-label="Breadcrumb">
                    <a href="/registry/">Registry</a><span aria-hidden="true">/</span><a href="${escapeHTML(ownerPath(route.owner))}">${escapeHTML(route.owner)}</a><span aria-hidden="true">/</span><span aria-current="page">${escapeHTML(route.name)}</span>
                </nav>
                <header class="registry-module-header">
                    <div>
                        <p class="registry-section-label"><a href="${escapeHTML(ownerPath(route.owner))}">${escapeHTML(route.owner)}</a> / module</p>
                        <h2>${escapeHTML(module.id)}</h2>
                        <p>${escapeHTML(version.description || module.description || "No description is available.")}</p>
                    </div>
                    <label class="registry-version-picker">
                        <span>Version</span>
                        <select id="registry-version" class="select-input">
                            ${module.versions.map((entry) => `<option value="${escapeHTML(entry.version)}"${entry.version === version.version ? " selected" : ""}>${escapeHTML(entry.version)}</option>`).join("")}
                        </select>
                    </label>
                </header>

                <div class="registry-module-layout">
                    <div class="registry-documentation-column">
                        ${documentationError ? `<section class="registry-state registry-state-inline registry-state-warning"><h3>Documentation unavailable</h3><p>${escapeHTML(documentationError)}</p></section>` : `
                            <nav id="registry-toc" class="registry-toc" aria-label="On this page" hidden></nav>
                            <section id="registry-documentation" class="content registry-documentation">${documentation}</section>`}
                    </div>
                    <aside class="registry-metadata" aria-label="Module metadata">
                        <div class="registry-status registry-status-${status.tone}">${escapeHTML(status.label)}</div>
                        <dl>
                            <div><dt>Version</dt><dd>${escapeHTML(version.version)}</dd></div>
                            <div><dt>Namespace</dt><dd><code>${escapeHTML(version.namespace || "Unavailable")}</code></dd></div>
                            <div><dt>License</dt><dd>${escapeHTML(version.license || "Unavailable")}</dd></div>
                            ${version.ferret ? `<div><dt>Ferret</dt><dd>${escapeHTML(version.ferret)}</dd></div>` : ""}
                        </dl>
                        <div class="registry-metadata-links">
                            ${repository ? `<a href="${escapeHTML(repository)}">Repository <span aria-hidden="true">↗</span></a>` : ""}
                            ${homepage ? `<a href="${escapeHTML(externalHTTPSURL(homepage[1]))}">Homepage <span aria-hidden="true">↗</span></a>` : ""}
                            ${otherLinks.map(([name, href]) => `<a href="${escapeHTML(externalHTTPSURL(href))}">${escapeHTML(name)} <span aria-hidden="true">↗</span></a>`).join("")}
                        </div>
                        <section class="registry-install" aria-labelledby="registry-install-title">
                            <h3 id="registry-install-title">Install</h3>
                            ${installCommand ? `<div class="registry-command"><code role="region" aria-label="Scrollable install command" tabindex="0">${escapeHTML(installCommand)}</code><button id="registry-copy" type="button" aria-describedby="registry-copy-feedback">Copy</button></div><p id="registry-copy-feedback" class="registry-copy-feedback" role="status"></p>` : `<p>Installation metadata is unavailable for this release.</p>`}
                        </section>
                    </aside>
                </div>
            </article>`;

        this.root.querySelector("#registry-version")?.addEventListener("change", (event) => {
            const nextPath = modulePath(module.id, event.target.value);
            window.history.pushState({}, "", nextPath);
            this.render();
        });

        const installCommandElement = this.root.querySelector(".registry-command code");
        installCommandElement?.addEventListener("keydown", (event) => {
            if (installCommandElement.scrollWidth <= installCommandElement.clientWidth) return;

            switch (event.key) {
                case "ArrowLeft":
                    installCommandElement.scrollLeft -= 40;
                    break;
                case "ArrowRight":
                    installCommandElement.scrollLeft += 40;
                    break;
                case "Home":
                    installCommandElement.scrollLeft = 0;
                    break;
                case "End":
                    installCommandElement.scrollLeft = installCommandElement.scrollWidth;
                    break;
                default:
                    return;
            }

            event.preventDefault();
        });

        const copyButton = this.root.querySelector("#registry-copy");
        copyButton?.addEventListener("click", async () => {
            const feedback = this.root.querySelector("#registry-copy-feedback");
            try {
                await copyText(installCommand);
                copyButton.textContent = "Copied";
                feedback.textContent = "Install command copied to the clipboard.";
            } catch {
                feedback.textContent = "Could not copy automatically. Select the command and copy it manually.";
            }
            window.setTimeout(() => {
                copyButton.textContent = "Copy";
            }, 2000);
        });

        this.wrapTables();
        this.buildTOC();
        this.scrollToFragment();
    }

    wrapTables() {
        const documentation = this.root.querySelector("#registry-documentation");
        if (!documentation) return;

        for (const table of documentation.querySelectorAll("table")) {
            if (table.parentElement?.classList.contains("table-scroll")) continue;

            const wrapper = document.createElement("div");
            wrapper.className = "table-scroll";
            wrapper.setAttribute("aria-label", "Scrollable table");
            wrapper.setAttribute("role", "region");
            wrapper.tabIndex = 0;
            table.before(wrapper);
            wrapper.append(table);
        }
    }

    buildTOC() {
        const documentation = this.root.querySelector("#registry-documentation");
        const toc = this.root.querySelector("#registry-toc");
        if (!documentation || !toc) return;
        const headings = [...documentation.querySelectorAll("h2[id], h3[id]")].filter((heading) => heading.textContent.trim());
        if (headings.length < 2) return;
        toc.innerHTML = `<p>On this page</p><ol>${headings.map((heading) => `<li class="registry-toc-${heading.tagName.toLowerCase()}"><a href="#${escapeHTML(heading.id)}">${escapeHTML(heading.textContent.trim())}</a></li>`).join("")}</ol>`;
        toc.hidden = false;
    }

    scrollToFragment() {
        const id = fragmentID(window.location.hash);
        if (!id) return;

        const documentation = this.root.querySelector("#registry-documentation");
        const target = document.getElementById(id);
        if (!documentation || !target || !documentation.contains(target)) return;

        target.scrollIntoView({ block: "start" });
    }
}

async function copyText(value) {
    if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(value);
        return;
    }
    const input = document.createElement("textarea");
    input.value = value;
    input.setAttribute("readonly", "");
    input.style.position = "fixed";
    input.style.opacity = "0";
    document.body.append(input);
    input.select();
    const copied = document.execCommand("copy");
    input.remove();
    if (!copied) throw new Error("Copy failed");
}

function validModuleID(id) {
    const parts = typeof id === "string" ? id.split("/") : [];
    return parts.length === 2 && parts.every((part) => ROUTE_SEGMENT.test(part));
}

if (typeof document !== "undefined") {
    const root = document.getElementById("registry-app");
    if (root) new RegistryApp(root).start();
}
