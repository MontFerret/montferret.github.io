import { canInspect, createPlaygroundClient, playgroundErrorMessage } from './playground-api.mjs';
import {
    applyFormattedSource, createRequestState, createShareURL, diagnosticRange,
    parameterRows, readShareURL, serializeParameters,
} from './playground-state.mjs';

const analysisDelay = 350;
const sourceName = 'query.fql';

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

function initPlayground(root) {
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
    const diagnostics = new Map();
    const tabs = ['result', 'bytecode'];
    let environment;
    let analysisTimer;
    let analysisError = '';

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
        if (output && !stale && !pending) summary.push(...describeTimings(output.timings));
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

    function renderDiagnostics() {
        const list = part('diagnostics-list');
        list.replaceChildren();
        const seen = new Set();
        const annotations = [];
        for (const diagnostic of Array.from(diagnostics.values()).flat()) {
            const { severity, message, hint, phase } = diagnostic;
            const range = diagnosticRange(diagnostic);
            const key = JSON.stringify([severity, message, hint, range]);
            if (seen.has(key)) continue;
            seen.add(key);

            const row = document.createElement('li');
            const label = document.createElement('span');
            label.className = 'playground-diagnostic-label';
            label.dataset.severity = severity;
            label.textContent = [severity, phase].filter(Boolean).join(' · ');
            row.append(label);
            if (range) {
                annotations.push({ ...range.start, text: message, type: severity });
                const { start, end } = range;
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
        }
        if (analysisError) {
            const row = document.createElement('li');
            row.textContent = `Analysis unavailable. ${analysisError}`;
            list.append(row);
        }
        editor.session.setAnnotations(annotations);
        part('diagnostics').hidden = !list.childElementCount;
        const summary = part('diagnostics-status');
        const count = seen.size;
        const label = count === 1 ? '1 diagnostic' : count ? `${count} diagnostics` : '';
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
        const suggestions = bytecode && requests.matches(bytecode.revision)
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
        diagnostics.delete('execute');
        const error = part('params-error');
        error.hidden = true;
        part('params-rows').querySelector('[aria-invalid]')?.removeAttribute('aria-invalid');
        part('params-empty').hidden = part('params-rows').childElementCount > 0;
        renderSuggestions();
        renderDiagnostics();
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

    function currentParameters() {
        const fields = rows();
        const params = serializeParameters(fields.map(pair => pair.map(field => field.value)));
        if (params.error) {
            const field = fields[params.error.index][0];
            const error = part('params-error');
            setDock(true);
            error.textContent = params.error.message;
            error.hidden = false;
            field.setAttribute('aria-invalid', 'true');
            field.focus();
            return null;
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

    async function operate(operation, onOutcome, onError, paramsJSON) {
        const request = requests.start(operation);
        if (!request) return;
        const source = { text: editor.getValue(), name: sourceName };
        renderControls();
        if (operation === 'execute') renderOutput('result');
        if (operation === 'compile') renderOutput('bytecode');
        try {
            const response = await client.request(operation, { source, paramsJSON, signal: request.controller.signal });
            if (!requests.canApply(request)) return;
            diagnostics.set(operation, response.diagnostics);
            onOutcome(response, request.revision);
            renderDiagnostics();
        } catch (error) {
            if (!requests.canApply(request)) return;
            diagnostics.delete(operation);
            onError(playgroundErrorMessage(error), request.revision);
            renderDiagnostics();
        } finally {
            requests.finish(request);
            renderControls();
            if (operation === 'execute') renderOutput('result');
            if (operation === 'compile') {
                renderOutput('bytecode');
                renderSuggestions();
            }
        }
    }

    function scheduleAnalysis() {
        clearTimeout(analysisTimer);
        if (!environment?.capabilities.analyze) return;
        analysisTimer = setTimeout(() => operate('analyze',
            () => { analysisError = ''; },
            message => { analysisError = message; },
        ), analysisDelay);
    }

    function formatSource() {
        if (!environment?.capabilities.format || readonly || requests.pending('format')) return;
        setNotice('Formatting source…');
        void operate('format', response => {
            if (response.status === 'invalid') {
                setNotice('Check the diagnostics before formatting.');
            } else {
                const changed = applyFormattedSource(editor, response);
                setNotice(changed ? 'Source formatted.' : 'Source is already formatted.');
            }
        }, setNotice);
    }

    function inspectSource() {
        if (!canInspect(environment) || requests.pending('compile')) return;
        void operate('compile', (response, revision) => {
            const artifacts = response.status === 'ok' ? response.artifacts : null;
            outputs.bytecode = {
                text: artifacts?.bytecode.text,
                parameters: artifacts?.parameters,
                message: 'Check the diagnostics before inspecting bytecode.',
                timings: response.timings,
                revision,
            };
        }, (message, revision) => { outputs.bytecode = { message, revision }; });
    }

    function run() {
        if (!environment?.capabilities.execute || requests.pending('execute')) return;
        const params = currentParameters();
        if (!params) return;
        selectTab('result');
        void operate('execute', (response, revision) => {
            const failure = response.diagnostics[0];
            const message = {
                ok: 'The query completed without a result.',
                invalid: 'Check the diagnostics before running the query.',
                failed: playgroundErrorMessage(failure ?? {}),
            }[response.status];
            outputs.result = { text: response.outputText, message, timings: response.timings, revision };
        }, (message, revision) => { outputs.result = { message, revision }; }, params.json);
    }

    async function discover() {
        const connection = part('connection');
        const retry = part('retry');
        retry.hidden = true;
        connection.textContent = 'Connecting to Playground…';
        try {
            environment = await client.request('environment');
            connection.textContent = environment.runtime.version ? `Ferret ${environment.runtime.version}` : 'Ferret Playground';
            const { maxTimeoutMs, maxOutputBytes } = environment.limits.execution;
            connection.title = `Execution limit ${maxTimeoutMs.toLocaleString()} ms · Output limit ${maxOutputBytes.toLocaleString()} bytes`;
            renderOutput('result');
            renderOutput('bytecode');
            scheduleAnalysis();
        } catch (error) {
            connection.textContent = playgroundErrorMessage(error);
            retry.hidden = false;
        }
        renderControls();
    }

    async function share() {
        const params = currentParameters();
        if (!params) return;
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
        diagnostics.clear();
        analysisError = '';
        setNotice('');
        renderDiagnostics();
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

    part('examples')?.addEventListener('change', event => {
        const example = event.target.value === '' ? null : examples[Number(event.target.value)];
        if (!example) return;
        editor.setValue(example.source, -1);
        loadParameters(parameterRows(example.params));
        event.target.value = '';
    });
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
            if (kind === 'bytecode' && (!outputs.bytecode || !requests.matches(outputs.bytecode.revision))) inspectSource();
        });
        button.addEventListener('keydown', event => {
            const neighbors = { ArrowLeft: 1 - index, ArrowRight: 1 - index, Home: 0, End: 1 };
            if (!Object.hasOwn(neighbors, event.key)) return;
            event.preventDefault();
            part(`${tabs[neighbors[event.key]]}-tab`).focus();
        });
    });

    void discover();
}

if (window.ace) document.querySelectorAll('[data-playground]').forEach(initPlayground);
