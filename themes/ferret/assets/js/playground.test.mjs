import assert from 'node:assert/strict';
import { test } from 'node:test';
import { canInspect, createPlaygroundClient, playgroundErrorMessage, readPlaygroundResponse } from './playground-api.mjs';
import { initPlayground } from './playground.mjs';
import {
    applyFormattedSource, createRequestState, createShareURL, diagnosticRange,
    parameterRows, readShareURL, serializeParameters,
} from './playground-state.mjs';

const origin = 'https://ferretlang.org/try/';
const runtime = {
    apiVersion: 'v1',
    runtime: { version: 'v2.0.0-alpha.62' },
    capabilities: { analyze: true, format: true, compile: true, execute: true, inspectionArtifacts: ['bytecode', 'parameters'] },
    limits: {
        tooling: { maxSourceBytes: 64 },
        http: { maxRequestBodyBytes: 256, writeTimeoutMs: 15000 },
        execution: { maxTimeoutMs: 5000, maxOutputBytes: 4096 },
    },
};

async function mountPlayground(t) {
    let focused;
    class Element extends EventTarget {
        children = [];
        attributes = new Map();
        dataset = {};
        textContent = '';
        classList = { toggle() {} };

        get childElementCount() { return this.children.length; }
        setAttribute(name, value) { this.attributes.set(name, value); }
        querySelector() { return null; }
        querySelectorAll() { return this.children.filter(child => ['input', 'textarea'].includes(child.tag)); }
        focus() { focused = this; }
        click() { this.dispatchEvent(new Event('click')); }
        append(...children) {
            for (const child of children) child.parent = this;
            this.children.push(...children);
        }
        replaceChildren() { this.children = []; }
        remove() { this.parent.children = this.parent.children.filter(child => child !== this); }
    }

    const parts = new Map();
    const part = name => {
        if (!parts.has(name)) parts.set(name, new Element());
        return parts.get(name);
    };
    const root = {
        dataset: { examples: '[]', url: 'https://playground.example/' },
        querySelector(selector) {
            const name = selector.match(/data-part="([^"]+)"/)[1];
            return ['examples', 'share', 'source-copy'].includes(name) ? null : part(name);
        },
    };
    part('result-tab').setAttribute('aria-selected', 'true');
    const editorInput = new Element();
    let source = 'RETURN 1';
    let changeSource;
    let run;
    let selectedRange;
    let annotations = [];
    let nextMarker = 0;
    const markers = new Map();
    const editor = {
        setTheme() {}, setOptions() {}, setReadOnly() {},
        getValue: () => source,
        focus: () => editorInput.focus(),
        textInput: { getElement: () => editorInput },
        renderer: { scrollCursorIntoView() {} },
        selection: { setSelectionRange(range) { selectedRange = range; } },
        commands: { addCommand(command) { run = command.exec; } },
        session: {
            setUseWorker() {}, setMode() {},
            on(_event, callback) { changeSource = callback; },
            setAnnotations(value) { annotations = value; },
            addMarker(range) { markers.set(++nextMarker, range); return nextMarker; },
            removeMarker(id) { markers.delete(id); },
        },
    };
    const document = {
        createElement(tag) { return Object.assign(new Element(), { tag }); },
        createTextNode(textContent) { return Object.assign(new Element(), { textContent }); },
    };
    const pending = [];
    t.mock.method(globalThis, 'fetch', url => {
        const request = Promise.withResolvers();
        pending.push({ operation: url.pathname.split('/').at(-1), ...request });
        return request.promise;
    });
    for (const [name, value] of Object.entries({
        document,
        ace: { edit: () => editor, require: () => ({ Range: { fromPoints: (start, end) => ({ start, end }) } }) },
        ResizeObserver: class { observe() {} },
    })) {
        const original = Object.getOwnPropertyDescriptor(globalThis, name);
        Object.defineProperty(globalThis, name, { value, configurable: true });
        t.after(() => {
            if (original) Object.defineProperty(globalThis, name, original);
            else delete globalThis[name];
        });
    }
    t.mock.timers.enable({ apis: ['setTimeout'] });

    async function respond(operation, body, status = 200) {
        const index = pending.findIndex(request => request.operation === operation);
        assert.ok(index >= 0, `${operation} request is pending`);
        const [request] = pending.splice(index, 1);
        request.resolve(Response.json(body, { status }));
        await new Promise(setImmediate);
    }

    initPlayground(root);
    await respond('environment', runtime);
    return {
        part, respond, editorInput, markers,
        analyze() { t.mock.timers.tick(350); },
        edit(text) { source = text; changeSource(); },
        run: () => run(),
        get focused() { return focused; },
        get selectedRange() { return selectedRange; },
        get annotations() { return annotations; },
        get tab() { return ['result', 'bytecode', 'problems'].find(kind => part(`${kind}-tab`).attributes.get('aria-selected') === 'true'); },
    };
}

