package cs

import (
	"fmt"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The UniverseGenerator generates a new universe based on some game settings and players.
type UniverseGenerator interface {
	Generate() (*Universe, error)
	Area() Vector
}

// A universe generator, used to generate starting universes for new games.
type universeGenerator struct {
	*FullGame
	area Vector
	log  zerolog.Logger
}

func NewUniverseGenerator(game *Game, players []*Player) UniverseGenerator {
	genLogger := log.With().Int64("GameID", game.ID).Str("GameName", game.Name).Logger()
	return &universeGenerator{
		FullGame: &FullGame{
			Game:    game,
			Players: players,
		},
		log: genLogger,
	}
}

func (ug *universeGenerator) Area() Vector {
	return ug.area
}

// Generate a new universe using a UniverseGenerator
func (ug *universeGenerator) Generate() (*Universe, error) {
	ug.log.Debug().Msgf("%s: Generating universe", ug.Size)

	for _, player := range ug.Players {
		player.Race.Spec = computeRaceSpec(&player.Race, &ug.Rules)
		player.discoverer = newDiscovererWithAllies(ug.log, player, ug.Players)
	}

	u := NewUniverse(ug.log, &ug.Rules)
	ug.Universe = &u
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
	if err := ug.generatePlayerShipDesigns(); err != nil {
		return nil, err
	}
	ug.generatePlayerRelations()

	if err := ug.generatePlayerHomeworlds(ug.area); err != nil {
		return nil, err
	}

	if err := ug.generatePlayerPlanetReports(); err != nil {
		return nil, err
	}

	ug.applyGameStartModeModifier()

	// setup all the specs for planets, fleets, etc
	// Normal games only need to compute player specs and whatnot, but max mode games require
	// complete re-computation due to changing techLevels, etc.
	if ug.StartMode != GameStartModeNormal {
		ug.computeSpecs()
	} else {
		for _, player := range ug.Players {
			player.Spec = computePlayerSpec(player, &ug.Rules, ug.Universe.Planets)
		}

		for _, planet := range ug.Universe.Planets {
			if planet.Owned() {
				player := ug.Players[planet.PlayerNum-1]
				planet.Spec = computePlanetSpec(&ug.Rules, player, planet)
				if err := planet.PopulateProductionQueueDesigns(player); err != nil {
					return nil, fmt.Errorf("planet %s failed to populate queue design; error: \n%w", planet, err)
				}
				if err := planet.PopulateProductionQueueEstimates(&ug.Rules, player); err != nil {
					return nil, fmt.Errorf("planet %s failed to populate queue estimates; error: \n%w", planet.Name, err)
				}
			}
		}
	}

	// TODO: chicken and egg problem. Player spec needs planet spec for resources, planet spec needs player spec for defense/scanner
	for _, player := range ug.Players {
		player.Spec = computePlayerSpec(player, &ug.Rules, ug.Universe.Planets)
	}

	// do one scan run
	if err := ug.generatePlayerIntel(); err != nil {
		return nil, err
	}

	return ug.Universe, nil
}

func (ug *universeGenerator) generatePlanets() error {

	numPlanets, err := ug.Rules.GetNumPlanets(ug.Size, ug.Density)
	if err != nil {
		return err
	}

	ug.log.Debug().Msgf("Generating %d planets in universe size %0.0fx%0.0f for ", numPlanets, ug.area.X, ug.area.Y)

	names := planetNames
	rules := &ug.Rules
	rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	ug.Universe.Planets = make([]*Planet, numPlanets)

	planetsByPosition := make(map[Vector]*Planet, numPlanets)
	occupiedLocations := make([]Vector, numPlanets)
	width, height := int(ug.area.X), int(ug.area.Y)

	for i := 0; i < numPlanets; i++ {

		// find a valid position for the planet
		posCheckCount := 0
		pos := Vector{X: float64(rules.random.Intn(width)), Y: float64(rules.random.Intn(height))}
		for !ug.Universe.isPositionValid(pos, &occupiedLocations, float64(rules.PlanetMinDistance)) {
			pos = Vector{X: float64(rules.random.Intn(width)), Y: float64(rules.random.Intn(height))}
			posCheckCount++
			if posCheckCount > 1000 {
				return fmt.Errorf("could not find a valid position for a planet in 1000 tries;\n min distance: %d, numPlanets: %d, area: %v", rules.PlanetMinDistance, numPlanets, ug.area)
			}
		}

		// we found a good position; setup a new planet
		planet := NewPlanet()
		planet.Name = names[i]
		planet.Num = i + 1
		planet.Position = pos
		planet.randomize(rules, ug.StartMode == GameStartModeAccBBS)

		if ug.MaxMinerals {
			planet.MineralConcentration = Mineral{100, 100, 100}
		}
		if ug.RandomEvents && rules.RandomEventChances[RandomEventAncientArtifact] >= rules.random.Float64() {
			// roll for a random artifact
			planet.RandomArtifact = true
		}

		ug.Universe.Planets[i] = planet
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

	planetPositions := make([]Vector, len(ug.Universe.Planets))
	wormholePositions := make([]Vector, len(wormholes))
	for i, planet := range ug.Universe.Planets {
		planetPositions[i] = planet.Position
	}

	for i := 0; i < numPairs*2; i++ {
		position, stability, err := generateWormhole(ug.Universe, ug.area, ug.Rules.random, planetPositions, wormholePositions, ug.Rules.WormholeMinPlanetDistance)

		if err != nil {
			return err
		}

		var companion *Wormhole
		if i%2 > 0 {
			companion = wormholes[i-1]
		}
		wormhole := ug.Universe.createWormhole(&ug.Rules, position, stability, companion)
		ug.log.Debug().Msgf("generated Wormhole at (%0.0f, %0.0f)", wormhole.Position.X, wormhole.Position.Y)

		wormholePositions[i] = wormhole.Position
		wormholes[i] = wormhole
	}

	ug.Universe.Wormholes = wormholes

	return nil
}

func (ug *universeGenerator) generateAIPlayers() {
	names := AINames
	cheaterNames := AICheaterNames
	ug.Rules.random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })
	ug.Rules.random.Shuffle(len(cheaterNames), func(i, j int) { cheaterNames[i], cheaterNames[j] = cheaterNames[j], cheaterNames[i] })
	for index, player := range ug.Players {
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
	for _, player := range ug.Players {
		player.TechLevels = TechLevel(player.Race.Spec.StartingTechLevels)
	}
}

func (ug *universeGenerator) generatePlayerPlans() {
	for _, player := range ug.Players {
		player.PlayerPlans = player.defaultPlans()
	}
}

// generate designs for each player
func (ug *universeGenerator) generatePlayerShipDesigns() error {
	var err error
	for _, player := range ug.Players {
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
				design, err := DesignShip(&ug.Game.Rules, hull, startingFleet.Name, player, num, player.DefaultHullSet, startingFleet.Purpose, FleetPurposeFromShipDesignPurpose(startingFleet.Purpose))
				if err != nil {
					return fmt.Errorf("DesignShip returned error %w", err)
				}
				if design == nil {
					return fmt.Errorf("failed to design ship for %s", hull)
				}
				player.Designs = append(player.Designs, design)
				designNames.Add(design.Name)
				num++
			}
		}

		starbaseDesigns := ug.createStartingStarbaseDesigns(ug.Rules.techs, player, num)

		for i := range starbaseDesigns {
			design := starbaseDesigns[i]
			design.Spec, err = ComputeShipDesignSpec(&ug.Rules, player.TechLevels, player.Race.Spec, design)
			if err != nil {
				return fmt.Errorf("ComputeShipDesignSpec returned error: %w", err)
			}
			player.Designs = append(player.Designs, design)
		}
	}
	return nil
}

