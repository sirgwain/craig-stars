package cs

import (
	"fmt"
	"math"
)

// Represents a TechLevel a player has or a tech requires, or the amount of research spent on a tech level
type TechLevel struct {
	Energy        int `json:"energy,omitempty"`
	Weapons       int `json:"weapons,omitempty"`
	Propulsion    int `json:"propulsion,omitempty"`
	Construction  int `json:"construction,omitempty"`
	Electronics   int `json:"electronics,omitempty"`
	Biotechnology int `json:"biotechnology,omitempty"`
}

// return true if tl has the required levels for required
// (i.e. tl >= required for all fields)
func (tl TechLevel) HasRequiredLevels(required TechLevel) bool {
	return tl.Energy >= required.Energy &&
		tl.Weapons >= required.Weapons &&
		tl.Propulsion >= required.Propulsion &&
		tl.Construction >= required.Construction &&
		tl.Electronics >= required.Electronics &&
		tl.Biotechnology >= required.Biotechnology
}

// return the total of all tech levels
func (tl TechLevel) Total() int {
	return tl.Energy +
		tl.Weapons +
		tl.Propulsion +
		tl.Construction +
		tl.Electronics +
		tl.Biotechnology
}

func (tl TechLevel) ToSlice() [6]int {
	return [6]int{tl.Energy,
		tl.Weapons,
		tl.Propulsion,
		tl.Construction,
		tl.Electronics,
		tl.Biotechnology,
	}
}

// return the lowest numerical value in a TechLevel struct, including 0
// Ties are broken by order of precedence (En>We>Pr>Co>El>Bi)
func (tl TechLevel) LowestLevel() int {
	a := tl.ToSlice()
	return Min(a[:]...)
}

// return the TechField with the lowest numerical value in a TechLevel struct, including 0
// Ties are broken by order of precedence (En>We>Pr>Co>El>Bi)
func (tl TechLevel) Lowest() TechField {
	lowest := Energy
	lowestLevel := math.MaxInt
	for _, field := range TechFields {
		level := tl.Get(field)
		if lowestLevel > level {
			lowestLevel = level
			lowest = field
		}
	}

	return lowest
}

func (tl TechLevel) Lowest_alt() TechField {
	return tl.GetFieldFromAmount(tl.LowestLevel())
} 

// return the lowest positive TechField in a TechLevel struct.
//
// Ties are broken by order of precedence (En>We>Pr>Co>El>Bi)
func (tl TechLevel) LowestPositive() TechField {
	lowest := Energy
	lowestLevel := math.MaxInt
	for _, field := range TechFields {
		level := tl.Get(field)
		if lowestLevel > level && level > 0 {
			lowestLevel = level
			lowest = field
		}
	}

	return lowest
}

// get the level for the specified TechField
func (tl TechLevel) Get(field TechField) int {
	switch field {
	case Energy:
		return tl.Energy
	case Weapons:
		return tl.Weapons
	case Propulsion:
		return tl.Propulsion
	case Construction:
		return tl.Construction
	case Electronics:
		return tl.Electronics
	case Biotechnology:
		return tl.Biotechnology
	}
	return None
}

// return the first valid TechField in a TechLevel struct with the given numerical value;
// panics if no TechField with the corresponding value exists
func (tl TechLevel) GetFieldFromAmount(amt int) TechField {
	switch amt {
	case tl.Energy:
		return Energy
	case tl.Weapons:
		return Weapons
	case tl.Propulsion:
		return Propulsion
	case tl.Construction:
		return Construction
	case tl.Electronics:
		return Electronics
	case tl.Biotechnology:
		return Biotechnology
	}
	panic(fmt.Sprintf("GetFieldFromAmount called with value %v but no corresponding TechField was found in struct; \nStruct values: %v",
		amt, tl))
}

func (tl *TechLevel) Set(field TechField, level int) {
	switch field {
	case Energy:
		tl.Energy = level
	case Weapons:
		tl.Weapons = level
	case Propulsion:
		tl.Propulsion = level
	case Construction:
		tl.Construction = level
	case Electronics:
		tl.Electronics = level
	case Biotechnology:
		tl.Biotechnology = level
	}
}

// add together 2 TechLevels and return the result
func (tl TechLevel) Add(other TechLevel) TechLevel {
	return TechLevel{
		tl.Energy + other.Energy,
		tl.Weapons + other.Weapons,
		tl.Propulsion + other.Propulsion,
		tl.Construction + other.Construction,
		tl.Electronics + other.Electronics,
		tl.Biotechnology + other.Biotechnology,
	}
}

