//go:build !wasi && !wasm

package cs

import (
	"slices"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
					WithSlots(slices.Clone(DesignSmallFreighter.Slots)).
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
					WithSlots(slices.Clone(DesignGalleon.Slots)).
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
					WithSlots(slices.Clone(DesignCloakedScout.Slots)).
					WithSpec(&rules, player)},
		},
		OrbitingPlanetNum: None,
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
					WithSlots(slices.Clone(DesignGatePrivateer.Slots)).
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

func testPotatoBug(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Potato Bug",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Potato Bug").
					WithHull(MidgetMiner.Name).
					WithSlots(slices.Clone(DesignPotatoBug.Slots)).
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

func testSantaMaria(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Santa Maria",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Santa Maria").
					WithHull(ColonyShip.Name).
					WithSlots(slices.Clone(DesignSantaMaria.Slots)).
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

func testSantaMariaIFE(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Santa Maria",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithName("Santa Maria").
					WithHull(ColonyShip.Name).
					WithSlots(slices.Clone(DesignSantaMariaIFE.Slots)).
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

func Test_computeFleetSpec(t *testing.T) {
	starterHumanoidPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithTechLevels(TechLevel{3, 3, 3, 3, 3, 3})
	starterHumanoidPlayer.Race.Spec = ComputeRaceSpec(&starterHumanoidPlayer.Race, &rules)

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
			universe := Universe{log: testLogger, Fleets: []*Fleet{tt.fleet}}
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
			universe := Universe{log: testLogger, Fleets: []*Fleet{fleet}}
			universe.buildMaps([]*Player{player})

			rules := NewRulesWithSeed(0)
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

func TestFleet_advanceWaypoints(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	planet := NewPlanet().WithNum(1).withPosition(Vector{10, 0})
	planet.Name = "Planet 1"
	a := NewPositionWaypoint(Vector{0, 0}, 5)
	b := NewPlanetWaypoint(Vector{10, 0}, 1, "Planet 1", 5)
	c := NewPositionWaypoint(Vector{20, 0}, 5)
	intercept := NewFleetWaypoint(Vector{10, 0}, 2, 2, "Enemy", 5)
	intercept.PartiallyComplete = true
	transportFleet := NewFleetWaypoint(Vector{10, 0}, 2, 1, "Freighter", 5)
	transportFleet.Task = WaypointTaskTransport
	scoutFleet := NewFleetWaypoint(Vector{10, 0}, 2, 1, "Scout", 5)

	type want struct {
		positions []Vector
		wp0Target MapObjectType
	}
	tests := []struct {
		name      string
		waypoints []Waypoint
		repeat    bool
		want      want
	}{
		{"no repeat", []Waypoint{a, b, c}, false, want{[]Vector{{10, 0}, {20, 0}}, MapObjectTypePlanet}},
		{"repeat sends reached waypoint to the end", []Waypoint{a, b, c}, true, want{[]Vector{{10, 0}, {20, 0}, {10, 0}}, MapObjectTypePlanet}},
		{"repeat doesn't repeat the last waypoint", []Waypoint{a, b}, true, want{[]Vector{{10, 0}}, MapObjectTypePlanet}},
		{"repeat doesn't repeat if the last waypoint is already there", []Waypoint{a, b, c, b}, true, want{[]Vector{{10, 0}, {20, 0}, {10, 0}}, MapObjectTypePlanet}},
		{"repeat doesn't repeat an intercept", []Waypoint{a, intercept, a}, true, want{[]Vector{{10, 0}, {0, 0}}, MapObjectTypePlanet}},
		{"reaching a fleet makes wp0 where we are", []Waypoint{a, scoutFleet}, false, want{[]Vector{{10, 0}}, MapObjectTypePlanet}},
		{"reaching a fleet to transport keeps the fleet", []Waypoint{a, transportFleet}, false, want{[]Vector{{10, 0}}, MapObjectTypeFleet}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet := testLongRangeScout(player).withWaypoints(tt.waypoints...)
			fleet.RepeatOrders = tt.repeat
			fleet.Position = Vector{10, 0}
			fleet.OrbitingPlanetNum = planet.Num
			universe := Universe{log: testLogger, Fleets: []*Fleet{fleet}, Planets: []*Planet{planet}}
			universe.buildMaps([]*Player{player})

			fleet.advanceWaypoints(&universe, tt.waypoints[1])

			positions := []Vector{}
			for _, wp := range fleet.Waypoints {
				positions = append(positions, wp.Position)
			}
			assert.Equal(t, tt.want.positions, positions)
			assert.Equal(t, tt.want.wp0Target, fleet.Waypoints[0].TargetType)
		})
	}
}

func TestFleet_moveStepFuel(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	fleet := testLongRangeScout(player).
		withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{200, 0}, 9)).
		withFuel(300)
	universe := Universe{log: testLogger, Fleets: []*Fleet{fleet}}
	universe.buildMaps([]*Player{player})
	move := newFleetMove(fleet)

	// a pursuer moves 1/5 of its move each pass, 17+17+17+17+13 = 81ly
	for _, step := range []float64{17, 17, 17, 17, 13} {
		fleet.moveStep(&rules, &universe, newTestPlayerGetter(player), move, step)
	}

	// charged for 81ly at once, 91.125mg rounded up, not 20+20+20+20+15mg for each step
	assert.Equal(t, Vector{81, 0}, fleet.Position)
	assert.Equal(t, 92, move.fuelUsed)
	assert.Equal(t, 300-92, fleet.Fuel)
}

