package cs

import (
	"fmt"
	"log/slog"
	"math"

	"slices"

	"golang.org/x/exp/maps"
)

// When all players submit their turns, the turnGenerator generator is used to generate a new turnGenerator
// This follows the Stars! order of events: https://wiki.starsautohost.org/wiki/Order_of_Events
type turnGenerator struct {
	game *FullGame
	log  *slog.Logger
}

func newTurnGenerator(game *FullGame) turnGenerator {
	turnLogger := slog.Default().With(
		slog.Int64("GameID", game.ID),
		slog.String("GameName", game.Name),
		slog.Int("Year", game.Year+1), // log for next turn
	)
	t := turnGenerator{game, turnLogger}

	t.game.Universe.setLogger(turnLogger)
	t.game.Universe.buildMaps(game.Players)

	return t
}

// generate a new turn
// TODO: add more error handling. A failed turn generation is easier to fix than
// a corrupt game
func (t *turnGenerator) generateTurn() error {
	t.log.Debug("begin generating turn")
	t.game.Year++

	// reset players for start of the turn
	for _, player := range t.game.Players {
		player.clearTransientIntel()
		player.incrementReportAge()
		player.Messages = []PlayerMessage{}
		player.BattleRecords = []BattleRecord{}
		player.leftoverResources = 0
		player.techLevelGained = false
		player.acquirablePartGained = false
		player.TechsJustGained = []*Tech{}
	}

	t.computeSpecs()
	t.packetInit()

	// wp0 tasks
	t.fleetInit()
	t.fleetByHandLoads()
	t.fleetByHandUnloads()
	t.fleetClearByHandCargoTransfers()
	t.fleetScrap()
	t.fleetUnload()
	t.fleetColonize()
	t.fleetLoad()
	t.fleetMerge()
	t.fleetRoute()
	t.fleetMarkWaypointsProcessed()

	// move stuff through space
	t.packetMove(false)
	t.mysteryTraderMove()
	t.fleetMove()
	t.fleetRadiatingEngineDieoff()
	t.fleetReproduce()
	t.decaySalvage()
	t.decayPackets(false)
	t.wormholeJiggle()
	t.detonateMines()
	t.fleetRemoteMineAR() // sort of a wp1 task, for AR races it happens before production and works on the turn of arrival
	t.planetMine()
	if err := t.planetProduction(); err != nil {
		return err
	}
	t.playerResearch()
	t.permaform()
	t.planetGrow()
	t.packetMove(true)   // move packets built this turn
	t.decayPackets(true) // decay packets built this turn
	t.fleetRefuel()      // refuel after production so fleets will refuel at planets that just built a starbase this turn
	t.randomCometStrike()
	t.randomMineralDeposit()
	t.randomPlanetaryChange()
	t.fleetBattle()
	t.fleetBomb()
	if err := t.mysteryTraderMeet(); err != nil {
		return err
	}
	t.mysteryTraderSpawn()

	// wp1 tasks
	t.fleetRemoteMine() // remote mine for normal miners
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
	if err := t.scan(); err != nil {
		return err
	}

	// notify about battles
	t.checkBattleReports()

	// as a last turn step, calculate scores and check for victories
	t.calculateScores()
	t.checkDeath()

	t.game.State = GameStateWaitingForPlayers

	t.log.Info("generated turn")
	return nil
}

// update all planet specs with the latest info
// useful before turn generation and after building
func (t *turnGenerator) computeSpecs() {
	t.game.computeSpecs()
}

// fleetInit will reset any fleet data before processing
func (t *turnGenerator) fleetInit() {
	for _, fleet := range t.game.Fleets {
		// age this fleet by 1 year
		fleet.Age++

		// remove previous position, it will be reset on move
		fleet.PreviousPosition = nil

		wp0 := &fleet.Waypoints[0]
		wp0.processed = false

		if wp0.Task == WaypointTaskTransport {
			wp0.WaitAtWaypoint = false
		}
	}
}

// fleetByHandLoads will do any by hand cargo transfer load orders
func (t *turnGenerator) fleetByHandLoads() {
	cargoTransferer := newCargoTransferer(t.log, t.game)
	for _, player := range t.game.Players {
		if len(player.CargoTransfers) == 0 {
			continue
		}

		// CargoTransfers are a map of transfer per location
		// process each transfer for a location in order
		for _, transfers := range player.CargoTransfers {
			results := cargoTransferer.loadByHands(player, transfers)

			// for by hand transfers, we only care if something went wrong
			for _, result := range results {
				if result.status == CargoTransferStatusNone {
					// the player assumes all by hand transfer go through, so if it works, don't send any messages
					if result.wanted != result.transferred {
						// we transferred some but not all, someone else got to it first perhaps
						messager.fleetByHandTransferIncomplete(player, result.fleet, result.dest, result.cargoType, result.transferred, result.wanted, result.status)
					}
					continue
				}
				// alert the player of any issues
				messager.fleetByHandTransferIncomplete(player, result.fleet, result.dest, result.cargoType, result.transferred, result.wanted, result.status)
			}
		}
	}
}

// fleetByHandUnloads will do any by hand cargo transfer unload orders
func (t *turnGenerator) fleetByHandUnloads() {
	cargoTransferer := newCargoTransferer(t.log, t.game)
	for _, player := range t.game.Players {
		if len(player.CargoTransfers) == 0 {
			continue
		}

		// CargoTransfers are a map of transfer per location
		// process each transfer for a location in order
		for _, transfers := range player.CargoTransfers {
			results := cargoTransferer.unloadByHands(player, transfers)

			for _, result := range results {
				if result.status == CargoTransferStatusNone {
					// the player assumes all by hand transfer go through, so if it works, don't send any messages
					if result.wanted != result.transferred {
						// we transferred some but not all, someone else got to it first perhaps
						messager.fleetByHandTransferIncomplete(player, result.fleet, result.dest, result.cargoType, result.transferred, result.wanted, result.status)
					}
					continue
				}

				// alert the player of any issues
				messager.fleetByHandTransferIncomplete(player, result.fleet, result.dest, result.cargoType, result.transferred, result.wanted, result.status)
			}
		}
	}

	// resolve any by hand invasions
	t.resolveInvasions(cargoTransferer.invader)
}

// resolveInvasions resolves all invasions for an invader helper
func (t *turnGenerator) resolveInvasions(invader invader) {
	invasions := invader.resolveInvasions(&t.game.Rules)
	for _, invasion := range invasions {
		planet := invasion.planet
		attacker := invasion.attacker
		defender := invasion.defender

		t.log.Debug("planet invaded",
			slog.Int("Defender", defender.Num),
			slog.Int("Attacker", attacker.Num),
			slog.String("Fleet", invasion.fleetDescription()),
			slog.String("Planet", planet.Name),
			slog.Int("Attackers", invasion.attackers),
			slog.Int("Defenders", invasion.defenders),
			slog.Int("RemainingAttackers", invasion.remainingAttackers),
			slog.Int("RemainingDefenders", invasion.remainingDefenders),
			slog.Bool("AttackerWon", invasion.successful),
		)

		// during invasion, even if the player loses the planet, they discover the invader
		for _, fleet := range invasion.fleets {

			for _, token := range fleet.Tokens {
				defender.discoverer.discoverDesign(token.design, defender.Race.Spec.DiscoverDesignOnScan)
			}
			defender.discoverer.discoverFleet(fleet, false)
		}

		// notify each player of the invasion
		messager.planetInvaded(defender, planet, invasion.fleetDescription(), attacker, defender, invasion.attackersKilled, invasion.defendersKilled, invasion.successful)
		messager.planetInvaded(attacker, planet, invasion.fleetDescription(), attacker, defender, invasion.attackersKilled, invasion.defendersKilled, invasion.successful)

		if !invasion.successful {
			// reduce the population to however many colonists remain and move on
			planet.setPopulation(invasion.remainingDefenders)
			continue
		}

		// empty this planet
		planet.emptyPlanet()

		// take over the planet.
		planet.PlayerNum = invasion.attacker.Num
		planet.setPopulation(invasion.remainingAttackers)

		// apply a production plan
		if len(attacker.ProductionPlans) > 0 {
			plan := attacker.ProductionPlans[0]
			plan.Apply(planet)
		}

		// make sure the defender knows about this new planet
		// the last dying colonist sends a report to their compatriots
		defender.discoverer.clearPlanetOwnerIntel(planet)
		defender.discoverer.discoverPlanet(&t.game.Rules, planet, true, true)

		// check for tech trades
		if !attacker.techLevelGained {
			tt := newTechTrader()
			field := tt.checkInvasionTechTrade(&t.game.Rules, attacker, defender.TechLevels)
			if field != TechFieldNone {
				// sweet, we gained a tech level
				attacker.techLevelGained = true
				attacker.TechLevels.Set(field, attacker.TechLevels.Get(field)+1) // add 1 to corresponding lvl

				messager.playerTechGainedInvasion(attacker, planet, field)
				attacker.updateTechsJustGained(t.game.Rules.techs, field)

				t.log.Debug("invader gained tech level",
					slog.Int("Attacker", attacker.Num),
					slog.Int("Defender", defender.Num),
					slog.String("Planet", planet.Name),
					slog.String("field", string(field)),
				)

			}
		}
	}
}

// fleetClearByHandCargoTransfers clear's out all by-hand style cargo transfers after they are processed
func (t *turnGenerator) fleetClearByHandCargoTransfers() {
	for _, p := range t.game.Players {
		p.CargoTransfers = CargoTransfers{}
	}
}

// scrap a fleet at wp0/wp1
func (t *turnGenerator) fleetScrap() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskScrapFleet {
			t.scrapFleet(fleet, false)
		}
	}
}

// scrap a fleet giving a planet resources or creating salvage
func (t *turnGenerator) scrapFleet(fleet *Fleet, colonize bool) {
	player := t.game.getPlayer(fleet.PlayerNum)
	planet := t.game.getOrbitingPlanet(fleet)

	cost := fleet.getScrapAmount(&t.game.Rules, player, planet, colonize)

	if planet != nil {
		// scrap over a planet
		planet.Cargo = planet.Cargo.AddMineral(cost.ToMineral())
		// UR bonus resources only come into play for normal scrapping
		// but fleet.getScrapAmount already sets it to 0 regardless
		planet.bonusResources += cost.Resources
		if planet.OwnedBy(player.Num) {
			// add colonists to planet cargo if it's our own planet
			planet.Cargo = planet.Cargo.Add(fleet.Cargo)
		} else {
			// if not our planet, only the minerals in cargo get transferred (bye bye colonists)
			planet.Cargo = planet.Cargo.AddMineral(fleet.Cargo.ToMineral())
		}

		// Check for level/component tech trading.
		// We do this for every token in the fleet - if it's the player's original ships, it won't lead
		// to a tech trade (they obviously have the tech levels required to build it),
		// but if a ship in the fleet was gifted to them before being scrapped,
		// they should be able to gain tech from it
		if planet.Owned() && planet.Spec.HasStarbase && !colonize {
			planetPlayer := t.game.getPlayer(planet.PlayerNum)
			tt := newTechTrader()
			field, acquiredPart := tt.checkFleetTechTrade(&t.game.Rules, planetPlayer, fleet.Tokens)
			if field != TechFieldNone {
				// we gained a level!
				player.techLevelGained = true
				player.TechLevels.Set(field, player.TechLevels.Get(field)+1)
				messager.playerTechGainedScrappedFleet(planetPlayer, planet, fleet.Name, field)

				planetPlayer.updateTechsJustGained(t.game.TechStore, field)

				t.log.Debug("gained tech level from scrapping fleet",
					slog.Int("Player", planetPlayer.Num),
					slog.String("Planet", planet.Name),
					slog.String("Fleet", fleet.Name),
					slog.String("field", string(field)),
				)
			}

			if acquiredPart != nil {
				// we gained a part!
				player.acquirablePartGained = true
				player.AcquiredTechs[acquiredPart.Name] = true
				messager.playerAcquirablePartGainedScrappedFleet(planetPlayer, planet, fleet.Name, acquiredPart.Name)
				if player.HasTech(acquiredPart) {
					player.TechsJustGained = append(player.TechsJustGained, acquiredPart)
				}

				t.log.Debug("gained tech part from scrapping",
					slog.Int("Player", planetPlayer.Num),
					slog.String("Planet", planet.Name),
					slog.String("Fleet", fleet.Name),
					slog.String("Tech", acquiredPart.Name),
				)
			}
		}
	} else {
		// create salvage
		t.game.getOrCreateSalvage(fleet.Position, player.Num, cost.ToCargo())
	}

	planetName := ""
	if planet != nil {
		planetName = planet.Name
	}

	t.log.Debug("fleet scrapped",
		slog.Int("Player", fleet.PlayerNum),
		slog.String("Planet", planetName),
		slog.String("Fleet", fleet.Name),
		slog.String("Cargo", fmt.Sprintf("%v", fleet.Cargo)),
		slog.String("Scrap", fmt.Sprintf("%v", cost)),
	)

	messager.fleetScrapped(player, fleet, cost, planet)
	t.game.deleteFleet(fleet)
}

