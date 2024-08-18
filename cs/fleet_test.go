package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

// create a new long rang scout fleet for testing
func testLongRangeScout(player *Player) *Fleet {
	return testLongRangeScoutWithQuantity(player, 1)
}

func testLongRangeScoutWithQuantity(player *Player, quantity int) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Long Range Scout",
		Tokens: []ShipToken{
			{
				Quantity:  quantity,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
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
				design: NewShipDesign(player, 1).
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
	return fleet
}

// create a new small freighter (with cargo pod) fleet for testing
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
				design: NewShipDesign(player, 1).
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

// create a new small freighter (with cargo pod) fleet for testing
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
				design: NewShipDesign(player, 1).
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
		BaseName: "Cloaked Scout",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
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
				design: NewShipDesign(player, 1).
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
				design: NewShipDesign(player, 1).
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
				Cost:           Cost{16, 2, 7, 19},
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
				Cost:           Cost{122, 263, 236, 752},
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
				Cost:           Cost{12, 2, 9, 20}.MultiplyInt(2),
				FuelCapacity:   50 * 2,
				ReduceCloaking: 1,
				ScanRange:      66,
				ScanRangePen:   30,
				Scanner:        true,
				Mass:           19 * 2,
				Armor:          20 * 2,
				CloakUnits:     70,
				CloakPercent:   35, // still 35%
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        19 * 2,
			BaseCloakedCargo: 0,
			TotalShips:       2,
		}},
		{"0 Scouts", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Long Range Scout",
			Tokens: []ShipToken{
				{
					DesignNum: 1,
					Quantity:  0,
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
						WithHull(MiniBomber.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: LadyFingerBomb.Name, HullSlotIndex: 2, Quantity: 2},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{22, 45, 10, 43},
				FuelCapacity:   120,
				Mass:           112,
				Armor:          50,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				Bomber:         true,
				Bombs: []Bomb{
					{Quantity: 2, KillRate: .6, MinKillRate: 300, StructureDestroyRate: .2},
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
					design: NewShipDesign(starterHumanoidPlayer, 1).
						WithHull(MiniBomber.Name).
						WithSlots([]ShipDesignSlot{
							{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
							{HullComponent: LadyFingerBomb.Name, HullSlotIndex: 2, Quantity: 2},
						}).
						WithSpec(&rules, starterHumanoidPlayer)},
			},
		}}, FleetSpec{
			ShipDesignSpec: ShipDesignSpec{
				Cost:           Cost{22, 45, 10, 43}.MultiplyInt(2),
				FuelCapacity:   120 * 2,
				Mass:           112 * 2,
				Armor:          50 * 2,
				ReduceCloaking: 1,
				Scanner:        true,
				ScanRangePen:   NoScanner,
				Bomber:         true,
				Bombs: []Bomb{
					{Quantity: 4, KillRate: .6, MinKillRate: 300, StructureDestroyRate: .2},
				},
				Engine: Engine{
					IdealSpeed:   QuickJump5.IdealSpeed,
					FreeSpeed:    QuickJump5.FreeSpeed,
					MaxSafeSpeed: QuickJump5.MaxSafeSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        112 * 2,
			BaseCloakedCargo: 112 * 2,
			TotalShips:       2,
		}},
		{"2 B52 Bombers with multiple bomb types", args{&rules, starterHumanoidPlayer, &Fleet{
			BaseName: "Bomber",
			Tokens: []ShipToken{
				{
					Quantity:  2,
					DesignNum: 1,
					design: NewShipDesign(starterHumanoidPlayer, 1).
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
					{Quantity: 8, KillRate: 2.5, MinKillRate: 300, StructureDestroyRate: 1},
					{Quantity: 8, KillRate: .3, StructureDestroyRate: 2.8},
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
			if got := ComputeFleetSpec(tt.args.rules, tt.args.player, tt.args.fleet); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("ComputeFleetSpec() = \n%v, want \n%v", got, tt.want)
			}
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
			universe.buildMaps([]*Player{player})

			tt.fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, tt.fleet.Position)
			assert.Equal(t, tt.want.position, tt.fleet.Waypoints[0].Position)
			assert.Equal(t, tt.want.fuelUsed, tt.fleet.Spec.FuelCapacity-tt.fleet.Fuel)
		})
	}
}

func TestFleet_moveFleetEngineFailure(t *testing.T) {

	player := NewPlayer(1, NewRace().WithLRT(CE).WithSpec(&rules))
	playerWithoutCE := NewPlayer(1, NewRace().WithSpec(&rules))

	type args struct {
		player *Player
		random rng
	}
	type want struct {
		position Vector
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  want
	}{
		{
			"move without engine failure",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 6)),
			args{player, newFloat64Random(0)},
			want{Vector{36, 0}},
		},
		{
			"move without engine failure high speed",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 7)),
			args{player, newFloat64Random(.2)}, // engine failure occurs 10% of the time, < 10/100
			want{Vector{49, 0}},
		},
		{
			"move with engine failure",
			testLongRangeScout(player).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 7)),
			args{player, newFloat64Random(.1)}, // engine failure at 10/100
			want{Vector{0, 0}},
		},
		{
			"move without engine failure, no CE",
			testLongRangeScout(playerWithoutCE).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 7)),
			args{playerWithoutCE, newFloat64Random(.1)},
			want{Vector{49, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := tt.args.player
			universe := Universe{Fleets: []*Fleet{tt.fleet}}
			universe.buildMaps([]*Player{player})

			rules := NewRules()
			rules.random = tt.args.random

			tt.fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, tt.fleet.Position)
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

	type args struct {
		player  *Player
		players []*Player
		planets []*Planet
	}
	type want struct {
		position          Vector
		orbitingPlanetNum int
	}
	tests := []struct {
		name        string
		fleet       *Fleet
		args        args
		want        want
		wantMessage bool
	}{
		{
			name: "gate between planets",
			fleet: testLongRangeScout(player).
				withOrbitingPlanetNum(sourcePlanet.Num).
				withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			args: args{player: player, players: []*Player{player}, planets: []*Planet{sourcePlanet, destPlanet}},
			want: want{position: Vector{50, 0}, orbitingPlanetNum: destPlanet.Num},
		},
		{
			name: "gate fail, no source",
			fleet: testLongRangeScout(player).
				withPosition(Vector{200, 0}).
				withWaypoints(NewPositionWaypoint(Vector{200, 0}, 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			args:        args{player: player, players: []*Player{player}, planets: []*Planet{sourcePlanet, destPlanet}},
			want:        want{position: Vector{200, 0}},
			wantMessage: true,
		},
		{
			name: "gate fail, no dest",
			fleet: testLongRangeScout(player).
				withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPositionWaypoint(Vector{200, 0}, StargateWarpSpeed)),
			args:        args{player: player, players: []*Player{player}, planets: []*Planet{sourcePlanet, destPlanet}},
			want:        want{position: Vector{0, 0}},
			wantMessage: true,
		},
		{
			name: "jump gate",
			fleet: testGatePrivateer(player, 1).
				withPosition(Vector{200, 0}).
				withCargo(Cargo{10, 10, 10, 10}). // jumpgates allow cargo, sweet!
				withWaypoints(NewPositionWaypoint(Vector{200, 0}, 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)),
			args: args{player: player, players: []*Player{player}, planets: []*Planet{sourcePlanet, destPlanet}},
			want: want{position: Vector{50, 0}, orbitingPlanetNum: destPlanet.Num},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, token := range tt.fleet.Tokens {
				player.Designs = append(player.Designs, token.design)
			}
			universe := Universe{Fleets: []*Fleet{tt.fleet}, Planets: tt.args.planets, designsByNum: map[playerObject]*ShipDesign{}}
			universe.buildMaps(tt.args.players)

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

func TestFleet_getCargoLoadAmount(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	planet := NewPlanet().WithCargo(Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 1000})

	type args struct {
		dest      cargoHolder
		cargoType CargoType
		task      WaypointTransportTask
	}
	tests := []struct {
		name               string
		fleet              *Fleet
		args               args
		wantTransferAmount int
		wantWaitAtWaypoint bool
	}{
		{
			name:               "load 1kt ironium",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAmount, Amount: 1}},
			wantTransferAmount: 1,
		},
		{
			name:               "load 1mg fuel",
			fleet:              testLongRangeScout(player).withFuel(0),
			args:               args{dest: testLongRangeScout(player).withFuel(10), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAmount, Amount: 1}},
			wantTransferAmount: 1,
		},
		{
			name:               "load all ironium we can fit",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 120, // small freighter has 120kT cargo capacity
		},
		{
			name:               "load all fuel we can fit",
			fleet:              testSmallFreighter(player).withFuel(0),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 130, // small freighter has 130mg fuel capacity
		},
		{
			name:               "load all ironium we can fit (we already loaded 20)",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Boranium: 20}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 100,
		},
		{
			name:               "load all fuel we can fit (we already loaded 20)",
			fleet:              testSmallFreighter(player).withFuel(20),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 110,
		},
		{
			name:               "load fill percent",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionFillPercent, Amount: 50}},
			wantTransferAmount: 60, // 50% of 120kT capacity
		},
		{
			name:               "load fill percent fuel",
			fleet:              testSmallFreighter(player).withFuel(0),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionFillPercent, Amount: 50}},
			wantTransferAmount: 65, // 50% of 130mg capacity
		},
		{
			name:               "load fill percent but wait",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Ironium: 50}), cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionWaitForPercent, Amount: 50}},
			wantTransferAmount: 50, // load all 50, wait for the additional 10
			wantWaitAtWaypoint: true,
		},
		{
			name:               "set amount to 20kT when we have 10kT already",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 10, // load 10kT more
		},
		{
			name:               "set fuel amount to 20mg when we have 10mg already",
			fleet:              testSmallFreighter(player).withFuel(10),
			args:               args{dest: testSmallFreighter(player), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 10, // load 10mg more
		},
		{
			name:               "set amount to 20kT when we have 10kT already, but planet only has 5k",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Ironium: 5}), cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 5,    // load 5kT more
			wantWaitAtWaypoint: true, // wait for remaining 5kT we want
		},
		{
			name:               "set amount to 20mg fuel when we have 10mg already, but dest fleet only has 5mg",
			fleet:              testSmallFreighter(player).withFuel(10),
			args:               args{dest: testSmallFreighter(player).withFuel(5), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 5,    // load 5mg more
			wantWaitAtWaypoint: true, // wait for remaining 5mg we want
		},
		{
			name:               "set amount to 20kT when we have 30kT already. We should do nothing as we already have > 20kT",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 30}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 0, // don't unload
		},
		{
			name:               "set amount to 20mg fuel when we have 30mg already. We should do nothing, we already have > 20mg",
			fleet:              testSmallFreighter(player).withFuel(30),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 0, // don't unload
		},
		{
			name:               "set waypoint to 2500kT colonists",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2600}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 100, // load 100kT
		},
		{
			name:               "set waypoint to 2500kT colonists, nothing to load",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 0, // no load, colonists too low
		},
		{
			name:               "set waypoint to 2500kT colonists, fleet has colonists, don't unload during fleetLoad",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 200}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 0, // no load, colonists too low
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotTransferAmount, gotWaitAtWaypoint := tt.fleet.getCargoLoadAmount(tt.args.dest, tt.args.cargoType, tt.args.task)
			if gotTransferAmount != tt.wantTransferAmount {
				t.Errorf("Fleet.getCargoLoadAmount() gotTransferAmount = %v, want %v", gotTransferAmount, tt.wantTransferAmount)
			}
			if gotWaitAtWaypoint != tt.wantWaitAtWaypoint {
				t.Errorf("Fleet.getCargoLoadAmount() gotWaitAtWaypoint = %v, want %v", gotWaitAtWaypoint, tt.wantWaitAtWaypoint)
			}
		})
	}
}

