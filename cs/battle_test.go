package cs

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStalwardDefender(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Stalwart Defender",
		Tokens: []ShipToken{
			{
				DesignNum: 1,
				Quantity:  1,
				design: NewShipDesign(player, 1).
					WithHull(Destroyer.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: BetaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: XRayLaser.Name, HullSlotIndex: 3, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
						{HullComponent: FuelTank.Name, HullSlotIndex: 6, Quantity: 1},
						{HullComponent: BattleComputer.Name, HullSlotIndex: 6, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: NotOrbitingPlanet,
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

// create a new small freighter (with cargo pod) fleet for testing
func testTeamster(player *Player) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{
			PlayerNum: player.Num,
		},
		BaseName: "Teamster",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
					WithHull(SmallFreighter.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
						{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: NotOrbitingPlanet,
	}

	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet

}

func Test_battle_regenerateShields(t *testing.T) {
	type args struct {
		player *Player
		token  battleToken
	}
	tests := []struct {
		name        string
		args        args
		wantShields int
	}{
		{
			name: "no regen",
			args: args{
				player: testPlayer().WithNum(1),
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					Shields:      50,
					TotalShields: 100,
				},
			},
			wantShields: 50,
		},
		{
			name: "regen",
			args: args{
				player: NewPlayer(1, NewRace().WithLRT(RS).WithSpec(&rules)).WithNum(1),
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					Shields:      50,
					TotalShields: 100,
				},
			},
			wantShields: 60,
		},
		{
			name: "no regen when shields gone",
			args: args{
				player: NewPlayer(1, NewRace().WithLRT(RS).WithSpec(&rules)).WithNum(1),
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					Shields:      0,
					TotalShields: 100,
				},
			},
			wantShields: 0,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			battle := battle{
				players: map[int]*Player{tt.args.player.Num: tt.args.player},
			}
			battle.regenerateShields(&tt.args.token)

			got := tt.args.token.Shields
			if got != tt.wantShields {
				t.Errorf("battle.regenerateShields() = %v, want %v", got, tt.wantShields)
			}

		})
	}
}

func Test_battle_RunBattle1(t *testing.T) {
	player1 := testPlayer().WithNum(1)
	player2 := testPlayer().WithNum(2)
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}

	fleets := []*Fleet{
		testStalwardDefender(player1),
		testTeamster(player2),
	}

	battle := newBattler(&rules, &StaticTechStore, map[int]*Player{1: player1, 2: player2}, fleets, nil)

	record := battle.RunBattle()

	// ran some number of turns
	assert.Greater(t, len(record.ActionsPerRound), 1)
}
