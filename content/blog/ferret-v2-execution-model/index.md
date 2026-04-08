---
title: "Inside Ferret v2: The New Execution Model"
subtitle: "A New Foundation"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-04-08"
---

Hello friends!

The [previous post](/blog/ferret-v2-announcement/) covered why Ferret v2 exists - why a rewrite made sense, what changed philosophically, and why I decided it was time to move the project into a new chapter.

This time I want to talk about something more concrete: execution.

The new runtime isn't just a faster version of the old one. It's built on a different execution model entirely, and that change touches almost everything: performance, embeddability, tooling, future language features, and the kinds of optimizations Ferret can support.

## The wall

As I explained in the previous post, I wasn't trying to build some grand system. I needed something practical that worked, and I needed it quickly, so I built it that way. That was the right call. Ferret probably wouldn't exist otherwise.

But over time, as the project grew, I started to feel the limits of those early decisions. New features were still possible, but they often felt harder to implement than they should have been. Performance work was possible too, but too much of it depended on working around the existing shape of the runtime instead of improving it directly.

I could have kept patching it, refactoring around the edges, and squeezing out incremental gains. But at some point it became clear that the real issue wasn't one bad subsystem or one annoying abstraction - it was the execution model itself.

In v1, higher-level language behavior and runtime execution were too tightly intertwined. Evaluating a construct often meant mixing syntax-level intent, runtime control flow, temporary value management, and host integration in the same path. That made the system harder to reason about, harder to inspect, and much harder to optimize in a principled way.

What v2 changes is not just the implementation. It changes the shape of the system.

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

That separation matters more than it may seem at first. Once a query becomes an explicit program, it stops being just "something the runtime is currently doing" and becomes a real artifact with structure: instructions, constants, source locations, control flow, and exception regions. That gives Ferret something stable to execute, inspect, serialize, cache, and optimize.

## Register-based VM

There are different ways to design a virtual machine. Two common approaches are stack-based and register-based execution. Ferret v2 uses registers - instructions work with explicit slots instead of pushing and popping values from a stack.

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

In Ferret's case, instructions operate on explicit operands, usually a destination plus one or two source operands. That makes the bytecode easy to analyze and relatively easy to transform. It also means the VM does not have to infer intent from stack position or reconstruct temporary state on the fly. The compiler makes those decisions up front, and the VM executes the result.

## A tiny example

A small example makes the shift clearer.

Take a trivial query like this:

```fql
LET x = 1 + 2
RETURN x * 3
```

At the source level, this is simple enough that you barely think about it. But internally, v2 can treat it as a proper compilation problem:

1. Parse the source into a structured representation.
2. Lower the expressions into explicit operations.
3. Assign registers to intermediate values.
4. Fold constant work where possible.
5. Emit a compact bytecode program for the VM.

Conceptually, an early lowered form might look something like this:

```text
r1 = const 1
r2 = const 2
r3 = add r1, r2
r4 = const 3
r5 = mul r3, r4
return r5
```

And after simple optimization, it may become closer to this:

```text
r1 = const 3
r2 = const 9
return r2
```

That is intentionally simplified, but it shows the important part: the compiler can see the work explicitly. Once the query becomes a real program, constant folding, register reuse, and instruction cleanup become straightforward compiler passes instead of runtime guesswork.

Even when the query is far more dynamic than this, the same principle holds. The compiler turns higher-level language constructs into explicit operations, and the VM runs those operations against a well-defined execution state.

## The compilation pipeline

One of the most important changes in v2 is that the compiler is no longer a thin step on the way to execution - it's now a real part of the architecture.

The pipeline works roughly like this:

1. **Parse the query.** The parser turns FQL source into a structured representation of the program. This is the stage that understands syntax.

2. **Lower it into executable operations.** The compiler takes that structure and translates it into bytecode instructions for the VM. Higher-level language constructs become a lower-level execution plan.

3. **Apply optimizations.** Once the program exists in explicit form, the compiler can optimize it before it ever reaches the VM. This is one of the biggest practical benefits of the new design.

4. **Execute the bytecode.** The VM runs the compiled program using a register file and explicit execution state: call frames, scratch memory, exception handling, and other runtime machinery that is now much easier to reason about.

This is also where Ferret starts to benefit from having a real program representation instead of a loosely structured runtime path. Calls are explicit. Jumps are explicit. Catch regions are explicit. Source locations can be preserved. The compiler can emit something the VM can execute directly, and tooling can inspect that same representation without needing to reverse-engineer runtime behavior after the fact.

