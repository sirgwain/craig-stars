package game

import (
	"fmt"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
)

type generatedUniverse struct {
	universe Universe
	players  []*Player
	size     Size
	density  Density
	rules    *Rules
}

type UniverseGenerator interface {
	Generate() (*Universe, error)
}

func NewUniverseGenerator(size Size, density Density, players []*Player, rules *Rules) UniverseGenerator {
	return &generatedUniverse{
		size:    size,
		density: density,
		players: players,
		rules:   rules,
	}
}

// Generate a new universe using a UniverseGenerator
func (gu *generatedUniverse) Generate() (*Universe, error) {
	log.Debug().Msgf("%s: Generating universe", gu.size)

	gu.universe = Universe{}
	area, err := gu.rules.GetArea(gu.size)
	if err != nil {
		return nil, err
	}

	if err := gu.generatePlanets(area); err != nil {
		return nil, err
	}

	// save our area
	gu.universe.Area = area

	gu.generateWormholes()

	gu.generatePlayerTechLevels()
	gu.generatePlayerPlans()
	gu.generatePlayerShipDesigns()

	if err := gu.generatePlayerHomeworlds(area); err != nil {
		return nil, err
	}

	if err := gu.generatePlayerPlanetReports(); err != nil {
		return nil, err
	}

	// generatePlayerFleets(g)
	gu.applyGameStartModeModifier()

	// setup all the specs for planets, fleets, etc
	for _, player := range gu.players {
		player.Spec = computePlayerSpec(player, gu.rules)
	}

	for _, planet := range gu.universe.Planets {
		if planet.Owned() {
			player := gu.players[planet.PlayerNum]
			planet.Spec = ComputePlanetSpec(gu.rules, planet, player)
		}
	}
	// disoverer := discover{g}
	for _, player := range gu.players {

		scanner := newPlayerScanner(&gu.universe, gu.rules, player)
		if err := scanner.scan(); err != nil {
			return nil, err
		}

		// todo: Do player discoverer to help player's disover about other players
		// disoverer.playerInfoDiscover(player)

		// TODO: check for AI player
		ai := NewAIPlayer(player, gu.universe.GetPlayerMapObjects(player.Num))
		ai.processTurn()
	}

	return &gu.universe, nil
}

// Check if a position vector is a mininum distance away from all other points
func (gu *generatedUniverse) isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool {
	minDistanceSquared := minDistance * minDistance

	for _, to := range *occupiedLocations {
		if pos.DistanceSquaredTo(to) <= minDistanceSquared {
			return false
		}
	}
	return true
}

func (gu *generatedUniverse) generatePlanets(area Vector) error {

	numPlanets, err := gu.rules.GetNumPlanets(gu.size, gu.density)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Generating %d planets in universe size %v for ", numPlanets, area)

	names := planetNames
	gu.rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	gu.universe.Planets = make([]*Planet, numPlanets)

	planetsByPosition := make(map[Vector]*Planet, numPlanets)
	occupiedLocations := make([]Vector, numPlanets)

	for i := 0; i < numPlanets; i++ {

		// find a valid position for the planet
		posCheckCount := 0
		pos := Vector{X: float64(gu.rules.random.Intn(int(area.X))), Y: float64(gu.rules.random.Intn(int(area.Y)))}
		for !gu.isPositionValid(pos, &occupiedLocations, float64(gu.rules.PlanetMinDistance)) {
			pos = Vector{X: float64(gu.rules.random.Intn(int(area.X))), Y: float64(gu.rules.random.Intn(int(area.Y)))}
			posCheckCount++
			if posCheckCount > 1000 {
				return fmt.Errorf("failed to find a valid position for a planet in 1000 tries, min: %d, numPlanets: %d, area: %v", gu.rules.PlanetMinDistance, numPlanets, area)
			}
		}

		// setup a new planet
		planet := NewPlanet()
		planet.Name = names[i]
		planet.Num = int(i + 1)
		planet.Position = pos
		planet.randomize(gu.rules)

		gu.universe.Planets[i] = planet
		planetsByPosition[pos] = planet
		occupiedLocations = append(occupiedLocations, pos)
	}

	// shuffle these so id 1 is not always the first planet in the list
	// later on we will add homeworlds based on first planet, second planet, etc
	// gu.rules.Random.Shuffle(len(gu.planets), func(i, j int) { gu.planets[i], gu.planets[j] = gu.planets[j], gu.planets[i] })

	return nil
}

func (gu *generatedUniverse) generateWormholes() {

}

func (gu *generatedUniverse) generatePlayerTechLevels() {
	for _, player := range gu.players {
		player.TechLevels = TechLevel(player.Race.Spec.StartingTechLevels)
	}
}

