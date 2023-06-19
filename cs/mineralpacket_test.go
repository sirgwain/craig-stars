package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMineralPacket_movePacket(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)

	type fields struct {
		Cargo         Cargo
		WarpFactor    int
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
			fields: fields{Cargo: Cargo{Ironium: 100}, WarpFactor: 5, builtThisTurn: false},
			args:   args{player: player, target: NewPlanet().WithNum(1).withPosition(Vector{100, 0})},
			want:   Vector{25, 0},
		},
		{
			name:   "move 18ly (just launched, warp 6)",
			fields: fields{Cargo: Cargo{Ironium: 100}, WarpFactor: 6, builtThisTurn: true},
			args:   args{player: player, target: NewPlanet().WithNum(1).withPosition(Vector{100, 0})},
			want:   Vector{18, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := newMineralPacket(tt.args.player, 1, tt.fields.WarpFactor, 5, tt.fields.Cargo, Vector{}, tt.args.target.Num)
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
		WarpFactor    int
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
		{"no decay", fields{WarpFactor: 5, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, 0},
		{"overwarp1", fields{WarpFactor: 6, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .1},
		{"overwarp2", fields{WarpFactor: 7, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .25},
		{"overwarp3", fields{WarpFactor: 8, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .5},
		{"overwarp4", fields{WarpFactor: 9, SafeWarpSpeed: 5}, args{NewRace().WithSpec(&rules)}, .5},
		// should be equivalent to 2 levels lower
		{"overwarp4 as PP", fields{WarpFactor: 9, SafeWarpSpeed: 5}, args{NewRace().WithPRT(PP).WithSpec(&rules)}, .25},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(1, tt.args.race).withSpec(&rules)
			packet := newMineralPacket(player, 1, tt.fields.WarpFactor, tt.fields.SafeWarpSpeed, Cargo{}, Vector{}, 1)
			if got := packet.getPacketDecayRate(&rules, tt.args.race); got != tt.want {
				t.Errorf("MineralPacket.getPacketDecayRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
