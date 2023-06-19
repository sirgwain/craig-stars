package cs

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type turn struct {
	game *FullGame
}

type turnGenerator interface {
	generateTurn() error
}

func newTurnGenerator(game *FullGame) turnGenerator {
	t := turn{game}

	t.game.Universe.buildMaps(game.Players)

	return &t
}

// generate a new turn
func (t *turn) generateTurn() error {
	log.Debug().
		Int64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year).
		Msgf("begin generating turn")
	t.game.Year++

	// reset players for start of the turn
	for _, player := range t.game.Players {
		player.Messages = []PlayerMessage{}
		player.Battles = []BattleRecord{}
		player.leftoverResources = 0
		player.Spec = computePlayerSpec(player, &t.game.Rules, t.game.Planets)
	}

	t.computePlanetSpecs()

	// wp0 tasks
	t.fleetInit()
	t.fleetScrap()
	t.fleetUnload()
	t.fleetColonize()
	t.fleetLoad()
	t.fleetMerge0()
	t.fleetRoute()
	t.fleetMarkWaypointsProcessed()

	// move stuff through space
	t.packetMove0()
	t.mysteryTraderMove()
	t.fleetMove()
	t.fleetReproduce()
	t.decaySalvage()
	t.decayPackets()
	t.wormholeJiggle()
	t.detonateMines()
	t.planetMine()
	t.fleetRemoteMineAR() // sort of a wp1 task, for AR races it happens before production
	t.planetProduction()
	t.playerResearch()
	t.permaform()
	t.planetGrow()
	t.packetMove1()
	t.fleetRefuel() // refuel after production so fleets will refuel at planets that just built a starbase this turn
	t.randomCometStrike()
	t.randomMineralDeposit()
	t.randomPlanetaryChange()
	t.fleetBattle()
	t.fleetBomb()
	t.mysteryTraderMeet()

	// wp1 tasks
	t.fleetRemoteMine()
	t.fleetUnload()
	t.fleetColonize() // colonize wp1 after arriving at a planet
	t.fleetScrap()
	t.fleetLoad()
	t.decayMines()
	t.fleetLayMines()
	t.fleetTransferOwner()
	t.fleetMerge1()
	t.fleetRoute()

	// do some final stuff like instaforming and repairing
	t.instaform()
	t.fleetSweepMines()
	t.fleetRepair()
	t.remoteTerraform()

	// reset all players
	// and do player specific things like scanning
	// and patrol orders
	for _, player := range t.game.Players {
		player.Spec = computePlayerSpec(player, &t.game.Rules, t.game.Planets)

		scanner := newPlayerScanner(t.game.Universe, t.game.Players, &t.game.Rules, player)
		if err := scanner.scan(); err != nil {
			return fmt.Errorf("scan universe and update player intel -> %w", err)
		}
		t.discoverPlayer(player)
		t.fleetPatrol(player)
		t.calculateScore(player)
		t.checkVictory(player)

		player.SubmittedTurn = false
	}

	log.Info().
		Int64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year-1).
		Msgf("generated turn")
	return nil
}

// update all planet specs with the latest info
// useful before turn generation and after building
func (t *turn) computePlanetSpecs() {
	for _, planet := range t.game.Planets {
		if planet.owned() {
			player := t.game.Players[planet.PlayerNum-1]
			planet.Spec = computePlanetSpec(&t.game.Rules, player, planet)
		}
	}

}

// fleetInit will reset any fleet data before processing
func (t *turn) fleetInit() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		wp0.processed = false

		if wp0.Task == WaypointTaskTransport {
			wp0.WaitAtWaypoint = false
		}
	}
}

// scrap a fleet at wp0/wp1
func (t *turn) fleetScrap() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskScrapFleet {
			t.scrapFleet(fleet)
		}
	}
}

// scrap a fleet giving a planet resources or creating salvage
func (t *turn) scrapFleet(fleet *Fleet) {
	player := t.game.getPlayer(fleet.PlayerNum)
	planet := t.game.getOrbitingPlanet(fleet)

	cost := fleet.getScrapAmount(&t.game.Rules, player, planet)

	if planet != nil {
		// scrap over a planet
		planet.Cargo = planet.Cargo.AddCostMinerals(cost)
		if planet.OwnedBy(player.Num) {
			planet.BonusResources += cost.Resources
		}
	} else {
		// create salvage
		t.game.createSalvage(fleet.Position, player.Num, cost.ToCargo())
	}

	messager.fleetScrapped(player, fleet, cost.ToCargo().Total(), cost.Resources, planet)
	t.game.deleteFleet(fleet)
}

