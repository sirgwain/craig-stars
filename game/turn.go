package game

func generateTurn(game *Game) error {
	game.Year++
	game.computeSpecs()
	mine(game)
	produce(game)
	research(game)
	grow(game)

	for i := range game.Players {
		player := &game.Players[i]
		player.Messages = []PlayerMessage{}
		player.LeftoverResources = 0
		player.Spec = computePlayerSpec(player, &game.Rules)
		playerScan(game, player)
	}

	// reset all players
	for i := range game.Players {
		player := &game.Players[i]
		player.SubmittedTurn = false
	}

	return nil
}

// grow all owned planets by some population
func grow(game *Game) {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() {
			// player := &game.Players[*planet.PlayerNum]
			planet.SetPopulation(planet.Population() + planet.Spec.GrowthAmount)
			planet.Dirty = true // flag for update
		}
	}
}

// mine all owned planets for minerals
func mine(game *Game) {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() {
			planet.Cargo = planet.Cargo.AddMineral(planet.Spec.MineralOutput)
		}
	}

}

func produce(game *Game) {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() && len(planet.ProductionQueue) > 0 {
			planet.produce(game)
		}
	}
}

func research(game *Game) {
}

// scan planets, fleets, etc for a player
func playerScan(game *Game, player *Player) {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.OwnedBy(player.Num) {
			discoverPlanet(&game.Rules, player, planet, false)
		}
	}
}
