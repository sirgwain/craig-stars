package cs

import (
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

func Test_orders_SplitFleetTokens(t *testing.T) {
	player := testPlayer().WithNum(1)
	scoutDesign := NewShipDesign(player, 1).
		WithName("Long Range Scout").
		WithHull(Scout.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	freighterDesign := NewShipDesign(player, 2).
		WithName("Teamster").
		WithHull(SmallFreighter.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	player.Designs = append(player.Designs, scoutDesign, freighterDesign)

	type args struct {
		player *Player
		source *Fleet
		tokens []ShipToken
	}
	tests := []struct {
		name            string
		args            args
		wantSourceFleet *Fleet
		wantNewFleet    *Fleet
		wantErr         bool
	}{
		{"nil fleet, should err", args{}, nil, nil, true},
		{"empty tokens, should err", args{source: testLongRangeScout(player), tokens: []ShipToken{}}, nil, nil, true},
		{"split into more tokens, should err", args{source: testLongRangeScout(player),
			tokens: []ShipToken{
				{DesignNum: 1, Quantity: 10}, // try and create 10 scouts from thin air
			},
		}, nil, nil, true},
		{"split into new token, should err", args{source: testLongRangeScout(player),
			tokens: []ShipToken{
				{DesignNum: 2, Quantity: 1}, // try and create 1 new freighter from thin air
			},
		}, nil, nil, true},
		{"split and leave no tokens in source, should err", args{source: testLongRangeScout(player),
			tokens: []ShipToken{
				{DesignNum: 1, Quantity: 1}, // try and create 10 scouts from thin air
			},
		}, nil, nil, true},
		{
			name: "split a scoutx2 into two fleets",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 2, // fully fueled
				},
				tokens: []ShipToken{
					{DesignNum: scoutDesign.Num, Quantity: 1},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
				Fuel: scoutDesign.Spec.FuelCapacity,
			},
			wantNewFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       2,
					PlayerNum: player.Num,
					Name:      "Long Range Scout #2",
				},
				BaseName: "Long Range Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
				Fuel: scoutDesign.Spec.FuelCapacity,
			},
			wantErr: false,
		},
		{
			name: "split a scoutx2 into two fleets with low fuel",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2},
					},
					Fuel: 7, // only 7mg of fuel, oh noes!
				},
				tokens: []ShipToken{
					{DesignNum: scoutDesign.Num, Quantity: 1},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
				Fuel: 4,
			},
			wantNewFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       2,
					PlayerNum: player.Num,
					Name:      "Long Range Scout #2",
				},
				BaseName: "Long Range Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
				Fuel: 3, // the dest fleet rounds down
			},
			wantErr: false,
		},
		{
			name: "split a scoutx2 and freighterx2 into two fleets",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2},
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2},
					},
					Fuel:  scoutDesign.Spec.FuelCapacity*2 + freighterDesign.Spec.FuelCapacity*2,
					Cargo: Cargo{10, 20, 30, 40},
				},
				tokens: []ShipToken{
					{DesignNum: scoutDesign.Num, Quantity: 1},
					{DesignNum: freighterDesign.Num, Quantity: 1},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
				},
				Fuel:  scoutDesign.Spec.FuelCapacity + freighterDesign.Spec.FuelCapacity,
				Cargo: Cargo{5, 10, 15, 20}, // half the cargo moves over
			},
			wantNewFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       2,
					PlayerNum: player.Num,
					Name:      "Long Range Scout #2",
				},
				BaseName: "Long Range Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
				},
				Fuel:  scoutDesign.Spec.FuelCapacity + freighterDesign.Spec.FuelCapacity, // half the fuel moves over
				Cargo: Cargo{5, 10, 15, 20},                                              // half the cargo moves over
			},
			wantErr: false,
		},
		{
			name: "split a single freighter out of 3",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Teamster #1",
					},
					BaseName: "Teamster",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 3},
					},
					Fuel:  freighterDesign.Spec.FuelCapacity * 3,
					Cargo: Cargo{12, 21, 33, 45},
				},
				tokens: []ShipToken{
					{DesignNum: freighterDesign.Num, Quantity: 1}, // split out one freighter
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Teamster #1",
				},
				BaseName: "Teamster",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2},
				},
				Fuel:  freighterDesign.Spec.FuelCapacity * 2, // keep 2/3rd of the fuel
				Cargo: Cargo{8, 14, 22, 30},                  // 2/3rd of the cargo
			},
			wantNewFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       2,
					PlayerNum: player.Num,
					Name:      "Teamster #2",
				},
				BaseName: "Teamster",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
				},
				Fuel:  freighterDesign.Spec.FuelCapacity, // keep 1/3rd of the fuel
				Cargo: Cargo{4, 7, 11, 15},               // 1/3rd of the cargo
			},
			wantErr: false,
		},

		{
			name: "split a scoutx2 with 1 damaged and freighterx3 with 3 damaged into two fleets",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2, Damage: 5, QuantityDamaged: 1},
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 3, Damage: 6, QuantityDamaged: 2},
					},
				},
				// split out one of each token
				tokens: []ShipToken{
					{DesignNum: scoutDesign.Num, Quantity: 1},
					{DesignNum: freighterDesign.Num, Quantity: 1},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},                                        // should leave no damaged scouts
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2, Damage: 6, QuantityDamaged: 1}, // should leave 1 damaged freighter
				},
			},
			wantNewFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       2,
					PlayerNum: player.Num,
					Name:      "Long Range Scout #2",
				},
				BaseName: "Long Range Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1, Damage: 5, QuantityDamaged: 1},         // gets the damaged scout
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1, Damage: 6, QuantityDamaged: 1}, // gets 1 damaged freighter
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}

			// we assume the player knows about the source fleet
			// and the fleets have specs computed
			playerFleets := []*Fleet{}

			if tt.args.source != nil {
				playerFleets = append(playerFleets, tt.args.source)
			}
			for _, fleet := range playerFleets {
				fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			}
			gotNewFleet, err := o.SplitFleetTokens(&rules, tt.args.player, playerFleets, tt.args.source, tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("orders.SplitFleetTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				// compute the spec for our wantSourceFleet. No need to pass this one in
				tt.wantSourceFleet.Spec = ComputeFleetSpec(&rules, player, tt.wantSourceFleet)
				tt.wantNewFleet.Spec = ComputeFleetSpec(&rules, player, tt.wantNewFleet)
				if !test.CompareAsJSON(t, tt.args.source, tt.wantSourceFleet) {
					t.Errorf("orders.SplitFleetTokens() gotSourceFleet = \n%v, wantSourceFleet = \n%v", tt.args.source, tt.wantSourceFleet)
				}
				if !test.CompareAsJSON(t, gotNewFleet, tt.wantNewFleet) {
					t.Errorf("orders.SplitFleetTokens() gotNewFleet = \n%v, wantNewFleet = \n%v", gotNewFleet, tt.wantNewFleet)
				}
			}
		})
	}
}

