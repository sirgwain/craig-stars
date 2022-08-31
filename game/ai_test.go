package game

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getClosestPlanet(t *testing.T) {
	rules := NewRules()
	player := NewPlayer(1, NewRace().WithSpec(&rules))
	aiPlayer := NewAIPlayer(player, PlayerMapObjects{})

	planetAt0_0 := PlanetIntel{
		MapObjectIntel: MapObjectIntel{Position: Vector{0, 0}},
	}
	planetAt50_50 := PlanetIntel{
		MapObjectIntel: MapObjectIntel{Position: Vector{50, 50}},
	}
	planetAt100_100 := PlanetIntel{
		MapObjectIntel: MapObjectIntel{Position: Vector{100, 100}},
	}

	type args struct {
		fleet               *Fleet
		unknownPlanetsByNum map[int]PlanetIntel
	}
	tests := []struct {
		name string
		args args
		want *PlanetIntel
	}{
		{"no planets, should be nil", args{testLongRangeScout(player, &rules), map[int]PlanetIntel{}}, nil},
		{"1 planet, should be it", args{testLongRangeScout(player, &rules), map[int]PlanetIntel{
			1: planetAt0_0,
		}}, &planetAt0_0},
		{"2 planets, should be closer one", args{testLongRangeScout(player, &rules), map[int]PlanetIntel{
			1: planetAt100_100,
			2: planetAt50_50,
		}}, &planetAt50_50},
		{"2 planets, should be closer one, regardless of order", args{testLongRangeScout(player, &rules), map[int]PlanetIntel{
			1: planetAt50_50,
			2: planetAt100_100,
		}}, &planetAt50_50},
		{"3 planets, should be closer one", args{testLongRangeScout(player, &rules), map[int]PlanetIntel{
			1: planetAt50_50,
			2: planetAt100_100,
			3: planetAt0_0,
		}}, &planetAt0_0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aiPlayer.getClosestPlanet(tt.args.fleet, tt.args.unknownPlanetsByNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClosestPlanet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAIPlayer_GetPlanet(t *testing.T) {
	player := NewAIPlayer(NewPlayer(1, NewRace()), PlayerMapObjects{})

	// no planet by that id
	assert.Nil(t, player.GetPlanet(1))

	// should have a planet by this id
	planet := NewPlanet()
	planet.Num = 1
	player.Planets = append(player.Planets, planet)
	player.buildMaps()

	assert.Same(t, planet, player.GetPlanet(1))

	assert.Nil(t, player.GetPlanet(2))
}