package game

type TechLevel struct {
	Energy        int `json:"energy,omitempty"`
	Weapons       int `json:"weapons,omitempty"`
	Propulsion    int `json:"propulsion,omitempty"`
	Construction  int `json:"construction,omitempty"`
	Electronics   int `json:"electronics,omitempty"`
	Biotechnology int `json:"biotechnology,omitempty"`
}

// return true if this techlevel has the required techlevels for a requirements
func (tl *TechLevel) HasRequiredLevels(required TechLevel) bool {
	return tl.Energy >= required.Energy &&
		tl.Weapons >= required.Weapons &&
		tl.Propulsion >= required.Propulsion &&
		tl.Construction >= required.Construction &&
		tl.Electronics >= required.Electronics &&
		tl.Biotechnology >= required.Biotechnology
}

// return the minimum tech level
func (tl *TechLevel) Min() int {
	return minInt(
		tl.Energy,
		tl.Weapons,
		tl.Propulsion,
		tl.Construction,
		tl.Electronics,
		tl.Biotechnology)
}
