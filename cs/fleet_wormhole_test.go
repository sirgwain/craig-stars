//go:build !wasi && !wasm

package cs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wormholeRoute(t *testing.T) (*Player, *Fleet) {
	t.Helper()
	player := testPlayer().WithNum(1)
	player.WormholeIntels = []*Wormhole{
		{MapObject: MapObject{Num: 1, Position: Vector{100, 100}}, DestinationNum: 2},
		{MapObject: MapObject{Num: 2, Position: Vector{1000, 100}}, DestinationNum: 1},
	}
	entrance := NewPositionWaypoint(Vector{100, 100}, 6)
	entrance.MapObjectTarget = MapObjectTarget{TargetType: MapObjectTypeWormhole, TargetNum: 1}
	fleet := testLongRangeScout(player).withWaypoints(
		NewPositionWaypoint(Vector{64, 100}, 0),
		entrance,
		NewPositionWaypoint(Vector{1036, 100}, 6),
	)
	return player, fleet
}

func TestFleet_WormholeFuelEstimates(t *testing.T) {
	t.Run("each leg departs from the wormhole exit", func(t *testing.T) {
		player, fleet := wormholeRoute(t)
		fleet.ComputeFuelUsage(player)
		assert.Equal(t, 5, fleet.Waypoints[1].EstFuelUsage) // 36 ly to the entrance
		assert.Equal(t, 5, fleet.Waypoints[2].EstFuelUsage) // 36 ly from the exit
		assert.Equal(t, 10, fleet.GetFuelAllocated(player, 2))
		assert.Equal(t, 2, fleet.getYearsToWaypoint(player, 2))
		assert.Equal(t, Vector{100, 100}, fleet.Waypoints[1].Position)
	})

	for _, missing := range []string{"entrance intel", "link", "exit intel", "wormhole target"} {
		t.Run("falls back without "+missing, func(t *testing.T) {
			player, fleet := wormholeRoute(t)
			switch missing {
			case "entrance intel":
				player.WormholeIntels = player.WormholeIntels[1:]
			case "link":
				player.WormholeIntels[0].DestinationNum = None
			case "exit intel":
				player.WormholeIntels = player.WormholeIntels[:1]
			case "wormhole target":
				fleet.Waypoints[1].MapObjectTarget = MapObjectTarget{}
			}
			fleet.ComputeFuelUsage(player)
			assert.Equal(t, fleet.GetFuelCost(player, 6, 936), fleet.Waypoints[2].EstFuelUsage)
		})
	}

	t.Run("current position does not teleport", func(t *testing.T) {
		player, fleet := wormholeRoute(t)
		fleet.Waypoints = fleet.Waypoints[1:]
		fleet.ComputeFuelUsage(player)
		assert.Equal(t, fleet.GetFuelCost(player, 6, 936), fleet.Waypoints[1].EstFuelUsage)
	})
}

func TestFleet_WormholeWaypointOrders(t *testing.T) {
	for _, update := range []bool{false, true} {
		name := "add"
		if update {
			name = "update"
		}
		t.Run(name, func(t *testing.T) {
			player, fleet := wormholeRoute(t)
			dest := WaypointDest{Position: Vector{1025, 100}}
			if update {
				require.Equal(t, UpdateWaypointResultUpdated, fleet.UpdateWaypoint(player, dest, 2, false))
			} else {
				fleet.Waypoints = fleet.Waypoints[:2]
				require.Equal(t, 2, fleet.AddWaypoint(player, dest, 1, false))
			}
			assert.Equal(t, 5, fleet.Waypoints[2].WarpSpeed) // 25 ly from the exit takes one year at warp 5
			assert.Equal(t, 4, fleet.Waypoints[2].EstFuelUsage)
		})
	}

	t.Run("can fly back to the entrance after emerging", func(t *testing.T) {
		player, fleet := wormholeRoute(t)
		fleet.Waypoints = fleet.Waypoints[:2]
		require.Equal(t, 2, fleet.AddWaypoint(player, WaypointDest{Position: Vector{100, 100}}, 1, false))
		require.Equal(t, UpdateWaypointResultUpdated, fleet.UpdateWaypoint(player, WaypointDest{Position: Vector{100, 100}}, 2, false))
	})
}