// fleetColonize will attempt to colonize planets for any fleets with the Colonize WaypointTask
func (t *turnGenerator) fleetColonize() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskColonize {
			player := t.game.Players[fleet.PlayerNum-1]

			if wp.TargetType != MapObjectTypePlanet {
				messager.fleetColonizeNonPlanet(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if wp.TargetNum == None {
				err := fmt.Errorf("%s attempted to colonize a planet but didn't target a planet", fleet.Name)
				t.log.Error("fleet attempted to colonize a planet but didn't target a planet", slog.String("Fleet", fleet.Name), slog.Any("err", err))
				messager.error(player, err)
				wp.Task = WaypointTaskNone
				continue
			}

			planet := t.game.getPlanet(wp.TargetNum)
			if planet.Owned() {
				messager.fleetColonizeOwnedPlanet(player, planet, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if !fleet.Spec.Colonizer {
				messager.fleetColonizeWithNoModule(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			if fleet.Cargo.Colonists == 0 {
				messager.fleetColonizeWithNoColonists(player, fleet)
				wp.Task = WaypointTaskNone
				continue
			}

			t.log.Debug("colonized planet",
				slog.Int("Player", player.Num),
				slog.String("Planet", planet.Name),
				slog.String("Fleet", fleet.Name),
				slog.Int("Colonists", fleet.Cargo.Colonists*100),
			)

			if fleet.Spec.OrbitalConstructionModule {
				design := player.GetLatestDesign(ShipDesignPurposeStarterColony)
				if design != nil {
					t.buildStarbase(player, planet, design)
				} else {
					t.log.Error("colonizer can't find Starter Colony design",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Planet", planet.Name),
					)
				}
			}

			// colonize the planet and scrap the fleet
			fleet.colonizePlanet(&t.game.Rules, player, planet)
			t.scrapFleet(fleet, true)
			messager.planetColonized(player, planet)
		}
	}
}

// fleetUnload executes wp0/wp1 unload transport tasks for fleets
func (t *turnGenerator) fleetUnload() {
	cargoTransferer := newCargoTransferer(t.log, t.game)

	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest, ok := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			if !ok {
				// unload to salvage in deep space
				dest = t.game.getOrCreateSalvage(fleet.Position, player.Num, Cargo{})
				t.log.Debug("created salvage",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.String("Position", fleet.Position.String()),
				)
			}

			results := cargoTransferer.unload(fleet, dest, wp.TransportTasks)

			for _, result := range results {
				if result.status != CargoTransferStatusNone {
					t.log.Debug("unload cargo failed",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Dest", dest.GetMapObject().Name),
						slog.Int("Transfered", result.transferred),
						slog.String("cargoType", result.cargoType.String()),
						slog.Any("status", result.status),
					)
					messager.fleetTransportInvalid(player, fleet, dest, result.cargoType, result.transferred, result.wanted, result.status)

					continue
				}
				wp.WaitAtWaypoint = wp.WaitAtWaypoint || result.waitAtWaypoint
				t.log.Debug("unloaded cargo",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.String("Dest", dest.GetMapObject().Name),
					slog.Int("Transfered", result.transferred),
					slog.String("cargoType", result.cargoType.String()),
				)
				if result.transferred != 0 {
					messager.fleetTransportedCargo(player, fleet, dest, result.cargoType, result.transferred)
				}
			}
			if planet, ok := dest.(*Planet); ok {
				planet.MarkDirty()
			}
		}
	}

	// resolve any by hand invasions
	t.resolveInvasions(cargoTransferer.invader)
}

func (t *turnGenerator) fleetLoad() {
	cargoTransferer := newCargoTransferer(t.log, t.game)
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskTransport {
			dest, ok := t.game.getCargoHolder(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
			if !ok || dest.Deleted() {
				// can't load from space
				continue
			}

			results := cargoTransferer.load(fleet, dest, wp.TransportTasks)
			for _, result := range results {
				if result.status != CargoTransferStatusNone {
					t.log.Debug("load cargo failed",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Dest", dest.GetMapObject().Name),
						slog.Int("Transfered", result.transferred),
						slog.String("cargoType", result.cargoType.String()),
						slog.Any("status", result.status),
					)
					messager.fleetTransportInvalid(player, fleet, dest, result.cargoType, result.transferred, result.wanted, result.status)

					continue
				}
				wp.WaitAtWaypoint = wp.WaitAtWaypoint || result.waitAtWaypoint
				t.log.Debug("loaded cargo",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.String("Dest", dest.GetMapObject().Name),
					slog.Int("Transfered", result.transferred),
					slog.String("cargoType", result.cargoType.String()),
				)
				if result.transferred != 0 {
					messager.fleetTransportedCargo(player, fleet, dest, result.cargoType, result.transferred)
				}
			}
			if planet, ok := dest.(*Planet); ok {
				planet.MarkDirty()
			}

			// after load, remove the transport task if this isn't a repeating order
			if !fleet.RepeatOrders && !wp.WaitAtWaypoint {
				wp.Task = WaypointTaskNone
				wp.TransportTasks = WaypointTransportTasks{}
			}
		}
	}

	// after load delete any empty salvages or packets
	for _, salvage := range t.game.Salvages {
		// delete this salvage if we emptied it
		salvage.Cargo.Colonists = 0 // make sure we kill off any colonists dumped into deep space
		if salvage.Cargo == (Cargo{}) {
			t.game.deleteSalvage(salvage)

			t.log.Debug("deleted salvage",
				slog.Int("Player", salvage.PlayerNum),
				slog.String("Salvage", salvage.Name),
			)

		}
	}
	for _, packet := range t.game.MineralPackets {
		// delete this packet if we emptied it
		if packet.Cargo == (Cargo{}) {
			t.game.deletePacket(packet)

			t.log.Debug("deleted packet",
				slog.Int("Player", packet.PlayerNum),
				slog.String("Packet", packet.Name),
			)
		}
	}
}

func (t *turnGenerator) fleetMerge() {
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
			t.log.Error("Failed to merge fleets",
				slog.Any("err", err),
				slog.Int("PlayerNum", player.Num),
				slog.Int("Num", fleet.Num),
				slog.Any("fleet", fleet),
				slog.Any("target", target),
			)
			messager.error(player, err)
			continue
		}

		messager.fleetMerged(player, fleet, target)

		t.log.Debug("fleet merged into target",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Target", target.Name),
		)

		// remove this fleet from the universe
		t.game.deleteFleet(fleet)
	}
}

func (t *turnGenerator) fleetRoute() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp := &fleet.Waypoints[0]

		if !wp.processed && wp.Task == WaypointTaskRoute {
			player := t.game.Players[fleet.PlayerNum-1]
			planet := t.game.getOrbitingPlanet(fleet)
			if planet == nil || planet.RouteTargetNum == None {
				// no route
				continue
			}

			mo := t.game.getMapObject(planet.RouteTargetType, planet.RouteTargetNum, planet.RouteTargetPlayerNum)
			if mo == nil {
				t.log.Warn("planet route target not found",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.String("Planet", planet.Name),
				)
				// wipe this planet's route since it's invalid
				planet.RouteTargetNum = None
				planet.RouteTargetPlayerNum = None
				planet.RouteTargetType = MapObjectTypeNone
				planet.MarkDirty()
				continue
			}

			fleet.AddWaypoint(player, WaypointDest{MO: *mo}, len(fleet.Waypoints)-1, false)
			fleet.Waypoints[len(fleet.Waypoints)-1].Task = WaypointTaskRoute

			messager.fleetRouted(player, fleet, planet, mo.Name)

			t.log.Debug("fleet routed to target",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
				slog.String("Planet", planet.Name),
				slog.String("Target", mo.Name),
			)

		}
	}
}

func (t *turnGenerator) fleetNotifyIdle() {
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
		if len(fleet.Waypoints) > 1 {
			continue
		}

		// if we don't have a previous position, we didn't move this round, don't notify
		if fleet.PreviousPosition == nil {
			continue
		}

		if fleet.Waypoints[0].Task == WaypointTaskNone {
			player := t.game.getPlayer(fleet.PlayerNum)
			messager.fleetCompletedAssignedOrders(player, fleet)

			t.log.Debug("fleet idle",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
			)

		}
	}
}

// mark all wp0 as processed so they won't be processed again during wp1 steps
func (t *turnGenerator) fleetMarkWaypointsProcessed() {
	for _, fleet := range t.game.Fleets {
		wp := &fleet.Waypoints[0]
		wp.processed = true
	}
}

// packetInit will reset any packet data before processing
func (t *turnGenerator) packetInit() {
	for _, packet := range t.game.MineralPackets {
		packet.builtThisTurn = false

		if packet.Cargo.Total() == 0 {
			// this packet was probably snatched away by a player
			t.log.Debug("packet empty",
				slog.Int("Player", packet.PlayerNum),
				slog.String("Packet", packet.Name),
			)
			t.game.deletePacket(packet)
		}
	}
}

// move packets through space
// if builtThisTurn is true, this will only move packets that were built this turn (i.e. just launched)
func (t *turnGenerator) packetMove(builtThisTurn bool) {

	for _, packet := range t.game.MineralPackets {
		if packet.Delete {
			continue
		}
		if packet.builtThisTurn != builtThisTurn {
			continue
		}
		player := t.game.getPlayer(packet.PlayerNum)
		planet := t.game.getPlanet(int(packet.TargetPlanetNum))
		var planetPlayer *Player
		var starbase *Fleet
		if planet.Owned() {
			planetPlayer = t.game.getPlayer(planet.PlayerNum)
			starbase = planet.Starbase
		}

		packet.movePacket(&t.game.Rules, player, planet, planetPlayer)

		t.log.Debug("moved packet",
			slog.Int("Player", packet.PlayerNum),
			slog.String("Packet", packet.Name),
			slog.String("Position", packet.Position.String()),
		)

		if planetPlayer != nil && planet.GetPopulation() == 0 {
			// this planet just got killed by a packet
			if starbase != nil {
				t.game.deleteStarbase(starbase)
				planet.Spec.PlanetStarbaseSpec = PlanetStarbaseSpec{}

				t.log.Debug("packet wiped out planet, deleting starbase",
					slog.Int("Player", planetPlayer.Num),
					slog.String("Packet", packet.Name),
					slog.String("Planet", planet.Name),
				)

			}
		}
	}
}

