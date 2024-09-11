package ai

import (
	"fmt"
	"math"

	"github.com/sirgwain/craig-stars/cs"
)

type aiPlayer struct {
	*cs.Player
	cs.PlayerMapObjects
	requests
	game                   *cs.Game
	techStore              *cs.TechStore
	config                 playerConfig
	client                 cs.Orderer
	planetsByNum           map[int]*cs.Planet
	fleetsByNum            map[int]*cs.Fleet
	fleetsByPlanetNum      map[int][]*cs.Fleet
	fleetIntelsByPlanetNum map[int][]*cs.FleetIntel
	designsByPurpose       map[cs.ShipDesignPurpose]*cs.ShipDesign
	fleetsByPurpose        map[cs.FleetPurpose]fleet
	targetedPlanets        map[int][]*cs.FleetIntel

	fuelDepotDesign       *cs.ShipDesign
	fortDesign            *cs.ShipDesign
	starbaseQuarterDesign *cs.ShipDesign
	starbaseHalfDesign    *cs.ShipDesign
	starbaseDesign        *cs.ShipDesign
}

type requests struct {
	fleetBuilds map[cs.FleetPurpose]int
}

type fleetBuildRequest struct {
	count    int
	fleet    fleet
	priority int
}

type playerConfig struct {
	colonizerPopulationDensity       float64
	colonistTransportDensity         float64
	invasionFactor                   float64
	fleetProductionCutoff            float64
	bomberProductionCutoff           float64
	minYearsToQueueStarbasePeaceTime int
	minYearsToQueueStarbaseWarTime   int
	minYearsToBuildScanner           int
	minYearsToBuildFort              int
	namesByPurpose                   map[cs.ShipDesignPurpose]string
	researchOrder                    []cs.TechLevel
}

// each AI has a personality that influences decisions
type Personality string

const (
	Neutral    Personality = ""
	Aggressive Personality = "Aggressive"
	Defensive  Personality = "Defensive"
	Sneaky     Personality = "Sneaky"
)

// each AI has a personality that influences decisions
type Stage string

const (
	Start       Stage = ""
	Explore     Stage = "Explore"
	Expand      Stage = "Expand"
	Exploit     Stage = "Exploit"
	Exterminate Stage = "Exterminate"
)

func NewAIPlayer(game *cs.Game, techStore *cs.TechStore, player *cs.Player, playerMapObjects cs.PlayerMapObjects) *aiPlayer {
	aiPlayer := aiPlayer{
		Player:    player,
		game:      game,
		techStore: techStore,
		requests: requests{
			fleetBuilds: make(map[cs.FleetPurpose]int),
		},
		config: playerConfig{
			colonizerPopulationDensity:       .25, // default to requiring 25% pop density before sending off colonizers
			colonistTransportDensity:         .25, // default to requiring 50% pop density before taking colonists from a feeder to a needer
			minYearsToQueueStarbasePeaceTime: 2,   // don't build starbases if it takes over 2 years to build it
			minYearsToQueueStarbaseWarTime:   4,   // don't build starbases if it takes over 2 years to build it
			minYearsToBuildFort:              10,  // if we are being threatened and need to throw up a fort, do it even if it takes a bit
			minYearsToBuildScanner:           1,
			invasionFactor:                   2,  // we invade if we have 2x the colonists to drop
			fleetProductionCutoff:            .5, // don't try and build starbases until we have 50% factories/mines built first
			bomberProductionCutoff:           .9, // don't try and build bombers until we have 90% factories/mines built first
			namesByPurpose: map[cs.ShipDesignPurpose]string{
				cs.ShipDesignPurposeScout:                 "Long Range Scout",
				cs.ShipDesignPurposeColonizer:             "Santa Maria",
				cs.ShipDesignPurposeBomber:                "Bomber",
				cs.ShipDesignPurposeStructureBomber:       "Structure Bomber",
				cs.ShipDesignPurposeSmartBomber:           "Smart Bomber",
				cs.ShipDesignPurposeStartingFighter:       "Stalwart Defender",
				cs.ShipDesignPurposeFighterScout:          "Armed Probe",
				cs.ShipDesignPurposeTorpedoFighter:        "Missiler",
				cs.ShipDesignPurposeBeamFighter:           "Beamer",
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
				cs.ShipDesignPurposeStarbaseQuarter:       "Tiny Base",
				cs.ShipDesignPurposeStarbaseHalf:          "Small Base",
				cs.ShipDesignPurposePacketThrower:         "Flinger",
				cs.ShipDesignPurposeStargater:             "Teleporter",
				cs.ShipDesignPurposeFort:                  "Orbital Fort",
				cs.ShipDesignPurposeStarterColony:         "Starter Colony",
				cs.ShipDesignPurposeFuelDepot:             "Fuel Depot",
			},
			researchOrder: []cs.TechLevel{
				{Propulsion: 2},
				{Biotechnology: 1},
				{Energy: 1},
				{Weapons: 1},
				{Construction: 4}, // destroyers/privateers
				{Electronics: 1},
				{Weapons: 6},
				{Energy: 4, Propulsion: 4, Electronics: 4, Biotechnology: 4},
				{Weapons: 8, Construction: 6},              // frigates/Phaser bazookas
				{Energy: 6, Propulsion: 6, Electronics: 6}, //+7 Terraform
				{Biotechnology: 7},                         //Organic Armor
				{Weapons: 10},                              //CP and Deltas
				{Construction: 9, Propulsion: 7},           //Cruisers/Warp 8 drives
				{Weapons: 12},                              //Jihads
				{Energy: 8, Propulsion: 8, Electronics: 8},
				{Weapons: 16, Construction: 13},               //Battleships/Juggernauts
				{Energy: 12, Propulsion: 12, Electronics: 11}, //Overthruster/SuperBC/LangstonShell
				{Weapons: 20, Construction: 16},               //Dreadnoughts
				{Energy: 14, Propulsion: 14, Electronics: 14}, //Backfill
				{Weapons: 24},                                 //Armageddon
				{Energy: 18, Propulsion: 16, Electronics: 19}, //Battle Nexus, Warp 10 RS
			},
		},
		PlayerMapObjects: playerMapObjects,
		client:           cs.NewOrderer(),
	}

	aiPlayer.buildMaps()

	return &aiPlayer
}

