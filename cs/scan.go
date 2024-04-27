package cs

import (
	"math"
	"slices"

	"golang.org/x/exp/maps"
)

// The scanner is used at the end of the turn generation to update player intels
// with their knowledge of the universe. It handles scanning planets, fleets, minefields, etc
type scanner struct {
	Position            Vector
	RangeSquared        int
	RangePenSquared     int
	DiscoverFleetCargo  bool
	DiscoverPlanetCargo bool
	CloakReduction      float64
}

type playerScan struct {
	universe          *Universe
	rules             *Rules
	player            *Player
	players           []*Player
	discoveredPlayers map[int]bool
	discoverer        discoverer
}

type playerScanner interface {
	scan() error
}

func newPlayerScanner(universe *Universe, players []*Player, rules *Rules, player *Player) playerScanner {
	discoverer := newDiscoverer(player)
	return &playerScan{universe, rules, player, players, make(map[int]bool, len(player.PlayerIntels.PlayerIntels)), discoverer}
}

// scan planets, fleets, etc for a player
func (scan *playerScan) scan() error {
	// clear out any reports that we recreate each year
	player := scan.player
	scan.discoverer.clearTransientReports()

	for i := range player.PlanetIntels {
		planet := &player.PlanetIntels[i]
		if planet.ReportAge != ReportAgeUnexplored {
			planet.ReportAge++
		}
	}

	for i := range player.WormholeIntels {
		wormhole := &player.WormholeIntels[i]
		if wormhole.ReportAge != ReportAgeUnexplored {
			wormhole.ReportAge++
		}
	}

	for i := range player.MineFieldIntels {
		mineField := &player.MineFieldIntels[i]
		if mineField.ReportAge != ReportAgeUnexplored {
			mineField.ReportAge++
		}
	}

	scanners := scan.getScanners()
	remoteMiningScanners := scan.getRemoteMiningScanners()
	cargoScanners := scan.getCargoScanners()

	// scan planets
	if err := scan.scanPlanets(scanners, append(cargoScanners, remoteMiningScanners...)); err != nil {
		return err
	}

	// scan universe
	scan.scanFleets(scanners, cargoScanners)
	scan.scanMineFields(scanners)
	scan.scanMineralPackets(scanners)
	scan.scanSalvages(scanners)
	scan.scanWormholes(scanners)
	scan.scanMysteryTraders(scanners)

	scan.discoverPlayers()

	scan.updateFleetTargets()

	return nil
}

// scan all planets with this player's scanners
func (scan *playerScan) scanPlanets(scanners []scanner, cargoScanners []scanner) error {
	for _, planet := range scan.universe.Planets {
		if planet.OwnedBy(scan.player.Num) {
			if err := scan.discoverer.discoverPlanet(scan.rules, scan.player, planet, false); err != nil {
				return err
			}
			continue
		}

		// try and scan the planet with this scanner
		for _, scanner := range scanners {
			scanned, err := scan.scanPlanet(planet, scanner)
			if err != nil {
				return err
			}
			if scanned {
				break
			}
		}

		// try and scan the planet with a cargo scanner
		for _, scanner := range cargoScanners {
			scanned, err := scan.scanPlanet(planet, scanner)
			if err != nil {
				return err
			}
			if scanned {
				break
			}
		}

		intel := scan.discoverer.getPlanetIntel(planet.Num)
		if intel.ReportAge != ReportAgeUnexplored {
			// TODO: remove this after initial test games are done
			// it's just here because some old games don't have basehab on intels
			if intel.BaseHab != planet.BaseHab {
				intel.BaseHab = planet.BaseHab
			}
			scan.discoverer.discoverPlanetTerraformability(scan.player, planet.Num)
		}
	}

	return nil
}

// scan this planet
func (scan *playerScan) scanPlanet(planet *Planet, scanner scanner) (bool, error) {
	if float64(scanner.RangePenSquared) >= scanner.Position.DistanceSquaredTo(planet.Position) {
		if planet.Owned() {
			scan.discoveredPlayers[planet.PlayerNum] = true
		}
		if err := scan.discoverer.discoverPlanet(scan.rules, scan.player, planet, true); err != nil {
			return false, err
		}
		if scanner.DiscoverPlanetCargo {
			if err := scan.discoverer.discoverPlanetCargo(scan.player, planet); err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanFleets(scanners []scanner, cargoScanners []scanner) {
	// scan fleets
	fleetsToScan := []*Fleet{}
	fleetsToCargoScan := []*Fleet{}
	for _, fleet := range scan.universe.Fleets {
		// skip deleted fleets
		if fleet.Delete {
			continue
		}
		if fleet.OwnedBy(scan.player.Num) {
			// The player already gets a copy of all their own fleets
			continue
		}

		// try and scan the planet with this scanner
		for _, scanner := range scanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToScan = append(fleetsToScan, fleet)
				break
			}
		}

		// try and scan the planet with a cargo scanner
		for _, scanner := range cargoScanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToCargoScan = append(fleetsToScan, fleet)
				break
			}
		}
	}

	for _, fleet := range fleetsToScan {
		scan.discoveredPlayers[fleet.PlayerNum] = true

		scan.discoverer.discoverFleet(scan.player, fleet)
		for _, token := range fleet.Tokens {
			scan.discoverer.discoverDesign(scan.player, token.design, scan.player.Race.Spec.DiscoverDesignOnScan)
		}
	}

	for _, fleet := range fleetsToCargoScan {
		scan.discoverer.discoverFleetCargo(scan.player, fleet)
	}
}