func TestFleet_moveFleetEngineStrain(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// a Trans-Star 10 freighter is safe at warp 10
	transStarFreighter := NewShipDesign(player.Num, 2).
		WithName("Trans-Star Freighter").
		WithHull(SmallFreighter.Name).
		WithSlots(slices.Clone(DesignSmallFreighter.Slots))
	transStarFreighter.Slots[0].HullComponent = TransStar10.Name
	transStarFreighter.WithSpec(&rules, player)
	player.Designs = append(player.Designs, transStarFreighter)

	type want struct {
		ships          int
		cargo          Cargo
		destroyed      bool
		messageType    PlayerMessageType
		explodedAmount int
	}
	tests := []struct {
		name       string
		transStars int // Trans-Star 10 freighters added to the Quick Jump 5 freighter
		cargo      Cargo
		random     rng
		want       want
	}{
		{
			name:   "safe roll",
			random: newFloat64Random(0.5),
			want:   want{ships: 1},
		},
		{
			name:   "fleet destroyed",
			random: newFloat64Random(0),
			want:   want{destroyed: true, messageType: PlayerMessageFleetEngineStrainDestroyed},
		},
		{
			// the Quick Jump 5 freighter has half the fleet's capacity and takes half the cargo with it
			name:       "lost ship takes its share of cargo",
			transStars: 1,
			cargo:      Cargo{Ironium: 100, Colonists: 20},
			random:     newFloat64Random(0),
			want:       want{ships: 1, cargo: Cargo{Ironium: 50, Colonists: 10}, messageType: PlayerMessageFleetExceededSafeSpeed, explodedAmount: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.Messages = nil
			fleet := testSmallFreighter(player)
			if tt.transStars > 0 {
				fleet.Tokens = append(fleet.Tokens, ShipToken{Quantity: tt.transStars, DesignNum: 2, design: transStarFreighter})
			}
			fleet.withCargo(tt.cargo).withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{999, 0}, 10))
			fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			fleet.Fuel = 10000 // plenty of fuel for warp 10
			universe := Universe{log: testLogger, Fleets: []*Fleet{fleet}}
			universe.buildMaps([]*Player{player})

			rules := NewRulesWithSeed(0)
			rules.random = tt.random

			interrupted := fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			if tt.want.destroyed {
				assert.Equal(t, fleetMoveInterruptedDestroyed, int(interrupted.reason))
				assert.Empty(t, fleet.Tokens)
				assert.Equal(t, Vector{0, 0}, fleet.Position)
			} else {
				assert.Equal(t, tt.want.ships, fleet.Spec.TotalShips)
				assert.Equal(t, tt.want.cargo, fleet.Cargo)
				assert.Equal(t, Vector{100, 0}, fleet.Position)
			}

			if tt.want.messageType == PlayerMessageNone {
				assert.Empty(t, player.Messages)
			} else {
				assert.Len(t, player.Messages, 1)
				assert.Equal(t, tt.want.messageType, player.Messages[0].Type)
				assert.Equal(t, tt.want.explodedAmount, player.Messages[0].Spec.Amount)
			}
		})
	}
}

func TestFleet_gateFleetChooseGate(t *testing.T) {
	gate := func(safeRange int) PlanetSpec {
		return PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{HasStargate: true, SafeRange: safeRange, SafeHullMass: 500, MaxRange: safeRange * 5, MaxHullMass: 2500}}
	}
	tests := []struct {
		name         string
		sourceRange  int
		destRange    int
		jumpgate     bool
		dest         Waypoint
		wantPosition Vector
		wantCargo    Cargo // cargo left on the fleet
		wantDamaged  bool
	}{
		{
			name:        "gate to a planet by position",
			sourceRange: 500, destRange: 500,
			dest:         NewPositionWaypoint(Vector{100, 0}, StargateWarpSpeed),
			wantPosition: Vector{100, 0},
		},
		{
			name:        "gate to a fleet at a planet",
			sourceRange: 500, destRange: 500,
			dest:         NewFleetWaypoint(Vector{100, 0}, 5, 1, "Other Fleet", StargateWarpSpeed),
			wantPosition: Vector{100, 0},
		},
		{
			name:        "jumpgate is better than the planet's gate, keep cargo",
			sourceRange: 50, destRange: 500, jumpgate: true,
			dest:         NewPlanetWaypoint(Vector{100, 0}, 2, "planet 2", StargateWarpSpeed),
			wantPosition: Vector{100, 0},
			wantCargo:    Cargo{Ironium: 10},
		},
		{
			name:        "planet's gate is better than the jumpgate, dump cargo",
			sourceRange: 500, destRange: 50, jumpgate: true,
			dest:         NewPlanetWaypoint(Vector{100, 0}, 2, "planet 2", StargateWarpSpeed),
			wantPosition: Vector{100, 0},
		},
		{
			name:        "planet's gate overgates without a jumpgate",
			sourceRange: 50, destRange: 500,
			dest:         NewPlanetWaypoint(Vector{100, 0}, 2, "planet 2", StargateWarpSpeed),
			wantPosition: Vector{100, 0},
			wantDamaged:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
			player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}}
			sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(1)
			sourcePlanet.Spec = gate(tt.sourceRange)
			destPlanet := NewPlanet().WithNum(2).WithPlayerNum(1).withPosition(Vector{100, 0})
			destPlanet.Spec = gate(tt.destRange)

			fleet := testLongRangeScout(player)
			if tt.jumpgate {
				fleet = testGatePrivateer(player, 1)
			}
			fleet.withOrbitingPlanetNum(sourcePlanet.Num).
				withCargo(Cargo{Ironium: 10}).
				withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), tt.dest)
			for _, token := range fleet.Tokens {
				player.Designs = append(player.Designs, token.design)
			}
			universe := Universe{
				log:          testLogger,
				Fleets:       []*Fleet{fleet},
				Planets:      []*Planet{sourcePlanet, destPlanet},
				designsByNum: map[playerObject]*ShipDesign{},
			}
			universe.buildMaps([]*Player{player})

			fleet.gateFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.wantPosition, fleet.Position)
			assert.Equal(t, destPlanet.Num, fleet.OrbitingPlanetNum)
			assert.Equal(t, tt.wantCargo, fleet.Cargo)
			assert.Equal(t, tt.wantDamaged, fleet.Tokens[0].Damage > 0)
		})
	}
}

