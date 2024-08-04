package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMineralPacket_movePacket(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)

	type fields struct {
		Cargo         Cargo
		WarpSpeed     int
		builtThisTurn bool
	}
	type args struct {
		rules        *Rules
		player       *Player
		target       *Planet
		planetPlayer *Player
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Vector
	}{
		{
			name:   "move 25ly",
			fields: fields{Cargo: Cargo{Ironium: 100}, WarpSpeed: 5, builtThisTurn: false},
			args:   args{player: player, target: NewPlanet().WithNum(1).withPosition(Vector{100, 0})},
			want:   Vector{25, 0},
		},
		{
			name:   "move 18ly (just launched, warp 6)",
			fields: fields{Cargo: Cargo{Ironium: 100}, WarpSpeed: 6, builtThisTurn: true},
			args:   args{player: player, target: NewPlanet().WithNum(1).withPosition(Vector{100, 0})},
			want:   Vector{18, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := newMineralPacket(tt.args.player, 1, tt.fields.WarpSpeed, 5, tt.fields.Cargo, Vector{}, tt.args.target.Num)
			packet.builtThisTurn = tt.fields.builtThisTurn

			packet.movePacket(tt.args.rules, tt.args.player, tt.args.target, tt.args.planetPlayer)

			if packet.Position != tt.want {
				t.Errorf("MineralPacket.movePacket() = %v, want %v", packet, tt.want)
			}
		})
	}
}

func TestMineralPacket_completeMoveEmptyPlanet(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1)

	packet := newMineralPacket(player, 1, 5, 5, Cargo{300, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, nil)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 100})
	assert.True(t, packet.Delete)
}

func TestMineralPacket_completeMoveUncaught(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1).WithPlayerNum(1).WithCargo(Cargo{Colonists: 10000})

	packet := newMineralPacket(player, 1, 5, 5, Cargo{480, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, player)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 160, Colonists: 2500})
	assert.True(t, packet.Delete)

}

func TestMineralPacket_completeMoveUncaughtAR(t *testing.T) {
	player := NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1).WithPlayerNum(1).WithCargo(Cargo{Colonists: 100})

	packet := newMineralPacket(player, 1, 5, 5, Cargo{100, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, player)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 33, Colonists: 100})
	assert.True(t, packet.Delete)

}
func TestMineralPacket_completeMoveCaught(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1).WithPlayerNum(1).WithCargo(Cargo{Colonists: 100})
	planet.Spec.HasStarbase = true
	planet.Spec.HasMassDriver = true
	planet.Spec.SafePacketSpeed = 5

	packet := newMineralPacket(player, 1, 5, 5, Cargo{100, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, player)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 100, Colonists: 100})
	assert.True(t, packet.Delete)

}

func TestMineralPacket_getPacketDecayRate(t *testing.T) {
	type fields struct {
		SafeWarpSpeed int
		WarpSpeed     int
	}
	type args struct {
		race *Race
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{"no decay", fields{WarpSpeed: 5, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, 0},
		{"overwarp1", fields{WarpSpeed: 6, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .1},
		{"overwarp2", fields{WarpSpeed: 7, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .25},
		{"overwarp3", fields{WarpSpeed: 8, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .5},
		{"overwarp4", fields{WarpSpeed: 9, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .5},
		// should be equivalent to 2 levels lower
		{"overwarp4 as PP", fields{WarpSpeed: 9, SafeWarpSpeed: 5}, args{NewRace().WithPRT(PP).WithSpec(&rules)}, .25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.args.race).withSpec(&rules)
			packet := newMineralPacket(player, 1, tt.fields.WarpSpeed, tt.fields.SafeWarpSpeed, Cargo{}, Vector{}, 1)
			if got := packet.getPacketDecayRate(&rules, tt.args.race); got != tt.want {
				t.Errorf("MineralPacket.getPacketDecayRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMineralPacket_estimateDamage(t *testing.T) {
	type fields struct {
		SafeWarpSpeed int
		WarpSpeed     int
	}
	type args struct {
		race              *Race
		planetDriverSpeed int
		planetPosition    Vector
		planetDefCoverage float64
		targetRace        *Race
		planetPop         int
		mass              Cargo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MineralPacketDamage
	}{
		{
			`1 yr away; vanishing packet`,
			fields{WarpSpeed: 5, SafeWarpSpeed: 5},
			args{
				race:              NewRace().WithSpec(&rules),
				planetDriverSpeed: 0,
				planetPosition:    Vector{25, 0},
				planetDefCoverage: 0,
				targetRace:        NewRace().WithSpec(&rules),
				planetPop:         1000000,
				mass:              Cargo{Ironium: 10, Boranium: 10, Germanium: 10},
			},
			MineralPacketDamage{Uncaught: -1},
		},
		{
			`3 lvls overwarp + 1 yr travel w/ min decay (300 dmg)`,
			fields{WarpSpeed: 10, SafeWarpSpeed: 7},
			args{
				race:              NewRace().WithSpec(&rules).WithPRT(PP),
				planetDriverSpeed: 0,
				planetPosition:    Vector{60, 80},
				// 60^2 + 80^2 = 100^2
				planetDefCoverage: 0,
				targetRace:        NewRace().WithSpec(&rules),
				planetPop:         100000,
				mass:              Cargo{Ironium: 5, Boranium: 5, Germanium: 640},
			},
			MineralPacketDamage{Killed: 30000, DefensesDestroyed: 10},
		},
		{
			`3 lvls overwarp + 0.33 yr travel (25.2 dmg)`,
			fields{WarpSpeed: 8, SafeWarpSpeed: 5},
			args{
				race:              NewRace().WithSpec(&rules),
				planetDriverSpeed: 0,
				planetPosition:    Vector{21.333333, 0},
				planetDefCoverage: 0,
				targetRace:        NewRace().WithSpec(&rules),
				planetPop:         1000000,
				mass:              Cargo{Germanium: 75},
			},
			MineralPacketDamage{Killed: 25200, DefensesDestroyed: 1},
		},
		{
			`3 lvls overwarp + 3.25 yr travel (70 dmg)`,
			fields{WarpSpeed: 6, SafeWarpSpeed: 3},
			args{
				race:              NewRace().WithSpec(&rules),
				planetDriverSpeed: 4,
				planetPosition:    Vector{117, 0},
				planetDefCoverage: 0.5,
				targetRace:        NewRace().WithSpec(&rules),
				planetPop:         100000,
				mass:              Cargo{Ironium: 10240},
			},
			MineralPacketDamage{Killed: 7000, DefensesDestroyed: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.args.race).withSpec(&rules).WithNum(1)
			planet := NewPlanet().withPosition(tt.args.planetPosition).WithPlayerNum(2)
			planetPlayer := NewPlayer(2, tt.args.targetRace).withSpec(&rules).WithNum(2)
			planet.setPopulation(tt.args.planetPop)
			planet.Defenses = 10
			planet.Spec.DefenseCoverage = tt.args.planetDefCoverage
			planet.Spec.PlanetStarbaseSpec.HasStarbase = true
			planet.Spec.PlanetStarbaseSpec.HasMassDriver = true
			planet.Spec.PlanetStarbaseSpec.SafePacketSpeed = tt.args.planetDriverSpeed
			packet := newMineralPacket(player, 1, tt.fields.WarpSpeed, tt.fields.SafeWarpSpeed, tt.args.mass, Vector{0, 0}, 1)
			if got := packet.estimateDamage(&rules, player, planet, planetPlayer); got != tt.want {
				t.Errorf("MineralPacket.estimateDamage() = %v, want %v;", got, tt.want)
			}
		})
	}
}