// fleetColonize will attempt to colonize planets for any fleets with the Colonize WaypointTask
func (t *turn) fleetColonize() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskColonize {
			player := t.game.Players[fleet.PlayerNum-1]

			if wp.TargetType != MapObjectTypePlanet {
				messager.colonizeNonPlanet(player, fleet)
				continue
			}

			if wp.TargetNum == None {
				err := fmt.Errorf("%s attempted to colonize a planet but didn't target a planet", fleet.Name)
				log.Err(err).
					Int64("GameID", t.game.ID).
					Str("Fleet", fleet.Name)
				messager.error(player, err)
				continue
			}

			planet := t.game.getPlanet(wp.TargetNum)
			if planet.owned() {
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
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", player.Num).
				Str("Planet", planet.Name).
				Str("Fleet", fleet.Name).
				Msgf("colonized planet")

			// colonize the planet and scrap the fleet
			fleet.colonizePlanet(&t.game.Rules, player, planet)
			t.scrapFleet(fleet)
			t.game.deleteFleet(fleet)
			messager.planetColonized(player, planet)
		}
	}
}

// fleetUnload executes wp0/wp1 unload transport tasks for fleets
func (t *turn) fleetUnload() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			var salvage *Salvage
			if dest == nil {

				salvage = t.game.createSalvage(fleet.Position, fleet.PlayerNum, Cargo{})
				dest = salvage
			}

			for cargoType, task := range wp.getTransportTasks() {
				transferAmount, waitAtWaypoint := fleet.getCargoUnloadAmount(dest, cargoType, task)

				wp.WaitAtWaypoint = wp.WaitAtWaypoint || waitAtWaypoint

				t.fleetTransferCargo(fleet, transferAmount, cargoType, dest)
			}

			// we tried to load/unload from empty space but we didn't deposit any cargo
			// into the salvage, so remove this empty salvage
			if salvage != nil && salvage.Cargo.Total() == 0 {
				t.game.deleteSalvage(salvage)
			}

		}
	}
}

func (t *turn) fleetLoad() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			if dest == nil {
				// can't load from space
				return
			}

			// dunnage tasks are done after regular tasks
			type dunnageTask struct {
				cargoType CargoType
				task      WaypointTransportTask
			}
			dunnageTasks := []dunnageTask{}

			// process regular load tasks
			for cargoType, task := range wp.getTransportTasks() {
				if task.Action == TransportActionLoadDunnage {
					dunnageTasks = append(dunnageTasks, dunnageTask{cargoType, task})
					continue
				}

				// process transport task
				transferAmount, waitAtWaypoint := fleet.getCargoLoadAmount(dest, cargoType, task)

				// if we need to wait for any task, wait
				wp.WaitAtWaypoint = wp.WaitAtWaypoint || waitAtWaypoint

				t.fleetTransferCargo(fleet, -transferAmount, cargoType, dest)
			}

			// process dunnage tasks
			for _, dunnageTask := range dunnageTasks {
				cargoType, task := dunnageTask.cargoType, dunnageTask.task

				transferAmount, waitAtWaypoint := fleet.getCargoLoadAmount(dest, cargoType, task)

				// if we need to wait for any task, wait
				wp.WaitAtWaypoint = wp.WaitAtWaypoint || waitAtWaypoint

				t.fleetTransferCargo(fleet, -transferAmount, cargoType, dest)
			}

		}
	}
}

// fleetTransferCargo transfers cargo from a fleet to a cargo holder
// this will send a player a message if they are not allowed to load from this cargoholder
// this will trigger an invasion if a player unloads colonists onto a planet
func (t *turn) fleetTransferCargo(fleet *Fleet, transferAmount int, cargoType CargoType, dest cargoHolder) {
	if transferAmount != 0 {
		player := t.game.Players[fleet.PlayerNum-1]
		planet, ok := dest.(*Planet)
		if transferAmount > 0 && cargoType == Colonists && ok && planet.owned() && !planet.OwnedBy(fleet.PlayerNum) {
			// invasion!
			invadePlanet(planet, fleet, t.game.Players[planet.PlayerNum-1], t.game.Players[fleet.PlayerNum-1], transferAmount*100, t.game.Rules.InvasionDefenseCoverageFactor)
		} else if transferAmount < 0 && !dest.canLoad(fleet.PlayerNum) {
			// can't load from things we don't own
			messager.fleetInvalidLoadCargo(player, fleet, dest, cargoType, transferAmount)
		} else {
			fleet.transferToDest(dest, cargoType, transferAmount)
			fleet.MarkDirty()
			dest.MarkDirty()
			messager.fleetTransportedCargo(player, fleet, dest, cargoType, transferAmount)
		}
	}
}

