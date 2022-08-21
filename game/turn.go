package game

import "github.com/rs/zerolog/log"

func generateTurn(game *Game) error {
	log.Debug().
		Uint("GameID", game.ID).
		Str("Name", game.Name).
		Int("Year", game.Year).
		Msgf("begin generating turn")
	game.Year++
	game.computeSpecs()
	game.buildMaps()

	// reset players for start of the turn
	for i := range game.Players {
		player := &game.Players[i]
		player.Messages = []PlayerMessage{}
		player.LeftoverResources = 0
	}

	game.fleetAge()
	game.fleetScrap()
	game.fleetUnload0()
	game.fleetColonize0()
	game.fleetLoad0()
	game.fleetLoadDunnage0()
	game.fleetMerge0()
	game.fleetRoute0()
	game.packetMove0()
	game.mysteryTraderMove()
	game.fleetMove()
	game.fleetReproduce()
	game.decaySalvage()
	game.decayPackets()
	game.wormholeJiggle()
	game.detonateMines()
	game.planetMine()
	game.fleetRemoteMineAR()
	game.planetProduction()
	game.playerResearch()
	game.researchStealer()
	game.permaform()
	game.planetGrow()
	game.packetMove1()
	game.fleetRefuel()
	game.randomCometStrike()
	game.randomMineralDeposit()
	game.randomPlanetaryChange()
	game.fleetBattle()
	game.fleetBomb()
	game.mysteryTraderMeet()
	game.fleetRemoteMine0()
	game.fleetUnload1()
	game.fleetColonize1()
	game.fleetLoad1()
	game.fleetLoadDunnage1()
	game.decayMines()
	game.fleetLayMines()
	game.fleetTransfer()
	game.fleetMerge1()
	game.fleetRoute1()
	game.instaform()
	game.fleetSweepMines()
	game.fleetRepair()
	game.remoteTerraform()

	// reset all players
	// and do player specific things like scanning
	// and patrol orders
	for i := range game.Players {
		player := &game.Players[i]

		player.Spec = computePlayerSpec(player, &game.Rules)

		game.updatePlayerOwnedObjects(player)

		if err := game.playerScan(player); err != nil {
			return err
		}
		game.playerInfoDiscover(player)
		game.fleetPatrol(player)
		game.calculateScore(player)
		game.checkVictory(player)

		player.SubmittedTurn = false
		processTurn(player)
	}

	log.Info().
		Uint("GameID", game.ID).
		Str("Name", game.Name).
		Int("Year", game.Year-1).
		Msgf("generated turn")
	return nil
}

func (game *Game) fleetAge() {

}

func (game *Game) fleetScrap() {

}

func (game *Game) fleetUnload0() {

}

func (game *Game) fleetColonize0() {

}

func (game *Game) fleetLoad0() {

}

func (game *Game) fleetLoadDunnage0() {

}

func (game *Game) fleetMerge0() {

}

func (game *Game) fleetRoute0() {

}

func (game *Game) packetMove0() {

}

func (game *Game) mysteryTraderMove() {

}

func (game *Game) fleetMove() {

	for i := range game.Fleets {
		fleet := &game.Fleets[i]
		if !fleet.Starbase {
			// remove the fleet from the list of map objects at it's current location
			originalPosition := fleet.Position

			if len(fleet.Waypoints) > 1 {
				player := &game.Players[fleet.PlayerNum]
				wp0 := fleet.Waypoints[0]
				wp1 := fleet.Waypoints[1]
				totalDist := fleet.Position.DistanceTo(wp1.Position)

				if wp1.WarpFactor == StargateWarpFactor {
					// yeah, gate!
					fleet.gateFleet(&game.Rules, player, wp0, wp1, totalDist)
				} else {
					fleet.moveFleet(game, &game.Rules, player, wp0, wp1, totalDist)
				}

				// remove the previous waypoint, it's been processed already
				if fleet.RepeatOrders && !wp0.PartiallyComplete {
					// if we are supposed to repeat orders,
					fleet.Waypoints = append(fleet.Waypoints, wp0)
				}

				// update the game dictionaries with this fleet's new position
				delete(game.FleetsByPosition, originalPosition)
				game.FleetsByPosition[fleet.Position] = fleet
				fleet.Dirty = true
			} else {
				fleet.PreviousPosition = &originalPosition
				fleet.WarpSpeed = 0
				fleet.Heading = Vector{}
			}
		}

	}
}

func (game *Game) fleetReproduce() {

}

func (game *Game) decaySalvage() {

}

func (game *Game) decayPackets() {

}

func (game *Game) wormholeJiggle() {

}

func (game *Game) detonateMines() {

}

// mine all owned planets for minerals
func (game *Game) planetMine() {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() {
			planet.Cargo = planet.Cargo.AddMineral(planet.Spec.MineralOutput)
			planet.MineYears.AddInt(planet.Mines)
			planet.reduceMineralConcentration(game)
		}
	}
}

func (game *Game) fleetRemoteMineAR() {

}

func (game *Game) planetProduction() {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() && len(planet.ProductionQueue) > 0 {
			planet.produce(game)
		}
	}
}

func (game *Game) playerResearch() {

}

func (game *Game) researchStealer() {

}

func (game *Game) permaform() {

}

// grow all owned planets by some population
func (game *Game) planetGrow() {
	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.Owned() {
			// player := &game.Players[*planet.PlayerNum]
			planet.SetPopulation(planet.Population() + planet.Spec.GrowthAmount)
			planet.Dirty = true // flag for update
		}
	}
}

func (game *Game) packetMove1() {

}

func (game *Game) fleetRefuel() {

}

func (game *Game) randomCometStrike() {

}

func (game *Game) randomMineralDeposit() {

}

func (game *Game) randomPlanetaryChange() {

}

func (game *Game) fleetBattle() {

}

func (game *Game) fleetBomb() {

}

func (game *Game) mysteryTraderMeet() {

}

func (game *Game) fleetRemoteMine0() {

}

func (game *Game) fleetUnload1() {

}

func (game *Game) fleetColonize1() {

}

func (game *Game) fleetLoad1() {

}

func (game *Game) fleetLoadDunnage1() {

}

func (game *Game) decayMines() {

}

func (game *Game) fleetLayMines() {

}

func (game *Game) fleetTransfer() {

}

func (game *Game) fleetMerge1() {

}

func (game *Game) fleetRoute1() {

}

func (game *Game) instaform() {

}

func (game *Game) fleetSweepMines() {

}

func (game *Game) fleetRepair() {

}

func (game *Game) remoteTerraform() {

}

// Update a player's information about other players
// If a player scanned another player's fleet or planet, they discover the player's
// race name
func (game *Game) playerInfoDiscover(player *Player) {

}

func (game *Game) fleetPatrol(player *Player) {

}

func (game *Game) calculateScore(player *Player) {

}

func (game *Game) checkVictory(player *Player) {

}
