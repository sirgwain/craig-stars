package cs

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

// create a new long range scout fleet for testing with pre-computed specs.
func testLongRangeScout(player *Player) *Fleet {
	return testLongRangeScoutWithQuantity(player, 1)
}

// create a new long range scout ship design for testing.
// Does NOT come with precomputed specs; those will have to be done manually.
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
		BaseName: "Stealing Freighter",
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
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	type args struct {
		player *Player
		planet *Planet
	}
	type want struct {
		position          Vector
		fuelUsed          int
		orbitingPlanetNum int
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  want
	}{
		{
			"move 25ly at warp5",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 5)),
			args{player: player},
			want{Vector{25, 0}, 4, None},
		},
		{
			"move 1ly at warp 1",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{1, 1}, 1)),
			args{player: player},
			want{Vector{1, 1}, 0, None},
		},
		{
			"overshoot waypoint at warp 5",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{5, 5}, 5)),
			args{player: player},
			want{Vector{5, 5}, 1, None},
		},
		{
			"end up at planet",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{5, 5}, 5)),
			args{player: player, planet: NewPlanet().WithNum(1).withPosition(Vector{5, 5})},
			want{Vector{5, 5}, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := tt.args.player
			universe := Universe{Fleets: []*Fleet{tt.fleet}}
			if tt.args.planet != nil {
				universe.Planets = []*Planet{tt.args.planet}
			}
			if err := universe.buildMaps([]*Player{player}); err != nil {
				t.Fatal(err)
			}

			tt.fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, tt.fleet.Position)
			assert.Equal(t, tt.want.position, tt.fleet.Waypoints[0].Position)
			assert.Equal(t, tt.want.fuelUsed, tt.fleet.Spec.FuelCapacity-tt.fleet.Fuel)
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
			if err := universe.buildMaps([]*Player{player}); err != nil {
				t.Fatal(err)
			}

			rules := NewRules()
			rules.random = tt.random

			fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			if (fleet.Position != Vector{}) != tt.wantMoved {
				if tt.wantMoved {
					t.Error("Fleet.moveFleet() did not move fleet when expected")
				} else {
					t.Errorf("Fleet.moveFleet() moved fleet to position %+v; expected engine failure", fleet.Position)
				}
			}
		})
	}
}

