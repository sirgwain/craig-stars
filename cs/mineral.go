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

// Set sets the value corresponding to minType to amount.
// Unlike all other Mineral functions, this _will_ mutate the original struct's values,
// and is best used for more complex cases not handled by simple addition.
func (m *Mineral) Set(minType MineralType, amount int) {
	switch minType {
	case Ironium:
		m.Ironium = amount
	case Boranium:
		m.Boranium = amount
	case Germanium:
		m.Germanium = amount
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
		Ironium:   max(m.Ironium, other.Ironium),
		Boranium:  max(m.Boranium, other.Boranium),
		Germanium: max(m.Germanium, other.Germanium),
	}
}

// MaxNum return the higher of num and this Mineral struct's values
// for each MineralType.
//
//	Mineral{1, 2, 3}.MaxNum(2) = Mineral{2, 2, 3}
func (m Mineral) MaxNum(num int) Mineral {
	return Mineral{
		Ironium:   max(m.Ironium, num),
		Boranium:  max(m.Boranium, num),
		Germanium: max(m.Germanium, num),
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
func (m Mineral) AddToAll(amount int) Mineral {
	return Mineral{
		Ironium:   m.Ironium + amount,
		Boranium:  m.Boranium + amount,
		Germanium: m.Germanium + amount,
	}
}

// Add an int to a single component of a mineral and return the result.
func (m Mineral) AddNum(minType MineralType, amount int) Mineral {
	switch minType {
	case Ironium:
		m.Ironium += amount
	case Boranium:
		m.Boranium += amount
	case Germanium:
		m.Germanium += amount
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

// HighestType returns the MineralType and numerical value of the
// Nth highest value in a Mineral struct.
// Negative indices count backwards from lowest value.
// (1 = highest, 2 = 2nd highest, -1 = lowest, etc etc).
//
// Ties are broken in order of precendence (I>B>G); tie order not affected by negative indices
//
// panics if ranking is 0 or if abs(ranking) is greater than 3/
//
// Also see [Cost.HighestType]
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
func (m Mineral) GetTypeFromAmount(amount int) MineralType {
	switch amount {
	case m.Ironium:
		return Ironium
	case m.Boranium:
		return Boranium
	case m.Germanium:
		return Germanium
	}
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v, but no corresponding MineralType was found in mineral struct: \n%#v", amount, m))
}
