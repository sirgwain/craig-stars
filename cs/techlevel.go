package cs

import "math"

// Represents a TechLevel the player has or a tech requires, or the amount of research spent on each tech level
type TechLevel struct {
	Energy        int `json:"energy,omitempty"`
	Weapons       int `json:"weapons,omitempty"`
	Propulsion    int `json:"propulsion,omitempty"`
	Construction  int `json:"construction,omitempty"`
	Electronics   int `json:"electronics,omitempty"`
	Biotechnology int `json:"biotechnology,omitempty"`
}

// return true if this techlevel has the required techlevels for a requirements
func (tl TechLevel) HasRequiredLevels(required TechLevel) bool {
	return tl.Energy >= required.Energy &&
		tl.Weapons >= required.Weapons &&
		tl.Propulsion >= required.Propulsion &&
		tl.Construction >= required.Construction &&
		tl.Electronics >= required.Electronics &&
		tl.Biotechnology >= required.Biotechnology
}

// return the minimum tech level
func (tl TechLevel) Min() int {
	return MinInt(
		tl.Energy,
		tl.Weapons,
		tl.Propulsion,
		tl.Construction,
		tl.Electronics,
		tl.Biotechnology)
}

// return the sum of all tech levels
func (tl TechLevel) Sum() int {
	return tl.Energy +
		tl.Weapons +
		tl.Propulsion +
		tl.Construction +
		tl.Electronics +
		tl.Biotechnology
}

func (tl TechLevel) Lowest() TechField {
	lowestField := Energy
	lowest := math.MaxInt
	for _, field := range TechFields {
		level := tl.Get(field)
		if level < lowest {
			lowestField = field
			lowest = level
		}
	}
	return lowestField
}

func (tl TechLevel) LowestNonZero() TechField {
	lowestField := Energy
	lowest := math.MaxInt
	for _, field := range TechFields {
		level := tl.Get(field)
		if level != 0 && level < lowest {
			lowestField = field
			lowest = level
		}
	}
	return lowestField
}

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

// add a tech level to this one
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

func (tl TechLevel) Minus(tl2 TechLevel) TechLevel {
	return TechLevel{
		tl.Energy - tl2.Energy,
		tl.Weapons - tl2.Weapons,
		tl.Propulsion - tl2.Propulsion,
		tl.Construction - tl2.Construction,
		tl.Electronics - tl2.Electronics,
		tl.Biotechnology - tl2.Biotechnology,
	}
}

func (tl TechLevel) Max(tl2 TechLevel) TechLevel {
	return TechLevel{
		Energy:        MaxInt(tl.Energy, tl2.Energy),
		Weapons:       MaxInt(tl.Weapons, tl2.Weapons),
		Propulsion:    MaxInt(tl.Propulsion, tl2.Propulsion),
		Construction:  MaxInt(tl.Construction, tl2.Construction),
		Electronics:   MaxInt(tl.Electronics, tl2.Electronics),
		Biotechnology: MaxInt(tl.Biotechnology, tl2.Biotechnology),
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

// Get the mininum levels above this tech
// i.e. if we just are a starter humanoid and just gained prop 5
// level ({ 3, 3, 6, 3, 3}) we are 0 levels above the Radiating Hyrdo-Ram Scoop
// returns maxInt if the tech is all 0
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
	return requirement.Minus(tl).MinZero().LowestNonZero()
}
