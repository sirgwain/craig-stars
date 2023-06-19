package cs

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStalwartDefender(player *Player) *Fleet {
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
						{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
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
					shields:      50,
					totalShields: 100,
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
					shields:      50,
					totalShields: 100,
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
					shields:      0,
					totalShields: 100,
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

			got := tt.args.token.shields
			if got != tt.wantShields {
				t.Errorf("battle.regenerateShields() = %v, want %v", got, tt.wantShields)
			}

		})
	}
}

func Test_battle_willTarget(t *testing.T) {

	type args struct {
		target BattleTarget
		token  battleToken
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// if our token has armed/starbase attributes, it should only target armed or starbases
		{args: args{BattleTargetAny, battleToken{}}, want: true},
		{args: args{BattleTargetStarbase, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: true},
		{args: args{BattleTargetArmedShips, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: true},
		{args: args{BattleTargetNone, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: false},
		{args: args{BattleTargetBombersFreighters, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: false},
		{args: args{BattleTargetUnarmedShips, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: false},
		{args: args{BattleTargetFuelTransports, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: false},
		{args: args{BattleTargetFreighters, battleToken{attributes: battleTokenAttributeArmed | battleTokenAttributeStarbase}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &battle{}
			if got := b.willTarget(tt.args.target, &tt.args.token); got != tt.want {
				t.Errorf("battle.willTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleToken_getDistanceAway(t *testing.T) {

	type args struct {
		position Vector
	}
	tests := []struct {
		name string
		bt   battleToken
		args args
		want int
	}{
		{"no distance", battleToken{position: Vector{0, 0}}, args{Vector{0, 0}}, 0},
		{"x distance greatest", battleToken{position: Vector{2, 1}}, args{Vector{4, 2}}, 2},
		{"y distance greatest", battleToken{position: Vector{1, 2}}, args{Vector{2, 5}}, 3},
		{"negative distance (token behind)", battleToken{position: Vector{1, 1}}, args{Vector{0, 0}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.getDistanceAway(tt.args.position); got != tt.want {
				t.Errorf("battleToken.getDistanceAway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleWeaponSlot_isInRangePosition(t *testing.T) {
	type args struct {
		position Vector
	}
	tests := []struct {
		name   string
		weapon battleWeaponSlot
		args   args
		want   bool
	}{
		{"no distance, in range", battleWeaponSlot{token: &battleToken{position: Vector{0, 0}}}, args{Vector{0, 0}}, true},
		{"distance 1, in range", battleWeaponSlot{token: &battleToken{position: Vector{0, 0}}, weaponRange: 1}, args{Vector{1, 1}}, true},
		{"distance 2, out of range", battleWeaponSlot{token: &battleToken{position: Vector{0, 0}}, weaponRange: 1}, args{Vector{1, 2}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.weapon.isInRangePosition(tt.args.position); got != tt.want {
				t.Errorf("battleWeaponSlot.isInRangePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battle_runBattle1(t *testing.T) {
	player1 := testPlayer().WithNum(1)
	player2 := testPlayer().WithNum(2)
	player1.Relations = []PlayerRelationship{{Relation: PlayerRelationFriend}, {Relation: PlayerRelationEnemy}}
	player2.Relations = []PlayerRelationship{{Relation: PlayerRelationEnemy}, {Relation: PlayerRelationFriend}}

	fleets := []*Fleet{
		testStalwartDefender(player1),
		testLongRangeScout(player1),
		testTeamster(player2),
	}

	designNum := 1
	for _, fleet := range fleets {
		for _, token := range fleet.Tokens {
			token.design.Num = designNum
			designNum += 1
		}
	}

	battle := newBattler(&rules, &StaticTechStore, map[int]*Player{1: player1, 2: player2}, fleets, nil)

	record := battle.runBattle()

	// ran some number of turns
	assert.Greater(t, len(record.ActionsPerRound), 1)

	gotJson, _ := json.MarshalIndent(record, "", "  ")
	_ = ioutil.WriteFile("../tmp/battle.json", gotJson, 0644)
}
