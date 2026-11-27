import { canInspect, createPlaygroundClient, playgroundErrorMessage } from './playground-api.mjs';
import {
    applyFormattedSource, createRequestState, createShareURL, diagnosticRange,
    parameterRows, readShareURL, serializeParameters,
} from './playground-state.mjs';

const analysisDelay = 350;
const sourceName = 'query.fql';
const operationLabels = { execute: 'Run', format: 'Format', compile: 'Bytecode' };
const operationFailures = {
    execute: 'The query could not run.',
    format: 'The source could not be formatted.',
    compile: 'The source could not be compiled for bytecode inspection.',
};

async function copyText(text) {
    if (navigator.clipboard?.writeText) {
        try {
            await navigator.clipboard.writeText(text);
            return;
        } catch {
            // Use the selection-based fallback when clipboard permission is unavailable.
        }
    }
    const focused = document.activeElement;
    const field = document.createElement('textarea');
    field.value = text;
    field.className = 'playground-clipboard';
    field.readOnly = true;
    field.tabIndex = -1;
    document.body.append(field);
    try {
        field.select();
        if (!document.execCommand('copy')) throw new Error('Clipboard unavailable');
    } finally {
        field.remove();
        focused?.focus({ preventScroll: true });
    }
}

function describeTimings(timings) {
    return [['compileMs', 'Compile'], ['executeMs', 'Execute'], ['totalMs', 'Total']]
        .filter(([key]) => Number.isFinite(timings?.[key]))
        .map(([key, label]) => `${label} ${timings[key].toLocaleString(undefined, { maximumFractionDigits: 2 })}\u00a0ms`);
}

function bindExamplesMenu(trigger, menu, selectExample) {
    const items = Array.from(menu.querySelectorAll('[role="menuitem"]'));
    const viewport = window.visualViewport;
    let openListeners;

    function positionMenu() {
        const gutter = 8;
        const viewportLeft = viewport?.offsetLeft ?? 0;
        const viewportTop = viewport?.offsetTop ?? 0;
        const width = viewport?.width ?? document.documentElement.clientWidth;
        const height = viewport?.height ?? document.documentElement.clientHeight;
        const bottom = viewportTop + height - gutter;
        const anchor = trigger.getBoundingClientRect();
        const top = anchor.bottom + 4;

        menu.style.maxWidth = `${Math.max(0, width - 2 * gutter)}px`;
        const left = Math.min(anchor.left, viewportLeft + width - menu.offsetWidth - gutter);
        menu.style.left = `${Math.max(viewportLeft + gutter, left)}px`;
        menu.style.top = `${top}px`;
        menu.style.setProperty('--examples-available-height', `${Math.max(0, bottom - top)}px`);
        return bottom;
    }

    function focusItem(index) {
        const item = items[index];
        item.focus({ preventScroll: true });
        const top = item.offsetTop;
        const bottom = top + item.offsetHeight;
        if (top < menu.scrollTop) {
            menu.scrollTop = top;
        } else if (bottom > menu.scrollTop + menu.clientHeight) {
            menu.scrollTop = bottom - menu.clientHeight;
        }
    }

    function close(restoreFocus = false) {
        menu.hidden = true;
        trigger.setAttribute('aria-expanded', 'false');
        openListeners.abort();
        if (restoreFocus) trigger.focus({ preventScroll: true });
    }

    function dismissOutside(event) {
        if (!trigger.contains(event.target) && !menu.contains(event.target)) close();
    }

    function open() {
        menu.hidden = false;
        trigger.setAttribute('aria-expanded', 'true');
        menu.scrollTop = 0;
        const viewportBottom = positionMenu();
        const clipped = Math.ceil(items[0].getBoundingClientRect().bottom + menu.clientTop - viewportBottom);
        if (clipped > 0) {
            window.scrollBy({ top: clipped, behavior: 'instant' });
            positionMenu();
        }
        openListeners = new AbortController();
        const { signal } = openListeners;
        document.addEventListener('pointerdown', dismissOutside, { capture: true, signal });
        document.addEventListener('focusin', dismissOutside, { capture: true, signal });
        window.addEventListener('scroll', positionMenu, { passive: true, signal });
        window.addEventListener('resize', positionMenu, { signal });
        viewport?.addEventListener('scroll', positionMenu, { passive: true, signal });
        viewport?.addEventListener('resize', positionMenu, { signal });
        focusItem(0);
    }

    trigger.addEventListener('click', () => {
        if (menu.hidden) open();
        else close(true);
    });
    items.forEach(item => {
        item.addEventListener('click', () => {
            close(true);
            selectExample(Number(item.dataset.exampleIndex));
        });
    });
    menu.addEventListener('keydown', event => {
        if (event.key === 'Tab') {
            // Let the browser move focus before hiding the item it starts from.
            const { signal } = openListeners;
            setTimeout(() => {
                if (!signal.aborted) close();
            }, 0);
            return;
        }
        if (event.key === 'Escape') {
            event.preventDefault();
            event.stopPropagation();
            close(true);
            return;
        }

        const index = items.indexOf(event.target);
        const next = {
            ArrowDown: (index + 1) % items.length,
            ArrowUp: (index + items.length - 1) % items.length,
            Home: 0,
            End: items.length - 1,
        };
        if (!Object.hasOwn(next, event.key)) return;
        event.preventDefault();
        focusItem(next[event.key]);
    });
}

