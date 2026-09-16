---
title: "Math collection functions"
sidebarTitle: "Math collections"
weight: 70
draft: false
description: "Canonical math functions, strict numeric inputs, and migration from deprecated globals."
---

# Math collection functions

The canonical Ferret v2 mathematical API uses `math::`. Use a runtime containing
this migration before running these examples. Existing global functions remain
available as deprecated compatibility APIs with their earlier behavior.

The generated [Math reference]({{< ref "/docs/standard-library/math" >}}) describes
the published runtime pinned by this website. It will reflect this migration
after a release containing it is published and the website runtime pin is updated.

{{< code lang="fql" >}}
LET values = [1, 2, 3, 100]

RETURN {
    sum: math::sum(values),
    mean: math::mean(values),
    median: math::median(values)
}
// { "sum": 106, "mean": 26.5, "median": 2.5 }
{{</ code >}}

## Scalar math

The following additions are available only under `math::`, with no global
aliases. They require a runtime containing the expanded canonical math API.
Arguments must be integers or floats; strings, booleans, and `none` are not
coerced into numbers.

| Function | Result |
| --- | --- |
| `math::clamp(value, min, max)` | Selects the original value or a bound, preserving its integer or float type. Requires `min <= max` and rejects NaN bounds. |
| `math::sign(value)` | Returns integer `-1`, `0`, or `1`. NaN raises an argument error. |
| `math::trunc(value)` | Truncates toward zero, preserving the integer or float input type. |
| `math::cbrt(value)` | Returns the real cube root as a float, including for negative inputs. |
| `math::hypot(x, y)` | Computes the distance `sqrt(x*x + y*y)` as a float, avoiding unnecessary overflow and underflow. |
| `math::log1p(value)` | Computes `ln(1+value)` with improved numerical accuracy near zero. |
| `math::expm1(value)` | Computes `exp(value)-1` with improved numerical accuracy near zero. |
| `math::e()` | Returns Euler's number, approximately `2.718281828459045`, as a float. |

{{< code lang="fql" >}}
RETURN [
    math::clamp(20, 0, 10),
    math::sign(-12),
    math::trunc(-3.9),
    math::cbrt(-8),
    math::hypot(3, 4)
]
// [10, -1, -3, -2, 5]
{{</ code >}}

Use `log1p` and `expm1` directly when working with very small values: evaluating
`log(1+x)` or `exp(x)-1` can lose precision through rounding or cancellation.
Both functions return floats. `log1p(-1)` returns negative infinity, and values
below `-1` return NaN.

Infinite clamp bounds are allowed. Runtime comparisons preserve exact numeric
ordering, including large integers and mixed integer/float inputs. Clamp returns
`min` when `value` is strictly below it, `max` when strictly above it, and `value`
otherwise. The selected input keeps its original type and exact value. For
example, `math::clamp(5, 0.5, 10.5)` returns integer `5`, while
`math::clamp(-1, 0.5, 10)` returns float `0.5`.

Equality preserves the original `value`, including its signed zero. A NaN input
value is returned unchanged, even with equal infinite bounds. Bounds are still
validated first: NaN or reversed bounds raise an argument error even when the
input value is NaN.

`sign` returns zero for either signed zero and returns `-1` or `1` for negative
or positive infinity. `trunc` preserves integer inputs exactly; for float inputs
it preserves NaN, infinities, and signed zero. `cbrt` also preserves those special
float values. `hypot` returns positive infinity if either coordinate is infinite,
even when the other is NaN. `expm1` preserves NaN, positive infinity, and signed
zero, while negative infinity returns `-1`.

## Numeric inputs

`math::sum`, `math::mean`, `math::min`, `math::max`, `math::median`,
`math::variance`, `math::variance_sample`, `math::stddev`,
`math::stddev_sample`, and `math::percentile` accept lists containing integers,
floats, or a mixture of both. Every element must be numeric. Strings such as
`"2"`, `none`, booleans, and other non-numeric values cause an error identifying
the argument and element index; values are not converted or silently filtered.

For example, `math::sum([1, "2", 3])` raises an error. Validate or explicitly
transform malformed data before calculating statistics.

These functions leave their inputs unchanged, including when an error occurs.
Host-provided lists are supported through the runtime list contract. Read and
cancellation errors returned by hosts propagate. Traversal and sorting pass
through the caller's context without polling. A completed operation keeps its
result even if cancellation occurs concurrently; the VM observes execution
cancellation after control returns to a cancellation boundary.

Non-finite floats remain numeric values. Arithmetic uses floating-point
calculations, so accepted numeric input does not guarantee a finite result.

## Empty and insufficient input

| Canonical function | Empty list |
| --- | --- |
| `math::sum` | Integer `0` |
| `math::min`, `math::max` | `none` |
| `math::mean`, `math::median`, `math::percentile` | `NaN` |
| `math::variance`, `math::variance_sample` | `NaN` |
| `math::stddev`, `math::stddev_sample` | `NaN` |

