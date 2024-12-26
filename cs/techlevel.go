package cs

import (
	"fmt"
	"math"
	"slices"
)

// Represents a TechLevel the player has or a tech requires, or the amount of research spent on each tech level
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

// return the TechField with the Nth highest numerical value in a TechLevel struct (1 = highest, 2 = 2nd highest, etc etc).
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (En>We>Pr>Co>El>Bi); tie order not affected by negative indices
func (tl TechLevel) HighestType(ranking int) TechField {
	return tl.GetTypeFromAmount(tl.HighestAmount(ranking))
}

// return the numerical value of the Nth highest TechField in a TechLevel struct (1 = highest, 2 = 2nd highest, etc etc)
// Negative indices count backwards from lowest value
//
// Ties are broken in order of precendence (En>We>Pr>Co>El>Bi); tie order not affected by negative indices
func (tl TechLevel) HighestAmount(ranking int) int {
	a := tl.ToSlice()
	slice := slices.Clone(a[:])
	if ranking < 0 {
		slices.SortStableFunc(slice, func(a, b int) int { return b - a })
		ranking = -ranking
	} else {
		slices.Sort(slice)
	}
	return slice[len(slice)-ranking]
}

// return the lowest positive TechField in a TechLevel struct
func (tl TechLevel) LowestPositive() TechField {
	a := tl.ToSlice()
	l := slices.Clone(a[:])
	l = slices.DeleteFunc(l, func(i int) bool { return i <= 0 })
	if len(l) == 0 {
		l = append(l, tl.Energy) // in the event we have nothing, return energy as a failsafe
	}
	return tl.GetTypeFromAmount(slices.Min(l))
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

// return the first valid TechField in a TechLevel struct with the given numerical value
// panics if no TechField with the corresponding value exists
func (tl TechLevel) GetTypeFromAmount(amt int) TechField {
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
	panic(fmt.Sprintf("GetTypeFromAmount called with value %v but no corresponding TechField was found in struct; \nStruct values: %v",
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
		Energy:        MaxInt(tl.Energy, other.Energy),
		Weapons:       MaxInt(tl.Weapons, other.Weapons),
		Propulsion:    MaxInt(tl.Propulsion, other.Propulsion),
		Construction:  MaxInt(tl.Construction, other.Construction),
		Electronics:   MaxInt(tl.Electronics, other.Electronics),
		Biotechnology: MaxInt(tl.Biotechnology, other.Biotechnology),
	}
}

// Return lesser of 2 TechLevel structs for all TechFields separately
func (tl TechLevel) Min(other TechLevel) TechLevel {
	return TechLevel{
		Energy:        MinInt(tl.Energy, other.Energy),
		Weapons:       MinInt(tl.Weapons, other.Weapons),
		Propulsion:    MinInt(tl.Propulsion, other.Propulsion),
		Construction:  MinInt(tl.Construction, other.Construction),
		Electronics:   MinInt(tl.Electronics, other.Electronics),
		Biotechnology: MinInt(tl.Biotechnology, other.Biotechnology),
	}
}

// return this TechLevel with a minimum of zero for each value
func (tl TechLevel) MinZero() TechLevel {
	return TechLevel{
		Energy:        MaxInt(tl.Energy, 0),
		Weapons:       MaxInt(tl.Weapons, 0),
		Propulsion:    MaxInt(tl.Propulsion, 0),
		Construction:  MaxInt(tl.Construction, 0),
		Electronics:   MaxInt(tl.Electronics, 0),
		Biotechnology: MaxInt(tl.Biotechnology, 0),
	}

}

// Get the lowest amount of levels tl is above other.
// 
// returns maxInt if other is all 0s
func (tl TechLevel) LevelsAbove(other TechLevel) int {
	levelsAbove := math.MaxInt
	if tl.Energy != 0 {
		levelsAbove = MinInt(levelsAbove, other.Energy-tl.Energy)
	}
	if tl.Weapons != 0 {
		levelsAbove = MinInt(levelsAbove, other.Weapons-tl.Weapons)
	}
	if tl.Propulsion != 0 {
		levelsAbove = MinInt(levelsAbove, other.Propulsion-tl.Propulsion)
	}
	if tl.Construction != 0 {
		levelsAbove = MinInt(levelsAbove, other.Construction-tl.Construction)
	}
	if tl.Electronics != 0 {
		levelsAbove = MinInt(levelsAbove, other.Electronics-tl.Electronics)
	}
	if tl.Biotechnology != 0 {
		levelsAbove = MinInt(levelsAbove, other.Biotechnology-tl.Biotechnology)
	}
	return levelsAbove
}

// LevelsAboveField returns the levels we are above a tech in a given field, or MaxInt if the field requirement is 0
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
	default:
		return 0
	}

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
