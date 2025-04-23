package ai

import (
	"fmt"
	"math"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/cs"
)

type aiPlayer struct {
	*cs.Player
	cs.PlayerMapObjects
	requests
	log                    zerolog.Logger
	game                   *cs.Game
	techStore              *cs.TechStore
	config                 playerConfig // TODO: Make config vary with race
	client                 cs.Orderer
	planetsByNum           map[int]*cs.Planet
	fleetsByNum            map[int]*cs.Fleet
	fleetsByPlanetNum      map[int][]*cs.Fleet
	fleetIntelsByPlanetNum map[int][]*cs.FleetIntel
	designsByPurpose       map[cs.ShipDesignPurpose]*cs.ShipDesign
	fleetsByPurpose        map[cs.FleetPurpose]fleet
	targetedPlanets        map[int][]*cs.FleetIntel
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
	colonizerPopulationDensity float64
	colonistTransportDensity   float64
	invasionFactor             float64
	// TODO: Make these a map of cutoffs per fleet type - colonizers are less of an expense than scouts
	fleetProductionCutoff            float64
	bomberProductionCutoff           float64
	minYearsToQueueStarbasePeaceTime int
	minYearsToQueueStarbaseWarTime   int
	minYearsToBuildScanner           int
	minYearsToBuildFort              int
	// year to start prioritizing about minerals over resources for warship building
	// TODO: Make this some sort of map[TechLevel]CostFloat64 containing weights for each stage of the game
	mineralConservationYear int
	startAttackingYear      int // year to start launching attacks on players
	namesByPurpose          map[cs.ShipDesignPurpose]string
	researchOrder           []cs.TechLevel
}

// each AI has a personality that influences decisions; currently WIP
type Personality string

const (
	Neutral    Personality = ""
	Aggressive Personality = "Aggressive"
	Defensive  Personality = "Defensive"
	Sneaky     Personality = "Sneaky"
)

// The stage of the ai's game plan; currently WIP
type Stage string

const (
	Start       Stage = ""
	Explore     Stage = "Explore"
	Expand      Stage = "Expand"
	Exploit     Stage = "Exploit"
	Exterminate Stage = "Exterminate"
)

type warshipCount struct {
	bombers        int
	warships       int
	fuelTransports int
}

type aiCutoffsByStartMode struct {
	attackYear  map[cs.GameStartMode]int // min year to start attacking
	mineralYear map[cs.GameStartMode]int // min year to care about minerals
}

// TODO: Make these cutoffs dynamic and configurable on a per-race basis
var defaultCutoffs = aiCutoffsByStartMode{
	attackYear: map[cs.GameStartMode]int{
		cs.GameStartModeNormal: 25,
		cs.GameStartModeAccBBS: 20,
		cs.GameStartModeMax:    0,
	},
	mineralYear: map[cs.GameStartMode]int{
		cs.GameStartModeNormal: 55,
		cs.GameStartModeAccBBS: 50,
		cs.GameStartModeMax:    0,
	},
}

