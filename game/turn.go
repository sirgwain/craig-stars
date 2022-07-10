package game

func generateTurn(game *Game) error {
	game.Year++
	game.computeSpecs()

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
	for i := range game.Players {
		player := &game.Players[i]

		player.Spec = computePlayerSpec(player, &game.Rules)

		game.playerScan(player)
		game.playerInfoDiscover(player)
		game.fleetPatrol(player)
		game.calculateScore(player)
		game.checkVictory(player)

		player.SubmittedTurn = false
	}

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

// scan planets, fleets, etc for a player
func (game *Game) playerScan(player *Player) {

	for i := range game.Planets {
		planet := &game.Planets[i]
		if planet.OwnedBy(player.Num) {
			discoverPlanet(&game.Rules, player, planet, false)
		}
	}

}

func (game *Game) playerInfoDiscover(player *Player) {

}

func (game *Game) fleetPatrol(player *Player) {

}

func (game *Game) calculateScore(player *Player) {

}

func (game *Game) checkVictory(player *Player) {

}
