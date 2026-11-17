---
title: "Inside Ferret v2: Rethinking Results"
subtitle: "Why we changed how Ferret returns results, and what it means for embedding, resource management, and the future of execution."
draft: true
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-11-17"
---

Hello, friends!

Over the past few months, a lot of work has gone into Ferret Core. Some of it has been fairly visible: language improvements, standard library changes, better diagnostics, and compiler optimizations.

But quite a bit of work has also happened underneath the surface.

One of the more significant changes is how Ferret handles the results of executing a program.

At first glance, it might look like a relatively small API refactoring. We replaced the old output structure with a new interface, added a couple of methods, and updated the code that consumes results.

But there is a bit more to it than that.

This change addresses some assumptions that made sense when Ferret was primarily a local query engine, but became increasingly limiting as it grew into an embeddable runtime with external execution, debugging, and communication protocols.

## The old approach

Until recently, running a Ferret program followed a fairly straightforward pattern.

You compiled a query, created a session, executed it, and received an output containing the encoded result.

Conceptually, it looked like this:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    return err
}

fmt.Println(string(output.Content))
{{< /code >}}

The output was already materialized. The caller received the encoded bytes and could immediately use them.

For many applications, this was perfectly fine.

If you're executing a small script and expecting a JSON document, there isn't much wrong with this approach.

The problems become more apparent when you start thinking about other execution environments.

What if Ferret runs in another process? What if the result is large? What if the caller wants to forward it somewhere without creating another copy? What happens if execution succeeds but encoding fails? Or if a resource cannot be closed properly?

More importantly, should the API assume that every result must be fully materialized before it can be delivered?

That last question is what eventually led to the redesign.

## An output is no longer just a collection of bytes

The new API makes a distinction between an **output handle** and the **content** it delivers.

Executing a program now returns an `Output` interface. Instead of exposing the data directly, it provides three main operations:

- `Collect` — obtains the encoded content as a detached value.
- `Consume` — delivers encoded bytes to a consumer callback.
- `Close` — abandons the output and releases resources associated with it.

There is also `Metadata`, which describes the encoded representation without requiring the caller to consume it.

For applications that simply want the entire result, the code is still straightforward:

{{< code lang="go" >}}
output, err := session.Run(ctx)
if err != nil {
    return err
}
defer output.Close()

content, err := output.Collect(ctx)

if content != nil {
    fmt.Println(string(content.Data))
}

return err
{{< /code >}}

The difference is that retrieving the result is now an explicit operation.

Once collected, the content belongs to the caller. It remains valid after the output handle has been closed.

This may seem like a small distinction, but it gives Ferret a much clearer boundary between execution, result delivery, and resource ownership.

## Two ways to handle results

Not every application needs to do the same thing with the result of a Ferret program.

Sometimes you want to collect everything into memory, inspect it, and process it.

Other times, you simply want to forward the result somewhere else: a file, an HTTP response, another service, or a remote client.

That's why the new API provides both `Collect` and `Consume`.

`Collect` is the convenient option. It produces a detached `Content` value containing the encoded bytes and their metadata.

`Consume` takes a different approach. Instead of returning the entire payload, it invokes a callback with chunks of encoded data.

This is especially useful for adapters that want to forward output through another interface.

There are a few rules worth mentioning.

The chunks are borrowed and remain valid only during the callback. If an application wants to keep them, it must copy them.

Chunk boundaries are also arbitrary. A chunk does not necessarily represent a complete JSON value, record, line, or even a complete character.

And importantly, **an output can only be consumed once**.

Calling `Collect` and then `Consume` on the same output is not supported. They are two alternatives for handling the same result.

This keeps ownership unambiguous and avoids introducing hidden buffering or replay behavior into the interface.

## Execution errors and output errors are different things

Another important change concerns error handling.

Previously, it was natural to interpret a successful `Run` call as meaning that the program had executed successfully and its result was ready.

That assumption no longer holds.

Today, a successful `Run` means that Ferret has accepted execution and returned a usable output handle.

The terminal outcome is reported when the output is collected or consumed.

This distinction is intentional.

There are several stages involved in producing a result: execution, encoding, delivery, and resource cleanup. Any of them can fail.

And those failures don't necessarily mean that no useful output exists.

For example, an encoder might produce part of a result before encountering an error. In that case, `Collect` can return both the available content and the error.

That's also why the example above checks `content != nil` independently of `err`.

A result can be present and empty, absent entirely, or partially available because something failed during encoding.

These are different situations, and the API now represents them explicitly.

The same principle applies to scripts that exist primarily for their side effects. Even if the application doesn't care about the returned bytes, it still needs to consume the output to observe the final execution outcome.

Simply calling `Close` abandons the result. It does not certify that execution succeeded.

## Resource ownership is now explicit

Ferret can work with values that represent external resources.

A host application might expose a file, a cursor, a network-backed object, or another value with its own lifecycle.

Those resources may also appear inside collections or values returned by a query.

That raises an important question: who is responsible for cleaning them up?

The recent runtime work has been focused on making ownership explicit throughout execution and result materialization.

Resources acquired during execution need to be released correctly, including when the query fails, encoding is interrupted, or the caller abandons the output.