const sourceProblem = {
    severity: 'error', phase: 'parse', message: 'Expected an expression', hint: 'Add a value after RETURN.',
    location: { range: { start: { line: 0, character: 0 }, end: { line: 0, character: 6 } } },
};
const invalidSource = { status: 'invalid', diagnostics: [sourceProblem] };
const validSource = { status: 'ok', diagnostics: [] };

test('background analysis preserves tab selection and focus while rendering source problems', async t => {
    const page = await mountPlayground(t);
    page.part('problems-tab').click();
    assert.match(page.part('analysis-status').textContent, /Checking source/);
    page.part('bytecode-tab').click();
    page.editorInput.focus();
    page.analyze();
    await page.respond('analyze', invalidSource);
    assert.equal(page.tab, 'bytecode');
    assert.equal(page.focused, page.editorInput);
    assert.equal(page.part('problems-tab').textContent, 'Problems (1)');
    assert.deepEqual(page.annotations, [{ row: 0, column: 0, text: sourceProblem.message, type: 'error' }]);
    assert.equal(page.markers.size, 1);

    page.part('problems-tab').click();
    const [label, link, hint] = page.part('problems-list').children[0].children;
    assert.equal(label.textContent, 'error · parse');
    assert.equal(link.textContent, `1:1 ${sourceProblem.message}`);
    assert.equal(hint.textContent, sourceProblem.hint);
    link.focus();
    link.click();
    assert.equal(page.focused, page.editorInput);
    assert.deepEqual(page.selectedRange, diagnosticRange(sourceProblem));

    page.edit('RETURN 2');
    assert.equal(page.markers.size, 0);
    assert.deepEqual(page.annotations, []);
    assert.match(page.part('analysis-status').textContent, /Checking source/);
    page.analyze();
    await page.respond('analyze', validSource);
    assert.equal(page.part('analysis-status').textContent, 'No problems detected');
    page.edit('RETURN 3');
    page.analyze();
    await page.respond('analyze', {}, 503);
    assert.match(page.part('analysis-status').textContent, /Analysis unavailable/);
    assert.equal(page.tab, 'problems');
});

test('Run selects the outcome tab and clean analysis keeps a current runtime failure and previous output', async t => {
    const page = await mountPlayground(t);
    const runtimeFailure = { severity: 'error', phase: 'execute', code: 'runtime_error', message: 'The program failed during execution.' };
    page.part('problems-tab').click();
    page.editorInput.focus();
    page.run();
    assert.equal(page.tab, 'problems');
    await page.respond('execute', { ...validSource, output: false });
    assert.equal(page.tab, 'result');
    assert.equal(page.part('result-text').textContent, 'false');
    page.run();
    await page.respond('execute', { status: 'failed', diagnostics: [runtimeFailure] });
    assert.equal(page.tab, 'problems');
    assert.equal(page.focused, page.editorInput);
    assert.equal(page.part('result-text').textContent, 'false');
    assert.equal(page.part('result-meta').textContent, 'Previous output');

    page.part('result-tab').click();
    page.analyze();
    await page.respond('analyze', validSource);
    assert.equal(page.tab, 'result');
    assert.equal(page.part('problems-tab').textContent, 'Problems (1)');
    assert.equal(page.part('problems-list').children[0].children[1].textContent, runtimeFailure.message);
    assert.deepEqual(page.annotations, []);
    assert.equal(page.markers.size, 0);

    page.part('params-add').click();
    assert.equal(page.part('problems-tab').textContent, 'Problems');
    assert.equal(page.part('result-meta').textContent, 'From earlier inputs');
    assert.equal(page.part('result-text').textContent, 'false');
});

