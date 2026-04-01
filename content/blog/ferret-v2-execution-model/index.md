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

This time I want to talk about something more concrete: execution.

If v2 has a heart, this is it.

The new runtime isn't just a faster version of the old one. It's built on a different execution model entirely, and that change touches almost everything: performance, embeddability, tooling, future language features, and the kinds of optimizations Ferret can support.

## The wall

Ferret v1 did its job.

It solved the original problem, it was useful, and it carried the project further than I expected when I first started it.

But it was also very much a product of how the project began. I wasn't trying to build some grand system from day one. I needed something practical that worked, and I needed it soon.

So I built it that way. And that was the right call. Ferret probably wouldn't exist otherwise.

Over time though, as the project grew, I started feeling the limits of those early decisions. New features were still possible, but they often felt harder to implement than they should have been. Performance work was possible too, but it felt too local - too dependent on working around the existing shape of the runtime instead of improving it directly.

v1 was still working, but it was starting to push back. That's usually a sign.

I could have kept patching it, refactoring around the edges, squeezing out incremental gains. But at some point it became clear that the real issue ran deeper. Not one bad subsystem or one annoying abstraction - the execution model itself.

The question shifted from "how do I keep evolving v1?" to "what do I actually want Ferret's runtime to look like?"

## The conceptual shift

Ferret v2 introduces a compiled execution pipeline built around bytecode and a register-based virtual machine.

At a high level, the flow now looks like this:

```text
source -> parser -> compiler -> bytecode program -> VM execution
```

That may sound like an implementation detail, but it reshapes how the whole system works.

In v1, execution and language behavior were tightly coupled - the system interpreted and processed queries in one tangled pass. In v2, there's a clear boundary between:

- understanding the source code
- compiling it into an executable form
- running that form efficiently

The compiler can spend time shaping the program once. The VM can focus on running it well. And tooling can inspect the result in ways that were hard to pull off before.

## Register-based VM

There are different ways to design a virtual machine. Two common approaches are stack-based and register-based execution.

Ferret v2 uses registers. Instructions work with explicit slots instead of pushing and popping values from a stack.

Conceptually, instead of this:

```text
PUSH a
PUSH b
ADD
PUSH c
MUL
```

you get something closer to this:

```text
r1 = a
r2 = b
r3 = add r1, r2
r4 = c
r5 = mul r3, r4
```

Why does this matter?

**Data flow becomes explicit.** You can see where values come from, where they go, and when they can be reused.

**Less instruction churn.** A stack machine spends a lot of effort shuffling values around indirectly. A register machine can express the same work with fewer intermediate steps.

**Optimization becomes practical.** This was the big one for me. Once values live in named registers, the compiler can track temporaries, reuse slots, avoid unnecessary moves, and simplify expressions earlier. It doesn't magically solve every performance problem, but it gives the runtime a shape that's far friendlier to optimization work.

## The compilation pipeline

One of the most important changes in v2 is that the compiler is no longer a thin step on the way to execution. It's a real part of the architecture now.

The pipeline works roughly like this:

- **Parse the query.** The parser turns FQL source into a structured representation of the program. This is the stage that understands syntax.

- **Lower it into executable operations.** The compiler takes that structure and translates it into bytecode instructions for the VM. Higher-level language constructs become a lower-level execution plan.

- **Apply optimizations.** Once the program exists in explicit form, the compiler can optimize it before it ever reaches the VM. This is one of the biggest practical benefits of the new design.

- **Execute the bytecode.** The VM runs the compiled program using a register file and explicit execution state - call frames, scratch memory, exception handling, and other runtime machinery that's now straightforward to reason about.

## Optimization

This is where the redesign starts to pay off.

Because the program compiles down to bytecode with explicit registers, v2 can support optimizations that were previously hard, brittle, or simply not worth attempting.

Some examples:

**Register reuse.** Temporary values don't have to live forever. Once the compiler knows a register is no longer needed, that slot gets recycled. Less waste, tighter execution.

**Constant propagation and folding.** If an expression can be partially or fully resolved at compile time, v2 does it before execution begins. Less work for the VM at runtime.

**Peephole optimization.** Small instruction sequences get simplified into better ones. Sounds minor, but across a whole program it adds up fast.

**Cleaner control flow.** Jumps, branches, calls, and catch regions are all represented explicitly. The runtime can reason about control flow cleanly - which matters for correctness and maintainability, not just speed.

**Cheaper function calls.** As Ferret grows more expressive, function calls need to be cheap and predictable. A register-based VM gives me a far better starting point for that than the old model ever could.

## Future optimization headroom

Some of the optimizations above are already in place. Others are still experimental, partial, or just ideas. But the important thing is that v2 has an architecture where those ideas actually make sense.

Faster function dispatch, smarter register allocation, more precomputation during warmup - these are all natural next steps now. In v1, most of them were either awkward to implement or too fragile to be worth forcing into the existing design.

v2 isn't just something I can optimize today. It's something I can keep optimizing tomorrow without fighting the runtime every step of the way.

## Beyond speed

It would be easy to reduce all of this to "v2 is faster." And yes, in many cases it is. But speed isn't really the point.

The point is that the runtime now has a better shape. Here's what that means in practice:

**A clearer internal architecture.** The compiler compiles. The VM executes. The boundaries are clean. This makes the whole project easier to work on and easier to reason about.

**Better embeddability.** Ferret isn't just a CLI tool - it's also a library that other systems can embed. A compiled program and a dedicated VM give that embedding story a proper foundation.

**Better tooling.** Disassembly, debugging, inspection, serialization, caching, analysis - all of these become natural once execution is backed by a real program representation instead of a loose runtime path. This matters for contributors today and eventually for users too.

**Better predictability.** When execution is explicit, behavior is easier to trace. That helps with debugging, performance work, and confidence when adding new features.

**More room for the language to grow.** A weak execution model limits the language above it. A strong one lets the language expand without every new feature turning into a special case. This is one of the biggest reasons I made this change. I don't want v2 to be the old system with cleaner code. I want it to support the next several years of ideas without constantly fighting back.

## The great rewrite

This is probably the best example of why v2 had to be a rewrite rather than a gradual cleanup.

A lot can be improved incrementally - you can refactor code, redesign APIs, move packages around, clean up internals over time.

But some changes are architectural. Once you hit that point, preserving the old structure too aggressively just makes the new design worse. Instead of building what you actually want, you end up negotiating with what you already have.

The old execution model served Ferret well. It helped the project grow and it got me pretty far. But v2 needed a different starting point - one not shaped by decisions that belonged to an earlier stage of the project, and an earlier version of me as an engineer.

## What's next

The move to a compiled pipeline and a register-based VM is one of the biggest technical changes in the project. It improves performance, but more importantly, it gives Ferret a real foundation to build on.

In the next post, I'll talk about the language side of v2 - syntax changes, the overall direction, and how Ferret is becoming more expressive without losing what made it Ferret in the first place.

Thanks for reading.
