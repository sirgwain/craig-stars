package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
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
		player.BattleRecords = []BattleRecord{}
		player.leftoverResources = 0
		player.Race.Spec = computeRaceSpec(&player.Race, &t.game.Rules) // this isn't necessary when the game is complete, but if I add features (like scanner cost) this will pick it up
		player.Spec = computePlayerSpec(player, &t.game.Rules, t.game.Planets)
	}

	t.computeSpecs()

	// wp0 tasks
	t.fleetInit()
	t.fleetScrap()
	t.fleetUnload()
	t.fleetColonize()
	t.fleetLoad()
	t.fleetMerge()
	t.fleetRoute()
	t.fleetMarkWaypointsProcessed()

	// move stuff through space
	t.packetInit()
	t.packetMove(false)
	t.mysteryTraderMove()
	t.fleetMove()
	t.fleetDieoff()
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
	t.packetMove(true) // move packets built this turn
	t.fleetRefuel()    // refuel after production so fleets will refuel at planets that just built a starbase this turn
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
	t.fleetMerge()
	t.fleetRoute()
	t.fleetNotifyIdle()

	// do some final stuff like instaforming and repairing
	t.instaform()
	t.fleetSweepMines()
	t.fleetRepair()
	t.fleetRemoteTerraform()

	// reset all players
	// and do player specific things like scanning
	// and patrol orders
	t.computeSpecs()           // make sure our specs are up to date
	t.game.updateTokenCounts() // update token counts
	for _, player := range t.game.Players {
		player.Spec = computePlayerSpec(player, &t.game.Rules, t.game.Planets)

		scanner := newPlayerScanner(t.game.Universe, t.game.Players, &t.game.Rules, player)
		if err := scanner.scan(); err != nil {
			return fmt.Errorf("scan universe and update player intel -> %w", err)
		}
		t.fleetPatrol(player)

		player.SubmittedTurn = false
	}

	// as a last turn step, calculate scores and check for victories
	t.calculateScores()
	t.checkDeath()

	t.game.State = GameStateWaitingForPlayers

	log.Info().
		Int64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year).
		Msgf("generated turn")
	return nil
}

// update all planet specs with the latest info
// useful before turn generation and after building
func (t *turn) computeSpecs() {
	t.game.computeSpecs()
}