// Create a new AI player
func NewAIPlayer(game *cs.Game, techStore *cs.TechStore, player *cs.Player, playerMapObjects cs.PlayerMapObjects) *aiPlayer {
	aiPlayer := aiPlayer{
		Player: player,
		game:   game,
		log: log.With().
			Int64("GameID", game.ID).
			Str("Game", game.Name).
			Int("PlayerNum", player.Num).
			Str("Player", player.Name).
			Int("Year", game.Year). // @sirgwain Should this be year or year+1?
			Logger(),
		techStore: techStore,
		requests: requests{
			fleetBuilds: make(map[cs.FleetPurpose]int),
		},
		config: playerConfig{
			// TODO: Make below configurable with AI difficulty/aggression mode
			colonizerPopulationDensity:       .25,                                        // default to requiring 25% pop density before sending off colonizers
			colonistTransportDensity:         .25,                                        // default to requiring 25% pop density before taking colonists from a feeder to a needer
			mineralConservationYear:          defaultCutoffs.mineralYear[game.StartMode], // TODO: Make this account for the AI's current progress in game (perhaps with a tech lvl cutoff or similar)
			minYearsToQueueStarbasePeaceTime: 2,                                          // only build starbases if it takes <=2 years to build it
			minYearsToQueueStarbaseWarTime:   4,                                          // only build starbases if it takes <=4 years to build it and the planet is threatened
			minYearsToBuildFort:              10,                                         // only build emergency panic forts if it takes <=10 years to build it
			minYearsToBuildScanner:           1,                                          // only build planetary scanners if we can finish it in 1 year
			invasionFactor:                   2,                                          // only invade if we have 2x the colonists to drop
			fleetProductionCutoff:            .5,                                         // don't try and build ships until we have 50% factories/mines built first
			bomberProductionCutoff:           .9,                                         // don't try and build bombers until we have 90% factories/mines built first
			startAttackingYear:               defaultCutoffs.attackYear[game.StartMode],
			namesByPurpose: map[cs.ShipDesignPurpose]string{
				// TODO: make this return a slice of strings/structs to allow for name variety
				cs.ShipDesignPurposeScout:                 "Long Range Scout",
				cs.ShipDesignPurposeColonizer:             "Santa Maria",
				cs.ShipDesignPurposeBomber:                "Bomber",
				cs.ShipDesignPurposeStructureBomber:       "Structure Bomber",
				cs.ShipDesignPurposeSmartBomber:           "Smart Bomber",
				cs.ShipDesignPurposeStartingFighter:       "Stalwart Defender",
				cs.ShipDesignPurposeFighterScout:          "Armed Probe",
				cs.ShipDesignPurposeTorpedoFighter:        "Missile Boat",
				cs.ShipDesignPurposeBeamFighter:           "Beam Ship",
				cs.ShipDesignPurposeFreighter:             "Teamster",
				cs.ShipDesignPurposeColonistFreighter:     "Colonist Freighter",
				cs.ShipDesignPurposeFuelFreighter:         "Fuel Freighter",
				cs.ShipDesignPurposeMultiPurposeFreighter: "Swashbuckler",
				cs.ShipDesignPurposeArmedFreighter:        "Blackbeard",
				cs.ShipDesignPurposeMiner:                 "Cotton Picker",
				cs.ShipDesignPurposeTerraformer:           "Change of Heart",
				cs.ShipDesignPurposeDamageMineLayer:       "Little Hen",
				cs.ShipDesignPurposeSpeedMineLayer:        "Speed Turtle",
				cs.ShipDesignPurposeStarbase:              "Starbase",
				cs.ShipDesignPurposeStarbaseUnarmed:       "Holder of Place",
				cs.ShipDesignPurposeStarbaseQuarter:       "Tiny Base",
				cs.ShipDesignPurposeStarbaseHalf:          "Small Base",
				cs.ShipDesignPurposePacketThrower:         "Flinger",
				cs.ShipDesignPurposeStargater:             "Gateway",
				cs.ShipDesignPurposeFort:                  "Bunker",
				cs.ShipDesignPurposeStarterColony:         "Starter Colony",
				cs.ShipDesignPurposeFuelDepot:             "Fuel Depot",
			},
			researchOrder: []cs.TechLevel{
				// TODO: Make this a race-specific trait to allow for funky AI races with different research paths
				{Propulsion: 2}, // FM
				{Biotechnology: 1},
				{Energy: 1},
				{Weapons: 1},
				{Construction: 4}, // destroyers/privateers
				{Electronics: 1},
				{Weapons: 6, Biotechnology: 2},             // yaks & +7 weps terraform
				{Energy: 3, Weapons: 8, Construction: 6},   // shielded frigates & Phaser bazookas
				{Energy: 5, Propulsion: 5},                 // +7 temp/grav terraforming
				{Construction: 8, Electronics: 3},          // better scanners + LFs
				{Energy: 6, Weapons: 10, Biotechnology: 3}, // CP and Deltas
				{Construction: 10, Propulsion: 7},          // (Battle) Cruisers/Warp 8 drives
				{Weapons: 12},                              // Jihads
				{Energy: 10, Propulsion: 10, Electronics: 7, Biotechnology: 4}, // organic + better engines
				{Construction: 13}, //Battleships
				{Weapons: 16},      // Juggernauts
				{Energy: 12, Propulsion: 12, Electronics: 11},    //Overthruster/SuperBC/LangstonShell
				{Weapons: 20, Construction: 16},                  //Dreadnoughts
				{Energy: 14, Electronics: 14, Biotechnology: 10}, // Gorilla delagators, mega poly
				{Weapons: 24}, // Armageddon Missiles
				{Energy: 18, Propulsion: 16, Electronics: 19}, // Battle Nexus, Warp 10 RS
				{Weapons: 26},      // omega torps
				{Construction: 26}, // nubians
			},
		},
		PlayerMapObjects: playerMapObjects,
		client:           cs.NewOrderer(),
	}

	aiPlayer.buildMaps()

	return &aiPlayer
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
					quantity: 3,
				},
				{
					purpose:  cs.ShipDesignPurposeFuelFreighter,
					quantity: 2,
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
					quantity: 2,
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
					quantity: 5,
				},
				{
					purpose:  cs.ShipDesignPurposeBeamFighter,
					quantity: 7,
				},
				{
					purpose:  cs.ShipDesignPurposeTorpedoFighter,
					quantity: 7,
				},
			},
		},
		cs.FleetPurposeCapitalShip: {
			purpose: cs.FleetPurposeCapitalShip,
			ships: []fleetShip{
				{
					purpose:  cs.ShipDesignPurposeBeamFighter,
					quantity: 7,
				},
				{
					purpose:  cs.ShipDesignPurposeTorpedoFighter,
					quantity: 7,
				},
			},
		},
	}

	ai.targetedPlanets = make(map[int][]*cs.FleetIntel)
	return nil
}

