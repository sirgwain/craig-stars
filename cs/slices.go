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

// FilterSlice removes any elements from s for which keep does NOT return true,
// returning the modified slice.
// FilterSlice zeroes the elements between the new length and the original length.
func FilterSlice[S ~[]V, V any](s S, keep func(V) bool) S {
	values := s[:0] // minimzes
	for _, v := range s {
		if keep(v) {
			values = append(values, v)
		}
	}
	clear(s[len(values):]) // zero out for GC
	return s[:len(values)]
}