func TestFleet_gateFleetDumpCargoOnAlly(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
	ally := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2)
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}
	ally.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationFriend}}

	stargate := PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{HasStargate: true, SafeRange: 150, SafeHullMass: 150, MaxRange: 500, MaxHullMass: 500}}
	sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(ally.Num)
	sourcePlanet.Spec = stargate
	destPlanet := NewPlanet().WithNum(2).WithPlayerNum(player.Num)
	destPlanet.Spec = stargate

	fleet := testSmallFreighter(player).
		withCargo(Cargo{Ironium: 10}).
		withOrbitingPlanetNum(sourcePlanet.Num).
		withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed))
	universe := Universe{
		log:          testLogger,
		Fleets:       []*Fleet{fleet},
		Planets:      []*Planet{sourcePlanet, destPlanet},
		designsByNum: map[playerObject]*ShipDesign{},
	}
	universe.buildMaps([]*Player{player, ally})

	fleet.gateFleet(&rules, &universe, newTestPlayerGetter(player, ally))

	assert.Equal(t, Vector{50, 0}, fleet.Position)
	assert.Equal(t, Cargo{Ironium: 10}, sourcePlanet.Cargo)
	assert.Len(t, player.Messages, 1)
	assert.Equal(t, PlayerMessageFleetDumpedCargo, player.Messages[0].Type)
	assert.Len(t, ally.Messages, 1)
	assert.Equal(t, PlayerMessageFleetDumpedCargo, ally.Messages[0].Type)
}

