package cs

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
