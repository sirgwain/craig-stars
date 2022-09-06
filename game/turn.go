package game

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type turn struct {
	game *FullGame
}

type turnGenerator interface {
	GenerateTurn() error
}

func NewTurnGenerator(game *FullGame) turnGenerator {
	t := turn{game}

	t.game.Universe.buildMaps()

	return &t
}

// generate a new turn
func (t *turn) GenerateTurn() error {
	log.Debug().
		Uint64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year).
		Msgf("begin generating turn")
	t.game.Year++

	// reset players for start of the turn
	for _, player := range t.game.Players {
		player.Messages = []PlayerMessage{}
		player.leftoverResources = 0
		player.Spec = computePlayerSpec(player, &t.game.Rules)
	}

	t.computePlanetSpecs()
	t.fleetAge()

	// wp0 tasks
	t.fleetScrap()
	t.fleetUnload0()
	t.fleetColonize()
	t.fleetLoad0()
	t.fleetLoadDunnage0()
	t.fleetMerge0()
	t.fleetRoute0()
	t.packetMove0()

	// move stuff through space
	t.mysteryTraderMove()
	t.fleetMove()
	t.fleetReproduce()
	t.decaySalvage()
	t.decayPackets()
	t.wormholeJiggle()
	t.detonateMines()
	t.planetMine()
	t.fleetRemoteMineAR()
	t.planetProduction()
	t.playerResearch()
	t.researchStealer()
	t.permaform()
	t.planetGrow()
	t.packetMove1()
	t.fleetRefuel()
	t.randomCometStrike()
	t.randomMineralDeposit()
	t.randomPlanetaryChange()
	t.fleetBattle()
	t.fleetBomb()
	t.mysteryTraderMeet()
	t.fleetRemoteMine0()

	// wp1 tasks
	t.fleetUnload1()
	t.fleetColonize() // colonize wp1 after arriving at a planet
	t.fleetLoad1()
	t.fleetLoadDunnage1()
	t.decayMines()
	t.fleetLayMines()
	t.fleetTransfer()
	t.fleetMerge1()
	t.fleetRoute1()
	t.instaform()
	t.fleetSweepMines()
	t.fleetRepair()
	t.remoteTerraform()

	// reset all players
	// and do player specific things like scanning
	// and patrol orders
	for _, player := range t.game.Players {
		player.Spec = computePlayerSpec(player, &t.game.Rules)

		scanner := newPlayerScanner(t.game.Universe, &t.game.Rules, player)
		if err := scanner.scan(); err != nil {
			return err
		}
		t.playerInfoDiscover(player)
		t.fleetPatrol(player)
		t.calculateScore(player)
		t.checkVictory(player)

		player.SubmittedTurn = false
		pmo := t.game.GetPlayerMapObjects(player.Num)
		ai := NewAIPlayer(player, pmo)
		ai.processTurn()

		for _, f := range pmo.Fleets {
			f.ComputeFuelUsage(ai.Player)
		}
	}

	log.Info().
		Uint64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year-1).
		Msgf("generated turn")
	return nil
}

// update all planet specs with the latest info
// useful before turn generation and after building
func (t *turn) computePlanetSpecs() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			player := t.game.Players[planet.PlayerNum]
			planet.Spec = ComputePlanetSpec(&t.game.Rules, planet, player)
		}
	}

}
func (t *turn) fleetAge() {

}

func (t *turn) fleetScrap() {

}

func (t *turn) fleetUnload0() {

}

func (t *turn) fleetColonize() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskColonize {
			player := t.game.Players[fleet.PlayerNum]

			if wp0.TargetType != MapObjectTypePlanet {
				messager.colonizeNonPlanet(player, fleet)
				continue
			}

			if wp0.TargetNum == NoTarget {
				err := fmt.Errorf("%s attempted to colonize a planet but didn't target a planet", fleet.Name)
				log.Err(err).
					Uint64("GameID", t.game.ID).
					Uint64("PlayerID", player.ID).
					Str("Fleet", fleet.Name)
				messager.error(player, err)
				continue
			}

			planet := t.game.GetPlanet(wp0.TargetNum)
			if planet.Owned() {
				messager.colonizeOwnedPlanet(player, fleet)
				continue
			}

			if !fleet.Spec.Colonizer {
				messager.colonizeWithNoModule(player, fleet)
				continue
			}

			if fleet.Cargo.Colonists == 0 {
				messager.colonizeWithNoColonists(player, fleet)
				continue
			}

			log.Debug().
				Uint64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", player.Num).
				Str("Planet", planet.Name).
				Str("Fleet", fleet.Name).
				Msgf("colonized planet")

			// colonize the planet and scrap the fleet
			fleet.colonizePlanet(&t.game.Rules, player, planet)
			fleet.scrap(&t.game.Rules, player, planet)
			messager.planetColonized(player, planet)
		}
	}
}

func (t *turn) fleetLoad0() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]

		if wp0.Task == WaypointTaskTransport {
			player := t.game.Players[fleet.PlayerNum]
			dest := t.game.GetCargoHolder(wp0.TargetType, wp0.TargetNum, wp0.TargetPlayerNum)
			if dest == nil {
				salvage := NewSalvage(fleet.PlayerNum, fleet.Position, Cargo{})
				dest = &salvage
			}

			// build a cargo object for each transfer
			transferCargo := Cargo{}
			for cargoType, task := range wp0.getTransportTasks() {
				capacity := fleet.Spec.CargoCapacity - fleet.Cargo.Total() - transferCargo.Total()
				transferAmount := 0
				availableToLoad := dest.GetCargo().GetAmount(cargoType)
				currentAmount := fleet.Cargo.GetAmount(cargoType)
				_ = currentAmount
				switch task.Action {
				case TransportActionLoadAll:
					// load all available, based on our constraints
					transferAmount = MinInt(availableToLoad, capacity)
				}

				if transferAmount != 0 {
					transferCargo = transferCargo.WithCargo(cargoType, transferAmount)
					t.game.Transfer(fleet, dest, cargoType, transferAmount)
					messager.fleetTransportedCargo(player, fleet, dest, cargoType, transferAmount)
				}
			}

		}
	}
}

