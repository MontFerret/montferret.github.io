---
title: "DateTime functions and migration"
sidebarTitle: "DateTime migration"
weight: 47
draft: false
description: "Use the canonical datetime namespace and migrate legacy date functions."
---

# DateTime functions and migration

In runtime releases containing the canonical v2 datetime library, new code uses
`datetime::`. Existing globals remain available as deprecated migration
functions. Earlier releases require the global names; the examples below require
a runtime containing the namespace.

```fql
let created = datetime::parse("2026-09-11T14:30:00Z")
let other = datetime::parse("2026-09-11T16:00:00Z")
return {
    year: datetime::year(created),
    month: datetime::month(created),
    formatted: datetime::format(created, "2006-01-02"),
    same_day: datetime::same(created, other, "day"),
    elapsed_hours: datetime::diff(created, other, "hours")
}
```

`parse(text)` accepts RFC3339; `parse(text, layout)` accepts a Go time layout.
`format(value, layout)` uses the same layout convention, such as `"2006-01-02"`.
Parsing requires a String and does not construct dates from numeric components.

## Canonical names

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

Accessors use each value's own calendar. Weekdays remain Sunday=0 through
Saturday=6, day-of-year starts at 1, and quarters range from 1 to 4.
`days_in_month` handles Gregorian leap years and all month lengths, including
the corrected 31 days in July.

## Compare calendar precision

`datetime::same(a, b, unit)` requires every enclosing calendar field to match.

| Unit | Fields compared |
| --- | --- |
| `year` | Year |
| `month` | Year and month |
| `week` | ISO week-year and ISO week number |
| `day` | Year, month, and day |
| `hour` | Calendar date and hour |
| `minute` | Calendar date, hour, and minute |
| `second` | Calendar date, hour, minute, and second |
| `millisecond` | All fields through second and the truncated millisecond |

Each input uses its own local calendar. Values are not normalized to UTC, and
offset/location identity is not part of equality. Two repeated local times
during a DST overlap can compare equal; the same instant represented on different
local dates can compare unequal. Use native DateTime equality for instant equality.

ISO week-year can differ from calendar year. December 31, 2020 and January 1,
2021 belong to the same ISO week, while matching week numbers in different ISO
week-years do not match.

The deprecated `date_compare` keeps its range signature and default end of
millisecond. Its inclusive component order is year, month, week, day, hour,
minute, second, millisecond; every selected component must match. Week includes
ISO week-year. Reversed ranges fail. This fixes its former behavior where any
matching component could make the result true. It is not a direct synonym for
`datetime::same`: a day-only range compares the day of month, while
`datetime::same(a, b, "day")` compares the whole calendar date.

## Measure elapsed time

`datetime::diff(a, b, unit)` returns **b minus a** as a Float, including fractional
units. Swapping arguments reverses the sign.

```fql
let start = datetime::parse("2026-09-11T14:30:00Z")
let end = datetime::parse("2026-09-11T16:00:00Z")
return [
    datetime::diff(start, end, "hours"),
    datetime::diff(end, start, "hours")
]
// [1.5, -1.5]
```

Supported units are millisecond, second, minute, and hour, including their
plurals. Day, week, month, and year are rejected even for identical inputs.
Calendar differences need separate calendar/timezone semantics and are deferred.
DST offsets affect elapsed time: 01:30 at -05:00 to 03:30 at -04:00 is one hour.

The calculation uses checked native DateTime subtraction. Intervals outside the
runtime Duration range (roughly 292 years) produce a range error.

`date_diff` now shares the signed calculation, supported units, and range
errors. Its default/false result is an Int truncated toward zero; passing true
returns the canonical Float. The old absolute result and fixed-duration
approximations for calendar units are intentionally removed.

## Add and subtract calendar units

`datetime::add(value, amount, unit)` and `datetime::subtract(value, amount, unit)`
retain integer amounts, including negative amounts. Units are millisecond,
second, minute, hour, day, week, month, and year; singular and plural names are
case-insensitive.

Subday units measure elapsed time. Calendar units use the input's location and
Go calendar arithmetic: a day can span 23 or 25 hours across DST. Month ends
normalize rather than clamp. January 31, 2023 plus one month becomes March 3,
2023; February 29, 2024 plus one year becomes March 1, 2025.

## Go callers

The Go package exports `Now`, `Parse`, `Format`, `Year`, `Month`, `Day`,
`Hour`, `Minute`, `Second`, `Millisecond`, `DayOfWeek`, `DayOfYear`,
`Quarter`, `DaysInMonth`, `IsLeapYear`, `Add`, `Subtract`, `Same`, `Diff`,
and `RegisterLib`. Replace the former `Date*` names with these canonical
functions. Exported unit machinery is removed; migration compatibility is
provided only through the deprecated FQL globals.

The DateTime runtime representation is unchanged. This refactor adds no timezone
conversion, timestamp helpers, boundary operations, or numeric construction APIs.
See the [DateTime reference]({{< ref "docs/standard-library/datetime" >}}) for the
functions in the website's selected published runtime.