// deduct the given TechLevel from another TechLevel and return the result
func (tl TechLevel) Subtract(other TechLevel) TechLevel {
	return TechLevel{
		tl.Energy - other.Energy,
		tl.Weapons - other.Weapons,
		tl.Propulsion - other.Propulsion,
		tl.Construction - other.Construction,
		tl.Electronics - other.Electronics,
		tl.Biotechnology - other.Biotechnology,
	}
}

// Return greater of 2 TechLevel structs for all TechFields separately
func (tl TechLevel) Max(other TechLevel) TechLevel {
	return TechLevel{
		Energy:        Max(tl.Energy, other.Energy),
		Weapons:       Max(tl.Weapons, other.Weapons),
		Propulsion:    Max(tl.Propulsion, other.Propulsion),
		Construction:  Max(tl.Construction, other.Construction),
		Electronics:   Max(tl.Electronics, other.Electronics),
		Biotechnology: Max(tl.Biotechnology, other.Biotechnology),
	}
}

// Return lesser of 2 TechLevel structs for all TechFields separately
func (tl TechLevel) Min(other TechLevel) TechLevel {
	return TechLevel{
		Energy:        Min(tl.Energy, other.Energy),
		Weapons:       Min(tl.Weapons, other.Weapons),
		Propulsion:    Min(tl.Propulsion, other.Propulsion),
		Construction:  Min(tl.Construction, other.Construction),
		Electronics:   Min(tl.Electronics, other.Electronics),
		Biotechnology: Min(tl.Biotechnology, other.Biotechnology),
	}
}

// return this TechLevel with a minimum of zero for each value
func (tl TechLevel) MinZero() TechLevel {
	return TechLevel{
		Energy:        Max(tl.Energy, 0),
		Weapons:       Max(tl.Weapons, 0),
		Propulsion:    Max(tl.Propulsion, 0),
		Construction:  Max(tl.Construction, 0),
		Electronics:   Max(tl.Electronics, 0),
		Biotechnology: Max(tl.Biotechnology, 0),
	}

}

// Get the lowest amount of levels tl is above other.
// This assumes tl is above other in all levels; it's just finding the lowest non-zero field above
//
// Returns maxInt if other is all 0s
func (tl TechLevel) LevelsAbove(other TechLevel) int {
	levelsAbove := math.MaxInt
	if other.Energy > 0 {
		levelsAbove = Min(levelsAbove, tl.Energy-other.Energy)
	}
	if other.Weapons > 0 {
		levelsAbove = Min(levelsAbove, tl.Weapons-other.Weapons)
	}
	if other.Propulsion > 0 {
		levelsAbove = Min(levelsAbove, tl.Propulsion-other.Propulsion)
	}
	if other.Construction > 0 {
		levelsAbove = Min(levelsAbove, tl.Construction-other.Construction)
	}
	if other.Electronics > 0 {
		levelsAbove = Min(levelsAbove, tl.Electronics-other.Electronics)
	}
	if other.Biotechnology > 0 {
		levelsAbove = Min(levelsAbove, tl.Biotechnology-other.Biotechnology)
	}

	return levelsAbove
}

// Return the number of levels other is above tl in the given field.
func (tl TechLevel) LevelsAboveField(other TechLevel, field TechField) int {
	switch field {
	case Energy:
		return other.Energy - tl.Energy
	case Weapons:
		return other.Weapons - tl.Weapons
	case Propulsion:
		return other.Propulsion - tl.Propulsion
	case Construction:
		return other.Construction - tl.Construction
	case Electronics:
		return other.Electronics - tl.Electronics
	case Biotechnology:
		return other.Biotechnology - tl.Biotechnology
	}
	return math.MaxInt
}

// get all the learnable tech fields for a player
func (tl TechLevel) LearnableTechFields(rules *Rules) []TechField {
	fields := make([]TechField, 0, len(TechFields))
	for _, field := range TechFields {
		if tl.Get(field) < rules.MaxTechLevel {
			fields = append(fields, field)
		}
	}
	return fields
}

// get the lowest field missing from tl for a requirement
func (tl TechLevel) LowestMissingLevel(requirement TechLevel) TechField {
	return requirement.Subtract(tl).LowestPositive()
}
