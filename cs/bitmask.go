package cs

import "math/bits"

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

// Break down an individual Bitmask or similar type into a slice of its constituent bits
func getBits[B ~uint32](mask B) []B {
	bits := make([]B, bits.OnesCount(uint(mask)))

	for num := B(1); num <= mask; num <<= 1 {
		if num&mask != 0 {
			bits = append(bits, num)
		}
	}
	return bits
}