func TestFleet_gateFleet(t *testing.T) {
	// TODO: Add IT test cases
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
	destPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2)
	// friends with ourselves and the target planet's owner
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}
	destPlayer.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}

	var middleOfNowhere = Vector{200, 0}

	// make 2 planets with 100/250 gates
	sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(player.Num).WithName("Source Planet")
	sourcePlanet.Starbase = newStarbase(player, sourcePlanet, NewShipDesign(player.Num, 1).
		WithHull(SpaceStation.Name).
		WithSlots([]ShipDesignSlot{{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1}}).
		WithSpec(&rules, player), "Source Base").withSpec(&rules, player)
	destPlanet := NewPlanet().WithNum(2).WithPlayerNum(destPlayer.Num).withPosition(Vector{50, 0}).WithName("Dest Planet")
	destPlanet.Starbase = newStarbase(destPlayer, destPlanet, NewShipDesign(destPlayer.Num, 1).
		WithHull(SpaceStation.Name).
		WithSlots([]ShipDesignSlot{{HullComponent: Stargate100_250.Name, HullSlotIndex: 1, Quantity: 1}}).
		WithSpec(&rules, destPlayer), "Dest Base").withSpec(&rules, player)
	sourcePlanet.Spec = computePlanetSpec(&rules, player, sourcePlanet)
	destPlanet.Spec = computePlanetSpec(&rules, player, destPlanet)

	heavyNubianDesign := NewShipDesign(player.Num, 2).WithHull(Nubian.Name).WithSlots([]ShipDesignSlot{
		{HullComponent: Interspace10.Name, HullSlotIndex: 1, Quantity: 3},    // 175kT with hull
		{HullComponent: JihadMissile.Name, HullSlotIndex: 2, Quantity: 3},    // 305kT
		{HullComponent: JihadMissile.Name, HullSlotIndex: 3, Quantity: 3},    // 410kT
		{HullComponent: MineDispenser50.Name, HullSlotIndex: 4, Quantity: 3}, // 495kT
	}).WithSpec(&rules, player)

	type want struct {
		position Vector
		exploded int
		message  bool
		cargo    Cargo
	}
	tests := []struct {
		name      string
		fleet     *Fleet
		waypoints []Waypoint
		rng       rng
		want      want
	}{
		{
			name:  "no source planet",
			fleet: testLongRangeScout(player).withPosition(middleOfNowhere),
			waypoints: []Waypoint{
				NewPositionWaypoint(middleOfNowhere, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			want: want{position: middleOfNowhere, message: true},
		},
		{
			name:  "no dest planet",
			fleet: testLongRangeScout(player),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPositionWaypoint(middleOfNowhere, StargateWarpSpeed),
			},
			want: want{position: sourcePlanet.Position, message: true},
		},
		{
			name: "ship explodes",
			fleet: newFleetForTokens(player, 1, "this gets overridden", []ShipToken{
				{design: heavyNubianDesign, Quantity: 1, Damage: 1000, QuantityDamaged: 1}}, nil),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			rng:  newFloat64Random(1), // doesn't vanish, but explodes regardless
			want: want{position: sourcePlanet.Position, exploded: 1, message: true},
		},
		{
			name: "can't dump colonists on others' worlds",
			fleet: testSmallFreighter(player).
				withCargo(Cargo{100, 100, 100, 100}),
			// start at other player's starbase and trying to return home
			waypoints: []Waypoint{
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, 5),
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, StargateWarpSpeed),
			},
			want: want{position: destPlanet.Position, cargo: Cargo{100, 100, 100, 100}, message: true},
		},
		{
			name: "can't gate; ships too heavy",
			fleet: newFleetForTokens(player, 69420, "this gets overridden", []ShipToken{
				{design: NewShipDesign(player.Num, 2).WithHull(Nubian.Name).WithSlots([]ShipDesignSlot{
					{HullComponent: Interspace10.Name, HullSlotIndex: 1, Quantity: 3}, // 175kT with hull
					{HullComponent: Carbonic.Name, HullSlotIndex: 3, Quantity: 3},     // 275kT
					{HullComponent: Tritanium.Name, HullSlotIndex: 3, Quantity: 3},    // 455kT
					{HullComponent: JihadMissile.Name, HullSlotIndex: 4, Quantity: 3}, // 560kT
				}).WithSpec(&rules, player), Quantity: 10},
				{design: testLongRangeScoutDesign(player.Num).WithSpec(&rules, player), Quantity: 10}}, nil),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			want: want{position: sourcePlanet.Position, message: true},
		},
		{
			name:  "gate succeeds",
			fleet: testLongRangeScout(player),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			want: want{position: destPlanet.Position},
		},
		{
			name: "gate succeeds, dump cargo",
			fleet: testSmallFreighter(player).
				withOrbitingPlanetNum(sourcePlanet.Num).
				withCargo(Cargo{100, 100, 100, 20}),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			want: want{position: destPlanet.Position, cargo: Cargo{}, message: true},
		},
		{
			name: "mixed fleet gates with losses",
			fleet: newFleetForTokens(player, 69420, "this gets overridden", []ShipToken{
				{design: heavyNubianDesign, Quantity: 10},
				{design: testLongRangeScoutDesign(player.Num).WithSpec(&rules, player), Quantity: 10}}, nil),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			// 1/3 mass vanish chance on nubians means the first 4 rolls pass.
			// Rest of them default to 0, but the scouts won't vanish anyways
			rng:  newFloat64Random(0.30, 0.31, 0.32, 0.33, 0.34, 0.35, 0.36, 0.37, 0.38, 0.39, 0.4),
			want: want{position: destPlanet.Position, exploded: 6, message: true},
		},
		{
			name: "vanish prioritizes damaged tokens",
			fleet: newFleetForTokens(player, 1, "this gets overridden", []ShipToken{
				// Exactly enough damage to kill an overgated nubian (2% of 5000)
				{design: heavyNubianDesign, Quantity: 6, Damage: 100, QuantityDamaged: 56}}, nil),
			waypoints: []Waypoint{
				NewPlanetWaypoint(sourcePlanet.Position, sourcePlanet.Num, sourcePlanet.Name, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			// first damaged ship vanishes, leaving only 1 ship surviving
			rng:  newFloat64Random(0, 1, 1, 1, 1, 1),
			want: want{position: destPlanet.Position, exploded: 5, message: true},
		},
		{
			name: "jump gate from middle of space",
			fleet: testGatePrivateer(player, 1).
				withPosition(middleOfNowhere).
				withCargo(Cargo{10, 10, 10, 10}),
			waypoints: []Waypoint{
				NewPositionWaypoint(middleOfNowhere, 5),
				NewPlanetWaypoint(destPlanet.Position, destPlanet.Num, destPlanet.Name, StargateWarpSpeed),
			},
			rng:  newFloat64Random(1), // make sure we get a highroll and don't explode
			want: want{position: destPlanet.Position, cargo: Cargo{10, 10, 10, 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rCopy := NewRules()
			if tt.rng != nil {
				rCopy.random = tt.rng
			}

			// reset player designs & messages from prior tests
			player.Messages = []PlayerMessage{}
			player.Designs = []*ShipDesign{}
			initQty := 0 // used for token headcount later
			for i, token := range tt.fleet.Tokens {
				tt.fleet.Tokens[i].DesignNum = token.design.Num // needed to prevent crash inside buildMaps
				player.Designs = append(player.Designs, token.design)
				initQty += token.Quantity
			}

			tt.fleet.Waypoints = tt.waypoints
			tt.fleet.OrbitingPlanetNum = tt.waypoints[0].TargetNum
			tt.fleet.Name = tt.name // for debugging test cases
			tt.fleet.Spec = ComputeFleetSpec(&rCopy, player, tt.fleet)
			universe := Universe{
				Fleets:  []*Fleet{tt.fleet},
				Planets: []*Planet{sourcePlanet, destPlanet},
			}
			if err := universe.buildMaps([]*Player{player, destPlayer}); err != nil {
				t.Fatal(err)
			}

			tt.fleet.gateFleet(&rCopy, &universe, newTestPlayerGetter(player, destPlayer))

			if tt.fleet.Position != tt.want.position {
				t.Errorf("Fleet.gateFleet() position = %v, want %v", tt.fleet.Position, tt.want.position)
			}
			if tt.fleet.Cargo != tt.want.cargo {
				t.Errorf("Fleet.gateFleet() produced fleet cargo \n%+v, want \n%+v", tt.fleet.Cargo, tt.want.cargo)
			}
			if exploded := initQty - tt.fleet.Spec.TotalShips; exploded != tt.want.exploded {
				t.Errorf("Fleet.gateFleet() destroyed %d tokens, want %d", exploded, tt.want.exploded)
			}
			if tt.want.message != (len(player.Messages) > 0) {
				if tt.want.message {
					t.Errorf("Fleet.gateFleet() expected to produce messages, got none")
				} else {
					messages, err := json.MarshalIndent(player.Messages, "", "\t")
					if err != nil {
						t.Fatalf("Error marshaling unexpected messages to JSON: \n%v", err)
					}
					t.Errorf("Fleet.gateFleet() produced messages unexpectedly; messages: \n%s", string(messages))
				}
			}

		})
	}
}

func TestFleet_repairFleet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	var midgetMiner = NewShipDesign(player.Num, 1).
		WithHull(MidgetMiner.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		})

	type args struct {
		prt    PRT
		tokens []ShipToken
		planet *Planet
	}
	type want struct {
		Damage          float64
		QuantityDamaged int
	}
	tests := []struct {
		name   string
		args   args
		moving bool // TODO: Add moving fleet test cases
		want   []want
	}{
		// TODO: Add super fuel xport tests
		{
			name: "no damage",
			args: args{
				prt: JoaT,
				tokens: []ShipToken{{
					Quantity:        1,
					QuantityDamaged: 0,
					Damage:          0,
					design:          testLongRangeScoutDesign(1),
				}},
				planet: nil,
			},
			want: []want{{
				QuantityDamaged: 0, // these are zero values but just for clarity
				Damage:          0,
			}},
		},
		{
			name: "repair min 1dp",
			args: args{
				prt: JoaT,
				tokens: []ShipToken{{
					Quantity:        1,
					QuantityDamaged: 1,
					Damage:          10,
					design:          testLongRangeScoutDesign(1),
				}},
				planet: nil,
			},
			want: []want{{
				QuantityDamaged: 1,
				Damage:          9,
			}},
		},
		{
			name: "repair 5% when orbiting our planet",
			args: args{
				prt: JoaT,
				tokens: []ShipToken{{
					Quantity:        3,
					QuantityDamaged: 2,
					Damage:          10,
					design:          midgetMiner,
				}},
				planet: NewPlanet().WithNum(1).WithPlayerNum(player.Num),
			},
			want: []want{{
				QuantityDamaged: 2,
				Damage:          5,
			}},
		},
		{
			name: "IS repair double",
			args: args{
				prt: IS,
				tokens: []ShipToken{{
					Quantity:        1,
					QuantityDamaged: 1,
					Damage:          20,
					design:          midgetMiner,
				}},
				planet: NewPlanet().WithNum(1).WithPlayerNum(player.Num),
			},
			want: []want{{
				QuantityDamaged: 1,
				Damage:          10,
			}},
		},
		{
			name: "repair fully",
			args: args{
				prt: JoaT,
				tokens: []ShipToken{{
					Quantity:        3,
					QuantityDamaged: 2,
					Damage:          5,
					design:          midgetMiner,
				}},
				planet: NewPlanet().WithNum(1).WithPlayerNum(player.Num),
			},
			want: []want{{
				QuantityDamaged: 0,
				Damage:          0,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.tokens) != len(tt.want) {
				t.Fatalf("Incorrect test setup for test %q; tt.want and tt.args.tokens must be same length", tt.name)
			}

			p := *player
			p.Race.PRT = tt.args.prt
			p.Race.Spec = computeRaceSpec(&p.Race, &rules)

			for i := range tt.args.tokens {
				tt.args.tokens[i].design = tt.args.tokens[i].design.WithSpec(&rules, player)
			}
			fleet := newFleetForTokens(&p, 1, tt.name, tt.args.tokens, []Waypoint{{Position: Vector{}}})
			if tt.moving {
				// add extra waypoint if we're supposed to move
				fleet.Waypoints = append(fleet.Waypoints, Waypoint{Position: Vector{1, 1}})
			}

			fleet.Spec = ComputeFleetSpec(&rules, &p, fleet)

			fleet.repairFleet(log.Logger, &rules, &p, tt.args.planet)

			for i, token := range fleet.Tokens {
				want := tt.want[i]
				if token.Damage != want.Damage {
					t.Errorf("Fleet.repairFleet() gotDamage for token #%d = %v, wantDamage %v", i, token.Damage, want.Damage)
				}
				if token.QuantityDamaged != want.QuantityDamaged {
					t.Errorf("Fleet.repairFleet() gotQuantityDamaged for token #%d = %v, wantQuantityDamaged %v", i, token.QuantityDamaged, want.QuantityDamaged)
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

			starbase.repairStarbase(log.Logger, &rules, player)

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
		// examples were all cross checked against ones from base game
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
			name: "multiple copies, same range",
			fleets: []testfleet{{
				design: testLongRangeScoutDesign(1),
				qty:    2,
			}},
			race:      NewRace().WithSpec(&rules),
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
			cargoMass: 250,
			want:      295,
		},
		{
			name: "mixed privateer fleet",
			fleets: []testfleet{
				{
					design: NewShipDesign(1, 1).WithHull(Privateer.Name).WithSlots([]ShipDesignSlot{
						{HullComponent: DaddyLongLegs7.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 2},
						{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: AlphaTorpedo.Name, HullSlotIndex: 5, Quantity: 1},
					}),
					qty: 1,
				},
				{
					design: NewShipDesign(1, 1).WithHull(Privateer.Name).WithSlots([]ShipDesignSlot{
						{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 5, Quantity: 1},
					}),
					qty: 1,
				},
			},
			race:      NewRace().WithPRT(IT).WithLRT(IFE).WithSpec(&rules),
			warpSpeed: 4,
			want:      3134,
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
				fleet.Fuel = tt.fuel
			}

			if got := fleet.getEstimatedRange(player, tt.warpSpeed, fleet.Spec.CargoCapacity); got != tt.want {
				t.Errorf("Fleet.getEstimatedRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_getFuelGeneration(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	fuelMizerScout := testLongRangeScout(player)
	fuelMizerScout.Tokens[0].design.Slots[0].HullComponent = FuelMizer.Name
	fuelMizerScout.Tokens[0].design.Spec, _ = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, fuelMizerScout.Tokens[0].design)
	fuelMizerScout.Spec = ComputeFleetSpec(&rules, player, fuelMizerScout)
	fuelMizerScoutX2 := testLongRangeScout(player)
	fuelMizerScoutX2.Tokens[0].design.Slots[0].HullComponent = FuelMizer.Name
	fuelMizerScoutX2.Tokens[0].design.Spec, _ = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, fuelMizerScoutX2.Tokens[0].design)
	fuelMizerScoutX2.Tokens[0].Quantity = 2
	fuelMizerScoutX2.Spec = ComputeFleetSpec(&rules, player, fuelMizerScoutX2)

	type args struct {
		warpSpeed int
		distance  float64
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  int
	}{
		{"normal warp, no fuel generation", testLongRangeScout(player), args{6, 36}, 0},
		{"warp1, 1mg fuel", testLongRangeScout(player), args{1, 1}, 1},
		{"fuel mizer, 16mg fuel", fuelMizerScout, args{4, 16}, 16},
		{"fuel mizer, warp3 27mg fuel", fuelMizerScout, args{3, 9}, 27},
		{"fuel mizerx2, double fuel", fuelMizerScoutX2, args{4, 16}, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fleet.getFuelGeneration(tt.args.warpSpeed, tt.args.distance); got != tt.want {
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
		{name: "empty", fields: fields{cargo: Cargo{}, cargoCapacity: 0}, wantCargo: Cargo{}, wantJettisoned: Cargo{}},
		{name: "no reduce", fields: fields{cargo: Cargo{Ironium: 10}, cargoCapacity: 10}, wantCargo: Cargo{Ironium: 10}, wantJettisoned: Cargo{}},
		{name: "no room left", fields: fields{cargo: Cargo{Ironium: 10}, cargoCapacity: 0}, wantCargo: Cargo{}, wantJettisoned: Cargo{Ironium: 10}},
		{name: "save the people!", fields: fields{cargo: Cargo{Ironium: 10, Colonists: 10}, cargoCapacity: 5}, wantCargo: Cargo{Colonists: 5}, wantJettisoned: Cargo{Ironium: 10, Colonists: 5}},
		{name: "save 50% of each", fields: fields{cargo: Cargo{10, 8, 6, 0}, cargoCapacity: 12}, wantCargo: Cargo{5, 4, 3, 0}, wantJettisoned: Cargo{5, 4, 3, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet := &Fleet{
				Cargo: tt.fields.cargo,
				Spec: FleetSpec{
					ShipDesignSpec: ShipDesignSpec{
						CargoCapacity: tt.fields.cargoCapacity,
					},
				},
			}

			if got := fleet.reduceCargoToMax(); got != tt.wantJettisoned {
				t.Errorf("Fleet.reduceCargoToMax() = %v, wantJettisoned %v", got, tt.wantJettisoned)
			}
			if got := fleet.Cargo; got != tt.wantCargo {
				t.Errorf("Fleet.reduceCargoToMax() cargo = %v, wantCargo %v", got, tt.wantCargo)
			}
		})
	}
}
