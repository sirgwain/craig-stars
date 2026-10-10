//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// create a new mini bomber fleet for testing
func testMiniBomber(player *Player, bomb TechHullComponent) *Fleet {
	fleet := &Fleet{
		Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num,
		BaseName: "Mini Bomber",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player.Num, 1).
					WithHull(MiniBomber.Name).
					WithSlots([]ShipDesignSlot{
						{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
						{HullComponent: bomb.Name, HullSlotIndex: 2, Quantity: 2},
					}).
					WithSpec(&rules, player)},
		},
		battlePlan:        &player.BattlePlans[0],
		OrbitingPlanetNum: None,
	}
	fleet.Spec = ComputeFleetSpec(&rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func Test_bomber_getUnterraformAmount(t *testing.T) {
	type args struct {
		retroBombAmount int
		baseHab         Hab
		hab             Hab
	}
	tests := []struct {
		name string
		args args
		want Hab
	}{
		{
			name: "no unterraform, planet isn't terraformed",
			args: args{retroBombAmount: 10, baseHab: Hab{50, 50, 50}, hab: Hab{50, 50, 50}},
			want: Hab{0, 0, 0},
		},
		{
			name: "one retro strength removes one point from each axis",
			args: args{retroBombAmount: 1, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-1, 1, -1},
		},
		{
			name: "unterraform by one, even though we bombed with 10",
			args: args{retroBombAmount: 10, baseHab: Hab{50, 49, 50}, hab: Hab{50, 50, 50}},
			want: Hab{0, -1, 0},
		},
		{
			name: "thirty retro strength restores all three axes",
			args: args{retroBombAmount: 30, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-30, 20, -10},
		},
		{
			name: "extra retro strength cannot overshoot base hab",
			args: args{retroBombAmount: 36, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-30, 20, -10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomber{
				rules: &rules,
				log:   testLogger,
			}
			if got := b.getUnterraformAmount(tt.args.retroBombAmount, tt.args.baseHab, tt.args.hab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bomb.getUnterraformAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bomber_bombPlanet(t *testing.T) {
	fleetOwner := testPlayer().WithNum(1)
	planetOwner := testPlayer().WithNum(2)
	pg := newTestPlayerGetter(fleetOwner, planetOwner)

	type want struct {
		population int
		mines      int
		factories  int
		defenses   int
		hab        Hab
	}
	type args struct {
		planet       *Planet
		enemyBombers []*Fleet
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Mini bomber, uses min kill rate",
			args: args{
				planet:       &Planet{PlayerNum: planetOwner.Num, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}},
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb)},
			},
			want: want{population: 9500, mines: 99, factories: 98, defenses: 9, hab: Hab{50, 50, 50}},
		},
		{
			name: "Cherry bombers nuke planet; reset partial pop",
			args: args{
				planet:       NewPlanet().WithPopulation(255),
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, CherryBomb)},
			},
			want: want{population: 0},
		},
		{
			name: "Two mini bombers, one with smart bombs",
			args: args{
				planet:       &Planet{PlayerNum: planetOwner.Num, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}},
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb), testMiniBomber(fleetOwner, SmartBomb)},
			},
			want: want{population: 9500, mines: 99, factories: 98, defenses: 9, hab: Hab{50, 50, 50}},
		},
		{
			name: "Three mini bombers, one with smart bombs, one with retro bombs",
			args: args{
				planet:       &Planet{PlayerNum: planetOwner.Num, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}, BaseHab: Hab{50, 49, 50}},
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb), testMiniBomber(fleetOwner, SmartBomb), testMiniBomber(fleetOwner, RetroBomb)},
			},
			want: want{population: 9500, mines: 99, factories: 98, defenses: 9, hab: Hab{50, 49, 50}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rulesCopy := rules
			rulesCopy.random = &testRandom{}
			b := &bomber{
				rules: &rulesCopy,
				log:   testLogger,
			}
			tt.args.planet.Spec = ComputePlanetSpec(&rules, planetOwner, tt.args.planet)
			b.bombPlanet(tt.args.planet, planetOwner, tt.args.enemyBombers, pg)

			got := want{
				population: tt.args.planet.GetPopulation(),
				mines:      tt.args.planet.Mines,
				factories:  tt.args.planet.Factories,
				defenses:   tt.args.planet.Defenses,
				hab:        tt.args.planet.Hab,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bomb.bombPlanet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bomber_bombPlayerCombined(t *testing.T) {
	r := NewRulesWithSeed(0)
	r.random = &testRandom{}
	defender := NewPlayer(2, NewRace().WithSpec(&r)).withSpec(&r)
	p := NewPlanet().WithPlayerNum(2).WithPopulation(100_000)
	p.Mines, p.Factories, p.Defenses = 100, 100, 100
	p.Hab, p.BaseHab = Hab{60, 60, 60}, Hab{50, 50, 50}
	p.Spec.DefenseCoverage = .5
	p.Spec.DefenseCoverageSmart = .25
	f := &Fleet{}
	f.Spec.Bombs = []Bomb{{Quantity: 1, KillRate: 20, MinKillRate: 300, StructureDestroyRate: 10}}
	f.Spec.SmartBombs = []Bomb{{Quantity: 1, KillRate: 50}}
	f.Spec.RetroBombs = []Bomb{{Quantity: 1, UnterraformRate: 4}}
	b := newBomber(testLogger, &r)

	result := b.bombPlayer(p, defender, []*Fleet{f})

	// smart bombs kill 375kT, then normal bombs kill 63kT of the remaining 625kT
	assert.Equal(t, 43800, result.ColonistsKilled)
	assert.Equal(t, 8, result.MinesDestroyed+result.FactoriesDestroyed+result.DefensesDestroyed)
	// retro bombs undo terraforming on every axis
	assert.Equal(t, Hab{57, 57, 57}, p.Hab)
	assert.Equal(t, p.Hab.Subtract(p.BaseHab), p.TerraformedAmount)
}