## What the VM actually owns

At runtime, the VM is no longer trying to "figure out" the query as it goes. Its job is narrower and cleaner.

It receives a compiled program and executes it against explicit state:
- a register file for values
- parameter and constant access
- call frames for function execution
- scratch state for intermediate runtime needs
- structured exception and catch handling
- source mapping data for diagnostics and tooling

That sounds like an internal detail, but it changes the feel of the entire system. The VM is no longer a place where language semantics, temporary bookkeeping, and host behavior blur together. It becomes an execution engine with a defined contract.

That makes behavior easier to trace, easier to test, and easier to evolve.

## Optimization

This is where the redesign starts to justify itself.

Because the program compiles down to bytecode with explicit registers, v2 can support optimizations that were previously hard, brittle, or simply not worth attempting.

Temporary values don't have to live forever - once the compiler knows a register is no longer needed, that slot gets recycled.

Expressions that can be partially or fully resolved at compile time get folded before execution begins.

Small instruction sequences get simplified into better equivalents through peephole optimization. That sounds minor on its own, but it adds up across a whole program.

Jumps, branches, calls, and catch regions are all represented explicitly, so the runtime can reason about control flow cleanly instead of treating it as an accidental byproduct of evaluation.

And as Ferret grows more expressive, function calls need to be cheap and predictable. A register-based VM gives me a much better starting point for that than the old model ever could.

The important part here is not just that optimizations are possible. It's that they now have a proper place to live. They can happen on a stable intermediate form before execution starts, rather than being forced into ad hoc runtime shortcuts.

## Future optimization headroom

Some of these optimizations are already in place. Others are still experimental or partial. But the important thing is that v2 has an architecture where they actually make sense to pursue.

Faster function dispatch, smarter register allocation, more precomputation during warmup, better call fast paths, and broader instruction-level simplification - these are all natural next steps now.

In v1, most of them were either awkward to implement or too fragile to force into the existing design.  
v2 isn't just something I can optimize today; it's something I can keep optimizing without fighting the runtime every step of the way.

That matters to me more than any single benchmark result. A runtime that improves once is nice. A runtime that can keep improving is much more valuable.

## Beyond speed

v2 is faster in many cases, but speed isn't really the point. The point is that the runtime now has a shape you can actually build on.

The internal architecture is clearer - the compiler compiles, the VM executes, and the boundary between them is clean. That makes the whole project easier to work on and easier to reason about.

Ferret isn't just a CLI tool. It's also a library that other systems embed, and a compiled program backed by a dedicated VM gives that embedding story a proper foundation.

Disassembly, debugging, inspection, serialization, caching, analysis - all of these become natural once execution is backed by a real program representation instead of a loose runtime path.

When execution is explicit, behavior is easier to trace, which helps with debugging, performance work, and confidence when adding new features.

And maybe most importantly, a weak execution model limits the language above it. A strong one lets the language expand without every new feature turning into a special case.

That is the real value of the rewrite. Not just that v2 can run queries better, but that it gives Ferret a foundation solid enough to support the language I want it to grow into.

## What's next

The move to a compiled pipeline and a register-based VM is one of the biggest technical changes in the project. It improves performance, but more importantly, it gives Ferret a real foundation to build on.

Next up: the language side of v2 - syntax changes, the overall direction, and how Ferret is becoming more expressive without losing what made it Ferret in the first place.

Thanks for reading.

### Useful links
- [The design of the Inferno virtual machine](https://inferno-os.org/inferno/papers/hotchips.pdf) - a short paper on register-based VM design that helped shape Ferret’s move to a register-based architecture.
- [The Implementation of Lua 5.0](https://www.lua.org/doc/jucs05.pdf) - a great paper on the design of the Lua 5 virtual machine and a very helpful resource for understanding how register-based VMs work.
- [Virtual Machine Showdown: Stack Versus Registers](https://www.usenix.org/legacy/events/vee05/full_papers/p153-yunhe.pdf) - a great paper on the tradeoffs between stack-based and register-based VMs.
- [Crafting Interpreters](https://craftinginterpreters.com/) - a great book on building interpreters and virtual machines, and a valuable resource for anyone interested in compilers and language implementation.
- [The Expr expression language](https://github.com/expr-lang/expr) - a very nice expression language for Go with a stack-based VM. I used it as a reference in the early stages of the project.
