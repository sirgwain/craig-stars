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

func NewMineral(values [3]int) Mineral {
	return Mineral{
		Ironium:   values[0],
		Boranium:  values[1],
		Germanium: values[2],
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

func (m *Mineral) Set(mineralType MineralType, value int) *Mineral {
	switch mineralType {
	case Ironium:
		m.Ironium = value
	case Boranium:
		m.Boranium = value
	case Germanium:
		m.Germanium = value
	}
	return m
}

// return higher of 2 Mineral structs for all MineralTypes separately
func (m Mineral) Max(other Mineral) Mineral {
	return Mineral{
		Ironium:   Max(m.Ironium, other.Ironium),
		Boranium:  Max(m.Boranium, other.Boranium),
		Germanium: Max(m.Germanium, other.Germanium),
	}
}

func (m Mineral) GetAmount(mineralType MineralType) int {
	var amt int
	switch mineralType {
	case Ironium:
		amt = m.Ironium
	case Boranium:
		amt = m.Boranium
	case Germanium:
		amt = m.Germanium
	}
	return amt
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

// convert a mineral to a cargo
func (m Mineral) ToCargo() Cargo {
	return Cargo{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// convert a mineral to a cost
func (m Mineral) ToCost() Cost {
	return Cost{
		Ironium:   m.Ironium,
		Boranium:  m.Boranium,
		Germanium: m.Germanium,
	}
}

// add two minerals
func (m Mineral) Add(m2 Mineral) Mineral {
	return Mineral{
		Ironium:   m.Ironium + m2.Ironium,
		Boranium:  m.Boranium + m2.Boranium,
		Germanium: m.Germanium + m2.Germanium,
	}
}

// add an int to all components of the mineral
func (m Mineral) AddInt(num int) Mineral {
	return Mineral{
		Ironium:   m.Ironium + num,
		Boranium:  m.Boranium + num,
		Germanium: m.Germanium + num,
	}
}

// subtract two minerals
func (m Mineral) Subtract(m2 Mineral) Mineral {
	return Mineral{
		Ironium:   m.Ironium - m2.Ironium,
		Boranium:  m.Boranium - m2.Boranium,
		Germanium: m.Germanium - m2.Germanium,
	}
}

// subtract the mineral components of a Cost
func (m Mineral) SubtractCost(m2 Cost) Mineral {
	return Mineral{
		Ironium:   m.Ironium - m2.Ironium,
		Boranium:  m.Boranium - m2.Boranium,
		Germanium: m.Germanium - m2.Germanium,
	}
}

func (c Mineral) MultiplyFloat64(factor float64) Mineral {
	return Mineral{
		Ironium:   int(float64(c.Ironium) * factor),
		Boranium:  int(float64(c.Boranium) * factor),
		Germanium: int(float64(c.Germanium) * factor),
	}
}

func (m Mineral) Clamp(min, max int) Mineral {
	return Mineral{
		Ironium:   Clamp(m.Ironium, min, max),
		Boranium:  Clamp(m.Boranium, min, max),
		Germanium: Clamp(m.Germanium, min, max),
	}
}

// return the MineralType with the Nth highest numerical value in a Mineral struct (1 = highest, 2 = 2nd highest, etc etc)
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (I>B>G); tie order not affected by negative indices
//
// panics if ranking is 0 or if abs(ranking) is greater than 3
func (m Mineral) HighestType(ranking int) MineralType {
	if ranking == 0 || Abs(ranking) > 3 {
		panic(fmt.Sprintf("HighestType called with incorrect ranking %d; must be non-zero integer between -3 and 3", ranking))
	}
	return m.GetTypeFromAmount(m.HighestAmount(ranking))
}

// return the numerical value of the Nth highest MineralType in a Mineral struct (1 = highest, 2 = 2nd highest, etc).
// Negative indices count backwards from lowest value  (-1 = lowest, -2 = 2nd lowest, etc).

// panics if ranking is 0 or abs(ranking) is greater than 3
func (m Mineral) HighestAmount(ranking int) int {
	if ranking == 0 || Abs(ranking) > 3 {
		panic(fmt.Sprintf("HighestAmount called with incorrect ranking %d; must be non-zero integer between -3 and 3", ranking))
	}
	a := m.ToSlice()
	slices.Sort(a[:])
	if ranking > 0 {
		return a[3-ranking] // Slice is ordered in ascending order, so biggest values will be at the end 
	} else {
		return a[-ranking-1] // negative indices count from the start (lowest first)
	}
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
