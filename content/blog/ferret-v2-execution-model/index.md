---
title: "Inside Ferret v2: The New Execution Model"
subtitle: "A New Foundation"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-04-07"
---

Hello friends!

The [previous post](/blog/ferret-v2-announcement/) covered why Ferret v2 exists - why a rewrite made sense, what changed philosophically, and why I decided it was time to move the project into a new chapter.

This time I want to talk about something more concrete: execution.

The new runtime isn't just a faster version of the old one. It's built on a different execution model entirely, and that change touches almost everything: performance, embeddability, tooling, future language features, and the kinds of optimizations Ferret can support.

## The wall

As I explained in the previous post, I wasn't trying to build some grand system. I needed something practical that worked, and I needed it quickly, so I built it that way, and that was the right call. Ferret probably wouldn't exist otherwise.

But over time, as the project grew, I started to feel the limits of those early decisions. New features were still possible, but they often felt harder to implement than they should have been. Performance work was possible too, but too much of it depended on working around the existing shape of the runtime instead of improving it directly.

I could have kept patching it, refactoring around the edges, and squeezing out incremental gains. But at some point it became clear that the real issue wasn't one bad subsystem or one annoying abstraction - it was the execution model itself.

## The conceptual shift

Ferret v2 introduces a compiled execution pipeline built around bytecode and a register-based virtual machine.

At a high level, the flow now looks like this:

```text
source -> parser -> compiler -> bytecode program -> VM execution
```

That may sound like an implementation detail, but it defines how the entire system works.

In v1, execution and language behavior were tightly coupled - the system interpreted and processed queries in one tangled pass.
v2 draws a clear boundary between understanding the source code, compiling it into an executable form, and running that form efficiently.

The compiler can spend time shaping the program once. The VM can focus on running it well. And tooling can inspect the result in ways that were much harder to support before.

## Register-based VM

There are different ways to design a virtual machine. Two common approaches are stack-based and register-based execution. 
Ferret v2 uses registers - instructions work with explicit slots instead of pushing and popping values from a stack.

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

The difference matters in a few ways.
Data flow becomes explicit - you can see where values come from, where they go, and when they can be reused. That makes the system easier to reason about.

A stack machine spends a lot of effort shuffling values around indirectly. A register machine can express the same work with fewer intermediate steps.

But the biggest reason I went in this direction is optimization. Once values live in named registers, the compiler can track temporaries, reuse slots, avoid unnecessary moves, and simplify expressions earlier.

It doesn't magically solve every performance problem, but it gives the runtime a shape that is far friendlier to optimization work down the road.

## The compilation pipeline

One of the most important changes in v2 is that the compiler is no longer a thin step on the way to execution - it's now a real part of the architecture.

The pipeline works roughly like this:

1. **Parse the query.** The parser turns FQL source into a structured representation of the program. This is the stage that understands syntax.

2. **Lower it into executable operations.** The compiler takes that structure and translates it into bytecode instructions for the VM. Higher-level language constructs become a lower-level execution plan.

3. **Apply optimizations.** Once the program exists in explicit form, the compiler can optimize it before it ever reaches the VM - one of the biggest practical benefits of the new design.

4. **Execute the bytecode.** The VM runs the compiled program using a register file and explicit execution state: call frames, scratch memory, exception handling, and other runtime machinery that's now straightforward to reason about.

## Optimization

This is where the redesign starts to pay off.

Because the program compiles down to bytecode with explicit registers, v2 can support optimizations that were previously hard, brittle, or simply not worth attempting.

Temporary values don't have to live forever - once the compiler knows a register is no longer needed, that slot gets recycled. 
Expressions that can be partially or fully resolved at compile time get folded before execution begins. 
Small instruction sequences get simplified into better equivalents through peephole optimization, which sounds minor on its own but adds up across a whole program. 
Jumps, branches, calls, and catch regions are all represented explicitly, so the runtime can reason about control flow cleanly. 

And as Ferret grows more expressive, function calls need to be cheap and predictable - a register-based VM gives me a far better starting point for that than the old model ever could.

## Future optimization headroom

Some of these optimizations are already in place, others are still experimental or partial. But the important thing is that v2 has an architecture where they actually make sense to pursue.

Faster function dispatch, smarter register allocation, more precomputation during warmup - these are all natural next steps now. 

In v1, most of them were either awkward to implement or too fragile to force into the existing design. 
v2 isn't just something I can optimize today; it's something I can keep optimizing without fighting the runtime every step of the way.

## Beyond speed

v2 is faster in many cases, but speed isn't really the point. The point is that the runtime now has a shape you can actually build on.

The internal architecture is clearer - the compiler compiles, the VM executes, and the boundaries between them are clean, which makes the whole project easier to work on and reason about. 

Ferret isn't just a CLI tool; it's also a library that other systems embed, and a compiled program backed by a dedicated VM gives that embedding story a proper foundation. 
Disassembly, debugging, inspection, serialization, caching, analysis - all of these become natural once execution is backed by a real program representation instead of a loose runtime path. 
When execution is explicit, behavior is easier to trace, which helps with debugging, performance work, and confidence when adding new features.

And maybe most importantly: a weak execution model limits the language above it. A strong one lets the language expand without every new feature turning into a special case. 

## What's next

The move to a compiled pipeline and a register-based VM is one of the biggest technical changes in the project. It improves performance, but more importantly, it gives Ferret a real foundation to build on.

Next up: the language side of v2 - syntax changes, the overall direction, and how Ferret is becoming more expressive without losing what made it Ferret in the first place.

Thanks for reading.

### Useful links
- [The design of the Inferno virtual machine](https://inferno-os.org/inferno/papers/hotchips.pdf) - a short paper on register-based VM design that helped shape Ferret's move to a register-based architecture.
- [The Implementation of Lua 5.0](https://www.lua.org/doc/jucs05.pdf) -  a great paper on the design of the Lua 5 virtual machine and a very helpful resource for understanding how register-based VMs work.
- [Virtual Machine Showdown: Stack Versus Registers](https://www.usenix.org/legacy/events/vee05/full_papers/p153-yunhe.pdf) - a great paper on the tradeoffs between stack-based and register-based VMs.
- [Crafting interpreters](https://craftinginterpreters.com/) - a great book on building interpreters and virtual machines, and a valuable resource for anyone interested in compilers and language implementation.
- [The Expr expression language](https://github.com/expr-lang/expr) - a very nice expression language for Go with a stack-based VM. I used it as a reference in the early stages of the project.