func (gu *generatedUniverse) generatePlayerPlans() {

}

// generate designs for each player
func (gu *generatedUniverse) generatePlayerShipDesigns() {
	for _, player := range gu.players {
		designNames := mapset.NewSet[string]()
		for _, startingPlanet := range player.Race.Spec.StartingPlanets {
			for _, startingFleet := range startingPlanet.StartingFleets {
				if designNames.Contains(startingFleet.Name) {
					// only one design per name, i.e. Scout, Armored Probe
					continue
				}
				techStore := gu.rules.techs
				hull := techStore.GetHull(string(startingFleet.HullName))
				design := designShip(techStore, hull, startingFleet.Name, player, player.DefaultHullSet, startingFleet.Purpose)
				design.HullSetNumber = int(startingFleet.HullSetNumber)
				design.Purpose = startingFleet.Purpose
				design.Spec = ComputeShipDesignSpec(gu.rules, player, design)
				player.Designs = append(player.Designs, *design)
			}
		}

		starbaseDesigns := gu.getStartingStarbaseDesigns(gu.rules.techs, player)

		for i := range starbaseDesigns {
			design := &starbaseDesigns[i]
			design.Purpose = ShipDesignPurposeStarbase
			design.Spec = ComputeShipDesignSpec(gu.rules, player, design)
			player.Designs = append(player.Designs, *design)
		}
	}

}