func TestFleet_gateFleet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}}
	sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(1)
	sourcePlanet.Spec = PlanetSpec{
		PlanetStarbaseSpec: PlanetStarbaseSpec{
			HasStargate:  true,
			SafeRange:    150,
			SafeHullMass: 150,
			MaxRange:     500,
			MaxHullMass:  500,
		},
	}
	destPlanet := NewPlanet().WithNum(2).WithPlayerNum(1)
	destPlanet.Spec = PlanetSpec{
		PlanetStarbaseSpec: PlanetStarbaseSpec{
			HasStargate:  true,
			SafeRange:    150,
			SafeHullMass: 150,
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
				log:          testLogger,
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
			p.Race.Spec = ComputeRaceSpec(&p.Race, &rules)

			tt.args.fleet.Spec = ComputeFleetSpec(&rules, player, tt.args.fleet)

			// if a planet is passed in, orbit it
			if tt.args.planet != nil {
				tt.args.fleet.OrbitingPlanetNum = tt.args.planet.Num
			}

			tt.args.fleet.repairFleet(testLogger, &rules, &p, tt.args.planet)

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
			want:      2666, // 2673 in base game
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

			if got := fleet.getEstimatedRange(player, tt.warpSpeed); got != tt.want {
				t.Errorf("Fleet.getEstimatedRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_GetFuelCost(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// two scouts in separate tokens
	scoutTokens := testLongRangeScout(player)
	scoutTokens.Tokens = append(scoutTokens.Tokens, scoutTokens.Tokens[0])
	scoutTokens.Spec = ComputeFleetSpec(&rules, player, scoutTokens)

	// a Quick Jump 5 freighter and a Long Hump 6 freighter
	mixedFreighters := testSmallFreighter(player)
	longHumpFreighter := NewShipDesign(player.Num, 2).
		WithName("Long Hump Freighter").
		WithHull(SmallFreighter.Name).
		WithSlots(slices.Clone(DesignSmallFreighter.Slots))
	longHumpFreighter.Slots[0].HullComponent = LongHump6.Name
	longHumpFreighter.WithSpec(&rules, player)
	mixedFreighters.Tokens = append(mixedFreighters.Tokens, ShipToken{Quantity: 1, DesignNum: 2, design: longHumpFreighter})
	mixedFreighters.Cargo = Cargo{Ironium: 120}
	mixedFreighters.Spec = ComputeFleetSpec(&rules, player, mixedFreighters)

	tests := []struct {
		name      string
		fleet     *Fleet
		warpSpeed int
		distance  float64
		want      int
	}{
		// 25kT * 100 * 25ly / 20000 = 3.125mg
		{"scout", testLongRangeScout(player), 5, 25, 4},
		// 2 * 3.125mg = 6.25mg, rounded once for the fleet, not per token
		{"two scout tokens round once", scoutTokens, 5, 25, 7},
		{"two scouts in one token", testLongRangeScoutWithQuantity(player, 2), 5, 25, 7},
		// all 120kT of cargo goes on the Long Hump 6 freighter, it uses less fuel at warp 6
		// (41kT + 120kT) * 105 * 100ly / 20000 + 36kT * 180 * 100ly / 20000 = 116.925mg
		{"cargo on most efficient ships", mixedFreighters, 6, 100, 117},
		{"stargate is free", testLongRangeScout(player), StargateWarpSpeed, 100, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fleet.GetFuelCost(player, tt.warpSpeed, tt.distance))
		})
	}
}

func TestFleet_moveFleetFuel(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	type want struct {
		position  Vector
		fuel      int
		warpSpeed int // wp1 warp speed after the move
	}
	tests := []struct {
		name      string
		engine    string // replaces the scout's Long Hump 6
		dest      Vector
		warpSpeed int
		fuel      int
		want      want
	}{
		{
			// warp 9 for 81ly costs 92mg (91.125mg), 46mg gets us 46/91.125 of the way
			name:      "out of fuel stops where the fuel runs out",
			dest:      Vector{200, 0},
			warpSpeed: 9,
			fuel:      46,
			want:      want{position: Vector{41, 0}, fuel: 0, warpSpeed: 1},
		},
		{
			name:      "no fuel doesn't move",
			dest:      Vector{200, 0},
			warpSpeed: 9,
			fuel:      0,
			want:      want{position: Vector{0, 0}, fuel: 0, warpSpeed: 1},
		},
		{
			// the fuel mizer travels free at warp 4, but makes no fuel on a turn it ran out
			name:      "ramscoop makes no fuel when out of fuel",
			engine:    FuelMizer.Name,
			dest:      Vector{200, 0},
			warpSpeed: 9,
			fuel:      1,
			want:      want{position: Vector{3, 0}, fuel: 0, warpSpeed: 4},
		},
		{
			// 51ly at warp 7 costs 29mg, but 49ly costs 28mg and the last 2ly cost 2mg
			name:      "top up fuel lost to rounding",
			dest:      Vector{51, 0},
			warpSpeed: 7,
			fuel:      29,
			want:      want{position: Vector{49, 0}, fuel: 2, warpSpeed: 7},
		},
		{
			name:      "no top up without enough fuel to arrive",
			dest:      Vector{200, 0},
			warpSpeed: 7,
			fuel:      200,
			want:      want{position: Vector{49, 0}, fuel: 172, warpSpeed: 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleet := testLongRangeScout(player)
			if tt.engine != "" {
				design := fleet.Tokens[0].design
				design.Slots[0].HullComponent = tt.engine
				design.Spec, _ = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, design)
				fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			}
			fleet.withWaypoints(NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(tt.dest, tt.warpSpeed)).
				withFuel(tt.fuel)
			universe := Universe{log: testLogger, Fleets: []*Fleet{fleet}}
			universe.buildMaps([]*Player{player})

			fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, fleet.Position)
			assert.Equal(t, tt.want.fuel, fleet.Fuel)
			assert.Equal(t, tt.want.warpSpeed, fleet.Waypoints[1].WarpSpeed)
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

func TestFleet_CanColonize(t *testing.T) {
	type fields struct {
		cargo     Cargo
		colonizer bool
	}
	type args struct {
		planet *Planet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "can colonize",
			fields: fields{cargo: Cargo{Colonists: 1}, colonizer: true},
			args:   args{planet: &Planet{Spec: PlanetSpec{TerraformedHabitability: 1}}},
			want:   true,
		},
		{
			name:   "cannot colonize no colonists",
			fields: fields{cargo: Cargo{Colonists: 0}, colonizer: true},
			args:   args{planet: &Planet{Spec: PlanetSpec{TerraformedHabitability: 1}}},
			want:   false,
		},
		{
			name:   "cannot colonize no colonizer",
			fields: fields{cargo: Cargo{Colonists: 1}, colonizer: false},
			args:   args{planet: &Planet{Spec: PlanetSpec{TerraformedHabitability: 1}}},
			want:   false,
		},
		{
			name:   "cannot colonize not habitable",
			fields: fields{cargo: Cargo{Colonists: 1}, colonizer: true},
			args:   args{planet: &Planet{Spec: PlanetSpec{TerraformedHabitability: -1}}},
			want:   false,
		},
		{
			name:   "cannot colonize planet owned",
			fields: fields{cargo: Cargo{Colonists: 1}, colonizer: true},
			args:   args{planet: &Planet{MapObject: MapObject{PlayerNum: 1}, Spec: PlanetSpec{TerraformedHabitability: 1}}},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fleet{
				Cargo: tt.fields.cargo,
				Spec: FleetSpec{
					ShipDesignSpec: ShipDesignSpec{
						Colonizer: tt.fields.colonizer,
					},
				},
			}
			if got := f.CanColonize(tt.args.planet); got != tt.want {
				t.Errorf("Fleet.CanColonize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_CanFuel(t *testing.T) {
	type args struct {
		player *Player
		planet *Planet
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "can fuel",
			args: args{
				player: testPlayer(),
				planet: &Planet{MapObject: MapObject{PlayerNum: 1}, Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{DockCapacity: 1}}},
			},
			want: true,
		},
		{
			name: "cannot fuel no planet",
			args: args{
				player: testPlayer(),
			},
			want: false,
		},
		{
			name: "cannot fuel no dock",
			args: args{
				player: testPlayer(),
				planet: &Planet{MapObject: MapObject{PlayerNum: 1}, Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{DockCapacity: 0}}},
			},
			want: false,
		},
		{
			name: "cannot fuel not friends",
			args: args{
				player: testPlayer(),
				planet: &Planet{MapObject: MapObject{PlayerNum: 2}, Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{DockCapacity: 0}}},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFleet(tt.args.player, 1, "fleet", []Waypoint{NewPositionWaypoint(Vector{}, 5)})
			if got := f.CanFuel(tt.args.player, tt.args.planet); got != tt.want {
				t.Errorf("Fleet.CanFuel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_CanRemoteMine(t *testing.T) {
	type fields struct {
		miningRate int
	}
	type args struct {
		player *Player
		planet *Planet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "can remote mine",
			fields: fields{miningRate: 1},
			args: args{
				player: NewPlayer(0, NewRace()).WithNum(1),
				planet: &Planet{},
			},
			want: true,
		},
		{
			name: "cannot remote mine no planet",
			args: args{
				player: NewPlayer(0, NewRace()).WithNum(1),
			},
			want: false,
		},
		{
			name: "cannot remote mine planet owned",
			args: args{
				player: NewPlayer(0, NewRace()).WithNum(1),
				planet: &Planet{MapObject: MapObject{PlayerNum: 1}},
			},
			want: false,
		},
		{
			name: "can remote mine own planets",
			args: args{
				player: NewPlayer(0, NewRace().WithPRT(AR).WithSpec(&rules)).WithNum(1),
				planet: &Planet{MapObject: MapObject{PlayerNum: 1}},
			},
			fields: fields{miningRate: 1},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFleet(tt.args.player, 1, "fleet", []Waypoint{NewPositionWaypoint(Vector{}, 5)})
			f.Spec.MiningRate = tt.fields.miningRate
			if got := f.CanRemoteMine(tt.args.player, tt.args.planet); got != tt.want {
				t.Errorf("Fleet.CanRemoteMine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_CanJump(t *testing.T) {
	type fields struct {
		cargo  Cargo
		tokens []ShipToken
	}
	type args struct {
		player   *Player
		orbiting *Planet
		target   *Planet
		dist     float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "can jump",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1)}},
			},
			args: args{
				player: testPlayer(),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: true,
		},
		{
			name: "can jump with gate",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: JumpGate.Name, HullSlotIndex: 3, Quantity: 1},
					}),
				}},
			},
			args: args{
				player: testPlayer(),
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: true,
		},
		{
			name: "can jump with cargo",
			fields: fields{
				cargo: Cargo{Ironium: 1},
				tokens: []ShipToken{{Quantity: 1, design: NewShipDesign(1, 1).
					WithName("Colony Ship").
					WithHull(ColonyShip.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
					})}},
			},
			args: args{
				player: NewPlayer(0, NewRace().WithPRT(IT).WithSpec(&rules)).WithNum(1).WithRelations([]PlayerRelationship{{Relation: PlayerRelationFriend}}),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: true,
		},
		{
			name: "cannot jump with cargo",
			fields: fields{
				cargo: Cargo{Ironium: 1},
				tokens: []ShipToken{{Quantity: 1, design: NewShipDesign(1, 1).
					WithName("Colony Ship").
					WithHull(ColonyShip.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
					})}},
			},
			args: args{
				player: testPlayer(),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: false,
		},
		{
			name: "cannot jump no orbiting",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1)}},
			},
			args: args{
				player: testPlayer(),
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: false,
		},
		{
			name: "cannot jump no target",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1)}},
			},
			args: args{
				player: testPlayer(),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: false,
		},
		{
			name: "cannot jump too heavy",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1)}},
			},
			args: args{
				player: testPlayer(),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 1,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 1,
					}}},
				dist: 1,
			},
			want: false,
		},
		{
			name: "cannot jump too far",
			fields: fields{
				tokens: []ShipToken{{Quantity: 1, design: testLongRangeScoutDesign(1)}},
			},
			args: args{
				player: testPlayer(),
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    1,
						SafeHullMass: 100,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    1,
						SafeHullMass: 100,
					}}},
				dist: 10,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFleet(tt.args.player, 1, "fleet", []Waypoint{NewPositionWaypoint(Vector{}, 5)})
			f.Cargo = tt.fields.cargo
			f.Tokens = tt.fields.tokens
			for i := range f.Tokens {
				f.Tokens[i].design.WithSpec(&rules, tt.args.player)
			}
			f.Spec = ComputeFleetSpec(&rules, tt.args.player, f)
			if got := f.CanJump(tt.args.player, tt.args.dist, tt.args.orbiting, tt.args.target); got != tt.want {
				t.Errorf("Fleet.CanJump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_GetMinimalWarp(t *testing.T) {
	player := testPlayer()
	type args struct {
		player               *Player
		fuelAlreadyAllocated int
		dist                 float64
		startSpeed           int
		freeSpeed            int
		maxSafeSpeed         int
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  int
	}{
		{
			name:  "49ly warp 7", // one year to go 49 ly
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 49,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 7,
		},
		{
			name:  "50ly warp 5", // two years at warps 7, 6 or 5; pick warp 5
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 50,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 5,
		},
		{
			name:  "73ly warp 7", // two years at warp 7, 3 at warp6
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 73,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 7,
		},
		{
			name:  "36ly warp 6",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 36,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 6,
		},
		{
			name:  "25ly warp 6",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 25,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 5,
		},
		{
			name:  "slow own if run out of fuel",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 200,
				dist:                 200,
				startSpeed:           7,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 6,
		},
		{
			name:  "do not exceed safe warp",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 100,
				startSpeed:           10,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 9,
		},
		{
			name:  "go w10 if w10 is safe",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 100,
				startSpeed:           10,
				freeSpeed:            1,
				maxSafeSpeed:         10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fleet.GetMinimalWarp(tt.args.player, tt.args.fuelAlreadyAllocated, tt.args.dist, tt.args.startSpeed, tt.args.freeSpeed, tt.args.maxSafeSpeed); got != tt.want {
				t.Errorf("Fleet.GetMinimalWarp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_GetMaxWarp(t *testing.T) {
	player := testPlayer()
	type args struct {
		player               *Player
		fuelAlreadyAllocated int
		dist                 float64
		freeSpeed            int
		maxSafeSpeed         int
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  int
	}{
		{
			name:  "91ly warp 9", // one year to go 81 ly, plenty of fuel
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 81,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 9,
		},
		{
			name:  "64ly warp 8", // one year to go 64 ly, plenty of fuel
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 64,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 8,
		},
		{
			name:  "25ly warp 9", // don't go faster than we need, go warp 5 for 25ly
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 25,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 5,
		},
		{
			// make sure we don't run out of fuel over long distances
			// 300 ly at warp 9 would take 338mg of fuel, so go warp 8 (using 282mg fuel)
			name:  "300ly warp 8",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 0,
				dist:                 300,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 8,
		},
		{
			// go even slower if we've already allocated fuel to previous waypoints
			name:  "300ly warp 8",
			fleet: testLongRangeScout(player),
			args: args{
				player:               player,
				fuelAlreadyAllocated: 100,
				dist:                 300,
				freeSpeed:            1,
				maxSafeSpeed:         9,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fleet.GetMaxWarp(tt.args.player, tt.args.fuelAlreadyAllocated, tt.args.dist, tt.args.freeSpeed, tt.args.maxSafeSpeed); got != tt.want {
				t.Errorf("Fleet.GetMaxWarp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_GetWarpSpeed(t *testing.T) {
	player := testPlayer()
	type args struct {
		player               *Player
		dist                 float64
		orbiting             *Planet
		target               *Planet
		fuelAlreadyAllocated int
		fastestWaypoint      bool
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  int
	}{
		{
			name:  "in space",
			fleet: testLongRangeScout(player),
			args: args{
				player: player,
				dist:   25,
			},
			want: 5,
		},
		{
			name:  "can jump",
			fleet: testLongRangeScout(player),
			args: args{
				player: player,
				orbiting: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				target: &Planet{
					MapObject: MapObject{PlayerNum: 1},
					Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
						HasStargate:  true,
						SafeRange:    100,
						SafeHullMass: 100,
					}}},
				dist: 1,
			},
			want: StargateWarpSpeed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fleet.GetWarpSpeed(tt.args.player, tt.args.dist, tt.args.orbiting, tt.args.target, tt.args.fuelAlreadyAllocated, tt.args.fastestWaypoint); got != tt.want {
				t.Errorf("Fleet.GetWarpSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_GetFuelAllocated(t *testing.T) {
	player := testPlayer()
	type args struct {
		player        *Player
		waypointIndex int
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  int
	}{
		{
			name:  "no waypoints no fuel allocated",
			fleet: testLongRangeScout(player),
			args:  args{player: player, waypointIndex: 0},
			want:  0,
		},
		{
			name: "one waypoint w6 5mg fuel used",
			fleet: testLongRangeScout(player).withWaypoints(
				NewPositionWaypoint(Vector{0, 0}, 0),
				NewPositionWaypoint(Vector{36, 0}, 6),
			),
			args: args{player: player, waypointIndex: 1},
			want: 5,
		},
		{
			name: "two waypoints w9 5mg fuel used",
			fleet: testLongRangeScout(player).withWaypoints(
				NewPositionWaypoint(Vector{0, 0}, 0),
				NewPositionWaypoint(Vector{81, 0}, 9),
				NewPositionWaypoint(Vector{162, 0}, 9),
			),
			args: args{player: player, waypointIndex: 2},
			want: 184,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fleet.ComputeFuelUsage(tt.args.player)
			if got := tt.fleet.GetFuelAllocated(tt.args.player, tt.args.waypointIndex); got != tt.want {
				t.Errorf("Fleet.GetFuelAllocated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFleet_AddWaypointStargate(t *testing.T) {
	stargate := PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{HasStargate: true, SafeRange: 500, SafeHullMass: 500}}
	newPlanet := func(num int, position Vector) *Planet {
		planet := NewPlanet().WithNum(num).WithPlayerNum(1).withPosition(position)
		planet.Spec = stargate
		return planet
	}
	tests := []struct {
		name         string
		dest         WaypointDest
		wantStargate bool
	}{
		{"gate to a planet", WaypointDest{MO: MapObject{Type: MapObjectTypePlanet, Num: 2, Position: Vector{100, 0}}}, true},
		{"gate to a fleet at a planet", WaypointDest{MO: MapObject{Type: MapObjectTypeFleet, Num: 5, PlayerNum: 1, Position: Vector{100, 0}}}, true},
		{"no gate for a fleet in space", WaypointDest{MO: MapObject{Type: MapObjectTypeFleet, Num: 5, PlayerNum: 1, Position: Vector{50, 0}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := testPlayer().withPlanetIntels([]*Planet{newPlanet(1, Vector{0, 0}), newPlanet(2, Vector{100, 0})})
			fleet := testLongRangeScout(player).withWaypoints(NewPlanetWaypoint(Vector{0, 0}, 1, "Planet 1", 5))

			fleet.AddWaypoint(player, tt.dest, 0, false)

			assert.Equal(t, tt.wantStargate, fleet.Waypoints[1].WarpSpeed == StargateWarpSpeed)
		})
	}
}

func TestFleet_AddWaypoint(t *testing.T) {
	player := testPlayer()
	ifePlayer := NewPlayer(0, NewRace().WithLRT(IFE).WithSpec(&rules)).WithNum(1)
	type args struct {
		player                       *Player
		dest                         WaypointDest
		currentSelectedWaypointIndex int
		fastestWaypoint              bool
	}
	tests := []struct {
		name         string
		fleet        *Fleet
		args         args
		want         int
		wantWaypoint Waypoint
	}{
		{
			name:  "add 36ly waypoint",
			fleet: testLongRangeScout(player),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{36, 0}},
				currentSelectedWaypointIndex: 0,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{36, 0},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{36, 0},
				},
				WarpSpeed:    6,
				EstFuelUsage: 5,
			},
		},
		{
			name:  "add 81ly waypoint",
			fleet: testLongRangeScout(player),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{81, 0}},
				currentSelectedWaypointIndex: 0,
				fastestWaypoint:              false,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{81, 0},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{81, 0},
				},
				WarpSpeed:    6,
				EstFuelUsage: 11,
			},
		},
		{
			name:  "add 81ly waypoint fastest",
			fleet: testLongRangeScout(player),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{81, 0}},
				currentSelectedWaypointIndex: 0,
				fastestWaypoint:              true,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{81, 0},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{81, 0},
				},
				WarpSpeed:    9,
				EstFuelUsage: 92,
			},
		},
		{
			// uses 201mg of fuel at w9 if dist is rounded up, so make sure we go w8 and don't run out of fuel
			// TODO: change this if we update the fuel usage in game to not count rounded up dists (i.e. 81.9ly dist travelled at W9)
			name: "add ife fastest",
			fleet: testSantaMariaIFE(ifePlayer).
				withPosition(Vector{174, 367}).
				withWaypoints(NewPositionWaypoint(Vector{174, 367}, 0)).
				withCargo(Cargo{Colonists: 25}),
			args: args{
				player:                       ifePlayer,
				dest:                         WaypointDest{Position: Vector{48, 462}},
				currentSelectedWaypointIndex: 0,
				fastestWaypoint:              true,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{48, 462},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{48, 462},
				},
				WarpSpeed:    8,
				EstFuelUsage: 132,
			},
		},
		{
			name: "copy transport tasks",
			fleet: testPrivateer(player, 1).withWaypoints(
				Waypoint{
					Task:           WaypointTaskTransport,
					TransportTasks: WaypointTransportTasks{Ironium: WaypointTransportTask{Action: TransportActionUnloadAll}},
				},
			),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{25, 0}},
				currentSelectedWaypointIndex: 0,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{25, 0},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{25, 0},
				},
				WarpSpeed:      5,
				Task:           WaypointTaskTransport,
				TransportTasks: WaypointTransportTasks{Ironium: WaypointTransportTask{Action: TransportActionUnloadAll}},
				EstFuelUsage:   19,
			},
		},
		{
			name:  "colonize",
			fleet: testSantaMaria(player).withCargo(Cargo{Colonists: 25}),
			args: args{
				player: NewPlayer(0, NewRace().WithSpec(&rules)).withPlanetIntels([]*Planet{
					{
						MapObject: MapObject{Type: MapObjectTypePlanet, Num: 1},
						Spec:      PlanetSpec{TerraformedHabitability: 100},
					},
				}),
				dest: WaypointDest{
					MO: MapObject{
						Position: Vector{49, 0},
						Type:     MapObjectTypePlanet,
						Num:      1,
					},
				},
				currentSelectedWaypointIndex: 0,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{49, 0},
				MapObjectTarget: MapObjectTarget{
					TargetType:     MapObjectTypePlanet,
					TargetPosition: Vector{49, 0},
					TargetNum:      1,
				},
				WarpSpeed:    7,
				Task:         WaypointTaskColonize,
				EstFuelUsage: 95,
			},
		},
		{
			name:  "remote mine",
			fleet: testPotatoBug(player),
			args: args{
				player: NewPlayer(0, NewRace().WithSpec(&rules)).withPlanetIntels([]*Planet{
					{
						MapObject: MapObject{Type: MapObjectTypePlanet, Num: 1},
					},
				}),
				dest: WaypointDest{
					MO: MapObject{
						Position: Vector{25, 0},
						Type:     MapObjectTypePlanet,
						Num:      1,
					},
				},
				currentSelectedWaypointIndex: 0,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{25, 0},
				MapObjectTarget: MapObjectTarget{
					TargetType:     MapObjectTypePlanet,
					TargetPosition: Vector{25, 0},
					TargetNum:      1,
				},
				WarpSpeed:    5,
				EstFuelUsage: 12,
				Task:         WaypointTaskRemoteMining,
			},
		},
		{
			name: "gate",
			fleet: testLongRangeScout(player).withWaypoints(
				// orbiting planet with stargate
				Waypoint{
					Position: Vector{},
					MapObjectTarget: MapObjectTarget{
						TargetType:     MapObjectTypePlanet,
						TargetPosition: Vector{},
						TargetNum:      1,
					},
				},
			),
			args: args{
				// player knows of two planets
				player: NewPlayer(0, NewRace().WithSpec(&rules)).
					WithRelations([]PlayerRelationship{{Relation: PlayerRelationFriend}}).
					withPlanetIntels([]*Planet{
						{
							MapObject: MapObject{Type: MapObjectTypePlanet, Num: 1, PlayerNum: 1},
							Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
								HasStargate:  true,
								SafeRange:    100,
								SafeHullMass: 100,
							}},
						},
						{
							MapObject: MapObject{Type: MapObjectTypePlanet, Num: 2, PlayerNum: 1, Position: Vector{100, 0}},
							Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
								HasStargate:  true,
								SafeRange:    100,
								SafeHullMass: 100,
							}},
						},
					}),
				dest: WaypointDest{
					MO: MapObject{
						Position: Vector{100, 0},
						Type:     MapObjectTypePlanet,
						Num:      2,
					},
				},
				currentSelectedWaypointIndex: 1,
			},
			want: 1,
			wantWaypoint: Waypoint{
				Position: Vector{100, 0},
				MapObjectTarget: MapObjectTarget{
					TargetType:     MapObjectTypePlanet,
					TargetPosition: Vector{100, 0},
					TargetNum:      2,
				},
				WarpSpeed:    StargateWarpSpeed,
				EstFuelUsage: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fleet.AddWaypoint(tt.args.player, tt.args.dest, tt.args.currentSelectedWaypointIndex, tt.args.fastestWaypoint)
			if got != tt.want {
				t.Errorf("Fleet.AddWaypoint() = %v, want %v", got, tt.want)
			}
			if got != 0 {
				test.CompareAsJSON(t, tt.fleet.Waypoints[got], tt.wantWaypoint)
			}
		})
	}
}