test('explicit failures select Problems and duplicate compiler diagnostics retain their focused link', async t => {
    const page = await mountPlayground(t);
    page.run();
    await page.respond('execute', invalidSource);
    assert.equal(page.tab, 'problems');
    const link = page.part('problems-list').children[0].children[1];
    page.part('params-add').click();
    assert.equal(page.part('problems-list').children[0].children[1], link);
    link.focus();
    page.analyze();
    await page.respond('analyze', invalidSource);
    assert.equal(page.part('problems-tab').textContent, 'Problems (1)');
    assert.equal(page.part('problems-list').children[0].children[1], link);
    assert.equal(page.focused, link);
    assert.equal(page.markers.size, 1);

    page.part('result-tab').click();
    page.part('format').click();
    await page.respond('format', invalidSource);
    assert.equal(page.tab, 'problems');
    page.part('bytecode-tab').click();
    await page.respond('compile', {}, 503);
    assert.equal(page.tab, 'problems');
    page.part('result-tab').click();
    page.run();
    await page.respond('execute', {}, 429);
    assert.equal(page.tab, 'problems');
    const operations = page.part('problems-list').children.slice(1);
    assert.match(operations[0].children[0].textContent, /Bytecode/);
    assert.match(operations[1].children[1].textContent, /All execution slots are busy/);
    assert.equal(page.markers.size, 1);
});

test('local Run validation opens Problems while editor focus stays in place', async t => {
    const page = await mountPlayground(t);
    page.part('params-add').click();
    const [name, value] = page.part('params-rows').children[0].children;
    value.value = 'true';
    value.dispatchEvent(new Event('input'));
    page.editorInput.focus();
    page.run();
    assert.equal(page.tab, 'problems');
    assert.equal(page.focused, page.editorInput);
    assert.equal(name.attributes.get('aria-invalid'), 'true');
    assert.equal(page.part('problems-list').children[0].children[1].textContent, 'Enter a name for this parameter.');
    assert.equal(page.markers.size, 0);
});

test('stale tooling and Run settlements leave current problems, Ace markers, and tabs alone', async t => {
    const page = await mountPlayground(t);
    page.analyze();
    page.part('bytecode-tab').click();
    page.run();
    page.edit('RETURN 2');
    page.part('problems-tab').click();
    page.editorInput.focus();
    page.analyze();
    await page.respond('analyze', { status: 'invalid', diagnostics: [{ ...sourceProblem, message: 'Old problem' }] });
    assert.equal(page.part('problems-tab').textContent, 'Problems');
    await page.respond('analyze', invalidSource);
    const [marker] = page.markers.keys();
    await page.respond('compile', {}, 503);
    await page.respond('execute', { ...validSource, output: 1 });
    assert.equal(page.tab, 'problems');
    assert.equal(page.part('problems-tab').textContent, 'Problems (1)');
    assert.deepEqual([...page.markers.keys()], [marker]);
    assert.equal(page.part('result-text').hidden, true);
    assert.equal(page.focused, page.editorInput);

    page.part('result-tab').click();
    page.run();
    page.part('params-add').click();
    await page.respond('execute', {}, 429);
    assert.equal(page.tab, 'result');
    assert.equal(page.part('problems-tab').textContent, 'Problems (1)');
    assert.deepEqual([...page.markers.keys()], [marker]);
});

