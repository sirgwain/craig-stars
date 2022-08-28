package game

import (
	"fmt"
	"math"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
)

// Check if a position vector is a mininum distance away from all other points
func isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool {
	minDistanceSquared := minDistance * minDistance

	for _, to := range *occupiedLocations {
		if pos.DistanceSquaredTo(to) <= minDistanceSquared {
			return false
		}
	}
	return true
}

func generatePlanets(g *Game, area Vector) error {

	numPlanets, err := g.Rules.GetNumPlanets(g.Size, g.Density)
	if err != nil {
		return err
	}

	log.Debug().Msgf("%s: Generating %d planets in universe size %v for ", g, numPlanets, area)

	names := planetNames
	g.Rules.Random.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	g.Planets = make([]Planet, numPlanets)

	planetsByPosition := make(map[Vector]*Planet, numPlanets)
	occupiedLocations := make([]Vector, numPlanets)

	for i := 0; i < numPlanets; i++ {

		// find a valid position for the planet
		posCheckCount := 0
		pos := Vector{X: float64(g.Rules.Random.Intn(int(area.X))), Y: float64(g.Rules.Random.Intn(int(area.Y)))}
		for !isPositionValid(pos, &occupiedLocations, float64(g.Rules.PlanetMinDistance)) {
			pos = Vector{X: float64(g.Rules.Random.Intn(int(area.X))), Y: float64(g.Rules.Random.Intn(int(area.Y)))}
			posCheckCount++
			if posCheckCount > 1000 {
				return fmt.Errorf("failed to find a valid position for a planet in 1000 tries, min: %d, numPlanets: %d, area: %v", g.Rules.PlanetMinDistance, numPlanets, area)
			}
		}

		// setup a new planet
		planet := NewPlanet(g.ID)
		planet.Name = names[i]
		planet.Num = int(i + 1)
		planet.Position = pos
		planet.randomize(&g.Rules)

		g.Planets[i] = *planet
		planetsByPosition[pos] = planet
		occupiedLocations = append(occupiedLocations, pos)
	}

	// shuffle these so id 1 is not always the first planet in the list
	// later on we will add homeworlds based on first planet, second planet, etc
	// g.Rules.Random.Shuffle(len(g.Planets), func(i, j int) { g.Planets[i], g.Planets[j] = g.Planets[j], g.Planets[i] })

	return nil
}

func generateWormholes(game *Game) {

}

func generatePlayerTechLevels(game *Game) {
	for i := range game.Players {
		// the first time we allocate an array of planets
		player := &game.Players[i]
		player.TechLevels = TechLevel(player.Race.Spec.StartingTechLevels)
	}
}

func generatePlayerPlans(game *Game) {

}

// generate designs for each player
func generatePlayerShipDesigns(game *Game) {
	for i := range game.Players {
		// the first time we allocate an array of planets
		player := &game.Players[i]
		designNames := mapset.NewSet[string]()
		for _, startingPlanet := range player.Race.Spec.StartingPlanets {
			for _, startingFleet := range startingPlanet.StartingFleets {
				if designNames.Contains(startingFleet.Name) {
					// only one design per name, i.e. Scout, Armored Probe
					continue
				}
				techStore := game.Rules.Techs
				hull := techStore.GetHull(string(startingFleet.HullName))
				design := designShip(techStore, hull, startingFleet.Name, player, player.DefaultHullSet, startingFleet.Purpose)
				design.HullSetNumber = int(startingFleet.HullSetNumber)
				design.Purpose = startingFleet.Purpose
				design.Spec = ComputeShipDesignSpec(&game.Rules, player, design)
				player.Designs = append(player.Designs, design)
			}
		}

		starbaseDesigns := getStartingStarbaseDesigns(game.Rules.Techs, player)

		for i := range starbaseDesigns {
			design := &starbaseDesigns[i]
			design.Purpose = ShipDesignPurposeStarbase
			design.Spec = ComputeShipDesignSpec(&game.Rules, player, design)
			player.Designs = append(player.Designs, design)
		}
	}

}

