package ai

import (
	"reflect"
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/stretchr/testify/assert"
)

func Test_aiPlayer_ProcessTurn(t *testing.T) {
	tests := []struct {
		name string
		prt  cs.PRT
	}{
		{"HE", cs.HE},
		{"SS", cs.SS},
		{"WM", cs.WM},
		{"CA", cs.CA},
		{"IS", cs.IS},
		{"SD", cs.SD},
		{"PP", cs.PP},
		{"IT", cs.IT},
		{"AR", cs.AR},
		{"JoaT", cs.JoaT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			race := cs.NewRace().WithPRT(tt.prt)
			gamer := cs.NewGamer()
			game := gamer.CreateGame(0, *cs.NewGameSettings().WithAIPlayer(cs.AIDifficultyEasy, 0))
			player := gamer.NewPlayer(0, *race, &game.Rules)
			player.Num = 1
			player.Name = cs.AINames[0][0]
			universe, err := gamer.GenerateUniverse(game, []*cs.Player{player})
			if err != nil {
				t.Fatalf("gamer.generateUniverse() failed: \n%v", err)
			}

			// process a turn
			ai := NewAIPlayer(game, &cs.StaticTechStore, player, universe.GetPlayerMapObjects(player.Num))
			if err := ai.ProcessTurn(); err != nil {
				t.Fatalf("ai turn processing failed: \n%v", err)
			}
			ai.SubmittedTurn = true

			// generate a new turn, make sure no errors
			gamer.GenerateTurn(game, universe, []*cs.Player{player})
		})
	}
}

func Test_aiPlayer_updateWarfleets(t *testing.T) {
	tests := []struct {
		name      string
		race      *cs.Race
		techLevel cs.TechLevel
		year      int
		want      []fleetShip
	}{
		{
			name:      "Tech 0 - default plan due to no warships",
			race:      cs.NewRace().WithPRT(cs.HE),
			techLevel: cs.TechLevel{},
			year:      25,
			want: []fleetShip{
				{purpose: cs.ShipDesignPurposeBomber, quantity: 5},
				{purpose: cs.ShipDesignPurposeBeamFighter, quantity: 7},
				{purpose: cs.ShipDesignPurposeTorpedoFighter, quantity: 7},
			},
		},
		{
			name:      "Tech 10 WM - beam BCs very good",
			race:      cs.NewRace().WithPRT(cs.WM),
			techLevel: cs.TechLevel{Energy: 10, Weapons: 10, Propulsion: 10, Construction: 10, Electronics: 10, Biotechnology: 10},
			year:      27,
			want: []fleetShip{
				{purpose: cs.ShipDesignPurposeBomber, quantity: 5},
				{purpose: cs.ShipDesignPurposeBeamFighter, quantity: 14},
				{purpose: cs.ShipDesignPurposeFuelFreighter, quantity: 3}, // 19/5
			},
		},
		{
			name:      "Tech 12 WM - Jihad BCs",
			race:      cs.NewRace().WithPRT(cs.WM).WithLRT(cs.RS),
			techLevel: cs.TechLevel{Energy: 6, Weapons: 12, Propulsion: 9, Construction: 10, Electronics: 11, Biotechnology: 7},
			year:      30,
			want: []fleetShip{
				{purpose: cs.ShipDesignPurposeBomber, quantity: 7},
				{purpose: cs.ShipDesignPurposeTorpedoFighter, quantity: 16},
				{purpose: cs.ShipDesignPurposeFuelFreighter, quantity: 4}, // 23/5
			},
		},
		{
			name:      "Tech 24 - Fairly even",
			race:      cs.NewRace().WithPRT(cs.JoaT),
			techLevel: cs.TechLevel{Energy: 24, Weapons: 24, Propulsion: 24, Construction: 24, Electronics: 24, Biotechnology: 24},
			year:      75, // 50 years after attack start yr
			want: []fleetShip{
				{purpose: cs.ShipDesignPurposeBomber, quantity: 40},
				{purpose: cs.ShipDesignPurposeBeamFighter, quantity: 36},
				{purpose: cs.ShipDesignPurposeTorpedoFighter, quantity: 24},
				{purpose: cs.ShipDesignPurposeFuelFreighter, quantity: 20}, // 100/5
			},
		},
		{
			name:      "Tech 26 HE - Beam leaning",
			race:      cs.NewRace().WithPRT(cs.HE),
			techLevel: cs.TechLevel{Energy: 26, Weapons: 26, Propulsion: 26, Construction: 26, Electronics: 26, Biotechnology: 26},
			year:      75, // 50 years after attack start yr
			want: []fleetShip{
				{purpose: cs.ShipDesignPurposeBomber, quantity: 40},
				{purpose: cs.ShipDesignPurposeBeamFighter, quantity: 60},
				{purpose: cs.ShipDesignPurposeFuelFreighter, quantity: 20}, // 100/5
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gamer := cs.NewGamer()
			game := gamer.CreateGame(0, *cs.NewGameSettings().WithAIPlayer(cs.AIDifficultyEasy, 0))
			game.Year = game.Rules.StartingYear + tt.year
			player := gamer.NewPlayer(0, *tt.race.WithSpec(&game.Rules), &game.Rules)
			player.Num = 1
			player.Name = tt.name
			universe, err := gamer.GenerateUniverse(game, []*cs.Player{player})
			if err != nil {
				t.Fatal(err)
			}
			ai := NewAIPlayer(game, &cs.StaticTechStore, player, universe.GetPlayerMapObjects(player.Num))
			ai.Player.TechLevels = tt.techLevel

			// update warship designs
			if ai.designsByPurpose[cs.ShipDesignPurposeFuelFreighter], err = ai.designShip("Fuel Pod", cs.ShipDesignPurposeFuelFreighter, cs.FleetPurposeFreighter); err != nil {
				t.Errorf("designing fuel ship for test returned error \n%v", err)
			}

			if err = ai.updateWarfleets(); err != nil {
				t.Errorf("aiPlayer.updateWarfleets() errored: \n%v", err)
			}
			got := ai.fleetsByPurpose[cs.FleetPurposeBomber].ships
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aiPlayer.updateWarfleets() returned incorrect fleet ratios; got: \n%v, want: \n%v", got, tt.want)
			}
		})
	}
}

