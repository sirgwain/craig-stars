package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTechStore_GetBestEngine(t *testing.T) {
	// TODO: Fix this someday
	type args struct {
		player  *Player
		hull    *TechHull
		purpose FleetPurpose
	}
	tests := []struct {
		name string
		args args
		want *TechEngine
	}{
		{"Base scout", args{testPlayer(), &Scout, FleetPurposeScout}, &QuickJump5},
		{"Mini Colonizer", args{NewPlayer(0, NewRace().WithPRT(HE).WithSpec(&rules)).withSpec(&rules), &MiniColonyShip, FleetPurposeColonizer}, &SettlersDelight},
		{"Max Techs", args{NewPlayer(0, NewRace().WithPRT(HE).WithLRT(IFE).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26}).withSpec(&rules), &Frigate, FleetPurposeColonizer}, &GalaxyScoop},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StaticTechStore.GetBestEngine(tt.args.player, tt.args.hull, tt.args.purpose); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetBestEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechStore_GetBestBattleEngine(t *testing.T) {
	// TODO: Fix this someday
	type args struct {
		player *Player
		hull   *TechHull
	}
	tests := []struct {
		name    string
		args    args
		want    *TechEngine
		mtTechs bool
	}{
		{"Base scout", args{testPlayer(), &Scout}, &QuickJump5, false},
		{"Prop 5 vs mizer", args{NewPlayer(0, NewRace().WithPRT(JoaT).WithLRT(IFE).WithSpec(&rules)).WithTechLevels(TechLevel{0,0,5,5,0,0}).withSpec(&rules), &Destroyer}, &DaddyLongLegs7, false},
		{"Max Techs", args{NewPlayer(0, NewRace().WithPRT(SD).WithLRT(IFE).WithSpec(&rules)).WithTechLevels(TechLevel{26, 26, 26, 26, 26, 26}).withSpec(&rules), &Nubian}, &EnigmaPulsar, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mtTechs {
				for _, tech := range MysteryTraderTechs {
					tt.args.player.AcquiredTechs[tech.Name] = true
				}
			}
			if got := StaticTechStore.GetBestBattleEngine(tt.args.player, tt.args.hull, 1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetBestBattleEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechStore_GetHullComponentsByCategory(t *testing.T) {
	type args struct {
		category TechCategory
	}
	tests := []struct {
		name string
		args args
		want TechHullComponent
	}{
		{"first shield", args{TechCategoryShield}, MoleSkinShield},
		{"first armor", args{TechCategoryArmor}, Tritanium},
		{"first beam", args{TechCategoryBeamWeapon}, Laser},
		{"first torpedo", args{TechCategoryTorpedo}, AlphaTorpedo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewTechStore()
			if got := store.GetHullComponentsByCategory(tt.args.category)[0]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetHullComponentsByCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTechStore_GetHullComponentsByHullSlotType(t *testing.T) {
	type fields struct {
		techLevels TechLevel
		race       *Race
		mtTechs    bool
	}
	type args struct {
		slot HullSlotType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*TechHullComponent
	}{
		{name: "max tech MT shields with IS", fields: fields{TechLevel{26, 26, 26, 26, 26, 26}, NewRace().WithPRT(IS), true},
			args: args{slot: HullSlotTypeShield},
			want: []*TechHullComponent{&MoleSkinShield, &CowHideShield, &WolverineDiffuseShield, &CrobySharmor, &BearNeutrinoBarrier, &LangstonShell, &GorillaDelagator, &ElephantHideFortress, &CompletePhaseShield}},
		{name: "Default Shields/Armors", fields: fields{TechLevel{26, 26, 26, 26, 26, 26}, NewRace().WithPRT(JoaT), false},
			args: args{slot: HullSlotTypeShieldArmor}, want: []*TechHullComponent{&Tritanium, &Crobmnium, &Carbonic, &Strobnium, &Organic, &Kelarium, &Neutronium, &Valanium, &Superlatanium,
				&MoleSkinShield, &CowHideShield, &WolverineDiffuseShield, &BearNeutrinoBarrier, &GorillaDelagator, &ElephantHideFortress, &CompletePhaseShield}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.fields.race.WithSpec(&rules)).WithTechLevels(tt.fields.techLevels)
			if tt.fields.mtTechs {
				for _, tech := range MysteryTraderTechs {
					player.AcquiredTechs[tech.Name] = true
				}
			}
			got := rules.techs.GetHullComponentsByHullSlotType(player, tt.args.slot, "Nubian")
			assert.ElementsMatch(t,got, tt.want)
		})
	}
}
