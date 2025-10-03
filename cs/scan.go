package cs

import (
	"math"
	"slices"

	"golang.org/x/exp/maps"
)

// The scanner is used at the end of the turn generation to update player intels
// with their knowledge of the universe. It handles scanning planets, fleets, minefields, etc
type scanner struct {
	Position             Vector
	Range                int
	RangePen             int
	DiscoverFleetCargo   bool
	DiscoverPlanetCargo  bool
	CloakReductionFactor float64
}

const NoCloakFactor = 1

// RangeSquared returns the cloak adjusted RangeSquared value
func (s scanner) RangeSquared(cloakFactor float64) int {
	if cloakFactor == NoCloakFactor {
		return s.Range * s.Range
	}
	r := float64(s.Range) * cloakFactor
	return int(math.Ceil(r * r))
}

// RangeSquared returns the cloak adjusted RangeSquared value
func (s scanner) RangePenSquared(cloakFactor float64) int {
	if cloakFactor == 1 {
		return s.RangePen * s.RangePen
	}
	r := float64(s.RangePen) * cloakFactor
	return int(math.Ceil(r * r))
}

type playerScanner struct {
	universe          *Universe
	rules             *Rules
	player            *Player
	players           []*Player
	discoveredPlayers map[int]bool
	discoverer        discoverer
}

func newPlayerScanner(universe *Universe, players []*Player, rules *Rules, player *Player) playerScanner {
	return playerScanner{universe, rules, player, players, make(map[int]bool, len(player.Intels.PlayerIntels)), player.discoverer}
}

// scan planets, fleets, etc for a player
func (scan *playerScanner) scan() error {
	scanners := scan.getScanners()
	remoteMiningScanners := scan.getRemoteMiningScanners()
	cargoScanners := scan.getCargoScanners()
	starGateScanners := scan.getStarGateScanners()

	// scan planets
	if err := scan.scanPlanets(scanners, append(cargoScanners, remoteMiningScanners...), starGateScanners); err != nil {
		return err
	}

	// scan universe
	scan.scanFleets(scanners, cargoScanners)
	scan.scanMinefields(scanners)
	scan.scanMineralPackets(scanners)
	scan.scanSalvages(scanners)
	scan.scanWormholes(scanners)
	scan.scanMysteryTraders()

	// scan ally objects
	scan.discoverAllies()

	// after we've discovered a bunch of stuff, make sure we discover the race
	// names of any players that owned scanned objects
	scan.discoverPlayers()

	// after our intel is updated, update the fleet targets to account for lost targets
	scan.updateFleetTargets()

	return nil
}