func (t *turnGenerator) mysteryTraderSpawn() {
	if !t.game.RandomEvents {
		// no mystery traders if no random events
		return
	}

	if len(t.game.MysteryTraders) >= t.game.Rules.MysteryTraderRules.MaxMysteryTraders {
		t.log.Debug("Max MysteryTraders reached, not generating")
		return
	}

	mt := generateMysteryTrader(&t.game.Rules, t.game.Game, t.game.getNextMysteryTraderNum())
	if mt != nil {
		// a mystery trader has spawned!
		t.game.addMysteryTrader(mt)

		// tell all the players
		for _, player := range t.game.Players {
			player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderDiscovered, mt))
		}

		t.log.Debug("MysteryTrader spawned",
			slog.Int("MysteryTrader", mt.Num),
			slog.String("Position", mt.Position.String()),
			slog.String("Destination", mt.Destination.String()),
		)
	}
}

func (t *turnGenerator) mysteryTraderMove() {
	for _, mt := range t.game.MysteryTraders {
		if mt.Delete {
			continue
		}

		// check for a course change
		if mt.change(&t.game.Rules, t.game.Game) {
			for _, player := range t.game.Players {
				player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderChangedCourse, mt))
			}
			t.log.Debug("mysteryTrader changed course",
				slog.Int("MysteryTrader", mt.Num),
			)
		}

		originalPosition := mt.Position
		mt.move()
		t.game.moveMysteryTrader(mt, originalPosition)

		t.log.Debug("moved mysteryTrader",
			slog.Int("MysteryTrader", mt.Num),
			slog.Int("WarpSpeed", mt.WarpSpeed),
			slog.String("Start", originalPosition.String()),
			slog.String("End", mt.Position.String()),
		)

		if mt.Position == mt.Destination {
			if mt.again(&t.game.Rules, t.game.Game, t.game.GetNumHumanPlayers()) {
				for _, player := range t.game.Players {
					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderAgain, mt))
				}
				t.log.Debug("mysteryTrader going again",
					slog.Int("MysteryTrader", mt.Num),
				)
			} else {
				t.log.Debug("mysteryTrader finished",
					slog.Int("MysteryTrader", mt.Num),
				)
				// all done, bye bye trader
				t.game.deleteMysteryTrader(mt)
			}
		}
	}
}

func (t *turnGenerator) fleetMove() {

	fleetsTargetingMovers := []*Fleet{}

	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}
		if fleet.Starbase {
			continue
		}

		if len(fleet.Waypoints) > 1 {
			wp0 := fleet.Waypoints[0]
			wp1 := fleet.Waypoints[1]

			// no move this turn, we wait
			if wp0.WaitAtWaypoint {
				continue
			}

			if wp1.TargetType == MapObjectTypeFleet || wp1.TargetType == MapObjectTypeMineralPacket || wp1.TargetType == MapObjectTypeMysteryTrader {
				// move this after all the fleets not targeting fleets move
				fleetsTargetingMovers = append(fleetsTargetingMovers, fleet)
				continue
			}

			t.moveFleet(fleet)
		} else {
			fleet.WarpSpeed = 0
			fleet.Heading = Vector{}
		}
	}

	// move all the fleets targeting other fleets
	// TODO: build a directed graph and detect cycles and all that jazz
	for _, fleet := range fleetsTargetingMovers {
		t.moveFleet(fleet)
	}
}

// move the actual fleet in the universe from a to b handling minefield destruction, engine strain, stargates, etc
func (t *turnGenerator) moveFleet(fleet *Fleet) {
	player := t.game.getPlayer(fleet.PlayerNum)
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
			t.log.Debug("fleet target gone, using position only",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
			)
		} else if target.Position != wp1.Position {
			// update the position
			wp1.Position = target.Position
			t.log.Debug("fleet target moved, updating position",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
			)
		}
	}

	if wp1.WarpSpeed == StargateWarpSpeed {
		// yeah, gate!
		fleet.gateFleet(&t.game.Rules, t.game.Universe, t.game)
	} else {
		interrupted := fleet.moveFleet(&t.game.Rules, t.game.Universe, t.game)
		if interrupted != nil {
			switch interrupted.reason {
			case fleetMoveInterruptedHitMinefield:
				// damage the fleet in the minefield
				minefield := interrupted.minefield
				minefieldPlayer := t.game.getPlayer(minefield.PlayerNum)
				stats := t.game.Rules.MinefieldStatsByType[minefield.MinefieldType]

				damage := minefield.damageFleet(fleet, player, stats)
				minefield.reduceMinefieldOnImpact()
				if minefieldPlayer.Race.Spec.MinefieldsAreScanners {
					// SD races discover the exact fleet makeup
					for _, token := range fleet.Tokens {
						// SD races discover the exact fleet makeup
						minefieldPlayer.discoverer.discoverDesign(token.design, true)
					}
				}

				// tell the fleet owner and the minefield owner the fleet was hit
				messager.fleetMinefieldHit(player, fleet, minefield, damage)
				if minefield.PlayerNum != player.Num {
					messager.fleetMinefieldHit(minefieldPlayer, fleet, minefield, damage)
				}

				t.log.Debug("minefield damaged fleet",
					slog.Int("Player", minefield.PlayerNum),
					slog.String("Minefield", minefield.Name),
					slog.String("Fleet", fleet.Name),
					slog.Int("FleetPlayer", fleet.PlayerNum),
					slog.Int("TotalDamage", damage.Damage),
					slog.Int("ShipsDestroyed", damage.ShipsDestroyed),
					slog.Bool("FleetDestroyed", damage.FleetDestroyed),
				)

			}
		} else {
			// check for exploded ships from overwarp
			explodedShips := fleet.applyOverwarpPenalty(&t.game.Rules)
			// tell the player they lost ships
			if explodedShips > 0 {
				t.log.Debug("fleet ships exploded due to unsafe warp",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.Int("ExplodedShips", explodedShips),
					slog.Int("Warp", wp1.WarpSpeed),
				)
				messager.fleetExceededSafeSpeed(player, fleet, explodedShips)
			}
		}
	}

	t.log.Debug("moved fleet",
		slog.Int("Player", fleet.PlayerNum),
		slog.String("Fleet", fleet.Name),
		slog.String("Fuel", fmt.Sprintf("%d/%d", fleet.Fuel, fleet.Spec.FuelCapacity)),
		slog.Int("WarpSpeed", fleet.WarpSpeed),
		slog.String("Start", wp0.Position.String()),
		slog.String("End", fleet.Position.String()),
	)

	// update the game dictionaries with this fleet's new position
	t.game.moveFleet(fleet, originalPosition)

	// make sure we have tokens left after move
	fleet.removeEmptyTokens()
	if len(fleet.Tokens) == 0 {
		t.log.Debug("deleted fleet after move",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
		)
		t.game.deleteFleet(fleet)
		return
	}

	// make sure we don't have extra fuel if we lost ships during movement
	fleet.Spec = ComputeFleetSpec(&t.game.Rules, player, fleet)
	fleet.reduceFuelToMax()

	// remove the previous waypoint, it's been processed already
	if fleet.RepeatOrders && !wp0.PartiallyComplete {
		// if we are supposed to repeat orders,
		wp0.processed = false
		wp0.WaitAtWaypoint = false
		wp0.PartiallyComplete = false
		fleet.Waypoints = append(fleet.Waypoints, wp0)

		t.log.Debug("repeating waypoint",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
			slog.String("Waypoint", fmt.Sprintf("%s: %s", wp0.TargetName, wp0.Task)),
		)
	}
}

// kill off colonists on fleets from radiation poisoning.
// TODO: Make this part of the TechHullComponent
func (t *turnGenerator) fleetRadiatingEngineDieoff() {
	// https://wiki.starsautohost.org/wiki/Radiating_Ramscoop
	// DeathRate/Year % = int ((86 - C)/2)
	// where C is the center of your Rad Hab range (mR)
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		// no radiation in this fleet or no colonists to kill
		if !fleet.Spec.Radiating || fleet.Cargo.Colonists == 0 {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		if player.Race.IsImmune(Rad) {
			// rad immune races could care less about engine radiation
			continue
		}

		habCenter := player.Race.Spec.HabCenter
		clicksAway := max(0, t.game.Rules.RadiatingImmune-habCenter.Rad)
		if clicksAway <= 0 {
			// race has high enough of a hab center to be unaffected by radiation
			continue
		}
		deathRate := math.Round(float64(clicksAway)/2) / 100

		killed := max(1, int(deathRate*float64(fleet.Cargo.Colonists)))
		fleet.Cargo.Colonists -= killed

		// Message the player
		messager.fleetRadiatingEngineDieoff(player, fleet, killed*100)

		t.log.Debug("fleet radiation dieoff",
			slog.Int("Player", fleet.PlayerNum),
			slog.String("Fleet", fleet.Name),
		)
	}
}

func (t *turnGenerator) fleetReproduce() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete || fleet.Cargo.Colonists == 0 {
			continue
		}

		// check if this player's freighters reproduce
		player := t.game.getPlayer(fleet.PlayerNum)
		fg := player.Race.Spec.FreighterGrowth
		if fg.GrowthFactor == 0 {
			continue
		}

		var growth int
		if fg.Absolute {
			// calculate absolute pop growth on fleets
			// TODO: Check rounding on this...?
			growth = int(fg.GrowthFactor * float64(fleet.Cargo.Colonists))
		} else {
			// Calculate relative pop growth based on growth rate
			growth = int(fg.GrowthFactor * float64(fleet.Cargo.Colonists*player.Race.GrowthRate) / 100)
		}
		fleet.Cargo.Colonists = fleet.Cargo.Colonists + growth
		over := max(0, fleet.Cargo.Total()-fleet.Spec.CargoCapacity)

		planet := t.game.getOrbitingPlanet(fleet)
		if over > 0 {
			// remove excess colonists, dumping them onto our own planets if possible
			fleet.Cargo.Colonists = fleet.Cargo.Colonists - over
			if planet != nil && planet.OwnedBy(fleet.PlayerNum) {
				// add colonists to the planet this fleet is orbiting
				planet.Cargo.Colonists = planet.Cargo.Colonists + over
			}
		}

		// Send the appropriate message to the player
		if growth > 0 {
			messager.fleetReproduce(player, fleet, growth*100, planet, over)

			t.log.Debug("fleet reproduced",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
				slog.Int("Growth", growth),
				slog.Int("Pop overflow", over),
			)
		} else {
			messager.fleetDieOff(player, fleet, growth*100)
			t.log.Debug("fleet died off",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
				slog.Int("Deaths", growth),
			)
		}
	}
}

// decay each salvage and remove it from the universe if it's empty
func (t *turnGenerator) decaySalvage() {
	for _, salvage := range t.game.Salvages {
		beforeCargo := salvage.Cargo
		salvage.decay(&t.game.Rules)

		t.log.Debug("decayed salvage",
			slog.Int("Player", salvage.PlayerNum),
			slog.String("Salvage", salvage.Name),
			slog.String("CargoBefore", beforeCargo.PrettyString()),
			slog.String("CargoAfter", salvage.Cargo.PrettyString()),
		)
		if (salvage.Cargo == Cargo{}) {
			t.game.deleteSalvage(salvage)

			t.log.Debug("deleted salvage",
				slog.Int("Player", salvage.PlayerNum),
				slog.String("Salvage", salvage.Name),
			)
		}
	}
}