// fleetInit will reset any fleet data before processing
func (t *turn) fleetInit() {
	for _, fleet := range t.game.Fleets {
		// age this fleet by 1 year
		fleet.Age++
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
		if fleet.Delete {
			continue
		}

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
		planet.Cargo = planet.Cargo.AddMineral(fleet.Cargo.ToMineral())
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
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskColonize {
			player := t.game.Players[fleet.PlayerNum-1]

			if wp.TargetType != MapObjectTypePlanet {
				messager.colonizeNonPlanet(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if wp.TargetNum == None {
				err := fmt.Errorf("%s attempted to colonize a planet but didn't target a planet", fleet.Name)
				log.Err(err).
					Int64("GameID", t.game.ID).
					Str("Fleet", fleet.Name)
				messager.error(player, err)
				wp.Task = WaypointTaskNone
				continue
			}

			planet := t.game.getPlanet(wp.TargetNum)
			if planet.Owned() {
				messager.colonizeOwnedPlanet(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if !fleet.Spec.Colonizer {
				messager.colonizeWithNoModule(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if fleet.Cargo.Colonists == 0 {
				messager.colonizeWithNoColonists(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", player.Num).
				Str("Planet", planet.Name).
				Str("Fleet", fleet.Name).
				Int("Colonists", fleet.Cargo.Colonists*100).
				Msgf("colonized planet")

			if fleet.Spec.OrbitalConstructionModule {
				design := player.GetLatestDesign(ShipDesignPurposeStarterColony)
				if design != nil {
					t.buildStarbase(player, planet, design)
				} else {
					log.Error().
						Int64("GameID", fleet.GameID).
						Int("Player", fleet.PlayerNum).
						Str("Fleet", fleet.Name).
						Str("Planet", planet.Name).
						Msgf("colonizer can't find Starter Colony design.")
				}
			}

			// colonize the planet and scrap the fleet
			fleet.colonizePlanet(&t.game.Rules, player, planet)
			t.scrapFleet(fleet)
			messager.planetColonized(player, planet)
		}
	}
}

// fleetUnload executes wp0/wp1 unload transport tasks for fleets
func (t *turn) fleetUnload() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest, found := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			var salvage *Salvage
			if !found {

				salvage = t.game.createSalvage(fleet.Position, fleet.PlayerNum, Cargo{})
				dest = salvage

				log.Debug().
					Int64("GameID", t.game.ID).
					Str("Name", t.game.Name).
					Int("Year", t.game.Year).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Str("Position", fleet.Position.String()).
					Msgf("created salvage")

			}

			for cargoType, task := range wp.getTransportTasks() {
				transferAmount, waitAtWaypoint := fleet.getCargoUnloadAmount(dest, cargoType, task)

				wp.WaitAtWaypoint = wp.WaitAtWaypoint || waitAtWaypoint

				t.fleetTransferCargo(fleet, transferAmount, cargoType, dest)

				log.Debug().
					Int64("GameID", t.game.ID).
					Str("Name", t.game.Name).
					Int("Year", t.game.Year).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Str("Dest", dest.getMapObject().Name).
					Int("Transfered", transferAmount).
					Str("cargoType", cargoType.String()).
					Msgf("unloaded cargo")
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
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest, ok := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			if !ok || dest.getMapObject().Delete {
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

				log.Debug().
					Int64("GameID", t.game.ID).
					Str("Name", t.game.Name).
					Int("Year", t.game.Year).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Str("Dest", dest.getMapObject().Name).
					Int("Transfered", transferAmount).
					Str("cargoType", cargoType.String()).
					Msgf("loaded cargo")

				t.fleetTransferCargo(fleet, -transferAmount, cargoType, dest)
			}

			// process dunnage tasks
			for _, dunnageTask := range dunnageTasks {
				cargoType, task := dunnageTask.cargoType, dunnageTask.task

				transferAmount, waitAtWaypoint := fleet.getCargoLoadAmount(dest, cargoType, task)

				// if we need to wait for any task, wait
				wp.WaitAtWaypoint = wp.WaitAtWaypoint || waitAtWaypoint

				log.Debug().
					Int64("GameID", t.game.ID).
					Str("Name", t.game.Name).
					Int("Year", t.game.Year).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Str("Dest", dest.getMapObject().Name).
					Int("Transfered", transferAmount).
					Str("cargoType", cargoType.String()).
					Msgf("dunnage loaded cargo")

				t.fleetTransferCargo(fleet, -transferAmount, cargoType, dest)
			}

			// delete this salvage if we emptied it
			if salvage, ok := dest.(*Salvage); ok && salvage.Cargo == (Cargo{}) {
				t.game.deleteSalvage(salvage)

				log.Debug().
					Int64("GameID", salvage.GameID).
					Int("Player", salvage.PlayerNum).
					Str("Salvage", salvage.Name).
					Msgf("deleted salvage")

			}
			// delete this packet if we emptied it
			if packet, ok := dest.(*MineralPacket); ok && packet.Cargo == (Cargo{}) {
				t.game.deletePacket(packet)

				log.Debug().
					Int64("GameID", packet.GameID).
					Int("Player", packet.PlayerNum).
					Str("Packet", packet.Name).
					Msgf("deleted salvage")

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
		if transferAmount > 0 && cargoType == Colonists && ok && planet.Owned() && !planet.OwnedBy(fleet.PlayerNum) {
			// invasion!
			attacker := player
			if !planet.Owned() || planet.population() == 0 {
				// can't invade uninhabited planets
				messager.planetInvadeEmpty(attacker, planet, fleet)
				return
			}
			defender := t.game.getPlayer(planet.PlayerNum)

			invadePlanet(planet, fleet, defender, player, transferAmount*100, t.game.Rules.InvasionDefenseCoverageFactor)
			fleet.Cargo.Colonists -= transferAmount
			fleet.MarkDirty()
			planet.MarkDirty()
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

func (t *turn) fleetMerge() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]
		if wp.processed || wp.Task != WaypointTaskMergeWithFleet {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)

		if wp.TargetType != MapObjectTypeFleet {
			messager.fleetInvalidMergeNotFleet(player, fleet)
			continue
		}

		target := t.game.getFleet(wp.TargetPlayerNum, wp.TargetNum)
		if target == nil {
			messager.fleetInvalidMergeNotFleet(player, fleet)
			continue
		}
		if target.PlayerNum != fleet.PlayerNum {
			messager.fleetInvalidMergeNotOwned(player, fleet)
			continue
		}

		orderer := NewOrderer()
		_, err := orderer.Merge(&t.game.Rules, player, []*Fleet{target, fleet})
		if err != nil {
			log.Err(err).Int64("GameID", t.game.ID).Int("PlayerNum", player.Num).Int("Num", fleet.Num).Msgf("Failed to merge %v with %v", fleet, target)
			messager.error(player, err)
			continue
		}

		messager.fleetMerged(player, fleet, target)

		log.Debug().
			Int64("GameID", t.game.ID).
			Str("Name", t.game.Name).
			Int("Year", t.game.Year).
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Target", target.Name).
			Msgf("fleet merged into target")

		// remove this fleet from the universe
		t.game.deleteFleet(fleet)
	}
}

func (t *turn) fleetRoute() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskRoute {
			player := t.game.Players[fleet.PlayerNum-1]
			planet := t.game.getOrbitingPlanet(fleet)
			if planet == nil {
				messager.fleetInvalidRouteNotPlanet(player, fleet)
			} else {
				if !player.IsFriend(planet.PlayerNum) {
					messager.fleetInvalidRouteNotFriendlyPlanet(player, fleet, planet)
				} else if planet.RouteTargetType == MapObjectTypeNone || planet.RouteTargetNum == 0 {
					messager.fleetInvalidRouteNoRouteTarget(player, fleet, planet)
				} else {
					mo := t.game.getMapObject(planet.RouteTargetType, planet.RouteTargetNum, planet.RouteTargetPlayerNum)
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
						TargetType:      planet.RouteTargetType,
						TargetNum:       planet.RouteTargetNum,
						TargetPlayerNum: planet.RouteTargetPlayerNum,
						WarpSpeed:       wp.WarpSpeed,
					}

					// if the new target is a planet and it has a target, keep routing
					if mo.Type == MapObjectTypePlanet {
						targetPlanet := t.game.getPlanet(mo.Num)
						if targetPlanet.RouteTargetNum != 0 && targetPlanet.RouteTargetType != MapObjectTypeNone {
							fleet.Waypoints[1].Task = WaypointTaskRoute
						}
					}

					messager.fleetRouted(player, fleet, planet, mo.Name)

					log.Debug().
						Int64("GameID", t.game.ID).
						Str("Name", t.game.Name).
						Int("Year", t.game.Year).
						Int("Player", fleet.PlayerNum).
						Str("Fleet", fleet.Name).
						Str("Planet", planet.Name).
						Str("Target", mo.Name).
						Msgf("fleet routed to target")

				}
			}
		}
	}
}

func (t *turn) fleetNotifyIdle() {
	// don't notify the first year
	if t.game.Year == t.game.Rules.StartingYear {
		return
	}

	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		// we were just built this turn
		if fleet.Age == 0 {
			continue
		}

		// we are moving/moved
		if len(fleet.Waypoints) > 1 || *fleet.PreviousPosition != fleet.Position {
			continue
		}

		if fleet.Waypoints[0].Task == WaypointTaskNone {
			player := t.game.getPlayer(fleet.PlayerNum)
			messager.fleetCompletedAssignedOrders(player, fleet)

			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Msgf("fleet idle")

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

// packetInit will reset any packet data before processing
func (t *turn) packetInit() {
	for _, packet := range t.game.MineralPackets {
		packet.builtThisTurn = false
	}
}

// move packets through space
// if builtThisTurn is true, this will only move packets that were built this turn (i.e. just launched)
func (t *turn) packetMove(builtThisTurn bool) {

	for _, packet := range t.game.MineralPackets {
		if packet.builtThisTurn != builtThisTurn {
			continue
		}
		player := t.game.getPlayer(packet.PlayerNum)
		planet := t.game.getPlanet(int(packet.TargetPlanetNum))
		var planetPlayer *Player
		if planet.Owned() {
			planetPlayer = t.game.getPlayer(planet.PlayerNum)
		}

		packet.movePacket(&t.game.Rules, player, planet, planetPlayer)

		log.Debug().
			Int64("GameID", t.game.ID).
			Str("Name", t.game.Name).
			Int("Year", t.game.Year).
			Int("Player", packet.PlayerNum).
			Str("Packet", packet.Name).
			Str("Position", packet.Position.String()).
			Msgf("moved packet")

	}
}

func (t *turn) mysteryTraderMove() {

}

func (t *turn) fleetMove() {

	fleetsTargetingFleets := []*Fleet{}

	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}
		if fleet.Starbase {
			continue
		}

		// remove the fleet from the list of map objects at it's current location
		originalPosition := fleet.Position

		if len(fleet.Waypoints) > 1 {
			wp1 := fleet.Waypoints[1]
			if wp1.TargetType == MapObjectTypeFleet {
				// move this after all the fleets not targeting fleets move
				fleetsTargetingFleets = append(fleetsTargetingFleets, fleet)
				continue
			}

			t.moveFleet(fleet)
		} else {
			fleet.PreviousPosition = &originalPosition
			fleet.WarpSpeed = 0
			fleet.Heading = Vector{}
		}
	}

	// move all the fleets targeting other fleets
	// TODO: build a directed graph and detect cycles and all that jazz
	for _, fleet := range fleetsTargetingFleets {
		t.moveFleet(fleet)
	}
}

func (t *turn) moveFleet(fleet *Fleet) {
	originalPosition := fleet.Position
	wp0 := fleet.Waypoints[0]
	wp1 := &fleet.Waypoints[1]
	if wp1.TargetNum != None {
		target := t.game.getMapObject(wp1.TargetType, wp1.TargetNum, wp1.TargetPlayerNum)
		if target == nil || target.Delete {
			// target went away
			wp1.TargetName = ""
			wp1.TargetNum = None
			wp1.TargetType = MapObjectTypeNone
			wp1.TargetPlayerNum = None
			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Msgf("fleet target gone, using position only")
		} else if target.Position != wp1.Position {
			// update the position
			wp1.Position = target.Position
			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Msgf("fleet target moved, updating position")
		}
	}

	if wp1.WarpSpeed == StargateWarpSpeed {
		// yeah, gate!
		fleet.gateFleet(&t.game.Rules, t.game.Universe, t.game)
	} else {
		fleet.moveFleet(&t.game.Rules, t.game.Universe, t.game)
	}

	log.Debug().
		Int64("GameID", t.game.ID).
		Str("Name", t.game.Name).
		Int("Year", t.game.Year).
		Int("Player", fleet.PlayerNum).
		Str("Fleet", fleet.Name).
		Str("Fuel", fmt.Sprintf("%d/%d", fleet.Fuel, fleet.Spec.FuelCapacity)).
		Int("WarpSpeed", fleet.WarpSpeed).
		Str("Start", wp0.Position.String()).
		Str("End", fleet.Position.String()).
		Msgf("moved fleet")

	// make sure we have tokens left after move
	if len(fleet.Tokens) == 0 {
		t.game.deleteFleet(fleet)
		return
	}

	// remove the previous waypoint, it's been processed already
	if fleet.RepeatOrders && !wp0.PartiallyComplete {
		// if we are supposed to repeat orders,
		wp0.processed = false
		wp0.WaitAtWaypoint = false
		wp0.PartiallyComplete = false
		fleet.Waypoints = append(fleet.Waypoints, wp0)

		log.Debug().
			Int64("GameID", t.game.ID).
			Str("Name", t.game.Name).
			Int("Year", t.game.Year).
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Waypoint", fmt.Sprintf("%s: %s", wp0.TargetName, wp0.Task)).
			Msgf("repeating waypoint")
	}

	// update the game dictionaries with this fleet's new position
	t.game.moveFleet(fleet, originalPosition)
}

// kill off colonists on fleets from radiation poisoning
// https://wiki.starsautohost.org/wiki/Radiating_Ramscoop
// DeathRate/Year % = int ((86 - C)/2)
// where C is the center of your Rad-Hab-Range (mR)

func (t *turn) fleetDieoff() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		if fleet.Cargo.Colonists == 0 {
			continue
		}

		// we're safe, no radiation in this fleet
		if !fleet.Spec.Radiating {
			continue
		}

		// check if this player's freighters reproduce
		player := t.game.getPlayer(fleet.PlayerNum)
		habCenter := player.Race.Spec.HabCenter
		deathRate := math.Max(0, float64(t.game.Rules.RadiatingImmune+1)-float64(habCenter.Rad)) / 2 / 100

		if deathRate > 0 {
			killed := maxInt(1, int(deathRate*float64(fleet.Cargo.Colonists)))
			fleet.Cargo.Colonists = maxInt(0, fleet.Cargo.Colonists-killed)
			fleet.MarkDirty()

			// Message the player
			messager.fleetColonistsDieoff(player, fleet, killed*100)

			log.Debug().
				Int64("GameID", t.game.ID).
				Str("Name", t.game.Name).
				Int("Year", t.game.Year).
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Msgf("fleet radiation dieoff")
		}

	}
}

func (t *turn) fleetReproduce() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

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

		log.Debug().
			Int64("GameID", t.game.ID).
			Str("Name", t.game.Name).
			Int("Year", t.game.Year).
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Int("Growth", growth).
			Int("Over", over).
			Msgf("fleet reproduced")

	}
}

// decay each salvage and remove it from the universe if it's empty
func (t *turn) decaySalvage() {
	for _, salvage := range t.game.Salvages {
		salvage.decay(&t.game.Rules)

		log.Debug().
			Int64("GameID", salvage.GameID).
			Int("Player", salvage.PlayerNum).
			Str("Salvage", salvage.Name).
			Str("Cargo", salvage.Cargo.PrettyString()).
			Msgf("decayed salvage")

		if (salvage.Cargo == Cargo{}) {
			t.game.deleteSalvage(salvage)

			log.Debug().
				Int64("GameID", salvage.GameID).
				Int("Player", salvage.PlayerNum).
				Str("Salvage", salvage.Name).
				Msgf("deleted salvage")

		}
	}
}

// Decay mineral packets in flight
func (t *turn) decayPackets() {
	for _, packet := range t.game.MineralPackets {
		player := t.game.getPlayer(packet.PlayerNum)
		// update the decay rate based on this distance traveled this turn
		decayRate := 1 - packet.getPacketDecayRate(&t.game.Rules, &player.Race)*(packet.distanceTravelled/float64(packet.WarpSpeed*packet.WarpSpeed))
		packet.Cargo = packet.Cargo.Multiply(decayRate)

		log.Debug().
			Int64("GameID", packet.GameID).
			Int("Player", packet.PlayerNum).
			Str("Packet", packet.Name).
			Str("Cargo", packet.Cargo.PrettyString()).
			Msgf("decayed packet")

		if (packet.Cargo == Cargo{}) {
			t.game.deletePacket(packet)

			log.Debug().
				Int64("GameID", packet.GameID).
				Int("Player", packet.PlayerNum).
				Str("Packet", packet.Name).
				Msgf("deleted salvage")

		}

	}
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
			newWormhole := t.game.createWormhole(&t.game.Rules, position, WormholeStabilityRockSolid, companion)

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

		stats := t.game.Rules.MineFieldStatsByType[mineField.MineFieldType]
		if !stats.CanDetonate {
			continue
		}

		player := t.game.getPlayer(mineField.PlayerNum)
		fleetsWithin := t.game.fleetsWithin(mineField.Position, mineField.Spec.Radius)
		for _, fleet := range fleetsWithin {
			fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
			mineField.damageFleet(player, fleet, fleetPlayer, stats)
		}

		log.Debug().
			Int64("GameID", t.game.ID).
			Int("Player", mineField.PlayerNum).
			Str("MineField", mineField.Name).
			Int("NumMines", mineField.NumMines).
			Msgf("detonated mineField")

	}
}