// update an AI player's warfleet quantities for attack/defense fleets
func (ai *aiPlayer) updateWarfleets() (err error) {
	// grab baseline warship counts for current year
	warshipCount := ai.getWarshipCount()

	// get our warship designs, updating the spec if needed
	beamDesign := ai.designsByPurpose[cs.ShipDesignPurposeBeamFighter]
	if beamDesign == nil {
		// if design is nil, try to make one from scratch
		ai.designsByPurpose[cs.ShipDesignPurposeBeamFighter], err = ai.designShip(ai.config.namesByPurpose[cs.ShipDesignPurposeBeamFighter], cs.ShipDesignPurposeBeamFighter, cs.FleetPurposeFighter) // fleet purpose unimportant as it's just used for radrams
		if err != nil {
			return fmt.Errorf("error designing beam fighter during warship quantity updating: %w", err)
		}
		beamDesign = ai.designsByPurpose[cs.ShipDesignPurposeBeamFighter]
	}
	torpDesign := ai.designsByPurpose[cs.ShipDesignPurposeTorpedoFighter]
	if torpDesign == nil {
		// if design is nil, try to make one from scratch
		ai.designsByPurpose[cs.ShipDesignPurposeTorpedoFighter], err = ai.designShip(ai.config.namesByPurpose[cs.ShipDesignPurposeTorpedoFighter], cs.ShipDesignPurposeTorpedoFighter, cs.FleetPurposeFighter) // fleet purpose unimportant as it's just used for radrams
		if err != nil {
			return fmt.Errorf("error designing torpedo fighter during warship quantity updating: %w", err)
		}
		torpDesign = ai.designsByPurpose[cs.ShipDesignPurposeTorpedoFighter]
	}

	// if design is STILL nil, assume we can't make a design of that type
	// if 1 design exists and the other doesn't, automatically use it and exit
	if beamDesign == nil {
		if torpDesign == nil {
			ai.log.Debug().Msgf("Skipping over choosing warship quantities due to nil designs")
		} else {
			ai.updateWarshipAmounts(warshipCount.bombers, 0, warshipCount.warships, warshipCount.fuelTransports)
		}
		return nil
	} else if torpDesign == nil {
		ai.updateWarshipAmounts(warshipCount.bombers, warshipCount.warships, 0, warshipCount.fuelTransports)
		return nil
	}

	// Compare ships' power rating and overall mineral/res expenditure
	// TODO: Add a less jank way of evaluating warship performance than ranking
	// and make the AI consider how much spare minerals it has
	scoreRatio := float64(beamDesign.Spec.PowerRating) / float64(torpDesign.Spec.PowerRating)

	ct := []cs.CostType{cs.Resources}
	if ai.game.YearsPassed() >= ai.config.mineralConservationYear {
		// care about minerals over resources after year 55 (50 for accBBS)
		ct = cs.MineralTypes[:]
	}
	costRatio := cs.GetCostEfficiencyRatio(beamDesign.Spec.Cost.ToCostFloat64(), torpDesign.Spec.Cost.ToCostFloat64(), ct...)

	if scoreRatio >= 1.5*costRatio { // beams are >50% more cost efficient than torps
		ai.updateWarshipAmounts(warshipCount.bombers, warshipCount.warships, 0, warshipCount.fuelTransports)
	} else if scoreRatio*1.5 <= costRatio { // torp ships are >50% more cost efficient than beams
		ai.updateWarshipAmounts(warshipCount.bombers, 0, warshipCount.warships, warshipCount.fuelTransports)
	} else {
		// mix fleets based on relative strength factor
		beamShips := int(scoreRatio / costRatio * float64(warshipCount.warships) / 2)
		ai.updateWarshipAmounts(warshipCount.bombers, beamShips, warshipCount.warships-beamShips, warshipCount.fuelTransports)
	}
	return nil
}

