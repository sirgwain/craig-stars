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

	d := newDiscoverer(player)
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
	d := newDiscoverer(player)
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

	d := newDiscoverer(player)
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
	d := newDiscoverer(player)

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
