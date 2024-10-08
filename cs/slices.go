package cs

// Compare 2 or more slices without order and return true if they are either equal or
// if slice 1 contains slice 2 
//
// identical determines what to check for - true requires the 2 slices to be strictly identical,
// while false merely requires that the first slice contains the other
func CompareSlicesUnordered[T comparable, S ~[]T](slice, other S, identical bool) bool {
	// set A is defined to be a subset of set B if there exists no element(s)
	// present in set A that are not also present in set B (B contains everything inside A)
	// If the two are the same length, it *necessarily* follows that the two are equal;
	// B must have everything in A (subset) and nothing more (as otherwise it'd be larger)
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
		if numItemsInFirst[item] < countInSecond { // there exist items in slice 2 not accounted for in slice 1
			return false
		}
	}

	return true
}

// remove duplicates from one or more slices and return the appended result;
// items appear in the order of the slices passed in (everything in slice 1, then everything in slice 2, etc.) 
//
// To pass in map objects or iterables, call maps.Values and/or slices.Collect on them first
func AppendWithoutDuplicates[T comparable, S ~[]T](slice ...S) S {
	checkedParts := map[T]bool{}
	var newSlice S
	// smush all our slices together into 1 big slice
	allSlices := slice[0]
	for i := 1; i < len(slice); i++ {
		allSlices = append(allSlices, slice[i]...)
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

// break down an individual bitmask into a slice of its component bits
func (mask Bitmask) getBits() []Bitmask {
	bits := []Bitmask{}

	for num := Bitmask(1); num <= mask; mask <<= 1 {
		if num&mask != 0 {
			bits = append(bits, num)
		}
	}
	return bits
}