// Decay mineral packets in flight
func (t *turnGenerator) decayPackets(builtThisTurn bool) {
	for _, packet := range t.game.MineralPackets {
		if packet.Delete {
			continue
		}
		if packet.builtThisTurn != builtThisTurn {
			continue
		}

		player := t.game.getPlayer(packet.PlayerNum)
		// update the decay amount based on this distance traveled this turn
		decayRate := packet.getPacketDecayRate(&t.game.Rules, &player.Race) * (packet.distanceTravelled / float64(packet.WarpSpeed*packet.WarpSpeed))

		// skip calcs if no decay
		if decayRate == 0 {
			continue
		}

		// loop through all 3 mineral types and reduce each one in turn
		for _, minType := range [3]CargoType{Ironium, Boranium, Germanium} {
			mineral := float64(packet.Cargo.GetAmount(minType))
			decayAmount := max(int(decayRate*mineral), int(float64(t.game.Rules.PacketMinDecay)*player.Race.Spec.PacketDecayFactor))
			packet.Cargo = packet.Cargo.SubtractAmount(minType, decayAmount)
			packet.Cargo = packet.Cargo.MinZero()
		}
		t.log.Debug("decayed packet",
			slog.Int("Player", packet.PlayerNum),
			slog.String("Packet", packet.Name),
			slog.String("Cargo", packet.Cargo.PrettyString()),
		)
		// delete empty packets
		if packet.Cargo.Total() == 0 {
			t.game.deletePacket(packet)
			t.log.Debug("deleted packet",
				slog.Int("Player", packet.PlayerNum),
				slog.String("Packet", packet.Name),
			)
		}
	}
}

// jiggle, degrade, and jump wormholes
func (t *turnGenerator) wormholeJiggle() {
	if len(t.game.Wormholes) == 0 {
		return
	}

	planetPositions := make([]Vector, len(t.game.Planets))
	wormholePositions := make([]Vector, len(t.game.Wormholes))

	for _, wormhole := range t.game.Wormholes {
		originalPosition := wormhole.Position
		wormhole.jiggle(t.game.Area, t.game.Universe, t.game.Rules.random)
		t.game.moveWormhole(wormhole, originalPosition)

		wormhole.degrade()
		if wormhole.shouldJump(t.game.Rules.random) {
			// this wormhole jumped. We actually delete the previous one and create a new one. This way scanner history is reset
			position, _, err := generateWormhole(t.game.Universe, t.game.Area, t.game.Rules.random, planetPositions, wormholePositions, t.game.Rules.WormholeMinPlanetDistance)
			if err != nil {
				// don't kill turn generation over this, just move on without a new wormhole
				t.log.Error("failed to generate new wormhole after wormhole jump", slog.Any("err", err))
				continue
			}

			// create the new wormhole
			companion := t.game.Universe.getWormhole(wormhole.DestinationNum)
			newWormhole := t.game.createWormhole(&t.game.Rules, position, WormholeStabilityRockSolid, companion)

			// queue the old wormhole for deletion and add the new wormhole to the universe
			t.game.deleteWormhole(wormhole)

			t.log.Debug("generated new wormhole after jump", slog.Int("num", newWormhole.Num), slog.Any("position", newWormhole.Position))
		}

		// update the spec
		wormhole.Spec = computeWormholeSpec(wormhole, &t.game.Rules)
	}
}

// SD races can detonate a minefield
func (t *turnGenerator) detonateMines() {
	for _, minefield := range t.game.Minefields {
		if !minefield.Detonate {
			continue
		}

		stats := t.game.Rules.MinefieldStatsByType[minefield.MinefieldType]
		if !stats.CanDetonate {
			continue
		}

		minefieldPlayer := t.game.getPlayer(minefield.PlayerNum)
		fleetsWithin := t.game.fleetsWithin(minefield.Position, minefield.Radius())
		for _, fleet := range fleetsWithin {
			fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
			damage := minefield.damageFleet(fleet, fleetPlayer, stats)

			if damage == (MinefieldDamage{}) {
				// no damage, probably immune
				continue
			}

			if minefieldPlayer.Race.Spec.MinefieldsAreScanners && minefieldPlayer.Num != fleet.PlayerNum {
				// SD races discover the exact fleet makeup
				for _, token := range fleet.Tokens {
					// SD races discover the exact fleet makeup
					minefieldPlayer.discoverer.discoverDesign(token.design, true)
				}
			}

			messager.fleetMinefieldHit(fleetPlayer, fleet, minefield, damage)
			if minefield.PlayerNum != fleetPlayer.Num {
				messager.fleetMinefieldHit(minefieldPlayer, fleet, minefield, damage)
			}

			// clear out any destroyed tokens
			fleet.removeEmptyTokens()

			t.log.Debug("minefield detonation damaged fleet",
				slog.Int("Player", minefield.PlayerNum),
				slog.String("Minefield", minefield.Name),
				slog.String("Fleet", fleet.Name),
				slog.Int("FleetPlayer", fleet.PlayerNum),
				slog.Int("TotalDamage", damage.Damage),
				slog.Int("ShipsDestroyed", damage.ShipsDestroyed),
				slog.Bool("FleetDestroyed", damage.FleetDestroyed),
			)
			if damage.FleetDestroyed {
				t.game.deleteFleet(fleet)
			}
		}

		// reduce minefield after detonation
		minefield.NumMines -= minefield.NumMines / 4

		t.log.Debug("detonated minefield",
			slog.Int("Player", minefield.PlayerNum),
			slog.String("Minefield", minefield.Name),
			slog.Int("NumMines", minefield.NumMines),
		)
	}
}

// mine all owned planets for minerals
func (t *turnGenerator) planetMine() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			planet.mine(&t.game.Rules, planet.Spec.MiningOutput, planet.Mines)
			t.log.Debug("planet mined",
				slog.Int("Player", planet.PlayerNum),
				slog.String("Planet", planet.Name),
				slog.String("Minerals", planet.Spec.MiningOutput.PrettyString()),
			)
		}
	}
}

// remote mine AR-owned planets with remote mining fleets in orbit
func (t *turnGenerator) fleetRemoteMineAR() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp0 := &fleet.Waypoints[0]
		if wp0.Task != WaypointTaskRemoteMining {
			continue
		}
		player := t.game.getPlayer(fleet.PlayerNum)
		planet := t.game.getOrbitingPlanet(fleet)

		// can't remote mine deep space
		if planet == nil {
			messager.fleetRemoteMineDeepSpace(player, fleet)
			wp0.Task = WaypointTaskNone
			continue
		}

		// no miners no minerals
		if fleet.Spec.MiningRate == 0 {
			messager.fleetRemoteMineNoMiners(player, fleet, planet)
			fleet.Waypoints[0].Task = WaypointTaskNone
			continue
		}

		// If this is our own planet, remote mine it (happens earlier than normal)
		if planet.OwnedBy(fleet.PlayerNum) && player.Race.Spec.CanRemoteMineOwnPlanets {
			t.remoteMine(fleet, player, planet)
		}
	}
}

// remote mine planets
func (t *turnGenerator) fleetRemoteMine() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp0 := &fleet.Waypoints[0]
		if wp0.Task != WaypointTaskRemoteMining {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		planet := t.game.getOrbitingPlanet(fleet)

		// Skip AR self remote mining (since it's already been covered prior)
		if planet != nil && planet.OwnedBy(fleet.PlayerNum) && player.Race.Spec.CanRemoteMineOwnPlanets {
			continue
		}

		// can't remote mine deep space
		if planet == nil {
			messager.fleetRemoteMineDeepSpace(player, fleet)
			wp0.Task = WaypointTaskNone
			continue
		}

		if planet.Owned() {
			messager.fleetRemoteMineInhabited(player, fleet, planet)
			wp0.Task = WaypointTaskNone
			continue
		}

		if fleet.Spec.MiningRate == 0 {
			messager.fleetRemoteMineNoMiners(player, fleet, planet)
			fleet.Waypoints[0].Task = WaypointTaskNone
			continue
		}

		if fleet.PreviousPosition != nil {
			// just got here this turn; don't mine
			continue
		}

		t.remoteMine(fleet, player, planet)
	}
}

// remote mine a planet
func (t *turnGenerator) remoteMine(fleet *Fleet, player *Player, planet *Planet) {
	output := planet.getMineralOutput(&t.game.Rules, fleet.Spec.MiningRate, t.game.Rules.RemoteMiningMineOutput)
	planet.mine(&t.game.Rules, output, fleet.Spec.MiningRate)
	planet.MarkDirty()

	// make sure we know about this planet's cargo after remote mining;
	// mark this fleet as having remote mined so it doesn't get counted twice
	fleet.remoteMined = true
	messager.fleetRemoteMined(player, fleet, planet, output)

	t.log.Debug("fleet remote mined planet",
		slog.Int("Player", fleet.PlayerNum),
		slog.String("Fleet", fleet.Name),
		slog.String("Planet", planet.Name),
		slog.String("Mineral output", output.PrettyString()),
	)
}

// process all owned planets' production queues
func (t *turnGenerator) planetProduction() error {
	for _, planet := range t.game.Planets {
		if !planet.Owned() {
			// unowned planets can't build anything (duh)
			continue
		}

		player := t.game.Players[planet.PlayerNum-1]
		producer := newProducer(t.log, &t.game.Rules, planet, player)
		result, err := producer.produce()
		if err != nil {
			return err
		}

		// Now it's time to handle the fruits of our labor!

		// Send any pre-made error messages as well as installation messages
		if len(result.messages) > 0 {
			player.Messages = append(player.Messages, result.messages...)
		}
		if result.mines > 0 {
			messager.planetBuiltMines(player, planet, result.mines)
		}
		if result.factories > 0 {
			messager.planetBuiltFactories(player, planet, result.factories)
		}
		if result.defenses > 0 {
			messager.planetBuiltDefenses(player, planet, result.defenses)
		}

		// message about mineral alchemy
		if result.alchemy > 0 {
			messager.planetBuiltMineralAlchemy(player, planet, result.alchemy)
		}

		// message about each terraform step
		if len(result.terraformResults) > 0 {
			for _, terraformResult := range result.terraformResults {
				messager.planetTerraform(player, planet, terraformResult.Type, terraformResult.Direction)
			}
		}

		// handle built fleets
		// TODO: Add option to merge with existing fleets at location
		for _, token := range result.tokens {
			design := token.design
			if design == nil {
				return fmt.Errorf("player %d has no design %d", player.Num, token.DesignNum)
			}
			design.Spec.NumBuilt += token.Quantity
			design.Spec.NumInstances += token.Quantity

			player.Stats.FleetsBuilt++
			player.Stats.TokensBuilt += token.Quantity

			var routeTarget *MapObject
			if planet.RouteTargetNum != None {
				routeTarget = t.game.getMapObject(planet.RouteTargetType, planet.RouteTargetNum, planet.RouteTargetPlayerNum)
			}
			fleet, err := t.buildFleet(player, planet, token.ShipToken, token.tags, routeTarget)
			if err != nil {
				return err
			}
			messager.fleetBuilt(player, planet, fleet, token.Quantity, routeTarget)
		}

		// yeet packets
		if result.packets != (Cargo{}) {
			target := t.game.getPlanet(planet.PacketTargetNum)
			packet := t.buildMineralPacket(player, planet, result.packets, target)
			messager.planetBuiltMineralPacket(player, planet, packet)
		}

		// build bases
		if result.starbase != nil {
			starbase, err := t.buildStarbase(player, planet, result.starbase)
			if err != nil {
				return err
			}
			planet.Starbase = starbase
			planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)
			messager.planetBuiltStarbase(player, planet, starbase)
		}

		// planetary scanner
		if result.scanner {
			planet.Scanner = true
			messager.planetBuiltScanner(player, planet, planet.Spec.Scanner)
		}

		// genesis device
		if result.reset {
			planet.randomize(&t.game.Rules, t.game.StartMode == GameStartModeAccBBS)
			planet.Mines = 0
			planet.Factories = 0
			messager.planetBuiltGenesisDevice(player, planet)
		}

		// re-compute spec once for brevity
		if result.scanner || result.reset {
			planet.Spec = ComputePlanetSpec(&t.game.Rules, player, planet)
		}

		// log what we actually did
		for _, itemBuilt := range result.itemsBuilt {
			if itemBuilt.numBuilt > 0 {
				t.log.Debug("built production item",
					slog.Int("Player", planet.PlayerNum),
					slog.String("Planet", planet.Name),
					slog.String("Item", string(itemBuilt.queueItemType)),
					slog.Int("DesignNum", itemBuilt.designNum),
					slog.Int("NumBuilt", itemBuilt.numBuilt),
				)
			}
		}

		// any leftover resources go back to the player for research
		// TODO: auto dump this in alchemy if techs are maxed
		player.leftoverResources += result.leftoverResources
	}
	return nil
}

