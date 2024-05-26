package cs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
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
		{
			name: "split a scoutx3 damaged with floats",
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
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 20, Damage: 22.799999999999997, QuantityDamaged: 20},
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 7},
					},
				},
				// split out 20 damaged scouts with float damage, and 5 undamaged freighters
				tokens: []ShipToken{
					{DesignNum: scoutDesign.Num, Quantity: 20},
					{DesignNum: freighterDesign.Num, Quantity: 5},
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
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2},
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
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 20, Damage: 22.799999999999997, QuantityDamaged: 20},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 5},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.Designs = append(player.Designs, scoutDesign, freighterDesign)

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
			gotNewFleet, err := o.splitFleetTokens(&rules, tt.args.player, playerFleets, tt.args.source, tt.args.tokens)
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
		{
			name: "merge a scout and freighter, scout, freighter, scout",
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
				},
				{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       3,
						PlayerNum: player.Num,
						Name:      "Scout #3",
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
						Num:       4,
						PlayerNum: player.Num,
						Name:      "Teamster #4",
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
						Num:       5,
						PlayerNum: player.Num,
						Name:      "Scout #5",
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
						Num:       6,
						PlayerNum: player.Num,
						Name:      "Teamster #6",
					},
					BaseName: "Teamster",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 1},
					},
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
					{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 3},
					{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 3},
				},
			},
			wantErr: false,
		}}
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

func Test_orders_TransferPlanetCargo(t *testing.T) {
	player := NewPlayer(0, NewRace().WithSpec(&rules)).withSpec(&rules)
	type args struct {
		source         *Fleet
		dest           *Planet
		transferAmount CargoTransferRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"transfer 10kT Ironium from planet",
			args{
				source:         testTeamster(player),
				dest:           NewPlanet().WithCargo(Cargo{Ironium: 10}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 10}, 0},
			},
			false,
		},
		{
			"fail to transfer 10kT Ironium from planet",
			args{
				source:         testTeamster(player),
				dest:           NewPlanet().WithCargo(Cargo{Ironium: 5}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 10}, 0},
			},
			true,
		},
		{
			"transfer 10kT Ironium to planet",
			args{
				source:         testTeamster(player).withCargo(Cargo{Ironium: 10}),
				dest:           NewPlanet(),
				transferAmount: CargoTransferRequest{Cargo{Ironium: -10}, 0},
			},
			false,
		},
		{
			"fail to transfer 10kT Ironium to planet",
			args{
				source:         testTeamster(player),
				dest:           NewPlanet(),
				transferAmount: CargoTransferRequest{Cargo{Ironium: -10}, 0},
			},
			true,
		},
		{
			"transfer 210kT Mixed Minerals from planet",
			args{
				source:         testTeamster(player),
				dest:           NewPlanet().WithCargo(Cargo{1000, 1000, 1000, 1000}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 70, Boranium: 70, Germanium: 70}, 0},
			},
			false,
		},
		{
			"fail to transfer 211kT Mixed Minerals from planet",
			args{
				source:         testTeamster(player),
				dest:           NewPlanet().WithCargo(Cargo{1000, 1000, 1000, 1000}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 70, Boranium: 70, Germanium: 70, Colonists: 1}, 0},
			},
			true,
		},
		{
			"transfer 4000kT Mixed Cargo from planet where planet is out of one mineral",
			args{
				source:         testPrivateer(player, 10),
				dest:           NewPlanet().WithCargo(Cargo{2726 + 366, 4763 + 414, 0, 1601 + 3027}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 366, Boranium: 414, Germanium: 193, Colonists: 3027}, 0},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}
			sourceCargo := tt.args.source.Cargo
			destCargo := tt.args.dest.Cargo
			err := o.TransferPlanetCargo(&rules, player, tt.args.source, tt.args.dest, tt.args.transferAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("orders.TransferPlanetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				// we should transfer from the dest to the soruce
				assert.Equal(t, sourceCargo.Add(tt.args.transferAmount.Cargo), tt.args.source.Cargo)
				assert.Equal(t, destCargo.Subtract(tt.args.transferAmount.Cargo), tt.args.dest.Cargo)
			}
		})
	}
}

