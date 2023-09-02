package ai

import (
	"math"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
	"golang.org/x/exp/slices"
)

type aiPlayer struct {
	*cs.Player
	cs.PlayerMapObjects
	game         *cs.Game
	techStore    *cs.TechStore
	config       aiPlayerConfig
	client       cs.Orderer
	planetsByNum map[int]*cs.Planet
	fleetsByNum  map[int]*cs.Fleet
}

type aiPlayerConfig struct {
	colonizerPopulationDensity float64
}

func NewAIPlayer(game *cs.Game, techStore *cs.TechStore, player *cs.Player, playerMapObjects cs.PlayerMapObjects) *aiPlayer {
	aiPlayer := aiPlayer{
		Player: player,
		game:   game,
		config: aiPlayerConfig{
			colonizerPopulationDensity: .25, // default to requiring 25% pop density before sending off colonizers
		},
		PlayerMapObjects: playerMapObjects,
		client:           cs.NewOrderer(),
	}

	aiPlayer.buildMaps()

	return &aiPlayer
}

// build maps used for quick lookups for various player objects
func (p *aiPlayer) buildMaps() {
	p.planetsByNum = make(map[int]*cs.Planet, len(p.Planets))
	for _, planet := range p.Planets {
		p.planetsByNum[planet.Num] = planet
	}

	p.fleetsByNum = make(map[int]*cs.Fleet, len(p.Fleets))
	for _, fleet := range p.Fleets {
		p.fleetsByNum[fleet.Num] = fleet
	}

}

// process an AI player's turn
func (ai *aiPlayer) ProcessTurn() error {
	ai.scout()
	ai.scoutPackets()
	ai.colonize()
	ai.layMines()
	return nil
}

// dispatch scouts to unknown planets
func (ai *aiPlayer) scout() {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeScout)
	unknownPlanetsByNum := map[int]cs.PlanetIntel{}
	buildablePlanets := ai.getBuildablePlanets(design.Spec.Mass)

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	// find all idle fleets that have scanners
	scannerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeScout]; contains && fleet.Spec.Scanner {
			if len(fleet.Waypoints) <= 1 {
				// this fleet can be sent to scan a planet
				scannerFleets = append(scannerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
						delete(unknownPlanetsByNum, wp.TargetNum)
					}
				}
			}
		}
	}

	for _, fleet := range scannerFleets {
		closestPlanet := ai.getClosestPlanetIntel(fleet.Position, unknownPlanetsByNum)
		if closestPlanet != nil {
			warpSpeed := ai.getWarpSpeed(fleet, closestPlanet.Position)
			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, warpSpeed))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(unknownPlanetsByNum, closestPlanet.Num)
		}
	}

	for _, planet := range buildablePlanets {
		existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool { return item.DesignNum == design.Num })
		if existingQueueItemIndex == -1 {
			// put a new scout at the front of the queue
			planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num}}, planet.ProductionQueue...)
			ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders)

		}
	}
}

// fling packets at planets to scout
func (ai *aiPlayer) scoutPackets() {
	if !ai.Player.Race.Spec.PacketBuiltInScanner {
		return
	}
	unknownPlanetsByNum := map[int]cs.PlanetIntel{}

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Unexplored() {
			unknownPlanetsByNum[planet.Num] = planet
		}
	}

	// remove any planets we're already targetting
	for _, packet := range ai.MineralPackets {
		// this packet travels along a path to the target
		// build a rectangle in the shape of the scan path
		//
		angle := math.Atan2(packet.Heading.Y, packet.Heading.X)
		target := ai.getPlanetIntel(packet.TargetPlanetNum)
		dist := packet.Position.DistanceTo(target.Position)

		// the rectangle is the height of the scan range and the width of the distance
		// from the packet to the target
		// i.e
		// rect origin is upper left corner
		//          *-------------
		// packet    |           |  target
		// (x,y)    *|>>>>>>>>>>>|* (x, y)
		//           |           |
		//           -------------
		height := float64(packet.ScanRangePen * 2)
		width := dist
		rect := cs.Rect{X: packet.Position.X, Y: packet.Position.Y - height/2, Width: width, Height: height}
		for num, planet := range unknownPlanetsByNum {
			// we will catch this in our scanners, remove it
			if rect.PointInRotatedRectangle(planet.Position, angle) {
				delete(unknownPlanetsByNum, num)
			}
		}
		delete(unknownPlanetsByNum, packet.TargetPlanetNum)
	}

	for _, planet := range ai.Planets {
		if planet.Spec.HasMassDriver {
			existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool { return item.Type.IsPacket() })
			if existingQueueItemIndex == -1 {

				// find the farthest planet
				farthest := ai.getFarthestPlanetIntel(planet.Position, unknownPlanetsByNum)
				if farthest == nil {
					continue
				}

				// fling a packet with the mineral we have the most of
				cargoType := planet.Cargo.GreatestMineralType()
				queueItemType := cs.QueueItemTypeMixedMineralPacket
				switch cargoType {
				case cs.Ironium:
					queueItemType = cs.QueueItemTypeIroniumMineralPacket
				case cs.Boranium:
					queueItemType = cs.QueueItemTypeBoraniumMineralPacket
				case cs.Germanium:
					queueItemType = cs.QueueItemTypeGermaniumMineralPacket
				}

				// Build a new packet targetted towards this planet
				planet.PacketTargetNum = farthest.Num
				planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: queueItemType, Quantity: 1}}, planet.ProductionQueue...)
				delete(unknownPlanetsByNum, farthest.Num)

				ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders)

			}
		}
	}
}

