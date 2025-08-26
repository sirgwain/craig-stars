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
		{"remove 5% from big field", 10_000, 9500},
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

func Test_checkForMinefieldCollision_Hit(t *testing.T) {
	// make a new fleet at -15x, and move it through the field
	fleetPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	fleet := testLongRangeScout(fleetPlayer)
	fleet.Position = Vector{-15, 0}

	radius := 10
	minefieldPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	minefield := newMinefield(minefieldPlayer, MinefieldTypeStandard, radius*radius, 1, Vector{})

	u := &Universe{
		Minefields: []*Minefield{minefield},
	}

	// make the speed minefield allow speed 5, 25% hit chance per warp
	// we'll go warp 9 to guarantee a hit
	rules := NewRules()
	stats := MinefieldStats{
		MaxSpeed:    5,
		ChanceOfHit: .25,
		// leave damage stuff the same as a standard minefield
		MinDamagePerFleetRS: 600,
		DamagePerEngineRS:   125,
		MinDamagePerFleet:   500,
		DamagePerEngine:     100,
	}
	rules.MinefieldStatsByType[MinefieldTypeStandard] = stats

	// send the fleet at warp 9, straight through the minefield
	dest := NewPositionWaypoint(Vector{20, 0}, 9)
	dist := float64(dest.WarpSpeed * dest.WarpSpeed)

	minefieldHit, actualDist := checkForMinefieldCollision(&rules, newTestPlayerGetter(fleetPlayer, minefieldPlayer), u, fleet, dest, dist)

	// we should come to a dead stop, ship destroyed
	assert.Equal(t, 5.0, actualDist)
	assert.Equal(t, minefield, minefieldHit)

}

func Test_checkForMinefieldCollision_Miss(t *testing.T) {
	// make a new fleet at -15x, and move it through the field
	fleetPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	fleet := testLongRangeScout(fleetPlayer)
	fleet.Position = Vector{-5, 0}

	radius := 10
	minefieldPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	minefield := newMinefield(minefieldPlayer, MinefieldTypeStandard, radius*radius, 1, Vector{})

	u := &Universe{
		Minefields: []*Minefield{minefield},
	}

	// send the fleet at warp 4, straight through the minefield, should be safe at warp 4
	dest := NewPositionWaypoint(Vector{20, 0}, 4)
	dist := float64(dest.WarpSpeed * dest.WarpSpeed)

	minefieldHit, actualDist := checkForMinefieldCollision(&rules, newTestPlayerGetter(fleetPlayer, minefieldPlayer), u, fleet, dest, dist)

	// we should come to a dead stop, ship destroyed
	assert.Nil(t, minefieldHit)
	assert.Equal(t, 16.0, actualDist)
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
			want:     MinefieldDamage{Damage: 1000, ShipsDestroyed: 3, FleetDestroyed: false},
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
