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
	}{
		{name: "move simple", fields: fields{7, Vector{}, Vector{100, 0}}, wantPosition: Vector{49, 0}},
		{name: "move done", fields: fields{7, Vector{}, Vector{45, 0}}, wantPosition: Vector{45, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := newMysteryTrader(tt.fields.position, 1, tt.fields.warpSpeed, tt.fields.destination, 0, MysteryTraderRewardNone)
			mt.move()

			assert.Equal(t, tt.wantPosition, mt.Position)
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
		game  *Game
		fleet *Fleet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MysteryTraderReward
	}{
		{"meet without cargo", fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}}, args{&rules, &Game{Year: 2500}, testLongRangeScout(player)}, MysteryTraderReward{}},
		{
			name:   "meet with cargo, research bonus",
			fields: fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type:       MysteryTraderRewardResearch,
				TechLevels: TechLevel{Energy: 6}, // rng will just pick energy levels each time...
			},
		},
		{
			name:   "meet with more cargo, research bonus",
			fields: fields{5000, MysteryTraderRewardResearch, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 6200})},
			want: MysteryTraderReward{
				Type:       MysteryTraderRewardResearch,
				TechLevels: TechLevel{Energy: 7}, // bonus level for more minerals
			},
		},
		{
			name:   "meet with cargo, get genesis",
			fields: fields{5000, MysteryTraderRewardGenesis, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type: MysteryTraderRewardGenesis,
				Tech: GenesisDevice.Name,
			},
		},
		{
			name:   "meet with cargo, get jump gate",
			fields: fields{5000, MysteryTraderRewardJumpGate, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type: MysteryTraderRewardJumpGate,
				Tech: JumpGate.Name,
			},
		},
		{
			name:   "meet with cargo, get mechanical",
			fields: fields{5000, MysteryTraderRewardMechanical, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type: MysteryTraderRewardMechanical,
				Tech: MultiCargoPod.Name,
			},
		},
		{
			name:   "meet with cargo, get ship",
			fields: fields{5000, MysteryTraderRewardLifeboat, TechLevel{}, &testRandom{}},
			args:   args{&rules, &Game{Year: 2500}, testGalleon(player).withCargo(Cargo{Ironium: 5000})},
			want: MysteryTraderReward{
				Type:      MysteryTraderRewardLifeboat,
				Ship:      MysteryTraderScout,
				ShipCount: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player.TechLevels = tt.fields.PlayerTechLevels
			rulesCopy := rules
			rulesCopy.random = tt.fields.random

			mt := newMysteryTrader(Vector{20, 20}, 1, 7, Vector{380, 20}, tt.fields.RequestedBoon, tt.fields.Reward)
			if got := mt.meet(&rulesCopy, tt.args.game, tt.args.fleet, player); !test.CompareAsJSON(t, got, tt.want) {
				t.Errorf("MysteryTrader.meet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateMysteryTraderReward(t *testing.T) {
	type args struct {
		rng       testRandom
		year      int
		warpSpeed int
	}
	tests := []struct {
		name string
		args args
		want MysteryTraderRewardType
	}{
		{"all 0s", args{rng: *newIntRandom(), year: 2400, warpSpeed: 7}, MysteryTraderRewardResearch},
		{"lifeboat 1 in 6 chance", args{rng: *newIntRandom(0, 5), year: 2400, warpSpeed: 7}, MysteryTraderRewardLifeboat},
		{"research1", args{rng: *newIntRandom(0, 1), year: 2400, warpSpeed: 7}, MysteryTraderRewardResearch},
		{"research2", args{rng: *newIntRandom(6, 0), year: 2400, warpSpeed: 7}, MysteryTraderRewardResearch},
		{"engine", args{rng: *newIntRandom(6, 1), year: 2400, warpSpeed: 7}, MysteryTraderRewardEngine},
		{"engine more likely later", args{rng: *newIntRandom(5, 1), year: 2500, warpSpeed: 7}, MysteryTraderRewardEngine},
		{"engine more likely faster", args{rng: *newIntRandom(4, 1), year: 2500, warpSpeed: 11}, MysteryTraderRewardEngine},
		{"torpedos are hard", args{rng: *newIntRandom(6, 7, 7, 0), year: 2400, warpSpeed: 7}, MysteryTraderRewardTorpedo},
		{"beams are harder", args{rng: *newIntRandom(6, 8, 8, 0), year: 2400, warpSpeed: 7}, MysteryTraderRewardBeamWeapon},
		{"almost a beam, research instead", args{rng: *newIntRandom(6, 8, 8, 1), year: 2400, warpSpeed: 7}, MysteryTraderRewardResearch},
		{"late game jump gate", args{rng: *newIntRandom(6, 12, 12), year: 2580, warpSpeed: 7}, MysteryTraderRewardJumpGate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := NewRules()
			rules.random = &tt.args.rng
			if got := generateMysteryTraderReward(&rules, tt.args.year, tt.args.warpSpeed); got != tt.want {
				t.Errorf("generateMysteryTraderReward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRandomMysteryTraderTech(t *testing.T) {
	type args struct {
		rng    rng
		player *Player
		reward MysteryTraderRewardType
	}
	tests := []struct {
		name           string
		args           args
		want           *Tech
		wantRewardType MysteryTraderRewardType
	}{
		{"get engine", args{newIntRandom(), testPlayer(), MysteryTraderRewardEngine}, &EnigmaPulsar.Tech, MysteryTraderRewardEngine},
		{"get armor", args{newIntRandom(), testPlayer(), MysteryTraderRewardArmor}, &MegaPolyShell.Tech, MysteryTraderRewardArmor},
		{"get torpedo", args{newIntRandom(), testPlayer(), MysteryTraderRewardTorpedo}, &AntiMatterTorpedo.Tech, MysteryTraderRewardTorpedo},
		{"get minimorph", args{newIntRandom(), testPlayer(), MysteryTraderRewardShipHull}, &MiniMorph.Tech, MysteryTraderRewardShipHull},
		{"lifeboat is special, not ever awarded randomly by category", args{newIntRandom(), testPlayer(), MysteryTraderRewardLifeboat}, nil, MysteryTraderRewardLifeboat},
		{"genesis is special, not ever awarded randomly by category", args{newIntRandom(), testPlayer(), MysteryTraderRewardGenesis}, &GenesisDevice.Tech, MysteryTraderRewardGenesis},
		{"jumpgate is special, not ever awarded randomly by category", args{newIntRandom(), testPlayer(), MysteryTraderRewardJumpGate}, &JumpGate.Tech, MysteryTraderRewardJumpGate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotRewardType := getMysteryTraderPart(tt.args.rng, tt.args.player, tt.args.reward)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMysteryTraderTech() = %v, want %v", got, tt.want)
			}
			if gotRewardType != tt.wantRewardType {
				t.Errorf("getMysteryTraderTech() rewardType = %v, want rewardType %v", gotRewardType, tt.wantRewardType)
			}
		})
	}
}