func TestFleet_UpdateWaypoint(t *testing.T) {
	player := testPlayer()
	type args struct {
		player                       *Player
		dest                         WaypointDest
		currentSelectedWaypointIndex int
		fastestWaypoint              bool
	}
	tests := []struct {
		name         string
		fleet        *Fleet
		args         args
		want         UpdateWaypointResult
		wantWaypoint Waypoint
	}{
		{
			name:  "no update wp0",
			fleet: testLongRangeScout(player),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{36, 0}},
				currentSelectedWaypointIndex: 0,
			},
			want: UpdateWaypointResultNone,
		},
		{
			name: "no update onto prev waypoint",
			fleet: testLongRangeScout(player).withWaypoints(
				NewPositionWaypoint(Vector{}, 5),
				NewPositionWaypoint(Vector{25, 0}, 5),
			),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{}},
				currentSelectedWaypointIndex: 1,
			},
			want: UpdateWaypointResultNone,
		},
		{
			name: "update 25ly w5 to 36ly w6 waypoint",
			fleet: testLongRangeScout(player).withWaypoints(
				NewPositionWaypoint(Vector{}, 5),
				NewPositionWaypoint(Vector{25, 0}, 5),
			),
			args: args{
				player:                       player,
				dest:                         WaypointDest{Position: Vector{36, 0}},
				currentSelectedWaypointIndex: 1,
			},
			want: UpdateWaypointResultUpdated,
			wantWaypoint: Waypoint{
				Position: Vector{36, 0},
				MapObjectTarget: MapObjectTarget{
					TargetPosition: Vector{36, 0},
				},
				WarpSpeed:    6,
				EstFuelUsage: 5,
			},
		},
		{
			name: "update to colonize planet",
			fleet: testSantaMaria(player).withCargo(Cargo{Colonists: 25}).withWaypoints(
				NewPositionWaypoint(Vector{}, 5),
				NewPositionWaypoint(Vector{25, 0}, 5),
			),
			args: args{
				player: NewPlayer(0, NewRace().WithSpec(&rules)).withPlanetIntels([]*Planet{
					{
						MapObject: MapObject{Type: MapObjectTypePlanet, Num: 1},
						Spec:      PlanetSpec{TerraformedHabitability: 100},
					},
				}),
				dest: WaypointDest{
					MO: MapObject{
						Position: Vector{49, 0},
						Type:     MapObjectTypePlanet,
						Num:      1,
					},
				},
				currentSelectedWaypointIndex: 1,
			},
			want: UpdateWaypointResultUpdated,
			wantWaypoint: Waypoint{
				Position: Vector{49, 0},
				MapObjectTarget: MapObjectTarget{
					TargetType:     MapObjectTypePlanet,
					TargetPosition: Vector{49, 0},
					TargetNum:      1,
				},
				WarpSpeed:    7,
				Task:         WaypointTaskColonize,
				EstFuelUsage: 95,
			},
		},
		{
			name: "update to gate",
			fleet: testLongRangeScout(player).withWaypoints(
				// orbiting planet with stargate
				Waypoint{
					Position: Vector{},
					MapObjectTarget: MapObjectTarget{
						TargetType:     MapObjectTypePlanet,
						TargetPosition: Vector{},
						TargetNum:      1,
					},
				},
				NewPositionWaypoint(Vector{25, 0}, 5),
			),
			args: args{
				// player knows of two planets
				player: NewPlayer(0, NewRace().WithSpec(&rules)).
					WithRelations([]PlayerRelationship{{Relation: PlayerRelationFriend}}).
					withPlanetIntels([]*Planet{
						{
							MapObject: MapObject{Type: MapObjectTypePlanet, Num: 1, PlayerNum: 1},
							Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
								HasStargate:  true,
								SafeRange:    100,
								SafeHullMass: 100,
							}},
						},
						{
							MapObject: MapObject{Type: MapObjectTypePlanet, Num: 2, PlayerNum: 1, Position: Vector{100, 0}},
							Spec: PlanetSpec{PlanetStarbaseSpec: PlanetStarbaseSpec{
								HasStargate:  true,
								SafeRange:    100,
								SafeHullMass: 100,
							}},
						},
					}),
				dest: WaypointDest{
					MO: MapObject{
						Position: Vector{100, 0},
						Type:     MapObjectTypePlanet,
						Num:      2,
					},
				},
				currentSelectedWaypointIndex: 1,
			},
			want: UpdateWaypointResultUpdated,
			wantWaypoint: Waypoint{
				Position: Vector{100, 0},
				MapObjectTarget: MapObjectTarget{
					TargetType:     MapObjectTypePlanet,
					TargetPosition: Vector{100, 0},
					TargetNum:      2,
				},
				WarpSpeed:    StargateWarpSpeed,
				EstFuelUsage: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fleet.UpdateWaypoint(tt.args.player, tt.args.dest, tt.args.currentSelectedWaypointIndex, tt.args.fastestWaypoint)
			if got != tt.want {
				t.Errorf("Fleet.UpdateWaypoint() = %v, want %v", got, tt.want)
			}
			if got != UpdateWaypointResultNone {
				test.CompareAsJSON(t, tt.fleet.Waypoints[tt.args.currentSelectedWaypointIndex], tt.wantWaypoint)
			}
		})
	}
}

