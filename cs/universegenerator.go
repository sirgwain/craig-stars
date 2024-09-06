package cs

import (
	"fmt"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
)

// The UniverseGenerator generates a new universe based on some game settings and players.
type UniverseGenerator interface {
	Generate() (*Universe, error)
	Area() Vector
}

type universeGenerator struct {
	*Game
	universe Universe
	players  []*Player
	area     Vector
}

func NewUniverseGenerator(game *Game, players []*Player) UniverseGenerator {
	return &universeGenerator{
		Game:    game,
		players: players,
	}
}

func (ug *universeGenerator) Area() Vector {
	return ug.area
}

// Generate a new universe using a UniverseGenerator
func (ug *universeGenerator) Generate() (*Universe, error) {
	log.Debug().Msgf("%s: Generating universe", ug.Size)

	for _, player := range ug.players {
		player.Race.Spec = computeRaceSpec(&player.Race, &ug.Rules)
		player.discoverer = newDiscovererWithAllies(player, ug.players)
	}

	ug.universe = NewUniverse(&ug.Rules)
	area, err := ug.Rules.GetArea(ug.Size)
	if err != nil {
		return nil, err
	}

	ug.area = area

	if err := ug.generatePlanets(); err != nil {
		return nil, err
	}

	if err := ug.generateWormholes(); err != nil {
		return nil, err
	}

	ug.generateAIPlayers()
	ug.generatePlayerTechLevels()
	ug.generatePlayerPlans()
	ug.generatePlayerShipDesigns()
	ug.generatePlayerRelations()

	if err := ug.generatePlayerHomeworlds(ug.area); err != nil {
		return nil, err
	}

	if err := ug.generatePlayerPlanetReports(); err != nil {
		return nil, err
	}

	ug.applyGameStartModeModifier()

	// setup all the specs for planets, fleets, etc
	for _, player := range ug.players {
		player.Spec = computePlayerSpec(player, &ug.Rules, ug.universe.Planets)
	}

	for _, planet := range ug.universe.Planets {
		if planet.Owned() {
			player := ug.players[planet.PlayerNum-1]
			planet.Spec = computePlanetSpec(&ug.Rules, player, planet)
			if err := planet.PopulateProductionQueueDesigns(player); err != nil {
				return nil, fmt.Errorf("%s failed to populate queue designs: %w", planet, err)
			}
			planet.PopulateProductionQueueEstimates(&ug.Rules, player)
		}
	}

	// TODO: chicken and egg problem. Player spec needs planet spec for resources, planet spec needs player spec for defense/scanner
	for _, player := range ug.players {
		player.Spec = computePlayerSpec(player, &ug.Rules, ug.universe.Planets)
	}

	// do one scan run
	ug.generatePlayerIntel()

	return &ug.universe, nil
}

func (ug *universeGenerator) generatePlanets() error {

	numPlanets, err := ug.Rules.GetNumPlanets(ug.Size, ug.Density)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Generating %d planets in universe size %0.0fx%0.0f for ", numPlanets, ug.area.X, ug.area.Y)

	names := planetNames
	rules := &ug.Rules
	rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	ug.universe.Planets = make([]*Planet, numPlanets)

	planetsByPosition := make(map[Vector]*Planet, numPlanets)
	occupiedLocations := make([]Vector, numPlanets)
	width, height := int(ug.area.X), int(ug.area.Y)

	for i := 0; i < numPlanets; i++ {

		// find a valid position for the planet
		posCheckCount := 0
		pos := Vector{X: float64(rules.random.Intn(width)), Y: float64(rules.random.Intn(height))}
		for !ug.universe.isPositionValid(pos, &occupiedLocations, float64(rules.PlanetMinDistance)) {
			pos = Vector{X: float64(rules.random.Intn(width)), Y: float64(rules.random.Intn(height))}
			posCheckCount++
			if posCheckCount > 1000 {
				return fmt.Errorf("find a valid position for a wormhole in 1000 tries, min: %d, numPlanets: %d, area: %v", rules.PlanetMinDistance, numPlanets, ug.area)
			}
		}

		// setup a new planet
		planet := NewPlanet()
		planet.Name = names[i]
		planet.Num = i + 1
		planet.Position = pos
		planet.randomize(rules)

		if ug.MaxMinerals {
			planet.MineralConcentration = Mineral{100, 100, 100}
		}
		if ug.RandomEvents && rules.RandomEventChances[RandomEventAncientArtifact] >= rules.random.Float64() {
			// check if this planet has a random artifact
			planet.RandomArtifact = true
		}

		ug.universe.Planets[i] = planet
		planetsByPosition[pos] = planet
		occupiedLocations = append(occupiedLocations, pos)
	}

	// TODO: to make it easier to develop and troubleshoot data, currently leaving this unshuffled
	// shuffle these so id 1 is not always the first planet in the list
	// later on we will add homeworlds based on first planet, second planet, etc
	// gu.rules.Random.Shuffle(len(gu.planets), func(i, j int) { gu.planets[i], gu.planets[j] = gu.planets[j], gu.planets[i] })

	return nil
}

