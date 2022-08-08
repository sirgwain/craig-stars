package game

import (
	"reflect"
	"testing"
)

func TestComputeFleetSpec(t *testing.T) {
	rules := NewRules()
	starterHumanoidPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3})
	starterHumanoidPlayer.Race.Spec = computeRaceSpec(&starterHumanoidPlayer.Race, &rules)

	type args struct {
		rules  *Rules
		player *Player
		fleet  *Fleet
	}
	tests := []struct {
		name string
		args args
		want FleetSpec
	}{
		{"empty", args{&rules, NewPlayer(1, NewRace().WithSpec(&rules)), &Fleet{}}, FleetSpec{ShipDesignSpec: ShipDesignSpec{
			ScanRange:    NoScanner,
			ScanRangePen: NoScanner,
			SpaceDock:    UnlimitedSpaceDock,
		}}},
		{"Starter Humanoid Long Range Scout", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Long Range Scout",
			Tokens: []ShipToken{
				{Quantity: 1, Design: NewShipDesign(starterHumanoidPlayer).
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:         Cost{15, 2, 7, 19},
				FuelCapacity: 300,
				ScanRange:    66,
				ScanRangePen: 30,
				Mass:         20,
				Armor:        20,
				Engine:       "Quick Jump 5",
			},
			MassEmpty:        20,
			BaseCloakedCargo: 20,
			TotalShips:       1,
		}},
		{"Starter Starbase", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Starbase",
			Tokens: []ShipToken{
				{Quantity: 1, Design: NewShipDesign(starterHumanoidPlayer).
					WithHull(SpaceStation.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
						{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
						{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
					}).
					WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:         Cost{122, 263, 236, 752},
				Mass:         48,
				Armor:        500,
				Shield:       400,
				SpaceDock:    UnlimitedSpaceDock,
				RepairBonus:  .15,
				MineSweep:    640,
				ScanRange:    NoScanner,
				ScanRangePen: NoScanner,
				HasWeapons:   true,
			},
			MassEmpty:        48,
			BaseCloakedCargo: 48,
			TotalShips:       1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := *ComputeFleetSpec(tt.args.rules, tt.args.player, tt.args.fleet); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeFleetSpec() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
