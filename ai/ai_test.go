package ai

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

func Test_aiPlayer_ProcessTurn(t *testing.T) {
	type fields struct {
		prt cs.PRT
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"HE", fields{cs.HE}},
		{"SS", fields{cs.SS}},
		{"WM", fields{cs.WM}},
		{"CA", fields{cs.CA}},
		{"IS", fields{cs.IS}},
		{"SD", fields{cs.SD}},
		{"PP", fields{cs.PP}},
		{"IT", fields{cs.IT}},
		{"AR", fields{cs.AR}},
		{"JoaT", fields{cs.JoaT}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			race := cs.NewRace().WithPRT(tt.fields.prt)
			gamer := cs.NewGamer()
			game := gamer.CreateGame(0, *cs.NewGameSettings().WithAIPlayer(cs.AIDifficultyEasy, 0))
			player := gamer.NewPlayer(0, *race, &game.Rules)
			player.Num = 1
			player.Name = cs.AINames[0][0]
			universe, err := gamer.GenerateUniverse(game, []*cs.Player{player})
			if err != nil {
				t.Error(err)
				return
			}

			// process a turn
			ai := NewAIPlayer(game, &cs.StaticTechStore, player, universe.GetPlayerMapObjects(player.Num))
			if err := ai.ProcessTurn(); err != nil {
				t.Errorf("process turn %v", err)
				return
			}
			ai.SubmittedTurn = true

			// generate a ne turn, make sure no errors
			gamer.GenerateTurn(game, universe, []*cs.Player{player})
		})
	}
}

func Test_getClosestPlanet(t *testing.T) {
	game := cs.NewGame()
	player := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
	aiPlayer := NewAIPlayer(game, &cs.StaticTechStore, player, cs.PlayerMapObjects{})

	planetAt0_0 := cs.PlanetIntel{
		MapObjectIntel: cs.MapObjectIntel{Position: cs.Vector{X: 0, Y: 0}},
	}
	planetAt50_50 := cs.PlanetIntel{
		MapObjectIntel: cs.MapObjectIntel{Position: cs.Vector{X: 50, Y: 50}},
	}
	planetAt100_100 := cs.PlanetIntel{
		MapObjectIntel: cs.MapObjectIntel{Position: cs.Vector{X: 100, Y: 100}},
	}

	type args struct {
		position            cs.Vector
		unknownPlanetsByNum map[int]cs.PlanetIntel
	}
	tests := []struct {
		name string
		args args
		want *cs.PlanetIntel
	}{
		{"no planets, should be nil", args{cs.Vector{}, map[int]cs.PlanetIntel{}}, nil},
		{"1 planet, should be it", args{cs.Vector{}, map[int]cs.PlanetIntel{
			1: planetAt0_0,
		}}, &planetAt0_0},
		{"2 planets, should be closer one", args{cs.Vector{}, map[int]cs.PlanetIntel{
			1: planetAt100_100,
			2: planetAt50_50,
		}}, &planetAt50_50},
		{"2 planets, should be closer one, regardless of order", args{cs.Vector{}, map[int]cs.PlanetIntel{
			1: planetAt50_50,
			2: planetAt100_100,
		}}, &planetAt50_50},
		{"3 planets, should be closer one", args{cs.Vector{}, map[int]cs.PlanetIntel{
			1: planetAt50_50,
			2: planetAt100_100,
			3: planetAt0_0,
		}}, &planetAt0_0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aiPlayer.getClosestPlanetIntel(tt.args.position, tt.args.unknownPlanetsByNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClosestPlanet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIPlayer_GetPlanet(t *testing.T) {
	game := cs.NewGame()
	player := NewAIPlayer(game, &cs.StaticTechStore, cs.NewPlayer(1, cs.NewRace()), cs.PlayerMapObjects{})

	// no planet by that id
	assert.Nil(t, player.getPlanet(1))

	// should have a planet by this id
	planet := cs.NewPlanet()
	planet.Num = 1
	player.Planets = append(player.Planets, planet)
	player.buildMaps()

	assert.Same(t, planet, player.getPlanet(1))

	assert.Nil(t, player.getPlanet(2))
}