// mine all owned planets for minerals
func (t *turn) planetMine() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			planet.mine(&t.game.Rules)
			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", planet.PlayerNum).
				Str("Planet", planet.Name).
				Str("Minerals", planet.Spec.MiningOutput.PrettyString()).
				Msgf("planet mined")

		}
	}
}

// for AR races, remote mine their own planets during this phase
func (t *turn) fleetRemoteMineAR() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

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
		if fleet.Delete {
			continue
		}

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

			if planet.Owned() {
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
		mineralOutput := planet.getMineralOutput(numMines, t.game.Rules.RemoteMiningMineOutput)
		planet.Cargo = planet.Cargo.AddMineral(mineralOutput)
		planet.MineYears = planet.MineYears.AddInt(numMines)
		planet.reduceMineralConcentration(&t.game.Rules)
		planet.MarkDirty()

		// make sure we know about this planet's cargo after remote mining
		d := newDiscoverer(player)
		d.discoverPlanetCargo(player, planet)
		messager.remoteMined(player, fleet, planet, mineralOutput)

		log.Debug().
			Int64("GameID", t.game.ID).
			Int("Player", fleet.PlayerNum).
			Str("Fleet", fleet.Name).
			Str("Planet", planet.Name).
			Str("Minerals", mineralOutput.PrettyString()).
			Msgf("planet remote mined")

	}
}

