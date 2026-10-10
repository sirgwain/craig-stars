//go:build !wasi && !wasm

package cs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolvePlanetInvasions(t *testing.T) {
	tests := []struct {
		name     string
		invasion invasion
		want     invasionResult
	}{
		{
			name: "10000 attackers 10000 defenders, attacker wins",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(10_000),
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 10_000,
			},
			want: invasionResult{
				defenders:          10_000,
				attackersKilled:    9_100,
				defendersKilled:    10_000,
				remainingAttackers: 900,
				remainingDefenders: 0,
				successful:         true,
			},
		},
		{
			name: "5000 attackers for 10000 undefended defenders, defenders win",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(10_000),
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 5000,
			},
			want: invasionResult{
				defenders:          10_000,
				attackersKilled:    5000,
				defendersKilled:    5500,
				remainingAttackers: 0,
				remainingDefenders: 4500,
				successful:         false,
			},
		},
		{
			name: "100,000 attackers for 100,000 well defended defenders, defenders win",
			invasion: invasion{
				planet: &Planet{
					Cargo: Cargo{}.WithPopulation(100_000),
					Spec:  PlanetSpec{DefenseCoverage: .9},
				},
				defender:  NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules),
				attacker:  NewPlayer(2, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules),
				attackers: 100_000,
			},
			want: invasionResult{
				defenders:                 100_000,
				attackersKilled:           100_000,
				defendersKilled:           35700,
				attackersKilledByDefenses: 67501,
				remainingAttackers:        0,
				remainingDefenders:        64_300,
				successful:                false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolvePlanetInvasions(&rules, []invasion{tt.invasion})[0]
			// zero out the invasion itself in the result, we only care about the result numbers
			got.invasion = invasion{}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolvePlanetInvasions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_turn_resolveColonistDropsMultipleAttackers(t *testing.T) {
	tests := []struct {
		name                                 string
		defenderPop, a, b, winner, survivors int
	}{
		{"equal invaders leave no colony", 100_000, 200_000, 200_000, 0, 0},
		{"combined attacks overcome defender, stronger attacker wins", 100_000, 50_000, 60_000, 3, 1700},
		{"attacker order doesn't matter", 100_000, 60_000, 50_000, 2, 1700},
		{"equal to defender power, attacker wins", 11_000, 10_000, 0, 2, 100},
		{"defender survives combined attacks", 100_000, 20_000, 20_000, 1, 56_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TestScenario{Players: []ScenarioPlayer{{}, {}, {}}, Planets: []ScenarioPlanet{{Name: "Target", Owner: 1, Cargo: Cargo{}.WithPopulation(tt.defenderPop)}}}
			u := newTestUniverse(t, s)
			planet := u.Planet("Target")
			drops := newInvader()
			drops.addInvasion(invasion{planet: planet, defender: u.Player(1), attacker: u.Player(2), attackers: tt.a})
			if tt.b > 0 {
				drops.addInvasion(invasion{planet: planet, defender: u.Player(1), attacker: u.Player(3), attackers: tt.b})
			}

			u.turn.queueColonistDrops(drops)
			u.turn.resolveColonistDrops()

			assert.Equal(t, tt.winner, planet.PlayerNum)
			assert.Equal(t, tt.survivors, planet.GetPopulation())
		})
	}
}

func Test_turn_fleetColonizeContested(t *testing.T) {
	colonizer := func() ScenarioPlayer {
		return ScenarioPlayer{Designs: Designs(DesignSantaMaria), Fleets: []ScenarioFleet{{Design: "Santa Maria", At: "Target", Cargo: Cargo{Colonists: 100}, Waypoints: []ScenarioWaypoint{{To: "Target", Task: WaypointTaskColonize}}}}}
	}
	s := TestScenario{Players: []ScenarioPlayer{colonizer(), colonizer()}, Planets: []ScenarioPlanet{{Name: "Target"}}}
	u := newTestUniverse(t, s)

	u.turn.fleetColonize()
	u.turn.resolveColonistDrops()

	// equal colonizers destroy each other
	planet := u.Planet("Target")
	assert.False(t, planet.Owned())
	assert.Equal(t, 0, planet.GetPopulation())
	for _, f := range u.Game.Fleets {
		assert.True(t, f.Delete)
	}
	for _, num := range []int{1, 2} {
		messages := u.Messages(num, PlayerMessagePlanetColonizeContested)
		if assert.Len(t, messages, 1) {
			assert.Equal(t, 2, messages[0].Spec.Amount)
			assert.Equal(t, 0, messages[0].Spec.TargetPlayerNum)
		}
	}
}

func Test_turn_byHandAndWaypointDropsResolveTogether(t *testing.T) {
	s := TestScenario{Players: []ScenarioPlayer{
		{},
		{Designs: Designs(DesignGalleon), Fleets: []ScenarioFleet{{Design: "Galleon", At: "Target", Cargo: Cargo{Colonists: 500}}}},
		{Designs: Designs(DesignGalleon), Fleets: []ScenarioFleet{{Design: "Galleon", At: "Target", Cargo: Cargo{Colonists: 500}, Waypoints: []ScenarioWaypoint{{To: "Target", Task: WaypointTaskTransport, TransportTasks: WaypointTransportTasks{Colonists: WaypointTransportTask{Action: TransportActionUnloadAll}}}}}}},
	}, Planets: []ScenarioPlanet{{Name: "Target", Owner: 1, Cargo: Cargo{Colonists: 500}}}}
	u := newTestUniverse(t, s)
	u.TransferByHand(2, "Galleon #1", "Target", Cargo{Colonists: -500})

	u.GenerateTurn()

	// player 2 drops by hand, player 3 by waypoint, they tie and nobody survives
	assert.False(t, u.Planet("Target").Owned())
	assert.Equal(t, 0, u.Planet("Target").GetPopulation())
}

func Test_turn_fleetColonizeContestedLoserLearnsOwner(t *testing.T) {
	colonizer := func(colonists int) ScenarioPlayer {
		return ScenarioPlayer{Designs: Designs(DesignSantaMaria), Fleets: []ScenarioFleet{{Design: "Santa Maria", At: "Target", Cargo: Cargo{Colonists: colonists}, Waypoints: []ScenarioWaypoint{{To: "Target", Task: WaypointTaskColonize}}}}}
	}
	s := TestScenario{Players: []ScenarioPlayer{colonizer(100), colonizer(200)}, Planets: []ScenarioPlanet{{Name: "Target"}}}
	u := newTestUniverse(t, s)

	u.turn.fleetColonize()
	u.turn.resolveColonistDrops()

	// the loser's colony ship is gone, but they still know who took the planet
	planet := u.Planet("Target")
	assert.True(t, planet.OwnedBy(2))
	assert.Equal(t, 2, u.Game.Players[0].PlanetIntels[planet.Num-1].PlayerNum)
}