// have each player discover all the planets in the universe
func (ug *universeGenerator) generatePlayerPlanetReports() error {
	for _, player := range ug.Players {
		player.initDefaultPlanetIntels(ug.Universe.Planets)
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

	for _, player := range ug.Players {
		minPlayerDistance := float64(area.X+area.Y) / 2.0 / float64(len(ug.Players)+1)
		fleetNum := 1
		var homeworld *Planet
		extraPoints, pointsType := player.Race.ComputeLeftoverRacePoints(rules.RaceStartingPoints)

		for _, startingPlanet := range player.Race.Spec.StartingPlanets {

			if !startingPlanet.Homeworld && homeworld == nil {
				// TODI: Do we want to support homeworlds in subsequent slots?
				return fmt.Errorf("first planet in player #%d's startingPlanets was not homeworld, exiting", player.Num)
			}

			// find a playerPlanet that is a min distance from other homeworlds

			var playerPlanet *Planet
			farthestDistance := float64(math.MinInt)
			closestDistance := math.MaxFloat64

			if startingPlanet.Homeworld && homeworld == nil { // planet is homeworld & we have no other
				// homeworld should be distant from other players
				for _, planet := range ug.Universe.Planets {
					if planet.Owned() {
						continue
					}

					// if we find a planet far enough away
					// from other players' planets, use it.
					// Otherwise, keep track of the furthest one away
					// and default to it at the end.
					shortedDistanceToPlanets := planet.shortestDistanceToPlanets(&ownedPlanets)
					if len(ownedPlanets) == 0 || shortedDistanceToPlanets > minPlayerDistance {
						playerPlanet = planet
						break
					}

					if shortedDistanceToPlanets >= farthestDistance {
						farthestDistance = shortedDistanceToPlanets
						playerPlanet = planet
					}
				}

				homeworld = playerPlanet
			} else {
				// extra planets are close to the homeworld
				for _, planet := range ug.Universe.Planets {
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
				return fmt.Errorf("could not find homeworld for player %v among %d planets, minDistance: %0.1f", player, len(ug.Universe.Planets), minPlayerDistance)
			}

			ownedPlanets = append(ownedPlanets, playerPlanet)

			var surface Mineral
			if startingPlanet.Homeworld {
				surface = homeworldSurfaceMinerals
			} else {
				surface = extraWorldSurfaceMinerals
			}

			// make a new starter world and spend leftover points
			ug.log.Debug().Msgf("Assigning %s to %s as homeworld", playerPlanet, player)
			playerPlanet.initStartingWorld(player, &ug.Rules, startingPlanet, homeworldMinConc, surface)
			if startingPlanet.Homeworld {
				ug.assignRaceStartingPointBonuses(&player.Race, playerPlanet, extraPoints, pointsType)
			} else if !ug.MaxMinerals {
				playerPlanet.MineralConcentration = randomizeMinerals(rules, playerPlanet.Hab.Rad, ug.StartMode == GameStartModeAccBBS)
			}

			// add a starbase to this planet
			if startingPlanet.StarbaseDesignName != "" {
				if err := ug.buildStarbase(player, playerPlanet, startingPlanet.StarbaseDesignName); err != nil {
					return fmt.Errorf("building starbase during universe gen failed: error %w", err)
				}
			}

			// tell the player about their homeworld
			if startingPlanet.Homeworld {
				messager.planetHomeworld(player, playerPlanet)
			}

			// generate some fleets on the homeworld
			if err := ug.generatePlayerFleets(player, playerPlanet, &fleetNum, startingPlanet.StartingFleets); err != nil {
				return fmt.Errorf("generating fleets for planet %s during universe gen failed: error %w", playerPlanet, err)
			}
		}
	}

	return nil
}

// Assign race starting point bonuses to a player's homeworld
func (ug *universeGenerator) assignRaceStartingPointBonuses(race *Race, planet *Planet, extraPoints int, pointsType SpendLeftoverPointsOn) {
	rules := ug.Rules

	// old games have this empty
	if pointsType == SpendLeftoverPointsOnNone {
		pointsType = SpendLeftoverPointsOnSurfaceMinerals
	}
	pointsThreshold := rules.RaceLeftoverPointsPerItem[pointsType]
	switch pointsType {
	case SpendLeftoverPointsOnDefenses:
		if !race.Spec.LivesOnStarbases && extraPoints >= pointsThreshold {
			planet.Defenses += (extraPoints / pointsThreshold)
			extraPoints -= extraPoints / pointsThreshold
		}
	case SpendLeftoverPointsOnFactories:
		if !race.Spec.InnateResources && extraPoints >= pointsThreshold {
			planet.Factories += (extraPoints / pointsThreshold)
			extraPoints -= extraPoints / pointsThreshold
		}
	case SpendLeftoverPointsOnMines:
		if !race.Spec.InnateMining && extraPoints >= pointsThreshold {
			planet.Mines += (extraPoints / pointsThreshold)
			extraPoints -= extraPoints / pointsThreshold
		}
	case SpendLeftoverPointsOnMineralConcentrations:
		// example situation: 25 unspent points; HW has 40I, 30B and 35G concs
		// first we start by increasing B up to 36, using 18 pts
		// G is now lowest, so we bump it up to 37, using 6 points
		// the remaining 1 point goes into surface minerals (since 1 < 3)
		for extraPoints >= pointsThreshold {
			conc := planet.MineralConcentration
			lowestType := conc.HighestType(-1)
			diff := conc.GetAmount(conc.HighestType(2)) - conc.GetAmount(lowestType)
			amtToAdd := Min(extraPoints/pointsThreshold, diff+1)
			planet.MineralConcentration.Set(lowestType, conc.GetAmount(lowestType)+amtToAdd)
			extraPoints -= pointsThreshold * amtToAdd
		}
	}

	// In the event the player has extra points leftover
	// (or selected surface mineral starting points), dump em in
	// _Technically_, we don't really know if Stars! actually did this, but I'm too lazy to check
	for extraPoints > 0 {
		// example situation: 10 points leftover HW with 300I, 400B, 350G starting mins
		// first we add 60kT of I using 6 pts;
		// then, since G is now the lowest mineral,
		// we alternate between adding G and I for the remaining 4 pts
		mins := planet.getCargo().ToMineral()
		lowestType := mins.HighestType(3)
		diff := mins.GetAmount(mins.HighestType(2)) - mins.GetAmount(lowestType)
		amtToAdd := Min(extraPoints, (diff/pointsThreshold)+1) // 70 difference / 10 mins/round => 8 rounds
		planet.Cargo.AddAmount(CargoType(int(lowestType)), amtToAdd*pointsThreshold)
		extraPoints -= amtToAdd
	}
}

// build a starbase on a planet
func (ug *universeGenerator) buildStarbase(player *Player, planet *Planet, designName string) error {
	// the homeworld gets a starbase
	design := player.GetDesignByName(designName)
	if design == nil {
		return fmt.Errorf("no design named %q found for player %s", designName, player)
	}

	design.Spec.NumBuilt++
	design.Spec.NumInstances++
	starbase := newStarbase(player, planet, design, design.Name)
	starbase.Spec = ComputeFleetSpec(&ug.Rules, player, &starbase)
	planet.setStarbase(&starbase)

	ug.Universe.Starbases = append(ug.Universe.Starbases, &starbase)

	return nil
}

func (ug *universeGenerator) generatePlayerFleets(player *Player, planet *Planet, fleetNum *int, startingFleets []StartingFleet) error {
	for _, startingFleet := range startingFleets {
		design := player.GetDesignByName(startingFleet.Name)
		if design == nil {
			return fmt.Errorf("no design named %q found for player %s", startingFleet.Name, player)
		}
		fleet := newFleetForDesign(player, design, 1, *fleetNum, startingFleet.Name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, design.Spec.Engine.IdealSpeed)})
		fleet.OrbitingPlanetNum = planet.Num
		fleet.Spec = ComputeFleetSpec(&ug.Rules, player, &fleet)
		fleet.Fuel = fleet.Spec.FuelCapacity
		fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
		purpose := FleetPurposeFromShipDesignPurpose(design.Purpose)
		fleet.SetTag(TagPurpose, string(purpose))
		ug.Universe.Fleets = append(ug.Universe.Fleets, &fleet)
		design.Spec.NumInstances++
		design.Spec.NumBuilt++
		(*fleetNum)++ // increment the fleet num
	}

	return nil
}

