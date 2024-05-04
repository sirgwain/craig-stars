package cs

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/test"
)

// create a new long rang scout fleet for testing
func testMiniBomber(player *Player, bomb TechHullComponent) *Fleet {
	fleet := &Fleet{
		MapObject: MapObject{Type: MapObjectTypeFleet, Num: 1, PlayerNum: player.Num},
		BaseName:  "Mini Bomber",
		Tokens: []ShipToken{
			{
				Quantity:  1,
				DesignNum: 1,
				design: NewShipDesign(player, 1).
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

func Test_bomb_getColonistsKilledForBombs(t *testing.T) {
	type args struct {
		population      int
		defenseCoverage float64
		bombs           []Bomb
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "No bombs",
			args: args{
				population:      1000,
				defenseCoverage: 0.5,
				bombs:           []Bomb{},
			},
			want: 0.0,
		},
		{
			name: "One bomb, no defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
				},
			},
			want: 2500,
		},
		{
			name: "One bomb, partial defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.9792,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
				},
			},
			want: 52.0,
		},
		{
			name: "Multiple bombs, no defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
					{Quantity: 5, KillRate: 1.2},
				},
			},
			want: 3100.0,
		},
		{
			name: "Multiple bombs, partial defense",
			args: args{
				population:      10000,
				defenseCoverage: 0.9792,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 2.5},
					{Quantity: 5, KillRate: 1.2},
				},
			},
			want: 64.48,
		},
		{
			name: "One bomb, full defense",
			args: args{
				population:      1000,
				defenseCoverage: 1.0,
				bombs: []Bomb{
					{Quantity: 1, KillRate: 50.0},
				},
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getColonistsKilledForBombs(tt.args.population, tt.args.defenseCoverage, tt.args.bombs); !test.WithinTolerance(got, tt.want, .1) {
				t.Errorf("bomb.getColonistsKilledForBombs() = %v, want %v", int(got), int(tt.want))
			}
		})
	}
}