// go through each player planet and process it's production queue
func (t *turn) planetProduction() {
	for _, planet := range t.game.Planets {
		if planet.Owned() && len(planet.ProductionQueue) > 0 {
			player := t.game.Players[planet.PlayerNum-1]
			planet.PopulateProductionQueueCosts(player)
			producer := newProducer(planet, player)
			result := producer.produce()

			// add any invalid messages we encountered
			if len(result.messages) > 0 {
				player.Messages = append(player.Messages, result.messages...)
			}

			// message about planetary installations
			if result.mines > 0 {
				messager.minesBuilt(player, planet, result.mines)
			}
			if result.factories > 0 {
				messager.factoriesBuilt(player, planet, result.factories)
			}
			if result.defenses > 0 {
				messager.defensesBuilt(player, planet, result.defenses)
			}

			// message about mineral alchemy
			if result.alchemy != (Mineral{}) {
				// alchemy builds evenly, but we only message the single amount
				numBuilt := result.alchemy.Ironium
				messager.mineralAlchemyBuilt(player, planet, numBuilt)
			}

			// message about each terraform step
			if len(result.terraformResults) > 0 {
				for _, terraformResult := range result.terraformResults {
					messager.terraform(player, planet, terraformResult.Type, terraformResult.Direction)
				}
			}
			for _, token := range result.tokens {
				design := token.design
				design.Spec.NumBuilt += token.Quantity
				design.Spec.NumInstances += token.Quantity
				design.MarkDirty()

				fleet := t.buildFleet(player, planet, token)
				messager.fleetBuilt(player, planet, fleet, token.Quantity)
			}
			for _, cargo := range result.packets {
				target := t.game.getPlanet(planet.PacketTargetNum)
				packet := t.buildMineralPacket(player, planet, cargo, target)
				messager.mineralPacket(player, planet, packet, target.Name)
			}
			if result.starbase != nil {
				starbase := t.buildStarbase(player, planet, result.starbase)
				messager.starbaseBuilt(player, planet, starbase)
			}
			if result.scanner {
				planet.Scanner = true
				planet.Spec = computePlanetSpec(&t.game.Rules, player, planet)
				planet.MarkDirty()
				messager.scannerBuilt(player, planet, planet.Spec.Scanner)
			}

			// log what we actually did
			for _, itemBuilt := range result.itemsBuilt {
				if itemBuilt.numBuilt > 0 {
					log.Debug().
						Int("Player", planet.PlayerNum).
						Str("Planet", planet.Name).
						Str("Item", string(itemBuilt.queueItemType)).
						Int("DesignNum", itemBuilt.designNum).
						Int("NumBuilt", itemBuilt.numBuilt).
						Msgf("built item")
				}
			}

			// any leftover resources go back to the player for research
			player.leftoverResources += result.leftoverResources
		}
	}
}