// build a fleet with some number of tokens
func (t *turnGenerator) buildFleet(player *Player, planet *Planet, token ShipToken, tags Tags, routeTarget *MapObject) (*Fleet, error) {
	fleet, err := t.addFleet(player, planet.Position, token, tags)
	if err != nil {
		return nil, err
	}
	fleet.OrbitingPlanetNum = planet.Num

	fleet.Waypoints[0] = NewPlanetWaypoint(planet.Position, planet.Num, planet.Name, token.design.Spec.Engine.IdealSpeed)
	var route string
	if routeTarget != nil {
		route = routeTarget.Name
		fleet.AddWaypoint(player, WaypointDest{MO: *routeTarget}, 0, false)
		fleet.Waypoints[len(fleet.Waypoints)-1].Task = WaypointTaskRoute // follow route by default

	}

	t.log.Debug("fleet built",
		slog.Int("Player", fleet.PlayerNum),
		slog.String("Planet", planet.Name),
		slog.String("Fleet", fleet.Name),
		slog.String("Route", route),
	)
	return fleet, nil
}

// add a new fleet to the universe
func (t *turnGenerator) addFleet(player *Player, position Vector, token ShipToken, tags Tags) (*Fleet, error) {
	playerFleets := t.game.getFleets(player.Num)
	fleetNum := player.GetNextFleetNum(playerFleets)
	fleet := newFleetForToken(player, fleetNum, token, []Waypoint{NewPositionWaypoint(position, token.design.Spec.Engine.IdealSpeed)})
	fleet.Position = position
	fleet.Spec = ComputeFleetSpec(&t.game.Rules, player, &fleet)
	fleet.Fuel = fleet.Spec.FuelCapacity
	fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
	fleet.Tags = tags

	t.game.Fleets = append(t.game.Fleets, &fleet)
	if err := t.game.Universe.addFleet(&fleet); err != nil {
		return nil, err
	}
	return &fleet, nil
}

// build a starbase on a planet
func (t *turnGenerator) buildStarbase(player *Player, planet *Planet, design *ShipDesign) (*Fleet, error) {
	player.Stats.StarbasesBuilt++
	player.Stats.TokensBuilt++
	design.Spec.NumBuilt++

	var prevDamage float64
	var prevArmor int
	// remove the old starbase, tracking its prior damage
	if planet.Starbase != nil {
		// TODO: Make this account for quantity if or when multi token starbases become a thing
		prevDamage = planet.Starbase.Tokens[0].Damage
		prevArmor = planet.Starbase.Tokens[0].design.Spec.Armor
		t.game.deleteStarbase(planet.Starbase)
		planet.Starbase = nil
		planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)
	}

	starbase := newStarbase(player, planet, design, design.Name)
	starbase.Spec = ComputeFleetSpec(&t.game.Rules, player, &starbase)

	// if the prior starbase was damaged, set the new base's damage proportional to the old base's dmg%
	if prevDamage > 0 && prevArmor > 0 {
		starbase.Tokens[0].QuantityDamaged = 1
		starbase.Tokens[0].Damage = (prevDamage / float64(prevArmor)) * float64(starbase.Tokens[0].design.Spec.Armor)
	}

	planet.setStarbase(&starbase)
	t.log.Debug("built starbase",
		slog.Int("Player", starbase.PlayerNum),
		slog.String("Planet", planet.Name),
		slog.String("Starbase", starbase.Name),
	)
	t.game.Starbases = append(t.game.Starbases, &starbase)
	if err := t.game.addStarbase(&starbase); err != nil {
		return nil, err
	}
	return &starbase, nil
}

// build a mineral packet with cargo
func (t *turnGenerator) buildMineralPacket(player *Player, planet *Planet, cargo Cargo, target *Planet) *MineralPacket {

	playerMineralPackets := t.game.getMineralPackets(player.Num)
	num := player.getNextMineralPacketNum(playerMineralPackets)
	packet := newMineralPacket(player, num, planet.PacketSpeed, planet.Spec.SafePacketSpeed, cargo, planet.Position, target.Num)
	packet.builtThisTurn = true

	t.log.Debug("mineral packet built",
		slog.Int("Player", packet.PlayerNum),
		slog.String("Planet", planet.Name),
		slog.String("Fleet", packet.Name),
	)
	t.game.MineralPackets = append(t.game.MineralPackets, packet)
	return packet
}

func (t *turnGenerator) playerResearch() error {
	r := newResearcher(&t.game.Rules)

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

		messager.playerGainTechLevel(player, field, player.TechLevels.Get(field), player.Researching)
		player.updateTechsJustGained(t.game.TechStore, field)
		playerGainedLevel[player.Num] = true

		t.log.Debug("player researched new tech level",
			slog.Int("Player", player.Num),
			slog.String("Field", string(field)),
			slog.Int("Level", player.TechLevels.Get(field)),
		)
	}

	// keep track of how many research resources are stealable by other players
	stealableResearchResources := TechLevel{}

	// handle bonus artifacts
	if t.game.RandomEvents {

		for _, planet := range t.game.Planets {
			if planet.RandomArtifact && planet.Owned() {
				/*
				   // must be owned, have an artifact, and random events must be allowed
				   if (pl->iPlayer != -1 && pl->fArtifact && !game.fNoRandom) {
				       // consume the artifact so it only fires once
				       pl->fArtifact = 0;

				       // pick a random tech field [0..5] and a base bonus [100..400]
				       int tech  = Random(6);
				       int bonus = Random(301) + 100;

				       // small/early colonies get a proportionally smaller bonus
				       if (colonyScale < 10) {
				           bonus = (colonyScale * bonus) / 10;
				       }

				       // notify the player (message 94). params: (owner, id, from=-2, planet id, tech, bonus, 0..)
				       FSendPlrMsg(pl->iPlayer, 94, -2, pl->id, tech, bonus, 0, 0, 0, 0);

				       // apply the research bonus
				       rgplr[(unsigned)pl->iPlayer].rgResSpent[tech] += (unsigned long)bonus;

				       // some rule-sets halve the effective amount (original checked (game.wCrap & 0x0002) != 0)
				       if (game.fSlowTech) {
				           bonus >>= 1;
				       }

				       // (nothing else to do here; the “halved” value only affected the local shown/returned amount)
				   }
				*/

				// score, we got a new artifact, but only once
				planet.RandomArtifact = false

				// figure out which field we research
				player := t.game.getPlayer(planet.PlayerNum)
				bonusRange := t.game.Rules.RandomArtifactResearchBonusRange
				amount := t.game.Rules.random.Intn(bonusRange[1]-bonusRange[0]) + bonusRange[0]
				field := TechFields[t.game.Rules.random.Intn(len(TechFields))]
				messager.planetBonusResearchArtifact(player, planet, amount, field)

				// research the field this random artifact came in
				r.researchField(player, field, amount, onLevelGained)
				stealableResearchResources.Set(field, stealableResearchResources.Get(field)+amount)
				player.ResearchSpentLastYear += amount

				t.log.Debug("player found a research bonus artifact",
					slog.Int("Player", player.Num),
					slog.String("Planet", planet.Name),
					slog.Int("Amount", amount),
					slog.String("Field", string(field)),
				)
			}
		}
	}

	// finally, do regular research for each player
	for _, player := range t.game.Players {
		primaryField := player.Researching
		resourcesToSpend := int(float64(resourcesToSpendByPlayer[player.Num])*player.Race.Spec.ResearchFactor + .5)
		player.ResearchSpentLastYear = resourcesToSpend

		// research tech levels until the resources run out
		spent := r.research(player, resourcesToSpend, onLevelGained)
		stealableResearchResources = stealableResearchResources.Add(spent)

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
		if stolenResearch.Total() > 0 {
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
		player.Spec = ComputePlayerSpec(player, &t.game.Rules)

		// update design spec
		for i := range player.Designs {
			var err error
			design := player.Designs[i]

			// store the numBuilt/numInstances because the spec resets them to 0
			numBuilt := design.Spec.NumBuilt
			numInstances := design.Spec.NumInstances
			design.Spec, err = ComputeShipDesignSpec(&t.game.Rules, player.TechLevels, player.Race.Spec, design)
			if err != nil {
				return fmt.Errorf("ComputeShipDesignSpec returned error: %w", err)
			}
			design.Spec.NumBuilt = numBuilt
			design.Spec.NumInstances = numInstances
		}
	}

	// update planet specs for players who gained a level
	for _, planet := range t.game.Planets {
		if !playerGainedLevel[planet.PlayerNum] {
			continue
		}
		planet.Spec = ComputePlanetSpec(&t.game.Rules, t.game.Players[planet.PlayerNum-1], planet)
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
	return nil
}

// for each planet, randomly check if the owner permaforms it
func (t *turnGenerator) permaform() {

	terraformer := NewTerraformer()

	for _, planet := range t.game.Planets {
		if planet.Owned() {
			player := t.game.Players[planet.PlayerNum-1]
			if player.Race.Spec.PermaformChance == 0 || player.Race.Spec.PermaformPopulation == 0 {
				continue
			}
			adjustedPermaformChance := player.Race.Spec.PermaformChance
			if planet.GetPopulation() <= player.Race.Spec.PermaformPopulation {
				adjustedPermaformChance *= float64(planet.GetPopulation() / player.Race.Spec.PermaformPopulation)
			}

			if adjustedPermaformChance >= t.game.Rules.random.Float64() {
				habType := HabTypes[t.game.Rules.random.Intn(len(HabTypes))]
				result := terraformer.PermaformOneStep(planet, player, habType)

				if result.Terraformed() {
					planet.Spec = ComputePlanetSpec(&t.game.Rules, player, planet)
					planet.MarkDirty()
					messager.planetPermaform(player, planet, result.Type, result.Direction)

					t.log.Debug("player permaformed planet",
						slog.Int("Player", player.Num),
						slog.Int("Planet", planet.Num),
						slog.String("HabType", result.Type.String()),
					)
				}
			}
		}
	}

}

// grow all owned planets by some population
func (t *turnGenerator) planetGrow() {
	for _, planet := range t.game.Planets {
		if !planet.Owned() {
			// can't grow what doesn't exist
			continue
		}
		player := t.game.getPlayer(planet.PlayerNum)
		prevPop := planet.exactPopulation()
		planet.grow(player)

		// tell players about dying colonists
		if planet.Spec.GrowthAmount < 0 {
			if player.Race.GetPlanetHabitability(planet.Hab) < 0 {
				// negative hab pop loss takes priority over overcrowding deaths, so the messages should too
				messager.planetPopulationDecreased(player, planet, prevPop, planet.GetPopulation())
			} else {
				messager.planetPopulationDecreasedOvercrowding(player, planet, planet.Spec.GrowthAmount)
			}
		}

		t.log.Debug("planet grew",
			slog.Int("Player", planet.PlayerNum),
			slog.String("Planet", planet.Name),
			slog.Int("Capacity", int(planet.Spec.PopulationDensity*100)),
			slog.Int("PrevPopulation", prevPop),
			slog.Int("GrowthAmount", planet.Spec.GrowthAmount),
			slog.Int("Population", planet.exactPopulation()),
		)
		if planet.GetPopulation() <= 0 { // should never happen, but covers our bases
			planet.emptyPlanet()
			messager.planetDiedOff(player, planet)

			t.log.Warn("planet pop died off after growth",
				slog.Int("Player", player.Num),
				slog.String("Planet", planet.Name),
			)
		}
	}
}

// refuel fleets if they are orbiting a planet with a friendly starbase
func (t *turnGenerator) fleetRefuel() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		if fleet.Fuel == fleet.Spec.FuelCapacity {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)

		if fleet.Spec.FuelGeneration > 0 {
			fleet.Fuel = Clamp(fleet.Fuel+fleet.Spec.FuelGeneration, 0, fleet.Spec.FuelCapacity)
			fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)
			t.log.Debug("fleet generated fuel",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
			)
		}

		planet := t.game.getOrbitingPlanet(fleet)
		if planet == nil {
			continue
		}

		// can only refuel on docks
		if planet.Spec.DockCapacity == 0 {
			continue
		}

		planetPlayer := t.game.getPlayer(planet.PlayerNum)
		if planetPlayer.IsFriend(fleet.PlayerNum) {
			fleet.Fuel = fleet.Spec.FuelCapacity
			fleet.Spec.EstimatedRange = fleet.getEstimatedRange(player, fleet.Spec.Engine.IdealSpeed, fleet.Spec.CargoCapacity)

			t.log.Debug("fleet refueled at starbase",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Planet", planet.Name),
				slog.Int("PlanetPlayer", planet.PlayerNum),
				slog.String("Fleet", fleet.Name),
			)
		}

	}
}

