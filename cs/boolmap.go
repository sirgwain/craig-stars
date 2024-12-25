package cs 

// a map[bool]float64, but with some extra fluff to circumvent json encoding jank
type BoolMap struct {
	valueIfTrue float64
	valueIfFalse float64
}

func (m BoolMap) Get(boolean bool) float64 {
	if boolean {
		return m.valueIfTrue
	}
	return m.valueIfFalse
}