Check undefined statistics with `IS_NAN` before encoding a result.

Population variance is the default: `math::variance` divides by `N`, and
`math::variance_sample` divides by `N - 1`. Standard deviation is the square root
of the corresponding variance. Welford's one-pass algorithm uses constant
additional storage, so host lists need not support repeated traversal.

For one finite number, population variance and standard deviation return zero.
Sample statistics require at least two observations and return `NaN` otherwise.
Non-finite numbers produce `NaN` for variance and standard deviation.

## Median and percentile

Median selects the central number after ordering the values. For an even count,
it averages only the two middle numbers. Both `math::median([1,2,3,4])` and
`math::median([1,2,3,100])` return `2.5`.

Canonical percentile takes exactly two arguments:

{{< code lang="fql" >}}
RETURN [
    math::percentile([0, 100, 200, 300, 400], 0),
    math::percentile([0, 100, 200, 300, 400], 50),
    math::percentile([0, 100, 200, 300, 400], 99.9),
    math::percentile([0, 100, 200, 300, 400], 100)
]
// [0, 200, 399.6, 400]
{{</ code >}}

The percentile must be a finite integer or float from `0` through `100`,
inclusive. Invalid percentiles raise an error even when the list is empty.
A method-string argument is not accepted.

The calculation uses linear interpolation at ascending zero-based position
`(p / 100) * (N - 1)`. Exact positions select the corresponding number;
fractional positions interpolate between adjacent values. Zero selects the
minimum and 100 selects the maximum. A singleton returns its only value for
every valid percentile.

Odd medians and exact percentile selections retain the selected integer or
float type. Averaged and interpolated results are floats. Both operations sort
a private numeric snapshot without copying, sorting, indexing, or mutating the
source list.

## Migrate global calls

| Deprecated global | Canonical replacement |
| --- | --- |
| `average(values)` | `math::mean(values)` |
| `variance_population(values)` | `math::variance(values)` |
| `stddev_population(values)` | `math::stddev(values)` |
| `variance_sample`, `stddev_sample` | Same name under `math::` |
| `abs`, `ceil`, `floor`, `round`, `min`, `max`, `sum`, `median` | Same name under `math::` |
| Power, exponential, logarithmic, trigonometric, and angle-conversion functions | Same name under `math::` |
| `percentile(values, p[, method])` | `math::percentile(values, p)` with canonical interpolation |
| `pi()` | `math::pi()` |
| `range(start, end[, step])` | `arrays::range(start, end[, step])` |

The deprecated global collection functions continue ignoring non-numeric
elements without coercion. This behavior exists for migration compatibility;
the canonical API rejects malformed input.

| Legacy operation | Empty list | Nonempty list with no numbers |
| --- | --- | --- |
| `sum` | Integer `0` | Float `0` |
| `average` | Float `0` | Float `0` |
| `min`, `max`, `median` | `none` | `none` |
| Variance, standard deviation, percentile | `NaN` | `NaN` |

Legacy percentile still requires an integer in `1..100` and defaults to nearest
rank, selecting one-based position `ceil(p * N / 100)`. Only the exact string
`"interpolation"` enables linear interpolation; other strings retain rank
fallback. Here `N` counts the numeric elements after filtering. If none remain,
legacy percentile returns `NaN` before validating the percentile argument, while
a supplied method must still be a string. These legacy rules do not apply to
`math::percentile`.

## Range, randomness, and constants

Range construction belongs to [Arrays]({{< ref "array-migration" >}}).
`arrays::range` preserves inclusive endpoints, a default step of positive one,
and explicit negative steps for descending ranges. Math-only embeddings retain
the deprecated global `range`; Arrays-only embeddings expose `arrays::range`.

The deprecated global `rand` keeps its coercion, rounded results, and maximum-first
argument order. Canonical value generation belongs to `random::`; see
[Random functions]({{< ref "/docs/language/functions/random-functions" >}}) for
intervals and migration. Both APIs consume the same Session source. Math-only
embeddings retain `rand`; Random-only embeddings expose the canonical functions.
There is no `math::rand` or `math::range`.

Pi remains a zero-argument function, `math::pi()`, because the current namespace
registry supports functions rather than constants. Euler's number follows the
same model as `math::e()`. The existing global `pi()` remains deprecated; there
is no global `e()`.

## Query aggregation

Built-in `COLLECT AGGREGATE` operations keep their language behavior independently
of public math functions. Unqualified `SUM`, `AVERAGE`, `MIN`, and `MAX` retain
numeric filtering, and `COUNT` counts observations. Generic and fused execution
use the same VM-owned reduction semantics. Explicit namespaced selectors such
as `math::sum` remain ordinary function calls with strict numeric validation.
