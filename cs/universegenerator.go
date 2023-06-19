package cs

import (
	"fmt"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
)

type universeGenerator struct {
	universe Universe
	players  []*Player
	size     Size
	density  Density
	rules    *Rules
	area     Vector
}

type UniverseGenerator interface {
	Generate() (*Universe, error)
	Area() Vector
}

func NewUniverseGenerator(size Size, density Density, players []*Player, rules *Rules) UniverseGenerator {
	return &universeGenerator{
		size:    size,
		density: density,
		players: players,
		rules:   rules,
	}
}

func (ug *universeGenerator) Area() Vector {
	return ug.area
}

// Generate a new universe using a UniverseGenerator
func (ug *universeGenerator) Generate() (*Universe, error) {
	log.Debug().Msgf("%s: Generating universe", ug.size)

	ug.universe = NewUniverse(ug.rules)
	area, err := ug.rules.GetArea(ug.size)
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
		player.Spec = computePlayerSpec(player, ug.rules, ug.universe.Planets)
	}

	for _, planet := range ug.universe.Planets {
		if planet.owned() {
			player := ug.players[planet.PlayerNum-1]
			planet.Spec = computePlanetSpec(ug.rules, player, planet)
		}
	}

	// TODO: chicken and egg problem. Player spec needs planet spec for resources, planet spec needs player spec for defense/scanner
	for _, player := range ug.players {
		player.Spec = computePlayerSpec(player, ug.rules, ug.universe.Planets)
	}

	// do one scan run
	ug.generatePlayerIntel()

	return &ug.universe, nil
}

func (ug *universeGenerator) generatePlanets() error {

	numPlanets, err := ug.rules.GetNumPlanets(ug.size, ug.density)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Generating %d planets in universe size %0.0fx%0.0f for ", numPlanets, ug.area.X, ug.area.Y)

	names := planetNames
	ug.rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	ug.universe.Planets = make([]*Planet, numPlanets)

	planetsByPosition := make(map[Vector]*Planet, numPlanets)
	occupiedLocations := make([]Vector, numPlanets)
	width, height := int(ug.area.X), int(ug.area.Y)

	for i := 0; i < numPlanets; i++ {

		// find a valid position for the planet
		posCheckCount := 0
		pos := Vector{X: float64(ug.rules.random.Intn(width)), Y: float64(ug.rules.random.Intn(height))}
		for !ug.universe.isPositionValid(pos, &occupiedLocations, float64(ug.rules.PlanetMinDistance)) {
			pos = Vector{X: float64(ug.rules.random.Intn(width)), Y: float64(ug.rules.random.Intn(height))}
			posCheckCount++
			if posCheckCount > 1000 {
				return fmt.Errorf("find a valid position for a wormhole in 1000 tries, min: %d, numPlanets: %d, area: %v", ug.rules.PlanetMinDistance, numPlanets, ug.area)
			}
		}

		// setup a new planet
		planet := NewPlanet()
		planet.Name = names[i]
		planet.Num = i + 1
		planet.Position = pos
		planet.randomize(ug.rules)

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
	numPairs := ug.rules.WormholePairsForSize[ug.size]
	wormholes := make([]*Wormhole, numPairs*2)

	planetPositions := make([]Vector, len(ug.universe.Planets))
	wormholePositions := make([]Vector, len(wormholes))
	for i, planet := range ug.universe.Planets {
		planetPositions[i] = planet.Position
	}

	for i := 0; i < numPairs*2; i++ {
		position, stability, err := generateWormhole(&ug.universe, i+1, ug.area, ug.rules.random, planetPositions, wormholePositions, ug.rules.WormholeMinPlanetDistance)

		if err != nil {
			return err
		}

		var companion *Wormhole
		if i%2 > 0 {
			companion = wormholes[i-1]
		}
		wormhole := ug.universe.createWormhole(position, stability, companion)
		log.Debug().Msgf("generated Wormhole at (%0.0f, %0.0f)", wormhole.Position.X, wormhole.Position.Y)

		wormholePositions[i] = wormhole.Position
		wormholes[i] = wormhole
	}

	ug.universe.Wormholes = wormholes

	return nil
}

func (ug *universeGenerator) generateAIPlayers() {
	names := AINames
	ug.rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })
	for index, player := range ug.players {
		if player.AIControlled {
			name := names[index%len(names)]
			player.Race.Name = name
			player.Race.PluralName = name + "s"
		}
	}
}