func (ug *universeGenerator) generateWormholes() error {
	numPairs := ug.Rules.WormholePairsForSize[ug.Size]
	wormholes := make([]*Wormhole, numPairs*2)

	planetPositions := make([]Vector, len(ug.universe.Planets))
	wormholePositions := make([]Vector, len(wormholes))
	for i, planet := range ug.universe.Planets {
		planetPositions[i] = planet.Position
	}

	for i := 0; i < numPairs*2; i++ {
		position, stability, err := generateWormhole(&ug.universe, ug.area, ug.Rules.random, planetPositions, wormholePositions, ug.Rules.WormholeMinPlanetDistance)

		if err != nil {
			return err
		}

		var companion *Wormhole
		if i%2 > 0 {
			companion = wormholes[i-1]
		}
		wormhole := ug.universe.createWormhole(&ug.Rules, position, stability, companion)
		log.Debug().Msgf("generated Wormhole at (%0.0f, %0.0f)", wormhole.Position.X, wormhole.Position.Y)

		wormholePositions[i] = wormhole.Position
		wormholes[i] = wormhole
	}

	ug.universe.Wormholes = wormholes

	return nil
}

func (ug *universeGenerator) generateAIPlayers() {
	names := AINames
	cheaterNames := AICheaterNames
	ug.Rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })
	ug.Rules.random.Shuffle(len(cheaterNames), func(i, j int) { cheaterNames[i], cheaterNames[j] = cheaterNames[j], cheaterNames[i] })
	for index, player := range ug.players {
		if player.AIControlled {

			name := names[index%len(names)]
			if player.AIDifficulty == AIDifficultyCheater {
				name = cheaterNames[index%len(cheaterNames)]
			}
			player.Race.Name = name[0]
			player.Race.PluralName = name[1]
		}
	}
}

func (ug *universeGenerator) generatePlayerTechLevels() {
	for _, player := range ug.players {
		player.TechLevels = TechLevel(player.Race.Spec.StartingTechLevels)
	}
}

func (ug *universeGenerator) generatePlayerPlans() {
	for _, player := range ug.players {
		player.PlayerPlans = player.defaultPlans()
	}
}

// generate designs for each player
func (ug *universeGenerator) generatePlayerShipDesigns() {
	for _, player := range ug.players {
		designNames := mapset.NewSet[string]()
		num := 1
		for _, startingPlanet := range player.Race.Spec.StartingPlanets {
			for _, startingFleet := range startingPlanet.StartingFleets {
				if designNames.Contains(startingFleet.Name) {
					// only one design per name, i.e. Scout, Armed Probe
					continue
				}
				techStore := ug.Rules.techs
				hull := techStore.GetHull(string(startingFleet.HullName))
				design := DesignShip(techStore, hull, startingFleet.Name, player, num, player.DefaultHullSet, startingFleet.Purpose, FleetPurposeFromShipDesignPurpose(startingFleet.Purpose))
				design.HullSetNumber = int(startingFleet.HullSetNumber)
				design.Purpose = startingFleet.Purpose
				design.Spec = ComputeShipDesignSpec(&ug.Rules, player.TechLevels, player.Race.Spec, design)
				player.Designs = append(player.Designs, design)
				designNames.Add(design.Name)
				num++
			}
		}

		starbaseDesigns := ug.getStartingStarbaseDesigns(ug.Rules.techs, player, num)

		for i := range starbaseDesigns {
			design := &starbaseDesigns[i]
			design.Spec = ComputeShipDesignSpec(&ug.Rules, player.TechLevels, player.Race.Spec, design)
			player.Designs = append(player.Designs, design)
		}
	}

}

// have each player discover all the planets in the universe
func (ug *universeGenerator) generatePlayerPlanetReports() error {
	for _, player := range ug.players {
		player.initDefaultPlanetIntels(ug.universe.Planets)
	}
	return nil
}