func Test_orders_TransferFleetCargo(t *testing.T) {
	player := NewPlayer(0, NewRace().WithSpec(&rules)).withSpec(&rules)
	type args struct {
		source         *Fleet
		dest           *Fleet
		transferAmount CargoTransferRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"transfer 10kT Ironium from fleet",
			args{
				source:         testTeamster(player),
				dest:           testTeamster(player).withCargo(Cargo{Ironium: 10}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 10}, 0},
			},
			false,
		},
		{
			"fail to transfer 10kT Ironium from fleet",
			args{
				source:         testTeamster(player),
				dest:           testTeamster(player).withCargo(Cargo{Ironium: 5}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 10}, 0},
			},
			true,
		},
		{
			"transfer 10kT Ironium to fleet",
			args{
				source:         testTeamster(player).withCargo(Cargo{Ironium: 10}),
				dest:           testTeamster(player),
				transferAmount: CargoTransferRequest{Cargo{Ironium: -10}, 0},
			},
			false,
		},
		{
			"fail to transfer 10kT Ironium to fleet",
			args{
				source:         testTeamster(player),
				dest:           testTeamster(player),
				transferAmount: CargoTransferRequest{Cargo{Ironium: -10}, 0},
			},
			true,
		},
		{
			"transfer 210kT Mixed Minerals from fleet",
			args{
				source:         testTeamster(player),
				dest:           testTeamster(player).withCargo(Cargo{Ironium: 70, Boranium: 70, Germanium: 70}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 70, Boranium: 70, Germanium: 70}, 0},
			},
			false,
		},
		{
			"fail to transfer 211kT Mixed Minerals from fleet",
			args{
				source:         testTeamster(player).withCargo(Cargo{Colonists: 1}),
				dest:           testTeamster(player).withCargo(Cargo{Ironium: 70, Boranium: 70, Germanium: 70}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: 70, Boranium: 70, Germanium: 70}, 0},
			},
			true,
		},
		{
			"fail to transfer 211kT Mixed Minerals to fleet",
			args{
				source:         testTeamster(player).withCargo(Cargo{Ironium: 70, Boranium: 70, Germanium: 70}),
				dest:           testTeamster(player).withCargo(Cargo{Colonists: 1}),
				transferAmount: CargoTransferRequest{Cargo{Ironium: -70, Boranium: -70, Germanium: -70}, 0},
			},
			true,
		},
		{
			"transfer 10mg Fuel from fleet",
			args{
				source:         testTeamster(player).withFuel(0),
				dest:           testTeamster(player).withFuel(10),
				transferAmount: CargoTransferRequest{Fuel: 10},
			},
			false,
		},
		{
			"fail to transfer 10mg Fuel from fleet",
			args{
				source:         testTeamster(player).withFuel(0),
				dest:           testTeamster(player).withFuel(5),
				transferAmount: CargoTransferRequest{Fuel: 10},
			},
			true,
		},
		{
			"transfer 10mg Fuel to fleet",
			args{
				source:         testTeamster(player).withFuel(10),
				dest:           testTeamster(player).withFuel(0),
				transferAmount: CargoTransferRequest{Fuel: -10},
			},
			false,
		},
		{
			"fail to transfer 10mg Fuel to fleet",
			args{
				source:         testTeamster(player).withFuel(0),
				dest:           testTeamster(player).withFuel(0),
				transferAmount: CargoTransferRequest{Fuel: -10},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}

			sourceCargo := tt.args.source.Cargo
			sourceFuel := tt.args.source.Fuel
			destCargo := tt.args.dest.Cargo
			destFuel := tt.args.dest.Fuel

			err := o.TransferFleetCargo(&rules, player, player, tt.args.source, tt.args.dest, tt.args.transferAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("orders.TransferFleetCargo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				// we should transfer from the dest to the soruce
				assert.Equal(t, sourceCargo.Add(tt.args.transferAmount.Cargo), tt.args.source.Cargo)
				assert.Equal(t, sourceFuel+tt.args.transferAmount.Fuel, tt.args.source.Fuel)
				assert.Equal(t, destCargo.Subtract(tt.args.transferAmount.Cargo), tt.args.dest.Cargo)
				assert.Equal(t, destFuel-tt.args.transferAmount.Fuel, tt.args.dest.Fuel)
			}

		})
	}
}

