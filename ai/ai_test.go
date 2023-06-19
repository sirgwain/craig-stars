package ai

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

// create a new long rang scout fleet for testing
func testLongRangeScout(player *cs.Player) *cs.Fleet {
	fleet := &cs.Fleet{
		BaseName:          "Long Range Scout",
		Tokens:            []cs.ShipToken{},
		OrbitingPlanetNum: cs.None,
	}
	// fleet.Spec = computeFleetSpec(rules, player, fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	return fleet
}

func Test_getClosestPlanet(t *testing.T) {
	game := cs.NewGame()
	player := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
	aiPlayer := NewAIPlayer(game, player, cs.PlayerMapObjects{})

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
		fleet               *cs.Fleet
		unknownPlanetsByNum map[int]cs.PlanetIntel
	}
	tests := []struct {
		name string
		args args
		want *cs.PlanetIntel
	}{
		{"no planets, should be nil", args{testLongRangeScout(player), map[int]cs.PlanetIntel{}}, nil},
		{"1 planet, should be it", args{testLongRangeScout(player), map[int]cs.PlanetIntel{
			1: planetAt0_0,
		}}, &planetAt0_0},
		{"2 planets, should be closer one", args{testLongRangeScout(player), map[int]cs.PlanetIntel{
			1: planetAt100_100,
			2: planetAt50_50,
		}}, &planetAt50_50},
		{"2 planets, should be closer one, regardless of order", args{testLongRangeScout(player), map[int]cs.PlanetIntel{
			1: planetAt50_50,
			2: planetAt100_100,
		}}, &planetAt50_50},
		{"3 planets, should be closer one", args{testLongRangeScout(player), map[int]cs.PlanetIntel{
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
	game := cs.NewGame()
	player := NewAIPlayer(game, cs.NewPlayer(1, cs.NewRace()), cs.PlayerMapObjects{})

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