// strike a random planet with a comet
func (t *turnGenerator) randomCometStrike() {
	if t.game.Year < t.game.Rules.StartingYear+t.game.Rules.RandomCometMinYear {
		// no comets in the first 10 years
		return
	}

	random := t.game.Rules.random
	chance := t.game.Rules.RandomEventChances[RandomEventComet]
	if chance == 0 || t.game.Rules.random.Float64() > chance {
		// no strike today
		return
	}

	planet := t.game.Planets[random.Intn(len(t.game.Planets))]
	if planet.Owned() && t.game.Year < t.game.Rules.StartingYear+t.game.Rules.RandomCometMinYearPlayerWorld {
		// don't hit a player world in the first 20 years
		return
	}

	// pick a random sizeIndex
	sizeIndex := random.Intn(len(CometSizes))
	size := CometSizes[sizeIndex]
	stats := t.game.Rules.CometStatsBySize[size]

	// we made it, time to tear it up
	var minerals [3]int
	var mineralConcentration [3]int
	var terraformAmount [3]int

	for i := range 3 {
		// every mineral gets a slight boost
		minerals[i] = (stats.AllMinerals + random.Intn(stats.AllRandomMinerals))

		// add a bonus to some other minerals, 1 to 3 of them depending on comet size
		if i < stats.BonusAffectsMinerals {
			minerals[i] += (stats.BonusMinerals + random.Intn(stats.BonusRandomMinerals))
			mineralConcentration[i] = stats.BonusMinConcentration + random.Intn(stats.BonusRandomConcentration)
		}
		// drop down the random numbers a bit
		minerals[i] = minerals[i] >> 4

		if i < stats.AffectsHabs {
			// terraform up or down randomly
			terraformFactor := t.game.Rules.random.Intn(2)*2 - 1
			terraformAmount[i] = (stats.MinTerraform + random.Intn(stats.RandomTerraform)) * terraformFactor
		}
	}

	// shuffle the amounts so comets don't always hit the same things
	random.Shuffle(len(minerals), func(i, j int) {
		minerals[i], minerals[j] = minerals[j], minerals[i]
		mineralConcentration[i], mineralConcentration[j] = mineralConcentration[j], mineralConcentration[i]
	})
	random.Shuffle(len(terraformAmount), func(i, j int) {
		terraformAmount[i], terraformAmount[j] = terraformAmount[j], terraformAmount[i]
	})

	mineralsAdded := Mineral{minerals[0], minerals[1], minerals[2]}
	mineralConcentrationIncreased := Mineral{mineralConcentration[0], mineralConcentration[1], mineralConcentration[2]}
	habChanged := Hab{terraformAmount[0], terraformAmount[1], terraformAmount[2]}
	colonistsKilled := 0

	planet.Cargo = planet.Cargo.AddMineral(mineralsAdded)
	planet.MineralConcentration = planet.MineralConcentration.Add(mineralConcentrationIncreased).Clamp(t.game.Rules.MinMineralConcentration, t.game.Rules.MaxMineralConcentration)
	planet.Hab = planet.Hab.Add(habChanged).Clamp(t.game.Rules.MinHab, t.game.Rules.MaxHab)
	planet.BaseHab = planet.BaseHab.Add(habChanged).Clamp(t.game.Rules.MinHab, t.game.Rules.MaxHab)
	if planet.Cargo.Colonists > 0 {
		pop := planet.GetPopulation()
		planet.Cargo.Colonists = int(float64(planet.Cargo.Colonists) * (1 - stats.PopKilledPercent))
		colonistsKilled = pop - planet.GetPopulation()
	}
	planet.MarkDirty()

	for _, player := range t.game.Players {
		messager.planetComet(player, planet, size, mineralsAdded, mineralConcentrationIncreased, habChanged, colonistsKilled)
	}

	t.log.Debug("planet struck by comet",
		slog.String("Planet", planet.Name),
		slog.Int("Player", planet.PlayerNum),
		slog.String("Size", string(size)),
		slog.String("MineralsAdded", fmt.Sprintf("%+v", mineralsAdded)),
		slog.String("MineralConcentrationIncreased", fmt.Sprintf("%+v", mineralConcentrationIncreased)),
		slog.String("HabChanged", fmt.Sprintf("%+v", habChanged)),
		slog.Int("ColonistsKilled", colonistsKilled),
	)
}

// TODO: Implement this
func (t *turnGenerator) randomMineralDeposit() {
}

// TODO: Implement this
func (t *turnGenerator) randomPlanetaryChange() {

}

func (t *turnGenerator) fleetBattle() {
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

		battler := newBattler(t.log, &t.game.Rules, battleNum, playersAtPosition, fleets, planet)

		if battler.findTargets() {
			// someone wants to fight, run the battle!
			record := battler.runBattle()

			// first discover each other and the planet
			for _, player := range playersAtPosition {

				// discover other players at the battle
				for _, otherplayer := range playersAtPosition {
					player.discoverer.discoverPlayer(otherplayer)
				}

				// discover parts of this planet's starbase
				if planet != nil {
					player.discoverer.discoverPlanet(&t.game.Rules, planet, false, false)
				}

			}

			var highestTechLevel TechLevel
			tokens := make([]ShipToken, len(record.DestroyedTokens)) // needed for component trading function call
			destroyedCost := Cost{}
			salvageOwner := 1
			for i, token := range record.DestroyedTokens {
				// figure out how much salvage this generates
				destroyedCost = destroyedCost.Add(MultiplyCost(token.design.Spec.Cost, token.Quantity))
				// TODO: who owns this salvage if there are destroyed ships from different players?
				salvageOwner = token.PlayerNum

				// record its tech level for tech trading
				highestTechLevel = highestTechLevel.Max(token.design.Spec.TechLevel)
				tokens[i].DesignNum = token.DesignNum
				tokens[i].Quantity = token.Quantity
				tokens[i].design = token.design
			}

			salvageMinerals := MultiplyCost(destroyedCost, t.game.Rules.SalvageFromBattleFactor).ToMineral()

			// every player should discover all designs in a battle as if they were penscanned.
			designsToDiscover := map[playerObject]*ShipDesign{}
			fleetsToDiscover := map[playerObject]*Fleet{}
			survivingPlayers := make(map[int]bool, len(playersAtPosition))
			for _, fleet := range fleets {
				updatedTokens := make([]ShipToken, 0, len(fleet.Tokens))
				for _, token := range fleet.Tokens {
					// add this design to our set of designs that should be discovered
					designsToDiscover[playerObjectKey(fleet.PlayerNum, token.DesignNum)] = token.design
					if token.Quantity > 0 {
						// keep this token
						updatedTokens = append(updatedTokens, token)
					}
				}
				fleet.Tokens = updatedTokens

				if len(fleet.Tokens) == 0 {
					// dead fleet, remove it
					if fleet.Starbase {
						// AR races live on starbases, so empty the planet
						player := t.game.getPlayer(fleet.PlayerNum)
						if player.Race.Spec.LivesOnStarbases {
							planet.emptyPlanet()
							messager.planetDiedOff(player, planet)
						}
						// remove this starbase from the planet
						t.game.deleteStarbase(fleet)
						planet.Starbase = nil
						planet.Spec.PlanetStarbaseSpec = computePlanetStarbaseSpec(planet)
					} else {
						t.game.deleteFleet(fleet)
					}

					// for any fleets targeting this dead fleet, update their target to the planet (or none)
					for _, otherFleet := range fleets {
						wp0 := &otherFleet.Waypoints[0]
						if wp0.TargetPlayerNum == fleet.PlayerNum && wp0.TargetNum == fleet.Num {
							wp0.clearTarget()
							if planet != nil {
								wp0.targetPlanet(planet)
							}
						}
					}
				} else {
					// this player survived
					survivingPlayers[fleet.PlayerNum] = true
					// every player discovers the remaining fleets after a battle
					fleetsToDiscover[playerObjectKey(fleet.PlayerNum, fleet.Num)] = fleet
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

			// discover all alien designs
			for _, design := range designsToDiscover {
				for _, player := range playersAtPosition {
					if player.Num != design.PlayerNum {
						player.discoverer.discoverDesign(design, true)
					}
				}
			}

			for _, fleet := range fleetsToDiscover {
				// all players discover each remaining fleet in the battle
				for _, player := range playersAtPosition {
					if fleet.PlayerNum == player.Num || fleet.Starbase {
						continue
					}
					player.discoverer.discoverFleet(fleet, false)
				}
			}

			// each player should discover starbase info after battle
			if planet != nil && planet.Starbase != nil {
				for _, player := range playersAtPosition {
					if planet.PlayerNum == player.Num {
						continue
					}
					player.discoverer.discoverPlanetStarbase(planet)
				}
			}

			if salvageMinerals.Total() > 0 {
				if planet == nil {
					t.game.getOrCreateSalvage(record.Position, salvageOwner, salvageMinerals.ToCargo())
				} else {
					planet.Cargo = planet.Cargo.AddMineral(salvageMinerals)
				}
			}

			// message each player
			for _, player := range playersAtPosition {
				// player knows about the battle
				player.BattleRecords = append(player.BattleRecords, *record)
				messager.battle(player, planet, record)

				// share battle records with allies
				for _, otherPlayer := range t.game.Players {
					if _, ok := playersAtPosition[otherPlayer.Num]; ok {
						// player is already here, no need to record the battle
						continue
					}
					if !player.IsSharingMap(otherPlayer.Num) {
						// not sharing with this player
						continue
					}

					// share the battle recording with this player
					otherPlayer.BattleRecords = append(otherPlayer.BattleRecords, *record)
					messager.battleAlly(otherPlayer, planet, record)
				}
			}

			// check for tech trades
			tt := newTechTrader()
			for playerNum, survived := range survivingPlayers {
				if !survived {
					continue // dead fleets tell no tales...
				}

				player := t.game.getPlayer(playerNum)
				field, acquiredPart := tt.checkFleetTechTrade(&t.game.Rules, player, tokens)
				if field != TechFieldNone {
					// we gained a level!
					player.techLevelGained = true
					player.TechLevels.Set(field, player.TechLevels.Get(field)+1)
					messager.playerTechGainedBattle(player, planet, record, field)
					player.updateTechsJustGained(t.game.TechStore, field)

					t.log.Debug("gained tech level from battle",
						slog.Int("Battle", battleNum),
						slog.Int("Player", player.Num),
						slog.String("field", string(field)),
					)
				}

				if acquiredPart != nil {
					player.AcquiredTechs[acquiredPart.Name] = true
					player.acquirablePartGained = true
					messager.playerAcquirablePartGainedBattle(player, planet, record, acquiredPart.Name)
					t.log.Debug("gained tech part from battle",
						slog.Int("Battle", battleNum),
						slog.Int("Player", player.Num),
						slog.String("tech", acquiredPart.Name),
					)
				}
			}

			t.log.Debug("battle finished",
				slog.Int("Battle", battleNum),
				slog.String("Players", fmt.Sprintf("%v", maps.Keys(playersAtPosition))),
			)
			battleNum++
		}
	}
}

func (t *turnGenerator) fleetBomb() {
	bomber := newBomber(t.log, &t.game.Rules)
	for _, planet := range t.game.Planets {
		if !planet.Owned() || planet.GetPopulation() == 0 || planet.Spec.HasStarbase {
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
				if fleet.Delete || fleet.OwnedBy(planetPlayer.Num) {
					continue
				}
				if fleet.Spec.Bomber && willBomb {
					enemyBombers = append(enemyBombers, fleet)
				}

				// in case our planet is destroyed by this bomber, discover any fleets in orbit
				for _, token := range fleet.Tokens {
					planetPlayer.discoverer.discoverDesign(token.design, planetPlayer.Race.Spec.DiscoverDesignOnScan)
				}
				planetPlayer.discoverer.discoverFleet(fleet, false)
			}
		}

		if len(enemyBombers) > 0 {
			// see if this planet has enemy bomber fleets, and if so, bomb it
			bomber.bombPlanet(planet, planetPlayer, enemyBombers, t.game)

			if planet.PlayerNum != planetPlayer.Num {
				// the planet was lost; discover the new (lack of an) owner and reset other intel
				planetPlayer.discoverer.clearPlanetOwnerIntel(planet)
				planetPlayer.discoverer.discoverPlanet(&t.game.Rules, planet, false, false)
			}
		}
	}
}

