---
title: "Inside Ferret v2: The New Execution Model"
subtitle: "A New Foundation"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-03-31"
---

Hello friends!

In the [previous post](/blog/ferret-v2-announcement/), I talked about why Ferret v2 exists at all - why a rewrite made sense, what changed philosophically, and why I decided it was time to move the project into a new chapter.

This time I want to talk about something much more concrete: execution.

Because if v2 has a heart, this is it.

The new runtime is not just a faster version of the old one. It is built on a different execution model entirely, and that change affects almost everything - performance, embeddability, tooling, future language features, and the kinds of optimizations Ferret can do.

## The wall

Ferret v1 did its job.

It solved the original problem, it was useful, and it carried the project a lot further than I expected when I first started it.

But it was also very much a product of how the project started. I wasn’t trying to build some grand system from day one. I just needed something practical that worked, and I needed it sooner rather than later.

So I built it that way.

And honestly, that was the right decision at the time. Ferret probably would not exist otherwise.

But over time, as the project grew, I started feeling the limits of those early decisions more and more. New features were still possible, but they often felt harder to implement than they should have been. Performance work was possible too, but it often felt too local, too fragile, too dependent on working around the existing shape of the runtime instead of improving it directly.

In other words, v1 was still working, but it was also starting to push back.

That is usually a sign.

I could have kept patching it, refactoring around the edges, and squeezing out incremental improvements. But at some point it became clear that the real issue was deeper than that. The problem was not one bad subsystem or one annoying abstraction. It was the execution model itself.

That was the moment the question changed from “how do I keep evolving v1?” to “what do I actually want Ferret’s runtime to look like?”

## The conceptual shift

Ferret v2 introduces a compiled execution pipeline built around bytecode and a register-based virtual machine.

At a high level, the flow now looks like this:

```text
source -> parser -> compiler -> bytecode program -> VM execution
```

That may sound like an implementation detail, but it is much more than that.

In v1, execution and language behavior were much more tightly coupled to the way the system interpreted and processed queries at runtime. In v2, there is a clearer boundary between:

- understanding the source code
- compiling it into an executable form
- running that executable form efficiently

This separation gives the runtime a much more stable and explicit foundation.

The compiler can spend time shaping the program once.
The VM can focus on executing that program well.
And tooling can inspect the result in a way that was much harder before.

## Register-based VM

There are different ways to design a virtual machine. Two common approaches are stack-based and register-based execution models.

Ferret v2 uses registers.

That means instructions work with explicit slots instead of pushing and popping values implicitly from a stack.

Conceptually, instead of doing something like this:

```text
PUSH a
PUSH b
ADD
PUSH c
MUL
```

you can do something closer to this:

```text
r1 = a
r2 = b
r3 = add r1, r2
r4 = c
r5 = mul r3, r4
```

This has a few important advantages.

First, it makes data flow more explicit. You can see where values come from, where they go, and when they can be reused.

Second, it tends to reduce instruction churn. A stack machine often spends a lot of effort moving values around indirectly. A register machine can express the same work with fewer intermediate steps.

Third, and this was a big one for me, it makes optimization much more practical.

Once values live in explicit registers, the compiler has far more room to reason about them. It can track temporary values better, reuse registers, avoid unnecessary moves, and simplify expressions earlier and more safely.

That does not magically solve every performance problem, of course. But it gives Ferret a runtime shape that is much friendlier to optimization.

## The compilation pipeline

One of the most important changes in v2 is that the compiler is no longer just a thin step on the way to execution. It is now a real part of the architecture.

The job of the pipeline is roughly this:

- **Parse the query**

  The parser turns FQL source into a structured representation of the program.

  This is the stage that understands syntax.

- **Lower it into executable operations**

  The compiler takes that structure and translates it into bytecode instructions for the VM.

  This is where higher-level language constructs become a lower-level execution plan.

- **Apply optimizations**

  Once the program exists in a more explicit form, the compiler can optimize it before it ever reaches the VM.

  This is one of the biggest practical benefits of the new design.

- **Execute the bytecode**

  The VM runs the compiled program using a register file and a much more explicit execution state.

  That includes not only the instruction stream itself, but also things like call frames, scratch memory, exception handling, and other runtime details that are now much easier to reason about directly.

