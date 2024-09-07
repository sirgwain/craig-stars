package cs

import (
	"reflect"
	"testing"
)

func TestTechStore_GetBestScanner(t *testing.T) {
	type args struct {
		player *Player
	}
	tests := []struct {
		name string
		args args
		want *TechHullComponent
	}{
		{"Get Lowest Scanner", args{NewPlayer(1, NewRace())}, &BatScanner},
		{"Get Nicer Scanner", args{NewPlayer(1, NewRace()).WithTechLevels(TechLevel{Electronics: 5})}, &PossumScanner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := StaticTechStore.GetBestScanner(tt.args.player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetBestScanner() = %v, want %v", got, tt.want)
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

func TestTechStore_GetBestEngine(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StaticTechStore.GetBestEngine(tt.args.player, tt.args.hull, tt.args.purpose); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TechStore.GetBestEngine() = %v, want %v", got, tt.want)
			}
		})
	}
}