// Determine the amount of warships used by an AI player based on years passed
func (ai *aiPlayer) getWarshipCount() (warshipQty warshipCount) {
	yearsAfterStart := ai.game.YearsPassed() - ai.config.startAttackingYear

	// determine ship counts by year
	// TODO: Make these values configurable per AI type
	switch {
	case yearsAfterStart < 5: // <2429 non-BBS; <2424 accBBS
		warshipQty.bombers = 5
		warshipQty.warships = 14
	case yearsAfterStart < 10: // 2430-2434 non-BBS; 2425-2429 accBBS
		warshipQty.bombers = 7
		warshipQty.warships = 16
	case yearsAfterStart < 15: // 2435-2439 non-BBS; 2430-2434 accBBS
		warshipQty.bombers = 9
		warshipQty.warships = 18
	case yearsAfterStart < 20: // 2440-2444 non-BBS; 2435-2439 accBBS
		warshipQty.bombers = 10
		warshipQty.warships = 20
	case yearsAfterStart < 30: // 2445-2454 non-BBS; 2440-2449 accBBS
		warshipQty.bombers = 15
		warshipQty.warships = 30
	case yearsAfterStart < 40: // 2455-2464 non-BBS; 2450-2459 accBBS
		warshipQty.bombers = 20
		warshipQty.warships = 50
	case yearsAfterStart < 50: // 2465-2474 non-BBS; 2460-2469 accBBS
		warshipQty.bombers = 30
		warshipQty.warships = 60
	default: // 2475+ non-BBS; 2470+ acc-BBS
		warshipQty.bombers = 40
		warshipQty.warships = min((yearsAfterStart/5)*6, 150)
	}

	// only add on fuel transports if we have them and they can repair our fleets
	if ai.designsByPurpose[cs.ShipDesignPurposeFuelFreighter] != nil &&
		ai.designsByPurpose[cs.ShipDesignPurposeFuelFreighter].Spec.RepairBonus > 0 {
		warshipQty.fuelTransports = min((warshipQty.bombers+warshipQty.warships)/5, 25)
	}

	return warshipQty
}

