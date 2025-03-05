package cs

import (
	"fmt"
	"slices"
	"strings"
)

// Minerals are mined from planets and turned into Cargo
type Mineral struct {
	Ironium   int `json:"ironium,omitempty"`
	Boranium  int `json:"boranium,omitempty"`
	Germanium int `json:"germanium,omitempty"`
}

type MineralType = ResourceType

var MineralTypes = [3]MineralType{
	Ironium,
	Boranium,
	Germanium,
}

func NewMineral(ironium, boranium, germanium int) Mineral {
	return Mineral{
		Ironium:   ironium,
		Boranium:  boranium,
		Germanium: germanium,
	}
}

func (m Mineral) String() string {
	return fmt.Sprintf("Ironium: %d, Boranium: %d, Germanium: %d", m.Ironium, m.Boranium, m.Germanium)
}

func (m Mineral) PrettyString() string {
	texts := make([]string, 0, 4)
	if m.Ironium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT ironium", m.Ironium))
	}
	if m.Boranium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT boranium", m.Boranium))
	}
	if m.Germanium > 0 {
		texts = append(texts, fmt.Sprintf("%dkT germanium", m.Germanium))
	}
	return strings.Join(texts, ", ")
}

// Set sets the value corresponding to minType to amt.
// Unlike all other Mineral functions, this _will_ mutate the original struct's values,
// and is best used for more complex cases not handled by simple addition.
func (m *Mineral) Set(minType MineralType, amt int) {
	switch minType {
	case Ironium:
		m.Ironium = amt
	case Boranium:
		m.Boranium = amt
	case Germanium:
		m.Germanium = amt
	default:
		panic(fmt.Sprintf("mineral.Set called with invalid MineralType %s", minType))
	}
}

// Max returns a Mineral struct containing the higher of
// m's and other's values for each MineralType.
//
//	Mineral{1, 2, 3}.Max(Mineral{4, 0, 5}) = Mineral{4, 2, 5}
func (m Mineral) Max(other Mineral) Mineral {
	return Mineral{
		Ironium:   Max(m.Ironium, other.Ironium),
		Boranium:  Max(m.Boranium, other.Boranium),
		Germanium: Max(m.Germanium, other.Germanium),
	}
}

// MaxNum return the higher of max and this Mineral struct's values
// for each MineralType.
//
//	Mineral{1, 2, 3}.MaxNum(2) = Mineral{2, 2, 3}
func (m Mineral) MaxNum(max int) Mineral {
	return Mineral{
		Ironium:   Max(m.Ironium, max),
		Boranium:  Max(m.Boranium, max),
		Germanium: Max(m.Germanium, max),
	}
}

func (m Mineral) GetAmount(minType MineralType) int {
	switch minType {
	case Ironium:
		return m.Ironium
	case Boranium:
		return m.Boranium
	case Germanium:
		return m.Germanium
	default:
		panic(fmt.Sprintf("Mineral.GetAmount called with invalid MineralType %q", minType))
	}
}

func (m Mineral) Total() int {
	return m.Ironium + m.Boranium + m.Germanium
}

func (m Mineral) ToSlice() [3]int {
	return [3]int{
		m.Ironium,
		m.Boranium,
		m.Germanium,
	}
}

