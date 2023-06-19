package game

func generateTurn(game *Game) error {
	game.Year++
	game.computeSpecs()
	mine(game)
	produce(game)
	research(game)
	grow(game)

	for i := range game.Players {
		playerScan(game, &game.Players[i])
	}

	// reset all players
	for i := range game.Players {
		player := &game.Players[i]
		player.SubmittedTurn = false
	}

	return nil
}

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

func mine(game *Game) {
}

func produce(game *Game) {
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
