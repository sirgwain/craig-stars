package cs

import "math"

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
	return minInt(
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
func (tl *TechLevel) Add(other TechLevel) {
	tl.Energy += other.Energy
	tl.Weapons += other.Weapons
	tl.Propulsion += other.Propulsion
	tl.Construction += other.Construction
	tl.Electronics += other.Electronics
	tl.Biotechnology += other.Biotechnology
}