// have each player discover all the planets in the universe
func (gu *generatedUniverse) generatePlayerPlanetReports() error {
	for _, player := range gu.players {
		discoverer := newDiscoverer(player)
		player.PlanetIntels = make([]PlanetIntel, len(gu.universe.Planets))
		for j := range gu.universe.Planets {
			if err := discoverer.discoverPlanet(gu.rules, player, gu.universe.Planets[j], false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (gu *generatedUniverse) generatePlayerHomeworlds(area Vector) error {

	ownedPlanets := []*Planet{}
	rules := gu.rules
	random := gu.rules.random

	// each player homeworld has the same random mineral concentration, for fairness
	homeworldMinConc := Mineral{
		Ironium:   rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
		Boranium:  rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
		Germanium: rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
	}

	homeworldSurfaceMinerals := Mineral{
		Ironium:   rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
		Boranium:  rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
		Germanium: rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
	}

	extraWorldSurfaceMinerals := Mineral{
		Ironium:   rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
		Boranium:  rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
		Germanium: rules.MinStartingMineralSurface + random.Intn(rules.MaxStartingMineralSurface),
	}

	for _, player := range gu.players {
		minPlayerDistance := float64(area.X+area.Y) / 2.0 / float64(len(gu.players)+1)
		fleetNum := 1
		var homeworld *Planet

		for startingPlanetIndex, startingPlanet := range player.Race.Spec.StartingPlanets {
			// find a playerPlanet that is a min distance from other homeworlds
			var playerPlanet *Planet
			if startingPlanetIndex > 0 {

				// extra planets are close to the homeworld
				for _, planet := range gu.universe.Planets {
					distToHomeworld := planet.Position.DistanceSquaredTo(homeworld.Position)
					if !planet.Owned() && (distToHomeworld <= float64(rules.MaxExtraWorldDistance*rules.MaxExtraWorldDistance) && distToHomeworld >= float64(rules.MinExtraWorldDistance*rules.MinExtraWorldDistance)) {
						playerPlanet = planet
						break
					}
				}

			} else {
				// homeworld should be distant from other players
				for _, planet := range gu.universe.Planets {
					if !planet.Owned() && (len(ownedPlanets) == 0 || planet.shortestDistanceToPlanets(&ownedPlanets) > minPlayerDistance) {
						playerPlanet = planet
						break
					}
				}

				homeworld = playerPlanet
			}

			if playerPlanet == nil {
				return fmt.Errorf("failed to find homeworld for player %v among %d planets, minDistance: %0.1f", player, len(gu.universe.Planets), minPlayerDistance)
			}

			ownedPlanets = append(ownedPlanets, playerPlanet)

			// our starting planet starts with default fleets
			surfaceMinerals := homeworldSurfaceMinerals
			if startingPlanetIndex != 0 {
				surfaceMinerals = extraWorldSurfaceMinerals
			}

			// first planet is a homeworld
			// make a new homeworld
			if err := playerPlanet.initStartingWorld(player, gu.rules, startingPlanet, homeworldMinConc, surfaceMinerals); err != nil {
				return err
			}
			// generate some fleets on the homeworld
			if err := gu.generatePlayerFleets(player, playerPlanet, &fleetNum, startingPlanet.StartingFleets); err != nil {
				return err
			}

			messager.longMessage(player)
		}
	}

	return nil
}

func (gu *generatedUniverse) generatePlayerFleets(player *Player, planet *Planet, fleetNum *int, startingFleets []StartingFleet) error {
	for _, startingFleet := range startingFleets {
		design := player.GetDesign(string(startingFleet.Name))
		if design == nil {
			return fmt.Errorf("no design named %s found for player %s", startingFleet.Name, player)
		}

		fleet := NewFleet(player, design, *fleetNum, startingFleet.Name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, design.Spec.IdealSpeed)})
		fleet.OrbitingPlanetNum = planet.Num
		fleet.Spec = ComputeFleetSpec(gu.rules, player, &fleet)
		fleet.Fuel = fleet.Spec.FuelCapacity
		gu.universe.Fleets = append(gu.universe.Fleets, &fleet)
		(*fleetNum)++ // increment the fleet num
	}

	return nil
}

func (gu *generatedUniverse) applyGameStartModeModifier() {

}

// get the initial starbase designs for a player
func (gu *generatedUniverse) getStartingStarbaseDesigns(techStore *TechStore, player *Player) []ShipDesign {
	designs := []ShipDesign{}

	if player.Race.Spec.LivesOnStarbases {
		// create a starter colony for AR races
		starterColony := NewShipDesign(player).
			WithName("Starter Colony").
			WithHull(OrbitalFort.Name).
			WithPurpose(ShipDesignPurposeStarterColony).
			WithHullSetNumber(player.DefaultHullSet)
		starterColony.CanDelete = false
		designs = append(designs, *starterColony)
	}

	startingPlanets := player.Race.Spec.StartingPlanets

	starbase := NewShipDesign(player).
		WithName(startingPlanets[0].StarbaseDesignName).
		WithHull(startingPlanets[0].StarbaseHull).
		WithPurpose(ShipDesignPurposeStarbase).
		WithHullSetNumber(player.DefaultHullSet)

	fillStarbaseSlots(techStore, starbase, &player.Race, startingPlanets[0])
	designs = append(designs, *starbase)

	// add an orbital fort for players that start with extra planets
	if len(startingPlanets) > 1 {
		for i := range startingPlanets {
			if i == 0 {
				continue
			}
			startingPlanet := startingPlanets[i]
			fort := NewShipDesign(player).
				WithName(startingPlanet.StarbaseDesignName).
				WithHull(startingPlanet.StarbaseHull).
				WithPurpose(ShipDesignPurposeFort).
				WithHullSetNumber(player.DefaultHullSet)
			// TODO: Do we want to support a PRT that includes more than 2 planets but only some of them with
			// stargates?
			fillStarbaseSlots(techStore, fort, &player.Race, startingPlanets[i])
			designs = append(designs, *fort)
		}
	}

	return designs
}

// Player starting starbases are all the same, regardless of starting tech level
// They get half filled with the starter beam, shield, and armor
func fillStarbaseSlots(techStore *TechStore, starbase *ShipDesign, race *Race, startingPlanet StartingPlanet) {
	hull := techStore.GetHull(starbase.Hull)
	beamWeapon := techStore.GetHullComponentsByCategory(TechCategoryBeamWeapon)[0]
	shield := techStore.GetHullComponentsByCategory(TechCategoryShield)[0]
	var massDriver TechHullComponent
	var stargate TechHullComponent
	for _, hc := range techStore.GetHullComponentsByCategory(TechCategoryOrbital) {
		if hc.PacketSpeed > 0 {
			massDriver = hc
			break
		}
	}

	for _, hc := range techStore.GetHullComponentsByCategory(TechCategoryOrbital) {
		if hc.SafeRange > 0 {
			stargate = hc
			break
		}
	}

	placedMassDriver := false
	placedStargate := false
	for index, slot := range hull.Slots {
		switch slot.Type {
		case HullSlotTypeWeapon:
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{beamWeapon.Name, index + 1, int(math.Round(float64(slot.Capacity) / 2))})
		case HullSlotTypeShield:
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{shield.Name, index + 1, int(math.Round(float64(slot.Capacity) / 2))})
		case HullSlotTypeOrbital:
		case HullSlotTypeOrbitalElectrical:
			if startingPlanet.HasStargate && !placedStargate {
				starbase.Slots = append(starbase.Slots, ShipDesignSlot{stargate.Name, index + 1, 1})
				placedStargate = true
			} else if startingPlanet.HasMassDriver && !placedMassDriver {
				starbase.Slots = append(starbase.Slots, ShipDesignSlot{massDriver.Name, index + 1, 1})
				placedMassDriver = true
			}
		}
	}
}
