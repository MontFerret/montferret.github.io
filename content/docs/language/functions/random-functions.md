---
title: "Random functions"
sidebarTitle: "Random"
weight: 72
draft: false
description: "Pseudo-random values, exact interval boundaries, and reproducible Session seeding."
---

# Random functions

Use `random::` for ordinary automation, sampling, randomized behavior, and tests.
These functions require a runtime containing the Random API. This website's
generated reference follows its published runtime pin and will include the API
after the corresponding release is published and the pin is updated.

{{< code lang="fql" >}}
RETURN {
    fraction: random::float(),
    temperature: random::float(-10, 30),
    die: random::int(1, 6),
    enabled: random::bool()
}
{{</ code >}}

| Call | Result |
| --- | --- |
| `random::float()` | Float in `[0, 1)` |
| `random::float(min, max)` | Float in `[min, max)` |
| `random::int(min, max)` | Int in `[min, max]` |
| `random::bool()` | Boolean, with equal probability for true and false |
| `random::choice(values)` | One uniformly selected original value, or `none` for an empty list |
| `random::shuffle(values)` | New list preserving the input backend, with the original values in randomized order |

`[min, max)` includes the minimum and excludes the maximum; `[min, max]`
includes both. Canonical bounds use minimum-first ordering. The Random standard
library group exposes these functions without global aliases; Full and Safe
include it.

## Float bounds

The ranged form accepts native Int and Float bounds, including mixed types.
Both must be finite and the minimum must not exceed the maximum. NaN,
infinities, strings, Booleans, and other nonnumeric arguments are rejected.
There is no one-argument form.

Equal bounds return ordinary Float conversion of the bound without drawing:
`random::float(5, 5)` returns Float `5`. Float conversion can round large integers.
For unequal bounds, every result stays inside the original numeric interval,
even when integer bounds cannot be represented exactly as Floats. If no Float
fits, the call raises an argument error. For example, the interval
`[9007199254740993, 9007199254740994)` contains no representable Float.

Adjacent Float bounds are valid: the lower endpoint is the only possible result.
Extreme finite bounds are supported without overflowing their difference into
an invalid result.

## Integer bounds

Both bounds must be Int; even whole-valued Float arguments are rejected.
Use `random::int(1, 6)` for a six-sided die and `random::int(0, 1)` for either
zero or one. The inclusive range covers the complete supported Int domain,
including its minimum and maximum, without overflow or modulo bias.
`random::int(5, 5)` returns Int `5` without drawing. Reversed bounds fail.

## Choosing a value

`random::choice(values)` selects one value uniformly from a list. It returns the
original value without copying or changing it, and leaves the source list
unchanged. An empty list returns `none`.

{{< code lang="fql" >}}
RETURN random::choice(["chromium", "firefox", "webkit"])
{{</ code >}}

Choice traverses the list once using reservoir sampling. It supports host Lists
that provide traversal without random access or copying, takes O(n) time, and
uses O(1) additional storage. Repeated values remain separate candidates.

## Shuffling a list

`random::shuffle(values)` returns a new list with the same values and
multiplicities in randomized order. It preserves the source list's implementation
family and backend configuration. The source list is not modified; nested objects
and host values retain their original identity. Array inputs produce Arrays.

{{< code lang="fql" >}}
LET original = [1, 2, 3, 4]
RETURN {
    original,
    shuffled: random::shuffle(original)
}
{{</ code >}}

Shuffle uses the source List's factory to create an independent destination,
materializes the complete input into that destination, and applies Fisher-Yates
shuffling. A storage-backed List retains its backend instead of being forced into
an in-memory Array. There are O(n) append/swap operations; their cost and the
destination's storage requirements depend on the backend. The original order is a
valid possible result.

Host Lists must support `New` and traversal on the source, and `Append`, `Length`,
and `Swap` on the destination. A traversal-only List without a working factory
can still be used with choice. Shuffle propagates host errors and cancellation;
after successful creation, a failed operation closes a closable destination and
preserves cleanup errors. Successful destinations remain open for normal result
ownership. The source and its values are never explicitly closed.

Both operations accept Lists, propagate traversal errors, and consume no random
numbers for empty or single-element inputs. Empty inputs still create an
independent empty list in the same implementation family; `random::shuffle([])`
returns an empty Array.

## Reproducible execution

Every Session owns an independent source shared by all these functions,
deprecated `rand`, and `WAITFOR` jitter. Go embedders can use
[`ferret.WithSessionRandomSeed`]({{< ref "/docs/embedding/go/configuration" >}}#reproducible-randomness)
for reproducible executions and tests. The same seed and executed inputs reproduce
a sequence within one Ferret version. Repeated runs of one Session and debugger
resumes continue it; a new equally seeded Session starts it again.

An unseeded source obtains operating-system entropy on its first actual draw.
Invalid arguments, equal canonical bounds, and empty or single-element choice
and shuffle inputs neither initialize it nor consume randomness. Queries that
never draw do not read entropy. Timing and external inputs may change which
calls execute. Exact sequences are not
promised across Ferret versions. No `random::seed()` mutation function exists.

## Migrating from deprecated rand

Global `rand` remains available through the Math group with its historical
contract. It uses the same Session source as the canonical functions.

| Legacy call | Preserved behavior | Migration |
| --- | --- | --- |
| `rand()` | Float in `[0, 1)` | `random::float()` |
| `rand(x)` | Rounded calculation using minimum `x/2` and maximum `x*2` | Choose explicit bounds and the desired canonical result type |
| `rand(max, min)` | Maximum-first rounded calculation | Reverse arguments; use `random::int(min, max)` for inclusive Int bounds or `random::float(min, max)` for continuous bounds |

Legacy ranged results are Floats computed as `floor(u*(max-min+1))+min`, where
`u` is in `[0, 1)`. They retain permissive Float conversion, fractional results
when the minimum is fractional, and historical reversed/non-finite behavior.
Every successful legacy call draws once, including equal bounds. Canonical
functions apply their own strict types and bounds; they are not direct aliases
for these historical calculations. No `math::rand` or `math::random` exists.

## Security-sensitive randomness

Do not use `random::` for passwords, authentication tokens, secrets,
cryptographic keys, or other security-sensitive values. Operating-system
seeding does not make its pseudo-random output secure. Secure randomness belongs
to Ferret's `crypto::` functionality.