func (ug *universeGenerator) generatePlayerHomeworlds(area Vector) error {

	ownedPlanets := []*Planet{}
	rules := &ug.Rules
	random := rules.random

	// each player homeworld has the same random mineral concentration, for fairness
	homeworldMinConc := Mineral{
		Ironium:   rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
		Boranium:  rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
		Germanium: rules.MinHomeworldMineralConcentration + random.Intn(rules.MaxStartingMineralConcentration),
	}
	if ug.MaxMinerals {
		homeworldMinConc = Mineral{100, 100, 100}
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

	for _, player := range ug.players {
		minPlayerDistance := float64(area.X+area.Y) / 2.0 / float64(len(ug.players)+1)
		fleetNum := 1
		var homeworld *Planet

		for _, startingPlanet := range player.Race.Spec.StartingPlanets {

			if !startingPlanet.Homeworld && homeworld == nil {
				return fmt.Errorf("first planet in startingPlanets not homeworld, exiting")
			}

			// find a playerPlanet that is a min distance from other homeworlds
			var playerPlanet *Planet
			farthestDistance := float64(math.MinInt)
			closestDistance := math.MaxFloat64

			if startingPlanet.Homeworld && homeworld == nil { // planet is homeworld & we have no other
				// homeworld should be distant from other players
				for _, planet := range ug.universe.Planets {
					if planet.Owned() {
						continue
					}

					// if we can't find a planet within tolerances, pick the farthest one
					shortedDistanceToPlanets := planet.shortestDistanceToPlanets(&ownedPlanets)
					if shortedDistanceToPlanets >= farthestDistance {
						farthestDistance = shortedDistanceToPlanets
						playerPlanet = planet
					}
					if len(ownedPlanets) == 0 || shortedDistanceToPlanets > minPlayerDistance {
						playerPlanet = planet
						break
					}
				}
				homeworld = playerPlanet

			} else {
				// extra planets are close to the homeworld
				for _, planet := range ug.universe.Planets {
					if planet.Owned() {
						continue
					}

					// if we can't find a planet within tolerances, pick the closest one
					distToHomeworld := planet.Position.DistanceSquaredTo(homeworld.Position)
					if distToHomeworld <= closestDistance {
						closestDistance = distToHomeworld
						playerPlanet = planet
					}
					if distToHomeworld <= float64(rules.MaxExtraWorldDistance*rules.MaxExtraWorldDistance) && distToHomeworld >= float64(rules.MinExtraWorldDistance*rules.MinExtraWorldDistance) {
						playerPlanet = planet
						break
					}
				}
			}

			if playerPlanet == nil {
				return fmt.Errorf("find homeworld for player %v among %d planets, minDistance: %0.1f", player, len(ug.universe.Planets), minPlayerDistance)
			}

			ownedPlanets = append(ownedPlanets, playerPlanet)

			var surface Mineral
			if startingPlanet.Homeworld {
				surface = homeworldSurfaceMinerals
			} else {
				surface = extraWorldSurfaceMinerals
			}

			// make a new starter world
			playerPlanet.initStartingWorld(player, &ug.Rules, startingPlanet, homeworldMinConc, surface)
			if !startingPlanet.Homeworld && !ug.MaxMinerals {
				// starting planets start with different mincons
				playerPlanet.MineralConcentration = randomizeMinerals(rules, playerPlanet.Hab.Rad)
			}

			// add a starbase to this planet
			if startingPlanet.StarbaseDesignName != "" {
				if err := ug.buildStarbase(player, playerPlanet, startingPlanet.StarbaseDesignName); err != nil {
					return err
				}
			}

			// tell the player about their homeworld
			if startingPlanet.Homeworld {
				messager.planetHomeworld(player, playerPlanet)
			}

			// generate some fleets on the homeworld
			if err := ug.generatePlayerFleets(player, playerPlanet, &fleetNum, startingPlanet.StartingFleets); err != nil {
				return err
			}
		}
	}

	return nil
}

// build a starbase on a planet
func (ug *universeGenerator) buildStarbase(player *Player, planet *Planet, designName string) error {
	// // the homeworld gets a starbase
	design := player.GetDesignByName(designName)
	if design == nil {
		return fmt.Errorf("no design named %s found", designName)
	}

	design.Spec.NumBuilt++
	design.Spec.NumInstances++
	starbase := newStarbase(player, planet, design, design.Name)
	starbase.Spec = ComputeFleetSpec(&ug.Rules, player, &starbase)
	planet.setStarbase(&ug.Rules, player, &starbase)

	ug.universe.Starbases = append(ug.universe.Starbases, &starbase)

	return nil
}

func (ug *universeGenerator) generatePlayerFleets(player *Player, planet *Planet, fleetNum *int, startingFleets []StartingFleet) error {
	for _, startingFleet := range startingFleets {
		design := player.GetDesignByName(startingFleet.Name)
		if design == nil {
			return fmt.Errorf("no design named %s found for player %s", startingFleet.Name, player)
		}
		fleet := newFleetForDesign(player, design, 1, *fleetNum, startingFleet.Name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, design.Spec.Engine.IdealSpeed)})
		fleet.OrbitingPlanetNum = planet.Num
		fleet.Spec = ComputeFleetSpec(&ug.Rules, player, &fleet)
		fleet.Fuel = fleet.Spec.FuelCapacity
		fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
		purpose := FleetPurposeFromShipDesignPurpose(design.Purpose)
		fleet.SetTag(TagPurpose, string(purpose))
		ug.universe.Fleets = append(ug.universe.Fleets, &fleet)
		design.Spec.NumInstances++
		design.Spec.NumBuilt++
		(*fleetNum)++ // increment the fleet num
	}

	return nil
}

