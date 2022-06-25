package game

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// Check if a position vector is a mininum distance away from all other points
func isPositionValid(pos Vector, occupiedLocations *[]Vector, minDistance float64) bool {
	minDistanceSquared := minDistance * minDistance

	for _, to := range *occupiedLocations {
		if pos.DistanceSquaredTo(&to) <= minDistanceSquared {
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
		planet.Randomize(&g.Rules)

		g.Planets[i] = planet
		planetsByPosition[pos] = &planet
		occupiedLocations = append(occupiedLocations, pos)
	}

	// shuffle these so id 1 is not always the first planet in the list
	// later on we will add homeworlds based on first planet, second planet, etc
	// g.Rules.Random.Shuffle(len(g.Planets), func(i, j int) { g.Planets[i], g.Planets[j] = g.Planets[j], g.Planets[i] })

	return nil
}

func generateWormholes(game *Game) {

}

func generatePlayerTechLevelss(game *Game) {

}

func generatePlayerPlanss(game *Game) {

}

func generatePlayerShipDesignss(game *Game) {

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

	ownedPlanets := []Planet{}
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

	for playerIndex := range game.Players {
		player := &game.Players[playerIndex]
		minPlayerDistance := (area.X + area.Y) / 2.0 / float64(len(game.Players))

		// find a homeworld that is a min distance from other homeworlds
		var homeworld *Planet
		for i := range game.Planets {
			planet := &game.Planets[i]
			if !planet.Owned() && (len(ownedPlanets) == 0 || planet.shortestDistanceToPlanets(&ownedPlanets) > minPlayerDistance) {
				homeworld = planet
				break
			}
		}

		if homeworld == nil {
			return fmt.Errorf("failed to find homeworld for player %v among %d planets", player, len(game.Planets))
		}

		ownedPlanets = append(ownedPlanets, *homeworld)

		// make a new homeworld

		if err := homeworld.initHomeworld(player, &game.Rules, homeworldMinConc, homeworldSurfaceMinerals); err != nil {
			return err
		}
		// player.Planets = append(player.Planets, *homeworld)
	}

	return nil
}

func generatePlayerFleets(game *Game) {

}

func applyGameStartModeModifier(game *Game) {

}
