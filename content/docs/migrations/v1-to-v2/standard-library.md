---
title: "Standard library migration"
sidebarTitle: "Standard library"
weight: 20
draft: false
description: "Migrate legacy function calls and review changed standard-library behavior."
aliases:
  - /docs/language/functions/string-migration/
  - /docs/language/functions/array-migration/
  - /docs/language/functions/object-migration/
  - /docs/language/functions/datetime-migration/
---

# Standard library migration

Review legacy calls by area after the [initial compatibility check]({{< ref "docs/migrations/v1-to-v2" >}}).
This guide includes the namespace and behavior changes introduced during the v2
alpha series. Use a runtime containing the relevant APIs and a CLI release
containing the documented migration rules. Older releases can support fewer rules.
The [generated reference]({{< ref "docs/standard-library" >}}) follows the website's
published runtime pin, which may precede these changes.

Preview replacements with `ferret migrate run --print path/to/project`.
Automatic replacements target immutable operations, never `arrays::mut::` or
`object::mut::`. Calls shadowed by declarations or aliases need manual review;
already-qualified targets are preserved. Other safe calls in the same file can
still migrate. See the [CLI reference]({{< ref "docs/tools/cli/migrate" >}}) for
supported rules, discovery boundaries, diagnostics, and file-level failures.

## Strings, encoding, crypto, and paths

| Former call | Current call |
| --- | --- |
| `concat_separator(separator, values)` | `join(values, separator)` |
| `substitute(text, search, replacement[, limit])` | `replace(text, search, replacement[, limit])` |
| `substitute(text, search)` | `replace(text, search, "")` |
| `contains(text, search, true)` | `find_first(text, search)` for an index |
| `regex_match(text, pattern)` | `regex_find(text, pattern)` or `regex_find_all(text, pattern)` |
| Regex Boolean case options | Inline pattern flags such as `(?i)` |
| `json_parse`, `json_stringify` | `encoding::json_parse`, `encoding::json_stringify` |
| `encode_uri_component`, `decode_uri_component` | `encoding::query_escape`, `encoding::query_unescape` |
| `to_base64`, `from_base64` | `encoding::base64_encode`, `encoding::base64_decode` |
| `escape_html`, `unescape_html` | `encoding::html_escape`, `encoding::html_unescape` |
| `md5`, `sha1`, `sha512`, `random_token` | The same names under `crypto::` |
| `base`, `clean`, `dir`, `ext`, `is_abs`, `separate`, `match` | The same names under `path::` |
| Global path `join(parts...)` | `path::join(parts...)` |

The CLI automatically rewrites the listed encoding, crypto, and path calls
except `join`. Global `join` can mean old path joining or current string joining,
so it is left for manual review. Apply the remaining string renames and option
changes manually; they are not covered by automatic rules.

String operations remain global. Removed names have no compatibility aliases.
`starts_with`, `ends_with`, `repeat`, and `crypto::sha256` are new. Identifiers
can contain underscores after digits, including `base64_encode`. All eight
path functions move to `path::` without global aliases; their behavior and
argument counts are preserved.

Review these behavior changes as well as names:

- Text, patterns, separators, cutsets, and `join` elements require Strings.
  Convert values explicitly with `to_string` where intended. Bounds, lengths,
  and limits require Int values.
- Text offsets count Unicode runes. Split limits count results, replacement
  limits count replacements, and explicit `-1` limits are invalid. Omit a limit
  for unlimited operation.
- `contains` returns Boolean; use `find_first` for an index. `like` uses glob
  syntax, with literal percent and underscore characters.
- Regex options use inline flags. `regex_find` returns `none` or a match object;
  `regex_find_all` returns an array. Review code expecting the old result shape.
  Boolean regex options and the old reserved `regex_split` argument are removed.
- `fmt` requires either automatic or indexed placeholders, with every supplied
  value referenced. Review mixed modes and unused arguments.
- Base64 decoding returns Binary, including non-text bytes. Query escaping uses
  query-form rules: spaces become `+` and a literal plus becomes `%2B`.