func Test_orders_SplitFleet(t *testing.T) {
	player := NewPlayer(0, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
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

	type args struct {
		source         *Fleet
		dest           *Fleet
		sourceTokens   []ShipToken
		destTokens     []ShipToken
		transferAmount CargoTransferRequest
	}

	type want struct {
		err          bool
		errContains  string
		deleteSource bool
		deleteDest   bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{"fail with missing source", args{source: nil, dest: nil}, want{err: true, errContains: "no source fleet"}},
		{"fail with missing tokens", args{source: testLongRangeScoutWithQuantity(player, 2), dest: nil}, want{err: true, errContains: "source fleet tokens and split request tokens don't match"}},
		{
			name: "fail trying to add extra ship to dest",
			args: args{
				source: testLongRangeScoutWithQuantity(player, 2),
				dest:   nil,
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				destTokens: []ShipToken{
					{
						Quantity:  2,
						DesignNum: 1,
					},
				},
			},
			want: want{err: true, errContains: "token in original fleet has different quantity/damage that token in split request"},
		},
		{
			name: "fail trying to add extra ship stack",
			args: args{
				source: testLongRangeScoutWithQuantity(player, 2),
				dest:   nil,
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
					{
						Quantity:  1,
						DesignNum: 2,
					},
				},
				destTokens: []ShipToken{
					{
						Quantity:  2,
						DesignNum: 1,
					},
					{
						Quantity:  1,
						DesignNum: 2,
					},
				},
			},
			want: want{err: true, errContains: "source fleet tokens and split request tokens don't match"},
		},
		{
			name: "split 2 scout fleet into two fleets",
			args: args{
				source: testLongRangeScoutWithQuantity(player, 2),
				dest:   nil,
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				destTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
			},
		},
		{
			name: "split damaged 2 scout fleet into two fleets",
			args: args{
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						// one of these scouts has 10 damage
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2, QuantityDamaged: 1, Damage: 10},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 2,
				},
				dest: nil,
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				// move the damaged token into a new fleet
				destTokens: []ShipToken{
					{
						Quantity:        1,
						DesignNum:       1,
						QuantityDamaged: 1,
						Damage:          10,
					},
				},
			},
		},
		{
			name: "split damaged 3 scout fleet",
			args: args{
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						// one of these scouts has 10 damage
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 3, QuantityDamaged: 2, Damage: 10},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 3,
				},
				dest: nil,
				// keep one of the damaged tokens in the old fleet
				sourceTokens: []ShipToken{
					{
						Quantity:        2,
						DesignNum:       1,
						QuantityDamaged: 1,
						Damage:          10,
					},
				},
				// move one of the damaged tokens into a new fleet
				destTokens: []ShipToken{
					{
						Quantity:        1,
						DesignNum:       1,
						QuantityDamaged: 1,
						Damage:          10,
					},
				},
			},
		},
		{
			name: "fail if split tries to remove damage",
			args: args{
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						// one of these scouts has 10 damage
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 3, QuantityDamaged: 2, Damage: 10},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 3,
				},
				dest: nil,
				// pretend like our source fleet is undamaged (cheater!, or more likely a UI bug...)
				sourceTokens: []ShipToken{
					{
						Quantity:  2,
						DesignNum: 1,
					},
				},
				// move one of the damaged tokens into a new fleet
				destTokens: []ShipToken{
					{
						Quantity:        1,
						DesignNum:       1,
						QuantityDamaged: 1,
						Damage:          10,
					},
				},
			},
			want: want{err: true, errContains: "token in original fleet has different quantity/damage that token in split request"},
		},
		{
			name: "merge two indepdently damaged fleets",
			args: args{
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						// one of these scouts has 10 damage
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2, QuantityDamaged: 1, Damage: 10},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 3,
				},
				dest: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						// one of these scouts has 5 damage
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2, QuantityDamaged: 1, Damage: 5},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 3,
				},
				// move the 10 dmg scout to the 5dmg scout fleet
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				// move one of the damaged tokens into a new fleet
				destTokens: []ShipToken{
					{
						Quantity:        3,
						DesignNum:       1,
						QuantityDamaged: 2,
						Damage:          (10. + 5.) / 2.,
					},
				},
			},
		},
		{
			name: "split 2 freighters",
			args: args{
				source: testSmallFreighterWithQuantity(player, 2).withCargo(Cargo{10, 10, 10, 10}),
				dest:   nil,
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				destTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: 1,
					},
				},
				transferAmount: CargoTransferRequest{Cargo: Cargo{-5, -5, -5, -5}, Fuel: -130},
			},
		},
		{
			name: "split mixed fleet of 2 scouts <-> 2 freighters into one of each",
			args: args{
				source: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       1,
						PlayerNum: player.Num,
						Name:      "Fleet #1",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: scoutDesign, DesignNum: scoutDesign.Num, Quantity: 2},
					},
					Fuel: scoutDesign.Spec.FuelCapacity * 2, // fully fueled
				},
				dest: &Fleet{
					MapObject: MapObject{
						Type:      MapObjectTypeFleet,
						Num:       2,
						PlayerNum: player.Num,
						Name:      "Fleet #2",
					},
					BaseName: "Fleet",
					FleetOrders: FleetOrders{
						Waypoints: []Waypoint{NewPositionWaypoint(Vector{}, 5)},
					},
					Tokens: []ShipToken{
						{design: freighterDesign, DesignNum: freighterDesign.Num, Quantity: 2},
					},
					Fuel: freighterDesign.Spec.FuelCapacity * 2, // fully fueled
				},
				sourceTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: scoutDesign.Num,
					},
					{
						Quantity:  1,
						DesignNum: freighterDesign.Num,
					},
				},
				destTokens: []ShipToken{
					{
						Quantity:  1,
						DesignNum: scoutDesign.Num,
					},
					{
						Quantity:  1,
						DesignNum: freighterDesign.Num,
					},
				},
				// give one freighter's worth of fuel but take one scout's worth of fuel
				transferAmount: CargoTransferRequest{Fuel: freighterDesign.Spec.FuelCapacity - scoutDesign.Spec.FuelCapacity},
			},
		},
		{
			name: "delete source",
			args: args{
				source: testLongRangeScoutWithQuantity(player, 1).withNum(1),
				dest:   testLongRangeScoutWithQuantity(player, 1).withNum(2),
				destTokens: []ShipToken{
					{
						Quantity:  2,
						DesignNum: 1,
					},
				},
			},
			want: want{deleteSource: true},
		},
		{
			name: "delete dest",
			args: args{
				source: testLongRangeScoutWithQuantity(player, 1).withNum(1),
				dest:   testLongRangeScoutWithQuantity(player, 1).withNum(2),
				sourceTokens: []ShipToken{
					{
						Quantity:  2,
						DesignNum: 1,
					},
				},
			},
			want: want{deleteDest: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &orders{}

			designsByNum := make(map[int]*ShipDesign)
			var sourceCargo, destCargo Cargo
			var sourceFuel, destFuel int
			if tt.args.source != nil {
				sourceCargo = tt.args.source.Cargo
				sourceFuel = tt.args.source.Fuel
				for _, token := range tt.args.source.Tokens {
					designsByNum[token.DesignNum] = token.design
				}
			}

			if tt.args.dest != nil {
				destCargo = tt.args.dest.Cargo
				destFuel = tt.args.dest.Fuel
				for _, token := range tt.args.dest.Tokens {
					if _, found := designsByNum[token.DesignNum]; !found {
						designsByNum[token.DesignNum] = token.design
					}
				}
			}

			player.Designs = make([]*ShipDesign, 0, len(designsByNum))
			for designNum := range designsByNum {
				player.Designs = append(player.Designs, designsByNum[designNum])
			}

			playerFleets := []*Fleet{tt.args.source}

			source, dest, err := o.SplitFleet(&rules, player, playerFleets, SplitFleetRequest{
				Source:         tt.args.source,
				Dest:           tt.args.dest,
				SourceTokens:   tt.args.sourceTokens,
				DestTokens:     tt.args.destTokens,
				TransferAmount: tt.args.transferAmount,
			})
			if (err != nil) != tt.want.err {
				t.Errorf("orders.SplitFleet() error = %v, wantErr %v", err, tt.want.err)
			}
			if err != nil && !strings.Contains(fmt.Sprint(err), tt.want.errContains) {
				t.Errorf("orders.SplitFleet() error = %v, wantErrContains %s", err, tt.want.errContains)
			}
			if err == nil {
				// the dest and source should have the passed in tokens
				assert.Equal(t, len(source.Tokens), len(tt.args.sourceTokens))
				assert.Equal(t, len(dest.Tokens), len(tt.args.destTokens))

				// we should transfer from the dest to the soruce
				assert.Equal(t, sourceCargo.Add(tt.args.transferAmount.Cargo), source.Cargo)
				assert.Equal(t, sourceFuel+tt.args.transferAmount.Fuel, source.Fuel)
				assert.Equal(t, destCargo.Subtract(tt.args.transferAmount.Cargo), dest.Cargo)
				assert.Equal(t, destFuel-tt.args.transferAmount.Fuel, dest.Fuel)

				if tt.want.deleteSource {
					assert.True(t, source.Delete)
				}
				if tt.want.deleteDest {
					assert.True(t, dest.Delete)
				}
			}

		})
	}
}
