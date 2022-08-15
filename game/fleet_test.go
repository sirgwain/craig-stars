package game

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// create a new long rang scout fleet for testing
func testLongRangeScout(player *Player, rules *Rules) *Fleet {
	fleet := &Fleet{
		BaseName: "Long Range Scout",
		Tokens: []ShipToken{
			{Quantity: 1, Design: NewShipDesign(player).
				WithHull(Scout.Name).
				WithSlots([]ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				}).
				WithSpec(rules, player)},
		},
	}
	fleet.Spec = ComputeFleetSpec(rules, player, fleet)
	return fleet
}

// create a new small freighter (with cargo pod) fleet for testing
func testSmallFreighter(player *Player, rules *Rules) *Fleet {
	fleet := &Fleet{
		BaseName: "Small Freighter",
		Tokens: []ShipToken{
			{Quantity: 1, Design: NewShipDesign(player).
				WithHull(SmallFreighter.Name).
				WithSlots([]ShipDesignSlot{
					{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
				}).
				WithSpec(rules, player)},
		},
	}

	fleet.Spec = ComputeFleetSpec(rules, player, fleet)
	return fleet

}

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

func TestFleet_TransferPlanetCargo(t *testing.T) {
	rules := NewRules()
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	// scout := testLongRangeScout(player, &rules)
	freighter := testSmallFreighter(player, &rules)
	type args struct {
		planet         *Planet
		transferAmount Cargo
	}
	tests := []struct {
		name    string
		fleet   *Fleet
		args    args
		wantErr bool
	}{
		{"Should transfer from planet", freighter, args{NewPlanet(1).WithCargo(Cargo{1, 2, 3, 4}), Cargo{1, 0, 0, 0}}, false},
		{"Should fail to transfer from planet", freighter, args{NewPlanet(1).WithCargo(Cargo{1, 2, 3, 4}), Cargo{2, 0, 0, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceCargo := tt.fleet.Cargo
			destCargo := tt.args.planet.Cargo
			if err := tt.fleet.TransferPlanetCargo(tt.args.planet, tt.args.transferAmount); (err != nil) != tt.wantErr {
				t.Errorf("Fleet.TransferPlanetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// if we didn't want an error, ensure we transferred from the planet to the fleet
				assert.Equal(t, tt.fleet.Cargo, sourceCargo.Add(tt.args.transferAmount))
				assert.Equal(t, tt.args.planet.Cargo, destCargo.Subtract(tt.args.transferAmount))
			}

		})
	}
}

func TestFleet_TransferFleetCargo(t *testing.T) {
	rules := NewRules()
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	// scout := testLongRangeScout(player, &rules)
	freighter := testSmallFreighter(player, &rules)
	type args struct {
		fleet         *Fleet
		transferAmount Cargo
	}
	tests := []struct {
		name    string
		fleet   *Fleet
		args    args
		wantErr bool
	}{
		{"Should transfer from fleet", freighter, args{testSmallFreighter(player, &rules).WithCargo(Cargo{1, 2, 3, 4}), Cargo{1, 0, 0, 0}}, false},
		{"Should fail to transfer from fleet", freighter, args{testSmallFreighter(player, &rules).WithCargo(Cargo{1, 2, 3, 4}), Cargo{2, 0, 0, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceCargo := tt.fleet.Cargo
			destCargo := tt.args.fleet.Cargo
			if err := tt.fleet.TransferFleetCargo(tt.args.fleet, tt.args.transferAmount); (err != nil) != tt.wantErr {
				t.Errorf("Fleet.TransferFleetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// if we didn't want an error, ensure we transferred from the fleet to the fleet
				assert.Equal(t, tt.fleet.Cargo, sourceCargo.Add(tt.args.transferAmount))
				assert.Equal(t, tt.args.fleet.Cargo, destCargo.Subtract(tt.args.transferAmount))
			}

		})
	}
}