// find all colonizable planets and send colony ships to them
func (ai *aiPlayer) colonize() {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeColonizer)
	colonizablePlanets := map[int]cs.PlanetIntel{}
	buildablePlanets := ai.getBuildablePlanets(design.Spec.Mass)

	// find all the unexplored planets
	for _, planet := range ai.Player.PlanetIntels {
		if planet.Explored() && !planet.Owned() && ai.Player.Race.GetPlanetHabitability(planet.Hab) > 0 {
			colonizablePlanets[planet.Num] = planet
		}
	}

	// find all idle fleets that are colonizers
	colonizerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeColonizer]; contains && fleet.Spec.Colonizer {
			if len(fleet.Waypoints) <= 1 {
				planet := ai.getPlanet(fleet.OrbitingPlanetNum)
				if planet != nil && planet.OwnedBy(ai.Player.Num) && planet.Spec.PopulationDensity > ai.config.colonizerPopulationDensity {
					// this fleet can be sent to colonize a planet
					colonizerFleets = append(colonizerFleets, fleet)
				}
			} else {
				// this fleet is already scanning a planet, remove the target from the colonizable planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
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
			// load colonists
			if err := ai.client.TransferPlanetCargo(&ai.game.Rules, ai.Player, fleet, ai.getPlanet(fleet.OrbitingPlanetNum), cs.Cargo{Colonists: fleet.Spec.CargoCapacity}); err != nil {
				// something went wrong, skipi this planet
				log.Error().Err(err).Msg("transferring colonists from planet, skipping")
				continue
			}

			warpSpeed := ai.getWarpSpeed(fleet, bestPlanet.Position)
			fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(bestPlanet.Position, bestPlanet.Num, bestPlanet.Name, warpSpeed).WithTask(cs.WaypointTaskColonize))
			ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
			delete(colonizablePlanets, bestPlanet.Num)
		}
	}

	for _, planet := range buildablePlanets {

		if planet.Spec.PopulationDensity >= ai.config.colonizerPopulationDensity {
			existingQueueItemIndex := slices.IndexFunc(planet.ProductionQueue, func(item cs.ProductionQueueItem) bool { return item.DesignNum == design.Num })
			if existingQueueItemIndex == -1 {
				// put a new scout at the front of the queue
				planet.ProductionQueue = append([]cs.ProductionQueueItem{{Type: cs.QueueItemTypeShipToken, Quantity: 1, DesignNum: design.Num}}, planet.ProductionQueue...)
				ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders)
			}
		}
	}
}

