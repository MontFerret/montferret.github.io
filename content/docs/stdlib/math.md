

---
title: "math"
weight: 1
draft: false
menu: [ABS,ACOS,ASIN,ATAN,ATAN2,AVERAGE,CEIL,COS,DEGREES,EXP,EXP2,FLOOR,LOG,LOG10,LOG2,MAX,MEDIAN,MIN,PERCENTILE,PI,POW,RADIANS,RAND,RANGE,ROUND,SIN,SQRT,STDDEV_POPULATION,STDDEV_SAMPLE,SUM,TAN,VARIANCE_POPULATION,VARIANCE_SAMPLE,]
---



{{< header >}}
ABS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/abs.go#L14)

ABS returns the absolute value of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The absolute value of a given number.
- - - -


{{< header >}}
ACOS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/acos.go#L14)

ACOS returns the arccosine, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arccosine, in radians, of a given number.
- - - -


{{< header >}}
ASIN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/asin.go#L14)

ASIN returns the arcsine, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arcsine, in radians, of a given number.
- - - -


{{< header >}}
ATAN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan.go#L14)

ATAN returns the arctangent, in radians, of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arctangent, in radians, of a given number.
- - - -


{{< header >}}
ATAN2
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan2.go#L15)

ATAN2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number1` | `Int` `Float`  |  | Input number.
`number2` | `Int` `Float`  |  | Input number.


**Returns** `Float` The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
- - - -


{{< header >}}
AVERAGE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/average.go#L13)

AVERAGE Returns the average (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The average of the values in array.
- - - -


{{< header >}}
CEIL
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/ceil.go#L14)

CEIL returns the least integer value greater than or equal to a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The least integer value greater than or equal to a given value.
- - - -


{{< header >}}
COS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/cos.go#L14)

COS returns the cosine of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The cosine of a given number.
- - - -


{{< header >}}
DEGREES
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/degrees.go#L13)

DEGREES returns the angle converted from radians to degrees.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | The input number.


**Returns** `Float` The angle in degrees
- - - -


{{< header >}}
EXP
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp.go#L14)

EXP returns Euler's constant (2.71828...) raised to the power of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` Euler's constant raised to the power of value.
- - - -


{{< header >}}
EXP2
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp2.go#L14)

EXP2 returns 2 raised to the power of value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` 2 raised to the power of value.
- - - -


{{< header >}}
FLOOR
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/floor.go#L14)

FLOOR returns the greatest integer value less than or equal to a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The greatest integer value less than or equal to a given value.
- - - -


{{< header >}}
LOG
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log.go#L14)

LOG returns the natural logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The natural logarithm of a given value.
- - - -


{{< header >}}
LOG10
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log10.go#L14)

LOG10 returns the decimal logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The decimal logarithm of a given value.
- - - -


{{< header >}}
LOG2
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log2.go#L14)

LOG2 returns the binary logarithm of a given value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The binary logarithm of a given value.
- - - -


{{< header >}}
MAX
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/max.go#L13)

MAX returns the greatest (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The greatest of the values in array.
- - - -


{{< header >}}
MEDIAN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/median.go#L14)

MEDIAN returns the median of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The median of the values in array.
- - - -


{{< header >}}
MIN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/min.go#L13)

MIN returns the smallest (arithmetic mean) of the values in array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The smallest of the values in array.
- - - -


{{< header >}}
PERCENTILE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/percentile.go#L17)

PERCENTILE returns the nth percentile of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`array` | `Int[]` `Float[]`  |  | Array of numbers.
`number` | `Int`  |  | A number which must be between 0 (excluded) and 100 (included).
`method` | `String`  | `"rank"` | "rank" or "interpolation".


**Returns** `Float` The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
- - - -


{{< header >}}
PI
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pi.go#L12)

PI returns Pi value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description


**Returns** `Float` Pi value.
- - - -


{{< header >}}
POW
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pow.go#L15)

POW returns the base to the exponent value.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`base` | `Int` `Float`  |  | The base value.
`exp` | `Int` `Float`  |  | The exponent value.


**Returns** `Float` The exponentiated value.
- - - -


{{< header >}}
RADIANS
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/radians.go#L13)

RADIANS returns the angle converted from degrees to radians.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | The input number.


**Returns** `Float` The angle in radians.
- - - -


{{< header >}}
RAND
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/rand.go#L15)

RAND return a pseudo-random number between 0 and 1.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`max` | `Int` `Float`  |  | Upper limit.
`min` | `Int` `Float`  |  | Lower limit.


**Returns** `Float` A number greater than 0 and less than 1.
- - - -


{{< header >}}
RANGE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/range.go#L15)

RANGE returns an array of numbers in the specified range, optionally with increments other than 1.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`start` | `Int` `Float`  |  | The value to start the range at (inclusive).
`end` | `Int` `Float`  |  | The value to end the range with (inclusive).
`step` | `Int` `Float`  | `1.0` | How much to increment in every step.


**Returns** `Int[]` `Float[]` Array of numbers in the specified range, optionally with increments other than 1.
- - - -


{{< header >}}
ROUND
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/round.go#L14)

ROUND returns the nearest integer, rounding half away from zero.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Int` The nearest integer, rounding half away from zero.
- - - -


{{< header >}}
SIN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sin.go#L14)

SIN returns the sine of the radian argument.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | Input number.


**Returns** `Float` The sin, in radians, of a given number.
- - - -


{{< header >}}
SQRT
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sqrt.go#L14)

SQRT returns the square root of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`value` | `Int` `Float`  |  | A number.


**Returns** `Float` The square root.
- - - -


{{< header >}}
STDDEV_POPULATION
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_population.go#L14)

STDDEV_POPULATION returns the population standard deviation of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The population standard deviation.
- - - -


{{< header >}}
STDDEV_SAMPLE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_sample.go#L14)

STDDEV_SAMPLE returns the sample standard deviation of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sample standard deviation.
- - - -


{{< header >}}
SUM
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sum.go#L13)

SUM returns the sum of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sum of the values.
- - - -


{{< header >}}
TAN
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/tan.go#L14)

TAN returns the tangent of a given number.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`number` | `Int` `Float`  |  | A number.


**Returns** `Float` The tangent.
- - - -


{{< header >}}
VARIANCE_POPULATION
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_population.go#L14)

VARIANCE_POPULATION returns the population variance of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The population variance.
- - - -


{{< header >}}
VARIANCE_SAMPLE
{{</ header >}}
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_sample.go#L14)

VARIANCE_SAMPLE returns the sample variance of the values in a given array.

|          |          |                |
---------- | -------- | -------------- | ----------
Argument   | Type     | Default value  | Description
`numbers` | `Int[]` `Float[]`  |  | Array of numbers.


**Returns** `Float` The sample variance.
- - - -