## Arrays

Existing global array functions remain temporary compatibility functions with
their legacy signatures and behavior. Canonical operations use `arrays::`.

| Global call | Canonical operation |
| --- | --- |
| `first`, `last`, `flatten`, `slice`, `unique`, `sorted` | Same name under `arrays::` |
| `range(start, end[, step])` | `arrays::range(start, end[, step])` |
| `nth(xs, index)` | `arrays::at(xs, index)` |
| `append(xs, value)` or `push(xs, value)` | `arrays::append(xs, value)` |
| `union(a, b, ...)` | `arrays::concat(a, b, ...)` |
| `union_distinct(a, b, ...)` | `arrays::union(a, b, ...)` |
| `intersection(a, b, ...)` | `arrays::intersection(a, b, ...)` |
| `minus(a, b, ...)` | `arrays::difference(a, b, ...)` |
| `position(xs, value)` or `position(xs, value, false)` | `arrays::contains(xs, value)` |
| `position(xs, value, true)` | `arrays::index_of(xs, value)` |
| `remove_value(xs, value)` | `arrays::remove(xs, value)` |
| `remove_values(xs, values)` | `arrays::remove_any(xs, values)` |
| `remove_nth(xs, index)` | Review host-list removal semantics before using `arrays::remove_at(xs, index)` |
| `sorted_unique(xs)` | `arrays::sorted(arrays::unique(xs))` |
| `shift(xs)` | `arrays::slice(xs, 1)` |
| `unshift(xs, value)` | `arrays::concat([value], xs)` |