func (ug *universeGenerator) applyGameStartModeModifier() {

}

// get the initial starbase designs for a player
func (ug *universeGenerator) getStartingStarbaseDesigns(techStore *TechStore, player *Player, designNumStart int) []ShipDesign {
	designs := []ShipDesign{}

	if player.Race.Spec.LivesOnStarbases {
		// create a starter colony for AR races
		starterColony := NewShipDesign(player, designNumStart).
			WithName("Starter Colony").
			WithHull(OrbitalFort.Name).
			WithPurpose(ShipDesignPurposeStarterColony).
			WithHullSetNumber(player.DefaultHullSet)
		starterColony.CannotDelete = true
		designNumStart++
		designs = append(designs, *starterColony)
	}

	startingPlanets := player.Race.Spec.StartingPlanets

	starbase := NewShipDesign(player, designNumStart).
		WithName(startingPlanets[0].StarbaseDesignName).
		WithHull(startingPlanets[0].StarbaseHull).
		WithPurpose(ShipDesignPurposeStarbase).
		WithHullSetNumber(player.DefaultHullSet)

	fillStarbaseSlots(techStore, starbase, &player.Race, startingPlanets[0])
	designNumStart++
	designs = append(designs, *starbase)

	// add an orbital fort for players that start with extra planets
	if len(startingPlanets) > 1 {
		for i := range startingPlanets {
			if i == 0 {
				continue
			}
			startingPlanet := startingPlanets[i]
			fort := NewShipDesign(player, designNumStart).
				WithName(startingPlanet.StarbaseDesignName).
				WithHull(startingPlanet.StarbaseHull).
				WithPurpose(ShipDesignPurposeFort).
				WithHullSetNumber(player.DefaultHullSet)
			// TODO: Do we want to support a PRT that includes more than 2 planets but only some of them with
			// stargates?
			fillStarbaseSlots(techStore, fort, &player.Race, startingPlanets[i])
			designNumStart++
			designs = append(designs, *fort)
		}
	}

	return designs
}

// Player starting starbases are all the same, regardless of starting tech level
// They get half filled with the starter beam & shield
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
		case HullSlotTypeGeneral: // No starting starbases (or any starbase) currently have GP slots, but this is a precaution if they did
			fallthrough
		case HullSlotTypeWeapon:
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{beamWeapon.Name, index + 1, int(math.Round(float64(slot.Capacity) / 2))})
		case HullSlotTypeShieldArmor:
			fallthrough
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

func (ug *universeGenerator) generatePlayerRelations() {
	for _, player := range ug.players {
		player.Relations = player.defaultRelationships(ug.players, ug.ComputerPlayersFormAlliances)
	}
}

func (ug *universeGenerator) generatePlayerIntel() error {
	for _, player := range ug.players {

		// discover other players
		player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels(ug.players)
		player.PlayerIntels.ScoreIntels = make([]ScoreIntel, len(ug.players))

		// do initial scans
		scanner := newPlayerScanner(&ug.universe, ug.players, &ug.Rules, player)
		if err := scanner.scan(); err != nil {
			return err
		}

	}

	return nil
}
