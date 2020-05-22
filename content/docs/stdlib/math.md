---
title: "math"
weight: 1
draft: false
---


## MIN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/min.go#L13)

Min returns the smallest (arithmetic mean) of the values in array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The smallest of the values in array.
- - - -

## RADIANS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/radians.go#L13)

Radians returns the angle converted from degrees to radians.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Float` `Int` | The input number.


**Returns** `Float` The angle in radians.
- - - -

## FLOOR
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/floor.go#L14)

Floor returns the greatest integer value less than or equal to a given value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Int` The greatest integer value less than or equal to a given value.
- - - -

## PERCENTILE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/percentile.go#L17)

Percentile returns the nth percentile of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.
`numb` | `Int` | A number which must be between 0 (excluded) and 100 (included).
`method` | `String, optional` | "rank" (default) or "interpolation".


**Returns** `Float` The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
- - - -

## LOG10
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log10.go#L14)

Log10 returns the decimal logarithm of a given value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The decimal logarithm of a given value.
- - - -

## SIN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sin.go#L14)

Sin returns the sine of the radian argument.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The sin, in radians, of a given number.
- - - -

## LOG2
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log2.go#L14)

Log2 returns the binary logarithm of a given value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The binary logarithm of a given value.
- - - -

## VARIANCE_POPULATION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_population.go#L14)

PopulationVariance returns the population variance of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The population variance.
- - - -

## TAN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/tan.go#L14)

Tan returns the tangent of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Int` `Float` | A number.


**Returns** `Float` The tangent.
- - - -

## AVERAGE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/average.go#L13)

Average Returns the average (arithmetic mean) of the values in array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The average of the values in array.
- - - -

## DEGREES
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/degrees.go#L13)

Degrees returns the angle converted from radians to degrees.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Float` `Int` | The input number.


**Returns** `Float` The angle in degrees.l
- - - -

## RANGE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/range.go#L14)

Range returns an array of numbers in the specified range, optionally with increments other than 1.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`start` | `Int` `Float` | The value to start the range at (inclusive).
`end` | `Int` `Float` | The value to end the range with (inclusive).
`step` | `Int` `Float, optional` | How much to increment in every step, the default is 1.0.


**Returns** `None`
- - - -

## ACOS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/acos.go#L14)

Acos returns the arccosine, in radians, of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The arccosine, in radians, of a given number.
- - - -

## RAND
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/rand.go#L15)

Rand return a pseudo-random number between 0 and 1.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`max` | `Float` `Int, optional` | Upper limit.
`min` | `Float` `Int, optional` | Lower limit.


**Returns** `Float` A number greater than 0 and less than 1.
- - - -

## STDDEV_SAMPLE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_sample.go#L14)

StandardDeviationSample returns the sample standard deviation of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The sample standard deviation.
- - - -

## SUM
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sum.go#L13)

Sum returns the sum of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The sum of the values.
- - - -

## MEDIAN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/median.go#L14)

Median returns the median of the values in array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The median of the values in array.
- - - -

## ABS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/abs.go#L14)

Abs returns the absolute value of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The absolute value of a given number.
- - - -

## ROUND
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/round.go#L14)

Round returns the nearest integer, rounding half away from zero.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Int` The nearest integer, rounding half away from zero.
- - - -

## STDDEV_POPULATION
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/stddev_population.go#L14)

StandardDeviationPopulation returns the population standard deviation of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The population standard deviation.
- - - -

## ATAN2
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan2.go#L15)

Atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number1` | `Int` `Float` | Input number.
`number2` | `Int` `Float` | Input number.


**Returns** `Float` The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
- - - -

## COS
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/cos.go#L14)

Cos returns the cosine of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The cosine of a given number.
- - - -

## PI
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pi.go#L12)

Pi returns Pi value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description


**Returns** `Float` Pi value.
- - - -

## VARIANCE_SAMPLE
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/variance_sample.go#L14)

SampleVariance returns the sample variance of the values in a given array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The sample variance.
- - - -

## ATAN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/atan.go#L14)

Atan returns the arctangent, in radians, of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The arctangent, in radians, of a given number.
- - - -

## ASIN
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/asin.go#L14)

Asin returns the arcsine, in radians, of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The arcsine, in radians, of a given number.
- - - -

## POW
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/pow.go#L15)

Pow returns the base to the exponent value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`base` | `Int` `Float` | The base value.
`exp` | `Int` `Float` | The exponent value.


**Returns** `Float` The exponentiated value.
- - - -

## EXP2
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp2.go#L14)

Exp2 returns 2 raised to the power of value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` 2 raised to the power of value.
- - - -

## CEIL
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/ceil.go#L14)

Ceil returns the least integer value greater than or equal to a given value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Int` The least integer value greater than or equal to a given value.
- - - -

## LOG
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/log.go#L14)

Log returns the natural logarithm of a given value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` The natural logarithm of a given value.
- - - -

## EXP
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/exp.go#L14)

Exp returns Euler's constant (2.71828...) raised to the power of value.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`number` | `Int` `Float` | Input number.


**Returns** `Float` Euler's constant raised to the power of value.
- - - -

## MAX
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/max.go#L13)

Max returns the greatest (arithmetic mean) of the values in array.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`array` | `Array` | Array of numbers.


**Returns** `Float` The greatest of the values in array.
- - - -

## SQRT
[Source](https://github.com/MontFerret/ferret/tree/master/pkg/stdlib/math/sqrt.go#L14)

Sqrt returns the square root of a given number.

|          |          |          |
---------- | -------- | ----------
Argument   | Type     | Description
`value` | `Int` `Float` | A number.


**Returns** `Float` The square root.
- - - -
