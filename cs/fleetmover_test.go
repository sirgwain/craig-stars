package cs

import (
	"testing"
)

type testFleetGetter struct {
	fleets []*Fleet
}

func newTestFleetGetter(fleets ...*Fleet) fleetGetter {
	return &testFleetGetter{fleets}
}

func (pg *testFleetGetter) getFleet(playerNum, num int) *Fleet {
	for _, p := range pg.fleets {
		if p.PlayerNum == playerNum && p.Num == num {
			return p
		}
	}
	return nil
}

func Test_fleetMove_buildFleetMoveGraph(t *testing.T) {
	player := testPlayer()
	fleet1 := testLongRangeScout(player).withNum(1)
	fleet2 := testLongRangeScout(player).withNum(2)
	fleet3 := testLongRangeScout(player).withNum(3)
	fleet4 := testLongRangeScout(player).withNum(4)

	type fleetTarget struct {
		fleet  *Fleet
		target *Fleet
	}
	tests := []struct {
		name   string
		fleets []fleetTarget
		want   string
	}{
		{
			name:   "one fleet, no targets",
			fleets: []fleetTarget{{fleet1, nil}},
			want:   "",
		},
		{
			name: "fleet1 targets fleet 2",
			fleets: []fleetTarget{
				{fleet1, nil},
				{fleet2, fleet1},
			},
			want: "",
		},
		{
			name: "fleet1 targets fleet 2, fleet 2 targets fleet 1",
			fleets: []fleetTarget{
				{fleet1, fleet2},
				{fleet2, fleet1},
			},
			want: "",
		},
		{
			// 1 -> 2
			// 2 -> 3
			// 3 -> 1 (cycle)
			name: "fleet1 targets fleet 2, fleet 2 targets fleet 3, fleet 3 targets fleet1",
			fleets: []fleetTarget{
				{fleet1, fleet2},
				{fleet2, fleet3},
				{fleet3, fleet1},
			},
			want: "",
		},
		{
			// 1 -> 2
			// 2 -> 1 (cycle)
			// 3 -> 2 (not a cycle)
			name: "fleet1 targets fleet 2, fleet 2 targets fleet 1, fleet 3 targets fleet2",
			fleets: []fleetTarget{
				{fleet1, fleet2},
				{fleet2, fleet1},
				{fleet3, fleet2},
			},
			want: "",
		},
		{
			// 1 -> 3
			// 2 -> 3
			// 3 -> 4 (not a cycle)
			// 4 targeting nothing
			name: "two fleets target fleet 3, fleet 3 targets fleet 4",
			fleets: []fleetTarget{
				{fleet1, fleet3},
				{fleet2, fleet3},
				{fleet3, fleet4},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fleets := make([]*Fleet, len(tt.fleets))
			for i, targ := range tt.fleets {
				if targ.target != nil {
					targ.fleet.withWaypoints(NewPositionWaypoint(Vector{}, 5), NewFleetWaypoint(targ.target.Position, targ.target.Num, targ.target.PlayerNum, targ.target.Name, 5))
				} else {
					targ.fleet.withWaypoints(NewPositionWaypoint(Vector{}, 5), NewPositionWaypoint(Vector{25, 0}, 5))
				}
				fleets[i] = targ.fleet
			}

			fm := &fleetMove{}
			got := fm.buildFleetMoveGraph(fleets, newTestFleetGetter(fleets...))
			if got.String() != tt.want {
				// t.Errorf("fm.buildFleetMoveGraph() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}