// build a fleet with some number of tokens
func (t *turn) buildFleet(player *Player, planet *Planet, token ShipToken) *Fleet {
	player.Stats.FleetsBuilt++
	player.Stats.TokensBuilt += token.Quantity

	playerFleets := t.game.getFleets(player.Num)
	fleetNum := player.getNextFleetNum(playerFleets)
	fleet := newFleetForToken(player, fleetNum, token, []Waypoint{NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, token.design.Spec.Engine.IdealSpeed)})
	fleet.Position = planet.Position
	fleet.Spec = ComputeFleetSpec(&t.game.Rules, player, &fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
	fleet.OrbitingPlanetNum = planet.Num

	log.Debug().
		Int64("GameID", t.game.ID).
		Int("Player", fleet.PlayerNum).
		Str("Planet", planet.Name).
		Str("Fleet", fleet.Name).
		Msgf("fleet built")

	t.game.Fleets = append(t.game.Fleets, &fleet)
	t.game.Universe.addFleet(&fleet)
	return &fleet
}

// build a starbase on a planet
func (t *turn) buildStarbase(player *Player, planet *Planet, design *ShipDesign) *Fleet {
	player.Stats.StarbasesBuilt++
	player.Stats.TokensBuilt++
	design.Spec.NumBuilt++

	// remove the old starbase
	if planet.Starbase != nil {
		t.game.deleteStarbase(planet.Starbase)
	}

	starbase := newStarbase(player, planet, design, design.Name)
	starbase.Spec = ComputeFleetSpec(&t.game.Rules, player, &starbase)
	planet.setStarbase(&t.game.Rules, player, &starbase)
	log.Debug().
		Int64("GameID", t.game.ID).
		Int("Player", starbase.PlayerNum).
		Str("Planet", planet.Name).
		Str("Starbase", starbase.Name).
		Msgf("starbase built")

	t.game.Starbases = append(t.game.Starbases, &starbase)
	t.game.addStarbase(&starbase)
	return &starbase
}

// build a mineral packet with cargo
func (t *turn) buildMineralPacket(player *Player, planet *Planet, cargo Cargo, target *Planet) *MineralPacket {

	playerMineralPackets := t.game.getMineralPackets(player.Num)
	num := player.getNextMineralPacketNum(playerMineralPackets)
	packet := newMineralPacket(player, num, planet.PacketSpeed, planet.Spec.SafePacketSpeed, cargo, planet.Position, target.Num)
	packet.builtThisTurn = true

	log.Debug().
		Int64("GameID", t.game.ID).
		Int("Player", packet.PlayerNum).
		Str("Planet", planet.Name).
		Str("Fleet", packet.Name).
		Msgf("mineral packet built")

	t.game.MineralPackets = append(t.game.MineralPackets, packet)
	return packet
}

func (t *turn) playerResearch() {
	r := NewResearcher(&t.game.Rules)

	// figure out how much each player can spend on research this turn
	resourcesToSpendByPlayer := make(map[int]int, len(t.game.Players))

	// start with leftover from production
	for _, player := range t.game.Players {
		resourcesToSpendByPlayer[player.Num] = player.leftoverResources
	}

	for _, planet := range t.game.Planets {
		if planet.Owned() {
			resourcesToSpendByPlayer[planet.PlayerNum] += planet.Spec.ResourcesPerYearResearch
		}
	}

	// create a map of player num to player who gained a level
	playerGainedLevel := make(map[int]bool, len(t.game.Players))

	onLevelGained := func(player *Player, field TechField) {

		messager.techLevel(player, field, player.TechLevels.Get(field), player.Researching)
		techsGained := t.game.TechStore.GetTechsJustGained(player, field)
		for _, tech := range techsGained {
			messager.techGained(player, field, tech)
		}
		playerGainedLevel[player.Num] = true

		log.Debug().
			Int64("GameID", t.game.ID).
			Int("Player", player.Num).
			Str("Field", string(field)).
			Int("Level", player.TechLevels.Get(field)).
			Msgf("player researched new tech level")

	}

	// keep track of how many research resources are stealable by other players
	stealableResearchResources := TechLevel{}
	for _, player := range t.game.Players {
		primaryField := player.Researching
		resourcesToSpend := int(float64(resourcesToSpendByPlayer[player.Num])*player.Race.Spec.ResearchFactor + .5)
		player.ResearchSpentLastYear = resourcesToSpend

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
		if fleet.Delete {
			continue
		}

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
		if planet.Owned() {
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
					planet.MarkDirty()
					messager.permaform(player, planet, result.Type, result.Direction)

					log.Debug().
						Int64("GameID", t.game.ID).
						Int("Player", player.Num).
						Int("Planet", planet.Num).
						Str("HabType", result.Type.String()).
						Msgf("player permaformed planet")

				}
			}
		}
	}

}

// grow all owned planets by some population
func (t *turn) planetGrow() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			player := t.game.getPlayer(planet.PlayerNum)
			prevPop := planet.population()
			planet.grow(player)
			planet.MarkDirty()

			// tell players about dieing colonists
			if planet.Spec.GrowthAmount < 0 {
				if planet.Spec.PopulationDensity > 1 {
					messager.planetPopulationDecreasedOvercrowding(player, planet, planet.Spec.GrowthAmount)
				} else {
					messager.planetPopulationDecreased(player, planet, prevPop, planet.population())
				}
			}

			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", planet.PlayerNum).
				Str("Planet", planet.Name).
				Int("Capacity", int(planet.Spec.PopulationDensity*100)).
				Int("PrevPopulation", prevPop).
				Int("GrowthAmount", planet.Spec.GrowthAmount).
				Int("Population", planet.population()).
				Msgf("planet grow")

			if planet.population() <= 0 {
				planet.emptyPlanet()
				messager.planetDiedOff(player, planet)

				log.Debug().
					Int64("GameID", t.game.ID).
					Int("Player", player.Num).
					Str("Planet", planet.Name).
					Msgf("planet pop died off")
			}
		}
	}
}