func TestFleet_WormholePursuit(t *testing.T) {
	for _, update := range []bool{false, true} {
		t.Run(map[bool]string{false: "add", true: "update"}[update], func(t *testing.T) {
			player, fleet := wormholeRoute(t)
			target := &Fleet{MapObject: MapObject{Type: MapObjectTypeFleet, PlayerNum: 2, Num: 1, Position: Vector{1043, 100}}, Heading: VectorFloat64{X: -1}, WarpSpeed: 1}
			player.FleetIntels = []*Fleet{target}
			if update {
				fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{1045, 100}, 6))
				require.Equal(t, UpdateWaypointResultUpdated, fleet.UpdateWaypoint(player, WaypointDest{MO: target.MapObject}, 3, false))
			} else {
				require.Equal(t, 3, fleet.AddWaypoint(player, WaypointDest{MO: target.MapObject}, 2, false))
			}
			// Two years to reach the waypoint after the exit. The target is then at 1041,
			// and moves to 1040 on the intercept turn: 4 ly from 1036, reachable at warp 2.
			assert.Equal(t, 2, fleet.Waypoints[3].WarpSpeed)
		})
	}
}

func TestFleet_ChainedWormholeEstimates(t *testing.T) {
	player, fleet := wormholeRoute(t)
	player.WormholeIntels = append(player.WormholeIntels,
		&Wormhole{MapObject: MapObject{Num: 3, Position: Vector{1036, 100}}, DestinationNum: 4},
		&Wormhole{MapObject: MapObject{Num: 4, Position: Vector{2000, 100}}, DestinationNum: 3},
	)
	fleet.Waypoints[2].MapObjectTarget = MapObjectTarget{TargetType: MapObjectTypeWormhole, TargetNum: 3}
	fleet.Waypoints = append(fleet.Waypoints, NewPositionWaypoint(Vector{2036, 100}, 6))
	fleet.ComputeFuelUsage(player)
	assert.Equal(t, 15, fleet.GetFuelAllocated(player, 3))
	assert.Equal(t, 3, fleet.getYearsToWaypoint(player, 3))
}

func TestFleet_WormholeEditRefreshesFollowingFuel(t *testing.T) {
	player, fleet := wormholeRoute(t)
	fleet.Waypoints[1] = NewPositionWaypoint(Vector{100, 100}, 6)
	fleet.ComputeFuelUsage(player)
	dest := WaypointDest{MO: MapObject{Type: MapObjectTypeWormhole, Num: 1, Position: Vector{100, 100}}}
	require.Equal(t, UpdateWaypointResultUpdated, fleet.UpdateWaypoint(player, dest, 1, false))
	assert.Equal(t, 5, fleet.Waypoints[2].EstFuelUsage)

	// Inserting the wormhole before an existing destination refreshes the following leg too.
	fleet.Waypoints = append(fleet.Waypoints[:1], NewPositionWaypoint(Vector{1036, 100}, 6))
	require.Equal(t, 1, fleet.AddWaypoint(player, dest, 0, false))
	assert.Equal(t, 5, fleet.Waypoints[2].EstFuelUsage)
}

func Test_updateFleetTargetsAfterWormhole(t *testing.T) {
	for _, atExit := range []bool{false, true} {
		t.Run(map[bool]string{false: "keep return to entrance", true: "drop redundant exit"}[atExit], func(t *testing.T) {
			player, fleet := wormholeRoute(t)
			position := Vector{100, 100}
			if atExit {
				position = Vector{1000, 100}
			}
			// A lost fleet target becomes a position waypoint. It is redundant only at the exit.
			fleet.Waypoints[2] = NewFleetWaypoint(position, 1, 2, "Lost fleet", 6)
			scan := playerScanner{player: player, universe: &Universe{Fleets: []*Fleet{fleet}}}
			scan.updateFleetTargets()
			if atExit {
				assert.Len(t, fleet.Waypoints, 2)
			} else {
				require.Len(t, fleet.Waypoints, 3)
				assert.Equal(t, MapObjectTypeNone, fleet.Waypoints[2].TargetType)
			}
		})
	}
}
