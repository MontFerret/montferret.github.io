---
title: "On the Road to Ferret v2"
subtitle: "A New Chapter"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-03-26"
---

Hello friends!

Long time no see! I hope you all have been doing well, it's been a while since my last post. 

Given the Ferret v2 alpha release, I thought it was a good time to share some thoughts about the journey that led us here, the design choices we made, and what you can expect from this new version.

## How We Got Here

When I first started working on Ferret, I didn’t have some grand plan. I also didn’t know much about parsers, compilers, or runtimes.
I just needed a solution to a very specific pain.

After spending some time looking for solutions and not finding any, I decided to create one.

I wanted something that would let me scrape websites without writing a thousand lines of glue code every time I needed to parse a new page.
I wanted to be able to say "give me all the links on this page, filter them by this condition, and transform them into this format" without dealing with the underlying mechanics.

So I built Ferret v1. A small language that let me describe what I needed.
Despite being pretty simple with a naive implementation, it worked pretty well. And to my big surprise, other people started relying on it for real work.

And this is when the real fun started.
Because once people start using your tool — they use it in ways you never expected.
And then they ask for features you never thought of.
They report bugs you never imagined.
They give you feedback that challenges your assumptions.

## The Slow Realization

After a while, I started noticing a pattern: people kept asking for things that made perfect sense, but adding them properly was more painful than it had any right to be. Not because the features were that hard, but because Ferret v1 was not built to grow in those directions.

I'm not a compiler or VM expert. I just wanted to solve my own problem and built a tool that worked for me. Ferret was my first serious attempt at building a language and a runtime, so a lot of the early design decisions came from limited knowledge and whatever felt right at the time.

As the tool grew, some of those decisions started getting in the way.

A big part of Ferret v1 was also shaped by AQL. That gave me a strong starting point and helped make the language practical early on. I did not have to invent everything from zero, which was definitely a good thing.

But those kinds of influences do not come for free. They shape how the language thinks. They define what feels natural, and over time they can make other ideas feel harder to express cleanly.

Eventually I came to an unpleasant conclusion: every time I wanted to make Ferret cleaner, faster, or more expressive, I was paying a tax to preserve decisions that made sense years ago.

And those kinds of taxes tend to compound.

If a tool is useful, people naturally try to do bigger things with it. At some point, they start running into the edges.

In Ferret v1, too many of those edges were structural:

- some performance ceilings were not about a missed optimization - they were about the execution model
- some ergonomics problems were not about nicer messages - they were about how the compiler understood code
- some language evolution problems were not about choosing syntax - they were about semantics already being crowded

At some point, patching stops being enough. You have to step back and rethink the foundation.

## Why a rewrite

Rewrites have a bad reputation, and not without reason. It is easy to spend a long time rebuilding something, only to end up with a cleaner system that is somehow missing the small practical details that made the old one useful.

And yes, I did disappear for a while.

But this rewrite was never about chasing perfection. It was about building a foundation that does not fight me every time I try to improve something.

Ferret v1 could be improved incrementally, and it was. But I eventually reached a point where every genuinely good change came with too much friction from the past: backwards compatibility quirks, subtle behavior coupling, performance tradeoffs, and special-case logic that I knew would come back again with the next feature.

Ferret v2 is my attempt to build a foundation where the next five years of work do not feel like surgery.

## What Ferret v2 Is

Ferret v1 started as a tool for scraping data from websites. That is still the core of what Ferret does.

But even a small language does not have to stay frozen. It can grow control flow, diagnostics, tooling, and extension points without losing its identity.

I am not trying to turn Ferret into JavaScript. I am trying to make it feel more like itself - just less fragile, more capable, and easier to evolve.

Ferret v2 is a complete rewrite of both the language and the runtime. It keeps the same name, much of the same syntax, and the same core ideas - but the implementation underneath is entirely different.

### Execution Model

Ferret v2 is built around a register-based VM and a proper compilation pipeline. That may sound like an internal detail, but it changes a lot of what the language can do well.

In v1, many things were effectively assembled as the program ran. That worked, but it also meant there were limits to how much the system could understand, optimize, and improve ahead of time.

v2 gives Ferret a much stronger foundation. Instead of treating execution as a chain of loosely connected steps, it can compile queries into a form the runtime can reason about more directly. That makes real optimization possible.

In practice, this means common patterns can be improved at compile time instead of being stitched together at runtime. Things like filtering, looping, expressions, and other repeated operations can be executed more efficiently without the user having to write code differently.

This is not just about speed, although speed matters. It is also about making the language more predictable, easier to evolve, and better suited for future features that would have been awkward or expensive to support in v1.

### Error Handling

One of the most frustrating parts of Ferret v1 was error reporting.

Sometimes something would go wrong and you would get a message like "no viable alternative" plus a line number that was not especially helpful. Technically, the tool was telling you that something failed. In practice, it often did very little to help you understand what actually went wrong.

Ferret v2 has a much better diagnostic system: span-based errors, labeled highlights, better context rendering, and smarter parse-time error rewriting.

The goal is simple - when something breaks, you should be able to see what happened, where it happened, and have a much better chance of fixing it without digging through guesswork.

### Extensibility

Ferret v1 had some extension points, but they were limited and ad hoc.

v2 leans toward a more capability-based runtime design, where values can expose behavior without everything collapsing into "just a map." 

This is what makes future things plausible — storage backends, remote objects, plugins, richer integrations without turning Ferret's core into a kitchen sink.

### Tooling

Ferret v1 didn't really have tooling infrastructure. Things like formatting or debugging were either missing or bolted on after the fact.

v2 has a compiler pipeline that actually makes those things buildable. A formatter, a debugger, a language server, and a REPL are all on the roadmap.

## What this means for you

v2 is not me abandoning v1. v1 proved the idea, and it's the reason I can design v2 with confidence instead of guesses.

But if you've ever felt that Ferret is almost right for your workflow — except for performance, tooling, clarity, or expressiveness — v2 is built to address those at the root, not at the surface.

Over time, you should see:
- better performance on real pipelines
- much better diagnostics and debugging
- cleaner ways to express complex logic without losing the "query" feel
- a foundation for modules and integrations that doesn't require hacks

And yes, a few breaking changes — because that's the price of making the system simpler.

## What's Next

As of today, Ferret v2 is in alpha.

It is not ready for production use yet, but it is ready to be explored, tested, and pushed on. If you want to try it, you can check out the core and CLI repositories on GitHub.

The core repository contains the language and runtime. The CLI repository contains the command-line interface for running Ferret scripts.

The next few posts will be more concrete. I want to write in more detail about:

- new language features such as MATCH, functions/UDFs, scoping rules, and query dialects
- the execution model, including VM choices and what gets optimized
- modules and capabilities, and how extensions plug into the system
- tooling, including the formatter, assembler, debugger, and language server

If you are already using Ferret - or tried it and bounced off - this is the best time to tell me what hurt, what felt awkward, and what you wanted but could not do.

v2 is still at the stage where that kind of feedback can shape the foundation.

Thank you for sticking around.

I am excited about what comes next.