// refuel fleets if they are orbiting a planet with a friendly starbase
func (t *turn) fleetRefuel() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		if fleet.Fuel == fleet.Spec.FuelCapacity {
			continue
		}

		planet := t.game.getOrbitingPlanet(fleet)
		if planet == nil {
			continue
		}
		if !planet.Spec.HasStarbase {
			continue
		}

		planetPlayer := t.game.getPlayer(planet.PlayerNum)
		if planetPlayer.IsFriend(fleet.PlayerNum) {
			player := t.game.getPlayer(fleet.PlayerNum)
			fleet.Fuel = fleet.Spec.FuelCapacity
			fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
			fleet.MarkDirty()

			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", fleet.PlayerNum).
				Str("Planet", planet.Name).
				Int("PlanetPlayer", planet.PlayerNum).
				Str("Fleet", fleet.Name).
				Msgf("fleet refueled at starbase")

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
			if fleet, ok := mo.(*Fleet); ok && !fleet.Delete {
				fleets = append(fleets, fleet)
				playersAtPosition[fleet.PlayerNum] = t.game.getPlayer(fleet.PlayerNum)
			} else if p, ok := mo.(*Planet); ok {
				planet = p
			}

		}

		if len(playersAtPosition) <= 1 {
			// not more than one player, no battle
			continue
		}

		battler := newBattler(&t.game.Rules, t.game.Rules.techs, battleNum, playersAtPosition, fleets, planet)

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
				if planet != nil && planet.Starbase != nil && planet.PlayerNum != player.Num {
					discoverer.discoverPlanetStarbase(player, planet)
				}

				// store this discoverer for discovering designs
				discoverersByPlayer[player.Num] = discoverer

			}

			// figure out how much salvage this generates
			destroyedCost := Cost{}
			salvageOwner := 1
			salvageFactor := t.game.Rules.SalvageFromBattleFactor
			if salvageFactor == 0 {
				// TODO: remove this when games are up to date
				salvageFactor = .3 // upgrade old games before this rule was available
			}
			for _, token := range record.DestroyedTokens {
				destroyedCost = destroyedCost.Add(token.design.Spec.Cost.MultiplyInt(token.Quantity))
				// TODO: who owns this salvage
				salvageOwner = token.PlayerNum
			}
			salvageMinerals := destroyedCost.MultiplyFloat64(t.game.Rules.SalvageFromBattleFactor).ToMineral()

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
					if fleet.Starbase {
						// AR races live on starbases, so empty the planet
						player := t.game.getPlayer(fleet.PlayerNum)
						if player.Race.Spec.LivesOnStarbases {
							planet.emptyPlanet()
							planet.MarkDirty()
							messager.planetDiedOff(player, planet)
						}
						// remove this starbase from the planet
						t.game.deleteStarbase(fleet)
						planet.Starbase = nil
						planet.MarkDirty()
					} else {
						t.game.deleteFleet(fleet)
					}
				}

				// recompute the spec of this fleet and make sure we don't have extra fuel sitting around
				fleet.Spec = ComputeFleetSpec(&t.game.Rules, t.game.getPlayer(fleet.PlayerNum), fleet)
				fleet.reduceFuelToMax()

				// jettison cargo
				jettisoned := fleet.reduceCargoToMax()
				record.Stats.CargoLostByPlayer[fleet.PlayerNum] = record.Stats.CargoLostByPlayer[fleet.PlayerNum].Add(jettisoned)
				jettisonedMinerals := jettisoned.ToMineral()
				if jettisonedMinerals.Total() > 0 {
					salvageMinerals = salvageMinerals.Add(jettisonedMinerals)
				}
			}

			if salvageMinerals.Total() > 0 {
				if planet == nil {
					t.game.createSalvage(record.Position, salvageOwner, salvageMinerals.ToCargo())
				} else {
					planet.Cargo = planet.Cargo.AddMineral(salvageMinerals)
					planet.MarkDirty()
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

			// message each player
			for _, player := range playersAtPosition {
				// player knows about the battle
				player.BattleRecords = append(player.BattleRecords, *record)
				messager.battle(player, planet, record)
			}
			battleNum++
		}

	}
}

func (t *turn) fleetBomb() {
	bomber := NewBomber(&t.game.Rules)
	for _, planet := range t.game.Planets {
		if !planet.Owned() || planet.population() == 0 || planet.Spec.HasStarbase {
			// can't bomb uninhabited planets, planets with starbases
			continue
		}
		planetPlayer := t.game.getPlayer(planet.PlayerNum)
		if planetPlayer.Race.Spec.LivesOnStarbases {
			// can't bomb planets that are owned by AR races
			continue
		}

		// find any enemy bombers orbiting this planet
		enemyBombers := []*Fleet{}
		for _, mo := range t.game.getMapObjectsAtPosition(planet.Position) {
			if fleet, ok := mo.(*Fleet); ok {
				fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
				willBomb := fleet.willAttack(fleetPlayer, planet.PlayerNum)
				if !fleet.Delete && fleet.Spec.Bomber && willBomb {
					enemyBombers = append(enemyBombers, fleet)
				}
			}
		}

		if len(enemyBombers) > 0 {
			// see if this planet has enemy bomber fleets, and if so, bomb it
			bomber.bombPlanet(planet, planetPlayer, enemyBombers, t.game)
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
		mineField.Spec = computeMinefieldSpec(&t.game.Rules, player, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
		mineField.MarkDirty()

		log.Debug().
			Int64("GameID", t.game.ID).
			Int("Player", mineField.PlayerNum).
			Str("MineField", mineField.Name).
			Int("NumMines", mineField.NumMines).
			Msgf("decayed mineField")

	}
}

func (t *turn) fleetLayMines() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

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
					mineField.moveTowardsMineLayer(fleet.Position, minesLaid)
				}

				// TODO (performance): the radius will be computed in the spec as well. hmmmm
				mineField.Spec = computeMinefieldSpec(&t.game.Rules, player, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
				mineField.MarkDirty()

				log.Debug().
					Int64("GameID", t.game.ID).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Int("MinesLaid", minesLaid).
					Str("MineField", mineField.Name).
					Int("NumMines", mineField.NumMines).
					Msgf("laid mines")

			}
		}
	}

}

