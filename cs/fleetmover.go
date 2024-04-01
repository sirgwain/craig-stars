package cs

import (
	"fmt"

	"slices"
)

// The FleetMover builds a movement graph to handle fleets chasing fleets chasing fleets
// in the case of a circular graph, each fleet will move one ly at a time
// TODO: not yet implemented
type fleetMover interface {
	buildFleetMoveGraph(fleets []*Fleet, fleetGetter fleetGetter) fleetMoveGraph
}

type fleetMoveGraph struct {
	fleetsNotTargetingFleets []*Fleet
	fleetsTargetingFleets    []fleetMoveNode
	fleetsInTargetCycle      map[playerObject][]fleetMoveNode
}

type fleetMoveNode struct {
	fleet  *Fleet
	target *fleetMoveNode
}

func (node fleetMoveNode) String() string {
	return fmt.Sprintf("fleet %d -> %d", node.fleet.Num, node.target.fleet.Num)
}

type fleetMove struct {
}

// render the graph as a string for easy testing
func (graph *fleetMoveGraph) String() string {
	return fmt.Sprintf(`
fleetsNotTargetingFleets: %v,\n
fleetsTargetingFleets: %v,\n
fleetsInTargetCycle: %v,\n
`,
		graph.fleetsNotTargetingFleets,
		graph.fleetsTargetingFleets,
		graph.fleetsInTargetCycle,
	)
}

func reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (fm *fleetMove) buildFleetMoveGraph(fleets []*Fleet, fleetGetter fleetGetter) fleetMoveGraph {
	graph := fleetMoveGraph{}
	nodes := make(map[*Fleet]*fleetMoveNode)

	for _, fleet := range fleets {
		if fleet.Delete || len(fleet.Waypoints) < 2 {
			continue
		}

		wp1 := fleet.Waypoints[1]
		if wp1.TargetType == MapObjectTypeFleet {
			target := fleetGetter.getFleet(wp1.TargetPlayerNum, wp1.TargetNum)
			nodes[fleet] = &fleetMoveNode{
				fleet:  fleet,
				target: &fleetMoveNode{fleet: target},
			}
		} else {
			graph.fleetsNotTargetingFleets = append(graph.fleetsNotTargetingFleets, fleet)
		}
	}

	for fleet, root := range nodes {
		// fmt.Printf("Checking fleet %d\n", fleet.Num)
		visitedForThisFleet := map[*Fleet]bool{fleet: true}
		fleetNodes := []fleetMoveNode{*root}
		targetNode := nodes[root.target.fleet]
		prevTargetNode := root
		_ = prevTargetNode
		cycle := false
		for targetNode != nil {

			targetFleet := targetNode.fleet
			// fmt.Printf("\t fleet %d -> fleet %d\n", prevTargetNode.fleet.Num, targetFleet.Num)
			if _, found := visitedForThisFleet[targetFleet]; found {
				// we've already visited this
				// fmt.Printf("\t fleet %d already visited\n", targetFleet.Num)
				if targetFleet == fleet {
					// fmt.Printf("\t found cycle\n")
					cycle = true
				}
				break
			}
			fleetNodes = append(fleetNodes, *targetNode)
			visitedForThisFleet[targetFleet] = true

			prevTargetNode = targetNode
			targetNode = nodes[targetNode.target.fleet]
		}
		if cycle {
			if graph.fleetsInTargetCycle == nil {
				graph.fleetsInTargetCycle = make(map[playerObject][]fleetMoveNode)
			}

			graph.fleetsInTargetCycle[playerObjectKey(fleet.PlayerNum, fleet.Num)] = fleetNodes
		} else {
			// we're just targeting a fleet, no cycle
			graph.fleetsTargetingFleets = append(graph.fleetsTargetingFleets, *root)
		}
	}

	// sort results
	slices.SortFunc(graph.fleetsTargetingFleets, func(f1, f2 fleetMoveNode) int {
		if f1.fleet.PlayerNum == f2.fleet.PlayerNum {
			return f1.fleet.Num - f2.fleet.Num
		}
		return f1.fleet.PlayerNum - f2.fleet.PlayerNum
	})

	return graph
}

// getFleetMoveOrder returns two slices, one with fleets in order of movement, another with fleets
// in cycles that should be moved 1 l.y. at a time
func (fm *fleetMove) getFleetMoveOrder(graph fleetMoveGraph) (fleets, cycleFleets []*Fleet) {

	// build a list of fleets in move order
	fleets = make([]*Fleet, 0, len(graph.fleetsNotTargetingFleets)+len(graph.fleetsTargetingFleets))

	// start with fleets not targeting fleets, they move first
	copy(fleets, graph.fleetsNotTargetingFleets)

	for _, node := range graph.fleetsTargetingFleets {
		chain := make([]*Fleet, 0, 2)
		target := node.target
		for target != nil {
			chain = append(chain, target.fleet)
			target = target.target
		}
		// we've built a chain of movement for this fleet, add it in reverse order
		reverse(chain)
		fleets = append(fleets, chain...)
	}

	return fleets, cycleFleets
}