func (t *turn) fleetLoadDunnage0() {

}

func (t *turn) fleetMerge0() {

}

func (t *turn) fleetRoute0() {

}

func (t *turn) packetMove0() {

}

func (t *turn) mysteryTraderMove() {

}

func (t *turn) fleetMove() {

	for _, fleet := range t.game.Fleets {
		if !fleet.Starbase {
			// remove the fleet from the list of map objects at it's current location
			originalPosition := fleet.Position

			if len(fleet.Waypoints) > 1 {
				player := t.game.Players[fleet.PlayerNum]
				wp0 := fleet.Waypoints[0]
				wp1 := fleet.Waypoints[1]

				if wp1.WarpFactor == StargateWarpFactor {
					// yeah, gate!
					fleet.gateFleet(&t.game.Rules, player)
				} else {
					fleet.moveFleet(t.game.Universe, &t.game.Rules, player)
				}

				// remove the previous waypoint, it's been processed already
				if fleet.RepeatOrders && !wp0.PartiallyComplete {
					// if we are supposed to repeat orders,
					fleet.Waypoints = append(fleet.Waypoints, wp0)
				}

				// update the game dictionaries with this fleet's new position
				t.game.MoveFleet(fleet, originalPosition)
			} else {
				fleet.PreviousPosition = &originalPosition
				fleet.WarpSpeed = 0
				fleet.Heading = Vector{}
			}
		}

	}
}

func (t *turn) fleetReproduce() {

}

func (t *turn) decaySalvage() {

}

func (t *turn) decayPackets() {

}

func (t *turn) wormholeJiggle() {

}

func (t *turn) detonateMines() {

}

// mine all owned planets for minerals
func (t *turn) planetMine() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			planet.Cargo = planet.Cargo.AddMineral(planet.Spec.MineralOutput)
			planet.MineYears.AddInt(planet.Mines)
			planet.reduceMineralConcentration(&t.game.Rules)
		}
	}
}

func (t *turn) fleetRemoteMineAR() {

}

func (t *turn) planetProduction() {
	for _, planet := range t.game.Planets {
		if planet.Owned() && len(planet.ProductionQueue) > 0 {
			player := t.game.Players[planet.PlayerNum]
			result := planet.produce(player)
			for _, token := range result.tokens {
				fleet := t.buildFleet(player, planet, token)
				messager.fleetBuilt(player, planet, &fleet, token.Quantity)
			}
		}
	}
}

// build a fleet with some number of tokens
func (t *turn) buildFleet(player *Player, planet *Planet, token ShipToken) Fleet {
	player.Stats.FleetsBuilt++
	player.Stats.TokensBuilt += token.Quantity

	fleetNum := t.game.getNextFleetNum(player.Num)
	fleet := NewFleetForToken(player, fleetNum, token, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, token.Design.Spec.IdealSpeed)})
	fleet.Position = planet.Position
	fleet.BattlePlanID = player.BattlePlans[0].ID
	fleet.Spec = ComputeFleetSpec(&t.game.Rules, player, &fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	fleet.OrbitingPlanetNum = planet.Num

	t.game.Fleets = append(t.game.Fleets, &fleet)
	return fleet
}

func (t *turn) playerResearch() {

}

func (t *turn) researchStealer() {

}

func (t *turn) permaform() {

}

// grow all owned planets by some population
func (t *turn) planetGrow() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			// player := t.game.Players[*planet.PlayerNum]
			planet.SetPopulation(planet.Population() + planet.Spec.GrowthAmount)
			planet.Dirty = true // flag for update
		}
	}
}

func (t *turn) packetMove1() {

}

func (t *turn) fleetRefuel() {

}

func (t *turn) randomCometStrike() {

}

func (t *turn) randomMineralDeposit() {

}

func (t *turn) randomPlanetaryChange() {

}

func (t *turn) fleetBattle() {

}

func (t *turn) fleetBomb() {

}

func (t *turn) mysteryTraderMeet() {

}

func (t *turn) fleetRemoteMine0() {

}

func (t *turn) fleetUnload1() {

}

func (t *turn) fleetLoad1() {

}

func (t *turn) fleetLoadDunnage1() {

}

func (t *turn) decayMines() {

}

func (t *turn) fleetLayMines() {

}

func (t *turn) fleetTransfer() {

}

func (t *turn) fleetMerge1() {

}

func (t *turn) fleetRoute1() {

}

func (t *turn) instaform() {

}

func (t *turn) fleetSweepMines() {

}

func (t *turn) fleetRepair() {

}

func (t *turn) remoteTerraform() {

}

// Update a player's information about other players
// If a player scanned another player's fleet or planet, they discover the player's
// race name
func (t *turn) playerInfoDiscover(player *Player) {

}

func (t *turn) fleetPatrol(player *Player) {

}

func (t *turn) calculateScore(player *Player) {

}

func (t *turn) checkVictory(player *Player) {

}