func (t *turn) fleetTransferOwner() {

}

func (t *turn) instaform() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			player := t.game.getPlayer(planet.PlayerNum)
			if player.Race.Spec.Instaforming {
				terraformAmount := planet.Spec.TerraformAmount
				if terraformAmount.absSum() > 0 {
					// Instantly terraform this planet (but don't update planet.TerraformAmount, this change doesn't stick if we leave)
					planet.Hab = planet.Hab.Add(terraformAmount)
					planet.Spec = computePlanetSpec(&t.game.Rules, player, planet)
					messager.instaform(player, planet, terraformAmount)

					log.Debug().
						Int64("GameID", t.game.ID).
						Int("Player", player.Num).
						Str("Planet", planet.Name).
						Str("Hab", planet.Hab.String()).
						Msgf("instaformed planet")

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

				// sweep mines
				if fleet.willAttack(fleetPlayer, mineField.PlayerNum) && isPointInCircle(fleet.Position, mineField.Position, mineField.Radius()) {
					mineFieldPlayer := t.game.getPlayer(mineField.PlayerNum)
					mineField.sweep(&t.game.Rules, fleet, fleetPlayer, mineFieldPlayer)

					log.Debug().
						Int64("GameID", t.game.ID).
						Int("Player", fleet.PlayerNum).
						Str("Fleet", fleet.Name).
						Str("MineField", mineField.Name).
						Int("MineFieldPlayer", mineField.PlayerNum).
						Int("NumMines", mineField.NumMines).
						Msgf("fleet swept mines")

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
			mineField.Spec = computeMinefieldSpec(&t.game.Rules, mineFieldPlayer, mineField, t.game.Universe.numPlanetsWithin(mineField.Position, mineField.Radius()))
		}
	}
}

// repair fleets and starbases
func (t *turn) fleetRepair() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		orbiting := t.game.getOrbitingPlanet(fleet)
		fleet.repairFleet(&t.game.Rules, player, orbiting)
	}

	for _, starbase := range t.game.Starbases {
		if starbase.Delete {
			continue
		}

		if starbase.Tokens[0].QuantityDamaged == 0 {
			continue
		}

		player := t.game.getPlayer(starbase.PlayerNum)
		starbase.repairStarbase(&t.game.Rules, player)
	}
}

func (t *turn) fleetRemoteTerraform() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		// can't remote terraform with this fleet
		if fleet.Spec.TerraformRate == 0 {
			continue
		}

		// not over a planet
		if fleet.OrbitingPlanetNum == None {
			continue
		}

		// don't remote terraform an unowned planet or a planet owned by us
		planet := t.game.getPlanet(fleet.OrbitingPlanetNum)
		if !planet.Owned() || planet.OwnedBy(fleet.PlayerNum) {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		planetPlayer := t.game.getPlayer(planet.PlayerNum)
		deterraform := fleet.willAttack(player, planet.PlayerNum)
		friend := player.IsFriend(planet.PlayerNum)

		// do nothing to netural planets
		if !friend && !deterraform {
			continue
		}

		terraformer := NewTerraformer()
		for i := 0; i < fleet.Spec.TerraformRate; i++ {
			result := terraformer.TerraformOneStep(planet, planetPlayer, player, deterraform)
			if result != (TerraformResult{}) {
				log.Debug().
					Int64("GameID", t.game.ID).
					Int("Player", fleet.PlayerNum).
					Str("Fleet", fleet.Name).
					Str("Planet", planet.Name).
					Str("HabType", result.Type.String()).
					Msgf("fleet remote terraformed planet")
			}
		}
	}
}

func (t *turn) fleetPatrol(player *Player) {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete || fleet.PlayerNum != player.Num {
			continue
		}

		if len(fleet.Waypoints) != 1 {
			continue
		}

		wp := &fleet.Waypoints[0]
		if wp.Task != WaypointTaskPatrol {
			continue
		}

		distSquared := float64(wp.PatrolRange * wp.PatrolRange)
		closestDistance := float64(math.MaxFloat32)
		var closest *FleetIntel

		for _, enemyFleet := range player.FleetIntels {
			if fleet.willAttack(player, enemyFleet.PlayerNum) {
				distSquaredToFleet := fleet.Position.DistanceSquaredTo(enemyFleet.Position)
				if distSquaredToFleet <= distSquared {
					if distSquaredToFleet <= closestDistance {
						closestDistance = distSquaredToFleet
						closest = &enemyFleet
					}
				}
			}
		}

		if closest != nil {

			if wp.PatrolWarpSpeed == PatrolWarpSpeedAutomatic {
				wp.PatrolWarpSpeed = fleet.Spec.Engine.IdealSpeed
			}
			fleet.Waypoints = append(fleet.Waypoints, NewFleetWaypoint(closest.Position, closest.Num, closest.PlayerNum, closest.Name, wp.PatrolWarpSpeed))

			messager.fleetPatrolTargeted(player, fleet, closest)

			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", fleet.PlayerNum).
				Str("Fleet", fleet.Name).
				Str("Target", closest.Name).
				Int("TargetPlayer", closest.PlayerNum).
				Msgf("fleet patrol targeted enemy")

		}
	}
}