test('successful execute renders the JSON value, including falsy output and absent content', () => {
    for (const value of [null, false, 0, '', [], {}]) {
        const output = JSON.stringify(value);
        const body = `{"status":"ok","diagnostics":[],"output":${output},"timings":{"totalMs":0}}`;
        const response = readPlaygroundResponse('execute', 200, body);
        assert.equal(response.status, 'ok');
        assert.equal(response.outputText, output);
        assert.equal(response.timings.totalMs, 0);
    }
    assert.equal(readPlaygroundResponse('execute', 200, '{"status":"ok","diagnostics":[]}').outputText, undefined);
});

test('execute preserves numeric lexemes and handles quoted delimiters and nested output keys', () => {
    const body = '{"before":{"output":0},"output":[{"output":9223372036854775807},"😀 \\"output\\": }",1e-9],"status":"ok","diagnostics":[]}';
    const response = readPlaygroundResponse('execute', 200, body);
    assert.equal(response.outputText, [
        '[',
        '  {',
        '    "output": 9223372036854775807',
        '  },',
        '  "😀 \\"output\\": }",',
        '  1e-9',
        ']',
    ].join('\n'));
    const longText = 'f'.repeat(512 * 1024);
    const longBody = JSON.stringify({ status: 'ok', diagnostics: [], output: longText });
    assert.equal(readPlaygroundResponse('execute', 200, longBody).outputText, JSON.stringify(longText));
});

test('program failure outcomes stay distinct from errors at the HTTP or transport boundary', () => {
    const diagnostic = { severity: 'error', phase: 'execute', code: 'execution_timeout', message: 'deadline' };
    const failed = readPlaygroundResponse('execute', 200, JSON.stringify({ status: 'failed', diagnostics: [diagnostic] }));
    assert.equal(failed.status, 'failed');
    assert.equal(failed.outputText, undefined);
    assert.match(playgroundErrorMessage(failed.diagnostics[0]), /execution time limit/);
    assert.match(playgroundErrorMessage({ code: 'output_too_large' }), /output limit/);
    assert.equal(readPlaygroundResponse('format', 200, '{"status":"invalid","diagnostics":[]}').status, 'invalid');

    const errors = [
        [429, '{"error":{"code":"execution_capacity_exceeded","message":"busy"}}', 'execution_capacity_exceeded', /slots are busy/],
        [413, 'Payload Too Large', 'request_too_large', /request limit/],
        [422, '{"error":{"code":"response_too_large","message":"size"}}', 'response_too_large', /response.*size limit/],
        [503, '', 'service_unavailable', /temporarily unavailable/],
        [200, 'not JSON', 'invalid_response', /unexpected response/],
    ];
    for (const [status, body, code, message] of errors) {
        assert.throws(() => readPlaygroundResponse('execute', status, body),
            error => error.code === code && message.test(playgroundErrorMessage(error)));
    }
    assert.throws(() => readPlaygroundResponse('format', 200, '{"status":"ok","diagnostics":[]}'), { code: 'invalid_response' });
});

