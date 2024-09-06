package cs

import (
	"math"
)

// the rng rules all
type rng interface {
	// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0)
	// from the default Source.
	Float64() float64

	// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
	// It panics if n <= 0.
	Intn(n int) int

	// Shuffle pseudo-randomizes the order of elements using the default Source.
	// n is the number of elements. Shuffle panics if n < 0.
	// swap swaps the elements with indexes i and j.
	Shuffle(n int, swap func(i, j int))
}

// GaussianPdf is a gaussian probability distribution function, please refer to the PDF on the sidebar
// https://en.wikipedia.org/wiki/Normal_distribution
// mean is where the distribution is set,
// variance is related to standard deviation - it is set by users based on values
func GaussianPdf(mean float64, variance float64, x float64) float64 {
	var eNum = -math.Pow(x-mean, 2)
	var eDen = 2 * math.Pow(variance, 2)

	return math.Exp(eNum/eDen) / math.Sqrt(2*math.Pi*math.Pow(variance, 2))
}

// NormalizedGaussianPdf returns a normalized result, where it is 1.0 on the mean
func NormalizedGaussianPdf(mean float64, variance float64, x int) float64 {
	// samples the mean as it is the highest value
	maxY := GaussianPdf(mean, variance, mean)

	// normalizes the result
	return GaussianPdf(mean, variance, float64(x)) / maxY
}

// NormalSample uses a rejection sampling algorithm to generate a random distribution equal to the normal distribution
// It works by generating and x and y, then looking up a value of y' from the distribution, and if y' < y then it
// accepts the value of x, otherwise it regenerates it until it finds one
func NormalSample(random rng, mean float64, variance float64, max int) int {
	x := random.Intn(max)
	y := random.Float64()
	y0 := NormalizedGaussianPdf(mean, variance, x)

	for y0 < y {
		x = random.Intn(max)
		y = random.Float64()
		y0 = NormalizedGaussianPdf(mean, variance, x)
	}

	return x
}