## Optimization

This is where the redesign starts paying rent.

Because the program is now compiled into bytecode with explicit registers, v2 can support kinds of optimization that were previously hard, brittle, or simply not worth attempting.

Some examples:

- **Register reuse**

  Temporary values do not have to live forever. Once the compiler knows a register is no longer needed, that slot can be reused later.

  This reduces unnecessary churn and keeps execution more efficient.

- **Constant propagation and folding**

  If an expression can be partially or fully resolved at compile time, v2 can do that before execution begins.

  That means less work for the VM and fewer unnecessary runtime operations.

- **Peephole optimization**

  Small instruction sequences can be simplified into better ones.

  This is one of those techniques that sounds minor, but across a whole program it adds up quickly.

- **Cleaner control flow**

  Because jumps, branches, calls, and catch regions are represented explicitly, the runtime can reason about control flow much more cleanly than before.

  That matters not just for speed, but also for correctness and maintainability.

- **A better function execution model**

  As Ferret grows more expressive, function calls need to be cheap and predictable.

  A register-based VM gives me a much better foundation for that than the old execution model did.

## Future optimization headroom

This part is important.

There are optimizations that are already in place, and there are others that are still experimental, partial, or just ideas for now. But the important thing is that v2 finally has an architecture where those ideas actually make sense.

Things like faster function calls, better register reuse, more precomputation during warmup, and other runtime improvements are now much more natural to implement. In v1, many of these things were either awkward, too fragile, or simply not worth forcing into the existing design.

That is a big difference.

v2 is not just something I can optimize today. It is something I can keep optimizing tomorrow without feeling like I am fighting the runtime every time.

## Beyond speed

It would be very easy to reduce all of this to “v2 is faster.”

And yes, in many cases, it is.

But that is not really the main point.

The bigger change is that the runtime now has a much better foundation.

A better execution model gives Ferret:

- **A clearer internal architecture**

  The compiler compiles. The VM executes. The boundaries are easier to understand.

  That sounds obvious, but in practice it makes the whole project easier to evolve.

- **Better embeddability**

  Ferret is not just a CLI tool. It is also a library and a runtime that other systems can embed.

  A compiled program and a dedicated VM give that embedding story a much cleaner foundation.

- **Better tooling**

  Disassembly, debugging, inspection, serialization, caching, and analysis all become much more natural once execution is backed by a real program representation instead of a looser runtime path.

  This matters not only for contributors, but eventually for users too.

- **Better predictability**

  When execution is explicit, behavior becomes easier to reason about.

  That matters for debugging, for performance work, and for confidence when adding new language features.

- **More room for the language to grow**

  A weak execution model limits the language above it.

  A stronger execution model gives the language room to expand without every new feature turning into a special case.

  That is one of the biggest reasons I made this change. I do not want Ferret v2 to merely inherit the old system with cleaner code.

  I want it to have a foundation that can support the next several years of ideas without constantly fighting back.

## The great rewrite

This is probably one of the best examples of why v2 had to be a rewrite rather than a gradual cleanup.

A lot can be improved incrementally. You can refactor code, redesign APIs, move packages around, and make the internals cleaner over time.

But there is a limit to that approach.

Some changes are architectural. Once you hit that point, preserving the old structure too much just makes the new design worse. Instead of building the system you actually want, you end up negotiating with the one you already have.

That was exactly the problem here.

The old execution model served Ferret well. It helped the project grow, and it got me pretty far.

But v2 needed a different foundation. I did not want it to remain shaped by decisions that belonged to an earlier stage of the project - and honestly, to an earlier stage of me as an engineer too.

## Closing thoughts

If the first post was about why Ferret v2 exists, this one was about what changed under the hood.

The move to a compiled pipeline and a register-based VM is one of the biggest technical changes in the project.

It improves performance, yes. But more importantly, it gives Ferret a much clearer foundation to build on.

In the next post, I’ll probably talk about the language side of v2 - the syntax changes, the overall direction, and how Ferret is becoming more expressive without losing what made it Ferret in the first place.

Thanks for reading.
