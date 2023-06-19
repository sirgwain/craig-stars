package cs

type Salvage struct {
	MapObject
	Cargo Cargo `json:"cargo,omitempty"`
}

// create a new salvage object
func newSalvage(position Vector, num int, playerNum int, cargo Cargo) *Salvage {
	return &Salvage{
		MapObject: MapObject{
			Type:      MapObjectTypeSalvage,
			Position:  position,
			Num:       num,
			PlayerNum: playerNum,
			Dirty:     true,
		},
		Cargo: cargo,
	}
}

// decay this salvage
// https://wiki.starsautohost.org/wiki/Guts_of_scrapping
// In deep space, each type of mineral decays 10%, or 10kT per year, whichever is higher. Salvage deposited on planets does not decay.
func (salvage *Salvage) decay(rules *Rules) {
	salvage.Cargo = Cargo{
		Ironium: maxInt(0, minInt(
			salvage.Cargo.Ironium-int(float64(salvage.Cargo.Ironium)*rules.SalvageDecayRate),
			salvage.Cargo.Ironium-rules.SalvageDecayMin,
		)),
		Boranium: maxInt(0, minInt(
			salvage.Cargo.Boranium-int(float64(salvage.Cargo.Boranium)*rules.SalvageDecayRate),
			salvage.Cargo.Boranium-rules.SalvageDecayMin,
		)),
		Germanium: maxInt(0, minInt(
			salvage.Cargo.Germanium-int(float64(salvage.Cargo.Germanium)*rules.SalvageDecayRate),
			salvage.Cargo.Germanium-rules.SalvageDecayMin,
		)),
	}
}
