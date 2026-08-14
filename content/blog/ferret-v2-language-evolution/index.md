---
title: "Ferret v2: Moving Beyond AQL"
subtitle: "What We’re Keeping, What We’re Changing, and Why"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-08-13"
---

One of the fun (and sometimes dangerous) parts of working on a language for a long time is that, at some point, it starts developing a personality of its own.

You add a feature because you need it, another one because the previous feature made something else possible, and then one day you look at the language and realize that some of the conventions you started with don't quite fit anymore.

That's more or less where Ferret is with v2.

Ferret started life with a pretty strong influence from ArangoDB's AQL, and you can still see that heritage immediately in a traditional Ferret script:

{{< code lang="fql" >}}
FOR user IN users
    FILTER user.active
    RETURN user.email
{{</ code >}}

That syntax has served Ferret very well. It made data transformations familiar, kept scripts concise, and gave the language the declarative foundation that is still very much at the heart of Ferret today.

But over the years, Ferret has grown quite a bit beyond that original model and v2 has pushed it much further.

It has functions and pattern matching now. It has mutable state when you need it, synchronization primitives for dealing with asynchronous systems, host-defined and queryable values, and an embedding API that increasingly treats Ferret as a language applications can extend rather than simply a query engine they can call.

And while working on v2, I've increasingly found myself asking a slightly different question. Not *"how would AQL express this?"*, but simply *"how should Ferret express this?"*

That's a small distinction, but I think it's an important one.

At some point, carrying every convention inherited from AQL stops making the language feel familiar and starts obscuring what it has actually become. So with v2, we're letting go of a few of those conventions — while keeping the parts of that declarative foundation that still make Ferret, well, Ferret.

## Lowercase is now canonical

Ferret keywords have long been case-insensitive, and that's not changing. Both of these are still perfectly valid:

{{< code lang="fql" >}}
FOR user IN users
    FILTER user.active
    RETURN user
{{</ code >}}

{{< code lang="fql" >}}
for user in users
    filter user.active
    return user
{{</ code >}}

Starting with Ferret v2, though, the second form becomes the canonical one. The formatter, documentation, website, and official examples will all use lowercase keywords:

{{< code lang="fql" >}}
let users = query "users" in db

return for user in users {
    filter user.active

    return {
        name: user.name,
        email: user.email
    }
}
{{</ code >}}

On the surface this is obviously a cosmetic change, but I think it changes the first impression of the language quite a bit.

Uppercase keywords carry a very strong association with SQL and other query languages, which made perfect sense when Ferret was much closer to its AQL origins. Today, that association can actually be a little misleading because querying is only one part of what the language does.

Ferret is designed to be embedded into applications, and perfectly normal Ferret code now looks like this:

{{< code lang="fql" >}}
func classify(response) {
    return match response.status {
        200 => "ok"
        404 => "missing"
        _ => "error"
    }
}
{{</ code >}}

or this:

{{< code lang="fql" >}}
{{< code lang="fql" >}}
click(submit)
waitfor exists confirmation
{{</ code >}}

or this:

{{< code lang="fql" >}}
for request in requests {
    filter request.amount > 1000
    return request.id
}
{{</ code >}}

With lowercase keywords, all of these constructs simply look like parts of the same small programming language instead of programming constructs gradually grafted onto something SQL-like.

And that's much closer to how I think about Ferret today.

## Declarative doesn't have to mean SQL-like

Changing the visual style doesn't mean changing Ferret's execution philosophy. Ferret is still very deliberately a declarative-first language; it just turns out that uppercase keywords aren't what makes a language declarative.

Consider:

{{< code lang="fql" >}}
for user in users {
    filter user.active
    return user.email
}
{{</ code >}}

This describes a selection and a transformation. It doesn't tell the runtime to iterate a collection, test a condition, skip an iteration, append a value to an accumulator, and then continue with the next item. That difference is important, and it's something I very much want to preserve.