// choose whether to use beamers or torps in combat for the AI
func (ai *aiPlayer) bestWarship() (cs.ShipDesignPurpose, error) {
	beamDesign := ai.designsByPurpose[cs.ShipDesignPurposeBeamFighter]
	torpDesign := ai.designsByPurpose[cs.ShipDesignPurposeTorpedoFighter]
	var err error
	beamDesign.Spec, err = cs.ComputeShipDesignSpec(&ai.game.Rules, ai.TechLevels, ai.Race.Spec, beamDesign)
	if err != nil {
		return cs.ShipDesignPurposeNone, fmt.Errorf("ComputeShipDesignSpec returned error %w for design %s", err, beamDesign.Name)
	}

	torpDesign.Spec, err = cs.ComputeShipDesignSpec(&ai.game.Rules, ai.TechLevels, ai.Race.Spec, torpDesign)
	if err != nil {
		return cs.ShipDesignPurposeNone, fmt.Errorf("ComputeShipDesignSpec returned error %w for design %s", err, torpDesign.Name)
	}

	if beamDesign.Spec.PowerRating > torpDesign.Spec.PowerRating {
		return cs.ShipDesignPurposeBeamFighter, nil
	}

	return cs.ShipDesignPurposeTorpedoFighter, nil
}

// build maps used for quick lookups for various player objects
func (ai *aiPlayer) buildMaps() error {
	ai.planetsByNum = make(map[int]*cs.Planet, len(ai.Planets))
	for _, planet := range ai.Planets {
		ai.planetsByNum[planet.Num] = planet
	}

	ai.fleetsByNum = make(map[int]*cs.Fleet, len(ai.Fleets))
	ai.fleetsByPlanetNum = make(map[int][]*cs.Fleet, len(ai.Planets))
	for _, fleet := range ai.Fleets {
		ai.fleetsByNum[fleet.Num] = fleet
		ai.fleetsByPlanetNum[fleet.OrbitingPlanetNum] = append(ai.fleetsByPlanetNum[fleet.OrbitingPlanetNum], fleet)
	}

	ai.fleetIntelsByPlanetNum = make(map[int][]*cs.FleetIntel, len(ai.PlanetIntels))
	for _, fleet := range ai.FleetIntels {
		ai.fleetIntelsByPlanetNum[fleet.OrbitingPlanetNum] = append(ai.fleetIntelsByPlanetNum[fleet.OrbitingPlanetNum], &fleet)
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

	bestWarshipPurpose, err := ai.bestWarship()
	if err != nil {
		return fmt.Errorf("Could not decide on whether using beams or torps when building maps; spec calc errored %w", err)
	}

	ai.fleetsByPurpose = map[cs.FleetPurpose]fleet{
		cs.FleetPurposeScout: {
			purpose: cs.FleetPurposeScout,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeScout,
					quantity: 1,
				},
			},
		},
		cs.FleetPurposeColonizer: {
			purpose: cs.FleetPurposeColonizer,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeColonistFreighter,
					quantity: 1,
				},
				{
					purpose:  cs.ShipDesignPurposeFuelFreighter,
					quantity: 1,
				},
				{
					purpose:  cs.ShipDesignPurposeColonizer,
					quantity: 1,
				},
			},
		},
		cs.FleetPurposeColonistFreighter: {
			purpose: cs.FleetPurposeColonistFreighter,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeColonistFreighter,
					quantity: 1,
				},
				{
					purpose:  cs.ShipDesignPurposeFuelFreighter,
					quantity: 1,
				},
			},
		},
		cs.FleetPurposeBomber: {
			purpose: cs.FleetPurposeBomber,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeBomber,
					quantity: cs.MinInt(5*(cs.MaxInt((ai.game.Year-25)/5, 1)), 40), // adds 5 ships every 5 years after 25, up to 40 at max
				},
				{
					purpose:  bestWarshipPurpose,
					quantity: cs.MinInt(5*(cs.MaxInt((ai.game.Year-25)/5, 1)), 50), // adds 5 ships every 5 years after 25, up to 50 at max
				},
			},
		},
	}

	ai.targetedPlanets = make(map[int][]*cs.FleetIntel)
	return nil
}