// Calculate the score for this year for each player
//
// Note: this depends on each player having updated player reports
//
// Here's how empires score:
// Planets:  From 1 to 6 points, scoring 1 point for each 100,000 colonists
// Starbases: 3 points each (doesn't include Orbital Forts)
// Unarmed Ships: An unarmed ship has a power rating of 0. You receive 1/2 point for each unarmed ship (up to the number of planets you own).
// Escort Ships: An escort ship has a power rating greater than 0 and less than 2000. You receive 2 points for each Escort ship (up to the number of planets you own).
// Capital Ships A Capital ship has a power rating of greater than 1999.  For each capital ship, you receive points calculated by the following formula:
// Points = (8 * #_capital_ships * #_planets) /( #_capital_ships + #_planets)
//
//	For example, if you have 20 capital ships and 30 planets, you receive (8 x 20 x 30) / (20 + 30) or 4.8 points for each ship.
//	      Tech Levels:  1 point for levels 1-3,
//	                    2 points for levels 4-6,
//	                    3 points for levels 7-9,
//	                    4 points for level 10 and above
//
// Resources: 1 point for every 30 resources
func (t *turn) calculateScores() {
	scores := make([]PlayerScore, len(t.game.Players))

	// Sum up planets
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			score := &scores[planet.PlayerNum-1]
			score.Planets++
			if planet.Spec.HasStarbase {
				score.Starbases++
			}
			// Planets: From 1 to 6 points, scoring 1 point for each 100,000 colonists
			score.Score += int(math.Min(float64(planet.population()/100000), 6))
			score.Resources += planet.Spec.ResourcesPerYear
		}
	}

	// Calculate ship counts
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		score := &scores[fleet.PlayerNum-1]
		for _, token := range fleet.Tokens {
			powerRating := token.design.Spec.PowerRating
			if powerRating <= 0 {
				score.UnarmedShips += token.Quantity
			} else if powerRating < 1999 {
				score.EscortShips += token.Quantity
			} else {
				score.CapitalShips += token.Quantity
			}
		}
	}

	for _, player := range t.game.Players {
		score := &scores[player.Num-1]

		// Calculate tech levels
		for _, field := range TechFields {
			level := player.TechLevels.Get(field)
			score.TechLevels += level
			switch {
			case level >= 1 && level <= 3:
				score.Score += 1
			case level >= 4 && level <= 6:
				score.Score += 2
			case level >= 7 && level <= 9:
				score.Score += 3
			case level >= 10:
				score.Score += 4
			}
		}

		// Calculate additional score components
		// Resources: 1 point for every 30 resources
		score.Score += score.Resources / 30
		// Starbases: 3 points each (doesn't include Orbital Forts)
		score.Score += score.Starbases * 3
		// Unarmed Ships: You receive 1/2 point for each unarmed ship (up to the number of planets you own).
		score.Score += int(math.Min(float64(score.UnarmedShips)*0.5+5, float64(score.Planets)))
		// Escort Ships: You receive 2 points for each Escort ship (up to the number of planets you own).
		score.Score += int(math.Min(float64(score.EscortShips)*2, float64(score.Planets)))
		// Capital Ships (8 * #_capital_ships * #_planets) /( #_capital_ships + #_planets)
		if score.CapitalShips+score.Planets > 0 {
			score.Score += int((8 * score.CapitalShips * score.Planets) / (score.CapitalShips + score.Planets))
		}

		// add this to the player's score history
		player.ScoreHistory = append(player.ScoreHistory, *score)

		// check for victory/death for this player
		t.checkVictory(player)
	}

	// sort players by score, highest to lowest
	scoreSortedPlayers := make([]*Player, len(t.game.Players))
	copy(scoreSortedPlayers, t.game.Players)
	slices.SortFunc(scoreSortedPlayers, func(p1, p2 *Player) bool {
		return p1.GetScore().Score > p2.GetScore().Score
	})

	// update rank for all scores
	rank := 1
	for i, player := range scoreSortedPlayers {
		if i > 0 {
			if scoreSortedPlayers[i-1].GetScore().Score != player.GetScore().Score {
				rank++
			}
		}

		if len(player.ScoreHistory) > 0 {
			score := &player.ScoreHistory[len(player.ScoreHistory)-1]
			score.Rank = rank
		}
	}

	// share score intel if show public scores is enabled, or if a victor has been found
	if (t.game.PublicPlayerScores && t.game.Rules.ShowPublicScoresAfterYears > 0 && t.game.YearsPassed() >= t.game.Rules.ShowPublicScoresAfterYears) || t.game.VictorDeclared {
		for _, player := range t.game.Players {
			discoverer := newDiscoverer(player)
			for _, otherPlayer := range t.game.Players {

				if player.Num == otherPlayer.Num {
					// our score is stored separately
					continue
				}
				discoverer.discoverPlayerScores(otherPlayer)
			}
		}
	}

}

// check if this player is victorious, and if so, notify everyone
func (t *turn) checkVictory(player *Player) {
	victoryChecker := newVictoryChecker(t.game)
	for _, player := range t.game.Players {
		if err := victoryChecker.checkForVictor(player); err != nil {
			log.Error().Err(err).Msg("error while checking for victory")
			return
		}
	}

	// we don't declare a victor until some time has passed
	if t.game.YearsPassed() >= t.game.VictoryConditions.YearsPassed && t.game.VictorDeclared {

		// if we won, tell everyone about it!
		if player.Victor {
			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", player.Num).
				Str("PlayerName", player.Name).
				Str("Race", player.Race.PluralName).
				Msgf("you are victorious your majesty!")

			for _, p := range t.game.Players {
				messager.victory(p, player)
			}
		}
	}
}

func (t *turn) checkDeath() {
	for _, player := range t.game.Players {

		numPlanets := 0
		numFleets := 0
		numColonists := 0
		for _, planet := range t.game.Planets {
			if planet.PlayerNum == player.Num {
				numPlanets++
				numColonists += planet.population()
			}
		}

		for _, fleet := range t.game.Fleets {
			if fleet.PlayerNum == player.Num && !fleet.Delete {
				numFleets++
				numColonists += fleet.Cargo.Colonists * 100
			}
		}

		// tell players they are dead
		if numPlanets == 0 && numFleets == 0 {
			player.Messages = append([]PlayerMessage{newMessage(PlayerMessagePlayerDead)}, player.Messages...)
			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", player.Num).
				Str("PlayerName", player.Name).
				Str("Race", player.Race.PluralName).
				Msgf("player is dead")
		} else if numPlanets == 0 && numFleets > 0 {
			player.Messages = append([]PlayerMessage{newMessage(PlayerMessagePlayerNoPlanets).withSpec(PlayerMessageSpec{Amount: numColonists})}, player.Messages...)
			log.Debug().
				Int64("GameID", t.game.ID).
				Int("Player", player.Num).
				Str("PlayerName", player.Name).
				Str("Race", player.Race.PluralName).
				Int("NumColonists", numColonists).
				Msgf("player has no planets")
		}
	}

}