func (ug *universeGenerator) generatePlayerTechLevels() {
	for _, player := range ug.players {
		player.TechLevels = TechLevel(player.Race.Spec.StartingTechLevels)
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
					// only one design per name, i.e. Scout, Armored Probe
					continue
				}
				techStore := ug.rules.techs
				hull := techStore.GetHull(string(startingFleet.HullName))
				design := DesignShip(techStore, hull, startingFleet.Name, player, num, player.DefaultHullSet, startingFleet.Purpose)
				design.HullSetNumber = int(startingFleet.HullSetNumber)
				design.Purpose = startingFleet.Purpose
				design.Spec = ComputeShipDesignSpec(ug.rules, player.TechLevels, player.Race.Spec, design)
				player.Designs = append(player.Designs, design)
				designNames.Add(design.Name)
				num++
			}
		}

		starbaseDesigns := ug.getStartingStarbaseDesigns(ug.rules.techs, player, num)

		for i := range starbaseDesigns {
			design := &starbaseDesigns[i]
			design.Purpose = ShipDesignPurposeStarbase
			design.Spec = ComputeShipDesignSpec(ug.rules, player.TechLevels, player.Race.Spec, design)
			player.Designs = append(player.Designs, design)
		}
	}

}

// have each player discover all the planets in the universe
func (ug *universeGenerator) generatePlayerPlanetReports() error {
	for _, player := range ug.players {
		discoverer := newDiscoverer(player)
		player.PlanetIntels = make([]PlanetIntel, len(ug.universe.Planets))
		for j := range ug.universe.Planets {
			// start with some defaults
			intel := &player.PlanetIntels[j]
			intel.ReportAge = ReportAgeUnexplored
			intel.Type = MapObjectTypePlanet
			intel.PlayerNum = Unowned

			if err := discoverer.discoverPlanet(ug.rules, player, ug.universe.Planets[j], false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ug *universeGenerator) generatePlayerHomeworlds(area Vector) error {

	ownedPlanets := []*Planet{}
	rules := ug.rules
	random := ug.rules.random

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

	for _, player := range ug.players {
		minPlayerDistance := float64(area.X+area.Y) / 2.0 / float64(len(ug.players)+1)
		fleetNum := 1
		var homeworld *Planet

		for startingPlanetIndex, startingPlanet := range player.Race.Spec.StartingPlanets {
			// find a playerPlanet that is a min distance from other homeworlds
			var playerPlanet *Planet
			if startingPlanetIndex > 0 {

				// extra planets are close to the homeworld
				for _, planet := range ug.universe.Planets {
					distToHomeworld := planet.Position.DistanceSquaredTo(homeworld.Position)
					if !planet.owned() && (distToHomeworld <= float64(rules.MaxExtraWorldDistance*rules.MaxExtraWorldDistance) && distToHomeworld >= float64(rules.MinExtraWorldDistance*rules.MinExtraWorldDistance)) {
						playerPlanet = planet
						break
					}
				}

			} else {
				// homeworld should be distant from other players
				for _, planet := range ug.universe.Planets {
					if !planet.owned() && (len(ownedPlanets) == 0 || planet.shortestDistanceToPlanets(&ownedPlanets) > minPlayerDistance) {
						playerPlanet = planet
						break
					}
				}

				homeworld = playerPlanet
			}

			if playerPlanet == nil {
				return fmt.Errorf("find homeworld for player %v among %d planets, minDistance: %0.1f", player, len(ug.universe.Planets), minPlayerDistance)
			}

			ownedPlanets = append(ownedPlanets, playerPlanet)

			// our starting planet starts with default fleets
			surfaceMinerals := homeworldSurfaceMinerals
			if startingPlanetIndex != 0 {
				surfaceMinerals = extraWorldSurfaceMinerals
			}

			// first planet is a homeworld
			// make a new homeworld
			if err := playerPlanet.initStartingWorld(player, ug.rules, startingPlanet, homeworldMinConc, surfaceMinerals); err != nil {
				return err
			}
			if playerPlanet.starbase != nil {
				ug.universe.Starbases = append(ug.universe.Starbases, playerPlanet.starbase)
			}
			// generate some fleets on the homeworld
			if err := ug.generatePlayerFleets(player, playerPlanet, &fleetNum, startingPlanet.StartingFleets); err != nil {
				return err
			}

			// // create a test minefield
			// testMineField := newMineField(player, MineFieldTypeStandard, 1200, ug.universe.getNextMineFieldNum(), playerPlanet.Position)
			// testMineField.Spec = computeMinefieldSpec(testMineField)
			// ug.universe.MineFields = append(ug.universe.MineFields, testMineField)

			messager.longMessage(player)
		}
	}

	return nil
}

func (ug *universeGenerator) generatePlayerFleets(player *Player, planet *Planet, fleetNum *int, startingFleets []StartingFleet) error {
	for _, startingFleet := range startingFleets {
		design := player.GetDesign(string(startingFleet.Name))
		if design == nil {
			return fmt.Errorf("no design named %s found for player %s", startingFleet.Name, player)
		}
		fleet := newFleet(player, design, *fleetNum, startingFleet.Name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, design.Spec.Engine.IdealSpeed)})
		fleet.OrbitingPlanetNum = planet.Num
		fleet.Spec = ComputeFleetSpec(ug.rules, player, &fleet)
		fleet.Fuel = fleet.Spec.FuelCapacity
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
		starterColony.CanDelete = false
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
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{beamWeapon.Name, index + 1, int(math.Round(float64(slot.Capacity) / 2)), &beamWeapon})
		case HullSlotTypeShield:
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{shield.Name, index + 1, int(math.Round(float64(slot.Capacity) / 2)), &shield})
		case HullSlotTypeOrbital:
		case HullSlotTypeOrbitalElectrical:
			if startingPlanet.HasStargate && !placedStargate {
				starbase.Slots = append(starbase.Slots, ShipDesignSlot{stargate.Name, index + 1, 1, &stargate})
				placedStargate = true
			} else if startingPlanet.HasMassDriver && !placedMassDriver {
				starbase.Slots = append(starbase.Slots, ShipDesignSlot{massDriver.Name, index + 1, 1, &massDriver})
				placedMassDriver = true
			}
		}
	}
}

