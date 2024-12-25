package cs

/*
Compare 2 or more slices without order and return true if they are either equal or
if slice 1 contains slice 2

Identical determines what criteria to check for - true requires the 2 slices to be strictly identical,
while false merely requires that slice contains other (other is a *subset* of slice).
*/
func CompareSlicesUnordered[T comparable, S ~[]T](slice, other S, identical bool) bool {
	/* If two sets have the same length and one is a subset of the other, it
	*necessarily* follows that the two are equal.
	 */
	if identical && len(slice) != len(other) {
		return false
	}

	// Tally up counters for items in both sets
	numItemsInFirst := map[T]int{}
	numItemsInSecond := map[T]int{}
	for _, item := range slice {
		numItemsInFirst[item]++
	}
	for _, item := range other {
		numItemsInSecond[item]++
	}

	for item, countInSecond := range numItemsInSecond {
		if numItemsInFirst[item] < countInSecond {
			// there exist items in slice 2 not accounted for in slice 1
			return false
		}
	}

	return true
}

// Remove duplicates from one or more slices and return the appended result.
// Items appear in the order of the slices passed in (everything in slice 1, then everything in slice 2, etc.)
//
// To pass in map objects or iterables, call maps.Values and/or slices.Collect on them first
func AppendWithoutDuplicates[T comparable, S ~[]T](slices ...S) S {
	checkedParts := map[T]bool{}
	var newSlice S
	// smush all our slices together into 1 big slice
	allSlices := slices[0]
	for i := 1; i < len(slices); i++ {
		allSlices = append(allSlices, slices[i]...)
	}

	// iterate over big slice and slap items onto new list if not already covered
	for _, item := range allSlices {
		if !checkedParts[item] {
			newSlice = append(newSlice, item)
			checkedParts[item] = true
		}
	}
	return newSlice
}

// If key does not exist in lookupMap, evaluates funcToCall on key and
// stores the resulting value in the map before returning it.
// Otherwise, simply returns the previous value in the map.
//
// Used for quick storage/lookup of repeatedly needed/computed values
func UpdateLookupMap[M ~map[K]V, K comparable, V any](lookupMap M, key K, funcToCall func(K) V) V {
	if val, ok := lookupMap[key]; ok {
		return val
	} else {
		val = funcToCall(key)
		lookupMap[key] = val
		return val
	}
}

/* break down an individual bitmask into a slice of its component bits
func (mask Bitmask) GetBits() []Bitmask {
	bits := []Bitmask{}

	for num := Bitmask(1); num <= mask; num <<= 1 {
		if num&mask != 0 {
			bits = append(bits, num)
		}
	}
	return bits
}*/