func Test_bomb_getColonistsKilledWithSmartBombs(t *testing.T) {
	type args struct {
		population           int
		defenseCoverageSmart float64
		bombs                []Bomb
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "No bombs",
			args: args{
				population: 10000,
				bombs:      []Bomb{},
			},
			want: 0.0,
		},
		{
			name: "Multiple bombs, high defense",
			args: args{
				population:           10000,
				defenseCoverageSmart: 0.8524,
				bombs: []Bomb{
					{Quantity: 10, KillRate: 7},
					{Quantity: 5, KillRate: 2.2},
				},
			},
			want: 837,
		},
		{
			name: "Many smart bombs, low pop",
			args: args{
				population: 1000,
				bombs: []Bomb{
					{Quantity: 17 * 4, KillRate: 2.2},
					{Quantity: 17 * 4, KillRate: 2.2},
				},
			},
			want: 1000,
		},
		{
			name: "Smart bombs, very low pop",
			args: args{
				population: 500,
				bombs: []Bomb{
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
				},
			},
			want: 49,
		},
		{
			name: "1 fleet of 5 B-17 bombers with smart bombs, very low pop",
			args: args{
				population: 500,
				bombs: []Bomb{
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
					{Quantity: 4, KillRate: 1.3},
				},
			},
			want: 203,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getColonistsKilledWithSmartBombs(tt.args.population, tt.args.defenseCoverageSmart, tt.args.bombs); !test.WithinTolerance(got, tt.want, .1) {
				t.Errorf("bomb.getColonistsKilledWithSmartBombs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bomb_getStructuresDestroyed(t *testing.T) {
	type args struct {
		defenseCoverage float64
		bombs           []Bomb
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "No bombs",
			args: args{
				defenseCoverage: 0.5,
				bombs:           []Bomb{},
			},
			want: 0.0,
		},
		{
			name: "One bomb, no defense",
			args: args{
				defenseCoverage: 0.0,
				bombs: []Bomb{
					{Quantity: 10, StructureDestroyRate: 1},
				},
			},
			want: 100,
		},
		{
			name: "Many bombs, good defense",
			args: args{
				defenseCoverage: .9792,
				bombs: []Bomb{
					{Quantity: 10, StructureDestroyRate: 1},
					{Quantity: 5, StructureDestroyRate: .6},
				},
			},
			want: 66,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getStructuresDestroyed(tt.args.defenseCoverage, tt.args.bombs); got != tt.want {
				t.Errorf("bomb.getStructuresDestroyed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bomb_getUnterraformAmount(t *testing.T) {
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
			name: "should unterraform 1 gravity because it's 30 points from the base",
			args: args{retroBombAmount: 1, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-1, 0, 0},
		},
		{
			name: "unterraform by one, even though we bombed with 10",
			args: args{retroBombAmount: 10, baseHab: Hab{50, 49, 50}, hab: Hab{50, 50, 50}},
			want: Hab{0, -1, 0},
		},
		{
			name: "if we bomb for -30 terraform points, we should equalize between the two highest",
			args: args{retroBombAmount: 30, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-20, 10, 0},
		},
		{
			name: "if we bomb for -36 terraform points, we should equalize between all three",
			args: args{retroBombAmount: 36, baseHab: Hab{20, 70, 40}, hab: Hab{50, 50, 50}},
			want: Hab{-22, 12, -2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			if got := b.getUnterraformAmount(tt.args.retroBombAmount, tt.args.baseHab, tt.args.hab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bomb.getUnterraformAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bomb_bombPlanet(t *testing.T) {
	fleetOwner := testPlayer().WithNum(1)
	planetOwner := testPlayer().WithNum(2)

	type want struct {
		population int
		mines      int
		factories  int
		defenses   int
		hab        Hab
	}
	type args struct {
		planet       *Planet
		planetOwner  *Player
		enemyBombers []*Fleet
		pg           playerGetter
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Mini bomber, 10k planet, 10 defenses, 100 mines and factories, uses min kill rate",
			args: args{
				planet:       &Planet{MapObject: MapObject{PlayerNum: planetOwner.Num}, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}},
				planetOwner:  planetOwner,
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb)},
				pg:           newTestPlayerGetter(fleetOwner, planetOwner),
			},
			want: want{population: 9500, mines: 98, factories: 98, defenses: 9, hab: Hab{50, 50, 50}},
		},
		{
			name: "Two mini bombers, one with smart bombs",
			args: args{
				planet:       &Planet{MapObject: MapObject{PlayerNum: planetOwner.Num}, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}},
				planetOwner:  planetOwner,
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb), testMiniBomber(fleetOwner, SmartBomb)},
				pg:           newTestPlayerGetter(fleetOwner, planetOwner),
			},
			want: want{population: 9300, mines: 98, factories: 98, defenses: 9, hab: Hab{50, 50, 50}},
		},
		{
			name: "Three mini bombers, one with smart bombs, one with retro bombs",
			args: args{
				planet:       &Planet{MapObject: MapObject{PlayerNum: planetOwner.Num}, Mines: 100, Factories: 100, Defenses: 10, Cargo: Cargo{Colonists: 100}, Hab: Hab{50, 50, 50}, BaseHab: Hab{50, 49, 50}},
				planetOwner:  planetOwner,
				enemyBombers: []*Fleet{testMiniBomber(fleetOwner, LadyFingerBomb), testMiniBomber(fleetOwner, SmartBomb), testMiniBomber(fleetOwner, RetroBomb)},
				pg:           newTestPlayerGetter(fleetOwner, planetOwner),
			},
			want: want{population: 9300, mines: 98, factories: 98, defenses: 9, hab: Hab{50, 49, 50}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bomb{
				rules: &rules,
			}
			tt.args.planet.Spec = computePlanetSpec(&rules, tt.args.planetOwner, tt.args.planet)
			b.bombPlanet(tt.args.planet, tt.args.planetOwner, tt.args.enemyBombers, tt.args.pg)

			got := want{
				population: tt.args.planet.population(),
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
