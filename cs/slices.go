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

// Finds and returns a slice containing all elements in s for which
// funcToCall returns true, or an empty slice
// if none are found.
func FilterSlice[S ~[]V, V any](s S, funcToCall func(V) bool) (allValues S) {
	values := S{}
	for _, v := range s {
		if funcToCall(v) {
			values = append(values, v)
		}
	}
	return values
}
