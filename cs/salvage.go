package cs

import "fmt"

type Salvage struct {
	GameDBObject `tstype:",extends"`
	MapObject    `tstype:",extends"`
	Cargo        Cargo `json:"cargo,omitempty"`
}

// create a new salvage object
func newSalvage(position Vector, num int, playerNum int, cargo Cargo) *Salvage {
	return &Salvage{
		MapObject: MapObject{
			Type:      MapObjectTypeSalvage,
			Position:  position,
			Num:       num,
			PlayerNum: playerNum,
			Name:      fmt.Sprintf("Salvage #%d", num),
		},
		Cargo: cargo,
	}
}

// decay this salvage
// https://wiki.starsautohost.org/wiki/Guts_of_scrapping
// In deep space, each type of mineral decays 10%, or 10kT per year, whichever is higher. Salvage deposited on planets does not decay.
func (salvage *Salvage) decay(rules *Rules) {
	salvage.Cargo = Cargo{
		Ironium: Max(0, Min(
			salvage.Cargo.Ironium-int(float64(salvage.Cargo.Ironium)*rules.SalvageDecayRate),
			salvage.Cargo.Ironium-rules.SalvageDecayMin,
		)),
		Boranium: Max(0, Min(
			salvage.Cargo.Boranium-int(float64(salvage.Cargo.Boranium)*rules.SalvageDecayRate),
			salvage.Cargo.Boranium-rules.SalvageDecayMin,
		)),
		Germanium: Max(0, Min(
			salvage.Cargo.Germanium-int(float64(salvage.Cargo.Germanium)*rules.SalvageDecayRate),
			salvage.Cargo.Germanium-rules.SalvageDecayMin,
		)),
	}
}
