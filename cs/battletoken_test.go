//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{"no distance", battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{0, 0}}}, args{Vector{0, 0}}, 0},
		{"x distance greatest", battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{2, 1}}}, args{Vector{4, 2}}, 2},
		{"y distance greatest", battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{1, 2}}}, args{Vector{2, 5}}, 3},
		{"negative distance (token behind)", battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{1, 1}}}, args{Vector{0, 0}}, 1},
		{"3,4 to 7,4", battleToken{BattleRecordToken: BattleRecordToken{Position: Vector{3, 4}}}, args{Vector{7, 4}}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.getDistanceAway(tt.args.position); got != tt.want {
				t.Errorf("battleToken.getDistanceAway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleToken_isTargetOf(t *testing.T) {

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
			if got := tt.args.token.isTargetOf(tt.args.target); got != tt.want {
				t.Errorf("battleToken.isTargetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleToken_regenerateShields(t *testing.T) {
	type args struct {
		token battleToken
	}
	tests := []struct {
		name        string
		args        args
		wantShields int
	}{
		{
			name: "no regen",
			args: args{
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					player:            testPlayer().WithNum(1),
					stackShields:      50,
					totalStackShields: 100,
				},
			},
			wantShields: 50,
		},
		{
			name: "regen",
			args: args{
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					player:            NewPlayer(1, NewRace().WithLRT(RS).WithSpec(&rules)).WithNum(1),
					stackShields:      50,
					totalStackShields: 100,
				},
			},
			wantShields: 60,
		},
		{
			name: "no regen when shields gone",
			args: args{
				token: battleToken{
					BattleRecordToken: BattleRecordToken{
						PlayerNum: 1,
					},
					player:            NewPlayer(1, NewRace().WithLRT(RS).WithSpec(&rules)).WithNum(1),
					stackShields:      0,
					totalStackShields: 100,
				},
			},
			wantShields: 0,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			tt.args.token.regenerateShields()
			got := tt.args.token.stackShields
			if got != tt.wantShields {
				t.Errorf("battle.regenerateShields() = %v, want %v", got, tt.wantShields)
			}

		})
	}
}

func Test_getCargoPerShip(t *testing.T) {
	type args struct {
		fleetCargo         int
		fleetCargoCapacity int
		tokenCargoCapacity int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"no cargo", args{fleetCargo: 0, fleetCargoCapacity: 0, tokenCargoCapacity: 0}, 0},
		{"full fleet one ship", args{fleetCargo: 10, fleetCargoCapacity: 10, tokenCargoCapacity: 10}, 10},
		{"fleet is full, has more cargo than a single ship fits", args{fleetCargo: 100, fleetCargoCapacity: 100, tokenCargoCapacity: 10}, 10},
		{"fleet is half full, ship should be half full", args{fleetCargo: 50, fleetCargoCapacity: 100, tokenCargoCapacity: 10}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCargoPerShip(tt.args.fleetCargo, tt.args.fleetCargoCapacity, tt.args.tokenCargoCapacity); got != tt.want {
				t.Errorf("getCargoPerShip() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_newBattleToken_orders verifies starting orders and starbase classification.
func Test_newBattleToken_orders(t *testing.T) {
	tests := []struct {
		name                       string
		hullType                   TechHullType
		armed                      bool
		wantTactic                 BattleTactic
		wantPrimary, wantSecondary BattleTarget
	}{
		{"armed ships retain orders", TechHullTypeFighter, true, BattleTacticMinimizeDamageToSelf, BattleTargetFreighters, BattleTargetArmedShips},
		{"unarmed ships retreat without a primary target", TechHullTypeFreighter, false, BattleTacticDisengage, BattleTargetNone, BattleTargetArmedShips},
		{"starbases attack any hostile ship", TechHullTypeStarbase, true, BattleTacticMaximizeDamage, BattleTargetAny, BattleTargetAny},
		{"orbital forts attack any hostile ship", TechHullTypeOrbitalFort, true, BattleTacticMaximizeDamage, BattleTargetAny, BattleTargetAny},
		{"unarmed starbases have no primary target", TechHullTypeStarbase, false, BattleTacticMaximizeDamage, BattleTargetNone, BattleTargetAny},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRulesWithSeed(1)
			design := &ShipDesign{Hull: SmallFreighter.Name, Spec: ShipDesignSpec{HullType: tt.hullType, HasWeapons: tt.armed}}
			token := &ShipToken{Quantity: 1, design: design}
			plan := BattlePlan{Tactic: BattleTacticMinimizeDamageToSelf, PrimaryTarget: BattleTargetFreighters, SecondaryTarget: BattleTargetArmedShips}
			got := newBattleToken(&r, 1, Vector{}, token, plan, testPlayer().WithNum(1))
			assert.Equal(t, tt.wantTactic, got.Tactic)
			assert.Equal(t, tt.wantPrimary, got.PrimaryTarget)
			assert.Equal(t, tt.wantSecondary, got.SecondaryTarget)
			assert.Equal(t, tt.hullType == TechHullTypeStarbase || tt.hullType == TechHullTypeOrbitalFort, got.isTargetOf(BattleTargetStarbase))
		})
	}
}

// testBattleToken returns a single ship with 100 armor that targets any ship and is
// hostile toward the attack players.
func testBattleToken(playerNum int, position Vector, tactic BattleTactic, attack ...int) *battleToken {
	token := &battleToken{
		BattleRecordToken: BattleRecordToken{PlayerNum: playerNum, Position: position, Tactic: tactic, PrimaryTarget: BattleTargetAny},
		ShipToken:         &ShipToken{Quantity: 1},
		armor:             100,
		attackPlayers:     map[int]bool{},
	}
	for _, playerNum := range attack {
		token.attackPlayers[playerNum] = true
	}
	return token
}

// withTestBeam arms the token with a slot holding one beam weapon.
func (token *battleToken) withTestBeam(power, weaponRange int) *battleToken {
	token.attributes |= battleTokenAttributeArmed
	token.weaponSlots = append(token.weaponSlots, &battleWeaponSlot{token: token, weaponType: battleWeaponTypeBeam, power: power, beamBonus: 1, slotQuantity: 1, weaponRange: weaponRange})
	return token
}