export function initPlayground(root) {
    const part = name => root.querySelector(`[data-part="${name}"]`);
    const sourceElement = part('editor');
    const editor = ace.edit(sourceElement);
    const { Range } = ace.require('ace/range');
    const readonly = root.dataset.readonly === 'true';
    editor.setTheme(root.dataset.theme);
    editor.session.setUseWorker(false);
    editor.session.setMode('ace/mode/fql');
    editor.setOptions({ displayIndentGuides: false, showPrintMargin: false, tabSize: 4, useSoftTabs: true });
    editor.setReadOnly(readonly);
    editor.textInput.getElement().setAttribute('aria-label', 'FQL source editor');
    new ResizeObserver(() => editor.resize()).observe(sourceElement);

    const examples = JSON.parse(root.dataset.examples);
    const client = createPlaygroundClient(root.dataset.url);
    const requests = createRequestState();
    const outputs = {};
    const problemSets = new Map();
    const problemRows = new Map();
    const tabs = ['result', 'bytecode', 'problems'];
    let environment;
    let analysisTimer;
    let analysisStatus = 'Checking source…';

    function setNotice(message) {
        part('notice').textContent = message;
        part('share-link').hidden = true;
    }

    function renderControls() {
        const running = requests.pending('execute');
        part('run').disabled = !environment?.capabilities.execute;
        part('run').setAttribute('aria-disabled', String(running || part('run').disabled));
        const format = part('format');
        if (format) {
            format.disabled = !environment?.capabilities.format;
            format.setAttribute('aria-disabled', String(requests.pending('format') || format.disabled));
        }
        part('inspect').disabled = !canInspect(environment);
        part('inspect').setAttribute('aria-disabled', String(requests.pending('compile') || part('inspect').disabled));
    }

    function renderOutput(kind) {
        const result = kind === 'result';
        const pending = requests.pending(result ? 'execute' : 'compile');
        const output = outputs[kind];
        const stale = output && !requests.matches(output.revision);
        const summary = [];
        if (pending) summary.push(result ? 'Running query…' : 'Inspecting source…');
        if (stale) summary.push(result ? 'From earlier inputs' : 'From earlier source');
        else if (output?.previous) summary.push(result ? 'Previous output' : 'Previous bytecode');
        else if (output && !pending) summary.push(...describeTimings(output.timings));
        const meta = part(`${kind}-meta`);
        const label = summary.join(' · ');
        if (meta.textContent !== label) meta.textContent = label;

        const text = part(`${kind}-text`);
        const empty = part(`${kind}-empty`);
        const hasText = output?.text !== undefined;
        text.hidden = !hasText;
        empty.hidden = hasText;
        part(`${kind}-copy`).disabled = !hasText;
        if (hasText && text.textContent !== output.text) {
            text.textContent = output.text;
            text.scrollTop = 0;
            text.scrollLeft = 0;
        }
        if (!hasText) text.textContent = '';
        if (hasText) return;

        if (output && !stale && !pending) {
            empty.textContent = output.message;
        } else if (result) {
            empty.textContent = environment && !environment.capabilities.execute
                ? 'This runtime does not support execution.'
                : 'Run a query to see its JSON result.';
        } else {
            empty.textContent = environment && !canInspect(environment)
                ? 'This runtime does not provide bytecode and parameter inspection.'
                : 'Inspect the source to see its bytecode.';
        }
    }

    function renderProblems() {
        const problems = new Map();
        for (const { revision, entries } of problemSets.values()) {
            if (!requests.matches(revision)) continue;
            for (const entry of entries) {
                const { severity, message, hint, phase, operation, code } = entry;
                const range = diagnosticRange(entry);
                const key = JSON.stringify([severity, message, hint || '', phase, operation, code, range]);
                problems.set(key, { ...entry, range });
            }
        }

        const list = part('problems-list');
        let annotationsChanged = false;
        for (const [key, { row, marker, annotation }] of problemRows) {
            if (problems.has(key)) continue;
            row.remove();
            if (marker !== undefined) editor.session.removeMarker(marker);
            if (annotation) annotationsChanged = true;
            problemRows.delete(key);
        }

        for (const [key, { severity, message, hint, phase, operation, range }] of problems) {
            if (problemRows.has(key)) continue;
            const row = document.createElement('li');
            const label = document.createElement('span');
            label.className = 'playground-diagnostic-label';
            label.dataset.severity = severity;
            label.textContent = [severity, operationLabels[operation] ?? phase].filter(Boolean).join(' · ');
            row.append(label);
            let annotation;
            let marker;
            if (range) {
                const { start, end } = range;
                annotation = { ...start, text: message, type: severity };
                annotationsChanged = true;
                if (start.row !== end.row || start.column !== end.column) {
                    marker = editor.session.addMarker(Range.fromPoints(start, end),
                        `playground-diagnostic-marker playground-diagnostic-marker-${severity}`, 'text');
                }
                const button = document.createElement('button');
                button.type = 'button';
                button.className = 'playground-diagnostic-link';
                button.textContent = `${start.row + 1}:${start.column + 1} ${message}`;
                button.addEventListener('click', () => {
                    editor.focus();
                    editor.selection.setSelectionRange(Range.fromPoints(start, end));
                    editor.renderer.scrollCursorIntoView();
                });
                row.append(button);
            } else {
                row.append(document.createTextNode(message));
            }
            if (hint) {
                const details = document.createElement('p');
                details.textContent = hint;
                row.append(details);
            }
            list.append(row);
            problemRows.set(key, { row, marker, annotation });
        }
        if (annotationsChanged) {
            editor.session.setAnnotations(Array.from(problemRows.values(), problem => problem.annotation).filter(Boolean));
        }
        const count = problems.size;
        list.hidden = count === 0;
        const status = part('analysis-status');
        const message = analysisStatus || (count ? '' : 'No problems detected');
        if (status.textContent !== message) status.textContent = message;
        status.hidden = !message;
        const tab = part('problems-tab');
        const title = count ? `Problems (${count})` : 'Problems';
        if (tab.textContent !== title) tab.textContent = title;
        const summary = part('problems-status');
        const label = count === 1 ? '1 problem' : count ? `${count} problems` : '';
        if (summary.textContent !== label) summary.textContent = label;
    }

    function rows() {
        return Array.from(part('params-rows').children, row => Array.from(row.querySelectorAll('input, textarea')));
    }

    function setDock(open) {
        part('dock').hidden = !open;
        part('params').setAttribute('aria-expanded', String(open));
        part('grid').classList.toggle('has-params', open);
    }

    function renderSuggestions() {
        const container = part('params-suggestions');
        const bytecode = outputs.bytecode;
        const names = new Set(rows().map(([name]) => name.value.trim()));
        const suggestions = bytecode && !bytecode.previous && requests.matches(bytecode.revision)
            ? bytecode.parameters?.filter(name => !names.has(name)) ?? []
            : [];
        container.replaceChildren();
        container.hidden = suggestions.length === 0;
        if (!suggestions.length) return;

        const heading = document.createElement('p');
        heading.textContent = 'Referenced parameters';
        container.append(heading);
        for (const name of suggestions) {
            const button = document.createElement('button');
            button.type = 'button';
            button.className = 'playground-button';
            button.textContent = `Add @${name}`;
            button.addEventListener('click', () => {
                const [, value] = addRow(name);
                parametersChanged();
                value.focus();
            });
            container.append(button);
        }
    }

    function parametersChanged() {
        requests.changeParams();
        const runProblems = problemSets.get('execute');
        if (runProblems && !requests.matches(runProblems.revision)) problemSets.delete('execute');
        const error = part('params-error');
        error.hidden = true;
        part('params-rows').querySelector('[aria-invalid]')?.removeAttribute('aria-invalid');
        part('params-empty').hidden = part('params-rows').childElementCount > 0;
        renderSuggestions();
        renderProblems();
        renderOutput('result');
    }

    function addRow(name = '', value = '') {
        const row = document.createElement('div');
        row.className = 'playground-param-row';
        const fields = [name, value].map((initial, index) => {
            const field = document.createElement(index === 0 ? 'input' : 'textarea');
            const label = index === 0 ? 'Name' : 'Value';
            if (index === 0) field.type = 'text';
            else field.rows = 2;
            field.value = initial;
            field.placeholder = label;
            field.setAttribute('aria-label', `Parameter ${label.toLowerCase()}`);
            field.autocomplete = 'off';
            field.spellcheck = false;
            field.addEventListener('input', parametersChanged);
            return field;
        });
        const remove = document.createElement('button');
        remove.type = 'button';
        remove.className = 'playground-text-button';
        remove.textContent = 'Remove';
        const labelRemove = () => {
            const name = fields[0].value.trim();
            remove.setAttribute('aria-label', name ? `Remove parameter ${name}` : 'Remove parameter row');
        };
        labelRemove();
        fields[0].addEventListener('input', labelRemove);
        remove.addEventListener('click', () => {
            const next = (row.nextElementSibling ?? row.previousElementSibling)?.querySelector('input');
            row.remove();
            parametersChanged();
            (next ?? part('params-add')).focus();
        });
        row.append(...fields, remove);
        part('params-rows').append(row);
        return fields;
    }

    function loadParameters(entries) {
        part('params-rows').replaceChildren();
        for (const [name, value] of entries) addRow(name, value);
        parametersChanged();
    }

    function currentParameters(revealError = true) {
        const fields = rows();
        const params = serializeParameters(fields.map(pair => pair.map(field => field.value)));
        if (params.error) {
            const field = fields[params.error.index][0];
            const error = part('params-error');
            error.textContent = params.error.message;
            error.hidden = false;
            field.setAttribute('aria-invalid', 'true');
            if (revealError) {
                setDock(true);
                field.focus();
            }
        }
        return params;
    }

    function selectTab(kind) {
        for (const tab of tabs) {
            const button = part(`${tab}-tab`);
            button.setAttribute('aria-selected', String(tab === kind));
            button.tabIndex = tab === kind ? 0 : -1;
            part(`${tab}-panel`).hidden = tab !== kind;
        }
    }

    function operationFailed(operation) {
        const output = operation === 'execute' ? outputs.result : outputs.bytecode;
        if (operation === 'format') setNotice('');
        else if (output) output.previous = true;
        selectTab('problems');
    }

    async function operate(operation, onSuccess, paramsJSON) {
        const request = requests.start(operation);
        if (!request) return;
        const source = { text: editor.getValue(), name: sourceName };
        renderControls();
        if (operation === 'execute') renderOutput('result');
        if (operation === 'compile') renderOutput('bytecode');

        let response;
        let requestError;
        try {
            response = await client.request(operation, { source, paramsJSON, signal: request.controller.signal });
        } catch (error) {
            requestError = playgroundErrorMessage(error);
        }

        const applies = requests.canApply(request);
        requests.finish(request);
        if (applies) {
            const failed = response?.status !== 'ok';
            let entries = response?.diagnostics ?? [];
            const revision = operation === 'execute' && response?.status === 'invalid'
                ? { source: request.revision.source } : request.revision;
            if (operation === 'analyze') {
                analysisStatus = requestError ? `Analysis unavailable. ${requestError}` : '';
                if (failed && !entries.length && !analysisStatus) analysisStatus = 'Analysis could not check this source.';
            } else if (failed) {
                if (!entries.length) {
                    entries = [{ severity: 'error', operation, message: requestError ?? operationFailures[operation] }];
                }
                operationFailed(operation);
            }
            problemSets.set(operation, { revision, entries });
            if (!failed) onSuccess?.(response, request.revision);
            renderProblems();
        }

        renderControls();
        if (operation === 'execute') renderOutput('result');
        if (operation === 'compile') renderOutput('bytecode');
        if (operation === 'compile' && applies) renderSuggestions();
    }

    function scheduleAnalysis() {
        clearTimeout(analysisTimer);
        if (environment?.capabilities.analyze) {
            analysisStatus = 'Checking source…';
            analysisTimer = setTimeout(() => operate('analyze'), analysisDelay);
        } else if (environment) {
            analysisStatus = 'This runtime does not support source analysis.';
        }
        renderProblems();
    }

    function formatSource() {
        if (!environment?.capabilities.format || readonly || requests.pending('format')) return;
        setNotice('Formatting source…');
        void operate('format', response => {
            const changed = applyFormattedSource(editor, response);
            setNotice(changed ? 'Source formatted.' : 'Source is already formatted.');
        });
    }

    function inspectSource() {
        if (!canInspect(environment) || requests.pending('compile')) return;
        void operate('compile', (response, revision) => {
            const { bytecode, parameters } = response.artifacts;
            outputs.bytecode = {
                text: bytecode.text,
                parameters,
                timings: response.timings,
                revision,
            };
        });
    }

    function run() {
        if (!environment?.capabilities.execute || requests.pending('execute')) return;
        const params = currentParameters(false);
        if (params.error) {
            problemSets.set('execute', {
                revision: requests.revision('execute'),
                entries: [{ severity: 'error', operation: 'execute', message: params.error.message }],
            });
            operationFailed('execute');
            renderProblems();
            renderOutput('result');
            return;
        }
        void operate('execute', (response, revision) => {
            outputs.result = {
                text: response.outputText,
                message: 'The query completed without a result.',
                timings: response.timings,
                revision,
            };
            selectTab('result');
        }, params.json);
    }

    async function discover() {
        const connection = part('connection');
        const retry = part('retry');
        retry.hidden = true;
        connection.textContent = 'Connecting to Playground…';
        analysisStatus = 'Checking source…';
        renderProblems();
        try {
            environment = await client.request('environment');
            connection.textContent = environment.runtime.version ? `Ferret ${environment.runtime.version}` : 'Ferret Playground';
            const { maxTimeoutMs, maxOutputBytes } = environment.limits.execution;
            connection.title = `Execution limit ${maxTimeoutMs.toLocaleString()} ms · Output limit ${maxOutputBytes.toLocaleString()} bytes`;
            renderOutput('result');
            renderOutput('bytecode');
            scheduleAnalysis();
        } catch (error) {
            const message = playgroundErrorMessage(error);
            connection.textContent = message;
            retry.hidden = false;
            analysisStatus = `Analysis unavailable. ${message}`;
            renderProblems();
        }
        renderControls();
    }

    async function share() {
        const params = currentParameters();
        if (params.error) return;
        let link;
        try {
            link = createShareURL(location.href, editor.getValue(), params.entries);
        } catch (error) {
            setNotice(error.message);
            return;
        }
        try {
            await copyText(link);
            setNotice('Share link copied.');
        } catch {
            setNotice('Copy the share link from the field below.');
            const field = part('share-link');
            field.value = link;
            field.hidden = false;
            field.select();
        }
    }

    function bindCopy(name, label, value) {
        part(name)?.addEventListener('click', async () => {
            try {
                await copyText(value());
                setNotice(`${label} copied.`);
            } catch {
                setNotice(`Could not copy ${label.toLowerCase()}.`);
            }
        });
    }

    function sharedProgram() {
        try {
            return readShareURL(location.href);
        } catch (error) {
            setNotice(error.message);
            return null;
        }
    }

    const shared = part('share') ? sharedProgram() : null;
    if (shared) editor.setValue(shared.source, -1);
    loadParameters(shared?.entries ?? parameterRows(JSON.parse(root.dataset.params || '{}')));
    editor.session.on('change', () => {
        requests.changeSource();
        problemSets.clear();
        setNotice('');
        renderSuggestions();
        renderControls();
        renderOutput('result');
        renderOutput('bytecode');
        scheduleAnalysis();
    });
    editor.commands.addCommand({
        name: 'runPlaygroundQuery',
        bindKey: { win: 'Ctrl-Enter', mac: 'Command-Enter|Ctrl-Enter' },
        readOnly: true,
        exec: run,
    });

    const examplesTrigger = part('examples');
    if (examplesTrigger) {
        bindExamplesMenu(examplesTrigger, part('examples-menu'), index => {
            const example = examples[index];
            editor.setValue(example.source, -1);
            loadParameters(parameterRows(example.params));
        });
    }
    part('run').addEventListener('click', run);
    part('format')?.addEventListener('click', formatSource);
    part('inspect').addEventListener('click', inspectSource);
    part('retry').addEventListener('click', discover);
    part('share')?.addEventListener('click', share);
    if (part('share')) {
        window.addEventListener('hashchange', () => {
            if (!location.hash.startsWith('#fql=')) return;
            const program = sharedProgram();
            if (!program) return;
            editor.setValue(program.source, -1);
            loadParameters(program.entries);
        });
    }
    part('params').addEventListener('click', () => setDock(part('dock').hidden));
    part('params-close').addEventListener('click', () => {
        setDock(false);
        part('params').focus();
    });
    part('params-add').addEventListener('click', () => {
        const [name] = addRow();
        parametersChanged();
        name.focus();
    });
    bindCopy('source-copy', 'Source', () => editor.getValue());
    bindCopy('result-copy', 'Result', () => outputs.result.text);
    bindCopy('bytecode-copy', 'Bytecode', () => outputs.bytecode.text);

    tabs.forEach((kind, index) => {
        const button = part(`${kind}-tab`);
        button.addEventListener('click', () => {
            selectTab(kind);
            if (kind === 'bytecode' && !problemSets.has('compile')) inspectSource();
        });
        button.addEventListener('keydown', event => {
            const neighbors = {
                ArrowLeft: (index + tabs.length - 1) % tabs.length,
                ArrowRight: (index + 1) % tabs.length,
                Home: 0,
                End: tabs.length - 1,
            };
            if (!Object.hasOwn(neighbors, event.key)) return;
            event.preventDefault();
            part(`${tabs[neighbors[event.key]]}-tab`).focus();
        });
    });

    void discover();
}

if (globalThis.window?.ace) document.querySelectorAll('[data-playground]').forEach(initPlayground);