// process an AI player's turn
func (ai *aiPlayer) ProcessTurn() error {
	ai.assignPurpose()
	ai.gatherIntel()
	ai.plan()
	ai.designStarbases()

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
	if err := ai.invade(); err != nil {
		return err
	}
	if err := ai.bomb(); err != nil {
		return err
	}
	if err := ai.updateFleetWarpSpeed(); err != nil {
		return err
	}
	if err := ai.transport(); err != nil {
		return err
	}
	if err := ai.produce(); err != nil {
		return err
	}

	// make sure our research is optimal
	ai.research()
	// cleanup any old designs we haven't built
	ai.removedUnusedDesigns()

	return nil
}

func (ai *aiPlayer) getWarpSpeed(fleet *cs.Fleet, position cs.Vector) int {
	dist := fleet.Position.DistanceTo(position)
	return cs.Clamp(ai.getMaxWarp(dist, fleet), 1, 10)
}

// get the maximum warp we can travel to reach the destination
// in the minimal number of years, within our fuel constraints
func (ai *aiPlayer) getMaxWarp(dist float64, fleet *cs.Fleet) int {
	freeSpeed := fleet.Spec.Engine.FreeSpeed

	// start at freespeed+1 and move up until we run out of fuel
	var speed int
	for speed = freeSpeed + 1; speed < fleet.Spec.Engine.MaxSafeSpeed; speed++ {
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

	return cs.Clamp(speed, 1, 10)
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

// get the closest planet to this fleet from a list of unknown planets
func (ai *aiPlayer) getClosestStarbasePlanet(fleet *cs.Fleet) *cs.Planet {
	shortestDist := math.MaxFloat64
	var closest *cs.Planet = nil

	for _, planet := range ai.Planets {
		if !planet.Spec.HasStarbase {
			continue
		}

		distSquared := fleet.Position.DistanceSquaredTo(planet.Position)
		if shortestDist > distSquared {
			shortestDist = distSquared
			closest = planet
		}
	}

	return closest
}

// check if a starbase is already in the queue
func (ai *aiPlayer) isStarbaseInQueue(planet *cs.Planet) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeStarbase {
			return true
		}
	}
	return false
}

// check if a the queue contains any ships with this fleet purpose
func (ai *aiPlayer) isFleetInQueue(planet *cs.Planet, fleetPurpose cs.FleetPurpose) bool {
	for _, item := range planet.ProductionQueue {
		if item.GetTag(cs.TagPurpose) == string(fleetPurpose) {
			return true
		}
	}
	return false
}

// check if a shipdesign with the given purpose is in the queue
func (ai *aiPlayer) isShipInQueue(planet *cs.Planet, fleetPurpose cs.FleetPurpose, purpose cs.ShipDesignPurpose, quantity int) bool {
	for _, item := range planet.ProductionQueue {
		if item.Type == cs.QueueItemTypeShipToken && item.GetTag(cs.TagPurpose) == string(fleetPurpose) {
			design := ai.GetDesign(item.DesignNum)
			if design != nil && design.Purpose == purpose && item.Quantity >= quantity {
				return true
			}
		}
	}
	return false
}

// get the number of idle ships above a planet, matching a fleet and ship purpose
func (ai *aiPlayer) getIdleShipCount(planet *cs.Planet, fleetPurpose cs.FleetPurpose, purpose cs.ShipDesignPurpose) int {
	count := 0
	for _, fleet := range ai.fleetsByPlanetNum[planet.Num] {
		if !fleet.Idle() {
			continue
		}

		if fleet.GetTag(cs.TagPurpose) != string(fleetPurpose) {
			continue
		}

		for _, token := range fleet.Tokens {
			design := ai.GetDesign(token.DesignNum)
			if design != nil && design.Purpose == purpose {
				count += token.Quantity
			}
		}
	}
	return count
}

// get all the enemy ships above a planet
func (ai *aiPlayer) enemyShipsAbovePlanet(planet *cs.Planet) []*cs.FleetIntel {
	fleets := []*cs.FleetIntel{}

	for _, fleet := range ai.fleetIntelsByPlanetNum[planet.Num] {
		if ai.IsEnemy(fleet.PlayerNum) {
			fleets = append(fleets, fleet)
		}
	}

	return fleets
}

// hasAttackShips returns true if the fleet intels likely contains hostile ships
func (ai *aiPlayer) hasAttackShips(fleets []*cs.FleetIntel) bool {
	for _, fleet := range fleets {
		for _, token := range fleet.Tokens {
			design := ai.GetForeignDesign(fleet.PlayerNum, token.DesignNum)
			if design != nil {
				if design.Spec.PowerRating > 0 {
					return true
				}
				hull := ai.techStore.GetHull(design.Hull)
				if hull != nil && hull.Type.IsAttackHull() {
					// we aren't positive this is hostile without checking slots, but it likely is
					return true
				}
			}
		}
	}
	return false
}