At the same time, Ferret does have imperative features. Mutable variables and `for while` exist because, well, real-world automation sometimes requires state, and browser automation is probably the easiest place to see why:

{{< code lang="fql" >}}
let previous_count = 0
var done = false

for while !done {
    scroll(page)

    waitfor exists page.items

    let current_count = count(page.items)

    // ...
}
{{</ code >}}

Infinite scrolling, pagination, changing DOM state, asynchronous events, and all the other weird things websites do don't always fit neatly into a purely declarative model, and I don't think Ferret should pretend otherwise.

So the goal isn't declarative purity. It's to provide enough state and iteration to deal with the real world without slowly turning Ferret into yet another general-purpose scripting language.

## Why `match` instead of `if`

That philosophy is also why Ferret v2 gained pattern matching rather than conventional `if`/`else` control flow.

For example:

{{< code lang="fql" >}}
let status = match response.status {
    200 => "ok"
    404 => "missing"
    status when status >= 500 => "server-error"
    _ => "unexpected"
}
{{</ code >}}

Pattern matching keeps branching value-oriented, and the same idea applies elsewhere in the language. If I want to filter a collection, for example, I don't need an `if` followed by `continue`:

{{< code lang="fql" >}}
return for user in users {
    filter user.active
    return user
}
{{</ code >}}

Ferret already has a construct that says exactly what I mean.

I think this gives us a useful rule for evolving the language: before adding another general-purpose control-flow mechanism, ask whether the operation we're trying to express deserves a more declarative primitive of its own.

Sometimes the answer will still be an imperative feature — `while` exists for a reason — but I don't want that to be the default answer.

## Scripts no longer have to return something

There is another piece of AQL inheritance we're dropping in v2: historically, every Ferret script had to produce a value.

For a query language, that makes perfect sense. You run a query because you want it to calculate and return something.

For this, though?

{{< code lang="fql" >}}
click(submit)
waitfor exists confirmation
{{</ code >}}

What exactly should it return?

Previously, scripts like this had to manufacture a result even when there wasn't a meaningful one:

{{< code lang="fql" >}}
click(submit)
waitfor exists confirmation

return none
{{</ code >}}

In Ferret v2, they don't have to. Reaching the end of a script simply means that it completed successfully with `none`.

You can still write the explicit version:

{{< code lang="fql" >}}
return none
{{</ code >}}

but there is no reason to require it when the script exists entirely for its effects.

It's a small change, but effect-oriented Ferret programs feel much more natural this way, especially when Ferret is embedded as a DSL inside another application.

## Functions follow the same rule

User-defined functions now behave the same way. If a function exists for its effects, it doesn't need a ceremonial `return none` either:

{{< code lang="fql" >}}
func notify(user) {
    send(user.email)
    audit(user.id)
}
{{</ code >}}

If execution reaches the end, the function returns `none`.

Functions that actually produce values continue to use `return`:

{{< code lang="fql" >}}
func normalize(user) {
    return {
        name: trim(user.name),
        email: lower(user.email)
    }
}
{{</ code >}}

And expression-bodied functions stay nice and concise:

{{< code lang="fql" >}}
func square(x) => x * x
{{</ code >}}

One thing we're deliberately *not* doing is implicitly returning the last arbitrary expression from a block. So this:

{{< code lang="fql" >}}
func square(x) {
    x * x
}
{{</ code >}}

doesn't return `x * x`.

If a block-bodied function produces a value, it says so explicitly.

## `for` is no longer secretly a script return

That brings us to one actual breaking change.

Historically, if a `for` loop was the final statement in a Ferret script, its result could implicitly become the result of the whole script:

{{< code lang="fql" >}}
for user in users {
    filter user.active
    return user.email
}
{{</ code >}}

This made a lot of sense when Ferret was much more query-shaped, but once scripts are allowed to complete without producing a value, making the final `for` magically special starts looking pretty strange.

A `for` already produces a value in Ferret, and you can capture that value just like any other:

{{< code lang="fql" >}}
let emails = (
    for user in users {
        filter user.active
        return user.email
    }
)

return emails
{{</ code >}}

The fact that the same construct also happened to return from the entire script just because it appeared at the end was really a separate rule, and v2 removes it.

A standalone `for` is still value-producing, but if you don't use that value, it's discarded. If you want the result of the loop to become the result of the script, just say so:

{{< code lang="fql" >}}
return for user in users {
    filter user.active
    return user.email
}
{{</ code >}}

I actually like how this reads:

> return, for every active user, their email.

More importantly, though, `for` now means exactly the same thing wherever it appears.

You can assign it:

{{< code lang="fql" >}}
let emails = (
    for user in users {
        return user.email
    }
)
{{</ code >}}

return it:

{{< code lang="fql" >}}
return for user in users {
    return user.email
}
{{</ code >}}

or use it standalone when you simply don't care about the produced value:

{{< code lang="fql" >}}
for user in users {
    log(user)
}
{{</ code >}}

Its position in a body no longer quietly changes what it means.

## One rule for executable bodies

The nice part is that all of these changes leave Ferret with a much simpler result model.

For scripts and block-bodied functions:

- `return value` explicitly returns a value.
- `return for ...` explicitly returns the result of a loop.
- Reaching the end returns `none`.
- The result of a standalone value-producing construct is discarded.
- The last arbitrary expression is never implicitly returned.

Expression-bodied functions remain explicitly value-producing:

{{< code lang="fql" >}}
func double(value) => value * 2
{{</ code >}}

That's pretty much it. We no longer need one rule saying that every script must return something and another special rule promoting a final `for` into the script result.

## What about existing Ferret programs?

The good news is that changing the canonical keyword casing doesn't break anything. Uppercase remains perfectly valid:

{{< code lang="fql" >}}
FOR item IN items {
    RETURN item
}
{{</ code >}}

The formatter will simply turn it into:

{{< code lang="fql" >}}
for item in items {
    return item
}
{{</ code >}}

The change to final `for` semantics, however, is intentionally breaking.

A v1-style script like this:

{{< code lang="fql" >}}
FOR item IN items
    RETURN item
{{</ code >}}

which relies on the loop becoming the script result needs to become the explicit v2 equivalent:

{{< code lang="fql" >}}
return for item in items {
    return item
}
{{</ code >}}

Fortunately, you won't have to go hunting through old Ferret scripts looking for these cases by hand. The CLI migration command will recognize affected constructs and migrate them to their explicit v2 equivalents where possible.

That's really the goal of the migration tooling in general: make semantic changes explicit without making upgrading existing Ferret programs unnecessarily painful.

## A language, not just a query

None of these changes individually redefine Ferret, but taken together they reflect something that's been happening gradually throughout v2.

Ferret isn't trying to become a general-purpose programming language, but it isn't just a query language anymore either. It's becoming a small, declarative-first, expression-oriented language designed to be embedded into applications and specialized through capabilities provided by the host.

It can query:

{{< code lang="fql" >}}
let users = query "active-users" in db
{{</ code >}}

transform:

{{< code lang="fql" >}}
let names = (
    for user in users {
        return user.name
    }
)
{{</ code >}}

make decisions:

{{< code lang="fql" >}}
let action = match request.kind {
    "create" => create(request)
    "update" => update(request)
    _ => none
}
{{</ code >}}

and synchronize with the outside world:

{{< code lang="fql" >}}
click(submit)
waitfor exists confirmation
{{</ code >}}

And if that's everything the program needed to do, it can simply end. No artificial result required.

Ferret's AQL heritage is still very visible in the language, and I don't want to erase it — a lot of Ferret's best ideas grew directly from that foundation. But v2 feels like the right time to separate the ideas that are actually fundamental to Ferret from conventions that survived mostly because, well, they've always been there.

Lowercase syntax is probably the most immediately visible sign of that transition, but the deeper change has been happening for a while now:

Ferret is growing into its own language.