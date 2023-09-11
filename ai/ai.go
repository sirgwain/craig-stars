package ai

import (
	"math"

	"github.com/sirgwain/craig-stars/cs"
	"golang.org/x/exp/slices"
)

type aiPlayer struct {
	*cs.Player
	cs.PlayerMapObjects
	requests
	game             *cs.Game
	techStore        *cs.TechStore
	config           playerConfig
	client           cs.Orderer
	planetsByNum     map[int]*cs.Planet
	fleetsByNum      map[int]*cs.Fleet
	designsByPurpose map[cs.ShipDesignPurpose]*cs.ShipDesign
	fleetsByPurpose  map[fleetPurpose]fleet
}

type requests struct {
	shipBuilds map[cs.ShipDesignPurpose]int
}

type playerConfig struct {
	colonizerPopulationDensity float64
	minYearsToQueueStarbase    int
	namesByPurpose             map[cs.ShipDesignPurpose]string
}

func NewAIPlayer(game *cs.Game, techStore *cs.TechStore, player *cs.Player, playerMapObjects cs.PlayerMapObjects) *aiPlayer {
	aiPlayer := aiPlayer{
		Player:    player,
		game:      game,
		techStore: techStore,
		requests: requests{
			shipBuilds: make(map[cs.ShipDesignPurpose]int),
		},
		config: playerConfig{
			colonizerPopulationDensity: .25, // default to requiring 25% pop density before sending off colonizers
			minYearsToQueueStarbase:    2,   // don't build starbases if it takes over 2 years to build it
			namesByPurpose: map[cs.ShipDesignPurpose]string{
				cs.ShipDesignPurposeScout:                 "Long Range Scout",
				cs.ShipDesignPurposeColonizer:             "Santa Maria",
				cs.ShipDesignPurposeBomber:                "Bomber",
				cs.ShipDesignPurposeStructureBomber:       "Structure Bomber",
				cs.ShipDesignPurposeSmartBomber:           "Smart Bomber",
				cs.ShipDesignPurposeFighter:               "Stalwart Defender",
				cs.ShipDesignPurposeFighterScout:          "Armed Probe",
				cs.ShipDesignPurposeCapitalShip:           "Warship",
				cs.ShipDesignPurposeFreighter:             "Teamster",
				cs.ShipDesignPurposeColonistFreighter:     "Colonist Freighter",
				cs.ShipDesignPurposeFuelFreighter:         "Fuel Freighter",
				cs.ShipDesignPurposeMultiPurposeFreighter: "Swashbuckler",
				cs.ShipDesignPurposeArmedFreighter:        "Armed Freighter",
				cs.ShipDesignPurposeMiner:                 "Cotton Picker",
				cs.ShipDesignPurposeTerraformer:           "Change of Heart",
				cs.ShipDesignPurposeDamageMineLayer:       "Little Hen",
				cs.ShipDesignPurposeSpeedMineLayer:        "Speed Turtle",
				cs.ShipDesignPurposeStarbase:              "Starbase",
				cs.ShipDesignPurposePacketThrower:         "Flinger",
				cs.ShipDesignPurposeStargater:             "Teleporter",
				cs.ShipDesignPurposeFort:                  "Orbital Fort",
				cs.ShipDesignPurposeStarterColony:         "Starter Colony",
				cs.ShipDesignPurposeFuelDepot:             "Fuel Depot",
			},
		},
		PlayerMapObjects: playerMapObjects,
		client:           cs.NewOrderer(),
	}

	aiPlayer.buildMaps()

	return &aiPlayer
}

// build maps used for quick lookups for various player objects
func (ai *aiPlayer) buildMaps() {
	ai.planetsByNum = make(map[int]*cs.Planet, len(ai.Planets))
	for _, planet := range ai.Planets {
		ai.planetsByNum[planet.Num] = planet
	}

	ai.fleetsByNum = make(map[int]*cs.Fleet, len(ai.Fleets))
	for _, fleet := range ai.Fleets {
		ai.fleetsByNum[fleet.Num] = fleet
	}

	ai.designsByPurpose = make(map[cs.ShipDesignPurpose]*cs.ShipDesign, len(ai.Designs))
	for _, design := range ai.Designs {
		if existing, ok := ai.designsByPurpose[design.Purpose]; ok {
			// add latest version
			if existing.Version < design.Version {
				ai.designsByPurpose[design.Purpose] = design
			}
		} else {
			ai.designsByPurpose[design.Purpose] = design
		}
	}

	// TODO: use this. :)
	ai.fleetsByPurpose = map[fleetPurpose]fleet{
		scout: {
			purpose: scout,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeScout,
					quantity: 1,
					hull:     ai.getHull(scout),
				},
			},
		},
		colonizer: {
			purpose: colonizer,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeColonizer,
					quantity: 1,
					hull:     ai.getHull(colonizer),
				},
				{
					purpose:  cs.ShipDesignPurposeColonistFreighter,
					quantity: 2,
					hull:     ai.getHull(transport),
				},
			},
		},
	}
}

