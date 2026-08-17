---
title: "Parameters"
sidebarTitle: "Parameters"
weight: 30
draft: false
description: "Pass JavaScript values into FQL programs and understand value conversion."
---

# Parameters

Bind parameters pass application data into an FQL program without generating source code. FQL reads each value through its `@name`.

## Pass parameters to a run

Provide a plain JavaScript object as `params`:

{{< code lang="javascript" >}}
const orders = [
    { id: "A-100", customer: "Ada", total: 125 },
    { id: "A-101", customer: "Lin", total: 48 },
];

const result = await engine.run(`
    return for order in @orders {
        filter order.total >= @minimum
        return {
            id: order.id,
            customer: order.customer,
        }
    }
`, {
    params: {
        orders,
        minimum: 100,
    },
});

console.log(result); // [{ id: "A-100", customer: "Ada" }]
{{</ code >}}

Use the same `params` option with `plan.run()`. Use `plan.createSession({ params })` when a reusable session should capture one parameter set.

The JavaScript API does not currently define engine-wide persistent parameters. Bind values for each plan run or session.

## Inspect declared parameters

After compilation, `plan.params` lists the parameter names referenced by the program:

{{< code lang="javascript" >}}
const plan = await engine.compile(
    "return { name: @name, enabled: @enabled }",
);

console.log(plan.params); // ["name", "enabled"]
{{</ code >}}

The array is read-only. Missing required parameters fail when a session is created or executed.

## Value conversion

Parameters, JavaScript function arguments, and JavaScript function results use the same input conversion rules:

| JavaScript value | Ferret value |
| --- | --- |
| `undefined` or `null` | `none` |
| `boolean` | Boolean |
| `string` | String |
| finite `number` | Number |
| `Array` | Array |
| plain object | Object |
| `Uint8Array` | Binary |

Nested arrays and objects are converted recursively. Object keys remain strings.

For example:

{{< code lang="javascript" >}}
const result = await engine.run(`
    return {
        name: @name,
        count: @count,
        active: @active,
        tags: @tags,
        profile: @profile,
        missing: @missing,
    }
`, {
    params: {
        name: "Ada",
        count: 3,
        active: true,
        tags: ["admin", "editor"],
        profile: { team: "runtime" },
        missing: null,
    },
});
{{</ code >}}

## Unsupported values

Conversion fails explicitly for values that do not have a stable Ferret representation, including:

- `NaN`, `Infinity`, and `-Infinity`
- cyclic arrays or objects
- class instances
- `Date`, `Map`, `Set`, and other non-plain objects
- functions used as parameter values
- symbols and big integers

Convert these values into supported primitives, arrays, plain objects, or `Uint8Array` before passing them to Ferret.

## Returned values

Ferret results are JSON-decoded. `none` becomes `null`; booleans, strings, numbers, arrays, and objects retain their corresponding shapes. Binary results are serialized as base64 strings rather than returned as `Uint8Array`.

## Next steps

{{< docs-related tiles="language-parameters,embedding-javascript-executing,embedding-javascript-custom-functions,language-types-serialization" >}}
