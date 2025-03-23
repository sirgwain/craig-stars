package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMineField_getDecayRate(t *testing.T) {

	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)

	type fields struct {
		MineFieldType MineFieldType
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
		{"1000 mines, no planets", fields{MineFieldTypeStandard, 1000}, args{rules: &rules, player: player, numPlanets: 0}, 20},
		{"100 mines, min decay 10", fields{MineFieldTypeStandard, 100}, args{rules: &rules, player: player, numPlanets: 0}, 10},
		{"1000 mines, decay 4% per planet", fields{MineFieldTypeStandard, 1000}, args{rules: &rules, player: player, numPlanets: 1}, 60},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mineField := newMineField(tt.args.player, tt.fields.MineFieldType, tt.fields.NumMines, 1, Vector{})
			mineField.Spec = computeMinefieldSpec(tt.args.rules, tt.args.player, mineField, tt.args.numPlanets)
			if got := mineField.getDecayRate(tt.args.rules, tt.args.player, tt.args.numPlanets); got != tt.want {
				t.Errorf("MineField.getDecayRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMineField_reduceMineFieldOnImpact(t *testing.T) {
	tests := []struct {
		name      string
		numMines  int
		numTokens int
		want      int
	}{
		{"10 mines, field cleared", 10, 1, 0},
		{"1000 mines, 9.75% clear", 1000, 2, 902},
		{"5K mines, 50 tokens", 5000, 50, 2500},
		{"reduces down tfo 5K", 20000, 30, 4656},
		// 28 hits reduce field down to 4756 (2000*0.95^28); last 2 knock off extra 100
		{"163 hits clear field", 20000, 163, 0},
		// 28 hits reduce field down to 4756; 76 hits reduce down to 956
		// 31 hits reduce down to 195 mines; last 20 clear the field
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mineField := &MineField{NumMines: tt.numMines}
			mineField.reduceMineFieldOnImpact(tt.numTokens)
			if mineField.NumMines != tt.want {
				t.Errorf("MineField.reduceMineFieldOnImpact() produced mine count %v, want %v", mineField.NumMines, tt.want)
			}
		})
	}
}

func Test_checkForMineFieldCollision_Hit(t *testing.T) {
	// make a new fleet at -15x, and move it through the field
	fleetPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	fleet := testLongRangeScout(fleetPlayer)
	fleet.Position = Vector{-15, 0}

	radius := 10
	mineFieldPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	mineField := newMineField(mineFieldPlayer, MineFieldTypeStandard, radius*radius, 1, Vector{})
	mineField.Spec = computeMinefieldSpec(&rules, mineFieldPlayer, mineField, 0)

	u := &Universe{
		MineFields: []*MineField{mineField},
	}

	// make the speed minefield allow speed 5, 25% hit chance per warp
	// we'll go warp 9 to guarantee a hit
	rules := NewRules()
	stats := MineFieldStats{
		MaxSpeed:    5,
		ChanceOfHit: .25,
		// leave damage stuff the same as a standard minefield
		MinDamagePerFleetRS: 600,
		DamagePerEngineRS:   125,
		MinDamagePerFleet:   500,
		DamagePerEngine:     100,
	}
	rules.MineFieldStatsByType[MineFieldTypeStandard] = stats

	// send the fleet at warp 9, straight through the minefield
	dest := NewPositionWaypoint(Vector{20, 0}, 9)
	dist := float64(dest.WarpSpeed * dest.WarpSpeed)

	mineFieldHit, actualDist := checkForMineFieldCollision(&rules, newTestPlayerGetter(fleetPlayer, mineFieldPlayer), u, fleet, dest, dist)

	// we should come to a dead stop, ship destroyed
	assert.Equal(t, 5.0, actualDist)
	assert.Equal(t, mineField, mineFieldHit)

}

func Test_checkForMineFieldCollision_Miss(t *testing.T) {
	// make a new fleet at -15x, and move it through the field
	fleetPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	fleet := testLongRangeScout(fleetPlayer)
	fleet.Position = Vector{-5, 0}

	radius := 10
	mineFieldPlayer := NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	mineField := newMineField(mineFieldPlayer, MineFieldTypeStandard, radius*radius, 1, Vector{})
	mineField.Spec = computeMinefieldSpec(&rules, mineFieldPlayer, mineField, 0)

	u := &Universe{
		MineFields: []*MineField{mineField},
	}

	// send the fleet at warp 4, straight through the minefield, should be safe at warp 4
	dest := NewPositionWaypoint(Vector{20, 0}, 4)
	dist := float64(dest.WarpSpeed * dest.WarpSpeed)

	mineFieldHit, actualDist := checkForMineFieldCollision(&rules, newTestPlayerGetter(fleetPlayer, mineFieldPlayer), u, fleet, dest, dist)

	// we should come to a dead stop, ship destroyed
	assert.Nil(t, mineFieldHit)
	assert.Equal(t, 16.0, actualDist)
}

