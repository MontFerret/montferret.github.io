export function createRequestState() {
    let source = 0;
    let params = 0;
    const pending = new Map();

    function matches(revision) {
        return revision.source === source && (revision.params === undefined || revision.params === params);
    }

    function revision(operation) {
        return { source, params: operation === 'execute' ? params : undefined };
    }

    return {
        changeSource() {
            source++;
            for (const [operation, request] of pending) {
                if (operation === 'execute') continue;
                pending.delete(operation);
                request.controller.abort();
            }
        },
        changeParams() {
            params++;
        },
        start(operation) {
            const previous = pending.get(operation);
            if (operation === 'execute' && previous) return null;
            previous?.controller.abort();
            const controller = new AbortController();
            const request = { operation, controller, revision: revision(operation) };
            pending.set(operation, request);
            return request;
        },
        finish(request) {
            if (pending.get(request.operation) === request) pending.delete(request.operation);
        },
        pending(operation) {
            return pending.has(operation);
        },
        revision,
        matches,
        canApply(request) {
            return pending.get(request.operation) === request && matches(request.revision);
        },
    };
}

export function diagnosticRange(diagnostic) {
    const { start, end } = diagnostic.location?.range ?? {};
    const coordinates = [start?.line, start?.character, end?.line, end?.character];
    if (!coordinates.every(value => Number.isSafeInteger(value) && value >= 0)
        || start.line > end.line || (start.line === end.line && start.character > end.character)) {
        return null;
    }
    return {
        start: { row: start.line, column: start.character },
        end: { row: end.line, column: end.character },
    };
}

export function applyFormattedSource(editor, response) {
    if (response.status !== 'ok' || response.formattedSource === editor.getValue()) return false;
    const session = editor.session;
    const lastRow = session.getLength() - 1;
    const start = { row: 0, column: 0 };
    const end = { row: lastRow, column: session.getLine(lastRow).length };
    // Ace groups session changes made in the same turn. Isolate this replacement
    // from the changes immediately before and after it.
    session.markUndoGroup();
    session.replace({ start, end }, response.formattedSource);
    session.markUndoGroup();
    return true;
}

function parameterValueJSON(value) {
    try {
        JSON.parse(value);
        return value.trim(); // Preserve numeric literals, including nested int64s.
    } catch {
        return JSON.stringify(value);
    }
}

export function serializeParameters(rows) {
    const entries = [];
    const fields = [];
    const names = new Set();
    for (const [index, [rawName, value]] of rows.entries()) {
        const name = rawName.trim();
        let message;
        if (!name && value.trim()) message = 'Enter a name for this parameter.';
        if (names.has(name)) message = `The parameter “${name}” appears more than once.`;
        if (message) return { error: { index, message } };
        if (!name) continue;
        names.add(name);
        entries.push([name, value]);
        fields.push(`${JSON.stringify(name)}:${parameterValueJSON(value)}`);
    }
    return { entries, json: `{${fields.join(',')}}` };
}

export function parameterRows(params) {
    return Object.entries(params ?? {}).map(([name, value]) => [name, JSON.stringify(value)]);
}

const sharePrefix = '#fql=v1.';
const shareLimit = 64 * 1024;

function checkShareSize(size) {
    if (size > shareLimit) throw new Error('This program is too large for a share link.');
}

export function createShareURL(href, source, entries) {
    const bytes = new TextEncoder().encode(JSON.stringify({ source, params: entries }));
    checkShareSize(Math.ceil(bytes.length * 4 / 3) + sharePrefix.length);
    const base64 = btoa(Array.from(bytes, byte => String.fromCharCode(byte)).join(''));
    const url = new URL(href);
    url.searchParams.delete('script');
    url.hash = sharePrefix + base64.replaceAll('+', '-').replaceAll('/', '_').replace(/=+$/, '');
    return url.href;
}

export function readShareURL(href) {
    const url = new URL(href);
    if (url.hash.startsWith('#fql=')) {
        checkShareSize(url.hash.length);
        try {
            const encoded = url.hash.slice(sharePrefix.length);
            if (!url.hash.startsWith(sharePrefix) || !/^[\w-]+$/.test(encoded)) throw new Error();
            const binary = atob(encoded.replaceAll('-', '+').replaceAll('_', '/'));
            const bytes = Uint8Array.from(binary, char => char.charCodeAt(0));
            const state = JSON.parse(new TextDecoder('utf-8', { fatal: true }).decode(bytes));
            if (typeof state?.source !== 'string' || !Array.isArray(state.params)
                || !state.params.every(row => Array.isArray(row) && row.length === 2 && row.every(value => typeof value === 'string'))) {
                throw new Error();
            }
            const params = serializeParameters(state.params);
            if (params.error) throw new Error();
            return { source: state.source, entries: params.entries };
        } catch {
            throw new Error('This share link is invalid or uses an unsupported version.');
        }
    }
    const legacy = url.searchParams.get('script');
    if (!legacy) return null;
    checkShareSize(legacy.length);
    try {
        // The old encoder used btoa directly and emitted literal plus signs.
        return { source: atob(legacy.replaceAll(' ', '+')), entries: [] };
    } catch {
        throw new Error('This legacy share link is invalid.');
    }
}
