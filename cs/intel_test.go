//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_discover_discoverWormhole1(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// create a wormhole pair
	wormhole1 := newWormhole(Vector{}, 1, WormholeStabilityStable)
	wormhole2 := newWormhole(Vector{}, 2, WormholeStabilityStable)
	wormhole1.DestinationNum = wormhole2.Num
	wormhole2.DestinationNum = wormhole1.Num

	d := newDiscoverer(testLogger, player)
	d.discoverWormhole(wormhole1)
	assert.Equal(t, 1, len(player.WormholeIntels))
	assert.Equal(t, 1, player.WormholeIntels[0].Num)
	assert.Equal(t, None, player.WormholeIntels[0].DestinationNum)

}

func Test_discover_discoverWormhole2(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// create a wormhole pair
	wormhole1 := newWormhole(Vector{}, 1, WormholeStabilityStable)
	wormhole2 := newWormhole(Vector{}, 2, WormholeStabilityStable)
	wormhole1.DestinationNum = wormhole2.Num
	wormhole2.DestinationNum = wormhole1.Num

	// discover both wormholes
	d := newDiscoverer(testLogger, player)
	d.discoverWormhole(wormhole1)
	d.discoverWormhole(wormhole2)
	d.discoverWormholeLink(wormhole1, wormhole2)
	assert.Equal(t, 2, len(player.WormholeIntels))
	assert.Equal(t, 1, player.WormholeIntels[0].Num)
	assert.Equal(t, 2, player.WormholeIntels[1].Num)
	assert.Equal(t, wormhole2.Num, player.WormholeIntels[0].DestinationNum)
	assert.Equal(t, wormhole1.Num, player.WormholeIntels[1].DestinationNum)

}

func Test_discover_forgetWormhole1(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// create a wormhole pair
	wormhole1 := newWormhole(Vector{}, 1, WormholeStabilityStable)
	wormhole2 := newWormhole(Vector{}, 2, WormholeStabilityStable)
	wormhole1.DestinationNum = wormhole2.Num
	wormhole2.DestinationNum = wormhole1.Num

	d := newDiscoverer(testLogger, player)
	d.discoverWormhole(wormhole1)
	assert.Equal(t, 1, len(player.WormholeIntels))
	assert.Equal(t, 1, player.WormholeIntels[0].Num)

	d.forgetWormhole(wormhole1.Num)
	assert.Equal(t, 0, len(player.WormholeIntels))
}

func Test_discover_forgetWormhole2(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)

	// create a wormhole pair
	wormhole1 := newWormhole(Vector{}, 1, WormholeStabilityStable)
	wormhole2 := newWormhole(Vector{}, 2, WormholeStabilityStable)
	wormhole1.DestinationNum = wormhole2.Num
	wormhole2.DestinationNum = wormhole1.Num

	// discover both wormholes so we know the link
	d := newDiscoverer(testLogger, player)

	// should do nothing
	d.forgetWormhole(wormhole1.Num)
	assert.Equal(t, 0, len(player.WormholeIntels))

	// should discover, then forget wormhole
	d.discoverWormhole(wormhole1)
	assert.Equal(t, 1, len(player.WormholeIntels))
	d.forgetWormhole(wormhole1.Num)
	assert.Equal(t, 0, len(player.WormholeIntels))

	// should discover forget link
	d.discoverWormhole(wormhole1)
	d.discoverWormhole(wormhole2)
	d.discoverWormholeLink(wormhole1, wormhole2)
	assert.Equal(t, 2, len(player.WormholeIntels))

	d.forgetWormhole(wormhole1.Num)
	assert.Equal(t, 1, len(player.WormholeIntels))
	assert.Equal(t, None, player.WormholeIntels[0].DestinationNum)
}

