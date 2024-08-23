package cs

import (
	"testing"
)

func Test_battleToken_getDistanceAway(t *testing.T) {

	type args struct {
		position BattleVector
	}
	tests := []struct {
		name string
		bt   battleToken
		args args
		want int
	}{
		{"no distance", battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{0, 0}}}, args{BattleVector{0, 0}}, 0},
		{"x distance greatest", battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{2, 1}}}, args{BattleVector{4, 2}}, 2},
		{"y distance greatest", battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{1, 2}}}, args{BattleVector{2, 5}}, 3},
		{"negative distance (token behind)", battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{1, 1}}}, args{BattleVector{0, 0}}, 1},
		{"3,4 to 7,4", battleToken{BattleRecordToken: BattleRecordToken{Position: BattleVector{3, 4}}}, args{BattleVector{7, 4}}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.getDistanceAway(tt.args.position); got != tt.want {
				t.Errorf("battleToken.getDistanceAway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_battleToken_willTarget(t *testing.T) {

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
				t.Errorf("battle.willTarget() = %v, want %v", got, tt.want)
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