func Test_getClosestPlanet(t *testing.T) {
	game := cs.NewGame()
	player := cs.NewPlayer(1, cs.NewRace().WithSpec(&game.Rules))
	aiPlayer := NewAIPlayer(game, &cs.StaticTechStore, player, cs.PlayerMapObjects{})

	planetAt0_0 := cs.PlanetIntel{
		MapObject: cs.MapObject{Position: cs.Vector{X: 0, Y: 0}},
	}
	planetAt50_50 := cs.PlanetIntel{
		MapObject: cs.MapObject{Position: cs.Vector{X: 50, Y: 50}},
	}
	planetAt100_100 := cs.PlanetIntel{
		MapObject: cs.MapObject{Position: cs.Vector{X: 100, Y: 100}},
	}

	tests := []struct {
		name                string
		position            cs.Vector
		unknownPlanetsByNum map[int]cs.PlanetIntel
		want                *cs.PlanetIntel
	}{
		{
			name:                "no planets",
			position:            cs.Vector{},
			unknownPlanetsByNum: map[int]cs.PlanetIntel{},
			want:                nil,
		},
		{
			name:     "1 planet",
			position: cs.Vector{},
			unknownPlanetsByNum: map[int]cs.PlanetIntel{
				1: planetAt0_0,
			},
			want: &planetAt0_0,
		},
		{
			name:     "2 planets - picks closest",
			position: cs.Vector{},
			unknownPlanetsByNum: map[int]cs.PlanetIntel{
				1: planetAt100_100,
				2: planetAt50_50,
			},
			want: &planetAt50_50,
		},
		{
			name:     "3 planets - picks closest",
			position: cs.Vector{},
			unknownPlanetsByNum: map[int]cs.PlanetIntel{
				12: planetAt0_0,
				1:  planetAt50_50,
				2:  planetAt100_100,
			},
			want: &planetAt0_0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aiPlayer.getClosestPlanetIntel(tt.position, tt.unknownPlanetsByNum); !reflect.DeepEqual(got, tt.want) {
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
