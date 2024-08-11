package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
	"github.com/stretchr/testify/assert"
)

func Test_generateMysteryTrader(t *testing.T) {
	type args struct {
		random *testRandom
		game   *Game
		num    int
	}
	tests := []struct {
		name string
		args args
		want *MysteryTrader
	}{
		{"no mystery trader, too early", args{&testRandom{}, &Game{Year: 2400}, 1}, nil},
		{"no mystery trader, odd year", args{&testRandom{}, &Game{Year: 2441}, 1}, nil},
		{"no mystery trader, random chance failed", args{newIntRandom(0, 1), &Game{Year: 2440}, 1}, nil},
		{
			"mystery trader, random always returns 0",
			args{newIntRandom(), &Game{Year: 2440}, 1},
			newMysteryTrader(Vector{X: -20, Y: 20}, 1, 7, Vector{X: 20, Y: -20}, 5000, MysteryTraderRewardResearch),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rulesCopy := rules
			rulesCopy.random = tt.args.random

			if got := generateMysteryTrader(&rulesCopy, tt.args.game, tt.args.num); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("generateMysteryTrader() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestMysteryTrader_move(t *testing.T) {
	type fields struct {
		warpSpeed   int
		position    Vector
		destination Vector
	}
	tests := []struct {
		name         string
		fields       fields
		wantPosition Vector
		wantDelete   bool
	}{
		{name: "move simple", fields: fields{7, Vector{}, Vector{100, 0}}, wantPosition: Vector{49, 0}, wantDelete: false},
		{name: "move done", fields: fields{7, Vector{}, Vector{45, 0}}, wantPosition: Vector{45, 0}, wantDelete: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newMysteryTrader(tt.fields.position, 1, tt.fields.warpSpeed, tt.fields.destination, 0, MysteryTraderRewardNone)
			mt.move()

			assert.Equal(t, tt.wantPosition, mt.Position)
			assert.Equal(t, tt.wantDelete, mt.Delete)
		})
	}
}

func TestMysteryTrader_meet(t *testing.T) {

	player := NewPlayer(1, NewRace())

	type fields struct {
		RequestedBoon    int
		Reward           MysteryTraderRewardType
		PlayerTechLevels TechLevel
		random           rng
	}
	type args struct {
		rules *Rules
		fleet *Fleet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MysteryTraderReward
	}{
		{"meet without cargo", fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}}, args{&rules, testLongRangeScout(player)}, MysteryTraderReward{}},
		{
			name:   "meet with cargo, research bonus",
			fields: fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}},
			args:   args{&rules, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type:       MysteryTraderRewardResearch,
				TechLevels: TechLevel{Energy: 6}, // rng will just pick energy levels each time...
			},
		},
		{
			name:   "meet with cargo, research bonus",
			fields: fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}},
			args:   args{&rules, testGalleon(player).withCargo(Cargo{Ironium: 6200})},
			want: MysteryTraderReward{
				Type:       MysteryTraderRewardResearch,
				TechLevels: TechLevel{Energy: 7}, // bonus level for more minerals
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.TechLevels = tt.fields.PlayerTechLevels
			rulesCopy := rules
			rulesCopy.random = tt.fields.random

			mt := newMysteryTrader(Vector{20, 20}, 1, 7, Vector{380, 20}, tt.fields.RequestedBoon, tt.fields.Reward)
			if got := mt.meet(&rulesCopy, tt.args.fleet, player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MysteryTrader.meet() = %v, want %v", got, tt.want)
			}
		})
	}
}
