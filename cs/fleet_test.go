package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

// create a new long rang scout fleet for testing
func testLongRangeScout(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Long Range Scout",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
					WithHull(Scout.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
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
func testSmallFreighter(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Num:       1,
		},
		BaseName: "Small Freighter",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
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
				ScanRange:    NoScanner,
				ScanRangePen: NoScanner,
				SpaceDock:    UnlimitedSpaceDock,
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
				Cost:         Cost{15, 2, 7, 19},
				FuelCapacity: 300,
				ScanRange:    66,
				ScanRangePen: 30,
				Scanner:      true,
				Mass:         20,
				Armor:        20,
				Engine: Engine{
					IdealSpeed: QuickJump5.IdealSpeed,
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
				Cost:         Cost{12, 2, 9, 20},
				FuelCapacity: 50,
				ScanRange:    66,
				ScanRangePen: 30,
				Scanner:      true,
				Mass:         19,
				Armor:        20,
				CloakUnits:   70,
				CloakPercent: 35,
				Engine: Engine{
					IdealSpeed: QuickJump5.IdealSpeed,
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
				Cost:         Cost{12, 2, 9, 20}.MultiplyInt(2),
				FuelCapacity: 50 * 2,
				ScanRange:    66,
				ScanRangePen: 30,
				Scanner:      true,
				Mass:         19 * 2,
				Armor:        20 * 2,
				CloakUnits:   70,
				CloakPercent: 35, // still 35%
				Engine: Engine{
					IdealSpeed: QuickJump5.IdealSpeed,
				},
			},
			Purposes:         map[ShipDesignPurpose]bool{},
			MassEmpty:        19 * 2,
			BaseCloakedCargo: 0,
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
	}
	type want struct {
		position Vector
		fuelUsed int
	}
	tests := []struct {
		name  string
		fleet *Fleet
		args  args
		want  want
	}{
		{
			"move 25ly at warp5",
			testLongRangeScout(player).withWaypoints([]Waypoint{NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{50, 0}, 5)}),
			args{player},
			want{Vector{25, 0}, 4},
		},
		{
			"move 1ly at warp 1",
			testLongRangeScout(player).withWaypoints([]Waypoint{NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{1, 1}, 1)}),
			args{player},
			want{Vector{1, 1}, 0},
		},
		{
			"overshoot waypoint at warp 5",
			testLongRangeScout(player).withWaypoints([]Waypoint{NewPositionWaypoint(Vector{0, 0}, 0), NewPositionWaypoint(Vector{5, 5}, 5)}),
			args{player},
			want{Vector{5, 5}, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			universe := Universe{Fleets: []*Fleet{tt.fleet}, rules: &rules}

			tt.fleet.moveFleet(&rules, &universe, newTestPlayerGetter(player))

			assert.Equal(t, tt.want.position, tt.fleet.Position)
			assert.Equal(t, tt.want.position, tt.fleet.Waypoints[0].Position)
			assert.Equal(t, tt.want.fuelUsed, tt.fleet.Spec.FuelCapacity-tt.fleet.Fuel)
		})
	}
}

func TestFleet_gateFleet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1)
	player.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}}
	sourcePlanet := NewPlanet().WithNum(1).WithPlayerNum(1)
	sourcePlanet.Spec = PlanetSpec{
		HasStargate:  true,
		SafeRange:    100,
		SafeHullMass: 100,
		MaxRange:     500,
		MaxHullMass:  500,
	}
	destPlanet := NewPlanet().WithNum(2).WithPlayerNum(1)
	destPlanet.Spec = PlanetSpec{
		HasStargate:  true,
		SafeRange:    100,
		SafeHullMass: 100,
		MaxRange:     500,
		MaxHullMass:  500,
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
		name  string
		fleet *Fleet
		args  args
		want  want
	}{
		{
			name: "gate between planets",
			fleet: testLongRangeScout(player).
				withOrbitingPlanetNum(sourcePlanet.Num).
				withWaypoints([]Waypoint{NewPlanetWaypoint(Vector{0, 0}, 1, "planet 1", 5), NewPlanetWaypoint(Vector{50, 0}, 2, "planet 2", StargateWarpSpeed)}),
			args: args{player: player, players: []*Player{player}, planets: []*Planet{sourcePlanet, destPlanet}},
			want: want{position: Vector{50, 0}, orbitingPlanetNum: destPlanet.Num},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, token := range tt.fleet.Tokens {
				player.Designs = append(player.Designs, token.design)
			}
			universe := Universe{Fleets: []*Fleet{tt.fleet}, Planets: tt.args.planets, rules: &rules, designsByNum: map[playerObject]*ShipDesign{}}
			universe.buildMaps(tt.args.players)

			tt.fleet.gateFleet(&rules, &universe, newTestPlayerGetter(player))

			if tt.fleet.Position != tt.want.position {
				t.Errorf("Fleet.gateFleet() position = %v, want %v", tt.fleet.Position, tt.want.position)
			}
			if tt.fleet.OrbitingPlanetNum != tt.want.orbitingPlanetNum {
				t.Errorf("Fleet.gateFleet() OrbitingPlanetNum = %v, want %v", tt.fleet.OrbitingPlanetNum, tt.want.orbitingPlanetNum)
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
			name:               "load all ironium we can fit",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 120, // small freighter has 120kT cargo capacity
		},
		{
			name:               "load all ironium we can fit (we already loaded 20)",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Boranium: 20}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionLoadAll}},
			wantTransferAmount: 100,
		},
		{
			name:               "load fill percent",
			fleet:              testSmallFreighter(player),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionFillPercent, Amount: 50}},
			wantTransferAmount: 60, // 50% of 120kT capacity
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
			name:               "set amount to 20kT when we have 10kT already, but planet only has 5k",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 10}),
			args:               args{dest: NewPlanet().WithCargo(Cargo{Ironium: 5}), cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: 5,    // load 5kT more
			wantWaitAtWaypoint: true, // wait for remaining 5kT we want
		},
		{
			name:               "set amount to 20kT when we have 30kT already. We should unload 10kT to the planet",
			fleet:              testSmallFreighter(player).withCargo(Cargo{Ironium: 30}),
			args:               args{dest: planet, cargoType: Ironium, task: WaypointTransportTask{Action: TransportActionSetAmountTo, Amount: 20}},
			wantTransferAmount: -10, // unload 10kT
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