// return true if this scanner successfully scans this fleet, taking into account cloaking
// and the fleet's cloak penetration
func (scan *playerScan) fleetInScannerRange(fleet *Fleet, scanner scanner) bool {
	var cloakFactor = 1 - (float64(fleet.Spec.CloakPercent) * (1 - scanner.CloakReduction) / 100.0)
	var distance = scanner.Position.DistanceSquaredTo(fleet.Position)

	// if we pen scanned this, update the report
	if float64(scanner.RangePenSquared)*cloakFactor >= distance {
		// update the fleet report with pen scanners
		return true
	}

	// if we aren't orbiting a planet, we can be seen with regular scanners
	if !fleet.Orbiting() && float64(scanner.RangeSquared)*cloakFactor >= distance {
		return true
	}
	return false
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanWormholes(scanners []scanner) {
	for _, wormhole := range scan.universe.Wormholes {
		intel := scan.discoverer.getWormholeIntel(wormhole.Num)
		cloakFactor := 1.0 - (float64(scan.rules.WormholeCloak) / 100)
		if intel != nil {
			cloakFactor = 1
		}

		for _, scanner := range scanners {
			if cloakFactor != 1 {
				// calculate cloak reduction for tachyon detectors if this wormhole is cloaked
				cloakFactor = 1 - (1-cloakFactor)*scanner.CloakReduction
			}

			// we only care about regular scanners for wormholes
			if float64(scanner.RangeSquared)*cloakFactor >= scanner.Position.DistanceSquaredTo(wormhole.Position) {
				if wormhole.Delete {
					// this wormhole went away, rmeove it from intel
					scan.discoverer.forgetWormhole(wormhole.Num)
				} else {
					scan.discoverer.discoverWormhole(wormhole)
				}

				break
			}
		}
	}

	intels := make([]WormholeIntel, len(scan.player.WormholeIntels))
	copy(intels, scan.player.WormholeIntels)
	for _, intel := range intels {
		for _, scanner := range scanners {
			// if we scanned this wormhole where we last saw it, but it no longer exists in the universe, forget it
			if float64(scanner.RangeSquared) >= scanner.Position.DistanceSquaredTo(intel.Position) {
				wormhole := scan.universe.getWormhole(intel.Num)
				if wormhole == nil || wormhole.Delete {
					// this wormhole went away, rmeove it from intel
					scan.discoverer.forgetWormhole(intel.Num)
				}
				break
			}
		}
	}
}

// scan Mystery Traders
func (scan *playerScan) scanMysteryTraders(scanners []scanner) {
	for _, mysteryTrader := range scan.universe.MysteryTraders {
		if mysteryTrader.Delete {
			continue
		}
		for _, scanner := range scanners {
			// we only care about regular scanners for mysteryTraders
			if float64(scanner.RangeSquared) >= scanner.Position.DistanceSquaredTo(mysteryTrader.Position) {
				scan.discoverer.discoverMysteryTrader(mysteryTrader)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanMineralPackets(scanners []scanner) {
	for _, packet := range scan.universe.MineralPackets {
		if packet.Delete {
			continue
		}
		// skip our own
		if scan.player.Num == packet.PlayerNum {
			continue
		}

		target := scan.universe.getPlanet(packet.TargetPlanetNum)
		packetPlayer := scan.players[packet.PlayerNum-1]

		// PP races detect all packets in flight
		if scan.player.Race.Spec.DetectAllPackets {
			scan.discoverer.discoverMineralPacket(scan.rules, scan.player, packet, packetPlayer, target)
			continue
		}

		for _, scanner := range scanners {
			// we only care about regular scanners for mineral packets
			if float64(scanner.RangeSquared) >= scanner.Position.DistanceSquaredTo(packet.Position) {
				scan.discoverer.discoverMineralPacket(scan.rules, scan.player, packet, packetPlayer, target)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanMineFields(scanners []scanner) {
	for _, mineField := range scan.universe.MineFields {
		if mineField.Delete {
			continue
		}

		if mineField.OwnedBy(scan.player.Num) {
			// The player already gets a copy of all their own mineFields
			continue
		}
		intel := scan.discoverer.getMineFieldIntel(mineField.PlayerNum, mineField.Num)

		cloakFactor := 1.0 - (float64(scan.rules.MineFieldCloak) / 100)
		if intel != nil {
			cloakFactor = 1
		}

		for _, scanner := range scanners {
			if cloakFactor != 1 {
				// calculate cloak reduction for tachyon detectors if this minefield is cloaked
				cloakFactor = 1 - (1-cloakFactor)*scanner.CloakReduction
			}

			// we only care about regular scanners for wormholes
			if float64(scanner.RangeSquared)*cloakFactor >= scanner.Position.DistanceSquaredTo(mineField.Position) {
				scan.discoverer.discoverMineField(scan.player, mineField)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanSalvages(scanners []scanner) {
	for _, salvage := range scan.universe.Salvages {
		if salvage.Delete {
			continue
		}
		for _, scanner := range scanners {
			// we only care about regular scanners for mineral packets
			if float64(scanner.RangeSquared) >= scanner.Position.DistanceSquaredTo(salvage.Position) {
				scan.discoverer.discoverSalvage(salvage)
				break
			}
		}
	}
}

func (scan *playerScan) discoverPlayers() {
	for player, discovered := range scan.discoveredPlayers {
		if discovered {
			scan.discoverer.discoverPlayer(scan.players[player-1])
		}
	}
}

// get a list of unique scanners per player.
// This is a minimal list only containing the best scanner values for each position
func (scan *playerScan) getScanners() []scanner {
	scanningFleetsByPosition := map[Vector]scanner{}
	for _, fleet := range scan.universe.Fleets {
		if fleet.Delete {
			continue
		}
		if fleet.PlayerNum == scan.player.Num && fleet.Spec.Scanner {
			scanner, found := scanningFleetsByPosition[fleet.Position]
			if !found {
				// start with NoScanner (-1)
				scanner.Position = fleet.Position
				scanner.RangeSquared = NoScanner
				scanner.RangePenSquared = NoScanner
			}
			scanner.RangeSquared = MaxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = MaxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
			scanner.CloakReduction = math.Max(scanner.CloakReduction, fleet.Spec.ReduceCloaking)
			scanningFleetsByPosition[fleet.Position] = scanner
		}
	}

	// build a list of scanners for this player
	scanners := []scanner{}
	for _, planet := range scan.universe.Planets {
		if planet.PlayerNum == scan.player.Num && planet.Scanner {
			scanner := scanner{
				Position:        planet.Position,
				RangeSquared:    planet.Spec.ScanRange * planet.Spec.ScanRange,
				RangePenSquared: planet.Spec.ScanRangePen * planet.Spec.ScanRangePen,
			}

			// use the fleet scanner if it's better
			if fleetScanner, ok := scanningFleetsByPosition[planet.Position]; ok {
				scanner.RangeSquared = MaxInt(scanner.RangeSquared, fleetScanner.RangeSquared)
				scanner.RangePenSquared = MaxInt(scanner.RangePenSquared, fleetScanner.RangePenSquared)
				scanner.CloakReduction = math.Max(scanner.CloakReduction, fleetScanner.CloakReduction)
			}
			scanners = append(scanners, scanner)
		}
	}

	// Space demolition minefields act as scanners
	if scan.player.Race.Spec.MineFieldsAreScanners {
		for _, mineField := range scan.universe.MineFields {
			if mineField.PlayerNum == scan.player.Num {
				scanner := scanner{
					Position:     mineField.Position,
					RangeSquared: int(mineField.Spec.Radius),
				}
				scanners = append(scanners, scanner)
			}
		}
	}

	// Packet Physics packets act as pen scanners
	for _, packet := range scan.universe.MineralPackets {
		if packet.PlayerNum == scan.player.Num && (packet.ScanRange != NoScanner || packet.ScanRangePen != NoScanner) {
			scanner := scanner{
				Position:        packet.Position,
				RangeSquared:    packet.ScanRange * packet.ScanRange,
				RangePenSquared: packet.ScanRangePen * packet.ScanRangePen,
			}
			scanners = append(scanners, scanner)
		}
	}

	// add in any fleet scanners that weren't on a planet
	if len(scanners) == 0 {
		// we have no planetary scanners (weird, but possible if all planets with scanners are destroyed)
		// so just add the fleet scanners to the list
		for _, fleetScanner := range scanningFleetsByPosition {
			scanners = append(scanners, fleetScanner)
		}
	} else {
		scannersByPosition := map[Vector]scanner{}
		for _, scanner := range scanners {
			scannersByPosition[scanner.Position] = scanner
		}
		for position, fleetScanner := range scanningFleetsByPosition {
			// if we don't find a scanner at this position, add the fleetScanner
			// to our existing scanners list
			if _, found := scannersByPosition[position]; !found {
				scanners = append(scanners, fleetScanner)
			}
		}
	}

	return scanners
}

// get a list of remote mining scanners by player
func (scan *playerScan) getRemoteMiningScanners() []scanner {
	scanningFleetsByPosition := map[Vector]scanner{}
	for _, fleet := range scan.universe.Fleets {
		if fleet.Delete {
			continue
		}

		// find any fleets that remote mined this turn, but only add one per position
		if fleet.PlayerNum == scan.player.Num && fleet.remoteMined {
			if scanner, found := scanningFleetsByPosition[fleet.Position]; !found {
				scanner.Position = fleet.Position
				scanner.RangeSquared = 0
				scanner.RangePenSquared = 0
				scanner.DiscoverPlanetCargo = true
				scanningFleetsByPosition[fleet.Position] = scanner
			}
		}
	}

	return maps.Values(scanningFleetsByPosition)
}

// get a list of scanners that can scan cargo from fleets or planets
func (scan *playerScan) getCargoScanners() []scanner {
	scanners := []scanner{}
	scanningFleetsByPosition := map[Vector]scanner{}

	for _, fleet := range scan.universe.Fleets {
		if fleet.Delete {
			continue
		}

		if fleet.PlayerNum == scan.player.Num && fleet.Spec.Scanner && (fleet.Spec.CanStealFleetCargo || fleet.Spec.CanStealPlanetCargo) {
			scanner, found := scanningFleetsByPosition[fleet.Position]
			if !found {
				// start with NoScanner (-1)
				scanner.Position = fleet.Position
				scanner.RangeSquared = NoScanner
				scanner.RangePenSquared = NoScanner
			}
			scanner.RangeSquared = MaxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = MaxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
			scanner.CloakReduction = math.Max(scanner.CloakReduction, fleet.Spec.ReduceCloaking)
			scanner.DiscoverFleetCargo = fleet.Spec.CanStealFleetCargo
			scanner.DiscoverPlanetCargo = fleet.Spec.CanStealPlanetCargo
			scanningFleetsByPosition[fleet.Position] = scanner
		}
	}

	for _, fleetScanner := range scanningFleetsByPosition {
		scanners = append(scanners, fleetScanner)
	}

	return scanners
}

// make sure our fleets are pointing to valid targets
func (scan *playerScan) updateFleetTargets() {
	for _, fleet := range scan.universe.Fleets {
		// skip deleted fleets
		if fleet.Delete {
			continue
		}
		if !fleet.OwnedBy(scan.player.Num) {
			// Skip fleets we don't own
			continue
		}

		if len(fleet.Waypoints) == 1 {
			wp0 := fleet.Waypoints[0]
			if fleet.OrbitingPlanetNum == None && wp0.TargetType != MapObjectTypeNone {
				// we arrived at our target, but it's not a planet. Keep it as wp1
				fleet.Waypoints = []Waypoint{NewPositionWaypoint(fleet.Position, fleet.WarpSpeed), wp0}
			} else {
				fleet.WarpSpeed = 0
				fleet.Heading = Vector{}
			}
		}

		for i := 1; i < len(fleet.Waypoints); i++ {
			wp := &fleet.Waypoints[i]

			// none and planet targets always work
			if wp.TargetType == MapObjectTypeNone || wp.TargetType == MapObjectTypePlanet {
				continue
			}

			if wp.TargetPlayerNum == scan.player.Num {
				// we own this and won't have intel for it
				mo := scan.universe.getMapObject(wp.TargetType, wp.TargetNum, wp.TargetPlayerNum)
				if mo == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				}
				continue
			}

			switch wp.TargetType {
			case MapObjectTypeFleet:
				target := scan.discoverer.getFleetIntel(wp.TargetPlayerNum, wp.TargetNum)
				if target == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				} else {
					// fleets move, make sure our position updates
					wp.Position = target.Position
				}

			case MapObjectTypeSalvage:
				target := scan.discoverer.getSalvageIntel(wp.TargetNum)
				if target == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				}

			case MapObjectTypeMineralPacket:
				target := scan.discoverer.getMineralPacketIntel(wp.TargetPlayerNum, wp.TargetNum)
				if target == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				} else {
					// fleets move, make sure our position updates
					wp.Position = target.Position
				}

			case MapObjectTypeWormhole:
				target := scan.discoverer.getWormholeIntel(wp.TargetNum)
				if target == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				} else {
					// wormholes move, make sure our position updates
					wp.Position = target.Position
				}
			}

			// this waypoint is now the same as the one before it, so delete it
			if wp.TargetType == MapObjectTypeNone && wp.Position == fleet.Waypoints[i-1].Position {
				fleet.Waypoints = slices.Delete(fleet.Waypoints, i, i+1)
				i--
			}
		}

	}
}