// have each player discover all the planets in the universe
func generatePlayerPlanetReports(game *Game) error {
	for i := range game.Players {
		// the first time we allocate an array of planets
		player := &game.Players[i]
		player.PlanetIntels = make([]PlanetIntel, len(game.Planets))
		for j := range game.Planets {
			if err := discoverPlanet(&game.Rules, player, &game.Planets[j], false); err != nil {
				return err
			}
		}
	}
	return nil
}

func generatePlayerHomeworlds(game *Game, area Vector) error {

	ownedPlanets := []*Planet{}
	rules := game.Rules
	random := game.Rules.Random

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

	for playerIndex := range game.Players {
		player := &game.Players[playerIndex]
		minPlayerDistance := float64(area.X+area.Y) / 2.0 / float64(len(game.Players)+1)
		fleetNum := 1
		var homeworld *Planet

		for startingPlanetIndex, startingPlanet := range player.Race.Spec.StartingPlanets {
			// find a playerPlanet that is a min distance from other homeworlds
			var playerPlanet *Planet
			if startingPlanetIndex > 0 {

				// extra planets are close to the homeworld
				for i := range game.Planets {
					planet := &game.Planets[i]
					distToHomeworld := planet.Position.DistanceSquaredTo(homeworld.Position)
					if !planet.Owned() && (distToHomeworld <= float64(rules.MaxExtraWorldDistance*rules.MaxExtraWorldDistance) && distToHomeworld >= float64(rules.MinExtraWorldDistance*rules.MinExtraWorldDistance)) {
						playerPlanet = planet
						break
					}
				}

			} else {
				// homeworld should be distant from other players
				for i := range game.Planets {
					planet := &game.Planets[i]
					if !planet.Owned() && (len(ownedPlanets) == 0 || planet.shortestDistanceToPlanets(&ownedPlanets) > minPlayerDistance) {
						playerPlanet = planet
						break
					}
				}

				homeworld = playerPlanet
			}

			if playerPlanet == nil {
				return fmt.Errorf("failed to find homeworld for player %v among %d planets, minDistance: %0.1f", player, len(game.Planets), minPlayerDistance)
			}

			ownedPlanets = append(ownedPlanets, playerPlanet)
			player.Planets = append(player.Planets, playerPlanet)

			// our starting planet starts with default fleets
			surfaceMinerals := homeworldSurfaceMinerals
			if startingPlanetIndex != 0 {
				surfaceMinerals = extraWorldSurfaceMinerals
			}

			// first planet is a homeworld
			// make a new homeworld
			if err := playerPlanet.initStartingWorld(player, &game.Rules, startingPlanet, homeworldMinConc, surfaceMinerals); err != nil {
				return err
			}
			// generate some fleets on the homeworld
			if err := generatePlayerFleets(game, player, playerPlanet, &fleetNum, startingPlanet.StartingFleets); err != nil {
				return err
			}

			messager.longMessage(player)
		}
	}

	return nil
}

func generatePlayerFleets(game *Game, player *Player, planet *Planet, fleetNum *int, startingFleets []StartingFleet) error {
	for _, startingFleet := range startingFleets {
		design := player.GetDesign(string(startingFleet.Name))
		if design == nil {
			return fmt.Errorf("no design named %s found for player %s", startingFleet.Name, player)
		}

		fleet := NewFleet(player, design, *fleetNum, startingFleet.Name, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, design.Spec.IdealSpeed)})
		fleet.OrbitingPlanetNum = planet.Num
		fleet.Spec = ComputeFleetSpec(&game.Rules, player, &fleet)
		fleet.Fuel = fleet.Spec.FuelCapacity
		game.Fleets = append(game.Fleets, fleet)
		(*fleetNum)++ // increment the fleet num
	}

	return nil
}

func applyGameStartModeModifier(game *Game) {

}

// get the initial starbase designs for a player
func getStartingStarbaseDesigns(techStore *TechStore, player *Player) []ShipDesign {
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
