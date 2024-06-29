package cs

type testRandom struct {
	floatCount int
	intCount   int

	floatsToReturn []float64
	intsToReturn   []int
}

func newIntRandom(intsToReturn ...int) *testRandom {
	return &testRandom{
		intsToReturn: intsToReturn,
	}
}

func newFloat64Random(floatsToReturn ...float64) *testRandom {
	return &testRandom{
		floatsToReturn: floatsToReturn,
	}
}

func (t *testRandom) Float64() float64 {
	var result float64
	if len(t.floatsToReturn) > t.floatCount {
		result = t.floatsToReturn[t.floatCount]
		t.floatCount++
	}
	return result
}

func (t *testRandom) Intn(n int) int {
	var result int
	if len(t.intsToReturn) > t.intCount {
		result = t.intsToReturn[t.intCount]
		t.intCount++
	}
	return result
}

func (t *testRandom) Shuffle(n int, swap func(i, j int)) {
	// do nothing
}