The output contract now makes the caller's responsibility clearer.

A returned output handle must be collected, consumed, or closed. Successful consumption finalizes its resources, while explicit closure provides a way to abandon it safely.

This is particularly important for embedding.

Ferret should be able to expose capabilities supplied by a host application without leaving that application guessing which objects are still alive or who is supposed to close them.

It also makes the behavior easier to reason about when a runtime is reused across multiple sessions.

## Does this mean Ferret now streams everything?

Not quite.

This is an important distinction.

The new API supports consuming output through a chunk-based interface, but native Ferret currently executes and encodes results eagerly.

In other words, the runtime still materializes the encoded result before returning the output handle.

Calling `Consume` does not suddenly turn a query into an incremental execution pipeline.

And that's okay.

The purpose of this change was to establish a better contract first, without requiring the entire execution engine to be redesigned at the same time.

A local runtime can continue using an eager implementation, while another implementation may deliver data incrementally.

Both can expose the same API.

This is particularly valuable for the Universal API and for projects such as Ferret Wire, where execution and result delivery may happen across process boundaries.

The caller should not need to know whether a runtime is local, remote, buffered, or capable of incremental delivery simply to retrieve its output.

Of course, the implementation still determines the actual memory and performance characteristics. A streaming-capable interface is not, by itself, a guarantee of constant-memory execution.

But it gives us a much better starting point.

## Encoding received some attention too

The output refactoring was accompanied by another substantial change: the codec infrastructure.

Ferret supports different output representations, including JSON and MessagePack.

Previously, encoding was largely oriented around producing byte slices. The newer codec contract uses context-aware readers and writers instead:

{{< code lang="go" >}}
Encode(
    ctx context.Context,
    dst io.Writer,
    value runtime.Value,
) error

Decode(
    ctx context.Context,
    src io.Reader,
) (runtime.Value, error)
{{< /code >}}

This lets codecs work directly with I/O abstractions instead of requiring a byte-slice-oriented interface.

The implementation also clarifies several important behaviors: codecs borrow their readers and writers, do not close them, propagate errors, and observe cancellation during supported operations.

Both JSON and MessagePack have been updated to follow this model, with additional work on buffering, allocations, and resource discovery during encoding.

Native output still uses an owned buffer today, but the encoding layer no longer needs to be designed exclusively around that assumption.

For most Ferret users, none of this changes how they write FQL.

For developers embedding Ferret or implementing custom codecs, however, it provides a more consistent foundation.

## Part of a larger cleanup

The output changes are not happening in isolation.

Recent Ferret Core releases have also brought improvements to engine construction, session ownership, lifecycle management, and the public API.

The standard library has been moving toward a more consistent namespaced API, while compiler work has continued to separate optimizations from runtime execution. Tail-call elimination, for example, now belongs to the compiler optimizer rather than being a special responsibility of the VM.

These changes might look unrelated, but they share a common direction.

I want Ferret Core to have clear responsibilities and predictable contracts.

The compiler should compile and optimize programs. The VM should execute them. Sessions should own their execution environment. Encoders should encode values. And the application embedding Ferret should have an explicit way to receive results and manage their lifetime.

That might sound obvious, but getting these boundaries right takes considerably more work than it seems.

Especially when the language is evolving alongside the runtime.

## What changes for applications embedding Ferret?

If you're using Ferret through the CLI, this change doesn't introduce new FQL syntax or alter how `return` works.

The main impact is on applications using the Go embedding API or the Universal API.

Code that previously accessed an output's `Content` field now needs to collect or consume the output explicitly.

Encoded data lives in `Content.Data`, while metadata such as the content type is available through `Output.Metadata()` or the collected content's metadata.

Applications also need to keep the execution context alive until output handling is complete.

And because the change is intentionally breaking, integrations need to adopt the updated API rather than expecting the old materialized output structure to continue working.

The underlying contract is shared with the Universal Ferret API, so adapters and other tools can follow the same rules without inventing their own result lifecycle.

These changes landed as part of the recent Ferret v2 alpha work, including [Ferret Core's output refactoring](https://github.com/MontFerret/ferret/pull/1057), the [codec I/O improvements](https://github.com/MontFerret/ferret/pull/1058), and the corresponding [Universal API changes](https://github.com/MontFerret/api/pull/3).

## Closing thoughts

One of the advantages of still being in alpha is that I can make changes like this before the API becomes something developers have to rely on for years.

And I'd rather make an inconvenient breaking change now than preserve an abstraction that will eventually become a limitation.

The old output model wasn't fundamentally wrong. It was appropriate for the execution model Ferret had at the time.

But Ferret has grown beyond that.

Today, it needs to support local embedding, reusable sessions, external runtimes, debugging, different encoding formats, and increasingly sophisticated integrations.

The new output model doesn't solve every problem related to streaming or execution. There is still work to do, particularly if we want to avoid eager materialization in more situations.

What it does give us is a cleaner foundation on which to build those capabilities without having to redesign the public interface again.

And that, more than any individual API improvement, is what much of the recent work on Ferret v2 has been about.

Not simply adding features, but making sure the foundations are solid enough for whatever comes next.

Thanks for reading!