package cs

// the rng rules all
type rng interface {
	Float64() float64
	Intn(n int) int
	Shuffle(n int, swap func(i, j int))
}
 