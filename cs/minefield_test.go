//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinefield_getDecayRate(t *testing.T) {

	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)

	type fields struct {
		MinefieldType MinefieldType
		NumMines      int
	}
	type args struct {
		rules      *Rules
		player     *Player
		numPlanets int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"1000 mines, no planets", fields{MinefieldTypeStandard, 1000}, args{rules: &rules, player: player, numPlanets: 0}, 20},
		{"100 mines, min decay 10", fields{MinefieldTypeStandard, 100}, args{rules: &rules, player: player, numPlanets: 0}, 10},
		{"1000 mines, decay 4% per planet", fields{MinefieldTypeStandard, 1000}, args{rules: &rules, player: player, numPlanets: 1}, 60},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minefield := newMinefield(tt.args.player, tt.fields.MinefieldType, tt.fields.NumMines, 1, Vector{})
			if got := minefield.getDecayRate(tt.args.rules, tt.args.player, tt.args.numPlanets); got != tt.want {
				t.Errorf("Minefield.getDecayRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinefield_reduceMinefieldOnImpact(t *testing.T) {
	tests := []struct {
		name     string
		numMines int
		want     int
	}{
		{"remove all", 10, 0},
		{"remove min", 20, 10},
		{"remove 5% from small field", 500, 475},
		{"remove 50 from medium field", 5000, 4950},
		{"remove 1% from big field", 10_000, 9900},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minefield := &Minefield{
				NumMines: tt.numMines,
			}

			minefield.reduceMinefieldOnImpact()
			if minefield.NumMines != tt.want {
				t.Errorf("Minefield.reduceMinefieldOnImpact() = %v, want %v", minefield.NumMines, tt.want)
			}
		})
	}
}

func Test_checkForMinefieldCollision(t *testing.T) {
	misses := func(n int) []float64 {
		rolls := make([]float64, n)
		for i := range rolls {
			rolls[i] = .99
		}
		return rolls
	}
	type field struct {
		position Vector
		radius   int
	}
	tests := []struct {
		name     string
		position Vector
		dest     Vector
		warp     int
		fields   []field
		rolls    []float64 // rolls after these always hit
		wantHit  int       // index of the field hit, or -1
		wantDist float64
	}{
		{
			name:     "hit where the fleet enters the field",
			position: Vector{-15, 0}, dest: Vector{100, 0}, warp: 9,
			fields:  []field{{Vector{0, 0}, 10}},
			wantHit: 0, wantDist: 5,
		},
		{
			name:     "roll once per ly in the field",
			position: Vector{-15, 0}, dest: Vector{100, 0}, warp: 9,
			fields:  []field{{Vector{0, 0}, 10}},
			rolls:   misses(10),
			wantHit: 0, wantDist: 15,
		},
		{
			name:     "roll for the whole path through the field, not its radius",
			position: Vector{-15, 0}, dest: Vector{100, 0}, warp: 9,
			fields:  []field{{Vector{0, 0}, 10}},
			rolls:   misses(19),
			wantHit: 0, wantDist: 24,
		},
		{
			name:     "overlapping fields are rolled for once",
			position: Vector{-15, 0}, dest: Vector{100, 0}, warp: 9,
			fields:  []field{{Vector{0, 0}, 10}, {Vector{5, 0}, 10}},
			rolls:   misses(25), // 25ly from -10 to 15
			wantHit: -1, wantDist: 81,
		},
		{
			name:     "fields are checked in path order",
			position: Vector{-15, 0}, dest: Vector{100, 0}, warp: 9,
			fields:  []field{{Vector{50, 0}, 5}, {Vector{0, 0}, 5}},
			wantHit: 1, wantDist: 10,
		},
		{
			name:     "safe at warp 4",
			position: Vector{-5, 0}, dest: Vector{20, 0}, warp: 4,
			fields:  []field{{Vector{0, 0}, 10}},
			wantHit: -1, wantDist: 16,
		},
		{
			// warp 9, but the fleet only goes 16ly, so it's checked at warp 4
			name:     "slow arrival through a field is safe",
			position: Vector{-5, 0}, dest: Vector{12, 0}, warp: 9,
			fields:  []field{{Vector{0, 0}, 10}},
			wantHit: -1, wantDist: 17,
		},
		{
			name:     "last ly before arriving isn't checked",
			position: Vector{-15, 0}, dest: Vector{6, 0}, warp: 9,
			fields:  []field{{Vector{8, 0}, 3}},
			wantHit: -1, wantDist: 21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fleetPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
			minefieldPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
			fleet := testLongRangeScout(fleetPlayer)
			fleet.Position = tt.position

			minefields := []*Minefield{}
			for i, f := range tt.fields {
				minefields = append(minefields, newMinefield(minefieldPlayer, MinefieldTypeStandard, f.radius*f.radius, i+1, f.position))
			}
			u := &Universe{Minefields: minefields}

			rules := NewRulesWithSeed(0)
			rules.random = newFloat64Random(tt.rolls...)
			dest := NewPositionWaypoint(tt.dest, tt.warp)
			dist := min(float64(tt.warp*tt.warp), tt.position.DistanceTo(tt.dest))

			minefieldHit, actualDist := checkForMinefieldCollision(&rules, newTestPlayerGetter(fleetPlayer, minefieldPlayer), u, fleet, dest, dist)

			if tt.wantHit == -1 {
				assert.Nil(t, minefieldHit)
			} else {
				assert.Equal(t, minefields[tt.wantHit], minefieldHit)
			}
			assert.Equal(t, tt.wantDist, actualDist)
		})
	}
}