func Test_orders_SplitAll(t *testing.T) {
	player := testPlayer().WithNum(1)
	scoutDesign := NewShipDesign(player, 1).
		WithName("Long Range Scout").
		WithHull(Scout.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	freighterDesign := NewShipDesign(player, 2).
		WithName("Teamster").
		WithHull(SmallFreighter.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	freighter2Design := NewShipDesign(player, 3).
		WithName("Teamster2").
		WithHull(SmallFreighter.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	player.Designs = append(player.Designs, scoutDesign, freighterDesign, freighter2Design)

	type args struct {
		player *Player
		source *Fleet
	}
	tests := []struct {
		name            string
		args            args
		wantSourceFleet *Fleet
		wantNewFleets   []*Fleet
		wantErr         bool
	}{
		{
			name: "split a scoutx3 into three fleets",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 3},
					},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
			},
			wantNewFleets: []*Fleet{
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Long Range Scout #2",
					},
					BaseName: "Long Range Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					},
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       3,
						PlayerNum: player.Num,
						Name:      "Long Range Scout #3",
					},
					BaseName: "Long Range Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "split one scout, one freighter1, one freighter2 into three fleets",
			args: args{
				player: player,
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
						{design: freighter2Design, DesignNum: freighter2Design.Num, Quantity: 1},
					},
				},
			},
			wantSourceFleet: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
				},
			},
			wantNewFleets: []*Fleet{
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Teamster #2",
					},
					BaseName: "Teamster",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
					},
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       3,
						PlayerNum: player.Num,
						Name:      "Teamster2 #3",
					},
					BaseName: "Teamster2",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighter2Design, DesignNum: freighter2Design.Num, Quantity: 1},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}

			// we assume the player knows about the source fleet
			// and the fleets have specs computed
			playerFleets := []*Fleet{}

			if tt.args.source != nil {
				playerFleets = append(playerFleets, tt.args.source)
			}
			for _, fleet := range playerFleets {
				fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			}

			gotNewFleets, err := o.SplitAll(&rules, tt.args.player, playerFleets, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("orders.SplitFleetTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				// compute the spec for our wantSourceFleet. No need to pass this one in
				tt.wantSourceFleet.Spec = ComputeFleetSpec(&rules, player, tt.wantSourceFleet)

				if !test.CompareAsJSON(t, tt.args.source, tt.wantSourceFleet) {
					t.Errorf("orders.SplitFleetTokens() gotSourceFleet = \n%v, wantSourceFleet = \n%v", tt.args.source, tt.wantSourceFleet)
				}

				for i, fleet := range tt.wantNewFleets {
					fleet.Spec = ComputeFleetSpec(&rules, player, fleet)

					if !test.CompareAsJSON(t, gotNewFleets[i], fleet) {
						t.Errorf("orders.SplitFleetTokens() gotNewFleet = \n%v, wantNewFleet = \n%v", gotNewFleets[i], fleet)
					}
				}
			}
		})
	}
}

