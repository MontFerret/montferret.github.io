const utf8 = new TextEncoder();
const inspectionArtifacts = ['bytecode', 'parameters'];
const jsonToken = /"[^"\\]*(?:\\[\s\S][^"\\]*)*"|[{}\[\]:,]|[^\s{}\[\]:,"]+/g;

class PlaygroundError extends Error {
    constructor(code, message = '') {
        super(message);
        this.code = code;
    }
}

const errorMessages = {
    execution_capacity_exceeded: 'All execution slots are busy. Try Run again in a moment.',
    rate_limited: 'The Playground is receiving too many requests. Try again in a moment.',
    source_too_large: 'Source exceeds the Playground’s byte limit.',
    request_too_large: 'Source and parameters exceed the Playground’s request limit.',
    execution_timeout: 'The query exceeded its execution time limit.',
    output_too_large: 'The result exceeded the Playground’s output limit.',
    response_too_large: 'The generated response exceeded the Playground’s size limit.',
    request_timeout: 'The Playground could not respond in time.',
    request_canceled: 'The Playground canceled the request.',
    browser_timeout: 'The browser stopped waiting for the Playground.',
    network: 'Could not reach the Playground. Check your connection and try again.',
    service_unavailable: 'The Playground is temporarily unavailable.',
    invalid_response: 'The Playground returned an unexpected response.',
};

export function playgroundErrorMessage(error) {
    return Object.hasOwn(errorMessages, error.code)
        ? errorMessages[error.code]
        : error.message || 'The Playground could not complete the request.';
}

export function canInspect(environment) {
    return environment?.capabilities.compile === true
        && inspectionArtifacts.every(artifact => environment.capabilities.inspectionArtifacts.includes(artifact));
}

function validateEnvironment(body) {
    const { runtime, capabilities, limits } = body ?? {};
    if (body?.apiVersion !== 'v1'
        || !runtime || !(typeof runtime.version === 'string' || runtime.version === null)
        || !['analyze', 'format', 'compile', 'execute'].every(key => typeof capabilities?.[key] === 'boolean')
        || !Array.isArray(capabilities.inspectionArtifacts)
        || ![limits?.tooling?.maxSourceBytes, limits?.http?.maxRequestBodyBytes, limits?.http?.writeTimeoutMs,
            limits?.execution?.maxOutputBytes, limits?.execution?.maxTimeoutMs].every(value => Number.isFinite(value) && value > 0)) {
        throw new PlaygroundError('invalid_response');
    }
    return body;
}

function httpError(operation, status, body) {
    if (typeof body?.error?.code === 'string' && typeof body.error.message === 'string') {
        return new PlaygroundError(body.error.code, body.error.message);
    }
    const codes = {
        408: 'request_timeout',
        413: 'request_too_large',
        422: 'response_too_large',
        429: operation === 'execute' ? 'execution_capacity_exceeded' : 'rate_limited',
        503: 'service_unavailable',
        504: 'request_timeout',
    };
    return new PlaygroundError(codes[status], `The Playground returned HTTP ${status}.`);
}

// Read the top-level output value from validated JSON so JavaScript's number
// representation does not round a Ferret int64 before display or copying.
function rawOutput(text) {
    let depth = 0;
    let previous;
    let outputKey = false;
    let start;
    for (const { 0: token, index } of text.matchAll(jsonToken)) {
        if (start !== undefined && depth === 1 && (token === ',' || token === '}')) {
            return text.slice(start, index).trim();
        }
        if (outputKey) {
            start = index + 1; // Skip the colon.
            outputKey = false;
        } else if (depth === 1 && (previous === '{' || previous === ',')
            && token[0] === '"' && JSON.parse(token) === 'output') {
            outputKey = true;
        }
        if (token === '[' || token === '{') depth++;
        if (token === ']' || token === '}') depth--;
        previous = token;
    }
    throw new PlaygroundError('invalid_response');
}

function indentJSON(text) {
    // Bound display expansion for large or deeply nested results.
    if (text.length > 128 * 1024) return text;
    const parts = [];
    let depth = 0;
    let previous;
    for (const [token] of text.matchAll(jsonToken)) {
        const opening = previous === '[' || previous === '{';
        const closing = token === ']' || token === '}';
        if (closing) depth--;
        if (previous === ',' || (opening && !closing) || (closing && !opening)) {
            parts.push('\n', '  '.repeat(Math.min(depth, 12)));
        } else if (previous === ':') {
            parts.push(' ');
        }
        parts.push(token);
        if (token === '[' || token === '{') depth++;
        previous = token;
    }
    return parts.join('');
}

export function readPlaygroundResponse(operation, status, text) {
    const ok = status >= 200 && status < 300;
    let body;
    try {
        body = JSON.parse(text);
    } catch {
        if (ok) throw new PlaygroundError('invalid_response');
    }
    if (!ok) throw httpError(operation, status, body);
    if (operation === 'environment') return validateEnvironment(body);

    if (!(body?.status === 'ok' || body?.status === 'invalid' || (operation === 'execute' && body?.status === 'failed'))
        || !Array.isArray(body.diagnostics)) {
        throw new PlaygroundError('invalid_response');
    }
    if (operation === 'format' && body.status === 'ok' && typeof body.formattedSource !== 'string') {
        throw new PlaygroundError('invalid_response');
    }
    if (operation === 'compile' && body.status === 'ok'
        && (body.artifacts?.bytecode?.format !== 'text' || typeof body.artifacts.bytecode.text !== 'string'
            || !Array.isArray(body.artifacts.parameters))) {
        throw new PlaygroundError('invalid_response');
    }
    const outputText = operation === 'execute' && body.status === 'ok' && Object.hasOwn(body, 'output')
        ? indentJSON(rawOutput(text))
        : undefined;
    return { ...body, outputText };
}

function requestBody(operation, source, paramsJSON) {
    if (operation === 'compile') {
        return JSON.stringify({ source, include: inspectionArtifacts });
    }
    if (operation === 'execute') {
        return `{"source":${JSON.stringify(source)},"params":${paramsJSON}}`;
    }
    return JSON.stringify({ source });
}

export function createPlaygroundClient(baseURL, fetcher = globalThis.fetch) {
    let environment;

    async function request(operation, { source, paramsJSON, signal } = {}) {
        const url = new URL(`api/v1/${operation}`, `${baseURL.replace(/\/+$/, '')}/`);
        const options = { credentials: 'omit', cache: 'no-store' };
        if (source) {
            if (utf8.encode(source.text).length > environment.limits.tooling.maxSourceBytes) {
                throw new PlaygroundError('source_too_large');
            }
            const body = requestBody(operation, source, paramsJSON);
            if (utf8.encode(body).length > environment.limits.http.maxRequestBodyBytes) {
                throw new PlaygroundError('request_too_large');
            }
            Object.assign(options, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body,
            });
        }
        const deadline = AbortSignal.timeout(environment ? Math.ceil(environment.limits.http.writeTimeoutMs + 5000) : 15000);
        options.signal = signal ? AbortSignal.any([signal, deadline]) : deadline;

        let response;
        let text;
        try {
            response = await fetcher(url, options);
            text = await response.text();
        } catch {
            signal?.throwIfAborted();
            throw new PlaygroundError(deadline.aborted ? 'browser_timeout' : 'network');
        }
        const data = readPlaygroundResponse(operation, response.status, text);
        if (operation === 'environment') environment = data;
        return data;
    }

    return { request };
}