func TestFleet_getCargoUnloadAmount(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	planet := NewPlanet().WithCargo(Cargo{Ironium: 1000, Boranium: 1000, Germanium: 1000, Colonists: 1000})

	type args struct {
		dest      cargoHolder
		cargoType CargoType
		task      WaypointTransportTask
	}
	tests := []struct {
		name               string
		fleet              *Fleet
		args               args
		wantTransferAmount int
		wantWaitAtWaypoint bool
	}{
		{
			name:               "unload 1kt ironium",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 1}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAmount, Amount: 1}},
			wantTransferAmount: 1,
		},
		{
			name:               "unload 1mg fuel",
			fleet:              testLongRangeScout(player).withFuel(10),
			args:               args{dest: testLongRangeScout(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAmount, Amount: 1}},
			wantTransferAmount: 1,
		},
		{
			name:               "unload all ironium we have",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 120,
		},
		{
			name:               "unload all fuel the dest can fit",
			fleet:              testSmallFreighter(player),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 130, // small freighter has 130mg fuel capacity
		},
		{
			name:               "unload all ironium we have",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 20}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 20,
		},
		{
			name:               "unload all fuel dest can fit (they already have 20)",
			fleet:              testSmallFreighter(player),
			args:               args{dest: testSmallFreighter(player).withFuel(20), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 110,
		},
		{
			name:               "unload all fuel at planet, does nothing",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionUnloadAll}},
			wantTransferAmount: 0,
		},
		{
			name:               "set amount to 20kT when we have 30kT already",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 30}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 10, // unload 10kT onto planet
		},
		{
			name:               "set fuel amount to 20mg when we have 30mg already",
			fleet:              testSmallFreighter(player).withFuel(30),
			args:               args{dest: testSmallFreighter(player).withFuel(0), cargoType: Fuel, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 10, // load 10mg more
		},
		{
			name:               "set amount to 20kT when we have 10kT, should unload nothing",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 0,
		},
		{
			name:               "set waypoint to 2500kT colonists, should unload 100kT",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 100}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 100, // unload 100kT
		},
		{
			name:               "set waypoint to 2500kT colonists, nothing to unload",
			fleet:              testSmallFreighter(player),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2400}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 0, // no unload
		},
		{
			name:               "set waypoint to 2500kT colonists, fleet has colonists, planet has enough, don't unload",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Colonists: 200}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Colonists: 2500}), cargoType: Colonists, task: WaypointTransportTask{Action: TransportActionSetWaypointTo, Amount: 2500}},
			wantTransferAmount: 0, // no unload, planet has enough
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotTransferAmount, gotWaitAtWaypoint := tt.fleet.getCargoUnloadAmount(tt.args.dest, tt.args.cargoType, tt.args.task)
			if gotTransferAmount != tt.wantTransferAmount {
				t.Errorf("Fleet.getCargoLoadAmount() gotTransferAmount = %v, want %v", gotTransferAmount, tt.wantTransferAmount)
			}
			if gotWaitAtWaypoint != tt.wantWaitAtWaypoint {
				t.Errorf("Fleet.getCargoLoadAmount() gotWaitAtWaypoint = %v, want %v", gotWaitAtWaypoint, tt.wantWaitAtWaypoint)
			}
		})
	}
}

