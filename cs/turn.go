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
	t.fleetUnload()
	t.fleetColonize() // colonize wp1 after arriving at a planet
	t.fleetScrap()
	t.fleetLoad()
	t.decayMines()
	t.fleetLayMines()
	t.fleetTransferOwner()
	t.fleetMerge1()
	t.fleetRoute()
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
	player := t.game.Players[fleet.PlayerNum-1]

	var planet *Planet
	if fleet.OrbitingPlanetNum != None {
		planet = t.game.getPlanet(fleet.OrbitingPlanetNum)
	}

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
		wp := fleet.Waypoints[0]
		if !wp.processed && wp.Task == WaypointTaskColonize {
			wp.processed = true

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
		wp := fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			wp.processed = true
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
		wp := fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			wp.processed = true
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
		planet := dest.(*Planet)
		if transferAmount > 0 && cargoType == Colonists && planet != nil && planet.owned() && !planet.OwnedBy(fleet.PlayerNum) {
			// invasion!
			invadePlanet(planet, fleet, t.game.Players[planet.PlayerNum-1], t.game.Players[fleet.PlayerNum-1], transferAmount*100, t.game.Rules.InvasionDefenseCoverageFactor)
		} else if transferAmount < 0 && dest.canLoad(fleet.PlayerNum) {
			// can't load from things we don't own
			messager.fleetInvalidLoadCargo(player, fleet, dest, cargoType, transferAmount)
		} else {
			fleet.transferToDest(dest, cargoType, transferAmount)
			messager.fleetTransportedCargo(player, fleet, dest, cargoType, transferAmount)
		}
	}
}

func (t *turn) fleetMerge0() {

}

func (t *turn) fleetRoute() {
	for _, fleet := range t.game.Fleets {
		wp := fleet.Waypoints[0]
		if wp.Task == WaypointTaskRoute {
			player := t.game.Players[fleet.PlayerNum-1]
			if fleet.OrbitingPlanetNum == None {
				messager.fleetInvalidRouteNotPlanet(player, fleet)
			} else {
				planet := t.game.Planets[fleet.OrbitingPlanetNum-1]
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
				player := t.game.Players[fleet.PlayerNum-1]
				wp0 := fleet.Waypoints[0]
				wp1 := fleet.Waypoints[1]

				if wp1.WarpFactor == StargateWarpFactor {
					// yeah, gate!
					fleet.gateFleet(t.game.Universe, t.game, &t.game.Rules, player)
				} else {
					fleet.moveFleet(t.game.Universe, &t.game.Rules, player)
				}

				// remove the previous waypoint, it's been processed already
				if fleet.RepeatOrders && !wp0.PartiallyComplete {
					// if we are supposed to repeat orders,
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

}

func (t *turn) decaySalvage() {

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

func (t *turn) detonateMines() {
}

// mine all owned planets for minerals
func (t *turn) planetMine() {
	for _, planet := range t.game.Planets {
		if planet.owned() {
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

	fleetNum := t.game.getNextFleetNum(player.Num)
	fleet := newFleetForToken(player, fleetNum, token, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, token.design.Spec.IdealSpeed)})
	fleet.Position = planet.Position
	fleet.BattlePlanName = player.BattlePlans[0].Name
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

func (t *turn) researchStealer() {

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

		battler := newBattler(t.game.rules, t.game.rules.techs, playersAtPosition, fleets, planet)

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

				// store this discoverer for discovering designs
				discoverersByPlayer[player.Num] = discoverer

				// player knows about the battle
				player.Battles = append(player.Battles, *record)
				messager.battle(player, planet, record)
			}
			for _, fleet := range fleets {
				for _, token := range fleet.Tokens {
					// add this design to our set of designs that should be discovered
					designsToDiscover[playerObjectKey(fleet.PlayerNum, token.DesignNum)] = token.design
					// TODO: remove tokens and fleets after battle
					if token.Quantity == 0 {

					}
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

func (t *turn) fleetRemoteMine0() {

}

func (t *turn) decayMines() {

}

func (t *turn) fleetLayMines() {

}

func (t *turn) fleetTransferOwner() {

}

func (t *turn) fleetMerge1() {

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
func (t *turn) discoverPlayer(player *Player) {

}

func (t *turn) fleetPatrol(player *Player) {

}

func (t *turn) calculateScore(player *Player) {

}

func (t *turn) checkVictory(player *Player) {

}