test('the client uses advertised limits, the exact request shapes, and no execution retry', async () => {
    const calls = [];
    const client = createPlaygroundClient('https://playground.example/', async (url, options) => {
        calls.push([url.pathname, options]);
        switch (url.pathname) {
        case '/api/v1/environment':
            return Response.json(runtime);
        case '/api/v1/analyze':
        case '/api/v1/format':
            return Response.json({ status: 'invalid', diagnostics: [] });
        case '/api/v1/compile':
            return Response.json({ status: 'ok', diagnostics: [], artifacts: { bytecode: { format: 'text', text: 'RET' }, parameters: ['payload'] } });
        default:
            throw new TypeError('Connection failed');
        }
    });
    const env = await client.request('environment');
    assert.equal(calls[0][1].method, undefined);
    assert.equal(canInspect(env), true);
    assert.equal(canInspect({ ...env, capabilities: { ...env.capabilities, inspectionArtifacts: ['bytecode'] } }), false);
    assert.throws(() => readPlaygroundResponse('environment', 200, '{"apiVersion":"v2"}'), { code: 'invalid_response' });
    const source = { text: 'RETURN @payload', name: 'query.fql' };
    for (const operation of ['analyze', 'format', 'compile']) {
        await client.request(operation, { source });
    }
    const paramsJSON = '{"payload":9223372036854775807}';
    await assert.rejects(client.request('execute', { source, paramsJSON }), { code: 'network' });
    assert.deepEqual(calls.slice(1).map(([url, options]) => [url, options.method, options.headers, options.body]), [
        ['/api/v1/analyze', 'POST', { 'Content-Type': 'application/json' }, JSON.stringify({ source })],
        ['/api/v1/format', 'POST', { 'Content-Type': 'application/json' }, JSON.stringify({ source })],
        ['/api/v1/compile', 'POST', { 'Content-Type': 'application/json' }, JSON.stringify({ source, include: ['bytecode', 'parameters'] })],
        ['/api/v1/execute', 'POST', { 'Content-Type': 'application/json' }, `{"source":${JSON.stringify(source)},"params":${paramsJSON}}`],
    ]);
    await assert.rejects(client.request('analyze', { source: { text: '😀'.repeat(17) } }), { code: 'source_too_large' });
    await assert.rejects(client.request('execute', { source, paramsJSON: JSON.stringify({ payload: 'x'.repeat(256) }) }), { code: 'request_too_large' });
    assert.equal(calls.length, 5);
});

test('superseded tooling cannot overwrite a newer request even if cancellation races with the response', async () => {
    const state = createRequestState();
    const analysis = state.start('analyze');
    const compile = state.start('compile');
    const deferred = Promise.withResolvers();
    const accepted = [];
    const first = deferred.promise.then(value => {
        if (state.canApply(analysis)) accepted.push(value);
        state.finish(analysis);
    });
    state.changeSource();
    assert.equal(analysis.controller.signal.aborted, true);
    assert.equal(compile.controller.signal.aborted, true);
    assert.equal(state.canApply(compile), false);
    const next = state.start('analyze');
    deferred.resolve('old diagnostics');
    await first;
    assert.deepEqual(accepted, []);
    assert.equal(state.canApply(next), true);
});

test('editing source or parameters makes a run stale and keeps its submission locked until settlement', () => {
    const state = createRequestState();
    const first = state.start('execute');
    const analysis = state.start('analyze');
    state.changeParams();
    assert.equal(state.canApply(first), false);
    assert.equal(state.canApply(analysis), true);
    assert.equal(state.start('execute'), null);
    assert.equal(first.controller.signal.aborted, false);
    state.finish(first);
    const second = state.start('execute');
    state.changeSource();
    assert.equal(state.canApply(second), false);
    assert.equal(state.matches(first.revision), false);
    assert.equal(state.start('execute'), null);
    state.finish(second);
    assert.ok(state.start('execute'));
});

test('only the current, valid formatter response applies one complete undoable replacement', () => {
    const state = createRequestState();
    const changes = [];
    const source = 'LET x=1\nreturn x';
    const editor = {
        getValue: () => source,
        session: {
            getLength: () => 2,
            getLine: row => source.split('\n')[row],
            markUndoGroup: () => changes.push('undo boundary'),
            replace: (range, text) => changes.push({ range, text }),
        },
    };
    const previous = state.start('format');
    state.changeSource();
    const current = state.start('format');
    const formattedSource = 'LET x = 1\nRETURN x';
    if (state.canApply(previous)) applyFormattedSource(editor, { status: 'ok', formattedSource });
    assert.deepEqual(changes, []);
    assert.equal(applyFormattedSource(editor, { status: 'invalid' }), false);
    assert.equal(applyFormattedSource(editor, { status: 'ok', formattedSource: source }), false);
    assert.deepEqual(changes, []);
    assert.equal(state.canApply(current), true);
    assert.equal(applyFormattedSource(editor, { status: 'ok', formattedSource }), true);
    assert.deepEqual(changes, [
        'undo boundary',
        { range: { start: { row: 0, column: 0 }, end: { row: 1, column: 8 } }, text: formattedSource },
        'undo boundary',
    ]);
});