func (ug *universeGenerator) generatePlayerRelations() error {
	for _, player := range ug.players {

		player.Relations = make([]PlayerRelationship, len(ug.players))

		for i, otherPlayer := range ug.players {
			relationship := &player.Relations[i]
			if otherPlayer.Num == player.Num {
				// we're friends with ourselves
				relationship.Relation = PlayerRelationFriend
			} else {
				relationship.Relation = PlayerRelationEnemy
			}

		}

	}

	return nil
}

func (ug *universeGenerator) generatePlayerIntel() error {
	for _, player := range ug.players {

		// discover other players
		player.PlayerIntels.PlayerIntels = make([]PlayerIntel, len(ug.players))
		for i, otherPlayer := range ug.players {
			relationship := &player.Relations[i]
			if otherPlayer.Num == player.Num {
				// we're friends with ourselves
				relationship.Relation = PlayerRelationFriend
			} else {
				relationship.Relation = PlayerRelationEnemy
			}

			playerIntel := &player.PlayerIntels.PlayerIntels[i]
			playerIntel.Color = otherPlayer.Color
			playerIntel.Name = otherPlayer.Name
			playerIntel.Num = otherPlayer.Num

			// we know about ourselves
			if otherPlayer.Num == player.Num {
				playerIntel.Seen = true
				playerIntel.RaceName = player.Race.Name
				playerIntel.RacePluralName = player.Race.PluralName
			}
		}

		// do initial scans
		scanner := newPlayerScanner(&ug.universe, ug.players, ug.rules, player)
		if err := scanner.scan(); err != nil {
			return err
		}

	}

	return nil
}