func TestFleet_RenameByPlayer(t *testing.T) {
	fleet := testLongRangeScout(NewPlayer(1, NewRace().WithSpec(&rules))).withNum(5)
	assert.Equal(t, "Long Range Scout #5", fleet.Name)

	// player renames don't include the fleet number
	fleet.RenameByPlayer("Scouts")
	assert.Equal(t, "Scouts", fleet.Name)
	assert.Equal(t, "Scouts", fleet.BaseName)
	assert.Equal(t, 5, fleet.Num)

	// system renames still include the fleet number
	fleet.Rename("Scouts")
	assert.Equal(t, "Scouts #5", fleet.Name)
}

func TestFleet_repairFleetRates(t *testing.T) {
	tests := []struct {
		name          string
		prt           PRT
		moved, noHeal bool
		tanker        float64
		want          float64
	}{
		{"moved this turn", JoaT, true, false, 0, 49},
		{"stationary in space", JoaT, false, false, 0, 48},
		{"no healing after gating, battle or mines", JoaT, false, true, 0, 50},
		{"IS doubles the base rate, not the tanker bonus", IS, false, false, .05, 41},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(0)
			p := NewPlayer(1, NewRace().WithPRT(tt.prt).WithSpec(&r)).withSpec(&r)
			design := &ShipDesign{Spec: ShipDesignSpec{Armor: 100}}
			f := &Fleet{Tokens: []ShipToken{{design: design, Quantity: 1, QuantityDamaged: 1, Damage: 50}}, warped: tt.moved, noHeal: tt.noHeal}
			f.Spec.RepairBonus = tt.tanker

			f.repairFleet(testLogger, &r, p, nil)
			assert.Equal(t, tt.want, f.Tokens[0].Damage)
		})
	}
}