func TestMinefield_moveTowardsMineLayer(t *testing.T) {
	player := NewPlayer(0, NewRace()).WithNum(1)
	type fields struct {
		minefieldPosition Vector
		numMines          int
	}
	type args struct {
		position  Vector
		minesLaid int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Vector
	}{
		{
			name:   "no movement",
			fields: fields{minefieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{}, minesLaid: 1000},
			want:   Vector{},
		},
		{
			name:   "move towards fleet completely",
			fields: fields{minefieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{5, 5}, minesLaid: 1000},
			want:   Vector{5, 5},
		},
		{
			name:   "move towards fleet halfway",
			fields: fields{minefieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{6, 6}, minesLaid: 500},
			want:   Vector{3, 3},
		},
		{
			name:   "move towards fleet halfway round up",
			fields: fields{minefieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{5, 5}, minesLaid: 500},
			want:   Vector{3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			minefield := newMinefield(player, MinefieldTypeStandard, tt.fields.numMines, 1, tt.fields.minefieldPosition)
			minefield.moveTowardsMineLayer(tt.args.position, tt.args.minesLaid)

			if minefield.Position != tt.want {
				t.Errorf("Minefield.moveTowardsMineLayer() = %v, want %v", minefield.Position, tt.want)
			}

		})
	}
}

func TestMinefield_damageFleet(t *testing.T) {
	minefieldPlayer := testPlayer().WithNum(1)
	fleetPlayer := testPlayer().WithNum(2)
	type args struct {
		fleet       *Fleet
		fleetPlayer *Player
		stats       MinefieldStats
	}
	tests := []struct {
		name     string
		detonate bool
		args     args
		want     MinefieldDamage
	}{
		{
			name:     "destroy scout",
			detonate: false,
			args:     args{testLongRangeScout(fleetPlayer), fleetPlayer, rules.MinefieldStatsByType[MinefieldTypeStandard]},
			want:     MinefieldDamage{Damage: 500, ShipsDestroyed: 1, FleetDestroyed: true},
		},
		{
			name:     "destroy big fleet",
			detonate: false,
			args:     args{testLongRangeScoutWithQuantity(fleetPlayer, 100), fleetPlayer, rules.MinefieldStatsByType[MinefieldTypeStandard]},
			want:     MinefieldDamage{Damage: 10000, ShipsDestroyed: 100, FleetDestroyed: true},
		},
		{
			name:     "damage stalwart defender",
			detonate: false,
			args:     args{testStalwartDefenderWithQuantity(fleetPlayer, 10), fleetPlayer, rules.MinefieldStatsByType[MinefieldTypeStandard]},
			want:     MinefieldDamage{Damage: 1000, ShipsDestroyed: 0, FleetDestroyed: false}, // spread over all 10 ships
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minefield := newMinefield(minefieldPlayer, MinefieldTypeStandard, 1000, 1, Vector{})
			if got := minefield.damageFleet(tt.args.fleet, tt.args.fleetPlayer, tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Minefield.damageFleet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinefield_sweep(t *testing.T) {
	type fields struct {
		minefieldPosition Vector
		numMines          int
	}
	type args struct {
		fleetPosition Vector
		mineSweep     int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"sweep mines in center", fields{minefieldPosition: Vector{}, numMines: 100}, args{fleetPosition: Vector{}, mineSweep: 10}, 10},
		{"sweep all mines", fields{minefieldPosition: Vector{}, numMines: 100}, args{fleetPosition: Vector{}, mineSweep: 1000}, 100},
		{
			"radius 10 minefield, we are 9 away so we can sweep it down to a 9ly mf (81 mines)",
			fields{minefieldPosition: Vector{}, numMines: 100},
			args{fleetPosition: Vector{9, 0}, mineSweep: 1000},
			100 - 81, // sweep from 10 radius to 9 radius so we are just on the edge
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minefield := newMinefield(testPlayer(), MinefieldTypeStandard, tt.fields.numMines, 1, tt.fields.minefieldPosition)
			if got := minefield.sweep(&rules, tt.args.fleetPosition, tt.args.mineSweep); got != tt.want {
				t.Errorf("Minefield.sweep() = %v, want %v", got, tt.want)
			}
		})
	}
}