func TestMineField_moveTowardsMineLayer(t *testing.T) {
	player := NewPlayer(0, NewRace()).WithNum(1)
	type fields struct {
		mineFieldPosition Vector
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
			fields: fields{mineFieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{}, minesLaid: 1000},
			want:   Vector{},
		},
		{
			name:   "move towards fleet completely",
			fields: fields{mineFieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{5, 5}, minesLaid: 1000},
			want:   Vector{5, 5},
		},
		{
			name:   "move towards fleet halfway",
			fields: fields{mineFieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{6, 6}, minesLaid: 500},
			want:   Vector{3, 3},
		},
		{
			name:   "move towards fleet halfway round up",
			fields: fields{mineFieldPosition: Vector{}, numMines: 1000},
			args:   args{position: Vector{5, 5}, minesLaid: 500},
			want:   Vector{3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mineField := newMineField(player, MineFieldTypeStandard, tt.fields.numMines, 1, tt.fields.mineFieldPosition)
			mineField.moveTowardsMineLayer(tt.args.position, tt.args.minesLaid)

			if mineField.Position != tt.want {
				t.Errorf("MineField.moveTowardsMineLayer() = %v, want %v", mineField.Position, tt.want)
			}

		})
	}
}

func TestMineField_damageFleet(t *testing.T) {
	mineFieldPlayer := testPlayer().WithNum(1)
	fleetPlayer := testPlayer().WithNum(2)
	type args struct {
		fleet       *Fleet
		fleetPlayer *Player
		stats       MineFieldStats
	}
	tests := []struct {
		name     string
		detonate bool
		args     args
		want     MineFieldDamage
	}{
		{
			name:     "destroy scout",
			detonate: false,
			args:     args{testLongRangeScout(fleetPlayer), fleetPlayer, rules.MineFieldStatsByType[MineFieldTypeStandard]},
			want:     MineFieldDamage{Damage: 500, ShipsDestroyed: 1, FleetDestroyed: true},
		},
		{
			name:     "destroy big fleet",
			detonate: false,
			args:     args{testLongRangeScoutWithQuantity(fleetPlayer, 100), fleetPlayer, rules.MineFieldStatsByType[MineFieldTypeStandard]},
			want:     MineFieldDamage{Damage: 10000, ShipsDestroyed: 100, FleetDestroyed: true},
		},
		{
			name:     "damage stalwart defender",
			detonate: false,
			args:     args{testStalwartDefenderWithQuantity(fleetPlayer, 10), fleetPlayer, rules.MineFieldStatsByType[MineFieldTypeStandard]},
			want:     MineFieldDamage{Damage: 1000, ShipsDestroyed: 3, FleetDestroyed: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mineField := newMineField(mineFieldPlayer, MineFieldTypeStandard, 1000, 1, Vector{})
			if got := mineField.damageFleet(tt.args.fleet, tt.args.fleetPlayer, tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MineField.damageFleet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMineField_sweep(t *testing.T) {
	tests := []struct {
		name          string
		numMines      int
		fleetPosition Vector
		mineSweep     int
		fieldType     MineFieldType
		want          int
	}{
		{
			name:          "sweep mines in center",
			numMines:      100,
			fleetPosition: Vector{},
			mineSweep:     10,
			fieldType:     MineFieldTypeStandard,
			want:          10,
		},
		{
			name:          "sweep all mines",
			numMines:      100,
			fleetPosition: Vector{},
			mineSweep:     1000,
			fieldType:     MineFieldTypeStandard,
			want:          100,
		},
		{
			name:          "sweep to current position",
			numMines:      100,
			fleetPosition: Vector{9, 0},
			mineSweep:     1000,
			fieldType:     MineFieldTypeStandard,
			want:          19, // 100-81
		},
		{
			name:          "speed bump harder to sweep",
			numMines:      10000,
			fleetPosition: Vector{},
			mineSweep:     3000,
			fieldType:     MineFieldTypeSpeedBump,
			want:          1000, // speed bumps take 3x longer to sweep
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mineField := newMineField(testPlayer(), tt.fieldType, tt.numMines, 1, Vector{})
			if got := mineField.sweep(tt.fleetPosition, tt.mineSweep,
				rules.MineFieldStatsByType[tt.fieldType].SweepFactor); got != tt.want {
				t.Errorf("MineField.sweep() = %v, want %v", got, tt.want)
			}
		})
	}
}
