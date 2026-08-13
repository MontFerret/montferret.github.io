---
title: "Ferret v2: Language Evolution"
subtitle: "Ferret Is Growing Into Its Own Language"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-08-13"
---
Ferret started life with a strong influence from ArangoDB's AQL.

You can still see that heritage immediately in a traditional Ferret script:

```fql
FOR user IN users
    FILTER user.active
    RETURN user.email
```

That syntax served Ferret well. It made data transformations familiar, kept scripts concise, and gave the language a strongly declarative foundation.

But Ferret has changed.

Especially with Ferret v2, the language is no longer best understood as a query language adapted for automation. It has functions, pattern matching, mutable state, synchronization primitives, host-defined values, queryable objects, and an increasingly capable embedding API.

At some point, carrying every convention inherited from AQL stops providing familiarity and starts obscuring what the language has become.

So we're changing a few of them.

## Lowercase is now canonical

Ferret keywords have long been case-insensitive. This isn't changing.

These are still valid:

```fql
FOR user IN users
    FILTER user.active
    RETURN user
```

```fql
for user in users
    filter user.active
    return user
```

But starting with Ferret v2, the second form is the canonical one.

The formatter, documentation, website, and official examples will use lowercase keywords:

```fql
let users = query "users" in db

return for user in users {
    filter user.active

    return {
        name: user.name,
        email: user.email
    }
}
```

This may look like a cosmetic change, but it changes the first impression of the language considerably.

Uppercase keywords carry a strong association with SQL and other query languages. That association was useful when Ferret was much closer to its AQL origins.

Today it can be misleading.

Ferret is designed to be embedded into applications and used for more than querying:

```fql
func classify(response) {
    return match response.status {
        200 => "ok"
        404 => "missing"
        _ => "error"
    }
}
```

```fql
click(submit)
waitfor exists confirmation
```

```fql
for request in requests {
    filter request.amount > 1000
    return request.id
}
```

In lowercase, these constructs look like parts of the same small programming language rather than programming constructs grafted onto SQL-like syntax.

That's much closer to how we think about Ferret today.

## Declarative doesn't have to mean SQL-like

Changing the visual style doesn't change Ferret's execution philosophy.

Ferret remains a declarative-first language.

The distinction matters.

Uppercase keywords don't make a language declarative, just as lowercase keywords don't make one imperative.

Consider:

```fql
for user in users {
    filter user.active
    return user.email
}
```

The program describes selection and transformation. It doesn't tell the runtime to iterate a collection, test a condition, skip an iteration, append a value to an accumulator, and continue.

That distinction is intentional.

Ferret does have imperative features. Mutable variables and `while` exist because real-world automation sometimes requires state.

Browser automation is a particularly good example:

```fql
let previous_count = 0
var done = false

for while !done {
    scroll(page)

    waitfor exists page.items

    let current_count = count(page.items)

    // ...
}
```

Infinite scrolling, pagination, changing DOM state, asynchronous events, and other interactions don't always fit neatly into a purely declarative model.

Ferret doesn't try to pretend otherwise.

The goal isn't declarative purity.

The goal is to provide enough state and iteration to deal with the real world without gradually turning Ferret into another general-purpose scripting language.

## Why `match` instead of `if`

This philosophy is also why Ferret v2 gained pattern matching rather than conventional `if`/`else` control flow.

For example:

```fql
let status = match response.status {
    200 => "ok"
    404 => "missing"
    status when status >= 500 => "server-error"
    _ => "unexpected"
}
```

Pattern matching makes branching value-oriented.

Similarly, filtering a collection doesn't require an `if` followed by `continue`:

```fql
return for user in users {
    filter user.active
    return user
}
```

Ferret already has a construct that directly expresses the intent.

This is a useful principle for the language going forward: before adding another general-purpose control-flow mechanism, ask whether the underlying operation deserves a more declarative primitive instead.

Sometimes the answer will still be an imperative feature. `while` is evidence of that.

But it shouldn't be the default.

## Scripts no longer have to return something

There is another AQL inheritance we're removing in v2.

Historically, every Ferret script had to produce a value.

That makes perfect sense for a query language: a query exists to calculate a result.

It makes considerably less sense for this:

```fql
click(submit)
waitfor exists confirmation
```

What meaningful value should this script return?

Previously, scripts like this had to manufacture one:

```fql
click(submit)
waitfor exists confirmation

return none
```

In Ferret v2, they don't.

Reaching the end of a script now means successful completion with `none`.

The explicit version remains valid:

```fql
return none
```

It simply isn't required.

This makes effect-oriented Ferret programs much more natural, particularly when Ferret is embedded as a DSL inside another application.

## Functions follow the same rule

User-defined functions now behave consistently with scripts.

A function that exists for its effects doesn't need a ceremonial return:

```fql
func notify(user) {
    send(user.email)
    audit(user.id)
}
```

If execution reaches the end of the function, it returns `none`.

Functions that produce values continue to use `return`:

```fql
func normalize(user) {
    return {
        name: trim(user.name),
        email: lower(user.email)
    }
}
```

And expression-bodied functions remain concise:

```fql
func square(x) => x * x
```

There is no implicit return of the last arbitrary expression in a block.

For example:

```fql
func square(x) {
    x * x
}
```

does not return `x * x`.

Block-bodied functions return values explicitly.

## `for` is no longer secretly a script return

This leads to one breaking change.

Historically, a final `for` loop could implicitly provide the result of a Ferret script:

```fql
for user in users {
    filter user.active
    return user.email
}
```

This was another useful behavior inherited from Ferret's query-language model.

But once scripts are allowed to complete without producing a value, making the final `for` special introduces an inconsistency.

A `for` already produces a value in Ferret. That value can be captured:

```fql
let emails = (
    for user in users {
        filter user.active
        return user.email
    }
)

return emails
```

The fact that the same construct happened to return from the script merely because it appeared at the end was a separate rule.

Ferret v2 removes that special case.

A standalone `for` remains value-producing, but an unused result is discarded.

If you want the loop result to become the result of the script, say so:

```fql
return for user in users {
    filter user.active
    return user.email
}
```

Besides being more explicit, this reads surprisingly well:

> return, for every active user, their email.

More importantly, `for` now means exactly the same thing wherever it appears.

It can be assigned:

```fql
let emails = (
    for user in users {
        return user.email
    }
)
```

returned:

```fql
return for user in users {
    return user.email
}
```

or used standalone when its result isn't needed:

```fql
for user in users {
    log(user)
}
```

Its position in a body no longer changes its semantics.

## One rule for executable bodies

These changes allow Ferret to have a much simpler result model.

For scripts and block-bodied functions:

- `return value` explicitly returns a value.
- `return for ...` explicitly returns the result of a loop.
- Reaching the end returns `none`.
- The result of a standalone value-producing construct is discarded.
- The last arbitrary expression is never implicitly returned.

Expression-bodied functions remain explicitly value-producing:

```fql
func double(value) => value * 2
```

There is no longer a special rule saying that every script must return something, nor another rule promoting a final `for` into the script result.

## What about existing Ferret programs?

Changing the canonical keyword casing does not break existing programs.

Uppercase remains valid:

```fql
FOR item IN items {
    RETURN item
}
```

The formatter will simply canonicalize it to:

```fql
for item in items {
    return item
}
```

The change to final `for` semantics, however, is intentionally breaking.

A v1-style script such as:

```fql
FOR item IN items
    RETURN item
```

that relies on the loop becoming the script result needs to become the v2 equivalent:

```fql
return for item in items {
    return item
}
```

You won't have to hunt these cases down manually.

The Ferret CLI migration command will recognize affected constructs and migrate them to their explicit v2 equivalents where possible.

The goal of the migration tooling is straightforward: **make semantic changes explicit without making upgrading existing Ferret programs unnecessarily painful.**

## A language, not just a query

None of these changes individually redefine Ferret.

Together, though, they reflect something that has been happening gradually throughout v2.

Ferret isn't trying to become a general-purpose programming language.

It also isn't just a query language anymore.

It's a small, declarative-first, expression-oriented language designed to be embedded into applications and specialized through host-provided capabilities.

It can query:

```fql
let users = query "active-users" in db
```

transform:

```fql
let names = (
    for user in users {
        return user.name
    }
)
```

make decisions:

```fql
let action = match request.kind {
    "create" => create(request)
    "update" => update(request)
    _ => none
}
```

and synchronize with the outside world:

```fql
click(submit)
waitfor exists confirmation
```

And when that's all the program needed to accomplish, it can simply end.

No artificial result required.

Ferret's AQL heritage is still visible in the language, and that's a good thing. Many of its best ideas grew from that foundation.

But v2 is also the right time to distinguish between the ideas that remain fundamental to Ferret and conventions that survived mostly because they've always been there.

Lowercase syntax is the most immediately visible sign of that transition.

The deeper change is simpler:

**Ferret is growing into its own language.**