func (ug *universeGenerator) applyGameStartModeModifier() {
	switch ug.StartMode {
	case GameStartModeAccBBS:
		ug.applyAccBBS()
	case GameStartModeMax:
		ug.maxPlayersAndPlanets()
	}
}

func (ug *universeGenerator) applyAccBBS() {
	for _, planet := range ug.Planets {
		// only owned planets will have surface mineral deposits or population,
		// so we can skip unowned ones
		if !planet.Owned() {
			continue
		}

		// Add 25% extra homeworld surface minerals
		// (the help manual lied when it said 20%)
		planet.Cargo = planet.Cargo.AddMineral(planet.Cargo.ToMineral().MultiplyFloat64(0.25))

		// AccBBS adds 20% addiional starting pop (+5K over the default 25K)
		// per 1% of a race's growth rate.
		race := ug.getPlayer(planet.PlayerNum).Race
		planet.Cargo.Colonists += int(float64(planet.Cargo.Colonists*race.GrowthRate) *
			race.Spec.GrowthFactor / 5)
	}
}

func (ug *universeGenerator) maxPlayersAndPlanets() {
	rules := &ug.Rules
	for _, player := range ug.Players {
		// max tech levels, acquire all mt techs
		player.TechLevels = TechLevel{rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel, rules.MaxTechLevel}
		for _, mtTech := range MysteryTraderTechs {
			player.AcquiredTechs[mtTech.Name] = true
		}
	}

	for _, planet := range ug.Planets {
		// max out min concs and add a lot of surface minerals
		planet.MineralConcentration = Mineral{rules.MaxMineralConcentration, rules.MaxMineralConcentration, rules.MaxMineralConcentration}
		planet.Cargo.Ironium = 1_000_000
		planet.Cargo.Boranium = 1_000_000
		planet.Cargo.Germanium = 1_000_000
		if !planet.Owned() {
			continue
		}

		// max out pop & installations on owned planets
		player := ug.Players[planet.PlayerNum-1]
		planet.setPopulation(planet.getMaxPopulation(rules, player, player.Race.GetPlanetHabitability(planet.Hab)))
		if player.Race.Spec.CanBuildDefenses {
			planet.Defenses = 100
		}
		if !player.Race.Spec.InnateMining {
			planet.Mines = getMaxInstallations(player.Race.NumMines, planet.population())
		}
		if !player.Race.Spec.InnateResources {
			planet.Factories = getMaxInstallations(player.Race.NumFactories, planet.population())
		}
	}
}

