package game

import (
	"math"

	"golang.org/x/exp/slices"
)

type aiPlayer struct {
	*Player
	PlayerMapObjects
	config       aiPlayerConfig
	planetsByNum map[int]*Planet
	fleetsByNum  map[int]*Fleet
}

type aiPlayerConfig struct {
	colonizerPopulationDensity float64
}

func NewAIPlayer(player *Player, playerMapObjects PlayerMapObjects) *aiPlayer {
	aiPlayer := aiPlayer{
		Player: player,
		config: aiPlayerConfig{
			colonizerPopulationDensity: .25, // default to requiring 25% pop density before sending off colonizers
		},
		PlayerMapObjects: playerMapObjects,
	}

	aiPlayer.buildMaps()

	return &aiPlayer
}

// build maps used for quick lookups for various player objects
func (p *aiPlayer) buildMaps() {
	p.planetsByNum = make(map[int]*Planet, len(p.Planets))
	for _, planet := range p.Planets {
		p.planetsByNum[planet.Num] = planet
	}

	p.fleetsByNum = make(map[int]*Fleet, len(p.Fleets))
	for _, fleet := range p.Fleets {
		p.fleetsByNum[fleet.Num] = fleet
	}

}

// process an AI player's turn
func (ai *aiPlayer) processTurn() {
	ai.scout()
	ai.colonize()
}

// dispatch scouts to unknown planets
func (ai *aiPlayer) scout() {
	design := ai.Player.GetLatestDesign(ShipDesignPurposeScout)
	unknownPlanetsByNum := map[int]PlanetIntel{}
	buildablePlanets := ai.GetBuildablePlanets(design.Spec.Mass)

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	// find all idle fleets that have scanners
	scannerFleets := []*Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[ShipDesignPurposeScout]; contains && fleet.Spec.Scanner {
			if len(fleet.Waypoints) <= 1 {
				// this fleet can be sent to scan a planet
				scannerFleets = append(scannerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != None {
						delete(unknownPlanetsByNum, wp.TargetNum)
					}
				}
			}
		}
	}

	for _, fleet := range scannerFleets {
		closestPlanet := ai.getClosestPlanet(fleet, unknownPlanetsByNum)
		if closestPlanet != nil {
			fleet.Waypoints = append(fleet.Waypoints, NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, fleet.Spec.IdealSpeed))
			fleet.Dirty = true
			delete(unknownPlanetsByNum, closestPlanet.Num)
		}
	}

	for _, planet := range buildablePlanets {
		existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item ProductionQueueItem) bool { return item.DesignName == design.Name })
		if existingQueueItemIndex == -1 {
			// put a new scout at the front of the queue
			planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignName: design.Name}}, planet.ProductionQueue...)
			planet.Dirty = true
		}
	}
}

// find all colonizable planets and send colony ships to them
func (ai *aiPlayer) colonize() {
	design := ai.Player.GetLatestDesign(ShipDesignPurposeColonizer)
	colonizablePlanets := map[int]PlanetIntel{}
	buildablePlanets := ai.GetBuildablePlanets(design.Spec.Mass)

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Explored() && ai.Player.Race.GetPlanetHabitability(planet.Hab) > 0 {
			colonizablePlanets[planet.Num] = planet
		}
	}

	// find all idle fleets that are colonizers
	colonizerFleets := []*Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[ShipDesignPurposeColonizer]; contains && fleet.Spec.Colonizer {
			if len(fleet.Waypoints) <= 1 {
				planet := ai.GetPlanet(fleet.OrbitingPlanetNum)
				if planet != nil && planet.OwnedBy(ai.Player.Num) && planet.Spec.PopulationDensity > ai.config.colonizerPopulationDensity {
					// this fleet can be sent to colonize a planet
					colonizerFleets = append(colonizerFleets, fleet)
				}
			} else {
				// this fleet is already scanning a planet, remove the target from the colonizable planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != None {
						delete(colonizablePlanets, wp.TargetNum)
					}
				}
			}
		}
	}

	// TODO: we should sort by best planet and find the closest fleet to that planet to send a colonizer
	for _, fleet := range colonizerFleets {
		bestPlanet := ai.getHighestHabPlanet(colonizablePlanets)
		if bestPlanet != nil {
			fleet.Waypoints[0].Task = WaypointTaskTransport
			fleet.Waypoints[0].TransportTasks = WaypointTransportTasks{
				Colonists: WaypointTransportTask{
					Action: TransportActionLoadAll,
				},
			}

			// todo: remove after transports are complete
			fleet.Cargo.Colonists = fleet.Spec.CargoCapacity

			fleet.Waypoints = append(fleet.Waypoints, NewPlanetWaypoint(bestPlanet.Position, bestPlanet.Num, bestPlanet.Name, fleet.Spec.IdealSpeed).WithTask(WaypointTaskColonize))
			fleet.Dirty = true
			delete(colonizablePlanets, bestPlanet.Num)
		}
	}

	for _, planet := range buildablePlanets {
		existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item ProductionQueueItem) bool { return item.DesignName == design.Name })
		if existingQueueItemIndex == -1 {
			// put a new scout at the front of the queue
			planet.ProductionQueue = append([]ProductionQueueItem{{Type: QueueItemTypeShipToken, Quantity: 1, DesignName: design.Name}}, planet.ProductionQueue...)
			planet.Dirty = true
		}
	}
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) GetPlanet(num int) *Planet {
	return p.planetsByNum[num]
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) GetFleet(num int) *Fleet {
	return p.fleetsByNum[num]
}

// get all planets the player owns that can build ships of mass mass
func (p *aiPlayer) GetBuildablePlanets(mass int) []*Planet {
	planets := []*Planet{}
	for _, planet := range p.Planets {
		if planet.CanBuild(mass) {
			planets = append(planets, planet)
		}
	}
	return planets
}

// get the closest planet to this fleet from a list of unknown planets
func (ai *aiPlayer) getClosestPlanet(fleet *Fleet, unknownPlanetsByNum map[int]PlanetIntel) *PlanetIntel {
	shortestDist := math.MaxFloat64
	var closest *PlanetIntel = nil

	for num := range unknownPlanetsByNum {
		intel := unknownPlanetsByNum[num]
		distSquared := fleet.Position.DistanceSquaredTo(intel.Position)
		if shortestDist > distSquared {
			shortestDist = distSquared
			closest = &intel
		}
	}

	return closest
}

// get the planet with the highest hab value
func (ai *aiPlayer) getHighestHabPlanet(colonizablePlanets map[int]PlanetIntel) *PlanetIntel {
	highestHabValue := math.MinInt
	var best *PlanetIntel = nil

	for num := range colonizablePlanets {
		intel := colonizablePlanets[num]
		habValue := ai.Player.Race.GetPlanetHabitability(intel.Hab)
		if habValue > highestHabValue {
			highestHabValue = habValue
			best = &intel
		}
	}

	return best
}