func Test_turn_noHealAfterGateAndMines(t *testing.T) {
	t.Run("gating prevents repair", func(t *testing.T) {
		s := ScenarioStargateTest()
		s.Players[0].Fleets[0].Waypoints[1].Warp = StargateWarpSpeed
		u := newTestUniverse(t, s)
		f := u.Game.Fleets[0]
		f.Tokens[0].Damage = 10
		f.Tokens[0].QuantityDamaged = 1

		u.GenerateTurn()

		assert.Equal(t, u.Planet("Planet 2").Position, f.Position)
		assert.Equal(t, float64(10), f.Tokens[0].Damage)
	})
	t.Run("mine damage prevents repair, even after a merge", func(t *testing.T) {
		s := SingleUnitScenario()
		s.Players[0].Fleets = append(s.Players[0].Fleets, s.Players[0].Fleets[0])
		u := newTestUniverse(t, s)
		f := u.Game.Fleets[0]
		m := Minefield{}
		m.damageFleet(f, u.Player(1), MinefieldStats{DamagePerEngine: 1})

		merged, err := NewOrderer().Merge(&u.Game.Rules, u.Player(1), []*Fleet{u.Game.Fleets[1], f})
		require.NoError(t, err)
		before := merged.Tokens[0].Damage
		merged.repairFleet(testLogger, &u.Game.Rules, u.Player(1), u.Planet("Planet 1"))
		assert.Greater(t, before, float64(0))
		assert.Equal(t, before, merged.Tokens[0].Damage)
	})
}