// scan all planets with this player's scanners
func (scan *playerScanner) scanPlanets(scanners []scanner, cargoScanners []scanner, starGateScanners []scanner) error {
	for _, planet := range scan.universe.Planets {
		if planet.OwnedBy(scan.player.Num) {
			// scan owned planets
			if err := scan.discoverer.discoverPlanet(scan.rules, planet, true, true); err != nil {
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

		// try and scan the planet's stargate
		if planet.Spec.PlanetStarbaseSpec.HasStargate {
			for _, scanner := range starGateScanners {
				if scan.fleetInScannerRange(planet.Starbase, scanner) {
					scanned, err := scan.scanPlanet(planet, scanner)
					if err != nil {
						return err
					}
					if scanned {
						break
					}
				}
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

		intel := scan.player.GetPlanetIntel(planet.Num)
		if intel.ReportAge != ReportAgeUnexplored {
			// TODO: remove this after initial test games are done
			// it's just here because some old games don't have basehab on intels
			if intel.BaseHab != planet.BaseHab {
				intel.BaseHab = planet.BaseHab
			}
			scan.discoverer.discoverPlanetTerraformability(planet.Num)
		}
		// TODO: another fix for planets we don't own still being ours in intel
		if (intel.PlayerNum == scan.player.Num && planet.PlayerNum != scan.player.Num) || (intel.PlayerNum == Unowned && intel.Cargo.Colonists > 0) {
			// we think we own this planet, but we don't. Remove ownership info
			scan.discoverer.clearPlanetOwnerIntel(planet)
		}

	}

	return nil
}

// scan this planet
func (scan *playerScanner) scanPlanet(planet *Planet, scanner scanner) (scanned bool, err error) {
	if scanner.RangePen != NoScanner && scanner.RangePenSquared(NoCloakFactor) >= scanner.Position.DistanceSquaredTo(planet.Position) {
		if planet.Owned() {
			scan.discoveredPlayers[planet.PlayerNum] = true
		}
		if err := scan.discoverer.discoverPlanet(scan.rules, planet, true, planet.OwnedBy(scan.player.Num)); err != nil {
			return false, err
		}
		if scanner.DiscoverPlanetCargo {
			if err := scan.discoverer.discoverPlanetCargo(planet); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	// non-pen scan this planet if we are right on top
	if scanner.Range != NoScanner && scanner.Position.DistanceSquaredTo(planet.Position) == 0 {
		if planet.Owned() {
			scan.discoveredPlayers[planet.PlayerNum] = true
		}
		if err := scan.discoverer.discoverPlanet(scan.rules, planet, false, false); err != nil {
			return false, err
		}
	}
	return false, nil
}

// scan all fleets and discover their designs if we should
func (scan *playerScanner) scanFleets(scanners []scanner, cargoScanners []scanner) {
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

		// try and scan the fleet with this scanner
		for _, scanner := range scanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToScan = append(fleetsToScan, fleet)
				break
			}
		}

		// try and scan the fleet with a cargo scanner
		for _, scanner := range cargoScanners {
			if scan.fleetInScannerRange(fleet, scanner) {
				fleetsToCargoScan = append(fleetsToScan, fleet)
				break
			}
		}
	}

	for _, fleet := range fleetsToScan {
		scan.discoveredPlayers[fleet.PlayerNum] = true

		for _, token := range fleet.Tokens {
			scan.discoverer.discoverDesign(token.design, scan.player.Race.Spec.DiscoverDesignOnScan)
		}

		scan.discoverer.discoverFleet(fleet, false)
	}

	for _, fleet := range fleetsToCargoScan {
		scan.discoverer.discoverFleetCargo(fleet)
	}
}

// return true if this scanner successfully scans this fleet, taking into account cloaking
// and the fleet's cloak penetration
func (scan *playerScanner) fleetInScannerRange(fleet *Fleet, scanner scanner) bool {
	cloakFactor := getCloakFactor(fleet.Spec.CloakPercent, scanner.CloakReductionFactor)
	distanceSquared := scanner.Position.DistanceSquaredTo(fleet.Position)
	scanRangePenSqaured := scanner.RangePenSquared(cloakFactor)

	// if we pen scanned this, update the report
	if scanner.RangePen != NoScanner && scanRangePenSqaured >= distanceSquared {
		// update the fleet report with pen scanners
		return true
	}

	// if we aren't orbiting a planet, we can be seen with regular scanners
	scanRangeSqaured := scanner.RangeSquared(cloakFactor)
	if scanner.Range != NoScanner && !fleet.Orbiting() && scanRangeSqaured >= distanceSquared {
		return true
	}
	return false
}

// scan all fleets and discover their designs if we should
func (scan *playerScanner) scanWormholes(scanners []scanner) {
	for _, wormhole := range scan.universe.Wormholes {
		intel := scan.player.GetWormholeIntel(wormhole.Num)

		for _, scanner := range scanners {
			if scanner.Range == NoScanner {
				continue
			}

			// calculate cloak reduction for tachyon detectors if this wormhole is cloaked
			cloakFactor := getCloakFactor(scan.rules.WormholeCloak, scanner.CloakReductionFactor)
			if intel != nil {
				cloakFactor = 1
			}

			distanceSquared := scanner.Position.DistanceSquaredTo(wormhole.Position)
			scanRangeSquared := scanner.RangeSquared(cloakFactor)
			// we only care about regular scanners for wormholes
			if scanRangeSquared >= distanceSquared {
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

	intels := make([]*Wormhole, len(scan.player.WormholeIntels))
	copy(intels, scan.player.WormholeIntels)
	for _, intel := range intels {
		for _, scanner := range scanners {
			if scanner.Range == NoScanner {
				continue
			}

			// if we scanned this wormhole where we last saw it, but it no longer exists in the universe, forget it
			if scanner.RangeSquared(NoCloakFactor) >= scanner.Position.DistanceSquaredTo(intel.Position) {
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
func (scan *playerScanner) scanMysteryTraders() {
	for _, mysteryTrader := range scan.universe.MysteryTraders {
		if mysteryTrader.Delete {
			continue
		}
		// every player discovers mystery traders
		scan.discoverer.discoverMysteryTrader(mysteryTrader)
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScanner) scanMineralPackets(scanners []scanner) {
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
			scan.discoverer.discoverMineralPacket(scan.rules, packet, packetPlayer, target)
			continue
		}

		for _, scanner := range scanners {
			if scanner.Range == NoScanner {
				continue
			}

			// we only care about regular scanners for mineral packets
			if scanner.RangeSquared(NoCloakFactor) >= scanner.Position.DistanceSquaredTo(packet.Position) {
				scan.discoverer.discoverMineralPacket(scan.rules, packet, packetPlayer, target)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScanner) scanMinefields(scanners []scanner) {
	for _, minefield := range scan.universe.Minefields {
		if minefield.Delete {
			continue
		}

		if minefield.OwnedBy(scan.player.Num) {
			// The player already gets a copy of all their own minefields
			continue
		}
		intel := scan.player.GetMinefieldIntel(minefield.PlayerNum, minefield.Num)

		for _, scanner := range scanners {
			if scanner.Range == NoScanner {
				continue
			}

			cloakFactor := getCloakFactor(scan.rules.MinefieldCloak, scanner.CloakReductionFactor)
			if intel != nil {
				cloakFactor = 1
			}

			distanceToEdge := max(0, scanner.Position.DistanceTo(minefield.Position)-minefield.Radius())
			scannerRange := float64(scanner.Range) * cloakFactor
			// we only care about regular scanners for wormholes
			if scannerRange >= distanceToEdge {
				scan.discoverer.discoverMinefield(minefield)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScanner) scanSalvages(scanners []scanner) {
	for _, salvage := range scan.universe.Salvages {
		if salvage.Delete {
			continue
		}
		for _, scanner := range scanners {
			if scanner.Range == NoScanner {
				continue
			}

			// we only care about regular scanners for mineral packets
			if scanner.RangeSquared(NoCloakFactor) >= scanner.Position.DistanceSquaredTo(salvage.Position) {
				scan.discoverer.discoverSalvage(salvage)
				break
			}
		}
	}
}

// discover any map sharing ally data
func (scan *playerScanner) discoverAllies() error {
	for _, player := range scan.players {
		if !player.IsSharingMap(scan.player.Num) {
			continue
		}

		// discover this ally's planets
		for _, planet := range scan.universe.Planets {
			if planet.PlayerNum != player.Num {
				continue
			}
			if err := scan.discoverer.discoverPlanet(scan.rules, planet, true, planet.OwnedBy(scan.player.Num)); err != nil {
				return err
			}
			if err := scan.discoverer.discoverPlanetCargo(planet); err != nil {
				return err
			}
			if err := scan.discoverer.discoverPlanetScanner(planet); err != nil {
				return err
			}
		}

		// discover any in use designs
		for _, design := range player.Designs {
			if design.Spec.NumInstances > 0 {
				scan.discoverer.discoverDesign(design, true)
			}
		}

		// discover this ally's fleets/designs
		for _, fleet := range scan.universe.Fleets {
			if fleet.PlayerNum != player.Num || fleet.Delete {
				continue
			}
			scan.discoverer.discoverFleet(fleet, true)
			scan.discoverer.discoverFleetCargo(fleet)
			scan.discoverer.discoverFleetScanner(fleet)
		}

		for _, mf := range scan.universe.Minefields {
			if mf.PlayerNum != player.Num {
				continue
			}
			scan.discoverer.discoverMinefield(mf)
		}

		for _, mp := range scan.universe.MineralPackets {
			if mp.PlayerNum != player.Num {
				continue
			}
			target := scan.universe.getPlanet(mp.TargetPlanetNum)
			scan.discoverer.discoverMineralPacket(scan.rules, mp, player, target)
			scan.discoverer.discoverMineralPacketScanner(mp)
		}

		// discover our ally and anyone they know about
		scan.discoveredPlayers[player.Num] = true
		for _, otherPlayer := range player.Intels.PlayerIntels {
			if otherPlayer.Seen {
				scan.discoveredPlayers[otherPlayer.Num] = true
			}
		}

	}

	return nil
}

func (scan *playerScanner) discoverPlayers() {
	for player, discovered := range scan.discoveredPlayers {
		if discovered {
			scan.discoverer.discoverPlayer(scan.players[player-1])
		}
	}
}

// get a list of unique scanners per player.
// This is a minimal list only containing the best scanner values for each position
func (scan *playerScanner) getScanners() []scanner {
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
				scanner.Range = NoScanner
				scanner.RangePen = NoScanner
				scanner.CloakReductionFactor = 1
			}
			if fleet.Spec.ScanRange != NoScanner {
				scanner.Range = max(scanner.Range, fleet.Spec.ScanRange)
			}
			if fleet.Spec.ScanRangePen != NoScanner {
				scanner.RangePen = max(scanner.RangePen, fleet.Spec.ScanRangePen)
			}
			scanner.CloakReductionFactor = min(scanner.CloakReductionFactor, fleet.Spec.ReduceCloaking)
			scanningFleetsByPosition[fleet.Position] = scanner
		}
	}

	// build a list of scanners for this player
	scanners := []scanner{}
	for _, planet := range scan.universe.Planets {
		if planet.PlayerNum == scan.player.Num {
			// planets we own without scanners act as range 0 scanners
			planetaryScanner := scanner{
				Position:             planet.Position,
				Range:                0,
				RangePen:             0,
				CloakReductionFactor: 1,
			}

			if planet.Scanner {
				// update this scanner to use the planetary scanner stats
				planetaryScanner = scanner{
					Position:             planet.Position,
					Range:                planet.Spec.ScanRange,
					RangePen:             planet.Spec.ScanRangePen,
					CloakReductionFactor: 1,
				}
			}
			// use the fleet scanner if it's better
			if fleetScanner, ok := scanningFleetsByPosition[planet.Position]; ok {
				planetaryScanner.Range = max(planetaryScanner.Range, fleetScanner.Range)
				planetaryScanner.RangePen = max(planetaryScanner.RangePen, fleetScanner.RangePen)
				planetaryScanner.CloakReductionFactor = min(planetaryScanner.CloakReductionFactor, fleetScanner.CloakReductionFactor)
			}
			scanners = append(scanners, planetaryScanner)
		}
	}

	// Space demolition minefields act as scanners
	if scan.player.Race.Spec.MinefieldsAreScanners {
		for _, minefield := range scan.universe.Minefields {
			if minefield.PlayerNum == scan.player.Num {
				scanner := scanner{
					Position:             minefield.Position,
					Range:                int(minefield.Radius()),
					CloakReductionFactor: 1,
				}
				// use the fleet scanner if it's better
				if fleetScanner, ok := scanningFleetsByPosition[minefield.Position]; ok {
					scanner.Range = max(scanner.Range, fleetScanner.Range)
					scanner.RangePen = max(scanner.RangePen, fleetScanner.RangePen)
					scanner.CloakReductionFactor = min(scanner.CloakReductionFactor, fleetScanner.CloakReductionFactor)
				}
				scanners = append(scanners, scanner)
			}
		}
	}

	// Packet Physics packets act as pen scanners
	for _, packet := range scan.universe.MineralPackets {
		if packet.PlayerNum == scan.player.Num && (packet.ScanRange != NoScanner || packet.ScanRangePen != NoScanner) {
			scanner := scanner{
				Position:             packet.Position,
				Range:                packet.ScanRange,
				RangePen:             packet.ScanRangePen,
				CloakReductionFactor: 1,
			}
			// use the fleet scanner if it's better
			if fleetScanner, ok := scanningFleetsByPosition[packet.Position]; ok {
				scanner.Range = max(scanner.Range, fleetScanner.Range)
				scanner.RangePen = max(scanner.RangePen, fleetScanner.RangePen)
				scanner.CloakReductionFactor = min(scanner.CloakReductionFactor, fleetScanner.CloakReductionFactor)
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
func (scan *playerScanner) getRemoteMiningScanners() []scanner {
	scanningFleetsByPosition := map[Vector]scanner{}
	for _, fleet := range scan.universe.Fleets {
		// find any fleets that remote mined this turn, but only add one per position
		if fleet.PlayerNum == scan.player.Num && fleet.remoteMined {
			if scanner, found := scanningFleetsByPosition[fleet.Position]; !found {
				scanner.Position = fleet.Position
				scanner.Range = 0
				scanner.RangePen = 0
				scanner.DiscoverPlanetCargo = true
				scanner.CloakReductionFactor = 1
				scanningFleetsByPosition[fleet.Position] = scanner
			}
		}
	}

	return maps.Values(scanningFleetsByPosition)
}

// get a list of scanners that can scan cargo from fleets or planets
func (scan *playerScanner) getCargoScanners() []scanner {
	scanners := []scanner{}
	scanningFleetsByPosition := map[Vector]scanner{}

	for _, fleet := range scan.universe.Fleets {
		if fleet.PlayerNum == scan.player.Num && fleet.Spec.Scanner && (fleet.Spec.CanStealFleetCargo || fleet.Spec.CanStealPlanetCargo) {
			scanner, found := scanningFleetsByPosition[fleet.Position]
			if !found {
				// start with NoScanner (-1)
				scanner.Position = fleet.Position
				scanner.Range = NoScanner
				scanner.RangePen = NoScanner
				scanner.CloakReductionFactor = 1
			}
			scanner.Range = max(scanner.Range, fleet.Spec.ScanRange)
			scanner.RangePen = max(scanner.RangePen, fleet.Spec.ScanRangePen)
			scanner.CloakReductionFactor = min(scanner.CloakReductionFactor, fleet.Spec.ReduceCloaking)
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

// get a list of star gates that can scan other star gates by player
func (scan *playerScanner) getStarGateScanners() []scanner {
	scanners := []scanner{}
	if !scan.player.Race.Spec.CanDetectStargatePlanets {
		return scanners
	}
	for _, planet := range scan.universe.Planets {
		if planet.PlayerNum == scan.player.Num && planet.Spec.PlanetStarbaseSpec.HasStargate {
			penRange := min(planet.Spec.PlanetStarbaseSpec.SafeRange, math.MaxInt16)
			scanner := scanner{
				Position:             planet.Position,
				RangePen:             penRange,
				CloakReductionFactor: 1,
			}
			scanners = append(scanners, scanner)
		}
	}
	return scanners
}

// make sure our fleets are pointing to valid targets
func (scan *playerScanner) updateFleetTargets() {
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
			if fleet.PreviousPosition != nil && fleet.OrbitingPlanetNum == None && wp0.TargetType != MapObjectTypeNone {
				// we arrived at our target, but it's not a planet. Keep it as wp1
				fleet.Waypoints = []Waypoint{NewPositionWaypoint(fleet.Position, fleet.WarpSpeed), wp0}
			} else {
				fleet.WarpSpeed = 0
				fleet.Heading = VectorFloat64{}
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
				target := scan.player.GetFleetIntel(wp.TargetPlayerNum, wp.TargetNum)
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
			case MapObjectTypeMysteryTrader:
				target := scan.player.GetMysteryTraderIntel(wp.TargetNum)
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
				target := scan.player.GetSalvageIntel(wp.TargetNum)
				if target == nil {
					messager.fleetTargetLost(scan.player, fleet, wp.TargetName, wp.TargetType)
					wp.TargetType = MapObjectTypeNone
					wp.TargetPlayerNum = None
					wp.TargetNum = None
					wp.TargetName = ""
				}

			case MapObjectTypeMineralPacket:
				target := scan.player.GetMineralPacketIntel(wp.TargetPlayerNum, wp.TargetNum)
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
				target := scan.player.GetWormholeIntel(wp.TargetNum)
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
