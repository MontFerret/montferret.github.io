import assert from 'node:assert/strict';
import { test } from 'node:test';
import { canInspect, createPlaygroundClient, playgroundErrorMessage, readPlaygroundResponse } from './playground-api.mjs';
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