func (t *turnGenerator) mysteryTraderMeet() error {

	for _, mt := range t.game.MysteryTraders {

		mapObjectsAtPosition := t.game.mapObjectsByPosition[mt.Position]
		if len(mapObjectsAtPosition) <= 1 {
			continue
		}

		for _, mo := range mapObjectsAtPosition {
			if fleet, ok := mo.(*Fleet); ok && !fleet.Delete && fleet.Waypoints[0].TargetType == MapObjectTypeMysteryTrader && fleet.Waypoints[0].TargetNum == mt.Num {

				player := t.game.getPlayer(fleet.PlayerNum)

				if mt.rewardedPlayer(player.Num) {
					// the same mystery trader can't give the same player a reward twice
					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderAlreadyRewarded, mt).withSpec(PlayerMessageSpec{}.withTargetFleet(fleet)))
					continue
				}

				reward := mt.meet(&t.game.Rules, t.game.Game, fleet, player)

				if reward.Type == MysteryTraderRewardNone {
					// fleet wasn't absorbed, move on
					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderMetWithoutReward, mt).withSpec(PlayerMessageSpec{}.withTargetFleet(fleet)))
					continue
				}

				// record that this mystery trader awarded this player
				mt.PlayersRewarded[player.Num] = true

				// player got a reward
				switch reward.Type {
				case MysteryTraderRewardResearch:
					// tech levels!
					// we gained a level!
					player.techLevelGained = true
					player.TechLevels = player.TechLevels.Add(reward.TechLevels)
					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderMetWithReward, mt).
						withSpec(PlayerMessageSpec{MysteryTrader: &PlayerMessageSpecMysteryTrader{reward, 0}}.
							withTargetFleet(fleet)))

					t.log.Debug("gained tech levels from mysteryTrader",
						slog.Int("MysteryTrader", mt.Num),
						slog.Int("Player", player.Num),
						slog.String("TechLevel", fmt.Sprintf("%v", reward.TechLevels)),
					)
				case MysteryTraderRewardLifeboat:
					design := player.GetDesignByName(reward.Ship.Name)
					if design != nil && !design.MysteryTrader {
						// uh oh, the player has their own design named the same as the mystery trader, they get nothing
						t.log.Debug("player had design with same name as mysteryTrader design",
							slog.Int("MysteryTrader", mt.Num),
							slog.Int("Player", player.Num),
							slog.String("Ship", reward.Ship.Name),
						)
						player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderMetWithoutReward, mt).withSpec(PlayerMessageSpec{MysteryTrader: &PlayerMessageSpecMysteryTrader{reward, 0}}.withTargetFleet(fleet)))
						continue
					}
					if design == nil {
						// give the player a new design
						design = reward.Ship
						num := player.GetNextDesignNum(player.Designs)
						design.PlayerNum = player.Num
						design.Num = num
						design.GameDBObject = GameDBObject{}

						var err error
						design.Spec, err = ComputeShipDesignSpec(&t.game.Rules, player.TechLevels, player.Race.Spec, design)
						if err != nil {
							return fmt.Errorf("ComputeShipDesignSpec returned error: %w", err)
						}
						player.Designs = append(player.Designs, design)
						t.game.addDesign(design)
					}
					token := ShipToken{design: design, DesignNum: design.Num, Quantity: reward.ShipCount}
					design.Spec.NumInstances += token.Quantity
					design.Spec.NumBuilt += token.Quantity

					rewardFleet, err := t.addFleet(player, fleet.Position, token, Tags{})
					if err != nil {
						return err
					}

					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderMetWithReward, mt).withSpec(PlayerMessageSpec{MysteryTrader: &PlayerMessageSpecMysteryTrader{reward, rewardFleet.Num}}.withTargetFleet(fleet)))

					t.log.Debug("new fleet created from mysteryTrader reward",
						slog.Int("Player", rewardFleet.PlayerNum),
						slog.String("Position", rewardFleet.Position.String()),
						slog.String("Fleet", rewardFleet.Name),
						slog.Int("ShipCount", reward.ShipCount),
					)
				default:
					// MT awarded tech
					mtTech := t.game.GetTech(reward.Tech)
					if mtTech == nil {
						return fmt.Errorf("mystery trader %d awarded unknown tech %s", mt.Num, reward.Tech)
					}
					if _, ok := player.AcquiredTechs[reward.Tech]; ok {
						t.log.Warn("player already has tech awarded by mysteryTrader",
							slog.Int("MysteryTrader", mt.Num),
							slog.Int("Player", player.Num),
							slog.String("Tech", reward.Tech),
						)
					}
					player.AcquiredTechs[reward.Tech] = true
					player.Messages = append(player.Messages, newMysteryTraderMessage(PlayerMessageMysteryTraderMetWithReward, mt).withSpec(PlayerMessageSpec{MysteryTrader: &PlayerMessageSpecMysteryTrader{reward, 0}}.withTargetFleet(fleet)))
				}

				// remove the absorbed fleet from the universe
				t.game.deleteFleet(fleet)
			}
		}
	}
	return nil
}

// decay Minefields and remove any minefields that are too small
func (t *turnGenerator) decayMines() {
	for _, minefield := range t.game.Minefields {
		player := t.game.getPlayer(minefield.PlayerNum)
		decayRate := minefield.getDecayRate(&t.game.Rules, player, NumMapObjectsWithin(t.game.Planets, minefield.Position, minefield.Radius()))
		minefield.NumMines -= decayRate
		if minefield.NumMines <= 10 {
			t.game.deleteMinefield(minefield)
			continue
		}

		t.log.Debug("decayed minefield",
			slog.Int("Player", minefield.PlayerNum),
			slog.String("Minefield", minefield.Name),
			slog.Int("NumMines", minefield.NumMines),
		)
	}
}

func (t *turnGenerator) fleetLayMines() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp0 := fleet.Waypoints[0]
		if wp0.Task == WaypointTaskLayMinefield {
			player := t.game.getPlayer(fleet.PlayerNum)

			if !fleet.Spec.CanLayMines {
				messager.fleetMinesLaidFailed(player, fleet)
				continue
			}

			for mineType, minesLaid := range fleet.Spec.MineLayingRateByMineType {
				if len(fleet.Waypoints) > 1 {
					minesLaid = int(float64(minesLaid) * player.Race.Spec.MinefieldRateMoveFactor)
				}

				// We aren't laying mines (probably because we're moving, skip it)
				if minesLaid == 0 {
					continue
				}

				// See if we are adding to an existing minefield
				minefield := t.game.getMinefieldNearPosition(player.Num, fleet.Position, mineType)
				if minefield == nil {
					minefield = newMinefield(player, mineType, minesLaid, t.game.getNextMinefieldNum(), fleet.Position)
					t.game.addMinefield(minefield)
				} else {
					// Add to it!
					minefield.NumMines += minesLaid
				}

				messager.fleetMinesLaid(player, fleet, minefield, minesLaid)

				if minefield.Position != fleet.Position {
					// Move this minefield closer to us (in case it's not in our location)
					// This was taken from the FreeStars codebase (like many other things)
					minefield.moveTowardsMineLayer(fleet.Position, minesLaid)
				}

				t.log.Debug("laid mines",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.Int("MinesLaid", minesLaid),
					slog.String("Minefield", minefield.Name),
					slog.Int("NumMines", minefield.NumMines),
				)
			}
		}
	}

}

// process transfer fleet orders to gift fleets to other players
func (t *turnGenerator) fleetTransferOwner() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		wp0 := &fleet.Waypoints[0]
		if wp0.Task == WaypointTaskTransferFleet {
			player := t.game.getPlayer(fleet.PlayerNum)
			targetPlayer := t.game.getPlayer(wp0.TransferToPlayer)

			if targetPlayer == nil {
				// can't find target player
				messager.fleetTransferInvalidPlayer(player, fleet)
				t.log.Error("tried to transfer fleet but target player doesn't exist",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.Int("TargetPlayer", wp0.TargetPlayerNum),
				)
				wp0.Task = WaypointTaskNone
				wp0.TransferToPlayer = None
				continue
			}

			if fleet.Cargo.Colonists > 0 {
				// can't give colonists
				messager.fleetTransferInvalidColonists(player, fleet, targetPlayer)
				t.log.Debug("transferring fleet failed, fleet has colonists",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
				)
				wp0.Task = WaypointTaskNone
				wp0.TransferToPlayer = None
				continue
			}

			if targetPlayer == player {
				t.log.Error("tried to transfer fleet to self",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
				)
				wp0.Task = WaypointTaskNone
				wp0.TransferToPlayer = None
				continue
			}

			if !targetPlayer.IsFriend(player.Num) {
				// they are not allies, they will refuse the offer
				messager.fleetTransferInvalidGiveRefused(player, fleet, targetPlayer)
				messager.fleetTransferInvalidReceiveRefused(targetPlayer, fleet, player)
				t.log.Debug("transferring fleet refused",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.Int("TargetPlayer", targetPlayer.Num),
				)
				wp0.Task = WaypointTaskNone
				wp0.TransferToPlayer = None

				continue
			}

			// give the gift of this fleet!
			t.log.Debug("transferring fleet",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
				slog.Int("TargetPlayer", targetPlayer.Num),
			)
			for i := range fleet.Tokens {
				token := &fleet.Tokens[i]
				design := token.design

				// give the player a copy of this design
				newName := fmt.Sprintf("%s %s", player.Race.PluralName, design.Name)
				targetPlayerDesign := targetPlayer.GetDesignByName(newName)
				if targetPlayerDesign != nil {
					if !targetPlayerDesign.SlotsEqual(design.Slots) {
						// uh oh, design has been updated since the last time it was transferred to us...
						// create a new design for the target player
						num := targetPlayer.GetNextDesignNum(targetPlayer.Designs)
						newDesign := *design
						newDesign.GameDBObject = GameDBObject{}
						newDesign.OriginalPlayerNum = player.Num
						newDesign.PlayerNum = targetPlayer.Num
						newDesign.Num = num
						// rev the version and append it to the name
						newDesign.Version++
						newDesign.Name = fmt.Sprintf("%s v%d", newName, newDesign.Version)
						targetPlayerDesign = &newDesign
						targetPlayer.Designs = append(targetPlayer.Designs, targetPlayerDesign)
					}
				} else {
					// create a new design for the target player
					num := targetPlayer.GetNextDesignNum(targetPlayer.Designs)
					newDesign := *design
					newDesign.GameDBObject = GameDBObject{}
					newDesign.Name = newName
					newDesign.OriginalPlayerNum = player.Num
					newDesign.PlayerNum = targetPlayer.Num
					newDesign.Num = num
					targetPlayerDesign = &newDesign
					targetPlayer.Designs = append(targetPlayer.Designs, targetPlayerDesign)
				}

				// make sure we don't update this spec
				targetPlayerDesign.Spec.NumBuilt = 0
				targetPlayerDesign.Spec.NumInstances = 0

				token.design = targetPlayerDesign
				token.DesignNum = targetPlayerDesign.Num
			}

			playerFleets := t.game.getFleets(targetPlayer.Num)
			fleet.Num = targetPlayer.GetNextFleetNum(playerFleets)
			fleet.PlayerNum = targetPlayer.Num

			// clear out the waypoints
			wp0.Task = WaypointTaskNone
			wp0.TransferToPlayer = None
			fleet.Waypoints = fleet.Waypoints[:1]

			// notify the player here (before we give it away and change the name)
			messager.fleetTransferGiven(player, fleet, targetPlayer)
			messager.fleetTransferReceived(targetPlayer, fleet, player)

			fleet.Rename(fmt.Sprintf("%s %s", player.Race.PluralName, fleet.BaseName))

		}
	}
}