func (t *turn) fleetMerge0() {

}

func (t *turn) fleetRoute() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskRoute {
			player := t.game.Players[fleet.PlayerNum-1]
			planet := t.game.getOrbitingPlanet(fleet)
			if planet == nil {
				messager.fleetInvalidRouteNotPlanet(player, fleet)
			} else {
				if !player.IsFriend(planet.PlayerNum) {
					messager.fleetInvalidRouteNotFriendlyPlanet(player, fleet, planet)
				} else if planet.TargetType == MapObjectTypeNone || planet.TargetNum == 0 {
					messager.fleetInvalidRouteNoRouteTarget(player, fleet, planet)
				} else {
					mo := t.game.getMapObject(planet.TargetType, planet.TargetNum, planet.TargetPlayerNum)
					if mo == nil {
						messager.fleetInvalidRouteNoRouteTarget(player, fleet, planet)
						continue
					}

					// insert a new waypoint after this one and route to it
					if len(fleet.Waypoints) > 1 {
						fleet.Waypoints = append(fleet.Waypoints[:1], fleet.Waypoints[1:]...)
					} else {
						fleet.Waypoints = append(fleet.Waypoints, Waypoint{})
					}
					fleet.Waypoints[1] = Waypoint{
						Position:        mo.Position,
						TargetType:      planet.TargetType,
						TargetNum:       planet.TargetNum,
						TargetPlayerNum: planet.TargetPlayerNum,
						WarpFactor:      wp.WarpFactor,
					}

					// if the new target is a planet and it has a target, keep routing
					if mo.Type == MapObjectTypePlanet {
						targetPlanet := t.game.getPlanet(mo.Num)
						if targetPlanet.TargetNum != 0 && targetPlanet.TargetType != MapObjectTypeNone {
							fleet.Waypoints[1].Task = WaypointTaskRoute
						}
					}

					messager.fleetRouted(player, fleet, planet, mo.Name)
				}
			}
		}
	}
}

// mark all wp0 as processed so they won't be processed again during wp1 steps
func (t *turn) fleetMarkWaypointsProcessed() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]
		wp.processed = true
	}
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
				wp0 := fleet.Waypoints[0]
				wp1 := fleet.Waypoints[1]

				if wp1.WarpFactor == StargateWarpFactor {
					// yeah, gate!
					fleet.gateFleet(&t.game.Rules, t.game.Universe, t.game)
				} else {
					fleet.moveFleet(&t.game.Rules, t.game.Universe, t.game)
				}

				// make sure we have tokens left after move
				if len(fleet.Tokens) == 0 {
					t.game.deleteFleet(fleet)
					continue
				}

				// remove the previous waypoint, it's been processed already
				if fleet.RepeatOrders && !wp0.PartiallyComplete {
					// if we are supposed to repeat orders,
					wp0.processed = false
					wp0.WaitAtWaypoint = false
					wp0.PartiallyComplete = false
					fleet.Waypoints = append(fleet.Waypoints, wp0)
				}

				// update the game dictionaries with this fleet's new position
				t.game.moveFleet(fleet, originalPosition)
			} else {
				fleet.PreviousPosition = &originalPosition
				fleet.WarpSpeed = 0
				fleet.Heading = Vector{}
			}
		}

	}
}

func (t *turn) fleetReproduce() {
	for _, fleet := range t.game.Fleets {
		if fleet.Cargo.Colonists == 0 {
			continue
		}

		// check if this player's freighters reproduce
		player := t.game.getPlayer(fleet.PlayerNum)
		if player.Race.Spec.FreighterGrowthFactor == 0 {
			continue
		}

		// load the orbiting planet
		planet := t.game.getOrbitingPlanet(fleet)

		growthFactor := player.Race.Spec.FreighterGrowthFactor
		growth := maxInt(1, int(growthFactor*float64(player.Race.GrowthRate)/100.0*float64(fleet.Cargo.Colonists)))
		fleet.Cargo.Colonists = fleet.Cargo.Colonists + growth
		fleet.MarkDirty()
		over := maxInt(0, fleet.Cargo.Total()-fleet.Spec.CargoCapacity)
		if over > 0 {
			// remove excess colonists
			fleet.Cargo.Colonists = fleet.Cargo.Colonists - over
			if planet != nil && planet.OwnedBy(fleet.PlayerNum) {
				// add colonists to the planet this fleet is orbiting
				planet.Cargo.Colonists = planet.Cargo.Colonists + over
				planet.MarkDirty()
			}
		}

		// Message the player
		messager.fleetReproduce(player, fleet, growth*100, planet, over)

	}
}