func TestFleet_transferToDest(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules))

	type args struct {
		dest           cargoHolder
		cargoType      CargoType
		transferAmount int
	}
	tests := []struct {
		name           string
		fleet          *Fleet
		args           args
		wantFleetCargo Cargo
		wantDestCargo  Cargo
		wantErr        bool
	}{
		{
			name:           "transfer 10kT to planet",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: 10},
			wantFleetCargo: Cargo{},
			wantDestCargo:  Cargo{Ironium: 1010},
		},
		{
			name:           "transfer 10kT to another fleet",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:           args{dest: testSmallFreighter(player).withCargo(Cargo{Ironium: 100}), cargoType: Ironium, transferAmount: 10},
			wantFleetCargo: Cargo{Ironium: 110},
			wantDestCargo:  Cargo{Ironium: 110},
		},

		{
			name:           "transfer 10kT from planet",
			fleet:          testSmallFreighter(player),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: -10},
			wantFleetCargo: Cargo{Ironium: 10},
			wantDestCargo:  Cargo{Ironium: 990},
		},
		{
			name:          "transfer 1000kT from planet, error",
			fleet:         testSmallFreighter(player),
			args:          args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: -1000},
			wantDestCargo: Cargo{Ironium: 1000},
			wantErr:       true,
		},
		{
			name:           "transfer 1000kT to planet, error",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:           args{dest: NewPlanet().WithCargo(Cargo{Ironium: 1000}), cargoType: Ironium, transferAmount: 1000},
			wantFleetCargo: Cargo{Ironium: 10},
			wantDestCargo:  Cargo{Ironium: 1000},
			wantErr:        true,
		},
		{
			name:           "transfer 120kT to another fleet with cargo, error",
			fleet:          testSmallFreighter(player).withCargo(Cargo{Ironium: 120}),
			args:           args{dest: testSmallFreighter(player).withCargo(Cargo{Ironium: 100}), cargoType: Ironium, transferAmount: 120},
			wantFleetCargo: Cargo{Ironium: 120},
			wantDestCargo:  Cargo{Ironium: 100},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.fleet.transferToDest(tt.args.dest, tt.args.cargoType, tt.args.transferAmount); (err != nil) != tt.wantErr {
				t.Errorf("Fleet.transferToDest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if *tt.args.dest.getCargo() != tt.wantDestCargo {
				t.Errorf("Fleet.transferToDest() destCargo = %v, wantDestCargo %v", *tt.args.dest.getCargo(), tt.wantDestCargo)
			}

			if tt.fleet.Cargo != tt.wantFleetCargo {
				t.Errorf("Fleet.transferToDest() fleet.Cargo = %v, wantFleetCargo %v", tt.fleet.Cargo, tt.wantFleetCargo)
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
						design: NewShipDesign(player, 1).
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
						design: NewShipDesign(player, 1).
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
						design: NewShipDesign(player, 1).
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
						design: NewShipDesign(player, 1).
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

			tt.args.fleet.repairFleet(&rules, &p, tt.args.planet)

			for i := range tt.args.fleet.Tokens {
				token := tt.args.fleet.Tokens[i]
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
			starbase := newStarbase(player, NewPlanet(), NewShipDesign(player, 1).WithHull(SpaceStation.Name).WithSpec(&rules, player), "Starbase")
			starbase.Tokens[0].QuantityDamaged = 1
			starbase.Tokens[0].Damage = tt.args.damage
			starbase.Tokens[0].design.Spec.Armor = tt.args.armor

			starbase.repairStarbase(&rules, player)

			if starbase.Tokens[0].Damage != tt.want {
				t.Errorf("Fleet.repairStarbase() got = %v, want %v", starbase.Tokens[0].Damage, tt.want)
			}
		})
	}
}

func TestFleet_getEstimatedRange(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	type args struct {
		fleet     *Fleet
		player    *Player
		warpSpeed int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"long range scout range", args{testLongRangeScout(player), player, 5}, 2400},
		{"mini miner", args{testMiniMineLayer(player), player, 5}, 689},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.fleet.getEstimatedRange(tt.args.player, tt.args.warpSpeed, tt.args.fleet.Spec.CargoCapacity); got != tt.want {
				t.Errorf("Fleet.getEstimatedRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_getFuelGeneration(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	fuelMizerScout := testLongRangeScout(player)
	fuelMizerScout.Tokens[0].design.Slots[0].HullComponent = FuelMizer.Name
	fuelMizerScout.Tokens[0].design.Spec = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, fuelMizerScout.Tokens[0].design)
	fuelMizerScout.Spec = ComputeFleetSpec(&rules, player, fuelMizerScout)
	fuelMizerScoutX2 := testLongRangeScout(player)
	fuelMizerScoutX2.Tokens[0].design.Slots[0].HullComponent = FuelMizer.Name
	fuelMizerScoutX2.Tokens[0].design.Spec = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, fuelMizerScoutX2.Tokens[0].design)
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
