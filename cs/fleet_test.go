package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

// create a new long range scout fleet for testing with pre-computed specs.
func testLongRangeScout(player *Player) *Fleet {
	return testLongRangeScoutWithQuantity(player, 1)
}

// create a new long range scout ship design for testing.
// Does NOT come with precomputed specs; those will have to be done manually.
//
// Note: If putting this inside a ShipToken in a fleet in the universe,
// make sure to set the player's Num field to the same as this
// (otherwise the token design will be nil)
func testLongRangeScoutDesign(playerNum int) *ShipDesign {
	return NewShipDesign(playerNum, 1).
		WithName("Long Range Scout").
		WithHull(Scout.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
		})
}

func testLongRangeScoutWithQuantity(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Long Range Scout",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Long Range Scout").
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player),
			},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	player.Designs = append(player.Designs, fleet.Tokens[0].design)
	return fleet
}

// create a new small freighter (with cargo pod) fleet for testing
func testSmallFreighter(player *Player) *Fleet {
	return testSmallFreighterWithQuantity(player, 1)
}

func testSmallFreighterWithQuantity(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "Small Freighter",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Small Freighter").
					WithHull(SmallFreighter.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	player.Designs = append(player.Designs, fleet.Tokens[0].design)
	return fleet
}

func testStealingFreighter(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "La Robba",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Stealing Freighter").
					WithHull(MediumFreighter.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: RobberBaronScanner.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

// Create a new Nubian design for fuel-cost related testing with the specified engine.
// The ship always has enough components to weigh exactly 200 kT,
// resulting in a clean 1mG/ly at 100% engine efficiency.
//
// As with [testLongRangeScoutDesign], the design's spec is left
// uncomputed and must be done by the caller.
func testFuelNubianDesign(playerNum int, engine *TechEngine) *ShipDesign {
	if engine.Mass > 34 {
		panic("engine too heavy (nubian weighs 100 kT baseline and needs 3 engines)")
	}

	design := NewShipDesign(playerNum, 1).
		WithName("My 200kT life").
		WithHull(Nubian.Name)

	slots := []ShipDesignSlot{
		{HullComponent: engine.Name, HullSlotIndex: 1, Quantity: 3},
	}

	remainingMass := 200 - (100 + engine.Mass*3) // 100 kT baseline + engine mass

	for i := 1; remainingMass > 0 && i < len(Nubian.Slots); i++ {
		slot := ShipDesignSlot{HullSlotIndex: i + 1} // start from 2nd slot, but also 1-indexed

		if remainingMass/5 > 0 {
			// add rhino scanner (5 kT)
			num := min(3, remainingMass/5)
			slot.HullComponent = RhinoScanner.Name
			slot.Quantity = num
			remainingMass -= 5 * num
		} else {
			// add lasers (1kT)
			num := min(3, remainingMass)
			slot.HullComponent = Laser.Name
			slot.Quantity = num
			remainingMass -= num
		}
		slots = append(slots, slot)
	}

	return design.WithSlots(slots)
}

func Test_testFuelNubianDesign(t *testing.T) {
	for i := range StaticTechStore.Engines {
		engine := &StaticTechStore.Engines[i]
		t.Run(engine.Name, func(t *testing.T) {
			design := testFuelNubianDesign(1, engine).WithSpec(&rules, testPlayer())
			if design.Spec.Mass != 200 {
				t.Errorf("testNubianDesign returned ship with mass %d (want 200); \nEngine: %s\nSlots: %v", design.Spec.Mass, engine, design.Slots)
			}
		})
	}
}

// create a new Galleon (with fuel scoop) fleet for testing
func testGalleon(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "Galleon",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Galleon").
					WithHull(Galleon.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: SubGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

// create a new mini mine layer fleet for testing
func testMiniMineLayer(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "Little Hen",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Little Hen").
					WithHull(MiniMineLayer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: MineDispenser40.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: MineDispenser40.Name, HullSlotIndex: 3, Quantity: 2},
						{HullComponent: BatScanner.Name, HullSlotIndex: 4, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

func testCloakedScout(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
		},
		BaseName: "Cloaked Scout",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player.Num, 1).
					WithName("Cloaked Scout").
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: StealthCloak.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		OrbitingPlanetNum: None,
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func testRemoteTerraformer(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Remote Terraformer",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Remote Terraformer").
					WithHull(MiniMiner.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: BatScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: OrbitalAdjuster.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: OrbitalAdjuster.Name, HullSlotIndex: 4, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func testGatePrivateer(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "Gate Privateer",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Gate Privateer").
					WithHull(Privateer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: JumpGate.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
		FleetOrders: FleetOrders{
			Waypoints: []Waypoint{
				NewPositionWaypoint(Vector{}, 5),
			},
		},
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func Test_computeFleetSpec(t *testing.T) {
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
		{"empty", args{&rules, NewPlayer(1, NewRace().WithSpec(&rules)), &Fleet{}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				ScanRangePen:   NoScanner,
				SpaceDock:      UnlimitedSpaceDock,
				EstimatedRange: Infinite,
				ReduceCloaking: 1,
			},
			Purposes: map[ShipDesignPurpose]bool{},
		}},
		{"Starter Humanoid Long Range Scout", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Long Range Scout",
			Tokens: []ShipToken{
				{
					DesignNum: 1,
					Quantity:  1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
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
				Cost:           Cost{15, 2, 7, 19},
				FuelCapacity:   300,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           20,
				Armor:          20,
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        20,
			BaseCloakedCargo: 20,
			TotalShips:       1,
		}},
		{"Starter Starbase", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Starbase",
			Tokens: []ShipToken{
				{
					Quantity:  1,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
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
				Starbase:       true,
				Cost:           Cost{114, 150, 228, 624},
				Mass:           48,
				Armor:          500,
				Shields:        400,
				SpaceDock:      UnlimitedSpaceDock,
				RepairBonus:    .15,
				MineSweep:      640,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				HasWeapons:     true,
				MaxPopulation:  1_000_000,
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        48,
			BaseCloakedCargo: 48,
			TotalShips:       1,
		}},

		{"Cloaked Scout", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Cloaked Scout",
			Tokens: []ShipToken{
				{
					Quantity:  1,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(Scout.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: StealthCloak.Name, HullSlotIndex: 3, Quantity: 1},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{12, 2, 9, 20},
				FuelCapacity:   50,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           19,
				Armor:          20,
				CloakUnits:     70,
				CloakPercent:   35,
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        19,
			BaseCloakedCargo: 0,
			TotalShips:       1,
		}},
		{"2 Cloaked Scouts", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Cloaked Scout",
			Tokens: []ShipToken{
				{
					Quantity:  2,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(Scout.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: StealthCloak.Name, HullSlotIndex: 3, Quantity: 1},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{24, 4, 18, 40},
				FuelCapacity:   100,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           38,
				Armor:          40,
				CloakUnits:     70,
				CloakPercent:   35, // still 35%
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        38,
			BaseCloakedCargo: 0,
			TotalShips:       2,
		}},
		{"1 Cloaked Scout + 1 uncloaked scout", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Cloaked Scout",
			Tokens: []ShipToken{
				{
					Quantity:  1,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(Scout.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: StealthCloak.Name, HullSlotIndex: 3, Quantity: 1},
						}).
						WithSpec(&rules, starterHumanoidPlayer),
				},
				{
					Quantity:  1,
					DesignNum: 2,
					design: NewShipDesign(starterHumanoidPlayer.Num, 2).
						WithHull(Scout.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
							{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
						}).
						WithSpec(&rules, starterHumanoidPlayer),
				},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{15, 2, 7, 19},
				FuelCapacity:   350,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           39,
				Armor:          40,
				CloakUnits:     70,
				CloakPercent:   23, // uncloaked ship counts as cargo
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        39,
			BaseCloakedCargo: 20,
			TotalShips:       2,
		}},
		{"0 Scouts", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Long Range Scout",
			Tokens: []ShipToken{
				{
					DesignNum: 1,
					Quantity:  0,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
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
				Cost:           Cost{},
				FuelCapacity:   0,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           0,
				Armor:          0,
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
				EstimatedRange: Infinite,
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        0,
			BaseCloakedCargo: 0,
			TotalShips:       0,
		}},
		{"Bomber", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Bomber",
			Tokens: []ShipToken{
				{
					Quantity:  1,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(MiniBomber.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: LadyFingerBomb.Name, HullSlotIndex: 2, Quantity: 2},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{22, 43, 9, 42},
				FuelCapacity:   120,
				Mass:           112,
				Armor:          50,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				Bomber:         true,
				Bombs: []Bomb{
					{Quantity: 2, KillRate: .6, MinKillRate: 300, StructureDestroyRate: 2},
				},
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        112,
			BaseCloakedCargo: 112,
			TotalShips:       1,
		}},
		{"2 Bombers", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Bomber",
			Tokens: []ShipToken{
				{
					Quantity:  2,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(MiniBomber.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: LadyFingerBomb.Name, HullSlotIndex: 2, Quantity: 2},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{44, 86, 18, 84},
				FuelCapacity:   240,
				Mass:           224,
				Armor:          100,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				Bomber:         true,
				Bombs: []Bomb{
					{Quantity: 4, KillRate: .6, MinKillRate: 300, StructureDestroyRate: 2},
				},
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        224,
			BaseCloakedCargo: 224,
			TotalShips:       2,
		}},
		{"2 B52 Bombers with multiple bomb types", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Bomber",
			Tokens: []ShipToken{
				{
					Quantity:  2,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer.Num, 1).
						WithHull(B52Bomber.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 4},
							{HullComponent: CherryBomb.Name, HullSlotIndex: 2, Quantity: 4},
							{HullComponent: LBU32Bomb.Name, HullSlotIndex: 3, Quantity: 4},
							{HullComponent: NeutronBomb.Name, HullSlotIndex: 4, Quantity: 4},
							{HullComponent: RetroBomb.Name, HullSlotIndex: 5, Quantity: 4},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{348, 782, 228, 1392},
				FuelCapacity:   1500,
				Mass:           1764,
				Armor:          900,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				Bomber:         true,
				Bombs: []Bomb{
					{Quantity: 8, KillRate: 2.5, MinKillRate: 300, StructureDestroyRate: 10},
					{Quantity: 8, KillRate: .3, StructureDestroyRate: 28},
				},
				SmartBombs: []Bomb{
					{Quantity: 8, KillRate: 2.2},
				},
				RetroBombs: []Bomb{
					{Quantity: 8, UnterraformRate: 1},
				},
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        1764,
			BaseCloakedCargo: 1764,
			TotalShips:       2,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeFleetSpec(tt.args.rules, tt.args.player, tt.args.fleet)
			test.CompareAsJSON(t, got, tt.want)
		})
	}
}

func TestFleet_moveFleet(t *testing.T) {
	type want struct {
		position          Vector
		fuelUsed          int
		orbitingPlanetNum int
	}
	tests := []struct {
		name        string
		design      *ShipDesign
		destination Vector
		warpSpeed   int
		initialFuel int
		planet      *Planet
		want        want
	}{
		{
			name:        "move 25ly at warp5",
			design:      testLongRangeScoutDesign(1),
			destination: Vector{50, 0},
			warpSpeed:   5,
			initialFuel: 300,
			want:        want{position: Vector{25, 0}, fuelUsed: 4, orbitingPlanetNum: None},
		},
		{
			name:        "move 1ly at warp 1",
			design:      testLongRangeScoutDesign(1),
			destination: Vector{1, 1},
			warpSpeed:   1,
			initialFuel: 1,
			want:        want{position: Vector{1, 1}, fuelUsed: -1, orbitingPlanetNum: None},
		},
		{
			name:        "overshoot waypoint at warp 5",
			design:      testLongRangeScoutDesign(1),
			destination: Vector{5, 5},
			warpSpeed:   5,
			initialFuel: 300,
			want:        want{position: Vector{5, 5}, fuelUsed: 1, orbitingPlanetNum: None},
		},
		{
			name:        "end up at planet",
			design:      testLongRangeScoutDesign(1),
			destination: Vector{5, 5},
			warpSpeed:   5,
			initialFuel: 300,
			planet:      NewPlanet().WithNum(1).withPosition(Vector{5, 5}),
			want:        want{position: Vector{5, 5}, fuelUsed: 1, orbitingPlanetNum: 1},
		},
		{
			name:        "W10 run out of fuel",
			design:      testFuelNubianDesign(1, &Interspace10),
			destination: Vector{0, 100},
			warpSpeed:   10,
			initialFuel: 50,
			planet:      NewPlanet().WithNum(1).withPosition(Vector{0, 100}),
			// 50 LY at W10, then 0.5 LY at W1
			// Generates 1.5 --> 1 mg
			want: want{position: Vector{0, 51}, fuelUsed: 49, orbitingPlanetNum: None},
		},
		{
			name: "Refueler out of fuel",
			design: NewShipDesign(1, 1).WithHull(SuperFuelXport.Name).WithSlots([]ShipDesignSlot{
				{HullComponent: TransGalacticDrive.Name, HullSlotIndex: 1, Quantity: 2},
			}),
			destination: Vector{0, 81},
			warpSpeed:   9,
			initialFuel: 0,
			planet:      NewPlanet().WithNum(1).withPosition(Vector{0, 81}),
			want:        want{position: Vector{0, 81}, fuelUsed: -134, orbitingPlanetNum: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
			tt.design = tt.design.WithSpec(&rules, player)
			player.Designs = append(player.Designs, tt.design)
			f := newFleetForDesign(player, tt.design, 1, 1, tt.name, []Waypoint{
				NewPositionWaypoint(Vector{0, 0}, 0),
				NewPositionWaypoint(tt.destination, tt.warpSpeed),
			})
			fleet := &f
			fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			fleet.Fuel = tt.initialFuel

			universe := Universe{Fleets: []*Fleet{fleet}}
			if tt.planet != nil {
				universe.Planets = []*Planet{tt.planet}
			}
			universe.buildMaps([]*Player{player})

			fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, fleet.Position)
			assert.Equal(t, tt.want.position, fleet.Waypoints[0].Position)
			assert.Equal(t, tt.want.fuelUsed, tt.initialFuel-fleet.Fuel)
		})
	}
}

func TestFleet_moveFleetEngineFailure(t *testing.T) {
	tests := []struct {
		name      string
		warpSpeed int
		player    *Player
		random    rng
		wantMoved bool
	}{
		{
			name:      "under failure speed",
			player:    NewPlayer(1, NewRace().WithLRT(CE).WithSpec(&rules)),
			warpSpeed: 6,
			random:    newFloat64Random(0),
			wantMoved: true,
		},
		{
			name:      "high roll; still moves",
			player:    NewPlayer(1, NewRace().WithLRT(CE).WithSpec(&rules)),
			warpSpeed: 7,
			random:    newFloat64Random(0.11), // engine failure occurs 10% of the time, < 10/100
			wantMoved: true,
		},
		{
			name:      "no CE; still moves",
			warpSpeed: 7,
			player:    NewPlayer(1, NewRace().WithSpec(&rules)),
			random:    newFloat64Random(0.1),
			wantMoved: true,
		},
		{
			name:      "low roll; failure",
			warpSpeed: 9,
			player:    NewPlayer(1, NewRace().WithLRT(CE).WithSpec(&rules)),
			random:    newFloat64Random(0.1), // engine failure occurs at 10/100
			wantMoved: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := tt.player
			fleet := testLongRangeScout(player).withWaypoints(
				NewPositionWaypoint(Vector{0, 0}, 0),
				NewPositionWaypoint(Vector{999, 0}, tt.warpSpeed),
			)
			universe := Universe{Fleets: []*Fleet{fleet}}
			universe.buildMaps([]*Player{player})

			rules := NewRules()
			rules.random = tt.random

			fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			if (fleet.Position != Vector{0, 0}) != tt.wantMoved {
				if tt.wantMoved {
					t.Errorf("Fleet.moveFleetEngineFailure() did not move fleet when expected")
				} else {
					t.Errorf("Fleet.moveFleetEngineFailure() moved fleet to position %+v; expected no movement", fleet.Position)
				}
			}
		})
	}
}

func TestFleet_gateFleet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}}
	sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(1)
	sourcePlanet.Spec = PlanetSpec{
		PlanetStarbaseSpec: PlanetStarbaseSpec{
			HasStargate:  true,
			SafeRange:    100,
			SafeHullMass: 100,
			MaxRange:     500,
			MaxHullMass:  500,
		},
	}
	destPlanet := NewPlanet().WithNum(2).WithPlayerNum(1)
	destPlanet.Spec = PlanetSpec{
		PlanetStarbaseSpec: PlanetStarbaseSpec{
			HasStargate:  true,
			SafeRange:    100,
			SafeHullMass: 100,
			MaxRange:     500,
			MaxHullMass:  500,
		},
	}

	type want struct {
		position          Vector
		orbitingPlanetNum int
	}
	tests := []struct {
		name        string
		fleet       *Fleet
		want        want
		wantMessage bool
	}{
		{
			name: "gate between planets",
			fleet: testLongRangeScout(player).
				withOrbitingPlanetNum(sourcePlanet.Num).
				withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			want: want{position: Vector{50, 0}, orbitingPlanetNum: destPlanet.Num},
		},
		{
			name: "gate fail, no source",
			fleet: testLongRangeScout(player).
				withPosition(Vector{200, 0}).
				withWaypoints(NewPositionWaypoint(Vector{200, 0}, 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			want:        want{position: Vector{200, 0}},
			wantMessage: true,
		},
		{
			name: "gate fail, no dest",
			fleet: testLongRangeScout(player).
				withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPositionWaypoint(Vector{200, 0}, StargateWarpSpeed)),
			want:        want{position: Vector{0, 0}},
			wantMessage: true,
		},
		{
			name: "jump gate",
			fleet: testGatePrivateer(player, 1).
				withPosition(Vector{200, 0}).
				withCargo(Cargo{10, 10, 10, 10}). // jumpgates allow cargo, sweet!
				withWaypoints(NewPositionWaypoint(Vector{200, 0}, 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			want: want{position: Vector{50, 0}, orbitingPlanetNum: destPlanet.Num},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, token := range tt.fleet.Tokens {
				player.Designs = append(player.Designs, token.design)
			}
			universe := Universe{
				Fleets:       []*Fleet{tt.fleet},
				Planets:      []*Planet{sourcePlanet, destPlanet},
				designsByNum: map[playerObject]*ShipDesign{},
			}
			universe.buildMaps([]*Player{player})

			tt.fleet.gateFleet(&rules, &universe, newTestPlayerGetter(player))

			if tt.fleet.Position != tt.want.position {
				t.Errorf("Fleet.gateFleet() position = %v, want %v", tt.fleet.Position, tt.want.position)
			}
			if tt.fleet.OrbitingPlanetNum != tt.want.orbitingPlanetNum {
				t.Errorf("Fleet.gateFleet() OrbitingPlanetNum = %v, want %v", tt.fleet.OrbitingPlanetNum, tt.want.orbitingPlanetNum)
			}
			if tt.wantMessage && len(player.Messages) == 0 {
				t.Errorf("Fleet.gateFleet() wantMessages, got none")
			}

		})
	}
}

func TestFleet_repairFleet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	type args struct {
		prt    PRT
		fleet  *Fleet
		planet *Planet
	}
	tests := []struct {
		name string
		args args
		want []ShipToken
	}{
		{"no damage", args{JoaT, testLongRangeScout(player), nil}, []ShipToken{{QuantityDamaged: 0, Damage: 0}}},
		{"repair min 1dp", args{JoaT,
			&Fleet{
				MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
				BaseName:  "Long Range Scout",
				Tokens: []ShipToken{
					{
						Quantity:        1,
						QuantityDamaged: 1,
						Damage:          10,
						DesignNum:       1,
						design: NewShipDesign(player.Num, 1).
							WithHull(Scout.Name).
							WithSlots([]ShipDesignSlot{
								{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							}).
							WithSpec(&rules, player)},
				},
				OrbitingPlanetNum: None,
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{
						NewPositionWaypoint(Vector{}, 5),
					},
				},
			},
			nil,
		},
			// should repair 2% (min 1dp)
			[]ShipToken{{QuantityDamaged: 1, Damage: 9}},
		},
		{"repair 5% when orbiting our planet", args{JoaT,
			&Fleet{
				MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
				BaseName:  "100dp Fleet",
				Tokens: []ShipToken{
					{
						Quantity:        3,
						QuantityDamaged: 2,
						Damage:          10,
						DesignNum:       1,
						design: NewShipDesign(player.Num, 1).
							WithHull(MidgetMiner.Name). // has 100dp armor
							WithSlots([]ShipDesignSlot{
								{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							}).
							WithSpec(&rules, player)},
				},
				OrbitingPlanetNum: 1,
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{
						NewPositionWaypoint(Vector{}, 5),
					},
				},
			},
			NewPlanet().WithNum(1).WithPlayerNum(player.Num),
		},
			// should repair 5% of 100, or 5 dp (on both damaged tokens)
			[]ShipToken{{QuantityDamaged: 2, Damage: 5}},
		},
		{"IS repair double (10%) when orbiting our planet", args{IS,
			&Fleet{
				MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
				BaseName:  "100dp Fleet",
				Tokens: []ShipToken{
					{
						Quantity:        1,
						QuantityDamaged: 1,
						Damage:          20,
						DesignNum:       1,
						design: NewShipDesign(player.Num, 1).
							WithHull(MidgetMiner.Name). // has 100dp armor
							WithSlots([]ShipDesignSlot{
								{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							}).
							WithSpec(&rules, player)},
				},
				OrbitingPlanetNum: 1,
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{
						NewPositionWaypoint(Vector{}, 5),
					},
				},
			},
			NewPlanet().WithNum(1).WithPlayerNum(player.Num),
		},
			// should repair 5% of 100, or 5 dp (on both damaged tokens)
			[]ShipToken{{QuantityDamaged: 1, Damage: 10}},
		},
		{"repair fully", args{JoaT,
			&Fleet{
				MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
				BaseName:  "100dp Fleet",
				Tokens: []ShipToken{
					{
						Quantity:        3,
						QuantityDamaged: 2,
						Damage:          5,
						DesignNum:       1,
						design: NewShipDesign(player.Num, 1).
							WithHull(MidgetMiner.Name). // has 100dp armor
							WithSlots([]ShipDesignSlot{
								{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							}).
							WithSpec(&rules, player)},
				},
				OrbitingPlanetNum: 1,
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{
						NewPositionWaypoint(Vector{}, 5),
					},
				},
			},
			NewPlanet().WithNum(1).WithPlayerNum(player.Num),
		},
			// should repair 5% of 100, or 5 dp (on both damaged tokens)
			[]ShipToken{{QuantityDamaged: 0, Damage: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := *player
			p.Race.PRT = tt.args.prt
			p.Race.Spec = computeRaceSpec(&p.Race, &rules)

			tt.args.fleet.Spec = ComputeFleetSpec(&rules, player, tt.args.fleet)

			// if a planet is passed in, orbit it
			if tt.args.planet != nil {
				tt.args.fleet.OrbitingPlanetNum = tt.args.planet.Num
			}

			tt.args.fleet.repairFleet(testLogger, &rules, &p, tt.args.planet)

			for i, token := range tt.args.fleet.Tokens {
				if token.Damage != tt.want[i].Damage {
					t.Errorf("Fleet.repairFleet() token %d gotDamage = %v, wantDamage %v", i, token.Damage, tt.want[i].Damage)
				}
				if token.QuantityDamaged != tt.want[i].QuantityDamaged {
					t.Errorf("Fleet.repairFleet() token %d gotQuantityDamaged = %v, wantQuantityDamaged %v", i, token.QuantityDamaged, tt.want[i].QuantityDamaged)
				}
			}
		})
	}
}

func TestFleet_repairStarbase(t *testing.T) {
	type args struct {
		prt    PRT
		armor  int
		damage float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"no damage", args{JoaT, 100, 0}, 0},
		{"20 damage, repair 10", args{JoaT, 100, 20}, 10},
		{"20 damage, PRT IS repair 15", args{IS, 100, 20}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithPRT(tt.args.prt).WithSpec(&rules)).withSpec(&rules)
			starbase := newStarbase(player, NewPlanet(), NewShipDesign(player.Num, 1).WithHull(SpaceStation.Name).WithSpec(&rules, player), "Starbase")
			starbase.Tokens[0].QuantityDamaged = 1
			starbase.Tokens[0].Damage = tt.args.damage
			starbase.Tokens[0].design.Spec.Armor = tt.args.armor

			starbase.repairStarbase(testLogger, &rules, player)

			if starbase.Tokens[0].Damage != tt.want {
				t.Errorf("Fleet.repairStarbase() got = %v, want %v", starbase.Tokens[0].Damage, tt.want)
			}
		})
	}
}

func TestFleet_getEstimatedRange(t *testing.T) {
	type testfleet struct {
		design *ShipDesign
		qty    int
	}
	tests := []struct {
		name      string
		fleets    []testfleet
		race      *Race
		cargoMass int
		fuel      int
		warpSpeed int
		want      int
	}{
		// examples are (mostly) cross checked against ones from base game
		{
			name: "long range scout",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1),
				qty:    1,
			}},
			race:      NewRace().WithSpec(&rules),
			warpSpeed: 5,
			want:      2400,
		},
		{
			name: "W6 scout with IFE",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1),
				qty:    1,
			}},
			race:      NewRace().WithLRT(IFE).WithSpec(&rules),
			warpSpeed: 6,
			want:      2654, // 2673 in base game
		},

		{
			name: "scout with less fuel",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1),
				qty:    1,
			}},
			race:      NewRace().WithSpec(&rules),
			warpSpeed: 5,
			fuel:      100,
			want:      800, // 1/3 the fuel = 1/3 the range
		},
		{
			name: "multiple copies, same range",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1),
				qty:    2,
			}},
			race:      NewRace().WithSpec(&rules),
			fuel:      600,
			warpSpeed: 5,
			want:      2400,
		},
		{
			name: "infinite range mizer scout",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Scout.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				}),
				qty: 1,
			}},
			race:      NewRace().WithLRT(IFE).WithSpec(&rules),
			fuel:      300,
			warpSpeed: 4,
			want:      Infinite,
		},
		{
			name: "IT swashbuckler with full cargo",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Privateer.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: DaddyLongLegs7.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 2},
					{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
					{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 1},
					{HullComponent: AlphaTorpedo.Name, HullSlotIndex: 5, Quantity: 1},
				}),
				qty: 1,
			}},
			race:      NewRace().WithPRT(IT).WithLRT(IFE).WithSpec(&rules),
			warpSpeed: 7,
			fuel:      650,
			cargoMass: 250,
			want:      295,
		},
		{
			name: "Galleon with minicols, no fuel",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Galleon.Name).WithSlots([]ShipDesignSlot{ // 125kT
					{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 4},       // 36kT
					{HullComponent: Tritanium.Name, HullSlotIndex: 2, Quantity: 2},       // 120kT
					{HullComponent: Tritanium.Name, HullSlotIndex: 3, Quantity: 2},       // 120kT
					{HullComponent: Tritanium.Name, HullSlotIndex: 4, Quantity: 3},       // 180kT
					{HullComponent: Tritanium.Name, HullSlotIndex: 5, Quantity: 3},       // 180kT
					{HullComponent: FuelTank.Name, HullSlotIndex: 6, Quantity: 2},        // 6kT
					{HullComponent: MineDispenser50.Name, HullSlotIndex: 7, Quantity: 1}, // 30kT
					{HullComponent: PossumScanner.Name, HullSlotIndex: 8, Quantity: 1},   // 3kT
				}),
				qty: 200,
				// 800kT * 1.05 = -840 mg/ly at warp 6
			}, {
				design: NewShipDesign(1, 2).WithHull(MiniColonyShip.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: SettlersDelight.Name, HullSlotIndex: 1, Quantity: 1},
				}),
				qty: 840,
				// +840 mg/ly at warp 6
			}},
			race:      NewRace().WithPRT(HE).WithSpec(&rules),
			fuel:      1,
			warpSpeed: 6,
			want:      Infinite, // minicol fuel production exactly offsets galleon fuel consumption (net zero)
		},
		{
			name: "200kT nubian with QJ5",
			fleets: []testfleet{{
				design: testFuelNubianDesign(1, &QuickJump5),
				qty:    1,
			}},
			race:      NewRace().WithSpec(&rules),
			warpSpeed: 5,
			fuel:      5000,
			want:      5000, // 1 mg/ly
		},
		{
			name: "200kT nubian with IS-10",
			fleets: []testfleet{{
				design: testFuelNubianDesign(1, &Interspace10),
				qty:    1,
			}},
			race:      NewRace().WithSpec(&rules),
			warpSpeed: 10,
			fuel:      5000,
			want:      5000, // 1 mg/ly
		},
		{
			name: "Net Positive Refueler",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(SuperFuelXport.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 2},
				}),
				qty: 1,
			}},
			race:      NewRace().WithSpec(&rules),
			warpSpeed: 9,
			fuel:      0,
			want:      Infinite, // Mizer consumes 180 mg/yr; refueler makes 200
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.race).withSpec(&rules)
			fleet := NewFleet(player, 1, tt.name, []Waypoint{{}})
			for _, f := range tt.fleets {
				fleet.Tokens = append(fleet.Tokens, ShipToken{
					Quantity: f.qty,
					design:   f.design.WithSpec(&rules, player),
				})
			}

			fleet.Cargo = Cargo{Ironium: tt.cargoMass}
			fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			if tt.fuel == 0 {
				fleet.Fuel = fleet.Spec.FuelCapacity
			} else {
				fleet.Fuel = max(tt.fuel, 0)
			}

			if got := fleet.getEstimatedRange(player, tt.warpSpeed); got != tt.want {
				t.Errorf("Fleet.getEstimatedRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_getFuelGeneration(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	type testfleet struct {
		design *ShipDesign
		qty    int
	}
	tests := []struct {
		name      string
		fleets    []testfleet
		warpSpeed int
		distance  float64
		want      int
	}{
		{
			name: "normal engine, no fuel generation",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1).WithSpec(&rules, player),
				qty:    1,
			}},
			warpSpeed: 6,
			distance:  36,
			want:      0,
		},
		{
			name: "warp1, 1mg fuel",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1).WithSpec(&rules, player),
				qty:    1,
			}},
			warpSpeed: 1,
			distance:  1,
			want:      1,
		},
		{
			name: "Mizer warp 4",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Scout.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				}).WithSpec(&rules, player),
				qty: 1,
			}},
			warpSpeed: 4,
			distance:  16,
			want:      16,
		},
		{
			name: "Mizer warp 3",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Scout.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				}).WithSpec(&rules, player),
				qty: 1,
			}},
			warpSpeed: 3,
			distance:  9,
			want:      27,
		},
		{
			name: "double ships, double fuel",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(Scout.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
					{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
					{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
				}).WithSpec(&rules, player),
				qty: 2,
			}},
			warpSpeed: 4,
			distance:  16,
			want:      32,
		},
		{
			name: "Super Fuel Xport + Ramscoop",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(SuperFuelXport.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 2},
				}).WithSpec(&rules, player),
				qty: 1,
			}},
			warpSpeed: 5,
			distance:  25,
			want:      350, // 25*3*2 + 200
		},
		{
			name: "Distance doesn't affect flat gen",
			fleets: []testfleet{{
				design: NewShipDesign(1, 1).WithHull(SuperFuelXport.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: TransGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 2},
				}).WithSpec(&rules, player),
				qty: 1,
			}},
			warpSpeed: 6,
			distance:  0,
			want:      200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet := NewFleet(player, 1, tt.name, []Waypoint{{}})
			for _, f := range tt.fleets {
				fleet.Tokens = append(fleet.Tokens, ShipToken{
					Quantity: f.qty,
					design:   f.design,
				})
			}

			fleet.Spec = ComputeFleetSpec(&rules, player, fleet)

			if got := fleet.getFuelGeneration(tt.warpSpeed, tt.distance); got != tt.want {
				t.Errorf("Fleet.getFuelGeneration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_reduceCargoToMax(t *testing.T) {
	type fields struct {
		cargo         Cargo
		cargoCapacity int
	}
	tests := []struct {
		name           string
		fields         fields
		wantCargo      Cargo
		wantJettisoned Cargo
	}{
		{
			name: "empty",
			fields: fields{
				cargo:         Cargo{},
				cargoCapacity: 0,
			},
			wantCargo: Cargo{}, wantJettisoned: Cargo{},
		},
		{
			name: "no reduce",
			fields: fields{
				cargo:         Cargo{Ironium: 10},
				cargoCapacity: 10,
			},
			wantCargo: Cargo{Ironium: 10}, wantJettisoned: Cargo{},
		},
		{
			name: "no room left",
			fields: fields{
				cargo:         Cargo{Ironium: 10},
				cargoCapacity: 0,
			},
			wantCargo: Cargo{}, wantJettisoned: Cargo{Ironium: 10},
		},
		{
			name: "save the people!",
			fields: fields{
				cargo:         Cargo{Ironium: 10, Colonists: 10},
				cargoCapacity: 5,
			},
			wantCargo: Cargo{Colonists: 5}, wantJettisoned: Cargo{Ironium: 10, Colonists: 5},
		},
		{
			name: "half of each",
			fields: fields{
				cargo:         Cargo{10, 8, 6, 0},
				cargoCapacity: 12,
			},
			wantCargo: Cargo{5, 4, 3, 0}, wantJettisoned: Cargo{5, 4, 3, 0},
		},
		{
			name: "low mins",
			fields: fields{
				cargo:         Cargo{1, 1, 1, 2},
				cargoCapacity: 4,
			},
			wantCargo: Cargo{1, 1, 0, 2}, wantJettisoned: Cargo{0, 0, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet := &Fleet{
				BaseName: tt.name,
				Cargo:    tt.fields.cargo,
				Spec: FleetSpec{
					ShipDesignSpec: ShipDesignSpec{
						CargoCapacity: tt.fields.cargoCapacity,
					},
				},
			}

			gotJettisoned := fleet.reduceCargoToMax()
			if gotCargo := fleet.Cargo; gotCargo != tt.wantCargo {
				t.Errorf("Fleet.reduceCargoToMax() returned fleet cargo %v, wantCargo %v", gotCargo, tt.wantCargo)
			}
			if gotJettisoned != tt.wantJettisoned {
				t.Errorf("Fleet.reduceCargoToMax() returned jettisonned cargo %v, wantJettisoned %v", gotJettisoned, tt.wantJettisoned)
			}
		})
	}
}