// create initial starbase designs for a player
func (ug *universeGenerator) createStartingStarbaseDesigns(techStore *TechStore, player *Player, designNum int) []*ShipDesign {
	designs := make([]*ShipDesign, len(player.Race.Spec.StartingPlanets))

	for i, startingPlanet := range player.Race.Spec.StartingPlanets {
		var starbase *ShipDesign
		var purpose ShipDesignPurpose
		switch {
		case i == 0:
			// first design is a starbase, rest are forts of various kinds
			purpose = ShipDesignPurposeStarbase
		case startingPlanet.HasMassDriver:
			purpose = ShipDesignPurposePacketThrower
		case startingPlanet.HasStargate:
			purpose = ShipDesignPurposeStargater
		default:
			purpose = ShipDesignPurposeFort
		}

		starbase = NewShipDesign(player.Num, designNum).
			WithName(startingPlanet.StarbaseDesignName).
			WithHull(startingPlanet.StarbaseHull).
			WithPurpose(purpose).
			WithHullSetNumber(player.DefaultHullSet)
		fillStarbaseSlots(techStore, starbase, startingPlanet)
		designNum++
		designs[i] = starbase
	}

	if player.Race.Spec.LivesOnStarbases {
		// create a starter colony for AR races
		starterColony := NewShipDesign(player.Num, designNum).
			WithName("Starter Colony").
			WithHull(OrbitalFort.Name).
			WithPurpose(ShipDesignPurposeStarterColony).
			WithHullSetNumber(player.DefaultHullSet)
		starterColony.CannotDelete = true
		designs = append(designs, starterColony) // add it to the back
	}

	return designs
}