// decay each salvage and remove it from the universe if it's empty
func (t *turn) decaySalvage() {
	for _, salvage := range t.game.Salvages {
		salvage.decay(t.game.rules)

		if (salvage.Cargo == Cargo{}) {
			t.game.deleteSalvage(salvage)
		}
	}
}

func (t *turn) decayPackets() {

}

// jiggle, degrade, and jump wormholes
func (t *turn) wormholeJiggle() {
	if len(t.game.Wormholes) == 0 {
		return
	}

	planetPositions := make([]Vector, len(t.game.Planets))
	wormholePositions := make([]Vector, len(t.game.Wormholes))

	latestNum := t.game.Wormholes[len(t.game.Wormholes)-1].Num
	for _, wormhole := range t.game.Wormholes {
		originalPosition := wormhole.Position
		wormhole.jiggle(t.game.Universe, t.game.Rules.random)
		t.game.moveWormhole(wormhole, originalPosition)

		wormhole.degrade()
		if wormhole.shouldJump(t.game.Rules.random) {
			// this wormhole jumped. We actually delete the previous one and create a new one. This way scanner history is reset
			position, _, err := generateWormhole(t.game.Universe, latestNum+1, t.game.Area, t.game.Rules.random, planetPositions, wormholePositions, t.game.Rules.WormholeMinPlanetDistance)
			if err != nil {
				// don't kill turn generation over this, just move on without a new wormhole
				log.Error().Err(err).Msgf("failed to generate new wormhole after wormhole jump")
				continue
			}

			// create the new wormhole
			companion := t.game.Universe.getWormhole(wormhole.DestinationNum)
			newWormhole := t.game.createWormhole(position, WormholeStabilityRockSolid, companion)

			// queue the old wormhole for deletion and add the new wormhole to the universe
			t.game.deleteWormhole(wormhole)

			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Msgf("generated new wormhole %d (%v) after jump", newWormhole.Num, newWormhole.Position)

		}

		// update the spec
		wormhole.Spec = computeWormholeSpec(wormhole, &t.game.Rules)
	}
}

// SD races can detonate a minefield
func (t *turn) detonateMines() {
	for _, mineField := range t.game.MineFields {
		if !mineField.Detonate {
			continue
		}

		stats := t.game.rules.MineFieldStatsByType[mineField.MineFieldType]
		if !stats.CanDetonate {
			continue
		}

		player := t.game.getPlayer(mineField.PlayerNum)
		fleetsWithin := t.game.fleetsWithin(mineField.Position, mineField.Spec.Radius)
		for _, fleet := range fleetsWithin {
			fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
			mineField.damageFleet(player, fleet, fleetPlayer, stats)
		}
	}
}

// mine all owned planets for minerals
func (t *turn) planetMine() {
	for _, planet := range t.game.Planets {
		if planet.owned() {
			planet.Cargo = planet.Cargo.AddMineral(planet.Spec.MiningOutput)
			planet.MineYears.AddInt(planet.Mines)
			planet.reduceMineralConcentration(&t.game.Rules)
		}
	}
}

// for AR races, remote mine their own planets during this phase
func (t *turn) fleetRemoteMineAR() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskRemoteMining {
			player := t.game.getPlayer(fleet.PlayerNum)
			planet := t.game.getOrbitingPlanet(fleet)

			// can't remote mine deep space
			if planet == nil {
				messager.remoteMineDeepSpace(player, fleet)
				continue
			}

			// we can remote mine our own planets, so remote mine this planet now (it happens earlier than normal remote mining)
			if planet.OwnedBy(fleet.PlayerNum) && player.Race.Spec.CanRemoteMineOwnPlanets {
				t.remoteMine(fleet, player, planet)
			}
		}
	}

}

// remote mine planets
func (t *turn) fleetRemoteMine() {

	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskRemoteMining {
			player := t.game.getPlayer(fleet.PlayerNum)
			planet := t.game.getOrbitingPlanet(fleet)

			// can't remote mine deep space
			if planet == nil {
				messager.remoteMineDeepSpace(player, fleet)
				continue
			}

			// we can remote mine our own planets, but that happens at an earlier step, so skip  during normal remote mining
			if planet.OwnedBy(fleet.PlayerNum) && player.Race.Spec.CanRemoteMineOwnPlanets {
				continue
			}

			if planet.owned() {
				messager.remoteMineInhabited(player, fleet, planet)
				continue
			}

			t.remoteMine(fleet, player, planet)
		}
	}
}

