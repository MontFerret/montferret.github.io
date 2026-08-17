import { readFile } from "node:fs/promises";
import { fileURLToPath, pathToFileURL } from "node:url";
import path from "node:path";

const nativeFetch = globalThis.fetch;
globalThis.fetch = async (input, init) => {
    const target = input instanceof URL
        ? input
        : new URL(typeof input === "string" ? input : input.url);

    if (target.protocol === "file:") {
        return new Response(await readFile(fileURLToPath(target)));
    }

    return nativeFetch(input, init);
};

const pagefindURL = pathToFileURL(path.join(process.cwd(), "public", "pagefind", "pagefind.js"));
const pagefind = await import(pagefindURL.href);
await pagefind.options({ baseUrl: "/" });
await pagefind.init();

const expectations = new Map([
    ["Math", "/docs/standard-library/math/"],
    ["FLATTEN", "/docs/standard-library/functions/flatten/"],
    ["TO_NUMBER", "/docs/standard-library/functions/to_number/"],
    ["IO", "/docs/standard-library/io/"],
    ["IO::FS::READ", "/docs/standard-library/io/fs/read/"],
    ["Testing", "/docs/standard-library/testing/"],
]);

for (const [term, expectedURL] of expectations) {
    const search = await pagefind.search(term);
    const results = await Promise.all(search.results.slice(0, 20).map((result) => result.data()));
    const found = results.some((result) => new URL(result.url, "https://ferretlang.org").pathname === expectedURL);
    if (!found) {
        const received = results.map((result) => result.url).join(", ");
        throw new Error(`Pagefind query ${JSON.stringify(term)} did not return ${expectedURL}; received ${received}`);
    }
}

await pagefind.destroy();