// Player starting starbases are all the same, regardless of starting tech level
// They get half filled with the starter beam & shield
func fillStarbaseSlots(techStore *TechStore, starbase *ShipDesign, startingPlanet StartingPlanet) {
	hull := techStore.GetHull(starbase.Hull)
	beamWeapon := techStore.GetHullComponentsByCategory(TechCategoryBeamWeapon)[0]
	shield := techStore.GetHullComponentsByCategory(TechCategoryShield)[0]
	var massDriver, stargate TechHullComponent
	var haveDriver, haveGate bool
	for _, hc := range techStore.GetHullComponentsByCategory(TechCategoryOrbital) {
		if hc.PacketSpeed > 0 {
			massDriver = hc
			haveDriver = true
		}
		if hc.SafeRange > 0 {
			stargate = hc
			haveGate = true
		}

		if (haveDriver || !startingPlanet.HasMassDriver) && // we either don't need a driver or have one already
			(haveGate || !startingPlanet.HasStargate) { // we either don't need a driver or have one already
			// we have all the orbital components we need; done with lookup
			break
		}
	}

	placedMassDriver := false
	placedStargate := false
	for index, slot := range hull.Slots {
		var item TechHullComponent
		qty := int(math.Round(float64(slot.Capacity) / 2))
		switch {
		case slot.Type&HullSlotTypeWeapon != 0:
			item = beamWeapon
		case slot.Type&HullSlotTypeShield != 0:
			item = shield
		case slot.Type&HullSlotTypeOrbital != 0:
			qty = 1
			if startingPlanet.HasStargate && !placedStargate {
				item = stargate
			} else if startingPlanet.HasMassDriver && !placedMassDriver {
				item = massDriver
			}
		}

		if item.Name != "" {
			starbase.Slots = append(starbase.Slots, ShipDesignSlot{
				HullComponent: item.Name,
				HullSlotIndex: index + 1,
				Quantity:      qty,
			})
			placedStargate = placedStargate || item.SafeRange > 0
			placedMassDriver = placedMassDriver || item.PacketSpeed > 0
		}
	}
}

func (ug *universeGenerator) generatePlayerRelations() {
	for _, player := range ug.Players {
		player.Relations = player.defaultRelationships(ug.Players, ug.ComputerPlayersFormAlliances)
	}
}

func (ug *universeGenerator) generatePlayerIntel() error {
	for _, player := range ug.Players {

		// discover other players
		player.PlayerIntels.PlayerIntels = player.defaultPlayerIntels(ug.Players)
		player.PlayerIntels.ScoreIntels = make([]ScoreIntel, len(ug.Players))

		// do initial scans
		scanner := newPlayerScanner(ug.Universe, ug.Players, &ug.Rules, player)
		if err := scanner.scan(); err != nil {
			return err
		}

	}

	return nil
}