// process an AI player's turn
func (ai *aiPlayer) ProcessTurn() error {
	if err := ai.scout(); err != nil {
		return err
	}
	if err := ai.scoutPackets(); err != nil {
		return err
	}
	if err := ai.colonize(); err != nil {
		return err
	}
	if err := ai.layMines(); err != nil {
		return err
	}
	if err := ai.transport(); err != nil {
		return err
	}
	if err := ai.produce(); err != nil {
		return err
	}
	return nil
}

// dispatch scouts to unknown planets
func (ai *aiPlayer) scout() error {
	unknownPlanetsByNum := map[int]cs.PlanetIntel{}

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

						target := ai.getPlanetIntel(wp.TargetNum)
						if target.ReportAge != cs.ReportAgeUnexplored {
							// we discovered this target in some other way, change targets
							fleet.Waypoints = fleet.Waypoints[:1]
							scannerFleets = append(scannerFleets, fleet)
							continue
						}
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

	// build scout ships where necessary
	if len(unknownPlanetsByNum) > 0 {
		ai.addShipBuildRequest(cs.ShipDesignPurposeScout, len(unknownPlanetsByNum))
	}

	return nil
}

// fling packets at planets to scout
func (ai *aiPlayer) scoutPackets() error {
	if !ai.Player.Race.Spec.PacketBuiltInScanner {
		return nil
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

				if err := ai.client.UpdatePlanetOrders(&ai.game.Rules, ai.Player, planet, planet.PlanetOrders); err != nil {
					return err
				}

			}
		}
	}
	return nil
}

func (ai *aiPlayer) layMines() error {
	design := ai.Player.GetLatestDesign(cs.ShipDesignPurposeDamageMineLayer)

	// no mine layers
	if design == nil {
		return nil
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

	return nil
}

func (ai *aiPlayer) getWarpSpeed(fleet *cs.Fleet, position cs.Vector) int {
	dist := fleet.Position.DistanceTo(position)
	return ai.getMaxWarp(dist, fleet)
}

// get the maximum warp we can travel to reach the destination
// in the minimal number of years, within our fuel constraints
func (ai *aiPlayer) getMaxWarp(dist float64, fleet *cs.Fleet) int {
	freeSpeed := fleet.Spec.Engine.FreeSpeed

	// start at freespeed+1 and move up until we run out of fuel
	var speed int
	for speed = freeSpeed + 1; speed <= fleet.Spec.Engine.MaxSafeSpeed; speed++ {
		fuelUsed := fleet.GetFuelCost(ai.Player, speed, dist)

		// we are using too much fuel, go to the previous speed
		if fuelUsed > fleet.Fuel {
			speed--
			break
		}
	}

	idealSpeed := fleet.Spec.Engine.IdealSpeed

	// if we are using a ramscoop, make sure we at least go the ideal
	// speed of the engine. If we run out, oh well, it'll drop to
	// the free speed
	if freeSpeed > 1 && speed < idealSpeed {
		speed = idealSpeed
	}

	// don't go faster than we need
	return ai.getMinimalWarp(dist, speed, fleet)
}

// get the minimal warp starting at an idealSpeed and working downward
// if we can travel the same amount of time at a lower speed, do it
func (ai *aiPlayer) getMinimalWarp(dist float64, idealSpeed int, fleet *cs.Fleet) int {
	speed := idealSpeed

	freeSpeed := fleet.Spec.Engine.FreeSpeed

	// travelling 49 ly at warp 7 takes one year
	yearsAtIdealSpeed := int(math.Ceil(dist / float64(idealSpeed*idealSpeed)))
	for i := idealSpeed; i > freeSpeed; i-- {
		yearsAtSpeed := int(math.Ceil(dist / float64(i*i)))

		// It takes the same time to go slower, so go slower
		if yearsAtIdealSpeed == yearsAtSpeed {
			speed = i
		}
	}

	return speed
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) getPlanet(num int) *cs.Planet {
	return p.planetsByNum[num]
}

// get a player owned planet by num, or nil if it doesn't exist
func (p *aiPlayer) getPlanetIntel(num int) cs.PlanetIntel {
	return p.Player.PlanetIntels[num-1]
}

// get all planets we own with space docks
func (p *aiPlayer) getPlanetsWithDocks() []*cs.Planet {
	planets := []*cs.Planet{}
	for _, planet := range p.Planets {
		if planet.Spec.HasStarbase && planet.Spec.DockCapacity != 0 {
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

func (ai *aiPlayer) isStarbaseInQueue(planet *cs.Planet) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeStarbase {
			return true
		}
	}
	return false
}