Most direct replacements in the table can be automated. `remove_nth` and
`unshift` require manual review: host-list removal behavior and argument
evaluation order can differ. Literal modes in `position`, `append`, and `push`
are handled only when the replacement preserves their meaning. Dynamic modes
remain unchanged. Review the [CLI support details]({{< ref "docs/tools/cli/migrate" >}}#calls-requiring-manual-review).

### Replace mode flags deliberately

Canonical `append` and `remove` take exactly two arguments. Canonical search
functions do not accept a Boolean return-type switch.

To deduplicate an appended result, use:

{{< code lang="fql" >}}
return arrays::unique(arrays::append([1, 1], 2))
{{</ code >}}

This returns `[1,2]`. It intentionally removes existing duplicates too. Legacy
`append(xs, value, true)` and `push(xs, value, true)` only skip adding a value
already present; they preserve existing duplicates. Keep those global calls
when that exact legacy behavior is required. A false flag is equivalent to the
canonical two-argument append.

`arrays::remove(xs, value)` removes all matches. `arrays::remove_any(xs, values)`
removes all matches to any listed value, preserving other duplicates and order.
The optional removal limit remains only on global `remove_value`: negative
means unlimited, zero removes nothing, and positive limits removals.

Legacy `unshift(xs, value, true)` corresponds to
`arrays::concat([value], arrays::remove(xs, value))`. It preserves duplicates of
other values. Global `pop` remains available and safely returns `[]` for an empty
array. To express it through canonical functions:

{{< code lang="fql" >}}
let xs = [1, 2, 3]
return length(xs) == 0 ? [] : arrays::slice(xs, 0, length(xs) - 1)
{{</ code >}}

The CLI drops a literal `false` append/push mode but leaves unique mode for
review. For `remove_value`, it can drop a negative integer literal limit;
zero, positive, and dynamic limits remain manual. Parenthesized literals are
recognized, but arbitrary constant expressions are not evaluated.

For `pop`, bind the array expression once before using the example above so
migration does not repeat side effects. For `unshift`, bind the source array
before the value if evaluating either expression can have side effects.
The CLI leaves both calls unchanged for these reasons.

### Review symmetric difference

`unique` and `union` retain first occurrences. Intersection and difference follow
the first input's order. Equal numeric representations and nested values use
Ferret equality, with hash collisions checked for equality.

`arrays::symmetric_difference` keeps values present in an odd number of input
arrays. Duplicates within one array count once. Output follows first encounter
across inputs and retains the first representation of an equal value.

{{< code lang="fql" >}}
return {
    xor: arrays::symmetric_difference([7, 3, 7], [7, 2], [7, 2, 5]),
    legacy: outersection([7, 3, 7], [7, 2], [7, 2, 5])
}
{{</ code >}}

The XOR result is `[7,3,5]`. Legacy `outersection` returns `[3,5]`: it keeps values
present in exactly one input array, excluding `7` even though it appears in three
inputs. The two operations agree for two inputs, but migrating three or more
inputs requires choosing the intended rule. There is no canonical
`arrays::outersection` or `arrays::exclusive`.

The CLI rewrites `outersection(a, b)` automatically. Three or more inputs
require manual review.

### Range and mutation

`range(start, end[, step])` migrates to `arrays::range` automatically. Inclusive
endpoints, a default positive-one step, explicit negative descending steps, and
floating-point elements are preserved. The deprecated global remains in Math-only
embeddings; Arrays-only embeddings expose `arrays::range`. There is no `math::range`.

Global `push`, `pop`, `shift`, and `unshift` remain immutable. Do not replace
them with `arrays::mut::` merely because the names match: mutable calls change
the target and can return an extracted value instead of an array.

### Go array callers

The exported functions in `pkg/stdlib/arrays` adopt the canonical v2 names and
signatures. In particular, `Union` becomes distinct union, `Concat` concatenates,
and `Append` and `Remove` accept exactly two values after the context. `Slice`
rejects extra arguments beyond its optional length. Obsolete Go names and
mode-bearing wrappers are removed. Migration compatibility applies to global
FQL registrations, not the old Go functions.

## Objects

The seven global aliases below are deprecated in generated Core API metadata,
with a message naming the canonical replacement. This metadata does not produce
compiler or runtime deprecation warnings. The compatibility layer is temporary;
no removal release is specified.

There are no global aliases for the new `object::entries`,
`object::from_entries`, or `object::omit_keys` APIs, or for mutable operations.
Mutation remains explicit under `object::mut::`.

| Legacy call | Replacement |
| --- | --- |
| `KEYS(value)` | `object::keys(value)` |
| `VALUES(value)` | `object::values(value)` |
| `HAS(value, key)` | `object::has_key(value, key)` |
| `KEEP_KEYS(value, keys...)` | `object::keep_keys(value, keys...)` |
| `MERGE(values...)` | `object::merge(values...)` |
| `MERGE_RECURSIVE(values...)` | `object::merge_deep(values...)` |
| `ZIP(keys, values)` | `object::zip(keys, values)` |

The CLI also converts `KEYS(value, false)` to `object::keys(value)` and
`KEYS(value, true)` to `arrays::sorted(object::keys(value))`. A dynamic mode or
unsupported arity remains manual. The deprecated v2 global `keys` accepts only
one argument, so review old sorted-key calls before executing them on v2.

In v1, `ZIP` kept the first value for duplicate keys. Both `object::zip` and its
deprecated v2 global alias keep the last value. Review affected data even when
the CLI can rename the call automatically.

See the [Objects reference]({{< ref "docs/standard-library/objects" >}})
for current constructors, copy behavior, and explicit mutation.

## Date and Time

Existing globals remain deprecated migration functions in runtimes containing
`datetime::`. Earlier releases require the global names.

| Legacy global | Canonical function |
| --- | --- |
| `now()` | `datetime::now()` |
| `date(text[, layout])` | `datetime::parse(text[, layout])` |
| `date_format(value, layout)` | `datetime::format(value, layout)` |
| `date_year(value)` | `datetime::year(value)` |
| `date_month(value)` | `datetime::month(value)` |
| `date_day(value)` | `datetime::day(value)` |
| `date_hour(value)` | `datetime::hour(value)` |
| `date_minute(value)` | `datetime::minute(value)` |
| `date_second(value)` | `datetime::second(value)` |
| `date_millisecond(value)` | `datetime::millisecond(value)` |
| `date_dayofweek(value)` | `datetime::day_of_week(value)` |
| `date_dayofyear(value)` | `datetime::day_of_year(value)` |
| `date_quarter(value)` | `datetime::quarter(value)` |
| `date_days_in_month(value)` | `datetime::days_in_month(value)` |
| `date_leapyear(value)` | `datetime::is_leap_year(value)` |
| `date_add(value, amount, unit)` | `datetime::add(value, amount, unit)` |
| `date_subtract(value, amount, unit)` | `datetime::subtract(value, amount, unit)` |
| `date_diff(a, b, unit[, asFloat])` | `datetime::diff(a, b, unit)`; see changed semantics below |
| `date_compare(a, b, start[, end])` | Use `datetime::same(a, b, unit)` when precision equality is intended |

The clock global is `now`; there is no `date_now`. Deprecation appears in API
metadata and documentation without compilation or execution warnings.

The CLI automates the direct renames above. It also rewrites
`date_diff(a, b, unit, true)` to `datetime::diff(a, b, unit)`. An omitted, false,
or dynamic floating flag requires manual review because the result type can
differ. `date_compare` always requires manual review.

### Comparison and difference changes

The deprecated `date_compare` keeps its range signature and default end of
millisecond. Its inclusive component order is year, month, week, day, hour,
minute, second, millisecond; every selected component must match. Week includes
ISO week-year. Reversed ranges fail. This fixes its former behavior where any
matching component could make the result true. It is not a direct synonym for
`datetime::same`: a day-only range compares the day of month, while
`datetime::same(a, b, "day")` compares the whole calendar date.

`date_diff` now shares the signed calculation, supported units, and range
errors. Its default/false result is an Int truncated toward zero; passing true
returns the canonical Float. The old absolute result and fixed-duration
approximations for calendar units are intentionally removed.

The canonical calculation is signed **b minus a** and always returns a Float.
Only millisecond, second, minute, and hour units (including plurals) are supported;
day, week, month, and year are rejected even for identical inputs. Intervals
outside the runtime Duration range, roughly 292 years, fail. `days_in_month`
also corrects July to 31 days.

### Go datetime callers

The Go package exports `Now`, `Parse`, `Format`, `Year`, `Month`, `Day`,
`Hour`, `Minute`, `Second`, `Millisecond`, `DayOfWeek`, `DayOfYear`,
`Quarter`, `DaysInMonth`, `IsLeapYear`, `Add`, `Subtract`, `Same`, `Diff`,
and `RegisterLib`. Replace the former `Date*` names with these canonical
functions. Exported unit machinery is removed; migration compatibility is
provided only through the deprecated FQL globals.

The DateTime runtime representation is unchanged. This refactor adds no timezone
conversion, timestamp helpers, boundary operations, or numeric construction APIs.
See the [Date & Time reference]({{< ref "docs/standard-library/datetime" >}})
for current calendar and elapsed-time operations.

## Math

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

The CLI automates scalar math renames and `pi`. Numeric collection operations
require manual review because canonical functions reject non-numbers instead
of filtering them. Review empty inputs and percentile interpolation as well as
function names. The [Math reference]({{< ref "docs/standard-library/math" >}})
documents current numeric inputs and results. Built-in `COLLECT AGGREGATE`
behavior is separate from this public-function migration.

## Randomness

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

The CLI rewrites only zero-argument `rand()` to `random::float()`.
Parameterized calls require manual review. Math-only embeddings retain `rand`;
Random-only embeddings expose the canonical functions.

See the [Random reference]({{< ref "docs/standard-library/random" >}})
for current intervals, sampling, and shuffling. Go embedders can configure
[reproducible Session seeding]({{< ref "docs/embedding/go/configuration" >}}#reproducible-randomness).

## Next steps

{{< docs-related tiles="migrations-v1-to-v2-go-embedding,migrations-v1-to-v2,stdlib,tools-cli-migrate" >}}