// Return the Cargo equivalent of a Mineral struct.
func (m Mineral) ToCargo() Cargo {
	return Cargo{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// Return the Cost equivalent of a Mineral struct.
func (m Mineral) ToCost() Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// add two minerals and return the result.
func (m Mineral) Add(other Mineral) Mineral {
	return Mineral{
		Ironium:   m.Ironium + other.Ironium,
		Boranium:  m.Boranium + other.Boranium,
		Germanium: m.Germanium + other.Germanium,
	}
}

// Add a number to all components of a Mineral and return the result.
func (m Mineral) AddToAll(amt int) Mineral {
	return Mineral{
		Ironium:   m.Ironium + amt,
		Boranium:  m.Boranium + amt,
		Germanium: m.Germanium + amt,
	}
}

// Add an int to a single component of a mineral and return the result.
func (m Mineral) AddNum(minType MineralType, amt int) Mineral {
	switch minType {
	case Ironium:
		m.Ironium += amt
	case Boranium:
		m.Boranium += amt
	case Germanium:
		m.Germanium += amt
	default:
		panic(fmt.Sprintf("mineral.AddNum called with invalid MineralType %q; \nmust be Ironium, Boranium or Germanium", minType))
	}
	return m
}

// Multiply all components of a mineral by a float64, round them using roundFunc and
// return the result truncated to an integer.
func (m Mineral) MultiplyFloat64(factor float64, roundFunc func(float64) float64) Mineral {
	return Mineral{
		Ironium:   int(roundFunc(float64(m.Ironium) * factor)),
		Boranium:  int(roundFunc(float64(m.Boranium) * factor)),
		Germanium: int(roundFunc(float64(m.Germanium) * factor)),
	}
}

func (m Mineral) Clamp(min, max int) Mineral {
	return Mineral{
		Ironium:   Clamp(m.Ironium, min, max),
		Boranium:  Clamp(m.Boranium, min, max),
		Germanium: Clamp(m.Germanium, min, max),
	}
}

// Attempt to equalize a Mineral's values as best as possible by repeatedly
// adding or subtracting amtToAdd in total.
// If amtToAdd is positive, it adds to the lowest values;
// if negative, it subtracts from the highest ones.
//
// Ties among equal values will be broken in order of precedence (I>B>G).
func (m Mineral) Equalize(amtToAdd int) Mineral {
	if amtToAdd == 0 {
		return m
	}

	/*
		Example scenario:
		19 Iron, 3 Bor & 31 Germ with 50 total.
		First, we add 17 Boranium to make it equal to Ironium.
		Next, we add 13 (31-19) to both Iron and Bor to equalize all 3.
		The remaining 11 is split evenly 3 ways (4 to I/B, 3 to G).
	*/

	mArray := m.ToSlice()
	mSlice := mArray[:]
	var origOrder = []int{0, 1, 2} // original value order; used to "un-shuffle" slice at the end
	// sort mineral values/types
	slices.SortFunc(mSlice, func(a, b int) int {
		diff := a - b
		if amtToAdd < 0 {
			diff = b - a // reverse sorting order for negative indices so we deduct from the highest
		}
		if diff < 0 {
			// shuffle around original order slice to keep it in sync
			// (3 values is small enough for go to use insertion sort)
			i := slices.Index(mSlice, a)
			origOrder[i], origOrder[i-1] = origOrder[i-1], origOrder[i]
		}
		return diff
	})

	addFunc := func(index, amt int) {
		mSlice[index] += amt
		amtToAdd -= amt
	}

	// attempt to equalize lowest 2 (highest 2 for negatives)
	diffLowest := AbsMin(amtToAdd, mSlice[1]-mSlice[0])
	if diffLowest != 0 {
		addFunc(0, diffLowest)
	}

	// lowest/middle now equal; try to equalize with highest
	diffMiddle := AbsMin(amtToAdd, (mSlice[2]-mSlice[1])*2)
	if diffMiddle != 0 {
		addFunc(0, diffMiddle/2)
		addFunc(1, diffMiddle/2)
	}

	// deal with any excess
	if third := amtToAdd / 3; third != 0 {
		for i := range mSlice {
			mSlice[i] += third
		}
		amtToAdd %= 3
	}
	for i := range Abs(amtToAdd) {
		if amtToAdd < 0 {
			// use original order so we add to iron first
			mSlice[origOrder[i]]--
		} else {
			mSlice[origOrder[i]]++
		}
	}
	return NewMineral(mSlice[origOrder[0]], mSlice[origOrder[1]], mSlice[origOrder[2]])
}

// HighestType returns the MineralType and numerical value of the
// Nth highest value in a Mineral struct.
// Negative indices count backwards from lowest value.
// (1 = highest, 2 = 2nd highest, -1 = lowest, etc etc).
//
// Ties are broken in order of precendence (I>B>G); tie order not affected by negative indices
//
// panics if ranking is 0 or if abs(ranking) is greater than 3
func (m Mineral) HighestType(ranking int) (minType MineralType, value int) {
	a := m.ToSlice()
	if ranking == 0 || Abs(ranking) > len(a) {
		panic(fmt.Sprintf("Mineral.HighestType() called with incorrect ranking %d; must be non-zero integer between -%d and %[2]d", ranking, len(a)))
	}
	slices.Sort(a[:])

	if ranking > 0 {
		value = a[len(a)-ranking] // Slice is ordered in ascending order, so biggest values will be at the end
	} else {
		value = a[-ranking-1] // negative indices count from the start (lowest first)
	}

	return m.GetTypeFromAmount(value), value
}

// return the first valid MineralType in a Mineral struct with the given numerical value;
// panics if no MineralType with the corresponding value exists
func (m Mineral) GetTypeFromAmount(amt int) MineralType {
	switch amt {
	case m.Ironium:
		return Ironium
	case m.Boranium:
		return Boranium
	case m.Germanium:
		return Germanium
	}
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v but no corresponding MineralType was found in mineral struct; Struct values: \n%#v", amt, m))
}