// update an AI player's fleetsByPurpose map with the provided warship quantities;
// effectively a large wrapper function
func (ai *aiPlayer) updateWarshipAmounts(bombers, beamShips, torpedoShips, fuelTransports int) {
	// reset the fleets
	ai.fleetsByPurpose[cs.FleetPurposeBomber] = fleet{
		purpose: cs.FleetPurposeBomber,
		ships:   []fleetShip{{purpose: cs.ShipDesignPurposeBomber, quantity: bombers}},
	}
	ai.fleetsByPurpose[cs.FleetPurposeCapitalShip] = fleet{
		purpose: cs.FleetPurposeCapitalShip,
		ships:   []fleetShip{},
	}

	// only add on ships if we want to add any
	if beamShips > 0 {
		// re-assign entire structs due to golang jank
		ai.fleetsByPurpose[cs.FleetPurposeBomber] = fleet{
			purpose: cs.FleetPurposeBomber,
			ships: append(ai.fleetsByPurpose[cs.FleetPurposeBomber].ships, fleetShip{
				purpose:  cs.ShipDesignPurposeBeamFighter,
				quantity: beamShips,
			}),
		}
		ai.fleetsByPurpose[cs.FleetPurposeCapitalShip] = fleet{
			purpose: cs.FleetPurposeCapitalShip,
			ships: append(ai.fleetsByPurpose[cs.FleetPurposeCapitalShip].ships, fleetShip{
				purpose:  cs.ShipDesignPurposeBeamFighter,
				quantity: beamShips,
			}),
		}
	}
	if torpedoShips > 0 {
		ai.fleetsByPurpose[cs.FleetPurposeBomber] = fleet{
			purpose: cs.FleetPurposeBomber,
			ships: append(ai.fleetsByPurpose[cs.FleetPurposeBomber].ships, fleetShip{
				purpose:  cs.ShipDesignPurposeTorpedoFighter,
				quantity: torpedoShips,
			}),
		}
		ai.fleetsByPurpose[cs.FleetPurposeCapitalShip] = fleet{
			purpose: cs.FleetPurposeCapitalShip,
			ships: append(ai.fleetsByPurpose[cs.FleetPurposeCapitalShip].ships, fleetShip{
				purpose:  cs.ShipDesignPurposeTorpedoFighter,
				quantity: torpedoShips,
			}),
		}
	}
	if fuelTransports > 0 {
		ai.fleetsByPurpose[cs.FleetPurposeBomber] = fleet{
			purpose: cs.FleetPurposeBomber,
			ships: append(ai.fleetsByPurpose[cs.FleetPurposeBomber].ships, fleetShip{
				purpose:  cs.ShipDesignPurposeFuelFreighter,
				quantity: fuelTransports,
			}),
		}
	}
}

// process an AI player's turn
func (ai *aiPlayer) ProcessTurn() error {
	ai.assignPurpose()
	ai.gatherIntel()
	ai.plan()
	ai.designStarbases()

	// TODO: Add packet defense checks
	// if we see a packet coming to us, try to load its cargo,
	// queue up defenses/drivers or evacuate pop in that order

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
	// only attack or invade after a while
	if ai.game.YearsPassed() >= ai.config.startAttackingYear {
		if err := ai.invade(); err != nil {
			return err
		}
		if err := ai.bomb(); err != nil {
			return err
		}
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

	if ai.game.Year%4 == 0 || len(ai.Player.Spec.TechsJustGained) > 0 {
		// only update warship amounts/designs every 4 years or if we just gained a tech level
		if err := ai.updateWarfleets(); err != nil {
			return err
		}
	}

	// make sure our research is optimal
	ai.research()
	// cleanup any old designs we haven't built
	ai.removeUnusedDesigns()

	return nil
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