func Test_discover_discoverPlayerScores(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(1).withSpec(&rules)
	otherPlayer := NewPlayer(1, NewRace().WithSpec(&rules)).WithNum(2).withSpec(&rules)
	player.Intels.ScoreIntels = make([]ScoreIntel, 2)

	// 25 years of scores
	otherPlayer.ScoreHistory = make([]PlayerScore, 25)
	for i := range otherPlayer.ScoreHistory {
		otherPlayer.ScoreHistory[i] = PlayerScore{Score: i + 1}
	}

	d := newDiscoverer(testLogger, player)
	d.discoverPlayerScores(otherPlayer, 20)

	history := player.Intels.ScoreIntels[1].ScoreHistory
	assert.Equal(t, 25, len(history))
	// the first 20 years are hidden
	for i := range 20 {
		assert.Equal(t, PlayerScore{}, history[i])
	}
	// the rest are visible
	for i := 20; i < 25; i++ {
		assert.Equal(t, i+1, history[i].Score)
	}

	// shorter history than the hidden years is all hidden
	otherPlayer.ScoreHistory = otherPlayer.ScoreHistory[:10]
	d.discoverPlayerScores(otherPlayer, 20)
	assert.Equal(t, make([]PlayerScore, 10), player.Intels.ScoreIntels[1].ScoreHistory)
}

func Test_discover_discoverDesignShipClass(t *testing.T) {
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	design := &ShipDesign{PlayerNum: 2, Num: 1, Hull: Destroyer.Name, Spec: ShipDesignSpec{Mass: 30, ShipClass: Destroyer.Class}}

	// a design known from before hulls had classes
	player.ShipDesignIntels = append(player.ShipDesignIntels, &ShipDesign{PlayerNum: 2, Num: 1, Hull: Destroyer.Name, Spec: ShipDesignSpec{Mass: 30}})

	d := newDiscoverer(testLogger, player)
	d.discoverDesign(design, false)
	assert.Equal(t, 1, len(player.ShipDesignIntels))
	assert.Equal(t, ShipClassScout, player.ShipDesignIntels[0].Spec.ShipClass)

	// a new design learns its class from the hull, without its components being scanned
	design2 := &ShipDesign{PlayerNum: 2, Num: 2, Hull: Battleship.Name, Spec: ShipDesignSpec{Mass: 222, ShipClass: Battleship.Class}}
	d.discoverDesign(design2, false)
	assert.Equal(t, 2, len(player.ShipDesignIntels))
	assert.Equal(t, ShipClassWarship, player.ShipDesignIntels[1].Spec.ShipClass)
	assert.Equal(t, 0, len(player.ShipDesignIntels[1].Slots))
}

func Test_discover_discoverPlanetStarbase(t *testing.T) {
	planet := NewPlanet().WithNum(1).WithPlayerNum(2)
	player := NewPlayer(1, NewRace().WithSpec(&rules)).withSpec(&rules)
	assert.NoError(t, player.initDefaultPlanetIntels([]*Planet{planet}))
	intel := player.PlanetIntels[0]

	// we last saw this planet with an older starbase
	intel.Spec.HasStarbase = true
	intel.Spec.StarbaseDesignName = "Old Fort"
	intel.Spec.StarbaseDesignNum = 1

	// fighting the upgraded starbase tells us which design is there now
	design := &ShipDesign{PlayerNum: 2, Num: 3, Name: "Death Star"}
	planet.Starbase = &Fleet{Tokens: []ShipToken{{DesignNum: design.Num, Quantity: 1, design: design}}}
	planet.Spec.HasStarbase = true

	d := newDiscoverer(testLogger, player)
	assert.NoError(t, d.discoverPlanetStarbase(planet))
	assert.True(t, intel.Spec.HasStarbase)
	assert.Equal(t, "Death Star", intel.Spec.StarbaseDesignName)
	assert.Equal(t, 3, intel.Spec.StarbaseDesignNum)

	// a starbase destroyed in the battle is gone from our intel
	planet.Starbase = nil
	planet.Spec.HasStarbase = false
	assert.NoError(t, d.discoverPlanetStarbase(planet))
	assert.False(t, intel.Spec.HasStarbase)
	assert.Equal(t, "", intel.Spec.StarbaseDesignName)
	assert.Equal(t, 0, intel.Spec.StarbaseDesignNum)
}