func (t *turnGenerator) instaform() {
	for _, planet := range t.game.Planets {
		if planet.Owned() {
			player := t.game.getPlayer(planet.PlayerNum)
			if player.Race.Spec.Instaforming {
				// find out how much our instaform would terraform this planet from base
				terraformer := NewTerraformer()
				instaformAmount := terraformer.GetTerraformAmount(planet.BaseHab, planet.BaseHab, player, player)
				newHab := planet.BaseHab.Add(instaformAmount)

				// see if we would change this planet's hab
				if newHab != planet.Hab {
					// Instantly terraform this planet (but don't update planet.TerraformAmount, this change doesn't stick if we leave)
					prevHab := planet.Hab
					planet.Hab = newHab
					planet.Spec = ComputePlanetSpec(&t.game.Rules, player, planet)
					messager.planetInstaform(player, planet, instaformAmount)

					t.log.Debug("instaformed planet",
						slog.Int("Player", player.Num),
						slog.String("Planet", planet.Name),
						slog.String("PreviousHab", prevHab.String()),
						slog.String("Hab", planet.Hab.String()),
					)
				}
			}
		}
	}
}

func (t *turnGenerator) fleetSweepMines() {

	// fleets and starbases sweep
	for _, fleet := range append(t.game.Fleets, t.game.Starbases...) {
		if !fleet.Delete && fleet.Spec.MineSweep > 0 {
			fleetPlayer := t.game.getPlayer(fleet.PlayerNum)
			for _, minefield := range t.game.Minefields {
				// don't sweep dead fields
				if minefield.Delete {
					continue
				}

				// sweep mines
				if fleet.willAttack(fleetPlayer, minefield.PlayerNum) && isPointInCircle(fleet.Position, minefield.Position, minefield.Radius()) {
					minefieldPlayer := t.game.getPlayer(minefield.PlayerNum)
					numSwept := minefield.sweep(&t.game.Rules, fleet.Position, fleet.Spec.MineSweep)

					if numSwept == 0 {
						t.log.Debug("no mines swept",
							slog.Int("Player", fleet.PlayerNum),
							slog.String("Fleet", fleet.Name),
							slog.String("Minefield", minefield.Name),
							slog.Int("MinefieldPlayer", minefield.PlayerNum),
							slog.Int("NumMines", minefield.NumMines),
						)
						continue
					}

					messager.fleetMinefieldSwept(fleetPlayer, fleet, minefield, numSwept)
					messager.fleetMinefieldSwept(minefieldPlayer, fleet, minefield, numSwept)

					t.log.Debug("fleet swept mines",
						slog.Int("Player", fleet.PlayerNum),
						slog.String("Fleet", fleet.Name),
						slog.String("Minefield", minefield.Name),
						slog.Int("MinefieldPlayer", minefield.PlayerNum),
						slog.Int("NumMines", minefield.NumMines),
					)
					if minefield.NumMines <= 10 {
						t.game.deleteMinefield(minefield)
						continue
					}
				}
			}
		}
	}

}

// repair fleets and starbases
func (t *turnGenerator) fleetRepair() {
	for _, fleet := range t.game.Fleets {
		if fleet.Delete {
			continue
		}

		player := t.game.getPlayer(fleet.PlayerNum)
		orbiting := t.game.getOrbitingPlanet(fleet)
		fleet.repairFleet(t.log, &t.game.Rules, player, orbiting)
	}

	for _, starbase := range t.game.Starbases {
		if starbase.Delete {
			continue
		}

		if starbase.Tokens[0].QuantityDamaged == 0 {
			continue
		}

		player := t.game.getPlayer(starbase.PlayerNum)
		starbase.repairStarbase(t.log, &t.game.Rules, player)
	}
}

func (t *turnGenerator) fleetRemoteTerraform() {
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
				t.log.Debug("fleet remote terraformed planet",
					slog.Int("Player", fleet.PlayerNum),
					slog.String("Fleet", fleet.Name),
					slog.String("Planet", planet.Name),
					slog.String("HabType", result.Type.String()),
				)
			}
		}
	}
}

func (t *turnGenerator) fleetPatrol(player *Player) {
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

		rangeDistanceSquared := float64(wp.PatrolRange * wp.PatrolRange)
		if wp.PatrolRange == PatrolRangeInfinite {
			rangeDistanceSquared = math.MaxFloat64
		}

		closestDistance := float64(math.MaxFloat32)
		var closest *Fleet

		for _, enemyFleet := range player.FleetIntels {
			if fleet.willAttack(player, enemyFleet.PlayerNum) {
				distSquaredToFleet := fleet.Position.DistanceSquaredTo(enemyFleet.Position)
				if distSquaredToFleet <= rangeDistanceSquared {
					if distSquaredToFleet < closestDistance {
						closestDistance = distSquaredToFleet
						closest = enemyFleet
					}
				}
			}
		}

		if closest != nil {

			if wp.PatrolWarpSpeed == PatrolWarpSpeedAutomatic {
				wp.PatrolWarpSpeed = fleet.Spec.Engine.IdealSpeed
			}

			// add a waypoint to the fleet
			wpTarget := NewFleetWaypoint(closest.Position, closest.Num, closest.PlayerNum, closest.Name, wp.PatrolWarpSpeed)
			// this is an ephemeral waypoint that isn't ever "completed" and so shouldn't be repeated. This should probably be
			// named differently...
			wpTarget.PartiallyComplete = true

			// for fleets that do Patrol + repeat orders, we let them intercept the fleet
			// and head back to base to patrol again
			// if they aren't repeating orders, we assume they just want to keep patroling
			// and auto intercepting the closest fleet they will attack. Roaming the universe
			// for all time.
			if !fleet.RepeatOrders {
				wpTarget.Task = WaypointTaskPatrol
				wpTarget.PatrolRange = wp.PatrolRange
				wpTarget.PatrolWarpSpeed = wp.PatrolWarpSpeed
				wpTarget.PartiallyComplete = false
			}

			fleet.Waypoints = append(fleet.Waypoints, wpTarget)

			messager.fleetPatrolTargeted(player, fleet, closest)

			t.log.Debug("fleet patrol targeted enemy",
				slog.Int("Player", fleet.PlayerNum),
				slog.String("Fleet", fleet.Name),
				slog.String("Target", closest.Name),
				slog.Int("TargetPlayer", closest.PlayerNum),
			)
		}
	}
}

func (t *turnGenerator) scan() error {
	for _, player := range t.game.Players {
		player.Spec = ComputePlayerSpec(player, &t.game.Rules)

		scanner := newPlayerScanner(t.game.Universe, t.game.Players, &t.game.Rules, player)
		if err := scanner.scan(); err != nil {
			return fmt.Errorf("scan universe and update player intel failed: %w", err)
		}
		t.fleetPatrol(player)

		player.SubmittedTurn = false
	}

	return nil
}

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

// Calculate the score for this year for each player.
func (t *turnGenerator) calculateScores() {
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
			score.Score += int(min(float64(planet.GetPopulation()/100000), 6))
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
			achievedLevel := player.TechLevels.Get(field)
			score.TechLevels += achievedLevel
			for level := 0; level <= achievedLevel; level++ {
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
		}

		// Calculate additional score components
		// Resources: 1 point for every 30 resources
		score.Score += score.Resources / 30
		// Starbases: 3 points each (doesn't include Orbital Forts)
		score.Score += score.Starbases * 3
		// Unarmed Ships: You receive 1/2 point for each unarmed ship (up to the number of planets you own).
		score.Score += int(min(float64(score.UnarmedShips)*0.5+5, float64(score.Planets)))
		// Escort Ships: You receive 2 points for each Escort ship (up to the number of planets you own).
		score.Score += int(min(float64(score.EscortShips)*2, float64(score.Planets)))
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
	slices.SortFunc(scoreSortedPlayers, func(p1, p2 *Player) int {
		return p2.GetScore().Score - p1.GetScore().Score
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
			discoverer := player.discoverer
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

func (t *turnGenerator) checkBattleReports() {
	for _, player := range t.game.Players {
		if len(player.BattleRecords) == 0 {
			continue
		}
		// notify this player there are battle reports
		messager.battleReports(player)
	}

}

// check if this player is victorious, and if so, notify everyone
func (t *turnGenerator) checkVictory(player *Player) {
	victoryChecker := newVictoryChecker(t.game)
	for _, player := range t.game.Players {
		if err := victoryChecker.checkForVictor(player); err != nil {
			t.log.Error("error while checking for victory", slog.Any("err", err))
			return
		}
	}

	// we don't declare a victor until some time has passed
	if t.game.YearsPassed() >= t.game.VictoryConditions.YearsPassed && t.game.VictorDeclared {

		// if we won, tell everyone about it!
		if player.Victor {
			t.log.Debug("you are victorious your majesty!",
				slog.Int("Player", player.Num),
				slog.String("PlayerName", player.Name),
				slog.String("Race", player.Race.PluralName),
			)
			for _, p := range t.game.Players {
				messager.playerVictory(p, player)
			}
		}
	}
}

func (t *turnGenerator) checkDeath() {
	for _, player := range t.game.Players {

		numPlanets := 0
		numFleets := 0
		numColonists := 0
		for _, planet := range t.game.Planets {
			if planet.PlayerNum == player.Num {
				numPlanets++
				numColonists += planet.GetPopulation()
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
			// let everyone know this player died
			for _, otherPlayer := range t.game.Players {
				messager.playerDead(otherPlayer, player)
			}
			t.log.Debug("player is dead",
				slog.Int("Player", player.Num),
				slog.String("PlayerName", player.Name),
				slog.String("Race", player.Race.PluralName),
			)
		} else if numPlanets == 0 && numFleets > 0 {
			messager.playerNoPlanets(player, numColonists)
			t.log.Debug("player has no planets",
				slog.Int("Player", player.Num),
				slog.String("PlayerName", player.Name),
				slog.String("Race", player.Race.PluralName),
				slog.Int("NumColonists", numColonists),
			)
		}
	}

}