func Test_orders_Merge(t *testing.T) {

	player := testPlayer().WithNum(1)
	scoutDesign := NewShipDesign(player, 1).
		WithName("Long Range Scout").
		WithHull(Scout.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	freighterDesign := NewShipDesign(player, 2).
		WithName("Teamster").
		WithHull(SmallFreighter.Name).
		WithSlots([]ShipDesignSlot{
			{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
			{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
			{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
		}).
		WithSpec(&rules, player)

	player.Designs = append(player.Designs, scoutDesign, freighterDesign)

	tests := []struct {
		name    string
		fleets  []*Fleet
		want    *Fleet
		wantErr bool
	}{
		{"merge with no fleets, should err", []*Fleet{}, nil, true},
		{"merge with one fleet, should err", []*Fleet{testLongRangeScout(player)}, nil, true},
		{
			name: "merge two scouts into each other with a little fuel",
			fleets: []*Fleet{
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					},
					Fuel: 10,
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Scout #2",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					},
					Fuel: 15,
				},
			},
			want: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2},
				},
				Fuel: 25,
			},
			wantErr: false,
		},
		{
			name: "merge a scout and freighter with cargo",
			fleets: []*Fleet{
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Scout #1",
					},
					BaseName: "Scout",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					},
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Teamster #2",
					},
					BaseName: "Teamster",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
					},
					Cargo: Cargo{10, 20, 30, 40},
				},
			},
			want: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Scout #1",
				},
				BaseName: "Scout",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
				},
				Cargo: Cargo{10, 20, 30, 40},
			},
			wantErr: false,
		},
		{
			name: "merge two fleets with damage",
			fleets: []*Fleet{
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Damaged #1",
					},
					BaseName: "Damaged",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 1},                                         // no damage
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 3, QuantityDamaged: 2, Damage: 25}, // 2@50 damage = 10 total damage
					},
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Damaged #2",
					},
					BaseName: "Damaged",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2, QuantityDamaged: 2, Damage: 10},         // 2@10 damage
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2, QuantityDamaged: 1, Damage: 10}, // 1@10 Damage = 10 total damage
					},
				},
			},
			want: &Fleet{
				MapObject: MapObject{
					Type:      MapObjectTypeFleet,
					Num:       1,
					PlayerNum: player.Num,
					Name:      "Damaged #1",
				},
				BaseName: "Damaged",
				FleetOrders: FleetOrders{
					Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
				},
				Tokens: []ShipToken{
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 3, QuantityDamaged: 2, Damage: 10},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 5, QuantityDamaged: 3, Damage: 20}, // 60 total damage is 3@20 Damage
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}

			// compute specs for fleets we pass in. It's assumed these will be there
			for _, fleet := range tt.fleets {
				fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
			}

			if tt.want != nil {
				tt.want.Spec = ComputeFleetSpec(&rules, player, tt.want)
			}

			got, err := o.Merge(&rules, player, tt.fleets)
			if (err != nil) != tt.wantErr {
				t.Errorf("orders.Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("orders.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
