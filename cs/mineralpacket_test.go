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

	packet := newMineralPacket(player, 1, 5, 5, Cargo{100, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, nil)
	assert.Equal(t, planet.Cargo, packet.Cargo)
	assert.True(t, packet.Delete)
}

func TestMineralPacket_completeMoveUncaught(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1).WithPlayerNum(1).WithCargo(Cargo{Colonists: 100})

	packet := newMineralPacket(player, 1, 5, 5, Cargo{100, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, player)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 100, Colonists: 84})
	assert.True(t, packet.Delete)

}

func TestMineralPacket_completeMoveUncaughtAR(t *testing.T) {
	player := NewPlayer(1, NewRace().WithPRT(AR).WithSpec(&rules)).WithNum(1).withSpec(&rules)
	planet := NewPlanet().withPosition(Vector{20, 0}).WithNum(1).WithPlayerNum(1).WithCargo(Cargo{Colonists: 100})

	packet := newMineralPacket(player, 1, 5, 5, Cargo{100, 0, 0, 0}, Vector{}, planet.Num)

	packet.movePacket(&rules, player, planet, player)
	assert.Equal(t, planet.Cargo, Cargo{Ironium: 100, Colonists: 100})
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
		{"1 yr away; no decay/dmg", fields{WarpSpeed: 5, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules), 5, Vector{25.0, 0.0}, 0.0, NewRace().WithSpec(&rules), 1000000, Cargo{10, 10, 10, 0}}, MineralPacketDamage{}},
		{"3 lvls overwarp + 1 yr travel (50% left); 1M pop; no driver/defenses", fields{WarpSpeed: 8, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules), 0, Vector{64.0, 0.0}, 0.0, NewRace().WithSpec(&rules), 1000000, Cargo{50, 0, 0, 0}}, MineralPacketDamage{Killed: 10000}},
		{`3 lvls overwarp + 3.25 yr travel (10.9% left); 1M pop; W6 driver/50% defenses`, fields{WarpSpeed: 8, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules), 6, Vector{208.0, 0.0}, 0.5, NewRace().WithSpec(&rules), 1000000, Cargo{1045, 0, 0, 0}}, MineralPacketDamage{Killed: 1000}},
	}
	//To sirgwain (or whoever happens to be unlucky enough to test this): I did not do an IT/PP example as of yet because I too lazy. Feel free to add one if/when you have the mental capacity to do so.

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.args.race).withSpec(&rules)
			planet := NewPlanet().withPosition(tt.args.planetPosition)
			planet.Spec.Population = tt.args.planetPop
			planet.Spec.DefenseCoverage = tt.args.planetDefCoverage
			planet.Spec.PlanetStarbaseSpec.SafePacketSpeed = tt.args.planetDriverSpeed
			planetPlayer := NewPlayer(2, tt.args.targetRace).withSpec(&rules)
			packet := newMineralPacket(player, 1, tt.fields.WarpSpeed, tt.fields.SafeWarpSpeed, tt.args.mass, Vector{0, 0}, 1)
			if got := packet.estimateDamage(&rules, player, planet, planetPlayer); got != tt.want {
				t.Errorf("MineralPacket.estimateDamage() = %v, want %v", got, tt.want)
			}
		})
	}
}
