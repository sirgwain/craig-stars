package cs

import "math"

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

	// TODO: add in player mineral packets, minefields, etc
	scanners := scan.getScanners()
	cargoScanners := scan.getCargoScanners()

	// scan planets
	if err := scan.scanPlanets(scanners, cargoScanners); err != nil {
		return err
	}

	// scan universe
	scan.scanFleets(scanners, cargoScanners)
	scan.scanMineFields(scanners)
	scan.scanWormholes(scanners)

	scan.discoverPlayers()

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
	}

	return nil
}

// scan this planet
func (scan *playerScan) scanPlanet(planet *Planet, scanner scanner) (bool, error) {
	dist := scanner.Position.DistanceSquaredTo(planet.Position)
	_ = dist
	if float64(scanner.RangePenSquared) >= scanner.Position.DistanceSquaredTo(planet.Position) {
		if planet.owned() {
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
				// calculate cloak reduction for tachyon detectors if this minefield is cloaked
				cloakFactor = 1 - (1-cloakFactor)*scanner.CloakReduction
			}

			// we only care about regular scanners for wormholes
			if float64(scanner.RangeSquared)*cloakFactor >= scanner.Position.DistanceSquaredTo(wormhole.Position) {
				scan.discoverer.discoverWormhole(scan.player, wormhole)
				break
			}
		}
	}
}

// scan all fleets and discover their designs if we should
func (scan *playerScan) scanMineFields(scanners []scanner) {
	for _, mineField := range scan.universe.MineFields {
		if mineField.OwnedBy(scan.player.Num) {
			// The player already gets a copy of all their own mineFields
			continue
		}
		intel := scan.discoverer.getMineFieldIntel(mineField.Num)

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
	planetaryScanner := scan.player.Spec.PlanetaryScanner
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
			scanner.RangeSquared = maxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
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
				RangeSquared:    planetaryScanner.ScanRange * planetaryScanner.ScanRange,
				RangePenSquared: planetaryScanner.ScanRangePen * planetaryScanner.ScanRangePen,
			}

			// use the fleet scanner if it's better
			if fleetScanner, ok := scanningFleetsByPosition[planet.Position]; ok {
				scanner.RangeSquared = maxInt(scanner.RangeSquared, fleetScanner.RangeSquared)
				scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleetScanner.RangePenSquared)
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
			scanner.RangeSquared = maxInt(scanner.RangeSquared, fleet.Spec.ScanRange*fleet.Spec.ScanRange)
			scanner.RangePenSquared = maxInt(scanner.RangePenSquared, fleet.Spec.ScanRangePen*fleet.Spec.ScanRangePen)
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
