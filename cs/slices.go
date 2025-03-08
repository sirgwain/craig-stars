package cs

// Finds and returns the first element in s for which
// funcToCall returns true, or the type's zero value
// if none are found.
func FindSlice[S ~[]V, V any](s S, funcToCall func(V) bool) (firstValue V) {
	for _, v := range s {
		if funcToCall(v) {
			return v
		}
	}
	return firstValue
}

// MapSlice applies mapFunc on every element in s, returning the modified slice.
// Any elements for which keep evaluates to false are instead removed entirely.
// MapSlice zeroes the elements between the new length and the original length.
//
//	MapSlice([]int{1, 2, 3, 4, 6}, func(i int) (int, bool) {
//	   return i + 1, i % 3 != 0
//	}) = []int{2, 3, 5}
func MapSlice[S ~[]V, V any](s S, mapFunc func(V) (new V, keep bool)) S {
	values := s[:0] // minimizes reallocation
	for _, v := range s {
		if new, keep := mapFunc(v); keep {
			values = append(values, new)
		}
	}
	s = values
	clear(s[len(values):]) // zero out for GC
	return s[:len(values)]
}