func (ai *aiPlayer) layMines() {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeDamageMineLayer)

	// no mine layers
	if design == nil {
		return
	}

	// lay mines to protect all our planets
	planetsToProtect := make([]*cs.Planet, len(ai.Planets))
	copy(planetsToProtect, ai.Planets)
	planetsToProtectByNum := map[int]*cs.Planet{}
	for _, planet := range ai.Planets {
		planetsToProtectByNum[planet.Num] = planet
	}

	// find all idle fleets that have scanners
	mineLayerFleets := []*cs.Fleet{}
	for _, fleet := range ai.Fleets {
		if _, contains := fleet.Spec.Purposes[cs.ShipDesignPurposeDamageMineLayer]; contains && fleet.Spec.MineLayingRateByMineType != nil {
			if len(fleet.Waypoints) <= 1 && fleet.Waypoints[0].Task != cs.WaypointTaskLayMineField {
				// this fleet can be sent to scan a planet
				mineLayerFleets = append(mineLayerFleets, fleet)
			} else {
				// this fleet is already scanning a planet, remove the target from the unknown planets list
				for _, wp := range fleet.Waypoints[1:] {
					if wp.TargetNum != cs.None {
						delete(planetsToProtectByNum, wp.TargetNum)
					}
				}
			}
		}
	}

	for _, fleet := range mineLayerFleets {
		closestPlanet := ai.getClosestPlanet(fleet, planetsToProtectByNum)
		if closestPlanet != nil {
			if fleet.Position == closestPlanet.Position {
				fleet.Waypoints[0].Task = cs.WaypointTaskLayMineField
				fleet.Waypoints[0].LayMineFieldDuration = cs.Indefinite
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			} else {
				warpSpeed := ai.getWarpSpeed(fleet, closestPlanet.Position)
				fleet.Waypoints = append(fleet.Waypoints, cs.NewPlanetWaypoint(closestPlanet.Position, closestPlanet.Num, closestPlanet.Name, warpSpeed))
				ai.client.UpdateFleetOrders(ai.Player, fleet, fleet.FleetOrders)
				delete(planetsToProtectByNum, closestPlanet.Num)
			}
		}
	}

	// TODO: build new minelayer fleets when some conditions are true...
	buildablePlanets := ai.getBuildablePlanets(design.Spec.Mass)
	_ = buildablePlanets

}

func (ai *aiPlayer) getWarpSpeed(fleet *cs.Fleet, position cs.Vector) int {
	dist := fleet.Position.DistanceTo(position)
	fuelUsage := fleet.GetFuelCost(ai.Player, fleet.Spec.Engine.IdealSpeed, dist)
	warpSpeed := fleet.Spec.Engine.IdealSpeed
	for ; fuelUsage > fleet.Fuel && warpSpeed > 1; warpSpeed-- {
		fuelUsage = fleet.GetFuelCost(ai.Player, warpSpeed, dist)
	}

	return warpSpeed
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) getPlanet(num int) *cs.Planet {
	return p.planetsByNum[num]
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) getPlanetIntel(num int) cs.PlanetIntel {
	return p.Player.PlanetIntels[num-1]
}

// get all planets the player owns that can build ships of mass mass
func (p *aiPlayer) getBuildablePlanets(mass int) []*cs.Planet {
	planets := []*cs.Planet{}
	for _, planet := range p.Planets {
		if planet.CanBuild(mass) {
			planets = append(planets, planet)
		}
	}
	return planets
}

// get the closest planet to this fleet from a list of unknown planets
func (ai *aiPlayer) getClosestPlanetIntel(position cs.Vector, planetIntelsByNum map[int]cs.PlanetIntel) *cs.PlanetIntel {
	shortestDist := math.MaxFloat64
	var closest *cs.PlanetIntel = nil

	for num := range planetIntelsByNum {
		intel := planetIntelsByNum[num]

		distSquared := position.DistanceSquaredTo(intel.Position)
		if shortestDist > distSquared {
			shortestDist = distSquared
			closest = &intel
		}
	}

	return closest
}

// get the farthest planet to this fleet from a list of unknown planets
func (ai *aiPlayer) getFarthestPlanetIntel(position cs.Vector, planetIntelsByNum map[int]cs.PlanetIntel) *cs.PlanetIntel {
	var longestDistance float64 = -1
	var farthest *cs.PlanetIntel = nil

	for _, intel := range planetIntelsByNum {
		distSquared := position.DistanceSquaredTo(intel.Position)
		if longestDistance < distSquared {
			longestDistance = distSquared
			farthest = &intel
		}
	}

	return farthest
}

// get the closest planet to this fleet from a list of unknown planets
func (ai *aiPlayer) getClosestPlanet(fleet *cs.Fleet, planetsByNum map[int]*cs.Planet) *cs.Planet {
	shortestDist := math.MaxFloat64
	var closest *cs.Planet = nil

	for _, planet := range planetsByNum {
		distSquared := fleet.Position.DistanceSquaredTo(planet.Position)
		if shortestDist > distSquared {
			shortestDist = distSquared
			closest = planet
		}
	}

	return closest
}

// get the planet with the highest hab value
func (ai *aiPlayer) getHighestHabPlanet(colonizablePlanets map[int]cs.PlanetIntel) *cs.PlanetIntel {
	highestHabValue := math.MinInt
	var best *cs.PlanetIntel = nil

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
