package cs

import (
	"math"
)

// population is often updated with floating point math, but we have to convert
// it back to Colonist Cargo values, which are stored in units of 100 colonists per 1kT of Colonist Cargo
func roundToNearest100f(value float64) int {
	return int(math.Round(value/100) * 100)
}

func roundToNearest100(value int) int {
	return int(math.Round(float64(value)/100) * 100)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// round a float to nearest whole number, rounding halves down
func roundHalfDown(x float64) float64 {
	if x > 0 {
		return math.Floor(x + 0.5)
	}
	return math.Ceil(x - 0.5)
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}

func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	} else {
		if value > max {
			return max
		}
	}
	return value
}

func MaxInt(nums ...int) int {
	result := math.MinInt
	for _, value := range nums {
		if value > result {
			result = value
		}
	}

	return result
}

func MinInt(nums ...int) int {
	result := math.MaxInt
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

func MinFloat64(nums ...float64) float64 {
	result := math.MaxFloat64
	for _, value := range nums {
		if value < result {
			result = value
		}
	}

	return result
}

// raise an integer to the power of another integer (apparently this is the fastYY)
//
// Does not support negative values (we *are* dealing with integers here after all)
func PowInt(base, exponent int) int {
	result := 1
	// According to internet, this is the fastest way to do int exponentiation
	for {
		if exponent&1 == 1 {
			result *= base
		}
		exponent >>= 1
		if exponent == 0 {
			break
		}
		base *= base
	}

	return result
}

func AbsInt(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

// return the absolutely greater of 2 integers
// (ie: the one furthest away from 0)
func MaxAbsInt(nums ...int) int {
	result := math.MinInt
	for _, value := range nums {
		if AbsInt(value) > AbsInt(result) {
			result = value
		}
	}

	return result
}