test('diagnostics use zero-based UTF-16 columns, with optional exclusive or zero-width ranges', () => {
    const line = 'LET value = "é😀" broken';
    const start = line.indexOf('broken');
    assert.equal(start, 18);
    assert.deepEqual(diagnosticRange({ location: { range: { start: { line: 1, character: start }, end: { line: 1, character: start + 6 } } } }), {
        start: { row: 1, column: 18 },
        end: { row: 1, column: 24 },
    });
    const eof = { line: 2, character: 0 };
    assert.deepEqual(diagnosticRange({ location: { range: { start: eof, end: eof } } }), {
        start: { row: 2, column: 0 }, end: { row: 2, column: 0 },
    });
    assert.equal(diagnosticRange({ location: { source: 'query.fql' } }), null);
});

test('params preserve JSON-or-string values, nested values, and special keys', () => {
    const result = serializeParameters([
        [' payload ', '{\n "nested": [true, false, null, 0, 9223372036854775807]\n}'],
        ['string', '"false"'], ['text', 'hello 🌿'], ['empty', ''], ['__proto__', '{"field":true}'], ['', ''],
    ]);
    const parsed = JSON.parse(result.json);
    assert.deepEqual(parsed.payload.nested.slice(0, 4), [true, false, null, 0]);
    assert.match(result.json, /9223372036854775807/);
    assert.equal(parsed.string, 'false');
    assert.equal(parsed.text, 'hello 🌿');
    assert.equal(parsed.empty, '');
    assert.deepEqual(parsed.__proto__, { field: true });
    assert.deepEqual(parameterRows({ text: 'null', value: null }), [['text', '"null"'], ['value', 'null']]);
    assert.match(serializeParameters([['x', '0'], [' x ', 'null']]).error.message, /more than once/);
});

test('versioned share fragments round-trip Unicode, raw typed params and empty source, and support old Latin-1 links', () => {
    const source = '// 日本語 and 🌿\r\nRETURN @payload';
    const entries = [['payload', '{"query":"café 😀", "id":9223372036854775807}'], ['flag', '"false"']];
    const link = createShareURL(origin + '?script=legacy&other=1', source, entries);
    assert.match(link, /^https:\/\/ferretlang\.org\/try\/\?other=1#fql=v1\.[\w-]+$/);
    assert.deepEqual(readShareURL(link), { source, entries });
    assert.deepEqual(readShareURL(createShareURL(origin, '', [])), { source: '', entries: [] });
    assert.deepEqual(readShareURL(origin + '?script=++k='), { source: 'ûé', entries: [] });
    assert.equal(readShareURL(origin + '#source'), null);
});

test('oversized links, malformed UTF-8 and invalid shared rows are rejected', () => {
    assert.throws(() => createShareURL(origin, '🌿'.repeat(20000), []), /too large/);
    assert.throws(() => readShareURL(origin + '#fql=v1.' + 'a'.repeat(65536)), /too large/);
    assert.throws(() => readShareURL(origin + '#fql=v2.Zm9v'), /unsupported version/);
    assert.throws(() => readShareURL(origin + '#fql=v1._w'), /invalid/);
    assert.throws(() => readShareURL(createShareURL(origin, 'RETURN @x', [['x', '1'], [' x', '2']])), /invalid/);
});