// remote mine a planet
func (t *turn) remoteMine(fleet *Fleet, player *Player, planet *Planet) {

	if fleet.Spec.MiningRate == 0 {
		messager.remoteMineNoMiners(player, fleet, planet)
		return
	}

	// don't mine if we moved here this round, otherwise mine
	if fleet.PreviousPosition == nil || *fleet.PreviousPosition == fleet.Position {
		numMines := fleet.Spec.MiningRate
		mineralOutput := planet.getMineralOutput(numMines, t.game.rules.RemoteMiningMineOutput)
		planet.Cargo = planet.Cargo.AddMineral(mineralOutput)
		planet.MineYears.AddInt(numMines)
		planet.reduceMineralConcentration(&t.game.Rules)
		planet.MarkDirty()

		// make sure we know about this planet's cargo after remote mining
		d := newDiscoverer(player)
		d.discoverPlanetCargo(player, planet)
		messager.remoteMined(player, fleet, planet, mineralOutput)
	}
}

func (t *turn) planetProduction() {
	for _, planet := range t.game.Planets {
		if planet.owned() && len(planet.ProductionQueue) > 0 {
			player := t.game.Players[planet.PlayerNum-1]
			producer := newProducer(planet, player)
			result := producer.produce()
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

	playerFleets := t.game.getFleets(player.Num)
	fleetNum := player.getNextFleetNum(playerFleets)
	fleet := newFleetForToken(player, fleetNum, token, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, token.design.Spec.Engine.IdealSpeed)})
	fleet.Position = planet.Position
	fleet.Spec = ComputeFleetSpec(&t.game.Rules, player, &fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	fleet.OrbitingPlanetNum = planet.Num

	t.game.Fleets = append(t.game.Fleets, &fleet)
	return fleet
}

func (t *turn) playerResearch() {
	r := NewResearcher(&t.game.Rules)

	// figure out how much each player can spend on research this turn
	resourcesToSpendByPlayer := make(map[int]int, len(t.game.Players))
	for _, planet := range t.game.Planets {
		if planet.owned() {
			resourcesToSpendByPlayer[planet.PlayerNum-1] += planet.Spec.ResourcesPerYearResearch
		}
	}

	// create a map of player num to player who gained a level
	playerGainedLevel := make(map[int]bool, len(t.game.Players))

	onLevelGained := func(player *Player, field TechField) {
		messager.techLevel(player, field, player.TechLevels.Get(field), player.Researching)
		playerGainedLevel[player.Num] = true
	}

	// keep track of how many research resources are stealable by other players
	stealableResearchResources := TechLevel{}
	for _, player := range t.game.Players {
		primaryField := player.Researching
		resourcesToSpend := int(float64(resourcesToSpendByPlayer[player.Num-1])*player.Race.Spec.ResearchFactor + .5)
		player.ResearchSpentLastYear += resourcesToSpend

		// research tech levels until the resources run out
		spent := r.research(player, resourcesToSpend, onLevelGained)
		stealableResearchResources.Add(spent)

		// some races research other techs in addition to their primary field
		if player.Race.Spec.ResearchSplashDamage > 0 {
			resourcesToSpendOnOtherFields := int(float64(resourcesToSpend)*player.Race.Spec.ResearchSplashDamage + .5)
			for _, field := range TechFields {
				if field != primaryField {
					r.researchField(player, field, resourcesToSpendOnOtherFields, onLevelGained)
					stealableResearchResources.Set(field, stealableResearchResources.Get(field)+resourcesToSpendOnOtherFields)
					player.ResearchSpentLastYear += resourcesToSpendOnOtherFields
				}
			}
		}
	}

	for _, player := range t.game.Players {
		// find out if this player should steal any percentage of this research
		stealsResearch := player.Race.Spec.StealsResearch
		stolenResearch := TechLevel{
			Energy:        int(float64(stealableResearchResources.Energy) * stealsResearch.Energy),
			Weapons:       int(float64(stealableResearchResources.Weapons) * stealsResearch.Weapons),
			Propulsion:    int(float64(stealableResearchResources.Propulsion) * stealsResearch.Propulsion),
			Construction:  int(float64(stealableResearchResources.Construction) * stealsResearch.Construction),
			Electronics:   int(float64(stealableResearchResources.Electronics) * stealsResearch.Electronics),
			Biotechnology: int(float64(stealableResearchResources.Biotechnology) * stealsResearch.Biotechnology),
		}

		// we have stolen research! yay!
		// we steal the average of each research
		if stolenResearch.Sum() > 0 {
			for _, field := range TechFields {
				stolenResourcesForField := stolenResearch.Get(field) / len(t.game.Players)
				r.researchField(player, field, stolenResourcesForField, onLevelGained)
			}
		}
	}

	// update player and design specs for players who gained a level
	for _, player := range t.game.Players {
		if !playerGainedLevel[player.Num] {
			continue
		}

		// update player spec for players who gained a level
		player.Spec = computePlayerSpec(player, &t.game.Rules, t.game.Planets)

		// update design spec
		for i := range player.Designs {
			design := player.Designs[i]

			// store the numBuilt/numInstances because they spec resets them to 0
			numBuilt := design.Spec.NumBuilt
			numInstances := design.Spec.NumInstances
			design.Spec = ComputeShipDesignSpec(&t.game.Rules, player.TechLevels, player.Race.Spec, design)
			design.Spec.NumBuilt = numBuilt
			design.Spec.NumInstances = numInstances
		}
	}

	// update planet specs for players who gained a level
	for _, planet := range t.game.Planets {
		if !playerGainedLevel[planet.PlayerNum] {
			continue
		}
		planet.Spec = computePlanetSpec(&t.game.Rules, t.game.Players[planet.PlayerNum-1], planet)
	}

	// update fleet specs for players who gained a level
	for _, fleet := range t.game.Fleets {
		if !playerGainedLevel[fleet.PlayerNum] {
			continue
		}
		fleet.Spec = ComputeFleetSpec(&t.game.Rules, t.game.Players[fleet.PlayerNum-1], fleet)
	}

}

// for each planet, randomly check if the owner permaforms it
func (t *turn) permaform() {

	terraformer := NewTerraformer()

	for _, planet := range t.game.Planets {
		if planet.owned() {
			player := t.game.Players[planet.PlayerNum-1]
			adjustedPermaformChance := player.Race.Spec.PermaformChance
			if planet.population() <= player.Race.Spec.PermaformPopulation {
				adjustedPermaformChance *= 1.0 - float64(player.Race.Spec.PermaformPopulation-planet.population())/float64(player.Race.Spec.PermaformPopulation)
			}

			if adjustedPermaformChance >= t.game.Rules.random.Float64() {
				habType := HabTypes[t.game.Rules.random.Intn(len(HabTypes))]
				result := terraformer.PermaformOneStep(planet, player, habType)

				if result.Terraformed() {
					planet.Spec = computePlanetSpec(&t.game.Rules, player, planet)
					messager.permaform(player, planet, result.Type, result.Direction)
				}
			}
		}
	}

}

// grow all owned planets by some population
func (t *turn) planetGrow() {
	for _, planet := range t.game.Planets {
		if planet.owned() {
			// player := t.game.Players[*planet.PlayerNum - 1]
			planet.setPopulation(planet.population() + planet.Spec.GrowthAmount)
			planet.MarkDirty() // flag for update
		}
	}
}

func (t *turn) packetMove1() {

}

// refuel fleets if they are orbiting a planet with a friendly starbase
func (t *turn) fleetRefuel() {
	for _, fleet := range t.game.Fleets {
		planet := t.game.getOrbitingPlanet(fleet)
		if planet == nil {
			continue
		}
		if !planet.Spec.HasStarbase {
			continue
		}

		planetPlayer := t.game.getPlayer(planet.PlayerNum)
		if planetPlayer.IsFriend(fleet.PlayerNum) {
			fleet.Fuel = fleet.Spec.FuelCapacity
			fleet.MarkDirty()
		}

	}
}

func (t *turn) randomCometStrike() {

}

func (t *turn) randomMineralDeposit() {

}

func (t *turn) randomPlanetaryChange() {

}

func (t *turn) fleetBattle() {
	battleNum := 1
	for _, mos := range t.game.mapObjectsByPosition {
		playersAtPosition := map[int]*Player{}
		fleets := make([]*Fleet, 0, len(mos))
		var planet *Planet

		// add all starbases and fleets at this location
		for _, mo := range mos {
			if fleet, ok := mo.(*Fleet); ok {
				fleets = append(fleets, fleet)
				playersAtPosition[fleet.PlayerNum] = t.game.getPlayer(fleet.PlayerNum)
			}
			if p, ok := mo.(*Planet); ok && p.starbase != nil {
				fleets = append(fleets, p.starbase)
				planet = p
				playersAtPosition[p.PlayerNum] = t.game.getPlayer(planet.PlayerNum)
			}
		}

		if len(playersAtPosition) <= 1 {
			// not more than one player, no battle
			continue
		}

		battler := newBattler(t.game.rules, t.game.rules.techs, battleNum, playersAtPosition, fleets, planet)

		if battler.hasTargets() {
			// someone wants to fight, run the battle!
			record := battler.runBattle()

			// every player should discover all designs in a battle as if they were penscanned.
			discoverersByPlayer := make(map[int]discoverer, len(t.game.Players))
			designsToDiscover := map[playerObject]*ShipDesign{}
			for _, player := range playersAtPosition {
				discoverer := newDiscoverer(player)
				// discover other players at the battle
				for _, otherplayer := range playersAtPosition {
					discoverer.discoverPlayer(otherplayer)
				}

				// discover parts of this planet's starbase
				if planet != nil && planet.starbase != nil && planet.PlayerNum != player.Num {
					discoverer.discoverPlanetStarbase(player, planet)
				}

				// store this discoverer for discovering designs
				discoverersByPlayer[player.Num] = discoverer

				// player knows about the battle
				player.Battles = append(player.Battles, *record)
				messager.battle(player, planet, record)
			}
			for _, fleet := range fleets {
				updatedTokens := make([]ShipToken, 0, len(fleet.Tokens))
				for _, token := range fleet.Tokens {
					// add this design to our set of designs that should be discovered
					designsToDiscover[playerObjectKey(fleet.PlayerNum, token.DesignNum)] = token.design
					// keep this token
					if token.Quantity > 0 {
						updatedTokens = append(updatedTokens, token)
					}
				}
				fleet.Tokens = updatedTokens

				if len(fleet.Tokens) == 0 {
					t.game.deleteFleet(fleet)
				}
			}

			// discover all enemy designs
			for _, design := range designsToDiscover {
				for _, player := range playersAtPosition {
					if player.Num != design.PlayerNum {
						d := discoverersByPlayer[player.Num]
						d.discoverDesign(player, design, true)
					}
				}

			}
			battleNum++
		}

	}
}

func (t *turn) fleetBomb() {
	bomber := NewBomber(&t.game.Rules)
	for _, planet := range t.game.Planets {
		if !planet.owned() || planet.population() == 0 || planet.Spec.HasStarbase {
			// can't bomb uninhabited planets, planets with starbases
			continue
		}
		player := t.game.getPlayer(planet.PlayerNum)
		if player.Race.Spec.LivesOnStarbases {
			// can't bomb planets that are owned by AR races
			continue
		}

		// find any enemy bombers orbiting this planet
		enemyBombers := []*Fleet{}
		for _, mo := range t.game.getMapObjectsAtPosition(planet.Position) {
			if fleet, ok := mo.(*Fleet); ok && fleet.Spec.Bomber && t.game.getPlayer(fleet.PlayerNum).IsEnemy(planet.PlayerNum) {
				enemyBombers = append(enemyBombers, fleet)
			}
		}

		if len(enemyBombers) > 0 {
			// see if this planet has enemy bomber fleets, and if so, bomb it
			bomber.bombPlanet(planet, player, enemyBombers, t.game)
		}
	}
}

func (t *turn) mysteryTraderMeet() {

}

// decay MineFields and remove any minefields that are too small
func (t *turn) decayMines() {
	for _, mineField := range t.game.MineFields {
		player := t.game.getPlayer(mineField.PlayerNum)
		mineField.NumMines -= mineField.Spec.DecayRate
		if mineField.NumMines <= 10 {
			t.game.deleteMineField(mineField)
			continue
		}
		mineField.Spec = computeMinefieldSpec(t.game.rules, player, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
		mineField.MarkDirty()
	}
}

func (t *turn) fleetLayMines() {
	for _, fleet := range t.game.Fleets {
		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskLayMineField {
			player := t.game.getPlayer(fleet.PlayerNum)

			if !fleet.Spec.CanLayMines {
				messager.fleetMinesLaidFailed(player, fleet)
				continue
			}

			for mineType, minesLaid := range fleet.Spec.MineLayingRateByMineType {
				if len(fleet.Waypoints) > 1 {
					minesLaid = int(float64(minesLaid) * player.Race.Spec.MineFieldRateMoveFactor)
				}

				// We aren't laying mines (probably because we're moving, skip it)
				if minesLaid == 0 {
					continue
				}

				// See if we are adding to an existing minefield
				mineField := t.game.getMineFieldNearPosition(player.Num, fleet.Position, mineType)
				if mineField == nil {
					mineField = newMineField(player, mineType, minesLaid, t.game.getNextMineFieldNum(), fleet.Position)
					t.game.addMineField(mineField)
				} else {
					// Add to it!
					mineField.NumMines += minesLaid
				}

				messager.fleetMinesLaid(player, fleet, mineField, minesLaid)

				if mineField.Position != fleet.Position {
					// Move this minefield closer to us (in case it's not in our location)
					// This was taken from the FreeStars codebase (like many other things)
					mineField.Position = Vector{
						X: float64(minesLaid)/float64(mineField.NumMines)*(float64(fleet.Position.X)-float64(mineField.Position.X)) + mineField.Position.X,
						Y: float64(minesLaid)/float64(mineField.NumMines)*(float64(fleet.Position.Y)-float64(mineField.Position.Y)) + mineField.Position.Y,
					}
				}

				// TODO (performance): the radius will be computed in the spec as well. hmmmm
				mineField.Spec = computeMinefieldSpec(t.game.rules, player, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
				mineField.MarkDirty()
			}
		}
	}

}

func (t *turn) fleetTransferOwner() {

}

func (t *turn) fleetMerge1() {

}

func (t *turn) instaform() {
	for _, planet := range t.game.Planets {
		if planet.owned() {
			player := t.game.getPlayer(planet.PlayerNum)
			if player.Race.Spec.Instaforming {
				terraformAmount := planet.Spec.TerraformAmount
				if terraformAmount.absSum() > 0 {
					// Instantly terraform this planet (but don't update planet.TerraformAmount, this change doesn't stick if we leave)
					planet.Hab = planet.Hab.Add(terraformAmount)
					planet.Spec = computePlanetSpec(&t.game.Rules, player, planet)
					messager.instaform(player, planet, terraformAmount)
				}
			}
		}
	}
}

func (t *turn) fleetSweepMines() {

	// fleets and starbases sweep
	for _, fleet := range append(t.game.Fleets, t.game.Starbases...) {
		if !fleet.Delete && fleet.Spec.MineSweep > 0 {
			fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
			for _, mineField := range t.game.MineFields {
				// don't sweep dead fields
				if mineField.Delete {
					continue
				}
				willSweep := false
				switch fleet.battlePlan.AttackWho {
				case BattleAttackWhoEnemies:
					willSweep = fleetPlayer.IsEnemy(mineField.PlayerNum)
				case BattleAttackWhoEnemiesAndNeutrals:
					willSweep = fleetPlayer.IsEnemy(mineField.PlayerNum) || fleetPlayer.IsNeutral(mineField.PlayerNum)
				case BattleAttackWhoEveryone:
					willSweep = true
				}

				// sweep mines
				if willSweep && isPointInCircle(fleet.Position, mineField.Position, mineField.Radius()) {
					mineFieldPlayer := t.game.getPlayer(mineField.PlayerNum)
					mineField.sweep(t.game.rules, fleet, fleetPlayer, mineFieldPlayer)

					if mineField.NumMines <= 10 {
						t.game.deleteMineField(mineField)
						continue
					}

					mineField.Spec.Radius = mineField.Radius()
					mineField.MarkDirty()
				}
			}
		}
	}

	// compute specs for any dirty minefields
	// computing minefield specs is intensive because we have to count planets
	for _, mineField := range t.game.MineFields {
		if mineField.Dirty && !mineField.Delete {
			mineFieldPlayer := t.game.getPlayer(mineField.PlayerNum)
			mineField.Spec = computeMinefieldSpec(t.game.rules, mineFieldPlayer, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
		}
	}
}

// repair fleets and starbases
func (t *turn) fleetRepair() {
	for _, fleet := range t.game.Fleets {
		if !fleet.Delete {
			player := t.game.getPlayer(fleet.PlayerNum)
			orbiting := t.game.getOrbitingPlanet(fleet)
			fleet.repairFleet(t.game.rules, player, orbiting)
		}
	}

	for _, starbase := range t.game.Starbases {
		if !starbase.Delete && starbase.Tokens[0].QuantityDamaged > 0 {
			player := t.game.getPlayer(starbase.PlayerNum)
			starbase.repairStarbase(t.game.rules, player)
		}
	}
}

func (t *turn) remoteTerraform() {

}

// Update a player's information about other players
// If a player scanned another player's fleet or planet, they discover the player's
// race name
func (t *turn) discoverPlayer(player *Player) {

}

func (t *turn) fleetPatrol(player *Player) {

}

func (t *turn) calculateScore(player *Player) {

}

func (t *turn) checkVictory(player *Player) {

}
