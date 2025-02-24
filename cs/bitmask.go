package cs

type Bitmask uint32

// Return number of non-zero bits in a bitmask
func (mask Bitmask) countBits() int {
	count := 0

	for mask > 0 {
		count += int(mask & 1)
		mask >>= 1
	}

	return count
}
